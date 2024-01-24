package clientfund

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type PaymentCodeStruct struct {
	Segment string `json:"segment"`
}

/*
Purpose: This function is used to fetch the client account paymentCode
parameters: "pUser" = "FT045679"
Response:

	============
	On Success:
	============
		[{"segment":"NSE_CASH"}]
	==========
	On Error:
	==========
		[{"segment":null}]

Author: Pavithra
Date: 19DEC2023
*/
func GetPaymentCode(pUser string) ([]PaymentCodeStruct, error) {
	log.Println("GetPaymentCode (+)")
	// Create instance for Parameter struct
	// var lLogInputRec Function.ParameterStruct
	// Create instance for RespStruct
	var lApiRespRec []PaymentCodeStruct
	// create instance to hold the last inserted id
	// var lId int
	//To link the toml file
	lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["PaymentCode_API"])
	lUrl = lUrl + pUser
	// To get the details for v1/login from database
	// Marshalling the structure for LogEntry method
	// lRequest, lErr1 := json.Marshal(lUrl)
	// if lErr1 != nil {
	// 	log.Println("CFGPC01", lErr1)
	// 	return lApiRespRec, lErr1
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "/query/?qc=CLSGMT&typ=v&qp1=" + pUser
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "GET"

	// 	// LogEntry method is used to store the Request in Database
	// 	lId, lErr1 = Function.LogEntry(lLogInputRec)
	// 	if lErr1 != nil {
	// 		log.Println("CFGPC02", lErr1)
	// 		return lApiRespRec, lErr1
	// 	} else {
	// ClientFundApi used to fetch the client fund details
	lResp, lErr2 := paymentCodeApi(lUrl)
	if lErr2 != nil {
		log.Println("CFGPC03", lErr2)
		return lApiRespRec, lErr2
	} else {
		lApiRespRec = lResp
		// Store thre Response in Log table
		// lResponse, lErr3 := json.Marshal(lResp)
		// if lErr3 != nil {
		// 	log.Println("CFGPC04", lErr3)
		// 	return lApiRespRec, lErr3
		// } else {
		// 	lLogInputRec.Response = string(lResponse)
		// 	lLogInputRec.LastId = lId
		// 	lLogInputRec.Flag = common.UPDATE

		// 	_, lErr3 = Function.LogEntry(lLogInputRec)
		// 	if lErr3 != nil {
		// 		log.Println("CFGPC05", lErr3)
		// 		return lApiRespRec, lErr3
		// 	}
		// }
		// }
		// }
	}
	log.Println("GetPaymentCode (-)")
	return lApiRespRec, nil
}

func paymentCodeApi(pUrl string) ([]PaymentCodeStruct, error) {
	log.Println("paymentCodeApi (+)")
	//create instance for PaymentCode response struct
	var lRespArr []PaymentCodeStruct
	//create array of instance to store the key value pairs
	var lHeaderArr []apiUtil.HeaderDetails

	lResp, lErr := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "Client_PaymentCode")
	if lErr != nil {
		log.Println("CFPCA01", lErr)
		return lRespArr, lErr
	} else {
		log.Println("Client fund response := ", lResp)
		// Unmarshalling json to struct
		lErr = json.Unmarshal([]byte(lResp), &lRespArr)
		if lErr != nil {
			log.Println("CFPCA02", lErr)
			return lRespArr, lErr
		}
	}
	log.Println("paymentCodeApi (-)")
	return lRespArr, nil
}
