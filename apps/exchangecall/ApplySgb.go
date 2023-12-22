package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/integration/bse/bsesgb"
	"fcs23pkg/integration/nse/nsesgb"
	"fmt"
	"log"
	"math/rand"
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

func SgbReqConstruct(pSBseReq bsesgb.SgbReqStruct) nsesgb.SgbAddReqStruct {
	log.Println("SgbReqConstruct(+)")

	var lSNseReq nsesgb.SgbAddReqStruct

	lSNseReq.Symbol = pSBseReq.ScripId
	lSNseReq.Depository = pSBseReq.Depository
	lSNseReq.DpId = pSBseReq.DpId
	lSNseReq.ClientBenId = pSBseReq.ClientBenfId

	lSNseReq.PhysicalDematFlag = "D"
	lSNseReq.ClientCode = ""
	lSNseReq.Pan = pSBseReq.PanNo

	lByte := make([]byte, 8)
	rand.Read(lByte)
	lSNseReq.ClientRefNumber = fmt.Sprintf("%x", lByte)

	for ldx := 0; ldx < len(pSBseReq.Bids); ldx++ {

		if pSBseReq.Bids[ldx].ActionCode == "N" {
			lSNseReq.ActivityType = "ER"
		} else if pSBseReq.Bids[ldx].ActionCode == "M" {
			lSNseReq.ActivityType = "MR"
		} else {
			lSNseReq.ActivityType = "CR"
		}

		lSubscriptionUnit, lErr1 := strconv.Atoi(pSBseReq.Bids[ldx].SubscriptionUnit)
		if lErr1 != nil {
			log.Println("Error: 1", lErr1)
		} else {

			lSNseReq.Quantity = lSubscriptionUnit
		}

		lRate, lErr2 := common.ConvertStringToFloat(pSBseReq.Bids[ldx].Rate)
		if lErr2 != nil {
			log.Println("Error: 2", lErr2)
		} else {
			lSNseReq.Price = lRate
		}

		OrderNo, lErr3 := strconv.Atoi(pSBseReq.Bids[ldx].OrderNo)
		if lErr3 != nil {
			log.Println("Error: 3", lErr3)
		} else {
			lSNseReq.OrderNumber = OrderNo
		}

	}

	// lSNseArr = append(lSNseArr, lSNseReq)

	log.Println("SgbReqConstruct(-)")
	return lSNseReq

}

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
