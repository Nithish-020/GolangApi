package ncbplaceorder

// import (
// 	"fcs23pkg/common"
// 	"fcs23pkg/integration/nse/nsencb"
// 	"log"
// )

// type ErrorStruct struct {
// 	ErrCode string `json:"errCode"`
// 	ErrMsg  string `json:"errMsg"`
// }

// func NcbNsePlaceOrder(pExhangeReq nsencb.NcbAddReqStruct, pReqRec NcbReqStruct, pClientId string, pMailId string) (NcbRespStruct, nsencb.NcbAddResStruct, ErrorStruct) {
// 	log.Println("NcbNsePlaceOrder(+)")

// 	var lRespRec NcbRespStruct
// 	var lError ErrorStruct

// 	lResp, lErr1 := ProcessNseReq(pExhangeReq, pReqRec, pClientId, pMailId)
// 	if lErr1 != nil {
// 		log.Println("NNPO01", lErr1.Error())
// 		lRespRec.Status = common.ErrorCode
// 		lError.ErrCode = "NNPO01"
// 		lError.ErrMsg = "Exchange Server is Busy right now,Try After Sometime."
// 		return lRespRec, lResp, lError
// 	} else {
// 		log.Println("lRespRec.status", lResp.Status)
// 		if lResp.Status != "" {
// 			// for lRespIdx := 0; lRespIdx < len(lResp); lRespIdx++ {
// 			if lResp.Status == "failed" {
// 				lRespRec.OrderStatus = lResp.Status
// 				lRespRec.OrderReason = lResp.Reason
// 				lRespRec.Status = common.ErrorCode
// 				lError.ErrCode = "NNPO02"
// 				lError.ErrMsg = lResp.Reason + "\n" + "Application Failed"
// 				log.Println("NNPO02", lRespRec.OrderStatus, lRespRec.OrderReason)
// 				return lRespRec, lResp, lError
// 			} else if lResp.Status == "success" {
// 				lRespRec.OrderStatus = lResp.Status
// 				lRespRec.OrderReason = lResp.Reason
// 				lRespRec.Status = common.SuccessCode
// 			}
// 			// }
// 		} else {
// 			log.Println("lRespRec.status", lResp.Status)
// 			log.Println("NNPO03", "Unable to proceed Application!!!")
// 			lError.ErrCode = "NNPO03"
// 			lError.ErrMsg = "Unable to proceed Application!!!"
// 			log.Println("lRespRec", lRespRec)
// 			return lRespRec, lResp, lError
// 		}
// 	}
// 	log.Println("NcbNsePlaceOrder(-)")
// 	return lRespRec, lResp, lError
// }

/*
Pupose:This Function is used to if there is no value in the array it show the string value from the ncb master table
Parameters:

masterArr

Response:

*On Sucess
=========

NCB_GSHistoryNoDataText = "You haven't invested in any GOIs."
NCB_TBHistoryNoDataText = "You haven't invested in any TBills."
NCB_SGHistoryNoDataText = "You haven't invested in any SDLs."

!On Error
========
nil

Author: KAVYA DHARSHANI M
Date: 02JAN2024
*/

// func processNcbHistoryArr(pmasterArr []ActiveNcbStruct, lBrokerId string, pClientId string, pMethod string, lReportStruct ReportReqStruct) {

// 	var lNcbResp NcbOrderHistoryResp
// 	var lNcbOrderHistoryRec NcbOrderHistoryStruct
// 	var lOrderCount int
// 	var lReqUnits, lReqUnitPrice, lReqAmount, lAppliedUnits, lAppliedAmount, lAppliedUnitPrice, lAllocatedUnits, lStatus, lAllocatedUnitPrice, lAllocatedAmount string

// 	log.Println("processNcbHistoryArr(+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("NOHR01", lErr1)
// 		lNcbResp.Status = common.ErrorCode
// 		return lNcbResp, lErr1
// 	} else {
// 		defer lDb.Close()

// 		for _, lMasterRec := range pmasterArr {
// 			lMasterid := strconv.Itoa(lMasterRec.Id)
// 			lAdditionalCoreString := ""

