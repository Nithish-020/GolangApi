package nsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
	"time"
)

type SgbTransactionResp struct {
	Status       string
	Reason       string
	Transactions []SgbAddResStruct
}

/*
Pupose: This method  allows user to download SGB transactions.
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e"
Response:
    *On Sucess
    =========

    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: KAVYA DHARSHANI
Date: 28SEP2023
*/

func SgbTransactionsMaster(pToken string, pDate string, lUser string) (SgbTransactionResp, error) {
	log.Println("SgbTransactionsMaster...(+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbTransactionResp
	var lApiRespRec SgbTransactionResp
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["FetchSgbTranscation"])

	log.Println("pDate", pDate)
	// Parse the input date
	inputDate, err := time.Parse("2006-01-02", pDate)
	if err != nil {
		log.Println("Error parsing date:", err)
		return lApiRespRec, err
	}

	// Format the date in "01-02-2006" format (month-day-year)
	// pDate = inputDate.Format("01-02-2006")
	pDate = inputDate.Format("02-01-2006")
	time := "%2009:30:00"

	lUrl = lUrl + pDate + time
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "/v1/sgb/" + pDate + time
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = lUser
	// lLogInputRec.Method = "GET"

	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr1 = Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("NSTM01", lErr1)
	// 	return lApiRespRec, lErr1
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := ExchangeSgbTransactionsMaster(pToken, lUrl)
	if lErr2 != nil {
		log.Println("NSTM02", lErr2)
		return lApiRespRec, lErr2
	} else {
		lApiRespRec = lResp
		log.Println("Response", lResp)
	}
	// Store thre Response in Log table
	// 	lResponse, lErr3 := json.Marshal(lResp)
	// 	if lErr3 != nil {
	// 		log.Println("NSTM03", lErr3)
	// 		return lApiRespRec, lErr3
	// 	} else {
	// 		lLogInputRec.Response = string(lResponse)
	// 		lLogInputRec.LastId = lId
	// 		lLogInputRec.Flag = common.UPDATE
	// 		// create instance to hold errors
	// 		var lErr4 error
	// 		lId, lErr4 = Function.LogEntry(lLogInputRec)
	// 		if lErr4 != nil {
	// 			log.Println("NSTM04", lErr4)
	// 			return lApiRespRec, lErr4
	// 		}
	// 	}
	// }
	log.Println("SgbTransactionsMaster...(-)")
	return lApiRespRec, nil
}
func ExchangeSgbTransactionsMaster(pToken string, pUrl string) (SgbTransactionResp, error) {
	log.Println("ExchangeSgbTransactionsMaster (+)")
	// create a new instance of IpoResponseStruct
	var lSgbRespRec SgbTransactionResp
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

	lResp, lErr1 := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "NSESGBDownloadOrders")
	if lErr1 != nil {
		log.Println("NESTM01", lErr1)
		return lSgbRespRec, lErr1
	} else {
		log.Println("Response JSon", string(lResp))
		lErr2 := json.Unmarshal([]byte(lResp), &lSgbRespRec)
		if lErr2 != nil {
			log.Println("NESTM02", lErr2)
			return lSgbRespRec, lErr2
		}
	}
	log.Println("ExchangeSgbTransactionsMaster (-)")
	return lSgbRespRec, nil
}
