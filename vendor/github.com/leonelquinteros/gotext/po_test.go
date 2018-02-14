package gotext

import (
	"os"
	"path"
	"testing"
)

func TestPo(t *testing.T) {
	// Set PO content
	str := `
msgid   ""
msgstr  ""

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

# Multi-line msgid
msgid "multi"
"line"
"id"
msgstr "id with multiline content"

# Multi-line msgid_plural
msgid "multi"
"line"
"plural"
"id"
msgstr "plural id with multiline content"

#Multi-line string
msgid "Multi-line"
msgstr "Multi "
"line"

msgid "One with var: %s"
msgid_plural "Several with vars: %s"
msgstr[0] "This one is the singular: %s"
msgstr[1] "This one is the plural: %s"
msgstr[2] "And this is the second plural form: %s"

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

msgid "Empty translation"
msgstr ""

msgid "Empty plural form singular"
msgid_plural "Empty plural form"
msgstr[0] "Singular translated"
msgstr[1] ""

msgid "More"
msgstr "More translation"

	`

	// Write PO content to file
	filename := path.Clean(os.TempDir() + string(os.PathSeparator) + "default.po")

	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't create test file: %s", err.Error())
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		t.Fatalf("Can't write to test file: %s", err.Error())
	}

	// Create po object
	po := new(Po)

	// Try to parse a directory
	po.ParseFile(path.Clean(os.TempDir()))

	// Parse file
	po.ParseFile(filename)

	// Test translations
	tr := po.Get("My text")
	if tr != "Translated text" {
		t.Errorf("Expected 'Translated text' but got '%s'", tr)
	}

	v := "Variable"
	tr = po.Get("One with var: %s", v)
	if tr != "This one is the singular: Variable" {
		t.Errorf("Expected 'This one is the singular: Variable' but got '%s'", tr)
	}

	// Test multi-line id
	tr = po.Get("multilineid")
	if tr != "id with multiline content" {
		t.Errorf("Expected 'id with multiline content' but got '%s'", tr)
	}

	// Test multi-line plural id
	tr = po.Get("multilinepluralid")
	if tr != "plural id with multiline content" {
		t.Errorf("Expected 'plural id with multiline content' but got '%s'", tr)
	}

	// Test multi-line
	tr = po.Get("Multi-line")
	if tr != "Multi line" {
		t.Errorf("Expected 'Multi line' but got '%s'", tr)
	}

	// Test plural
	tr = po.GetN("One with var: %s", "Several with vars: %s", 2, v)
	if tr != "This one is the plural: Variable" {
		t.Errorf("Expected 'This one is the plural: Variable' but got '%s'", tr)
	}

	// Test inexistent translations
	tr = po.Get("This is a test")
	if tr != "This is a test" {
		t.Errorf("Expected 'This is a test' but got '%s'", tr)
	}

	tr = po.GetN("This is a test", "This are tests", 100)
	if tr != "This are tests" {
		t.Errorf("Expected 'This are tests' but got '%s'", tr)
	}

	// Test context translations
	v = "Test"
	tr = po.GetC("One with var: %s", "Ctx", v)
	if tr != "This one is the singular in a Ctx context: Test" {
		t.Errorf("Expected 'This one is the singular in a Ctx context: Test' but got '%s'", tr)
	}

	// Test plural
	tr = po.GetNC("One with var: %s", "Several with vars: %s", 17, "Ctx", v)
	if tr != "This one is the plural in a Ctx context: Test" {
		t.Errorf("Expected 'This one is the plural in a Ctx context: Test' but got '%s'", tr)
	}

	// Test default plural vs singular return responses
	tr = po.GetN("Original", "Original plural", 4)
	if tr != "Original plural" {
		t.Errorf("Expected 'Original plural' but got '%s'", tr)
	}
	tr = po.GetN("Original", "Original plural", 1)
	if tr != "Original" {
		t.Errorf("Expected 'Original' but got '%s'", tr)
	}

	// Test empty translation strings
	tr = po.Get("Empty translation")
	if tr != "Empty translation" {
		t.Errorf("Expected 'Empty translation' but got '%s'", tr)
	}

	tr = po.Get("Empty plural form singular")
	if tr != "Singular translated" {
		t.Errorf("Expected 'Singular translated' but got '%s'", tr)
	}

	tr = po.GetN("Empty plural form singular", "Empty plural form", 1)
	if tr != "Singular translated" {
		t.Errorf("Expected 'Singular translated' but got '%s'", tr)
	}

	tr = po.GetN("Empty plural form singular", "Empty plural form", 2)
	if tr != "Empty plural form" {
		t.Errorf("Expected 'Empty plural form' but got '%s'", tr)
	}

	// Test last translation
	tr = po.Get("More")
	if tr != "More translation" {
		t.Errorf("Expected 'More translation' but got '%s'", tr)
	}
}