// 			lCoreString := `select  h.Id Id,h.Symbol,h.Series, h.exchange,h.ClientId, d.ReqOrderNo,
// 								d.RespOrderNo , date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p')OrderDate, nvl(d.ReqUnit,0) as RequestedUnits,
// 								nvl(d.Reqprice,0) as RequestedUnitPrice,
// 								nvl(nvl(d.ReqUnit,0) *nvl(d.Reqprice,0),0) as RequestedAmount,
// 								nvl(d.ReqUnit,0) as AppliedUnits,
// 								nvl(d.Reqprice,0 )as AppliedUnitPrice,
// 								nvl(nvl(d.RespUnit ,0)  *nvl(d.Respprice ,0) ,0) as AppliedAmount,
// 								nvl(d.AllotedUnit ,0) as AllocatedUnits,
// 								nvl(d.AllotedRate ,0 )as AllocatedUnitPrice,
// 								nvl (nvl(d.AllotedUnit,0)  *nvl(d.AllotedRate,0)  ,0) as AllocatedAmount,
// 								nvl(h.SIvalue, '') SIvalue,
// 								nvl(h.SItext,'') SItext ,
// 								nvl(h.RbiRemarks,'') RbiRemarks,
// 								nvl(h.DpRemarks,'') DpRemarks,
// 								(case when  h.status ='failed' or h.cancelFlag  = 'Y' and h.status = 'success' then 'R' else 'G' end ) statuscolor,
// 								(case when h.status='failed' then 'F' when h.cancelFlag = 'Y' and h.status = 'success' then 'BC' else 'S' end )Status
// 						from a_ncb_orderdetails d, a_ncb_orderheader h
// 						where  d.HeaderId = h.Id
// 						and h.brokerId =  '` + lBrokerId + `'
// 					   and h.MasterId =  '` + lMasterid + `'`

// 			if pMethod == "getNcbOrderHistory" {
// 				lAdditionalCoreString = ` and h.clientId = '` + pClientId + `'
// 						and h.status is not null
// 						order by h.Id desc`
// 			} else if pMethod == "GetReport" {
// 				lAdditionalCoreString = ` AND ('` + lReportStruct.Symbol + `' = '' OR h.Symbol = '` + lReportStruct.Symbol + `')
// 						AND ('` + lReportStruct.ClientId + `' = '' OR  h.clientId = '` + lReportStruct.ClientId + `')
// 						AND ('` + lReportStruct.FromDate + `' = '' OR h.CreatedDate BETWEEN CONCAT('` + lReportStruct.FromDate + `',' 00:00:00.000') AND CONCAT('` + lReportStruct.ToDate + `',' 23:59:59.000'))`

// 			} else if pMethod == "GetDefault" {
// 				lAdditionalCoreString = ` and d.CreatedDate between concat (date(now()), ' 00:00:00.000')
// 						 and concat (date(now()), ' 23:59:59.000')`
// 			}

// 			lCoreString = lCoreString + lAdditionalCoreString

// 			lRows, lErr2 := lDb.Query(lCoreString)
// 			if lErr2 != nil {
// 				log.Println("NOHR02", lErr2)
// 				lNcbResp.Status = common.ErrorCode
// 				return lNcbResp, lErr2
// 			} else {
// 				for lRows.Next() {

// 					lErr3 := lRows.Scan(&lNcbOrderHistoryRec.Id, &lNcbOrderHistoryRec.Symbol, &lNcbOrderHistoryRec.Series, &lNcbOrderHistoryRec.Exchange, &lNcbOrderHistoryRec.ClientId, &lNcbOrderHistoryRec.OrderNo, &lNcbOrderHistoryRec.OrderDate, &lNcbOrderHistoryRec.RequestedUnit, &lNcbOrderHistoryRec.RequestedUnitPrice, &lNcbOrderHistoryRec.RequestedAmount, &lNcbOrderHistoryRec.AppliedUnit, &lNcbOrderHistoryRec.AppliedUnitPrice, &lNcbOrderHistoryRec.AppliedAmount, &lNcbOrderHistoryRec.AllotedUnit, &lNcbOrderHistoryRec.AllotedUnitPrice, &lNcbOrderHistoryRec.AllotedAmount, &lNcbOrderHistoryRec.SIValue, &lNcbOrderHistoryRec.SIText, &lNcbOrderHistoryRec.RBIStatus, &lNcbOrderHistoryRec.DPStatus, &lNcbOrderHistoryRec.StatusColor, &lStatus)

// 					if lErr3 != nil {
// 						log.Println("NOHR03", lErr3)
// 						lNcbResp.Status = common.ErrorCode
// 						lNcbResp.ErrMsg = lErr3.Error()
// 						return lNcbResp, lErr3
// 					} else {
// 						if lStatus == "S" {
// 							lOrderCount++
// 						}

