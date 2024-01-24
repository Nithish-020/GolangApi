package iposchedule

import (
	"fcs23pkg/apps/Ipo/placeorder"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nseipo"
	"log"
	"strconv"
	"strings"
	"time"
)

// this struct is response struct for OfflineSchedulePrgm
// type OfflineAppStruct struct {
// 	Status string `json:"status"`
// 	ErrMsg string `json:"errMsg"`
// }
type ScheduleStruct struct {
	SINo         string `json:"sino"`
	Method       string `json:"method"`
	TotalCount   int    `json:"totalCount"`
	SuccessCount int    `json:"successCount"`
	ErrCount     int    `json:"errCount"`
	Status       string `json:"status"`
	ErrMsg       string `json:"errMsg"`
}

type ScheduleRespStruct struct {
	ResponseArr []ScheduleStruct `json:"responseArr"`
	Status      string           `json:"status"`
	ErrMsg      string           `json:"errMsg"`
}

// this struct is used for grtting if values from order header table
type IdValueStruct struct {
	HeaderId int    `json:"headerId"`
	ClientId string `json:"clientId"`
}

// this struct is used to get bidRefno and their id from order details table
type BidRefStruct struct {
	BidId    int `json:"bidId"`
	BidRefNo int `json:"bidRefNo"`
}

