package ncbplaceorder

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
Pupose:This API is used to fetch the NCB master details from exchange
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
		"reason": "Can't able to get the ncb details"
	}
Author: KAVYA DHARSHAMI
Date: 22NOV2023
*/

func NcbFetchClientFund(w http.ResponseWriter, r *http.Request) {
	log.Println("NcbFetchClientFund (+)", r.Method)

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
		lClientId, lErr1 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")

		if lErr1 != nil {
			log.Println("NFCF01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NFCF01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NFCF01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "NFCF02 / UserDetails not Found"
				fmt.Fprintf(w, helpers.GetErrorString("NFCF02", "UserDetails not Found"))
				return
			}
		}

		lFundResp, lErr3 := VerifyFundDetails(lClientId)
		if lErr3 != nil {
			log.Println("NFCF03", lErr3.Error())
			lRespRec.Status = common.ErrorCode
			fmt.Fprintf(w, helpers.GetErrorString("NFCF03", "Unable to process your request now. Please try after sometime"))
			return
		} else {
			log.Println("lFundResp", lFundResp)
			lRespRec = lFundResp

			log.Println("LFUNDDD -->", lRespRec)
		}

		lData, lErr4 := json.Marshal(lRespRec)
		if lErr4 != nil {
			log.Println("NFCF04", lErr4)
			fmt.Fprintf(w, helpers.GetErrorString("NFCF04", "Unable to getting response.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}

	log.Println("NcbFetchClientFund (-)", r.Method)
}

func VerifyFundDetails(pClientId string) (FundRespStruct, error) {

	log.Println("VerifyFundDetails(+)")

	var lRespRec FundRespStruct

	lClientFundDetail, lErr1 := clientfund.VerifyMaxPayout(pClientId)
	if lErr1 != nil {
		log.Println("NVFD02", lErr1)
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

	log.Println("VerifyFundDetails(-)")
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

Author: Kavya Dharshani M
Date: 19Nov2023
*/
func NcbGetFOToken() (string, error) {
	log.Println("NcbGetFOToken (+)")
	lToken := ""

	lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
	lTokenKey := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["AdminApiSUserToken"])

	lToken = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, lTokenKey)
	// lToken := "01fe9896de0dc7fb82e9a61bbbf9ead8d6502f1de34acdc16f1c0848bc4cfc2f"
	// }
	log.Println("NcbGetFOToken (-)")
	return lToken, nil
}

func NcbPaymentCode(pUser string) (string, error) {
	log.Println("NcbPaymentCode(+)")

	var lCode string
	lRespRec, lErr1 := clientfund.GetPaymentCode(pUser)
	if lErr1 != nil {
		log.Println("NBFCF01", lErr1.Error())
		return lCode, lErr1
	} else {
		log.Println("lRespRec", lRespRec)
		for i := 0; i < len(lRespRec); i++ {
			if lRespRec[i].Segment != "" {
				lCode = lRespRec[i].Segment
			}
		}
	}
	log.Println("NcbPaymentCode(-)")
	return lCode, nil
}
