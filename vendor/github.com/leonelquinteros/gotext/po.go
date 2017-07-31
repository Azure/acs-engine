package gotext

import (
	"bufio"
	"fmt"
	"github.com/mattn/anko/vm"
	"io/ioutil"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"sync"
)

type translation struct {
	id       string
	pluralID string
	trs      map[int]string
}

func newTranslation() *translation {
	tr := new(translation)
	tr.trs = make(map[int]string)

	return tr
}

func (t *translation) get() string {
	// Look for translation index 0
	if _, ok := t.trs[0]; ok {
		return t.trs[0]
	}

	// Return unstranlated id by default
	return t.id
}

func (t *translation) getN(n int) string {
	// Look for translation index
	if _, ok := t.trs[n]; ok {
		return t.trs[n]
	}

	// Return unstranlated plural by default
	return t.pluralID
}

/*
Po parses the content of any PO file and provides all the translation functions needed.
It's the base object used by all package methods.
And it's safe for concurrent use by multiple goroutines by using the sync package for locking.

Example:

    import "github.com/leonelquinteros/gotext"

    func main() {
        // Create po object
        po := new(gotext.Po)

        // Parse .po file
        po.ParseFile("/path/to/po/file/translations.po")

        // Get translation
        println(po.Get("Translate this"))
    }

*/
type Po struct {
	// Headers
	RawHeaders string

	// Language header
	Language string

	// Plural-Forms header
	PluralForms string

	// Parsed Plural-Forms header values
	nplurals int
	plural   string

	// Storage
	translations map[string]*translation
	contexts     map[string]map[string]*translation

	// Sync Mutex
	sync.RWMutex

	// Parsing buffers
	trBuffer  *translation
	ctxBuffer string
}

// ParseFile tries to read the file by its provided path (f) and parse its content as a .po file.
func (po *Po) ParseFile(f string) {
	// Check if file exists
	info, err := os.Stat(f)
	if err != nil {
		return
	}

	// Check that isn't a directory
	if info.IsDir() {
		return
	}

	// Parse file content
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}

	po.Parse(string(data))
}

// Parse loads the translations specified in the provided string (str)
func (po *Po) Parse(str string) {
	// Lock while parsing
	po.Lock()
	defer po.Unlock()

	// Init storage
	if po.translations == nil {
		po.translations = make(map[string]*translation)
		po.contexts = make(map[string]map[string]*translation)
	}

	// Get lines
	lines := strings.Split(str, "\n")

	// Init buffer
	po.trBuffer = newTranslation()
	po.ctxBuffer = ""

	for _, l := range lines {
		// Trim spaces
		l = strings.TrimSpace(l)

		// Skip invalid lines
		if !po.isValidLine(l) {
			continue
		}

		// Buffer context and continue
		if strings.HasPrefix(l, "msgctxt") {
			po.parseContext(l)
			continue
		}

		// Buffer msgid and continue
		if strings.HasPrefix(l, "msgid") && !strings.HasPrefix(l, "msgid_plural") {
			po.parseID(l)
			continue
		}

		// Check for plural form
		if strings.HasPrefix(l, "msgid_plural") {
			po.parsePluralID(l)
			continue
		}

		// Save translation
		if strings.HasPrefix(l, "msgstr") {
			po.parseMessage(l)
			continue
		}

		// Multi line strings and headers
		if strings.HasPrefix(l, "\"") && strings.HasSuffix(l, "\"") {
			po.parseString(l)
			continue
		}
	}

	// Save last translation buffer.
	po.saveBuffer()

	// Parse headers
	po.parseHeaders()
}

// saveBuffer takes the context and translation buffers
// and saves it on the translations collection
func (po *Po) saveBuffer() {
	// If we have something to save...
	if po.trBuffer.id != "" {
		// With no context...
		if po.ctxBuffer == "" {
			po.translations[po.trBuffer.id] = po.trBuffer
		} else {
			// With context...
			if _, ok := po.contexts[po.ctxBuffer]; !ok {
				po.contexts[po.ctxBuffer] = make(map[string]*translation)
			}
			po.contexts[po.ctxBuffer][po.trBuffer.id] = po.trBuffer
		}

		// Flush buffer
		po.trBuffer = newTranslation()
		po.ctxBuffer = ""
	}
}

// parseContext takes a line starting with "msgctxt",
// saves the current translation buffer and creates a new context.
func (po *Po) parseContext(l string) {
	// Save current translation buffer.
	po.saveBuffer()

	// Buffer context
	po.ctxBuffer, _ = strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(l, "msgctxt")))
}

// parseID takes a line starting with "msgid",
// saves the current translation and creates a new msgid buffer.
func (po *Po) parseID(l string) {
	// Save current translation buffer.
	po.saveBuffer()

	// Set id
	po.trBuffer.id, _ = strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(l, "msgid")))
}

// parsePluralID saves the plural id buffer from a line starting with "msgid_plural"
func (po *Po) parsePluralID(l string) {
	po.trBuffer.pluralID, _ = strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(l, "msgid_plural")))
}

