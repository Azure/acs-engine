package naq

import (
	"math"
	"strconv"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/currency"
)

type naq struct {
	locale             string
	pluralsCardinal    []locales.PluralRule
	pluralsOrdinal     []locales.PluralRule
	pluralsRange       []locales.PluralRule
	decimal            string
	group              string
	minus              string
	percent            string
	perMille           string
	timeSeparator      string
	inifinity          string
	currencies         []string // idx = enum of currency code
	monthsAbbreviated  []string
	monthsNarrow       []string
	monthsWide         []string
	daysAbbreviated    []string
	daysNarrow         []string
	daysShort          []string
	daysWide           []string
	periodsAbbreviated []string
	periodsNarrow      []string
	periodsShort       []string
	periodsWide        []string
	erasAbbreviated    []string
	erasNarrow         []string
	erasWide           []string
	timezones          map[string]string
}

// New returns a new instance of translator for the 'naq' locale
func New() locales.Translator {
	return &naq{
		locale:             "naq",
		pluralsCardinal:    []locales.PluralRule{2, 3, 6},
		pluralsOrdinal:     nil,
		pluralsRange:       nil,
		timeSeparator:      ":",
		currencies:         []string{"ADP", "AED", "AFA", "AFN", "ALK", "ALL", "AMD", "ANG", "AOA", "AOK", "AON", "AOR", "ARA", "ARL", "ARM", "ARP", "ARS", "ATS", "AUD", "AWG", "AZM", "AZN", "BAD", "BAM", "BAN", "BBD", "BDT", "BEC", "BEF", "BEL", "BGL", "BGM", "BGN", "BGO", "BHD", "BIF", "BMD", "BND", "BOB", "BOL", "BOP", "BOV", "BRB", "BRC", "BRE", "BRL", "BRN", "BRR", "BRZ", "BSD", "BTN", "BUK", "BWP", "BYB", "BYN", "BYR", "BZD", "CAD", "CDF", "CHE", "CHF", "CHW", "CLE", "CLF", "CLP", "CNX", "CNY", "COP", "COU", "CRC", "CSD", "CSK", "CUC", "CUP", "CVE", "CYP", "CZK", "DDM", "DEM", "DJF", "DKK", "DOP", "DZD", "ECS", "ECV", "EEK", "EGP", "ERN", "ESA", "ESB", "ESP", "ETB", "EUR", "FIM", "FJD", "FKP", "FRF", "GBP", "GEK", "GEL", "GHC", "GHS", "GIP", "GMD", "GNF", "GNS", "GQE", "GRD", "GTQ", "GWE", "GWP", "GYD", "HKD", "HNL", "HRD", "HRK", "HTG", "HUF", "IDR", "IEP", "ILP", "ILR", "ILS", "INR", "IQD", "IRR", "ISJ", "ISK", "ITL", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRH", "KRO", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LTT", "LUC", "LUF", "LUL", "LVL", "LVR", "LYD", "MAD", "MAF", "MCF", "MDC", "MDL", "MGA", "MGF", "MKD", "MKN", "MLF", "MMK", "MNT", "MOP", "MRO", "MTL", "MTP", "MUR", "MVP", "MVR", "MWK", "MXN", "MXP", "MXV", "MYR", "MZE", "MZM", "MZN", "$", "NGN", "NIC", "NIO", "NLG", "NOK", "NPR", "NZD", "OMR", "PAB", "PEI", "PEN", "PES", "PGK", "PHP", "PKR", "PLN", "PLZ", "PTE", "PYG", "QAR", "RHD", "ROL", "RON", "RSD", "RUB", "RUR", "RWF", "SAR", "SBD", "SCR", "SDD", "SDG", "SDP", "SEK", "SGD", "SHP", "SIT", "SKK", "SLL", "SOS", "SRD", "SRG", "SSP", "STD", "SUR", "SVC", "SYP", "SZL", "THB", "TJR", "TJS", "TMM", "TMT", "TND", "TOP", "TPE", "TRL", "TRY", "TTD", "TWD", "TZS", "UAH", "UAK", "UGS", "UGX", "USD", "USN", "USS", "UYI", "UYP", "UYU", "UZS", "VEB", "VEF", "VND", "VNN", "VUV", "WST", "XAF", "XAG", "XAU", "XBA", "XBB", "XBC", "XBD", "XCD", "XDR", "XEU", "XFO", "XFU", "XOF", "XPD", "XPF", "XPT", "XRE", "XSU", "XTS", "XUA", "XXX", "YDD", "YER", "YUD", "YUM", "YUN", "YUR", "ZAL", "ZAR", "ZMK", "ZMW", "ZRN", "ZRZ", "ZWD", "ZWL", "ZWR"},
		monthsAbbreviated:  []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		monthsNarrow:       []string{"", "J", "F", "M", "A", "M", "J", "J", "A", "S", "O", "N", "D"},
		monthsWide:         []string{"", "ǃKhanni", "ǃKhanǀgôab", "ǀKhuuǁkhâb", "ǃHôaǂkhaib", "ǃKhaitsâb", "Gamaǀaeb", "ǂKhoesaob", "Aoǁkhuumûǁkhâb", "Taraǀkhuumûǁkhâb", "ǂNûǁnâiseb", "ǀHooǂgaeb", "Hôasoreǁkhâb"},
		daysAbbreviated:    []string{"Son", "Ma", "De", "Wu", "Do", "Fr", "Sat"},
		daysNarrow:         []string{"S", "M", "E", "W", "D", "F", "A"},
		daysWide:           []string{"Sontaxtsees", "Mantaxtsees", "Denstaxtsees", "Wunstaxtsees", "Dondertaxtsees", "Fraitaxtsees", "Satertaxtsees"},
		periodsAbbreviated: []string{"ǁgoagas", "ǃuias"},
		periodsWide:        []string{"ǁgoagas", "ǃuias"},
		erasAbbreviated:    []string{"BC", "AD"},
		erasNarrow:         []string{"", ""},
		erasWide:           []string{"Xristub aiǃâ", "Xristub khaoǃgâ"},
		timezones:          map[string]string{"HECU": "HECU", "CDT": "CDT", "AWDT": "AWDT", "ADT": "ADT", "HAT": "HAT", "AST": "AST", "GMT": "GMT", "ART": "ART", "GYT": "GYT", "EAT": "EAT", "AWST": "AWST", "HADT": "HADT", "MESZ": "MESZ", "AKST": "AKST", "HEOG": "HEOG", "HNPM": "HNPM", "CHAST": "CHAST", "CHADT": "CHADT", "CAT": "CAT", "ACWST": "ACWST", "HEPMX": "HEPMX", "CLST": "CLST", "WESZ": "WESZ", "ACST": "ACST", "HNNOMX": "HNNOMX", "ChST": "ChST", "HEPM": "HEPM", "HNOG": "HNOG", "EDT": "EDT", "IST": "IST", "CLT": "CLT", "WAT": "WAT", "BT": "BT", "EST": "EST", "HEEG": "HEEG", "UYT": "UYT", "∅∅∅": "∅∅∅", "SGT": "SGT", "HAST": "HAST", "ARST": "ARST", "JST": "JST", "WAST": "WAST", "ACDT": "ACDT", "SAST": "SAST", "WIT": "WIT", "ACWDT": "ACWDT", "WEZ": "WEZ", "TMT": "TMT", "WARST": "WARST", "HKST": "HKST", "COST": "COST", "WIB": "WIB", "CST": "CST", "MEZ": "MEZ", "VET": "VET", "OEZ": "OEZ", "HNT": "HNT", "AKDT": "AKDT", "HNCU": "HNCU", "MYT": "MYT", "WITA": "WITA", "BOT": "BOT", "PST": "PST", "PDT": "PDT", "OESZ": "OESZ", "TMST": "TMST", "ECT": "ECT", "JDT": "JDT", "WART": "WART", "COT": "COT", "HNEG": "HNEG", "LHDT": "LHDT", "MDT": "MDT", "AEST": "AEST", "AEDT": "AEDT", "SRT": "SRT", "HKT": "HKT", "GFT": "GFT", "MST": "MST", "HENOMX": "HENOMX", "UYST": "UYST", "LHST": "LHST", "HNPMX": "HNPMX", "NZST": "NZST", "NZDT": "NZDT"},
	}
}