func TestPlural(t *testing.T) {
	// Set PO content
	str := `
msgid   ""
msgstr  ""
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

msgid "Singular: %s"
msgid_plural "Plural: %s"
msgstr[0] "TR Singular: %s"
msgstr[1] "TR Plural: %s"
msgstr[2] "TR Plural 2: %s"

	
`
	// Create po object
	po := new(Po)
	po.Parse(str)

	v := "Var"
	tr := po.GetN("Singular: %s", "Plural: %s", 2, v)
	if tr != "TR Plural: Var" {
		t.Errorf("Expected 'TR Plural: Var' but got '%s'", tr)
	}

	tr = po.GetN("Singular: %s", "Plural: %s", 1, v)
	if tr != "TR Singular: Var" {
		t.Errorf("Expected 'TR Singular: Var' but got '%s'", tr)
	}
}

func TestPoHeaders(t *testing.T) {
	// Set PO content
	str := `
msgid   ""
msgstr  ""
# Initial comment
# Headers below
"Language: en\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

# Some comment
msgid "Example"
msgstr "Translated example"
	`

	// Create po object
	po := new(Po)

	// Parse
	po.Parse(str)

	// Check headers expected
	if po.Language != "en" {
		t.Errorf("Expected 'Language: en' but got '%s'", po.Language)
	}

	// Check headers expected
	if po.PluralForms != "nplurals=2; plural=(n != 1);" {
		t.Errorf("Expected 'Plural-Forms: nplurals=2; plural=(n != 1);' but got '%s'", po.PluralForms)
	}
}

func TestMissingPoHeadersSupport(t *testing.T) {
	// Set PO content
	str := `
msgid "Example"
msgstr "Translated example"
	`

	// Create po object
	po := new(Po)

	// Parse
	po.Parse(str)

	// Check translation expected
	if po.Get("Example") != "Translated example" {
		t.Errorf("Expected 'Translated example' but got '%s'", po.Get("Example"))
	}
}

func TestPluralFormsSingle(t *testing.T) {
	// Single form
	str := `
msgid   ""
msgstr  ""
"Plural-Forms: nplurals=1; plural=0;"

# Some comment
msgid "Singular"
msgid_plural "Plural"
msgstr[0] "Singular form"
msgstr[1] "Plural form 1"
msgstr[2] "Plural form 2"
msgstr[3] "Plural form 3"
	`

	// Create po object
	po := new(Po)

	// Parse
	po.Parse(str)

	// Check plural form
	n := po.pluralForm(0)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(0), got %d", n)
	}
	n = po.pluralForm(1)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(1), got %d", n)
	}
	n = po.pluralForm(2)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(2), got %d", n)
	}
	n = po.pluralForm(3)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(3), got %d", n)
	}
	n = po.pluralForm(50)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(50), got %d", n)
	}
}

