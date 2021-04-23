package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//APIError ...
type APIError struct {
	StatusCode int
	Error      string
}

//Do returns response body in bytes with status code or error with status code...
func (dao *Service) Do(ctx context.Context, request *http.Request) ([]byte, int, error) {
	log.Printf("making api call to %s", request.URL.String())

	request.WithContext(ctx)

	start := time.Now()
	response, err := dao.HTTPClient.Do(request)
	log.Printf("received response from api call to %s with duration %d ms", request.URL.String(), time.Now().Sub(start).Milliseconds())
	if err != nil {
		log.Printf("an error occurred while trying to make api request with error %s", err)
		return nil, http.StatusInternalServerError, err
	}
	if response == nil {
		log.Printf("received a nil response from %s", request.URL.String())
		err = fmt.Errorf("an error occurred trying to call http do request")
		return nil, http.StatusInternalServerError, err
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("an error occurred while trying to read response body from url %s with error %s", request.URL.String(), err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return bodyBytes, response.StatusCode, nil
}

//WithTimeout ... adds timeout to context for http requests
func (dao *Service) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Duration(dao.Config.ClientTimeoutInSeconds)*time.Second)
}
