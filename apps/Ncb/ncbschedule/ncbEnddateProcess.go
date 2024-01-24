package ncbschedule

import (
	"bytes"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type NcbJvStatusStruct struct {
	BoJvStatus    string `json:"boJvstatus"`
	BoJvStatement string `json:"boJvstatement"`
	BoJvAmount    string `json:"boJvamount"`
	BoJvType      string `json:"boJvtype"`
	FoJvStatus    string `json:"foJvstatus"`
	FoJvStatement string `json:"foJvstatement"`
	FoJvAmount    string `json:"foJvamount"`
	FoJvType      string `json:"foJvtype"`
}

type NcbIssueList struct {
	Symbol                  string `json:"symbol"`
	Series                  string `json:"string"`
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

//----------------------------------------------------------------
// this method processing the NCB orders on NCB enddate
//----------------------------------------------------------------

func PlacingNcbbOrder(r *http.Request) ([]NcbIssueList, error) {
	log.Println("PlacingNcbbOrder (+)")

	var lncbIssue NcbIssueList
	var lNcbIssuesList []NcbIssueList
	lValidNcbArr, lErr1 := GoodTimeForApplyNcb()

	if lErr1 != nil {
		log.Println("PNO01", lErr1)
		return lNcbIssuesList, lErr1
	} else {
		log.Println("NCB Valid Bond array", lValidNcbArr)

		if lValidNcbArr != nil {
			for lNcbIdx := 0; lNcbIdx < len(lValidNcbArr); lNcbIdx++ {
				if len(lValidNcbArr) != 0 {

					lBrokerIdArr, lErr2 := GetNcbBrokers(lValidNcbArr[lNcbIdx].Exchange)
					if lErr2 != nil {
						log.Println("PNO02", lErr2)
						return lNcbIssuesList, lErr2
					} else {
						log.Println("NCB Broker array", lBrokerIdArr)

						if len(lBrokerIdArr) != 0 {

							for lBrokerIdx := 0; lBrokerIdx < len(lBrokerIdArr); lBrokerIdx++ {

								lCountStruct, lErr3 := ProcessNcbOrder(lValidNcbArr[lNcbIdx], lBrokerIdArr[lBrokerIdx], r)

								if lErr3 != nil {
									log.Println("PNO03", lErr3)
									lCountStruct.Symbol = lErr3.Error()
									lNcbIssuesList = append(lNcbIssuesList, lCountStruct)
								} else {
									lCountStruct.Symbol = lValidNcbArr[lNcbIdx].Symbol
									lNcbIssuesList = append(lNcbIssuesList, lCountStruct)
								}

							}

						} else {
							lncbIssue.Symbol = "No Brokers to Process"
							lNcbIssuesList = append(lNcbIssuesList, lncbIssue)
						}

					}

				} else {

				}

			}

		} else {
			log.Println("No NCB issuses ending today")
			lncbIssue.Symbol = "No NCB issuses ending today"
			lNcbIssuesList = append(lNcbIssuesList, lncbIssue)
		}

	}

	log.Println("PlacingNcbOrder (-)")
	return lNcbIssuesList, nil
}

//----------------------------------------------------------------
// this method is used to get the valid NCB record from database
//----------------------------------------------------------------
func GoodTimeForApplyNcb() ([]NcbStruct, error) {
	log.Println("GoodTimeForApplyNcb (+)")

	var lValidNcbArr []NcbStruct
	var lValidNcbRec NcbStruct

	config := common.ReadTomlConfig("toml/NcbConfig.toml")
	lApplyFlag := fmt.Sprintf("%v", config.(map[string]interface{})["NCB_ApplyDay_Flag"])

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GTAN01", lErr1)
		return lValidNcbArr, lErr1
	} else {
		defer lDb.Close()
		lCondition := ""
		lCoreString := `select n.id, n.Symbol, n.Exchange
		                from a_ncb_master n
						where n.DailyEndTime > now()`

		// where n.BiddingEndDate  = curdate()
		// and n.DailyEndTime > now()

		if lApplyFlag == "Y" {
			lCondition = `and n.BiddingEndDate >= curdate()`
		} else {
			lCondition = `and n.BiddingEndDate = curdate()`
		}

		lFinalQueryString := lCoreString + lCondition

		lRows, lErr2 := lDb.Query(lFinalQueryString)
		if lErr2 != nil {
			log.Println("GTAN02", lErr2)
			return lValidNcbArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lValidNcbRec.MasterId, &lValidNcbRec.Symbol, &lValidNcbRec.Exchange)
				if lErr3 != nil {
					log.Println("GTAN03", lErr3)
					return lValidNcbArr, lErr3
				} else {
					lValidNcbArr = append(lValidNcbArr, lValidNcbRec)
				}
			}
			log.Println("GoodTimeForApplyNcb -->", lValidNcbArr)
		}
	}
	log.Println("GoodTimeForApplyNcb (-)")
	return lValidNcbArr, nil
}

