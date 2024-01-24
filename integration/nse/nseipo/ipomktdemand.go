package nseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type IpoMktDemandStruct struct {
	CutOffIndicator    bool    `json:"cutOffIndicator"`
	Series             string  `json:"series"`
	Price              float32 `json:"price"`
	AbsoluteQuantity   int     `json:"absoluteQuantity"`
	CumulativeQuantity int     `json:"cumulativeQuantity"`
	AbsoluteBidCount   int     `json:"absoluteBidCount"`
	CumulativeBidCount int     `json:"cumulativeBidCount"`
}

type IpoMktDemandRespStruct struct {
	Symbol string               `json:"symbol"`
	Isin   string               `json:"isin"`
	Price  int                  `json:"price"`
	Status string               `json:"status"`
	Reason string               `json:"reason"`
	Data   []IpoMktDemandStruct `json:"demand"`
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
				"cutOffIndicator": true,
				"absoluteQuantity": 100,
				"cumulativeQuantity": 100,
				"absoluteBidCount": 1,
				"cumulativeBidCount": 1
				},
				{
				"cutOffIndicator": false,
				"price": 121.00,
				"absoluteQuantity": 200,
				"cumulativeQuantity": 300,
				"absoluteBidCount": 2,
				"cumulativeBidCount": 3
				},
				{
				"cutOffIndicator": false,
				"price": 120.00,
				"absoluteQuantity": 300,
				"cumulativeQuantity": 600,
				"absoluteBidCount": 3,
				"cumulativeBidCount": 6
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
func IpoMktDemand(pToken string, pUser string, pSymbol string) (IpoMktDemandRespStruct, error) {
	log.Println("IpoMktDemand...(+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive iporespStruct
	var lApiRespRec IpoMktDemandRespStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["MarketDemand"])
	lUrl = lUrl + pSymbol
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "nse/mktdata/v1/demand/"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = pUser
	// lLogInputRec.Method = "GET"

	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr = Function.LogEntry(lLogInputRec)

	// if lErr != nil {
	// 	log.Println("NIMD01", lErr)
	// 	return lApiRespRec, lErr
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr := ExchangeMktDemand(pToken, lUrl)
	if lErr != nil {
		log.Println("NIMD02", lErr)
		return lApiRespRec, lErr
	} else {
		lApiRespRec = lResp
	}
	// Store thre Response in Log table
	// lResponse, lErr := json.Marshal(lResp)
	// if lErr != nil {
	// 	log.Println("NIMD03", lErr)
	// 	return lApiRespRec, lErr
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NIMD04", lErr)
	// 		return lApiRespRec, lErr
	// 	}
	// }
	// }
	// log.Println(lApiRespRec)
	log.Println("IpoMktDemand...(-)")
	return lApiRespRec, nil
}

func ExchangeMktDemand(pToken string, pUrl string) (IpoMktDemandRespStruct, error) {
	log.Println("ExchangeMktDemand....(+)")
	// create a new instance of IpoMktDemandRespStruct
	var lIpoRespRec IpoMktDemandRespStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lCredential := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Credential"])

	// lConsHeadRec.Key = "Access-Token"
	// lConsHeadRec.Value = pToken
	// lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "Content-Type"
	lConsHeadRec.Value = "application/json"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "User-Agent"
	lConsHeadRec.Value = "Flattrade-golang"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lConsHeadRec.Key = "Authorization"
	lConsHeadRec.Value = "Basic " + lCredential
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	lResp, lErr := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "FetchIpoMktDemand")
	if lErr != nil {
		log.Println("NEMD01", lErr)
		return lIpoRespRec, lErr
	} else {
		lErr = json.Unmarshal([]byte(lResp), &lIpoRespRec)
		if lErr != nil {
			log.Println("NEMD02", lErr)
			return lIpoRespRec, lErr
		}
	}
	log.Println("ExchangeMktDemand....(-)")
	return lIpoRespRec, nil
}