// Locale returns the current translators string locale
func (naq *naq) Locale() string {
	return naq.locale
}

// PluralsCardinal returns the list of cardinal plural rules associated with 'naq'
func (naq *naq) PluralsCardinal() []locales.PluralRule {
	return naq.pluralsCardinal
}

// PluralsOrdinal returns the list of ordinal plural rules associated with 'naq'
func (naq *naq) PluralsOrdinal() []locales.PluralRule {
	return naq.pluralsOrdinal
}

// PluralsRange returns the list of range plural rules associated with 'naq'
func (naq *naq) PluralsRange() []locales.PluralRule {
	return naq.pluralsRange
}

// CardinalPluralRule returns the cardinal PluralRule given 'num' and digits/precision of 'v' for 'naq'
func (naq *naq) CardinalPluralRule(num float64, v uint64) locales.PluralRule {

	n := math.Abs(num)

	if n == 1 {
		return locales.PluralRuleOne
	} else if n == 2 {
		return locales.PluralRuleTwo
	}

	return locales.PluralRuleOther
}

// OrdinalPluralRule returns the ordinal PluralRule given 'num' and digits/precision of 'v' for 'naq'
func (naq *naq) OrdinalPluralRule(num float64, v uint64) locales.PluralRule {
	return locales.PluralRuleUnknown
}

