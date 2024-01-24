package bseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type BseIpoMasterStruct struct {
	Symbol                string `json:"symbol"`
	Name                  string `json:"name"`
	ISIN                  string `json:"isin"`
	Category              string `json:"category"`
	IssueType             string `json:"issuetype"`
	OpenDateTime          string `json:"openDateTime"`
	CloseDateTime         string `json:"closeDateTime"`
	FloorPrice            string `json:"floorprice"`
	CeilingPrice          string `json:"ceilingprice"`
	CutOff                string `json:"cuttoff"`
	TickPrice             string `json:"tickprice"`
	MinBidQty             string `json:"minbidqty"`
	MaxBidQty             string `json:"maxbidqty"`
	TradingLot            string `json:"tradinglot"`
	MinValue              string `json:"minvalue"`
	MaxValue              string `json:"maxvalue"`
	DiscountType          string `json:"discounttype"`
	DiscountValue         string `json:"discountvalue"`
	AsbaNonAsba           string `json:"asbanonasba"`
	TplusModificationFrom string `json:"tplusmodificationfrom"`
	TplusModificationTo   string `json:"tplusmodificationto"`
	IssueSize             string `json:"issuesize"`
	ErrorCode             string `json:"errorcode"`
	Message               string `json:"message"`
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
func BseIpoMaster(pToken string, pUser string, pBrokerId int) ([]BseIpoMasterStruct, error) {
	log.Println("BseIpoMaster (+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbRespStruct
	var lApiRespArr []BseIpoMasterStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseIpoMaster"])
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "bse/v1/ipomaster"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = pUser
	// lLogInputRec.Method = "GET"
	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr1 = Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("BIM01", lErr1)
	// 	return lApiRespArr, lErr1
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := BseIpoApi(pToken, lUrl, pBrokerId)
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
	log.Println("BseSgbMaster (-)")
	return lApiRespArr, nil
}

func BseIpoApi(pToken string, pUrl string, pBrokerId int) ([]BseIpoMasterStruct, error) {
	log.Println("BseIpoApi (+)")
	// create a new instance of IpoResponseStruct
	var lIpoRespArr []BseIpoMasterStruct
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

		lResp, lErr1 := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "Bse-IPOMaster")
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
	log.Println("BseIpoApi (-)")
	return lIpoRespArr, nil
}
