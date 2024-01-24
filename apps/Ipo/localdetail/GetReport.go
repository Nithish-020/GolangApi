package localdetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/SGB/localdetails"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// This  Structure is used to get Request details

// This Structure is used to get Report details
type ReportRespStruct struct {
	MasterId      int    `json:"masterId"`
	Symbol        string `json:"symbol"`
	ApplicationNo string `json:"applicationNo"`
	ApplyDate     string `json:"applyDate"`
	AppliedTime   string `json:"appliedTime"`
	Status        string `json:"status"`
	ClientId      string `json:"clientId"`
	Exchange      string `json:"exchange"`
	Category      string `json:"category"`
}

// Response Sturcture for GetReport API
type RespStruct struct {
	IpoArr []ReportRespStruct `json:"ipoArr"`
	// SgbArr  []ReportRespStruct             `json:"sgbArr"`
	SgbArr []localdetails.SgbOrderHistoryStruct `json:"sgbArr"`
	Status string                               `json:"status"`
	ErrMsg string                               `json:"errMsg"`
}

/*
Purpose:This Function is used to Get reports for filterations
Parameters:

	{
		"clientId": "FT12345678",
		"fromDate": "2023-06-06",
		"toDate": "2023-06-07",
		"symbol": null,
	}

Response:

		=========
		*On Sucess
		=========
		{"reportArr": [
	        {
	            "symbol": "MMIPO26",
	            "applicationNo": "FT00006730554",
	            "bidRefNo": "2023060700000004",
	            "activityType": "new",
	            "quantity": 10,
	            "price": 2000,
	            "applyDate": "2023-06-07",
	            "status": "success"
	        }
			],
			"status": "S",
			"errMsg":""
		}
		=========
		!On Error
		=========
		{
			"status": "E",
			"errMsg": "Can't able to get data from database"
		}

Author: Nithish Kumar
Date: 07JUNE2023
*/
func GetReport(w http.ResponseWriter, r *http.Request) {
	log.Println("GetReport (+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {

		// This variable is used to pass the request struct to the query
		var lReportRec localdetails.ReportReqStruct
		// This variable helps to store the response and send it to front
		var lRespRec RespStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/report")
		if lErr1 != nil {
			log.Println("LGIH01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIH01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGR01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGR02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// Read the Request body in lBody
		lBody, lErr2 := ioutil.ReadAll(r.Body)
		if lErr2 != nil {
			log.Println("LGR03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGR03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGR03", "Error Occur in Reading inputs!"))
			return
		} else {
			// Unmarshal the Request body in lReportRec structure
			lErr3 := json.Unmarshal(lBody, &lReportRec)
			if lErr3 != nil {
				log.Println("LGR04", lErr3)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "LGR04" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("LGR04", "Error Occur in Reading inputs!.Please try after sometime"))
				return
			} else {
				// this method is added to check the given symbol is valid or not
				//and also added validation for empty array in the case of invalid clientId
				//by Pavithra
				log.Println("lReportRec", lReportRec)
				if lReportRec.Symbol != "" {
					lFlag, lErr4 := CheckSymbolValid(lReportRec)
					if lErr4 != nil {
						log.Println("LGR05", lErr4)
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "LGR05" + lErr4.Error()
						fmt.Fprintf(w, helpers.GetErrorString("LGR05", "Error Occur in  getting datas!.Please try after sometime"))
						return
					} else {
						if lFlag == common.ErrorCode {
							log.Println("LGR06", "")
							lRespRec.Status = common.ErrorCode
							lRespRec.ErrMsg = "Symbol is Not Valid"
							fmt.Fprintf(w, helpers.GetErrorString("LGR06", "Symbol is Not Valid"))
							return
						} else {
							IpoArr, SgbArr, lErr5 := GetReportModule(lReportRec, lBrokerId)
							if lErr5 != nil {
								log.Println("LGR07", lErr5)
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "LGR07" + lErr5.Error()
								fmt.Fprintf(w, helpers.GetErrorString("LGR07", "Error Occur in getting datas!.Please try after sometime"))
								return
							} else {
								lRespRec.IpoArr = IpoArr
								lRespRec.SgbArr = SgbArr
							}
						}
					}
				} else {
					IpoArr, SgbArr, lErr6 := GetReportModule(lReportRec, lBrokerId)
					if lErr6 != nil {
						log.Println("LGR08", lErr6)
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "LGR08" + lErr6.Error()
						fmt.Fprintf(w, helpers.GetErrorString("LGR08", "Error Occur in getting datas!.Please try after sometime"))
						return
					} else {
						lRespRec.IpoArr = IpoArr
						lRespRec.SgbArr = SgbArr
						// if IpoArr == nil {
						// 	lRespRec.Status = common.ErrorCode
						// 	lRespRec.ErrMsg = "Records not Found"
						// } else {
						// 	lRespRec.IpoArr = IpoArr
						// }
						// if SgbArr == nil {
						// 	lRespRec.Status = common.ErrorCode
						// 	lRespRec.ErrMsg = "Records not Found"

						// } else {
						// 	lRespRec.SgbArr = SgbArr
						// }

						if lReportRec.Module == "Ipo" && lReportRec.ClientId != "" {
							if len(lRespRec.IpoArr) == 0 || lRespRec.IpoArr == nil {
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "Records not Found"
							}
						}
						if lReportRec.Module == "Sgb" && lReportRec.ClientId != "" {
							if len(lRespRec.SgbArr) == 0 || lRespRec.SgbArr == nil {
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "Records not Found"
							}
						}
					}
				}

			}
		}
		ldata, lErr7 := json.Marshal(lRespRec)
		if lErr7 != nil {
			log.Println("LGR09", lErr7)
			fmt.Fprintf(w, helpers.GetErrorString("LGR09", "Issue in Getting your Repots!"))
			return
		} else {
			fmt.Fprintf(w, string(ldata))
		}
		log.Println("GetReport (-)", r.Method)
	}
}