// this struct is used to get bidtracking id and details from the bidtracking table
type BidTrackStruct struct {
	ClientId string         `json:"clientId"`
	AppNo    string         `json:"appNo"`
	BidIdArr []BidRefStruct `json:"bidIdArr"`
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------
type IpoOrderStruct struct {
	ApplicationNumber string `json:"applicationNumber"`
	ClientEmail       string `json:"clientEmail"`
	ClientId          string `json:"clientId"`
	ClientName        string `json:"clientName"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	OrderDate         string `json:"orderdate"`
	Amount            int    `json:"amount"`
	ActivityType      string `json:"activityType"`
}

type IpoStruct struct {
	MasterId int    `json:"masterId"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

// func IpoOfflineProcess(pUser string) {

// 	log.Println("offlineScheduler (+)")

// 	lValidIpoArr, lErr1 := GoodTimeForApplyIpo()
// 	if lErr1 != nil {
// 		log.Println("IIOP01", lErr1)
// 	} else {
// 		log.Println("IPO valid array", lValidIpoArr)
// 		if len(lValidIpoArr) != 0 {
// 			for lIpoIdx := 0; lIpoIdx < len(lValidIpoArr); lIpoIdx++ {

// 				lBrokerIdArr, lErr2 := GetIpoBrokers(lValidIpoArr[lIpoIdx].Exchange)
// 				if lErr2 != nil {
// 					log.Println("IIOP02", lErr2)
// 					// lRespRec.ErrMsg = "NSOS01" + lErr1.Error()
// 					// lRespRec.Status = common.ErrorCode
// 				} else {
// 					log.Println("Ipo Valid brokers", lBrokerIdArr)
// 					for lBrokerIdx := 0; lBrokerIdx < len(lBrokerIdArr); lBrokerIdx++ {
// 						// This variable is to store the exhange respose in the form of array
// 						var lRespArr []nseipo.ExchangeRespStruct

// 						//TODO lExchange, lErr3 := memberdetail.BseNsePercentCalc(lBrokerIdArr[lBrokerIdx], "IPO")

// 						lFinalReqArr, lIdValuesArr, lBidTrackIdArr, lAppArr, lErr3 := CommonConstruct(lValidIpoArr[lIpoIdx].Exchange, lBrokerIdArr[lBrokerIdx])
// 						if lErr3 != nil {
// 							log.Println("IIOP05", lErr3)
// 							// lRespRec.Status = common.ErrorCode
// 							// lRespRec.ErrMsg = "ISOAS05" + lErr5.Error()
// 						} else {

// 							if lFinalReqArr != nil {
// 								if lValidIpoArr[lIpoIdx].Exchange == common.NSE {
// 									var lErr4 error
// 									lRespArr, lErr4 = exchangecall.ApplyNseIpo(lFinalReqArr, pUser, lBrokerIdArr[lBrokerIdx])
// 									if lErr4 != nil {
// 										log.Println("IIOP05", lErr4)
// 										// lRespRec.Status = common.ErrorCode
// 										// lRespRec.ErrMsg = "ISOAS05" + lErr5.Error()

// 									} else {
// 										log.Println("lRespArr NSE", lRespArr)
// 									}
// 								} else if lValidIpoArr[lIpoIdx].Exchange == common.BSE {

// 									for lReqIdx := 0; lReqIdx < len(lFinalReqArr); lReqIdx++ {
// 										//
// 										var lNseArr []nseipo.ExchangeReqStruct
// 										lNseArr = append(lNseArr, lFinalReqArr[lReqIdx])
// 										var lErr5 error
// 										lRespArr, lErr5 = BseOrderPlace(lNseArr, pUser, lBrokerIdArr[lBrokerIdx])
// 										if lErr5 != nil {
// 											log.Println("IIOP05", lErr5)
// 											// lRespRec.Status = common.ErrorCode
// 											// lRespRec.ErrMsg = "ISBOS05" + lErr5.Error()
// 										} else {
// 											log.Println("lRespArr BSE", lRespArr)
// 										}
// 									}
// 								}
// 								ConstructCommonResp(pUser, lFinalReqArr, lRespArr, lIdValuesArr, lBidTrackIdArr, lAppArr)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		} else {
// 			log.Println("No IPO's are valid")
// 		}
// 	}

// }

func CommonConstruct(pExchange string, pBrokerId int) ([]nseipo.ExchangeReqStruct, []IdValueStruct, []BidTrackStruct, []placeorder.IpoOrderStruct, error) {
	log.Println("CommonConstruct (+)")
	//This variable is used to store the response from PendingList method
	// var lRespRec OfflineAppStruct
	var lBidTrackIdArr []BidTrackStruct
	var lApplicationArr []placeorder.IpoOrderStruct
	var lFinalReqArr []nseipo.ExchangeReqStruct

	// Call PendingList method to fetch the PendingApplications from database.
	lPendingBidsArr, lIdValuesArr, lErr1 := PendingListSch(pBrokerId, pExchange)
	if lErr1 != nil {
		log.Println("ICC01", lErr1)
		return lFinalReqArr, lIdValuesArr, lBidTrackIdArr, lApplicationArr, lErr1
	} else {

		//call GetBidTrackIdSch method to get the bidtracking details from the database
		lBidIdArr, lErr2 := GetBidTrackIdSch(pBrokerId, lPendingBidsArr)
		if lErr2 != nil {
			log.Println("ICC02", lErr2)
			return lFinalReqArr, lIdValuesArr, lBidTrackIdArr, lApplicationArr, lErr2
		} else {
			lBidTrackIdArr = lBidIdArr

			// log.Println("lPendingBidsArr", lPendingBidsArr)
			lReqArr, lErr3 := ConstructPendingSch(lPendingBidsArr)
			if lErr3 != nil {
				log.Println("ICC03", lErr3)
				return lFinalReqArr, lIdValuesArr, lBidTrackIdArr, lApplicationArr, lErr3
			} else {
				lFinalReqArr = lReqArr
				// ---------------->>>>> get the application no. ------------>>>>>>>>
				lAppArr, lErr4 := GetApiMaster(lFinalReqArr)
				if lErr4 != nil {
					log.Println("ICC04", lErr4)
					return lFinalReqArr, lIdValuesArr, lBidTrackIdArr, lApplicationArr, lErr4
				} else {
					lApplicationArr = lAppArr
					// log.Println(lFinalReqArr, "lFinalReqArr")
				}
			}
		}
	}
	log.Println("CommonConstruct (-)")
	return lFinalReqArr, lIdValuesArr, lBidTrackIdArr, lApplicationArr, nil

}

// func ConstructCommonResp(pUser string, pFinalReqArr []nseipo.ExchangeReqStruct, pRespArr []nseipo.ExchangeRespStruct, pIdValuesArr []IdValueStruct, pBidTrackIdArr []BidTrackStruct,
// 	pAppArr []placeorder.IpoOrderStruct) error {
// 	log.Println("ConstructCommonResp (+)")
// 	lFinalRespArr, lErr1 := StatusValidationSch(pRespArr)
// 	if lErr1 != nil {
// 		log.Println("ICCR01", lErr1)
// 		return lErr1
// 	} else {
// 		log.Println("lFinalRespArr", lFinalRespArr)
// 		// this looping is used to request array and response array to update the records
// 		for lpReqArrIdx := 0; lpReqArrIdx < len(pFinalReqArr); lpReqArrIdx++ {
// 			for lpRespArrIdx := 0; lpRespArrIdx < len(lFinalRespArr); lpRespArrIdx++ {
// 				// check whether the Req application no and response application no is equal or not
// 				if pFinalReqArr[lpReqArrIdx].ApplicationNo == lFinalRespArr[lpRespArrIdx].ApplicationNo {
// 					//Call the UpdateRecordSch to Update the Response array
// 					lErr2 := UpdateRecordSch(pFinalReqArr[lpReqArrIdx], lFinalRespArr[lpRespArrIdx], pIdValuesArr, pBidTrackIdArr, pUser)
// 					if lErr2 != nil {
// 						log.Println("ICCR02", lErr2)
// 						return lErr2
// 					} else {
// 					}
// 				}
// 			}
// 		}
// 		//------------------------>>>>>>>>>>>.emailsent <<<<<<<<<<<----------------
// 		for i := 0; i < len(lFinalRespArr); i++ {
// 			for j := 0; j < len(pAppArr); j++ {
// 				if lFinalRespArr[i].ApplicationNo == pAppArr[j].ApplicationNumber {
// 					if lFinalRespArr[i].Status == "success" {
// 						lString := "IpoOrderMail"
// 						lStatus := "success"
// 						lIpoEmailContent, lErr3 := placeorder.ConstructMail(pAppArr[j], lStatus)
// 						if lErr3 != nil {
// 							log.Println("ICCR03", lErr3)
// 							return lErr3
// 						} else {
// 							lErr4 := emailUtil.SendEmail(lIpoEmailContent, lString)
// 							if lErr4 != nil {
// 								log.Println("ICCR04", lErr4)
// 								return lErr4
// 							}
// 						}
// 					} else {
// 						lString := "IpoOrderMail"
// 						lStatus := "failed"
// 						lIpoEmailContent, lErr5 := placeorder.ConstructMail(pAppArr[j], lStatus)
// 						if lErr5 != nil {
// 							log.Println("ICCR05", lErr5)
// 							return lErr5
// 						} else {
// 							lErr6 := emailUtil.SendEmail(lIpoEmailContent, lString)
// 							// EmailInput.Subject = IpoOrderStruct.OrderDate
// 							if lErr6 != nil {
// 								log.Println("ICCR06", lErr6)
// 								return lErr6
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("ConstructCommonResp (-)")
// 	return nil
// }

//----------------------------------------------------------------
// this method is used to get the Brokerid from database
//----------------------------------------------------------------
func GetIpoBrokers(pExchange string) ([]int, error) {
	log.Println("GetIpoBrokers (+)")

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
						where md.AllowedModules  like '%Ipo%'
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
	log.Println("GetIpoBrokers (-)")
	return lBrokerArr, nil
}

//----------------------------------------------------------------
// this method is used to get the valid ipo record from database
//----------------------------------------------------------------
func GoodTimeForApplyIpo() ([]IpoStruct, error) {
	log.Println("GoodTimeForApplyIpo (+)")

	var lValidIpoArr []IpoStruct
	var lValidIpoRec IpoStruct
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SGTAS01", lErr1)
		return lValidIpoArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select (case when (a.value = 'Y' or a.value = 'I') then a.Id else 0 end ) value,a.Symbol,a.Exchange
						from (
						select m.Id,(case when (time(now()) between c.StartTime and c.EndTime) then 'Y' else 
						(
						(case when (c.StartTime is null and c.EndTime is null) then 
						(case when (time(now()) between m.DailyStartTime and m.DailyEndTime) then 'I' else 'N' end) else 'S' end)
						) end) value , m.Symbol ,m.Exchange 
						from a_ipo_master m,a_ipo_categories c,a_ipo_subcategory s
						where m.Id = c.MasterId
						and s.MasterId = c.MasterId 
						and c.Code = 'RETAIL'
						AND m.IssueType = "EQUITY"
						and s.CaCode = "RETAIL" 
						and s.SubCatCode = "IND"
						AND s.AllowUpi = 1) a`

		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("SGTAS02", lErr2)
			return lValidIpoArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lValidIpoRec.MasterId, &lValidIpoRec.Symbol, &lValidIpoRec.Exchange)
				if lErr3 != nil {
					log.Println("SGTAS03", lErr3)
					return lValidIpoArr, lErr3
				} else {
					lValidIpoArr = append(lValidIpoArr, lValidIpoRec)
				}
			}
		}
	}
	log.Println("GoodTimeForApplyIpo (-)")
	return lValidIpoArr, nil
}

/*
Purpose:This Method is used to process the pending application to the exchange
Parameters:

        User

Response:

	==========
	*On Sucess
	==========
	{
		"status": "success"
	}
	=========
	!On Error
	=========
	{
		"status": "E",
		"errMsg": "Error : Invalid Pending Bids"
	}

Author: KavyaDharshani
Date: 12SEP2023
*/

// func NseOfflineSch(lUser string) OfflineAppStruct {
// 	log.Println("OfflineAppSch (+)")

// //This variable is used to store the response from PendingList method
// var lRespRec OfflineAppStruct
// lRespRec.Status = common.SuccessCode
// lExchange := "NSE"
// // Call PendingList method to fetch the PendingApplications from database.
// lPendingBidsArr, lIdValuesArr, lErr1 := PendingListSch(lExchange)
// if lErr1 != nil {
// 	log.Println("ISOAS01", lErr1)
// 	lRespRec.Status = common.ErrorCode
// 	lRespRec.ErrMsg = "ISOAS01" + lErr1.Error()
// } else {

// 	//call GetBidTrackIdSch method to get the bidtracking details from the database
// 	lBidTrackIdArr, lErr2 := GetBidTrackIdSch(lPendingBidsArr)
// 	if lErr2 != nil {
// 		log.Println("ISOAS02", lErr2)
// 		lRespRec.Status = common.ErrorCode
// 		lRespRec.ErrMsg = "ISOAS02" + lErr2.Error()
// 	} else {
// 		// log.Println("lPendingBidsArr", lPendingBidsArr)
// 		//todo:-------------------------------
// 		lFinalReqArr, lErr3 := ConstructPendingSch(lPendingBidsArr)
// 		if lErr3 != nil {
// 			log.Println("ISOAS03", lErr3)
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "ISOAS03" + lErr3.Error()
// 		} else {
// 			// ---------------->>>>> get the application no. ------------>>>>>>>>
// 			lAppArr, lErr4 := GetApiMaster(lFinalReqArr)
// 			if lErr4 != nil {
// 				log.Println("ISOAS04", lErr4)
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = "ISOAS04" + lErr4.Error()
// 			} else {
// 				// log.Println(lAppArr)
// 				log.Println(lFinalReqArr, "lFinalReqArr")
// 				if lFinalReqArr != nil {
//-------------------------------------------------------------->>>>>>>>>>>
// Call ApplyIpo to pass the applications to Exchange.
// lRespArr, lErr5 := exchangecall.ApplyNseIpo(lFinalReqArr, lUser, 4)
// if lErr5 != nil {
// 	log.Println("ISOAS05", lErr5)
// 	lRespRec.Status = common.ErrorCode
// 	lRespRec.ErrMsg = "ISOAS05" + lErr5.Error()
// } else {
// 	log.Println(lRespArr)
//call StatusValidationSch method to split the records by status and its give the finalized Array Value
// lFinalRespArr, lErr6 := StatusValidationSch(lRespArr)
// if lErr6 != nil {
// 	log.Println("ISOAS06", lErr6)
// 	lRespRec.Status = common.ErrorCode
// 	lRespRec.ErrMsg = "ISOAS06" + lErr6.Error()
// } else {
// 	log.Println("lFinalRespArr", lFinalRespArr)
// 	// this looping is used to request array and response array to update the records
// 	for lpReqArrIdx := 0; lpReqArrIdx < len(lFinalReqArr); lpReqArrIdx++ {
// 		for lpRespArrIdx := 0; lpRespArrIdx < len(lFinalRespArr); lpRespArrIdx++ {
// 			// check whether the Req application no and response application no is equal or not
// 			if lFinalReqArr[lpReqArrIdx].ApplicationNo == lFinalRespArr[lpRespArrIdx].ApplicationNo {
// 				//Call the UpdateRecordSch to Update the Response array
// 				lErr7 := UpdateRecordSch(lFinalReqArr[lpReqArrIdx], lFinalRespArr[lpRespArrIdx], lIdValuesArr, lBidTrackIdArr, lUser)
// 				if lErr7 != nil {
// 					log.Println("ISOAS07", lErr7)
// 					lRespRec.Status = common.ErrorCode
// 					lRespRec.ErrMsg = "ISOAS07" + lErr7.Error()
// 				} else {
// 					lRespRec.Status = common.SuccessCode
// 					lRespRec.ErrMsg = "Processed Successfully...."
// 				}
// 			}
// 		}
// 	}
// 	//------------------------>>>>>>>>>>>.emailsent <<<<<<<<<<<----------------
// 	for i := 0; i < len(lFinalRespArr); i++ {
// 		for j := 0; j < len(lAppArr); j++ {
// 			if lFinalRespArr[i].ApplicationNo == lAppArr[j].ApplicationNumber {
// 				if lFinalRespArr[i].Status == "success" {
// 					lString := "IpoOrderMail"
// 					lStatus := "success"
// 					lIpoEmailContent, lErr8 := placeorder.ConstructMail(lAppArr[j], lStatus)
// 					if lErr8 != nil {
// 						log.Println("ISOAS08", lErr8)
// 						lRespRec.Status = common.ErrorCode
// 						lRespRec.ErrMsg = "ISOAS08" + lErr8.Error()
// 					} else {
// 						lErr9 := emailUtil.SendEmail(lIpoEmailContent, lString)
// 						if lErr9 != nil {
// 							log.Println("ISOAS09", lErr9)
// 							lRespRec.Status = common.ErrorCode
// 							lRespRec.ErrMsg = "ISOAS09" + lErr9.Error()
// 						}
// 					}
// 				} else {
// 					lString := "IpoOrderMail"
// 					lStatus := "failed"
// 					lIpoEmailContent, lErr10 := placeorder.ConstructMail(lAppArr[j], lStatus)
// 					if lErr10 != nil {
// 						log.Println("ISOAS10", lErr10)
// 						lRespRec.Status = common.ErrorCode
// 						lRespRec.ErrMsg = "ISOAS10" + lErr10.Error()
// 					} else {
// 						lErr11 := emailUtil.SendEmail(lIpoEmailContent, lString)
// 						// EmailInput.Subject = IpoOrderStruct.OrderDate
// 						if lErr11 != nil {
// 							log.Println("ISOAS11", lErr11)
// 							lRespRec.Status = common.ErrorCode
// 							lRespRec.ErrMsg = "ISOAS11" + lErr11.Error()
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// 		}
// 	} else {
// 		lRespRec.ErrMsg = "No Records found for process application"
// 		lRespRec.Status = common.ErrorCode
// 	}
// }
// }
// }
// 	}
// 	log.Println("OfflineAppSch (-)")
// 	return lRespRec
// }

/*
Purpose: this method is used to split the records for updating and finalize the array value to updating in database.
Parameters:

	pRespArr

Response:

	    ==========
	    *On Sucess
	    ==========
		[
		{
			"symbol": "UPIEMP",
			"applicationNumber": "FT000069172957",
			"category": "IND",
			"clientName": "LAKSHMANAN ASHOK KUMAR",
			"depository": "CDSL",
			"dpId": "",
			"clientBenId": "1208030000262661",
			"nonASBA": false,
			"pan": "AGMPA8575C",
			"referenceNumber": "bab61fa65db7c901",
			"allotmentMode": "demat",
			"upiFlag": "Y",
			"upi": "test@kmbl",
			"bankCode": "null",
			"locationCode": "null",
			"timestamp": "",
			"bids": [
				{
				    "activityType": "modify",
				    "bidReferenceNumber": "123456",
				    "quantity": 10,
				    "atCutOff": true,
				    "price": 10,
				    "amount": 100,
				    "remark": "FT000069172951",
				    "lotSize": 10
				},
				{
					"activityType": "cancel",
					"bidReferenceNumber": "2023051900001001",
					"quantity": 10,
					"atCutOff": false,
					"price": 9,
					"amount": 90,
					"remark": "FT000069172951",
					"lotSize": 10
				}

			]
		},
		]
	    =========
	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 07JUNE2023
*/
// func StatusValidationSch(pRespArr []nseipo.ExchangeRespStruct) ([]nseipo.ExchangeRespStruct, error) {
// 	log.Println("StatusValidationSch (+)")

// 	// this variable is used to get the finalArray value to store the responses in database
// 	var lFinalArr []nseipo.ExchangeRespStruct
// 	// this method is used to split the Records by status
// 	lSuccessArr, lFailedArr, lErr1 := SplitRecordSch(pRespArr)
// 	if lErr1 != nil {
// 		log.Println("ISSV01", lErr1)
// 		return lFinalArr, lErr1
// 	} else {
// 		// log.Println("lSuccessArr", lSuccessArr)
// 		// log.Println("lFailedArr", lFailedArr)
// 		if lSuccessArr != nil {
// 			for lSuccessIdx := 0; lSuccessIdx < len(lSuccessArr); lSuccessIdx++ {
// 				lFinalArr = append(lFinalArr, lSuccessArr[lSuccessIdx])
// 			}
// 			// log.Println("lFinalArr", lFinalArr)
// 		}
// 		if lFailedArr != nil {
// 			lFailUpdateArr, lErr2 := FailedUpdateSch(lFailedArr)
// 			if lErr2 != nil {
// 				log.Println("ISSV01", lErr2)
// 				return lFinalArr, lErr2
// 			} else {
// 				if lFailUpdateArr != nil {
// 					for lFailIdx := 0; lFailIdx < len(lFailUpdateArr); lFailIdx++ {
// 						lFinalArr = append(lFinalArr, lFailUpdateArr[lFailIdx])
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("StatusValidationSch (-)")
// 	return lFinalArr, nil
// }

func StatusValidationSch(pRespArr []nseipo.ExchangeRespStruct, pExchange string) ([]nseipo.ExchangeRespStruct, int, int, error) {
	log.Println("StatusValidationSch (+)")

	var lErrCount int
	var lSuccessCount int

	// this variable is used to get the finalArray value to store the responses in database
	var lFinalArr []nseipo.ExchangeRespStruct
	// this method is used to split the Records by status
	lSuccessArr, lFailedArr, lErr1 := SplitRecordSch(pRespArr)
	if lErr1 != nil {
		log.Println("ISSV01", lErr1)
		return lFinalArr, lSuccessCount, lErrCount, lErr1
	} else {
		// log.Println("lSuccessArr", lSuccessArr)
		// log.Println("lFailedArr", lFailedArr)
		if lSuccessArr != nil {
			for lSuccessIdx := 0; lSuccessIdx < len(lSuccessArr); lSuccessIdx++ {
				lFinalArr = append(lFinalArr, lSuccessArr[lSuccessIdx])
			}
			// log.Println("lFinalArr", lFinalArr)
		}
		if lFailedArr != nil {
			if pExchange == common.NSE {
				lFailUpdateArr, lFailedCount, lErr2 := FailedUpdateSch(lFailedArr)
				if lErr2 != nil {
					log.Println("ISSV02", lErr2)
					return lFinalArr, lSuccessCount, lErrCount, lErr2
				} else {
					if lFailUpdateArr != nil {
						for lFailIdx := 0; lFailIdx < len(lFailUpdateArr); lFailIdx++ {
							lFinalArr = append(lFinalArr, lFailUpdateArr[lFailIdx])
						}

					}
				}
				lErrCount = lFailedCount
			} else {
				lErrCount = len(lFailedArr)
			}
		}
	}
	lSuccessCount = len(lSuccessArr)
	// log.Println("failedArr ", lFailedArr)
	log.Println("StatusValidationSch (-)")
	return lFinalArr, lSuccessCount, lErrCount, nil
}

/*
Purpose: this method is filter the failed records for updating in database
Parameters:

	pFailedArr

Response:

	    ==========
	    *On Sucess
	    ==========
		[{
			"symbol": "UPIEMP",
			"applicationNumber": "FT000069172957",
			"category": "IND",
			"clientName": "LAKSHMANAN ASHOK KUMAR",
			"depository": "CDSL",
			"dpId": "",
			"clientBenId": "1208030000262661",
			"nonASBA": false,
			"pan": "AGMPA8575C",
			"referenceNumber": "bab61fa65db7c901",
			"allotmentMode": "demat",
			"upiFlag": "Y",
			"upi": "test@kmbl",
			"bankCode": "null",
			"locationCode": "null",
			"timestamp": "",
			"bids": [
				{
				    "activityType": "modify",
				    "bidReferenceNumber": "123456",
				    "quantity": 10,
				    "atCutOff": true,
				    "price": 10,
				    "amount": 100,
				    "remark": "FT000069172951",
				    "lotSize": 10
				},
				{
					"activityType": "cancel",
					"bidReferenceNumber": "2023051900001001",
					"quantity": 10,
					"atCutOff": false,
					"price": 9,
					"amount": 90,
					"remark": "FT000069172951",
					"lotSize": 10
				}

			]
		}]
	    =========
	    !On Error
	    =========
		[],error

Author: Pavithra
Date: 07JUNE2023
*/
// func FailedUpdateSch(pFailedArr []nseipo.ExchangeRespStruct) ([]nseipo.ExchangeRespStruct, error) {
// 	log.Println("FailedUpdateSch (+)")

// 	//this variable is used to get the lookuptable value from the table
// 	var lFlag string
// 	//this variable is used to get the updatable records
// 	var lUpdateRec []nseipo.ExchangeRespStruct
// 	//this variable is used to get records of non updatable
// 	var lNoUpdate []nseipo.ExchangeRespStruct
// 	//Establish a database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("ISFU01", lErr1)
// 		return lUpdateRec, lErr1
// 	} else {
// 		defer lDb.Close()
// 		// split the records to check the updatable records
// 		for lFailedIdx := 0; lFailedIdx < len(pFailedArr); lFailedIdx++ {
// 			log.Println("pFailedArr[lFailedIdx].ReasonCode", pFailedArr[lFailedIdx].ReasonCode)

// 			lCoreString := `select d.Attribute1
// 							from xx_lookup_header h,xx_lookup_details d
// 							where h.id = d.headerid
// 							and h.Code = 'IpoMasterCode'
// 							and d.Code = ?`

// 			lRows, lErr2 := lDb.Query(lCoreString, pFailedArr[lFailedIdx].ReasonCode)
// 			if lErr2 != nil {
// 				log.Println("ISFU02", lErr2)
// 				return lUpdateRec, lErr2
// 			} else {
// 				for lRows.Next() {
// 					lErr3 := lRows.Scan(&lFlag)
// 					if lErr3 != nil {
// 						log.Println("ISFU03", lErr3)
// 						return lUpdateRec, lErr3
// 					} else {
// 						log.Println(lFlag)
// 						// check the failed reasoncode is updatable or else
// 						if lFlag == "Y" {
// 							lUpdateRec = append(lUpdateRec, pFailedArr[lFailedIdx])
// 						} else if lFlag == "N" {
// 							lNoUpdate = append(lNoUpdate, pFailedArr[lFailedIdx])
// 						} else {
// 							log.Println("Failed To Get Flag from database")
// 						}
// 					}
// 				}
// 			}
// 		}
// 		// log.Println("lUpdateRec", lUpdateRec)
// 		// log.Println("lNoUpdate", lNoUpdate)
// 	}
// 	log.Println("FailedUpdateSch (-)")
// 	return lUpdateRec, nil
// }

func FailedUpdateSch(pFailedArr []nseipo.ExchangeRespStruct) ([]nseipo.ExchangeRespStruct, int, error) {
	log.Println("FailedUpdateSch (+)")

	//this variable is used to get the lookuptable value from the table
	var lFlag string
	//this variable is used to get the updatable records
	var lUpdateRec []nseipo.ExchangeRespStruct
	//this variable is used to get records of non updatable
	var lNoUpdate []nseipo.ExchangeRespStruct
	//this variable is used to get records of non updatable
	var lErrorCount int
	//Establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISFU01", lErr1)
		return lUpdateRec, lErrorCount, lErr1
	} else {
		defer lDb.Close()
		// split the records to check the updatable records
		for lFailedIdx := 0; lFailedIdx < len(pFailedArr); lFailedIdx++ {
			// log.Println("pFailedArr[lFailedIdx].ReasonCode", pFailedArr[lFailedIdx].ReasonCode)

			lCoreString := `select d.Attribute1
							from xx_lookup_header h,xx_lookup_details d
							where h.id = d.headerid 
							and h.Code = 'IpoMasterCode'
							and d.Code = ?`

			lRows, lErr2 := lDb.Query(lCoreString, pFailedArr[lFailedIdx].ReasonCode)
			if lErr2 != nil {
				log.Println("ISFU02", lErr2)
				return lUpdateRec, lErrorCount, lErr2
			} else {
				for lRows.Next() {
					lErr3 := lRows.Scan(&lFlag)
					if lErr3 != nil {
						log.Println("ISFU03", lErr3)
						return lUpdateRec, lErrorCount, lErr3
					} else {
						// log.Println(lFlag)
						// check the failed reasoncode is updatable or else
						if lFlag == "Y" {
							lUpdateRec = append(lUpdateRec, pFailedArr[lFailedIdx])
						} else if lFlag == "N" {
							lNoUpdate = append(lNoUpdate, pFailedArr[lFailedIdx])
						} else {
							log.Println("Failed To Get Flag from database")
						}
					}
				}
			}
		}
		// log.Println("lUpdateRec", lUpdateRec)
		lErrorCount = len(lNoUpdate) + len(lUpdateRec)
	}
	log.Println("FailedUpdateSch (-)")
	return lUpdateRec, lErrorCount, nil
}

/*
Purpose: this method is splitting the records based on the application status
Parameters:

	pRespArr

Response:

	    ==========
	    *On Sucess
	    ==========
		[{
			"symbol": "UPIEMP",
			"applicationNumber": "FT000069172957",
			"category": "IND",
			"clientName": "LAKSHMANAN ASHOK KUMAR",
			"depository": "CDSL",
			"dpId": "",
			"clientBenId": "1208030000262661",
			"nonASBA": false,
			"pan": "AGMPA8575C",
			"referenceNumber": "bab61fa65db7c901",
			"allotmentMode": "demat",
			"upiFlag": "Y",
			"upi": "test@kmbl",
			"bankCode": "null",
			"locationCode": "null",
			"timestamp": "",
			"bids": [
				{
				    "activityType": "modify",
				    "bidReferenceNumber": "123456",
				    "quantity": 10,
				    "atCutOff": true,
				    "price": 10,
				    "amount": 100,
				    "remark": "FT000069172951",
				    "lotSize": 10
				},
				{
					"activityType": "cancel",
					"bidReferenceNumber": "2023051900001001",
					"quantity": 10,
					"atCutOff": false,
					"price": 9,
					"amount": 90,
					"remark": "FT000069172951",
					"lotSize": 10
				}

			]
		}]
	    =========
	    !On Error
	    =========
		[],[],error

Author: Pavithra
Date: 07JUNE2023
*/
func SplitRecordSch(pRespArr []nseipo.ExchangeRespStruct) ([]nseipo.ExchangeRespStruct, []nseipo.ExchangeRespStruct, error) {
	log.Println("SplitRecordSch (+)")
	//This variable is used to store the pending Bids list in Array
	var lSuccessArr []nseipo.ExchangeRespStruct
	//
	var lFailedArr []nseipo.ExchangeRespStruct

	// for get the current time available ipos
	for lRespIdx := 0; lRespIdx < len(pRespArr); lRespIdx++ {
		// log.Println(pRespArr[lRespIdx].Status, "lRespIdx")
		//
		if pRespArr[lRespIdx].Status == "success" {
			lSuccessArr = append(lSuccessArr, pRespArr[lRespIdx])
		} else if pRespArr[lRespIdx].Status == "failed" {
			lFailedArr = append(lFailedArr, pRespArr[lRespIdx])
		}
	}
	log.Println("SplitRecordSch (-)")
	return lSuccessArr, lFailedArr, nil
}

/*
Purpose: this method is used to get the available ipo applications on currenttime
Parameters:

	pPendingArr

Response:

	    ==========
	    *On Sucess
	    ==========
	   	[{
			"symbol": "UPIEMP",
			"applicationNumber": "FT000069172957",
			"category": "IND",
			"clientName": "LAKSHMANAN ASHOK KUMAR",
			"depository": "CDSL",
			"dpId": "",
			"clientBenId": "1208030000262661",
			"nonASBA": false,
			"pan": "AGMPA8575C",
			"referenceNumber": "bab61fa65db7c901",
			"allotmentMode": "demat",
			"upiFlag": "Y",
			"upi": "test@kmbl",
			"bankCode": "null",
			"locationCode": "null",
			"timestamp": "",
			"bids": [
				{
				    "activityType": "modify",
				    "bidReferenceNumber": "123456",
				    "quantity": 10,
				    "atCutOff": true,
				    "price": 10,
				    "amount": 100,
				    "remark": "FT000069172951",
				    "lotSize": 10
				},
				{
					"activityType": "cancel",
					"bidReferenceNumber": "2023051900001001",
					"quantity": 10,
					"atCutOff": false,
					"price": 9,
					"amount": 90,
					"remark": "FT000069172951",
					"lotSize": 10
				}

			]
		}]
	    =========
	    !On Error
	    =========
		[],error

Author: Pavithra
Date: 07JUNE2023
*/
func ConstructPendingSch(pPendingArr []nseipo.ExchangeReqStruct) ([]nseipo.ExchangeReqStruct, error) {
	log.Println("ConstructPendingSch (+)")
	//This variable is used to store the pending Bids list in Array
	var lReqFinalArr []nseipo.ExchangeReqStruct

	//call the filteripoSch to get the list of IPO available to get applications
	lSymbolsArr, lErr1 := FilterIpoSch()
	if lErr1 != nil {
		log.Println("ISCP01", lErr1)
		return lReqFinalArr, lErr1
	} else {
		// for looping the symbols array
		for lpendingIdx := 0; lpendingIdx < len(pPendingArr); lpendingIdx++ {
			for lSymIdx := 0; lSymIdx < len(lSymbolsArr); lSymIdx++ {
				// for get the current time available ipos
				if lSymbolsArr[lSymIdx] == pPendingArr[lpendingIdx].Symbol {
					lReqFinalArr = append(lReqFinalArr, pPendingArr[lpendingIdx])
					// log.Println("lReqFinalArr", lReqFinalArr)
					break
				}
			}
		}
	}
	log.Println("ConstructPendingSch (-)")
	return lReqFinalArr, nil
}

/*
Purpose: this method filter the currently available IPO list to getting application
Parameters:

	not applicable

Response:

	    ==========
	    *On Sucess
	    ==========
		[DMART,JDIAL,YESBANK]
	    =========
	    !On Error
	    =========
	    [],error

Author: Pavithra
Date: 07JUNE2023
*/
func FilterIpoSch() ([]string, error) {
	log.Println("FilterIpoSch (+)")

	//this variable is used to get the symbol
	var lSymbol string
	//this variable is used to get the symbols in an array
	var lSymbolArr []string

	//Establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISFI01", lErr1)
		return lSymbolArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select list.Symbol
		from 
		(select m.Symbol Symbol,(case when (time(now()) between c.StartTime and c.EndTime) then 'Y' else 
		(
		(case when (c.StartTime is null and c.EndTime is null) then 
		(case when (time(now()) between m.DailyStartTime and m.DailyEndTime) then 'I' else 'N' end) else 'S' end)
		) end) value
		from a_ipo_master m,a_ipo_categories c
		where m.Id = c.MasterId
		and c.Code = 'RETAIL'
		and m.BiddingStartDate <= curdate() 
		and m.BiddingEndDate >= curdate() 
		) list
		where list.value = 'Y' or list.value = 'I'`

		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("ISFI02", lErr2)
			return lSymbolArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSymbol)
				if lErr3 != nil {
					log.Println("ISFI03", lErr3)
					return lSymbolArr, lErr3
				} else {
					lSymbolArr = append(lSymbolArr, lSymbol)
				}
			}
		}
	}
	log.Println("FilterIpoSch (-)")
	return lSymbolArr, nil
}

/*
Purpose: This Method fetch the pending order from the order table according to the Segmant(NSE/BSE) and
Broker(novo.flattrade.in).
Parameters:

	not applicable

Response:

	==========
	*On Sucess
	==========

	=========
	!On Error
	=========
	In case of any exception during the execution of this method you will get the
	error details. the calling program should handle the error

Author: Nithish Kumar
Date: 07JUNE2023
*/
func PendingListSch(pBrokerId int, pExchange string) ([]nseipo.ExchangeReqStruct, []IdValueStruct, error) {
	log.Println("PendingListSch (+)")
	//This variable is used to store the pending Bids list in Array
	var lReqArr []nseipo.ExchangeReqStruct
	//This variable is used to store the Bid Header in  Struct
	var lHeadRec nseipo.ExchangeReqStruct
	//This variable is used to store the Bid details in Struct
	var lBidRec nseipo.RequestBidStruct

	// this variable is used to store the Id value of header and details table
	var lIdArr []IdValueStruct

	// this variable is used to get the values of header and details id from the database
	var lIdRec IdValueStruct

	// this variable is used to get the detail id from the database.
	// var lDetailId int

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISPL01", lErr1)
		return lReqArr, lIdArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `SELECT oh.id,oh.clientId ,nvl(oh.Symbol ,'') ,nvl(oh.applicationNo,'') ,nvl(oh.category,'') ,
		nvl(oh.clientName,'') ,nvl(oh.depository,'') ,nvl(oh.dpId,'') ,nvl(oh.clientBenId,'') ,
		nvl(oh.nonASBA,'') ,nvl(oh.pan,'') ,nvl(oh.referenceNumber,'') ,nvl(oh.allotmentMode,'') ,
		nvl(oh.upiFlag,'') ,nvl(oh.upi,'') ,nvl(oh.bankCode,'') ,nvl(oh.locationCode,'') ,
		nvl(oh.time_Stamp,''),nvl(od.activityType,'') ,nvl(od.bidReferenceNo,'') ,nvl(od.req_quantity,'') ,
		nvl(od.atCutOff,'') ,nvl(od.req_price,'') ,nvl(od.req_amount,'') ,nvl(od.remark,'') ,nvl(od.series,'') 
		FROM a_ipo_order_header oh, a_ipo_orderdetails od 
		WHERE oh.Id  = od.HeaderId 
		and od.status = 'pending'
		and oh.Exchange = ?
		and oh.brokerId = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pExchange, pBrokerId)
		if lErr2 != nil {
			log.Println("ISPL02", lErr2)
			return lReqArr, lIdArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lIdRec.HeaderId, &lIdRec.ClientId, &lHeadRec.Symbol, &lHeadRec.ApplicationNo, &lHeadRec.Category, &lHeadRec.ClientName, &lHeadRec.Depository,
					&lHeadRec.DpId, &lHeadRec.ClientBenId, &lHeadRec.NonASBA, &lHeadRec.Pan, &lHeadRec.ReferenceNo, &lHeadRec.AllotmentMode,
					&lHeadRec.UpiFlag, &lHeadRec.Upi, &lHeadRec.BankCode, &lHeadRec.LocationCode, &lHeadRec.TimeStamp, &lBidRec.ActivityType,
					&lBidRec.BidReferenceNo, &lBidRec.Quantity, &lBidRec.AtCutOff, &lBidRec.Price, &lBidRec.Amount, &lBidRec.Remark, &lBidRec.Series)
				if lErr3 != nil {
					log.Println("ISPL03", lErr3)
					return lReqArr, lIdArr, lErr3
				} else {

					if len(lReqArr) == 0 {

						lHeadRec.Bids = append(lHeadRec.Bids, lBidRec)
						//------------------***------------------------

						lReqArr = append(lReqArr, lHeadRec)
						// log.Println("lReqArr", lReqArr)
						//-------------------***-----------------------
						lIdArr = append(lIdArr, lIdRec)
						// log.Println("lIdArr", lIdArr)
						lHeadRec.Bids = []nseipo.RequestBidStruct{}

					} else {

						//-------------------------------------
						// log.Println("Query Readed...")
						for lReqArrIdx := 0; lReqArrIdx < len(lReqArr); lReqArrIdx++ {
							// log.Println("lReqArrIdx", lReqArrIdx)
							//-------------------------------------=----------------------
							// log.Println("lReqArr[lReqArrIdx].ApplicationNo == lHeadRec.ApplicationNo", lReqArr[lReqArrIdx].ApplicationNo, lHeadRec.ApplicationNo)
							//-=----------------------------------------------------------
							if lReqArr[lReqArrIdx].ApplicationNo == lHeadRec.ApplicationNo {
								// lBidArr = append(lBidArr, lBidRec)
								lReqArr[lReqArrIdx].Bids = append(lReqArr[lReqArrIdx].Bids, lBidRec)
								// log.Println(" if", lReqArr[lReqArrIdx])
								break
							}
							if lReqArrIdx == len(lReqArr)-1 {
								// -----------------------------
								lHeadRec.Bids = append(lHeadRec.Bids, lBidRec)
								lReqArr = append(lReqArr, lHeadRec)
								// log.Println("else", lReqArr)

								lIdArr = append(lIdArr, lIdRec)
								// log.Println("else lIdArr", lIdArr)
								lHeadRec = nseipo.ExchangeReqStruct{}
								break
							}
						}
					}
				}
			}
		}
	}
	log.Println("PendingListSch (-)")
	return lReqArr, lIdArr, nil
}