// ----------------------------------------------------------------
// this method is used to get the Brokerid from database
// ----------------------------------------------------------------
func GetNcbBrokers(pExchange string) ([]int, error) {
	log.Println("GetNcbBrokers(+)")

	var lBrokerArr []int
	var lBrokerRec int
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GNB01", lErr1)
		return lBrokerArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := ` select md.BrokerId  
		                 from a_ipo_memberdetails md
		                 where md.AllowedModules  like '%Gsec%'
		                 and md.Flag = 'Y'
		                 and md.OrderPreference = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pExchange)
		if lErr2 != nil {
			log.Println("GNB02", lErr2)
			return lBrokerArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lBrokerRec)
				if lErr3 != nil {
					log.Println("GNB03", lErr3)
					return lBrokerArr, lErr3
				} else {
					lBrokerArr = append(lBrokerArr, lBrokerRec)
				}
			}
		}
	}
	log.Println("GetNcbBrokers(-)")
	return lBrokerArr, nil
}

func ProcessNcbOrder(pValidNcb NcbStruct, pBrokerId int, pR *http.Request) (NcbIssueList, error) {
	log.Println("ProcessNcbOrder (+)")

	var lCountRec NcbIssueList

	lNcbReqArr, lErr1 := fetchNcbOrder(pValidNcb, pBrokerId)
	lCountRec.TotalReqCount = len(lNcbReqArr)
	if lErr1 != nil {
		log.Println("PNO01", lErr1)
		return lCountRec, lErr1
	} else {
		log.Println("lNcbReqArr", lNcbReqArr)
		log.Println("lCountRec.TotalReqCount ", lCountRec.TotalReqCount)

		for lNcblReqIdx := 0; lNcblReqIdx < len(lNcbReqArr); lNcblReqIdx++ {

			pResponseRec, lErr2 := NcbPostJvandProcessOrder(lNcbReqArr[lNcblReqIdx], pR, pValidNcb, pBrokerId)
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
				log.Println("PNO02", lErr1)
				lErr3 := NcbUpdateResponse(pResponseRec, pValidNcb.Exchange, pBrokerId)
				if lErr3 != nil {
					log.Println("PNO03", lErr3)
				} else {
					lErr3 := NcbMailSendingProcess(pResponseRec)
					if lErr3 != nil {
						log.Println("PSO04", lErr3)

					}
				}
			} else {
				lErr5 := NcbUpdateResponse(pResponseRec, pValidNcb.Exchange, pBrokerId)
				if lErr5 != nil {
					log.Println("PNO05", lErr5)
				} else {
					lErr6 := NcbMailSendingProcess(pResponseRec)
					if lErr6 != nil {
						log.Println("PNO06", lErr6)
					}
				}
			}

		}
	}
	log.Println("ProcessNcbOrder (-)")
	return lCountRec, nil
}

// ----------------------------------------------------------------
// this method is used to fetch the local placed orders from database
// ----------------------------------------------------------------
func fetchNcbOrder(pValidNcbRec NcbStruct, BrokerId int) ([]exchangecall.NcbJvReqStruct, error) {
	log.Println("fetchNcbOrder (+)")

	var lReqNcbArr []exchangecall.NcbJvReqStruct
	var lReqNcbData exchangecall.NcbJvReqStruct

	config := common.ReadTomlConfig("toml/NcbConfig.toml")
	lConfigTime := fmt.Sprintf("%v", config.(map[string]interface{})["NCB_SchConfig_time"])
	lSpecificClient := fmt.Sprintf("%v", config.(map[string]interface{})["NCB_Specific_Client"])
	lApplyFlag := fmt.Sprintf("%v", config.(map[string]interface{})["NCB_ApplyDay_Flag"])
	lProcessFlag := fmt.Sprintf("%v", config.(map[string]interface{})["Process_NCB_Flag"])
	lScheduleFlag := fmt.Sprintf("%v", config.(map[string]interface{})["Schedule_NCB_Flag"])

	// Parse the string into a time.Time object
	lTime, lErr1 := time.Parse("15:04:05", lConfigTime)
	if lErr1 != nil {
		log.Println("FNO01", lErr1)
		return lReqNcbArr, lErr1
	} else {
		//Get Current day

		currentDate := time.Now().Local()

		// Set the date component of lTime to today's date
		lConfigTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), lTime.Hour(), lTime.Minute(), lTime.Second(), 0, time.Local)
		lConfigUnixTime := lConfigTime.Unix()
		//current date unix time value
		lCurrentTime := time.Now()
		lCurrentTimeUnix := lCurrentTime.Unix()

		if lCurrentTimeUnix > lConfigUnixTime {
			log.Println("Time Validation", lCurrentTimeUnix, lConfigUnixTime)

			lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
			if lErr1 != nil {
				log.Println("FNO01", lErr1)
				return lReqNcbArr, lErr1
			} else {
				defer lDb.Close()

				lSqlString := `select  h.Symbol,h.Series,h.pan,
				                       h.depository,h.dpId, h.clientBenId, date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date,	
				                       h.ClientEmail, d.ReqUnit , d.Reqprice ,
									   d.ReqOrderNo, h.ReqapplicationNo,h.clientId,
			                          (CASE
				                        WHEN d.activityType = 'M' THEN 'M'
				                        WHEN d.activityType = 'N' THEN 'N'
				                        WHEN d.activityType = 'D' THEN 'D'
				                      ELSE d.activityType
				                      END ) AS ActionDescription
	                          from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
	                          where h.MasterId = n.id
	                          and d.HeaderId = h.Id
	                          and h.Status = 'success'
	                          and n.BiddingStartDate <= curdate()
	                          and time(now()) between n.DailyStartTime and n.DailyEndTime
	                          and h.cancelFlag != 'Y'
		                      and d.activityType != 'D'
	                          and h.clientId  = ? 
		                      and h.Exchange = ?
	 	                      and h.brokerId = ?
	                          and n.id = ?`

				lCondition := ""
				if lApplyFlag == "Y" {
					lCondition = ` and n.BiddingEndDate >= curdate()`
				} else {
					lCondition = ` and n.BiddingEndDate = curdate()`
				}

				lInclause := ""
				if lSpecificClient != "" {
					lInclause = `and h.clientId  in ('` + lSpecificClient + `')`
				}

				lProcessFlagString := ` and h.ProcessFlag = '` + lProcessFlag + `' and h.ScheduleStatus = '` + lScheduleFlag + `'`

				lFinalQueryString := lSqlString + lProcessFlagString + lCondition + lInclause
				log.Println("lFinalQueryString", lFinalQueryString)

				lRows, lErr2 := lDb.Query(lFinalQueryString, pValidNcbRec.Exchange, BrokerId, pValidNcbRec.MasterId)
				if lErr2 != nil {
					log.Println("FNO02", lErr2)
					return lReqNcbArr, lErr2
				} else {

					for lRows.Next() {

						lErr3 := lRows.Scan(&lReqNcbData.Symbol, &lReqNcbData.Series, &lReqNcbData.PanNo, &lReqNcbData.Depository, &lReqNcbData.DpId, &lReqNcbData.ClientBenId, &lReqNcbData.OrderDate, &lReqNcbData.Mail, &lReqNcbData.Unit, &lReqNcbData.Price, &lReqNcbData.ReqOrderNo, &lReqNcbData.ApplicationNumber, &lReqNcbData.ClientId, &lReqNcbData.ActivityType)

						// d.ReqAmount , &lReqNcbData.Amount,

						if lErr3 != nil {
							log.Println("FNO03", lErr3)
							return lReqNcbArr, lErr3
						} else {

							if lProcessFlag == "N" && lScheduleFlag == "N" {
								lReqNcbData.ActivityType = "N"
							}

							lPrice := int(math.Round(lReqNcbData.Price))

							lReqNcbData.Amount = float64(lReqNcbData.Unit * lPrice)
							lReqNcbArr = append(lReqNcbArr, lReqNcbData)

						}
					}

				}
			}
		} else {
			return lReqNcbArr, common.CustomError("The time for Configuration not yet come to process applications.")
		}

	}

	log.Println("fetchNcbOrder (-)")
	return lReqNcbArr, nil
}

//---------------------------------------------------------------------------------
// this method update the whole response
//---------------------------------------------------------------------------------

func NcbUpdateResponse(pJvRec exchangecall.NcbJvReqStruct, pExchange string, pBrokerId int) error {

	log.Println("NcbUpdateResponse (+)")
	lErr1 := NcbUpdateHeader2(pJvRec, pExchange, pBrokerId)
	if lErr1 != nil {
		log.Println("NEPUR01", lErr1)
		return lErr1
	} else {
		log.Println("Response Record Updated Successfully")
	}

	log.Println("NcbUpdateResponse (-)")
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

Author: KAVYA DHARSHANI
Date: 08JAN2024
*/

func NcbUpdateHeader2(pRespRec exchangecall.NcbJvReqStruct, pExchange string, pBrokerId int) error {
	log.Println("NcbUpdateHeader2 (+)")
	var lStatus string

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NEPNUH201", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lHeaderId, lErr2 := GetOrderId(pRespRec, pExchange)
		if lErr2 != nil {
			log.Println("NEPNUH202", lErr2)
			return lErr2
		} else {

			if lHeaderId != 0 {
				if pRespRec.StatusCode != "0" {
					lStatus = common.FAILED
				} else {
					lStatus = common.SUCCESS
				}

				lSqlString := `update a_ncb_orderheader h
				               set h.RespapplicationNo = ?,h.pan = ?, h.depository = ?, h.dpId = ?, 
							       h.clientBenId , h.ClientRefNumber, h.StatusCode=? , 
					               h.StatusMessage = ?, h.ErrorCode = ?, 
							       h.ErrorMessage = ?, h.status = ?, h.ProcessFlag  =?,h.ScheduleStatus  =?,
					               h.UpdatedBy  =?, h.UpdatedDate = now()
				               where h.clientId  = ?
				               and h.brokerId  = ?
				               and h.id = ?`

				_, lErr3 := lDb.Exec(lSqlString, pRespRec.ApplicationNumber, pRespRec.PanNo, pRespRec.Depository, pRespRec.DpId, pRespRec.ClientBenId, pRespRec.ClientRefNumber, pRespRec.StatusCode, pRespRec.StatusMessage, pRespRec.ErrorCode, pRespRec.ErrorMessage, lStatus, pRespRec.ProcessStatus, pRespRec.ErrorStage, common.AUTOBOT, pRespRec.ClientId, pBrokerId, lHeaderId)

				if lErr3 != nil {
					log.Println("NEPNUH203", lErr3)
					return lErr3
				} else {
					// call InsertDetails method to inserting the order details in order details table
					lErr4 := NcbUpdateDetail2(pRespRec, lHeaderId)
					if lErr4 != nil {
						log.Println("NEPNUH204", lErr4)
						return lErr4
					}
				}

			} else {
				return common.CustomError("Record Not Found")
			}

		}

	}

	log.Println("NcbUpdateHeader2 (-)")
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
Date: 08JAN2024
*/
func NcbUpdateDetail2(pRespRec exchangecall.NcbJvReqStruct, pHeaderId int) error {
	log.Println("NcbUpdateDetail2 (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NEPNUD201", lErr1)
		return lErr1
	} else {

		defer lDb.Close()

		lSqlString := `update a_ncb_orderdetails d
		               set d.RespOrderNo = ?, d.RespUnit = ?,d.Respprice = ?, d.RespAmount = ?,
			               d.ErrorCode = ?, d.ErrorMessage  = ?,
			               d.BoJvStatus = ?, d.BoJvAmount = ?,
			               d.BoJvStatement , d.BoJvType = ?,
			               d.FoJvStatus = ?, d.FoJvAmount  =?,
			               d.FoJvStatement = ?, d.FoJvType  = ?,
			               d.AddedDate  = ?, d.ModifiedDate  = ?,
			               d.UpdatedBy = ?, d.UpdatedDate = now()
	                  where d.headerId =?
	                  and d.ReqOrderNo = ?`

		_, lErr2 := lDb.Exec(lSqlString, pRespRec.RespOrderNo, pRespRec.Unit, pRespRec.Price, pRespRec.Amount, pRespRec.ErrorCode, pRespRec.ErrorMessage, pRespRec.BoJvStatus, pRespRec.BoJvAmount, pRespRec.BoJvStatement, pRespRec.BoJvType, pRespRec.FoJvStatus, pRespRec.FoJvAmount, pRespRec.FoJvStatement, pRespRec.FoJvType, pRespRec.EntryTime, pRespRec.LastActionTime, common.AUTOBOT, pHeaderId, pRespRec.ReqOrderNo)

		if lErr2 != nil {
			log.Println("NEPNUD202", lErr2)
			return lErr2
		}

	}

	log.Println("NcbUpdateDetail2 (-)")
	return nil

}

func NcbMailSendingProcess(pNcbRec exchangecall.NcbJvReqStruct) error {
	log.Println("NcbMailSendingProcess (+)")

	if pNcbRec.ProcessStatus == common.ErrorCode {
		lEmailRec, lErr1 := NcbConstructMailContent(pNcbRec)
		if lErr1 != nil {
			log.Println("NMSP01", lErr1)
			return lErr1
		} else {
			lErr2 := emailUtil.SendEmail(lEmailRec, "NCB")
			if lErr2 != nil {
				log.Println("NMSP02", lErr2)
				return lErr2
			}
		}
	} else if pNcbRec.ProcessStatus == "Y" {
		lEmailRec, lErr3 := NcbConstructSuccessmail(pNcbRec)
		if lErr3 != nil {
			log.Println("NMSP03", lErr3)
			return lErr3
		} else {
			lErr4 := emailUtil.SendEmail(lEmailRec, "NCB")
			if lErr4 != nil {
				log.Println("NMSP04", lErr4)
				return lErr4
			}
		}
	} else {
		return common.CustomError("Unable to Find the Issue")
	}
	log.Println("NcbMailSendingProcess (-)")
	return nil
}

func NcbConstructMailContent(pNcbRec exchangecall.NcbJvReqStruct) (emailUtil.EmailInput, error) {
	log.Println("NcbConstructMailContent (+)")

	type dynamicEmailStruct struct {
		Date     string  `json:"date"`
		Name     string  `json:"name"`
		Title    string  `json:"title"`
		Content  string  `json:"content"`
		Footer   string  `json:"footer"`
		ClientId string  `json:"clientId"`
		OrderNo  int     `json:"orderNo"`
		JvUnit   string  `json:"jvUnit"`
		JvAmount float64 `json:"jvAmount"`
	}

	var lJVEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("toml/emailconfig.toml")
	lJVEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lJVEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lJVEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])

	lNovoConfig := common.ReadTomlConfig("toml/novoConfig.toml")
	lJVEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["NCB_ClientEmail_Subject"])
	htmlTemplate := fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["NCB_CommonTemplate"])

	lNCBconfig := common.ReadTomlConfig("toml/NcbConfig.toml")
	lEmailIdCode := fmt.Sprintf("%v", lNCBconfig.(map[string]interface{})["NCB_EmailId_LookupCode"])
	lEmailContentCode := fmt.Sprintf("%v", lNCBconfig.(map[string]interface{})["NCB_EmailContent_LookupCode"])

	var lDynamicEmailVal dynamicEmailStruct

	lEmailId, lErr1 := NcbGetEmailContent(lEmailIdCode, pNcbRec.EmailDist)
	if lErr1 != nil {
		log.Println("NNCMC01", lErr1)
		return lJVEmailContent, lErr1
	} else {

		lJVEmailContent.ToEmailId = lEmailId

		lEmailContent, lErr2 := NcbGetEmailContent(lEmailContentCode, pNcbRec.ErrorStage)
		if lErr2 != nil {
			log.Println("NNCMC02", lErr2)
			return lJVEmailContent, lErr2
		} else {
			// lEmailContent = strings.Replace(lEmailContent, "{{Param1}}", "Sovereign Gold Bond (SGB)", 1)
			// lEmailContent = strings.Replace(lEmailContent, "{{Param2}}", "Sovereign Gold Bond (SGB)", 1)

			lDynamicEmailVal.Content = lEmailContent

			lTemp, lErr := template.ParseFiles(htmlTemplate)
			if lErr != nil {
				log.Println("NNCMC02", lErr)
				return lJVEmailContent, lErr
			} else {

				var lTpl bytes.Buffer
				currentTime := time.Now()
				currentDate := currentTime.Format("02-01-2006")

				lDynamicEmailVal.Date = currentDate
				lDynamicEmailVal.ClientId = pNcbRec.ClientId
				lDynamicEmailVal.OrderNo = pNcbRec.ReqOrderNo
				lDynamicEmailVal.JvUnit = strconv.Itoa(pNcbRec.Unit)
				lDynamicEmailVal.JvAmount = pNcbRec.Amount
				lDynamicEmailVal.Name = "Team"

				switch pNcbRec.ErrorStage {
				case "V":
					//  IT
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER FO ERROR REPORT"
					lDynamicEmailVal.Footer = "We kindly request you to Identify the solution on Processing of verifying Client's Account Details."
				case "F":
					// IT, RMS
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER FO FAILED REPORT"
					lDynamicEmailVal.Footer = "We kindly request you to Identify the solution on Processing of Deducting Amount Client's Account Details."
				case "B":
					// IT, ACC
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER BO FAILED REPORT"
					lDynamicEmailVal.Footer = "We kindly request you to Identify the solution on Processing of Deducting Amount Client's Account Details."
				case "E":
					// IT
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."
				case "A":
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."
				case "R":
					//IT,RMS,ACC
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."
				case "H":
					// IT,RMS
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER DETAILS"
					lDynamicEmailVal.Footer = "We kindly request you to review the details mentioned above and proceed with the Payment transfers to the respective clients."

				case "I":
					// client
					// lDynamicEmailVal.Content = lEmailContent
					lDynamicEmailVal.Title = "NCB ORDER DETAILS"
					lDynamicEmailVal.Footer = "Please maintain an sufficient fund balance from your next NBC orders."
					lDynamicEmailVal.Name = pNcbRec.ClientName
					lJVEmailContent.ToEmailId = pNcbRec.Mail
				}

				lTemp.Execute(&lTpl, lDynamicEmailVal)
				lEmailbody := lTpl.String()

				lJVEmailContent.Body = lEmailbody

			}

		}

	}

	log.Println("NcbConstructMailContent (-)")
	return lJVEmailContent, nil
}

/*
Pupose:This method is used to get the email Input for success Ncb Place Order  .
Parameters:

	pNcbClientDetails {},pStatus string

Response:

	==========
	*On Sucess
	==========
	{
     Name: clientName
     Status: S
     OrderDate: 8SEP2023
     OrderNumber: 121345687
     Symbol: Ncb test
     Unit: 5
     Price: 500
     Amount:2500
     Activity : M
	},nil

	==========
	*On Error
	==========
	"",error

Author:KAVYA DHARSHANI
Date: 8SEP2023
*/

func NcbConstructSuccessmail(pNcbClientDetails exchangecall.NcbJvReqStruct) (emailUtil.EmailInput, error) {
	log.Println("NcbConstructSuccessmail (+)")

	type dynamicEmailStruct struct {
		Name        string
		Status      string
		OrderDate   string
		Symbol      string
		Series      string
		OrderNumber int
		Unit        int
		Price       float64
		Amount      float64
		// Activity    string
	}

	var lEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("toml/emailconfig.toml")

	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])

	lNovoConfig := common.ReadTomlConfig("toml/novoConfig.toml")
	lEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["NCB_ClientEmail_Subject"])
	html := fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["NCB_OrderSuccess_html"])

	lTemp, lErr := template.ParseFiles(html)
	if lErr != nil {
		log.Println("NCSM01", lErr)
		return lEmailContent, lErr
	} else {

		var lTpl bytes.Buffer

		var lDynamicEmailVal dynamicEmailStruct
		lDynamicEmailVal.Name = pNcbClientDetails.ClientName
		lDynamicEmailVal.Amount = pNcbClientDetails.Amount
		lDynamicEmailVal.Unit = pNcbClientDetails.Unit
		lDynamicEmailVal.Price = pNcbClientDetails.Price
		lDynamicEmailVal.OrderDate = pNcbClientDetails.OrderDate
		lDynamicEmailVal.OrderNumber = pNcbClientDetails.RespOrderNo
		lDynamicEmailVal.Symbol = pNcbClientDetails.Symbol
		lDynamicEmailVal.Series = pNcbClientDetails.Series

		if pNcbClientDetails.StatusCode == "0" {
			lDynamicEmailVal.Status = common.SUCCESS
		} else {
			lDynamicEmailVal.Status = common.FAILED
		}

		// lEmailContent.ToEmailId = "kavyadharshani.m@fcsonline.co.in"
		lEmailContent.ToEmailId = pNcbClientDetails.Mail
		lTemp.Execute(&lTpl, lDynamicEmailVal)
		lEmailbody := lTpl.String()

		lEmailContent.Body = lEmailbody

	}

	log.Println("NcbConstructSuccessmail (-)")
	return lEmailContent, nil
}

func NcbGetEmailContent(pRequestCode string, pRequestString string) (string, error) {
	log.Println("NcbGetEmailContent (+)")

	//this variable is used to get the lookuptable value from the table
	var lRespString string
	//Establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NGEC01", lErr1)
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
			log.Println("NGEC02", lErr2)
			return lRespString, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRespString)
				if lErr3 != nil {
					log.Println("NGEC03", lErr3)
					return lRespString, lErr3
				} else {
					log.Println("lRespString", lRespString)
				}
			}
		}
	}
	log.Println("NcbGetEmailContent (-)")
	return lRespString, nil
}
