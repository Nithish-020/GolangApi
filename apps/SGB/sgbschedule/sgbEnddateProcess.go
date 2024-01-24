package sgbschedule

// import (
// 	"fcs23pkg/apps/SGB/sgbplaceorder"
// 	"fcs23pkg/common"
// 	"fcs23pkg/ftdb"
// 	"fcs23pkg/integration/bse/bsesgb"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"
// )

// type JvReqStruct struct {
// 	OrderNo       string `json:"orderno"`
// 	ClientId      string `json:"clientid"`
// 	BoJvStatus    string `json:"boJvstatus"`
// 	BoJvStatement string `json:"boJvstatement"`
// 	BoJvAmount    string `json:"boJvamount"`
// 	BoJvType      string `json:"boJvtype"`
// 	FoJvStatus    string `json:"foJvstatus"`
// 	FoJvStatement string `json:"foJvstatement"`
// 	FoJvAmount    string `json:"foJvamount"`
// 	FoJvType      string `json:"foJvtype"`
// 	Unit          string `json:"unit"`
// 	Price         string `json:"price"`
// 	ActionCode    string `json:"actioncode"`
// 	Symbol        string `json:"symbol"`
// 	OrderDate     string `json:"orderdate"`
// 	Amount        string `json:"amount"`
// 	Mail          string `json:"mail"`
// 	ClientName    string `json:"clientname"`
// 	EmailDist     string `json:"emailDist"`
// 	EmailMsg      string `json:"emailMsg"`
// 	ErrorStage    string `json:"errorStage"`
// 	Status        string `json:"status"`
// }
// type FailedJvStruct struct {
// 	SINo        int
// 	ClientId    string
// 	Amount      string
// 	Transaction string
// }

// type dynamicEmailStruct struct {
// 	Date        string
// 	FailedJvArr []FailedJvStruct
// }

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
// }

//----------------------------------------------------------------
// this method processing the SGB orders on SGB enddate
//----------------------------------------------------------------

// func PlacingSgbOrder(r *http.Request) ([]SgbIssueList, error) {
// 	log.Println("PlacingSgbOrder (+)")

// 	var lsgbIssue SgbIssueList
// 	var lSgbIssuesList []SgbIssueList
// 	lValidSgbArr, lErr1 := GoodTimeForApplySgb()
// 	if lErr1 != nil {
// 		log.Println("SPSO01", lErr1)
// 		return lSgbIssuesList, lErr1
// 	} else {
// 		log.Println("SGbValid Bond array", lValidSgbArr)

//  if no issue
// 		if lValidSgbArr != nil {
// for lSgbIdx := 0; lSgbIdx < len(lValidSgbArr); lSgbIdx++ {
// 			for lSgbIdx, lsgbValues := range lValidSgbArr {

// 				if len(lValidSgbArr) != 0 {

