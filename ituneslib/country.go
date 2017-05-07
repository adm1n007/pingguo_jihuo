package ituneslib

import (
    "fmt"
)

type CountryID int

const (
    CountryID_China         = CountryID(143465)
    CountryID_Taiwan        = CountryID(143470)
    CountryID_Japan         = CountryID(143462)
    CountryID_India         = CountryID(143467)
    CountryID_NewZealand    = CountryID(143461)
    CountryID_Vietnam       = CountryID(143471)
    CountryID_UnitedState   = CountryID(143441)
)

type countryInfo struct {
    storeFront  string
    shortName   string
    name        string
    timeZone    string
}

/*++

    time zone = 8 * 3600 = 28800

--*/

var countryData = map[CountryID]*countryInfo{
    CountryID_China         : &countryInfo{"143465-19,32",   "CN",   "China",       "28800" },
    CountryID_Taiwan        : &countryInfo{"143470-18,32",   "TW",   "Taiwan",      "28800" },
    CountryID_Japan         : &countryInfo{"143462-9,32",    "JP",   "Japan",       "32400" },
    CountryID_India         : &countryInfo{"143467,32",      "IN",   "India",       "19800" },
    CountryID_NewZealand    : &countryInfo{"143461,32",      "NZ",   "NewZealand",  "46800" },
    CountryID_Vietnam       : &countryInfo{"143471-2,32",    "VN",   "Vietnam",     "25200" },
    CountryID_UnitedState   : &countryInfo{"143441-1,32",    "US",   "UnitedState", "-18000" },
}

func (self CountryID) String() string {
    if n, ok := countryData[self]; ok {
        return n.name
    }

    return fmt.Sprintf("UNKNOWN_COUNTRY_ID_%d", self)
}

func (self CountryID) Valid() bool {
    _, exists := countryData[self]
    return exists
}

func (self CountryID) ShortName() string {
    return countryData[self].shortName
}

func (self CountryID) CountryCode() int {
    return int(self)
}

func (self CountryID) StoreFront() string {
    return countryData[self].storeFront
}

func (self CountryID) TimeZone() string {
    return countryData[self].timeZone
}
