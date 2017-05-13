package fake

import (
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/leonelquinteros/gotext"
)

// This file is used to generate the test translation file
// 1. go-xgettext -o i18ntestinput.pot --keyword=translator.T --keyword-plural=translator.NT --msgid-bugs-address="" --sort-output test/i18n/i18ntestinput.go
// 2. msginit -l en_US -o i18ntestinput.po -i i18ntestinput.pot
// 3. Modify i18ntestinput.po using poedit as necessary
// Or msgfmt -c -v -o i18ntestinput.mo i18ntestinput.po
// 4. for d in "default" "en_US"; do cp i18ntestinput.mo translations/test/$d/LC_MESSAGES/acsengine.mo; cp i18ntestinput.po translations/test/$d/LC_MESSAGES/acsengine.po; done
// 5. rm i18ntestinput.*

var (
	locale     = gotext.NewLocale("d", "l")
	translator = &i18n.Translator{
		Locale: locale,
	}
	world    = "World"
	resource = "Foo"
)

func aloha() {
	translator.T("Aloha")
}

func foo() {
	translator.T("Hello %s", world)
}

func bar() {
	translator.NT("There is %d parameter in resource %s", "There are %d parameters in resource %s", 9, 9, resource)
}
