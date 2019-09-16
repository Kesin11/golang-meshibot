package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HotPepper APIの仕様
// 1: 300m
// 2: 500m
// 3: 1000m (初期値)
// 4: 2000m
// 5: 3000m
const searchRange = "3"

// "宮益坂": "X095"
// "渋谷東口": "XA0U"
// "桜丘": "XA0V"
// "渋谷新南口": "XA9G"
var searchSmallArea = []string{"X095", "XA0U", "XA0V", "XA9G"}

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

// NewHotPepperClient return HotPepper with default value
func NewHotPepperClient(key string) *HotPepper {
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

// Fetch API結果を全て返す
func (h *HotPepper) Fetch(keyword string) ([]Restaurant, error) {
	values := url.Values{}
	values.Add("key", h.key)       // APIキー
	values.Add("keyword", keyword) // 検索ワード
	values.Add("lat", "35.659104")
	values.Add("lng", "139.703742")                              // ヒカリエ中心
	values.Add("range", searchRange)                             // 検索範囲
	values.Add("small_area", strings.Join(searchSmallArea, ",")) // 検索エリア限定
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

// FetchRandom API結果を先頭からlimit件ランダムで返す
func (h *HotPepper) FetchRandom(keyword string, limit int) ([]Restaurant, error) {
	restaurants, err := h.Fetch(keyword)
	if len(restaurants) < 1 {
		return []Restaurant{}, err
	}

	// shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(restaurants), func(i, j int) {
		restaurants[i], restaurants[j] = restaurants[j], restaurants[i]
	})

	// 先頭からlimit分だけをreturn。要素数がlimit以下の場合は、存在する分だけreturn
	i := int(math.Min(float64(cap(restaurants)), float64(limit)))
	return restaurants[:i], err
}
