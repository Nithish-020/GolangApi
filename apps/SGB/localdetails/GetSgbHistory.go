package localdetails

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"

	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// version 1
// type SgbHistoryStruct struct {
// 	Id               int    `json:"id"`
// 	Name             string `json:"name"`
// 	ReqOrderNo       string `json:"reqOrderNo"`
// 	OrderNo          string `json:"orderNo"`
// 	OrderDate        string `json:"orderDate"`
// 	Isin             string `json:"isin"`
// 	StartDate        string `json:"startDate"`
// 	EndDate          string `json:"endDate"`
// 	DateTime         string `json:"dateTime"`
// 	Unit             int    `json:"unit"`
// 	Price            int    `json:"price"`
// 	Total            int    `json:"total"`
// 	Flag             string `json:"flag"`
// 	Status           string `json:"status"`
// 	Subscriptionunit int    `json:"subscriptionunit"`
// 	BlockedAmount    int    `json:"blockedAmount"`
// 	ToolTip          bool   `json:"toolTip"`
// }

// Response Structure for GetSgbMaster API
// type SgbHistoryResp struct {
// 	SgbHistoryArr []SgbHistoryStruct `json:"sgbHistoryArr"`
// 	Status        string             `json:"status"`
// 	ErrMsg        string             `json:"errMsg"`
// }

// version 2
type SgbOrderHistoryStruct struct {
	Id                 int    `json:"id"`
	Symbol             string `json:"symbol"`
	Name               string `json:"name"`
	OrderNo            string `json:"orderNo"`
	ExchOrderNo        string `json:"exchOrderNo"`
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
	SIValue            bool   `json:"SIvalue"`
	SIText             string `json:"SItext"`
	RBIStatus          string `json:"rbiStatus"`
	DPStatus           string `json:"dpStatus"`
	Exchange           string `json:"exchange"`
	ClientId           string `json:"clientId"`
}

