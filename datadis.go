package datadis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client  *http.Client
	token   string
	baseURL string
}

type Supply struct {
	Address         string `json:"address"`
	Cups            string `json:"cups"`
	Province        string `json:"province"`
	DistributorCode string `json:"distributorCode"`
	PostalCode      string `json:"postalCode"`
	Municipality    string `json:"municipality"`
	Distributor     string `json:"distributor"`
	PointType       int    `json:"pointType"`
	ValidDateFrom   string `json:"validDateFrom"`
	ValidDateTo     string `json:"validDateTo"`
}

type Measurement struct {
	Cups         string  `json:"cups"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
	Consumption  float32 `json:"consumptionKWh"`
	ObtainMethod string  `json:"obtainMethod"`
}

func NewClient() *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: "https://datadis.es/api-private/api",
	}
}

func (c *Client) Login(username, password string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://datadis.es/nikola-auth/tokens/login?username=%s&password=%s", username, password), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	token, err := ioutil.ReadAll(resp.Body)
	c.token = string(token)

	return err
}

func (c *Client) ConsumptionData(supply *Supply, from, to time.Time) ([]Measurement, error) {
	sto := to.Format("2006/01/02")
	sfrom := from.Format("2006/01/02")
	query := fmt.Sprintf("cups=%s&distributorCode=%s&startDate=%s&endDate=%s&measurementType=%d&pointType=%d", supply.Cups, supply.DistributorCode, url.QueryEscape(sfrom), url.QueryEscape(sto), 0, supply.PointType)
	u := fmt.Sprintf("%s/get-consumption-data?%s", c.baseURL, query)
	fmt.Println(u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var measurements []Measurement
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return measurements, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &measurements)

	return measurements, err
}

func (c *Client) Supplies() ([]Supply, error) {
	u := fmt.Sprintf("%s/get-supplies", c.baseURL)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var supplies []Supply
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return supplies, err
	}
	err = json.Unmarshal(body, &supplies)

	return supplies, err
}
