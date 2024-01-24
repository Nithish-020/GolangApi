// package scheduler

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// type ScheduleStruct struct {
// 	SINo   string `json:"sino"`
// 	Method string `json:"method"`
// 	Status string `json:"status"`
// 	ErrMsg string `json:"errMsg"`
// }

// type ScheduleRespStruct struct {
// 	ResponseArr []ScheduleStruct `json:"responseArr"`
// 	Status      string           `json:"status"`
// 	ErrMsg      string           `json:"errMsg"`
// }

// func OfflineScheduler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("offlineScheduler (+)", r.Method)
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "GET" {

// 		// lUser := r.Header.Get("USER")

// 		// var lRespRec ScheduleStruct

// 		// var err error
// 		var lRespArr ScheduleRespStruct

// 		lRespArr.Status = "S"
// 		// lRespArr.ErrMsg = common.SUCCESS

// 		//This Methods is used for OfflineAppsch
// 		// lOfflinesch := iposchedule.IpoOfflineProcess(lUser)
// 		// iposchedule.IpoOfflineProcess(lUser)
// 		// log.Println("lNseOfflinesch", lOfflinesch)
// 		// if lOfflinesch.Status == "S" {
// 		// 	lRespRec.Status = common.SuccessCode
// 		// 	lRespRec.ErrMsg = common.SUCCESS
// 		// 	lRespRec.SINo = "1"
// 		// 	lRespRec.Method = "NSE- IPO Offline Application"
// 		// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
// 		// } else {
// 		// 	lRespRec.Status = lOfflinesch.Status
// 		// 	lRespRec.ErrMsg = lOfflinesch.ErrMsg
// 		// 	lRespRec.SINo = "1"
// 		// 	lRespRec.Method = "NSE- IPO Offline Application"
// 		// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
// 		// }

// 		// lBseOfflinesch := iposchedule.BseOfflineSch(lUser)
// 		// log.Println("lBseOfflinesch", lBseOfflinesch)
// 		// if lBseOfflinesch.Status == "S" {
// 		// 	lRespRec.Status = common.SuccessCode
// 		// 	lRespRec.ErrMsg = common.SUCCESS
// 		// 	lRespRec.SINo = "2"
// 		// 	lRespRec.Method = "BSE- IPO Offline Application"
// 		// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
// 		// } else {
// 		// 	lRespRec.Status = lBseOfflinesch.Status
// 		// 	lRespRec.ErrMsg = lBseOfflinesch.ErrMsg
// 		// 	lRespRec.SINo = "2"
// 		// 	lRespRec.Method = "BSE- IPO Offline Application"
// 		// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
// 		// }

// 		//This Methods is used for sgbofflinedata
// 		// lBseSgbsch := sgbschedule.SgbOfflineSch(lUser, r)
// 		// if lBseSgbsch.Status == "S" {
// 		// 	log.Println("lBseSgbsch", lBseSgbsch)
// 		// 	lRespRec.Status = common.SuccessCode
// 		// 	lRespRec.ErrMsg = common.SUCCESS
// 		// 	lRespRec.SINo = "3"
// 		// 	lRespRec.Method = "BSE- SGB Offline Application"
// 		// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
// 		// } else {
// 		// 	lRespRec.Status = lBseSgbsch.Status
// 		// 	lRespRec.ErrMsg = lBseSgbsch.ErrMsg
// 		// 	lRespRec.SINo = "3"
// 		// 	lRespRec.Method = "BSE- SGB Offline Application"
// 		// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
// 		// }

// 		//Marshal the lRespArr Value
// 		lData, lErr3 := json.Marshal(lRespArr)
// 		if lErr3 != nil {
// 			log.Println("SOAOS03", lErr3)
// 			fmt.Fprintf(w, "SOAOS03"+lErr3.Error())
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 	}
// 	log.Println("offlineScheduler (-)", r.Method)

// }
package scheduler

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/iposchedule"
	"fcs23pkg/common"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ScheduleStruct struct {
	SINo         string `json:"sino"`
	Method       string `json:"method"`
	TotalCount   int    `json:"totalCount"`
	SuccessCount int    `json:"successCount"`
	ErrCount     int    `json:"errCount"`
	Status       string `json:"status"`
	ErrMsg       string `json:"errMsg"`
}

type ScheduleRespStruct struct {
	ResponseArr []ScheduleStruct `json:"responseArr"`
	Status      string           `json:"status"`
	ErrMsg      string           `json:"errMsg"`
}