// 					lBrokerIdArr, lErr2 := GetSgbBrokers(lValidSgbArr[lSgbIdx].Exchange)
// 					if lErr2 != nil {
// 						log.Println("SPSO02", lErr2)
// 						return lSgbIssuesList, lErr2
// 					} else {
// 						log.Println("SGB Broker array", lBrokerIdArr)
// 						if len(lBrokerIdArr) != 0 {
// 							// var lWg sync.WaitGroup
// 							for lBrokerIdx := 0; lBrokerIdx < len(lBrokerIdArr); lBrokerIdx++ {
// lWg.Add(1)
// 								lCountStruct, lErr3 := ProcessSgbOrder(lValidSgbArr[lSgbIdx], lBrokerIdArr[lBrokerIdx], r)
// 								if lErr3 != nil {
// 									log.Println("SPSO03", lErr3)
// 									lsgbIssue.Symbol = "Getting Error in ProcessSgbOrder"
// 									lsgbIssue.TotalReqCount = lCountStruct.TotalReqCount
// 									lsgbIssue.JvSuccesssCount = lCountStruct.JvSuccesssCount
// 									lsgbIssue.JvFailedCount = lCountStruct.JvFailedCount
// 									lsgbIssue.ExchangeSuccessCount = lCountStruct.ExchangeSuccessCount
// 									lsgbIssue.ExchangeFailedCount = lCountStruct.ExchangeFailedCount
// 									lsgbIssue.ReverseJvSuccessCount = lCountStruct.ReverseJvSuccessCount
// 									lsgbIssue.ReverseJvFailedCount = lCountStruct.ReverseJvFailedCount
// 									lsgbIssue.FoJvSuccessCount = lCountStruct.FoJvSuccessCount
// 									lsgbIssue.FoJvFailedCount = lCountStruct.FoJvFailedCount
// 									lsgbIssue.FoVerifySuccessCount = lCountStruct.FoVerifySuccessCount
// 									lsgbIssue.FoVerifyFailedCount = lCountStruct.FoVerifyFailedCount
// 									lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
// 								} else {
// 									lsgbIssue.Symbol = lsgbValues.Symbol
// 									lsgbIssue.TotalReqCount = lCountStruct.TotalReqCount
// 									lsgbIssue.JvSuccesssCount = lCountStruct.JvSuccesssCount
// 									lsgbIssue.JvFailedCount = lCountStruct.JvFailedCount
// 									lsgbIssue.ExchangeSuccessCount = lCountStruct.ExchangeSuccessCount
// 									lsgbIssue.ExchangeFailedCount = lCountStruct.ExchangeFailedCount
// 									lsgbIssue.ReverseJvSuccessCount = lCountStruct.ReverseJvSuccessCount
// 									lsgbIssue.ReverseJvFailedCount = lCountStruct.ReverseJvFailedCount
// 									lsgbIssue.FoJvSuccessCount = lCountStruct.FoJvSuccessCount
// 									lsgbIssue.FoJvFailedCount = lCountStruct.FoJvFailedCount
// 									lsgbIssue.FoVerifySuccessCount = lCountStruct.FoVerifySuccessCount
// 									lsgbIssue.FoVerifyFailedCount = lCountStruct.FoVerifyFailedCount
// 									lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
// 								}
// 							}
// lWg.Wait()
// 						}
// 					}
// 				} else {
// 					lsgbIssue.Symbol = lsgbValues.Symbol
// 					lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
// 				}
// 			}
// 		} else {
// 			log.Println("No SGB issuses ending today")
// 			lsgbIssue.Symbol = "No SGB issuses ending today"
// 			lSgbIssuesList = append(lSgbIssuesList, lsgbIssue)
// 		}
// 	}
// 	log.Println("PlacingSgbOrder (-)")
// 	return lSgbIssuesList, nil
// }

// func ProcessSgbOrder(pValidSgb SgbStruct, pBrokerId int, pApiRequest *http.Request) (SgbIssueList, error) {
// 	log.Println("ProcessSgbOrder (+)")

// 	var lCountRec SgbIssueList

// 	lSgbReqArr, lJvDetailArr, lErr1 := fetchSgbOrder(pValidSgb.Exchange, pBrokerId, pValidSgb.MasterId)
// 	lCountRec.TotalReqCount = len(lSgbReqArr)
// 	if lErr1 != nil {
// 		log.Println("PSO01", lErr1)
// 		return lCountRec, lErr1
// 	} else {

// 		log.Println("lSgbReqArr", lSgbReqArr)
// 		log.Println("lJvDetailArr", lJvDetailArr)
// 		log.Println("lCountRec.TotalReqCount ", lCountRec.TotalReqCount)
// 		for lSqlReqIdx := 0; lSqlReqIdx < len(lSgbReqArr); lSqlReqIdx++ {
// 			for lBidIdx := 0; lBidIdx < len(lSgbReqArr[lSqlReqIdx].Bids); lBidIdx++ {

// 				for lJvReqIdx := 0; lJvReqIdx < len(lJvDetailArr); lJvReqIdx++ {

// log.Println("lsgbClients", lsgbclients)
// 					if lSgbReqArr[lSqlReqIdx].Bids[lBidIdx].OrderNo == lJvDetailArr[lJvReqIdx].OrderNo {

