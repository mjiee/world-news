package locale

import (
	"fmt"
	"testing"
)

// TestAppLocale tests the app locale
func TestAppLocale(t *testing.T) {
	SetAppLocalizer("en")

	msg := AppLocalize("test.message")

	fmt.Println(msg)

	msg = AppLocalize("test.template", map[string]string{
		"Name": "John",
	})

	fmt.Println(msg)

	SetAppLocalizer("zh")

	msg = AppLocalize("test.message")

	fmt.Println(msg)

	msg = AppLocalize("test.notfound")

	fmt.Println(msg)
}