// ------ commented by kavya -------

// func GetReportModule(pReportRec ReportReqStruct, pBrokerId int) ([]ReportRespStruct, []ReportRespStruct, error) {
// 	log.Println("GetReportModule (+)")

// 	var lIpo, lSgb ReportRespStruct
// 	var lIpoArr, lSgbArr []ReportRespStruct

// 	// to Establish A database conncetion, call the LocalDbconnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("LGRM01", lErr1)
// 		return lIpoArr, lSgbArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		if pReportRec.Module == "Ipo" {

// 			// lCoreString := `SELECT aioh.MasterId,aioh.Symbol,aioh.applicationNo,nvl(DATE_FORMAT(aioh.CreatedDate , '%d-%b-%Y'),'') applyDate,
// 			// 				nvl(TIME_FORMAT(aioh.CreatedDate, '%h:%i:%s %p'),'') applytime,
// 			// 				(case when aioh.cancelFlag = 'Y' then 'user cancelled' else nvl(aioh.status ,'') end) flag,aioh.clientId
// 			// 				from a_ipo_order_header aioh
// 			// 				where aioh.brokerId = ?
// 			// 				and aioh.CreatedDate between concat (?, ' 00:00:00.000') and concat (?, ' 23:59:59.000')`

// 			// if pReportRec.Symbol != "" && pReportRec.ClientId == "" {
// 			// 	ladditionalCon1 = `or aioh.Symbol  = '` + pReportRec.Symbol + `'`
// 			// } else if pReportRec.ClientId != "" && pReportRec.Symbol == "" {
// 			// 	ladditionalCon1 = `or aioh.clientId = '` + pReportRec.ClientId + `'`
// 			// } else if pReportRec.Symbol != "" && pReportRec.ClientId != "" {
// 			// 	ladditionalCon1 = `or (aioh.Symbol  = '` + pReportRec.Symbol + `' and aioh.clientId = '` + pReportRec.ClientId + `')`
// 			// }

// 			// velmurugan

