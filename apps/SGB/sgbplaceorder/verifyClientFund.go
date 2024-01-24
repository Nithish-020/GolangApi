package sgbplaceorder

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/coresettings"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fcs23pkg/integration/clientfund"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type FundRespStruct struct {
	AccountBalance float64 `json:"accountBalance"`
	Status         string  `json:"status"`
	ErrMsg         string  `json:"errMsg"`
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
func FetchClientFund(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchClientFund (+)", r.Method)
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		var lRespRec FundRespStruct
		lRespRec.Status = common.SuccessCode
		lClientId, lErr1 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
		if lErr1 != nil {
			log.Println("SPFC01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "SPFC01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("SPFC01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "SPFC02 / UserDetails not Found"
				fmt.Fprintf(w, helpers.GetErrorString("SPFC02", "UserDetails not Found"))
				return
			}
		}

		// commented by Pavithra , old method for Fund checking
		// lFundResp, lErr2 := VerifyFundDetails(lClientId)
		lFundResp, lErr2 := VerifyFOFundDetails(lClientId)
		if lErr2 != nil {
			log.Println("SPFC03", lErr2.Error())
			lRespRec.Status = common.ErrorCode
			fmt.Fprintf(w, helpers.GetErrorString("SPFC03", "Unable to process your request now. Please try after sometime"))
			return
		} else {
			log.Println("lFundResp", lFundResp)
			lRespRec = lFundResp
		}

		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("SPFC04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("SPFC04", "Unable to getting response.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("FetchClientFund (-)", r.Method)
	}
}

// Written by Nithish
// func VerifyFOFundDetails(pClientId string) (FundRespStruct, error) {
// 	log.Println("VerifyFOFundDetails (+)")
// 	//
// 	var lRespRec FundRespStruct

// 	lToken, lErr := GetFOToken()
// 	if lErr != nil {
// 		log.Println("SPVFD01", lErr)
// 		return lRespRec, lErr
// 	} else {
// 		//
// 		lClientFundDetail, lErr1 := clientfund.VerifyFOFund(pClientId, lToken)
// 		if lErr1 != nil {
// 			log.Println("SPVFD02", lErr1)
// 			return lRespRec, lErr1
// 		} else {
// 			log.Println(lClientFundDetail)

// 			if lClientFundDetail.Cash == "" {
// 				lRespRec.AccountBalance = 0.0
// 				lRespRec.Status = common.ErrorCode
// 			} else {
// 				lCash, _ := strconv.ParseFloat(lClientFundDetail.Cash, 64)
// 				lPayin, _ := strconv.ParseFloat(lClientFundDetail.PayIn, 64)
// 				lPayout, _ := strconv.ParseFloat(lClientFundDetail.PayOut, 64)
// 				lBrokCollAmt, _ := strconv.ParseFloat(lClientFundDetail.BrkCollAmt, 64)

// 				lRespRec.AccountBalance = lCash + lPayin + lPayout - lBrokCollAmt
// 				lRespRec.Status = common.SuccessCode
// 			}
// 		} else {
// 			lRespRec.AccountBalance = 0.0
// 			lRespRec.Status = common.ErrorCode
// 		}
// 	}
// 	log.Println("VerifyFOFundDetails (-)")
// 	return lRespRec, nil
// }

// Written by Nithish
func VerifyFOFundDetails(pClientId string) (FundRespStruct, error) {
	log.Println("VerifyFOFundDetails (+)")

	var lRespRec FundRespStruct

	lClientFundDetail, lErr1 := clientfund.VerifyMaxPayout(pClientId)
	if lErr1 != nil {
		log.Println("SPVFD02", lErr1)
		return lRespRec, lErr1
	} else {

		log.Println(lClientFundDetail)
		if lClientFundDetail.Status == "Ok" {
			lRespRec.Status = common.SuccessCode
			lFloatVal, lErr3 := strconv.ParseFloat(lClientFundDetail.PayOut, 64)
			if lErr3 != nil {
				log.Println("SPVFD03", lErr3)
				return lRespRec, lErr3
			} else {
				log.Println("lFloatVal", lFloatVal)
				lRespRec.AccountBalance = lFloatVal
			}
		} else {
			lRespRec.Status = common.ErrorCode
			return lRespRec, nil
		}

	}
	log.Println("VerifyFOFundDetails (-)")
	return lRespRec, nil
}

/*
Purpose: This function is used to GET the FO token from the database

	Parameters:

	Response:

	On Success :
	===============
	{
		"susertoken": "6df6a7d90917b81e615b251bce4db2505a7df7b76127b61663c800dbbc6ca888",
	}
	===============

	On Error:
	===============
		{},error
	===============

Author:Nithish kumar
Date: 21Nov2023
*/
func GetFOToken() (string, error) {
	log.Println("GetFOToken (+)")
	lToken := ""
	// To Establish A database connection,call LocalDbConnect Method
	// commented by pavithra
	// lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	// if lErr1 != nil {
	// 	log.Println("SPISD01", lErr1)
	// 	return lToken, lErr1
	// } else {
	// 	defer lDb.Close()

	lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
	lTokenKey := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["AdminApiSUserToken"])

	lToken = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, lTokenKey)
	// lToken := "01fe9896de0dc7fb82e9a61bbbf9ead8d6502f1de34acdc16f1c0848bc4cfc2f"
	// }
	log.Println("GetFOToken (-)")
	return lToken, nil
}

