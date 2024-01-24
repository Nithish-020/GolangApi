package nsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type SgbDetailStruct struct {
	Symbol               string  `json:"symbol"`
	Series               string  `json:"series"`
	Name                 string  `json:"name"`
	IssueType            string  `json:"issueType"`
	LotSize              int     `json:"lotSize"`
	FaceValue            float32 `json:"faceValue"`
	MinBidQuantity       int     `json:"minBidQuantity"`
	MinPrice             float32 `json:"minPrice"`
	MaxPrice             float32 `json:"maxPrice"`
	TickSize             float32 `json:"tickSize"`
	BiddingStartDate     string  `json:"biddingStartDate"`
	BiddingEndDate       string  `json:"biddingEndDate"`
	DailyStartTime       string  `json:"dailyStartTime"`
	DailyEndTime         string  `json:"dailyEndTime"`
	T1ModStartDate       string  `json:"t1ModStartDate"`
	T1ModEndDate         string  `json:"t1ModEndDate"`
	T1ModStartTime       string  `json:"t1ModStartTime"`
	T1ModEndTime         string  `json:"t1ModEndTime"`
	ISIN                 string  `json:"isin"`
	IssueSize            int     `json:"issueSize"`
	IssueValueSize       int     `json:"issueValueSize"`
	MaxQuantity          int     `json:"maxQuantity"`
	AllotmentDate        string  `json:"allotmentDate"`
	IncompleteModEndDate string  `json:"incompleteModEndDate"`
}

type SgbRespStruct struct {
	Data   []SgbDetailStruct `json:"data"`
	Reason string            `json:"reason"`
	Status string            `json:"status"`
}

/*
Pupose: This method returns the active ipo list from the NSE Exchange
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e"
Response:
    *On Sucess
    =========

    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 05JUNE2023
*/
func SgbMaster(pToken string, lUser string) (SgbRespStruct, error) {
	log.Println("SgbMaster...(+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbRespStruct
	var lApiRespRec SgbRespStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SgbMaster"])
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "nse/v1/sgbmaster"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = lUser
	// lLogInputRec.Method = "GET"

	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr1 = Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("NSM01", lErr1)
	// 	return lApiRespRec, lErr1
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := ExchangeSgbMaster(pToken, lUrl)
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
	// log.Println(lApiRespRec)
	log.Println("SgbMaster...(-)")
	return lApiRespRec, nil
}

func ExchangeSgbMaster(pToken string, pUrl string) (SgbRespStruct, error) {
	log.Println("ExchangeSgbMaster....(+)")
	// create a new instance of IpoResponseStruct
	var lSgbRespRec SgbRespStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lConsHeadRec.Key = "Access-Token"
	lConsHeadRec.Value = pToken
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "Content-Type"
	lConsHeadRec.Value = "application/json"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "User-Agent"
	lConsHeadRec.Value = "Flattrade-golang"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lResp, lErr1 := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "SGBMaster")
	if lErr1 != nil {
		log.Println("NESM01", lErr1)
		return lSgbRespRec, lErr1
	} else {
		lErr2 := json.Unmarshal([]byte(lResp), &lSgbRespRec)
		if lErr2 != nil {
			log.Println("NESM02", lErr2)
			return lSgbRespRec, lErr2
		}
	}
	log.Println("ExchangeSgbMaster....(-)")
	return lSgbRespRec, nil
}
