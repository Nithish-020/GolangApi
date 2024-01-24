package sgbschedule

import (
	"encoding/json"
	"fcs23pkg/apigate"
	"fcs23pkg/apps/SGB/sgbplaceorder"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bsesgb"
	"fcs23pkg/integration/clientfund"
	"fcs23pkg/integration/nse/nsesgb"
	"fcs23pkg/integration/techexcel"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type SgbClientDetail struct {
	OrderNo    string `json:"orderno"`
	Unit       string `json:"unit"`
	Price      string `json:"price"`
	Symbol     string `json:"symbol"`
	OrderDate  string `json:"orderdate"`
	Amount     int    `json:"amount"`
	Mail       string `json:"mail"`
	ClientName string `json:"clientname"`
	Activity   string `json:"activity"`
}

type JvDataStruct struct {
	Unit        string `json:"unit"`
	Price       string `json:"price"`
	JVamount    string `json:"jvAmount"`
	OrderNo     string `json:"orderNo"`
	ClientId    string `json:"clientId"`
	ActionCode  string `json:"actionCode"`
	Transaction string `json:"transaction"`
	BidId       string `json:"bidId"`
}

type DynamicEmailStruct struct {
	Date        string
	FailedJvArr []FailedJvStruct
}

type SgbStruct struct {
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

// ---------------------------------------------------------------------------------
// this method is used to filter the Jv posted records
// ---------------------------------------------------------------------------------
func PostJvandProcessOrder(pSgbReqRec exchangecall.JvReqStruct, r *http.Request, pValidSgb SgbStruct, pBrokerId int) (exchangecall.JvReqStruct, error) {
	log.Println("PostJvandProcessOrder (+)")

	lClientFund, lErr1 := sgbplaceorder.VerifyFOFundDetails(pSgbReqRec.ClientId)
	log.Println("lClientFund", lClientFund)
	if lErr1 != nil || lClientFund.Status == common.ErrorCode {
		log.Println("SPJPO01", lErr1)
		// ====================================while verifying the fund getting error================================
		pSgbReqRec.EmailDist = "4"
		pSgbReqRec.ProcessStatus = common.ErrorCode
		pSgbReqRec.ErrorStage = "V"
		pSgbReqRec.FoVerifyFailedCount++
		return pSgbReqRec, lErr1
	} else if lClientFund.Status == common.SuccessCode {
		pSgbReqRec.FoVerifySuccessCount++

		//----------------------------------------------------------------
		// Fund Balance Verification Success
		//----------------------------------------------------------------

		if lClientFund.AccountBalance < float64(pSgbReqRec.Unit*pSgbReqRec.Price) {
			// ==================INSUFFICIENT ACC BALANCE============================
			if lClientFund.AccountBalance < float64(pSgbReqRec.Price) {
				pSgbReqRec.EmailDist = "0"
				pSgbReqRec.ProcessStatus = common.ErrorCode
				pSgbReqRec.ErrorStage = "I"
				pSgbReqRec.InsufficientFund++
				return pSgbReqRec, nil
			} else {
				// =====================INSUFFICIENT AMOUNT CHANGING UNIT======================

				//changing the unit as per the SI text
				MaximumUnits := lClientFund.AccountBalance / float64(pSgbReqRec.Price)
				// Convert float to int using Ceil (round up)
				pSgbReqRec.Unit = int(math.Floor(MaximumUnits))
				pSgbReqRec.Amount = strconv.Itoa(pSgbReqRec.Price * pSgbReqRec.Unit)

				// Convert the integer to string
				for lBidIdx := 0; lBidIdx < len(pSgbReqRec.Bids); lBidIdx++ {
					pSgbReqRec.Bids[lBidIdx].SubscriptionUnit = strconv.Itoa(pSgbReqRec.Unit)
				}

				FoStatus, lErr2 := BlockFOClientFund(pSgbReqRec, "", "D")
				log.Println("FoStatus", FoStatus)
				pSgbReqRec.FoJvAmount = FoStatus.FoJvAmount
				pSgbReqRec.FoJvStatement = FoStatus.FoJvStatement
				pSgbReqRec.FoJvStatus = FoStatus.FoJvStatus
				pSgbReqRec.FoJvType = FoStatus.FoJvType
				// while processing FO Jv when we getting error it returns error and send mail to IT,RMS
				if lErr2 != nil || FoStatus.FoJvStatus == common.ErrorCode {
					log.Println("SPJPO02", lErr2)
					pSgbReqRec.EmailDist = "4"
					pSgbReqRec.ProcessStatus = common.ErrorCode
					pSgbReqRec.ErrorStage = "F"
					pSgbReqRec.FoJvFailedCount++
					return pSgbReqRec, lErr2
				} else if FoStatus.FoJvStatus == common.SuccessCode {
					pSgbReqRec.FoJvSuccessCount++

					//BackOffice Jv Processing
					BoStatus, lErr3 := BlockBOClientFund(pSgbReqRec, r, "C")
					log.Println("BoStatus", BoStatus)
					pSgbReqRec.BoJvAmount = BoStatus.BoJvAmount
					pSgbReqRec.BoJvStatement = BoStatus.BoJvStatement
					pSgbReqRec.BoJvStatus = BoStatus.BoJvStatus
					pSgbReqRec.BoJvType = BoStatus.BoJvType
					if lErr3 != nil || BoStatus.BoJvStatus == common.ErrorCode {
						pSgbReqRec.BoJvFailedCount++
						log.Println("SPJPO03", lErr3)
						RevFoStatus, lErr4 := BlockFOClientFund(pSgbReqRec, "R", "D")
						if lErr4 != nil || RevFoStatus.FoJvStatus == common.ErrorCode {
							pSgbReqRec.ReverseFoJvFailedCount++
							log.Println("SPJPO04", lErr4)
							pSgbReqRec.EmailDist = "3"
							pSgbReqRec.ProcessStatus = common.ErrorCode
							pSgbReqRec.ErrorStage = "A"
							return pSgbReqRec, lErr4
						} else if RevFoStatus.FoJvStatus == common.SuccessCode {
							pSgbReqRec.ReverseFoJvSuccessCount++
							pSgbReqRec.EmailDist = "2"
							pSgbReqRec.ProcessStatus = common.ErrorCode
							pSgbReqRec.ErrorStage = "B"
							return pSgbReqRec, lErr3
						}
					} else if BoStatus.BoJvStatus == common.SuccessCode {
						pSgbReqRec.BoJvSuccesssCount++
						//----------------------------------------------------------------
						//exchange calling place
						//----------------------------------------------------------------
						lRespRec, lErr5 := ExchangeProcess(pSgbReqRec, pValidSgb.Exchange, pBrokerId, r)
						pSgbReqRec = lRespRec
						if lErr5 != nil || lRespRec.ProcessStatus == common.ErrorCode {
							pSgbReqRec.TotalFailedCount++
							log.Println("SPJPO05", lErr5)
							RevBoStatus, lErr6 := BlockBOClientFund(pSgbReqRec, r, "R")
							if lErr6 != nil || RevBoStatus.BoJvStatus == common.ErrorCode {
								pSgbReqRec.ReverseBoJvFailedCount++
								log.Println("SPJPO06", lErr6)
								pSgbReqRec.EmailDist = "3"
								pSgbReqRec.ProcessStatus = common.ErrorCode
								pSgbReqRec.ErrorStage = "R"
								return pSgbReqRec, lErr6
							} else if BoStatus.BoJvStatus == common.SuccessCode {
								pSgbReqRec.ReverseBoJvSuccessCount++
								RevFoStatus, lErr7 := BlockFOClientFund(pSgbReqRec, "R", "C")
								if lErr7 != nil || RevFoStatus.FoJvStatus == common.ErrorCode {
									pSgbReqRec.ReverseFoJvFailedCount++
									log.Println("SPJPO07", lErr7)
									pSgbReqRec.EmailDist = "2"
									pSgbReqRec.ProcessStatus = common.ErrorCode
									pSgbReqRec.ErrorStage = "H"
									return pSgbReqRec, lErr7
								} else if RevFoStatus.FoJvStatus == common.SuccessCode {
									pSgbReqRec.ReverseFoJvSuccessCount++
									pSgbReqRec.ExchangeFailedCount++
									pSgbReqRec.EmailDist = lRespRec.EmailDist
									pSgbReqRec.ProcessStatus = common.ErrorCode
									pSgbReqRec.ErrorStage = lRespRec.ErrorStage
									return pSgbReqRec, lErr6
								}
							}
						} else {
							pSgbReqRec.TotalSuccessCount++
							pSgbReqRec.ExchangeSuccessCount++
							return pSgbReqRec, nil
						}
					}
				}
			}
			// ===============================INSUFFICIENT AMOUNT CHANGING UNIT================================
		} else {
			// =================================BACK OFFICE JV SUFFICIENT AMOUNT================================
			FoStatus, lErr8 := BlockFOClientFund(pSgbReqRec, "", "D")
			log.Println("FoStatus", FoStatus)
			pSgbReqRec.FoJvAmount = FoStatus.FoJvAmount
			pSgbReqRec.FoJvStatement = FoStatus.FoJvStatement
			pSgbReqRec.FoJvStatus = FoStatus.FoJvStatus
			pSgbReqRec.FoJvType = FoStatus.FoJvType
			// while processing FO Jv when we getting error it returns error and send mail to IT,RMS
			if lErr8 != nil || FoStatus.FoJvStatus == common.ErrorCode {
				log.Println("SPJPO08", lErr8)
				pSgbReqRec.EmailDist = "4"
				pSgbReqRec.ProcessStatus = common.ErrorCode
				pSgbReqRec.ErrorStage = "F"
				pSgbReqRec.FoJvFailedCount++
				return pSgbReqRec, lErr8
			} else if FoStatus.FoJvStatus == common.SuccessCode {
				pSgbReqRec.FoJvSuccessCount++

				//BackOffice Jv Processing
				BoStatus, lErr9 := BlockBOClientFund(pSgbReqRec, r, "C")
				pSgbReqRec.BoJvAmount = BoStatus.BoJvAmount
				pSgbReqRec.BoJvStatement = BoStatus.BoJvStatement
				pSgbReqRec.BoJvStatus = BoStatus.BoJvStatus
				pSgbReqRec.BoJvType = BoStatus.BoJvType
				if lErr9 != nil || BoStatus.BoJvStatus == common.ErrorCode {
					pSgbReqRec.BoJvFailedCount++
					log.Println("SPJPO09", lErr9)

					RevFoStatus, lErr10 := BlockFOClientFund(pSgbReqRec, "R", "D")
					if lErr10 != nil || RevFoStatus.FoJvStatus != common.ErrorCode {
						pSgbReqRec.ReverseFoJvFailedCount++
						log.Println("SPJPO10", lErr10)
						pSgbReqRec.EmailDist = "3"
						pSgbReqRec.ProcessStatus = common.ErrorCode
						pSgbReqRec.ErrorStage = "A"
						return pSgbReqRec, lErr10
					} else if RevFoStatus.FoJvStatus == common.SuccessCode {
						pSgbReqRec.ReverseFoJvSuccessCount++
						pSgbReqRec.EmailDist = "2"
						pSgbReqRec.ProcessStatus = common.ErrorCode
						pSgbReqRec.ErrorStage = "B"
						return pSgbReqRec, lErr9
					}
				} else if BoStatus.BoJvStatus == common.SuccessCode {
					pSgbReqRec.BoJvSuccesssCount++
					//----------------------------------------------------------------
					//exchange calling place
					//----------------------------------------------------------------
					lRespRec, lErr11 := ExchangeProcess(pSgbReqRec, pValidSgb.Exchange, pBrokerId, r)
					pSgbReqRec = lRespRec
					if lErr11 != nil || lRespRec.ProcessStatus == common.ErrorCode {
						pSgbReqRec.TotalFailedCount++
						log.Println("SPJPO11", lErr11)
						RevBoStatus, lErr12 := BlockBOClientFund(pSgbReqRec, r, "R")
						if lErr12 != nil || RevBoStatus.BoJvStatus == common.ErrorCode {
							pSgbReqRec.ReverseBoJvFailedCount++
							log.Println("SPJPO12", lErr12)
							pSgbReqRec.EmailDist = "3"
							pSgbReqRec.ProcessStatus = common.ErrorCode
							pSgbReqRec.ErrorStage = "R"
							return pSgbReqRec, lErr12
						} else if BoStatus.BoJvStatus == common.SuccessCode {
							pSgbReqRec.ReverseBoJvSuccessCount++
							RevFoStatus, lErr13 := BlockFOClientFund(pSgbReqRec, "R", "C")
							if lErr13 != nil || RevFoStatus.FoJvStatus == common.ErrorCode {
								log.Println("SPJPO13", lErr13)
								pSgbReqRec.ReverseFoJvFailedCount++
								pSgbReqRec.EmailDist = "2"
								pSgbReqRec.ProcessStatus = common.ErrorCode
								pSgbReqRec.ErrorStage = "H"
								return pSgbReqRec, lErr13
							} else if RevFoStatus.FoJvStatus == common.SuccessCode {
								pSgbReqRec.ReverseFoJvSuccessCount++
								pSgbReqRec.ExchangeFailedCount++
								pSgbReqRec.TotalFailedCount++
								pSgbReqRec.EmailDist = lRespRec.EmailDist
								pSgbReqRec.ProcessStatus = common.ErrorCode
								pSgbReqRec.ErrorStage = lRespRec.ErrorStage
								return pSgbReqRec, lErr12
							}
						}
					} else {
						pSgbReqRec.TotalSuccessCount++
						pSgbReqRec.ExchangeSuccessCount++
						return pSgbReqRec, nil
					}
				}
			}
		}
		// =================================BACK OFFICE JV SUFFICIENT AMOUNT================================
	} else {
		// ====================================while verifying the fund getting error================================
		pSgbReqRec.EmailDist = "4"
		pSgbReqRec.ProcessStatus = common.ErrorCode
		pSgbReqRec.ErrorStage = "V"
		pSgbReqRec.FoVerifyFailedCount++
		log.Println("pSgbReqRec.FoVerifyFailedCount", pSgbReqRec.FoVerifyFailedCount)
		return pSgbReqRec, nil
	}
	log.Println("PostJvandProcessOrder (-)")
	return pSgbReqRec, nil
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

Author:Pavithra
Date: 30DEC2023
*/
func BlockBOClientFund(pJvDetailRec exchangecall.JvReqStruct, pRequest *http.Request, pFlag string) (JvStatusStruct, error) {
	log.Println("BlockClientFund (+)")

	// this variables is used to get Status
	var lJVReq techexcel.JvInputStruct
	var lReqDtl apigate.RequestorDetails
	var lJvStatusRec JvStatusStruct

	lReqDtl = apigate.GetRequestorDetail(pRequest)
	lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")

	lCocd, lErr1 := sgbplaceorder.PaymentCode(pJvDetailRec.ClientId)
	if lErr1 != nil {
		log.Println("SPOPJV02", lErr1.Error())
		lJvStatusRec.BoJvStatus = common.ErrorCode
		return lJvStatusRec, lErr1
	} else {
		lJVReq.COCD = lCocd
		lJVReq.VoucherDate = time.Now().Format("02/01/2006")
		lJVReq.BillNo = "SGB" + pJvDetailRec.ClientId
		lJVReq.SourceTable = "a_sgb_orderdetails"
		lJVReq.SourceTableKey = pJvDetailRec.ReqOrderNo
		lJVReq.Amount = pJvDetailRec.Amount
		lJVReq.WithGST = "N"

		lUnit := strconv.Itoa(pJvDetailRec.Unit)
		lPrice := strconv.Itoa(pJvDetailRec.Price)

		lFTCaccount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseSGBorderAccount"])

		if pFlag == "R" {
			lJVReq.AccountCode = lFTCaccount
			lJVReq.CounterAccount = pJvDetailRec.ClientId
			lJVReq.Narration = "Failed of SGB Order Refund  " + lUnit + " * " + lPrice + " on " + lJVReq.VoucherDate + " - " + pJvDetailRec.ClientId
			lJvStatusRec.BoJvStatement = lJVReq.Narration
			lJvStatusRec.BoJvType = "C"
		} else if pFlag == "C" {
			lJVReq.AccountCode = pJvDetailRec.ClientId
			lJVReq.CounterAccount = lFTCaccount
			lJVReq.Narration = "SGB Purchase  " + lUnit + " * " + lPrice + " on " + lJVReq.VoucherDate + " - " + pJvDetailRec.ClientId
			lJvStatusRec.BoJvStatement = lJVReq.Narration
			lJvStatusRec.BoJvType = "D"
		}

		lJvStatusRec.BoJvStatus = "S"
		lJvStatusRec.BoJvAmount = pJvDetailRec.Amount
		//JV processing method
		lErr2 := clientfund.BOProcessJV(lJVReq, lReqDtl)
		if lErr2 != nil {
			log.Println("SPOPJV02", lErr2.Error())
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

Author:Nithish kumar
Date: 21Nov2023
*/
func BlockFOClientFund(pRequest exchangecall.JvReqStruct, pStatusFlag string, pPaymentType string) (JvStatusStruct, error) {
	log.Println("BlockFOClientFund (+)")
	var lFOJvReqStruct FOJvStruct
	var lJvStatusRec JvStatusStruct
	var lUrl string

	config := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
	lFOJvReqStruct.UserId = fmt.Sprintf("%v", config.(map[string]interface{})["VerifyUser"])
	lFOJvReqStruct.SourceUserId = lFOJvReqStruct.UserId

	lToken, lErr1 := sgbplaceorder.GetFOToken()
	if lErr1 != nil {
		log.Println("SBFCF01", lErr1.Error())
		lJvStatusRec.FoJvStatus = common.ErrorCode
		return lJvStatusRec, lErr1
	} else {
		// commented by pavithra
		// lPrice, lErr2 := strconv.Atoi(pRequest.Price)
		// if lErr2 != nil {
		// 	log.Println("SBFCF02", lErr2)
		// 	lJvStatusRec.FoJvStatus = common.ErrorCode
		// 	return lJvStatusRec, lErr2
		// }

		// lUnit, lErr3 := strconv.Atoi(pRequest.Unit)
		// if lErr3 != nil {
		// 	log.Println("SBFCF03", lErr3)
		// 	lJvStatusRec.FoJvStatus = common.ErrorCode
		// 	return lJvStatusRec, lErr3
		// }
		// Choose the url for the api based on the mode Debit or Credit
		// if pRequest.BoJvType == "D" {
		if pPaymentType == "D" {
			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayoutUrl"])
			lJvStatusRec.FoJvType = pRequest.BoJvType
			if pStatusFlag != "R" {
				lFOJvReqStruct.Remarks = "AMOUNT HOLD FOR SGB ORDER"
				lJvStatusRec.FoJvStatement = lFOJvReqStruct.Remarks
				lFOJvReqStruct.Amount = "-" + strconv.Itoa(pRequest.Price*pRequest.Unit)
			} else {
				lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM SGB ORDER"
				lJvStatusRec.FoJvStatement = lFOJvReqStruct.Remarks
				lFOJvReqStruct.Amount = strconv.Itoa(pRequest.Price * pRequest.Unit)
			}
			// } else if pRequest.BoJvType == "C" {
		} else if pPaymentType == "C" {
			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayinUrl"])
			lJvStatusRec.FoJvType = pRequest.BoJvType
			lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM SGB ORDER"
			lFOJvReqStruct.Amount = strconv.Itoa(pRequest.Price * pRequest.Unit)
		}
		lFOJvReqStruct.AccountId = pRequest.ClientId
		lJvStatusRec.FoJvAmount = lFOJvReqStruct.Amount

		lRequest, lErr4 := json.Marshal(lFOJvReqStruct)
		if lErr4 != nil {
			log.Println("SBFCF04", lErr4.Error())
			lJvStatusRec.FoJvStatus = common.ErrorCode
			return lJvStatusRec, lErr4
		} else {
			// construct the request body
			lBody := `jData=` + string(lRequest) + `&jKey=` + lToken
			lResp, lErr5 := clientfund.FOProcessJV(lUrl, lBody)
			if lErr5 != nil {
				log.Println("SBFCF05", lErr5.Error())
				lJvStatusRec.FoJvStatus = common.ErrorCode
				return lJvStatusRec, lErr5
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
	log.Println("BlockFOClientFund (-)")
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

Author:Pavithra
Date: 12JUNE2023
*/
func UpdateHeader(pNseRespRec nsesgb.SgbAddResStruct, pJvRec exchangecall.JvReqStruct, pBseRespRec bsesgb.SgbRespStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateHeader (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PIH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lHeaderId, lErr2 := GetOrderId(pJvRec, pExchange)
		if lErr2 != nil {
			log.Println("SPGOI02", lErr2)
			return lErr2
		} else {
			if lHeaderId != 0 {
				if pExchange == common.BSE {
					lSqlString := `update a_sgb_orderheader h
							set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
							h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.Status = ?,
							h.StatusCode = ?,h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,
							h.UpdatedBy = ?,h.UpdatedDate = now()
							where h.ClientId = ?
							and h.brokerId = ?
							and h.Id = ?`

					_, lErr5 := lDb.Exec(lSqlString, pBseRespRec.ScripId, pBseRespRec.PanNo, pBseRespRec.InvestorCategory, pBseRespRec.ApplicantName, pBseRespRec.Depository, pBseRespRec.DpId, pBseRespRec.ClientBenfId, pBseRespRec.GuardianName, pBseRespRec.GuardianPanno, pBseRespRec.GuardianRelation, common.SUCCESS, pBseRespRec.StatusCode, pBseRespRec.StatusMessage, pBseRespRec.ErrorCode, pBseRespRec.ErrorMessage, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
					if lErr5 != nil {
						log.Println("PIH05", lErr5)
						return lErr5
					} else {
						// call InsertDetails method to inserting the order details in order details table
						lErr6 := UpdateDetail(pJvRec, pBseRespRec.Bids, lHeaderId)
						if lErr6 != nil {
							log.Println("PIH06", lErr6)
							return lErr6
						}
					}
				} else if pExchange == common.NSE {
					lSqlString := `update a_sgb_orderheader h
					set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
					h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.Status = ?,h.ClientReferenceNo = ?,
					h.StatusCode = ?,h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,
					h.UpdatedBy = ?,h.UpdatedDate = now()
					where h.ClientId = ?
					and h.brokerId = ?
					and h.Id = ?`

					_, lErr5 := lDb.Exec(lSqlString, pBseRespRec.ScripId, pBseRespRec.PanNo, pBseRespRec.InvestorCategory, pBseRespRec.ApplicantName, pBseRespRec.Depository, pBseRespRec.DpId, pBseRespRec.ClientBenfId, pBseRespRec.GuardianName, pBseRespRec.GuardianPanno, pBseRespRec.GuardianRelation, common.SUCCESS, pNseRespRec.ClientRefNumber, pBseRespRec.StatusCode, pBseRespRec.StatusMessage, pBseRespRec.ErrorCode, pBseRespRec.ErrorMessage, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
					if lErr5 != nil {
						log.Println("PIH05", lErr5)
						return lErr5
					} else {
						// call InsertDetails method to inserting the order details in order details table
						lErr6 := UpdateDetail(pJvRec, pBseRespRec.Bids, lHeaderId)
						if lErr6 != nil {
							log.Println("PIH05", lErr6)
							return lErr6
						} else {
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

Author:Pavithra
Date: 12JUNE2023
*/
func UpdateDetail(pJvRec exchangecall.JvReqStruct, pRespBidArr []bsesgb.RespSgbBidStruct, pHeaderId int) error {
	log.Println("UpdateDetail (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PIH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		for lIdx := 0; lIdx < len(pRespBidArr); lIdx++ {

			lSqlString := `update a_sgb_orderdetails d
							set d.RespApplicationNo = ?,d.RespOrderNo = ?,d.ActionCode = ?,d.RespSubscriptionunit = ?,d.RespRate = ?,
							d.ErrorCode = ?,d.Message = ?,d.BoJvStatus = ?,d.BoJvAmount = ?,d.BoJvStatement = ?,d.BoJvType = ?,
							d.FoJvStatus = ?,d.FoJvAmount = ?,d.FoJvStatement = ?,d.FoJvType = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
							where d.HeaderId = ?
							and d.ReqOrderNo = ?`

			_, lErr5 := lDb.Exec(lSqlString, pRespBidArr[lIdx].OrderNo, pRespBidArr[lIdx].BidId, pRespBidArr[lIdx].ActionCode, pRespBidArr[lIdx].SubscriptionUnit, pRespBidArr[lIdx].Rate, pRespBidArr[lIdx].ErrorCode, pRespBidArr[lIdx].Message, pJvRec.BoJvStatus, pJvRec.BoJvAmount, pJvRec.BoJvStatement, pJvRec.BoJvType, pJvRec.FoJvStatus, pJvRec.FoJvAmount, pJvRec.FoJvStatement, pJvRec.FoJvType, common.AUTOBOT, pHeaderId, pJvRec.ReqOrderNo)
			if lErr5 != nil {
				log.Println("PIH05", lErr5)
				return lErr5
			} else {
				// call InsertDetails method to inserting the order details in order details table

			}
		}
	}
	log.Println("UpdateDetail (-)")
	return nil
}

func GetOrderId(pJvRec exchangecall.JvReqStruct, pExchange string) (int, error) {
	log.Println("GetOrderId (+)")

	var lHeaderId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPGOI01", lErr1)
		return lHeaderId, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select (case when count(1) > 0 then h.Id  else 0 end) Id
						from a_sgb_orderheader h,a_sgb_orderdetails d
						where h.Id = d.HeaderId 
						and h.ClientId = ?
						and d.ReqOrderNo = ?
						and h.Exchange  = ?
						and h.MasterId = (		
							select m.id 
							from a_sgb_master m
							where m.Symbol = ?
						)`
		lRows, lErr2 := lDb.Query(lCoreString, pJvRec.ClientId, pJvRec.ReqOrderNo, pExchange, pJvRec.Symbol)
		if lErr2 != nil {
			log.Println("SPGOI02", lErr2)
			return lHeaderId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lHeaderId)
				if lErr3 != nil {
					log.Println("SPGOI03", lErr3)
					return lHeaderId, lErr3
				}
				log.Println(lHeaderId)
			}
		}
	}
	log.Println("GetOrderId (-)")
	return lHeaderId, nil
}

func ExchangeProcess(pSgbReqRec exchangecall.JvReqStruct, pExchange string, pBrokerId int, r *http.Request) (exchangecall.JvReqStruct, error) {
	log.Println("ExchangeProcess (+)")

	if pExchange == common.BSE {
		_, lBseReqRec := exchangecall.ConstructExchReq(pSgbReqRec, pExchange)
		lRespRec, lErr1 := exchangecall.ApplyBseSgb(lBseReqRec, common.AUTOBOT, pBrokerId)
		var lNseResp nsesgb.SgbAddResStruct
		pSgbReqRec = exchangecall.ConstructExchResp(lNseResp, lRespRec, pExchange, pSgbReqRec)
		if lErr1 != nil || lRespRec.ErrorCode != "0" {
			// ================================ EXCHANGE FAILED====================================
			log.Println("SSEP01", lErr1)
			pSgbReqRec.EmailDist = "1"
			pSgbReqRec.ProcessStatus = common.ErrorCode
			pSgbReqRec.ErrorStage = "E"
			return pSgbReqRec, lErr1
		} else {
			pSgbReqRec.EmailDist = "0"
			pSgbReqRec.ProcessStatus = "Y"
			pSgbReqRec.ErrorStage = "Y"
			return pSgbReqRec, nil
		}
	} else if pExchange == common.NSE {
		lNseReqRec, _ := exchangecall.ConstructExchReq(pSgbReqRec, pExchange)
		lRespRec, lErr2 := exchangecall.ApplyNseSgb(lNseReqRec, common.AUTOBOT, pBrokerId)
		var lBseResp bsesgb.SgbRespStruct
		pSgbReqRec = exchangecall.ConstructExchResp(lRespRec, lBseResp, pExchange, pSgbReqRec)
		if lErr2 != nil || lRespRec.Status != "success" {
			// ================================ EXCHANGE FAILED ====================================
			log.Println("SSEP02", lErr2)
			pSgbReqRec.EmailDist = "1"
			pSgbReqRec.ProcessStatus = common.ErrorCode
			pSgbReqRec.ErrorStage = "E"
			return pSgbReqRec, lErr2
		} else {
			pSgbReqRec.EmailDist = "0"
			pSgbReqRec.ProcessStatus = "Y"
			pSgbReqRec.ErrorStage = "Y"
			return pSgbReqRec, nil
		}
	}
	log.Println("ExchangeProcess (-)")
	return pSgbReqRec, nil
}

//  This Method is used to update Nse Sgb Pending Details
// author prashanth

// func updateNseSgbPending(pRespSgbdata nsesgb.SgbAddResStruct, pReqSgbdata nsesgb.SgbAddReqStruct, pUser string) error {
// 	log.Println("updatePendingSgbStatus (+)")
// 	lStatus := common.SUCCESS
// 	lErr1 := updateNseSgbHeader(pRespSgbdata, pUser, lStatus)
// 	if lErr1 != nil {
// 		log.Println("SSUNSP01", lErr1)
// 		return lErr1
// 	} else {
// 		lErr2 := updateNseSgbDetails(pRespSgbdata, pUser)
// 		if lErr2 != nil {
// 			log.Println("SSUNSP02", lErr2)
// 			return lErr2
// 		} else {
// 			lErr3 := updateNseBidTracking(pRespSgbdata, pReqSgbdata, pUser)
// 			if lErr3 != nil {
// 				log.Println("SSUNSP03", lErr3)
// 				return lErr3
// 			}
// 		}
// 	}
// 	log.Println("updatePendingSgbStatus (+)")
// 	return nil

// }

// Pupose:This Method is used to Construct JvReq Struct into an JvData Struct .
// Parameters:

// 	Jv JvReqStruct

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 		{
// 			MasterId:
// 			ActionCode
// 			BidId : 005156
// 			JVamount : 5000
// 			ClientId :FT000069
// 			OrderNo : 005156
// 			Price :8000
// 			Transaction :
// 		},nil

// 	==========
// 	*On Error
// 	==========
// 	"",error

// Author:PRASHANTH
// Date: 27SEP2023
// */
//
//	This Method is used to Construct JvReq Struct into an JvData Struct
//
// author prashanth
// commeted by pavithra
// func JvConstructor(Jv JvReqStruct, Flag string) (JvDataStruct, error) {
// 	log.Println("NseJvConstruct (+)")
// 	var JvData JvDataStruct
// 	// JvData.MasterId = Jv.,
// 	JvData.ActionCode = Jv.ActionCode
// 	JvData.BidId = Jv.OrderNo
// 	// JvData.JVamount = Jv.Amount
// 	JvData.ClientId = Jv.ClientId
// 	JvData.OrderNo = Jv.OrderNo
// 	JvData.Price = Jv.Price
// 	JvData.Unit = Jv.Unit

// 	if Flag == "C" {
// 		JvData.Transaction = "C"
// 	} else if Flag == "R" {
// 		JvData.Transaction = "F"
// 	}

// 	lPrice, lErr1 := strconv.Atoi(Jv.Price)
// 	if lErr1 != nil {
// 		log.Println("SSJVC01", lErr1)
// 		return JvData, lErr1
// 	}

// 	lUnit, lErr2 := strconv.Atoi(Jv.Unit)
// 	if lErr2 != nil {
// 		log.Println("SSJVC02", lErr2)
// 		return JvData, lErr2
// 	}
// 	JvData.JVamount = strconv.Itoa(lPrice * lUnit)

// 	log.Println("NseJvConstruct (-)")
// 	return JvData, nil
// }

/* Pupose:This Method is used to Construct JvReq Struct into an JvData Struct .
Parameters:

	pRespSgbdata nsesgb.SgbAddResStruct

re


Response:

	==========
	*On Sucess
	==========
	{
		"status": S,
		"reason": ""

	},
	nil

	==========
	*On Error
	{
		"status": E,
		"reason": "Can't able to Update the data in database"
	}
	==========
	"",error

Author:PRASHANTH
Date: 27SEP2023

author prashanth */
//  This Method is used to Update SgbDetails For Nse

// func updateNseSgbDetails(pRespSgbdata nsesgb.SgbAddResStruct, pUser string) error {
// 	log.Println("updateNseSgbDetails(+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUNSDO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `update a_sgb_orderdetails d
// 							set d.RespSubscriptionunit = ?,d.RespRate = ?,d.ActionCode = ?,UpdatedDate = now(),UpdatedBy = ?
// 							where d.OrderNo = ?`

// 		_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Quantity, pRespSgbdata.Price,
// 			pRespSgbdata.Status, pUser, pRespSgbdata.OrderNumber)
// 		if lErr1 != nil {
// 			log.Println("SSUNSDO2", lErr1)
// 			return lErr1
// 		}
// 	}
// 	log.Println("updateNseSgbDetails(-)")
// 	return nil
// }

/* Pupose:This Method is used to Update Sgb Header For Nse.
Parameters:

	pRespSgbdata nsesgb.SgbAddResStruct,user,Status

Response:

	==========
	*On Sucess
	==========
	{
		"status": S,
		"reason": ""

	},
	nil

	==========
	*On Error
	{
		"status": E,
		"reason": "Can't able to Update the data from database"
	}
	==========
	"",error

Author:PRASHANTH
Date: 27SEP2023
*/
// This Method is used to Update Sgb Header For Nse

// func updateNseSgbHeader(pRespSgbdata nsesgb.SgbAddResStruct, pUser string, pStatus string) error {
// 	log.Println("updateSgbHeader(+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUNSHO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `update a_sgb_orderheader  h
// 							set h.ScripId = ?,h.PanNo = ?,h.Depository = ?,
// 							h.DpId=?,h.ClientBenfId=?
// 							,h.status=?,
// 							h.UpdatedDate = now(),h.UpdatedBy = ?
// 							where h.Id in (select d.HeaderId
// 							from a_sgb_orderdetails d,a_sgb_orderheader h
// 							where d.HeaderId = h.Id and d.OrderNo = ?)`

// 		_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Symbol, pRespSgbdata.Pan,
// 			pRespSgbdata.Depository, pRespSgbdata.DpId, pRespSgbdata.ClientBenId, pStatus, pUser, pRespSgbdata.OrderNumber)
// 		if lErr1 != nil {
// 			log.Println("SSUNSHO2", lErr1)
// 			return lErr1
// 		}
// 	}

// 	log.Println("updateSgbHeader(-)")
// 	return nil
// }

/* Pupose:This Method is used to Update Sgb Header For Nse.
Parameters:

	pRespSgbdata nsesgb.SgbAddResStruct,user,Status

Response:

	==========
	*On Sucess
	==========
	{
		"status": S,
		"reason": ""

	},
	nil

	==========
	*On Error
	{
		"status": E,
		"reason": "Can't able to Update the data from database"
	}
	==========
	"",error

Author:PRASHANTH
Date: 27SEP2023
*/
// This Method is used to update Bid Tracking Details of Nse

// func updateNseBidTracking(pRespSgbdata nsesgb.SgbAddResStruct, pReqSgbData nsesgb.SgbAddReqStruct, pUser string) error {
// 	log.Println("updateBidTracking(+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUNBO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `update a_sgbtracking_table  b
// 							set b.ApplicationStatus =?,UpdatedDate = now(),UpdatedBy = ?
// 							where b.OrderNo = ? and b.ActivityType = ? and b.Unit = ?`

// 		_, lErr1 = lDb.Exec(lCoreString, common.SUCCESS, pUser, pRespSgbdata.OrderNumber, pReqSgbData.ActivityType, pRespSgbdata.Quantity)
// 		if lErr1 != nil {
// 			log.Println("SSUNBO2", lErr1)
// 			return lErr1
// 		}
// 	}

// 	log.Println("updateBidTracking(-)")
// 	return nil
// }

// func NseReqConstruct(BSEStruct []bsesgb.SgbReqStruct, JVStruct []JvReqStruct) (nsesgb.SgbAddReqStruct, error) {
// 	log.Println("NseReqConstruct (+)")
// 	var NseStruct nsesgb.SgbAddReqStruct
// 	var NseReqArr []nsesgb.SgbAddReqStruct
// 	for _, JV := range JVStruct {
// 		NseStruct.Symbol = JV.Symbol
// 		var lErr1 error
// 		NseStruct.Price, lErr1 = common.ConvertStringToFloat(JV.Price)
// 		if lErr1 != nil {
// 			log.Println("lErr1", lErr1)
// 			return NseStruct, lErr1
// 		}
// 		NseStruct.Quantity = strconv.Atoi(JV.Unit)
// 		NseStruct.ActivityType = JV.ActionCode
// 		NseStruct.PhysicalDematFlag = "D"
// 		NseStruct.ClientCode = ""
// 		NseStruct.OrderNumber = strconv.Atoi(JV.OrderNo)
// 		NseStruct.ClientRefNumber = ""

// 	}

// 	for _, Bse := range BSEStruct {
// 		NseStruct.ClientBenId = Bse.ClientBenfId
// 		NseStruct.Depository = Bse.Depository
// 		NseStruct.DpId = Bse.DpId
// 		NseStruct.Pan = Bse.PanNo
// 	}

// 	NseReqArr = append(NseReqArr, NseStruct)

// 	log.Println("NseReqConstruct (-)")
// 	return NseStruct, nil
// }

// func NseReqConstruct() ([]NseReqStruct, error) {
// 	log.Println("NseReqConstruct (+)")
// 	var NseStruct NseReqStruct
// 	var NseReqArr []NseReqStruct
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSFP01", lErr1)
// 		return NseReqArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select m.Symbol ,d.OrderNo ,d.ReqRate e,d.ReqSubscriptionUnit ,'D', '',h.PanNo,h.Depository ,h.ClientBenfId ,h.DpId, (CASE
// 		WHEN d.ActionCode = 'M' THEN 'Modify'
// 		WHEN d.ActionCode = 'N' THEN 'New'
// 		WHEN d.ActionCode = 'D' THEN 'Delete'
// 		ELSE d.ActionCode
// 	END ) AS ActionDescription,h.clientId
// 		from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
// 		where h.MasterId = m.id
// 		and d.HeaderId = h.Id
// 		and h.Status = "pending"
// 		and m.BiddingStartDate <= curdate()
// 		or m.BiddingEndDate >= curdate()
// 		and time(now()) between m.DailyStartTime and m.DailyEndTime
// 		and h.cancelFlag != 'Y' and m.Exchange = "NSE"`

// 		lRows, lErr2 := lDb.Query(lCoreString)
// 		if lErr2 != nil {
// 			log.Println("SSFP02", lErr2)
// 			return NseReqArr, lErr2
// 		} else {

// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&NseStruct.Symbol, &NseStruct.OrderNumber, &NseStruct.Price, &NseStruct.Quantity, &NseStruct.PhysicalDematFlag, &NseStruct.ClientCode, &NseStruct.Pan, &NseStruct.Depository, &NseStruct.ClientBenId, &NseStruct.DpId, &NseStruct.ActivityType, &NseStruct.ClientId)
// 				if lErr3 != nil {
// 					log.Println("SSFP03", lErr3)
// 					return NseReqArr, lErr3
// 				} else {

// 					NseReqArr = append(NseReqArr, NseStruct)
// 				}

// 			}
// 		}
// 	}
// 	log.Println("NseReqArr", NseReqArr)

// 	log.Println("NseReqConstruct (-)")
// 	return NseReqArr, nil
// }
/* Pupose:This method is used to update order status during schedule for Sgb Place Order  .
Parameters:

	processFlag string, ScheduleStatus string, pHeaderId int, pBrokerId int, pJv {}

Response:

	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:PRASHANTH
Date: 19DEC2023
*/
func updateSchStatus(processFlag string, ScheduleStatus string, pHeaderId int, pBrokerId int, pJv exchangecall.JvReqStruct) error {
	log.Println("updateSchStatus(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SUSS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// log.Println(" processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId", processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId)
		lCoreString := ` update a_sgb_orderheader 
		                 set  ProcessFlag =  ? , ScheduleStatus = ?
	                     where ClientId = ?
	                     and  brokerId = ?
	                     and Id = ?`
		_, lErr2 := lDb.Exec(lCoreString, processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId)
		// log.Println("lCoreString", lCoreString)
		if lErr2 != nil {
			log.Println("SUSS02", lErr2)
			return lErr2

		}
	}

	log.Println("updateSchStatus (-)")
	return nil
}

func updateOrderDetails(pJvRec exchangecall.JvReqStruct, pExchange string, pBrokerId int, processFlag string, ScheduleStatus string) error {
	log.Println("updateOrderDetails (+)")

	lHeaderId, lErr1 := GetOrderId(pJvRec, pExchange)
	if lErr1 != nil {
		log.Println("SUOD01", lErr1)
		return lErr1
	} else {
		lErr2 := updateSchStatus(processFlag, ScheduleStatus, lHeaderId, pBrokerId, pJvRec)
		if lErr2 != nil {
			log.Println("SUOD02", lErr2)
			return lErr2
		}

	}

	log.Println("updateOrderDetails (-)")
	return nil
}

func UpdateExchFailed(pJvRec exchangecall.JvReqStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateExchFailed (+)")

	lHeaderId, lErr1 := GetOrderId(pJvRec, pExchange)
	if lErr1 != nil {
		log.Println("SUEF01", lErr1)
		return lErr1
	} else {
		lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr1 != nil {
			log.Println("SUSS01", lErr1)
			return lErr1
		} else {
			defer lDb.Close()

			// log.Println(" processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId", processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId)
			lCoreString := ` update a_sgb_orderheader h
							set h.Status = ?,h.UpdatedBy = ?,h.UpdatedDate = now()
							where h.ClientId = ?
							and h.brokerId = ?
							and h.Id = ?`
			_, lErr2 := lDb.Exec(lCoreString, common.FAILED, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
			// log.Println("lCoreString", lCoreString)
			if lErr2 != nil {
				log.Println("SUSS02", lErr2)
				return lErr2

			}
		}
	}
	log.Println("UpdateExchFailed (-)")
	return nil
}
