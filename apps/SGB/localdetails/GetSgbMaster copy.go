package localdetails

// import (
// 	"encoding/json"
// 	"fcs23pkg/apps/Ipo/brokers"
// 	"fcs23pkg/apps/validation/apiaccess"
// 	"fcs23pkg/common"
// 	"fcs23pkg/ftdb"
// 	"fcs23pkg/helpers"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// //COMMENTED BY NITHISH
// // THIS BELOW STRUCT CONTAINS UNUSEFUL JSON AND IT'S USE FOR
// // OLD FLOW OF QUERY

// // type ActiveSgbStruct struct {
// // 	Id        int    `json:"id"`
// // 	Symbol    string `json:"symbol"`
// // 	Name      string `json:"name"`
// // 	MinBidQty int    `json:"minBidQty"`
// // 	MaxBidQty int    `json:"maxBidQty"`
// // 	Isin      string `json:"isin"`
// // 	MinPrice  int    `json:"minPrice"`
// // 	CloseDate string `json:"closeDate"`
// // 	MaxPrice  int    `json:"maxPrice"`
// // 	DateTime  string `json:"dateTime"`
// // 	Upcoming  string `json:"upcoming"`
// // 	LastDay   bool   `json:"lastDay"`
// // 	Flag      string `json:"flag"`
// // 	Pending   string `json:"pending"`
// // 	Unit      int    `json:"unit"`
// // 	OrderNo   string `json:"orderNo"`
// // 	Status    string `json:"status"`
// // 	StartTime string `json:"startTime"`
// // 	EndTime   string `json:"endTime"`
// // 	CloseTime string `json:"closeTime"`
// // }
// // type SgbStruct struct {
// // 	SgbDetails []ActiveSgbStruct `json:"sgbDetail"`
// // 	Status     string            `json:"status"`
// // 	ErrMsg     string            `json:"errMsg"`
// // }

// type ActiveSgbStruct struct {
// 	Id                int    `json:"id"`
// 	Symbol            string `json:"symbol"`
// 	Name              string `json:"name"`
// 	MinBidQty         int    `json:"minBidQty"`
// 	MaxBidQty         int    `json:"maxBidQty"`
// 	Isin              string `json:"isin"`
// 	UnitPrice         int    `json:"unitPrice"`
// 	DateRange         string `json:"dateRange"`
// 	StartDateWithTime string `json:"startDateWithTime"`
// 	EndDateWithTime   string `json:"endDateWithTime"`
// 	DiscountAmt       int    `json:"discountAmt"`
// 	DiscountText      string `json:"discountText"`
// 	ActionFlag        string `json:"actionFlag"`
// 	ButtonText        string `json:"buttonText"`
// 	AppliedUnit       int    `json:"appliedUnit"`
// 	DisableActionBtn  bool   `json:"diableActionBtn"`
// 	CancelAllowed     bool   `json:"cancelAllowed"`
// 	ModifyAllowed     bool   `json:"modifyAllowed"`
// 	OrderNo           string `json:"orderNo"`
// 	ShowSI            bool   `json:"showSI"`
// 	SIText            string `json:"SItext"`
// 	SIRefundText      string `json:"SIrefundText"`
// 	SIValue           bool   `json:"SIvalue"`
// 	InfoText          string `json:"infoText"`
// }

// // Response Structure for GetSgbMaster API
// type SgbStruct struct {
// 	// SgbDetails  []ActiveSgbStruct `json:"sgbDetail"`
// 	Disclaimer  string `json:"disclaimer"`
// 	MasterFound string `json:"masterFound"`
// 	NoDataText  string `json:"noDataText"`
// 	InvestCount int    `json:"investCount"`
// 	Status      string `json:"status"`
// 	ErrMsg      string `json:"errMsg"`
// }

// type DisclaimerStruct struct {
// 	Disclaimer   string
// 	SItext       string
// 	SIrefundText string
// }

// /*
// Pupose:This Function is used to Get the Active Ipo Details in our database table ....
// Parameters:

// not Applicable

// Response:

// *On Sucess
// =========

// 	{
// 		"sgbDetail": [
// 			{
// 			"id": 10,
// 			"symbol": "SGB201605",
// 			"name": "SOVEREIGN GOLD BONDS SCHEME 2016-17 (TRANCHE 5)",
// 			"minBidQty": 1,
// 			"maxBidQty": 4000,
// 			"isin": "IN0020150085",
// 			"unitPrice": 3149,
// 			"dateRange": "14 Sep 2013 -  19 Dec 2024",
// 			"dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
// 			"discountAmt": 1,
// 			"discountText": "Discount ₹",
// 			"actionFlag": "M",
// 			"buttonText": "Modify",
// 			"appliedUnit": 2,
// 			"diableActionBtn": false,
// 			"cancelAllowed": true,
// 			"modifyAllowed": true,
// 			"orderNo": "170375639931286",
// 			"SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
// 			"SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
// 			"SIvalue": true
// 			}
// 		],
// 		"disclaimer": "Your order will be placed on the exchange (NSE/BSE) at the end of the subscription period. Ensure to keep sufficient balance in your Trading account on the last day of the issue. Credit from \r\nstocks sold on the closing day of the issue will not be considered towards the purchase of the SGB. ",
// 		"masterFound": "",
// 		"noDataText": "",
// 		"investCount": 1,
// 		"status": "S",
// 		"errMsg": ""
// 	}

// !On Error
// ========

// 	{
// 		"status": E,
// 		"reason": "Can't able to get the data from database"
// 	}

// Author: Nithish Kumar
// Date: 05JUNE2023
// */
// func GetSgbMaster(w http.ResponseWriter, r *http.Request) {
// 	log.Println("GetSgbMaster (+)", r.Method)
// 	origin := r.Header.Get("Origin")
// 	var lBrokerId int
// 	var lErr error
// 	for _, allowedOrigin := range common.ABHIAllowOrigin {
// 		if allowedOrigin == origin {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
// 			log.Println(lErr, origin)
// 			break
// 		}
// 	}

// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
// 	if r.Method == "GET" {
// 		// create the instance for IpoStruct
// 		var lRespRec SgbStruct

// 		lRespRec.Status = common.SuccessCode

// 		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
// 		lClientId := ""
// 		var lErr1 error
// 		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
// 		if lErr1 != nil {
// 			log.Println("LGSM01", lErr1)
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "LGSM01" + lErr1.Error()
// 			fmt.Fprintf(w, helpers.GetErrorString("LGSM01", "UserDetails not Found"))
// 			return
// 		} else {
// 			if lClientId == "" {
// 				fmt.Fprintf(w, helpers.GetErrorString("LGSM02", "UserDetails not Found"))
// 				return
// 			}
// 		}
// 		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

// 		lResp, lErr2 := GetSGBdetail(lClientId, lBrokerId)
// 		if lErr2 != nil {
// 			log.Println("LGSM03", lErr2)
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "LGSM03" + lErr2.Error()
// 			fmt.Fprintf(w, helpers.GetErrorString("LGSM03", "Error Occur in getting Datas.."))
// 			return
// 		} else {
// 			lRespRec = lResp
// 		}

// 		// Marshal the Response Structure into lData
// 		lData, lErr3 := json.Marshal(lRespRec)
// 		if lErr3 != nil {
// 			log.Println("LGSM04", lErr3)
// 			fmt.Fprintf(w, helpers.GetErrorString("LGSM04", "Error Occur in getting Datas.."))
// 			return
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("GetSgbMaster (-)", r.Method)
// 	}
// }

// func GetSGBdetail_OldVersion(pClientId string, pBrokerId int) (SgbStruct, error) {
// 	log.Println("GetSGBdetail_OldVersion (+)")

// 	var lSgbMasterRec ActiveSgbStruct
// 	var lSgbMasterArr []ActiveSgbStruct
// 	var lSgbRespRec SgbStruct
// 	var lUnit string

// 	lSgbRespRec.Status = common.SuccessCode

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("GSGSD01", lErr1)
// 		lSgbRespRec.Status = common.ErrorCode
// 		return lSgbRespRec, lErr1
// 	} else {
// 		defer lDb.Close()

// 		// lDirectory, lErr2 := adminaccess.FetchDirectory()
// 		// if lErr2 != nil {
// 		// 	log.Println("GSGSD01", lErr2)
// 		// 	return lSgbRespRec, lErr2
// 		// } else {
// 		lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
// 		// lStartTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_StartTime"])
// 		// lEndTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_EndTime"])
// 		lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])
// 		lMaxQty := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_MAXQUANTITY"])
// 		lMasterNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_MasterNoDataText"])
// 		lDiscountTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_DiscountText"])
// 		lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CancelAllowed"])

// 		lConfigFile2 := common.ReadTomlConfig("toml/SgbConfig.toml")
// 		lMode := fmt.Sprintf("%v", lConfigFile2.(map[string]interface{})["SGB_Immediate_Flag"])

// 		//COMMENTED BY NITHISH
// 		// THIS QUERY CONTAINS UNUSEFUL RECORDS AND IT IS THE OLD FLOW

// 		// lCoreString := `select tab.id,tab.Symbol,tab.Name,tab.minBidqty,
// 		// tab.maxQty,tab.isin,tab.minprice,tab.maxprice,
// 		// tab.endDate,tab.formatted_datetime,tab.LastDay,
// 		// tab.Upcoming,tab.Flag,tab.Pending,lower(tab.Status),
// 		// nvl((case when tab.headerId = d.HeaderId and tab.Flag = 'N' then "" else d.ReqSubscriptionUnit end) ,"")unit,
// 		// nvl(d.ReqOrderNo,''),tab.startTime,tab.endTime
// 		//  from a_sgb_orderdetails d
// 		// 	right join (
// 		// 	select nvl(h.Id,0) headerId,master.id,master.Symbol,
// 		// 	master.Name,master.minBidqty,master.maxQty,master.isin,
// 		// 	master.minprice,master.maxprice,master.startDate,
// 		// 	master.endDate,master.Upcoming,master.formatted_datetime,master.LastDay,
// 		// 	(case when h.MasterId = master.Id and h.CancelFlag = 'N' and h.Status = 'success' then 'Y' else 'N' end) Flag,
// 		// 	(case when h.MasterId = master.Id and h.status = "pending" then "P" else "-" end) Pending,
// 		// 	(case when h.MasterId = master.Id and h.Status = "success" then h.Status else "" end) Status,
// 		// 	master.startTime,master.EndTime
// 		// 	 from (
// 		// 	 select sm.id Id,sm.Symbol Symbol,sm.Name Name,sm.MinBidQuantity minBidqty,
// 		// 	 sm.MaxQuantity maxQty,sm.Isin isin,sm.MinPrice minprice,sm.MaxPrice maxprice,
// 		// 	sm.BiddingStartDate startDate,
// 		// 	concat(
// 		// 			DATE_FORMAT(sm.BiddingStartDate, '%d %b %y'),' -  ',
// 		// 			DATE_FORMAT(sm.BiddingEndDate, '%d %b %y')
// 		// 		) as endDate,
// 		// 	(case
// 		// 	 when date_sub(sm.BiddingStartDate,interval 1 day ) = curdate() then 'P'
// 		// 	 when sm.BiddingStartDate > curdate() then 'U'
// 		// 	 else 'C' end) Upcoming,
// 		// 		CONCAT( case
// 		// 		 WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
// 		// 		WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
// 		// 		WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
// 		// 		ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
// 		// 		end,' ',
// 		// 		DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
// 		// 		TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime,
// 		// 	(case when  sm.BiddingEndDate = Date(now()) and Time(now()) > '` + lCloseTime + `' then 1 else 0 end ) as LastDay,
// 		// 	sm.DailyStartTime startTime,sm.DailyEndTime endTime
// 		// 		from a_sgb_master sm
// 		// 		where sm.BiddingEndDate >= curdate()
// 		// 		and sm.Exchange = 'NSE'
// 		// 		and sm.Redemption = 'N'
// 		// 		and not exists (
// 		// 		select 1
// 		// 		from a_sgb_master m
// 		// 		where m.BiddingEndDate = date(now())
// 		// 		and m.id = sm.id
// 		// 		and m.DailyEndTime <= time(now()))
// 		// 		) master
// 		// 	LEFT JOIN a_sgb_orderheader h
// 		// 	on master.Id = h.MasterId
// 		// 	and h.CancelFlag = 'N'
// 		// 	and h.ClientId = ?
// 		// 	and h.brokerId = ?
// 		// 	and h.Status is not null) tab
// 		// on tab.headerid = d.HeaderId
// 		// group by tab.id
// 		// order by (case when Flag = 'Y' then 1
// 		// when Upcoming = 'C' then 2
// 		// else 3 end),tab.startDate,tab.Symbol`

