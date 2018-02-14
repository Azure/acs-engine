package main

import (
	"fmt"

	"github.com/gosexy/gettext"
)

func main() {
	gettext.BindTextdomain("example", "./")
	gettext.Textdomain("example")

	gettext.SetLocale(gettext.LcAll, "es_MX.utf8")
	fmt.Println(gettext.Gettext("Hello, world!"))

	gettext.SetLocale(gettext.LcAll, "de_DE.utf8")
	fmt.Println(gettext.Gettext("Hello, world!"))

	gettext.SetLocale(gettext.LcAll, "en_US.utf8")
	fmt.Println(gettext.Gettext("Hello, world!"))
}