/*
Purpose:This Method is used to get the bid tracking id for all penidng applications
Parameters:

	pReqArr

Response:

	    ==========
	    *On Sucess
	    ==========

	    =========
	    !On Error
	    =========
		[],error

Author: Pavithra
Date: 07JUNE2023
*/
func GetBidTrackIdSch(pBrokerId int, pReqArr []nseipo.ExchangeReqStruct) ([]BidTrackStruct, error) {
	log.Println("GetBidTrackIdSch (+)")

	//this variable is used to get the bid tracking id for all applications
	var lBidTrackArr []BidTrackStruct
	//this variable is used to get the bid details for an application is this struct
	var lBidTrackRec BidTrackStruct
	//get the bidtracking id
	var lGetBidId BidRefStruct

	//Establish a dtatabase connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISGBI01", lErr1)
		return lBidTrackArr, lErr1
	} else {
		defer lDb.Close()

		for lReqArrIdx := 0; lReqArrIdx < len(pReqArr); lReqArrIdx++ {
			for lBidArrIdx := 0; lBidArrIdx < len(pReqArr[lReqArrIdx].Bids); lBidArrIdx++ {

				lCoreString := `select b.applicationNo,b.clientId,b.bidRefNo,b.Id 
							from a_bidtracking_table b
							where b.applicationNo = ?
							and b.bidRefNo = ?
							and b.activityType = ?
							and b.brokerId = ?
							order by b.UpdatedDate desc
							limit 1`

				lRows, lErr2 := lDb.Query(lCoreString, pReqArr[lReqArrIdx].ApplicationNo, pReqArr[lReqArrIdx].Bids[lBidArrIdx].BidReferenceNo,
					pReqArr[lReqArrIdx].Bids[lBidArrIdx].ActivityType, pBrokerId)
				if lErr2 != nil {
					log.Println("ISGBI02", lErr2)
					return lBidTrackArr, lErr2
				} else {
					for lRows.Next() {
						lErr3 := lRows.Scan(&lBidTrackRec.AppNo, &lBidTrackRec.ClientId, &lGetBidId.BidRefNo, &lGetBidId.BidId)
						if lErr3 != nil {
							log.Println("ISGBI03", lErr3)
							return lBidTrackArr, lErr3
						} else {
							// log.Println("Inside Else in getbidTrackId")
							//-------------------------------------
							//check the length of array to store the value
							if len(lBidTrackArr) == 0 {
								lBidTrackRec.BidIdArr = append(lBidTrackRec.BidIdArr, lGetBidId)
								lBidTrackArr = append(lBidTrackArr, lBidTrackRec)
							} else {
								//-------------------------------------
								//loop the lBidTrackArr to store if the application Details is already exists or not
								for lBidIdx := 0; lBidIdx < len(lBidTrackArr); lBidIdx++ {
									//if application deatils is already exists store the details on the corresponding index
									if lBidTrackArr[lBidIdx].AppNo == lBidTrackRec.AppNo && lBidTrackArr[lBidIdx].ClientId == lBidTrackRec.ClientId {
										lBidTrackArr[lBidIdx].BidIdArr = append(lBidTrackArr[lBidIdx].BidIdArr, lGetBidId)
										break
									} else {
										//if application deatils is not exists newly append the details in lBidTRackArr
										lBidTrackRec.BidIdArr = append(lBidTrackRec.BidIdArr, lGetBidId)
										lBidTrackArr = append(lBidTrackArr, lBidTrackRec)
										break
									}
								}
							}
						}
					}
				}
			}
		}
	}
	log.Println("GetBidTrackIdSch (-)")
	return lBidTrackArr, nil
}

