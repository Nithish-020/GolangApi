package bseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type PswdReqStruct struct {
	OldPswd     string `json:"oldpwd"`
	NewPswd     string `json:"newpwd"`
	ConfirmPswd string `json:"confirmpwd"`
}

type PswdRespStruct struct {
	ErrorCode string `json:"errorcode"`
	Message   string `json:"message"`
}

func BsePassword(pReqRec PswdReqStruct, pUser string, pToken string, pBrokerId int) (PswdRespStruct, error) {
	log.Println("BsePassword (+)")
	// Create instance for Parameter struct
	// var lLogInputRec Function.ParameterStruct
	// Create instance for loinRespStruct
	var lApiRespRec PswdRespStruct
	// create instance to hold the last inserted id
	// var lId int
	//To link the toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BsePassword"])
	// To get the details for v1/login from database

	// Marshalling the structure for LogEntry method
	// lRequest, lErr1 := json.Marshal(pReqRec)
	// if lErr1 != nil {
	// 	log.Println("BIBP01", lErr1)
	// 	return lApiRespRec, lErr1
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "bse/v1/password"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// 	// LogEntry method is used to store the Request in Database
	// 	var lErr2 error
	// 	lId, lErr2 = Function.LogEntry(lLogInputRec)
	// 	if lErr2 != nil {
	// 		log.Println("BIBP02", lErr2)
	// 		return lApiRespRec, lErr2
	// 	} else {
	// TokenApi method used to call exchange API
	lResp, lErr3 := BsePasswordApi(pReqRec, lUrl, pToken, pBrokerId)
	if lErr3 != nil {
		log.Println("BIBP03", lErr3)
		return lApiRespRec, lErr3
	} else {
		lApiRespRec = lResp
	}
	// Store thre Response in Log table
	// lResponse, lErr4 := json.Marshal(lResp)
	// if lErr4 != nil {
	// 	log.Println("BIBP04", lErr4)
	// 	return lApiRespRec, lErr4
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	var lErr5 error
	// 	lId, lErr5 = Function.LogEntry(lLogInputRec)
	// 	if lErr5 != nil {
	// 		log.Println("BIBP05", lErr5)
	// 		return lApiRespRec, lErr5
	// 	}
	// }
	// log.Println("lApiRespRec", lApiRespRec)
	// }
	// }
	log.Println("BsePassword (-)")
	return lApiRespRec, nil
}

func BsePasswordApi(pUser PswdReqStruct, pUrl string, pToken string, pBrokerId int) (PswdRespStruct, error) {
	log.Println("BsePasswordApi (+)")
	//create instance for loginResp struct
	var lUserRespRec PswdRespStruct
	//create array of instance to store the key value pairs
	var lHeaderArr []apiUtil.HeaderDetails
	var lConsHeadRec apiUtil.HeaderDetails
	lDetail, lErr := AccessBseCredential(pBrokerId)
	if lErr != nil {
		log.Println("BIBPA01", lErr)
		return lUserRespRec, lErr
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

		//added content type
		lConsHeadRec.Key = "Content-Type"
		lConsHeadRec.Value = "application/json"
		lHeaderArr = append(lHeaderArr, lConsHeadRec)
		log.Println("headervalue", lHeaderArr)

		// Marshall the structure parameter into json
		ljsonData, lErr1 := json.Marshal(pUser)
		if lErr1 != nil {
			log.Println("BIBPA02", lErr1)
			return lUserRespRec, lErr1
		} else {
			// convert json data into string
			lReqstring := string(ljsonData)
			lResp, lErr2 := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "BsePassword")
			if lErr2 != nil {
				log.Println("BIBPA03", lErr2)
				return lUserRespRec, lErr2
			} else {
				// Unmarshalling json to struct
				lErr3 := json.Unmarshal([]byte(lResp), &lUserRespRec)
				if lErr3 != nil {
					log.Println("BIBPA04", lErr3)
					return lUserRespRec, lErr3
				}
			}
		}
	}
	log.Println("BsePasswordApi (-)")
	return lUserRespRec, nil
}
