package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"errors"

	_ "github.com/lib/pq"
)

//////////////////////////////////////////////////////////////
/////////////////// IMDB CREDENTIALS /////////////////////////
//////////////////////////////////////////////////////////////

const (
	baseURL  = "http://www.omdbapi.com/?"
	plot     = "full"
	tomatoes = "true"

	MovieSearch   = "movie"
	SeriesSearch  = "series"
	EpisodeSearch = "episode"
)

////////////////////////////////////////////////////////////////////
/////////////////// IMDB STRUCTURE & FUNCTIONS /////////////////////
////////////////////////////////////////////////////////////////////

type OmdbApi struct {
	apiKey string
}

func Init(apiKey string) *OmdbApi {
	return &OmdbApi{apiKey: apiKey}
}

type QueryData struct {
	Title      string
	Year       string
	ImdbId     string
	SearchType string
}

type SearchResponse struct {
	Search       []MovieResult
	Response     string
	Error        string
	totalResults int
}

type MovieResult struct {
	Title             string
	Year              string
	Rated             string
	Genre             string
	ImdbID			  string
	Response          string
	Error             string
}

//MovieByTitle returns a MovieResult given Title
func (api *OmdbApi) MovieByTitle(query *QueryData) (*MovieResult, error) {
	resp, err := api.requestAPI("title", query.Title, query.Year, query.SearchType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := new(MovieResult)
	err = json.NewDecoder(resp.Body).Decode(r)

	if err != nil {
		return nil, err
	}
	if r.Response == "False" {
		return r, errors.New(r.Error)
	}
	return r, nil
}

func (api *OmdbApi) requestAPI(apiCategory string, params ...string) (resp *http.Response, err error) {
	var URL *url.URL
	URL, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if len(params) > 1 && params[2] != "" {
		if params[2] != MovieSearch &&
			params[2] != SeriesSearch &&
			params[2] != EpisodeSearch {
			return nil, errors.New("Invalid search category- " + params[2])
		}
	}
	URL.Path += "/"
	parameters := url.Values{}
	parameters.Add("apikey", api.apiKey)

	switch apiCategory {
	case "search":
		parameters.Add("s", params[0])
		parameters.Add("y", params[1])
		parameters.Add("type", params[2])
	case "title":
		parameters.Add("t", params[0])
		parameters.Add("y", params[1])
		parameters.Add("type", params[2])
		parameters.Add("plot", plot)
		parameters.Add("tomatoes", tomatoes)
	case "id":
		parameters.Add("i", params[0])
		parameters.Add("plot", plot)
		parameters.Add("tomatoes", tomatoes)
	}

	URL.RawQuery = parameters.Encode()
	res, err := http.Get(URL.String())
	err = checkErr(res.StatusCode)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func checkErr(status int) error {
	if status != 200 {
		return fmt.Errorf("Status Code %d received from IMDB", status)
	}
	return nil
}

//Stringer Interface for MovieResult
func (mr MovieResult) String() string {
	return fmt.Sprintf("#%s: %s (%s)", mr.ImdbID, mr.Title, mr.Year)
}