type SgbOrderHistoryResp struct {
	SgbOrderHistoryArr []SgbOrderHistoryStruct `json:"sgbOrderHistoryArr"`
	OrderCount         int                     `json:"orderCount"`
	HistoryFound       string                  `json:"historyFound"`
	HistoryNoDataText  string                  `json:"historynoDataText"`
	Status             string                  `json:"status"`
	ErrMsg             string                  `json:"errMsg"`
}
type ReportReqStruct struct {
	Module   string `json:"module"`
	ClientId string `json:"clientId"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
	Symbol   string `json:"symbol"`
}

// /*
// Pupose:This Function is used to Get the Active SGbOrderHistory  in our database table ....
// Parameters:

// not Applicable

// Response:
// *On Sucess
// =========

// 	{
// 		"sgbHistoryArr": [
// "sgbOrderHistoryArr": [
// 	{
// 	  "id": 73,
// 	  "name": "SOVEREIGN GOLD BONDS SCHEME 2016-17 (TRANCHE 5)",
// 	  "reqOrderNo": "170322432532287",
// 	  "orderNo": "170322432532287", @click="item.webToolTip =!item.webToolTip"
// 	  "orderDate": "22-Dec-23, 11:15 AM",
// 	  "isin": "IN0020150085",
// 	  "dateRange": "28-Dec-13 - 19-Dec-24",
// 	  "dateRangeWithTime": "28th Dec 2013 09:00AM - 28th Dec 2024 11:00PM",
// 	  "requestedUnit": 2,
// 	  "requestedUnitPrice": 3149,
// 	  "requestedAmount": 6298,
// 	  "appliedUnit": 0,
// 	  "appliedUnitPrice": 0,
// 	  "appliedAmount": 0,
// 	  "allotedUnit": 0,
// 	  "allotedUnitPrice": 0,
// 	  "allotedAmount": 0,
// 	  "orderStatus": "Success",
// 	  "discountAmt": 1,
// 	  "discountText": "Discount ₹",
// 	  "statusColor": "G",
// 	  "siValue": 0,
// 	  "siText": " Place order to the extent of the available balance in my account in case of insufficient balance.",
// 	  "webToolTip": false
// 	},
//
// 		],
//   "orderCount": 1,
//   "historyFound": "Y",
//   "historynoDataText": "",
//   "status": "S",
//   "errMsg": ""
// }

// !On Error
// ========
// 	{
//   "sgbOrderHistoryArr": null,
//   "orderCount": 0,
//   "historyFound": "N",
//   "historynoDataText": "You haven't invested in any SGBs.",
//   "status": "E",
//   "errMsg": "Can't able to get the data from database"
// }
// 	Author: Rajkumar M
// 	Date: 26 DECEMBER 2023
// */
func GetSgbOrderHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSgbOrderHistory (+)", r.Method)
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
		var lRespRec SgbOrderHistoryResp

		lRespRec.Status = common.SuccessCode
		// This struct is used for the purpose of report
		var lreportStruct ReportReqStruct
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
		if lErr1 != nil {
			log.Println("LGSH01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSH01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSH01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGSH02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lSgbOrderHistoryResp, lErr2 := GetSGBOrderHistorydetail(lClientId, lBrokerId, lreportStruct, "GetSgbOrderHistory")
		if lErr2 != nil {
			log.Println("LGSH03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSH03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSH03", "Error Occur in getting Datas.."))
			return
		} else {
			lRespRec = lSgbOrderHistoryResp
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGSH04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGSH04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetSgbOrderHistory (-)", r.Method)
	}
}

// version 1
// func GetSGBHistorydetail(pClientId string, pBrokerId int) ([]SgbHistoryStruct, error) {
// 	log.Println("GetSGBHistorydetail (+)")

// 	var lSgbHistoryRec SgbHistoryStruct
// 	var lSgbHistoryArr []SgbHistoryStruct
// 	var lUnit, lPrice, lRespUnit, lRespRate string
// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("LGSHD01", lErr1)
// 		return lSgbHistoryArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := ` select h.Id Id,
// 						sm.Name Name,
// 						d.ReqOrderNo,
// 						d.RespOrderNo,
// 						date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p'),
// 						sm.Isin isin,
// 						sm.BiddingStartDate startDate,
// 						sm.BiddingEndDate endDate,
// 						CONCAT( case
// 						WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
// 						WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
// 						WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
// 						ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
// 						end,' ',
// 						DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
// 						TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime ,
// 						d.ReqSubscriptionunit ,
// 						nvl(d.RespSubscriptionunit,0),
// 						d.ReqRate ,
// 						nvl(d.RespRate,0),
// 						(case when h.MasterId = sm.Id and h.CancelFlag = 'Y' and h.Status = 'success' then 'Y' else 'N' end) Flag,
// 						lower(h.Status) Status
// 					from a_sgb_master sm,a_sgb_orderheader h ,a_sgb_orderdetails d
// 					where sm.Id = h.MasterId
// 					and d.HeaderId = h.Id
// 					and h.ClientId = ?
// 					and h.brokerId = ?
// 					and h.Status is not null
// 					order by h.Id desc`

// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("LGSHD02", lErr2)
// 			return lSgbHistoryArr, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lSgbHistoryRec.Id, &lSgbHistoryRec.Name, &lSgbHistoryRec.ReqOrderNo, &lSgbHistoryRec.OrderNo, &lSgbHistoryRec.OrderDate, &lSgbHistoryRec.Isin, &lSgbHistoryRec.StartDate, &lSgbHistoryRec.EndDate, &lSgbHistoryRec.DateTime, &lUnit, &lRespUnit, &lPrice, &lRespRate, &lSgbHistoryRec.Flag, &lSgbHistoryRec.Status)
// 				if lErr3 != nil {
// 					log.Println("LGSHD03", lErr3)
// 					return lSgbHistoryArr, lErr3
// 				} else {
// 					lSgbHistoryRec.ToolTip = false
// 					lSgbHistoryRec.Unit, _ = strconv.Atoi(lUnit)
// 					lSgbHistoryRec.Price, _ = strconv.Atoi(lPrice)
// 					lSgbHistoryRec.Subscriptionunit, _ = strconv.Atoi(lRespUnit)
// 					lRespPrice, _ := strconv.Atoi(lRespRate)
// 					lSgbHistoryRec.BlockedAmount = lSgbHistoryRec.Subscriptionunit * lRespPrice
// 					lSgbHistoryRec.Total = lSgbHistoryRec.Unit * lSgbHistoryRec.Price
// 					// Append Upi End Point in lRespRec.SgbHistoryArr array
// 					lSgbHistoryArr = append(lSgbHistoryArr, lSgbHistoryRec)
// 				}
// 			}
// 			// log.Println(lSgbHistoryArr)
// 		}
// 	}
// 	log.Println("GetSGBHistorydetail (-)")
// 	return lSgbHistoryArr, nil
// }

func GetSGBOrderHistorydetail(pClientId string, pBrokerId int, lreportStruct ReportReqStruct, pMethod string) (SgbOrderHistoryResp, error) {
	log.Println("GetSGBOrderHistorydetail (+)")
	// version 1 variables
	var lSgbResp SgbOrderHistoryResp
	// var lOrderCount int
	// var lSgbOrderHistoryRec SgbOrderHistoryStruct
	// var lSgbOrderHistoryArr []SgbOrderHistoryStruct
	// seperately added
	// version 2
	var lMethod = "getSgbOrderHistory"
	lSgbResp.Status = common.SuccessCode

	// version 1 variables
	// var lReqUnits, lReqUnitPrice, lReqAmount, lAppliedUnits, lAppliedAmount, lAppliedUnitPrice, lAllocatedUnits, lStatus, lAllocatedUnitPrice, lAllocatedAmount string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGSHD01", lErr1)
		lSgbResp.Status = common.ErrorCode
		lSgbResp.ErrMsg = lErr1.Error()

		return lSgbResp, lErr1
	} else {
		defer lDb.Close()
		// version 2
		lSgbMaster, lErr2 := SgbMasterRec(lMethod, pBrokerId, pClientId)
		if lErr2 != nil {
			log.Println("LGSHD02", lErr1)
			lSgbResp.Status = common.ErrorCode
			lSgbResp.ErrMsg = lErr2.Error()

			return lSgbResp, lErr2

		} else {
			// fmt.Println(pBrokerId, "clentid", lreportStruct, "lreportStruct")

			lSgbRespRec, lErr3 := SgbOrderHistoryRec(pBrokerId, pClientId, lSgbMaster, lreportStruct, pMethod)
			if lErr3 != nil {
				log.Println("LGSHD03", lErr3)
				lSgbResp.Status = common.ErrorCode
				lSgbResp.ErrMsg = lErr3.Error()
				return lSgbResp, lErr3

			} else {
				lSgbResp = lSgbRespRec

				// version 1
				// 	lConfigFile := common.ReadTomlConfig("toml/debug.toml")
				// 	lHistoryNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_HistoryNoDataText"])
				// 	lDiscountTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_DiscountText"])
				// 	lCoreString := ` select h.Id Id,
				// 				sm.Name Name,
				// 				sm.Isin isin,
				// 				d.RespOrderNo,
				// 				d.ReqOrderNo ,
				// 				date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p')OrderDate,
				// 				concat(date_format(sm.BiddingStartDate, '%d-%b-%y'),' - ',date_format(sm.BiddingEndDate , '%d-%b-%y'))DateRange,						CONCAT( case
				// 				WHEN DAY(sm.BiddingStartDate) % 10 = 1 AND DAY(sm.BiddingStartDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingStartDate), 'st')
				// 				WHEN DAY(sm.BiddingStartDate) % 10 = 2 AND DAY(sm.BiddingStartDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingStartDate), 'nd')
				// 				WHEN DAY(sm.BiddingStartDate) % 10 = 3 AND DAY(sm.BiddingStartDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingStartDate), 'rd')
				// 				ELSE CONCAT(DAY(sm.BiddingStartDate), 'th')
				// 				end,' ',
				// 				DATE_FORMAT(sm.BiddingStartDate , '%b %Y'),
				// 				TIME_FORMAT(sm.DailyStartTime , ' %h:%i%p'),' - ', case
				// 				WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
				// 				WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
				// 				WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
				// 				ELSE CONCAT(DAY(sm.BiddingStartDate), 'th')
				// 				end,' ',
				// 				DATE_FORMAT(sm.BiddingEndDate  , '%b %Y'),
				// 				TIME_FORMAT(sm.DailyEndTime  , ' %h:%i%p')) DateRangeWithTime,
				// 				nvl(d.ReqSubscriptionUnit,0) as RequestedUnits,
				// 				nvl(d.ReqRate,0) as RequestedUnitPrice,
				// 				nvl(nvl(d.ReqSubscriptionUnit,0) *nvl(d.ReqRate,0),0) as RequestedAmount,
				// 				nvl(d.RespSubscriptionunit,0) as AppliedUnits,
				// 				nvl(d.RespRate,0 )as AppliedUnitPrice,
				// 				nvl(nvl(d.RespSubscriptionunit,0)  *nvl(d.RespRate,0) ,0) as AppliedAmount,
				// 				nvl(d.AllotedUnit ,0) as AllocatedUnits,
				// 				nvl(d.AllotedRate ,0 )as AllocatedUnitPrice,
				// 				nvl (nvl(d.AllotedUnit,0)  *nvl(d.AllotedRate,0)  ,0) as AllocatedAmount,
				// 				(case when (nvl(sm.FaceValue,0) <=0 or nvl(sm.MinPrice,0) <= 0) then 0 else nvl(sm.FaceValue-sm.MinPrice,0) end) DiscountAmount,
				// 				h.SIvalue ,
				// 				h.SItext ,
				// 			(case when  h.Status='failed' or h.CancelFlag = 'Y' and h.Status = 'success' then 'R' else 'G' end ) textcolor,
				// 			  nvl((select ld.description from xx_lookup_details ld,xx_lookup_header lh
				// 			where lh.id = ld.headerId and lh.Code = 'SgbAction' and ld.code = (case when h.status='failed' then 'F' when h.CancelFlag = 'Y' and h.Status = 'success' then 'BC' else 'S' end )),'') orderStatus,
				// 			(case when h.status='failed' then 'F' when h.CancelFlag = 'Y' and h.Status = 'success' then 'BC' else 'S' end )Status
				// 			from a_sgb_master sm,a_sgb_orderheader h ,a_sgb_orderdetails d
				// 			where sm.Id = h.MasterId
				// 			and d.HeaderId = h.Id
				// 			and h.ClientId = ?
				// 			and h.brokerId = ?
				// 			and h.Status is not null
				// 			order by h.Id desc`

				// 	lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
				// 	if lErr2 != nil {
				// 		log.Println("LGSHD02", lErr2)
				// 		return lSgbResp, lErr2
				// 	} else {
				// 		//This for loop is used to collect the records from the database and store them in structure
				// 		for lRows.Next() {
				// 			lErr3 := lRows.Scan(&lSgbOrderHistoryRec.Id, &lSgbOrderHistoryRec.Name, &lSgbOrderHistoryRec.Isin, &lSgbOrderHistoryRec.OrderNo, &lSgbOrderHistoryRec.ReqOrderNo, &lSgbOrderHistoryRec.OrderDate, &lSgbOrderHistoryRec.DateRange, &lSgbOrderHistoryRec.DateRangeWithTime, &lReqUnits, &lReqUnitPrice, &lReqAmount, &lAppliedUnits, &lAppliedUnitPrice, &lAppliedAmount, &lAllocatedUnits, &lAllocatedUnitPrice, &lAllocatedAmount, &lSgbOrderHistoryRec.DiscountAmt, &lSgbOrderHistoryRec.SIvalue, &lSgbOrderHistoryRec.SItext, &lSgbOrderHistoryRec.StatusColor, &lSgbOrderHistoryRec.OrderStatus, &lStatus)
				// 			if lErr3 != nil {
				// 				log.Println("LGSHD03", lErr3)
				// 				return lSgbResp, lErr3
				// 			} else {
				// 				if lStatus == "S" {
				// 					lOrderCount++
				// 				}
				// 				lSgbOrderHistoryRec.RequestedUnit, _ = strconv.Atoi(lReqUnits)
				// 				lSgbOrderHistoryRec.RequestedUnitPrice, _ = strconv.Atoi(lReqUnitPrice)
				// 				lSgbOrderHistoryRec.RequestedAmount, _ = strconv.Atoi(lReqAmount)
				// 				lSgbOrderHistoryRec.AppliedUnit, _ = strconv.Atoi(lAppliedUnits)
				// 				lSgbOrderHistoryRec.AppliedUnitPrice, _ = strconv.Atoi(lAppliedUnitPrice)
				// 				lSgbOrderHistoryRec.AppliedAmount, _ = strconv.Atoi(lAppliedAmount)
				// 				lSgbOrderHistoryRec.AllotedUnit, _ = strconv.Atoi(lAllocatedUnits)
				// 				lSgbOrderHistoryRec.AllotedUnitPrice, _ = strconv.Atoi(lAllocatedUnitPrice)
				// 				lSgbOrderHistoryRec.AllotedAmount, _ = strconv.Atoi(lAllocatedAmount)
				// 				lSgbOrderHistoryRec.DiscountText = lDiscountTxt
				// 				// Append Upi End Point in lRespRec.SgbHistoryArr array
				// 				lSgbOrderHistoryArr = append(lSgbOrderHistoryArr, lSgbOrderHistoryRec)

				// 			}
				// 		}

				// 		if lSgbOrderHistoryArr != nil {
				// 			lSgbResp.SgbOrderHistoryArr = lSgbOrderHistoryArr
				// 			lSgbResp.OrderCount = lOrderCount
				// 			lSgbResp.Status = common.SuccessCode
				// 			lSgbResp.HistoryFound = "Y"
				// 		} else {
				// 			if len(lSgbOrderHistoryArr) == 0 {
				// 				lSgbResp.HistoryFound = "N"
				// 				lSgbResp.HistoryNoDataText = lHistoryNoDataTxt
				// 			}
				// 		}
				// 		// log.Println(lSgbHistoryArr)
				// 	}
				// }
			}
		}
	}
	log.Println("GetSGBOrderHistorydetail (-)")
	return lSgbResp, nil
}

