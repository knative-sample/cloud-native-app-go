package detail

import "github.com/knative-sample/cloud-native-app-go/weather/pkg/db"

type DetailManager struct {
	Port             string
	TableStoreConfig *db.TableStoreConfig
}

type AreaDetail struct {
	Adcode       string `json:"citycode"`
	Name         string `json:"name"`
	Date         string `json:"date"`
	Daypower     string `json:"daypower"`
	Daytemp      string `json:"daytemp"`
	Dayweather   string `json:"dayweather"`
	Daywind      string `json:"daywind"`
	Nightpower   string `json:"nightpower"`
	Nighttemp    string `json:"nighttemp"`
	Nightweather string `json:"nightweather"`
	Nightwind    string `json:"nightwind"`
	Province     string `json:"province"`
	Reporttime   string `json:"reporttime"`
	Week         string `json:"week"`
}

type AreaDetailResponse struct {
	Status string
	Errmsg string
	Data   *AreaDetail
}
