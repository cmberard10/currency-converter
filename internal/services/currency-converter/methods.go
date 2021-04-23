package currency_converter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const freeCurrencyConverterEndpoint = "api/v7/convert"

func (dao *Service) GetCurrencyRate(ctx context.Context, conversion string) (float64, error) {
	url := fmt.Sprintf("%s/%s?q=%s&compact=y&apiKey=%s", dao.Config.HostName, freeCurrencyConverterEndpoint, conversion, dao.Config.APIKey)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("An error occurred while building http request with error: %s", err.Error())
		return 0, err
	}

	timeout, cancel := dao.httpDAO.WithTimeout(ctx)
	defer cancel()

	bodyBytes, statusCode, err := dao.httpDAO.Do(timeout, request)
	if err != nil {
		log.Printf("An error occurred while making api call to %s with error: %s", request.URL.String(), err.Error())
		return 0, err
	}


	if statusCode < 200 && statusCode > 299 {
		log.Printf("Recieved an unexpected status code %d with body: %s", statusCode, string(bodyBytes))
		return 0, err
	}

	response := map[string]map[string]float64{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Printf("An error occurred while unmarshalling response body with err %s and body %s", err.Error(), string(bodyBytes))
		return 0, err
	}

	if val, ok := response[conversion]; ok {
		if rate, rateOK := val["val"]; rateOK {
			return rate, nil
		}
	}
	err = fmt.Errorf("conversion rate was not found. Expected response object is not correct from body %s", string(bodyBytes))
	log.Printf(err.Error())
	return 0, err
}
