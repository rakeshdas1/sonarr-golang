package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//struct for the system status
type sysStatus struct {
	Version   string `json:"version"`
	OsName    string `json:"osName"`
	OsVersion string `json:"osVersion"`
}

//struct for a single series listing
type series struct {
	Title             string   `json:"title"`
	SeasonCount       uint8    `json:"seasonCount"`
	TotalEpisodeCount uint8    `json:"totalEpisodeCount"`
	EpisodeFileCount  uint8    `json:"episodeFileCount"`
	Status            string   `json:"status"`
	Overview          string   `json:"overview"`
	PreviousAiring    string   `json:"previousAiring"`
	Network           string   `json:"network"`
	Seasons           []season `json:"seasons"`
}

//struct for a single season listing
type season struct {
	SeasonNumber uint8  `json:"seasonNumber"`
	Monitored    string `json:"monitored"`
}
type Client struct {
	httpClient *http.Client
	BaseURL    string
	APIKey     string
}

func (c *Client) MakeRequest(url string) ([]byte, error) {
	fmt.Println("Making the get req to " + url)
	req, err := http.NewRequest("GET", url, nil)

	// set the api key in the request header
	req.Header.Set("X-Api-Key", c.APIKey)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(req, c.httpClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}

var APIKEY = "4399fb75326b41bb8422e1731046f157"
var APIBaseURL = "http://192.168.1.150:8989/"

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		APIKey:     APIKEY,
		BaseURL:    APIBaseURL}
}

func doReq(req *http.Request, client *http.Client) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func (c *Client) getAllSeries() ([]series, error) {
	var APISeriesPath = "api/series"
	//container for all the series objs
	var seriesArr []series
	//construct request url
	reqURL := c.constructAPIRequestURL(APISeriesPath)
	resp, err := c.MakeRequest(reqURL)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp), &seriesArr)
	return seriesArr, nil
}
func (c *Client) getSeriesByID(seriesID int) (*series, error) {
	var APISeriesPath = "api/series"

	//container for the returned series obj
	var returnedSeries *series
	reqURL := c.constructAPIRequestURL(APISeriesPath + "/" + strconv.Itoa(seriesID))
	resp, err := c.MakeRequest(reqURL)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp), &returnedSeries)
	return returnedSeries, nil
}

func (c *Client) constructAPIRequestURL(endpoint string) (apiReqURL string) {
	u, _ := url.ParseRequestURI(c.BaseURL)
	u.Path = endpoint
	return u.String()
}

func main() {
	//make http client to execute the request
	client := &http.Client{}

	requestHandler := NewClient(client)

	returnedData, _ := requestHandler.getAllSeries()
	fmt.Println(returnedData)

	returnedSeries, _ := requestHandler.getSeriesByID(8)
	fmt.Printf("%v", returnedSeries)
}
