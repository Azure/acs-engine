package sq_MK

import (
	"math"
	"strconv"
	"time"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/currency"
)

type sq_MK struct {
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
	currencyPositiveSuffix string
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

// New returns a new instance of translator for the 'sq_MK' locale
func New() locales.Translator {
	return &sq_MK{
		locale:                 "sq_MK",
		pluralsCardinal:        []locales.PluralRule{2, 6},
		pluralsOrdinal:         []locales.PluralRule{2, 5, 6},
		pluralsRange:           []locales.PluralRule{2, 6},
		decimal:                ",",
		group:                  " ",
		minus:                  "-",
		percent:                "%",
		perMille:               "‰",
		timeSeparator:          ":",
		inifinity:              "∞",
		currencies:             []string{"ADP", "AED", "AFA", "AFN", "ALK", "ALL", "AMD", "ANG", "AOA", "AOK", "AON", "AOR", "ARA", "ARL", "ARM", "ARP", "ARS", "ATS", "AUD", "AWG", "AZM", "AZN", "BAD", "BAM", "BAN", "BBD", "BDT", "BEC", "BEF", "BEL", "BGL", "BGM", "BGN", "BGO", "BHD", "BIF", "BMD", "BND", "BOB", "BOL", "BOP", "BOV", "BRB", "BRC", "BRE", "BRL", "BRN", "BRR", "BRZ", "BSD", "BTN", "BUK", "BWP", "BYB", "BYN", "BYR", "BZD", "CAD", "CDF", "CHE", "CHF", "CHW", "CLE", "CLF", "CLP", "CNX", "CNY", "COP", "COU", "CRC", "CSD", "CSK", "CUC", "CUP", "CVE", "CYP", "CZK", "DDM", "DEM", "DJF", "DKK", "DOP", "DZD", "ECS", "ECV", "EEK", "EGP", "ERN", "ESA", "ESB", "ESP", "ETB", "EUR", "FIM", "FJD", "FKP", "FRF", "GBP", "GEK", "GEL", "GHC", "GHS", "GIP", "GMD", "GNF", "GNS", "GQE", "GRD", "GTQ", "GWE", "GWP", "GYD", "HKD", "HNL", "HRD", "HRK", "HTG", "HUF", "IDR", "IEP", "ILP", "ILR", "ILS", "INR", "IQD", "IRR", "ISJ", "ISK", "ITL", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRH", "KRO", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LTT", "LUC", "LUF", "LUL", "LVL", "LVR", "LYD", "MAD", "MAF", "MCF", "MDC", "MDL", "MGA", "MGF", "den", "MKN", "MLF", "MMK", "MNT", "MOP", "MRO", "MTL", "MTP", "MUR", "MVP", "MVR", "MWK", "MXN", "MXP", "MXV", "MYR", "MZE", "MZM", "MZN", "NAD", "NGN", "NIC", "NIO", "NLG", "NOK", "NPR", "NZD", "OMR", "PAB", "PEI", "PEN", "PES", "PGK", "PHP", "PKR", "PLN", "PLZ", "PTE", "PYG", "QAR", "RHD", "ROL", "RON", "RSD", "RUB", "RUR", "RWF", "SAR", "SBD", "SCR", "SDD", "SDG", "SDP", "SEK", "SGD", "SHP", "SIT", "SKK", "SLL", "SOS", "SRD", "SRG", "SSP", "STD", "SUR", "SVC", "SYP", "SZL", "THB", "TJR", "TJS", "TMM", "TMT", "TND", "TOP", "TPE", "TRL", "TRY", "TTD", "TWD", "TZS", "UAH", "UAK", "UGS", "UGX", "USD", "USN", "USS", "UYI", "UYP", "UYU", "UZS", "VEB", "VEF", "VND", "VNN", "VUV", "WST", "XAF", "XAG", "XAU", "XBA", "XBB", "XBC", "XBD", "XCD", "XDR", "XEU", "XFO", "XFU", "XOF", "XPD", "XPF", "XPT", "XRE", "XSU", "XTS", "XUA", "XXX", "YDD", "YER", "YUD", "YUM", "YUN", "YUR", "ZAL", "ZAR", "ZMK", "ZMW", "ZRN", "ZRZ", "ZWD", "ZWL", "ZWR"},
		currencyPositiveSuffix: " ",
		currencyNegativePrefix: "(",
		currencyNegativeSuffix: " )",
		monthsAbbreviated:      []string{"", "jan", "shk", "mar", "pri", "maj", "qer", "kor", "gsh", "sht", "tet", "nën", "dhj"},
		monthsNarrow:           []string{"", "j", "s", "m", "p", "m", "q", "k", "g", "s", "t", "n", "d"},
		monthsWide:             []string{"", "janar", "shkurt", "mars", "prill", "maj", "qershor", "korrik", "gusht", "shtator", "tetor", "nëntor", "dhjetor"},
		daysAbbreviated:        []string{"Die", "Hën", "Mar", "Mër", "Enj", "Pre", "Sht"},
		daysNarrow:             []string{"D", "H", "M", "M", "E", "P", "S"},
		daysShort:              []string{"Die", "Hën", "Mar", "Mër", "Enj", "Pre", "Sht"},
		daysWide:               []string{"e diel", "e hënë", "e martë", "e mërkurë", "e enjte", "e premte", "e shtunë"},
		periodsAbbreviated:     []string{"e paradites", "e pasdites"},
		periodsNarrow:          []string{"e paradites", "e pasdites"},
		periodsWide:            []string{"e paradites", "e pasdites"},
		erasAbbreviated:        []string{"p.K.", "mb.K."},
		erasNarrow:             []string{"p.K.", "mb.K."},
		erasWide:               []string{"para Krishtit", "mbas Krishtit"},
		timezones:              map[string]string{"MDT": "Ora verore amerikane e Brezit Malor", "HNOG": "Ora standarde e Grenlandës Perëndimore", "IST": "Ora standarde e Indisë", "ACWDT": "Ora verore e Australisë Qendroro-Perëndimore", "HNNOMX": "Ora standarde e Meksikës Veriperëndimore", "MST": "Ora standarde amerikane e Brezit Malor", "UYST": "Ora verore e Uruguait", "ADT": "Ora verore e Atlantikut", "SRT": "Ora e Surinamit", "MEZ": "Ora standarde e Evropës Qendrore", "WART": "Ora standarde e Argjentinës Perëndimore", "AKST": "Ora standarde e Alaskës", "CAT": "Ora e Afrikës Qendrore", "WESZ": "Ora verore e Evropës Perëndimore", "CLST": "Ora verore e Kilit", "AST": "Ora standarde e Atlantikut", "ACST": "Ora standarde e Australisë Qendrore", "HEEG": "Ora verore e Grenlandës Lindore", "WITA": "Ora e Indonezisë Qendrore", "EDT": "Ora verore e SHBA-së Lindore", "EAT": "Ora e Afrikës Lindore", "OESZ": "Ora verore e Evropës Lindore", "WARST": "Ora verore e Argjentinës Perëndimore", "WAST": "Ora verore e Afrikës Perëndimore", "GFT": "Ora e Guajanës Franceze", "HEPM": "Ora verore e Shën-Pier dhe Mikuelon", "HECU": "Ora verore e Kubës", "MYT": "Ora e Malajzisë", "WAT": "Ora standarde e Afrikës Perëndimore", "ACDT": "Ora verore e Australisë Qendrore", "CHADT": "Ora verore e Katamit", "HADT": "Ora verore e Ishujve Hauai-Aleutian", "JDT": "Ora verore e Japonisë", "HENOMX": "Ora verore e Meksikës Veriperëndimore", "HNPM": "Ora standarde e Shën-Pier dhe Mikuelon", "HNPMX": "Ora standarde meksikane e Bregut të Paqësorit", "∅∅∅": "Ora verore e Ejkrit [Ako]", "LHST": "Ora standarde e Lord-Houit", "CLT": "Ora standarde e Kilit", "HEOG": "Ora verore e Grenlandës Perëndimore", "TMST": "Ora verore e Turkmenistanit", "SGT": "Ora e Singaporit", "OEZ": "Ora standarde e Evropës Lindore", "TMT": "Ora standarde e Turkmenistanit", "ARST": "Ora verore e Argjentinës", "HKST": "Ora verore e Hong-Kongut", "AEDT": "Ora verore e Australisë Lindore", "UYT": "Ora standarde e Uruguait", "HEPMX": "Ora verore meksikane e Bregut të Paqësorit", "GYT": "Ora e Guajanës", "CHAST": "Ora standarde e Katamit", "AWST": "Ora standarde e Australisë Perëndimore", "AWDT": "Ora verore e Australisë Perëndimore", "JST": "Ora standarde e Japonisë", "GMT": "Ora e Meridianit të Grinuiçit", "ChST": "Ora e Kamorros", "CST": "Ora standarde e SHBA-së Qendrore", "CDT": "Ora verore e SHBA-së Qendrore", "HNCU": "Ora standarde e Kubës", "BOT": "Ora e Bolivisë", "ACWST": "Ora standarde e Australisë Qendroro-Perëndimore", "NZST": "Ora standarde e Zelandës së Re", "WEZ": "Ora standarde e Evropës Perëndimore", "COT": "Ora standarde e Kolumbisë", "COST": "Ora verore e Kolumbisë", "HAT": "Ora verore e Njufaundlendit [Tokës së Re]", "HKT": "Ora standarde e Hong-Kongut", "EST": "Ora standarde e SHBA-së Lindore", "PST": "Ora standarde amerikane e Bregut të Paqësorit", "PDT": "Ora verore amerikane e Bregut të Paqësorit", "MESZ": "Ora verore e Evropës Qendrore", "HAST": "Ora standarde e Ishujve Hauai-Aleutian", "NZDT": "Ora verore e Zelandës së Re", "HNEG": "Ora standarde e Grenlandës Lindore", "AEST": "Ora standarde e Australisë Lindore", "SAST": "Ora standarde e Afrikës Jugore", "LHDT": "Ora verore e Lord-Houit", "WIB": "Ora e Indonezisë Perëndimore", "WIT": "Ora e Indonezisë Lindore", "ECT": "Ora e Ekuadorit", "VET": "Ora e Venezuelës", "HNT": "Ora standarde e Njufaundlendit [Tokës së Re]", "BT": "Ora e Butanit", "AKDT": "Ora verore e Alsaskës", "ART": "Ora standarde e Argjentinës"},
	}
}

// Locale returns the current translators string locale
func (sq *sq_MK) Locale() string {
	return sq.locale
}

// PluralsCardinal returns the list of cardinal plural rules associated with 'sq_MK'
func (sq *sq_MK) PluralsCardinal() []locales.PluralRule {
	return sq.pluralsCardinal
}

// PluralsOrdinal returns the list of ordinal plural rules associated with 'sq_MK'
func (sq *sq_MK) PluralsOrdinal() []locales.PluralRule {
	return sq.pluralsOrdinal
}

// PluralsRange returns the list of range plural rules associated with 'sq_MK'
func (sq *sq_MK) PluralsRange() []locales.PluralRule {
	return sq.pluralsRange
}

// CardinalPluralRule returns the cardinal PluralRule given 'num' and digits/precision of 'v' for 'sq_MK'
func (sq *sq_MK) CardinalPluralRule(num float64, v uint64) locales.PluralRule {

	n := math.Abs(num)

	if n == 1 {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther
}

// OrdinalPluralRule returns the ordinal PluralRule given 'num' and digits/precision of 'v' for 'sq_MK'
func (sq *sq_MK) OrdinalPluralRule(num float64, v uint64) locales.PluralRule {

	n := math.Abs(num)
	nMod10 := math.Mod(n, 10)
	nMod100 := math.Mod(n, 100)

	if n == 1 {
		return locales.PluralRuleOne
	} else if nMod10 == 4 && nMod100 != 14 {
		return locales.PluralRuleMany
	}

	return locales.PluralRuleOther
}

// RangePluralRule returns the ordinal PluralRule given 'num1', 'num2' and digits/precision of 'v1' and 'v2' for 'sq_MK'
func (sq *sq_MK) RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) locales.PluralRule {

	start := sq.CardinalPluralRule(num1, v1)
	end := sq.CardinalPluralRule(num2, v2)

	if start == locales.PluralRuleOne && end == locales.PluralRuleOther {
		return locales.PluralRuleOther
	} else if start == locales.PluralRuleOther && end == locales.PluralRuleOne {
		return locales.PluralRuleOne
	}

	return locales.PluralRuleOther

}

// MonthAbbreviated returns the locales abbreviated month given the 'month' provided
func (sq *sq_MK) MonthAbbreviated(month time.Month) string {
	return sq.monthsAbbreviated[month]
}

// MonthsAbbreviated returns the locales abbreviated months
func (sq *sq_MK) MonthsAbbreviated() []string {
	return sq.monthsAbbreviated[1:]
}

// MonthNarrow returns the locales narrow month given the 'month' provided
func (sq *sq_MK) MonthNarrow(month time.Month) string {
	return sq.monthsNarrow[month]
}

// MonthsNarrow returns the locales narrow months
func (sq *sq_MK) MonthsNarrow() []string {
	return sq.monthsNarrow[1:]
}

// MonthWide returns the locales wide month given the 'month' provided
func (sq *sq_MK) MonthWide(month time.Month) string {
	return sq.monthsWide[month]
}

// MonthsWide returns the locales wide months
func (sq *sq_MK) MonthsWide() []string {
	return sq.monthsWide[1:]
}

// WeekdayAbbreviated returns the locales abbreviated weekday given the 'weekday' provided
func (sq *sq_MK) WeekdayAbbreviated(weekday time.Weekday) string {
	return sq.daysAbbreviated[weekday]
}

// WeekdaysAbbreviated returns the locales abbreviated weekdays
func (sq *sq_MK) WeekdaysAbbreviated() []string {
	return sq.daysAbbreviated
}

// WeekdayNarrow returns the locales narrow weekday given the 'weekday' provided
func (sq *sq_MK) WeekdayNarrow(weekday time.Weekday) string {
	return sq.daysNarrow[weekday]
}

// WeekdaysNarrow returns the locales narrow weekdays
func (sq *sq_MK) WeekdaysNarrow() []string {
	return sq.daysNarrow
}

// WeekdayShort returns the locales short weekday given the 'weekday' provided
func (sq *sq_MK) WeekdayShort(weekday time.Weekday) string {
	return sq.daysShort[weekday]
}

// WeekdaysShort returns the locales short weekdays
func (sq *sq_MK) WeekdaysShort() []string {
	return sq.daysShort
}

// WeekdayWide returns the locales wide weekday given the 'weekday' provided
func (sq *sq_MK) WeekdayWide(weekday time.Weekday) string {
	return sq.daysWide[weekday]
}

// WeekdaysWide returns the locales wide weekdays
func (sq *sq_MK) WeekdaysWide() []string {
	return sq.daysWide
}

// FmtNumber returns 'num' with digits/precision of 'v' for 'sq_MK' and handles both Whole and Real numbers based on 'v'
func (sq *sq_MK) FmtNumber(num float64, v uint64) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 2 + 2*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, sq.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				for j := len(sq.group) - 1; j >= 0; j-- {
					b = append(b, sq.group[j])
				}
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, sq.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

// FmtPercent returns 'num' with digits/precision of 'v' for 'sq_MK' and handles both Whole and Real numbers based on 'v'
// NOTE: 'num' passed into FmtPercent is assumed to be in percent already
func (sq *sq_MK) FmtPercent(num float64, v uint64) string {
	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	l := len(s) + 3
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, sq.decimal[0])
			continue
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, sq.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	b = append(b, sq.percent...)

	return string(b)
}

