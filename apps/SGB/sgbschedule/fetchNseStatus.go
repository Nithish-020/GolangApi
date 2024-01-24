package sgbschedule

import (
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bsesgb"
	"fcs23pkg/integration/nse/nsesgb"
	"fmt"
	"log"
	"strconv"
)

type RespStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

type BiddingDateStruct struct {
	Date   string `json:"startDate"`
	Symbol string `json:"symbol"`
}

/*
Pupose:This API method is used to place a order
Request:

	lReqRec

Response:

	==========
	*On Sucess
	==========
	{
		"appResponse":[
			{

			},
			{

			}
		],
		"status":"S",
		"errMsg":""
	}
	==========
	*On Error
	==========

		{
		"appResponse":[],
		"status":"E",
		"errMsg":""
	}

Author:Kavya Dharshani
Date: 20SEP2023
*/
// func NseSgbFetchStatus(lwg *sync.WaitGroup, pBrokerId int, pUser string) error {
func NseSgbFetchStatus(pBrokerId int, pUser string) (string, error) {
	log.Println("NseSgbFetchStatus (+)")

	//commented by pavithra
	// defer lwg.Done()
	var lDateResp []BiddingDateStruct
	var lStatusFlag string

	lDRes, lErr1 := GetSgbEndDate(lDateResp)
	if lErr1 != nil {
		log.Println("SSNSFS01", lErr1)
		lStatusFlag = common.ErrorCode
		return lStatusFlag, lErr1
	} else {
		// log.Println("lDRes", lDRes)
		if len(lDRes) != 0 {
			//dateloop
			for Idx := 0; Idx < len(lDRes); Idx++ {
				//bsetoken
				lToken, lErr2 := exchangecall.GetToken(pUser, pBrokerId)
				if lErr2 != nil {
					log.Println("SSNSFS02", lErr2)
					lStatusFlag = common.ErrorCode
					return lStatusFlag, lErr2
				} else {
					if lToken != "" {
						lResp, lErr3 := nsesgb.SgbTransactionsMaster(lToken, lDRes[Idx].Date, pUser)
						if lErr3 != nil {
							log.Println("SSNSFS03", lErr3)
							lStatusFlag = common.ErrorCode
							return lStatusFlag, lErr3
						} else {
							if lResp.Status == "success" {
								if len(lResp.Transactions) != 0 {
									for i := 0; i < len(lResp.Transactions); i++ {
										if lDRes[Idx].Symbol == lResp.Transactions[i].Symbol {
											//commented by Pavithra
											// log.Println(" lDRes[Idx].Symbol", lDRes[Idx].Symbol)
											// log.Println("lResp.Transactions[i].Symbol", lResp.Transactions[i].Symbol)
											lErr4 := UpdateHeaderAndDetails(lResp.Transactions[i], pBrokerId)
											if lErr4 != nil {
												log.Println("SSNSFS04", lErr4)
												lStatusFlag = common.ErrorCode
												return lStatusFlag, lErr4
											} else {
												lStatusFlag = common.SuccessCode
												// log.Println("Updated", lResp.Transactions[i])
											}
										}
									}
								}
							} else {
								// log.Println("Failed")
								lStatusFlag = common.ErrorCode
								return lStatusFlag, nil
							}
						}
					} else {
						// log.Println("Token Not Found", pBrokerId)
						lStatusFlag = common.ErrorCode
						return lStatusFlag, nil
					}
				}
			}
		} else {
			lStatusFlag = common.ErrorCode
			return lStatusFlag, nil
		}
	}
	log.Println("NseSgbFetchStatus (-)")
	return lStatusFlag, nil
}

func GetSgbEndDate(pApiRespRec []BiddingDateStruct) ([]BiddingDateStruct, error) {
	log.Println("GetSgbEndDate  (+)")

	var lGetResp BiddingDateStruct

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SGSEDO1", lErr1)
		return pApiRespRec, lErr1
	} else {
		defer lDb.Close()

		// added by pavithra
		lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
		lDay := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_Status_Fetch_Count"])

		// commented by pavithra,fetching all sgb bond added enddate condition
		// lCoreString := `SELECT asm.BiddingEndDate,asm.symbol
		//                 FROM a_sgb_master asm
		//                 WHERE asm.Exchange = 'NSE'
		// 				and asm.Redemption = 'N'
		//                 AND DATE(asm.BiddingEndDate) + INTERVAL ` + lDay + ` DAY >= CURDATE();`

		lCoreString := `SELECT asm.BiddingEndDate,asm.symbol
						FROM a_sgb_master asm
						WHERE asm.Exchange = 'NSE'
						and asm.Redemption = 'N'
						and asm.BiddingEndDate <= curdate() 
						and DATE(asm.BiddingEndDate) + INTERVAL ` + lDay + ` DAY >= CURDATE();`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("SGSED02", lErr2)
			return pApiRespRec, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lGetResp.Date, &lGetResp.Symbol)
				if lErr3 != nil {
					log.Println("SGSED03", lErr3)
					return pApiRespRec, lErr3
				} else {
					pApiRespRec = append(pApiRespRec, lGetResp)
				}
			}
		}
	}
	log.Println("GetSgbEndDate (-)")
	return pApiRespRec, nil
}

