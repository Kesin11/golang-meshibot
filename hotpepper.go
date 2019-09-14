package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
)

// TODO: package分けてインターフェースにする
// Restaurant サービスに依存しない汎用的なデータクラス
type Restaurant struct {
	Name        string
	ImageURL    string
	Description string
	URL         string
	LunchURL    string
	Lat         string
	Lng         string
}

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
		Shops []responseShop `json:"shop"`
	} `json:"results"`
}

// Shop 店舗の情報
type responseShop struct {
	ID    string
	Name  string
	Lat   string
	Lng   string
	Photo struct {
		PC struct {
			M string
		}
	}
	Catch string
	URLs  struct {
		PC string
	}
}

func (h *HotPepper) fetch(keyword string) ([]Restaurant, error) {
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

	restaurants := make([]Restaurant, len(data.Result.Shops))
	for i, shop := range data.Result.Shops {
		restaurants[i] = Restaurant{
			Name:        shop.Name,
			ImageURL:    shop.Photo.PC.M,
			Description: shop.Catch,
			URL:         shop.URLs.PC,
			LunchURL:    fmt.Sprintf("https://www.hotpepper.jp/str%v/lunch/", shop.ID),
			Lat:         shop.Lat,
			Lng:         shop.Lng,
		}
	}

	return restaurants, nil
}

func (h *HotPepper) fetchRandom(keyword string, limit int) ([]Restaurant, error) {
	shops, err := h.fetch(keyword)
	// TODO: shuffleする
	i := int(math.Min(float64(cap(shops)), float64(limit))) - 1
	return shops[:i], err
}
