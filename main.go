package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	scheme   = flag.String("scheme", "http", "Sensu Enterprise URL scheme (http or https)")
	host     = flag.String("host", "localhost", "Sensu Enterprise API host.")
	port     = flag.Int("port", 4567, "Sensu Enterprise API port.")
	user     = flag.String("user", "", "Sensu Enterprise API user.")
	pass     = flag.String("password", "", "Sensu Enterprise API password.")
	timeout  = flag.Int("timeout", 15, "Sensu Enterprise API request timeout (in seconds).")
	latest   = flag.Bool("latest", false, "Only return the latest point per Enterprise metric.")
	caCert   = flag.String("ca-cert", "", "Path to CA certificate")
	insecure = flag.Bool("insecure-skip-verify", false, "Don't verify TLS hostnames")
)

type Metric struct {
	Metric string
	Points [][]int `json:"points"`
}

func main() {
	flag.Parse()

	output, err := getMetrics()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println(output)
}

// tlsOptionsSet only needs to return true if the default TLS settings need
// adjusting. This is usually not required if the OS' cert pool is set up
// correctly.
func tlsOptionsSet() bool {
	return *insecure || len(*caCert) > 0
}

func getRootCerts() (*x509.CertPool, error) {
	certs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if len(*caCert) > 0 {
		b, err := ioutil.ReadFile(*caCert)
		if err != nil {
			return nil, fmt.Errorf("error reading ca-cert: %s", err)
		}
		cert, err := x509.ParseCertificate(b)
		if err != nil {
			return nil, fmt.Errorf("error parsing ca-cert: %s", err)
		}
		certs.AddCert(cert)
	}
	return certs, nil
}

func newClient() (*http.Client, error) {
	var client http.Client
	client.Timeout = time.Duration(*timeout) * time.Second
	var transport http.Transport
	if tlsOptionsSet() {
		rootCerts, err := getRootCerts()
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: *insecure,
			RootCAs:            rootCerts,
		}
	}
	client.Transport = &transport

	return &client, nil
}

func apiRequest(client *http.Client, url string, user string, pass string, timeout int) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, pass)

	res, err := client.Do(req)
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

func getMetrics() (string, error) {
	metrics := []string{"check_requests", "clients", "events", "keepalives", "results"}
	output := ""

	client, err := newClient()
	if err != nil {
		return "", fmt.Errorf("couldn't create HTTP client: %s", err)
	}

	for _, m := range metrics {
		url := fmt.Sprintf("%s://%s:%d/metrics/check_requests", *scheme, *host, *port)
		body, err := apiRequest(client, url, *user, *pass, *timeout)
		if err != nil {
			return "", err
		}

		var metric Metric
		json.Unmarshal(body, &metric)

		if *latest {
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
