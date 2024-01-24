package bseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type BseExchangeReqStruct struct {
	ScripId            string                `json:"scripid"`
	ApplicationNo      string                `json:"applicationno"`
	Category           string                `json:"category"`
	ApplicantName      string                `json:"applicantname"`
	Depository         string                `json:"depository"`
	DpId               string                `json:"dpid"`
	ClientBenfId       string                `json:"clientbenfid"`
	ChequeReceivedFlag string                `json:"chequereceivedflag"`
	ChequeAmount       string                `json:"chequeamount"`
	PanNo              string                `json:"panno"`
	BankName           string                `json:"bankname"`
	Location           string                `json:"location"`
	AccountNo_UpiId    string                `json:"accountnumber_upiid"`
	IfscCode           string                `json:"ifsccode"`
	ReferenceNo        string                `json:"referenceno"`
	Asba_UpiId         string                `json:"asba_upiid"`
	Bids               []BseRequestBidStruct `json:"bids"`
}

type BseRequestBidStruct struct {
	BidId       string `json:"bidid"`
	Quantity    string `json:"quantity"`
	Rate        string `json:"rate"`
	CuttOffFlag string `json:"cuttoffflag"`
	OrderNo     string `json:"orderno"`
	ActionCode  string `json:"actioncode"`
}

type BseExchangeRespStruct struct {
	ScripId            string                 `json:"scripid"`
	ApplicationNo      string                 `json:"applicationno"`
	Category           string                 `json:"category"`
	ApplicantName      string                 `json:"applicantname"`
	Depository         string                 `json:"depository"`
	DpId               string                 `json:"dpid"`
	ClientBenfId       string                 `json:"clientbenfid"`
	ChequeReceivedFlag string                 `json:"chequereceivedflag"`
	ChequeAmount       string                 `json:"chequeamount"`
	PanNo              string                 `json:"panno"`
	BankName           string                 `json:"bankname"`
	Location           string                 `json:"location"`
	AccountNo_UpiId    string                 `json:"accountnumber_upiid"`
	IfscCode           string                 `json:"ifsccode"`
	ReferenceNo        string                 `json:"referenceno"`
	Asba_UpiId         string                 `json:"asba_upiid"`
	StatusCode         string                 `json:"statuscode"`
	StatusMessage      string                 `json:"statusmessage"`
	ErrorCode          string                 `json:"errorcode"`
	ErrorMessage       string                 `json:"errormessage"`
	Bids               []BseResponseBidStruct `json:"bids"`
}

type BseResponseBidStruct struct {
	BidId       string `json:"bidid"`
	Quantity    string `json:"quantity"`
	Rate        string `json:"rate"`
	CuttOffFlag string `json:"cuttoffflag"`
	OrderNo     string `json:"orderno"`
	ActionCode  string `json:"actioncode"`
	ErrorCode   string `json:"errorcode"`
	Message     string `json:"message"`
}

func BseIpoOrder(pToken string, pRequestArr BseExchangeReqStruct, pUser string, pBrokerId int) (BseExchangeRespStruct, error) {
	log.Println("BseIpoOrder (+)")

	//For Exchagnge response
	var lExchangeResArr BseExchangeRespStruct
	// var lLogInputRec Function.ParameterStruct
	// create instance to hold the last inserted id
	// var lId int
	// To establish the connection between toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseIpoOrder"])
	log.Println(lUrl, "endpoint")

	// lRequest, lErr := json.Marshal(pRequestArr)
	// if lErr != nil {
	// 	log.Println("BBIO01", lErr)
	// 	return lExchangeResArr, lErr
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "bse/v1/ipoorder"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// LogEntry method is used to store the Request in Database
	// lId, lErr = Function.LogEntry(lLogInputRec)
	// if lErr != nil {
	// 	log.Println("BBIO02", lErr)
	// 	return lExchangeResArr, lErr
	// } else {
	// ExchangeOrder method used to call exchange API
	lResp, lErr2 := BseExchangeOrder(pToken, lUrl, pRequestArr, pBrokerId)
	if lErr2 != nil {
		log.Println("BBIO03", lErr2)
		return lExchangeResArr, lErr2
	} else {
		lExchangeResArr = lResp
	}
	// Store thre Response in Log table
	// lResponse, lErr3 := json.Marshal(lResp)
	// if lErr3 != nil {
	// 	log.Println("BBIO04", lErr3)
	// 	return lExchangeResArr, lErr3
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	lId, lErr3 = Function.LogEntry(lLogInputRec)
	// 	if lErr3 != nil {
	// 		log.Println("BBIO05", lErr3)
	// 		return lExchangeResArr, lErr3
	// 	}
	// }
	// }
	// }
	log.Println(lExchangeResArr)
	log.Println("BseIpoOrder (-)")
	return lExchangeResArr, nil
}

func BseExchangeOrder(pToken string, pUrl string, pExchangeReqRec BseExchangeReqStruct, pBrokerId int) (BseExchangeRespStruct, error) {
	log.Println("BseExchangeOrder (+)")

	var lRespArr BseExchangeRespStruct
	//for constructing Header details
	var lConsHeadRec apiUtil.HeaderDetails
	var lHeaderArr []apiUtil.HeaderDetails

	lDetail, lErr := AccessBseCredential(pBrokerId)
	if lErr != nil {
		log.Println("NESM01", lErr)
		return lRespArr, lErr
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

		lConsHeadRec.Key = "Content-Type"
		lConsHeadRec.Value = "application/json"
		lHeaderArr = append(lHeaderArr, lConsHeadRec)

		log.Println("lHeaderArr", lHeaderArr)

		ljsonData, lErr1 := json.Marshal(pExchangeReqRec)
		if lErr1 != nil {
			log.Println("BBEO01", lErr1)
			return lRespArr, lErr1
		} else {
			lReqstring := string(ljsonData)
			lResp, lErr2 := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "BseIpoOrder")
			if lErr2 != nil {
				log.Println("BBEO02", lErr2)
				return lRespArr, lErr2
			} else {
				log.Println("resp", lResp)

				// lRespArr = lResp
				lErr2 = json.Unmarshal([]byte(lResp), &lRespArr)
				if lErr2 != nil {
					log.Println("BBEO03", lErr2, lResp)
					return lRespArr, lErr2
				} else {
					log.Println("ExchangeResponse---->>>>", lRespArr)
				}
			}
		}
	}
	log.Println("BseExchangeOrder (-)")
	return lRespArr, nil
}
