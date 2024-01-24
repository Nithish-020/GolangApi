package bsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/integration/bse/bseipo"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type SgbStatusReqStruct struct {
	ScripId   string `json:"scripid"`
	StartDate string `json:"startdate"`
	StartTime string `json:"starttime"`
	EndDate   string `json:"enddate"`
	EndTime   string `json:"endtime"`
	PanNo     string `json:"panno"`
}

type SgbDataStruct struct {
	ScripId string   `json:"scripid"`
	PanNo   []string `json:"panno"`
}
type SgbDownDataRespStruct struct {
	ScripId              string                   `json:"scripid"`
	PanNo                string                   `json:"panno"`
	InvestorCategory     string                   `json:"invcategory"`
	ApplicantName        string                   `json:"applicantname"`
	Depository           string                   `json:"depository"`
	DpId                 string                   `json:"dpid"`
	ClientBenfId         string                   `json:"clientid"`
	GuardianName         string                   `json:"guardianname"`
	GuardianRelationShip string                   `json:"guardianrelationship"`
	GuardianPanno        string                   `json:"guardianpanno"`
	RbiInvstId           string                   `json:"rbiinvstid"`
	DpStatus             string                   `json:"dpstatus"`
	DpRemarks            string                   `json:"dpremarks"`
	RbiStatus            string                   `json:"rbistatus"`
	RbiRemarks           string                   `json:"rbiremarks"`
	LoginId              string                   `json:"loginid"`
	StatusCode           string                   `json:"statuscode"`
	BranchCode           string                   `json:"branchcode"`
	ErrorCode            string                   `json:"errorcode"`
	Message              string                   `json:"message"`
	BankAccNo            string                   `json:"bankaccno"`
	Ifsc                 string                   `json:"ifsc"`
	Addr1                string                   `json:"addr1"`
	Addr2                string                   `json:"addr2"`
	Addr3                string                   `json:"addr3"`
	State                string                   `json:"state"`
	Pincode              string                   `json:"pincode"`
	Email                string                   `json:"email"`
	Mobile               string                   `json:"mobile"`
	Bids                 []RespSgbStatusBidStruct `json:"bids"`
}
type RespSgbStatusBidStruct struct {
	BidId            string `json:"bidid"`
	SubscriptionUnit string `json:"subscriptionunit"`
	Rate             string `json:"rate"`
	OrderNo          string `json:"orderno"`
	OrderStatus      string `json:"orderStatus"`
	Addeddate        string `json:"addeddate"`
	ModifiedDate     string `json:"modifyDate"`
}

func BseSgbDownOrder(pToken string, pUser string, pReqRec SgbStatusReqStruct, pBrokerId int) ([]SgbDownDataRespStruct, error) {
	log.Println("BseSgbDownOrder (+)")
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	//create instance to receive SgbRespStruct
	var lApiRespRec []SgbDownDataRespStruct
	// create instance to hold the last inserted id
	// var lId int
	// create instance to hold errors
	// var lErr1 error

	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseSgbDownload"])
	log.Println(lUrl, "endpoint")

	// lRequest, lErr := json.Marshal(pReqRec)
	// if lErr != nil {
	// 	log.Println("NT01", lErr)
	// 	return lApiRespRec, lErr
	// } else {
	// 	lLogInputRec.Request = string(lRequest)
	// 	lLogInputRec.EndPoint = "/v1/sgbdownload"
	// 	lLogInputRec.Flag = common.INSERT
	// 	lLogInputRec.ClientId = pUser
	// 	lLogInputRec.Method = "POST"

	// 	// ! LogEntry method is used to store the Request in Database
	// 	lId, lErr1 = Function.LogEntry(lLogInputRec)
	// 	if lErr1 != nil {
	// 		log.Println("NSM01", lErr1)
	// 		return lApiRespRec, lErr1
	// 	} else {
	// TokenApi method used to call exchange API
	lResp, lErr2 := BseSgbDownOrderApi(pToken, lUrl, pReqRec, pBrokerId)
	if lErr2 != nil {
		log.Println("NSM02", lErr2)
		return lApiRespRec, lErr2
	} else {
		lApiRespRec = lResp
	}
	log.Println("Response", lResp)
	// Store thre Response in Log table
	// lResponse, lErr3 := json.Marshal(lResp)
	// if lErr3 != nil {
	// 	log.Println("NSM03", lErr3)
	// 	return lApiRespRec, lErr3
	// } else {
	// lLogInputRec.Response = string(lResponse)
	// lLogInputRec.LastId = lId
	// lLogInputRec.Flag = common.UPDATE
	// // create instance to hold errors
	// var lErr4 error
	// lId, lErr4 = Function.LogEntry(lLogInputRec)
	// if lErr4 != nil {
	// 	log.Println("NSM04", lErr4)
	// 	return lApiRespRec, lErr4
	// }
	// }
	// }
	// }
	log.Println("BseSgbDownOrder (-)")
	return lApiRespRec, nil
}

func BseSgbDownOrderApi(pToken string, pUrl string, pReqRec SgbStatusReqStruct, pBrokerId int) ([]SgbDownDataRespStruct, error) {
	log.Println("BseSgbOrderApi (+)")
	// create a new instance of IpoResponseStruct
	var lSgbRespRec []SgbDownDataRespStruct
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails
	// create a Array instance of HeaderDetails struct
	var lHeaderArr []apiUtil.HeaderDetails

	lDetail, lErr := bseipo.AccessBseCredential(pBrokerId)
	if lErr != nil {
		log.Println("NESM01", lErr)
		return lSgbRespRec, lErr
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
		ljsonData, lErr := json.Marshal(pReqRec)
		if lErr != nil {
			log.Println("NEO01", lErr)
			return lSgbRespRec, lErr
		} else {
			lReqstring := string(ljsonData)
			lResp, lErr1 := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "SGBOrder")
			if lErr1 != nil {
				log.Println("NESM02", lErr1)
				return lSgbRespRec, lErr1
			} else {
				lErr2 := json.Unmarshal([]byte(lResp), &lSgbRespRec)
				if lErr2 != nil {
					log.Println("NESM03", lErr2)
					return lSgbRespRec, lErr2
				}
			}
		}
	}
	log.Println("BseSgbOrderApi (-)")
	return lSgbRespRec, nil
}
