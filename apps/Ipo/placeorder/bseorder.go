package placeorder

import (
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/integration/bse/bseipo"
	"fcs23pkg/integration/nse/nseipo"
	"fmt"
	"log"
	"strconv"
	"time"
)

//added by naveen:add one parameter source to insert in bidtracking table
//func BsePlaceOrder(pExhangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int)
func BsePlaceOrder(pExhangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int, pSource string) (OrderResStruct, []nseipo.ExchangeRespStruct, ErrorStruct) {
	log.Println("BsePlaceOrder (+)")

	var lRespRec OrderResStruct
	var lError ErrorStruct
	//added by naveen:add one argument source to insert in bidtracking
	//lResp, lErr9 := ProcessBseReq(pExhangeReq, pReqRec, pClientId, pMailId, pBrokerId)
	lResp, lErr9 := ProcessBseReq(pExhangeReq, pReqRec, pClientId, pMailId, pBrokerId, pSource)

	log.Println("lResp", lResp)
	if lErr9 != nil {
		log.Println("PBPO01", lErr9.Error())
		lRespRec.Status = common.ErrorCode
		lError.ErrCode = "PBPO01"
		lError.ErrMsg = "Exchange Server is Busy right now,Try After Sometime."
		return lRespRec, lResp, lError
	} else {
		if lResp != nil {
			for lRespIdx := 0; lRespIdx < len(lResp); lRespIdx++ {
				if lResp[lRespIdx].Status == "failed" {

					log.Println("lResp[lRespIdx].Status", lResp[lRespIdx].Status)
					lRespRec.AppStatus = lResp[lRespIdx].Status
					lRespRec.AppReason = lResp[lRespIdx].Reason
					lRespRec.Status = common.ErrorCode
					log.Println("PBPO02", lRespRec.AppStatus, lRespRec.AppReason)

					lError.ErrCode = "PBPO02"
					lError.ErrMsg = lResp[lRespIdx].Reason + "\n" + "Application Failed"
					return lRespRec, lResp, lError
				} else if lResp[lRespIdx].Status == "success" {

					log.Println("lResp[lRespIdx].Status", lResp[lRespIdx].Status)
					lRespRec.AppStatus = lResp[lRespIdx].Status
					lRespRec.AppReason = lResp[lRespIdx].Reason
					lRespRec.Status = common.SuccessCode
				} else {
					log.Println("lResp[lRespIdx].Status", lResp[lRespIdx].Status)
					lRespRec.AppStatus = lResp[lRespIdx].Status
					lRespRec.AppReason = lResp[lRespIdx].Reason
					lRespRec.Status = common.ErrorCode
					log.Println("PBPO03", lRespRec.AppStatus, lRespRec.AppReason)

					lError.ErrCode = "PBPO03"
					lError.ErrMsg = lResp[lRespIdx].Reason + "\n" + "Exchange Server Busy"
					return lRespRec, lResp, lError
				}
			}
		} else {
			log.Println("PBPO04", "Unable to proceed Application!!!")
			lError.ErrCode = "PBPO04"
			lError.ErrMsg = "Unable to proceed Application!!!"

			return lRespRec, lResp, lError
		}
	}
	log.Println("BsePlaceOrder (-)")
	return lRespRec, lResp, lError
}

/*
Pupose:This method is used to process the application request to exchange
Parameters:

	pExchangeReq,pMasterId,PClientId

Response:

	==========
	*On Sucess
	==========
	2,[1,2,3],[1,2,3],nil

	==========
	*On Error
	==========
	0,[],[],error

Author:Pavithra
Date: 12JUNE2023
*/
//added by naveen:add one parameter source to insert in bidtracking table
//func ProcessBseReq(pExchangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int)
func ProcessBseReq(pExchangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int, pSource string) ([]nseipo.ExchangeRespStruct, error) {
	log.Println("ProcessBseReq (+)")

	//changing the structure nse to bse
	lBseExchangeReq := BseReqConstruct(pExchangeReq)
	// for retruning value for this method
	var lRespArr []nseipo.ExchangeRespStruct
	//added by naveen:add one parameter source to insert in bidtracking table
	//lBidTrackId, lErr1 := InsertBidTrack(pExchangeReq, pClientId, common.BSE, pBrokerId)
	lBidTrackId, lErr1 := InsertBidTrack(pExchangeReq, pClientId, common.BSE, pBrokerId, pSource)
	if lErr1 != nil {
		log.Println("PPBR01", lErr1)
		return lRespArr, lErr1
	} else {
		log.Println("bidtrackid-------------------------------", lBidTrackId)
		// ---------------------------------------------------------
		// Call the ApplyIpo method to Process the Application Request.
		lResponse, lErr2 := exchangecall.ApplyBseIpo(lBseExchangeReq, pClientId, pBrokerId)
		if lErr2 != nil {
			log.Println("PPBR02", lErr2.Error())
			return lRespArr, lErr2
		} else {
			// if lResponse.ScripId != "" {
			//chnaging the structure bse to nse
			lRespArr = BseRespConstruct(lResponse)
			// update the bid tracking table when application status success or failed
			lErr3 := UpdateBidTrack(lRespArr, pClientId, lBidTrackId, pExchangeReq, common.BSE, pBrokerId)
			if lErr3 != nil {
				log.Println("PPBR03", lErr3)
				return lRespArr, lErr3
			} else {
				for lRespIdx := 0; lRespIdx < len(lRespArr); lRespIdx++ {
					// check whether the application status is success or not
					if lRespArr[lRespIdx].Status == "success" {
						// if the status is success update the response in header and details table.
						lErr4 := InsertHeader(lRespArr, pExchangeReq, pReqRec, pClientId, common.BSE, pMailId, pBrokerId)
						if lErr4 != nil {
							log.Println("PPBR04", lErr4)
							return lRespArr, lErr4
						}
					}
				}
			}
			// }
		}
	}
	log.Println("ProcessBseReq (-)")
	return lRespArr, nil
}

