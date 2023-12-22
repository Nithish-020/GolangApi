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
func NseSgbFetchStatus(pBrokerId int, pUser string) error {
	log.Println("NseSgbFetchStatus (+)")

	//commented by pavithra
	// defer lwg.Done()
	var lDateResp []BiddingDateStruct

	lDRes, lErr1 := GetSgbEndDate(lDateResp)
	if lErr1 != nil {
		log.Println("SSNSFS01", lErr1)
		return lErr1
	} else {
		log.Println("lDRes", lDRes)
		//dateloop
		for Idx := 0; Idx < len(lDRes); Idx++ {
			//bsetoken
			lToken, lErr2 := exchangecall.GetToken(pUser, pBrokerId)
			if lErr2 != nil {
				log.Println("SSNSFS02", lErr2)
				return lErr2
			} else {
				if lToken != "" {
					lResp, lErr3 := nsesgb.SgbTransactionsMaster(lToken, lDRes[Idx].Date, pUser)
					if lErr3 != nil {
						log.Println("SSNSFS03", lErr3)
						return lErr3
					} else {
						for i := 0; i < len(lResp.Transactions); i++ {
							if lDRes[Idx].Symbol == lResp.Transactions[i].Symbol {
								log.Println(" lDRes[Idx].Symbol", lDRes[Idx].Symbol)
								log.Println("lResp.Transactions[i].Symbol", lResp.Transactions[i].Symbol)
								lErr4 := UpdatedSgbConstruct(lResp.Transactions[i], pBrokerId)
								if lErr4 != nil {
									log.Println("SSNSFS04", lErr4)
									return lErr4
								} else {
									log.Println("Updated", lResp.Transactions[i])
								}
							}
						}
					}
				}
			}
		}
	}
	log.Println("NseSgbFetchStatus (-)")
	return nil
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
		lConfigFile := common.ReadTomlConfig("toml/debug.toml")
		lDay := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_Status_Fetch_Count"])

		// commented by pavithra added redemption flag in where condition
		// lCoreString := `SELECT asm.BiddingEndDate,asm.symbol
		//                 FROM a_sgb_master asm
		//                 WHERE asm.Exchange = 'NSE'
		//                 AND DATE(asm.BiddingEndDate) + INTERVAL ` + lDay + ` DAY >= CURDATE();`

		lCoreString := `SELECT asm.BiddingEndDate,asm.symbol
		                FROM a_sgb_master asm
		                WHERE asm.Exchange = 'NSE'
						and asm.Redemption = 'N'
		                AND DATE(asm.BiddingEndDate) + INTERVAL ` + lDay + ` DAY >= CURDATE();`
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

func UpdatedSgbConstruct(pTransaction nsesgb.SgbAddResStruct, pBrokerId int) error {
	log.Println("UpdatedSgbConstruct (+)")

	var lSBseResp bsesgb.SgbRespStruct
	var lSBseBid bsesgb.RespSgbBidStruct

	if pTransaction.OrderStatus == "ES" || pTransaction.OrderStatus == "EF" {
		lSBseBid.ActionCode = "N"
	} else if pTransaction.OrderStatus == "MS" || pTransaction.OrderStatus == "MF" {
		lSBseBid.ActionCode = "M"
	} else {
		lSBseBid.ActionCode = "C"
	}

	if pTransaction.Status == "success" {
		lSBseResp.StatusCode = "0"
		lSBseBid.ErrorCode = "0"
	} else {
		lSBseResp.StatusCode = "1"
		lSBseBid.ErrorCode = "1"
	}

	if pTransaction.OrderStatus == "ES" || pTransaction.OrderStatus == "MS" || pTransaction.OrderStatus == "CS" {
		lSBseBid.ErrorCode = "0"
	} else {
		lSBseBid.ErrorCode = "1"
	}

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SUSCO1", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// update a_sgb_orderdetails Table

		lSqlString1 := `update a_sgb_orderdetails d
		set d.RespSubscriptionunit = ?,d.RespRate = ?,d.ErrorCode = ?,d.Message = ?, d.ModifiedDate = ?,d.AddedDate = ?
		where d.RespOrderNo =?`

		// commenetd by pavithra
		// intVal := pTransaction.OrderNumber

		// Convert int to string
		strVal := strconv.Itoa(pTransaction.OrderNumber)

		_, lErr2 := lDb.Exec(lSqlString1, pTransaction.Quantity, pTransaction.Price, pTransaction.Status, pTransaction.Reason, pTransaction.LastActionTime, pTransaction.EntryTime, strVal)
		if lErr2 != nil {
			log.Println("SUSCO2", lErr2)
			return lErr2
		} else {
			log.Println("Details updated Successfully")

			llookup, lErr3 := SgbLookTransaction(pTransaction)
			if lErr3 != nil {
				log.Println("SUSCO3", lErr3)
				return lErr3
			} else {
				// COMMENTED BY PAVITHRA
				// log.Println("llookup", llookup)

				lSqlString2 := `update a_sgb_orderheader h
			                 set h.ScripId = ? ,h.Depository =?  , h.DpId = ? ,  h.PanNo = ?, h.Status = ? , h.StatusMessage = ? , h.ErrorCode = ? , h.ErrorMessage = ?,   h.DpStatus = ? , h.DpRemarks = ?
			          where h.ClientReferenceNo= ? and h.brokerid= ?`

				_, lErr4 := lDb.Exec(lSqlString2, pTransaction.Symbol, pTransaction.Depository, pTransaction.DpId, pTransaction.Pan, pTransaction.Status, pTransaction.Reason, pTransaction.Status, pTransaction.Reason, llookup.VerificationStatus, llookup.VerificationReason, pTransaction.ClientRefNumber, pBrokerId)
				if lErr4 != nil {
					log.Println("SUSCO4", lErr4)
					return lErr4
				} else {
					log.Println("OrderHeader Details updated Successfully")
				}
			}
		}
	}
	log.Println("UpdatedSgbConstruct (-)")
	return nil
}

func SgbLookTransaction(pTransaction nsesgb.SgbAddResStruct) (nsesgb.SgbAddResStruct, error) {
	log.Println("SgbLookTransaction (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSLTO1", lErr1)
		return pTransaction, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select d.description  
						from xx_lookup_header h,xx_lookup_details d
				        where h.id = d.headerid 
						and h.Code = 'Sgbstatus'
						and d.Code = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pTransaction.ClearingStatus)
		if lErr2 != nil {
			log.Println("SSLTO2", lErr2)
			return pTransaction, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&pTransaction.ClearingStatus)
				if lErr3 != nil {
					log.Println("SSLTO3", lErr3)
					return pTransaction, lErr3
				}
			}
		}
	}
	log.Println("SgbLookTransaction (-)")
	return pTransaction, nil
}
