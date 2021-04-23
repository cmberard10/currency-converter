package validations

import "fmt"

var validCurrencyCodes = map[string]string{
	"USD" : "USD",
	"EUR" : "EUR",
	"JPY" : "JPY",
}

func ValidateCurrencyCode(code string) error {
	if _, ok := validCurrencyCodes[code]; !ok {
		return fmt.Errorf("%s is not a valid currency code. List includes: USD, EUR, JPY", code)
	}
	return nil
}

func ValidateAmount (value float64) error {
	if value < 0 {
		return fmt.Errorf("value cannot be negative")
	}
	return nil
}
