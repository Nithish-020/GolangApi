package nsencb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
	"strconv"
)

type NcbDetailStruct struct {
	Symbol           string  `json:"symbol"`
	Series           string  `json:"series"`
	Name             string  `json:"name"`
	LotSize          int     `json:"lotSize"`
	FaceValue        float64 `json:"faceValue"`
	MinBidQuantity   int     `json:"minBidQuantity"`
	MinPrice         float64 `json:"minPrice"`
	MaxPrice         float64 `json:"maxPrice"`
	TickSize         float64 `json:"tickSize"`
	CutoffPrice      float64 `json:"cutoffPrice"`
	BiddingStartDate string  `json:"biddingStartDate"`
	BiddingEndDate   string  `json:"biddingEndDate"`
	DailyStartTime   string  `json:"dailyStartTime"`
	DailyEndTime     string  `json:"dailyEndTime"`
	T1ModStartDate   string  `json:"t1ModStartDate"`
	T1ModEndDate     string  `json:"t1ModEndDate"`
	T1ModStartTime   string  `json:"t1ModStartTime"`
	T1ModEndTime     string  `json:"t1ModEndTime"`
	Isin             string  `json:"isin"`
	IssueSize        int64   `json:"issueSize"`
	IssueValueSize   float64 `json:"issueValueSize"`
	// MaxQuantity           string  `json:"maxQuantity"`
	MaxQuantity           interface{} `json:"maxQuantity"`
	AllotmentDate         string      `json:"allotmentDate"`
	LastDayBiddingEndTime string      `json:"lastDayBiddingEndTime"`
}

type NcbRespStruct struct {
	Data   []NcbDetailStruct `json:"data"`
	Reason string            `json:"reason"`
	Status string            `json:"status"`
}

/*
Pupose: This method returns the active NCB list from the NSE Exchange
Parameters:
    "token":"1a5beb81-6e84-4efa-9c5e-3756869c482e"
Response:
    *On Sucess
    =========

    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Kavya Dharshani
Date: 03OCT2023
*/
func NcbMaster(pToken string, lUser string) (NcbRespStruct, error) {
	log.Println("NcbMaster...(+)")
	log.Println("pToken : NcbMaster", pToken)
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct

	//create instance to receive NcbRespStruct
	var lNcbRespRec NcbRespStruct
	// create instance to hold the last inserted id

	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NcbMaster"])
	// log.Println(lUrl, "endpoint")

	// lLogInputRec.Request = ""
	// lLogInputRec.EndPoint = "/v1/ncbmaster"
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = lUser
	// lLogInputRec.Method = "GET"

	// // ! LogEntry method is used to store the Request in Database
	// lId, lErr1 = Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("NNM01", lErr1)
	// 	return lNcbRespRec, lErr1
	// } else {
	// TokenApi method used to call exchange API
	lNcbResp, lErr2 := ExchangeNcbMaster(pToken, lUrl)
	if lErr2 != nil {
		log.Println("NNM02", lErr2)
		return lNcbRespRec, lErr2
	} else {
		lNcbRespRec = lNcbResp
		// log.Println("lNcbResp Response", lNcbResp)
	}
	// Store thre Response in Log table
	// lResponse, lErr3 := json.Marshal(lNcbResp)
	// if lErr3 != nil {
	// 	log.Println("NNM03", lErr3)
	// 	return lNcbRespRec, lErr3
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	// create instance to hold errors
	// 	var lErr4 error
	// 	lId, lErr4 = Function.LogEntry(lLogInputRec)
	// 	if lErr4 != nil {
	// 		log.Println("NNM04", lErr4)
	// 		return lNcbRespRec, lErr4
	// 	}
	// }
	// }

	log.Println("NcbMaster...(-)")
	return lNcbRespRec, nil
}

func ExchangeNcbMaster(pToken string, pUrl string) (NcbRespStruct, error) {
	log.Println("ExchangeNcbMaster....(+)")
	// log.Println("pToken :ExchangeNcbMaster ", pToken)
	// create a new instance of NcbResponseStruct
	var lNcbRespRec NcbRespStruct
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

	lResp, lErr1 := apiUtil.Api_call(pUrl, "GET", "", lHeaderArr, "NcbMaster")
	if lErr1 != nil {
		log.Println("NENM01", lErr1)
		return lNcbRespRec, lErr1
	} else {
		// log.Println("lResp API_CAll NCb", lResp)
		lErr2 := json.Unmarshal([]byte(lResp), &lNcbRespRec)
		if lErr2 != nil {
			log.Println("NENM02", lErr2)
			return lNcbRespRec, lErr2
		} else {
			log.Println("lNcbRespRec", lNcbRespRec)
		}

		// for Idx := 0; Idx < len(lNcbRespRec.Data); Idx++ {
		// 	maxQuantityStr := lNcbRespRec.Data[Idx].MaxQuantity
		// 	if maxQuantityStr == "" {
		// 		maxQuantityStr = "0"
		// 	} else {
		// 		log.Println("maxQuantityStr", maxQuantityStr)
		// 		maxQuantity, lErr3 := strconv.ParseFloat(maxQuantityStr, 64)
		// 		if lErr3 != nil {
		// 			log.Println("Error parsing MaxQuantity:", lErr3)
		// 			//return lErr3
		// 		} else {
		// 			log.Println("maxQuantity", maxQuantity)
		// 		}
		// 	}
		// }

		for Idx := 0; Idx < len(lNcbRespRec.Data); Idx++ {
			maxQuantityValue := lNcbRespRec.Data[Idx].MaxQuantity

			// Determine the data type and value
			switch maxQuantity := maxQuantityValue.(type) {
			case string:
				// Handle string value
				// log.Println("maxQuantity (string):", maxQuantity)
				// You can convert it to a numeric type if needed
				maxQuantityFloat, err := strconv.ParseFloat(maxQuantity, 64)
				if err != nil {
					log.Println("Error parsing MaxQuantity:", err)
				} else {
					lNcbRespRec.Data[Idx].MaxQuantity = maxQuantityFloat
					// log.Println("maxQuantity (float64):", maxQuantityFloat)
				}

			case float64:
				// Handle float64 value
				// log.Println("maxQuantity (float64):", maxQuantity)
				lNcbRespRec.Data[Idx].MaxQuantity = maxQuantity

			default:
				log.Println("Unknown type for maxQuantity")
			}
		}

	}
	log.Println("ExchangeNcbMaster....(-)")
	return lNcbRespRec, nil
}