// 			var ladditionalCon1 string
// 			if pReportRec.FromDate != "" && pReportRec.ToDate != "" {
// 				ladditionalCon1 = `and aioh.CreatedDate between concat ('` + pReportRec.FromDate + ` 00:00:00.000') and concat ('` + pReportRec.ToDate + ` 23:59:59.000')`
// 			}

// 			if pReportRec.Symbol != "" {
// 				ladditionalCon1 = ladditionalCon1 + `and aioh.Symbol  = '` + pReportRec.Symbol + `'`
// 			}
// 			if pReportRec.ClientId != "" {
// 				ladditionalCon1 = ladditionalCon1 + ` and clientId = '` + pReportRec.ClientId + `'`
// 			}

// 			lCoreString := `SELECT aioh.MasterId,aioh.Symbol,aioh.applicationNo,nvl(DATE_FORMAT(aioh.CreatedDate , '%d-%b-%Y'),'') applyDate,
// 							nvl(TIME_FORMAT(aioh.CreatedDate, '%h:%i:%s %p'),'') applytime,
// 							(case when aioh.cancelFlag = 'Y' then 'user cancelled' else nvl(aioh.status ,'') end) flag,aioh.clientId,nvl(aioh.exchange,''),nvl(category,'')
// 							from a_ipo_order_header aioh
// 							where aioh.brokerId = ?
// 							` + ladditionalCon1 + ``

// 			lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
// 			if lErr2 != nil {
// 				log.Println("LGRM02", lErr2)
// 				return lIpoArr, lSgbArr, lErr2
// 			} else {
// 				//Reading the Records
// 				for lRows.Next() {
// 					lErr3 := lRows.Scan(&lIpo.MasterId, &lIpo.Symbol, &lIpo.ApplicationNo, &lIpo.ApplyDate,
// 						&lIpo.AppliedTime, &lIpo.Status, &lIpo.ClientId, &lIpo.Exchange, &lIpo.Category)
// 					if lErr3 != nil {
// 						log.Println("LGRM03", lErr3)
// 						return lIpoArr, lSgbArr, lErr3
// 					} else {
// 						lIpoArr = append(lIpoArr, lIpo)
// 					}
// 				}
// 			}
// 		} else {
// 			ladditionalCon1 := ""
// 			// COMMENTED BY NITHISH BECAUSE THE ORDER NO COLUMN WAS INCORRECT
// 			// lCoreString := `select oh.MasterId,oh.ScripId ,od.OrderNo,nvl(DATE_FORMAT(oh.CreatedDate , '%d-%b-%Y'),'') applyDate,nvl(TIME_FORMAT(oh.CreatedDate, '%h:%i:%s %p'),'') applytime,(case when oh.cancelFlag = 'Y' then 'user cancelled' else nvl(oh.status ,'') end) flag,
// 			// NVL(oh.ClientId,'') ,nvl(oh.exchange,'')
// 			// from a_sgb_orderheader oh JOIN
// 			// a_sgb_orderdetails od ON oh.Id = od.HeaderId
// 			// where oh.brokerId = ?
// 			lCoreString := `select oh.MasterId,oh.ScripId ,od.ReqOrderNo ,nvl(DATE_FORMAT(oh.CreatedDate , '%d-%b-%Y'),'') applyDate,nvl(TIME_FORMAT(oh.CreatedDate, '%h:%i:%s %p'),'') applytime,(case when oh.cancelFlag = 'Y' then 'user cancelled' else nvl(oh.status ,'') end) flag,
// 			NVL(oh.ClientId,'') ,nvl(oh.exchange,''),nvl(od.RespOrderNo,'')
// 			from a_sgb_orderheader oh JOIN
// 			a_sgb_orderdetails od ON oh.Id = od.HeaderId
// 			where oh.brokerId = ?
// 			`

// 			// if pReportRec.Symbol != "" && pReportRec.ClientId == "" {
// 			// 	ladditionalCon1 = `or oh.ScripId  = '` + pReportRec.Symbol + `'`
// 			// } else if pReportRec.ClientId != "" && pReportRec.Symbol == "" {
// 			// 	ladditionalCon1 = `or oh.clientId = '` + pReportRec.ClientId + `'`
// 			// } else if pReportRec.Symbol != "" && pReportRec.ClientId != "" {
// 			// 	ladditionalCon1 = `or (oh.ScripId  = '` + pReportRec.Symbol + `' and oh.clientId = '` + pReportRec.ClientId + `')`
// 			// }