// 		// ==============================================================================
// 		lBrokerId := strconv.Itoa(pBrokerId)

// 		lCoreString := `select tab.id,tab.Symbol,tab.Name,tab.minBidqty,
// 						tab.maxQty,tab.isin,tab.unitPrice,
// 						tab.dateRange,tab.dateRangeWithTime,tab.discount,tab.ActionFlag,
// 						nvl((select ld.description from xx_lookup_details ld,xx_lookup_header lh
// 							where lh.id = ld.headerId and lh.Code = 'SgbAction' and ld.code = tab.ActionFlag),''
// 							) ActionBtn,
// 						nvl((case when tab.headerId = d.HeaderId and tab.ActionFlag = 'B' then 0 else d.ReqSubscriptionUnit end) ,0)unit,
// 						nvl((case when tab.ActionFlag = 'A' then 1 else tab.DisableActionBtn end),0	) DisableActionBtn,
// 						(case when tab.ActionFlag = 'B' then 0 else tab.CancelledAllowed end) CancelledAllowed,
// 						nvl(d.ReqOrderNo,'') orderNo,
// 						nvl((select bm.Sgb_SItext from a_ipo_brokermaster bm where bm.id = ` + lBrokerId + `),'') SItext,
// 						nvl((select bm.Sgb_SIRefundtext from a_ipo_brokermaster bm where bm.id = ` + lBrokerId + `),'') SIRefundtext,
// 						nvl(tab.SIvalue,0) siValue,
// 						nvl((select bm.SgbDisclaimer from a_ipo_brokermaster bm where bm.id = ` + lBrokerId + `),'') disclaimer
// 						from a_sgb_orderdetails d
// 							right join (
// 							select nvl(h.Id,0) headerId,h.brokerId,master.id,master.Symbol,
// 							master.Name,master.minBidqty,master.maxQty,master.isin,
// 							master.unitPrice,master.dateRange,master.dateRangeWithTime,
// 							(case
// 								when h.ProcessFlag = 'Y' and h.ScheduleStatus = 'Y' then 'A'
// 								when h.MasterId = master.Id and h.CancelFlag = 'N' and h.Status = 'success' then 'M'
// 								when master.Upcoming = '' then 'B'
// 								else master.Upcoming
// 							end) ActionFlag,
// 							master.DisableActionBtn,master.CancelledAllowed,h.SIvalue,master.discount
// 							from (
// 								select sm.id Id,sm.Symbol Symbol,sm.Name Name,sm.MinBidQuantity minBidqty,
// 								sm.MaxQuantity maxQty,sm.Isin isin,sm.MinPrice unitPrice,
// 								concat(
// 									DATE_FORMAT(sm.BiddingStartDate, '%d %b %Y'),' -  ',DATE_FORMAT(sm.BiddingEndDate, '%d %b %Y')
// 								) as dateRange,
// 								concat(DATE_FORMAT(sm.BiddingStartDate, '%d %b %Y'),' ',TIME_FORMAT(sm.DailyStartTime , '%h:%i%p')) as startDateWithTime,
// 								concat(DATE_FORMAT(sm.BiddingEndDate, '%d %b %Y'),' ',TIME_FORMAT('` + lCloseTime + `' , '%h:%i%p')) as endDateWithTime,
// 								(case
// 									when date_sub(sm.BiddingStartDate,interval 1 day ) = curdate() then 'P'
// 									when sm.BiddingStartDate > curdate() then 'U'
// 									when sm.BiddingEndDate = Date(now()) and Time(now()) > '` + lCloseTime + `' then 'C'
// 									else ''
// 								end) Upcoming,
// 								(case when  sm.BiddingEndDate = Date(now()) and Time(now()) > '` + lCloseTime + `' then 1 else 0 end ) as DisableActionBtn,
// 								(case when '` + lCancelAllowed + `' = '' then 0 else 1 end) CancelledAllowed,
// 								(case when (sm.FaceValue-sm.minPrice) < 0 then 0 else (sm.FaceValue-sm.minPrice) end )discount,
// 								sm.DailyStartTime startTime,sm.DailyEndTime endTime
// 									from a_sgb_master sm
// 									where sm.BiddingEndDate >= curdate()
// 									and sm.Exchange = 'NSE'
// 									and sm.Redemption = 'N'
// 									and not exists (
// 									select 1
// 									from a_sgb_master m
// 									where m.BiddingEndDate = date(now())
// 									and m.id = sm.id
// 									and m.DailyEndTime <= time(now()))
// 									) master
// 							LEFT JOIN a_sgb_orderheader h
// 							on master.Id = h.MasterId
// 							and h.CancelFlag = 'N'
// 							and h.ClientId = ?
// 							and h.brokerId = ?
// 							and h.Status is not null) tab
// 						on tab.headerid = d.HeaderId
// 						group by tab.id
// 						order by (case when ActionFlag = 'M' then 1
// 						when ActionFlag = 'C' then 2
// 						else 3 end)`

// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("GSGSD02", lErr2)
// 			lSgbRespRec.Status = common.ErrorCode
// 			return lSgbRespRec, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lSgbMasterRec.Id, &lSgbMasterRec.Symbol, &lSgbMasterRec.Name, &lSgbMasterRec.MinBidQty, &lSgbMasterRec.MaxBidQty, &lSgbMasterRec.Isin, &lSgbMasterRec.UnitPrice, &lSgbMasterRec.DateRange, &lSgbMasterRec.StartDateWithTime, &lSgbMasterRec.EndDateWithTime, &lSgbMasterRec.DiscountAmt, &lSgbMasterRec.ActionFlag, &lSgbMasterRec.ButtonText, &lUnit, &lSgbMasterRec.DisableActionBtn, &lSgbMasterRec.CancelAllowed, &lSgbMasterRec.OrderNo, &lSgbMasterRec.SIText, &lSgbMasterRec.SIRefundText, &lSgbMasterRec.SIValue, &lSgbRespRec.Disclaimer)
// 				if lErr3 != nil {
// 					log.Println("GSGSD03", lErr3)
// 					lSgbRespRec.Status = common.ErrorCode
// 					return lSgbRespRec, lErr3
// 				} else {
// 					// lSgbMasterRec.StartTime = lStartTime
// 					// lSgbMasterRec.EndTime = lEndTime
// 					// lSgbMasterRec.CloseTime = lCloseTime
// 					lSgbMasterRec.AppliedUnit, _ = strconv.Atoi(lUnit)
// 					lSgbMasterRec.DiscountText = lDiscountTxt

// 					if lMaxQty != "" {
// 						lSgbMasterRec.MaxBidQty, _ = strconv.Atoi(lMaxQty)
// 					}
// 					if lSgbMasterRec.ActionFlag == "M" {
// 						lSgbRespRec.InvestCount++
// 					}
// 					if lMode != "" {
// 						lSgbMasterRec.ModifyAllowed = false
// 					} else {
// 						lSgbMasterRec.ModifyAllowed = true
// 					}
// 					lSgbMasterArr = append(lSgbMasterArr, lSgbMasterRec)
// 				}
// 			}
// 			if lSgbMasterArr != nil {
// 				lSgbRespRec.SgbDetails = lSgbMasterArr
// 				// lSgbRespRec.MasterFound = "Y"
// 			} else {
// 				if len(lSgbMasterArr) == 0 {
// 					// lSgbRespRec.MasterFound = "N"
// 					lSgbRespRec.NoDataText = lMasterNoDataTxt
// 				}
// 			}
// 		}

// 	}
// 	log.Println("GetSGBdetail_OldVersion (-)")
// 	return lSgbRespRec, nil
// }

// /*
// Pupose:This Function is used to construct the final SGB master response
// Parameters:

// ClientId ,BrokerId

// Response:

// *On Sucess
// =========
// 	{
// 		"sgbDetail": [
// 			{
// 			"id": 10,
// 			"symbol": "SGB201605",
// 			"name": "SOVEREIGN GOLD BONDS SCHEME 2016-17 (TRANCHE 5)",
// 			"minBidQty": 1,
// 			"maxBidQty": 4000,
// 			"isin": "IN0020150085",
// 			"unitPrice": 3149,
// 			"dateRange": "14 Sep 2013 -  19 Dec 2024",
// 			"dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
// 			"discountAmt": 1,
// 			"discountText": "Discount ₹",
// 			"actionFlag": "M",
// 			"buttonText": "Modify",
// 			"appliedUnit": 2,
// 			"diableActionBtn": false,
// 			"cancelAllowed": true,
// 			"modifyAllowed": true,
// 			"orderNo": "170375639931286",
// 			"SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
// 			"SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
// 			"SIvalue": true
// 			}
// 		],
// 		"disclaimer": "Your order will be placed on the exchange (NSE/BSE) at the end of the subscription period. Ensure to keep sufficient balance in your Trading account on the last day of the issue. Credit from \r\nstocks sold on the closing day of the issue will not be considered towards the purchase of the SGB. ",
// 		"masterFound": "",
// 		"noDataText": "",
// 		"investCount": 1,
// 		"status": "S",
// 		"errMsg": ""
// 	}

// !On Error
// ========

// 	{
// 		"status": E,
// 		"reason": "Can't able to get the data from database"
// 	}

// Author: Nithish Kumar
// Date: 28DEC2023
// */
// func GetSGBdetail(pClientId string, pBrokerId int) (SgbStruct, error) {
// 	log.Println("GetSGBdetail (+)")
// 	var lSgbRespRec SgbStruct
// 	lSgbRespRec.Status = common.SuccessCode

// 	lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
// 	lMasterNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_MasterNoDataText"])

// 	lSgbMasterArr, lErr1 := SgbMasterRec("getSgbMaster", pBrokerId, pClientId)
// 	if lErr1 != nil {
// 		log.Println("LGSD01", lErr1)
// 		lSgbRespRec.Status = common.ErrorCode
// 	} else {
// 		lDisclaimerResp, lErr2 := GetSgbDisclaimer(pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("LGSD02", lErr2)
// 			lSgbRespRec.Status = common.ErrorCode
// 		} else {

// 			lSgbRespRec.SgbDetails = lSgbMasterArr
// 			lSgbRespRec.Disclaimer = lDisclaimerResp.Disclaimer

// 			if len(lSgbMasterArr) > 0 {
// 				// lSgbRespRec.MasterFound = "N"
// 				lSgbRespRec.NoDataText = lMasterNoDataTxt

// 				for _, lCount := range lSgbMasterArr {
// 					// When the client placed a order successfully the the investCount will be increased
// 					// based on the no.of bid they placed.
// 					if lCount.ActionFlag == "M" || lCount.ActionFlag == "A" || (lCount.ActionFlag == "C" &&
// 						lCount.AppliedUnit > 0) {
// 						lSgbRespRec.InvestCount++
// 					}
// 				}
// 			} else {
// 				// lSgbRespRec.MasterFound = "N"
// 				lSgbRespRec.NoDataText = lMasterNoDataTxt
// 			}
// 		}
// 	}
// 	log.Println("GetSGBdetail (-)")
// 	return lSgbRespRec, nil
// }

// /*
// Pupose:This Function gives the information about the active SGB master detail along with the client applied status
// Parameters:

