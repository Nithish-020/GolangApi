package nseipo

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type IpoMktCatwiseStruct struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Quantity    int    `json:"quantity"`
	BidCount    int    `json:"bidCount"`
}

type IpoMktCatwiseRespStruct struct {
	Symbol string                `json:"symbol"`
	Isin   string                `json:"isin"`
	Status string                `json:"status"`
	Reason string                `json:"reason"`
	Data   []IpoMktCatwiseStruct `json:"demand"`
}

/*
Pupose: This method returns the ipo market demand list from the NSE Exchange
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e",
	"user":"AUTOBOT",
	"symbol":"JDIAL"
Response:
    *On Sucess
    =========
	{
		"symbol": "JDIAL",
		"status": "success",
		"demand": [
			{
			"category": "QIB",
			"subCategory": "IC",
			"quantity": 100,
			"bidCount": 1
			},
			{
			"category": "QIB",
			"subCategory": "FII",
			"quantity": 200,
			"bidCount": 2
			},
			{
			"category": "QIB",
			"subCategory": "MF",
			"quantity": 250,
			"bidCount": 3
			},
			{
			"category": "NIB",
			"subCategory": "MF",
			"quantity": 250,
			"bidCount": 3
			},
			{
			"category": "NIB",
			"subCategory": "CO",
			"quantity": 250,
			"bidCount": 3
			},
			{
			"category": "NIB",
			"subCategory": "IND",
			"quantity": 20,
			"bidCount": 1
			},
			{
			"category": "RETAIL",
			"subCategory": "IND",
			"quantity": 500,
			"bidCount": 5
			}
		]
	}
    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Nithish kumar
Date: 11DEC2023
*/
func IpoMktCatwise(pToken string, pUser string, pSymbol string) (IpoMktCatwiseRespStruct, error) {
	log.Println("IpoMktCatwise...(+)")
	//create parameters struct for LogEntry method
	var lLogInputRec Function.ParameterStruct
	//create instance to receive iporespStruct
	var lApiRespRec IpoMktCatwiseRespStruct
	// create instance to hold the last inserted id
	var lId int
	// create instance to hold errors
	var lErr error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["MarketCatwise"])
	lUrl = lUrl + pSymbol
	log.Println(lUrl, "endpoint")

	lLogInputRec.Request = ""
	lLogInputRec.EndPoint = "nse/mktdata/v1/catwise/"
	lLogInputRec.Flag = common.INSERT
	lLogInputRec.ClientId = pUser
	lLogInputRec.Method = "GET"

	// ! LogEntry method is used to store the Request in Database
	lId, lErr = Function.LogEntry(lLogInputRec)

	if lErr != nil {
		log.Println("NIMC01", lErr)
		return lApiRespRec, lErr
	} else {
		// TokenApi method used to call exchange API
		lResp, lErr := ExchangeMktCatwise(pToken, lUrl)
		if lErr != nil {
			log.Println("NIMC02", lErr)
			return lApiRespRec, lErr
		}
		lApiRespRec = lResp
		// Store thre Response in Log table
		lResponse, lErr := json.Marshal(lResp)
		if lErr != nil {
			log.Println("NIMC03", lErr)
			return lApiRespRec, lErr
		} else {
			lLogInputRec.Response = string(lResponse)
			lLogInputRec.LastId = lId
			lLogInputRec.Flag = common.UPDATE
			lId, lErr = Function.LogEntry(lLogInputRec)
			if lErr != nil {
				log.Println("NIMC04", lErr)
				return lApiRespRec, lErr
			}
		}
	}
	// log.Println(lApiRespRec)
	log.Println("IpoMktCatwise...(-)")
	return lApiRespRec, nil
}

func ExchangeMktCatwise(pToken string, pUrl string) (IpoMktCatwiseRespStruct, error) {
	log.Println("ExchangeMktCatwise....(+)")
	// create a new instance of IpoMktDemandRespStruct
	var lIpoRespRec IpoMktCatwiseRespStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lCredential := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Credential"])

	lConsHeadRec.Key = "Access-Token"
	lConsHeadRec.Value = pToken
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "Content-Type"
	lConsHeadRec.Value = "application/json"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "User-Agent"
	lConsHeadRec.Value = "Flattrade-golang"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "Authorization"
	lConsHeadRec.Value = "Basic " + lCredential
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lResp, lErr := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "FetchIpoMktCatwise")
	if lErr != nil {
		log.Println("NEMC01", lErr)
		return lIpoRespRec, lErr
	} else {
		lErr = json.Unmarshal([]byte(lResp), &lIpoRespRec)
		if lErr != nil {
			log.Println("NEMC02", lErr)
			return lIpoRespRec, lErr
		}
	}
	log.Println("ExchangeMktCatwise....(-)")
	return lIpoRespRec, nil
}
