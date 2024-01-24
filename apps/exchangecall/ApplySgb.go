package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bsesgb"
	"fcs23pkg/integration/nse/nsesgb"
	"log"
	"strconv"
)

// func FetchSGBmaster(pReqJson bsesgb.SgbReqStruct, pUser string) error {
// 	log.Println("FetchSGBmaster (+)")
// 	// var pExchangeReq bsesgb.SgbReqStruct

// 	lErr1 := ApplySgb(pExchangeReq, pClientId)
// 	if lErr1 != nil {
// 		log.Println("EFSM01", lErr1)
// 		return lErr1
// 	} else {
// 		lErr2 := (pUser)
// 		if lErr2 != nil {
// 			log.Println("EFSM02", lErr2)
// 			return lErr2
// 		}
// 	}
// 	log.Println("FetchSGBmaster (-)")
// 	return nil
// }

/*
Pupose:  This method is used to inserting the collection of data to the  a_ipo_order_header ,
a_ipo_orderdetails , a_bidTracking tables in  database and also used to place the Bid in NSE
Parameters:
   (ExchangeReqStruct )
Response:
    *On Sucess
    =========
    In case of a successful execution of this method, you will apply for the Bid
	in Exchange using /v1/transaction/addbulk endpoint ad get the response struct from NSE Exchange
    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 09JUNE2023
*/
func ApplyBseSgb(pReqJson bsesgb.SgbReqStruct, pUser string, pBrokerId int) (bsesgb.SgbRespStruct, error) {
	log.Println("ApplyBseSgb (+)")

	var lRespJsonRec bsesgb.SgbRespStruct

	lToken, lErr1 := BseGetToken(pUser, pBrokerId)
	log.Println("lToken", lToken)
	if lErr1 != nil {
		log.Println("EAI01", lErr1)
		return lRespJsonRec, lErr1
	} else {
		if lToken != "" {
			lResp, lErr2 := bsesgb.BseSgbOrder(lToken, pUser, pReqJson)
			if lErr2 != nil {
				log.Println("EAI02", lErr2)
				return lRespJsonRec, lErr2
			} else {
				lRespJsonRec = lResp
			}
		}
	}
	log.Println("ApplyBseSgb (-)")
	return lRespJsonRec, nil
}

func ApplyNseSgb(pReqJson nsesgb.SgbAddReqStruct, pUser string, pBrokerId int) (nsesgb.SgbAddResStruct, error) {
	log.Println("ApplyNseSgb (+)")

	var lRespJsonRec nsesgb.SgbAddResStruct

	lToken, lErr3 := GetToken(pUser, pBrokerId)
	if lErr3 != nil {
		log.Println("EAI03", lErr3)
		return lRespJsonRec, lErr3
	} else {
		if lToken != "" {
			lResp, lErr5 := nsesgb.SgbOrderTransaction(lToken, pReqJson, pUser)
			if lErr5 != nil {
				log.Println("EAI05", lErr5)
				return lRespJsonRec, lErr5
			} else {
				lRespJsonRec = lResp
			}
		}
	}
	log.Println("ApplyNseSgb (-)")
	return lRespJsonRec, nil
}

