package city

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cm *CityManager) cityList(w http.ResponseWriter, r *http.Request) {
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

	cl := []*City{{
		Citycode: "010",
		Name:     "北京市",
	},
		{
			Citycode: "0571",
			Name:     "杭州市",
		},
		{
			Citycode: "0512",
			Name:     "苏州市",
		},
	}

	rs := &ListResponse{
		Status: "success",
		Errmsg: "",
		Data:   cl,
	}
	clb, _ := json.Marshal(rs)

	fmt.Fprintf(w, "%s", clb)
}
