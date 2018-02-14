package gettext

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	spanishMexico      = "es_MX.utf8"
	deutschDeutschland = "de_DE.utf8"
	frenchFrance       = "fr_FR.utf8"
)

// a setUp would be nice
func init() {
	textDomain := "example"
	BindTextdomain(textDomain, "_examples/")
	Textdomain(textDomain)
}

func TestSpanish(t *testing.T) {
	os.Setenv("LANGUAGE", spanishMexico)
	SetLocale(LcAll, "")

	assert.Equal(t, "¡Hola mundo!", Gettext("Hello, world!"))

	assert.Equal(t, "Una manzana", Sprintf(NGettext("An apple", "%d apples", 1), 1, "garbage"))

	assert.Equal(t, "3 manzanas", Sprintf(NGettext("An apple", "%d apples", 3), 3))

	assert.Equal(t, "Buenos días", Gettext("Good morning"))

	assert.Equal(t, "¡Hasta luego!", Gettext("Good bye!"))
}

func TestDeutsch(t *testing.T) {
	os.Setenv("LANGUAGE", deutschDeutschland)
	SetLocale(LcAll, "")

	assert.Equal(t, "Hallo, Welt!", Gettext("Hello, world!"))

	assert.Equal(t, "Ein Apfel", Sprintf(NGettext("An apple", "%d apples", 1), 1, "garbage"))

	assert.Equal(t, "3 Äpfel", Sprintf(NGettext("An apple", "%d apples", 3), 3))

	assert.Equal(t, "Guten morgen", Gettext("Good morning"))

	assert.Equal(t, "Auf Wiedersehen!", Gettext("Good bye!"))
}

func TestFrench(t *testing.T) {
	// Note that we don't have a french translation.
	os.Setenv("LANGUAGE", frenchFrance)
	SetLocale(LcAll, "")

	assert.Equal(t, "Hello, world!", Gettext("Hello, world!"))

	assert.Equal(t, "An apple", Sprintf(NGettext("An apple", "%d apples", 1), 1, "garbage"))

	assert.Equal(t, "3 apples", Sprintf(NGettext("An apple", "%d apples", 3), 3))

	assert.Equal(t, "Good morning", Gettext("Good morning"))

	assert.Equal(t, "Good bye!", Gettext("Good bye!"))
}
