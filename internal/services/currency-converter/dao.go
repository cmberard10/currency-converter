package currency_converter

import "context"

//DAO ...
type DAO interface {
	GetCurrencyRate(ctx context.Context, conversion string) (float64, error)
}