// this method is changed to commonized for NSE and BSE ,commented by pavithra
// func SgbReqConstruct(pSBseReq bsesgb.SgbReqStruct) nsesgb.SgbAddReqStruct {
// 	log.Println("SgbReqConstruct(+)")
// 	var lSNseReq nsesgb.SgbAddReqStruct
// 	lSNseReq.Symbol = pSBseReq.ScripId
// 	lSNseReq.Depository = pSBseReq.Depository
// 	lSNseReq.DpId = pSBseReq.DpId
// 	lSNseReq.ClientBenId = pSBseReq.ClientBenfId
// 	lSNseReq.PhysicalDematFlag = "D"
// 	lSNseReq.ClientCode = ""
// 	lSNseReq.Pan = pSBseReq.PanNo
// 	lByte := make([]byte, 8)
// 	rand.Read(lByte)
// 	lSNseReq.ClientRefNumber = fmt.Sprintf("%x", lByte)
// 	for ldx := 0; ldx < len(pSBseReq.Bids); ldx++ {
// 		if pSBseReq.Bids[ldx].ActionCode == "N" {
// 			lSNseReq.ActivityType = "ER"
// 		} else if pSBseReq.Bids[ldx].ActionCode == "M" {
// 			lSNseReq.ActivityType = "MR"
// 		} else {
// 			lSNseReq.ActivityType = "CR"
// 		}
// 		lSubscriptionUnit, lErr1 := strconv.Atoi(pSBseReq.Bids[ldx].SubscriptionUnit)
// 		if lErr1 != nil {
// 			log.Println("Error: 1", lErr1)
// 		} else {
// 			lSNseReq.Quantity = lSubscriptionUnit
// 		}
// 		lRate, lErr2 := common.ConvertStringToFloat(pSBseReq.Bids[ldx].Rate)
// 		if lErr2 != nil {
// 			log.Println("Error: 2", lErr2)
// 		} else {
// 			lSNseReq.Price = lRate
// 		}
// 		OrderNo, lErr3 := strconv.Atoi(pSBseReq.Bids[ldx].OrderNo)
// 		if lErr3 != nil {
// 			log.Println("Error: 3", lErr3)
// 		} else {
// 			lSNseReq.OrderNumber = OrderNo
// 		}
// 	}
// 	// lSNseArr = append(lSNseArr, lSNseReq)
// 	log.Println("SgbReqConstruct(-)")
// 	return lSNseReq
// }

func SgbRespConstruct(pSNseResp nsesgb.SgbAddResStruct) bsesgb.SgbRespStruct {
	log.Println("SgbRespConstruct(+)")

	var lSBseResp bsesgb.SgbRespStruct
	var lSBseBid bsesgb.RespSgbBidStruct

	// for _, lNseResp := range pSNseResp {
	lSBseResp.ScripId = pSNseResp.Symbol
	// log.Println("lSBseResp.ScripId", lSBseResp.ScripId)
	lSBseResp.Depository = pSNseResp.Depository
	// log.Println("lSBseResp.Depository", lSBseResp.Depository)

	lSBseResp.DpId = ""
	// log.Println("lSBseResp.DpId", lSBseResp.DpId)

	lSBseResp.ClientBenfId = pSNseResp.ClientBenId
	// log.Println("lSBseResp.ClientBenfId", lSBseResp.ClientBenfId)

	lSBseResp.PanNo = pSNseResp.Pan
	// log.Println("lSBseResp.PanNo", lSBseResp.PanNo)

	if pSNseResp.Status == "success" {
		log.Println("if")

		lSBseResp.StatusCode = "0"
		lSBseBid.ErrorCode = "0"
		// log.Println("lSBseBid.ErrorCode", lSBseBid.ErrorCode)

		lSBseResp.StatusMessage = pSNseResp.Reason
		// log.Println("lSBseResp.StatusMessage", lSBseResp.StatusMessage)

	} else {
		log.Println("elseif")
		lSBseResp.StatusCode = "1"
		lSBseBid.ErrorCode = "1"
		// log.Println("lSBseBid.ErrorCode ", lSBseBid.ErrorCode)

		lSBseResp.ErrorMessage = pSNseResp.Reason
		// log.Println("lSBseResp.ErrorMessage", lSBseResp.ErrorMessage)

	}

	// pSNseResp.Status = lSBseResp.StatusCode
	// log.Println("pSNseResp.Status", pSNseResp.Status)

	// pSNseResp.Status = lSBseBid.ErrorCode
	// log.Println("pSNseResp.Status ", pSNseResp.Status)

	// if pSNseReq.ActivityType == "ER" {
	// 	lSBseBid.ActionCode = "N"
	// } else if pSNseReq.ActivityType == "MR" {
	// 	lSBseBid.ActionCode = "M"
	// } else {
	// 	lSBseBid.ActionCode = "C"
	// }

	if pSNseResp.OrderStatus == "ES" || pSNseResp.OrderStatus == "EF" {
		lSBseBid.ActionCode = "N"
	} else if pSNseResp.OrderStatus == "MS" || pSNseResp.OrderStatus == "MF" {
		lSBseBid.ActionCode = "M"
	} else {
		lSBseBid.ActionCode = "C"
	}

	lQuantity := strconv.Itoa(pSNseResp.Quantity)
	lSBseBid.SubscriptionUnit = lQuantity

	lOrderNumber := strconv.Itoa(pSNseResp.OrderNumber)
	lSBseBid.BidId = lOrderNumber

	lSBseBid.OrderNo = pSNseResp.ApplicationNumber

	lPrice := strconv.FormatFloat(float64(pSNseResp.Price), 'f', -1, 32)
	lSBseBid.Rate = lPrice

	if pSNseResp.OrderStatus == "ES" || pSNseResp.OrderStatus == "MS" || pSNseResp.OrderStatus == "CS" {
		lSBseBid.ErrorCode = "0"
		lSBseBid.Message = pSNseResp.RejectionReason
	} else {
		lSBseBid.ErrorCode = "1"
		lSBseBid.ErrorCode = pSNseResp.OrderStatus
		lSBseBid.Message = pSNseResp.RejectionReason
	}
	lSBseResp.Bids = append(lSBseResp.Bids, lSBseBid)

	log.Println("lSBseResp", lSBseResp)

	log.Println("SgbRespConstruct(-)")
	return lSBseResp
}

