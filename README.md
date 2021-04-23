## Prerequisites 
1. Install Golang https://golang.org/doc/install
2. Place Unzipped project at root of go/src
3. Install dependencies running `go get ./...`
4. Sign up at `https://free.currencyconverterapi.com/free-api-key` to get a free api key 
5. Set API key in `./stylight/etc/config/config.dev.env` at field `SERVICE_CURRENCY_CONVERSION_CONFIG_API_KEY`. Project will throw errors if this value is not set.

### CLI
1. A built application and example json file has been added at `./stylight/cmd/cli`
2. To rebuild application run `go build` from `./stylight/cmd/cli`
3. To run the application run `./cli -file=<path/to/json/file> -target-currency=<Target-Currance>` for example: `./cli -file=conversion.json -target-currency=JPY`

## HTTP Server
1. To run locally cd into `./stylight/cmd/server` and run `go run main.go`
2. You can hit the endpoint with curl request: 
    ```curl --location --request POST 'http://localhost:50051/stylight/conversion' \
       --header 'Content-Type: application/json' \
       --data-raw '{
       	"conversion" : "USD",
       	"data" : [
       		{ "value" : 1.00, "currency" : "EUR" },
       { "value" : 1.00, "currency" : "USD" },
       { "value" : 1.00, "currency" : "JPY" },
       { "value" : 1.00, "currency" : "EUR" }
       		]
       }'```