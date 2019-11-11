package detail

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cm *DetailManager) getAreaDetail(w http.ResponseWriter, r *http.Request) {
	//client := tablestore.NewClient(cm.TableStoreConfig.Endpoint, cm.TableStoreConfig.InstanceName, cm.TableStoreConfig.AccessKeyId, cm.TableStoreConfig.AccessKeySecret)
	//getRowRequest := &tablestore.GetRowRequest{}
	//criteria := &tablestore.SingleRowQueryCriteria{}
	//
	//putPk := &tablestore.PrimaryKey{}
	//putPk.AddPrimaryKeyColumn("id", surl)
	//criteria.PrimaryKey = putPk
	//
	//getRowRequest.SingleRowQueryCriteria = criteria
	//getRowRequest.SingleRowQueryCriteria.TableName = tableName
	//getRowRequest.SingleRowQueryCriteria.MaxVersion = 1
	//
	//getResp, _ := client.GetRow(getRowRequest)
	//colmap := getResp.GetColumnMap()

	ad := &AreaDetail{
		Adcode:       "330700",
		Name:         "西湖区",
		Date:         "2019-09-26",
		Daypower:     "5",
		Daytemp:      "30",
		Dayweather:   "晴",
		Daywind:      "东",
		Nightpower:   "5",
		Nighttemp:    "19",
		Nightweather: "晴",
		Nightwind:    "东",
		Province:     "浙江",
		Reporttime:   "2019-09-26 22:49:20",
		Week:         "4",
	}

	rs := &AreaDetailResponse{
		Status: "success",
		Errmsg: "",
		Data:   ad,
	}
	clb, _ := json.Marshal(rs)

	fmt.Fprintf(w, "%s", clb)
}