type JvReqStruct struct {
	ReqOrderNo              string         `json:"reqOrderNo"`
	RespOrderNo             string         `json:"respOrderNo"`
	ClientId                string         `json:"clientId"`
	BoJvStatus              string         `json:"boJvstatus"`
	BoJvStatement           string         `json:"boJvstatement"`
	BoJvAmount              string         `json:"boJvamount"`
	BoJvType                string         `json:"boJvtype"`
	FoJvStatus              string         `json:"foJvstatus"`
	FoJvStatement           string         `json:"foJvstatement"`
	FoJvAmount              string         `json:"foJvamount"`
	FoJvType                string         `json:"foJvtype"`
	Unit                    int            `json:"unit"`
	Price                   int            `json:"price"`
	ActivityType            string         `json:"activityType"`
	Symbol                  string         `json:"symbol"`
	OrderDate               string         `json:"orderdate"`
	Amount                  string         `json:"amount"`
	Mail                    string         `json:"mail"`
	ClientName              string         `json:"clientname"`
	EmailDist               string         `json:"emailDist"`
	EmailMsg                string         `json:"emailMsg"`
	ErrorStage              string         `json:"errorStage"`
	ProcessStatus           string         `json:"processStatus"`
	PanNo                   string         `json:"panno"`
	InvestorCategory        string         `json:"invcategory"`
	ApplicantName           string         `json:"applicantname"`
	Depository              string         `json:"depository"`
	DpId                    string         `json:"dpid"`
	ClientBenfId            string         `json:"clientid"`
	GuardianName            string         `json:"guardianname"`
	GuardianPanno           string         `json:"guardianpanno"`
	GuardianRelation        string         `json:"guardianrelation"`
	StatusCode              string         `json:"statuscode"`
	StatusMessage           string         `json:"statusmessage"`
	ErrorCode               string         `json:"errorcode"`
	ErrorMessage            string         `json:"errormessage"`
	PhysicalDematFlag       string         `json:"physicalDematFlag"`
	ClientCode              string         `json:"clientCode"`
	ClientRefNumber         string         `json:"clientRefNumber"`
	Series                  string         `json:"series"`
	ApplicationNumber       string         `json:"applicationNumber"`
	EntryTime               string         `json:"entryTime"`
	VerificationStatus      string         `json:"verificationStatus"`
	VerificationReason      string         `json:"verificationReason"`
	ClearingStatus          string         `json:"clearingStatus"`
	ClearingReason          string         `json:"clearingReason"`
	LastActionTime          string         `json:"lastActionTime"`
	RejectionReason         string         `json:"rejectionReason"`
	Bids                    []SgbBidStruct `json:"bids"`
	TotalReqCount           int            `json:"totalReqCount"`
	FoVerifySuccessCount    int            `json:"foVerifySuccessCount"`
	FoVerifyFailedCount     int            `json:"foVerifyFailedCount"`
	BoJvSuccesssCount       int            `json:"jvSuccesssCount"`
	BoJvFailedCount         int            `json:"jvFailedCount"`
	ExchangeSuccessCount    int            `json:"exchangeSuccessCount"`
	ExchangeFailedCount     int            `json:"exchangeFailedCount"`
	ReverseBoJvSuccessCount int            `json:"reverseBoJvSuccessCount"`
	ReverseBoJvFailedCount  int            `json:"reverseBoJvFailedCount"`
	ReverseFoJvSuccessCount int            `json:"reverseFoJvSuccessCount"`
	ReverseFoJvFailedCount  int            `json:"reverseFoJvFailedCount"`
	FoJvSuccessCount        int            `json:"foJvSuccessCount"`
	FoJvFailedCount         int            `json:"foJvFailedCount"`
	TotalSuccessCount       int            `json:"totalSuccessCount"`
	TotalFailedCount        int            `json:"totalFailedCount"`
	InsufficientFund        int            `json:"insufficientFund"`
}

