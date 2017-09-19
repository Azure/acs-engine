package gotext

import (
	"os"
	"path"
	"testing"
)

func TestLocale(t *testing.T) {
	// Set PO content
	str := `
msgid ""
msgstr ""
# Initial comment
# Headers below
"Language: en\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

# Some comment
msgid "My text"
msgstr "Translated text"

# More comments
msgid "Another string"
msgstr ""

msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular: %s"
msgstr[1] "This one is the plural: %s"
msgstr[2] "And this is the second plural form: %s"

msgid "This one has invalid syntax translations"
msgid_plural "Plural index"
msgstr[abc] "Wrong index"
msgstr[1 "Forgot to close brackets"
msgstr[0] "Badly formatted string'

msgctxt "Ctx"
msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular in a Ctx context: %s"
msgstr[1] "This one is the plural in a Ctx context: %s"

msgid "Some random"
msgstr "Some random translation"

msgctxt "Ctx"
msgid "Some random in a context"
msgstr "Some random translation in a context"

msgid "More"
msgstr "More translation"

	`

	// Create Locales directory with simplified language code
	dirname := path.Join("/tmp", "en", "LC_MESSAGES")
	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		t.Fatalf("Can't create test directory: %s", err.Error())
	}

	// Write PO content to file
	filename := path.Join(dirname, "my_domain.po")

	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't create test file: %s", err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		t.Fatalf("Can't write to test file: %s", err.Error())
	}

	// Create Locale with full language code
	l := NewLocale("/tmp", "en_US")

	// Force nil domain storage
	l.domains = nil

	// Add domain
	l.AddDomain("my_domain")

	// Test translations
	tr := l.GetD("my_domain", "My text")
	if tr != "Translated text" {
		t.Errorf("Expected 'Translated text' but got '%s'", tr)
	}

	v := "Variable"
	tr = l.GetD("my_domain", "One with var: %s", v)
	if tr != "This one is the singular: Variable" {
		t.Errorf("Expected 'This one is the singular: Variable' but got '%s'", tr)
	}

	// Test plural
	tr = l.GetND("my_domain", "One with var: %s", "Several with vars: %s", 7, v)
	if tr != "This one is the plural: Variable" {
		t.Errorf("Expected 'This one is the plural: Variable' but got '%s'", tr)
	}

	// Test context translations
	v = "Test"
	tr = l.GetDC("my_domain", "One with var: %s", "Ctx", v)
	if tr != "This one is the singular in a Ctx context: Test" {
		t.Errorf("Expected 'This one is the singular in a Ctx context: Test' but got '%s'", tr)
	}

	// Test plural
	tr = l.GetNDC("my_domain", "One with var: %s", "Several with vars: %s", 3, "Ctx", v)
	if tr != "This one is the plural in a Ctx context: Test" {
		t.Errorf("Expected 'This one is the plural in a Ctx context: Test' but got '%s'", tr)
	}

	// Test last translation
	tr = l.GetD("my_domain", "More")
	if tr != "More translation" {
		t.Errorf("Expected 'More translation' but got '%s'", tr)
	}
}

func TestLocaleFails(t *testing.T) {
	// Set PO content
	str := `
msgid ""
msgstr ""
# Initial comment
# Headers below
"Language: en\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

# Some comment
msgid "My text"
msgstr "Translated text"

# More comments
msgid "Another string"
msgstr ""

msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular: %s"
msgstr[1] "This one is the plural: %s"
msgstr[2] "And this is the second plural form: %s"

msgid "This one has invalid syntax translations"
msgid_plural "Plural index"
msgstr[abc] "Wrong index"
msgstr[1 "Forgot to close brackets"
msgstr[0] "Badly formatted string'

msgid "Invalid formatted id[] with no translations

msgctxt "Ctx"
msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular in a Ctx context: %s"
msgstr[1] "This one is the plural in a Ctx context: %s"

msgid "Some random"
msgstr "Some random translation"

msgctxt "Ctx"
msgid "Some random in a context"
msgstr "Some random translation in a context"

msgid "More"
msgstr "More translation"

	`

	// Create Locales directory with simplified language code
	dirname := path.Join("/tmp", "en", "LC_MESSAGES")
	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		t.Fatalf("Can't create test directory: %s", err.Error())
	}

	// Write PO content to file
	filename := path.Join(dirname, "my_domain.po")

	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't create test file: %s", err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		t.Fatalf("Can't write to test file: %s", err.Error())
	}

	// Create Locale with full language code
	l := NewLocale("/tmp", "en_US")

	// Force nil domain storage
	l.domains = nil

	// Add domain
	l.AddDomain("my_domain")

	// Test non-existent "deafult" domain responses
	tr := l.Get("My text")
	if tr != "My text" {
		t.Errorf("Expected 'My text' but got '%s'", tr)
	}

	v := "Variable"
	tr = l.GetN("One with var: %s", "Several with vars: %s", 2, v)
	if tr != "Several with vars: Variable" {
		t.Errorf("Expected 'Several with vars: Variable' but got '%s'", tr)
	}

	// Test inexistent translations
	tr = l.Get("This is a test")
	if tr != "This is a test" {
		t.Errorf("Expected 'This is a test' but got '%s'", tr)
	}

	tr = l.GetN("This is a test", "This are tests", 1)
	if tr != "This are tests" {
		t.Errorf("Expected 'This are tests' but got '%s'", tr)
	}

	// Test syntax error parsed translations
	tr = l.Get("This one has invalid syntax translations")
	if tr != "This one has invalid syntax translations" {
		t.Errorf("Expected 'This one has invalid syntax translations' but got '%s'", tr)
	}

	tr = l.GetN("This one has invalid syntax translations", "This are tests", 1)
	if tr != "This are tests" {
		t.Errorf("Expected 'Plural index' but got '%s'", tr)
	}
}

func TestLocaleRace(t *testing.T) {
	// Set PO content
	str := `# Some comment
msgid "My text"
msgstr "Translated text"

# More comments
msgid "Another string"
msgstr ""

msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular: %s"
msgstr[1] "This one is the plural: %s"
msgstr[2] "And this is the second plural form: %s"

	`

	// Create Locales directory with simplified language code
	dirname := path.Join("/tmp", "es")
	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		t.Fatalf("Can't create test directory: %s", err.Error())
	}

	// Write PO content to file
	filename := path.Join(dirname, "race.po")

	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't create test file: %s", err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		t.Fatalf("Can't write to test file: %s", err.Error())
	}

	// Create Locale with full language code
	l := NewLocale("/tmp", "es")

	// Init sync channels
	ac := make(chan bool)
	rc := make(chan bool)

	// Add domain in goroutine
	go func(l *Locale, done chan bool) {
		l.AddDomain("race")
		done <- true
	}(l, ac)

	// Get translations in goroutine
	go func(l *Locale, done chan bool) {
		l.GetD("race", "My text")
		done <- true
	}(l, rc)

	// Get translations at top level
	l.GetD("race", "My text")

	// Wait for goroutines to finish
	<-ac
	<-rc
}
