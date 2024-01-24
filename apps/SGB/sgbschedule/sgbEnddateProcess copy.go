package sgbschedule

import (
	"bytes"
	"fcs23pkg/apps/SGB/sgbplaceorder"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bsesgb"
	"fcs23pkg/integration/nse/nsesgb"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type FailedJvStruct struct {
	SINo        int
	ClientId    string
	Amount      string
	Transaction string
}

// commented by pavithra
// type SgbIssueList struct {
// 	Symbol                string `json:"symbol"`
// 	TotalReqCount         int    `json:"totalReqCount"`
// 	FoVerifySuccessCount  int    `json:"foVerifySuccessCount"`
// 	FoVerifyFailedCount   int    `json:"foVerifyFailedCount"`
// 	JvSuccesssCount       int    `json:"jvSuccesssCount"`
// 	JvFailedCount         int    `json:"jvFailedCount"`
// 	ExchangeSuccessCount  int    `json:"exchangeSuccessCount"`
// 	ExchangeFailedCount   int    `json:"exchangeFailedCount"`
// 	ReverseJvSuccessCount int    `json:"reverseJvSuccessCount"`
// 	ReverseJvFailedCount  int    `json:"reverseJvFailedCount"`
// 	FoJvSuccessCount      int    `json:"foJvSuccessCount"`
// 	FoJvFailedCount       int    `json:"foJvFailedCount"`
// 	TotalSuccessCount     int    `json:"totalSuccessCount"`
// 	TotalFailedCount      int    `json:"totalFailedCount"`
// }
type SgbIssueList struct {
	Symbol                  string `json:"symbol"`
	TotalReqCount           int    `json:"totalReqCount"`
	FoVerifySuccessCount    int    `json:"foVerifySuccessCount"`
	FoVerifyFailedCount     int    `json:"foVerifyFailedCount"`
	JvSuccesssCount         int    `json:"jvSuccesssCount"`
	JvFailedCount           int    `json:"jvFailedCount"`
	ExchangeSuccessCount    int    `json:"exchangeSuccessCount"`
	ExchangeFailedCount     int    `json:"exchangeFailedCount"`
	ReverseBoJvSuccessCount int    `json:"reverseBoJvSuccessCount"`
	ReverseBoJvFailedCount  int    `json:"reverseBoJvFailedCount"`
	ReverseFoJvSuccessCount int    `json:"reverseFoJvSuccessCount"`
	ReverseFoJvFailedCount  int    `json:"reverseFoJvFailedCount"`
	FoJvSuccessCount        int    `json:"foJvSuccessCount"`
	FoJvFailedCount         int    `json:"foJvFailedCount"`
	TotalSuccessCount       int    `json:"totalSuccessCount"`
	TotalFailedCount        int    `json:"totalFailedCount"`
	InsufficientFund        int    `json:"insufficientFund"`
}

type JvStatusStruct struct {
	BoJvStatus    string `json:"boJvstatus"`
	BoJvStatement string `json:"boJvstatement"`
	BoJvAmount    string `json:"boJvamount"`
	BoJvType      string `json:"boJvtype"`
	FoJvStatus    string `json:"foJvstatus"`
	FoJvStatement string `json:"foJvstatement"`
	FoJvAmount    string `json:"foJvamount"`
	FoJvType      string `json:"foJvtype"`
}

//----------------------------------------------------------------
// this method processing the SGB orders on SGB enddate
//----------------------------------------------------------------

func PlacingSgbOrder(r *http.Request) ([]SgbIssueList, error) {
	log.Println("PlacingSgbOrder (+)")

	var lsgbIssue SgbIssueList
	var lSgbIssuesList []SgbIssueList
	lValidSgbArr, lErr1 := GoodTimeForApplySgb()
	if lErr1 != nil {
		log.Println("SPSO01", lErr1)
		return lSgbIssuesList, lErr1
	} else {
		log.Println("SGbValid Bond array", lValidSgbArr)

		//  if no issue
		if lValidSgbArr != nil {
			for lSgbIdx := 0; lSgbIdx < len(lValidSgbArr); lSgbIdx++ {
				// for lSgbIdx, lsgbValues := range lValidSgbArr {

				if len(lValidSgbArr) != 0 {

					lBrokerIdArr, lErr2 := GetSgbBrokers(lValidSgbArr[lSgbIdx].Exchange)
					if lErr2 != nil {
						log.Println("SPSO02", lErr2)
						return lSgbIssuesList, lErr2
					} else {
						log.Println("SGB Broker array", lBrokerIdArr)
						if len(lBrokerIdArr) != 0 {
							// var lWg sync.WaitGroup
							for lBrokerIdx := 0; lBrokerIdx < len(lBrokerIdArr); lBrokerIdx++ {
								// lWg.Add(1)
								lCountStruct, lErr3 := ProcessSgbOrder(lValidSgbArr[lSgbIdx], lBrokerIdArr[lBrokerIdx], r)
								if lErr3 != nil {
									log.Println("SPSO03", lErr3)
									lCountStruct.Symbol = lErr3.Error()
									lSgbIssuesList = append(lSgbIssuesList, lCountStruct)
								} else {
									lCountStruct.Symbol = lValidSgbArr[lSgbIdx].Symbol
									lSgbIssuesList = append(lSgbIssuesList, lCountStruct)
								}
							}
							// lWg.Wait()
						} else {
							lsgbIssue.Symbol = "No Brokers to Process"
							lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
						}
					}
				} else {
					lsgbIssue.Symbol = lValidSgbArr[lSgbIdx].Symbol
					lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
				}
			}
		} else {
			log.Println("No SGB issuses ending today")
			lsgbIssue.Symbol = "No SGB issuses ending today"
			lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
		}
	}
	log.Println("PlacingSgbOrder (-)")
	return lSgbIssuesList, nil
}

// ----------------------------------------------------------------
// this method is used to get the valid sgb record from database
// ----------------------------------------------------------------
func GoodTimeForApplySgb() ([]SgbStruct, error) {
	log.Println("GoodTimeForApplySgb (+)")

	var lValidSgbArr []SgbStruct
	var lValidSgbRec SgbStruct

	config := common.ReadTomlConfig("toml/SgbConfig.toml")
	lApplyFlag := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_ApplyDay_Flag"])
	log.Println("lApplyFlag", lApplyFlag)

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SGTAS01", lErr1)
		return lValidSgbArr, lErr1
	} else {
		defer lDb.Close()
		lCondition := ""
		lCoreString := `select m.id,m.Symbol ,m.Exchange
						from a_sgb_master m
						where m.Redemption = 'N'`

		if lApplyFlag == "Y" {
			lCondition = `and m.BiddingEndDate >= curdate()`
		} else {
			lCondition = `and m.BiddingEndDate = curdate()`
		}

		lFinalQueryString := lCoreString + lCondition
		// select m.id,m.Symbol ,m.Exchange
		// from a_sgb_master m
		// where m.BiddingEndDate = curdate()
		// and m.DailyEndTime > now()
		lRows, lErr2 := lDb.Query(lFinalQueryString)
		if lErr2 != nil {
			log.Println("SGTAS02", lErr2)
			return lValidSgbArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lValidSgbRec.MasterId, &lValidSgbRec.Symbol, &lValidSgbRec.Exchange)
				if lErr3 != nil {
					log.Println("SGTAS03", lErr3)
					return lValidSgbArr, lErr3
				} else {
					lValidSgbArr = append(lValidSgbArr, lValidSgbRec)
				}
			}
		}
	}
	log.Println("GoodTimeForApplySgb (-)")
	return lValidSgbArr, nil
}

