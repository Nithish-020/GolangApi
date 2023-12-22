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

type ActiveSgbStruct struct {
	Id        int    `json:"id"`
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	MinBidQty int    `json:"minBidQty"`
	MaxBidQty int    `json:"maxBidQty"`
	Isin      string `json:"isin"`
	MinPrice  int    `json:"minPrice"`
	CloseDate string `json:"closeDate"`
	MaxPrice  int    `json:"maxPrice"`
	DateTime  string `json:"dateTime"`
	Upcoming  string `json:"upcoming"`
	LastDay   bool   `json:"lastDay"`
	Flag      string `json:"flag"`
	Pending   string `json:"pending"`
	Unit      int    `json:"unit"`
	OrderNo   string `json:"orderNo"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	CloseTime string `json:"closeTime"`
}

// Response Structure for GetSgbMaster API
type SgbStruct struct {
	SgbDetails []ActiveSgbStruct `json:"sgbDetail"`
	Status     string            `json:"status"`
	ErrMsg     string            `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Active Ipo Details in our database table ....
Parameters:

not Applicable

Response:

*On Sucess
=========

	{
		"IpoDetails": [
			{
				"id": 18,
				"symbol": "MMIPO26",
				"startDate": "2023-06-02",
				"endDate": "2023-06-30",
				"priceRange": "1000 - 2000",
				"cutOffPrice": 2000,
				"minBidQuantity": 10,
				"applicationStatus": "Pending",
				"upiStatus": "Accepted by Investor"
			},
			{
				"id": 10,
				"symbol": "fixed",
				"startDate": "2023-05-10",
				"endDate": "2023-08-29",
				"priceRange": "755 - 755",
				"cutOffPrice": 755,
				"minBidQuantity": 100,
				"applicationStatus": "-",
				"upiStatus": "-"
			}
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

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetSgbMaster(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSgbMaster (+)", r.Method)
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
		var lRespRec SgbStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
		if lErr1 != nil {
			log.Println("LGSM01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGSM02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lRespArr, lErr2 := GetSGBdetail(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("LGSM03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSM03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSM03", "Error Occur in getting Datas.."))
			return
		} else {
			if lRespArr != nil {
				lRespRec.SgbDetails = lRespArr
			}
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGSM04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGSM04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetSgbMaster (-)", r.Method)
	}
}

func GetSGBdetail(pClientId string, pBrokerId int) ([]ActiveSgbStruct, error) {
	log.Println("GetSGBdetail (+)")

	var lSgbMasterRec ActiveSgbStruct
	var lSgbMasterArr []ActiveSgbStruct
	var lUnit string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GSGSD01", lErr1)
		return lSgbMasterArr, lErr1
	} else {
		defer lDb.Close()

		// lDirectory, lErr2 := adminaccess.FetchDirectory()
		// if lErr2 != nil {
		// 	log.Println("GSGSD01", lErr2)
		// 	return lSgbMasterArr, lErr2
		// } else {
		lConfigFile := common.ReadTomlConfig("toml/debug.toml")
		lStartTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_StartTime"])
		lEndTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_EndTime"])
		lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])

		//COMMENTED BY NITHISH
		// THIS QUERY DOESN'T TAKE DAILYSTARTTIME & ENDTIME AND THE CLOSETIME WAS HARDCODED
		// lCoreString := `select
		// 				tab.id,
		// 				tab.Symbol,
		// 				tab.Name,
		// 				tab.minBidqty,
		// 				tab.maxQty,
		// 				tab.isin,
		// 				tab.minprice,
		// 				tab.maxprice,
		// 				tab.endDate,
		// 				tab.formatted_datetime,
		// 				tab.LastDay,
		// 				tab.Upcoming,
		// 				tab.Flag,
		// 				tab.Pending,
		// 				lower(tab.Status),
		// 				nvl((case when tab.headerId = d.HeaderId and tab.Flag = 'N' then "" else d.ReqSubscriptionUnit end) ,"")unit,
		// 				nvl(d.OrderNo,'')
		// 			from a_sgb_orderdetails d
		// 			right join (select
		// 				nvl(h.Id,0) headerId,
		// 				master.id,
		// 				master.Symbol,
		// 				master.Name,
		// 				master.minBidqty,
		// 				master.maxQty,
		// 				master.isin,
		// 				master.minprice,
		// 				master.maxprice,
		// 				master.startDate,
		// 				master.endDate,
		// 				master.Upcoming,
		// 				master.formatted_datetime,
		// 				master.LastDay,
		// 				(case when h.MasterId = master.Id and h.CancelFlag = 'N' and h.Status = 'success' then 'Y' else 'N' end) Flag,
		// 				(case when h.MasterId = master.Id and h.status = "pending" then "P" else "-" end) Pending,
		// 				(case when h.MasterId = master.Id and h.Status = "success" then h.Status else "" end) Status
		// 				from (select sm.id Id,
		// 					sm.Symbol Symbol,
		// 					sm.Name Name,
		// 					sm.MinBidQuantity minBidqty,
		// 					sm.MaxQuantity maxQty,
		// 					sm.Isin isin,
		// 					sm.MinPrice minprice,
		// 					sm.MaxPrice maxprice,
		// 					sm.BiddingStartDate startDate,
		// 					concat(
		// 						DATE_FORMAT(sm.BiddingStartDate, '%d %b %y'),' -  ',
		// 						DATE_FORMAT(sm.BiddingEndDate, '%d %b %y')
		// 						)
		// 						 as endDate,
		// 					(case
		// 						when date_sub(sm.BiddingStartDate,interval 1 day ) = curdate() then 'P'
		// 						when sm.BiddingStartDate > curdate() then 'U'
		// 					else 'C' end) Upcoming,
		// 					CONCAT( case
		// 					WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
		// 					WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
		// 					WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
		// 					ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
		// 					end,' ',
		// 					DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
		// 					TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime,
		// 					(case when  sm.BiddingEndDate = Date(now()) and Time(now()) > '15:30:00' then 1 else 0 end ) as LastDay
		// 					from a_sgb_master sm
		// 					where sm.BiddingEndDate >= curdate()
		// 					and sm.Exchange = 'NSE'
		// 					and sm.Redemption = 'N'
		// 					and not exists (
		// 					select 1
		// 					from a_sgb_master m
		// 					where m.BiddingEndDate = date(now())
		// 					and m.id = sm.id
		// 					and m.DailyEndTime <= time(now()))
		// 					) master
		// 				LEFT JOIN a_sgb_orderheader h
		// 				on master.Id = h.MasterId
		// 				and h.CancelFlag = 'N'
		// 				and h.ClientId = ?
		// 				and h.brokerId = ?
		// 				and h.Status is not null) tab
		// 			on tab.headerid = d.HeaderId
		// 			group by tab.id
		// 			order by (case when Flag = 'Y' then 1
		// 			when Upcoming = 'C' then 2
		// 			else 3 end),tab.startDate,tab.Symbol`

		lCoreString := `select tab.id,tab.Symbol,tab.Name,tab.minBidqty,
		tab.maxQty,tab.isin,tab.minprice,tab.maxprice,
		tab.endDate,tab.formatted_datetime,tab.LastDay,
		tab.Upcoming,tab.Flag,tab.Pending,lower(tab.Status),
		nvl((case when tab.headerId = d.HeaderId and tab.Flag = 'N' then "" else d.ReqSubscriptionUnit end) ,"")unit,
		nvl(d.ReqOrderNo,''),tab.startTime,tab.endTime
		 from a_sgb_orderdetails d
			right join (
			select nvl(h.Id,0) headerId,master.id,master.Symbol,
			master.Name,master.minBidqty,master.maxQty,master.isin,
			master.minprice,master.maxprice,master.startDate,
			master.endDate,master.Upcoming,master.formatted_datetime,master.LastDay,
			(case when h.MasterId = master.Id and h.CancelFlag = 'N' and h.Status = 'success' then 'Y' else 'N' end) Flag,
			(case when h.MasterId = master.Id and h.status = "pending" then "P" else "-" end) Pending,
			(case when h.MasterId = master.Id and h.Status = "success" then h.Status else "" end) Status,
			master.startTime,master.EndTime
			 from (
			 select sm.id Id,sm.Symbol Symbol,sm.Name Name,sm.MinBidQuantity minBidqty,
			 sm.MaxQuantity maxQty,sm.Isin isin,sm.MinPrice minprice,sm.MaxPrice maxprice,
			sm.BiddingStartDate startDate,
			concat(
					DATE_FORMAT(sm.BiddingStartDate, '%d %b %y'),' -  ',
					DATE_FORMAT(sm.BiddingEndDate, '%d %b %y')
				) as endDate,
			(case
			 when date_sub(sm.BiddingStartDate,interval 1 day ) = curdate() then 'P'
			 when sm.BiddingStartDate > curdate() then 'U'
			 else 'C' end) Upcoming,
				CONCAT( case
				 WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
				WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
				WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
				ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
				end,' ',
				DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
				TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime,
			(case when  sm.BiddingEndDate = Date(now()) and Time(now()) > '` + lCloseTime + `' then 1 else 0 end ) as LastDay,
			sm.DailyStartTime startTime,sm.DailyEndTime endTime
				from a_sgb_master sm
				where sm.BiddingEndDate >= curdate()
				and sm.Exchange = 'NSE'
				and sm.Redemption = 'N'
				and not exists (
				select 1
				from a_sgb_master m
				where m.BiddingEndDate = date(now())
				and m.id = sm.id
				and m.DailyEndTime <= time(now()))
				) master
			LEFT JOIN a_sgb_orderheader h
			on master.Id = h.MasterId
			and h.CancelFlag = 'N'
			and h.ClientId = ?
			and h.brokerId = ?
			and h.Status is not null) tab
		on tab.headerid = d.HeaderId
		group by tab.id
		order by (case when Flag = 'Y' then 1
		when Upcoming = 'C' then 2
		else 3 end),tab.startDate,tab.Symbol`

		// ==============================================================================

		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
		if lErr2 != nil {
			log.Println("GSGSD02", lErr2)
			return lSgbMasterArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSgbMasterRec.Id, &lSgbMasterRec.Symbol, &lSgbMasterRec.Name, &lSgbMasterRec.MinBidQty, &lSgbMasterRec.MaxBidQty, &lSgbMasterRec.Isin, &lSgbMasterRec.MinPrice, &lSgbMasterRec.MaxPrice, &lSgbMasterRec.CloseDate, &lSgbMasterRec.DateTime, &lSgbMasterRec.LastDay, &lSgbMasterRec.Upcoming, &lSgbMasterRec.Flag, &lSgbMasterRec.Pending, &lSgbMasterRec.Status, &lUnit, &lSgbMasterRec.OrderNo, &lSgbMasterRec.StartTime, &lSgbMasterRec.EndTime)
				if lErr3 != nil {
					log.Println("GSGSD03", lErr3)
					return lSgbMasterArr, lErr3
				} else {
					lSgbMasterRec.StartTime = lStartTime
					lSgbMasterRec.EndTime = lEndTime
					lSgbMasterRec.CloseTime = lCloseTime
					lSgbMasterRec.Unit, _ = strconv.Atoi(lUnit)

					lConfigFile := common.ReadTomlConfig("toml/debug.toml")
					lMaxQty := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_MAXQUANTITY"])
					if lMaxQty != "" {
						lSgbMasterRec.MaxBidQty, _ = strconv.Atoi(lMaxQty)
					}

					// Append Upi End Point in lRespRec.UpiArr array
					lSgbMasterArr = append(lSgbMasterArr, lSgbMasterRec)
				}
			}
			// log.Println(lSgbMasterArr)
		}

	}
	log.Println("GetSGBdetail (-)")
	return lSgbMasterArr, nil
}
