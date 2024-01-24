package nsencb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type NcbAddReqStruct struct {
	Symbol            string  `json:"symbol"`
	InvestmentValue   int     `json:"investmentValue"`
	ApplicationNumber string  `json:"applicationNumber"`
	OrderNumber       int     `json:"orderNumber"`
	Price             float64 `json:"price"`
	PhysicalDematFlag string  `json:"physicalDematFlag"`
	Pan               string  `json:"pan"`
	Depository        string  `json:"depository"`
	DpId              string  `json:"dpId"`
	ClientBenId       string  `json:"clientBenId"`
	ClientRefNumber   string  `json:"clientRefNumber"`
	ActivityType      string  `json:"activityType"`
}

type NcbAddResStruct struct {
	Symbol             string  `json:"symbol"`
	OrderNumber        int     `json:"orderNumber"`
	Series             string  `json:"series"`
	ApplicationNumber  string  `json:"applicationNumber"`
	InvestmentValue    int     `json:"investmentValue"`
	Price              float64 `json:"price"`
	TotalAmountPayable float64 `json:"totalAmountPayable"`
	PhysicalDematFlag  string  `json:"physicalDematFlag"`
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

func NcbAddOrder(pToken string, pNcbRequestArr NcbAddReqStruct, pUser string) (NcbAddResStruct, error) {
	log.Println("NcbAddOrder....(+)")

	//For Exchagnge response
	var lNcbExhangeResArr NcbAddResStruct
	// var lLogInputRec Function.ParameterStruct
	// create instance to hold the last inserted id
	// var lId int
	// To establish the connection between toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NseNcbOrder"])
	log.Println(lUrl, "endpoint")

	// lRequest, lErr := json.Marshal(pNcbRequestArr)
	// if lErr != nil {
	// 	log.Println("NNAO01", lErr)
	// 	return lNcbExhangeResArr, lErr
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "/v1/ncb/add"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// 	// LogEntry method is used to store the Request in Database
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NNAO02", lErr)
	// 		return lNcbExhangeResArr, lErr
	// 	} else {
	// ExchangeOrder method used to call exchange API
	lResp, lErr := NcbExchangeOrder(pToken, lUrl, pNcbRequestArr)
	if lErr != nil {
		log.Println("NNAO03", lErr)
		return lNcbExhangeResArr, lErr
	} else {
		log.Println("lResp ADD NCB", lResp)
		lNcbExhangeResArr = lResp
	}

	// Store thre Response in Log table
	// lResponse, lErr := json.Marshal(lResp)
	// if lErr != nil {
	// 	log.Println("NNAO04", lErr)
	// 	return lNcbExhangeResArr, lErr
	// } else {
	// 	log.Println("lResponse", lResponse)
	// 	lLogInputRec.Response = string(lResponse)
	// 	lLogInputRec.LastId = lId
	// 	lLogInputRec.Flag = common.UPDATE
	// 	lId, lErr = Function.LogEntry(lLogInputRec)
	// 	if lErr != nil {
	// 		log.Println("NNAO05", lErr)
	// 		return lNcbExhangeResArr, lErr
	// 	}
	// }
	// }
	// }
	log.Println(lNcbExhangeResArr)
	log.Println("NcbAddOrder....(-)")
	return lNcbExhangeResArr, nil
}

func NcbExchangeOrder(pToken string, pUrl string, pNcbExchangeReqArr NcbAddReqStruct) (NcbAddResStruct, error) {
	log.Println("NcbExchangeOrder....(+)")

	var lNcbRespArr NcbAddResStruct
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

	ljsonData, lErr := json.Marshal(pNcbExchangeReqArr)
	if lErr != nil {
		log.Println("NNEO01", lErr)
		return lNcbRespArr, lErr
	} else {
		lReqstring := string(ljsonData)
		log.Println(pUrl, "pUrl", lReqstring, "lReqstring", lHeaderArr, "lHeaderArr")
		lNcbResp, lErr := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "Token")
		if lErr != nil {
			log.Println("NNEO02", lErr)
			return lNcbRespArr, lErr
		} else {
			log.Println(pUrl, "pUrl", lReqstring, "lReqstring", lHeaderArr, "lHeaderArr")
			log.Println("lNcbResp", lNcbResp)
			lErr = json.Unmarshal([]byte(lNcbResp), &lNcbRespArr)
			if lErr != nil {
				log.Println("NNEO03", lErr, lNcbResp)
				return lNcbRespArr, lErr
			} else {
				log.Println("NcbExchangeResponse---->>>>", lNcbRespArr)
			}
		}
	}
	log.Println("NcbExchangeOrder....(-)")
	return lNcbRespArr, nil
}