// ----------------------------------------------------------------
// this method is used to get the Brokerid from database
// ----------------------------------------------------------------
func GetSgbBrokers(pExchange string) ([]int, error) {
	log.Println("GetSgbBrokers(+)")

	var lBrokerArr []int
	var lBrokerRec int
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GSB01", lErr1)
		return lBrokerArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select md.BrokerId  
						from a_ipo_memberdetails md
						where md.AllowedModules  like '%Sgb%'
						and md.Flag = 'Y'
						and md.OrderPreference = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pExchange)
		if lErr2 != nil {
			log.Println("GSB02", lErr2)
			return lBrokerArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lBrokerRec)
				if lErr3 != nil {
					log.Println("GSB03", lErr3)
					return lBrokerArr, lErr3
				} else {
					lBrokerArr = append(lBrokerArr, lBrokerRec)
				}
			}
		}
	}
	log.Println("GetSgbBrokers(-)")
	return lBrokerArr, nil
}

func ProcessSgbOrder(pValidSgb SgbStruct, pBrokerId int, pR *http.Request) (SgbIssueList, error) {
	log.Println("ProcessSgbOrder (+)")

	var lCountRec SgbIssueList

	lSgbReqArr, lErr1 := fetchSgbOrder(pValidSgb, pBrokerId)
	lCountRec.TotalReqCount = len(lSgbReqArr)
	if lErr1 != nil {
		log.Println("PSO01", lErr1)
		return lCountRec, lErr1
	} else {

		log.Println("lSgbReqArr", lSgbReqArr)
		log.Println("lCountRec.TotalReqCount ", lCountRec.TotalReqCount)
		for lSqlReqIdx := 0; lSqlReqIdx < len(lSgbReqArr); lSqlReqIdx++ {

			pResponseRec, lErr2 := PostJvandProcessOrder(lSgbReqArr[lSqlReqIdx], pR, pValidSgb, pBrokerId)
			log.Println("lCountRec before", lCountRec)
			log.Println("pResponseRec", pResponseRec)
			lCountRec.FoVerifySuccessCount += pResponseRec.FoVerifySuccessCount
			lCountRec.FoVerifyFailedCount += pResponseRec.FoVerifyFailedCount
			lCountRec.JvSuccesssCount += pResponseRec.BoJvSuccesssCount
			lCountRec.JvFailedCount += pResponseRec.BoJvFailedCount
			lCountRec.ExchangeSuccessCount += pResponseRec.ExchangeSuccessCount
			lCountRec.ExchangeFailedCount += pResponseRec.ExchangeFailedCount
			lCountRec.ReverseBoJvSuccessCount += pResponseRec.ReverseBoJvSuccessCount
			lCountRec.ReverseBoJvFailedCount += pResponseRec.ReverseBoJvFailedCount
			lCountRec.ReverseFoJvSuccessCount += pResponseRec.ReverseFoJvSuccessCount
			lCountRec.ReverseFoJvFailedCount += pResponseRec.ReverseFoJvFailedCount
			lCountRec.FoJvSuccessCount += pResponseRec.FoJvSuccessCount
			lCountRec.FoJvFailedCount += pResponseRec.FoJvFailedCount
			lCountRec.TotalSuccessCount += pResponseRec.TotalSuccessCount
			lCountRec.TotalFailedCount += pResponseRec.TotalFailedCount
			lCountRec.InsufficientFund += pResponseRec.InsufficientFund
			log.Println("lCountRec after", lCountRec)
			if lErr2 != nil || pResponseRec.ProcessStatus == common.ErrorCode {
				log.Println("PSO02", lErr1)
				lErr3 := UpdateResponse(pResponseRec, pValidSgb.Exchange, pBrokerId)
				if lErr3 != nil {
					log.Println("PSO03", lErr3)
					// return lCountRec, lErr3
				} else {
					lErr3 := MailSendingProcess(pResponseRec)
					if lErr3 != nil {
						log.Println("PSO04", lErr3)
						// return lCountRec, lErr3
					}
				}
			} else {
				lErr5 := UpdateResponse(pResponseRec, pValidSgb.Exchange, pBrokerId)
				if lErr5 != nil {
					log.Println("PSO05", lErr5)
					// return lCountRec, lErr5
				} else {
					lErr6 := MailSendingProcess(pResponseRec)
					if lErr6 != nil {
						log.Println("PSO06", lErr6)
						// return lCountRec, lErr6
					}
				}
			}
		}
	}
	log.Println("ProcessSgbOrder (-)")
	return lCountRec, nil
}

