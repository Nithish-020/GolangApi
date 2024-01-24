package scheduler

import (
	"encoding/json"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

/*
Pupose:This Function is used to Get the Active Ipo's Current marketData from Exchange and insert it in our database table ....
Parameters:

	not Applicable

Response:

	*On Sucess
	=========
	{
		"responseArr": [
			{
				"sino": "1",
				"method": "NSE Fetch-IPO-MktDemand",
				"totalCount": 14,
				"successCount": 9,
				"errCount": 5,
				"status": "S",
				"errMsg": "success"
			},
			{
				"sino": "2",
				"method": "NSE Fetch-IPO-MktCatwise",
				"totalCount": 14,
				"successCount": 9,
				"errCount": 5,
				"status": "S",
				"errMsg": "success"
			}
		],
		"status": "S",
		"errMsg": ""
	}

	!On Error
	========
	{
		"status": "E",
		"reason": "Can't able to get the ipo mktdata details"
	}

Author: Nithish Kumar
Date: 20DEC2023
*/
func FetchIpoMktDataSch(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {
		log.Println("FetchIpoMktDataSch (+)", r.Method)
		//create instance for manualStruct
		var lRespRec FetchRespStruct

		var lPrgrmResp ProgramStruct
		//get header value
		lUser := r.Header.Get("USER")
		// validate the user and check if the user role is admin and then allow for next process
		lRespRec.Status = common.SuccessCode

		// hardcoded the broker id value to fetch Master datas by pavithra
		lConfigFile := common.ReadTomlConfig("./toml/debug.toml")
		lBroker := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["FetchMasterBrokerId"])

		lBrokerId, lErr1 := strconv.Atoi(lBroker)
		if lErr1 != nil {
			log.Println("Error in Convverting string to int", lErr1)
			lPrgrmResp.SINo = "1"
			lPrgrmResp.MethodName = "Error in Converting BrokerId "
			lPrgrmResp.Status = common.ErrorCode
			lPrgrmResp.ErrMsg = "ISFIMDS01" + lErr1.Error()
			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
		} else {

			// Calling the FetchIpomaster method to get the Active Ipo's MarketDemand details From exchange and
			// then store the details  in the database
			lNoToken := "Access Token not found"
			lNSEToken, lCount, lErr2 := exchangecall.NseIpoMktDemand(lUser, lBrokerId)
			if lErr2 != nil {
				log.Println("ISFIMDS02", lErr2)
				lPrgrmResp.SINo = "1"
				lPrgrmResp.MethodName = "NSE Fetch-IPO-MktDemand"
				lPrgrmResp.Status = common.ErrorCode
				lPrgrmResp.ErrMsg = "ISFIMDS02" + lErr2.Error()
				lPrgrmResp.TotalCount = lCount.TotalCount
				lPrgrmResp.SuccessCount = lCount.SuccessCount
				lPrgrmResp.ErrCount = lCount.ErrCount
				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
			} else {
				if lNSEToken == common.ErrorCode {
					lPrgrmResp.SINo = "1"
					lPrgrmResp.MethodName = "NSE Fetch-IPO-MktDemand"
					lPrgrmResp.Status = common.ErrorCode
					lPrgrmResp.ErrMsg = lNoToken
					lPrgrmResp.TotalCount = lCount.TotalCount
					lPrgrmResp.SuccessCount = lCount.SuccessCount
					lPrgrmResp.ErrCount = lCount.ErrCount
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					log.Println("NSE IPO MktDemand Details FETCHING ERROR")
				} else {
					lPrgrmResp.SINo = "1"
					lPrgrmResp.MethodName = "NSE Fetch-IPO-MktDemand"
					lPrgrmResp.Status = common.SuccessCode
					lPrgrmResp.ErrMsg = common.SUCCESS
					lPrgrmResp.TotalCount = lCount.TotalCount
					lPrgrmResp.SuccessCount = lCount.SuccessCount
					lPrgrmResp.ErrCount = lCount.ErrCount
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					log.Println("NSE IPO  MktDemand Details Fetched Successfully")
				}
			}

			lNSEToken, lCount1, lErr3 := exchangecall.NseIpoMktCatwise(lUser, lBrokerId)
			if lErr3 != nil {
				log.Println("ISFIMDS03", lErr3)
				lPrgrmResp.SINo = "2"
				lPrgrmResp.MethodName = "NSE Fetch-IPO-MktCatwise"
				lPrgrmResp.Status = common.ErrorCode
				lPrgrmResp.ErrMsg = "ISFIMDS03" + lErr3.Error()
				lPrgrmResp.TotalCount = lCount1.TotalCount
				lPrgrmResp.SuccessCount = lCount1.SuccessCount
				lPrgrmResp.ErrCount = lCount1.ErrCount
				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
			} else {
				if lNSEToken == common.ErrorCode {
					lPrgrmResp.SINo = "2"
					lPrgrmResp.MethodName = "NSE Fetch-IPO-MktCatwise"
					lPrgrmResp.Status = common.ErrorCode
					lPrgrmResp.ErrMsg = lNoToken
					lPrgrmResp.TotalCount = lCount1.TotalCount
					lPrgrmResp.SuccessCount = lCount1.SuccessCount
					lPrgrmResp.ErrCount = lCount1.ErrCount
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					log.Println("NSE IPO MktCatwise Details FETCHING ERROR")
				} else {
					lPrgrmResp.SINo = "2"
					lPrgrmResp.MethodName = "NSE Fetch-IPO-MktCatwise"
					lPrgrmResp.Status = common.SuccessCode
					lPrgrmResp.ErrMsg = common.SUCCESS
					lPrgrmResp.TotalCount = lCount1.TotalCount
					lPrgrmResp.SuccessCount = lCount1.SuccessCount
					lPrgrmResp.ErrCount = lCount1.ErrCount
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					log.Println("NSE IPO MktCatwise Details Fetched Successfully")
				}
			}
		}
		// Marshal the Response Structure into lData
		lData, lErr4 := json.Marshal(lRespRec)
		if lErr4 != nil {
			log.Println("ISFIMDS04", lErr4)
			fmt.Fprintf(w, "ISFIMDS04"+lErr4.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("FetchIpoMktDataSch (-)")

	}
}