// 			if pReportRec.Symbol != "" && pReportRec.ClientId == "" {
// 				ladditionalCon1 = `and oh.ScripId  = '` + pReportRec.Symbol + `'`
// 			}
// 			if pReportRec.ClientId != "" && pReportRec.Symbol == "" {
// 				ladditionalCon1 = ladditionalCon1 + `and oh.clientId = '` + pReportRec.ClientId + `'`
// 			}
// 			if pReportRec.Symbol != "" && pReportRec.ClientId != "" {
// 				ladditionalCon1 = ladditionalCon1 + `and (oh.ScripId  = '` + pReportRec.Symbol + `' and oh.clientId = '` + pReportRec.ClientId + `')`
// 			}
// 			if pReportRec.FromDate != "" && pReportRec.ToDate != "" {
// 				ladditionalCon1 = ladditionalCon1 + `and oh.CreatedDate between concat ('` + pReportRec.FromDate + ` 00:00:00.000') and concat ('` + pReportRec.ToDate + ` 23:59:59.000')`
// 			}

// 			lCoreString2 := lCoreString + ladditionalCon1
// 			lRows, lErr4 := lDb.Query(lCoreString2, pBrokerId)
// 			if lErr4 != nil {
// 				log.Println("LGRM04", lErr4)
// 				return lIpoArr, lSgbArr, lErr4
// 			} else {
// 				//Reading the Records
// 				for lRows.Next() {
// 					lErr5 := lRows.Scan(&lSgb.MasterId, &lSgb.Symbol, &lSgb.ApplicationNo, &lSgb.ApplyDate, &lSgb.AppliedTime, &lSgb.Status, &lSgb.ClientId, &lSgb.Exchange, &lSgb.ExchOrderNo)
// 					if lErr5 != nil {
// 						log.Println("LGRM05", lErr5)
// 						return lIpoArr, lSgbArr, lErr5
// 					} else {
// 						lSgbArr = append(lSgbArr, lSgb)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetReportModule (-)")
// 	return lIpoArr, lSgbArr, nil
// }

// if pReportRec.Symbol != "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.ScripId  = '` + pReportRec.Symbol + `'`
// }
// if pReportRec.ClientId != "" {
// 	ladditionalCon1 = ladditionalCon1 + ` and h.ClientId= '` + pReportRec.ClientId + `'`
// }

// if pReportRec.Symbol != "" && pReportRec.ClientId == "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.ScripId  = '` + pReportRec.Symbol + `'`
// }
// if pReportRec.ClientId != "" && pReportRec.Symbol == "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.ClientId = '` + pReportRec.ClientId + `'`
// }

// if pReportRec.Symbol != "" && pReportRec.ClientId != "" {
// 	ladditionalCon1 = ladditionalCon1 + `and (h.ScripId  = '` + pReportRec.Symbol + `' and h.ClientId = '` + pReportRec.ClientId + `')`
// }

// if pReportRec.FromDate != "" && pReportRec.ToDate != "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.CreatedDate between concat ('` + pReportRec.FromDate + ` 00:00:00.000') and concat ('` + pReportRec.ToDate + ` 23:59:59.000')`
// }
// lCoreString2 := lCoreString + ladditionalCon1
// lRows, lErr4 := lDb.Query(lCoreString2, pBrokerId)
/// ====
// if pReportRec.Symbol != "" && pReportRec.ClientId == "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.ScripId  = '` + pReportRec.Symbol + `'`
// }
// if pReportRec.ClientId != "" && pReportRec.Symbol == "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.ClientId = '` + pReportRec.ClientId + `'`
// }
// if pReportRec.Symbol != "" && pReportRec.ClientId != "" {
// 	ladditionalCon1 = ladditionalCon1 + `and (h.ScripId  = '` + pReportRec.Symbol + `' and h.ClientId = '` + pReportRec.ClientId + `')`
// }
// if pReportRec.FromDate != "" && pReportRec.ToDate != "" {
// 	ladditionalCon1 = ladditionalCon1 + `and h.CreatedDate between concat ('` + pReportRec.FromDate + ` 00:00:00.000') and concat ('` + pReportRec.ToDate + ` 23:59:59.000')`
// }

