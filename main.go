package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jroimartin/gocui"
)

// makes a GET request to the provided URL
func (c *Client) MakeGetRequest(url string) ([]byte, error) {
	fmt.Println("Making the get req to " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// set the api key in the request header
	req.Header.Set("X-Api-Key", c.APIKey)
	resp, err := doReq(req, c.httpClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// makes a POST request to the provided URL with the provided request body
func (c *Client) MakePostRequest(url string, reqBody []byte) ([]byte, error) {
	fmt.Println("Making post req to %v, with req body %v", url, reqBody)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// set the api key in the request header
	req.Header.Set("X-Api-Key", c.APIKey)

	resp, err := doReq(req, c.httpClient)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// move this to an env var file later
var APIKEY = "4399fb75326b41bb8422e1731046f157"
var APIBaseURL = "http://192.168.1.150:8989/"

// NewClient returns a client obj with the apikey and baseurl attached
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
	resp, err := c.MakeGetRequest(reqURL)
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
	resp, err := c.MakeGetRequest(reqURL)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp), &returnedSeries)
	return returnedSeries, nil
}
func (c *Client) searchForSeries(seriesName string) ([]series, error) {
	var APISeriesSearchPath = "api/series/lookup?"

	// container for returned series objs
	var returnedSearchResults []series
	reqURL := c.constructAPIRequestURL(APISeriesSearchPath + "term=" + url.QueryEscape(seriesName))

	resp, err := c.MakeGetRequest(reqURL)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp), &returnedSearchResults)
	return returnedSearchResults, nil
}

func (c *Client) getEpisodesBySeriesID(seriesID int) ([]episode, error) {
	var apiEpisodePath = "api/episode?"

	// container for returned episode objs
	var returnedEpisodeResults []episode

	reqURL := c.constructAPIRequestURL(apiEpisodePath + "seriesId=" + strconv.Itoa(seriesID))

	resp, err := c.MakeGetRequest(reqURL)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp), &returnedEpisodeResults)
	return returnedEpisodeResults, nil
}

func (c *Client) constructAPIRequestURL(endpoint string) (apiReqURL string) {
	u, _ := url.ParseRequestURI(c.BaseURL)
	u.Path = endpoint
	return u.String()
}

// layout func for gocui
func layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("shows", -1, -1, 50, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorBlue
		v.SelFgColor = gocui.ColorBlack
		if seriesList != nil {
			for _, currSeries := range seriesList {
				fmt.Fprintln(v, currSeries.Title)
				fmt.Fprintln(v, "Item")
			}
		}
	}
	return nil
}

// keybindings for gocui
func keybindingsCui(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitCui); err != nil {
		return err
	}
	return nil
}

func createLoadingView(g *gocui.Gui) error {
	tw, th := g.Size()
	v, err := g.SetView("LOADING_VIEW", tw/6, (th/2)-1, (tw*5)/6, (th/2)+1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Wrap = true
	fmt.Fprintf(v, "Loading...")
	_, err = g.SetCurrentView("LOADING_VIEW")
	return err
}

// quit func for gocui
func quitCui(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

var seriesList []series

func main() {
	//make http client to execute the request
	client := &http.Client{}

	requestHandler := NewClient(client)

	/* returnedData, _ := requestHandler.getAllSeries()
	fmt.Println(returnedData)

	returnedSeries, _ := requestHandler.getSeriesByID(8)
	fmt.Printf("%v", returnedSeries) */

	seriesList, _ := requestHandler.getAllSeries()
	// print the titles
	for _, currSeries := range seriesList {
		fmt.Println(currSeries.Title)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)
	if err := createLoadingView(g); err != nil {
		panic(err)
	}

	/* _, err = g.SetViewOnTop("LOADING_VIEW")
	if err != nil {
		panic(err)
	} */

	if err := keybindingsCui(g); err != nil {
		panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}
