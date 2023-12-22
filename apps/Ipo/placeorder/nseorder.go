package placeorder

import (
	"fcs23pkg/common"
	"fcs23pkg/integration/nse/nseipo"
	"log"
)

type ErrorStruct struct {
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

//added by naveen:add one parameter source to insert in bidtracking table
//func NsePlaceOrder(pExhangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int) (OrderResStruct, []nseipo.ExchangeRespStruct, ErrorStruct) {
func NsePlaceOrder(pExhangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int, pSource string) (OrderResStruct, []nseipo.ExchangeRespStruct, ErrorStruct) {
	log.Println("NsePlaceOrder (+)")

	var lRespRec OrderResStruct
	var lError ErrorStruct

	//added by naveen:add one argument source to insert in bidtracking
	//lResp, lErr9 := ProcessNseReq(pExhangeReq, pReqRec, pClientId, pMailId, pBrokerId)
	lResp, lErr9 := ProcessNseReq(pExhangeReq, pReqRec, pClientId, pMailId, pBrokerId, pSource)
	if lErr9 != nil {
		log.Println("PNPO01", lErr9.Error())
		lRespRec.Status = common.ErrorCode
		lError.ErrCode = "PNPO01"
		lError.ErrMsg = "Exchange Server is Busy right now,Try After Sometime."
		return lRespRec, lResp, lError
	} else {
		if lResp != nil {
			for lRespIdx := 0; lRespIdx < len(lResp); lRespIdx++ {
				if lResp[lRespIdx].Status == "failed" {
					lRespRec.AppStatus = lResp[lRespIdx].Status
					lRespRec.AppReason = lResp[lRespIdx].Reason
					lRespRec.Status = common.ErrorCode
					lError.ErrCode = "PNPO02"
					lError.ErrMsg = lResp[lRespIdx].Reason + "\n" + "Application Failed"
					log.Println("PNPO02", lRespRec.AppStatus, lRespRec.AppReason)
					return lRespRec, lResp, lError
				} else if lResp[lRespIdx].Status == "success" {
					lRespRec.AppStatus = lResp[lRespIdx].Status
					lRespRec.AppReason = lResp[lRespIdx].Reason
					lRespRec.Status = common.SuccessCode
				}
			}
		} else {
			log.Println("PNPO03", "Unable to proceed Application!!!")
			lError.ErrCode = "PNPO03"
			lError.ErrMsg = "Unable to proceed Application!!!"
			return lRespRec, lResp, lError
		}
	}
	log.Println("NsePlaceOrder (-)")
	return lRespRec, lResp, lError
}
