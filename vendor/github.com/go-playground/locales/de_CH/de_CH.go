package de_CH

import (
	"math"
	"strconv"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/currency"
)

type de_CH struct {
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

// New returns a new instance of translator for the 'de_CH' locale
func New() locales.Translator {
	return &de_CH{
		locale:                 "de_CH",
		pluralsCardinal:        []locales.PluralRule{2, 6},
		pluralsOrdinal:         []locales.PluralRule{6},
		pluralsRange:           []locales.PluralRule{2, 6},
		decimal:                ".",
		group:                  "'",
		minus:                  "-",
		percent:                "%",
		perMille:               "‰",
		timeSeparator:          ":",
		inifinity:              "∞",
		currencies:             []string{"ADP", "AED", "AFA", "AFN", "ALK", "ALL", "AMD", "ANG", "AOA", "AOK", "AON", "AOR", "ARA", "ARL", "ARM", "ARP", "ARS", "ATS", "AUD", "AWG", "AZM", "AZN", "BAD", "BAM", "BAN", "BBD", "BDT", "BEC", "BEF", "BEL", "BGL", "BGM", "BGN", "BGO", "BHD", "BIF", "BMD", "BND", "BOB", "BOL", "BOP", "BOV", "BRB", "BRC", "BRE", "BRL", "BRN", "BRR", "BRZ", "BSD", "BTN", "BUK", "BWP", "BYB", "BYN", "BYR", "BZD", "CAD", "CDF", "CHE", "CHF", "CHW", "CLE", "CLF", "CLP", "CNX", "CNY", "COP", "COU", "CRC", "CSD", "CSK", "CUC", "CUP", "CVE", "CYP", "CZK", "DDM", "DEM", "DJF", "DKK", "DOP", "DZD", "ECS", "ECV", "EEK", "EGP", "ERN", "ESA", "ESB", "ESP", "ETB", "EUR", "FIM", "FJD", "FKP", "FRF", "GBP", "GEK", "GEL", "GHC", "GHS", "GIP", "GMD", "GNF", "GNS", "GQE", "GRD", "GTQ", "GWE", "GWP", "GYD", "HKD", "HNL", "HRD", "HRK", "HTG", "HUF", "IDR", "IEP", "ILP", "ILR", "ILS", "INR", "IQD", "IRR", "ISJ", "ISK", "ITL", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRH", "KRO", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LTT", "LUC", "LUF", "LUL", "LVL", "LVR", "LYD", "MAD", "MAF", "MCF", "MDC", "MDL", "MGA", "MGF", "MKD", "MKN", "MLF", "MMK", "MNT", "MOP", "MRO", "MTL", "MTP", "MUR", "MVP", "MVR", "MWK", "MXN", "MXP", "MXV", "MYR", "MZE", "MZM", "MZN", "NAD", "NGN", "NIC", "NIO", "NLG", "NOK", "NPR", "NZD", "OMR", "PAB", "PEI", "PEN", "PES", "PGK", "PHP", "PKR", "PLN", "PLZ", "PTE", "PYG", "QAR", "RHD", "ROL", "RON", "RSD", "RUB", "RUR", "RWF", "SAR", "SBD", "SCR", "SDD", "SDG", "SDP", "SEK", "SGD", "SHP", "SIT", "SKK", "SLL", "SOS", "SRD", "SRG", "SSP", "STD", "SUR", "SVC", "SYP", "SZL", "THB", "TJR", "TJS", "TMM", "TMT", "TND", "TOP", "TPE", "TRL", "TRY", "TTD", "TWD", "TZS", "UAH", "UAK", "UGS", "UGX", "USD", "USN", "USS", "UYI", "UYP", "UYU", "UZS", "VEB", "VEF", "VND", "VNN", "VUV", "WST", "XAF", "XAG", "XAU", "XBA", "XBB", "XBC", "XBD", "XCD", "XDR", "XEU", "XFO", "XFU", "XOF", "XPD", "XPF", "XPT", "XRE", "XSU", "XTS", "XUA", "XXX", "YDD", "YER", "YUD", "YUM", "YUN", "YUR", "ZAL", "ZAR", "ZMK", "ZMW", "ZRN", "ZRZ", "ZWD", "ZWL", "ZWR"},
		currencyPositivePrefix: " ",
		currencyNegativePrefix: " ",
		monthsAbbreviated:      []string{"", "Jan.", "Feb.", "März", "Apr.", "Mai", "Juni", "Juli", "Aug.", "Sep.", "Okt.", "Nov.", "Dez."},
		monthsNarrow:           []string{"", "J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"},
		monthsWide:             []string{"", "Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"},
		daysAbbreviated:        []string{"So.", "Mo.", "Di.", "Mi.", "Do.", "Fr.", "Sa."},
		daysNarrow:             []string{"S", "M", "D", "M", "D", "F", "S"},
		daysShort:              []string{"So.", "Mo.", "Di.", "Mi.", "Do.", "Fr.", "Sa."},
		daysWide:               []string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"},
		periodsAbbreviated:     []string{"vorm.", "nachm."},
		periodsNarrow:          []string{"vm.", "nm."},
		periodsWide:            []string{"vorm.", "nachm."},
		erasAbbreviated:        []string{"v. Chr.", "n. Chr."},
		erasNarrow:             []string{"v. Chr.", "n. Chr."},
		erasWide:               []string{"v. Chr.", "n. Chr."},
		timezones:              map[string]string{"HNPMX": "Mexiko Pazifikzone-Normalzeit", "WAT": "Westafrikanische Normalzeit", "AKST": "Alaska-Normalzeit", "AEST": "Ostaustralische Normalzeit", "UYT": "Uruguyanische Normalzeit", "EAT": "Ostafrikanische Zeit", "SGT": "Singapur-Zeit", "COT": "Kolumbianische Normalzeit", "NZST": "Neuseeland-Normalzeit", "WESZ": "Westeuropäische Sommerzeit", "HAT": "Neufundland-Sommerzeit", "∅∅∅": "Azoren-Sommerzeit", "ECT": "Ecuadorianische Zeit", "PDT": "Nordamerikanische Westküsten-Sommerzeit", "CST": "Nordamerikanische Inland-Normalzeit", "HENOMX": "Mexiko Nordwestliche Zone-Sommerzeit", "HNT": "Neufundland-Normalzeit", "CLT": "Chilenische Normalzeit", "HEOG": "Westgrönland-Sommerzeit", "AEDT": "Ostaustralische Sommerzeit", "LHDT": "Lord-Howe-Sommerzeit", "VET": "Venezuela-Zeit", "MEZ": "Mitteleuropäische Normalzeit", "WART": "Westargentinische Normalzeit", "ADT": "Atlantik-Sommerzeit", "ART": "Argentinische Normalzeit", "HKST": "Hongkong-Sommerzeit", "LHST": "Lord-Howe-Normalzeit", "AWST": "Westaustralische Normalzeit", "AWDT": "Westaustralische Sommerzeit", "PST": "Nordamerikanische Westküsten-Normalzeit", "IST": "Indische Zeit", "ACWDT": "Zentral-/Westaustralische Sommerzeit", "MESZ": "Mitteleuropäische Sommerzeit", "MYT": "Malaysische Zeit", "HEPM": "Saint-Pierre-und-Miquelon-Sommerzeit", "HEPMX": "Mexiko Pazifikzone-Sommerzeit", "WIB": "Westindonesische Zeit", "EST": "Nordamerikanische Ostküsten-Normalzeit", "EDT": "Nordamerikanische Ostküsten-Sommerzeit", "WEZ": "Westeuropäische Normalzeit", "WIT": "Ostindonesische Zeit", "HAST": "Hawaii-Aleuten-Normalzeit", "JDT": "Japanische Sommerzeit", "SRT": "Suriname-Zeit", "JST": "Japanische Normalzeit", "COST": "Kolumbianische Sommerzeit", "HNOG": "Westgrönland-Normalzeit", "HKT": "Hongkong-Normalzeit", "BT": "Bhutan-Zeit", "GFT": "Französisch-Guayana-Zeit", "ChST": "Chamorro-Zeit", "UYST": "Uruguayanische Sommerzeit", "ACDT": "Zentralaustralische Sommerzeit", "MDT": "Rocky-Mountain-Sommerzeit", "HNPM": "Saint-Pierre-und-Miquelon-Normalzeit", "GYT": "Guyana-Zeit", "OESZ": "Osteuropäische Sommerzeit", "TMT": "Turkmenistan-Normalzeit", "ARST": "Argentinische Sommerzeit", "AKDT": "Alaska-Sommerzeit", "NZDT": "Neuseeland-Sommerzeit", "OEZ": "Osteuropäische Normalzeit", "HADT": "Hawaii-Aleuten-Sommerzeit", "CAT": "Zentralafrikanische Zeit", "ACST": "Zentralaustralische Normalzeit", "HEEG": "Ostgrönland-Sommerzeit", "MST": "Rocky Mountain-Normalzeit", "SAST": "Südafrikanische Zeit", "CDT": "Nordamerikanische Inland-Sommerzeit", "HNNOMX": "Mexiko Nordwestliche Zone-Normalzeit", "ACWST": "Zentral-/Westaustralische Normalzeit", "WARST": "Westargentinische Sommerzeit", "WAST": "Westafrikanische Sommerzeit", "CHAST": "Chatham-Normalzeit", "CHADT": "Chatham-Sommerzeit", "BOT": "Bolivianische Zeit", "CLST": "Chilenische Sommerzeit", "AST": "Atlantik-Normalzeit", "GMT": "Mittlere Greenwich-Zeit", "TMST": "Turkmenistan-Sommerzeit", "WITA": "Zentralindonesische Zeit", "HNCU": "Kubanische Normalzeit", "HECU": "Kubanische Sommerzeit", "HNEG": "Ostgrönland-Normalzeit"},
	}
}

// Locale returns the current translators string locale
func (de *de_CH) Locale() string {
	return de.locale
}

// PluralsCardinal returns the list of cardinal plural rules associated with 'de_CH'
func (de *de_CH) PluralsCardinal() []locales.PluralRule {
	return de.pluralsCardinal
}

// PluralsOrdinal returns the list of ordinal plural rules associated with 'de_CH'
func (de *de_CH) PluralsOrdinal() []locales.PluralRule {
	return de.pluralsOrdinal
}

// PluralsRange returns the list of range plural rules associated with 'de_CH'
func (de *de_CH) PluralsRange() []locales.PluralRule {
	return de.pluralsRange
}

// CardinalPluralRule returns the cardinal PluralRule given 'num' and digits/precision of 'v' for 'de_CH'
func (de *de_CH) CardinalPluralRule(num float64, v uint64) locales.PluralRule {

	n := math.Abs(num)
	i := int64(n)

	if i == 1 && v == 0 {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther
}

// OrdinalPluralRule returns the ordinal PluralRule given 'num' and digits/precision of 'v' for 'de_CH'
func (de *de_CH) OrdinalPluralRule(num float64, v uint64) locales.PluralRule {
	return locales.PluralRuleOther
}

// RangePluralRule returns the ordinal PluralRule given 'num1', 'num2' and digits/precision of 'v1' and 'v2' for 'de_CH'
func (de *de_CH) RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) locales.PluralRule {

	start := de.CardinalPluralRule(num1, v1)
	end := de.CardinalPluralRule(num2, v2)

	if start == locales.PluralRuleOne && end == locales.PluralRuleOther {
		return locales.PluralRuleOther
	} else if start == locales.PluralRuleOther && end == locales.PluralRuleOne {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther

}

// MonthAbbreviated returns the locales abbreviated month given the 'month' provided
func (de *de_CH) MonthAbbreviated(month time.Month) string {
	return de.monthsAbbreviated[month]
}

// MonthsAbbreviated returns the locales abbreviated months
func (de *de_CH) MonthsAbbreviated() []string {
	return de.monthsAbbreviated[1:]
}

// MonthNarrow returns the locales narrow month given the 'month' provided
func (de *de_CH) MonthNarrow(month time.Month) string {
	return de.monthsNarrow[month]
}

// MonthsNarrow returns the locales narrow months
func (de *de_CH) MonthsNarrow() []string {
	return de.monthsNarrow[1:]
}

// MonthWide returns the locales wide month given the 'month' provided
func (de *de_CH) MonthWide(month time.Month) string {
	return de.monthsWide[month]
}

// MonthsWide returns the locales wide months
func (de *de_CH) MonthsWide() []string {
	return de.monthsWide[1:]
}

// WeekdayAbbreviated returns the locales abbreviated weekday given the 'weekday' provided
func (de *de_CH) WeekdayAbbreviated(weekday time.Weekday) string {
	return de.daysAbbreviated[weekday]
}

// WeekdaysAbbreviated returns the locales abbreviated weekdays
func (de *de_CH) WeekdaysAbbreviated() []string {
	return de.daysAbbreviated
}

// WeekdayNarrow returns the locales narrow weekday given the 'weekday' provided
func (de *de_CH) WeekdayNarrow(weekday time.Weekday) string {
	return de.daysNarrow[weekday]
}

// WeekdaysNarrow returns the locales narrow weekdays
func (de *de_CH) WeekdaysNarrow() []string {
	return de.daysNarrow
}

// WeekdayShort returns the locales short weekday given the 'weekday' provided
func (de *de_CH) WeekdayShort(weekday time.Weekday) string {
	return de.daysShort[weekday]
}

// WeekdaysShort returns the locales short weekdays
func (de *de_CH) WeekdaysShort() []string {
	return de.daysShort
}

// WeekdayWide returns the locales wide weekday given the 'weekday' provided
func (de *de_CH) WeekdayWide(weekday time.Weekday) string {
	return de.daysWide[weekday]
}

// WeekdaysWide returns the locales wide weekdays
func (de *de_CH) WeekdaysWide() []string {
	return de.daysWide
}

// FmtNumber returns 'num' with digits/precision of 'v' for 'de_CH' and handles both Whole and Real numbers based on 'v'
func (de *de_CH) FmtNumber(num float64, v uint64) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 2 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, de.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, de.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, de.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

// FmtPercent returns 'num' with digits/precision of 'v' for 'de_CH' and handles both Whole and Real numbers based on 'v'
// NOTE: 'num' passed into FmtPercent is assumed to be in percent already
func (de *de_CH) FmtPercent(num float64, v uint64) string {
	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 3
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, de.decimal[0])
			continue
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, de.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	b = append(b, de.percent...)

	return string(b)
}

// FmtCurrency returns the currency representation of 'num' with digits/precision of 'v' for 'de_CH'
func (de *de_CH) FmtCurrency(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := de.currencies[currency]
	l := len(s) + len(symbol) + 4 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, de.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, de.group[0])
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

	for j := len(de.currencyPositivePrefix) - 1; j >= 0; j-- {
		b = append(b, de.currencyPositivePrefix[j])
	}

	if num < 0 {
		b = append(b, de.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, de.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	return string(b)
}

// FmtAccounting returns the currency representation of 'num' with digits/precision of 'v' for 'de_CH'
// in accounting notation.
func (de *de_CH) FmtAccounting(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := de.currencies[currency]
	l := len(s) + len(symbol) + 4 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, de.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, de.group[0])
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

		for j := len(de.currencyNegativePrefix) - 1; j >= 0; j-- {
			b = append(b, de.currencyNegativePrefix[j])
		}

		b = append(b, de.minus[0])

	} else {

		for j := len(symbol) - 1; j >= 0; j-- {
			b = append(b, symbol[j])
		}

		for j := len(de.currencyPositivePrefix) - 1; j >= 0; j-- {
			b = append(b, de.currencyPositivePrefix[j])
		}

	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, de.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	return string(b)
}

// FmtDateShort returns the short date representation of 't' for 'de_CH'
func (de *de_CH) FmtDateShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Day() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2e}...)

	if t.Month() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Month()), 10)

	b = append(b, []byte{0x2e}...)

	if t.Year() > 9 {
		b = append(b, strconv.Itoa(t.Year())[2:]...)
	} else {
		b = append(b, strconv.Itoa(t.Year())[1:]...)
	}

	return string(b)
}