// RangePluralRule returns the ordinal PluralRule given 'num1', 'num2' and digits/precision of 'v1' and 'v2' for 'naq'
func (naq *naq) RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) locales.PluralRule {
	return locales.PluralRuleUnknown
}

// MonthAbbreviated returns the locales abbreviated month given the 'month' provided
func (naq *naq) MonthAbbreviated(month time.Month) string {
	return naq.monthsAbbreviated[month]
}

// MonthsAbbreviated returns the locales abbreviated months
func (naq *naq) MonthsAbbreviated() []string {
	return naq.monthsAbbreviated[1:]
}

// MonthNarrow returns the locales narrow month given the 'month' provided
func (naq *naq) MonthNarrow(month time.Month) string {
	return naq.monthsNarrow[month]
}

// MonthsNarrow returns the locales narrow months
func (naq *naq) MonthsNarrow() []string {
	return naq.monthsNarrow[1:]
}

// MonthWide returns the locales wide month given the 'month' provided
func (naq *naq) MonthWide(month time.Month) string {
	return naq.monthsWide[month]
}

// MonthsWide returns the locales wide months
func (naq *naq) MonthsWide() []string {
	return naq.monthsWide[1:]
}

// WeekdayAbbreviated returns the locales abbreviated weekday given the 'weekday' provided
func (naq *naq) WeekdayAbbreviated(weekday time.Weekday) string {
	return naq.daysAbbreviated[weekday]
}

// WeekdaysAbbreviated returns the locales abbreviated weekdays
func (naq *naq) WeekdaysAbbreviated() []string {
	return naq.daysAbbreviated
}

// WeekdayNarrow returns the locales narrow weekday given the 'weekday' provided
func (naq *naq) WeekdayNarrow(weekday time.Weekday) string {
	return naq.daysNarrow[weekday]
}

// WeekdaysNarrow returns the locales narrow weekdays
func (naq *naq) WeekdaysNarrow() []string {
	return naq.daysNarrow
}

// WeekdayShort returns the locales short weekday given the 'weekday' provided
func (naq *naq) WeekdayShort(weekday time.Weekday) string {
	return naq.daysShort[weekday]
}

// WeekdaysShort returns the locales short weekdays
func (naq *naq) WeekdaysShort() []string {
	return naq.daysShort
}

// WeekdayWide returns the locales wide weekday given the 'weekday' provided
func (naq *naq) WeekdayWide(weekday time.Weekday) string {
	return naq.daysWide[weekday]
}

// WeekdaysWide returns the locales wide weekdays
func (naq *naq) WeekdaysWide() []string {
	return naq.daysWide
}

// FmtNumber returns 'num' with digits/precision of 'v' for 'naq' and handles both Whole and Real numbers based on 'v'
func (naq *naq) FmtNumber(num float64, v uint64) string {

	return strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
}