// 	Method - (Need to tell where the method calling from "getSgbMaster" / "getSgbOrderHistory" ),
// 			  ClientId ,BrokerId
// Response:

// *On Sucess
// =========
// 	"sgbDetail": [
// 		{
// 		"id": 10,
// 		"symbol": "SGB201605",
// 		"name": "SOVEREIGN GOLD BONDS SCHEME 2016-17 (TRANCHE 5)",
// 		"minBidQty": 1,
// 		"maxBidQty": 4000,
// 		"isin": "IN0020150085",
// 		"unitPrice": 3149,
// 		"dateRange": "14 Sep 2013 -  19 Dec 2024",
// 		"dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
// 		"discountAmt": 1,
// 		"discountText": "Discount ₹",
// 		"actionFlag": "M",
// 		"buttonText": "Modify",
// 		"appliedUnit": 2,
// 		"diableActionBtn": false,
// 		"cancelAllowed": true,
// 		"modifyAllowed": true,
// 		"orderNo": "170375639931286",
// 		"SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
// 		"SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
// 		"SIvalue": true
// 		},
// 	]

// !On Error
// ========

// 	{
// 		"status": E,
// 		"reason": "Can't able to get the data from database"
// 	}

// Author: Nithish Kumar
// Date: 28DEC2023
// */
// func SgbMasterRec(pMethod string, pBrokerId int, pClientId string) ([]ActiveSgbStruct, error) {
// 	log.Println("SgbMasterRec (+)")
// 	var lSgbMasterRec ActiveSgbStruct
// 	var lSgbMasterArr []ActiveSgbStruct
// 	var lSgbRespRec SgbStruct

// 	lEndDate := ""

// 	lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
// 	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])
// 	lMaxQty := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_MAXQUANTITY"])
// 	lDiscountTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_DiscountText"])
// 	lShowDiscount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_ShowDiscount"])
// 	lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CancelAllowed"])
// 	lProcessingMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_Immediate_Flag"])
// 	lDefaultRefundText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_DefaultRefundText"])

// 	lDisclaimerResp, lErr1 := GetSgbDisclaimer(pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("LSMR01", lErr1)
// 		lSgbRespRec.Status = common.ErrorCode
// 	} else {

// 		// Calling LocalDbConect method in ftdb to estabish the database connection
// 		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
// 		if lErr2 != nil {
// 			log.Println("LSMR02", lErr2)
// 			return lSgbMasterArr, lErr2
// 		} else {
// 			defer lDb.Close()

// 			// Get the current time
// 			lCurrentTime := time.Now()

// 			// Format the time as a string
// 			lTimeString := lCurrentTime.Format("15:04:05")

// 			// Trim leading and trailing spaces
// 			lCurrentTimeString := strings.TrimSpace(lTimeString)

// 			ladditionalString := ""
// 			// If you change the unitprice in the query
// 			// you shoul compulory change the discount calculation first value
// 			// (eg) := sm.MaxPrice unitPrice = 6199 ,(sm.MaxPrice-sm.MinPrice)discount = 6199 - 6149,(Final discount = 50)
// 			// (eg) := sm.MinPrice unitPrice = 6149 ,(sm.MinPrice-sm.MinPrice)discount = 6149 - 6149,(Final discount = 0)

// 			lCoreString := `select
// 							sm.id Id,sm.Symbol Symbol,sm.Name Name,sm.MinBidQuantity minBidqty,
// 							sm.MaxQuantity maxQty
// 							,sm.Isin isin
// 							,(case when '` + lShowDiscount + `' = 'Y' then sm.MaxPrice  else sm.MinPrice end ) unitPrice,
// 							concat(
// 								DATE_FORMAT(sm.BiddingStartDate, '%d %b %Y'),' -  ',DATE_FORMAT(sm.BiddingEndDate, '%d %b %Y')
// 							) as dateRange,
// 							concat(DATE_FORMAT(sm.BiddingStartDate, '%d %b %Y'),' ',TIME_FORMAT(sm.DailyStartTime , '%h:%i%p')) as startDateWithTime,
// 							concat(DATE_FORMAT(sm.BiddingEndDate, '%d %b %Y'),' ',TIME_FORMAT('` + lCloseTime + `' , '%h:%i%p')) as endDateWithTime,
// 							(case
// 									when date_sub(sm.BiddingStartDate,interval 1 day ) = curdate() then 'P'
// 									when sm.BiddingStartDate > curdate() then 'U'
// 									when sm.BiddingEndDate = Date(now()) and '` + lCurrentTimeString + `' > '` + lCloseTime + `' then 'C'
// 									else ''
// 								end) ActionFlag,
// 							(case when  sm.BiddingEndDate = Date(now()) and '` + lCurrentTimeString + `' > '` + lCloseTime + `' then 1 else 0 end ) as DisableActionBtn,
// 							(case when '` + lCancelAllowed + `' = 'Y' then 0 else 1 end) CancelledAllowed,
// 							(case when '` + lShowDiscount + `' = 'Y' then sm.MaxPrice-sm.MinPrice else 0 end )discount,
// 							sm.BiddingEndDate
// 						from a_sgb_master sm
// 						where sm.Exchange = 'NSE'
// 						and sm.Redemption = 'N'`

// 			//below condition will added when this method is called for master data that are currently active and upcoming
// 			if pMethod == "getSgbMaster" {
// 				ladditionalString = ` and sm.BiddingEndDate >= curdate()
// 			and not exists (
// 				select 1
// 				from a_sgb_master m
// 				where m.BiddingEndDate = date(now())
// 				and m.id = sm.id
// 				and m.DailyEndTime <= time(now())
// 			)`
// 			}

// 			lCoreFinalString := lCoreString + ladditionalString

