package bseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type BseIpoCategoryStruct struct {
	ScripId   string `json:"scripid"`
	Category  string `json:"category"`
	Series    string `json:"series"`
	Type      string `json:"type"`
	OpenDate  string `json:"opendate"`
	OpenTime  string `json:"opentime"`
	CloseDate string `json:"closedate"`
	CloseTime string `json:"closetime"`
	Message   string `json:"message"`
}

/*
Pupose: This method returns the active sgb list from the BSE Exchange
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e"
Response:
    *On Sucess
    ==========

    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 05JUNE2023
*/
func BseIpoCategory(pToken string, pUser string, pBrokerId int) ([]BseIpoCategoryStruct, error) {
	log.Println("BseIpoCategory (+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbRespStruct
	var lApiRespArr []BseIpoCategoryStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseIpoCategory"])
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "bse/v1/ipocategory"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = pUser
	// lLogInputRec.Method = "GET"
	// ! LogEntry method is used to store the Request in Database
	// lId, lErr1 = Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("BIM01", lErr1)
	// 	return lApiRespArr, lErr1
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := BseIpoCategoryApi(pToken, lUrl, pBrokerId)
	if lErr2 != nil {
		log.Println("BIM02", lErr2)
		return lApiRespArr, lErr2
	} else {
		lApiRespArr = lResp
	}
	// log.Println("Response", lResp)
	// Store thre Response in Log table
	// lResponse, lErr3 := json.Marshal(lResp)
	// if lErr3 != nil {
	// 	log.Println("BIM03", lErr3)
	// 	return lApiRespArr, lErr3
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	// create instance to hold errors
	// 	var lErr4 error
	// 	lId, lErr4 = Function.LogEntry(lLogInputRec)
	// 	if lErr4 != nil {
	// 		log.Println("BIM04", lErr4)
	// 		return lApiRespArr, lErr4
	// 	}
	// }
	// }
	// log.Println(lApiRespArr)
	log.Println("BseIpoCategory (-)")
	return lApiRespArr, nil
}

func BseIpoCategoryApi(pToken string, pUrl string, pBrokerId int) ([]BseIpoCategoryStruct, error) {
	log.Println("BseIpoCategoryApi (+)")
	// create a new instance of IpoResponseStruct
	var lIpoRespArr []BseIpoCategoryStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lDetail, lErr := AccessBseCredential(pBrokerId)
	if lErr != nil {
		log.Println("NESM01", lErr)
		return lIpoRespArr, lErr
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

		lResp, lErr1 := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "Bse-IPOCategory")
		if lErr1 != nil {
			log.Println("NESM02", lErr1)
			return lIpoRespArr, lErr1
		} else {
			lErr2 := json.Unmarshal([]byte(lResp), &lIpoRespArr)
			if lErr2 != nil {
				log.Println("NESM03", lErr2)
				return lIpoRespArr, lErr2
			}
		}
	}
	log.Println("BseIpoCategoryApi (-)")
	return lIpoRespArr, nil
}