// 						lNcbOrderHistoryRec.Symbol = lMasterRec.Symbol
// 						lNcbOrderHistoryRec.Series = lMasterRec.Series
// 						lNcbOrderHistoryRec.Isin = lMasterRec.Isin
// 						lNcbOrderHistoryRec.DiscountAmt = lMasterRec.DiscountAmt
// 						lNcbOrderHistoryRec.DiscountText = lMasterRec.DiscountText
// 						lNcbOrderHistoryRec.DateRange = lMasterRec.DateRange
// 						lNcbOrderHistoryRec.StartDateWithTime = lMasterRec.StartDateWithTime
// 						lNcbOrderHistoryRec.EndDateWithTime = lMasterRec.EndDateWithTime
// 						lNcbOrderHistoryRec.RequestedUnit, _ = strconv.Atoi(lReqUnits)
// 						lNcbOrderHistoryRec.RequestedUnitPrice, _ = strconv.Atoi(lReqUnitPrice)
// 						lNcbOrderHistoryRec.RequestedAmount, _ = strconv.Atoi(lReqAmount)
// 						lNcbOrderHistoryRec.AppliedUnit, _ = strconv.Atoi(lAppliedUnits)
// 						lNcbOrderHistoryRec.AppliedUnitPrice, _ = strconv.Atoi(lAppliedUnitPrice)
// 						lNcbOrderHistoryRec.AppliedAmount, _ = strconv.Atoi(lAppliedAmount)
// 						lNcbOrderHistoryRec.AllotedUnit, _ = strconv.Atoi(lAllocatedUnits)
// 						lNcbOrderHistoryRec.AllotedUnitPrice, _ = strconv.Atoi(lAllocatedUnitPrice)
// 						lNcbOrderHistoryRec.AllotedAmount, _ = strconv.Atoi(lAllocatedAmount)
// 						// Append Upi End Point in lRespRec.SgbHistoryArr array
// 						OrderStatus, lErr4 := NcbDiscription(lStatus)
// 						if lErr4 != nil {
// 							log.Println("LGSHD04", lErr4)
// 							lNcbResp.Status = common.ErrorCode
// 							return lNcbResp, lErr4
// 						} else {
// 							lNcbOrderHistoryRec.OrderStatus = OrderStatus
// 						}
// 					}

// 					if lNcbOrderHistoryRec.Series == "GS" {
// 						lGSNcbRespArr = append(lGSNcbRespArr, lNcbOrderHistoryRec)
// 					} else if lNcbOrderHistoryRec.Series == "TB" {
// 						lTBNcbRespArr = append(lTBNcbRespArr, lNcbOrderHistoryRec)
// 					} else {
// 						lSGNcbRespArr = append(lSGNcbRespArr, lNcbOrderHistoryRec)
// 					}

// 				}

// 			}

// 		}
// 		if lGSNcbRespArr != nil {
// 			lNcbResp.GSecOrderHistoryArr = lGSNcbRespArr
// 			lNcbResp.OrderCount = lOrderCount
// 			lNcbResp.HistoryFound = "Y"
// 		} else {
// 			if len(lGSNcbRespArr) == 0 {
// 				lNcbResp.HistoryFound = "N"
// 				lNcbResp.HistoryNoDataText = lHistoryNoDataTxt
// 			}
// 		}
// 	}

// 	log.Println("processNcbHistoryArr(-)")
// }

// lExchangeReq, lErr5 := ConstructNCBReqStruct(lReqRec, lClientId)
// 					if lErr5 != nil {
// 						log.Println("PNPO05", lErr5.Error())
// 						lRespRec.Status = common.ErrorCode
// 						fmt.Fprintf(w, helpers.GetErrorString("PNPO05", "Unable to process your request now. Please try after sometime"))
// 						return
// 					} else {
// 						log.Println(lExchangeReq)
// 						log.Println(lReqRec.MasterId, "lReqRec.MasterId")
// 						lTodayAvailable, lErr6 := validatencb.CheckNcbEndDate(lReqRec.MasterId)
// 						if lErr6 != nil {
// 							log.Println("PNPO06", lErr6.Error())
// 							lRespRec.Status = common.ErrorCode
// 							fmt.Fprintf(w, helpers.GetErrorString("PNPO06", "Unable to process your request now. Please try after sometime"))
// 							return
// 						} else {
// 							if lTodayAvailable == "True" {

