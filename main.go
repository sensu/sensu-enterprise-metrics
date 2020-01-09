package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Metric struct {
	Metric string
	Points [][]int `json:"points"`
}

func main() {
	host := flag.String("host", "localhost", "Sensu Enterprise API host.")
	port := flag.Int("port", 4567, "Sensu Enterprise API port.")
	user := flag.String("user", "", "Sensu Enterprise API user.")
	pass := flag.String("password", "", "Sensu Enterprise API password.")
	timeout := flag.Int("timeout", 15, "Sensu Enterprise API request timeout (in seconds).")
	latest := flag.Bool("latest", false, "Only return the latest point per Enterprise metric.")
	flag.Parse()

	output, err := getMetrics(*host, *port, *user, *pass, *timeout, *latest)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println(output)
}

func apiRequest(url string, user string, pass string, timeout int) ([]byte, error) {
	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, pass)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected API response status: %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return body, err
}

func getMetrics(host string, port int, user string, pass string, timeout int, latest bool) (string, error) {
	metrics := []string{"check_requests", "clients", "events", "keepalives", "results"}
	output := ""

	for _, m := range metrics {
		url := "http://" + host + ":" + strconv.Itoa(port) + "/metrics/check_requests"
		body, err := apiRequest(url, user, pass, timeout)
		if err != nil {
			return "", err
		}

		var metric Metric
		json.Unmarshal(body, &metric)

		if latest {
			p := metric.Points[len(metric.Points)-1]
			line := fmt.Sprintf("sensu_enterprise_%s %v %v\n", m, p[1], p[0])
			output = output + line
		} else {
			for _, p := range metric.Points {
				line := fmt.Sprintf("sensu_enterprise_%s %v %v\n", m, p[1], p[0])
				output = output + line
			}
		}
	}

	return output, nil
}