type SgbBidStruct struct {
	BidId            string `json:"bidid"`
	SubscriptionUnit string `json:"subscriptionunit"`
	Rate             string `json:"rate"`
	OrderNo          string `json:"orderno"`
	ActionCode       string `json:"actioncode"`
	ErrorCode        string `json:"errorcode"`
	Message          string `json:"message"`
}

/*
Pupose: this method is used to construct the Exchange structure for placing SGB order in NSE and BSE
Parameters:
  JvdetailRec,Exchange
Response:
    *On Sucess
    =========
   If the Exchange is NSE - it returns NSE Struct
   If the Exchange is BSE -  it returns BSE Struct
    !On Error
    ========
   In Case of Error, Error will be Return

Author: Pavithra
Date: 30Dec2023
*/
func ConstructExchReq(pReqStruct JvReqStruct, pExchange string) (nsesgb.SgbAddReqStruct, bsesgb.SgbReqStruct) {
	log.Println("ConstructExchReq (+)")

	var lSNseReq nsesgb.SgbAddReqStruct
	var lBseReq bsesgb.SgbReqStruct

	if pExchange == common.NSE {
		lSNseReq.Symbol = pReqStruct.Symbol
		lSNseReq.Depository = pReqStruct.Depository
		lSNseReq.DpId = pReqStruct.DpId
		lSNseReq.ClientBenId = pReqStruct.ClientBenfId
		lSNseReq.PhysicalDematFlag = "D"
		lSNseReq.ClientCode = ""
		lSNseReq.Pan = pReqStruct.PanNo
		lRandomNo, lErr1 := GetSGB_SequenceNo()
		if lErr1 != nil {
			log.Println("ECER01", lErr1)
			lSNseReq.OrderNumber = 0
			// return lSNseReq, lBseReq, lErr1
		} else {
			var lTrimmedString string
			if len(strconv.Itoa(lRandomNo)) >= 5 {
				lTrimmedString = strconv.Itoa(lRandomNo)[len(strconv.Itoa(lRandomNo))-5:]
			}
			lSNseReq.ClientRefNumber = pReqStruct.ClientId + lTrimmedString
		}
		lSNseReq.Quantity = pReqStruct.Unit
		lSNseReq.Price = float32(pReqStruct.Price)
		for ldx := 0; ldx < len(pReqStruct.Bids); ldx++ {
			if pReqStruct.Bids[ldx].ActionCode == "N" {
				lSNseReq.ActivityType = "ER"
			} else if pReqStruct.Bids[ldx].ActionCode == "M" {
				lSNseReq.ActivityType = "MR"
			} else {
				lSNseReq.ActivityType = "CR"
			}
		}
		if pReqStruct.ReqOrderNo != "" {
			OrderNo, lErr1 := strconv.Atoi(pReqStruct.ReqOrderNo)
			if lErr1 != nil {
				log.Println("ECER02", lErr1)
				lSNseReq.OrderNumber = 0
				// return lSNseReq, lBseReq, lErr1
			} else {
				lSNseReq.OrderNumber = OrderNo
			}
		} else {
			lSNseReq.OrderNumber = 0
		}
		log.Println("NSE Req Struct", lSNseReq)
	} else {
		lBseReq.ScripId = pReqStruct.Symbol
		lBseReq.PanNo = pReqStruct.PanNo
		lBseReq.InvestorCategory = pReqStruct.InvestorCategory
		lBseReq.ApplicantName = pReqStruct.ApplicantName
		lBseReq.ClientBenfId = pReqStruct.ClientBenfId
		lBseReq.Depository = pReqStruct.Depository
		lBseReq.DpId = pReqStruct.DpId
		lBseReq.GuardianName = pReqStruct.GuardianName
		lBseReq.GuardianPanno = pReqStruct.GuardianPanno
		lBseReq.GuardianRelation = pReqStruct.GuardianRelation

		var lBidRec bsesgb.ReqSgbBidStruct
		for lIdx := 0; lIdx < len(pReqStruct.Bids); lIdx++ {
			lBidRec.BidId = pReqStruct.Bids[lIdx].BidId
			lBidRec.SubscriptionUnit = pReqStruct.Bids[lIdx].SubscriptionUnit
			lBidRec.Rate = pReqStruct.Bids[lIdx].Rate
			lBidRec.OrderNo = pReqStruct.Bids[lIdx].OrderNo
			lBidRec.ActionCode = pReqStruct.Bids[lIdx].ActionCode
		}
		lBseReq.Bids = append(lBseReq.Bids, lBidRec)
		log.Println("BSE Req Struct", lSNseReq)
	}
	log.Println("ConstructExchReq (-)")
	return lSNseReq, lBseReq
}

