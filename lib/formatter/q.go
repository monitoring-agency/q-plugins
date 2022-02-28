package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/myOmikron/q-plugins/lib/state"
	"os"
	"time"
)

type IntDataPoint struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type Float64DataPoint struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type StringDataPoint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TimeDataPoint struct {
	Key   string        `json:"key"`
	Value time.Duration `json:"value"`
}

type DateTimeDataPoint struct {
	Key   string    `json:"key"`
	Value time.Time `json:"value"`
}

type dataPoint struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type out struct {
	Stdout   string        `json:"stdout"`
	Datasets []interface{} `json:"datasets"`
}

func FormatOutputQ(stdout string, st state.State, dataPoints ...interface{}) {
	var dps []interface{}

	for i := 0; i < len(dataPoints); i++ {
		if v, ok := dataPoints[i].(*TimeDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: float64(v.Value) / float64(time.Millisecond),
				Type:  "float",
			})
		} else if v, ok := dataPoints[i].(*Float64DataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: v.Value,
				Type:  "float",
			})
		} else if v, ok := dataPoints[i].(*IntDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: v.Value,
				Type:  "int",
			})
		} else if v, ok := dataPoints[i].(*StringDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: v.Value,
				Type:  "string",
			})
		} else if v, ok := dataPoints[i].(*DateTimeDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: v.Value,
				Type:  "time",
			})
		}
	}

	j, err := json.Marshal(out{
		Stdout:   stdout,
		Datasets: dps,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(int(state.UNKNOWN))
	}
	fmt.Println(string(j))
	os.Exit(int(st))
}
