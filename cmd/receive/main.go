package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"os"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/knative-sample/dingtalk-weather-service/pkg/dingding"
	"github.com/knative-sample/dingtalk-weather-service/pkg/kncloudevents"
	"encoding/json"
	"strings"
)

/*
Example Output:
 CloudEvent: valid ✅
Context Attributes,
  SpecVersion: 0.2
  Type: alicloud.tablestore
  Source: weather
  ID: 9a735ece-fcc8-4432-b053-162761c5bcdf
  Time: 2019-10-09T09:34:03.996828001Z
  ContentType: application/json
Transport Context,
  URI: /
  Host: event-display.default.svc.cluster.local
  Method: POST
Data,
  {
    "adcode": "650103",
    "city": "沙依巴克区",
    "date": "2019-10-04",
    "daypower": "≤3",
    "daytemp": "16",
    "dayweather": "晴",
    "daywind": "无风向",
    "id": "9a735ece-fcc8-4432-b053-162761c5bcdf",
    "nightpower": "≤3",
    "nighttemp": "5",
    "nightweather": "晴",
    "nightwind": "无风向",
    "province": "新疆",
    "reporttime": "2019-10-04 22:45:33",
    "week": "5"
  }
*/
type WeatherInfo struct {
	Adcode string `json:"adcode"`
	City string `json:"city"`
	Date string `json:"date"`
	Daypower string `json:"daypower"`
	Daytemp string `json:"daytemp"`
	Dayweather string `json:"dayweather"`
	Daywind string `json:"daywind"`
	Nightpower string `json:"nightpower"`
	Nighttemp string `json:"nighttemp"`
	Nightweather string `json:"nightweather"`
	Nightwind string `json:"nightwind"`
	Province string `json:"province"`
	Reporttime string `json:"reporttime"`
	Week string `json:"week"`
}
func dispatch(ctx context.Context, event cloudevents.Event) {
	fmt.Printf(event.String())
	payload := &WeatherInfo{}
	if event.Data == nil {
		fmt.Printf("cloudevents.Event\n  Type:%s\n  Data is empty", event.Context.GetType())
	}

	data, ok := event.Data.([]byte)
	if !ok {
		var err error
		data, err = json.Marshal(event.Data)
		if err != nil {
			data = []byte(err.Error())
		}
	}
	json.Unmarshal(data, payload)
	//城市、日期
	if payload.Adcode == adcode && payload.Date == date {
		//天气提醒
		if strings.Contains(payload.Dayweather, dayweather) {
			dingding.SendDingDingReqest(url, http.MethodPost, dingding.BuildTextContext("城市: "+payload.City+" >> 日期：" + payload.Date + " >> 有: "+ payload.Dayweather + ", 出门请注意带伞"))
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "ok")
}
var (
	url string
	adcode string
	date string
	dayweather string
)

func init() {
	flag.StringVar(&url, "dingtalkurl", "", "dingtalk url.")
	flag.StringVar(&url, "adcode", "", "adcode.")
	flag.StringVar(&url, "date", "", "date.")
	flag.StringVar(&url, "dayweather", "", "dayweather.")
}
func main() {
	flag.Parse()

	log.Println(adcode)
	log.Println(date)
	log.Println(dayweather)
	go func() {
		http.HandleFunc("/health", handler)
		port := os.Getenv("PORT")
		if port == "" {
			port = "8022"
		}
		http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}()

	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), dispatch))
}

