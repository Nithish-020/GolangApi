package ncbschedule

import (
	"encoding/json"
	"fcs23pkg/apigate"
	"fcs23pkg/apps/Ncb/ncbplaceorder"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/clientfund"
	"fcs23pkg/integration/nse/nsencb"
	"fcs23pkg/integration/techexcel"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type NcbStruct struct {
	MasterId int
	Symbol   string
	Exchange string
}

// this struct is used to get the Frontoffice JV details.
type FOJvStruct struct {
	UserId       string `json:"uid"`
	AccountId    string `json:"actid"`
	Amount       string `json:"amt"`
	SourceUserId string `json:"src_uid"`
	Remarks      string `json:"remarks"`
}

// ------------------------------------------------------
// this method is used to filter the Jv posted records
// ------------------------------------------------------

func NcbPostJvandProcessOrder(pNcbReqRec exchangecall.NcbJvReqStruct, r *http.Request, pValidNcb NcbStruct, pBrokerId int) (exchangecall.NcbJvReqStruct, error) {

	log.Println("NcbPostJvandProcessOrder (+)")

	lClientFund, lErr1 := ncbplaceorder.VerifyFundDetails(pNcbReqRec.ClientId)
	log.Println("lClientFund", lClientFund)
	if lErr1 != nil || lClientFund.Status == common.ErrorCode {
		log.Println("NPJPO01", lErr1)
		pNcbReqRec.EmailDist = "4"
		pNcbReqRec.ProcessStatus = common.ErrorCode
		pNcbReqRec.ErrorStage = "V"
		pNcbReqRec.FoVerifyFailedCount++
		return pNcbReqRec, lErr1
	} else if lClientFund.Status == common.SuccessCode {
		pNcbReqRec.FoVerifySuccessCount++

		//----------------------------------------------------------------
		// Fund Balance Verification Success
		//----------------------------------------------------------------

		lPrice := int(pNcbReqRec.Price)
		lAmount := strconv.FormatFloat(pNcbReqRec.Amount, 'f', -1, 64)

		if lClientFund.AccountBalance < float64(pNcbReqRec.Unit*lPrice) {

			// ========== INSUFFICIENT ACC BALANCE ==============

			if lClientFund.AccountBalance < pNcbReqRec.Price {
				pNcbReqRec.EmailDist = "0"
				pNcbReqRec.ProcessStatus = common.ErrorCode
				pNcbReqRec.ErrorStage = "I"
				pNcbReqRec.InsufficientFund++
				return pNcbReqRec, nil
			} else {

				// ======== INSUFFICIENT AMOUNT CHANGING UNIT ==============

				//changing the unit as per the SI text
				MaximumUnits := lClientFund.AccountBalance / pNcbReqRec.Price

				pNcbReqRec.Unit = int(math.Floor(MaximumUnits))
				lAmount = strconv.Itoa(lPrice * pNcbReqRec.Unit)

				log.Println("lAmount", lAmount)

				FoStatus, lErr2 := NcbBlockFOClientFund(pNcbReqRec, "", "D")
				log.Println("FoStatus", FoStatus)
				pNcbReqRec.FoJvAmount = FoStatus.FoJvAmount
				pNcbReqRec.FoJvStatement = FoStatus.FoJvStatement
				pNcbReqRec.FoJvStatus = FoStatus.FoJvStatus
				pNcbReqRec.FoJvType = FoStatus.FoJvType

				// while processing FO Jv when we getting error it returns error and send mail to IT,RMS
				if lErr2 != nil || FoStatus.FoJvStatus == common.ErrorCode {
					log.Println("NPJPO02", lErr2)
					pNcbReqRec.EmailDist = "4"
					pNcbReqRec.ProcessStatus = common.ErrorCode
					pNcbReqRec.ErrorStage = "F"
					pNcbReqRec.FoJvFailedCount++
					return pNcbReqRec, lErr2
				} else if FoStatus.FoJvStatus == common.SuccessCode {
					pNcbReqRec.FoJvSuccessCount++

					//BackOffice Jv Processing
					BoStatus, lErr3 := NcbBlockBOClientFund(pNcbReqRec, r, "C")
					log.Println("BoStatus", BoStatus)
					pNcbReqRec.BoJvAmount = BoStatus.BoJvAmount
					pNcbReqRec.BoJvStatement = BoStatus.BoJvStatement
					pNcbReqRec.BoJvStatus = BoStatus.BoJvStatus
					pNcbReqRec.BoJvType = BoStatus.BoJvType
					if lErr3 != nil || BoStatus.BoJvStatus == common.ErrorCode {
						pNcbReqRec.BoJvFailedCount++
						log.Println("NPJPO03", lErr3)

						RevFoStatus, lErr4 := NcbBlockFOClientFund(pNcbReqRec, "R", "D")

						if lErr4 != nil || RevFoStatus.FoJvStatus == common.ErrorCode {
							pNcbReqRec.ReverseFoJvFailedCount++
							log.Println("NPJPO04", lErr4)
							pNcbReqRec.EmailDist = "3"
							pNcbReqRec.ProcessStatus = common.ErrorCode
							pNcbReqRec.ErrorStage = "A"
							return pNcbReqRec, lErr4
						} else if RevFoStatus.FoJvStatus == common.SuccessCode {
							pNcbReqRec.ReverseFoJvSuccessCount++
							pNcbReqRec.EmailDist = "2"
							pNcbReqRec.ProcessStatus = common.ErrorCode
							pNcbReqRec.ErrorStage = "B"
							return pNcbReqRec, lErr3
						}

					} else if BoStatus.BoJvStatus == common.SuccessCode {
						pNcbReqRec.BoJvSuccesssCount++

						//--------------------------
						//exchange calling place
						//-------------------------

						lRespRec, lErr5 := NcbExchangeProcess(pNcbReqRec, pValidNcb.Exchange, pBrokerId, r)

						pNcbReqRec = lRespRec
						if lErr5 != nil || lRespRec.ProcessStatus == common.ErrorCode {

							pNcbReqRec.TotalFailedCount++
							log.Println("NPJPO05", lErr5)

							RevBoStatus, lErr6 := NcbBlockBOClientFund(pNcbReqRec, r, "R")

							if lErr6 != nil || RevBoStatus.BoJvStatus == common.ErrorCode {
								pNcbReqRec.ReverseBoJvFailedCount++
								log.Println("NPJPO06", lErr6)
								pNcbReqRec.EmailDist = "3"
								pNcbReqRec.ProcessStatus = common.ErrorCode
								pNcbReqRec.ErrorStage = "R"
								return pNcbReqRec, lErr6

							} else if BoStatus.BoJvStatus == common.SuccessCode {
								pNcbReqRec.ReverseBoJvSuccessCount++
								RevFoStatus, lErr7 := NcbBlockFOClientFund(pNcbReqRec, "R", "C")
								if lErr7 != nil || RevFoStatus.FoJvStatus == common.ErrorCode {
									pNcbReqRec.ReverseFoJvFailedCount++
									log.Println("NPJPO07", lErr7)
									pNcbReqRec.EmailDist = "2"
									pNcbReqRec.ProcessStatus = common.ErrorCode
									pNcbReqRec.ErrorStage = "H"
									return pNcbReqRec, lErr7
								} else if RevFoStatus.FoJvStatus == common.SuccessCode {
									pNcbReqRec.ReverseFoJvSuccessCount++
									pNcbReqRec.ExchangeFailedCount++
									pNcbReqRec.EmailDist = lRespRec.EmailDist
									pNcbReqRec.ProcessStatus = common.ErrorCode
									pNcbReqRec.ErrorStage = lRespRec.ErrorStage
									return pNcbReqRec, lErr6
								}
							}

						} else {
							pNcbReqRec.TotalSuccessCount++
							pNcbReqRec.ExchangeSuccessCount++
							return pNcbReqRec, nil
						}

					}

				}

			}

			// ========= INSUFFICIENT AMOUNT CHANGING UNIT ==========

		} else {
			// ======== BACK OFFICE JV SUFFICIENT AMOUNT =========

			FoStatus, lErr8 := NcbBlockFOClientFund(pNcbReqRec, "", "D")
			log.Println("FoStatus", FoStatus)
			pNcbReqRec.FoJvAmount = FoStatus.FoJvAmount
			pNcbReqRec.FoJvStatement = FoStatus.FoJvStatement
			pNcbReqRec.FoJvStatus = FoStatus.FoJvStatus
			pNcbReqRec.FoJvType = FoStatus.FoJvType

			// while processing FO Jv when we getting error it returns error and send mail to IT,RMS
			if lErr8 != nil || FoStatus.FoJvStatus == common.ErrorCode {
				log.Println("NPJPO08", lErr8)
				pNcbReqRec.EmailDist = "4"
				pNcbReqRec.ProcessStatus = common.ErrorCode
				pNcbReqRec.ErrorStage = "F"
				pNcbReqRec.FoJvFailedCount++
				return pNcbReqRec, lErr8

			} else if FoStatus.FoJvStatus == common.SuccessCode {

				pNcbReqRec.FoJvSuccessCount++

				//BackOffice Jv Processing
				BoStatus, lErr9 := NcbBlockBOClientFund(pNcbReqRec, r, "C")
				pNcbReqRec.BoJvAmount = BoStatus.BoJvAmount
				pNcbReqRec.BoJvStatement = BoStatus.BoJvStatement
				pNcbReqRec.BoJvStatus = BoStatus.BoJvStatus
				pNcbReqRec.BoJvType = BoStatus.BoJvType

				if lErr9 != nil || BoStatus.BoJvStatus == common.ErrorCode {
					pNcbReqRec.BoJvFailedCount++
					log.Println("NPJPO09", lErr9)

					RevFoStatus, lErr10 := NcbBlockFOClientFund(pNcbReqRec, "R", "D")
					if lErr10 != nil || RevFoStatus.FoJvStatus != common.ErrorCode {
						pNcbReqRec.ReverseFoJvFailedCount++
						log.Println("NPJPO10", lErr10)
						pNcbReqRec.EmailDist = "3"
						pNcbReqRec.ProcessStatus = common.ErrorCode
						pNcbReqRec.ErrorStage = "A"
						return pNcbReqRec, lErr10
					} else if RevFoStatus.FoJvStatus == common.SuccessCode {
						pNcbReqRec.ReverseFoJvSuccessCount++
						pNcbReqRec.EmailDist = "2"
						pNcbReqRec.ProcessStatus = common.ErrorCode
						pNcbReqRec.ErrorStage = "B"
						return pNcbReqRec, lErr9
					}

				} else if BoStatus.BoJvStatus == common.SuccessCode {

					pNcbReqRec.BoJvSuccesssCount++
					//-------------------------
					//exchange calling place
					//-------------------------

					lRespRec, lErr11 := NcbExchangeProcess(pNcbReqRec, pValidNcb.Exchange, pBrokerId, r)
					pNcbReqRec = lRespRec

					if lErr11 != nil || lRespRec.ProcessStatus == common.ErrorCode {

						pNcbReqRec.TotalFailedCount++

						log.Println("NPJPO11", lErr11)
						RevBoStatus, lErr12 := NcbBlockBOClientFund(pNcbReqRec, r, "R")

						if lErr12 != nil || RevBoStatus.BoJvStatus == common.ErrorCode {
							pNcbReqRec.ReverseBoJvFailedCount++
							log.Println("NPJPO12", lErr12)
							pNcbReqRec.EmailDist = "3"
							pNcbReqRec.ProcessStatus = common.ErrorCode
							pNcbReqRec.ErrorStage = "R"
							return pNcbReqRec, lErr12
						} else if BoStatus.BoJvStatus == common.SuccessCode {
							pNcbReqRec.ReverseBoJvSuccessCount++
							RevFoStatus, lErr13 := NcbBlockFOClientFund(pNcbReqRec, "R", "C")
							if lErr13 != nil || RevFoStatus.FoJvStatus == common.ErrorCode {
								log.Println("NPJPO13", lErr13)
								pNcbReqRec.ReverseFoJvFailedCount++
								pNcbReqRec.EmailDist = "2"
								pNcbReqRec.ProcessStatus = common.ErrorCode
								pNcbReqRec.ErrorStage = "H"
								return pNcbReqRec, lErr13
							} else if RevFoStatus.FoJvStatus == common.SuccessCode {
								pNcbReqRec.ReverseFoJvSuccessCount++
								pNcbReqRec.ExchangeFailedCount++
								pNcbReqRec.TotalFailedCount++
								pNcbReqRec.EmailDist = lRespRec.EmailDist
								pNcbReqRec.ProcessStatus = common.ErrorCode
								pNcbReqRec.ErrorStage = lRespRec.ErrorStage
								return pNcbReqRec, lErr12
							}
						}

					} else {
						pNcbReqRec.TotalSuccessCount++
						pNcbReqRec.ExchangeSuccessCount++
						return pNcbReqRec, nil
					}

				}

			}

		}
		// ============= BACK OFFICE JV SUFFICIENT AMOUNT ==============

	} else {

		// ========while verifying the fund getting error ============
		pNcbReqRec.EmailDist = "4"
		pNcbReqRec.ProcessStatus = common.ErrorCode
		pNcbReqRec.ErrorStage = "V"
		pNcbReqRec.FoVerifyFailedCount++
		log.Println("VerifyFundDetails", pNcbReqRec.FoVerifyFailedCount)
		return pNcbReqRec, nil

	}

	log.Println("NcbPostJvandProcessOrder (-)")
	return pNcbReqRec, nil

}

/*
Pupose:This method is used to get the pan number for order input.
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	AGMPA45767,nil

	==========
	*On Error
	==========
	"",error

Author:KAVYADHARSHANI
Date: 05JAN2024
*/
func NcbBlockBOClientFund(pJvDetailRec exchangecall.NcbJvReqStruct, pRequest *http.Request, pFlag string) (NcbJvStatusStruct, error) {
	log.Println("NcbBlockBOClientFund (+)")

	// this variables is used to get Status
	var lJVReq techexcel.JvInputStruct
	var lReqDtl apigate.RequestorDetails
	var lJvStatusRec NcbJvStatusStruct

	lReqDtl = apigate.GetRequestorDetail(pRequest)
	lConfigFile := common.ReadTomlConfig("./toml/techXLAPI_UAT.toml")

	lCocd, lErr1 := ncbplaceorder.NcbPaymentCode(pJvDetailRec.ClientId)
	if lErr1 != nil {
		log.Println("NBBCF01", lErr1.Error())
		lJvStatusRec.BoJvStatus = common.ErrorCode
		return lJvStatusRec, lErr1
	} else {

		lOrderNo := strconv.Itoa(pJvDetailRec.ReqOrderNo)
		lAmount := strconv.FormatFloat(pJvDetailRec.Amount, 'f', -1, 64)

		lJVReq.COCD = lCocd
		lJVReq.VoucherDate = time.Now().Format("02/01/2006")
		lJVReq.BillNo = "NCB" + pJvDetailRec.ClientId
		lJVReq.SourceTable = "a_ncb_orderdetails"
		lJVReq.SourceTableKey = lOrderNo
		lJVReq.Amount = lAmount
		lJVReq.WithGST = "N"

		lUnit := strconv.Itoa(pJvDetailRec.Unit)
		lPrice := strconv.FormatFloat(pJvDetailRec.Price, 'f', -1, 64)

		lFTCaccount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NseNCBorderAccount"])

		if pFlag == "R" {
			lJVReq.AccountCode = lFTCaccount
			lJVReq.CounterAccount = pJvDetailRec.ClientId
			lJVReq.Narration = "Failed of NCB Order Refund  " + lUnit + " * " + lPrice + " on " + lJVReq.VoucherDate + " - " + pJvDetailRec.ClientId
			lJvStatusRec.BoJvStatement = lJVReq.Narration
			lJvStatusRec.BoJvType = "C"
		} else if pFlag == "C" {
			lJVReq.AccountCode = pJvDetailRec.ClientId
			lJVReq.CounterAccount = lFTCaccount
			lJVReq.Narration = "NCB Purchase  " + lUnit + " * " + lPrice + " on " + lJVReq.VoucherDate + " - " + pJvDetailRec.ClientId
			lJvStatusRec.BoJvStatement = lJVReq.Narration
			lJvStatusRec.BoJvType = "D"
		}

		lJvStatusRec.BoJvStatus = "S"
		lJvStatusRec.BoJvAmount = lAmount

		//JV processing method
		lErr2 := clientfund.BOProcessJV(lJVReq, lReqDtl)
		if lErr2 != nil {
			log.Println("NBBCF02", lErr2.Error())
			lJvStatusRec.BoJvStatus = common.ErrorCode
			return lJvStatusRec, lErr2
		} else {
			lJvStatusRec.BoJvStatus = common.SuccessCode
		}
	}

	log.Println("BlockClientFund (-)")
	return lJvStatusRec, nil
}

/*
Purpose: This function is used to construct the details for the BO Jv API

	Parameters:
		{
			JvAmount    string `json:"jvAmount"`
			JvStatus    string `json:"jvStatus"`
			JvStatement string `json:"jvStatement"`
			JvType      string `json:"jvType"`
		}, "FT000069"

	Response:

	On Success :
	===============

	===============

	On Error:
	===============
		{},error
	===============

Author:KAVYA DHARSHANI
Date: 05JAN2024
*/

func NcbBlockFOClientFund(pRequest exchangecall.NcbJvReqStruct, pStatusFlag string, pPaymentType string) (NcbJvStatusStruct, error) {
	log.Println("NcbBlockFOClientFund (+)")

	var lFOJvReqStruct FOJvStruct
	var lJvStatusRec NcbJvStatusStruct
	var lUrl string

	config := common.ReadTomlConfig("./toml/techXLAPI_UAT.toml")
	lFOJvReqStruct.UserId = fmt.Sprintf("%v", config.(map[string]interface{})["VerifyUser"])
	lFOJvReqStruct.SourceUserId = lFOJvReqStruct.UserId

	lToken, lErr1 := ncbplaceorder.NcbGetFOToken()
	if lErr1 != nil {
		log.Println("NBFCF01", lErr1.Error())
		lJvStatusRec.FoJvStatus = common.ErrorCode
		return lJvStatusRec, lErr1
	} else {

		lPrice := int(pRequest.Price)

		if pPaymentType == "D" {
			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayoutUrl"])
			lJvStatusRec.FoJvType = pRequest.BoJvType
			if pStatusFlag != "R" {
				lFOJvReqStruct.Remarks = "AMOUNT HOLD FOR NCB ORDER"
				lJvStatusRec.FoJvStatement = lFOJvReqStruct.Remarks
				lFOJvReqStruct.Amount = "-" + strconv.Itoa(lPrice*pRequest.Unit)
			} else {
				lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM NCB ORDER"
				lJvStatusRec.FoJvStatement = lFOJvReqStruct.Remarks
				lFOJvReqStruct.Amount = strconv.Itoa(lPrice * pRequest.Unit)
			}
			// } else if pRequest.BoJvType == "C" {
		} else if pPaymentType == "C" {
			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayinUrl"])
			lJvStatusRec.FoJvType = pRequest.BoJvType
			lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM NCB ORDER"
			lFOJvReqStruct.Amount = strconv.Itoa(lPrice * pRequest.Unit)
		}
		lFOJvReqStruct.AccountId = pRequest.ClientId
		lJvStatusRec.FoJvAmount = lFOJvReqStruct.Amount

		lRequest, lErr2 := json.Marshal(lFOJvReqStruct)
		if lErr2 != nil {
			log.Println("NBFCF02", lErr2.Error())
			lJvStatusRec.FoJvStatus = common.ErrorCode
			return lJvStatusRec, lErr2
		} else {
			// construct the request body
			lBody := `jData=` + string(lRequest) + `&jKey=` + lToken
			lResp, lErr3 := clientfund.FOProcessJV(lUrl, lBody)
			if lErr3 != nil {
				log.Println("NBFCF03", lErr3.Error())
				lJvStatusRec.FoJvStatus = common.ErrorCode
				return lJvStatusRec, lErr3
			} else {
				if lResp.Status == "Ok" {
					log.Println("FO JV processed successfully for : ", lRequest)
					lJvStatusRec.FoJvStatus = common.SuccessCode
				} else {
					log.Println("FO JV processed Failed for : ", lRequest)
					lJvStatusRec.FoJvStatus = common.ErrorCode
				}
			}
		}
	}

	log.Println("NcbBlockFOClientFund (-)")
	return lJvStatusRec, nil
}

/*
Pupose:This method inserting the order head values in order header table.
Parameters:

	pReqArr,pMasterId,PClientId

Response:

	==========
	*On Sucess
	==========


	==========
	*On Error
	==========

Author:KAVYA DHARSHANI
Date: 05JAN2024
*/
func UpdateHeader(pNseRespRec nsencb.NcbAddResStruct, pJvRec exchangecall.NcbJvReqStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateHeader (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NUH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lHeaderId, lErr2 := GetOrderId(pJvRec, pExchange)
		if lErr2 != nil {
			log.Println("NUH02", lErr2)
			return lErr2
		} else {
			if lHeaderId != 0 {
				if pExchange == common.BSE {
					log.Println("pExchange", pExchange)
				} else if pExchange == common.NSE {

					log.Println("pExchange", pExchange)

					lSqlString := `update a_ncb_orderheader h
					               set  h.Symbol = ?, h.pan  = ?,h.RespInvestmentunit =?,  h.depository  =?, h.dpId =? ,h.RespapplicationNo=?, 
					                    h.clientBenId  =?, h.clientRefNumber =? ,h.status = ?,  h.StatusMessage =?,
					                    h.ErrorCode  =?, h.ErrorMessage =?, h.lastActionTime = ?, h.UpdatedBy = ?,
					                    h.UpdatedDate = now()
				                  where h.clientId = ?
				                  and h.brokerId  = ?
					              and h.Id = ?`

					_, lErr3 := lDb.Exec(lSqlString, pNseRespRec.Symbol, pNseRespRec.Pan, pNseRespRec.InvestmentValue, pNseRespRec.Depository, pNseRespRec.DpId, pNseRespRec.ApplicationNumber, pNseRespRec.ClientBenId, pNseRespRec.ClientRefNumber, pNseRespRec.Status, pNseRespRec.Reason, pNseRespRec.Status, pNseRespRec.Reason, pNseRespRec.LastActionTime, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
					if lErr3 != nil {
						log.Println("NUH03", lErr3)
						return lErr3
					} else {

						lErr4 := UpdateDetail(pJvRec, pNseRespRec, lHeaderId)
						if lErr4 != nil {
							log.Println("NUH04", lErr4)
							return lErr4
						} else {
							log.Println("Header Update SuccessFully")
						}

					}

				}

			}
		}

	}

	log.Println("UpdateHeader (-)")
	return nil
}

/*
Pupose:This method inserting the order head values in order header table.
Parameters:

	pReqArr,pMasterId,PClientId

Response:

	==========
	*On Sucess
	==========


	==========
	*On Error
	==========

Author:KAVYA DHARSHANI
Date: 05JAN2024
*/

func UpdateDetail(pJvRec exchangecall.NcbJvReqStruct, pNseRespRec nsencb.NcbAddResStruct, pHeaderId int) error {
	log.Println("UpdateDetail (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NUD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lSqlString := ` update a_ncb_orderdetails d
		                set d.RespOrderNo =?, d.RespUnit = ?, 
		                    d.Respprice = ?,d.RespAmount =?,d.ErrorCode =?, d.ErrorMessage = ?, 
		                    d.BoJvStatus = ?,d.BoJvAmount = ?,d.BoJvStatement = ?,d.BoJvType = ?,
		                    d.FoJvStatus = ?,d.FoJvAmount = ?,d.FoJvStatement = ?,d.FoJvType = ?, d.UpdatedBy = ?,d.UpdatedDate = now()
	                   where d.headerId = ?
	                   and d.ReqOrderNo = ?`

		_, lErr2 := lDb.Exec(lSqlString, pNseRespRec.OrderNumber, pNseRespRec.InvestmentValue, pNseRespRec.Price, pNseRespRec.TotalAmountPayable, pNseRespRec.Status, pNseRespRec.Reason, pJvRec.BoJvStatus, pJvRec.BoJvAmount, pJvRec.BoJvStatement, pJvRec.BoJvType, pJvRec.FoJvStatus, pJvRec.FoJvAmount, pJvRec.FoJvStatement, pJvRec.FoJvType, common.AUTOBOT, pHeaderId, pJvRec.ReqOrderNo)

		if lErr2 != nil {
			log.Println("NUD02", lErr2)
			return lErr2
		} else {
			// call InsertDetails method to inserting the order details in order details table
			log.Println("Details Updated SuccessFullly")
		}

	}
	log.Println("UpdateDetail (-)")
	return nil
}

func GetOrderId(pJvRec exchangecall.NcbJvReqStruct, pExchange string) (int, error) {
	log.Println("GetOrderId (+)")

	var lHeaderId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NGOI01", lErr1)
		return lHeaderId, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select (case when count(1) > 0 then h.Id  else 0 end) Id
	 	                 from a_ncb_orderheader h, a_ncb_orderdetails d 
		                 where h.Id = d.HeaderId 
		                 and h.clientId  = ?
		                 and d.ReqOrderNo = ?
		                 and h.Exchange  = ?
		                 and h.MasterId = (		
			                    select n.id
			                    from a_ncb_master n
			                    where n.symbol = ? )`

		lRows, lErr2 := lDb.Query(lCoreString, pJvRec.ClientId, pJvRec.ReqOrderNo, pExchange, pJvRec.Symbol)

		if lErr2 != nil {
			log.Println("NGOI02", lErr2)
			return lHeaderId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lHeaderId)
				if lErr3 != nil {
					log.Println("NGOI03", lErr3)
					return lHeaderId, lErr3
				}
				log.Println(lHeaderId)
			}
		}
	}
	log.Println("GetOrderId (-)")
	return lHeaderId, nil

}