/*
Purpose: This method is used to update the records
Parameters:

	pReqRec,pRespRec,pIdArr,pBidTrack,pUser

Response:

	    ==========
	    *On Sucess
	    ==========
		nil
	    =========
	    !On Error
	    =========
		error

Author: Pavithra
Date: 19JUNE2023
*/
func UpdateRecordSch(pReqRec nseipo.ExchangeReqStruct, pRespRec nseipo.ExchangeRespStruct, pIdArr []IdValueStruct, pBidTrack []BidTrackStruct, pUser string) error {
	log.Println("UpdateRecordSch (+)")

	// Changing the date format DD-MM-YYYY into YYYY-MM-DD of timestamp value.
	if pRespRec.TimeStamp != "" {
		lDate := strings.Split(pRespRec.TimeStamp, " ")
		lTimeStamp, _ := time.Parse("02-01-2006", lDate[0])
		pRespRec.TimeStamp = lTimeStamp.Format("2006-01-02") + " " + lDate[1]
	} else if pRespRec.TimeStamp == "" {
		lTime := time.Now()
		pRespRec.TimeStamp = lTime.Format("2006-01-02 15:04:05")
	}

	//this method is used to updating the header values in Response struct in orderheader table
	lUpdatedHeadId, lErr1 := UpdateHeaderSch(pRespRec, pUser)
	if lErr1 != nil {
		log.Println("ISURS01", lErr1)
		return lErr1
	} else {
		// this loop is used to get the header details
		for lBidIdx := 0; lBidIdx < len(pIdArr); lBidIdx++ {
			// check the updatedheader id equals with already taken id Array
			if lUpdatedHeadId == pIdArr[lBidIdx].HeaderId {
				// call update details to update the bid detail records in orderdetails table
				lErr2 := UpdateDetailSch(pReqRec.Bids, pUser, lUpdatedHeadId, pRespRec)
				if lErr2 != nil {
					log.Println("ISURS02", lErr2)
					return lErr2
				}
			}
		}

		// this loop is used to get the  bidtracking table id
		for lBidTrackIdx := 0; lBidTrackIdx < len(pBidTrack); lBidTrackIdx++ {
			//check the application No with already taken Array Vale
			if pRespRec.ApplicationNo == pBidTrack[lBidTrackIdx].AppNo {
				//call the updateBidTrackSch method to update bidtracking table
				lErr3 := UpdateBidTrackSch(pReqRec.Bids, pUser, pBidTrack[lBidTrackIdx], pRespRec)
				if lErr3 != nil {
					log.Println("ISURS03", lErr3)
					return lErr3
				}
			}
		}
	}
	log.Println("UpdateRecordSch (-)")
	return nil
}

