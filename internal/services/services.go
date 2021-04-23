package services

import (
	"context"
	"fmt"
	"log"
	currency_converter "stylight/internal/services/currency-converter"
	conversionModels "stylight/internal/services/currency-converter/models"
	httpService "stylight/internal/services/http"
	"sync"
)

var service StylightService

//StylightService ...
type StylightService struct {
	httpdao               httpService.DAO
	currencyConversion currency_converter.DAO
	config                *ServiceConfig
}

//ServiceConfig ...
type ServiceConfig struct {
	HTTPConfig                *httpService.Config               `required:"true" split_words:"true"`
	CurrencyConversionConfig *currency_converter.Config `required:"true" split_words:"true"`
}

//Initialize ...
func Initialize(c *ServiceConfig) (err error) {
	log.Println("initializing service")

	service = StylightService{
		config: c,
	}

	service.httpdao = httpService.NewHTTPDAO(c.HTTPConfig)

	service.currencyConversion = currency_converter.NewCurrencyConverterDAO(c.CurrencyConversionConfig, service.httpdao)
	return nil
}

//using a cache (would probably use redis or add some type of expiration on this)
var cachedConversions =  make(map[string]float64, 0)

type KeyedMutex struct {
	mutexes sync.Map // Zero value is empty and ready for use
}

func (m *KeyedMutex) Lock(key string) func() {
	value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}

var km = KeyedMutex{}

func GetCurrencyRate(ctx context.Context, convertTo string, valuesToBeConverted []conversionModels.ConvertItems) (convertedValues []conversionModels.ConvertItems){
	for _, valueToBeConverted := range valuesToBeConverted {
		conversionID := fmt.Sprintf("%s_%s", valueToBeConverted.Currency, convertTo)
		//Assumption that its okay to cache the conversion rate per job
		unlock := km.Lock(conversionID)
		if cachedConversion, ok := cachedConversions[conversionID]; ok {
			unlock()
			convertedValues = append(convertedValues, conversionModels.ConvertItems{
				Currency:convertTo,
				Value:cachedConversion*valueToBeConverted.Value,
			})
		} else {
			if cachedConversion, ok := cachedConversions[conversionID]; ok {
				convertedValues = append(convertedValues, conversionModels.ConvertItems{
					Currency:convertTo,
					Value:cachedConversion*valueToBeConverted.Value,
				})
				continue
			}
			rate,  err := service.currencyConversion.GetCurrencyRate(ctx, conversionID)
			unlock()
			if err != nil {
				log.Printf("an error occurred while trying to get conversion %s", conversionID)
				//going to ignore errors and continue (can return an error if thats what the user wants)
				continue
			}
			cachedConversions[conversionID] = rate
			convertedValues = append(convertedValues, conversionModels.ConvertItems{
				Currency:convertTo,
				Value:rate*valueToBeConverted.Value,
			})
		}
	}
	return convertedValues
}
