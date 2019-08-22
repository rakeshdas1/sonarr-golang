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
	ID uint8 `json:"id"`
}

// struct for a single season listing
type season struct {
	SeasonNumber uint8  `json:"seasonNumber"`
	Monitored    string `json:"monitored"`
	Statistics statistics `json:"statistics"`
}

// struct for stats obj found in season listing
type statistics struct {
	PreviousAiring    string `json:"previousAiring"`
	EpisodeFileCount  int       `json:"episodeFileCount"`
	EpisodeCount      int       `json:"episodeCount"`
	TotalEpisodeCount int       `json:"totalEpisodeCount"`
	SizeOnDisk        int64     `json:"sizeOnDisk"`
	PercentOfEpisodes int       `json:"percentOfEpisodes"`
}
// struct for episode listing
type episode struct {
	SeriesID              int    `json:"seriesId"`
	EpisodeFileID         int    `json:"episodeFileId"`
	SeasonNumber          int    `json:"seasonNumber"`
	EpisodeNumber         int    `json:"episodeNumber"`
	Title                 string `json:"title"`
	AirDate               string `json:"airDate"`
	AirDateUtc            string `json:"airDateUtc"`
	Overview              string `json:"overview"`
	HasFile               bool   `json:"hasFile"`
	Monitored             bool   `json:"monitored"`
	SceneEpisodeNumber    int    `json:"sceneEpisodeNumber"`
	SceneSeasonNumber     int    `json:"sceneSeasonNumber"`
	TvDbEpisodeID         int    `json:"tvDbEpisodeId"`
	AbsoluteEpisodeNumber int    `json:"absoluteEpisodeNumber"`
	Downloading           bool   `json:"downloading"`
	ID                    int    `json:"id"`
}

// struct for the main http client
type Client struct {
	httpClient *http.Client
	BaseURL    string
	APIKey     string
}