// 						BojvsuccessCount, BojvfailedCount, ExchangeSuccess, ExchangeFailed, ReverseJvSuccess, ReverseJvFailed, FojvsuccessCount, FojvfailedCount, lFundVerifySuccess, lFundVerifyFailed, lErr2 := PostJvForOrder(lSgbReqArr[lSqlReqIdx], lJvDetailArr[lJvReqIdx], pApiRequest, pValidSgb, pBrokerId)
// 						if lErr2 != nil {
// 							log.Println("PSO02", lErr2)
// 							lCountRec.JvSuccesssCount += BojvsuccessCount
// 							lCountRec.JvFailedCount += BojvfailedCount
// 							lCountRec.ExchangeSuccessCount += ExchangeSuccess
// 							lCountRec.ExchangeFailedCount += ExchangeFailed
// 							lCountRec.ReverseJvSuccessCount += ReverseJvSuccess
// 							lCountRec.ReverseJvFailedCount += ReverseJvFailed
// 							lCountRec.FoJvSuccessCount += FojvsuccessCount
// 							lCountRec.FoVerifySuccessCount += lFundVerifySuccess
// 							lCountRec.FoVerifyFailedCount += lFundVerifyFailed
// 							return lCountRec, lErr2
// 						} else {
// 							lCountRec.JvSuccesssCount += BojvsuccessCount
// 							lCountRec.JvFailedCount += BojvfailedCount
// 							lCountRec.ExchangeSuccessCount += ExchangeSuccess
// 							lCountRec.ExchangeFailedCount += ExchangeFailed
// 							lCountRec.ReverseJvSuccessCount += ReverseJvSuccess
// 							lCountRec.ReverseJvFailedCount += ReverseJvFailed
// 							lCountRec.FoJvSuccessCount += FojvsuccessCount
// 							lCountRec.FoJvFailedCount += FojvfailedCount
// 							lCountRec.FoVerifySuccessCount += lFundVerifySuccess
// 							lCountRec.FoVerifyFailedCount += lFundVerifyFailed
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("ProcessSgbOrder (-)")
// return ljvsuccessCount, ljvfailedCount, lExchangeSuccess, lExchangeFailed, lReverseJvSuccess, lReverseJvFailed, lTotalReqCount, lFoJvSuccess, lFoJvFailed, nil
// 	return lCountRec, nil
// }

// ----------------------------------------------------------------
// this method is used to get the valid sgb record from database
// ----------------------------------------------------------------
// func GoodTimeForApplySgb() ([]SgbStruct, error) {
// 	log.Println("GoodTimeForApplySgb (+)")

// 	var lValidSgbArr []SgbStruct
// 	var lValidSgbRec SgbStruct

// 	config := common.ReadTomlConfig("toml/SgbConfig.toml")
// 	lApplyFlag := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_ApplyDay_Flag"])
// 	log.Println("lApplyFlag", lApplyFlag)

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SGTAS01", lErr1)
// 		return lValidSgbArr, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCondition := ""
// 		lCoreString := `select m.id,m.Symbol ,m.Exchange
// 						from a_sgb_master m
// 						where m.Redemption = 'N'`

// 		if lApplyFlag == "Y" {
// 			lCondition = `and m.BiddingEndDate >= curdate()`
// 		} else {
// 			lCondition = `and m.BiddingEndDate = curdate()`
// 		}

// 		lFinalQueryString := lCoreString + lCondition
// 		// select m.id,m.Symbol ,m.Exchange
// 		// from a_sgb_master m
// 		// where m.BiddingEndDate = curdate()
// 		// and m.DailyEndTime > now()
// 		lRows, lErr2 := lDb.Query(lFinalQueryString)
// 		if lErr2 != nil {
// 			log.Println("SGTAS02", lErr2)
// 			return lValidSgbArr, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lValidSgbRec.MasterId, &lValidSgbRec.Symbol, &lValidSgbRec.Exchange)
// 				if lErr3 != nil {
// 					log.Println("SGTAS03", lErr3)
// 					return lValidSgbArr, lErr3
// 				} else {
// 					lValidSgbArr = append(lValidSgbArr, lValidSgbRec)
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GoodTimeForApplySgb (-)")
// 	return lValidSgbArr, nil
// }

// ----------------------------------------------------------------
// this method is used to fetch the local placed orders from database
// ----------------------------------------------------------------
// func fetchSgbOrder(pExchange string, BrokerId int, pMasterId int) ([]bsesgb.SgbReqStruct, []JvReqStruct, error) {
// 	log.Println("fetchSgbOrder (+)")

// 	var lReqSgbBid bsesgb.ReqSgbBidStruct
// 	var lReqSgbData bsesgb.SgbReqStruct
// 	var lReqJVData JvReqStruct
// 	var lReqSgbDataArr []bsesgb.SgbReqStruct
// 	var lReqJVDataArr []JvReqStruct

