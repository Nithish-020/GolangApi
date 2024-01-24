package localdetail

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

// This Structure is used to Collect the Active IPO informations
type ActiveIpoStruct struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Symbol         string  `json:"symbol"`
	BidDate        string  `json:"bidDate"`
	EndDate        string  `json:"endDate"`
	PriceRange     string  `json:"priceRange"`
	MinPrice       float64 `json:"minPrice"`
	CutOffPrice    float64 `json:"cutOffPrice"`
	MinBidQuantity int     `json:"minBidQuantity"`
	LotSize        float64 `json:"lotSize"`
	IssueSize      float64 `json:"issueSize"`
	CutOffFlag     string  `json:"cutOffFlag"`
	Flag           string  `json:"flag"`
	Pending        string  `json:"pending"`
	Upcoming       string  `json:"upcoming"`
	BlogLink       string  `json:"blogLink"`
	DrhpLink       string  `json:"drhpLink"`
	Sme            bool    `json:"sme"`
	PreApply       string  `json:"preApply"`
	Category       string  `json:"category"`
	Code           string  `json:"code"`
	ApplicationNo  string  `json:"applicationNo"`
}

// Response Structure for GetIpoMaster API
type IpoStruct struct {
	IpoDetails []ActiveIpoStruct `json:"ipoDetail"`
	Status     string            `json:"status"`
	ErrMsg     string            `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Active Ipo Details in our database table ....
Parameters:

not Applicable

Response:

*ON Sucess
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
				"upiStatus": "Accepted BY Investor"
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

!ON Error
========

	{
		"status": E,
		"reason": "Can't able to get the data FROM database"
	}

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetIpoMaster(w http.ResponseWriter, r *http.Request) {
	log.Println("GetIpoMaster(+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	// var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, _ = brokers.GetBrokerId(origin) // TO get brokerId
			// log.Println(lErr, origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// create the instance for IpoStruct
		var lRespRec IpoStruct
		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("LGIM01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGIM02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lIpoMasterArr, lErr2 := GetIpoMasterDetail(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("LGIM02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIM02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIM02", "Unable to fetch the master details"))
			return
		} else {
			lRespRec.IpoDetails = lIpoMasterArr
		}

		// Marshaling the response structure to lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGIM03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGIM03", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetIpoMaster (-)", r.Method)
	}
}

func GetIpoMasterDetail(pClientId string, pBrokerId int) ([]ActiveIpoStruct, error) {
	log.Println("GetIpoMasterDetail (+)")
	// create the instance for activeIpoStruct
	var lIpoDataRec ActiveIpoStruct
	var lIpoDataArr []ActiveIpoStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMD01", lErr1)
		return lIpoDataArr, lErr1

	} else {
		defer lDb.Close()

		lBrokerId := strconv.Itoa(pBrokerId)

		lCoreString := `SELECT tab.Id Id,tab.Name Name,tab.Symbol Symbol,tab.BidDates,
		tab.EndDate,tab.PriceRange PriceRange,tab.MinPrice,tab.MaxPrice MaxPrice,
		tab.MinQty MinQty,tab.Issue Issue,tab.LotSize Lot,tab.AllowCutOff AllowCutOff,
		(CASE WHEN d.activityType = 'cancel' || tab.Flag = 'N' THEN "N" else "Y" END ) Flag,
		tab.Pending Pending,tab.Upcoming,tab.dLink,tab.drhpLink,tab.SubType,
		(CASE WHEN curdate() = date_sub(tab.StartDate,interval 1 day) THEN 'pre' else 'NA' END) Apply,
			(
				SELECT
				CASE
					WHEN COUNT(*) = 1 THEN MAX(NVL(
					CASE WHEN vs.SubCatCode = "IND" THEN "Individual Investor"
					WHEN vs.SubCatCode = "EMP" THEN "Employee"
				   else "Existing Shareholder" END
						, ''))
						ELSE ''
				END AS code
				FROM v_ipo_subcategory vs
				where vs.id = tab.id
			) category,
			(
				SELECT
				case WHEN COUNT(*) = 1 THEN MAX(NVL(vs.SubCatCode,''))
					else '' end
				FROM v_ipo_subcategory vs
				where vs.id = tab.id
			)code,
			(
				SELECT
				CASE
					WHEN COUNT(*) = 1 THEN MAX(NVL(tab.applicationNo, ''))ELSE ''
				END AS appNo
				FROM v_ipo_subcategory vs
				where vs.id = tab.id
			) applicationNo
			FROM a_ipo_orderdetails d RIGHT JOIN
		(
		SELECT NVL(oh.Id,0) HeadId,	mas.Id Id,mas.Name Name,mas.Symbol Symbol,mas.StartDate,mas.BidDates,
		mas.EndDate,mas.PriceRange PriceRange,mas.MinPrice,mas.MaxPrice MaxPrice,mas.MinQty MinQty,
		mas.Issue Issue,mas.LotSize LotSize,mas.AllowCutOff AllowCutOff,
		(SELECT (CASE WHEN oh.MasterId = mas.Id AND oh.cancelFlag = "N" AND oh.status =  "success" THEN "Y" else "N" END)) Flag,
		(SELECT (CASE WHEN oh.MasterId = mas.Id AND oh.status =  "pending" THEN "P" else "-" END)) Pending,
		mas.Upcoming,mas.dLink,mas.drhpLink,mas.SubType,mas.Exchange,oh.applicationNo
			FROM 
		(
			SELECT im.Id,im.Name, im.Symbol,im.StartDate,im.BidDates ,im.EndDate,im.PriceRange ,
			im.MinPrice,im.MaxPrice ,im.MinQty,im.Issue,im.LotSize,im.SubType,im.Exchange,im.AllowCutOff,
			im.Upcoming,im.dLink,im.drhpLink 
				FROM (
					SELECT m.Id Id,m.Name Name,m.Symbol Symbol,m.BiddingStartDate StartDate,
					CONCAT(DATE_FORMAT(m.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(m.BiddingEndDate , '%d %b')) AS BidDates,
					m.BiddingEndDate EndDate,concat(m.MinPrice," - ", m.MaxPrice) PriceRange,m.MinPrice MinPrice,
					m.MaxPrice MaxPrice,m.MinBidQuantity MinQty,m.IssueSize Issue,m.LotSize LotSize,
					NVL(
						(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.category else 1 END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),1
					)SubType,
					m.Exchange Exchange,
					(CASE WHEN m.Id = s.MasterId AND s.AllowCutOff = 1 THEN "Y" else "N" END)AllowCutOff,
					(CASE WHEN m.BiddingStartDate > curdate() THEN 'U'else 'C' END)Upcoming,
					NVL(
						(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.detailsLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
					)dLink,
					NVL(
						(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.drhpLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
					)drhpLink											
				FROM a_ipo_master m ,a_ipo_categories c,a_ipo_subcategory s
					WHERE m.Id= s.MasterId AND m.Id = c.MasterId AND m.IssueType = "EQUITY"
					AND c.code= "RETAIL" AND s.CaCode = "RETAIL" AND s.SubCatCode = "IND"
					AND s.AllowUpi = 1 AND m.BiddingEndDate >= Curdate() AND m.Exchange = 'NSE'
					AND not exists (
						SELECT 1
						FROM(
							SELECT NVL(c.EndTime, m.DailyEndTime) endTime, c.MasterId
							FROM a_ipo_master m,a_ipo_categories c
							WHERE m.Id = c.MasterId
							AND c.Code = 'RETAIL'
							AND m.BiddingEndDate = date(now())
							) a, a_ipo_master m1
						WHERE m1.Id = a.masterId
						AND a.masterId = m.Id
						AND a.endTime <= time(now())
						)
				union 
					SELECT m.Id Id,m.Name Name,m.Symbol Symbol,m.BiddingStartDate StartDate,
					CONCAT(DATE_FORMAT(m.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(m.BiddingEndDate , '%d %b')) AS BidDates,
					m.BiddingEndDate EndDate,concat(m.MinPrice," - ", m.MaxPrice) PriceRange,m.MinPrice MinPrice,
					m.MaxPrice MaxPrice,m.MinBidQuantity MinQty,m.IssueSize Issue,m.LotSize LotSize,
					NVL(
						(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.category else 1 END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),1
					)SubType,
					m.Exchange Exchange,
					(CASE WHEN m.Id = s.MasterId AND s.AllowCutOff = 1 THEN "Y" else "N" END)AllowCutOff,
					(CASE WHEN m.BiddingStartDate > curdate() THEN 'U'else 'C' END)Upcoming,
					NVL(
						(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.detailsLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
					)dLink,
					NVL(
						(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.drhpLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
					)drhpLink									
				FROM a_ipo_master m ,a_ipo_categories c,a_ipo_subcategory s
					WHERE m.Id= s.MasterId AND m.Id = c.MasterId AND m.IssueType = "EQUITY"
					AND c.code= "RETAIL" AND s.CaCode = "RETAIL" AND s.SubCatCode = "IND"
					AND s.AllowUpi = 1 AND m.BiddingEndDate >= Curdate() AND m.Exchange = 'BSE'
					AND not exists (
						SELECT 1
						FROM (
							SELECT NVL(c.EndTime, m.DailyEndTime) endTime, c.MasterId
							FROM a_ipo_master m,a_ipo_categories c
							WHERE m.Id = c.MasterId
							AND c.Code = 'RETAIL'
							AND m.BiddingEndDate = date(now())
							) a, a_ipo_master m1
						WHERE m1.Id = a.masterId
						AND a.masterId = m.Id
						AND a.endTime <= time(now())
						)
					) im
				)mas 
			LEFT JOIN a_ipo_order_header oh 
			ON mas.Id = oh.MasterId AND oh.clientId = '` + pClientId + `' AND oh.brokerId = ` + lBrokerId + `
			AND oh.status is not null AND oh.status <> 'failed' AND oh.cancelFlag = 'N'
			GROUP BY mas.Symbol
		) tab
		ON tab.HeadId = d.headerId 
		GROUP BY tab.Id
		ORDER BY (CASE
		WHEN Flag = 'Y' THEN 1
		WHEN Pending = 'P' THEN 2
		WHEN Upcoming = 'U'  THEN 3
		else 4 END) ,tab.StartDate desc`

		lRows, lErr3 := lDb.Query(lCoreString)
		if lErr3 != nil {
			log.Println("LGIMD03", lErr3)
			return lIpoDataArr, lErr3

		} else {
			//This for loop is used to collect the records FROM the database AND store them in structure
			for lRows.Next() {
				lErr4 := lRows.Scan(&lIpoDataRec.Id, &lIpoDataRec.Name, &lIpoDataRec.Symbol, &lIpoDataRec.BidDate, &lIpoDataRec.EndDate, &lIpoDataRec.PriceRange, &lIpoDataRec.MinPrice, &lIpoDataRec.CutOffPrice, &lIpoDataRec.MinBidQuantity, &lIpoDataRec.IssueSize, &lIpoDataRec.LotSize, &lIpoDataRec.CutOffFlag, &lIpoDataRec.Flag, &lIpoDataRec.Pending, &lIpoDataRec.Upcoming, &lIpoDataRec.BlogLink, &lIpoDataRec.DrhpLink, &lIpoDataRec.Sme, &lIpoDataRec.PreApply, &lIpoDataRec.Category, &lIpoDataRec.Code, &lIpoDataRec.ApplicationNo)
				// log.Println("lRows", lIpoDataRec)
				if lErr4 != nil {
					log.Println("LGIMD04", lErr4)
					return lIpoDataArr, lErr4

				} else {
					// Append the IPO Records in lRespRec.IpoDetails Array
					lIpoDataArr = append(lIpoDataArr, lIpoDataRec)
				}
			}
		}
	}
	log.Println("GetIpoMasterDetail (-)")
	return lIpoDataArr, nil
}
