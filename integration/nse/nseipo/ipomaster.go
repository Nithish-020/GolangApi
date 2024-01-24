package nseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type IpoDetailsStruct struct {
	Symbol                string                      `json:"symbol"`
	BiddingStartDate      string                      `json:"biddingStartDate"`
	BiddingEndDate        string                      `json:"biddingEndDate"`
	Registrar             string                      `json:"registrar"`
	T1ModStartDate        string                      `json:"t1ModStartDate"`
	T1ModEndDate          string                      `json:"t1ModEndDate"`
	T1ModStartTime        string                      `json:"t1ModStartTime"`
	T1ModEndTime          string                      `json:"t1ModEndTime"`
	DailyStartTime        string                      `json:"dailyStartTime"`
	DailyEndTime          string                      `json:"dailyEndTime"`
	MaxPrice              float32                     `json:"maxPrice"`
	MinPrice              float32                     `json:"minPrice"`
	MinBidQuantity        int                         `json:"minBidQuantity"`
	LotSize               int                         `json:"lotSize"`
	TickSize              float32                     `json:"tickSize"`
	FaceValue             float32                     `json:"faceValue"`
	Name                  string                      `json:"name"`
	IssueSize             int                         `json:"issueSize"`
	CutOffPrice           float32                     `json:"cutOffPrice"`
	ISIN                  string                      `json:"isin"`
	IssueType             string                      `json:"issueType"`
	SubType               string                      `json:"subType"`
	CategoryDetailsArr    []CategoryDetailsStruct     `json:"categoryDetails"`
	SubCategoryDetailsArr []SubCategorySettingsStruct `json:"subCategorySettings"`
	SeriesDetailsArr      []SeriesDetailsStruct       `json:"seriesDetails"`
}

type CategoryDetailsStruct struct {
	Code      string `json:"code"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type SubCategorySettingsStruct struct {
	SubCatCode    string      `json:"subCatCode"`
	MinValue      float32     `json:"minValue"`
	MaxUpiLimit   float32     `json:"maxUpiLimit"`
	AllowCutOff   bool        `json:"allowCutOff"`
	AllowUpi      bool        `json:"allowUpi"`
	MaxValue      float32     `json:"maxValue"`
	DiscountPrice float32     `json:"discountPrice"`
	MaxQuantity   interface{} `json:"maxQuantity"`
	MaxPrice      float32     `json:"maxPrice"`
	DiscountType  string      `json:"discountType"`
	CaCode        string      `json:"caCode"`
}

type SeriesDetailsStruct struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type IpoResponseStruct struct {
	Data []IpoDetailsStruct `json:"data"`
}

/*
Pupose: This method returns the active ipo list from the NSE Exchange
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e"
Response:
    *On Sucess
    =========
	{
            "biddingStartDate": "02-06-2023",
            "symbol": "DB2601",
            "minBidQuantity": 10,
            "registrar": "ALANKIT",
            "lotSize": 1,
            "t1ModEndDate": "",
            "dailyStartTime": "09:00:00",
            "t1ModStartTime": "08:00:00",
            "categoryDetails": [
                {
                    "code": "I",
                    "startTime": "09:00:00",
                    "endTime": "18:00:00"
                },
            ],
            "biddingEndDate": "30-06-2023",
            "t1ModEndTime": "23:00:00",
            "seriesDetails": [
                {
                    "code": "S9",
                    "desc": "gfytut"
                },
                {
                    "code": "S2",
                    "desc": "Tenor 24 months, Payable at Maturity, Rs. 1,180.75 for Category I & II Investors, Tenor 24 months, Payable at Maturity, Rs. 1,189.47 for Category III & IV Investors*"
                },
                {
                    "code": "S1",
                    "desc": "8.65 %, Tenor 24 months, for Category I & II Investors Annual Payment, 9.05 %, Tenor 24 months, for Category III & IV Investors Annual Payment*"
                },
            ],
            "subCategorySettings": [

                {
                    "subCatCode": "14",
                    "minValue": 10000.0,
                    "maxUpiLimit": null,
                    "allowCutOff": false,
                    "allowUpi": false,
                    "maxValue": 2.44E9,
                    "discountPrice": null,
                    "discountType": "",
                    "maxPrice": null,
                    "caCode": "I"
                },
                {
                    "subCatCode": "15",
                    "minValue": 10000.0,
                    "maxUpiLimit": null,
                    "allowCutOff": false,
                    "allowUpi": false,
                    "maxValue": 2.44E9,
                    "discountPrice": null,
                    "discountType": "",
                    "maxPrice": null,
                    "caCode": "II"
                }
            ],
            "dailyEndTime": "22:00:00",
            "tickSize": 1.0,
            "issueType": "DEBT",
            "faceValue": 10.0,
            "minPrice": 1000.0,
            "t1ModStartDate": "",
            "name": "DB2601",
            "issueSize": 50000000,
            "maxPrice": 2000.0,
            "cutOffPrice": 2000.0,
            "isin": "INE686Y01026"
        },

    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 05JUNE2023
*/
func IpoMaster(pToken string, lUser string) (IpoResponseStruct, error) {
	log.Println("IpoMaster...(+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive iporespStruct
	var lApiRespRec IpoResponseStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Master"])
	log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "nse/v1/ipomaster"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = lUser
	// lLogInputRec.Method = "GET"

	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr = Function.LogEntry(lLogInputRec)

	// if lErr != nil {
	// 	log.Println("NIM01", lErr)
	// 	return lApiRespRec, lErr
	// } else {
	// TokenApi method used to call exchange API
	lResp, lErr := ExchangeIpoMaster(pToken, lUrl)
	if lErr != nil {
		log.Println("NIM02", lErr)
		return lApiRespRec, lErr
	} else {
		lApiRespRec = lResp
	}
	// Store thre Response in Log table
	// lResponse, lErr := json.Marshal(lResp)
	// if lErr != nil {
	// 	log.Println("NIM03", lErr)
	// 	return lApiRespRec, lErr
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NIM04", lErr)
	// 		return lApiRespRec, lErr
	// 	}
	// }
	// }
	// log.Println(lApiRespRec)
	log.Println("IpoMaster...(-)")
	return lApiRespRec, nil
}

func ExchangeIpoMaster(pToken string, pUrl string) (IpoResponseStruct, error) {
	log.Println("ExchangeIpoMaster....(+)")
	// create a new instance of IpoResponseStruct
	var lIpoRespRec IpoResponseStruct
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

	lResp, lErr := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "FetchMaster")
	if lErr != nil {
		log.Println("NEIM01", lErr)
		return lIpoRespRec, lErr
	} else {
		lErr = json.Unmarshal([]byte(lResp), &lIpoRespRec)
		if lErr != nil {
			log.Println("NEIM02", lErr)
			return lIpoRespRec, lErr
		}
	}
	log.Println("ExchangeIpoMaster....(-)")
	return lIpoRespRec, nil
}
