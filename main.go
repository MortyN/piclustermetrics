package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var tempnode = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "home_temperature_celsius",
	Help: "The current temperature in degrees Celsius.",
})

var nodeArr = []string{}

// var gaugeSlice = []*prometheus.Gauge{}
var gaugeMap = make(map[string]prometheus.Gauge)

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	tempval, temp := r.URL.Query()["temp"]
	nodeval, node := r.URL.Query()["node"]

	if temp && node {

		tempString := tempval[0][:2] + "." + tempval[0][:len(tempval[0])-2]
		tempFloat, err := strconv.ParseFloat(tempString, 32)

		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"status": "could not parse temp float number"})
		}

		fmt.Printf("Node: %s, Temp: %s", nodeval[0], tempval[0])
		if !containsString(nodeArr, nodeval[0]) {
			nodeArr = append(nodeArr, nodeval[0])
			tempnodes := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "temp_" + nodeval[0],
				Help: "The current temperature in degrees Celsius.",
			})
			gaugeMap[nodeval[0]] = tempnodes
			prometheus.MustRegister(tempnodes)

			tempnodes.Set(tempFloat)

		} else {
			temp_, _ := gaugeMap[nodeval[0]]
			temp_.Set(tempFloat)
		}
		fmt.Println(nodeArr)

		json.NewEncoder(w).Encode(map[string]string{"status": "200", "temp": tempString, "node": nodeval[0]})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "500"})
	}
}

func main() {

	fmt.Println("Started backend on port 8080")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
