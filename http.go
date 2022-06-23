package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

type HttpClient struct {
	client *http.Client
	jar    *Jar
}

var Http *HttpClient = &HttpClient{}

func (client *HttpClient) NewClient() *HttpClient {
	var newClient = &HttpClient{}

	newClient.jar = NewJar(nil)

	if client.jar != nil {
		newClient.jar.SetEntries(client.jar.GetEntries())
	}

	newClient.client = &http.Client{
		Jar: newClient.jar,
	}

	return newClient
}

func (client *HttpClient) ReadBody(res *http.Response) string {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	var contentType = res.Header.Get("Content-Type")

	var regex, _ = regexp.Compile("charset=(.+)")

	var match = regex.FindStringSubmatch(contentType)
	if match == nil {
		return string(body)
	}

	encoding, _ := charset.Lookup(match[1])

	body, _ = encoding.NewDecoder().Bytes(body)

	return string(body)
}

func (client *HttpClient) Request(req *http.Request) *http.Response {
	var httpClient = client.client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("get %s err", err)
		return nil
	}

	return res
}

func (client *HttpClient) Get(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}

	var httpClient = client.client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	log.Printf("get %s", url)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("get %s err %v ", url, err)
		return ""
	}

	log.Printf("get %s code %s", url, res.Status)

	return client.ReadBody(res)
}

func (client *HttpClient) PostForm(url string, formData url.Values) string {
	req, err := http.NewRequest(
		http.MethodPost, url,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return ""
	}

	var httpClient = client.client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	log.Printf("post %s", url)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("post %s err %v ", url, err)
		return ""
	}

	log.Printf("post %s code %s", url, res.Status)

	return client.ReadBody(res)
}