// 			lRows, lErr3 := lDb.Query(lCoreFinalString)
// 			if lErr3 != nil {
// 				log.Println("LSMR03", lErr3)
// 				lSgbRespRec.Status = common.ErrorCode
// 				return lSgbMasterArr, lErr3
// 			} else {
// 				//This for loop is used to collect the records from the database and store them in structure
// 				for lRows.Next() {
// 					lErr4 := lRows.Scan(&lSgbMasterRec.Id, &lSgbMasterRec.Symbol, &lSgbMasterRec.Name, &lSgbMasterRec.MinBidQty, &lSgbMasterRec.MaxBidQty, &lSgbMasterRec.Isin, &lSgbMasterRec.UnitPrice, &lSgbMasterRec.DateRange, &lSgbMasterRec.StartDateWithTime, &lSgbMasterRec.EndDateWithTime, &lSgbMasterRec.ActionFlag, &lSgbMasterRec.DisableActionBtn, &lSgbMasterRec.CancelAllowed,
// 						&lSgbMasterRec.DiscountAmt, &lEndDate)
// 					if lErr4 != nil {
// 						log.Println("LSMR04", lErr4)
// 						lSgbRespRec.Status = common.ErrorCode
// 						return lSgbMasterArr, lErr4
// 					} else {
// 						// lSgbMasterRec.CloseTime = lCloseTime

// 						// in any case if exchange sends max value lesser than min value, then the discount computed will
// 						// be in negative. so resting it to zero
// 						if lSgbMasterRec.DiscountAmt < 0 {
// 							lSgbMasterRec.DiscountAmt = 0
// 						}

// 						lSgbMasterRec.DiscountText = lDiscountTxt

// 						//if the processing of the order is set as immediate
// 						if lProcessingMode == "I" {
// 							//standing instruction to be displayed while cancelling the order
// 							lSgbMasterRec.SIRefundText = lDisclaimerResp.SIrefundText

// 							//if the processing of the order is set as lastday
// 						} else {
// 							//standing instruction to be displayed while pacing the order
// 							lSgbMasterRec.SIText = lDisclaimerResp.SItext
// 							// Show the default text if client try to cancel the bid
// 							lSgbMasterRec.SIRefundText = lDefaultRefundText
// 						}

// 						// Alter the MaximumBidQty allowed for the SGB
// 						if lMaxQty != "" {
// 							lSgbMasterRec.MaxBidQty, _ = strconv.Atoi(lMaxQty)
// 						}
// 						// This method is called only for getSgbMaster
// 						if pMethod == "getSgbMaster" {

// 							lSgbRec, lErr5 := sgbOrderRec(lSgbMasterRec, pClientId, pBrokerId, lEndDate)
// 							if lErr5 != nil {
// 								log.Println("LSMR05", lErr5)
// 								lSgbRespRec.Status = common.ErrorCode
// 								return lSgbMasterArr, lErr5
// 							} else {
// 								lSgbMasterArr = append(lSgbMasterArr, lSgbRec)
// 							}
// 						} else if pMethod == "getSgbOrderHistory" {
// 							lSgbMasterArr = append(lSgbMasterArr, lSgbMasterRec)

// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("SgbMasterRec (-)")
// 	return lSgbMasterArr, nil
// }

// /*
// Pupose:This Function is used to Get client order information for each active SGB
// Parameters:

// 	ActiveSgbStruct, ClientId , BrokerId , EndDate

// Response:

// *On Sucess
// =========
// 	{
//       "id": 10,
//       "symbol": "SGB201605",
//       "name": "SOVEREIGN GOLD BONDS SCHEME 2016-17 (TRANCHE 5)",
//       "minBidQty": 1,
//       "maxBidQty": 4000,
//       "isin": "IN0020150085",
//       "unitPrice": 3149,
//       "dateRange": "14 Sep 2013 -  19 Dec 2024",
//       "dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
//       "discountAmt": 1,
//       "discountText": "Discount ₹",
//       "actionFlag": "M",
//       "buttonText": "",
//       "appliedUnit": 2,
//       "diableActionBtn": false,
//       "cancelAllowed": true,
//       "modifyAllowed": true,
//       "orderNo": "170375639931286",
//       "SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
//       "SIrefundText": "",
//       "SIvalue": true
//     },

// !On Error
// ========

// 	{
// 		"status": E,
// 		"reason": "Can't able to get the data from database"
// 	}

// Author: Nithish Kumar
// Date: 28DEC2023
// */
// func sgbOrderRec(pSgbMasterRec ActiveSgbStruct, pClientId string, pBrokerId int, pEndDate string) (ActiveSgbStruct, error) {
// 	log.Println("SgbOrderRec (+)")
// 	var lUnit, lProcessFlag, lScheduleFlag, lCancelFlag, lStatus, lReqNo, lSIText string
// 	var lSIValue bool

// 	lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
// 	lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CancelAllowed"])
// 	lInfoText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_InfoText"])
// 	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])
// 	lProcessingMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_Immediate_Flag"])
// 	lModifyAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_ModifyAllowed"])

// 	if lModifyAllowed == "Y" {
// 		pSgbMasterRec.ModifyAllowed = true
// 	} else {
// 		pSgbMasterRec.ModifyAllowed = false
// 	}

// 	//if order exists or not, set the defaults for the record
// 	pSgbMasterRec.CancelAllowed = false
// 	// pSgbMasterRec.ModifyAllowed = true
// 	pSgbMasterRec.DisableActionBtn = true
// 	pSgbMasterRec.ShowSI = true
// 	pSgbMasterRec.AppliedUnit = 0

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("LSOR01", lErr1)
// 		return pSgbMasterRec, lErr1
// 	} else {
// 		defer lDb.Close()

// 		// Get the current time
// 		lCurrentTime := time.Now()

// 		// Format the time as a string
// 		lTimeString := lCurrentTime.Format("15:04:05")

// 		// Get the current date
// 		lCurrentDate := time.Now().Format("2006-01-02") // YYYY-MM-DD

// 		// Trim leading and trailing spaces
// 		lCurrentTimeString := strings.TrimSpace(lTimeString)

