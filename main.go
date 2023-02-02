package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	golog "log"

	"github.com/oschwald/geoip2-golang"
	lz4 "github.com/pierrec/lz4/v4"
)

//go:embed GeoLite2-City_20210608/GeoLite2-City.mmdb.lz4
var geoliteDb []byte

// CLI flags
var locale string = "en"

func parseFlags(supportedLocales []string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: geoip [OPTIONS] [IP_ADDRESS]\n")
		flag.PrintDefaults()
	}

	localeList := strings.Join(supportedLocales, ", ")

	flag.StringVar(&locale, "locale", "en", "Locale to use when displaying result ("+localeList+")")
	flag.Parse()

	// Verify that our locale is valid
	if !validLocale(locale, supportedLocales) {
		die("Error: Invalid locale %q; valid locales: %s\n", locale, localeList)
	}

	// Verify that we have exactly 1 remaining argument (the IP)
	if flag.NArg() == 0 {
		die("Error: Missing IP_ADDRESS argument. See geoip --help\n")
	} else if flag.NArg() > 1 {
		die("Error: Too many arguments. See geoip --help\n")
	}
}

func main() {
	db := openGeoLiteDb()
	defer db.Close()

	parseFlags(db.Metadata().Languages)

	ip := ipArg()
	location, err := locateIP(ip, db, locale)
	if err != nil {
		die(err.Error())
	}

	fmt.Println(location)
}

func locateIP(ip net.IP, db *geoip2.Reader, locale string) (string, error) {
	record, err := db.City(ip)
	if err != nil {
		return "", err
	}

	places := []string{}

	if record.City.Names[locale] != "" {
		places = append(places, record.City.Names[locale])
	}

	for _, place := range record.Subdivisions {
		log("code: %v - %v", place.IsoCode, place.Names[locale])
		if place.Names[locale] != "" {
			places = append(places, place.Names[locale])
		}
	}

	places = append(places, record.Country.Names[locale])

	return strings.Join(places, ", "), nil
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func log(format string, args ...interface{}) {
	// For debugging
	return
	golog.Printf(format, args...)
}

func ipArg() net.IP {
	arg := flag.Arg(0)
	log("Have IP address %v", arg)
	if ip := net.ParseIP(arg); ip != nil {
		return ip
	}

	die("Unable to parse %q as an IP address\n", arg)
	return nil
}

func validLocale(locale string, supportedLocales []string) bool {
	for _, language := range supportedLocales {
		if language == locale {
			return true
		}
	}

	return false
}

// Decompresses the embedded database and instantiates a geoip2.Reader
// lz4 provides the best speed:compression ratio.
// Embedding the uncompressed file directly adds ~60mb to the binary and parsing it adds ~.05 seconds to startup
// Embedding gzip (--best or --fast) adds ~30mb (2x improvement) but takes ~.5 seconds, ~10x slower startup
// Embedding bzip adds ~30mb but takes ~100x longer to decompress, adding ~3.5 seconds to startup (very bad)
// Embedding lz4 archive adds ~37mb but adds only ~0.1 seconds to startup, which seems like a fair trade for the extra ~7mb
// Using lz4 makes it fast enough to load the DB that we can load it *before* processing CLI flags, which
// allows us to ask the db for its supported languages to validate the --locale flag.
func openGeoLiteDb() *geoip2.Reader {
	reader := lz4.NewReader(bytes.NewReader(geoliteDb))

	decompressedDb, err := io.ReadAll(reader)
	if err != nil {
		die("Error decompressing GeoLite2-City: %s\n", err)
	}

	db, err := geoip2.FromBytes(decompressedDb)
	if err != nil {
		die("Unable to load GeoLite2 DB: %s\n", err.Error())
	}

	return db
}
