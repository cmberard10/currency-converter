package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"stylight/internal/server/handlers/validations"
	"stylight/internal/services"
	conversionModels "stylight/internal/services/currency-converter/models"
)



//CurrencyConversionHandler ...
func CurrencyConversionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requestBody := conversionModels.CurrencyConversionRequestObj{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error while unmarshalling currency conversion error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = validations.ValidateCurrencyCode(requestBody.Conversion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, conversion := range requestBody.Data {
		err := validations.ValidateCurrencyCode(conversion.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = validations.ValidateAmount(conversion.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	resp := services.GetCurrencyRate(ctx, requestBody.Conversion, requestBody.Data)

	responseBodyBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error while marshalling currency conversion error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseBodyBytes)
	if err != nil {
		log.Printf("Error while marshalling healthcheck error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