// 		lCoreString := `select (case when d.ReqSubscriptionUnit = '' then 0 else d.ReqSubscriptionUnit end)unit,
// 						nvl(h.ProcessFlag,'N'),nvl(h.ScheduleStatus,'N'),
// 						nvl(h.CancelFlag,'N'),nvl(h.Status,''),nvl(d.ReqOrderNo,'') orderNo,
// 						nvl(h.SIvalue,0) siValue,nvl(h.SItext,'') siText
// 						from a_sgb_orderheader h,a_sgb_orderdetails d
// 						where h.id = d.HeaderId
// 						and h.brokerId = ?
// 						and h.clientId = ?
// 						and h.MasterId = ?
// 						and nvl(h.CancelFlag,'N') = 'N'
// 						and h.Status = 'success'`

// 		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId, pClientId, pSgbMasterRec.Id)
// 		if lErr2 != nil {
// 			log.Println("LSOR02", lErr2)
// 			return pSgbMasterRec, lErr2
// 		} else {

// 			//This for loop is used to collect the records from the database and store them in respective structure
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lUnit, &lProcessFlag, &lScheduleFlag, &lCancelFlag, &lStatus, &lReqNo, &lSIValue, &lSIText)
// 				if lErr3 != nil {
// 					log.Println("LSOR03", lErr3)
// 					return pSgbMasterRec, lErr3
// 				} else {

// 					//It gives the client applied unit.
// 					pSgbMasterRec.AppliedUnit, _ = strconv.Atoi(lUnit)

// 					// &lUnit, &pSgbMasterRec.ActionFlag, &pSgbMasterRec.OrderNo, &pSgbMasterRec.SIValue, &pSgbMasterRec.SIText
// 					pSgbMasterRec.SIText = lSIText
// 					pSgbMasterRec.SIValue = lSIValue
// 					pSgbMasterRec.OrderNo = lReqNo

// 					// Manipulate the flag based on the conditions
// 					// CASE 1: when the order is placed to exchange then lProcessFlag and lScheduleFlag will be Y in the database
// 					if lProcessFlag == "Y" && lScheduleFlag == "Y" {
// 						// Set the flag A when the Jv process is completed
// 						//Mark the record as Applied
// 						pSgbMasterRec.ActionFlag = "A"
// 						// CASE 2: when the order is not placed to exchange, but eligble for modify, then show modify
// 					} else if pSgbMasterRec.ActionFlag != "C" {
// 						pSgbMasterRec.ActionFlag = "M"
// 					}

// 					// 	} else if lStatus == common.SUCCESS && pEndDate == lCurrentDate && lCurrentTimeString > lCloseTime {
// 					// 	// Set the flag C when the client applied order Endtime is exceeded current time
// 					// 	pSgbMasterRec.ActionFlag = "C"
// 					// } else if lCancelFlag == "N" && lStatus == common.SUCCESS {
// 					// 	// Set the flag M when the client already applied for the SGB
// 					// 	pSgbMasterRec.ActionFlag = "M"
// 					// } else if pSgbMasterRec.ActionFlag == "" {
// 					// 	// Set the flag B when the client doesn't applied for the SGB
// 					// 	pSgbMasterRec.ActionFlag = "B"
// 					// }

// 					//check if the global configuration is set to disable the cancel option
// 					if lCancelAllowed == "N" {
// 						pSgbMasterRec.CancelAllowed = false
// 					} else {
// 						//if the processing of the order is set as immediate
// 						if lProcessingMode == "I" {
// 							pSgbMasterRec.ModifyAllowed = false

// 							//incase the order is placed in immediate mode, then user can cancel and reapply
// 							pSgbMasterRec.InfoText = lInfoText

// 							pSgbMasterRec.CancelAllowed = false

// 							pSgbMasterRec.ShowSI = false
// 							//if the processing of the order is set as lastday
// 						} else {

// 							//if the record is already placed to exchange
// 							if lProcessFlag != "N" && lScheduleFlag != "N" {
// 								pSgbMasterRec.CancelAllowed = false
// 								pSgbMasterRec.ModifyAllowed = false

// 								// if today is the last day of the bid and current time is less than close time
// 								// then allow to cancel
// 							} else if pEndDate == lCurrentDate && lCurrentTimeString < lCloseTime {
// 								pSgbMasterRec.CancelAllowed = true
// 								pSgbMasterRec.ModifyAllowed = true

// 								//if today is not the last day of the bid
// 							} else if pEndDate > lCurrentDate {
// 								pSgbMasterRec.CancelAllowed = true
// 								pSgbMasterRec.ModifyAllowed = true

// 							} else {
// 								pSgbMasterRec.CancelAllowed = false
// 								pSgbMasterRec.ModifyAllowed = false
// 							}
// 						}
// 					}
// 					/*****************
// 					if pSgbMasterRec.ActionFlag == "A" {
// 						// ( A - Applied )
// 						// Modification will be prohibited when the client's order processed on the Exchange.
// 						pSgbMasterRec.ModifyAllowed = false

// 						// Cancellation will be allowed only when the placeorder mode in (I - ImmediateProcess)
// 						// and also the button should be clickable only on the immediate process
// 						// otherwise it will be prohibited.
// 						// The information text only visible when the placeorder mode is in "I"
// 						if lProcessingMode == "I" {
// 							pSgbMasterRec.DisableActionBtn = false
// 							pSgbMasterRec.CancelAllowed = true
// 							pSgbMasterRec.InfoText = lInfoText
// 						} else {
// 							pSgbMasterRec.DisableActionBtn = true
// 							pSgbMasterRec.CancelAllowed = false
// 						}
// 					} else if pSgbMasterRec.ActionFlag == "U" || pSgbMasterRec.ActionFlag == "C" {
// 						// ( C - Closed, U - Upcoming)
// 						// The Action button stays disable when the application is closed or upcoming
// 						pSgbMasterRec.DisableActionBtn = true
// 					}

// 					if lProcessingMode == "I" {
// 						pSgbMasterRec.SIText = ""
// 					}

