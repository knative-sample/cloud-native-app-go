package city

import "github.com/knative-sample/cloud-native-app-go/weather/pkg/db"

type CityManager struct {
	Port             string
	TableStoreConfig *db.TableStoreConfig
}

type City struct {
	Citycode string `json:"citycode"`
	Name     string `json:"name"`
}

type ListResponse struct {
	Status string
	Errmsg string
	Data   []*City
}
