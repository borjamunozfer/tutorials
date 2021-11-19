package todos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Post struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserId int    `json:"userId"`
}

type todoClient struct {
	baseURL *url.URL
	client  *http.Client
}

func (t *todoClient) GetPost() (*http.Response, error) {
	response, err := t.client.Get(fmt.Sprintf(t.baseURL.Host+"%s", t.baseURL.Path))
	if err != nil {
		return &http.Response{}, err
	}

	return response, nil
}

func (t *todoClient) GetPosts() (*http.Response, error) {
	response, err := t.client.Get(t.baseURL.Host + t.baseURL.Path)
	if err != nil {
		return &http.Response{}, err
	}
	return response, nil
}

func (t *todoClient) CreatePost(body map[string]interface{}) (*http.Response, error) {

	//we marshal/transform from an interface/object (map[string]interface in our case) to a []byte JSON representation.
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &http.Response{}, err
	}
	//post request with json content type. Body should implements io.Reader interface
	response, err := t.client.Post(t.baseURL.Host+t.baseURL.Path, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return &http.Response{}, err
	}
	return response, nil
}
