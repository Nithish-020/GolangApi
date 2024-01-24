package ncblocaldetails

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

type NcbOrderHistoryStruct struct {
	Id                 int    `json:"id"`
	Symbol             string `json:"symbol"`
	Series             string `json:"series"`
	ApplicationNo      string `json:"applicationNo"`
	OrderNo            string `json:"orderNo"`
	RespOrderNo        string `json:"respOrderNo"`
	OrderDate          string `json:"orderDate"`
	Isin               string `json:"isin"`
	DateRange          string `json:"dateRange"`
	StartDateWithTime  string `json:"startDateWithTime"`
	EndDateWithTime    string `json:"endDateWithTime"`
	RequestedUnit      int    `json:"requestedUnit"`
	RequestedUnitPrice int    `json:"requestedUnitPrice"`
	RequestedAmount    int    `json:"requestedAmount"`
	AppliedUnit        int    `json:"appliedUnit"`
	AppliedUnitPrice   int    `json:"appliedUnitPrice"`
	AppliedAmount      int    `json:"appliedAmount"`
	AllotedUnit        int    `json:"allotedUnit"`
	AllotedUnitPrice   int    `json:"allotedUnitPrice"`
	AllotedAmount      int    `json:"allotedAmount"`
	OrderStatus        string `json:"orderStatus"`
	DiscountAmt        int    `json:"discountAmt"`
	DiscountText       string `json:"discountText"`
	StatusColor        string `json:"statusColor"`
	SIValue            int    `json:"SIvalue"`
	SIText             string `json:"SItext"`
	RBIStatus          string `json:"rbiStatus"`
	DPStatus           string `json:"dpStatus"`
	Exchange           string `json:"exchange"`
	ClientId           string `json:"clientId"`
}

// Response Structure for GetNcbMaster API
type NcbOrderHistoryResp struct {
	GSecOrderHistoryArr    []NcbOrderHistoryStruct `json:"gSecOrderHistoryArr"`
	TBillOrderHistoryArr   []NcbOrderHistoryStruct `json:"tBillOrderHistoryArr"`
	SdlOrderHistoryArr     []NcbOrderHistoryStruct `json:"sdlOrderHistoryArr"`
	OrderCount             int                     `json:"orderCount"`
	GoiHistoryFound        string                  `json:"goihistoryFound"`
	GoiHistoryNoDataText   string                  `json:"goihistorynoDataText"`
	TbillHistoryFound      string                  `json:"tbillhistoryFound"`
	TbillHistoryNoDataText string                  `json:"tbillhistorynoDataText"`
	SdlHistoryFound        string                  `json:"sdlhistoryFound"`
	SdlHistoryNoDataText   string                  `json:"sdlhistorynoDataText"`
	Status                 string                  `json:"status"`
	ErrMsg                 string                  `json:"errMsg"`
}

