package nseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type ExchangeReqStruct struct {
	Symbol        string             `json:"symbol"`
	ApplicationNo string             `json:"applicationNumber"`
	Category      string             `json:"category"`
	ClientName    string             `json:"clientName"`
	Depository    string             `json:"depository"`
	DpId          string             `json:"dpId"`
	ClientBenId   string             `json:"clientBenId"`
	NonASBA       bool               `json:"nonASBA"`
	Pan           string             `json:"pan"`
	ReferenceNo   string             `json:"referenceNumber"`
	AllotmentMode string             `json:"allotmentMode"`
	UpiFlag       string             `json:"upiFlag"`
	Upi           string             `json:"upi"`
	BankCode      string             `json:"bankCode"`
	LocationCode  string             `json:"locationCode"`
	BankAccount   string             `json:"bankAccount"`
	IFSC          string             `json:"ifsc"`
	SubBrokerCode string             `json:"subBrokerCode"`
	TimeStamp     string             `json:"timestamp"`
	Bids          []RequestBidStruct `json:"bids"`
}

type RequestBidStruct struct {
	ActivityType   string `json:"activityType"`
	BidReferenceNo int    `json:"bidReferenceNumber"`
	Series         string `json:"series"`
	Quantity       int    `json:"quantity"`
	AtCutOff       bool   `json:"atCutOff"`
	Price          int    `json:"price"`
	Amount         int    `json:"amount"`
	Remark         string `json:"remark"`
	LotSize        int    `json:"lotSize"` //------------
}

type ExchangeRespStruct struct {
	Symbol               string              `json:"symbol"`
	ApplicationNo        string              `json:"applicationNumber"`
	Category             string              `json:"category"`
	ClientName           string              `json:"clientName"`
	Depository           string              `json:"depository"`
	DpId                 string              `json:"dpId"`
	ClientBenId          string              `json:"clientBenId"`
	NonASBA              bool                `json:"nonASBA"`
	ChequeNo             string              `json:"chequeNumber"`
	Pan                  string              `json:"pan"`
	ReferenceNo          string              `json:"referenceNumber"`
	AllotmentMode        string              `json:"allotmentMode"`
	UpiFlag              string              `json:"upiFlag"`
	Upi                  string              `json:"upi"`
	BankCode             string              `json:"bankCode"`
	LocationCode         string              `json:"locationCode"`
	BankAccount          string              `json:"bankAccount"`
	IFSC                 string              `json:"ifsc"`
	SubBrokerCode        string              `json:"subBrokerCode"`
	TimeStamp            string              `json:"timeStamp"`
	DpVerStatusFlag      string              `json:"dpVerStatusFlag"`
	DpVerFailCode        string              `json:"dpVerFailCode"`
	DpVerReason          string              `json:"dpVerReason"`
	UpiPaymentStatusFlag int                 `json:"upiPaymentStatusFlag"`
	UpiAmtBlocked        int                 `json:"upiAmtBlocked"`
	Status               string              `json:"status"`
	ReasonCode           int                 `json:"reasonCode"`
	Reason               string              `json:"reason"`
	Bids                 []ResponseBidStruct `json:"bids"`
}

type ResponseBidStruct struct {
	ActivityType   string  `json:"activityType"`
	BidReferenceNo int     `json:"bidReferenceNumber"`
	Series         string  `json:"series"`
	Quantity       int     `json:"quantity"`
	AtCutOff       bool    `json:"atCutOff"`
	Price          float32 `json:"price"`
	Amount         float64 `json:"amount"`
	Remark         string  `json:"remark"`
	Status         string  `json:"status"`
	ReasonCode     int     `json:"reasonCode"`
	Reason         string  `json:"reason"`
}

func Transaction(pToken string, pRequestArr []ExchangeReqStruct, pUser string) ([]ExchangeRespStruct, error) {
	log.Println("Transaction....(+)")

	//For Exchagnge response
	var lExhangeResArr []ExchangeRespStruct
	// var lLogInputRec Function.ParameterStruct
	// create instance to hold the last inserted id
	// var lId int
	// To establish the connection between toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Transcation"])
	log.Println(lUrl, "endpoint")

	// lRequest, lErr := json.Marshal(pRequestArr)
	// if lErr != nil {
	// 	log.Println("NT01", lErr)
	// 	return lExhangeResArr, lErr
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "/v1/transcations/addbulk"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// 	// LogEntry method is used to store the Request in Database
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NT02", lErr)
	// 		return lExhangeResArr, lErr
	// 	} else {
	// ExchangeOrder method used to call exchange API
	lResp, lErr := ExchangeOrder(pToken, lUrl, pRequestArr)
	if lErr != nil {
		log.Println("NT03", lErr)
		return lExhangeResArr, lErr
	} else {
		lExhangeResArr = lResp
	}
	// Store thre Response in Log table
	// lResponse, lErr := json.Marshal(lResp)
	// if lErr != nil {
	// 	log.Println("NT04", lErr)
	// 	return lExhangeResArr, lErr
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	_, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NT05", lErr)
	// 		return lExhangeResArr, lErr
	// 	}
	// }
	// }
	// }
	// log.Println(lExhangeResArr)
	log.Println("Transaction....(-)")
	return lExhangeResArr, nil
}

func ExchangeOrder(pToken string, pUrl string, pExchangeReqArr []ExchangeReqStruct) ([]ExchangeRespStruct, error) {
	log.Println("ExchangeOrder....(+)")

	var lRespArr []ExchangeRespStruct
	//for constructing Header details
	var lConstructHead apiUtil.HeaderDetails
	var lHeaderArr []apiUtil.HeaderDetails

	lConstructHead.Key = "Access-Token"
	lConstructHead.Value = pToken
	lHeaderArr = append(lHeaderArr, lConstructHead)

	lConstructHead.Key = "Content-Type"
	lConstructHead.Value = "application/json"
	lHeaderArr = append(lHeaderArr, lConstructHead)

	lConstructHead.Key = "User-Agent"
	lConstructHead.Value = "Flattrade-golang"
	lHeaderArr = append(lHeaderArr, lConstructHead)

	ljsonData, lErr := json.Marshal(pExchangeReqArr)
	if lErr != nil {
		log.Println("NEO01", lErr)
		return lRespArr, lErr
	} else {
		lReqstring := string(ljsonData)
		lResp, lErr := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "Token")
		if lErr != nil {
			log.Println("NEO02", lErr)
			return lRespArr, lErr
		} else {
			// lRespArr = lResp
			lErr = json.Unmarshal([]byte(lResp), &lRespArr)
			if lErr != nil {
				log.Println("NEO03", lErr, lResp)
				return lRespArr, lErr
			} else {
				log.Println("ExchangeResponse---->>>>", lRespArr)
			}
		}
	}
	log.Println("ExchangeOrder....(-)")
	return lRespArr, nil
}
