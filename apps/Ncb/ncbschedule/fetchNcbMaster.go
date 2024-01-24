package ncbschedule

// import (
// 	"encoding/json"
// 	"fcs23pkg/apps/Ipo/brokers"
// 	"fcs23pkg/apps/exchangecall"
// 	"fcs23pkg/common"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// type ProgramStruct struct {
// 	Slno       string `json:"slno"`
// 	MethodName string `json:"methodName"`
// 	Status     string `json:"status"`
// 	ErrMsg     string `json:"errMsg"`
// }

// // this response struct for the fetchNcbMaster API
// type NcbMasterRespStruct struct {
// 	ResponseArr []ProgramStruct `json:"responseArr"`
// 	Status      string          `json:"status"`
// 	ErrMsg      string          `json:"errMsg"`
// }

// /*
// Pupose:This Function is used to Get the Active NCB Details from Exchange and insert it in our database table ....
// Parameters:

// 	not Applicable

// Response:

// 	*On Sucess
// 	=========
// 	    Success: {
//            "data": [
//                     {
//                         "biddingStartDate": "01-11-2018",
//                         "symbol": "TEST",
//                         "minBidQuantity": 10,
//                         "maxQuantity": 20000000,
//                         "lotSize": 10,
//                         "t1ModEndDate": "02-11-2018",
//                         "dailyStartTime": "15:30:00",
//                         "allotmentDate": "",
//                         "t1ModStartTime": "20:30:00",
//                         "biddingEndDate": "10-12-2018",
// 			            "t1ModEndTime": "21:00:00",
//                         "dailyEndTime": "11:30:00",
//                         "tickSize": 100,
//                         "cutoffPrice": "0205GS",
//                         "series": "GS",
//                         "faceValue": 1000,
//                         "minPrice": 0,
//                         "t1ModStartDate": "02-11-2018",
//                         "issueValueSize": 0,
//                         "name": "0205GS",
//                         "issueSize": 100000000,
//                         "lastDayBiddingEndTime": "",
//                         "maxPrice": 0,
//                         "isin": "IN0020150044"
//                     },

// 		        ],
// 		    "status": "success"
// 	    }

// 	!On Error
// 	========
// 	{
// 		"status": "E",
// 		"reason": "Can't able to get the Ncb details"
// 	}

// Author: Kavya Dharshani
// Date: 04 OCT 2023
// */
// func FetchNcbMasterSch(w http.ResponseWriter, r *http.Request) {
// 	log.Println("FetchNcbMasterSch (+)", r.Method)
// 	origin := r.Header.Get("Origin")
// 	var lBrokerId int
// 	var lErr error
// 	// for _, allowedOrigin := range common.ABHIAllowOrigin {
// 	// 	if allowedOrigin == origin {
// 	w.Header().Set("Access-Control-Allow-Origin", origin)
// 	lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
// 	// log.Println(lErr, origin)
// 	if lErr != nil {
// 		log.Println(lErr, origin)
// 	}
// 	// 	}
// 	// }
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "GET" {

// 		log.Println("r.Method", r.Method)
// 		//create instance for manualStruct
// 		var lRespRec NcbMasterRespStruct

// 		var lPrgrmResp ProgramStruct
// 		//get header value
// 		lUser := r.Header.Get("USER")

// 		log.Println("lUser", lUser)
// 		// validate the user and check if the user role is admin and then allow for next process
// 		lRespRec.Status = common.SuccessCode

// 		// Calling the FetchNcbMasterSch method to get the Active NCB details From exchange and
// 		// then store the details  in the database
// 		lNoToken := "Access Token not found"
// 		lNSEToken, lErr1 := exchangecall.FetchNCBMaster(lUser, lBrokerId)
// 		if lErr1 != nil {
// 			log.Println("NSFNM01", lErr1)
// 			lPrgrmResp.Slno = "1"
// 			lPrgrmResp.MethodName = "NSE Fetch-NSB-Master"
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "NSFNM01" + lErr1.Error()
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 		} else {
// 			if lNSEToken == common.ErrorCode {
// 				lPrgrmResp.Slno = "1"
// 				lPrgrmResp.MethodName = "NSE Fetch-NCB-Master"
// 				lPrgrmResp.Status = common.ErrorCode
// 				lPrgrmResp.ErrMsg = lNoToken
// 				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 				log.Println("NSE NCB Details FETCHING ERROR")

// 			} else {

// 				lPrgrmResp.Slno = "1"
// 				lPrgrmResp.MethodName = "NSE Fetch-NCB-Master"
// 				lPrgrmResp.Status = common.SuccessCode
// 				lPrgrmResp.ErrMsg = common.SUCCESS
// 				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 				log.Println("NSE NCBs Details Fetched Successfully")
// 			}
// 		}

// 		// Marshal the Response Structure into lData
// 		lData, lErr2 := json.Marshal(lRespRec)
// 		if lErr2 != nil {
// 			log.Println("NSFNM02", lErr2)
// 			fmt.Fprintf(w, "NSFNM02"+lErr2.Error())
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("FetchNcbMasterSch (-)", r.Method)
// 	}
// }