func TestPluralForms2(t *testing.T) {
	// 2 forms
	str := `
msgid   ""
msgstr  ""
"Plural-Forms: nplurals=2; plural=n != 1;"

# Some comment
msgid "Singular"
msgid_plural "Plural"
msgstr[0] "Singular form"
msgstr[1] "Plural form 1"
msgstr[2] "Plural form 2"
msgstr[3] "Plural form 3"
	`

	// Create po object
	po := new(Po)

	// Parse
	po.Parse(str)

	// Check plural form
	n := po.pluralForm(0)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(0), got %d", n)
	}
	n = po.pluralForm(1)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(1), got %d", n)
	}
	n = po.pluralForm(2)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(2), got %d", n)
	}
	n = po.pluralForm(3)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(3), got %d", n)
	}
}

func TestPluralForms3(t *testing.T) {
	// 3 forms
	str := `
msgid   ""
msgstr  ""
"Plural-Forms: nplurals=3; plural=n%10==1 && n%100!=11 ? 0 : n != 0 ? 1 : 2;"

# Some comment
msgid "Singular"
msgid_plural "Plural"
msgstr[0] "Singular form"
msgstr[1] "Plural form 1"
msgstr[2] "Plural form 2"
msgstr[3] "Plural form 3"
	`

	// Create po object
	po := new(Po)

	// Parse
	po.Parse(str)

	// Check plural form
	n := po.pluralForm(0)
	if n != 2 {
		t.Errorf("Expected 2 for pluralForm(0), got %d", n)
	}
	n = po.pluralForm(1)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(1), got %d", n)
	}
	n = po.pluralForm(2)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(2), got %d", n)
	}
	n = po.pluralForm(3)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(3), got %d", n)
	}
	n = po.pluralForm(100)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(100), got %d", n)
	}
	n = po.pluralForm(49)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(3), got %d", n)
	}
}

func TestPluralFormsSpecial(t *testing.T) {
	// 3 forms special
	str := `
msgid   ""
msgstr  ""
"Plural-Forms: nplurals=3;"
"plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;"

# Some comment
msgid "Singular"
msgid_plural "Plural"
msgstr[0] "Singular form"
msgstr[1] "Plural form 1"
msgstr[2] "Plural form 2"
msgstr[3] "Plural form 3"
	`

	// Create po object
	po := new(Po)

	// Parse
	po.Parse(str)

	// Check plural form
	n := po.pluralForm(1)
	if n != 0 {
		t.Errorf("Expected 0 for pluralForm(1), got %d", n)
	}
	n = po.pluralForm(2)
	if n != 1 {
		t.Errorf("Expected 1 for pluralForm(2), got %d", n)
	}
	n = po.pluralForm(4)
	if n != 1 {
		t.Errorf("Expected 4 for pluralForm(4), got %d", n)
	}
	n = po.pluralForm(0)
	if n != 2 {
		t.Errorf("Expected 2 for pluralForm(2), got %d", n)
	}
	n = po.pluralForm(1000)
	if n != 2 {
		t.Errorf("Expected 2 for pluralForm(1000), got %d", n)
	}
}

func TestTranslationObject(t *testing.T) {
	tr := newTranslation()
	str := tr.get()

	if str != "" {
		t.Errorf("Expected '' but got '%s'", str)
	}

	// Set id
	tr.id = "Text"

	// Get again
	str = tr.get()

	if str != "Text" {
		t.Errorf("Expected 'Text' but got '%s'", str)
	}
}

func TestPoRace(t *testing.T) {
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

	// Create Po object
	po := new(Po)

	// Create sync channels
	pc := make(chan bool)
	rc := make(chan bool)

	// Parse po content in a goroutine
	go func(po *Po, done chan bool) {
		po.Parse(str)
		done <- true
	}(po, pc)

	// Read some translation on a goroutine
	go func(po *Po, done chan bool) {
		po.Get("My text")
		done <- true
	}(po, rc)

	// Read something at top level
	po.Get("My text")

	// Wait for goroutines to finish
	<-pc
	<-rc
}
