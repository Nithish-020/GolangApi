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

type SgbHistoryStruct struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ReqOrderNo       string `json:"reqOrderNo"`
	OrderNo          string `json:"orderNo"`
	OrderDate        string `json:"orderDate"`
	Isin             string `json:"isin"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	DateTime         string `json:"dateTime"`
	Unit             int    `json:"unit"`
	Price            int    `json:"price"`
	Total            int    `json:"total"`
	Flag             string `json:"flag"`
	Status           string `json:"status"`
	Subscriptionunit int    `json:"subscriptionunit"`
	BlockedAmount    int    `json:"blockedAmount"`
	ToolTip          bool   `json:"toolTip"`
}

// Response Structure for GetSgbMaster API
type SgbHistoryResp struct {
	SgbHistoryArr []SgbHistoryStruct `json:"sgbHistoryArr"`
	Status        string             `json:"status"`
	ErrMsg        string             `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Active Ipo Details in our database table ....
Parameters:

not Applicable

Response:

*On Sucess
=========

	{
		"sgbHistoryArr": [
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
Date: 22AUG2023
*/
func GetSgbHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSgbHistory (+)", r.Method)
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
		var lRespRec SgbHistoryResp

		lRespRec.Status = common.SuccessCode

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

		lRespArr, lErr2 := GetSGBHistorydetail(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("LGSH03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSH03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSH03", "Error Occur in getting Datas.."))
			return
		} else {
			lRespRec.SgbHistoryArr = lRespArr
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
		log.Println("GetSgbHistory (-)", r.Method)
	}
}

func GetSGBHistorydetail(pClientId string, pBrokerId int) ([]SgbHistoryStruct, error) {
	log.Println("GetSGBHistorydetail (+)")

	var lSgbHistoryRec SgbHistoryStruct
	var lSgbHistoryArr []SgbHistoryStruct
	var lUnit, lPrice, lRespUnit, lRespRate string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGSHD01", lErr1)
		return lSgbHistoryArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := ` select h.Id Id,
						sm.Name Name,
						d.ReqOrderNo,
						d.RespOrderNo,
						date_format(h.CreatedDate, '%d-%b-%y, %l:%i %p'),
						sm.Isin isin,
						sm.BiddingStartDate startDate,
						sm.BiddingEndDate endDate,
						CONCAT( case
						WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
						WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
						WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
						ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
						end,' ',
						DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
						TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime ,
						d.ReqSubscriptionunit ,
						nvl(d.RespSubscriptionunit,0),
						d.ReqRate ,
						nvl(d.RespRate,0),
						(case when h.MasterId = sm.Id and h.CancelFlag = 'Y' and h.Status = 'success' then 'Y' else 'N' end) Flag,						
						lower(h.Status) Status						
					from a_sgb_master sm,a_sgb_orderheader h ,a_sgb_orderdetails d
					where sm.Id = h.MasterId 
					and d.HeaderId = h.Id 
					and h.ClientId = ?
					and h.brokerId = ?
					and h.Status is not null
					order by h.Id desc`

		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
		if lErr2 != nil {
			log.Println("LGSHD02", lErr2)
			return lSgbHistoryArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSgbHistoryRec.Id, &lSgbHistoryRec.Name, &lSgbHistoryRec.ReqOrderNo, &lSgbHistoryRec.OrderNo, &lSgbHistoryRec.OrderDate, &lSgbHistoryRec.Isin, &lSgbHistoryRec.StartDate, &lSgbHistoryRec.EndDate, &lSgbHistoryRec.DateTime, &lUnit, &lRespUnit, &lPrice, &lRespRate, &lSgbHistoryRec.Flag, &lSgbHistoryRec.Status)
				if lErr3 != nil {
					log.Println("LGSHD03", lErr3)
					return lSgbHistoryArr, lErr3
				} else {
					lSgbHistoryRec.ToolTip = false
					lSgbHistoryRec.Unit, _ = strconv.Atoi(lUnit)
					lSgbHistoryRec.Price, _ = strconv.Atoi(lPrice)
					lSgbHistoryRec.Subscriptionunit, _ = strconv.Atoi(lRespUnit)
					lRespPrice, _ := strconv.Atoi(lRespRate)
					lSgbHistoryRec.BlockedAmount = lSgbHistoryRec.Subscriptionunit * lRespPrice
					lSgbHistoryRec.Total = lSgbHistoryRec.Unit * lSgbHistoryRec.Price
					// Append Upi End Point in lRespRec.SgbHistoryArr array
					lSgbHistoryArr = append(lSgbHistoryArr, lSgbHistoryRec)
				}
			}
			// log.Println(lSgbHistoryArr)
		}
	}
	log.Println("GetSGBHistorydetail (-)")
	return lSgbHistoryArr, nil
}
