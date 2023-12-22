package clientfund

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

//commented by pavithra
// THIS THE THE OLDER VERSION TO FETCH CLIENT FUND

// type ClientFundStruct struct {
// 	AccountCode string `json:"ACCOUNTCODE"`
// 	Amount      string `json:"amount"`
// }

// func VerifyFund(pUser string) ([]ClientFundStruct, error) {
// 	log.Println("VerifyFund (+)")
// 	// Create instance for Parameter struct
// 	var lLogInputRec Function.ParameterStruct
// 	// Create instance for loinRespStruct
// 	var lApiRespRec []ClientFundStruct
// 	// create instance to hold the last inserted id
// 	var lId int
// 	//To link the toml file
// 	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
// 	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["VerifyFund"])
// 	lUrl = lUrl + pUser
// 	// To get the details for v1/login from database
// 	// Marshalling the structure for LogEntry method
// 	lRequest, lErr := json.Marshal(lUrl)
// 	if lErr != nil {
// 		log.Println("CVF01", lErr)
// 		return lApiRespRec, lErr
// 	} else {
// 		lLogInputRec.Request = string(lRequest)
// 		lLogInputRec.EndPoint = "/query/?qc=CFNDBL&typ=v&qp1=" + pUser
// 		lLogInputRec.Flag = common.INSERT
// 		lLogInputRec.ClientId = pUser
// 		lLogInputRec.Method = "GET"

// 		// LogEntry method is used to store the Request in Database
// 		lId, lErr = Function.LogEntry(lLogInputRec)
// 		if lErr != nil {
// 			log.Println("CVF02", lErr)
// 			return lApiRespRec, lErr
// 		} else {
// 			// ClientFundApi used to fetch the client fund details
// 			lResp, lErr := ClientFundApi(lUrl)
// 			if lErr != nil {
// 				log.Println("CVF03", lErr)
// 				return lApiRespRec, lErr
// 			}
// 			lApiRespRec = lResp
// 			// Store thre Response in Log table
// 			lResponse, lErr := json.Marshal(lResp)
// 			if lErr != nil {
// 				log.Println("CVF04", lErr)
// 				return lApiRespRec, lErr
// 			} else {
// 				lLogInputRec.Response = string(lResponse)
// 				lLogInputRec.LastId = lId
// 				lLogInputRec.Flag = common.UPDATE

// 				lId, lErr = Function.LogEntry(lLogInputRec)
// 				if lErr != nil {
// 					log.Println("CVF05", lErr)
// 					return lApiRespRec, lErr
// 				}
// 			}
// 		}
// 	}
// 	log.Println("VerifyFund (-)")
// 	return lApiRespRec, nil
// }

// func ClientFundApi(pUrl string) ([]ClientFundStruct, error) {
// 	log.Println("ClientFundApi (+)")
// 	//create instance for loginResp struct
// 	var lFundRespRec []ClientFundStruct
// 	//create array of instance to store the key value pairs
// 	var lHeaderArr []apiUtil.HeaderDetails

// 	lResp, lErr := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "ClientFund")
// 	if lErr != nil {
// 		log.Println("CCFA01", lErr)
// 		return lFundRespRec, lErr
// 	} else {
// 		// Unmarshalling json to struct
// 		lErr = json.Unmarshal([]byte(lResp), &lFundRespRec)
// 		if lErr != nil {
// 			log.Println("CCFA02", lErr)
// 			return lFundRespRec, lErr
// 		}
// 	}
// 	log.Println("ClientFundApi (-)")
// 	return lFundRespRec, nil
// }

type ClientFundReq struct {
	UserId    string `json:"uid"`
	AccountId string `json:"actid"`
}
type ClientFundResp struct {
	RequestTime      string `json:"request_time"`
	Status           string `json:"stat"`
	AccountId        string `json:"actid"`
	PrfName          string `json:"prfname"`
	Cash             string `json:"cash"`
	PayIn            string `json:"payin"`
	PayOut           string `json:"payout"`
	BrkCollAmt       string `json:"brkcollamt"`
	UnclearedCash    string `json:"unclearedcash"`
	AuxDayCash       string `json:"aux_daycash"`
	AuxBrkCollAmt    string `json:"aux_brkcollamt"`
	AuxUnclearedCash string `json:"aux_unclearedcash"`
	DayCash          string `json:"daycash"`
	TurnoverLmt      string `json:"turnoverlmt"`
	PendordValLmt    string `json:"pendordvallmt"`
	BlockAmt         string `json:"blockamt"`
	ErrMsg           string `json:"emsg"`
}
type ClientFund struct {
	AccountCode string `json:"ACCOUNTCODE"`
	Amount      string `json:"amount"`
}