// parseMessage takes a line starting with "msgstr" and saves it into the current buffer.
func (po *Po) parseMessage(l string) {
	l = strings.TrimSpace(strings.TrimPrefix(l, "msgstr"))

	// Check for indexed translation forms
	if strings.HasPrefix(l, "[") {
		idx := strings.Index(l, "]")
		if idx == -1 {
			// Skip wrong index formatting
			return
		}

		// Parse index
		i, err := strconv.Atoi(l[1:idx])
		if err != nil {
			// Skip wrong index formatting
			return
		}

		// Parse translation string
		po.trBuffer.trs[i], _ = strconv.Unquote(strings.TrimSpace(l[idx+1:]))

		// Loop
		return
	}

	// Save single translation form under 0 index
	po.trBuffer.trs[0], _ = strconv.Unquote(l)
}

// parseString takes a well formatted string without prefix
// and creates headers or attach multi-line strings when corresponding
func (po *Po) parseString(l string) {
	// Check for multiline from previously set msgid
	if po.trBuffer.id != "" {
		// Append to last translation found
		uq, _ := strconv.Unquote(l)
		po.trBuffer.trs[len(po.trBuffer.trs)-1] += uq

		return
	}

	// Otherwise is a header
	h, err := strconv.Unquote(strings.TrimSpace(l))
	if err != nil {
		return
	}

	po.RawHeaders += h
}

// isValidLine checks for line prefixes to detect valid syntax.
func (po *Po) isValidLine(l string) bool {
	// Skip empty lines
	if l == "" {
		return false
	}

	// Check prefix
	if !strings.HasPrefix(l, "\"") && !strings.HasPrefix(l, "msgctxt") && !strings.HasPrefix(l, "msgid") && !strings.HasPrefix(l, "msgid_plural") && !strings.HasPrefix(l, "msgstr") {
		return false
	}

	return true
}

// parseHeaders retrieves data from previously parsed headers
func (po *Po) parseHeaders() {
	// Make sure we end with 2 carriage returns.
	po.RawHeaders += "\n\n"

	// Read
	reader := bufio.NewReader(strings.NewReader(po.RawHeaders))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return
	}

	// Get/save needed headers
	po.Language = mimeHeader.Get("Language")
	po.PluralForms = mimeHeader.Get("Plural-Forms")

	// Parse Plural-Forms formula
	if po.PluralForms == "" {
		return
	}

	// Split plural form header value
	pfs := strings.Split(po.PluralForms, ";")

	// Parse values
	for _, i := range pfs {
		vs := strings.SplitN(i, "=", 2)
		if len(vs) != 2 {
			continue
		}

		switch strings.TrimSpace(vs[0]) {
		case "nplurals":
			po.nplurals, _ = strconv.Atoi(vs[1])

		case "plural":
			po.plural = vs[1]
		}
	}
}

// pluralForm calculates the plural form index corresponding to n.
// Returns 0 on error
func (po *Po) pluralForm(n int) int {
	po.RLock()
	defer po.RUnlock()

	// Failsafe
	if po.nplurals < 1 {
		return 0
	}
	if po.plural == "" {
		return 0
	}

	// Init compiler
	env := vm.NewEnv()
	env.Define("n", n)

	plural, err := env.Execute(po.plural)
	if err != nil {
		return 0
	}
	if plural.Type().Name() == "bool" {
		if plural.Bool() {
			return 1
		}
		// Else
		return 0
	}

	if int(plural.Int()) > po.nplurals {
		return 0
	}

	return int(plural.Int())
}

// Get retrieves the corresponding translation for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (po *Po) Get(str string, vars ...interface{}) string {
	// Sync read
	po.RLock()
	defer po.RUnlock()

	if po.translations != nil {
		if _, ok := po.translations[str]; ok {
			return fmt.Sprintf(po.translations[str].get(), vars...)
		}
	}

	// Return the same we received by default
	return fmt.Sprintf(str, vars...)
}

// GetN retrieves the (N)th plural form of translation for the given string.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (po *Po) GetN(str, plural string, n int, vars ...interface{}) string {
	// Sync read
	po.RLock()
	defer po.RUnlock()

	if po.translations != nil {
		if _, ok := po.translations[str]; ok {
			return fmt.Sprintf(po.translations[str].getN(po.pluralForm(n)), vars...)
		}
	}

	// Return the plural string we received by default
	return fmt.Sprintf(plural, vars...)
}

// GetC retrieves the corresponding translation for a given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (po *Po) GetC(str, ctx string, vars ...interface{}) string {
	// Sync read
	po.RLock()
	defer po.RUnlock()

	if po.contexts != nil {
		if _, ok := po.contexts[ctx]; ok {
			if po.contexts[ctx] != nil {
				if _, ok := po.contexts[ctx][str]; ok {
					return fmt.Sprintf(po.contexts[ctx][str].get(), vars...)
				}
			}
		}
	}

	// Return the string we received by default
	return fmt.Sprintf(str, vars...)
}

// GetNC retrieves the (N)th plural form of translation for the given string in the given context.
// Supports optional parameters (vars... interface{}) to be inserted on the formatted string using the fmt.Printf syntax.
func (po *Po) GetNC(str, plural string, n int, ctx string, vars ...interface{}) string {
	// Sync read
	po.RLock()
	defer po.RUnlock()

	if po.contexts != nil {
		if _, ok := po.contexts[ctx]; ok {
			if po.contexts[ctx] != nil {
				if _, ok := po.contexts[ctx][str]; ok {
					return fmt.Sprintf(po.contexts[ctx][str].getN(po.pluralForm(n)), vars...)
				}
			}
		}
	}

	// Return the plural string we received by default
	return fmt.Sprintf(plural, vars...)
}
