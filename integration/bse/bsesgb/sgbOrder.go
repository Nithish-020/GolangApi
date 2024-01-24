package bsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type SgbReqStruct struct {
	ScripId          string            `json:"scripid"`
	PanNo            string            `json:"panno"`
	InvestorCategory string            `json:"invcategory"`
	ApplicantName    string            `json:"applicantname"`
	Depository       string            `json:"depository"`
	DpId             string            `json:"dpid"`
	ClientBenfId     string            `json:"clientid"`
	GuardianName     string            `json:"guardianname"`
	GuardianPanno    string            `json:"guardianpanno"`
	GuardianRelation string            `json:"guardianrelation"`
	Bids             []ReqSgbBidStruct `json:"bids"`
}

type ReqSgbBidStruct struct {
	BidId            string `json:"bidid"`
	SubscriptionUnit string `json:"subscriptionunit"`
	Rate             string `json:"rate"`
	OrderNo          string `json:"orderno"`
	ActionCode       string `json:"actioncode"`
}

type SgbRespStruct struct {
	ScripId          string             `json:"scripid"`
	PanNo            string             `json:"panno"`
	InvestorCategory string             `json:"invcategory"`
	ApplicantName    string             `json:"applicantname"`
	Depository       string             `json:"depository"`
	DpId             string             `json:"dpid"`
	ClientBenfId     string             `json:"clientid"`
	GuardianName     string             `json:"guardianname"`
	GuardianPanno    string             `json:"guardianpanno"`
	GuardianRelation string             `json:"guardianrelation"`
	StatusCode       string             `json:"statuscode"`
	StatusMessage    string             `json:"statusmessage"`
	ErrorCode        string             `json:"errorcode"`
	ErrorMessage     string             `json:"errormessage"`
	Bids             []RespSgbBidStruct `json:"bids"`
}

type RespSgbBidStruct struct {
	BidId            string `json:"bidid"`
	SubscriptionUnit string `json:"subscriptionunit"`
	Rate             string `json:"rate"`
	OrderNo          string `json:"orderno"`
	ActionCode       string `json:"actioncode"`
	ErrorCode        string `json:"errorcode"`
	Message          string `json:"message"`
}

/*
Pupose: This method returns the active sgb list from the BSE Exchange
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e"
Response:
    *On Sucess
    =========
	{
		"scripid": "SGBTEST",
		"panno": "AGMPA85hdh75C",
		"invcategory": "CTZ",
		"applicantname": " Kumar",
		"depository": "cdsl",
		"dpid": "0",
		"clientid": "fsff12314",
		"guardianname": "",
		"guardianpanno": "",
		"guardianrelation": "",
		"bids": [
			{
				"bidid": "56618",
				"subscriptionunit": "1",
				"rate": "6000",
				"orderno": "12534",
				"actioncode": "D"
			}
		]
	}
    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 05JUNE2023
*/
func BseSgbOrder(pToken string, pUser string, pReqRec SgbReqStruct) (SgbRespStruct, error) {
	log.Println("BseSgbOrder (+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbRespStruct
	var lApiRespRec SgbRespStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseSgbOrder"])
	log.Println(lUrl, "endpoint")

	// lRequest, lErr := json.Marshal(pReqRec)
	// if lErr != nil {
	// 	log.Println("NT01", lErr)
	// 	return lApiRespRec, lErr
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "/v1/sgbOrder"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// 	// ! LogEntry method is used to store the Request in Database
	// 	lId, lErr1 = Function.LogEntry(lLogInputRec)
	// 	if lErr1 != nil {
	// 		log.Println("NSM01", lErr1)
	// 		return lApiRespRec, lErr1
	// 	} else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := BseSgbOrderApi(pToken, lUrl, pReqRec)
	if lErr2 != nil {
		log.Println("NSM02", lErr2)
		return lApiRespRec, lErr2
	} else {
		lApiRespRec = lResp
		log.Println("Response", lResp)
	}
	// Store thre Response in Log table
	// lResponse, lErr3 := json.Marshal(lResp)
	// if lErr3 != nil {
	// 	log.Println("NSM03", lErr3)
	// 	return lApiRespRec, lErr3
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	// create instance to hold errors
	// 	var lErr4 error
	// 	lId, lErr4 = Function.LogEntry(lLogInputRec)
	// 	if lErr4 != nil {
	// 		log.Println("NSM04", lErr4)
	// 		return lApiRespRec, lErr4
	// 	}
	// }
	// }
	// }
	log.Println("BseSgbOrder (-)")
	return lApiRespRec, nil
}

func BseSgbOrderApi(pToken string, pUrl string, pReqRec SgbReqStruct) (SgbRespStruct, error) {
	log.Println("BseSgbOrderApi (+)")
	// create a new instance of IpoResponseStruct
	var lSgbRespRec SgbRespStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lDetail, lErr := AccessDetail()
	if lErr != nil {
		log.Println("NESM01", lErr)
		return lSgbRespRec, lErr
	} else {

		lConsHeadRec.Key = "MemberCode"
		lConsHeadRec.Value = lDetail.MemberCode
		lHeaderArr = append(lHeaderArr, lConsHeadRec)

		lConsHeadRec.Key = "Login"
		lConsHeadRec.Value = lDetail.LoginId
		lHeaderArr = append(lHeaderArr, lConsHeadRec)

		lConsHeadRec.Key = "Token"
		lConsHeadRec.Value = pToken
		lHeaderArr = append(lHeaderArr, lConsHeadRec)

		lConsHeadRec.Key = "Content-Type"
		lConsHeadRec.Value = "application/json"
		lHeaderArr = append(lHeaderArr, lConsHeadRec)

		ljsonData, lErr := json.Marshal(pReqRec)
		if lErr != nil {
			log.Println("NEO01", lErr)
			return lSgbRespRec, lErr
		} else {
			lReqstring := string(ljsonData)
			lResp, lErr1 := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "SGBOrder")
			if lErr1 != nil {
				log.Println("NESM02", lErr1)
				return lSgbRespRec, lErr1
			} else {
				lErr2 := json.Unmarshal([]byte(lResp), &lSgbRespRec)
				if lErr2 != nil {
					log.Println("NESM03", lErr2)
					return lSgbRespRec, lErr2
				}
			}
		}
	}
	log.Println("BseSgbOrderApi (-)")
	return lSgbRespRec, nil
}