func UpdateHeaderAndDetails(pTransaction nsesgb.SgbAddResStruct, pBrokerId int) error {
	log.Println("UpdateHeaderAndDetails (+)")

	var lSBseResp bsesgb.SgbRespStruct
	var lSBseBid bsesgb.RespSgbBidStruct

	if pTransaction.OrderStatus != "" {
		if pTransaction.OrderStatus == "ES" || pTransaction.OrderStatus == "EF" {
			lSBseBid.ActionCode = "N"
		} else if pTransaction.OrderStatus == "MS" || pTransaction.OrderStatus == "MF" {
			lSBseBid.ActionCode = "M"
		} else {
			lSBseBid.ActionCode = "C"
		}
	}

	if pTransaction.Status != "" {
		if pTransaction.Status == "success" {
			lSBseResp.StatusCode = "0"
			lSBseBid.ErrorCode = "0"
		} else {
			lSBseResp.StatusCode = "1"
			lSBseBid.ErrorCode = "1"
		}
	}

	if pTransaction.OrderStatus == "ES" || pTransaction.OrderStatus == "MS" || pTransaction.OrderStatus == "CS" {
		lSBseBid.ErrorCode = "0"
	} else {
		lSBseBid.ErrorCode = "1"
	}

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSFUHD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// update a_sgb_orderdetails Table

		lSqlString1 := `update a_sgb_orderdetails d
						set d.AllotedUnit = ?,d.AllotedRate = ?,d.ErrorCode = ?,
						d.Message = ?, d.ModifiedDate = ?,d.AddedDate = ?,
						d.UpdatedBy = ?,d.UpdatedDate = now() 
						where d.RespOrderNo = ?`

		// Convert int to string
		strVal := strconv.Itoa(pTransaction.OrderNumber)

		_, lErr2 := lDb.Exec(lSqlString1, pTransaction.Quantity, pTransaction.Price, lSBseBid.ErrorCode, pTransaction.Reason, pTransaction.LastActionTime, pTransaction.EntryTime, common.AUTOBOT, strVal)
		if lErr2 != nil {
			log.Println("SSFUHD02", lErr2)
			return lErr2
		} else {
			// log.Println("Details updated Successfully")
			if pTransaction.VerificationReason == "" {
				lVerificationReason, lErr3 := GetSgbLookup(pTransaction.VerificationStatus)
				if lErr3 != nil {
					log.Println("SSFUHD03", lErr3)
					return lErr3
				} else {
					pTransaction.VerificationReason = lVerificationReason
				}
			}
			if pTransaction.ClearingReason == "" {
				lClearingReason, lErr4 := GetSgbLookup(pTransaction.ClearingStatus)
				if lErr4 != nil {
					log.Println("SSFUHD04", lErr4)
					return lErr4
				} else {
					pTransaction.ClearingReason = lClearingReason
				}
			}

			lSqlString2 := `update a_sgb_orderheader h
								set h.Status = ? ,h.StatusMessage = ? ,
								h.ErrorCode = ? , h.ErrorMessage = ?,
								h.DpStatus = ? ,h.DpRemarks = ?,
								h.RbiStatus = ?,h.RbiRemarks = ?,
								h.UpdatedBy = ?,h.UpdatedDate = now()
								where h.ClientReferenceNo = ?
								and h.ClientBenfId = ?
								and h.brokerid = ?
								and h.ScripId = ?`

			_, lErr5 := lDb.Exec(lSqlString2, lSBseResp.StatusCode, pTransaction.Reason, lSBseBid.ErrorCode, pTransaction.RejectionReason, pTransaction.VerificationStatus, pTransaction.VerificationReason, pTransaction.ClearingStatus, pTransaction.ClearingReason, common.AUTOBOT, pTransaction.ClientRefNumber, pTransaction.ClientBenId, pBrokerId, pTransaction.Symbol)
			if lErr5 != nil {
				log.Println("SSFUHD05", lErr5)
				return lErr5
			} else {
				log.Println("OrderHeader updated Successfully")
			}
		}
	}
	log.Println("UpdateHeaderAndDetails (-)")
	return nil
}

func GetSgbLookup(pReqCode string) (string, error) {
	log.Println("GetSgbLookup (+)")

	var lString string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSLTO1", lErr1)
		return lString, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select nvl(d.description,'') description  
						from xx_lookup_header h,xx_lookup_details d
						where h.id = d.headerid 
						and h.Code = 'SgbStatus'
						and d.Code = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pReqCode)
		if lErr2 != nil {
			log.Println("SSLTO2", lErr2)
			return lString, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lString)
				if lErr3 != nil {
					log.Println("SSLTO3", lErr3)
					return lString, lErr3
				}
			}
		}
	}
	log.Println("GetSgbLookup (-)")
	return lString, nil
}