// FmtDateMedium returns the medium date representation of 't' for 'de_CH'
func (de *de_CH) FmtDateMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Day() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2e}...)

	if t.Month() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Month()), 10)

	b = append(b, []byte{0x2e}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateLong returns the long date representation of 't' for 'de_CH'
func (de *de_CH) FmtDateLong(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2e, 0x20}...)
	b = append(b, de.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateFull returns the full date representation of 't' for 'de_CH'
func (de *de_CH) FmtDateFull(t time.Time) string {

	b := make([]byte, 0, 32)

	b = append(b, de.daysWide[t.Weekday()]...)
	b = append(b, []byte{0x2c, 0x20}...)
	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2e, 0x20}...)
	b = append(b, de.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtTimeShort returns the short time representation of 't' for 'de_CH'
func (de *de_CH) FmtTimeShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, de.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)

	return string(b)
}

// FmtTimeMedium returns the medium time representation of 't' for 'de_CH'
func (de *de_CH) FmtTimeMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, de.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, de.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)

	return string(b)
}

// FmtTimeLong returns the long time representation of 't' for 'de_CH'
func (de *de_CH) FmtTimeLong(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, de.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, de.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()
	b = append(b, tz...)

	return string(b)
}

// FmtTimeFull returns the full time representation of 't' for 'de_CH'
func (de *de_CH) FmtTimeFull(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, de.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, de.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()

	if btz, ok := de.timezones[tz]; ok {
		b = append(b, btz...)
	} else {
		b = append(b, tz...)
	}

	return string(b)
}
