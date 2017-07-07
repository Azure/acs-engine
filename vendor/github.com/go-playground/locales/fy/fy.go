package fy

import (
	"math"
	"strconv"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/currency"
)

type fy struct {
	locale                 string
	pluralsCardinal        []locales.PluralRule
	pluralsOrdinal         []locales.PluralRule
	pluralsRange           []locales.PluralRule
	decimal                string
	group                  string
	minus                  string
	percent                string
	perMille               string
	timeSeparator          string
	inifinity              string
	currencies             []string // idx = enum of currency code
	currencyPositivePrefix string
	currencyNegativePrefix string
	currencyNegativeSuffix string
	monthsAbbreviated      []string
	monthsNarrow           []string
	monthsWide             []string
	daysAbbreviated        []string
	daysNarrow             []string
	daysShort              []string
	daysWide               []string
	periodsAbbreviated     []string
	periodsNarrow          []string
	periodsShort           []string
	periodsWide            []string
	erasAbbreviated        []string
	erasNarrow             []string
	erasWide               []string
	timezones              map[string]string
}

// New returns a new instance of translator for the 'fy' locale
func New() locales.Translator {
	return &fy{
		locale:                 "fy",
		pluralsCardinal:        []locales.PluralRule{2, 6},
		pluralsOrdinal:         []locales.PluralRule{6},
		pluralsRange:           nil,
		decimal:                ",",
		group:                  ".",
		minus:                  "-",
		percent:                "%",
		perMille:               "‰",
		timeSeparator:          ":",
		inifinity:              "∞",
		currencies:             []string{"ADP", "AED", "AFA", "AFN", "ALK", "ALL", "AMD", "ANG", "AOA", "AOK", "AON", "AOR", "ARA", "ARL", "ARM", "ARP", "ARS", "ATS", "AU$", "AWG", "AZM", "AZN", "BAD", "BAM", "BAN", "BBD", "BDT", "BEC", "BEF", "BEL", "BGL", "BGM", "BGN", "BGO", "BHD", "BIF", "BMD", "BND", "BOB", "BOL", "BOP", "BOV", "BRB", "BRC", "BRE", "R$", "BRN", "BRR", "BRZ", "BSD", "BTN", "BUK", "BWP", "BYB", "BYN", "BYR", "BZD", "C$", "CDF", "CHE", "CHF", "CHW", "CLE", "CLF", "CLP", "CNX", "CN¥", "COP", "COU", "CRC", "CSD", "CSK", "CUC", "CUP", "CVE", "CYP", "CZK", "DDM", "DEM", "DJF", "DKK", "DOP", "DZD", "ECS", "ECV", "EEK", "EGP", "ERN", "ESA", "ESB", "ESP", "ETB", "€", "FIM", "FJ$", "FKP", "FRF", "£", "GEK", "GEL", "GHC", "GHS", "GIP", "GMD", "GNF", "GNS", "GQE", "GRD", "GTQ", "GWE", "GWP", "GYD", "HK$", "HNL", "HRD", "HRK", "HTG", "HUF", "IDR", "IEP", "ILP", "ILR", "₪", "₹", "IQD", "IRR", "ISJ", "ISK", "ITL", "JMD", "JOD", "JP¥", "KES", "KGS", "KHR", "KMF", "KPW", "KRH", "KRO", "₩", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LTT", "LUC", "LUF", "LUL", "LVL", "LVR", "LYD", "MAD", "MAF", "MCF", "MDC", "MDL", "MGA", "MGF", "MKD", "MKN", "MLF", "MMK", "MNT", "MOP", "MRO", "MTL", "MTP", "MUR", "MVP", "MVR", "MWK", "MX$", "MXP", "MXV", "MYR", "MZE", "MZM", "MZN", "NAD", "NGN", "NIC", "NIO", "NLG", "NOK", "NPR", "NZ$", "OMR", "PAB", "PEI", "PEN", "PES", "PGK", "PHP", "PKR", "PLN", "PLZ", "PTE", "PYG", "QAR", "RHD", "ROL", "RON", "RSD", "RUB", "RUR", "RWF", "SAR", "SI$", "SCR", "SDD", "SDG", "SDP", "SEK", "SGD", "SHP", "SIT", "SKK", "SLL", "SOS", "SRD", "SRG", "SSP", "STD", "SUR", "SVC", "SYP", "SZL", "฿", "TJR", "TJS", "TMM", "TMT", "TND", "TOP", "TPE", "TRL", "TRY", "TTD", "NT$", "TZS", "UAH", "UAK", "UGS", "UGX", "US$", "USN", "USS", "UYI", "UYP", "UYU", "UZS", "VEB", "VEF", "₫", "VNN", "VUV", "WST", "FCFA", "XAG", "XAU", "XBA", "XBB", "XBC", "XBD", "EC$", "XDR", "XEU", "XFO", "XFU", "CFA", "XPD", "XPF", "XPT", "XRE", "XSU", "XTS", "XUA", "XXX", "YDD", "YER", "YUD", "YUM", "YUN", "YUR", "ZAL", "ZAR", "ZMK", "ZMW", "ZRN", "ZRZ", "ZWD", "ZWL", "ZWR"},
		currencyPositivePrefix: " ",
		currencyNegativePrefix: "( ",
		currencyNegativeSuffix: ")",
		monthsAbbreviated:      []string{"", "Jan", "Feb", "Mrt", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Des"},
		monthsNarrow:           []string{"", "J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"},
		monthsWide:             []string{"", "Jannewaris", "Febrewaris", "Maart", "April", "Maaie", "Juny", "July", "Augustus", "Septimber", "Oktober", "Novimber", "Desimber"},
		daysAbbreviated:        []string{"si", "mo", "ti", "wo", "to", "fr", "so"},
		daysNarrow:             []string{"S", "M", "T", "W", "T", "F", "S"},
		daysShort:              []string{"si", "mo", "ti", "wo", "to", "fr", "so"},
		daysWide:               []string{"snein", "moandei", "tiisdei", "woansdei", "tongersdei", "freed", "sneon"},
		periodsAbbreviated:     []string{"AM", "PM"},
		periodsNarrow:          []string{"AM", "PM"},
		periodsWide:            []string{"AM", "PM"},
		erasAbbreviated:        []string{"f.Kr.", "n.Kr."},
		erasNarrow:             []string{"f.K.", "n.K."},
		erasWide:               []string{"Foar Kristus", "nei Kristus"},
		timezones:              map[string]string{"HKST": "Hongkongse simmertiid", "AKST": "Alaska-standerttiid", "UYST": "Uruguayaanske simmertiid", "WIT": "East-Yndonezyske tiid", "TMT": "Turkmeense standerttiid", "HNT": "Newfoundlânske-standerttiid", "PDT": "Pasifik-simmertiid", "ACWST": "Midden-Australyske westelijke standerttiid", "WEZ": "West-Europeeske standerttiid", "SRT": "Surinaamske tiid", "ADT": "Atlantic-simmertiid", "ACDT": "Midden-Australyske simmertiid", "GFT": "Frâns-Guyaanske tiid", "OEZ": "East-Europeeske standerttiid", "HEEG": "East-Groenlânske simmertiid", "BOT": "Boliviaanske tiid", "NZST": "Nij-Seelânske standerttiid", "AST": "Atlantic-standerttiid", "WART": "West-Argentynske standerttiid", "WAST": "West-Afrikaanske simmertiid", "GYT": "Guyaanske tiid", "EAT": "East-Afrikaanske tiid", "WIB": "West-Yndonezyske tiid", "MEZ": "Midden-Europeeske standerttiid", "MYT": "Maleisyske tiid", "WARST": "West-Argentynske simmertiid", "ARST": "Argentynske simmertiid", "HEPM": "Saint Pierre en Miquelon-simmertiid", "LHST": "Lord Howe-eilânske standerttiid", "HNCU": "Kubaanske standerttiid", "HNNOMX": "HNNOMX", "BT": "Bhutaanske tiid", "∅∅∅": "Azoren-simmertiid", "CAT": "Sintraal-Afrikaanske tiid", "HKT": "Hongkongse standerttiid", "OESZ": "East-Europeeske simmertiid", "WAT": "West-Afrikaanske standerttiid", "SGT": "Singaporese standerttiid", "CLST": "Sileenske simmertiid", "HNOG": "West-Groenlânske standerttiid", "EST": "Eastern-standerttiid", "HENOMX": "HENOMX", "AKDT": "Alaska-simmertiid", "HNPM": "Saint Pierre en Miquelon-standerttiid", "HEPMX": "HEPMX", "CDT": "Central-simmertiid", "JST": "Japanske standerttiid", "CLT": "Sileenske standerttiid", "COT": "Kolombiaanske standerttiid", "HNEG": "East-Groenlânske standerttiid", "MST": "Mountain-standerttiid", "SAST": "Sûd-Afrikaanske tiid", "HNPMX": "HNPMX", "ECT": "Ecuadoraanske tiid", "NZDT": "Nij-Seelânske simmertiid", "COST": "Kolombiaanske simmertiid", "WITA": "Sintraal-Yndonezyske tiid", "MDT": "Mountain-simmertiid", "AWDT": "West-Australyske simmertiid", "ACWDT": "Midden-Australyske westelijke simmertiid", "HAT": "Newfoundlânske-simmertiid", "HECU": "Kubaanske simmertiid", "AEST": "East-Australyske standerttiid", "ChST": "Chamorro-tiid", "UYT": "Uruguayaanske standerttiid", "ART": "Argentynske standerttiid", "EDT": "Eastern-simmertiid", "AEDT": "East-Australyske simmertiid", "IST": "Yndiaaske tiid", "HEOG": "West-Groenlânske simmertiid", "CHAST": "Chatham standerttiid", "CHADT": "Chatham simmertiid", "HAST": "Hawaii-Aleoetyske standerttiid", "WESZ": "West-Europeeske simmertiid", "LHDT": "Lord Howe-eilânske simmertiid", "CST": "Central-standerttiid", "AWST": "West-Australyske standerttiid", "VET": "Fenezolaanske tiid", "JDT": "Japanske simmertiid", "GMT": "Greenwich Mean Time", "TMST": "Turkmeense simmertiid", "ACST": "Midden-Australyske standerttiid", "PST": "Pasifik-standerttiid", "HADT": "Hawaii-Aleoetyske simmertiid", "MESZ": "Midden-Europeeske simmertiid"},
	}
}

// Locale returns the current translators string locale
func (fy *fy) Locale() string {
	return fy.locale
}

// PluralsCardinal returns the list of cardinal plural rules associated with 'fy'
func (fy *fy) PluralsCardinal() []locales.PluralRule {
	return fy.pluralsCardinal
}

// PluralsOrdinal returns the list of ordinal plural rules associated with 'fy'
func (fy *fy) PluralsOrdinal() []locales.PluralRule {
	return fy.pluralsOrdinal
}

// PluralsRange returns the list of range plural rules associated with 'fy'
func (fy *fy) PluralsRange() []locales.PluralRule {
	return fy.pluralsRange
}

// CardinalPluralRule returns the cardinal PluralRule given 'num' and digits/precision of 'v' for 'fy'
func (fy *fy) CardinalPluralRule(num float64, v uint64) locales.PluralRule {

	n := math.Abs(num)
	i := int64(n)

	if i == 1 && v == 0 {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther
}

// OrdinalPluralRule returns the ordinal PluralRule given 'num' and digits/precision of 'v' for 'fy'
func (fy *fy) OrdinalPluralRule(num float64, v uint64) locales.PluralRule {
	return locales.PluralRuleOther
}

// RangePluralRule returns the ordinal PluralRule given 'num1', 'num2' and digits/precision of 'v1' and 'v2' for 'fy'
func (fy *fy) RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) locales.PluralRule {
	return locales.PluralRuleUnknown
}

// MonthAbbreviated returns the locales abbreviated month given the 'month' provided
func (fy *fy) MonthAbbreviated(month time.Month) string {
	return fy.monthsAbbreviated[month]
}

// MonthsAbbreviated returns the locales abbreviated months
func (fy *fy) MonthsAbbreviated() []string {
	return fy.monthsAbbreviated[1:]
}

// MonthNarrow returns the locales narrow month given the 'month' provided
func (fy *fy) MonthNarrow(month time.Month) string {
	return fy.monthsNarrow[month]
}

// MonthsNarrow returns the locales narrow months
func (fy *fy) MonthsNarrow() []string {
	return fy.monthsNarrow[1:]
}

// MonthWide returns the locales wide month given the 'month' provided
func (fy *fy) MonthWide(month time.Month) string {
	return fy.monthsWide[month]
}

// MonthsWide returns the locales wide months
func (fy *fy) MonthsWide() []string {
	return fy.monthsWide[1:]
}

// WeekdayAbbreviated returns the locales abbreviated weekday given the 'weekday' provided
func (fy *fy) WeekdayAbbreviated(weekday time.Weekday) string {
	return fy.daysAbbreviated[weekday]
}

// WeekdaysAbbreviated returns the locales abbreviated weekdays
func (fy *fy) WeekdaysAbbreviated() []string {
	return fy.daysAbbreviated
}

// WeekdayNarrow returns the locales narrow weekday given the 'weekday' provided
func (fy *fy) WeekdayNarrow(weekday time.Weekday) string {
	return fy.daysNarrow[weekday]
}

// WeekdaysNarrow returns the locales narrow weekdays
func (fy *fy) WeekdaysNarrow() []string {
	return fy.daysNarrow
}

// WeekdayShort returns the locales short weekday given the 'weekday' provided
func (fy *fy) WeekdayShort(weekday time.Weekday) string {
	return fy.daysShort[weekday]
}

// WeekdaysShort returns the locales short weekdays
func (fy *fy) WeekdaysShort() []string {
	return fy.daysShort
}

// WeekdayWide returns the locales wide weekday given the 'weekday' provided
func (fy *fy) WeekdayWide(weekday time.Weekday) string {
	return fy.daysWide[weekday]
}

// WeekdaysWide returns the locales wide weekdays
func (fy *fy) WeekdaysWide() []string {
	return fy.daysWide
}

// FmtNumber returns 'num' with digits/precision of 'v' for 'fy' and handles both Whole and Real numbers based on 'v'
func (fy *fy) FmtNumber(num float64, v uint64) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 2 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, fy.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, fy.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, fy.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

// FmtPercent returns 'num' with digits/precision of 'v' for 'fy' and handles both Whole and Real numbers based on 'v'
// NOTE: 'num' passed into FmtPercent is assumed to be in percent already
func (fy *fy) FmtPercent(num float64, v uint64) string {
	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 3
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, fy.decimal[0])
			continue
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, fy.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	b = append(b, fy.percent...)

	return string(b)
}

// FmtCurrency returns the currency representation of 'num' with digits/precision of 'v' for 'fy'
func (fy *fy) FmtCurrency(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := fy.currencies[currency]
	l := len(s) + len(symbol) + 4 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, fy.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, fy.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	for j := len(symbol) - 1; j >= 0; j-- {
		b = append(b, symbol[j])
	}

	for j := len(fy.currencyPositivePrefix) - 1; j >= 0; j-- {
		b = append(b, fy.currencyPositivePrefix[j])
	}

	if num < 0 {
		b = append(b, fy.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, fy.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	return string(b)
}

// FmtAccounting returns the currency representation of 'num' with digits/precision of 'v' for 'fy'
// in accounting notation.
func (fy *fy) FmtAccounting(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := fy.currencies[currency]
	l := len(s) + len(symbol) + 6 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, fy.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, fy.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {

		for j := len(symbol) - 1; j >= 0; j-- {
			b = append(b, symbol[j])
		}

		for j := len(fy.currencyNegativePrefix) - 1; j >= 0; j-- {
			b = append(b, fy.currencyNegativePrefix[j])
		}

	} else {

		for j := len(symbol) - 1; j >= 0; j-- {
			b = append(b, symbol[j])
		}

		for j := len(fy.currencyPositivePrefix) - 1; j >= 0; j-- {
			b = append(b, fy.currencyPositivePrefix[j])
		}

	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, fy.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	if num < 0 {
		b = append(b, fy.currencyNegativeSuffix...)
	}

	return string(b)
}

// FmtDateShort returns the short date representation of 't' for 'fy'
func (fy *fy) FmtDateShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Day() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2d}...)

	if t.Month() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Month()), 10)

	b = append(b, []byte{0x2d}...)

	if t.Year() > 9 {
		b = append(b, strconv.Itoa(t.Year())[2:]...)
	} else {
		b = append(b, strconv.Itoa(t.Year())[1:]...)
	}

	return string(b)
}

// FmtDateMedium returns the medium date representation of 't' for 'fy'
func (fy *fy) FmtDateMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, fy.monthsAbbreviated[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateLong returns the long date representation of 't' for 'fy'
func (fy *fy) FmtDateLong(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, fy.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateFull returns the full date representation of 't' for 'fy'
func (fy *fy) FmtDateFull(t time.Time) string {

	b := make([]byte, 0, 32)

	b = append(b, fy.daysWide[t.Weekday()]...)
	b = append(b, []byte{0x20}...)
	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, fy.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtTimeShort returns the short time representation of 't' for 'fy'
func (fy *fy) FmtTimeShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)

	return string(b)
}

// FmtTimeMedium returns the medium time representation of 't' for 'fy'
func (fy *fy) FmtTimeMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)

	return string(b)
}

// FmtTimeLong returns the long time representation of 't' for 'fy'
func (fy *fy) FmtTimeLong(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()
	b = append(b, tz...)

	return string(b)
}

// FmtTimeFull returns the full time representation of 't' for 'fy'
func (fy *fy) FmtTimeFull(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, fy.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()

	if btz, ok := fy.timezones[tz]; ok {
		b = append(b, btz...)
	} else {
		b = append(b, tz...)
	}

	return string(b)
}