/*
Purpose:This method updating the order head values in order header table.
Parameters:

	pRespArr,pHeaderId,PClientId,pDetailIdArr,pBidTrackIdArr

Response:

	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:Pavithra
Date: 19JUNE2023
*/
func UpdateHeaderSch(pRespRec nseipo.ExchangeRespStruct, pUser string) (int, error) {
	log.Println("UpdateHeaderSch (+)")

	//this variable is used to get the updated header value and error
	var lId int
	var lErr2 error

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISUHS01", lErr1)
		return lId, lErr1
	} else {
		defer lDb.Close()

		//call getHeaderId method to get the headerId for the application No
		lId, lErr2 = getHeaderId(pRespRec.ApplicationNo, pRespRec.Symbol)
		if lErr2 != nil {
			log.Println("ISUHS02", lErr2)
			return lId, lErr2
		} else {

			lSqlString := `update a_ipo_order_header h
							set h.Symbol = ?,h.applicationNo = ?,h.category = ?,h.clientName = ?,h.depository = ?,
							h.dpId = ?,h.clientBenId = ?,h.nonASBA = ?,h.pan = ?,h.referenceNumber = ?,
							h.allotmentMode = ?,h.upiFlag = ?,h.upi = ?,h.bankCode = ?,h.locationCode = ?,
							h.bankAccount = ?,h.ifsc = ?,h.subBrokerCode = ?,h.time_Stamp = ?,
							h.status = ?,h.dpVerStatusFlag = ?,h.dpVerFailCode = ?,h.dpVerReason = ?,
							h.upiPaymentStatusFlag = ?,h.upiAmtBlocked = ?,h.reasonCode = ?,h.reason = ?,
							h.UpdatedBy = ?,h.UpdatedDate = now()
							where h.Id = ? and h.applicationNo = ?`
			_, lErr3 := lDb.Exec(lSqlString, pRespRec.Symbol, pRespRec.ApplicationNo, pRespRec.Category, pRespRec.ClientName, pRespRec.Depository, pRespRec.DpId,
				pRespRec.ClientBenId, pRespRec.NonASBA, pRespRec.Pan, pRespRec.ReferenceNo, pRespRec.AllotmentMode, pRespRec.UpiFlag,
				pRespRec.Upi, pRespRec.BankCode, pRespRec.LocationCode, pRespRec.BankAccount, pRespRec.IFSC, pRespRec.SubBrokerCode,
				pRespRec.TimeStamp, pRespRec.Status, pRespRec.DpVerStatusFlag, pRespRec.DpVerFailCode, pRespRec.DpVerReason,
				pRespRec.UpiPaymentStatusFlag, pRespRec.UpiAmtBlocked, pRespRec.ReasonCode,
				pRespRec.Reason, pUser, lId, pRespRec.ApplicationNo)

			if lErr3 != nil {
				log.Println("ISUHS02", lErr3)
				return lId, lErr3
			}
		}
	}
	log.Println("UpdateHeaderSch (-)")
	return lId, nil
}

