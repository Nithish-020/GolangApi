package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/integration/nse/nsencb"
	"log"
)

type NcbJvReqStruct struct {
	ReqOrderNo              int     `json:"reqOrderNo"`
	RespOrderNo             int     `json:"respOrderNo"`
	ClientId                string  `json:"clientId"`
	BoJvStatus              string  `json:"boJvstatus"`
	BoJvStatement           string  `json:"boJvstatement"`
	BoJvAmount              string  `json:"boJvamount"`
	BoJvType                string  `json:"boJvtype"`
	FoJvStatus              string  `json:"foJvstatus"`
	FoJvStatement           string  `json:"foJvstatement"`
	FoJvAmount              string  `json:"foJvamount"`
	FoJvType                string  `json:"foJvtype"`
	Unit                    int     `json:"unit"`
	Price                   float64 `json:"price"`
	Amount                  float64 `json:"amount"`
	ActivityType            string  `json:"activityType"`
	Symbol                  string  `json:"symbol"`
	OrderDate               string  `json:"orderdate"`
	Mail                    string  `json:"mail"`
	ClientName              string  `json:"clientname"`
	EmailDist               string  `json:"emailDist"`
	EmailMsg                string  `json:"emailMsg"`
	ErrorStage              string  `json:"errorStage"`
	ProcessStatus           string  `json:"processStatus"`
	PanNo                   string  `json:"panno"`
	Depository              string  `json:"depository"`
	DpId                    string  `json:"dpid"`
	ClientBenId             string  `json:"clientBenId"`
	Status                  string  `json:"status"`
	Reason                  string  `json:"reason"`
	StatusCode              string  `json:"statuscode"`
	StatusMessage           string  `json:"statusmessage"`
	ErrorCode               string  `json:"errorcode"`
	ErrorMessage            string  `json:"errormessage"`
	PhysicalDematFlag       string  `json:"physicalDematFlag"`
	ClientCode              string  `json:"clientCode"`
	ClientRefNumber         string  `json:"clientRefNumber"`
	Series                  string  `json:"series"`
	ApplicationNumber       string  `json:"applicationNumber"`
	EntryTime               string  `json:"entryTime"`
	VerificationStatus      string  `json:"verificationStatus"`
	VerificationReason      string  `json:"verificationReason"`
	ClearingStatus          string  `json:"clearingStatus"`
	ClearingReason          string  `json:"clearingReason"`
	LastActionTime          string  `json:"lastActionTime"`
	RejectionReason         string  `json:"rejectionReason"`
	TotalReqCount           int     `json:"totalReqCount"`
	FoVerifySuccessCount    int     `json:"foVerifySuccessCount"`
	FoVerifyFailedCount     int     `json:"foVerifyFailedCount"`
	BoJvSuccesssCount       int     `json:"jvSuccesssCount"`
	BoJvFailedCount         int     `json:"jvFailedCount"`
	ExchangeSuccessCount    int     `json:"exchangeSuccessCount"`
	ExchangeFailedCount     int     `json:"exchangeFailedCount"`
	ReverseBoJvSuccessCount int     `json:"reverseBoJvSuccessCount"`
	ReverseBoJvFailedCount  int     `json:"reverseBoJvFailedCount"`
	ReverseFoJvSuccessCount int     `json:"reverseFoJvSuccessCount"`
	ReverseFoJvFailedCount  int     `json:"reverseFoJvFailedCount"`
	FoJvSuccessCount        int     `json:"foJvSuccessCount"`
	FoJvFailedCount         int     `json:"foJvFailedCount"`
	TotalSuccessCount       int     `json:"totalSuccessCount"`
	TotalFailedCount        int     `json:"totalFailedCount"`
	InsufficientFund        int     `json:"insufficientFund"`
}

func ApplyNseNcb(pReqJson nsencb.NcbAddReqStruct, pUser string, pBrokerId int) (nsencb.NcbAddResStruct, error) {
	log.Println("ApplyNseNcb(+)")

	var lRespJsonRec nsencb.NcbAddResStruct

	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EANN01", lErr1)
		return lRespJsonRec, lErr1
	} else {
		if lToken != "" {
			lResp, lErr2 := nsencb.NcbAddOrder(lToken, pReqJson, pUser)
			if lErr2 != nil {
				log.Println("EANN02", lErr2)
				return lRespJsonRec, lErr2
			} else {
				lRespJsonRec = lResp
				log.Println("lRespJsonRec", lRespJsonRec, lResp, "lResp")
			}
		}
	}
	log.Println("ApplyNseNcb(-)")
	return lRespJsonRec, nil
}