// version 2
func SgbOrderHistoryRec(pBrokerId int, pClientId string, lMasterArr []ActiveSgbStruct, lReportStruct ReportReqStruct, pMethod string) (SgbOrderHistoryResp, error) {
	log.Println("SgbOrderHistoryRec (+)")
	var lSgbResp SgbOrderHistoryResp
	var lSgbRespArr []SgbOrderHistoryStruct
	var lSgbOrderHistoryRec SgbOrderHistoryStruct
	var lOrderCount int
	var lReqUnits, lReqUnitPrice, lReqAmount, lAppliedUnits, lAppliedAmount, lAppliedUnitPrice, lAllocatedUnits, lStatus, lAllocatedUnitPrice, lAllocatedAmount string

	lSgbResp.Status = common.SuccessCode

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GSGHD01", lErr1)
		lSgbResp.Status = common.ErrorCode
		return lSgbResp, lErr1
	} else {
		defer lDb.Close()
		lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
		lHistoryNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_HistoryNoDataText"])
		lBrokerId := strconv.Itoa(pBrokerId)
		for _, lMasterRec := range lMasterArr {
			lMasterid := strconv.Itoa(lMasterRec.Id)
			lAdditionalCoreString := ""

			lCoreString := `select  h.Id Id,h.ScripId, h.exchange,h.ClientId, d.ReqOrderNo,
			d.RespOrderNo, date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p')OrderDate,nvl(d.ReqSubscriptionUnit,0) as RequestedUnits,
			nvl(d.ReqRate,0) as RequestedUnitPrice,
			nvl(nvl(d.ReqSubscriptionUnit,0) *nvl(d.ReqRate,0),0) as RequestedAmount,
			nvl(d.RespSubscriptionunit,0) as AppliedUnits,
			nvl(d.RespRate,0 )as AppliedUnitPrice,
			nvl(nvl(d.RespSubscriptionunit,0)  *nvl(d.RespRate,0) ,0) as AppliedAmount,	
			nvl(d.AllotedUnit ,0) as AllocatedUnits,
			nvl(d.AllotedRate ,0 )as AllocatedUnitPrice,
			nvl (nvl(d.AllotedUnit,0)  *nvl(d.AllotedRate,0)  ,0) as AllocatedAmount,
			h.SIvalue ,
			h.SItext ,
			nvl(h.RbiRemarks,'-'),
			nvl(h.DpRemarks,'-'),
			(case when  h.Status='failed' or h.CancelFlag = 'Y' and h.Status = 'success' then 'R' else 'G' end ) statuscolor,	(case when h.status='failed' then 'F' when h.CancelFlag = 'Y' and h.Status = 'success' then 'BC' else 'S' end )Status from a_sgb_orderheader h ,a_sgb_orderdetails d
			where  d.HeaderId = h.Id 
			and h.brokerId =  '` + lBrokerId + `' 
			and h.MasterId =  '` + lMasterid + `'`
			// and h.ClientId = ?
			// and h.MasterId = ?
			// and h.Status is not null
			// order by h.Id desc
			if pMethod == "GetSgbOrderHistory" {
				lAdditionalCoreString = ` and h.ClientId = '` + pClientId + `'
				and h.Status is not null
				order by h.Id desc`
			} else if pMethod == "GetReport" {
				lAdditionalCoreString = ` AND ('` + lReportStruct.Symbol + `' = '' OR h.ScripId = '` + lReportStruct.Symbol + `')
				AND ('` + lReportStruct.ClientId + `' = '' OR h.ClientId = '` + lReportStruct.ClientId + `')
				AND ('` + lReportStruct.FromDate + `' = '' OR h.CreatedDate BETWEEN CONCAT('` + lReportStruct.FromDate + `',' 00:00:00.000') AND CONCAT('` + lReportStruct.ToDate + `',' 23:59:59.000'))`

				// ` AND ('` + lReportStruct.Symbol + `' = '' OR h.ScripId = '` + lReportStruct.Symbol + `')
				// AND '` + lReportStruct.ClientId + `' = '' OR h.ClientId = '` + lReportStruct.ClientId + `')
				// AND ('` + lReportStruct.FromDate + `' = '' OR h.CreatedDate BETWEEN CONCAT('` + lReportStruct.FromDate + `',' 00:00:00.000') AND CONCAT('` + lReportStruct.ToDate + `',' 23:59:59.000'))`
			} else if pMethod == "GetDefault" {
				lAdditionalCoreString = ` and d.CreatedDate between concat (date(now()), ' 00:00:00.000')
				 and concat (date(now()), ' 23:59:59.000')`
			}

			lCoreString = lCoreString + lAdditionalCoreString

			lRows, lErr2 := lDb.Query(lCoreString)
			if lErr2 != nil {
				log.Println("LGSHD02", lErr2)
				lSgbResp.Status = common.ErrorCode
				return lSgbResp, lErr2
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows.Next() {
					lErr3 := lRows.Scan(&lSgbOrderHistoryRec.Id, &lSgbOrderHistoryRec.Symbol, &lSgbOrderHistoryRec.Exchange, &lSgbOrderHistoryRec.ClientId, &lSgbOrderHistoryRec.OrderNo, &lSgbOrderHistoryRec.ExchOrderNo, &lSgbOrderHistoryRec.OrderDate, &lReqUnits, &lReqUnitPrice, &lReqAmount, &lAppliedUnits, &lAppliedUnitPrice, &lAppliedAmount, &lAllocatedUnits, &lAllocatedUnitPrice, &lAllocatedAmount, &lSgbOrderHistoryRec.SIValue, &lSgbOrderHistoryRec.SIText, &lSgbOrderHistoryRec.RBIStatus, &lSgbOrderHistoryRec.DPStatus, &lSgbOrderHistoryRec.StatusColor, &lStatus)
					if lErr3 != nil {
						log.Println("LGSHD03", lErr3)
						lSgbResp.Status = common.ErrorCode
						lSgbResp.ErrMsg = lErr3.Error()
						return lSgbResp, lErr3
					} else {
						if lStatus == "S" {
							lOrderCount++
						}
						lSgbOrderHistoryRec.Name = lMasterRec.Name
						lSgbOrderHistoryRec.Isin = lMasterRec.Isin
						lSgbOrderHistoryRec.DiscountAmt = lMasterRec.DiscountAmt
						lSgbOrderHistoryRec.DiscountText = lMasterRec.DiscountText
						lSgbOrderHistoryRec.DateRange = lMasterRec.DateRange
						lSgbOrderHistoryRec.StartDateWithTime = lMasterRec.StartDateWithTime
						lSgbOrderHistoryRec.EndDateWithTime = lMasterRec.EndDateWithTime
						lSgbOrderHistoryRec.RequestedUnit, _ = strconv.Atoi(lReqUnits)
						lSgbOrderHistoryRec.RequestedUnitPrice, _ = strconv.Atoi(lReqUnitPrice)
						lSgbOrderHistoryRec.RequestedAmount, _ = strconv.Atoi(lReqAmount)
						lSgbOrderHistoryRec.AppliedUnit, _ = strconv.Atoi(lAppliedUnits)
						lSgbOrderHistoryRec.AppliedUnitPrice, _ = strconv.Atoi(lAppliedUnitPrice)
						lSgbOrderHistoryRec.AppliedAmount, _ = strconv.Atoi(lAppliedAmount)
						lSgbOrderHistoryRec.AllotedUnit, _ = strconv.Atoi(lAllocatedUnits)
						lSgbOrderHistoryRec.AllotedUnitPrice, _ = strconv.Atoi(lAllocatedUnitPrice)
						lSgbOrderHistoryRec.AllotedAmount, _ = strconv.Atoi(lAllocatedAmount)
						// Append Upi End Point in lRespRec.SgbHistoryArr array
						OrderStatus, lErr4 := SgbDiscription(lStatus)
						if lErr4 != nil {
							log.Println("LGSHD04", lErr4)
							lSgbResp.Status = common.ErrorCode
							return lSgbResp, lErr4
						} else {
							lSgbOrderHistoryRec.OrderStatus = OrderStatus
						}
					}
					lSgbRespArr = append(lSgbRespArr, lSgbOrderHistoryRec)

				}

			}

		}

		if lSgbRespArr != nil {
			lSgbResp.SgbOrderHistoryArr = lSgbRespArr
			lSgbResp.OrderCount = lOrderCount
			lSgbResp.HistoryFound = "Y"
		} else {
			if len(lSgbRespArr) == 0 {
				lSgbResp.HistoryFound = "N"
				lSgbResp.HistoryNoDataText = lHistoryNoDataTxt
			}
		}
		log.Println("SgbOrderHistoryRec (-)")
		return lSgbResp, lErr1
	}
}
