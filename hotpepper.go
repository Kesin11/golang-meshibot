package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var KEY = ""

// HotPepper クライアントのためのtype
type HotPepper struct {
	key     string
	baseURL string
}

// NewClient return HotPepper with default value
func NewClient(key string) *HotPepper {
	client := new(HotPepper)
	client.key = key
	client.baseURL = "http://webservice.recruit.co.jp/hotpepper/gourmet/v1/"

	return client
}

type responseResults struct {
	Result struct {
		Shops []Shop `json:"shop"`
	} `json:"results"`
}

// Shop 店舗の情報
type Shop struct {
	ID        string
	Name      string
	LogoImage string `json:"logo_image"`
	Catch     string
	URLs      struct {
		PC string
	}
}

// LunchURL ランチメニューのurl。実際には404のことがある
func (s *Shop) LunchURL() string {
	return fmt.Sprintf("https://www.hotpepper.jp/str%v/lunch/", s.ID)
}

func (h *HotPepper) fetch(keyword string) ([]Shop, error) {
	values := url.Values{}
	values.Add("key", h.key)       // APIキー
	values.Add("keyword", keyword) // 検索ワード
	values.Add("lat", "35.659104")
	values.Add("lng", "139.703742") // ヒカリエ中心
	values.Add("range", "3")        // 1000m
	values.Add("lunch", "1")        // true
	values.Add("format", "json")

	res, err := http.Get(h.baseURL + "?" + values.Encode())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var data responseResults
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Result.Shops, nil
}

// func main() {
// 	client := NewClient(KEY)
// 	shops, _ := client.fetch("和食")

// 	for _, shop := range shops {
// 		pretty.Println(shop)
// 		pretty.Println(shop.LunchURL())
// 	}
// }
