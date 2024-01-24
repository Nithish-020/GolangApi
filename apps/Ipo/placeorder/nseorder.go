package placeorder

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
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

func Ipo_EligibleToOrder(pRequest OrderReqStruct, pBrokerId int, pClientId string) (bool, error) {
	log.Println("Ipo_EligibleToOrder (+)")

	var lAllowToOrder bool
	lString := ""
	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PIETO01", lErr1)
		return lAllowToOrder, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select (case when count(1) > 0 then 'Y' else 'N' end)AllowFlag
						from a_ipo_order_header h
						where h.cancelFlag = 'N'
						and h.status  = 'success'
						and h.MasterId = ?
						and h.clientId = ?
						and h.category = ?
						and h.brokerId = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pRequest.MasterId, pClientId, pRequest.Category, pBrokerId)
		if lErr2 != nil {
			log.Println("PIETO02", lErr2)
			return lAllowToOrder, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lString)
				if lErr3 != nil {
					log.Println("PIETO03", lErr3)
					return lAllowToOrder, lErr3
				} else {
					for lIdx := 0; lIdx < len(pRequest.BidDetails); lIdx++ {
						if pRequest.BidDetails[lIdx].ActivityType == "new" {
							if lString == "N" {
								lAllowToOrder = true
							} else if lString == "Y" {
								lAllowToOrder = false
							}
						} else if pRequest.BidDetails[lIdx].ActivityType == "modify" || pRequest.BidDetails[lIdx].ActivityType == "cancel" {
							if lString == "Y" {
								lAllowToOrder = true
							} else {
								lAllowToOrder = false
							}
						}
					}
				}
			}
		}
	}

	log.Println("Ipo_EligibleToOrder (-)")
	return lAllowToOrder, nil
}