// 								lErr7 := LocalUpdate(lReqRec, lExchangeReq, lClientId, lExchange, r, lBrokerId)
// 								if lErr7 != nil {
// 									log.Println("PNPO07", lErr7.Error())
// 									lRespRec.Status = common.ErrorCode
// 									fmt.Fprintf(w, helpers.GetErrorString("PNPO07", "Unable to process your request now. Please try after sometime"))
// 									return
// 								} else {
// 									if lReqRec.ActionCode == "N" {
// 										lRespRec.OrderStatus = "Order Placed Successfully"
// 										lRespRec.Status = common.SuccessCode
// 									} else if lReqRec.ActionCode == "M" {
// 										lRespRec.OrderStatus = "Order Modified Successfully"
// 										lRespRec.Status = common.SuccessCode
// 									} else if lReqRec.ActionCode == "D" {
// 										lRespRec.OrderStatus = "Order Deleted Successfully"
// 										lRespRec.Status = common.SuccessCode
// 									}
// 								}

// 							} else {
// 								lRespRec.Status = common.ErrorCode
// 								lRespRec.ErrMsg = "Timing Closed for NCB"
// 								log.Println("Timing Closed for NCB")
// 							}
// 						}
// 					}

/*
Pupose:This method is used to get the client BenId and ClientName for order input.
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	123456789012,Lakshmanan Ashok Kumar,nil

	==========
	*On Error
	==========
	"","",error

Author:KAVYADHARSHANI M
Date: 12OCT2023
*/

// func getDpId(lClientId string) (string, string, error) {
// 	log.Println("getDpId (+)")

// 	// this variables is used to get DpId and ClientName from the database.
// 	var lDpId string
// 	var lClientName string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
// 	if lErr1 != nil {
// 		log.Println("NGDI01", lErr1)
// 		return lDpId, lClientName, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select  idm.CLIENT_DP_CODE, idm.CLIENT_DP_NAME
// 						from   TECHEXCELPROD.CAPSFO.DBO.IO_DP_MASTER idm
// 						where idm.CLIENT_ID = ?
// 						and DEFAULT_ACC = 'Y'
// 						and DEPOSITORY = 'CDSL' `
// 		lRows, lErr2 := lDb.Query(lCoreString, lClientId)
// 		if lErr2 != nil {
// 			log.Println("NGDI02", lErr2)
// 			return lDpId, lClientName, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lDpId, &lClientName)
// 				if lErr3 != nil {
// 					log.Println("NGDI03", lErr3)
// 					return lDpId, lClientName, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("getDpId (-)")
// 	return lDpId, lClientName, nil
// }

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
	 `select h.MasterId,nvl(d.ReqAmount,0.0), d.Reqprice, d.activityType activitytype, CAST((h.ReqInvestmentunit /(n.Lotsize/100)) AS SIGNED) Lotsize, d.ReqUnit unit, h.Symbol,h.cancelFlag flag, h.ReqapplicationNo, d.ReqOrderNo
		                from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
		                where h.id = d.headerId
		                and n.id  = h.MasterId
		                and d.status <> 'failed' and h.cancelFlag <> 'Y'
                        and h.clientId = ? and h.MasterId = ? and d.ReqOrderNo = ?`

Author:KAVYA DHARSHANI
Date: 12OCT2023
*/
// func getPan(pClientId string) (string, error) {
// 	log.Println("getPan(+)")

// 	// this variables is used to get Pan number of the client from the database.
// 	var lPanNo string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
// 	if lErr1 != nil {
// 		log.Println("NGP01", lErr1)
// 		return lPanNo, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select pan_no
// 						from TECHEXCELPROD.CAPSFO.DBO.client_details
// 						where client_Id = ? `
// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
// 		if lErr2 != nil {
// 			log.Println("NGP02", lErr2)
// 			return lPanNo, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lPanNo)
// 				if lErr3 != nil {
// 					log.Println("NGP03", lErr3)
// 					return lPanNo, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("getPan(-)")
// 	return lPanNo, nil
// }

// func GetOldNcbOrderHistorydetail(pClientId string) ([]NcbOrderHistoryStruct, []NcbOrderHistoryStruct, []NcbOrderHistoryStruct, error) {
// 	log.Println("GetNcbOrderHistorydetail(+)")

// 	var lNcbHistoryRec NcbOrderHistoryStruct