// FmtPercent returns 'num' with digits/precision of 'v' for 'naq' and handles both Whole and Real numbers based on 'v'
// NOTE: 'num' passed into FmtPercent is assumed to be in percent already
func (naq *naq) FmtPercent(num float64, v uint64) string {
	return strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
}

// FmtCurrency returns the currency representation of 'num' with digits/precision of 'v' for 'naq'
func (naq *naq) FmtCurrency(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := naq.currencies[currency]
	l := len(s) + len(symbol) + 0 + 0*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, naq.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, naq.group[0])
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

	if num < 0 {
		b = append(b, naq.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, naq.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	return string(b)
}

// FmtAccounting returns the currency representation of 'num' with digits/precision of 'v' for 'naq'
// in accounting notation.
func (naq *naq) FmtAccounting(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := naq.currencies[currency]
	l := len(s) + len(symbol) + 0 + 0*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, naq.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, naq.group[0])
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

		b = append(b, naq.minus[0])

	} else {

		for j := len(symbol) - 1; j >= 0; j-- {
			b = append(b, symbol[j])
		}

	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, naq.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	return string(b)
}

// FmtDateShort returns the short date representation of 't' for 'naq'
func (naq *naq) FmtDateShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Day() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2f}...)

	if t.Month() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Month()), 10)

	b = append(b, []byte{0x2f}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateMedium returns the medium date representation of 't' for 'naq'
func (naq *naq) FmtDateMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, naq.monthsAbbreviated[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateLong returns the long date representation of 't' for 'naq'
func (naq *naq) FmtDateLong(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, naq.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateFull returns the full date representation of 't' for 'naq'
func (naq *naq) FmtDateFull(t time.Time) string {

	b := make([]byte, 0, 32)

	b = append(b, naq.daysWide[t.Weekday()]...)
	b = append(b, []byte{0x2c, 0x20}...)
	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, naq.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtTimeShort returns the short time representation of 't' for 'naq'
func (naq *naq) FmtTimeShort(t time.Time) string {

	b := make([]byte, 0, 32)

	h := t.Hour()

	if h > 12 {
		h -= 12
	}

	b = strconv.AppendInt(b, int64(h), 10)
	b = append(b, naq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, []byte{0x20}...)

	if t.Hour() < 12 {
		b = append(b, naq.periodsAbbreviated[0]...)
	} else {
		b = append(b, naq.periodsAbbreviated[1]...)
	}

	return string(b)
}

// FmtTimeMedium returns the medium time representation of 't' for 'naq'
func (naq *naq) FmtTimeMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	h := t.Hour()

	if h > 12 {
		h -= 12
	}

	b = strconv.AppendInt(b, int64(h), 10)
	b = append(b, naq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, naq.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	if t.Hour() < 12 {
		b = append(b, naq.periodsAbbreviated[0]...)
	} else {
		b = append(b, naq.periodsAbbreviated[1]...)
	}

	return string(b)
}

// FmtTimeLong returns the long time representation of 't' for 'naq'
func (naq *naq) FmtTimeLong(t time.Time) string {

	b := make([]byte, 0, 32)

	h := t.Hour()

	if h > 12 {
		h -= 12
	}

	b = strconv.AppendInt(b, int64(h), 10)
	b = append(b, naq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, naq.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	if t.Hour() < 12 {
		b = append(b, naq.periodsAbbreviated[0]...)
	} else {
		b = append(b, naq.periodsAbbreviated[1]...)
	}

	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()
	b = append(b, tz...)

	return string(b)
}

// FmtTimeFull returns the full time representation of 't' for 'naq'
func (naq *naq) FmtTimeFull(t time.Time) string {

	b := make([]byte, 0, 32)

	h := t.Hour()

	if h > 12 {
		h -= 12
	}

	b = strconv.AppendInt(b, int64(h), 10)
	b = append(b, naq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, naq.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	if t.Hour() < 12 {
		b = append(b, naq.periodsAbbreviated[0]...)
	} else {
		b = append(b, naq.periodsAbbreviated[1]...)
	}

	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()

	if btz, ok := naq.timezones[tz]; ok {
		b = append(b, btz...)
	} else {
		b = append(b, tz...)
	}

	return string(b)
}
