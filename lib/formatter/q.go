package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/myOmikron/q-plugins/lib/state"
	"os"
	"time"
)

type Unit string

const (
	BYTES        Unit = "b"
	MILLISECONDS Unit = "ms"
	NIL          Unit = ""
	DATETIME     Unit = "datetime"
)

type IntDataPoint struct {
	Key       string `json:"key"`
	CurveName string `json:"curveName"`
	Value     int    `json:"value"`
	Unit      Unit   `json:"unit"`
}

type Float64DataPoint struct {
	Key       string  `json:"key"`
	CurveName string  `json:"curveName"`
	Value     float64 `json:"value"`
	Unit      Unit    `json:"unit"`
}

type StringDataPoint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TimeDataPoint struct {
	Key       string        `json:"key"`
	CurveName string        `json:"curveName"`
	Value     time.Duration `json:"value"`
}

type DateTimeDataPoint struct {
	Key   string    `json:"key"`
	Value time.Time `json:"value"`
}

type dataPoint struct {
	Key       string      `json:"key"`
	CurveName string      `json:"curveName"`
	Value     interface{} `json:"value"`
	Type      string      `json:"type"`
	Unit      Unit        `json:"unit"`
}

type out struct {
	Stdout   string      `json:"stdout"`
	Datasets []dataPoint `json:"datasets"`
}

func Error(err error) {
	FormatOutputQ(err.Error(), state.UNKNOWN)
}

func FormatOutputQ(stdout string, st state.State, dataPoints ...interface{}) {
	var dps = make([]dataPoint, 0)

	for i := 0; i < len(dataPoints); i++ {
		if v, ok := dataPoints[i].(*TimeDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:       v.Key,
				CurveName: v.CurveName,
				Value:     v.Value.Milliseconds(),
				Type:      "int",
				Unit:      MILLISECONDS,
			})
		} else if v, ok := dataPoints[i].(*Float64DataPoint); ok {
			dps = append(dps, dataPoint{
				Key:       v.Key,
				Value:     v.Value,
				Type:      "float",
				Unit:      v.Unit,
				CurveName: v.CurveName,
			})
		} else if v, ok := dataPoints[i].(*IntDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:       v.Key,
				Value:     v.Value,
				Type:      "int",
				Unit:      v.Unit,
				CurveName: v.CurveName,
			})
		} else if v, ok := dataPoints[i].(*StringDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: v.Value,
				Type:  "string",
				Unit:  NIL,
			})
		} else if v, ok := dataPoints[i].(*DateTimeDataPoint); ok {
			dps = append(dps, dataPoint{
				Key:   v.Key,
				Value: v.Value,
				Type:  "time",
				Unit:  DATETIME,
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
