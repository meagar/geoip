package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-City_20210608/GeoLite2-City.mmdb.gz
var geoliteDb []byte

func main() {
	ip := ipArg()
	db := openGeoLiteDb()
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		die(err.Error())
	}

	places := []string{}

	if record.City.Names["en"] != "" {
		places = append(places, record.City.Names["en"])
	}

	for _, place := range record.Subdivisions {
		places = append(places, place.Names["en"])
	}

	places = append(places, record.Country.Names["en"])

	fmt.Println(strings.Join(places, ", "))
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func ipArg() net.IP {
	if len(os.Args) != 2 {
		die("Usage: geoip <IP ADDRESS>\n")
	}

	arg := os.Args[1]
	if ip := net.ParseIP(arg); ip != nil {
		return ip
	}

	die("Unable to parse %q as an IP address\n", arg)
	return nil
}

func openGeoLiteDb() *geoip2.Reader {
	reader, err := gzip.NewReader(bytes.NewReader(geoliteDb))
	if err != nil {
		die("Error creating reader: %s", err)
	}
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
