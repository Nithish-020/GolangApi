package ncblocaldetails

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type NcbHistoryStruct struct {
	Id            int    `json:"id"`
	Symbol        string `json:"symbol"`
	Series        string `json:"series"`
	Name          string `json:"name"`
	ApplicationNo string `json:"applicationNo"`
	OrderNo       string `json:"orderNo"`
	OrderDate     string `json:"orderDate"`
	Isin          string `json:"isin"`
	StartDate     string `json:"startDate"`
	EndDateTime   string `json:"dateTime"`
	Unit          int    `json:"unit"`
	Price         int    `json:"price"`
	Total         int    `json:"total"`
	Flag          string `json:"flag"`
	// CancelFlag    string `json:"cancelFlag"`
	Status string `json:"status"`
}

// Response Structure for GetSgbMaster API
type NcbHistoryResp struct {
	GSecHisDetails  []NcbHistoryStruct `json:"gSecHisDetail"`
	TBillHisDetails []NcbHistoryStruct `json:"tBillHisDetail"`
	SdlHisDetails   []NcbHistoryStruct `json:"sdlHisDetail"`
	Status          string             `json:"status"`
	ErrMsg          string             `json:"errMsg"`
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

func GetNcbHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNcbHistory(+)", r.Method)

	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// create the instance for IpoStruct
		var lRespRec NcbHistoryResp

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
		if lErr1 != nil {
			log.Println("NGNH01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGNH01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGNH01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("NGNH02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lGsecRespHisArr, lTbilRespHisArr, lSdlRespHisArr, lErr2 := GetNcbHistorydetail(lClientId)
		if lErr2 != nil {
			log.Println("NGNH03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGNH03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGNH03", "Error Occur in getting Datas.."))
			return
		} else {

			if lGsecRespHisArr != nil || lTbilRespHisArr != nil || lSdlRespHisArr != nil {
				lRespRec.GSecHisDetails = lGsecRespHisArr
				log.Println("lRespRec.GSecHisDetails", lRespRec.GSecHisDetails)
				lRespRec.TBillHisDetails = lTbilRespHisArr
				log.Println("lRespRec.TBillHisDetails", lRespRec.TBillHisDetails)
				lRespRec.SdlHisDetails = lSdlRespHisArr
				log.Println("lRespRec.SdlHisDetails", lRespRec.SdlHisDetails)
			}
			// lRespRec.SgbHistoryArr = lRespArr
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NGNH04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("NGNH04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetNcbHistory(-)", r.Method)
	}
}

func GetNcbHistorydetail(pClientId string) ([]NcbHistoryStruct, []NcbHistoryStruct, []NcbHistoryStruct, error) {
	log.Println("GetNcbHistorydetail(+)")

	var lNcbHistoryRec NcbHistoryStruct

	var lGsecHistoryArr []NcbHistoryStruct
	var lTbillHistoryArr []NcbHistoryStruct
	var lSdlHistoryArr []NcbHistoryStruct

	var lUnit, lPrice string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NGNHD01", lErr1)
		return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select nvl(h.Id,'') id, nvl(h.ApplicationNo,'') ,n.Name Name, nvl(d.price,'') price ,n.Symbol, n.Series,nvl(DATE_FORMAT(h.CreatedDate , '%d-%b-%Y'),'') CreatedBy , CONCAT( case WHEN DAY(n.BiddingEndDate) % 10 = 1 AND DAY(n.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(n.BiddingEndDate), 'st')
				       WHEN DAY(n.BiddingEndDate) % 10 = 2 AND DAY(n.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(n.BiddingEndDate), 'nd')
				       WHEN DAY(n.BiddingEndDate) % 10 = 3 AND DAY(n.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(n.BiddingEndDate), 'rd')
				       ELSE CONCAT(DAY(n.BiddingEndDate), 'th')
	                   end,' ',
			              DATE_FORMAT(n.BiddingEndDate, '%b %Y'),' | ',
			              TIME_FORMAT(n.DailyEndTime , '%h:%i%p')) AS formatted_datetime,n.BiddingStartDate,
				          (case when h.MasterId = n.Id and h.CancelFlag = 'Y' and h.Status = 'success' then 'Y' else 'N' end) Flag,						
				         lower(h.Status) Status	
                       from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
                       where n.id  = h.MasterId  and d.HeaderId = h.Id and h.ClientId = ?  
                       group by h.applicationNo 
                     order by h.Id desc`

		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
		if lErr2 != nil {
			log.Println("NGNHD02", lErr2)
			return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {

				lErr3 := lRows.Scan(&lNcbHistoryRec.Id, &lNcbHistoryRec.ApplicationNo, &lNcbHistoryRec.Name, &lPrice, &lNcbHistoryRec.Symbol, &lNcbHistoryRec.Series, &lNcbHistoryRec.OrderDate, &lNcbHistoryRec.EndDateTime, &lNcbHistoryRec.StartDate, &lNcbHistoryRec.Flag, &lNcbHistoryRec.Status)
				if lErr3 != nil {
					log.Println("NGNHD03", lErr3)
					return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, lErr3
				} else {
					lNcbHistoryRec.Unit, _ = strconv.Atoi(lUnit)
					lNcbHistoryRec.Price, _ = strconv.Atoi(lPrice)

					lNcbHistoryRec.Total = lNcbHistoryRec.Unit * lNcbHistoryRec.Price
					// Append Upi End Point in lRespRec.SgbHistoryArr array

					if lNcbHistoryRec.Series == "GS" {
						lGsecHistoryArr = append(lGsecHistoryArr, lNcbHistoryRec)
						log.Println("lGsecHistoryArr--->GS", lGsecHistoryArr)
					} else if lNcbHistoryRec.Series == "TB" {
						lTbillHistoryArr = append(lTbillHistoryArr, lNcbHistoryRec)
						log.Println("lTbillHistoryArr--->TB", lTbillHistoryArr)
					} else {
						lSdlHistoryArr = append(lSdlHistoryArr, lNcbHistoryRec)
						log.Println("lSdlHistoryArr--->TB", lSdlHistoryArr)
					}

				}
			}

		}
	}
	log.Println("GetNcbHistorydetail(-)")
	return lGsecHistoryArr, lTbillHistoryArr, lSdlHistoryArr, nil
}
