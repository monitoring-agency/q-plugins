package formatter

import (
	"encoding/json"
	"fmt"
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

func FormatOutputQ(stdout string, state State, dataPoints ...interface{}) {
	j, err := json.Marshal(out{
		Stdout:   stdout,
		State:    int(state) - 1,
		Datasets: dataPoints,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}
	fmt.Println(string(j))
	os.Exit(int(state) - 1)
}
