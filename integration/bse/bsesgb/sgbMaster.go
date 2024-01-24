package bsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type BseSgbMasterStruct struct {
	IsssueName      string `json:"isssuename"`
	Symbol          string `json:"symbol"`
	SerialNo        string `json:"serialno"`
	Isin            string `json:"isin"`
	OpenDate        string `json:"opendate"`
	CloseDate       string `json:"closedate"`
	IssuePrice      string `json:"issueprice"`
	MinQty          string `json:"minqty"`
	MaxQty          string `json:"maxqty"`
	MaxQtyTrust     string `json:"maxqtytrust"`
	DateOfAllotment string `json:"dateofallotment"`
	ErrorCode       string `json:"errorcode"`
	Message         string `json:"message"`
}

/*
Pupose: This method returns the active sgb list from the BSE Exchange
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
func BseSgbMaster(pToken string, lUser string) ([]BseSgbMasterStruct, error) {
	log.Println("BseSgbMaster (+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbRespStruct
	var lApiRespArr []BseSgbMasterStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseSgbMaster"])
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "bse/v1/sgbmaster"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = lUser
	// lLogInputRec.Method = "GET"

	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr1 = Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("NSM01", lErr1)
	// 	return lApiRespArr, lErr1
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := BseSgbApi(pToken, lUrl)
	if lErr2 != nil {
		log.Println("NSM02", lErr2)
		return lApiRespArr, lErr2
	} else {
		lApiRespArr = lResp
		log.Println("Response", lResp)
	}
	// Store thre Response in Log table
	// lResponse, lErr3 := json.Marshal(lResp)
	// if lErr3 != nil {
	// 	log.Println("NSM03", lErr3)
	// 	return lApiRespArr, lErr3
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	// create instance to hold errors
	// 	var lErr4 error
	// 	lId, lErr4 = Function.LogEntry(lLogInputRec)
	// 	if lErr4 != nil {
	// 		log.Println("NSM04", lErr4)
	// 		return lApiRespArr, lErr4
	// 	}
	// }
	// }
	// log.Println(lApiRespArr)
	log.Println("BseSgbMaster (-)")
	return lApiRespArr, nil
}

func BseSgbApi(pToken string, pUrl string) ([]BseSgbMasterStruct, error) {
	log.Println("BseSgbApi (+)")
	// create a new instance of IpoResponseStruct
	var lSgbRespArr []BseSgbMasterStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lDetail, lErr := AccessDetail()
	if lErr != nil {
		log.Println("NESM01", lErr)
		return lSgbRespArr, lErr
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

		lResp, lErr1 := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "SGBMaster")
		if lErr1 != nil {
			log.Println("NESM02", lErr1)
			return lSgbRespArr, lErr1
		} else {
			lErr2 := json.Unmarshal([]byte(lResp), &lSgbRespArr)
			if lErr2 != nil {
				log.Println("NESM03", lErr2)
				return lSgbRespArr, lErr2
			}
		}
	}
	log.Println("BseSgbApi (-)")
	return lSgbRespArr, nil
}