func GetReportModule(pReportRec localdetails.ReportReqStruct, pBrokerId int) ([]ReportRespStruct, []localdetails.SgbOrderHistoryStruct, error) {
	log.Println("GetReportModule (+)")

	var lIpo ReportRespStruct
	var lIpoArr []ReportRespStruct
	// var lSgbOrder localdetails.SgbOrderHistoryStruct
	var lSgbResp localdetails.SgbOrderHistoryResp
	var SgbArr2 []localdetails.SgbOrderHistoryStruct

	// lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
	// lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])

	// to Establish A database conncetion, call the LocalDbconnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGRM01", lErr1)
		return lIpoArr, SgbArr2, lErr1
	} else {
		defer lDb.Close()

		if pReportRec.Module == "Ipo" {

			var ladditionalCon1 string
			if pReportRec.FromDate != "" && pReportRec.ToDate != "" {
				ladditionalCon1 = `and aioh.CreatedDate between concat ('` + pReportRec.FromDate + ` 00:00:00.000') and concat ('` + pReportRec.ToDate + ` 23:59:59.000')`
			}

			if pReportRec.Symbol != "" {
				ladditionalCon1 = ladditionalCon1 + `and aioh.Symbol  = '` + pReportRec.Symbol + `'`
			}
			if pReportRec.ClientId != "" {
				ladditionalCon1 = ladditionalCon1 + ` and clientId = '` + pReportRec.ClientId + `'`
			}

			lCoreString := `SELECT aioh.MasterId,aioh.Symbol,aioh.applicationNo,nvl(DATE_FORMAT(aioh.CreatedDate , '%d-%b-%Y'),'') applyDate,
							nvl(TIME_FORMAT(aioh.CreatedDate, '%h:%i:%s %p'),'') applytime,
							(case when aioh.cancelFlag = 'Y' then 'user cancelled' else nvl(aioh.status ,'') end) flag,aioh.clientId,nvl(aioh.exchange,''),nvl(category,'')
							from a_ipo_order_header aioh	 
							where aioh.brokerId = ?
							` + ladditionalCon1 + ``

			lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
			if lErr2 != nil {
				log.Println("LGRM02", lErr2)
				return lIpoArr, SgbArr2, lErr2
			} else {
				//Reading the Records
				for lRows.Next() {
					lErr3 := lRows.Scan(&lIpo.MasterId, &lIpo.Symbol, &lIpo.ApplicationNo, &lIpo.ApplyDate,
						&lIpo.AppliedTime, &lIpo.Status, &lIpo.ClientId, &lIpo.Exchange, &lIpo.Category)
					if lErr3 != nil {
						log.Println("LGRM03", lErr3)
						return lIpoArr, SgbArr2, lErr3
					} else {
						lIpoArr = append(lIpoArr, lIpo)
					}
				}
			}
		} else {

			// version 2

			lSgbOrderHistoryResp, lErr2 := localdetails.GetSGBOrderHistorydetail("", pBrokerId, pReportRec, "GetReport")
			if lErr2 != nil {
				log.Println("LGSH03", lErr2)
				lSgbResp.Status = common.ErrorCode
				lSgbResp.ErrMsg = "LGSH03" + lErr2.Error()
				return lIpoArr, SgbArr2, lErr2
			} else {
				SgbArr2 = lSgbOrderHistoryResp.SgbOrderHistoryArr
			}

			// lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
			// lDiscountTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_DiscountText"])

			// // ladditionalCon1 := ""
			// var lReqUnit, lReqUnitPrice, lAppliedUnit, lAppliedUnitPrice, lAllocatedUnit, lAllocatedUnitPrice string

			// lCoreString := `select h.MasterId,h.ScripId,m.Name,m.Isin,d.ReqOrderNo OrderNo,date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p') orderDate,
			//     concat(DATE_FORMAT(m.BiddingStartDate, '%d %b %y'),' -  ',DATE_FORMAT(m.BiddingEndDate, '%d %b %y')) as DateRange,
			// 	concat(DATE_FORMAT(sm.BiddingStartDate, '%d %b %Y'),' ',TIME_FORMAT(sm.DailyStartTime , '%h:%i%p')) as startDateWithTime,
			// 	concat(DATE_FORMAT(sm.BiddingEndDate, '%d %b %Y'),' ',TIME_FORMAT('` + lCloseTime + `' , '%h:%i%p')) as endDateWithTime,
			//     d.ReqSubscriptionUnit RequestedUnit, d.ReqRate RequestedUnitPrice, nvl(d.RespSubscriptionunit,'') AppliedUnit,
			//     nvl(d.RespRate,'') AppliedUnitPrice,nvl(d.RespSubscriptionunit,'') AllocatedUnit, nvl(d.RespRate,'') AllocatedUnitPrice,
			//     lower(h.Status) OrderStatus,(m.FaceValue - m.MinPrice) DiscountAmount, h.SIvalue,h.SItext,  d.Exchange,d.RespOrderNo,h.ClientId
			//     from a_sgb_master m,a_sgb_orderheader h ,a_sgb_orderdetails d
			//     where m.id = h.MasterId and d.HeaderId = h.Id
			//     and h.brokerId = ?
			// 	AND (? = '' OR h.ScripId = ?)
			// 	AND (? = '' OR h.ClientId = ?)
			// 	AND (? = '' OR h.CreatedDate BETWEEN CONCAT(?,' 00:00:00.000') AND CONCAT(?,' 23:59:59.000'))`

			// lCoreString := `select h.MasterId,h.ScripId,m.Name,m.Isin,d.ReqOrderNo OrderNo,date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p') orderDate,
			//     concat(DATE_FORMAT(m.BiddingStartDate, '%d %b %y'),' -  ',DATE_FORMAT(m.BiddingEndDate, '%d %b %y')) as DateRange,
			// 	concat(DATE_FORMAT(m.BiddingStartDate, '%d %b %Y'),' ',TIME_FORMAT(m.DailyStartTime , '%h:%i%p'),' -  ',
			// 		   DATE_FORMAT(m.BiddingEndDate, '%d %b %Y'),' ',TIME_FORMAT(m.DailyEndTime , '%h:%i%p')) as dateRangeWithTime2,
			//     d.ReqSubscriptionUnit RequestedUnit, d.ReqRate RequestedUnitPrice, nvl(d.RespSubscriptionunit,'') AppliedUnit,
			//     nvl(d.RespRate,'') AppliedUnitPrice,nvl(d.RespSubscriptionunit,'') AllocatedUnit, nvl(d.RespRate,'') AllocatedUnitPrice,
			//     lower(h.Status) OrderStatus,(m.FaceValue - m.MinPrice) DiscountAmount, h.SIvalue,h.SItext,  d.Exchange,d.RespOrderNo,h.ClientId
			//     from a_sgb_master m,a_sgb_orderheader h ,a_sgb_orderdetails d
			//     where m.id = h.MasterId and d.HeaderId = h.Id
			//     and h.brokerId = ?
			// 	AND (? = '' OR h.ScripId = ?)
			// 	AND (? = '' OR h.ClientId = ?)
			// 	AND (? = '' OR h.CreatedDate BETWEEN CONCAT(?,' 00:00:00.000') AND CONCAT(?,' 23:59:59.000'))`

			// log.Println("Symbol", pReportRec.Symbol, lSgbOrder.Symbol)

			// lRows, lErr4 := lDb.Query(lCoreString, pBrokerId, pReportRec.Symbol, pReportRec.Symbol, pReportRec.ClientId, pReportRec.ClientId, pReportRec.FromDate, pReportRec.FromDate, pReportRec.ToDate)
			// if lErr4 != nil {
			// 	log.Println("LGRM04", lErr4)
			// 	return lIpoArr, SgbArr2, lErr4
			// } else {
			// 	//Reading the Records
			// 	for lRows.Next() {
			// 		lErr5 := lRows.Scan(&lSgbOrder.Id, &lSgbOrder.Symbol, &lSgbOrder.Name, &lSgbOrder.Isin, &lSgbOrder.OrderNo, &lSgbOrder.OrderDate, &lSgbOrder.DateRange, &lSgbOrder.DateRangeWithTime, &lReqUnit, &lReqUnitPrice, &lAppliedUnit, &lAppliedUnitPrice, &lAllocatedUnit, &lAllocatedUnitPrice, &lSgbOrder.OrderStatus, &lSgbOrder.DiscountAmt, &lSgbOrder.SIValue, &lSgbOrder.SIText, &lSgbOrder.Exchange, &lSgbOrder.ExchOrderNo, &lSgbOrder.ClientId)
			// 		if lErr5 != nil {
			// 			log.Println("LGRM05", lErr5)
			// 			return lIpoArr, SgbArr2, lErr5
			// 		} else {
			// 			lSgbOrder.RequestedUnit, _ = strconv.Atoi(lReqUnit)
			// 			lSgbOrder.RequestedUnitPrice, _ = strconv.Atoi(lReqUnitPrice)
			// 			lSgbOrder.RequestedAmount = lSgbOrder.RequestedUnit * lSgbOrder.RequestedUnitPrice

			// 			lSgbOrder.AppliedUnit, _ = strconv.Atoi(lAppliedUnit)
			// 			lSgbOrder.AppliedUnitPrice, _ = strconv.Atoi(lAppliedUnitPrice)
			// 			lSgbOrder.AppliedAmount = lSgbOrder.AppliedUnit * lSgbOrder.AppliedUnitPrice

			// 			lSgbOrder.AllotedUnit, _ = strconv.Atoi(lAllocatedUnit)
			// 			lSgbOrder.AllotedUnitPrice, _ = strconv.Atoi(lAllocatedUnitPrice)
			// 			lSgbOrder.AppliedAmount = lSgbOrder.AllotedUnit * lSgbOrder.AllotedUnitPrice

			// lSgbOrder.WebToolTip = false
			// lSgbOrder.DiscountText = lDiscountTxt

			// 			SgbArr2 = append(SgbArr2, lSgbOrder)
			// 		}
			// 	}
			// }
			// log.Println("SgbArr2", SgbArr2)
		}
	}
	log.Println("GetReportModule (-)")
	return lIpoArr, SgbArr2, nil
}

