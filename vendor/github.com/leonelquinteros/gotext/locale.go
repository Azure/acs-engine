package gotext

import (
	"fmt"
	"os"
	"path"
	"sync"
)

/*
Locale wraps the entire i18n collection for a single language (locale)
It's used by the package functions, but it can also be used independently to handle
multiple languages at the same time by working with this object.

Example:

    import (
	    "fmt"
	    "github.com/leonelquinteros/gotext"
    )

    func main() {
        // Create Locale with library path and language code
        l := gotext.NewLocale("/path/to/i18n/dir", "en_US")

        // Load domain '/path/to/i18n/dir/en_US/LC_MESSAGES/default.po'
        l.AddDomain("default")

        // Translate text from default domain
        fmt.Println(l.Get("Translate this"))

        // Load different domain ('/path/to/i18n/dir/en_US/LC_MESSAGES/extras.po')
        l.AddDomain("extras")

        // Translate text from domain
        fmt.Println(l.GetD("extras", "Translate this"))
    }

*/
type Locale struct {
	// Path to locale files.
	path string

	// Language for this Locale
	lang string

	// List of available domains for this locale.
	domains map[string]*Po

	// Sync Mutex
	sync.RWMutex
}

// NewLocale creates and initializes a new Locale object for a given language.
// It receives a path for the i18n files directory (p) and a language code to use (l).
func NewLocale(p, l string) *Locale {
	return &Locale{
		path:    p,
		lang:    l,
		domains: make(map[string]*Po),
	}
}

func (l *Locale) findPO(dom string) string {
	filename := path.Join(l.path, l.lang, "LC_MESSAGES", dom+".po")
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	if len(l.lang) > 2 {
		filename = path.Join(l.path, l.lang[:2], "LC_MESSAGES", dom+".po")
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}

	filename = path.Join(l.path, l.lang, dom+".po")
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	if len(l.lang) > 2 {
		filename = path.Join(l.path, l.lang[:2], dom+".po")
	}

	return filename
}

// AddDomain creates a new domain for a given locale object and initializes the Po object.
// If the domain exists, it gets reloaded.
func (l *Locale) AddDomain(dom string) {
	po := new(Po)

	// Parse file.
	po.ParseFile(l.findPO(dom))

	// Save new domain
	l.Lock()
	defer l.Unlock()

	if l.domains == nil {
		l.domains = make(map[string]*Po)
	}
	l.domains[dom] = po
}

// Get uses a domain "default" to return the corresponding translation of a given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) Get(str string, vars ...interface{}) string {
	return l.GetD("default", str, vars...)
}

// GetN retrieves the (N)th plural form of translation for the given string in the "default" domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetN(str, plural string, n int, vars ...interface{}) string {
	return l.GetND("default", str, plural, n, vars...)
}

// GetD returns the corresponding translation in the given domain for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetD(dom, str string, vars ...interface{}) string {
	return l.GetND(dom, str, str, 1, vars...)
}

// GetND retrieves the (N)th plural form of translation in the given domain for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetND(dom, str, plural string, n int, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.domains != nil {
		if _, ok := l.domains[dom]; ok {
			if l.domains[dom] != nil {
				return l.domains[dom].GetN(str, plural, n, vars...)
			}
		}
	}

	// Return the same we received by default
	return fmt.Sprintf(plural, vars...)
}

// GetC uses a domain "default" to return the corresponding translation of the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetC(str, ctx string, vars ...interface{}) string {
	return l.GetDC("default", str, ctx, vars...)
}

// GetNC retrieves the (N)th plural form of translation for the given string in the given context in the "default" domain.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	return l.GetNDC("default", str, plural, n, ctx, vars...)
}

// GetDC returns the corresponding translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetDC(dom, str, ctx string, vars ...interface{}) string {
	return l.GetNDC(dom, str, str, 1, ctx, vars...)
}

// GetNDC retrieves the (N)th plural form of translation in the given domain for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (l *Locale) GetNDC(dom, str, plural string, n int, ctx string, vars ...interface{}) string {
	// Sync read
	l.RLock()
	defer l.RUnlock()

	if l.domains != nil {
		if _, ok := l.domains[dom]; ok {
			if l.domains[dom] != nil {
				return l.domains[dom].GetNC(str, plural, n, ctx, vars...)
			}
		}
	}

	// Return the same we received by default
	return fmt.Sprintf(plural, vars...)
}