// FmtCurrency returns the currency representation of 'num' with digits/precision of 'v' for 'sq_MK'
func (sq *sq_MK) FmtCurrency(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := sq.currencies[currency]
	l := len(s) + len(symbol) + 4 + 2*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, sq.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				for j := len(sq.group) - 1; j >= 0; j-- {
					b = append(b, sq.group[j])
				}
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {
		b = append(b, sq.minus[0])
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, sq.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	b = append(b, sq.currencyPositiveSuffix...)

	b = append(b, symbol...)

	return string(b)
}

// FmtAccounting returns the currency representation of 'num' with digits/precision of 'v' for 'sq_MK'
// in accounting notation.
func (sq *sq_MK) FmtAccounting(num float64, v uint64, currency currency.Type) string {

	s := strconv.FormatFloat(math.Abs(num), 'f', int(v), 64)
	symbol := sq.currencies[currency]
	l := len(s) + len(symbol) + 6 + 2*len(s[:len(s)-int(v)-1])/3
	count := 0
	inWhole := v == 0
	b := make([]byte, 0, l)

	for i := len(s) - 1; i >= 0; i-- {

		if s[i] == '.' {
			b = append(b, sq.decimal[0])
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				for j := len(sq.group) - 1; j >= 0; j-- {
					b = append(b, sq.group[j])
				}
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if num < 0 {

		b = append(b, sq.currencyNegativePrefix[0])

	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	if int(v) < 2 {

		if v == 0 {
			b = append(b, sq.decimal...)
		}

		for i := 0; i < 2-int(v); i++ {
			b = append(b, '0')
		}
	}

	if num < 0 {
		b = append(b, sq.currencyNegativeSuffix...)
		b = append(b, symbol...)
	} else {

		b = append(b, sq.currencyPositiveSuffix...)
		b = append(b, symbol...)
	}

	return string(b)
}

// FmtDateShort returns the short date representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtDateShort(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x2e}...)
	b = strconv.AppendInt(b, int64(t.Month()), 10)
	b = append(b, []byte{0x2e}...)

	if t.Year() > 9 {
		b = append(b, strconv.Itoa(t.Year())[2:]...)
	} else {
		b = append(b, strconv.Itoa(t.Year())[1:]...)
	}

	return string(b)
}

// FmtDateMedium returns the medium date representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtDateMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, sq.monthsAbbreviated[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateLong returns the long date representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtDateLong(t time.Time) string {

	b := make([]byte, 0, 32)

	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, sq.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtDateFull returns the full date representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtDateFull(t time.Time) string {

	b := make([]byte, 0, 32)

	b = append(b, sq.daysWide[t.Weekday()]...)
	b = append(b, []byte{0x2c, 0x20}...)
	b = strconv.AppendInt(b, int64(t.Day()), 10)
	b = append(b, []byte{0x20}...)
	b = append(b, sq.monthsWide[t.Month()]...)
	b = append(b, []byte{0x20}...)

	if t.Year() > 0 {
		b = strconv.AppendInt(b, int64(t.Year()), 10)
	} else {
		b = strconv.AppendInt(b, int64(t.Year()*-1), 10)
	}

	return string(b)
}

// FmtTimeShort returns the short time representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtTimeShort(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)

	return string(b)
}

// FmtTimeMedium returns the medium time representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtTimeMedium(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)

	return string(b)
}

// FmtTimeLong returns the long time representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtTimeLong(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()
	b = append(b, tz...)

	return string(b)
}

// FmtTimeFull returns the full time representation of 't' for 'sq_MK'
func (sq *sq_MK) FmtTimeFull(t time.Time) string {

	b := make([]byte, 0, 32)

	if t.Hour() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Hour()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Minute() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Minute()), 10)
	b = append(b, sq.timeSeparator...)

	if t.Second() < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(t.Second()), 10)
	b = append(b, []byte{0x20}...)

	tz, _ := t.Zone()

	if btz, ok := sq.timezones[tz]; ok {
		b = append(b, btz...)
	} else {
		b = append(b, tz...)
	}

	return string(b)
}