func NcbExchangeProcess(pNcbReqRec exchangecall.NcbJvReqStruct, pExchange string, pBrokerId int, r *http.Request) (exchangecall.NcbJvReqStruct, error) {
	log.Println("NcbExchangeProcess (+)")

	if pExchange == common.BSE {
		log.Println("pExchange", pExchange)
	} else if pExchange == common.NSE {

		lNseReqRec := exchangecall.NcbConstructExchReq(pNcbReqRec, pExchange)

		lRespRec, lErr1 := exchangecall.ApplyNseNcb(lNseReqRec, common.AUTOBOT, pBrokerId)

		pNcbReqRec = exchangecall.NcbConstructExchResp(lRespRec, pExchange)

		if lErr1 != nil || lRespRec.Status != "success" {
			// ================================ EXCHANGE FAILED ====================================
			log.Println("NEP01", lErr1)
			pNcbReqRec.EmailDist = "1"
			pNcbReqRec.ProcessStatus = common.ErrorCode
			pNcbReqRec.ErrorStage = "E"
			return pNcbReqRec, lErr1
		} else {
			pNcbReqRec.EmailDist = "0"
			pNcbReqRec.ProcessStatus = "Y"
			pNcbReqRec.ErrorStage = "Y"
			return pNcbReqRec, nil
		}
	}

	log.Println("ExchangeProcess (-)")
	return pNcbReqRec, nil
}
