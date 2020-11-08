package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// API struct
type API struct {
	authToken string
	zoneID    string
}

// New return new API base on enviroments CLOUDFLARE_AUTH_KEY, CLOUDFLARE_AUTH_EMAIL
func New() *API {
	return &API{
		os.Getenv("CLOUDFLARE_AUTH_TOKEN"),
		os.Getenv("CLOUDFLARE_ZONE_ID")}
}

// PurgeCache by urls in cloudflare cdn server
func (api *API) PurgeCache(urls ...string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", api.zoneID)
	body := struct {
		Files []string `json:"files"`
	}{
		urls,
	}
	reader, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(reader))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+api.authToken)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		log.Println(string(resqBody))
		return nil
	}
	return errors.New(string(resqBody))
}
