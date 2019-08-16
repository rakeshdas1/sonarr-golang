package main

import "net/http"

// struct for the system status
type sysStatus struct {
	Version   string `json:"version"`
	OsName    string `json:"osName"`
	OsVersion string `json:"osVersion"`
}

// struct for a single series listing
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

// struct for a single season listing
type season struct {
	SeasonNumber uint8  `json:"seasonNumber"`
	Monitored    string `json:"monitored"`
}

// struct for the main http client
type Client struct {
	httpClient *http.Client
	BaseURL    string
	APIKey     string
}
