package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Station struct {
	Id                    int64   `json:"id"`
	StationName           string  `json:"stationName"`
	AvailableDocks        int64   `json:"availableDocks"`
	TotalDocks            int64   `json:"totalDocks"`
	Latitude              float64 `json:"latitude"`
	Longitude             float64 `json:"longitude"`
	StatusValue           string  `json:"statusValue"`
	StatusKey             int64   `json:"statusKey"`
	AvailableBikes        int64   `json:"availableBikes"`
	StAddress1            string  `json:"stAddress1"`
	StAddress2            string  `json:"stAddress2"`
	City                  string  `json:"city"`
	PostalCode            string  `json:"postalCode"`
	Location              string  `json:"location"`
	Altitude              string  `json:"altitude"`
	TestStation           bool    `json:"testStation"`
	LastCommunicationTime string  `json:"lastCommunicationTime"`
	LandMark              string  `json:"landMark"`
}
type CoinData struct {
	Code                 string  `json:"code"`
	CandleDateTime       string  `json:"candleDateTime"`    //time.Time
	CandleDateTimeKst    string  `json:"candleDateTimeKst"` //time.Time
	OpeningPrice         float64 `json:"openingPrice"`
	HighPrice            float64 `json:"highPrice"`
	LowPrice             float64 `json:"lowPrice"`
	TradePrice           float64 `json:"tradePrice"`
	CandleAccTradeVolume float64 `json:"candleAccTradeVolume"`
	CandleAccTradePrice  float64 `json:"candleAccTradePrice"`
	Timestamp            uint64  `json:"timestamp"`
	PrevClosingPrice     float64 `json:"prevClosingPrice"`
	Change               string  `json:"change"` //bool
	ChangePrice          float64 `json:"changePrice"`
	ChangeRate           float64 `json:"changeRate"`
	SignedChangePrice    float64 `json:"signedChangePrice"`
	SignedChangeRate     float64 `json:"signedChangeRate"`
}

type StationAPIResponse struct {
	ExecutionTime   string    `json:"executionTime"`
	StationBeanList []Station `json:"stationBeanList"`
}

type coinDataAPIResponse struct {
	CoinDatas []CoinData `json:"CoinDatas"`
}

func getStations(body []byte) (*StationAPIResponse, error) {
	var s = new(StationAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func getCoindatas(body []byte) (*coinDataAPIResponse, error) {
	var s = new(coinDataAPIResponse)

	str := string(body)

	//str = strings.Join([]string{"{\"CoinDatas\":", str, "}"}, "")
	str = "{\"CoinDatas\":" + str + "}"
	fmt.Println(str)

	//jsonVal, err := json.MarshalIndent(body, "aaa", "bbb")

	err := json.Unmarshal([]byte(str), &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	fmt.Println("==== 결과 출력 ====")
	fmt.Println(s)
	return s, err
}

func ttttmain() {

	res, err := http.Get("https://www.citibikenyc.com/stations/json")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(body))
	s, err := getStations([]byte(body))

	fmt.Println(s.StationBeanList[0].City)
	fmt.Println(s.StationBeanList[0].AvailableBikes)
	fmt.Println(s.StationBeanList[0].Location)

	res, err = http.Get("https://crix-api-endpoint.upbit.com/v1/crix/candles/days?code=CRIX.UPBIT.BTC-SBD")
	if err != nil {
		panic(err.Error())
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	c, err := getCoindatas([]byte(body))

	fmt.Println(c.CoinDatas[0])
}
