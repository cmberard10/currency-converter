package currency_converter

import (
	httpService "stylight/internal/services/http"
)

//Service ...
type Service struct {
	httpDAO httpService.DAO
	Config     *Config
}

//Config holds the db connection information
type Config struct {
	HostName                 string `required:"true" split_words:"true" `
	APIKey string `required:"true" split_words:"true" `
}

//NewCurrencyConverterDAO ....
func NewCurrencyConverterDAO(config *Config, httpDAO httpService.DAO) *Service {
	s := &Service{
		Config: config,
		httpDAO: httpDAO,
	}

	return s
}
