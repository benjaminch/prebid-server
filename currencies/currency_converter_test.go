package currencies_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/prebid/prebid-server/currencies"
)

func TestFetch_Success(t *testing.T) {

	// Setup:
	responseBody := `{
		"dataAsOf":"2018-09-12",
		"conversions":{
			"USD":{
				"GBP":0.77208
			},
			"GBP":{
				"USD":1.2952
			}
		}
	}`
	currencySrcResponse := &http.Response{
		StatusCode:    200,
		Body:          ioutil.NopCloser(bytes.NewBufferString(responseBody)),
		ContentLength: int64(len(responseBody)),
	}

	mockedHttpClient := &HttpClientMock{
		response: currencySrcResponse,
		err:      nil,
	}

	currencyConverter, _ := currencies.NewCurrencyConverter(
		"currency.fake.com/latest.json",
		time.Duration(1)*time.Minute,
		mockedHttpClient,
	)

	expectedDataAsOf := time.Date(2018, time.September, 12, 0, 0, 0, 0, time.UTC)
	expectedConversions := currencies.Conversion{
		DataAsOf:    &expectedDataAsOf,
		DataAsOfRaw: "2018-09-12",
		Conversions: map[string]map[string]float32{
			"USD": {
				"GBP": 0.77208,
			},
			"GBP": {
				"USD": 1.2952,
			},
		},
	}

	// Execute:
	beforeExecution := time.Now()
	err := currencyConverter.Fetch()

	// Verify:
	if err != nil {
		t.Errorf("err should be nil but was: %s", err)
	}
	if currencyConverter.Conversions == nil {
		t.Errorf("Conversions shouldn't be nil")
	}
	if currencyConverter.LastFetched == nil {
		t.Errorf("LastFetched shouldn't be nil")
	}
	if currencyConverter.LastFetched.After(beforeExecution) == false {
		t.Errorf("LastFetched date should be after Fetch() execution date (%s) but was: %s", beforeExecution, *currencyConverter.LastFetched)
	}
	if currencyConverter.Conversions == nil {
		t.Errorf("Conversions shouldn't be nil")
	}
	if currencyConverter.Conversions.DataAsOf == nil {
		t.Errorf("Conversions.DataAsOf shouldn't be nil")
	}
	if !expectedConversions.DataAsOf.Equal(*currencyConverter.Conversions.DataAsOf) {
		t.Errorf("Conversions.DataAsOf should be %s but was %s", expectedConversions.DataAsOf, *currencyConverter.Conversions.DataAsOf)
	}
	if currencyConverter.Conversions.Conversions == nil {
		t.Errorf("Conversions.Conversions shouldn't be nil")
	}
	if !reflect.DeepEqual(expectedConversions.Conversions, currencyConverter.Conversions.Conversions) {
		t.Errorf("Conversions.Conversions weren't the expected ones")
	}
}

func TestFetch_Fail404(t *testing.T) {
	// TODO
}

func TestStart(t *testing.T) {
	// TODO
}

func TestStop(t *testing.T) {
	// TODO
}

func TestConvert(t *testing.T) {
	// TODO
}

type HttpClientMock struct {
	response *http.Response
	err      error
}

func (m *HttpClientMock) Get(url string) (*http.Response, error) {
	return m.response, m.err
}