// 	config := common.ReadTomlConfig("toml/SgbConfig.toml")
// 	lConfigTime := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_SchConfig_time"])
// 	lSpecificClient := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_Specific_Client"])
// 	lApplyFlag := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_ApplyDay_Flag"])
// 	lProcessFlag := fmt.Sprintf("%v", config.(map[string]interface{})["Process_SGB_Flag"])
// 	lScheduleFlag := fmt.Sprintf("%v", config.(map[string]interface{})["Schedule_SGB_Flag"])

// 	// Parse the string into a time.Time object
// 	lTime, lErr1 := time.Parse("15:04:05", lConfigTime)
// 	if lErr1 != nil {
// 		log.Println("SSFP01", lErr1)
// 		return lReqSgbDataArr, lReqJVDataArr, lErr1
// 	} else {
// 		// Get current date
// 		currentDate := time.Now().Local()

// 		// Set the date component of lTime to today's date
// 		lTime = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), lTime.Hour(), lTime.Minute(), lTime.Second(), 0, time.Local)

// 		lConfigUnixTime := lTime.Unix()
// 		log.Println("Unix time:", lConfigUnixTime)

// 		lCurrentTime := time.Now()
// 		lCurrentTimeUnix := lCurrentTime.Unix()
// 		log.Println("lCurrentTimeUnix", lCurrentTimeUnix)

// 		if lCurrentTimeUnix > lConfigUnixTime {
// 			log.Println("Time Validation", lCurrentTimeUnix, lConfigUnixTime)
// 			// Perform actions if the current time is greater than the configured time for today

// 			lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 			if lErr1 != nil {
// 				log.Println("SSFP02", lErr1)
// 				return lReqSgbDataArr, lReqJVDataArr, lErr1
// 			} else {
// 				defer lDb.Close()

// 				//commented by pavithra
// 				// lSqlString := `	select h.ScripId ,h.PanNo,h.InvestorCategory,h.ApplicantName,h.Depository,h.DpId,h.ApplicantName ,date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date,h.ClientEmail ,m.Symbol ,h.ClientBenfId ,h.GuardianName ,h.GuardianPanNo ,h.GuardianRelation,d.BidId ,
// 				// d.ReqSubscriptionUnit ,d.ReqRate ,d.OrderNo ,d.ActionCode,h.ClientId,d.ReqSubscriptionUnit,d.ReqRate,d.OrderNo
// 				// (CASE
// 				// 	WHEN d.ActionCode = 'M' THEN 'Modify'
// 				// 	WHEN d.ActionCode = 'N' THEN 'New'
// 				// 	WHEN d.ActionCode = 'D' THEN 'Delete'
// 				// 	ELSE d.ActionCode
// 				// 	END ) AS ActionDescription
// 				// from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
// 				// where h.MasterId = m.id
// 				// and d.HeaderId = h.Id
// 				// and h.Status = "success"
// 				// and m.BiddingStartDate <= curdate()
// 				// and m.BiddingEndDate >= curdate()
// 				// and time(now()) between m.DailyStartTime and m.DailyEndTime
// 				// and h.cancelFlag != 'Y'
// 				// and m.Exchange = ?
// 				// and h.brokerId = ?
// 				// and m.id = ?`

// 				lSqlString := `select h.ScripId ,
// 								h.PanNo,
// 								h.InvestorCategory,
// 								h.ApplicantName,
// 								h.Depository,
// 								h.DpId,
// 								h.ApplicantName ,
// 								date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date,h.ClientEmail ,
// 								m.Symbol ,
// 								h.ClientBenfId ,
// 								h.GuardianName ,
// 								h.GuardianPanNo ,
// 								h.GuardianRelation,
// 								d.BidId ,
// 								d.ReqSubscriptionUnit ,
// 								d.ReqRate ,
// 								d.ReqOrderNo ,
// 								d.ActionCode,
// 								h.ClientId,
// 								d.ReqSubscriptionUnit,
// 								d.ReqRate,
// 								d.ReqOrderNo,
// 								(CASE
// 									WHEN d.ActionCode = 'M' THEN 'Modify'
// 									WHEN d.ActionCode = 'N' THEN 'New'
// 									WHEN d.ActionCode = 'D' THEN 'Delete'
// 									ELSE d.ActionCode
// 									END ) AS ActionDescription
// 								from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
// 								where h.MasterId = m.id
// 								and d.HeaderId = h.Id
// 								and h.Status = 'success'
// 								and m.BiddingStartDate <= curdate()
// 								and time(now()) between m.DailyStartTime and m.DailyEndTime
// 								and h.cancelFlag != 'Y'
// 								and m.Exchange = ?
// 								and h.brokerId = ?
// 								and m.id = ?`

