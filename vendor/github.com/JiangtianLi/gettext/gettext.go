// Copyright (c) 2012-2016 José Carlos Nieto, https://menteslibres.net/xiam
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package gettext provides bindings for GNU Gettext.
package gettext

/*
// #cgo LDFLAGS: -lintl // Use this if: /usr/bin/ld: cannot find -lintl, see https://github.com/gosexy/gettext/issues/1

#include <libintl.h>

#include <locale.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

var (
	// LcAll is for all of the locale.
	LcAll = uint(C.LC_ALL)

	// LcCollate is for regular expression matching (it determines the meaning of
	// range expressions and equivalence classes) and string collation.
	LcCollate = uint(C.LC_COLLATE)

	// LcCtype is for regular expression matching, character classification,
	// conversion, case-sensitive comparison, and wide character functions.
	LcCtype = uint(C.LC_CTYPE)

	// LcMessages is for localizable natural-language messages.
	LcMessages = uint(C.LC_MESSAGES)

	// LcMonetary is for monetary formatting.
	LcMonetary = uint(C.LC_MONETARY)

	// LcNumeric is for number formatting (such as the decimal point and the
	// thousands separator).
	LcNumeric = uint(C.LC_NUMERIC)

	// LcTime is for time and date formatting.
	LcTime = uint(C.LC_TIME)
)

// Deprecated but kept for backwards compatibility.
var (
	LC_ALL      = LcAll
	LC_COLLATE  = LcCollate
	LC_CTYPE    = LcCtype
	LC_MESSAGES = LcMessages
	LC_MONETARY = LcMonetary
	LC_NUMERIC  = LcNumeric
	LC_TIME     = LcTime
)

// SetLocale sets the program's current locale.
func SetLocale(category uint, locale string) string {
	clocale := C.CString(locale)
	defer C.free(unsafe.Pointer(clocale))

	return C.GoString(C.setlocale(C.int(category), clocale))
}

// BindTextdomain sets the directory containing message catalogs.
func BindTextdomain(domainname string, dirname string) string {
	cdirname := C.CString(dirname)
	defer C.free(unsafe.Pointer(cdirname))

	cdomainname := C.CString(domainname)
	defer C.free(unsafe.Pointer(cdomainname))

	return C.GoString(C.bindtextdomain(cdomainname, cdirname))
}

// BindTextdomainCodeset sets the output codeset for message catalogs on the
// given domainname.
func BindTextdomainCodeset(domainname string, codeset string) string {
	cdomainname := C.CString(domainname)
	defer C.free(unsafe.Pointer(cdomainname))

	ccodeset := C.CString(codeset)
	defer C.free(unsafe.Pointer(ccodeset))

	return C.GoString(C.bind_textdomain_codeset(cdomainname, ccodeset))
}

// Textdomain sets or retrieves the current message domain.
func Textdomain(domainname string) string {
	cdomainname := C.CString(domainname)
	defer C.free(unsafe.Pointer(cdomainname))

	return C.GoString(C.textdomain(cdomainname))
}

// Gettext attempts to translate a text string into the user's system language,
// by looking up the translation in a message catalog.
func Gettext(msgid string) string {
	cmsgid := C.CString(msgid)
	defer C.free(unsafe.Pointer(cmsgid))

	return C.GoString(C.gettext(cmsgid))
}

// DGettext is like Gettext(), but looks up the message in the specified
// domain.
func DGettext(domain string, msgid string) string {
	cdomain := cDomainName(domain)
	defer C.free(unsafe.Pointer(cdomain))

	cmsgid := C.CString(msgid)
	defer C.free(unsafe.Pointer(cmsgid))

	return C.GoString(C.dgettext(cdomain, cmsgid))
}

// DCGettext is like Gettext(), but looks up the message in the specified
// domain and category.
func DCGettext(domain string, msgid string, category uint) string {
	cdomain := cDomainName(domain)
	defer C.free(unsafe.Pointer(cdomain))

	cmsgid := C.CString(msgid)
	defer C.free(unsafe.Pointer(cmsgid))

	return C.GoString(C.dcgettext(cdomain, cmsgid, C.int(category)))
}

// NGettext attempts to translate a text string into the user's system
// language, by looking up the appropriate plural form of the translation in a
// message catalog.
func NGettext(msgid string, msgidPlural string, n uint64) string {
	cmsgid := C.CString(msgid)
	defer C.free(unsafe.Pointer(cmsgid))

	cmsgidPlural := C.CString(msgidPlural)
	defer C.free(unsafe.Pointer(cmsgidPlural))

	return C.GoString(C.ngettext(cmsgid, cmsgidPlural, C.ulong(n)))
}

// Sprintf is like fmt.Sprintf() but without %!(EXTRA) errors.
func Sprintf(format string, a ...interface{}) string {
	expects := strings.Count(format, "%") - strings.Count(format, "%%")

	if expects > 0 {
		arguments := make([]interface{}, expects)
		for i := 0; i < expects; i++ {
			if len(a) > i {
				arguments[i] = a[i]
			}
		}
		return fmt.Sprintf(format, arguments...)
	}

	return format
}

// DNGettext is like NGettext(), but looks up the message in the specified
// domain.
func DNGettext(domainname string, msgid string, msgidPlural string, n uint64) string {
	cdomainname := cDomainName(domainname)
	cmsgid := C.CString(msgid)
	cmsgidPlural := C.CString(msgidPlural)

	defer func() {
		C.free(unsafe.Pointer(cdomainname))
		C.free(unsafe.Pointer(cmsgid))
		C.free(unsafe.Pointer(cmsgidPlural))
	}()

	return C.GoString(C.dngettext(cdomainname, cmsgid, cmsgidPlural, C.ulong(n)))
}

// DCNGettext is like NGettext(), but looks up the message in the specified
// domain and category.
func DCNGettext(domainname string, msgid string, msgidPlural string, n uint64, category uint) string {
	cdomainname := cDomainName(domainname)
	cmsgid := C.CString(msgid)
	cmsgidPlural := C.CString(msgidPlural)

	defer func() {
		C.free(unsafe.Pointer(cdomainname))
		C.free(unsafe.Pointer(cmsgid))
		C.free(unsafe.Pointer(cmsgidPlural))
	}()

	return C.GoString(C.dcngettext(cdomainname, cmsgid, cmsgidPlural, C.ulong(n), C.int(category)))
}

// cDomainName returns the domain name CString that can be nil.
func cDomainName(domain string) *C.char {
	if domain == "" {
		return nil
	}
	// The caller is responsible for freeing this up.
	return C.CString(domain)
}
