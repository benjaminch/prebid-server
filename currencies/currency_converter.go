package currencies

import (
	"encoding/json"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

// CurrencyConverter allows to converts a currency to another one.
// It periodically fetch the syncSourceURL to update currencies rates.
type CurrencyConverter struct {
	httpClient    HttpClient
	SyncSourceURL string
	RefreshDelay  time.Duration
	LastFetched   *time.Time
	Conversions   *Conversion
	tasks         chan bool
}

func NewCurrencyConverter(syncSourceURL string, refreshDelay time.Duration, httpClient HttpClient) (*CurrencyConverter, error) {
	cc := &CurrencyConverter{
		httpClient:    httpClient,
		SyncSourceURL: syncSourceURL,
		RefreshDelay:  refreshDelay,
	}

	cc.tasks = make(chan bool)

	return cc, nil
}

func (cc *CurrencyConverter) Fetch() error {
	response, err := cc.httpClient.Get(cc.SyncSourceURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var updatedConversions *Conversion
	err = json.NewDecoder(response.Body).Decode(&updatedConversions)
	if err != nil {
		return err
	}
	updatedConversions.DataAsOf = TryParseDate(updatedConversions.DataAsOfRaw)

	cc.Conversions = updatedConversions
	fetchedTime := time.Now()
	cc.LastFetched = &fetchedTime

	return nil
}

func (cc *CurrencyConverter) Start() error {
	// Run a first fetch

	// Start the scheduling
	go func() {
		for {
			cc.Fetch()
			select {
			case <-time.After(cc.RefreshDelay):
			case <-cc.tasks:
				return
			}
		}
	}()

	return nil
}

func (cc *CurrencyConverter) Stop() error {
	if _, ok := (<-cc.tasks); ok == false {
		cc.tasks <- true
		close(cc.tasks)
	}

	return nil
}

func (cc *CurrencyConverter) Convert(fromCur string, toCur string) (float32, error) {
	return 0, nil
}

func TryParseDate(rawDate string) *time.Time {
	if rawDate != "" {
		layout := "2006-01-02"
		t, err := time.Parse(layout, rawDate)
		if err == nil {
			return &t
		}
	}
	return nil
}