func BseReqConstruct(pReqArr []nseipo.ExchangeReqStruct) bseipo.BseExchangeReqStruct {
	log.Println("BseReqConstruct (+)")
	log.Println("pReqArr", pReqArr)
	// var lNseReq nseipo.ExchangeReqStruct

	var lBsepReq bseipo.BseExchangeReqStruct
	var lBseval bseipo.BseRequestBidStruct

	//for loop for pReqArr

	for _, lNseReq := range pReqArr {

		lBsepReq.ScripId = lNseReq.Symbol
		lBsepReq.AccountNo_UpiId = lNseReq.Upi
		lBsepReq.ApplicationNo = lNseReq.ApplicationNo
		lBsepReq.Category = lNseReq.Category
		lBsepReq.ApplicantName = lNseReq.ClientName
		lBsepReq.Depository = lNseReq.Depository
		lBsepReq.DpId = "0"
		lBsepReq.ClientBenfId = lNseReq.ClientBenId
		lBsepReq.PanNo = lNseReq.Pan
		lBsepReq.ReferenceNo = lNseReq.ReferenceNo
		lBsepReq.IfscCode = lNseReq.IFSC
		lBsepReq.ChequeReceivedFlag = "N"
		lBsepReq.ChequeAmount = "0"
		lBsepReq.BankName = "8888"
		lBsepReq.Location = "upiidl"
		lBsepReq.Asba_UpiId = "1"

		// Initialize Bids slice for the current request
		// var lbids []bseipo.BseRequestBidStruct

		for lIdx := 0; lIdx < len(lNseReq.Bids); lIdx++ {

			if lNseReq.Bids[lIdx].ActivityType == "new" {
				lBseval.ActionCode = "N"
				lBseval.BidId = ""
			} else if lNseReq.Bids[lIdx].ActivityType == "modify" {
				lBseval.ActionCode = "M"
				lBseval.BidId = strconv.Itoa(lNseReq.Bids[lIdx].BidReferenceNo)
			} else {
				lBseval.ActionCode = "D"
				lBseval.BidId = strconv.Itoa(lNseReq.Bids[lIdx].BidReferenceNo)
			}

			lBseval.OrderNo = strconv.Itoa(lNseReq.Bids[lIdx].BidReferenceNo)

			lQuantity := strconv.Itoa(lNseReq.Bids[lIdx].Quantity)
			lBseval.Quantity = lQuantity

			lPrice := strconv.Itoa(lNseReq.Bids[lIdx].Price)
			lBseval.Rate = lPrice

			// lBseval.CuttOffFlag = strconv.FormatBool(lNseReq.Bids[lIdx].AtCutOff)

			if lNseReq.Bids[lIdx].AtCutOff == true {
				lBseval.CuttOffFlag = "1"
			} else {
				lBseval.CuttOffFlag = "0"
			}

			lNseReq.Bids[lIdx].Remark = lBseval.OrderNo
			lBsepReq.Bids = append(lBsepReq.Bids, lBseval)
		}
	}
	log.Println("lBsepReq", lBsepReq)
	log.Println("BseReqConstruct (+)")
	return lBsepReq

}