// /*
// Purpose: This function is used to construct the details for the BO Jv API
// 	Parameters:
// 		{
// 			JvAmount    string `json:"jvAmount"`
// 			JvStatus    string `json:"jvStatus"`
// 			JvStatement string `json:"jvStatement"`
// 			JvType      string `json:"jvType"`
// 		}, "FT000069"

// 	Response:

// 	On Success :
// 	===============

// 	===============

// 	On Error:
// 	===============
// 		{},error
// 	===============

// Author:Nithish kumar
// Date: 21Nov2023

// */
// func BlockFOClientFund(pRequest SgbReqStruct, pJvType string, pClientId string, pStatusFlag string) (JvStatusStruct, error) {
// 	log.Println("BlockFOClientFund (+)")
// 	var lFOJvReqStruct FOJvStruct
// 	var lJvStatusRec JvStatusStruct
// 	var lUrl string

// 	config := common.ReadTomlConfig("./toml/config.toml")
// 	lFOJvReqStruct.UserId = fmt.Sprintf("%v", config.(map[string]interface{})["VerifyUser"])
// 	lFOJvReqStruct.SourceUserId = lFOJvReqStruct.UserId

// 	lToken, lErr1 := GetFOToken()
// 	if lErr1 != nil {
// 		log.Println("SBFCF01", lErr1.Error())
// 		return lJvStatusRec, lErr1
// 	} else {
// 		// Choose the url for the api based on the mode Debit or Credit
// 		if pJvType == "D" {
// 			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayoutUrl"])
// 			lJvStatusRec.JvType = pJvType
// 			if pStatusFlag != "R" {
// 				lFOJvReqStruct.Remarks = "AMOUNT HOLD FOR SGB ORDER"
// 				lFOJvReqStruct.Amount = "-" + strconv.Itoa(pRequest.Price*pRequest.Unit)
// 			} else {
// 				lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM SGB ORDER"
// 				lFOJvReqStruct.Amount = strconv.Itoa(pRequest.Price * pRequest.Unit)
// 			}
// 		} else if pJvType == "C" {
// 			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayinUrl"])
// 			lJvStatusRec.JvType = pJvType
// 			lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM SGB ORDER"
// 			lFOJvReqStruct.Amount = strconv.Itoa(pRequest.Price * pRequest.Unit)
// 		}

// 		lFOJvReqStruct.AccountId = pClientId
// 		lJvStatusRec.JvAmount = lFOJvReqStruct.Amount

// 		lRequest, lErr2 := json.Marshal(lFOJvReqStruct)
// 		if lErr2 != nil {
// 			log.Println("SBFCF02", lErr2.Error())
// 			return lJvStatusRec, lErr2
// 		} else {
// 			// construct the request body
// 			lBody := `jData=` + string(lRequest) + `&jKey=` + lToken

// 			lResp, lErr3 := clientfund.FOProcessJV(lUrl, lBody)
// 			if lErr3 != nil {
// 				log.Println("SBFCF03", lErr3.Error())
// 				return lJvStatusRec, lErr3

// 			} else {
// 				if lResp.Status == "Ok" {
// 					log.Println("FO JV processed successfully for : ", lRequest)
// 					lJvStatusRec.JvStatus = "S"
// 				} else {
// 					log.Println("FO JV processed Failed for : ", lRequest)
// 					lJvStatusRec.JvStatus = "E"
// 				}
// 			}

// 		}
// 	}
// 	log.Println("BlockFOClientFund (-)")
// 	return lJvStatusRec, nil
// }

func PaymentCode(pUser string) (string, error) {
	log.Println("PaymentCode (-)")
	var lCode string
	lRespRec, lErr2 := clientfund.GetPaymentCode(pUser)
	if lErr2 != nil {
		log.Println("SBFCF02", lErr2.Error())
		return lCode, lErr2
	} else {
		log.Println("lRespRec", lRespRec)
		for i := 0; i < len(lRespRec); i++ {
			if lRespRec[i].Segment != "" {
				lCode = lRespRec[i].Segment
			}
		}
	}
	log.Println("PaymentCode (-)")
	return lCode, nil
}