/*
Purpose:This method fetch the application no id.
Parameters:

	pAppNo,pSymbol

Response:

	==========
	*On Sucess
	==========
	1,nil

	==========
	*On Error
	==========
	0,error

Author:Pavithra
Date: 19JUNE2023
*/
func getHeaderId(pAppNo string, pSymbol string) (int, error) {
	log.Println("getHeaderId (+)")

	//get the application number id from orderheader table
	var lHeaderId int
	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISGHI01", lErr1)
		return lHeaderId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select h.Id
		from a_ipo_order_header h
		where h.applicationNo = ? and h.Symbol = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pAppNo, pSymbol)
		if lErr2 != nil {
			log.Println("ISGHI02", lErr2)
			return lHeaderId, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in Array of struct
			for lRows.Next() {
				lErr3 := lRows.Scan(&lHeaderId)
				if lErr3 != nil {
					log.Println("ISGHI03", lErr3)
					return lHeaderId, lErr3
				}
			}
		}
	}
	log.Println("getHeaderId (-)")
	return lHeaderId, nil
}

/*
Purpose:This method updating the bid details in order detail table.
Parameters:

	pRespBidArr,pClientId,pDetailIdArr,pBidTrackIdArr

Response:

	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:Pavithra
Date: 19JUNE2023
*/
func UpdateDetailSch(pReqBidArr []nseipo.RequestBidStruct, pUser string, pHeaderId int, pRespRec nseipo.ExchangeRespStruct) error {
	log.Println("UpdateDetailSch (+)")

	//get the bid array in pRespBidArr variable
	pRespBidArr := pRespRec.Bids

	// establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISUDS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// looping the req array
		for lReqRecIdx := 0; lReqRecIdx < len(pReqBidArr); lReqRecIdx++ {
			// if the bid is canceled update the status as success and change the cancel flag to 'Y'
			if len(pRespRec.Bids) == 0 && pRespRec.Status == "success" {
				if pReqBidArr[lReqRecIdx].ActivityType == "cancel" {

					// ! ----------------
					lSqlString1 := `update a_ipo_orderdetails d
										set d.activityType = ?,d.status = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
										where d.headerId = ? 
										and d.bidReferenceNo = ?`

					_, lErr2 := lDb.Exec(lSqlString1, pReqBidArr[lReqRecIdx].ActivityType, "success",
						pUser, pHeaderId, pReqBidArr[lReqRecIdx].BidReferenceNo)
					if lErr2 != nil {
						log.Println("ISUDS02", lErr2)
						return lErr2
					}
				}
			}
			// looping the response bid array
			for lRespRecIdx := 0; lRespRecIdx < len(pRespBidArr); lRespRecIdx++ {
				if pRespBidArr[lRespRecIdx].ActivityType == "new" {
					if pRespBidArr[lRespRecIdx].ActivityType == pReqBidArr[lReqRecIdx].ActivityType &&
						pRespBidArr[lRespRecIdx].Remark == pReqBidArr[lReqRecIdx].Remark {

						lSqlString2 := `update a_ipo_orderdetails d
							set d.activityType = ?,d.bidReferenceNo = ?,d.series = ?,d.atCutOff = ?,
							d.resp_quantity = ?,d.resp_price = ?,d.resp_amount = ?,d.remark = ?,
							d.status = ?,d.reasonCode = ?,d.reason = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
							where d.remark = ? and d.headerId = ?`

						_, lErr3 := lDb.Exec(lSqlString2, pRespBidArr[lRespRecIdx].ActivityType, pRespBidArr[lRespRecIdx].BidReferenceNo,
							pRespBidArr[lRespRecIdx].Series, pRespBidArr[lRespRecIdx].AtCutOff, pRespBidArr[lRespRecIdx].Quantity,
							pRespBidArr[lRespRecIdx].Price, pRespBidArr[lRespRecIdx].Amount, pRespBidArr[lRespRecIdx].Remark,
							pRespBidArr[lRespRecIdx].Status, pRespBidArr[lRespRecIdx].ReasonCode, pRespBidArr[lRespRecIdx].Reason,
							pUser, pRespBidArr[lRespRecIdx].Remark, pHeaderId)
						if lErr3 != nil {
							log.Println("ISUDS03", lErr3)
							return lErr3
						}
					}
				} else if pRespBidArr[lRespRecIdx].ActivityType == "modify" {
					if pRespBidArr[lRespRecIdx].ActivityType == pReqBidArr[lReqRecIdx].ActivityType &&
						pRespBidArr[lRespRecIdx].Remark == pReqBidArr[lReqRecIdx].Remark {

						for lRespArrIdx := 0; lRespArrIdx < len(pRespBidArr); lRespArrIdx++ {

							// log.Println("pHeaderId == pIdRec.Head ", pHeaderId, "==", pIdArr[lIdArrIdx].Head)
							// if pHeaderId == pIdArr[lIdArrIdx].Head {

							// for lDetailIdx := 0; lDetailIdx < len(pIdArr[lRespBidArrIdx].Detail); lDetailIdx++ {

							lSqlString3 := `update a_ipo_orderdetails d
												set d.activityType = ?,d.bidReferenceNo = ?,d.series = ?,d.atCutOff = ?,
												d.resp_quantity = ?,d.resp_price = ?,d.resp_amount = ?,d.remark = ?,
												d.status = ?,d.reasonCode = ?,d.reason = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
												where d.remark = ? and d.headerId = ? and d.bidReferenceNo = ?`

							_, lErr4 := lDb.Exec(lSqlString3, pRespBidArr[lRespRecIdx].ActivityType, pRespBidArr[lRespRecIdx].BidReferenceNo,
								pRespBidArr[lRespRecIdx].Series, pRespBidArr[lRespRecIdx].AtCutOff, pRespBidArr[lRespRecIdx].Quantity,
								pRespBidArr[lRespRecIdx].Price, pRespBidArr[lRespRecIdx].Amount, pRespBidArr[lRespRecIdx].Remark,
								pRespBidArr[lRespRecIdx].Status, pRespBidArr[lRespRecIdx].ReasonCode, pRespBidArr[lRespRecIdx].Reason,
								pUser, pRespBidArr[lRespRecIdx].Remark, pHeaderId, pRespBidArr[lRespRecIdx].BidReferenceNo)
							if lErr4 != nil {
								log.Println("ISUDS04", lErr4)
								return lErr4
							}
						}
					}
				}
			}
		}
	}
	log.Println("UpdateDetailSch (-)")
	return nil
}

