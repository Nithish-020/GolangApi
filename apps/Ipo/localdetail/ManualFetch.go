package localdetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// this struct is used to give response for ManualFetch API
type manualStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

/*
Pupose:This API is used to fetch the ipo master details from exchange
Parameters:
	not Applicable
Response:
	==========
	*On Sucess
	==========
	{
		"status": "S",
		"errMsg": ""
	}
	=========
	!On Error
	=========
	{
		"status": "E",
		"reason": "Can't able to get the ipo details"
	}
Author: Nithish Kumar
Date: 08JUNE2023
*/
func ManualFetch(w http.ResponseWriter, r *http.Request) {
	log.Println("ManualFetch (+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}
	log.Println("lBrokerId", lBrokerId)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		//create instance for manualStruct
		var lManualResp manualStruct
		lManualResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/fetchmaster")
		if lErr1 != nil {
			log.Println("LMF01", lErr1)
			lManualResp.Status = common.ErrorCode
			lManualResp.ErrMsg = "LMF01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LMF01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lManualResp.Status = common.ErrorCode
				lManualResp.ErrMsg = "LMF02 / UserDetails not Found"
				fmt.Fprintf(w, helpers.GetErrorString("LMF02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// added by pavithra FOR GETTING hardcoded brokerid
		lConfigFile := common.ReadTomlConfig("./toml/debug.toml")
		lBroker := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["FetchMasterBrokerId"])

		lBrokerId1, lErr5 := strconv.Atoi(lBroker)
		if lErr5 != nil {
			log.Println("LMF05", lErr5)
			lManualResp.Status = common.ErrorCode
			lManualResp.ErrMsg = "LMF05" + lErr5.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LMF05", "Unable to Get Details. Please try after sometime"))
			return
		} else {

			// Calling the FetchIpomaster method to get the Active Ipo details From exchange and
			// then store the details  in the database
			lFetchIpoErrCode, lErr2 := exchangecall.FetchIPOmaster(lClientId, lBrokerId1)
			if lErr2 != nil {
				log.Println("LMF03", lErr2)
				lManualResp.Status = common.ErrorCode
				lManualResp.ErrMsg = "LMF03" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("LMF03", "Unable to Fetch IPO Details. Please try after sometime"))
				return
			} else {
				if lFetchIpoErrCode != common.ErrorCode {
					lManualResp.Status = common.SuccessCode
				} else {
					lManualResp.Status = common.ErrorCode
				}
				//fetch sgb master from exchange
				lFetchSgbErrCode, lErr3 := exchangecall.FetchSGBMaster(lClientId, lBrokerId1)
				if lErr3 != nil {
					log.Println("LMF04", lErr3)
					lManualResp.Status = common.ErrorCode
					lManualResp.ErrMsg = "LMF04" + lErr3.Error()
					fmt.Fprintf(w, helpers.GetErrorString("LMF04", "Unable to Fetch SGB Details. Please try after sometime"))
					return
				} else {
					if lFetchSgbErrCode != common.ErrorCode {
						lManualResp.Status = common.SuccessCode
					} else {
						lManualResp.Status = common.ErrorCode
					}
					// fetch ncb master
					lNCBFetchErrCode, lErr2 := exchangecall.FetchNCBMaster(lClientId, lBrokerId)
					if lErr2 != nil {
						log.Println("LMF03", lErr2)
						lManualResp.Status = common.ErrorCode
						lManualResp.ErrMsg = "LMF03" + lErr2.Error()
						fmt.Fprintf(w, helpers.GetErrorString("LMF03", "Unable to Fetch NCB Details. Please try after sometime"))
						return
					} else {
						if lNCBFetchErrCode != common.ErrorCode {
							lManualResp.Status = common.SuccessCode
						} else {
							lManualResp.Status = common.ErrorCode
						}
					}
				}
			}
			// 	//fetch sgb master from exchange
			// 	// lErr3 := exchangecall.FetchSGBMaster(lClientId, lBrokerId)
			// 	// if lErr3 != nil {
			// 	// 	log.Println("LMF04", lErr3)
			// 	// 	lManualResp.Status = common.ErrorCode
			// 	// 	lManualResp.ErrMsg = "LMF04" + lErr3.Error()
			// 	// 	fmt.Fprintf(w, helpers.GetErrorString("LMF04", "Unable to Fetch SGB Details. Please try after sometime"))
			// 	// 	return
			// 	// } else {
			// }

		}
		// Marshal the Response Structure into lData
		lData, lErr4 := json.Marshal(lManualResp)
		if lErr4 != nil {
			log.Println("LMF05", lErr4)
			fmt.Fprintf(w, helpers.GetErrorString("LMF05", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("ManualFetch (-)", r.Method)

	}
}