// ----------------------------------------------------------------
// this method is used to fetch the local placed orders from database
// ----------------------------------------------------------------
func fetchSgbOrder(pValidSgbRec SgbStruct, BrokerId int) ([]exchangecall.JvReqStruct, error) {
	log.Println("fetchSgbOrder (+)")

	var lReqSgbBid exchangecall.SgbBidStruct
	var lReqSgbArr []exchangecall.JvReqStruct
	var lReqSgbData exchangecall.JvReqStruct

	config := common.ReadTomlConfig("toml/SgbConfig.toml")
	lConfigTime := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_SchConfig_time"])
	lSpecificClient := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_Specific_Client"])
	lApplyFlag := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_ApplyDay_Flag"])
	lProcessFlag := fmt.Sprintf("%v", config.(map[string]interface{})["Process_SGB_Flag"])
	lScheduleFlag := fmt.Sprintf("%v", config.(map[string]interface{})["Schedule_SGB_Flag"])

	// Parse the string into a time.Time object
	lTime, lErr1 := time.Parse("15:04:05", lConfigTime)
	if lErr1 != nil {
		log.Println("SSFP01", lErr1)
		return lReqSgbArr, lErr1
	} else {
		// Get current date
		currentDate := time.Now().Local()

		// Set the date component of lTime to today's date
		lConfigTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), lTime.Hour(), lTime.Minute(), lTime.Second(), 0, time.Local)
		lConfigUnixTime := lConfigTime.Unix()
		//current date unix time value
		lCurrentTime := time.Now()
		lCurrentTimeUnix := lCurrentTime.Unix()

		if lCurrentTimeUnix > lConfigUnixTime {
			log.Println("Time Validation", lCurrentTimeUnix, lConfigUnixTime)
			// Perform actions if the current time is greater than the configured time for today

			lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
			if lErr1 != nil {
				log.Println("SSFP02", lErr1)
				return lReqSgbArr, lErr1
			} else {
				defer lDb.Close()

				//commented by pavithra
				// lSqlString := `	select h.ScripId ,h.PanNo,h.InvestorCategory,h.ApplicantName,h.Depository,h.DpId,h.ApplicantName ,date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date,h.ClientEmail ,m.Symbol ,h.ClientBenfId ,h.GuardianName ,h.GuardianPanNo ,h.GuardianRelation,d.BidId ,
				// d.ReqSubscriptionUnit ,d.ReqRate ,d.OrderNo ,d.ActionCode,h.ClientId,d.ReqSubscriptionUnit,d.ReqRate,d.OrderNo
				// (CASE
				// 	WHEN d.ActionCode = 'M' THEN 'Modify'
				// 	WHEN d.ActionCode = 'N' THEN 'New'
				// 	WHEN d.ActionCode = 'D' THEN 'Delete'
				// 	ELSE d.ActionCode
				// 	END ) AS ActionDescription
				// from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
				// where h.MasterId = m.id
				// and d.HeaderId = h.Id
				// and h.Status = "success"
				// and m.BiddingStartDate <= curdate()
				// and m.BiddingEndDate >= curdate()
				// and time(now()) between m.DailyStartTime and m.DailyEndTime
				// and h.cancelFlag != 'Y'
				// and m.Exchange = ?
				// and h.brokerId = ?
				// and m.id = ?`

				lSqlString := `select h.ScripId ,
								h.PanNo,
								h.InvestorCategory,
								h.ApplicantName,
								h.Depository,
								h.DpId,
								date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date,h.ClientEmail ,
								h.ClientBenfId ,
								h.GuardianName ,
								h.GuardianPanNo ,
								h.GuardianRelation,
								d.BidId ,
								d.ReqSubscriptionUnit ,
								d.ReqRate ,
								d.ReqOrderNo ,
								d.ActionCode,
								h.ClientId,
								d.ReqOrderNo,
								(CASE
									WHEN d.ActionCode = 'M' THEN 'Modify'
									WHEN d.ActionCode = 'N' THEN 'New'
									WHEN d.ActionCode = 'D' THEN 'Delete'
									ELSE d.ActionCode
									END ) AS ActionDescription
								from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
								where h.MasterId = m.id
								and d.HeaderId = h.Id
								and h.Status = 'success'
								and m.BiddingStartDate <= curdate()
								and time(now()) between m.DailyStartTime and m.DailyEndTime
								and h.cancelFlag != 'Y'
								and m.Exchange = ?
								and h.brokerId = ?
								and m.id = ?`

				lCondition := ""
				if lApplyFlag == "Y" {
					lCondition = ` and m.BiddingEndDate >= curdate()`
				} else {
					lCondition = ` and m.BiddingEndDate = curdate()`
				}

				lInclause := ""
				if lSpecificClient != "" {
					lInclause = `and h.ClientId in ('` + lSpecificClient + `')`
				}
				lProcessFlagString := ` and h.ProcessFlag = '` + lProcessFlag + `' and h.ScheduleStatus = '` + lScheduleFlag + `'`

				lFinalQueryString := lSqlString + lProcessFlagString + lCondition + lInclause
				log.Println("lFinalQueryString", lFinalQueryString)

				lRows, lErr2 := lDb.Query(lFinalQueryString, pValidSgbRec.Exchange, BrokerId, pValidSgbRec.MasterId)
				if lErr2 != nil {
					log.Println("SSFP02", lErr2)
					return lReqSgbArr, lErr2
				} else {
					for lRows.Next() {
						lErr3 := lRows.Scan(&lReqSgbData.Symbol, &lReqSgbData.PanNo, &lReqSgbData.InvestorCategory, &lReqSgbData.ApplicantName, &lReqSgbData.Depository, &lReqSgbData.DpId, &lReqSgbData.OrderDate, &lReqSgbData.Mail, &lReqSgbData.ClientBenfId, &lReqSgbData.GuardianName, &lReqSgbData.GuardianPanno, &lReqSgbData.GuardianRelation, &lReqSgbBid.BidId, &lReqSgbBid.SubscriptionUnit, &lReqSgbBid.Rate, &lReqSgbBid.OrderNo, &lReqSgbBid.ActionCode, &lReqSgbData.ClientId, &lReqSgbData.ReqOrderNo, &lReqSgbData.ActivityType)
						if lErr3 != nil {
							log.Println("SSFP03", lErr3)
							return lReqSgbArr, lErr3
						} else {
							if lProcessFlag == "N" && lScheduleFlag == "N" {
								lReqSgbBid.ActionCode = "N"
							}
							lReqSgbData.Unit, lErr3 = strconv.Atoi(lReqSgbBid.SubscriptionUnit)
							if lErr3 != nil {
								log.Println("SSFP04", lErr3)
								return lReqSgbArr, lErr3
							}
							lReqSgbData.Price, lErr3 = strconv.Atoi(lReqSgbBid.Rate)
							if lErr3 != nil {
								log.Println("SSFP05", lErr3)
								return lReqSgbArr, lErr3
							}
							lReqSgbData.Amount = strconv.Itoa(lReqSgbData.Unit * lReqSgbData.Price)
							lReqSgbData.ClientName = lReqSgbData.ApplicantName

							lReqSgbData.Bids = append(lReqSgbData.Bids, lReqSgbBid)
							lReqSgbArr = append(lReqSgbArr, lReqSgbData)
							lReqSgbData.Bids = []exchangecall.SgbBidStruct{}
						}
					}
				}
			}
		} else {
			return lReqSgbArr, common.CustomError("The time for Configuration not yet come to process applications.")
		}
	}
	log.Println("fetchSgbOrder (-)")
	return lReqSgbArr, nil
}

