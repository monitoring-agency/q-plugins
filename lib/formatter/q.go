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

type out struct {
	Stdout   string        `json:"stdout"`
	State    int           `json:"state"`
	Datasets []interface{} `json:"datasets"`
}

func FormatOutputQ(stdout string, st state.State, dataPoints ...interface{}) {
	var dps = make([]interface{}, len(dataPoints))

	for i := 0; i < len(dataPoints); i++ {
		if v, ok := dataPoints[i].(*TimeDataPoint); ok {
			dps = append(dps, Float64DataPoint{
				Key:   v.Key,
				Value: float64(v.Value) / float64(time.Millisecond),
			})
		} else if v, ok := dataPoints[i].(*Float64DataPoint); ok {
			dps = append(dps, v)
		} else if v, ok := dataPoints[i].(*IntDataPoint); ok {
			dps = append(dps, v)
		} else if v, ok := dataPoints[i].(*StringDataPoint); ok {
			dps = append(dps, v)
		}
	}

	j, err := json.Marshal(out{
		Stdout:   stdout,
		State:    int(st),
		Datasets: dps,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(int(state.UNKNOWN))
	}
	fmt.Println(string(j))
	os.Exit(int(st))
}
