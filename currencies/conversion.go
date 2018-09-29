package currencies

import "time"

type Conversion struct {
	DataAsOfRaw string                        `json:"dataAsOf"`
	DataAsOf    *time.Time                    `json:"-"`
	Conversions map[string]map[string]float32 `json:"conversions"`
}