// this method is added to check the given symbol is valid or not
//by Pavithra
func CheckSymbolValid(pReportRec localdetails.ReportReqStruct) (string, error) {
	log.Println("CheckSymbolValid (+)")

	lStatus := common.ErrorCode
	var lCoreString string
	var lFlag string

	// to Establish A database conncetion, call the LocalDbconnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGRCS01", lErr1)
		return lStatus, lErr1
	} else {
		defer lDb.Close()

		if pReportRec.Module == "Ipo" {
			lCoreString = `select (case when count(1) > 0 then 'Y' else 'N' end) flag 
							from a_ipo_master 
							where Symbol = ?`
		} else {
			lCoreString = `select (case when count(1) > 0 then 'Y' else 'N' end) flag 
							from a_sgb_master 
							where Symbol = ?`
		}

		lRows, lErr2 := lDb.Query(lCoreString, pReportRec.Symbol)
		if lErr2 != nil {
			log.Println("LGRCS02", lErr2)
			return lStatus, lErr2
		} else {
			//Reading the Records
			for lRows.Next() {
				lErr3 := lRows.Scan(&lFlag)
				if lErr3 != nil {
					log.Println("LGRCS03", lErr3)
					return lStatus, lErr3
				} else {
					if lFlag == "Y" {
						lStatus = common.SuccessCode
					} else {
						lStatus = common.ErrorCode
					}
				}
			}
		}
	}
	log.Println("CheckSymbolValid (-)")
	return lStatus, nil
}
