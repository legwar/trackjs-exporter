package main

import (
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
	"encoding/json"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var TOKEN = "{{API_KEY}}"

type TrackJSResponce struct {
	Data []struct {
		Key        time.Time `json:"key"`
		Count      int       `json:"count"`
		UserCount  int       `json:"userCount"`
		TrackJsURL string    `json:"trackJsUrl"`
	} `json:"data"`
	Metadata struct {
		Page       int    `json:"page"`
		Size       int    `json:"size"`
		HasMore    bool   `json:"hasMore"`
		TrackJsURL string `json:"trackJsUrl"`
	} `json:"metadata"`
}

func getAlertsCounts() float64 {
	j := 0
        currentTime := time.Now()
        client := &http.Client{}
        req, err := http.NewRequest("GET", "https://api.trackjs.com/{{CUSTOMER_ID}}/v1/errors/daily", nil)
        if err != nil {
                fmt.Println(err)
        }
        req.Header.Add("Authorization", TOKEN)
        resp, err := client.Do(req)
        temp, err := ioutil.ReadAll(resp.Body)

        var s = TrackJSResponce{}
        json.Unmarshal(temp, &s)
        for _, row := range s.Data {
                if currentTime.Format("2006-01-02") == row.Key.Format("2006-01-02") {
                        j = row.Count
                }
        }
	return float64(j)
}

var (
	TotalAlerts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "TrackJS_total_alerts",
		Help: "TrackJS total alerts",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(TotalAlerts)
}

func main() {

        go func() {
                for {
                        time.Sleep(60 * time.Second)
			TotalAlerts.Set(getAlertsCounts())
                }
        }()
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)

}