// 				lCondition := ""
// 				if lApplyFlag == "Y" {
// 					lCondition = ` and m.BiddingEndDate >= curdate()`
// 				} else {
// 					lCondition = ` and m.BiddingEndDate = curdate()`
// 				}

// 				lInclause := ""
// 				if lSpecificClient != "" {
// 					lInclause = `and h.ClientId in ('` + lSpecificClient + `')`
// 				}
// 				lProcessFlagString := ` and h.ProcessFlag = '` + lProcessFlag + `' and h.ScheduleStatus = '` + lScheduleFlag + `'`

// 				lFinalQueryString := lSqlString + lProcessFlagString + lCondition + lInclause
// 				log.Println("lFinalQueryString", lFinalQueryString)

// 				lRows, lErr2 := lDb.Query(lFinalQueryString, pExchange, BrokerId, pMasterId)
// 				if lErr2 != nil {
// 					log.Println("SSFP02", lErr2)
// 					return lReqSgbDataArr, lReqJVDataArr, lErr2
// 				} else {
// 					for lRows.Next() {
// 						lErr3 := lRows.Scan(&lReqSgbData.ScripId, &lReqSgbData.PanNo, &lReqSgbData.InvestorCategory, &lReqSgbData.ApplicantName, &lReqSgbData.Depository, &lReqSgbData.DpId, &lReqJVData.ClientName, &lReqJVData.OrderDate, &lReqJVData.Mail, &lReqJVData.Symbol, &lReqSgbData.ClientBenfId, &lReqSgbData.GuardianName, &lReqSgbData.GuardianPanno, &lReqSgbData.GuardianRelation, &lReqSgbBid.BidId, &lReqSgbBid.SubscriptionUnit, &lReqSgbBid.Rate, &lReqSgbBid.OrderNo, &lReqSgbBid.ActionCode, &lReqJVData.ClientId, &lReqJVData.Unit, &lReqJVData.Price, &lReqJVData.OrderNo, &lReqJVData.ActionCode)
// 						if lErr3 != nil {
// 							log.Println("SSFP03", lErr3)
// 							return lReqSgbDataArr, lReqJVDataArr, lErr3
// 						} else {
// 							if lProcessFlag == "N" && lScheduleFlag == "N" {
// 								lReqSgbBid.ActionCode = "N"
// 							}
// 							lReqJVData.OrderNo = lReqSgbBid.OrderNo
// 							lReqSgbData.Bids = append(lReqSgbData.Bids, lReqSgbBid)
// 							lReqSgbDataArr = append(lReqSgbDataArr, lReqSgbData)
// 							lReqJVDataArr = append(lReqJVDataArr, lReqJVData)
// 							lReqSgbData.Bids = []bsesgb.ReqSgbBidStruct{}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("fetchSgbOrder (-)")
// 	return lReqSgbDataArr, lReqJVDataArr, nil
// }

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

// func updateSgbHeader(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pStatus string) error {
// 	log.Println("updateSgbHeader(+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUSHO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		for i := 0; i < len(pRespSgbdata.Bids); i++ {
// 			lCoreString := `update a_sgb_orderheader  h
// 							set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,
// 							h.DpId=?,h.ClientBenfId=?,h.GuardianName=?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.StatusCode = ?,
// 							h.StatusMessage = ?,h.UpdatedDate = now(),h.UpdatedBy = ?,h.status=?
// 							where h.Id in (select d.HeaderId
// 							from a_sgb_orderdetails d,a_sgb_orderheader h
// 							where d.HeaderId = h.Id and d.ReqOrderNo = ?)`

