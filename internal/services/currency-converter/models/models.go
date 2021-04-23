package models

type CurrencyConversionRequestObj struct {
	Conversion string `json:"conversion"`
	Data []ConvertItems `json:"data"`
}

type ConvertItems struct {
	Value float64 `json:"value"`
	Currency string `json:"currency"`
}