func OfflineScheduler(w http.ResponseWriter, r *http.Request) {
	log.Println("offlineScheduler (+)", r.Method)
	// origin := r.Header.Get("Origin")
	// var lBrokerId int
	// var lErr error
	// w.Header().Set("Access-Control-Allow-Origin", origin)
	// lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
	// if lErr != nil {
	// 	log.Println(lErr, origin)
	// }
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		lUser := r.Header.Get("USER")

		var lRespRec ScheduleStruct

		// var err error
		var lRespArr ScheduleRespStruct
		lRespArr.Status = "S"
		// lRespArr.ErrMsg = common.SUCCESS

		lConfigFile := common.ReadTomlConfig("./toml/debug.toml")
		lBroker := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["FetchMasterBrokerId"])

		lBrokerId, lErr1 := strconv.Atoi(lBroker)
		if lErr1 != nil {
			log.Println("Error in Convverting string to int", lErr1)
			lRespRec.SINo = "1"
			lRespRec.Method = "Error in Converting BrokerId "
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "ISFIMS01" + lErr1.Error()
			lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
		} else {
			//This Methods is used for OfflineAppsch
			lNseOfflinesch := iposchedule.NseOfflineSch(lUser, lBrokerId)
			log.Println("lNseOfflinesch", lNseOfflinesch)
			if lNseOfflinesch.Status == "S" {
				lRespRec.TotalCount = lNseOfflinesch.TotalCount
				lRespRec.SuccessCount = lNseOfflinesch.SuccessCount
				lRespRec.ErrCount = lNseOfflinesch.ErrCount
				lRespRec.Status = common.SuccessCode
				lRespRec.ErrMsg = common.SUCCESS
				lRespRec.SINo = "1"
				lRespRec.Method = "NSE- IPO Offline Applications"
				lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
			} else {
				lRespRec.TotalCount = lNseOfflinesch.TotalCount
				lRespRec.SuccessCount = lNseOfflinesch.SuccessCount
				lRespRec.ErrCount = lNseOfflinesch.ErrCount
				lRespRec.Status = lNseOfflinesch.Status
				lRespRec.ErrMsg = lNseOfflinesch.ErrMsg
				lRespRec.SINo = "1"
				lRespRec.Method = "NSE- IPO Offline Applications"
				lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
			}

			lBseOfflinesch := iposchedule.BseOfflineSch(lUser, lBrokerId)
			log.Println("lBseOfflinesch", lBseOfflinesch)
			if lBseOfflinesch.Status == "S" {
				lRespRec.TotalCount = lBseOfflinesch.TotalCount
				lRespRec.SuccessCount = lBseOfflinesch.SuccessCount
				lRespRec.ErrCount = lBseOfflinesch.ErrCount
				lRespRec.Status = common.SuccessCode
				lRespRec.ErrMsg = common.SUCCESS
				lRespRec.SINo = "2"
				lRespRec.Method = "BSE- IPO Offline Applications"
				lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
			} else {
				lRespRec.TotalCount = lBseOfflinesch.TotalCount
				lRespRec.SuccessCount = lBseOfflinesch.SuccessCount
				lRespRec.ErrCount = lBseOfflinesch.ErrCount
				lRespRec.Status = lBseOfflinesch.Status
				lRespRec.ErrMsg = lBseOfflinesch.ErrMsg
				lRespRec.SINo = "2"
				lRespRec.Method = "BSE- IPO Offline Applications"
				lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
			}

			//This Methods is used for sgbofflinedata
			// lBseSgbsch := sgbschedule.SgbOfflineSch(lUser, r)
			// if lBseSgbsch.Status == "S" {
			// 	log.Println("lBseSgbsch", lBseSgbsch)
			// 	lRespRec.Status = common.SuccessCode
			// 	lRespRec.ErrMsg = common.SUCCESS
			// 	lRespRec.SINo = "3"
			// 	lRespRec.Method = "BSE- SGB Offline Application"
			// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
			// } else {
			// 	lRespRec.Status = lBseSgbsch.Status
			// 	lRespRec.ErrMsg = lBseSgbsch.ErrMsg
			// 	lRespRec.SINo = "3"
			// 	lRespRec.Method = "BSE- SGB Offline Application"
			// 	lRespArr.ResponseArr = append(lRespArr.ResponseArr, lRespRec)
			// }
		}
		log.Println("RespArray of OFfline Sch", lRespArr)
		//Marshal the lRespArr Value
		lData, lErr3 := json.Marshal(lRespArr)
		if lErr3 != nil {
			log.Println("SOAOS03", lErr3)
			fmt.Fprintf(w, "SOAOS03"+lErr3.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("offlineScheduler (-)", r.Method)

}