/*
Purpose: This function is used to fetch the client fund information

	and it is the latest version available

parameters: "pUser" = "FT000069"
Response:

	============
	On Success:
	============
		{
			"request_time": "10:37:23 22-11-2023",
			"stat": "Ok",
			"actid": "FT000069",
			"prfname": "3TIMES",
			"cash": "0.00",
			"payin": "100.00",
			"payout": "0.00",
			"brkcollamt": "0.00",
			"unclearedcash": "0.00",
			"aux_daycash": "0.00",
			"aux_brkcollamt": "0.00",
			"aux_unclearedcash": "0.00",
			"daycash": "0.00",
			"turnoverlmt": "100000000000.00",
			"pendordvallmt": "100000000000.00",
			"blockamt": "0.00"
		}
	==========
	On Error:
	==========
		{
			"stat": "Not_Ok",
			"emsg": "Error Occurred"
		}

Author: Nithish Kumar
Date: 21NOV2023
*/
func VerifyFOFund(pUser string, pToken string) (ClientFundResp, error) {
	log.Println("VerifyFOFund (+)")
	// Create instance for Parameter struct
	var lLogInputRec Function.ParameterStruct
	// Create instance for clientReqStruct
	var lApiReqRec ClientFundReq
	// Create instance for clientRespStruct
	var lApiRespRec ClientFundResp
	// create instance to hold the last inserted id
	//commented by Lakshmanan on 17 Dec 2023
	//bug to have below hardcoded value in prod
	// pUser = "FT000069"
	//To link the toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["VerifyFOFund"])
	lUser := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["VerifyUser"])

	lLogInputRec.Request = lUrl
	lLogInputRec.EndPoint = "/NorenWCAdmin/Limits"
	lLogInputRec.Flag = common.INSERT
	lLogInputRec.ClientId = pUser
	lLogInputRec.Method = "POST"

	// LogEntry method is used to store the Request in Database
	lId, lErr := Function.LogEntry(lLogInputRec)
	if lErr != nil {
		log.Println("CVF02", lErr)
		return lApiRespRec, lErr
	} else {
		lApiReqRec.UserId = lUser
		lApiReqRec.AccountId = pUser

		lRequest, lErr := json.Marshal(lApiReqRec)
		if lErr != nil {
			log.Println("CVF04", lErr)
			return lApiRespRec, lErr
		} else {

			lBody := `jData=` + string(lRequest) + `&jKey=` + pToken

			// ClientFundApi used to fetch the client fund details
			lResp, lErr := ClientFundApi(lUrl, lBody)
			if lErr != nil {
				log.Println("CVF03", lErr)
				return lApiRespRec, lErr
			}
			lApiRespRec = lResp
			// Store thre Response in Log table
			lResponse, lErr := json.Marshal(lResp)
			if lErr != nil {
				log.Println("CVF04", lErr)
				return lApiRespRec, lErr
			} else {
				lLogInputRec.Response = string(lResponse)
				lLogInputRec.LastId = lId
				lLogInputRec.Flag = common.UPDATE

				lId, lErr = Function.LogEntry(lLogInputRec)
				if lErr != nil {
					log.Println("CVF05", lErr)
					return lApiRespRec, lErr
				}
			}
		}
	}
	log.Println("VerifyFOFund (-)")
	return lApiRespRec, nil
}

func ClientFundApi(pUrl string, pBody string) (ClientFundResp, error) {
	log.Println("ClientFundApi (+)")
	//create instance for clientFundResp struct
	var lFundRespRec ClientFundResp
	//create array of instance to store the key value pairs
	var lHeaderArr []apiUtil.HeaderDetails

	lResp, lErr := apiUtil.Api_call(pUrl, "POST", pBody, lHeaderArr, "ClientFund")
	if lErr != nil {
		log.Println("CCFA01", lErr)
		return lFundRespRec, lErr
	} else {
		// log.Println("Client fund response := ", lResp)
		// Unmarshalling json to struct
		lErr = json.Unmarshal([]byte(lResp), &lFundRespRec)
		if lErr != nil {
			log.Println("CCFA02", lErr)
			return lFundRespRec, lErr
		}
	}
	log.Println("ClientFundApi (-)")
	return lFundRespRec, nil
}