// 			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.ScripId, pRespSgbdata.PanNo, pRespSgbdata.InvestorCategory, pRespSgbdata.ApplicantName,
// 				pRespSgbdata.Depository, pRespSgbdata.DpId, pRespSgbdata.ClientBenfId, pRespSgbdata.GuardianName, pRespSgbdata.GuardianPanno,
// 				pRespSgbdata.GuardianRelation, pRespSgbdata.StatusCode, pRespSgbdata.StatusMessage, pUser, pStatus, pRespSgbdata.Bids[i].OrderNo)
// 			if lErr1 != nil {
// 				log.Println("SSUSHO2", lErr1)
// 				return lErr1
// 			}
// 		}
// 	}
// 	log.Println("updateSgbHeader(-)")
// 	return nil
// }

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
// func updateDetailsWithJv(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pJvRespRec sgbplaceorder.JvStatusStruct) error {
// 	log.Println("updateDetailsWithJv (+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUSDO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		for i := 0; i < len(pRespSgbdata.Bids); i++ {
// 			lCoreString := `update a_sgb_orderdetails d
// 			set d.RespSubscriptionunit = ?,d.RespRate = ?,d.ActionCode = ?,d.ErrorCode = ?,
// 			d.Message = ?,d.JvStatus=?,d.JvAmount=?,d.JvStatement=?,d.JvType=?,UpdatedDate = now(),UpdatedBy = ?
// 			where d.ReqOrderNo = ?`

// 			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Bids[i].SubscriptionUnit, pRespSgbdata.Bids[i].Rate,
// 				pRespSgbdata.Bids[i].ActionCode, pRespSgbdata.Bids[i].ErrorCode, pRespSgbdata.Bids[i].Message, pJvRespRec.JvStatus, pJvRespRec.JvAmount, pJvRespRec.JvStatement, pJvRespRec.JvType, pUser, pRespSgbdata.Bids[i].OrderNo)
// 			if lErr1 != nil {
// 				log.Println("SSUSDO2", lErr1)
// 				return lErr1
// 			}
// 		}
// 	}
// 	log.Println("updateDetailsWithJv (-)")
// 	return nil
// }
// func updateBidWithJv(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pJvRespRec sgbplaceorder.JvStatusStruct, pBrokerId int) error {
// 	log.Println("updateBidWithJv (+)")
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSUBO1", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		for i := 0; i < len(pRespSgbdata.Bids); i++ {
// 			lCoreString := `update a_sgbtracking_table  b
// 			set b.ErrorCode = ?,b.JvAmount=?,b.JvStatus=?,JvStatement=?,JvType=?,b.Message = ?,UpdatedDate = now(),UpdatedBy = ?
// 			where b.ReqOrderNo = ? and b.ActivityType = ? and b.Unit = ? and b.brokerId = ?`

// 			_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Bids[i].ErrorCode, pJvRespRec.JvAmount,
// 				pJvRespRec.JvStatus, pJvRespRec.JvStatement, pJvRespRec.JvType,
// 				pRespSgbdata.Bids[i].Message, pUser, pRespSgbdata.Bids[i].OrderNo, pRespSgbdata.Bids[i].ActionCode, pRespSgbdata.Bids[i].SubscriptionUnit, pBrokerId)
// 			if lErr1 != nil {
// 				log.Println("SSUBO2", lErr1)
// 				return lErr1
// 			}
// 		}
// 	}
// 	log.Println("updateBidWithJv (-)")
// 	return nil
// }
// func updatePendingErrorStatus(pRespSgbdata bsesgb.SgbRespStruct, pUser string, pJvRespRec sgbplaceorder.JvStatusStruct, pBrokerId int) error {
// 	log.Println("updatePendingErrorStatus (+)")
// 	lStatus := common.FAILED
// 	lErr1 := updateSgbHeader(pRespSgbdata, pUser, lStatus)
// 	if lErr1 != nil {
// 		log.Println("SS01", lErr1)
// 		return lErr1
// 	} else {
// 		lErr1 = updateDetailsWithJv(pRespSgbdata, pUser, pJvRespRec)
// 		if lErr1 != nil {
// 			log.Println("SS02", lErr1)
// 			return lErr1
// 		} else {
// 			lErr1 = updateBidWithJv(pRespSgbdata, pUser, pJvRespRec, pBrokerId)
// 			if lErr1 != nil {
// 				log.Println("SS03", lErr1)
// 				return lErr1
// 			}
// 		}
// 	}
// 	log.Println("updatePendingErrorStatus (-)")
// 	return nil
// }
