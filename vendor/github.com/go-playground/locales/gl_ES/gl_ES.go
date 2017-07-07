package gl_ES

import (
	"math"
	"strconv"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/currency"
)

type gl_ES struct {
	locale                 string
	pluralsCardinal        []locales.PluralRule
	pluralsOrdinal         []locales.PluralRule
	pluralsRange           []locales.PluralRule
	decimal                string
	group                  string
	minus                  string
	percent                string
	percentSuffix          string
	perMille               string
	timeSeparator          string
	inifinity              string
	currencies             []string // idx = enum of currency code
	currencyPositiveSuffix string
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

// New returns a new instance of translator for the 'gl_ES' locale
func New() locales.Translator {
	return &gl_ES{
		locale:                 "gl_ES",
		pluralsCardinal:        []locales.PluralRule{2, 6},
		pluralsOrdinal:         []locales.PluralRule{6},
		pluralsRange:           []locales.PluralRule{2, 6},
		decimal:                ",",
		group:                  ".",
		minus:                  "-",
		percent:                "%",
		perMille:               "‰",
		timeSeparator:          ":",
		inifinity:              "∞",
		currencies:             []string{"ADP", "AED", "AFA", "AFN", "ALK", "ALL", "AMD", "ANG", "AOA", "AOK", "AON", "AOR", "ARA", "ARL", "ARM", "ARP", "ARS", "ATS", "AUD", "AWG", "AZM", "AZN", "BAD", "BAM", "BAN", "BBD", "BDT", "BEC", "BEF", "BEL", "BGL", "BGM", "BGN", "BGO", "BHD", "BIF", "BMD", "BND", "BOB", "BOL", "BOP", "BOV", "BRB", "BRC", "BRE", "BRL", "BRN", "BRR", "BRZ", "BSD", "BTN", "BUK", "BWP", "BYB", "BYN", "BYR", "BZD", "CAD", "CDF", "CHE", "CHF", "CHW", "CLE", "CLF", "CLP", "CNX", "CNY", "COP", "COU", "CRC", "CSD", "CSK", "CUC", "CUP", "CVE", "CYP", "CZK", "DDM", "DEM", "DJF", "DKK", "DOP", "DZD", "ECS", "ECV", "EEK", "EGP", "ERN", "ESA", "ESB", "ESP", "ETB", "EUR", "FIM", "FJD", "FKP", "FRF", "GBP", "GEK", "GEL", "GHC", "GHS", "GIP", "GMD", "GNF", "GNS", "GQE", "GRD", "GTQ", "GWE", "GWP", "GYD", "HKD", "HNL", "HRD", "HRK", "HTG", "HUF", "IDR", "IEP", "ILP", "ILR", "ILS", "INR", "IQD", "IRR", "ISJ", "ISK", "ITL", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRH", "KRO", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LTT", "LUC", "LUF", "LUL", "LVL", "LVR", "LYD", "MAD", "MAF", "MCF", "MDC", "MDL", "MGA", "MGF", "MKD", "MKN", "MLF", "MMK", "MNT", "MOP", "MRO", "MTL", "MTP", "MUR", "MVP", "MVR", "MWK", "MXN", "MXP", "MXV", "MYR", "MZE", "MZM", "MZN", "NAD", "NGN", "NIC", "NIO", "NLG", "NOK", "NPR", "NZD", "OMR", "PAB", "PEI", "PEN", "PES", "PGK", "PHP", "PKR", "PLN", "PLZ", "PTE", "PYG", "QAR", "RHD", "ROL", "RON", "RSD", "RUB", "RUR", "RWF", "SAR", "SBD", "SCR", "SDD", "SDG", "SDP", "SEK", "SGD", "SHP", "SIT", "SKK", "SLL", "SOS", "SRD", "SRG", "SSP", "STD", "SUR", "SVC", "SYP", "SZL", "THB", "TJR", "TJS", "TMM", "TMT", "TND", "TOP", "TPE", "TRL", "TRY", "TTD", "TWD", "TZS", "UAH", "UAK", "UGS", "UGX", "USD", "USN", "USS", "UYI", "UYP", "UYU", "UZS", "VEB", "VEF", "VND", "VNN", "VUV", "WST", "XAF", "XAG", "XAU", "XBA", "XBB", "XBC", "XBD", "XCD", "XDR", "XEU", "XFO", "XFU", "XOF", "XPD", "XPF", "XPT", "XRE", "XSU", "XTS", "XUA", "XXX", "YDD", "YER", "YUD", "YUM", "YUN", "YUR", "ZAL", "ZAR", "ZMK", "ZMW", "ZRN", "ZRZ", "ZWD", "ZWL", "ZWR"},
		percentSuffix:          " ",
		currencyPositiveSuffix: " ",
		currencyNegativeSuffix: " ",
		monthsAbbreviated:      []string{"", "xan.", "feb.", "mar.", "abr.", "maio", "xuño", "xul.", "ago.", "set.", "out.", "nov.", "dec."},
		monthsNarrow:           []string{"", "x.", "f.", "m.", "a.", "m.", "x.", "x.", "a.", "s.", "o.", "n.", "d."},
		monthsWide:             []string{"", "xaneiro", "febreiro", "marzo", "abril", "maio", "xuño", "xullo", "agosto", "setembro", "outubro", "novembro", "decembro"},
		daysAbbreviated:        []string{"dom.", "luns", "mar.", "mér.", "xov.", "ven.", "sáb."},
		daysNarrow:             []string{"d.", "l.", "m.", "m.", "x.", "v.", "s."},
		daysShort:              []string{"dom.", "luns", "mar.", "mér.", "xov.", "ven.", "sáb."},
		daysWide:               []string{"domingo", "luns", "martes", "mércores", "xoves", "venres", "sábado"},
		periodsAbbreviated:     []string{"a.m.", "p.m."},
		periodsNarrow:          []string{"a.m.", "p.m."},
		periodsWide:            []string{"a.m.", "p.m."},
		erasAbbreviated:        []string{"a.C.", "d.C."},
		erasNarrow:             []string{"", ""},
		erasWide:               []string{"antes de Cristo", "despois de Cristo"},
		timezones:              map[string]string{"PDT": "Horario de verán do Pacífico, Norteamérica", "HAST": "Horario estándar de Hawai-Aleutiano", "TMT": "Horario estándar de Turcomenistán", "BT": "Horario de Bután", "ChST": "Horario estándar chamorro", "SAST": "Horario estándar de África do Sur", "WIT": "Horario de Indonesia Oriental", "PST": "Horario estándar do Pacífico, Norteamérica", "OEZ": "Horario estándar de Europa Oriental", "WEZ": "Horario estándar de Europa Occidental", "WESZ": "Horario de verán de Europa Occidental", "HNOG": "Horario estándar de Groenlandia Occidental", "HAT": "Horario de verán de Terranova", "AEDT": "Horario de verán de Australia Oriental", "CDT": "Horario de verán central, Norteamérica", "VET": "Horario de Venezuela", "CLST": "Horario de verán de Chile", "ADT": "Horario de verán do Atlántico", "WIB": "Horario de Indonesia Occidental", "MEZ": "Horario estándar de Europa Central", "MDT": "MDT", "EDT": "Horario de verán do leste, Norteamérica", "HNEG": "Horario estándar de Groenlandia Oriental", "AEST": "Horario estándar de Australia Oriental", "HEPM": "Horario de verán de Saint Pierre e Miquelon", "EAT": "Horario de África Oriental", "MST": "MST", "COT": "Horario estándar de Colombia", "CHADT": "Horario de verán de Chatham", "WART": "Horario estándar de Arxentina Occidental", "WAST": "Horario de verán de África Occidental", "HNT": "Horario estándar de Terranova", "∅∅∅": "Horario de verán de Brasilia", "IST": "Horario estándar da India", "COST": "Horario de verán de Colombia", "GFT": "Horario da Güiana Francesa", "UYT": "Horario estándar de Uruguai", "ACWDT": "Horario de verán de Australia Occidental Central", "NZST": "Horario estándar de Nova Zelandia", "OESZ": "Horario de verán de Europa Oriental", "HNNOMX": "Horario estándar de México Noroeste", "ACST": "Horario estándar de Australia Central", "WITA": "Horario de Indonesia Central", "HEPMX": "Horario de verán do Pacífico mexicano", "SGT": "Horario estándar de Singapur", "ACWST": "Horario estándar de Australia Occidental Central", "HKT": "Horario estándar de Hong Kong", "HENOMX": "Horario de verán de México Noroeste", "ACDT": "Horario de verán de Australia Central", "HEEG": "Horario de verán de Groenlandia Oriental", "SRT": "Horario de Suriname", "ECT": "Horario de Ecuador", "NZDT": "Horario de verán de Nova Zelandia", "HKST": "Horario de verán de Hong Kong", "WAT": "Horario estándar de África Occidental", "CHAST": "Horario estándar de Chatham", "MESZ": "Horario de verán de Europa Central", "CAT": "Horario de África Central", "JDT": "Horario de verán do Xapón", "HEOG": "Horario de verán de Groenlandia Occidental", "HNPM": "Horario estándar de Saint Pierre e Miquelon", "HNPMX": "Horario estándar do Pacífico mexicano", "AWDT": "Horario de verán de Australia Occidental", "CLT": "Horario estándar de Chile", "EST": "Horario estándar do leste, Norteamérica", "WARST": "Horario de verán de Arxentina Occidental", "ART": "Horario estándar de Arxentina", "HADT": "Horario de verán de Hawai-Aleutiano", "GMT": "Horario do meridiano de Greenwich", "AWST": "Horario estándar de Australia Occidental", "BOT": "Horario de Bolivia", "MYT": "Horario de Malaisia", "ARST": "Horario de verán de Arxentina", "AKDT": "Horario de verán de Alasca", "UYST": "Horario de verán de Uruguai", "HNCU": "Horario estándar de Cuba", "CST": "Horario estándar central, Norteamérica", "LHST": "Horario estándar de Lord Howe", "GYT": "Horario da Güiana", "AST": "Horario estándar do Atlántico", "TMST": "Horario de verán de Turcomenistán", "AKST": "Horario estándar de Alasca", "LHDT": "Horario de verán de Lord Howe", "HECU": "Horario de verán de Cuba", "JST": "Horario estándar do Xapón"},
	}
}

// Locale returns the current translators string locale
func (gl *gl_ES) Locale() string {
	return gl.locale
}

// PluralsCardinal returns the list of cardinal plural rules associated with 'gl_ES'
func (gl *gl_ES) PluralsCardinal() []locales.PluralRule {
	return gl.pluralsCardinal
}

// PluralsOrdinal returns the list of ordinal plural rules associated with 'gl_ES'
func (gl *gl_ES) PluralsOrdinal() []locales.PluralRule {
	return gl.pluralsOrdinal
}

// PluralsRange returns the list of range plural rules associated with 'gl_ES'
func (gl *gl_ES) PluralsRange() []locales.PluralRule {
	return gl.pluralsRange
}

// CardinalPluralRule returns the cardinal PluralRule given 'num' and digits/precision of 'v' for 'gl_ES'
func (gl *gl_ES) CardinalPluralRule(num float64, v uint64) locales.PluralRule {

	n := math.Abs(num)
	i := int64(n)

	if i == 1 && v == 0 {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther
}

// OrdinalPluralRule returns the ordinal PluralRule given 'num' and digits/precision of 'v' for 'gl_ES'
func (gl *gl_ES) OrdinalPluralRule(num float64, v uint64) locales.PluralRule {
	return locales.PluralRuleOther
}

// RangePluralRule returns the ordinal PluralRule given 'num1', 'num2' and digits/precision of 'v1' and 'v2' for 'gl_ES'
func (gl *gl_ES) RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) locales.PluralRule {

	start := gl.CardinalPluralRule(num1, v1)
	end := gl.CardinalPluralRule(num2, v2)

	if start == locales.PluralRuleOne && end == locales.PluralRuleOther {
		return locales.PluralRuleOther
	} else if start == locales.PluralRuleOther && end == locales.PluralRuleOne {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther

}

// MonthAbbreviated returns the locales abbreviated month given the 'month' provided
func (gl *gl_ES) MonthAbbreviated(month time.Month) string {
	return gl.monthsAbbreviated[month]
}

// MonthsAbbreviated returns the locales abbreviated months
func (gl *gl_ES) MonthsAbbreviated() []string {
	return gl.monthsAbbreviated[1:]
}

// MonthNarrow returns the locales narrow month given the 'month' provided
func (gl *gl_ES) MonthNarrow(month time.Month) string {
	return gl.monthsNarrow[month]
}

// MonthsNarrow returns the locales narrow months
func (gl *gl_ES) MonthsNarrow() []string {
	return gl.monthsNarrow[1:]
}

// MonthWide returns the locales wide month given the 'month' provided
func (gl *gl_ES) MonthWide(month time.Month) string {
	return gl.monthsWide[month]
}

// MonthsWide returns the locales wide months
func (gl *gl_ES) MonthsWide() []string {
	return gl.monthsWide[1:]
}

// WeekdayAbbreviated returns the locales abbreviated weekday given the 'weekday' provided
func (gl *gl_ES) WeekdayAbbreviated(weekday time.Weekday) string {
	return gl.daysAbbreviated[weekday]
}

// WeekdaysAbbreviated returns the locales abbreviated weekdays
func (gl *gl_ES) WeekdaysAbbreviated() []string {
	return gl.daysAbbreviated
}

// WeekdayNarrow returns the locales narrow weekday given the 'weekday' provided
func (gl *gl_ES) WeekdayNarrow(weekday time.Weekday) string {
	return gl.daysNarrow[weekday]
}

// WeekdaysNarrow returns the locales narrow weekdays
func (gl *gl_ES) WeekdaysNarrow() []string {
	return gl.daysNarrow
}

// WeekdayShort returns the locales short weekday given the 'weekday' provided
func (gl *gl_ES) WeekdayShort(weekday time.Weekday) string {
	return gl.daysShort[weekday]
}

// WeekdaysShort returns the locales short weekdays
func (gl *gl_ES) WeekdaysShort() []string {
	return gl.daysShort
}

// WeekdayWide returns the locales wide weekday given the 'weekday' provided
func (gl *gl_ES) WeekdayWide(weekday time.Weekday) string {
	return gl.daysWide[weekday]
}

// WeekdaysWide returns the locales wide weekdays
func (gl *gl_ES) WeekdaysWide() []string {
	return gl.daysWide
}

// FmtNumber returns 'num' with digits/precision of 'v' for 'gl_ES' and handles both Whole and Real numbers based on 'v'
func (gl *gl_ES) FmtNumber(num float64, v uint64) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 2 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, gl.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, gl.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, gl.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

// FmtPercent returns 'num' with digits/precision of 'v' for 'gl_ES' and handles both Whole and Real numbers based on 'v'
// NOTE: 'num' passed into FmtPercent is assumed to be in percent already
func (gl *gl_ES) FmtPercent(num float64, v uint64) string {
	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 5
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, gl.decimal[0])
			continue
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, gl.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	b = append(b, gl.percentSuffix...)

	b = append(b, gl.percent...)

	return string(b)
}

// FmtCurrency returns the currency representation of 'num' with digits/precision of 'v' for 'gl_ES'
func (gl *gl_ES) FmtCurrency(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := gl.currencies[currency]
	l := len(s) + len(symbol) + 4 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, gl.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, gl.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, gl.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, gl.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	b = append(b, gl.currencyPositiveSuffix...)

	b = append(b, symbol...)

	return string(b)
}

// FmtAccounting returns the currency representation of 'num' with digits/precision of 'v' for 'gl_ES'
// in accounting notation.
func (gl *gl_ES) FmtAccounting(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := gl.currencies[currency]
	l := len(s) + len(symbol) + 4 + 1*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, gl.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				b = append(b, gl.group[0])
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {

		b = append(b, gl.minus[0])

	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, gl.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	if num < 0 {
		b = append(b, gl.currencyNegativeSuffix...)
		b = append(b, symbol...)
	} else {

		b = append(b, gl.currencyPositiveSuffix...)
		b = append(b, symbol...)
	}

	return string(b)
}

// FmtDateShort returns the short date representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtDateShort(t time.Time) string {

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

	if t.Year() > 9 {
		b = append(b, strconv.Itoa(t.Year())[2:]...)
	} else {
		b = append(b, strconv.Itoa(t.Year())[1:]...)
	}

	return string(b)
}

// FmtDateMedium returns the medium date representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtDateMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20, 0x64, 0x65}...)
	b = append(b, []byte{0x20}...)
	b = append(b, gl.monthsAbbreviated[t.Month()]...)
	b = append(b, []byte{0x20, 0x64, 0x65}...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateLong returns the long date representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtDateLong(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20, 0x64, 0x65}...)
	b = append(b, []byte{0x20}...)
	b = append(b, gl.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20, 0x64, 0x65}...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateFull returns the full date representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtDateFull(t time.Time) string {

	b := make([]byte, 0, 32)

	b = append(b, gl.daysWide[t.Weekday()]...)
	b = append(b, []byte{0x2c, 0x20}...)
	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20, 0x64, 0x65}...)
	b = append(b, []byte{0x20}...)
	b = append(b, gl.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20, 0x64, 0x65}...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtTimeShort returns the short time representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtTimeShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)

	return string(b)
}

// FmtTimeMedium returns the medium time representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtTimeMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)

	return string(b)
}

// FmtTimeLong returns the long time representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtTimeLong(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()
	b = append(b, tz...)

	return string(b)
}

// FmtTimeFull returns the full time representation of 't' for 'gl_ES'
func (gl *gl_ES) FmtTimeFull(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, gl.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()

	if btz, ok := gl.timezones[tz]; ok {
		b = append(b, btz...)
	} else {
		b = append(b, tz...)
	}

	return string(b)
}