/*
Purpose:This method updating the bid details in bid tracking table.
Parameters:

	pRespBidRec,pClientId,pBidTrackIdRec

Response:

	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:Pavithra
Date: 19JUNE2023
*/
func UpdateBidTrackSch(pReqBidArr []nseipo.RequestBidStruct, pUser string, pBidIdRec BidTrackStruct, pRespRec nseipo.ExchangeRespStruct) error {
	log.Println("UpdateBidTrackSch (+)")

	//get the bid array in pRespBidArr variable
	pRespBid := pRespRec.Bids

	// establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ISUBT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for lReqRecIdx := 0; lReqRecIdx < len(pReqBidArr); lReqRecIdx++ {
			//
			if len(pRespRec.Bids) == 0 && pRespRec.Status == "success" {
				if pReqBidArr[lReqRecIdx].ActivityType == "cancel" {

					// ! ----------------
					lSqlString1 := `update a_bidtracking_table b
									set b.activityType = ?,b.status = ?,b.UpdatedBy = ?,b.UpdatedDate = now()
									where b.bidRefNo = ?`

					_, lErr2 := lDb.Exec(lSqlString1, pReqBidArr[lReqRecIdx].ActivityType, "success",
						pUser, pReqBidArr[lReqRecIdx].BidReferenceNo)
					if lErr2 != nil {
						// log.Println("Error updating into database (details)")
						log.Println("ISUBT01", lErr2)
						return lErr2
					}
				}
			}
			//get the bid track id array in lTrack
			lTrack := pBidIdRec.BidIdArr
			//
			for lRespRecIdx := 0; lRespRecIdx < len(pRespBid); lRespRecIdx++ {
				//
				for lTrackIdIdx := 0; lTrackIdIdx < len(lTrack); lTrackIdIdx++ {
					//
					if strconv.Itoa(lTrack[lTrackIdIdx].BidRefNo) == pRespBid[lRespRecIdx].Remark {

						if pReqBidArr[lReqRecIdx].ActivityType == "new" {
							if pRespBid[lRespRecIdx].ActivityType == pReqBidArr[lReqRecIdx].ActivityType &&
								pRespBid[lRespRecIdx].Remark == pReqBidArr[lReqRecIdx].Remark {

								lSqlString2 := `update a_bidtracking_table b
									set b.bidRefNo = ?,b.status = ?,b.UpdatedBy = ?,b.UpdatedDate = now() 
									where b.applicationNo = ? and b.clientId = ?
									and b.bidRefNo = ? and b.activityType = ? and b.Id = ?`

								_, lErr3 := lDb.Exec(lSqlString2, pRespBid[lRespRecIdx].BidReferenceNo, pRespBid[lRespRecIdx].Status,
									pUser, pBidIdRec.AppNo, pBidIdRec.ClientId, pRespBid[lRespRecIdx].Remark,
									pRespBid[lRespRecIdx].ActivityType, lTrack[lTrackIdIdx].BidId)
								if lErr3 != nil {
									log.Println("ISUBT02", lErr3)
									return lErr3
								}
							}
							//
						} else if pReqBidArr[lReqRecIdx].ActivityType == "modify" {
							if pRespBid[lRespRecIdx].ActivityType == pReqBidArr[lReqRecIdx].ActivityType &&
								pRespBid[lRespRecIdx].Remark == pReqBidArr[lReqRecIdx].Remark {

								if strconv.Itoa(lTrack[lTrackIdIdx].BidRefNo) == pRespBid[lRespRecIdx].Remark {

									lSqlString3 := `update a_bidtracking_table b
									set b.bidRefNo = ?,b.status = ?,b.UpdatedBy = ?,b.UpdatedDate = now() 
									where b.applicationNo = ? and b.clientId = ?
									and b.bidRefNo = ? and b.activityType = ? and b.Id = ?`

									_, lErr4 := lDb.Exec(lSqlString3, pRespBid[lRespRecIdx].BidReferenceNo, pRespBid[lRespRecIdx].Status,
										pUser, pBidIdRec.AppNo, pBidIdRec.ClientId, pRespBid[lRespRecIdx].Remark,
										pRespBid[lRespRecIdx].ActivityType, lTrack[lTrackIdIdx].BidId)
									if lErr4 != nil {
										log.Println("ISUBT03", lErr4)
										return lErr4
									}
								}
							}
						}
					}
				}
			}
		}
	}
	log.Println("UpdateBidTrackSch (-)")
	return nil
}

