package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getBalance() (s string, err error) {
	resp, err := http.PostForm("https://www.tarsnap.com/manage.cgi",
		url.Values{"address": {tarsnapEmail},
			"password": {tarsnapPassword},
			"action":   {"verboseactivity"},
			"format":   {"csv"},
		})
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	r := csv.NewReader(strings.NewReader(string(body)))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Get the last available balance
	// Format is roughly
	// Balance,2020-05-07,,,,,8.180860481360047296
	// Usage,2020-05-07,device,Daily storage,259688348,0.002094260837459568,
	length := len(records)
	balanceStr := ""
	for i := length - 1; i < length; i-- {
		if records[i][0] == "Balance" {
			balanceStr = records[i][len(records[i])-1]
			break
		}
	}

	return balanceStr, nil
}

func setbalance() {
	balanceStr, err := getBalance()
	if err != nil {
		log.Fatal(err)
	}
	//balanceStr := "8.180860481360047296"
	balance, err := strconv.ParseFloat(balanceStr, 32)
	if err != nil {
		log.Fatal(err)
	}
	//todo figure out how to add account label here
	balanceGauge.WithLabelValues(tarsnapEmail).Set(balance)
}

func recordMetrics() {
	go func() {
		for {
			setbalance()
			time.Sleep(1 * time.Hour)
		}
	}()
}

var (
	tarsnapEmail    = ""
	tarsnapPassword = ""
	opsProcessed    = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	balanceGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "account_balance",
		Namespace: "tarsnap",
		Help:      "",
	}, []string{"account"})
)

func main() {
	tarsnapEmail = os.Getenv("TARSNAP_EMAIL")
	tarsnapPassword = os.Getenv("TARSNAP_PASSWORD")
	if tarsnapEmail == "" || tarsnapPassword == "" {
		log.Fatal("TARSNAP_EMAIL or TARSNAP_PASSWORD env not set")
	}
	// start hourly check to update metrics
	recordMetrics()
	opsProcessed.Inc()
	prometheus.MustRegister(balanceGauge)
	http.Handle("/-/metrics", promhttp.Handler())
	log.Printf("Server listening on http://localhost:9823/-/metrics")
	http.ListenAndServe(":9823", nil)

}
