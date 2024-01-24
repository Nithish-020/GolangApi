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

type ProgramStruct struct {
	SINo         string `json:"sino"`
	MethodName   string `json:"method"`
	TotalCount   int    `json:"totalCount"`
	SuccessCount int    `json:"successCount"`
	ErrCount     int    `json:"errCount"`
	Status       string `json:"status"`
	ErrMsg       string `json:"errMsg"`
}

// this response struct for the fetchIpoMAster API
type FetchRespStruct struct {
	ResponseArr []ProgramStruct `json:"responseArr"`
	Status      string          `json:"status"`
	ErrMsg      string          `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Active Ipo Details from Exchange and insert it in our database table ....
Parameters:

	not Applicable

Response:

	*On Sucess
	=========
	{
		"IpoDetails": [
			{
				"id": 18,
				"symbol": "MMIPO26",
				"startDate": "2023-06-02",
				"endDate": "2023-06-30",
				"priceRange": "1000 - 2000",
				"cutOffPrice": 2000,
				"minBidQuantity": 10,
				"applicationStatus": "Pending",
				"upiStatus": "Accepted by Investor"
			},
			{
				"id": 10,
				"symbol": "fixed",
				"startDate": "2023-05-10",
				"endDate": "2023-08-29",
				"priceRange": "755 - 755",
				"cutOffPrice": 755,
				"minBidQuantity": 100,
				"applicationStatus": "-",
				"upiStatus": "-"
			}
		],
		"status": "S",
		"errMsg": ""
	}

	!On Error
	========
	{
		"status": "E",
		"reason": "Can't able to get the ipo details"
	}

Author: Nithish Kumar
Date: 08JUNE2023
*/
func FetchIpoMasterSch(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchIpoMasterSch (+)", r.Method)
	// commented by pavithra
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
			log.Println("ISFIMS01", lErr1)
			lPrgrmResp.SINo = "1"
			lPrgrmResp.MethodName = "Error in Converting BrokerId "
			lPrgrmResp.Status = common.ErrorCode
			lPrgrmResp.ErrMsg = "ISFIMS01" + lErr1.Error()
			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
		} else {

			// Calling the FetchIpomaster method to get the Active Ipo details From exchange and
			// then store the details  in the database
			lNoToken := "Access Token not found"
			lNSEToken, lErr2 := exchangecall.FetchNseIPOmaster(lUser, lBrokerId)
			if lErr2 != nil {
				log.Println("ISFIMS02", lErr2)
				lPrgrmResp.SINo = "1"
				lPrgrmResp.MethodName = "NSE Fetch-IPO-Master"
				lPrgrmResp.Status = common.ErrorCode
				lPrgrmResp.ErrMsg = "ISFIMS02" + lErr2.Error()
				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
			} else {
				if lNSEToken == common.ErrorCode {
					lPrgrmResp.SINo = "1"
					lPrgrmResp.MethodName = "NSE Fetch-IPO-Master"
					lPrgrmResp.Status = common.ErrorCode
					lPrgrmResp.ErrMsg = lNoToken
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("NSE IPO Details FETCHING ERROR")
				} else {
					lPrgrmResp.SINo = "1"
					lPrgrmResp.MethodName = "NSE Fetch-IPO-Master"
					lPrgrmResp.Status = common.SuccessCode
					lPrgrmResp.ErrMsg = common.SUCCESS
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("NSE IPO Details Fetched Successfully")
				}
			}

			lBseToken, lErr3 := exchangecall.FetchBseIPOmaster(lUser, lBrokerId)
			if lErr3 != nil {
				log.Println("ISFIMS03", lErr3)
				lPrgrmResp.SINo = "2"
				lPrgrmResp.MethodName = "BSE Fetch-IPO-Master"
				lPrgrmResp.Status = common.ErrorCode
				lPrgrmResp.ErrMsg = "ISFIMS03" + lErr3.Error()
				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
			} else {
				if lBseToken == common.ErrorCode {
					lPrgrmResp.SINo = "2"
					lPrgrmResp.MethodName = "BSE Fetch-IPO-Master"
					lPrgrmResp.Status = common.ErrorCode
					lPrgrmResp.ErrMsg = lNoToken
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("BSE IPO Details Fetching Error")
				} else {
					lPrgrmResp.SINo = "2"
					lPrgrmResp.MethodName = "BSE Fetch-IPO-Master"
					lPrgrmResp.Status = common.SuccessCode
					lPrgrmResp.ErrMsg = common.SUCCESS
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("BSE IPO Details Fetched Successfully")
				}
			}

			//sgb master fetch in bse exchange
			lSgbNseToken, lErr4 := exchangecall.FetchSgbMasterNSE(lUser, lBrokerId)
			if lErr4 != nil {
				log.Println("ISFIMS04", lErr4)
				lPrgrmResp.SINo = "3"
				lPrgrmResp.MethodName = "NSE Fetch-SGB-Master"
				lPrgrmResp.Status = common.ErrorCode
				lPrgrmResp.ErrMsg = "ISFIMS04" + lErr4.Error()
				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
			} else {
				if lSgbNseToken == common.ErrorCode {
					lPrgrmResp.SINo = "3"
					lPrgrmResp.MethodName = "NSE Fetch-SGB-Master"
					lPrgrmResp.Status = common.ErrorCode
					lPrgrmResp.ErrMsg = lNoToken
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("NSE Fetch SGB Master Error")
				} else {
					lPrgrmResp.SINo = "3"
					lPrgrmResp.MethodName = "NSE Fetch-SGB-Master"
					lPrgrmResp.Status = common.SuccessCode
					lPrgrmResp.ErrMsg = common.SUCCESS
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("NSE SGB Master Fetched Successfully")
				}
			}

			//NCB master fetch in Nse exchange
			lNcbNseToken, lErr5 := exchangecall.FetchNcbMasterNSE(lUser, lBrokerId)
			if lErr5 != nil {
				log.Println("ISFIMS05", lErr5)
				lPrgrmResp.SINo = "4"
				lPrgrmResp.MethodName = "NSE Fetch-NCB-Master"
				lPrgrmResp.Status = common.ErrorCode
				lPrgrmResp.ErrMsg = "ISFIMS05" + lErr5.Error()
				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
			} else {
				if lNcbNseToken == common.ErrorCode {
					lPrgrmResp.SINo = "4"
					lPrgrmResp.MethodName = "NSE Fetch-NCB-Master"
					lPrgrmResp.Status = common.ErrorCode
					lPrgrmResp.ErrMsg = lNoToken
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("NSE Fetch NCB Master Error")
				} else {
					lPrgrmResp.SINo = "4"
					lPrgrmResp.MethodName = "NSE Fetch-NCB-Master"
					lPrgrmResp.Status = common.SuccessCode
					lPrgrmResp.ErrMsg = common.SUCCESS
					lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
					// log.Println("NSE NCB Master Fetched Successfully")
				}
			}
		}
		// Marshal the Response Structure into lData
		lData, lErr6 := json.Marshal(lRespRec)
		if lErr6 != nil {
			log.Println("ISFIMS06", lErr6)
			fmt.Fprintf(w, "ISFIMS06"+lErr6.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("FetchIpoMasterSch (-)", r.Method)
	}
}