//---------------------------------------------------------------------------------
// this method update the whole response
//---------------------------------------------------------------------------------

func UpdateResponse(pJvRec exchangecall.JvReqStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateResponse (+)")
	lErr1 := UpdateHeader2(pJvRec, pExchange, pBrokerId)
	if lErr1 != nil {
		log.Println("SEPUR01", lErr1)
		return lErr1
	} else {
		log.Println("Response Record Updated Successfully")
	}

	log.Println("UpdateResponse (-)")
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
Date: 03jan2024
*/
func UpdateHeader2(pRespRec exchangecall.JvReqStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateHeader2 (+)")
	var lStatus string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SEPUH201", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lHeaderId, lErr2 := GetOrderId(pRespRec, pExchange)
		if lErr2 != nil {
			log.Println("SEPUH202", lErr2)
			return lErr2
		} else {
			if lHeaderId != 0 {
				if pRespRec.StatusCode != "0" {
					lStatus = common.FAILED
				} else {
					lStatus = common.SUCCESS
				}
				lSqlString := `update a_sgb_orderheader h
								set h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,
								h.Depository = ?,h.DpId = ?,h.ClientBenfId = ?,h.GuardianName = ?,
								h.GuardianPanNo = ?,h.GuardianRelation = ?,h.ClientReferenceNo = ?,
								h.StatusCode = ?,h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,
								h.Status = ?,h.ProcessFlag = ?,h.ScheduleStatus = ?,
								h.UpdatedBy = ?,h.UpdatedDate = now()
								where h.ClientId = ?
								and h.brokerId = ?
								and h.Id = ?`

				_, lErr3 := lDb.Exec(lSqlString, pRespRec.PanNo, pRespRec.InvestorCategory, pRespRec.ApplicantName,
					pRespRec.Depository, pRespRec.DpId, pRespRec.ClientBenfId, pRespRec.GuardianName, pRespRec.GuardianPanno, pRespRec.GuardianRelation, pRespRec.ClientRefNumber, pRespRec.StatusCode, pRespRec.StatusMessage, pRespRec.ErrorCode, pRespRec.ErrorMessage, lStatus, pRespRec.ProcessStatus, pRespRec.ErrorStage, common.AUTOBOT, pRespRec.ClientId, pBrokerId, lHeaderId)
				if lErr3 != nil {
					log.Println("SEPUH203", lErr3)
					return lErr3
				} else {
					// call InsertDetails method to inserting the order details in order details table
					lErr4 := UpdateDetail2(pRespRec, lHeaderId)
					if lErr4 != nil {
						log.Println("SEPUH204", lErr4)
						return lErr4
					}
				}
			} else {
				return common.CustomError("Record Not Found")
			}
		}
	}
	log.Println("UpdateHeader2 (-)")
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
Date: 03JAN2024
*/
func UpdateDetail2(pRespRec exchangecall.JvReqStruct, pHeaderId int) error {
	log.Println("UpdateDetail2 (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SEPUD201", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		for lIdx := 0; lIdx < len(pRespRec.Bids); lIdx++ {

			lSqlString := `update a_sgb_orderdetails d
							set d.RespApplicationNo = ?,d.RespOrderNo = ?,d.ActionCode = ?,
							d.RespSubscriptionunit = ?,d.RespRate = ?,
							d.ErrorCode = ?,d.Message = ?,
							d.BoJvStatus = ?,d.BoJvAmount = ?,
							d.BoJvStatement = ?,d.BoJvType = ?,
							d.FoJvStatus = ?,d.FoJvAmount = ?,
							d.FoJvStatement = ?,d.FoJvType = ?,
							d.AddedDate = ?,d.ModifiedDate = ?,
							d.UpdatedBy = ?,d.UpdatedDate = now()
							where d.HeaderId = ?
							and d.ReqOrderNo = ?`

			_, lErr2 := lDb.Exec(lSqlString, pRespRec.ApplicationNumber, pRespRec.Bids[lIdx].OrderNo, pRespRec.Bids[lIdx].ActionCode, pRespRec.Bids[lIdx].SubscriptionUnit, pRespRec.Bids[lIdx].Rate, pRespRec.Bids[lIdx].ErrorCode, pRespRec.Bids[lIdx].Message, pRespRec.BoJvStatus, pRespRec.BoJvAmount, pRespRec.BoJvStatement, pRespRec.BoJvType, pRespRec.FoJvStatus, pRespRec.FoJvAmount, pRespRec.FoJvStatement, pRespRec.FoJvType, pRespRec.EntryTime, pRespRec.LastActionTime, common.AUTOBOT, pHeaderId, pRespRec.ReqOrderNo)
			if lErr2 != nil {
				log.Println("SEPUD202", lErr2)
				return lErr2
			}
		}
	}
	log.Println("UpdateDetail2 (-)")
	return nil
}

// ---------------------------------------------------------------------------------
// this method is used to update the reocrds od NSE in db
// ---------------------------------------------------------------------------------
func UpdateNseRecord(pRespRec nsesgb.SgbAddResStruct, pJvRec exchangecall.JvReqStruct, pBseConvertedRec bsesgb.SgbRespStruct, pBrokerId int) error {
	log.Println("UpdateNseRecord (+)")
	// for lJvReqIdx := 0; lJvReqIdx < len(pJvRecArr); lJvReqIdx++ {
	// if strconv.Itoa(pRespRec.OrderNumber) == pJvRecArr[lJvReqIdx].OrderNo {
	lErr1 := UpdateHeader(pRespRec, pJvRec, pBseConvertedRec, common.NSE, pBrokerId)
	if lErr1 != nil {
		log.Println("SGBOFUNR01", lErr1)
		return lErr1
	} else {
		log.Println("Nse Record Updated Successfully")
	}
	// }
	// }
	log.Println("UpdateNseRecord (-)")
	return nil
}

//---------------------------------------------------------------------------------
// this method is used to update the reocrds od BSE in db
//---------------------------------------------------------------------------------
func UpdateBseRecord(pRespRec bsesgb.SgbRespStruct, pJvRec exchangecall.JvReqStruct, pBrokerId int) error {
	log.Println("UpdateBseRecord (+)")
	// for lJvReqIdx := 0; lJvReqIdx < len(pJvRecArr); lJvReqIdx++ {
	for lBidIdx := 0; lBidIdx < len(pRespRec.Bids); lBidIdx++ {
		if pRespRec.Bids[lBidIdx].OrderNo == pJvRec.ReqOrderNo {
			var lEmptyRec nsesgb.SgbAddResStruct
			lErr1 := UpdateHeader(lEmptyRec, pJvRec, pRespRec, common.BSE, pBrokerId)
			if lErr1 != nil {
				log.Println("SGBOFUBR01", lErr1)
				return lErr1
			} else {
				log.Println("Bse Record Updated Successfully")
			}
		}
	}
	// }
	log.Println("UpdateBseRecord (-)")
	return nil
}

func MailSendingProcess(pSgbRec exchangecall.JvReqStruct) error {
	log.Println("MailSendingProcess (+)")

	if pSgbRec.ProcessStatus == common.ErrorCode {
		lEmailRec, lErr1 := ConstructMailContent(pSgbRec)
		if lErr1 != nil {
			log.Println("SMSP01", lErr1)
			return lErr1
		} else {
			lErr2 := emailUtil.SendEmail(lEmailRec, "SGB")
			if lErr2 != nil {
				log.Println("SMSP02", lErr2)
				return lErr2
			}
		}
	} else if pSgbRec.ProcessStatus == "Y" {
		lEmailRec, lErr3 := ConstructSuccessmail(pSgbRec)
		if lErr3 != nil {
			log.Println("SMSP03", lErr3)
			return lErr3
		} else {
			lErr4 := emailUtil.SendEmail(lEmailRec, "SGB")
			if lErr4 != nil {
				log.Println("SMSP04", lErr4)
				return lErr4
			}
		}
	} else {
		return common.CustomError("Unable to Find the Issue")
	}

	log.Println("MailSendingProcess (-)")
	return nil
}

func ConstructMailContent(pSgbRec exchangecall.JvReqStruct) (emailUtil.EmailInput, error) {
	log.Println("ConstructMailContent (+)")

	type dynamicEmailStruct struct {
		Date     string `json:"date"`
		Name     string `json:"name"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Footer   string `json:"footer"`
		ClientId string `json:"clientId"`
		OrderNo  string `json:"orderNo"`
		JvUnit   string `json:"jvUnit"`
		JvAmount string `json:"jvAmount"`
	}
	var lJVEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("toml/emailconfig.toml")
	lJVEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lJVEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lJVEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	// Mail Id for It support is Not Added in Toml
	// lJVEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["ToEmailId"])

	lNovoConfig := common.ReadTomlConfig("toml/novoConfig.toml")
	lJVEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_ClientEmail_Subject"])
	htmlTemplate := fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_CommonTemplate"])

	lSGBconfig := common.ReadTomlConfig("toml/SgbConfig.toml")
	lEmailIdCode := fmt.Sprintf("%v", lSGBconfig.(map[string]interface{})["SGB_EmailId_LookupCode"])
	lEmailContentCode := fmt.Sprintf("%v", lSGBconfig.(map[string]interface{})["SGB_EmailContent_LookupCode"])

	var lDynamicEmailVal dynamicEmailStruct

	lEmailId, lErr1 := GetEmailContent(lEmailIdCode, pSgbRec.EmailDist)
	if lErr1 != nil {
		log.Println("SSCFJ01", lErr1)
		return lJVEmailContent, lErr1
	} else {
		lJVEmailContent.ToEmailId = lEmailId

		lEmailContent, lErr2 := GetEmailContent(lEmailContentCode, pSgbRec.ErrorStage)
		if lErr2 != nil {
			log.Println("SSCFJ01", lErr2)
			return lJVEmailContent, lErr2
		} else {
			lDynamicEmailVal.Content = lEmailContent

			lTemp, lErr := template.ParseFiles(htmlTemplate)
			if lErr != nil {
				log.Println("SSCFJ02", lErr)
				return lJVEmailContent, lErr
			} else {
				var lTpl bytes.Buffer
				currentTime := time.Now()
				currentDate := currentTime.Format("02-01-2006")

				lDynamicEmailVal.Date = currentDate
				lDynamicEmailVal.ClientId = pSgbRec.ClientId
				lDynamicEmailVal.OrderNo = pSgbRec.ReqOrderNo
				lDynamicEmailVal.JvUnit = strconv.Itoa(pSgbRec.Unit)
				lDynamicEmailVal.JvAmount = pSgbRec.Amount
				lDynamicEmailVal.Name = "Team"
				switch pSgbRec.ErrorStage {
				case "V":
					//  IT
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER FO ERROR REPORT"
					lDynamicEmailVal.Footer = "We kindly request you to Identify the solution on Processing of verifying Client's Account Details."
				case "F":
					// IT, RMS
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER FO FAILED REPORT"
					lDynamicEmailVal.Footer = "We kindly request you to Identify the solution on Processing of Deducting Amount Client's Account Details."
				case "B":
					// IT, ACC
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER BO FAILED REPORT"
					lDynamicEmailVal.Footer = "We kindly request you to Identify the solution on Processing of Deducting Amount Client's Account Details."
				case "E":
					// IT
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."
				case "A":
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."
				case "R":
					//IT,RMS,ACC
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."
				case "H":
					// IT,RMS
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."

				case "I":
					// client
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "SGB ORDER DETAILS"
					lDynamicEmailVal.Footer = "Please maintain an sufficient fund balance from your next SGB orders."
					lDynamicEmailVal.Name = pSgbRec.ClientName
					lJVEmailContent.ToEmailId = pSgbRec.Mail
				}

				lTemp.Execute(&lTpl, lDynamicEmailVal)
				lEmailbody := lTpl.String()

				lJVEmailContent.Body = lEmailbody
			}
		}
	}
	log.Println("ConstructMailContent (-)")
	return lJVEmailContent, nil
}

/*
Pupose:This method is used to get the email Input for success Sgb Place Order  .
Parameters:

	pSgbClientDetails {},pStatus string

Response:

	==========
	*On Sucess
	==========
	{
     Name: clientName
     Status: S
     OrderDate: 8SEP2023
     OrderNumber: 121345687
     Symbol: sgb test
     Unit: 5
     Price: 500
     Amount:2500
     Activity : M
	},nil

	==========
	*On Error
	==========
	"",error

Author:PRASHANTH
Date: 8SEP2023
*/
func ConstructSuccessmail(pSgbClientDetails exchangecall.JvReqStruct) (emailUtil.EmailInput, error) {
	log.Println("ConstructSuccessmail (+)")
	type dynamicEmailStruct struct {
		Name        string
		Status      string
		OrderDate   string
		Symbol      string
		OrderNumber string
		Unit        string
		Price       string
		Amount      string
		// Activity    string
	}

	var lEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("toml/emailconfig.toml")

	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	// commented by pavithra
	// lEmailContent.Subject = "SGB Order"
	// html := "html/SgbOrderTemplate.html"

	lNovoConfig := common.ReadTomlConfig("toml/novoConfig.toml")
	lEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_ClientEmail_Subject"])
	html := fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_OrderSuccess_html"])

	lTemp, lErr := template.ParseFiles(html)
	if lErr != nil {
		log.Println("IMP03", lErr)
		return lEmailContent, lErr
	} else {
		var lTpl bytes.Buffer
		var lDynamicEmailVal dynamicEmailStruct
		lDynamicEmailVal.Name = pSgbClientDetails.ClientName
		lDynamicEmailVal.Amount = pSgbClientDetails.Amount
		lDynamicEmailVal.Unit = strconv.Itoa(pSgbClientDetails.Unit)
		lDynamicEmailVal.Price = strconv.Itoa(pSgbClientDetails.Price)
		lDynamicEmailVal.OrderDate = pSgbClientDetails.OrderDate
		lDynamicEmailVal.OrderNumber = pSgbClientDetails.RespOrderNo
		lDynamicEmailVal.Symbol = pSgbClientDetails.Symbol
		// lDynamicEmailVal.Activity = pSgbClientDetails.ActionCode
		//   ================================================
		if pSgbClientDetails.StatusCode == "0" {
			lDynamicEmailVal.Status = common.SUCCESS
		} else {
			lDynamicEmailVal.Status = common.FAILED
		}
		//   ================================================||

		// lEmailContent.ToEmailId = "prashanth.s@fcsonline.co.in"
		lEmailContent.ToEmailId = pSgbClientDetails.Mail
		// log.Println("lDynamicEmailVal Succes Mail", lDynamicEmailVal)
		lTemp.Execute(&lTpl, lDynamicEmailVal)
		lEmailbody := lTpl.String()

		lEmailContent.Body = lEmailbody
	}
	log.Println("ConstructSuccessmail (-)")
	return lEmailContent, nil
}

func GetEmailContent(pRequestCode string, pRequestString string) (string, error) {
	log.Println("GetEmailContent (+)")

	//this variable is used to get the lookuptable value from the table
	var lRespString string
	//Establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SGEC01", lErr1)
		return lRespString, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select d.description 
						from xx_lookup_header h,xx_lookup_details d
						where h.id = d.headerid 
						and h.Code = ?
						and d.Code = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pRequestCode, pRequestString)
		if lErr2 != nil {
			log.Println("SGEC02", lErr2)
			return lRespString, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRespString)
				if lErr3 != nil {
					log.Println("SGEC03", lErr3)
					return lRespString, lErr3
				} else {
					log.Println("lRespString", lRespString)
				}
			}
		}
	}
	log.Println("GetEmailContent (-)")
	return lRespString, nil
}

//------------------------------------------------------------------------------
// below commented code is only process application to exchange and update to db
//------------------------------------------------------------------------------

// func SgbOfflineSch(lUser string, r *http.Request) SchRespStruct {
// 	log.Println("SgbOfflineSch (+)")
// 	var lRespRec SchRespStruct
// 	var lFailRecArr []bsesgb.SgbRespStruct
// 	lRespRec.Status = common.SuccessCode
// 	lExchange := "BSE"
// 	BrokerId := 4
// 	MasterId := 1
// 	lExchangeReqArr, lJvDataArr, lErr1 := fetchPendingSgbOrder(lExchange, BrokerId, MasterId)
// 	if lErr1 != nil {
// 		log.Println("SS01", lErr1)
// 		lRespRec.ErrMsg = "SS01" + lErr1.Error()
// 		lRespRec.Status = common.ErrorCode
// 	} else {
// 		lToken, lErr2 := exchangecall.BseGetToken(lUser)
// 		if lErr2 != nil {
// 			log.Println("SS02", lErr2)
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "SS02" + lErr2.Error()
// 		} else {
// 			if lExchangeReqArr != nil {
// 				for lIdx := 0; lIdx < len(lExchangeReqArr); lIdx++ {
// 					lRespSgbRec, lErr3 := bsesgb.BseSgbOrder(lToken, lUser, lExchangeReqArr[lIdx])
// 					if lErr3 != nil {
// 						log.Println("SS03", lErr3)
// 						lRespRec.Status = common.ErrorCode
// 						lRespRec.ErrMsg = "SS03" + lErr3.Error()
// 					} else {
// 						log.Println("lrespsgbrec", lRespSgbRec)
// 						for lIdx := 0; lIdx < len(lRespSgbRec.Bids); lIdx++ {
// 							//check the status code and error code
// 							if lRespSgbRec.Bids[lIdx].ErrorCode == "0" && lRespSgbRec.StatusCode == "0" {
// 								log.Println("success--------------")
// 								lErr4 := updatePendingSgbStatus(lRespSgbRec, lUser)
// 								if lErr4 != nil {
// 									log.Println("SS04", lErr4)
// 									lRespRec.Status = common.ErrorCode
// 									lRespRec.ErrMsg = "SS04" + lErr4.Error()
// 								} else {
// 									//var lEmailContent emailutil.EmailInput
// 									lString := "SgbStatusMail"
// 									lStatus := common.SUCCESS
// 									lEmailContent, lErr4 := constructMail(lJvDataArr, lRespSgbRec.Bids[lIdx].OrderNo, lStatus)
// 									if lErr4 != nil {
// 										log.Println("SS05", lErr4)
// 										lRespRec.Status = common.ErrorCode
// 										lRespRec.ErrMsg = "SS05" + lErr4.Error()
// 									}
// 									lErr5 := emailUtil.SendEmail(lEmailContent, lString)
// 									if lErr5 != nil {
// 										log.Println("SS06", lErr5)
// 										lRespRec.Status = common.ErrorCode
// 										lRespRec.ErrMsg = "SS06" + lErr5.Error()
// 									}
// 								}
// 							} else {
// 								lFailRecArr = append(lFailRecArr, lRespSgbRec)
// 								log.Println(lFailRecArr, "lFailRecArr")
// 							}
// 						}
// 					}
// 				}
// 				lErr3 := updateFailedRec(lFailRecArr, lJvDataArr, lUser, r)
// 				if lErr3 != nil {
// 					log.Println("SS07", lErr3)
// 					lRespRec.Status = common.ErrorCode
// 					lRespRec.ErrMsg = "SS07" + lErr3.Error()
// 				} else {
// 					//---------------->>>>>>>>>>>>>
// 					//accounts email
// 				}
// 			} else {
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = "No Records found to Process"
// 			}
// 		}
// 	}
// 	log.Println("SgbOfflineSch (-)")
// 	return lRespRec
// }

func updateSgbHeader(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pStatus string) error {
	log.Println("updateSgbHeader(+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSUSHO1", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for i := 0; i < len(pRespSgbdata.Bids); i++ {
			lCoreString := `update a_sgb_orderheader  h
							set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,
							h.DpId=?,h.ClientBenfId=?,h.GuardianName=?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.StatusCode = ?,
							h.StatusMessage = ?,h.UpdatedDate = now(),h.UpdatedBy = ?,h.status=?
							where h.Id in (select d.HeaderId 
							from a_sgb_orderdetails d,a_sgb_orderheader h
							where d.HeaderId = h.Id and d.ReqOrderNo = ?)`

			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.ScripId, pRespSgbdata.PanNo, pRespSgbdata.InvestorCategory, pRespSgbdata.ApplicantName,
				pRespSgbdata.Depository, pRespSgbdata.DpId, pRespSgbdata.ClientBenfId, pRespSgbdata.GuardianName, pRespSgbdata.GuardianPanno,
				pRespSgbdata.GuardianRelation, pRespSgbdata.StatusCode, pRespSgbdata.StatusMessage, pUser, pStatus, pRespSgbdata.Bids[i].OrderNo)
			if lErr1 != nil {
				log.Println("SSUSHO2", lErr1)
				return lErr1
			}
		}
	}
	log.Println("updateSgbHeader(-)")
	return nil
}

// func updateSgbDetails(pRespSgbdata bsesgb.SgbRespStruct, pUser string) error {
// 	log.Println("updateSgbDetails(+)")
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUSDO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		for i := 0; i < len(pRespSgbdata.Bids); i++ {
// 			lCoreString := `update a_sgb_orderdetails d
// 							set d.RespSubscriptionunit = ?,d.RespRate = ?,d.ActionCode = ?,d.ErrorCode = ?,
// 							d.Message = ?,UpdatedDate = now(),UpdatedBy = ?
// 							where d.OrderNo = ?`

// 			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Bids[i].SubscriptionUnit, pRespSgbdata.Bids[i].Rate,
// 				pRespSgbdata.Bids[i].ActionCode, pRespSgbdata.Bids[i].ErrorCode, pRespSgbdata.Bids[i].Message, pUser, pRespSgbdata.Bids[i].OrderNo)
// 			if lErr1 != nil {
// 				log.Println("SSUSDO2", lErr1)
// 				return lErr1
// 			}
// 		}
// 	}
// 	log.Println("updateSgbDetails(-)")
// 	return nil
// }
// func updateBidTracking(pRespSgbdata bsesgb.SgbRespStruct, pUser string) error {
// 	log.Println("updateBidTracking(+)")
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUBO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		for i := 0; i < len(pRespSgbdata.Bids); i++ {
// 			lCoreString := `update a_sgbtracking_table  b
// 							set b.ErrorCode = ?,b.Message = ?,b.ApplicationStatus =?,UpdatedDate = now(),UpdatedBy = ?
// 							where b.OrderNo = ? and b.ActivityType = ? and b.Unit = ?`
// 			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Bids[i].ErrorCode,
// 				pRespSgbdata.Bids[i].Message, common.SUCCESS, pUser, pRespSgbdata.Bids[i].OrderNo, pRespSgbdata.Bids[i].ActionCode, pRespSgbdata.Bids[i].SubscriptionUnit)
// 			if lErr1 != nil {
// 				log.Println("SSUBO2", lErr1)
// 				return lErr1
// 			}
// 		}
// 	}
// 	log.Println("updateBidTracking(-)")
// 	return nil
// }

func updateDetailsWithJv(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pJvRespRec sgbplaceorder.JvStatusStruct) error {
	log.Println("updateDetailsWithJv (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSUSDO1", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for i := 0; i < len(pRespSgbdata.Bids); i++ {
			lCoreString := `update a_sgb_orderdetails d
			set d.RespSubscriptionunit = ?,d.RespRate = ?,d.ActionCode = ?,d.ErrorCode = ?,
			d.Message = ?,d.JvStatus=?,d.JvAmount=?,d.JvStatement=?,d.JvType=?,UpdatedDate = now(),UpdatedBy = ?
			where d.ReqOrderNo = ?`

			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Bids[i].SubscriptionUnit, pRespSgbdata.Bids[i].Rate,
				pRespSgbdata.Bids[i].ActionCode, pRespSgbdata.Bids[i].ErrorCode, pRespSgbdata.Bids[i].Message, pJvRespRec.JvStatus, pJvRespRec.JvAmount, pJvRespRec.JvStatement, pJvRespRec.JvType, pUser, pRespSgbdata.Bids[i].OrderNo)
			if lErr1 != nil {
				log.Println("SSUSDO2", lErr1)
				return lErr1
			}
		}
	}
	log.Println("updateDetailsWithJv (-)")
	return nil
}

func updateBidWithJv(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pJvRespRec sgbplaceorder.JvStatusStruct, pBrokerId int) error {
	log.Println("updateBidWithJv (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSUBO1", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for i := 0; i < len(pRespSgbdata.Bids); i++ {
			lCoreString := `update a_sgbtracking_table  b
			set b.ErrorCode = ?,b.JvAmount=?,b.JvStatus=?,JvStatement=?,JvType=?,b.Message = ?,UpdatedDate = now(),UpdatedBy = ?
			where b.ReqOrderNo = ? and b.ActivityType = ? and b.Unit = ? and b.brokerId = ?`

			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Bids[i].ErrorCode, pJvRespRec.JvAmount,
				pJvRespRec.JvStatus, pJvRespRec.JvStatement, pJvRespRec.JvType,
				pRespSgbdata.Bids[i].Message, pUser, pRespSgbdata.Bids[i].OrderNo, pRespSgbdata.Bids[i].ActionCode, pRespSgbdata.Bids[i].SubscriptionUnit, pBrokerId)
			if lErr1 != nil {
				log.Println("SSUBO2", lErr1)
				return lErr1
			}
		}
	}
	log.Println("updateBidWithJv (-)")
	return nil
}

func updatePendingErrorStatus(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pJvRespRec sgbplaceorder.JvStatusStruct, pBrokerId int) error {
	log.Println("updatePendingErrorStatus (+)")
	lStatus := common.FAILED
	lErr1 := updateSgbHeader(pRespSgbdata, pUser, lStatus)
	if lErr1 != nil {
		log.Println("SS01", lErr1)
		return lErr1
	} else {
		lErr1 = updateDetailsWithJv(pRespSgbdata, pUser, pJvRespRec)
		if lErr1 != nil {
			log.Println("SS02", lErr1)
			return lErr1
		} else {
			lErr1 = updateBidWithJv(pRespSgbdata, pUser, pJvRespRec, pBrokerId)
			if lErr1 != nil {
				log.Println("SS03", lErr1)
				return lErr1
			}
		}
	}
	log.Println("updatePendingErrorStatus (-)")
	return nil
}
