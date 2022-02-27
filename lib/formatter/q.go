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
	j, err := json.Marshal(out{
		Stdout:   stdout,
		State:    int(st),
		Datasets: dataPoints,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(int(state.UNKNOWN))
	}
	fmt.Println(string(j))
	os.Exit(int(st))
}