/*
Pupose: this method is used to construct the Exchange structure for placing NCB order in NSE and BSE
Parameters:
  JvdetailRec,Exchange
Response:
    *On Sucess
    =========
   If the Exchange is NSE - it returns NSE Struct
   If the Exchange is BSE -  nil
    !On Error
    ========
   In Case of Error, Error will be Return

Author: KAVYA DHARSHANI
Date: 05JAN2023
*/

func NcbConstructExchReq(pReqStruct NcbJvReqStruct, pExchange string) nsencb.NcbAddReqStruct {
	log.Println("NcbConstructExchReq (+)")

	var lNseNcbReq nsencb.NcbAddReqStruct

	if pExchange == common.NSE {

		log.Println("pExchange", pExchange)

		lNseNcbReq.Symbol = pReqStruct.Symbol
		lNseNcbReq.Depository = pReqStruct.Depository
		lNseNcbReq.DpId = pReqStruct.DpId
		lNseNcbReq.ClientBenId = pReqStruct.ClientBenId
		lNseNcbReq.PhysicalDematFlag = "D"
		lNseNcbReq.Pan = pReqStruct.PanNo
		lNseNcbReq.ClientRefNumber = common.Generate16DigitString(pReqStruct.ClientId)
		lNseNcbReq.Price = pReqStruct.Price

		if pReqStruct.ActivityType == "N" {
			lNseNcbReq.ActivityType = "ER"
		} else if pReqStruct.ActivityType == "M" {
			lNseNcbReq.ActivityType = "MR"
		} else {
			lNseNcbReq.ActivityType = "CR"
		}

		lNseNcbReq.OrderNumber = pReqStruct.ReqOrderNo
		log.Println("NSE Req Struct", lNseNcbReq)

	} else {

		log.Println("pExchange", pExchange)
	}

	log.Println("ConstructExchReq (-)")
	return lNseNcbReq
}

/*
Pupose: this method is used to construct the JvReqStruct for updating all the NSE and BSE deatils in DB
Parameters:
  NSEResponseStruct,Exchange
Response:
    *On Sucess
    =========
	NcbJvReqStruct
    !On Error
    ========
  {}
Author: KAVYADHARSHANI
Date: 5JAN2023
*/

func NcbConstructExchResp(pNseRespRec nsencb.NcbAddResStruct, pExchange string) NcbJvReqStruct {

	log.Println("NcbConstructExchResp (+)")

	var lRespRec NcbJvReqStruct
	if pExchange == common.NSE {
		log.Println("pExchange", pExchange)

		lRespRec.Depository = pNseRespRec.Depository
		lRespRec.DpId = "0"
		lRespRec.ClientBenId = pNseRespRec.ClientBenId
		lRespRec.PanNo = pNseRespRec.Pan
		lRespRec.PhysicalDematFlag = pNseRespRec.PhysicalDematFlag
		lRespRec.Unit = pNseRespRec.InvestmentValue
		lRespRec.RespOrderNo = pNseRespRec.OrderNumber
		lRespRec.ApplicationNumber = pNseRespRec.ApplicationNumber

		if pNseRespRec.Status == "success" {
			lRespRec.StatusCode = "0"
			lRespRec.StatusMessage = pNseRespRec.Reason
		} else {
			log.Println("elseif")
			lRespRec.StatusCode = "1"
			lRespRec.ErrorMessage = pNseRespRec.Reason
		}

		if pNseRespRec.OrderStatus == "ES" || pNseRespRec.OrderStatus == "EF" {
			lRespRec.ActivityType = "N"
		} else if pNseRespRec.OrderStatus == "MS" || pNseRespRec.OrderStatus == "MF" {
			lRespRec.ActivityType = "M"
		} else {
			lRespRec.ActivityType = "C"
		}

		if pNseRespRec.OrderStatus == "ES" || pNseRespRec.OrderStatus == "MS" || pNseRespRec.OrderStatus == "CS" {
			lRespRec.ErrorCode = "0"
			lRespRec.Reason = pNseRespRec.RejectionReason
		} else {
			lRespRec.ErrorCode = "1"
			lRespRec.Reason = pNseRespRec.RejectionReason
		}

	} else {
		log.Println("pExchange", pExchange)
	}

	log.Println("ConstructExchResp (-)")
	return lRespRec
}
