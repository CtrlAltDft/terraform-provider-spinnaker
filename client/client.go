package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ErrInvalidDecodeResponseParameter invalid parameter for decodeResponse
var ErrInvalidDecodeResponseParameter = errors.New("nil interface provided to decodeResponse")

// Config for Client
type Config struct {
	Address string
	Auth    *Auth
}

func NewConfig() *Config {
	return &Config{
		Auth: NewAuth(),
	}
}

// Client to talk to Spinnaker
type Client struct {
	Config *Config
	client *http.Client
}

// NewClient Return a new client with loaded configuration
func NewClient(config *Config) *Client {

	var httpClient *http.Client
	if config.Auth.Enabled {
		httpClient = newTlsHttpClient(config)
	} else {
		httpClient = http.DefaultClient
	}

	return &Client{
		Config: config,
		client: httpClient,
	}
}

func newTlsHttpClient(config *Config) *http.Client {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("[ERROR]", err)
	}
	certPath := strings.Replace(config.Auth.CertPath, "~", homeDir, 1)
	if certPath == "" {
		log.Fatal("[ERROR] Missing Cert Path")
	}
	keyPath := strings.Replace(config.Auth.KeyPath, "~", homeDir, 1)
	if keyPath == "" {
		log.Fatal("[ERROR] Missing Cert Key Path")
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: config.Auth.Insecure,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}

// NewRequest create http request
func (client *Client) NewRequest(method string, path string) (*http.Request, error) {
	return client.NewRequestWithBody(method, path, nil)
}

// NewRequestWithBody create http request with data as body
func (client *Client) NewRequestWithBody(method string, path string, data interface{}) (*http.Request, error) {
	reqURL, urlErr := url.Parse(client.Config.Address + path)
	if urlErr != nil {
		return nil, urlErr
	}

	jsonValue, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	log.Printf("[INFO] Sending %s %s with body %s\n", method, reqURL, jsonValue)
	req, err := http.NewRequest(method, reqURL.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}

// Do send http request
func (client *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := client.do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	return resp, nil
}

// DoWithRetry send http request with retry
func (client *Client) DoWithRetry(retryOnStatus int, maxAttempts int, createReq func() (*http.Request, error)) (*http.Response, error) {
	attempts := 0
	req, err := createReq()
	if err != nil {
		return nil, err
	}
	resp, respErr := client.Do(req)
	for respErr != nil && attempts < maxAttempts {
		spinnakerError, ok := respErr.(*SpinnakerError)
		if !ok {
			return nil, respErr
		}
		log.Println("[INFO] spinnakerError.Status", spinnakerError.Status)
		if spinnakerError.Status != retryOnStatus {
			return nil, spinnakerError
		}
		time.Sleep(time.Duration(attempts*attempts) * time.Second)

		req, err := createReq()
		if err != nil {
			return nil, err
		}
		log.Printf("[INFO] retry attempt %v for request %v\n", attempts+2, req)
		resp, respErr = client.Do(req)
		attempts++
	}
	return resp, respErr
}

// DoWithResponse send http request and parse response body
func (client *Client) DoWithResponse(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := client.do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	err = decodeResponse(resp, v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// do internal function used by Do and DoWithResponse to validate response
func (client *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := client.client.Do(req)
	if err != nil {
		return resp, err
	}

	err = validateResponse(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func decodeResponse(r *http.Response, v interface{}) error {
	if v == nil {
		return ErrInvalidDecodeResponseParameter
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	log.Println("[DEBUG] Got response body", bodyString)

	return json.Unmarshal([]byte(bodyString), &v)
}

func validateResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	log.Println("[INFO] Error response body", bodyString)

	spinnakerError := SpinnakerError{}
	err := json.Unmarshal([]byte(bodyString), &spinnakerError)
	if err != nil {
		return err
	}

	return &spinnakerError
}