type ReportReqStruct struct {
	Module   string `json:"module"`
	ClientId string `json:"clientId"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
	Symbol   string `json:"symbol"`
}

/*
Pupose:This Function is used to Get the NcbDetailStruct in our database table ....
Parameters:

not Applicable

Response:

*On Sucess
=========


	{
		"NcbDetailStruct": [
			{
				"id": 18,
				"symbol": "GJ20392502",
				"startDate": "2023-12-13",
				"endDate": "2023-12-29",
				"priceRange": "100 - 20000000",
				"cutOffPrice": 100,
				"minBidQuantity": 10,
				"applicationStatus": "Pending",
				"upiStatus": "Accepted by Investor"
			},

		],
		"status": "S",
		"errMsg": ""
	}

!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYA DHARSHANI M
Date: 11 OCT 2023
*/

func GetNcbOrderHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNcbOrderHistory(+)", r.Method)

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
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// create the instance for IpoStruct
		var lRespRec NcbOrderHistoryResp

		lRespRec.Status = common.SuccessCode
		// This struct is used for the purpose of report
		var lreportStruct ReportReqStruct

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
		if lErr1 != nil {
			log.Println("NGOH01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGOH01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGOH01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("NGOH02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lNcbOrderHistoryResp, lErr2 := GetNcbOrderHistorydetail(lClientId, lBrokerId, lreportStruct, "getNcbOrderHistory")

		if lErr2 != nil {
			log.Println("NGOH03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGOH03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGOH03", "Error Occur in getting Datas.."))
			return
		} else {

			lRespRec = lNcbOrderHistoryResp

			log.Println("Historyvalue --> lRespRec", lRespRec)
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NGOH04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("NGOH04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetNcbOrderHistory(-)", r.Method)
	}
}

func GetNcbOrderHistorydetail(pClientId string, pBrokerId int, lreportStruct ReportReqStruct, pMethod string) (NcbOrderHistoryResp, error) {
	log.Println("GetNcbOrderHistorydetail (+)")
	log.Println("pMethgod", pMethod)
	var lNcbResp NcbOrderHistoryResp
	// var lNcbMaster NcbOrderHistoryStruct

	var lMethod = "getNcbOrderHistory"
	lNcbResp.Status = common.SuccessCode

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGNOHD01", lErr1)
		lNcbResp.Status = common.ErrorCode
		lNcbResp.ErrMsg = lErr1.Error()

		return lNcbResp, lErr1
	} else {
		defer lDb.Close()

		lGsNcbMasterArr, lTBNcbMasterArr, lSgNcbMasterArr, lErr2 := NcbMasterRec(lMethod, pBrokerId, pClientId)
		if lErr2 != nil {
			log.Println("LGNOHD02", lErr1)
			lNcbResp.Status = common.ErrorCode
			lNcbResp.ErrMsg = lErr2.Error()

			return lNcbResp, lErr2

		} else {
			log.Println("History master", lGsNcbMasterArr)

			// if lNcbResp.GSecOrderHistoryArr != nil || lNcbResp.TBillOrderHistoryArr != nil || lNcbResp.SdlOrderHistoryArr != nil {
			lNcbRespRec, lErr3 := NcbOrderHistoryRec(pBrokerId, pClientId, lGsNcbMasterArr, lTBNcbMasterArr, lSgNcbMasterArr, lreportStruct, pMethod)
			if lErr3 != nil {
				log.Println("LGNOHD03", lErr3)
				lNcbResp.Status = common.ErrorCode
				lNcbResp.ErrMsg = lErr3.Error()
				return lNcbResp, lErr3

			} else {
				log.Println("lNcbRespRec", lNcbRespRec)
				lNcbResp = lNcbRespRec
			}

			// lTBNcbRespRec, lErr4 := NcbOrderHistoryRec(pBrokerId, pClientId, lTBNcbMasterArr, lreportStruct, pMethod)
			// if lErr4 != nil {
			// 	log.Println("LGNOHD04", lErr4)
			// 	lNcbResp.Status = common.ErrorCode
			// 	lNcbResp.ErrMsg = lErr4.Error()
			// 	return lNcbResp, lErr4

			// } else {
			// 	log.Println("lTBNcbRespRec", lTBNcbRespRec)
			// 	lNcbResp = lTBNcbRespRec
			// }

			// lSDlNcbRespRec, lErr5 := NcbOrderHistoryRec(pBrokerId, pClientId, lSgNcbMasterArr, lreportStruct, pMethod)
			// if lErr5 != nil {
			// 	log.Println("LGNOHD05", lErr5)
			// 	lNcbResp.Status = common.ErrorCode
			// 	lNcbResp.ErrMsg = lErr5.Error()
			// 	return lNcbResp, lErr5

			// } else {
			// 	log.Println("lSDlNcbRespRec", lSDlNcbRespRec)
			// 	lNcbResp = lSDlNcbRespRec
			// }
			// } else {
			// 	log.Println("lNcbResp", lNcbResp)

			// 	// lGsNcbMasterArr =  lGsecNcbRespRec
			// }

		}

	}

	log.Println("GetNcbOrderHistorydetail (-)")
	return lNcbResp, nil
}

// version 2
func NcbOrderHistoryRec(pBrokerId int, pClientId string, pGsecmasterArr []ActiveNcbStruct, pTbillmasterArr []ActiveNcbStruct, pSdlmasterArr []ActiveNcbStruct, lReportStruct ReportReqStruct, pMethod string) (NcbOrderHistoryResp, error) {

	log.Println("NcbOrderHistoryRec (+)")
	var lNcbResp NcbOrderHistoryResp
	var lGSNcbRespArr []NcbOrderHistoryStruct
	var lTBNcbRespArr []NcbOrderHistoryStruct
	var lSGNcbRespArr []NcbOrderHistoryStruct
	var lNcbOrderHistoryRec NcbOrderHistoryStruct
	var lOrderCount int
	var lReqUnits, lAppliedUnits, lAllocatedUnits, lStatus string
	var lReqUnitPrice, lReqAmount, lAppliedAmount, lAppliedUnitPrice, lAllocatedUnitPrice, lAllocatedAmount float64

	lNcbResp.Status = common.SuccessCode

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NOHR01", lErr1)
		lNcbResp.Status = common.ErrorCode
		return lNcbResp, lErr1
	} else {
		defer lDb.Close()
		lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
		lGSHistoryNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_GSHistoryNoDataText"])
		lTBHistoryNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_TBHistoryNoDataText"])
		lSdlHistoryNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_SGHistoryNoDataText"])
		lBrokerId := strconv.Itoa(pBrokerId)

		for _, lMasterRec := range pGsecmasterArr {

			lMasterid := strconv.Itoa(lMasterRec.Id)
			lAdditionalCoreString := ""

			lCorestring := `select  h.Id Id,h.Symbol,h.Series, h.exchange,h.ClientId, d.ReqOrderNo,
							d.RespOrderNo ,nvl(h.ReqapplicationNo,'') applicationNo, date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p')OrderDate,
							nvl(d.ReqUnit,0) as RequestedUnits,
							nvl(d.Reqprice,0) as RequestedUnitPrice,
							nvl(d.ReqAmount,0) as RequestedAmount, 
							nvl(d.RespUnit,0) as AppliedUnits,
							nvl(d.Respprice,0) as AppliedUnitPrice,
							nvl(d.RespAmount, 0) as AppliedAmount,
							nvl(d.AllotedUnit ,0) as AllocatedUnits,
							nvl(d.AllotedPrice ,0) as AllocatedUnitPrice,
							nvl(d.AllotedAmount, 0) as AllocatedAmount,
							nvl(h.SIvalue, '') SIvalue,
							nvl(h.SItext,'') SItext ,
							nvl(h.RbiRemarks,'') RbiRemarks,
							nvl(h.DpRemarks,'') DpRemarks,
							(case when  h.status ='failed' or h.cancelFlag  = 'Y' and h.status = 'success' then 'R' else 'G' end ) statuscolor,	
							(case when h.status='failed' then 'F' when h.cancelFlag = 'Y' and h.status = 'success' then 'BC' else 'S' end )Status 
					from a_ncb_orderdetails d, a_ncb_orderheader h
					where  d.HeaderId = h.Id 
					and h.brokerId =  '` + lBrokerId + `' 
				   and h.MasterId =  '` + lMasterid + `'`

			if pMethod == "getNcbOrderHistory" {
				lAdditionalCoreString = ` and h.clientId = '` + pClientId + `'
					and h.status is not null
					order by h.Id desc`
			} else if pMethod == "GetReport" {
				lAdditionalCoreString = ` AND ('` + lReportStruct.Symbol + `' = '' OR h.Symbol = '` + lReportStruct.Symbol + `')
					AND ('` + lReportStruct.ClientId + `' = '' OR  h.clientId = '` + lReportStruct.ClientId + `')
					AND ('` + lReportStruct.FromDate + `' = '' OR h.CreatedDate BETWEEN CONCAT('` + lReportStruct.FromDate + `',' 00:00:00.000') AND CONCAT('` + lReportStruct.ToDate + `',' 23:59:59.000'))`

			} else if pMethod == "getDefaultReport" {
				lAdditionalCoreString = ` and d.CreatedDate between concat (date(now()), ' 00:00:00.000')
					 and concat (date(now()), ' 23:59:59.000')`
			}

			lCoreString := lCorestring + lAdditionalCoreString

			lRows, lErr2 := lDb.Query(lCoreString)
			if lErr2 != nil {
				log.Println("NOHR02", lErr2)
				lNcbResp.Status = common.ErrorCode
				return lNcbResp, lErr2
			} else {
				for lRows.Next() {

					lErr3 := lRows.Scan(&lNcbOrderHistoryRec.Id, &lNcbOrderHistoryRec.Symbol, &lNcbOrderHistoryRec.Series, &lNcbOrderHistoryRec.Exchange, &lNcbOrderHistoryRec.ClientId, &lNcbOrderHistoryRec.OrderNo, &lNcbOrderHistoryRec.RespOrderNo, &lNcbOrderHistoryRec.ApplicationNo, &lNcbOrderHistoryRec.OrderDate, &lReqUnits, &lReqUnitPrice, &lReqAmount, &lAppliedUnits, &lAppliedUnitPrice, &lAppliedAmount, &lAllocatedUnits, &lAllocatedUnitPrice, &lAllocatedAmount, &lNcbOrderHistoryRec.SIValue, &lNcbOrderHistoryRec.SIText, &lNcbOrderHistoryRec.RBIStatus, &lNcbOrderHistoryRec.DPStatus, &lNcbOrderHistoryRec.StatusColor, &lStatus)

					if lErr3 != nil {
						log.Println("NOHR03", lErr3)
						lNcbResp.Status = common.ErrorCode
						lNcbResp.ErrMsg = lErr3.Error()
						return lNcbResp, lErr3
					} else {
						if lStatus == "S" {
							lOrderCount++
						}
						log.Println("units122", lReqUnits, lReqUnitPrice, lReqAmount)

						lNcbOrderHistoryRec.Symbol = lMasterRec.Symbol
						lNcbOrderHistoryRec.Series = lMasterRec.Series
						lNcbOrderHistoryRec.Isin = lMasterRec.Isin
						lNcbOrderHistoryRec.DiscountAmt = lMasterRec.DiscountAmt
						lNcbOrderHistoryRec.DiscountText = lMasterRec.DiscountText
						lNcbOrderHistoryRec.DateRange = lMasterRec.DateRange
						lNcbOrderHistoryRec.StartDateWithTime = lMasterRec.StartDateWithTime
						lNcbOrderHistoryRec.EndDateWithTime = lMasterRec.EndDateWithTime
						lNcbOrderHistoryRec.RequestedUnit, _ = strconv.Atoi(lReqUnits)
						lNcbOrderHistoryRec.RequestedUnitPrice = int(math.Round(lReqUnitPrice))
						lNcbOrderHistoryRec.RequestedAmount = int(math.Round(lReqAmount))
						lNcbOrderHistoryRec.AppliedUnit, _ = strconv.Atoi(lAppliedUnits)
						lNcbOrderHistoryRec.AppliedUnitPrice = int(math.Round(lAppliedUnitPrice))
						lNcbOrderHistoryRec.AppliedAmount = int(math.Round(lAppliedAmount))
						lNcbOrderHistoryRec.AllotedUnit, _ = strconv.Atoi(lAllocatedUnits)
						lNcbOrderHistoryRec.AllotedUnitPrice = int(math.Round(lAllocatedUnitPrice))
						lNcbOrderHistoryRec.AllotedAmount = int(math.Round(lAllocatedAmount))

						log.Println("Unit", lNcbOrderHistoryRec.RequestedUnit, lNcbOrderHistoryRec.RequestedUnitPrice, lNcbOrderHistoryRec.RequestedAmount)
						// Append Upi End Point in lRespRec.NcbHistoryArr array
						OrderStatus, lErr4 := NcbDiscription(lStatus)
						if lErr4 != nil {
							log.Println("LGSHD04", lErr4)
							lNcbResp.Status = common.ErrorCode
							return lNcbResp, lErr4
						} else {
							lNcbOrderHistoryRec.OrderStatus = OrderStatus
						}
					}

					if lNcbOrderHistoryRec.Series == "GS" {
						lGSNcbRespArr = append(lGSNcbRespArr, lNcbOrderHistoryRec)

						log.Println("Gs history121212", lNcbOrderHistoryRec.Series, lGSNcbRespArr)
					}
					// else if lNcbOrderHistoryRec.Series == "TB" {
					// 	lTBNcbRespArr = append(lTBNcbRespArr, lNcbOrderHistoryRec)
					// } else {
					// 	lSGNcbRespArr = append(lSGNcbRespArr, lNcbOrderHistoryRec)
					// }

				}
			}

		}

		for _, lMasterRec := range pTbillmasterArr {

			lMasterid := strconv.Itoa(lMasterRec.Id)
			lAdditionalCoreString := ""

			lCorestring := `select  h.Id Id,h.Symbol,h.Series, h.exchange,h.ClientId, d.ReqOrderNo,
							d.RespOrderNo ,nvl(h.ReqapplicationNo,'') applicationNo, date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p')OrderDate,
							nvl(d.ReqUnit,0) as RequestedUnits,
							nvl(d.Reqprice,0) as RequestedUnitPrice,
							nvl(d.ReqAmount,0) as RequestedAmount, 
							nvl(d.RespUnit,0) as AppliedUnits,
							nvl(d.Respprice,0 )as AppliedUnitPrice,
							nvl(d.RespAmount, 0) as AppliedAmount,
							nvl(d.AllotedUnit ,0) as AllocatedUnits,
							nvl(d.AllotedPrice ,0) as AllocatedUnitPrice,
							nvl(d.AllotedAmount, 0) as AllocatedAmount,
							nvl(h.SIvalue, '') SIvalue,
							nvl(h.SItext,'') SItext ,
							nvl(h.RbiRemarks,'') RbiRemarks,
							nvl(h.DpRemarks,'') DpRemarks,
							(case when  h.status ='failed' or h.cancelFlag  = 'Y' and h.status = 'success' then 'R' else 'G' end ) statuscolor,	
							(case when h.status='failed' then 'F' when h.cancelFlag = 'Y' and h.status = 'success' then 'BC' else 'S' end )Status 
					from a_ncb_orderdetails d, a_ncb_orderheader h
					where  d.HeaderId = h.Id 
					and h.brokerId =  '` + lBrokerId + `' 
				   and h.MasterId =  '` + lMasterid + `'`

			if pMethod == "getNcbOrderHistory" {
				lAdditionalCoreString = ` and h.clientId = '` + pClientId + `'
					and h.status is not null
					order by h.Id desc`
			} else if pMethod == "GetReport" {
				lAdditionalCoreString = ` AND ('` + lReportStruct.Symbol + `' = '' OR h.Symbol = '` + lReportStruct.Symbol + `')
					AND ('` + lReportStruct.ClientId + `' = '' OR  h.clientId = '` + lReportStruct.ClientId + `')
					AND ('` + lReportStruct.FromDate + `' = '' OR h.CreatedDate BETWEEN CONCAT('` + lReportStruct.FromDate + `',' 00:00:00.000') AND CONCAT('` + lReportStruct.ToDate + `',' 23:59:59.000'))`

			} else if pMethod == "getDefaultReport" {
				lAdditionalCoreString = ` and d.CreatedDate between concat (date(now()), ' 00:00:00.000')
					 and concat (date(now()), ' 23:59:59.000')`
			}

			lCoreString := lCorestring + lAdditionalCoreString

			lRows, lErr5 := lDb.Query(lCoreString)
			if lErr5 != nil {
				log.Println("NOHR05", lErr5)
				lNcbResp.Status = common.ErrorCode
				return lNcbResp, lErr5
			} else {
				for lRows.Next() {

					lErr6 := lRows.Scan(&lNcbOrderHistoryRec.Id, &lNcbOrderHistoryRec.Symbol, &lNcbOrderHistoryRec.Series, &lNcbOrderHistoryRec.Exchange, &lNcbOrderHistoryRec.ClientId, &lNcbOrderHistoryRec.OrderNo, &lNcbOrderHistoryRec.RespOrderNo, &lNcbOrderHistoryRec.ApplicationNo, &lNcbOrderHistoryRec.OrderDate, &lReqUnits, &lReqUnitPrice, &lReqAmount, &lAppliedUnits, &lAppliedUnitPrice, &lAppliedAmount, &lAllocatedUnits, &lAllocatedUnitPrice, &lAllocatedAmount, &lNcbOrderHistoryRec.SIValue, &lNcbOrderHistoryRec.SIText, &lNcbOrderHistoryRec.RBIStatus, &lNcbOrderHistoryRec.DPStatus, &lNcbOrderHistoryRec.StatusColor, &lStatus)

					if lErr6 != nil {
						log.Println("NOHR06", lErr6)
						lNcbResp.Status = common.ErrorCode
						lNcbResp.ErrMsg = lErr6.Error()
						return lNcbResp, lErr6
					} else {
						if lStatus == "S" {
							lOrderCount++
						}
						log.Println("units122", lReqUnits, lReqUnitPrice, lReqAmount)

						lNcbOrderHistoryRec.Symbol = lMasterRec.Symbol
						lNcbOrderHistoryRec.Series = lMasterRec.Series
						lNcbOrderHistoryRec.Isin = lMasterRec.Isin
						lNcbOrderHistoryRec.DiscountAmt = lMasterRec.DiscountAmt
						lNcbOrderHistoryRec.DiscountText = lMasterRec.DiscountText
						lNcbOrderHistoryRec.DateRange = lMasterRec.DateRange
						lNcbOrderHistoryRec.StartDateWithTime = lMasterRec.StartDateWithTime
						lNcbOrderHistoryRec.EndDateWithTime = lMasterRec.EndDateWithTime
						lNcbOrderHistoryRec.RequestedUnit, _ = strconv.Atoi(lReqUnits)
						lNcbOrderHistoryRec.RequestedUnitPrice = int(math.Round(lReqUnitPrice))
						lNcbOrderHistoryRec.RequestedAmount = int(math.Round(lReqAmount))
						lNcbOrderHistoryRec.AppliedUnit, _ = strconv.Atoi(lAppliedUnits)
						lNcbOrderHistoryRec.AppliedUnitPrice = int(math.Round(lAppliedUnitPrice))
						lNcbOrderHistoryRec.AppliedAmount = int(math.Round(lAppliedAmount))
						lNcbOrderHistoryRec.AllotedUnit, _ = strconv.Atoi(lAllocatedUnits)
						lNcbOrderHistoryRec.AllotedUnitPrice = int(math.Round(lAllocatedUnitPrice))
						lNcbOrderHistoryRec.AllotedAmount = int(math.Round(lAllocatedAmount))

						// Append Upi End Point in lRespRec.NcbHistoryArr array
						OrderStatus, lErr7 := NcbDiscription(lStatus)
						if lErr7 != nil {
							log.Println("LGSHD07", lErr7)
							lNcbResp.Status = common.ErrorCode
							return lNcbResp, lErr7
						} else {
							lNcbOrderHistoryRec.OrderStatus = OrderStatus
						}
					}

					// if lNcbOrderHistoryRec.Series == "GS" {
					// 	lGSNcbRespArr = append(lGSNcbRespArr, lNcbOrderHistoryRec)

					// 	log.Println("Gs history121212", lNcbOrderHistoryRec.Series, lGSNcbRespArr)
					// } else
					if lNcbOrderHistoryRec.Series == "TB" {
						lTBNcbRespArr = append(lTBNcbRespArr, lNcbOrderHistoryRec)
					}
					// else {
					// 	lSGNcbRespArr = append(lSGNcbRespArr, lNcbOrderHistoryRec)
					// }

				}
			}

		}

		for _, lMasterRec := range pSdlmasterArr {

			lMasterid := strconv.Itoa(lMasterRec.Id)
			lAdditionalCoreString := ""

			lCorestring := `select  h.Id Id,h.Symbol,h.Series, h.exchange,h.ClientId, d.ReqOrderNo,
							d.RespOrderNo ,nvl(h.ReqapplicationNo,'') applicationNo, date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p')OrderDate,
							nvl(d.ReqUnit,0) as RequestedUnits,
							nvl(d.Reqprice,0) as RequestedUnitPrice,
							nvl(d.ReqAmount,0) as RequestedAmount, 
							nvl(d.RespUnit,0) as AppliedUnits,
							nvl(d.Respprice,0 )as AppliedUnitPrice,
							nvl(d.RespAmount, 0) as AppliedAmount,
							nvl(d.AllotedUnit ,0) as AllocatedUnits,
							nvl(d.AllotedPrice ,0) as AllocatedUnitPrice,
							nvl(d.AllotedAmount, 0) as AllocatedAmount,
							nvl(h.SIvalue, '') SIvalue,
							nvl(h.SItext,'') SItext ,
							nvl(h.RbiRemarks,'') RbiRemarks,
							nvl(h.DpRemarks,'') DpRemarks,
							(case when  h.status ='failed' or h.cancelFlag  = 'Y' and h.status = 'success' then 'R' else 'G' end ) statuscolor,	
							(case when h.status='failed' then 'F' when h.cancelFlag = 'Y' and h.status = 'success' then 'BC' else 'S' end )Status 
					from a_ncb_orderdetails d, a_ncb_orderheader h
					where  d.HeaderId = h.Id 
					and h.brokerId =  '` + lBrokerId + `' 
				   and h.MasterId =  '` + lMasterid + `'`

			if pMethod == "getNcbOrderHistory" {
				lAdditionalCoreString = ` and h.clientId = '` + pClientId + `'
					and h.status is not null
					order by h.Id desc`
			} else if pMethod == "GetReport" {
				lAdditionalCoreString = ` AND ('` + lReportStruct.Symbol + `' = '' OR h.Symbol = '` + lReportStruct.Symbol + `')
					AND ('` + lReportStruct.ClientId + `' = '' OR  h.clientId = '` + lReportStruct.ClientId + `')
					AND ('` + lReportStruct.FromDate + `' = '' OR h.CreatedDate BETWEEN CONCAT('` + lReportStruct.FromDate + `',' 00:00:00.000') AND CONCAT('` + lReportStruct.ToDate + `',' 23:59:59.000'))`

			} else if pMethod == "getDefaultReport" {
				lAdditionalCoreString = ` and d.CreatedDate between concat (date(now()), ' 00:00:00.000')
					 and concat (date(now()), ' 23:59:59.000')`
			}

			lCoreString := lCorestring + lAdditionalCoreString

			lRows, lErr8 := lDb.Query(lCoreString)
			if lErr8 != nil {
				log.Println("NOHR08", lErr8)
				lNcbResp.Status = common.ErrorCode
				return lNcbResp, lErr8
			} else {
				for lRows.Next() {

					lErr9 := lRows.Scan(&lNcbOrderHistoryRec.Id, &lNcbOrderHistoryRec.Symbol, &lNcbOrderHistoryRec.Series, &lNcbOrderHistoryRec.Exchange, &lNcbOrderHistoryRec.ClientId, &lNcbOrderHistoryRec.OrderNo, &lNcbOrderHistoryRec.RespOrderNo, &lNcbOrderHistoryRec.ApplicationNo, &lNcbOrderHistoryRec.OrderDate, &lReqUnits, &lReqUnitPrice, &lReqAmount, &lAppliedUnits, &lAppliedUnitPrice, &lAppliedAmount, &lAllocatedUnits, &lAllocatedUnitPrice, &lAllocatedAmount, &lNcbOrderHistoryRec.SIValue, &lNcbOrderHistoryRec.SIText, &lNcbOrderHistoryRec.RBIStatus, &lNcbOrderHistoryRec.DPStatus, &lNcbOrderHistoryRec.StatusColor, &lStatus)

					if lErr9 != nil {
						log.Println("NOHR09", lErr9)
						lNcbResp.Status = common.ErrorCode
						lNcbResp.ErrMsg = lErr9.Error()
						return lNcbResp, lErr9
					} else {
						if lStatus == "S" {
							lOrderCount++
						}
						log.Println("units122", lReqUnits, lReqUnitPrice, lReqAmount)

						lNcbOrderHistoryRec.Symbol = lMasterRec.Symbol
						lNcbOrderHistoryRec.Series = lMasterRec.Series
						lNcbOrderHistoryRec.Isin = lMasterRec.Isin
						lNcbOrderHistoryRec.DiscountAmt = lMasterRec.DiscountAmt
						lNcbOrderHistoryRec.DiscountText = lMasterRec.DiscountText
						lNcbOrderHistoryRec.DateRange = lMasterRec.DateRange
						lNcbOrderHistoryRec.StartDateWithTime = lMasterRec.StartDateWithTime
						lNcbOrderHistoryRec.EndDateWithTime = lMasterRec.EndDateWithTime
						lNcbOrderHistoryRec.RequestedUnit, _ = strconv.Atoi(lReqUnits)
						lNcbOrderHistoryRec.RequestedUnitPrice = int(math.Round(lReqUnitPrice))
						lNcbOrderHistoryRec.RequestedAmount = int(math.Round(lReqAmount))
						lNcbOrderHistoryRec.AppliedUnit, _ = strconv.Atoi(lAppliedUnits)
						lNcbOrderHistoryRec.AppliedUnitPrice = int(math.Round(lAppliedUnitPrice))
						lNcbOrderHistoryRec.AppliedAmount = int(math.Round(lAppliedAmount))
						lNcbOrderHistoryRec.AllotedUnit, _ = strconv.Atoi(lAllocatedUnits)
						lNcbOrderHistoryRec.AllotedUnitPrice = int(math.Round(lAllocatedUnitPrice))
						lNcbOrderHistoryRec.AllotedAmount = int(math.Round(lAllocatedAmount))

						// Append Upi End Point in lRespRec.NcbHistoryArr array
						OrderStatus, lErr10 := NcbDiscription(lStatus)
						if lErr10 != nil {
							log.Println("LGSHD10", lErr10)
							lNcbResp.Status = common.ErrorCode
							return lNcbResp, lErr10
						} else {
							lNcbOrderHistoryRec.OrderStatus = OrderStatus
						}
					}

					// if lNcbOrderHistoryRec.Series == "GS" {
					// 	lGSNcbRespArr = append(lGSNcbRespArr, lNcbOrderHistoryRec)

					// 	log.Println("Gs history121212", lNcbOrderHistoryRec.Series, lGSNcbRespArr)
					// } else
					// if lNcbOrderHistoryRec.Series == "TB" {
					// 	lTBNcbRespArr = append(lTBNcbRespArr, lNcbOrderHistoryRec)
					// } else {
					if lNcbOrderHistoryRec.Series == "SG" {
						lSGNcbRespArr = append(lSGNcbRespArr, lNcbOrderHistoryRec)
					}

				}
			}
		}

		if lGSNcbRespArr != nil {
			lNcbResp.GSecOrderHistoryArr = lGSNcbRespArr
			lNcbResp.OrderCount = lOrderCount
			lNcbResp.GoiHistoryFound = "Y"
		} else {
			if len(lGSNcbRespArr) == 0 {
				lNcbResp.GoiHistoryFound = "N"
				lNcbResp.GoiHistoryNoDataText = lGSHistoryNoDataTxt
			}
		}

		if lTBNcbRespArr != nil {
			lNcbResp.TBillOrderHistoryArr = lTBNcbRespArr
			lNcbResp.OrderCount = lOrderCount
			lNcbResp.TbillHistoryFound = "Y"
		} else {
			if len(lTBNcbRespArr) == 0 {
				lNcbResp.TbillHistoryFound = "N"
				lNcbResp.TbillHistoryNoDataText = lTBHistoryNoDataTxt
			}
		}

		if lSGNcbRespArr != nil {
			lNcbResp.SdlOrderHistoryArr = lSGNcbRespArr
			lNcbResp.OrderCount = lOrderCount
			lNcbResp.SdlHistoryFound = "Y"
		} else {
			if len(lSGNcbRespArr) == 0 {
				lNcbResp.SdlHistoryFound = "N"
				lNcbResp.SdlHistoryNoDataText = lSdlHistoryNoDataTxt
			}
		}

		log.Println("NcbOrderHistoryRec (-)")
		return lNcbResp, nil
	}
}