/*
Pupose: this method is used to construct the JvReqStruct for updating all the NSE and BSE deatils in DB
Parameters:
  NSEResponseStruct,BSEResponseStruct,Exchange
Response:
    *On Sucess
    =========
	JvReqStruct
    !On Error
    ========
  {}
Author: Pavithra
Date: 30Dec2023
*/
func ConstructExchResp(pNseRespRec nsesgb.SgbAddResStruct, pBseRespRec bsesgb.SgbRespStruct, pExchange string, pSgbRec JvReqStruct) JvReqStruct {
	log.Println("ConstructExchResp (+)")

	if pExchange == common.NSE {
		pSgbRec.Depository = pNseRespRec.Depository
		pSgbRec.DpId = "0"
		pSgbRec.ClientBenfId = pNseRespRec.ClientBenId
		pSgbRec.PanNo = pNseRespRec.Pan
		pSgbRec.PhysicalDematFlag = pNseRespRec.PhysicalDematFlag
		pSgbRec.RespOrderNo = strconv.Itoa(pNseRespRec.OrderNumber)
		pSgbRec.ApplicationNumber = pNseRespRec.ApplicationNumber

		for lIdx := 0; lIdx < len(pSgbRec.Bids); lIdx++ {

			pSgbRec.Bids[lIdx].SubscriptionUnit = strconv.Itoa(pNseRespRec.Quantity)

			pSgbRec.Bids[lIdx].Rate = strconv.FormatFloat(float64(pNseRespRec.Price), 'f', -1, 32)
			if pNseRespRec.Status == "success" {
				pSgbRec.StatusCode = "0"
				pSgbRec.StatusMessage = pNseRespRec.Reason
			} else {
				log.Println("elseif")
				pSgbRec.StatusCode = "1"
				pSgbRec.ErrorMessage = pNseRespRec.Reason
			}

			if pNseRespRec.OrderStatus == "ES" || pNseRespRec.OrderStatus == "EF" {
				pSgbRec.Bids[lIdx].ActionCode = "N"
			} else if pNseRespRec.OrderStatus == "MS" || pNseRespRec.OrderStatus == "MF" {
				pSgbRec.Bids[lIdx].ActionCode = "M"
			} else {
				pSgbRec.Bids[lIdx].ActionCode = "C"
			}
			if pNseRespRec.OrderStatus == "ES" || pNseRespRec.OrderStatus == "MS" || pNseRespRec.OrderStatus == "CS" {
				pSgbRec.Bids[lIdx].ErrorCode = "0"
				pSgbRec.Bids[lIdx].Message = pNseRespRec.RejectionReason
			} else {
				pSgbRec.Bids[lIdx].ErrorCode = "1"
				pSgbRec.Bids[lIdx].Message = pNseRespRec.RejectionReason
			}
		}

	} else {
		pSgbRec.Depository = pBseRespRec.Depository
		pSgbRec.DpId = pBseRespRec.DpId
		pSgbRec.PanNo = pBseRespRec.PanNo
		pSgbRec.InvestorCategory = pBseRespRec.InvestorCategory
		pSgbRec.ApplicantName = pBseRespRec.ApplicantName
		pSgbRec.ClientBenfId = pBseRespRec.ClientBenfId
		pSgbRec.GuardianName = pBseRespRec.GuardianName
		pSgbRec.GuardianPanno = pBseRespRec.GuardianPanno
		pSgbRec.GuardianRelation = pBseRespRec.GuardianRelation
		pSgbRec.StatusCode = pBseRespRec.StatusCode
		pSgbRec.StatusMessage = pBseRespRec.StatusMessage
		pSgbRec.ErrorCode = pBseRespRec.ErrorCode
		pSgbRec.ErrorMessage = pBseRespRec.ErrorMessage
		for lBidIdx := 0; lBidIdx < len(pBseRespRec.Bids); lBidIdx++ {
			pSgbRec.Bids[lBidIdx].ActionCode = pBseRespRec.Bids[lBidIdx].ActionCode
			pSgbRec.Bids[lBidIdx].BidId = pBseRespRec.Bids[lBidIdx].BidId
			pSgbRec.Bids[lBidIdx].OrderNo = pBseRespRec.Bids[lBidIdx].OrderNo
			pSgbRec.Bids[lBidIdx].SubscriptionUnit = pBseRespRec.Bids[lBidIdx].SubscriptionUnit
			pSgbRec.Bids[lBidIdx].Rate = pBseRespRec.Bids[lBidIdx].Rate
			pSgbRec.Bids[lBidIdx].ErrorCode = pBseRespRec.Bids[lBidIdx].ErrorCode
			pSgbRec.Bids[lBidIdx].Message = pBseRespRec.Bids[lBidIdx].Message
		}
	}
	log.Println("ConstructExchResp (-)")
	return pSgbRec
}

/*
Pupose:This method is used to get sequence no for unique ref no
Parameters:
nil
Response:

	==========
	*On Sucess
	==========
	123456789012,nil

	==========
	*On Error
	==========
	0,error

Author:Pavithra
Date: 08JAN2024
*/
func GetSGB_SequenceNo() (int, error) {
	log.Println("GetSGB_SequenceNo (+)")

	// this variables is used to get DpId and ClientName from the database.
	var lSequenceNo int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPGSN01", lErr1)
		return lSequenceNo, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `SELECT NEXT VALUE FOR a_sgb_reference_s;`

		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("SPGSN02", lErr2)
			return lSequenceNo, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSequenceNo)
				if lErr3 != nil {
					log.Println("SPGSN03", lErr3)
					return lSequenceNo, lErr3
				}
			}
		}
	}
	log.Println("GetSGB_SequenceNo (-)")
	return lSequenceNo, nil
}
