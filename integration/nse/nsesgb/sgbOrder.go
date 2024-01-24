package nsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type SgbAddReqStruct struct {
	Symbol            string  `json:"symbol"`
	OrderNumber       int     `json:"orderNumber"`
	Price             float32 `json:"price"`
	Quantity          int     `json:"quantity"`
	PhysicalDematFlag string  `json:"physicalDematFlag"`
	ClientCode        string  `json:"clientCode"`
	Pan               string  `json:"pan"`
	Depository        string  `json:"depository"`
	DpId              string  `json:"dpId"`
	ClientBenId       string  `json:"clientBenId"`
	ClientRefNumber   string  `json:"clientRefNumber"`
	ActivityType      string  `json:"activityType"`
}

type SgbAddResStruct struct {
	Symbol             string  `json:"symbol"`
	OrderNumber        int     `json:"orderNumber"`
	Series             string  `json:"series"`
	ApplicationNumber  string  `json:"applicationNumber"`
	Quantity           int     `json:"quantity"`
	Price              float32 `json:"price"`
	PhysicalDematFlag  string  `json:"physicalDematFlag"`
	ClientCode         string  `json:"clientCode"`
	Pan                string  `json:"pan"`
	Depository         string  `json:"depository"`
	DpId               string  `json:"dpId"`
	ClientBenId        string  `json:"clientBenId"`
	ClientRefNumber    string  `json:"clientRefNumber"`
	OrderStatus        string  `json:"orderStatus"`
	RejectionReason    string  `json:"rejectionReason"`
	EnteredBy          string  `json:"enteredBy"`
	EntryTime          string  `json:"entryTime"`
	VerificationStatus string  `json:"verificationStatus"`
	VerificationReason string  `json:"verificationReason"`
	ClearingStatus     string  `json:"clearingStatus"`
	ClearingReason     string  `json:"clearingReason"`
	LastActionTime     string  `json:"lastActionTime"`
	Status             string  `json:"status"`
	Reason             string  `json:"reason"`
}

func SgbOrderTransaction(pToken string, pSgbRequestArr SgbAddReqStruct, pUser string) (SgbAddResStruct, error) {
	log.Println("SgbOrderTransaction....(+)")

	//For Exchagnge response
	var lsgbExhangeResArr SgbAddResStruct
	// var lLogInputRec Function.ParameterStruct
	// create instance to hold the last inserted id
	// var lId int
	// To establish the connection between toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NseSgbOrder"])
	log.Println(lUrl, "endpoint")

	// lRequest, lErr := json.Marshal(pSgbRequestArr)
	// if lErr != nil {
	// 	log.Println("NSOT01", lErr)
	// 	return lsgbExhangeResArr, lErr
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "/v1/sgb/add"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// 	// LogEntry method is used to store the Request in Database
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NSOT02", lErr)
	// 		return lsgbExhangeResArr, lErr
	// 	} else {
	// ExchangeOrder method used to call exchange API
	lResp, lErr := SgbExchangeOrder(pToken, lUrl, pSgbRequestArr)
	if lErr != nil {
		log.Println("NSOT03", lErr)
		return lsgbExhangeResArr, lErr
	} else {
		lsgbExhangeResArr = lResp
	}
	// Store thre Response in Log table
	// lResponse, lErr := json.Marshal(lResp)
	// if lErr != nil {
	// 	log.Println("NSOT04", lErr)
	// 	return lsgbExhangeResArr, lErr
	// } else {
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NSOT05", lErr)
	// 		return lsgbExhangeResArr, lErr
	// 	}
	// }
	// }
	// }
	log.Println(lsgbExhangeResArr)
	log.Println("SgbOrderTransaction....(-)")
	return lsgbExhangeResArr, nil
}

func SgbExchangeOrder(pToken string, pUrl string, pSgbExchangeReqArr SgbAddReqStruct) (SgbAddResStruct, error) {
	log.Println("SgbExchangeOrder....(+)")

	var lSgbRespArr SgbAddResStruct
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

	ljsonData, lErr := json.Marshal(pSgbExchangeReqArr)
	if lErr != nil {
		log.Println("NSEO01", lErr)
		return lSgbRespArr, lErr
	} else {
		lReqstring := string(ljsonData)
		lSgbResp, lErr := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "Token")
		if lErr != nil {
			log.Println("NSEO02", lErr)
			return lSgbRespArr, lErr
		} else {
			// lSgbRespArr = lSgbResp
			lErr = json.Unmarshal([]byte(lSgbResp), &lSgbRespArr)
			if lErr != nil {
				log.Println("NSEO03", lErr, lSgbResp)
				return lSgbRespArr, lErr
			} else {
				log.Println("SgbExchangeResponse---->>>>", lSgbRespArr)
			}
		}
	}
	log.Println("SgbExchangeOrder....(-)")
	return lSgbRespArr, nil
}