// 					// To get the Button text based on the actionFlag
// 					lButtonText, lErr4 := SgbDiscription(pSgbMasterRec.ActionFlag)
// 					if lErr4 != nil {
// 						log.Println("LSOR03", lErr4)
// 						lSgbRespRec.Status = common.ErrorCode
// 						return pSgbMasterRec, lErr4
// 					} else {
// 						pSgbMasterRec.ButtonText = lButtonText
// 					}
// 					return pSgbMasterRec, nil
// 					***************/
// 				}
// 			}
// 			// if there is no order placed earlier and the action flag is blank then ready to bid
// 			if pSgbMasterRec.ActionFlag == "" {
// 				pSgbMasterRec.ActionFlag = "B"
// 			}
// 			//get the description of the action button
// 			lButtonText, lErr4 := SgbDiscription(pSgbMasterRec.ActionFlag)
// 			if lErr4 != nil {
// 				log.Println("LSOR03", lErr4)
// 				return pSgbMasterRec, lErr4
// 			} else {
// 				pSgbMasterRec.ButtonText = lButtonText
// 			}

// 			if pSgbMasterRec.ActionFlag == "P" || pSgbMasterRec.ActionFlag == "B" || pSgbMasterRec.ActionFlag == "M" {
// 				pSgbMasterRec.DisableActionBtn = false
// 			}
// 			if lProcessingMode == "I" {
// 				//if the processing of the order is set as immediate
// 				pSgbMasterRec.ShowSI = false
// 			}

// 			/************************8
// 			// If the query does't gives any record then construct
// 			// the default record by our own
// 			if lUnit == "" && lProcessFlag == "" && lScheduleFlag == "" && lCancelFlag == "" && lStatus == "" &&
// 				lReqNo == "" && lSIText == "" && lSIValue == false {
// 				pSgbMasterRec.AppliedUnit = 0
// 				lProcessFlag = "N"
// 				lScheduleFlag = "N"
// 				lCancelFlag = "N"

// 				// If the actionFlag is null then change it to (B - PlaceOrder)
// 				if pSgbMasterRec.ActionFlag == "" {
// 					pSgbMasterRec.ActionFlag = "B"
// 					pSgbMasterRec.ModifyAllowed = true
// 				} else if pSgbMasterRec.ActionFlag == "P" {
// 					pSgbMasterRec.ModifyAllowed = true
// 				} else if pSgbMasterRec.ActionFlag == "U" || pSgbMasterRec.ActionFlag == "C" {
// 					pSgbMasterRec.DisableActionBtn = true
// 				}
// 				// SI text shouldn't be shown the placeOrder process is in (I - Immediate)
// 				if lProcessingMode == "I" {
// 					pSgbMasterRec.SIText = ""
// 				} else {
// 					pSgbMasterRec.SIText = lSIText
// 				}

// 				// To get the Button text based on the actionFlag
// 				lButtonText, lErr4 := SgbDiscription(pSgbMasterRec.ActionFlag)
// 				if lErr4 != nil {
// 					log.Println("LSOR03", lErr4)
// 					lSgbRespRec.Status = common.ErrorCode
// 					return pSgbMasterRec, lErr4
// 				} else {
// 					pSgbMasterRec.ButtonText = lButtonText
// 				}
// 				return pSgbMasterRec, nil
// 			}
// 			*******************/
// 		}
// 	}
// 	log.Println("SgbOrderRec (-)")
// 	return pSgbMasterRec, nil
// }

// /*
// Pupose:This Function is used to Get SGB ActionFlag discription from the database
// Parameters:

// 	"actionFlag": "M",

// Response:

// *On Sucess
// =========
// 	"buttonText": "Modify",

// !On Error
// ========

// 	{
// 		"status": E,
// 		"reason": "Can't able to get the data from database"
// 	}

// Author: Nithish Kumar
// Date: 28DEC2023
// */
// func SgbDiscription(pActionFlag string) (string, error) {
// 	log.Println("SgbDiscription (+)")
// 	lButtonText := ""

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("LSD01", lErr1)
// 		return lButtonText, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `
// 					Select nvl(
// 						(select ld.description
// 						from xx_lookup_details ld,xx_lookup_header lh
// 						where lh.id = ld.headerId and lh.Code = 'SgbAction' and ld.code = '` + pActionFlag + `'),''
// 					) ActionBtn`
// 		lRows, lErr2 := lDb.Query(lCoreString)
// 		if lErr2 != nil {
// 			log.Println("LSD02", lErr2)
// 			return lButtonText, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in respective structure
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lButtonText)
// 				if lErr3 != nil {
// 					log.Println("LSD03", lErr3)
// 					return lButtonText, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("SgbDiscription (-)")
// 	return lButtonText, nil
// }

// /*
// Pupose:This Function is used to Get SGB disclamier,SIText and SIRefundText accoring to the Broker
// Parameters:

//  BrokerId

// Response:

// *On Sucess
// =========

// 	{
// 		"SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
// 		"SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
// 		"disclaimer": "Your order will be placed on the exchange (NSE/BSE) at the end of the subscription period. Ensure to keep sufficient balance
// 					   in your Trading account on the last day of the issue. Credit from \r\nstocks sold on the closing day of the issue will not be considered towards the purchase of the SGB. ",
// 	}

// !On Error
// ========

// 	{
// 		"status": E,
// 		"reason": "Can't able to get the data from database"
// 	}

// Author: Nithish Kumar
// Date: 27DEC2023
// */
// func GetSgbDisclaimer(pBrokerId int) (DisclaimerStruct, error) {
// 	log.Println("GetSgbDisclaimer (+)")
// 	var lSgbDisclaimer DisclaimerStruct

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("LGSD01", lErr1)
// 		return lSgbDisclaimer, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `
// 					SELECT nvl(bm.Sgb_SItext,'')SItext,nvl(bm.Sgb_SIRefundtext,'') SIRefundtext,nvl(bm.SgbDisclaimer,'') disclaimer
// 					FROM a_ipo_brokermaster bm
// 					WHERE bm.id = ?`
// 		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("LGSD02", lErr2)
// 			return lSgbDisclaimer, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in respective structure
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lSgbDisclaimer.SItext, &lSgbDisclaimer.SIrefundText, &lSgbDisclaimer.Disclaimer)
// 				if lErr3 != nil {
// 					log.Println("LGSD03", lErr3)
// 					return lSgbDisclaimer, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetSgbDisclaimer (-)")
// 	return lSgbDisclaimer, nil
// }