/*
Purpose:This method fetch the application number from  a_ipo_order_header table
Parameters:

	pPendingArr


Response:

	==========
	*On Sucess
	==========


	==========
	*On Error
	==========

Author:Kavya Dharshani
Date: 07SEP2023
*/
func GetApiMaster(pPendingArr []nseipo.ExchangeReqStruct) ([]placeorder.IpoOrderStruct, error) {
	log.Println("GetApiMaster (+)")

	//This variable is used to get  the ApplicationNumber  in Array
	var lAppRec placeorder.IpoOrderStruct

	//This variable is used to append the array in the IpoOrderStruct Structure
	var lorderArr []placeorder.IpoOrderStruct

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)

	if lErr1 != nil {
		log.Println("ISGAM01", lErr1)
		return lorderArr, lErr1
	} else {
		defer lDb.Close()

		for lpendingIdx := 0; lpendingIdx < len(pPendingArr); lpendingIdx++ {
			// log.Println("pPendingArr[lpendingIdx].ApplicationNo", pPendingArr[lpendingIdx].ApplicationNo)

			lCoreString := `select h.applicationNo, h.ClientEmail , h.clientId  , h.clientName, h.status , m.Symbol ,
			max(d.req_amount),date(d.CreatedDate), d.activityType  
			from a_ipo_order_header h , a_ipo_master m,a_ipo_orderdetails d
			where h.MasterId  = m.Id
			and h.Id = d.headerId 
			and h.applicationNo  = ?
			group by d.headerId `

			lRows, lErr2 := lDb.Query(lCoreString, pPendingArr[lpendingIdx].ApplicationNo)
			if lErr2 != nil {
				log.Println("ISGAM02", lErr2)
				return lorderArr, lErr2
			} else {
				for lRows.Next() {
					lErr3 := lRows.Scan(&lAppRec.ApplicationNumber, &lAppRec.ClientEmail, &lAppRec.ClientId, &lAppRec.ClientName, &lAppRec.Status, &lAppRec.Symbol, &lAppRec.Amount, &lAppRec.OrderDate, &lAppRec.ActivityType)
					if lErr3 != nil {
						log.Println("ISGAM03", lErr3)
						return lorderArr, lErr3
					} else {
						lorderArr = append(lorderArr, lAppRec)
					}
				}
			}
		}
	}
	log.Println("GetApiMaster (-)")
	return lorderArr, nil
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------
// func ConstructMail(lAppArr IpoOrderStruct, pStatus string) (emailUtil.EmailInput, error) {
// 	type dynamicEmailStruct struct {
// 		Name              string
// 		Symbol            string
// 		ApplicationNumber string
// 		OrderDate         string
// 		Unit              string
// 		Price             string
// 		Amount            int
// 		Activity          string
// 		Status            string
// 	}

// 	var lIpoEmailContent emailUtil.EmailInput
// 	config := common.ReadTomlConfig("./toml/emailconfig.toml")
// 	lIpoEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// 	lIpoEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
// 	lIpoEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 	lIpoEmailContent.Subject = "IPO Order"
// 	html := "html/IpoOrderTemplate.html"

// 	lTemp, lErr := template.ParseFiles(html)
// 	if lErr != nil {
// 		log.Println("IMP03", lErr)
// 		return lIpoEmailContent, lErr
// 	} else {
// 		var lTpl bytes.Buffer
// 		var lDynamicEmailVal dynamicEmailStruct

// 		lDynamicEmailVal.Name = lAppArr.ClientName
// 		lDynamicEmailVal.Symbol = lAppArr.Symbol
// 		lDynamicEmailVal.Amount = lAppArr.Amount
// 		lDynamicEmailVal.OrderDate = lAppArr.OrderDate
// 		lDynamicEmailVal.ApplicationNumber = lAppArr.ApplicationNumber
// 		lDynamicEmailVal.Activity = lAppArr.ActivityType
// 		lDynamicEmailVal.Status = pStatus
// 		// lIpoEmailContent.ToEmailId = lAppArr.ClientEmail
// 		log.Println("lAppArr.ClientEmail", lAppArr.ClientEmail)
// 		lClientEmail, lErr1 := clientDetail.GetClientEmailId(lAppArr.ClientId)
// 		if lErr1 != nil {
// 			return lIpoEmailContent, lErr1
// 		}
// 		// lIpoEmailContent.ToEmailId = "prashanth.s@fcsonline.co.in"
// 		lIpoEmailContent.ToEmailId = lClientEmail

// 		// log.Println(lDynamicEmailVal, "dynamic struct")
// 		lTemp.Execute(&lTpl, lDynamicEmailVal)
// 		lEmailbody := lTpl.String()

// 		lIpoEmailContent.Body = lEmailbody
// 	}
// 	return lIpoEmailContent, nil
// }
