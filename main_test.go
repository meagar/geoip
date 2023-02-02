package main

import (
	"net"
	"testing"
)

func Test_locateIP(t *testing.T) {

	db := openGeoLiteDb()
	testCases := []struct {
		ip     string
		locale string
		want   string
	}{
		{"194.60.38.225", "en", "Kensington, England, Royal Kensington and Chelsea, United Kingdom"},
		{"194.60.38.225", "ru", "Кенсингтон, Англия, Великобритания"},
		{"8.8.8.8", "en", "United States"},
		{"8.8.8.8", "fr", "États-Unis"},
	}

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				panic("Bad test case")
			}
			got, err := locateIP(ip, db, tc.locale)
			t.Logf("%v, %v", got, err)
			if err != nil {
				t.Errorf("Locating %s: Got unexpected error %v", tc.ip, err)
				t.FailNow()
			}
			if got != tc.want {
				t.Errorf("locateIP(%q, %q): Got %v, want %v", tc.ip, tc.locale, got, tc.want)
			}
		})
	}

}