// 	var lGsecHistoryArr []NcbOrderHistoryStruct
// 	var lTbillHistoryArr []NcbOrderHistoryStruct
// 	var lSdlHistoryArr []NcbOrderHistoryStruct

// 	// var lUnit, lPrice string
// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("NGNHD01", lErr1)
// 		return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select nvl(h.Id,'') id, nvl(h.ApplicationNo,'') ,n.Name Name, nvl(d.price,'') price ,n.Symbol, n.Series,nvl(DATE_FORMAT(h.CreatedDate , '%d-%b-%Y'),'') CreatedBy , CONCAT( case WHEN DAY(n.BiddingEndDate) % 10 = 1 AND DAY(n.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(n.BiddingEndDate), 'st')
// 				       WHEN DAY(n.BiddingEndDate) % 10 = 2 AND DAY(n.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(n.BiddingEndDate), 'nd')
// 				       WHEN DAY(n.BiddingEndDate) % 10 = 3 AND DAY(n.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(n.BiddingEndDate), 'rd')
// 				       ELSE CONCAT(DAY(n.BiddingEndDate), 'th')
// 	                   end,' ',
// 			              DATE_FORMAT(n.BiddingEndDate, '%b %Y'),' | ',
// 			              TIME_FORMAT(n.DailyEndTime , '%h:%i%p')) AS formatted_datetime,n.BiddingStartDate,
// 				          (case when h.MasterId = n.Id and h.CancelFlag = 'Y' and h.Status = 'success' then 'Y' else 'N' end) Flag,
// 				         lower(h.Status) Status
//                        from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
//                        where n.id  = h.MasterId  and d.HeaderId = h.Id and h.ClientId = ?
//                        group by h.applicationNo
//                      order by h.Id desc`

// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
// 		if lErr2 != nil {
// 			log.Println("NGNHD02", lErr2)
// 			return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows.Next() {

// 				lErr3 := lRows.Scan(&lNcbHistoryRec.Id, &lNcbHistoryRec.ApplicationNo, &lNcbHistoryRec.Symbol, &lNcbHistoryRec.Symbol, &lNcbHistoryRec.Series, &lNcbHistoryRec.OrderDate)
// 				if lErr3 != nil {
// 					log.Println("NGNHD03", lErr3)
// 					return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, lErr3
// 				} else {
// 					// lNcbHistoryRec.Unit, _ = strconv.Atoi(lUnit)
// 					// lNcbHistoryRec.Price, _ = strconv.Atoi(lPrice)

// 					// lNcbHistoryRec.Total = lNcbHistoryRec.Unit * lNcbHistoryRec.Price
// 					// Append Upi End Point in lRespRec.SgbHistoryArr array

// 					if lNcbHistoryRec.Series == "GS" {
// 						lGsecHistoryArr = append(lGsecHistoryArr, lNcbHistoryRec)
// 						log.Println("lGsecHistoryArr--->GS", lGsecHistoryArr)
// 					} else if lNcbHistoryRec.Series == "TB" {
// 						lTbillHistoryArr = append(lTbillHistoryArr, lNcbHistoryRec)
// 						log.Println("lTbillHistoryArr--->TB", lTbillHistoryArr)
// 					} else {
// 						lSdlHistoryArr = append(lSdlHistoryArr, lNcbHistoryRec)
// 						log.Println("lSdlHistoryArr--->TB", lSdlHistoryArr)
// 					}

// 				}
// 			}

// 		}
// 	}
// 	log.Println("GetNcbOrderHistorydetail(-)")
// 	return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, nil
// }

// func NcbLookTransaction(pTransaction nsencb.NcbAddResStruct) (nsencb.NcbAddResStruct, error) {

// 	log.Println("NcbLookTransaction (+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("NNLTO1", lErr1)
// 		return pTransaction, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select d.description
// 		                from xx_lookup_details d, xx_lookup_header h
// 		                where h.id = d.headerid
// 		                and h.Code = 'Sgbstatus'
// 		                and d.Code  = ?	`

// 		lRows, lErr2 := lDb.Query(lCoreString, pTransaction.ClearingStatus)
// 		if lErr2 != nil {
// 			log.Println("NNLTO2", lErr2)
// 			return pTransaction, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&pTransaction.ClearingStatus)
// 				if lErr3 != nil {
// 					log.Println("NNLTO3", lErr3)
// 					return pTransaction, lErr3
// 				}
// 			}
// 		}

// 	}

// 	log.Println("NcbLookTransaction (-)")
// 	return pTransaction, nil
// }
