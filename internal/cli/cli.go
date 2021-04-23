package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/remeh/sizedwaitgroup"
	"log"
	"os"
	"stylight/internal/config"
	"stylight/internal/server/handlers/validations"
	"stylight/internal/services"
	"stylight/internal/services/currency-converter/models"
	"sync"
)

//RunCLI ....
func RunCLI() (err error) {
	webServerConfig, err := config.FromEnv()
	if err != nil {
		return err
	}

	log.Printf("Starting HTTPS server on port %s", webServerConfig.Port)

	err = services.Initialize(webServerConfig.Service)
	if err != nil {
		log.Printf("an error occurred while initializing services: %s", err.Error())
		return err
	}


	filePath := flag.String("file", "", "Full path to json file")
	targetCurrency := flag.String("target-currency", "", "Currency file should be converted to")
	maxGoRoutines := flag.Int("max-routines", 2, "how many go routines can run at once")
	maxBatch := flag.Int("max-chunk", 100, "max amount of records in a batch")
	flag.Parse()

	if *filePath == "" {
		return fmt.Errorf("file param is required. Pass in path as --file=<path-to-file>")
	}
	if *targetCurrency == "" {
		return fmt.Errorf("target-currency param is required. Pass in path as --target-currency=USD")
	}

	err = validations.ValidateCurrencyCode(*targetCurrency)
	if err != nil {
		return
	}

	//byteValue, err := ioutil.ReadFile(*filePath)
	//if err != nil {
	//	log.Printf("an error occurred while reading file with error: %s", err.Error())
	//	return err
	//}
	//dec := json.NewDecoder(strings.NewReader(string(byteValue)))

	f, err := os.Open(*filePath)
	if err != nil {
		log.Printf("an error occurred while reading file with error: %s", err.Error())
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//does not read file into memory, reads line by line
	s := bufio.NewScanner(f)

	var convertedItemsValues []models.ConvertItems
	var batchedCoveredItemValues []models.ConvertItems

	i := 1

	swg := sizedwaitgroup.New(*maxGoRoutines)

	ctx := context.Background()

	mu := &sync.Mutex{}

	for s.Scan() {
		var conversion models.ConvertItems
		err := json.Unmarshal([]byte(s.Text()), &conversion)
		if err != nil {
			log.Printf("an error occurred while unmarshalling with error: %s", err.Error())
			return err
		}
		if i % *maxBatch == 0 && i != 0{
			swg.Add()
			go func(batchedCoveredItemValues []models.ConvertItems) {
				convertedItemsValue := services.GetCurrencyRate(ctx, *targetCurrency, batchedCoveredItemValues)
				mu.Lock()
				convertedItemsValues = append(convertedItemsValues, convertedItemsValue...)
				mu.Unlock()
				swg.Done()
			}(batchedCoveredItemValues)
			batchedCoveredItemValues = []models.ConvertItems{conversion}
		} else {
			batchedCoveredItemValues = append(batchedCoveredItemValues, conversion)
		}
		i++
	}
	if len(batchedCoveredItemValues) > 0 {
		swg.Add()
		go func(batchedCoveredItemValues []models.ConvertItems) {
			defer swg.Done()
			convertedItemsValue := services.GetCurrencyRate(ctx, *targetCurrency, batchedCoveredItemValues)
			mu.Lock()
			convertedItemsValues = append(convertedItemsValues, convertedItemsValue...)
			mu.Unlock()
		}(batchedCoveredItemValues)
	}

	swg.Wait()

	outputBytes, err := json.Marshal(convertedItemsValues)
	if err != nil {
		log.Printf("an error occurred while marshalling output with error: %s", err.Error())
		return err
	}

	log.Println(string(outputBytes))
	fmt.Println(fmt.Sprintf("count of converted Items: %d",len(convertedItemsValues)))
	return nil
}