func BseRespConstruct(pRespRec bseipo.BseExchangeRespStruct) []nseipo.ExchangeRespStruct {
	log.Println("BseRespConstruct (+)")
	log.Println("pRespRec", pRespRec)
	// this variables is used to append the value of Response

	var lNseRespArr []nseipo.ExchangeRespStruct

	//response struct
	var lNseResp nseipo.ExchangeRespStruct

	var lNse nseipo.ResponseBidStruct

	currentTime := time.Now()
	// Format the time as a string using a custom format
	lFormattedTime := currentTime.Format("02-01-2006 15:04:05")

	lNseResp.Symbol = pRespRec.ScripId
	lNseResp.Upi = pRespRec.AccountNo_UpiId
	lNseResp.ApplicationNo = pRespRec.ApplicationNo
	lNseResp.Category = pRespRec.Category
	lNseResp.ClientName = pRespRec.ApplicantName
	lNseResp.Depository = pRespRec.Depository
	lNseResp.DpId = pRespRec.DpId
	lNseResp.ClientBenId = pRespRec.ClientBenfId
	lNseResp.Pan = pRespRec.PanNo
	lNseResp.ReferenceNo = pRespRec.ReferenceNo
	lNseResp.IFSC = pRespRec.IfscCode
	lNseResp.TimeStamp = lFormattedTime

	lErrCode, lErr1 := strconv.Atoi(pRespRec.ErrorCode)
	if lErr1 != nil {
		log.Println("Error:1", lErr1)
	}
	lNseResp.ReasonCode = lErrCode
	if lNseResp.ReasonCode == 0 {
		lNseResp.Status = common.SUCCESS
	} else {
		lNseResp.Status = common.FAILED
	}

	lNseResp.LocationCode = pRespRec.Location
	lNseResp.ChequeNo = pRespRec.ChequeReceivedFlag
	lNseResp.Status = pRespRec.StatusMessage

	if pRespRec.Asba_UpiId == "1" {
		lNseResp.NonASBA = true
	} else {
		lNseResp.NonASBA = false
	}

	for ldx := 0; ldx < len(pRespRec.Bids); ldx++ {

		if pRespRec.ErrorMessage != "" {
			lNseResp.Reason = pRespRec.ErrorMessage
		} else {
			lNseResp.Reason = pRespRec.Bids[ldx].Message
		}

		if pRespRec.Bids[ldx].ActionCode == "N" {
			lNse.ActivityType = "new"
		} else if pRespRec.Bids[ldx].ActionCode == "M" {
			lNse.ActivityType = "modify"
		} else {
			lNse.ActivityType = "cancel"
		}

		lBidId, lErr2 := strconv.Atoi(pRespRec.Bids[ldx].BidId)
		if lErr2 != nil {
			fmt.Println("Error:2", lErr2)
		} else {
			lNse.BidReferenceNo = lBidId
		}

		lQuantity, lErr3 := strconv.Atoi(pRespRec.Bids[ldx].Quantity)
		if lErr3 != nil {
			fmt.Println("Error:3", lErr3)
		} else {
			lNse.Quantity = lQuantity
		}

		lRate, lErr4 := common.ConvertStringToFloat(pRespRec.Bids[ldx].Rate)
		if lErr4 != nil {
			fmt.Println("Error:4", lErr3)
		} else {
			lNse.Price = lRate
		}

		lCuttoff, lErr5 := strconv.ParseBool(pRespRec.Bids[ldx].CuttOffFlag)
		if lErr5 != nil {
			fmt.Println("Error:5", lErr5)
		} else {
			lNse.AtCutOff = lCuttoff
		}

		lErrorCode, lErr6 := strconv.Atoi(pRespRec.Bids[ldx].ErrorCode)
		if lErr6 != nil {
			fmt.Println("Error:6", lErr6)
		} else {
			lNse.ReasonCode = lErrorCode
		}

		if pRespRec.Bids[ldx].ErrorCode == "0" {
			lNse.Status = common.SUCCESS
		} else {
			lNse.Status = "failed"
		}
		lNse.Reason = pRespRec.Bids[ldx].Message
		lNse.Remark = pRespRec.Bids[ldx].OrderNo

		if pRespRec.StatusCode == "0" && pRespRec.Bids[ldx].ErrorCode == "0" {
			lNseResp.Status = common.SUCCESS
		} else {
			lNseResp.Status = "failed"
		}
		lNseResp.Bids = append(lNseResp.Bids, lNse)
	}
	for ldx := 0; ldx < len(pRespRec.Bids); ldx++ {
		if lNse.ActivityType == "cancel" {
			lNseResp.Bids = []nseipo.ResponseBidStruct{}
		}
	}
	lNseRespArr = append(lNseRespArr, lNseResp)
	log.Println("lNseRespArr", lNseRespArr)

	log.Println("BseRespConstruct (-)")
	return lNseRespArr

}
