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
)

type ActiveNcbStruct struct {
	Id              int     `json:"id"`
	Symbol          string  `json:"symbol"`
	Series          string  `json:"series"`
	Name            string  `json:"name"`
	MinBidQuantity  int     `json:"minBidQuantity"`
	MaxQuantity     string  `json:"maxQuantity"`
	TotalQuantity   string  `json:"totalQuantity"`
	Isin            string  `json:"isin"`
	MinPrice        float64 `json:"minPrice"`
	MaxPrice        float64 `json:"maxPrice"`
	CloseDate       string  `json:"closeDate"`
	DateTime        string  `json:"dateTime"`
	CutOffFlag      string  `json:"cutOffFlag"`
	ActivityType    string  `json:"activityType"`
	Flag            string  `json:"flag"`
	Unit            int     `json:"unit"`
	OrderNo         int     `json:"orderNo"`
	Status          string  `json:"status"`
	LotSize         float32 `json:"lotSize"`
	ModifiedLotSize int     `json:"modifiedLotSize"`
	Lotvalue        int     `json:"lotValue"`
	FaceValue       float64 `json:"faceValue"`
	CutoffPrice     float64 `json:"cutoffPrice"`
	Amount          float64 `json:"amount"`
}

type NcbSeriesStruct struct {
	GSecDetails  []ActiveNcbStruct `json:"gSecDetail"`
	TBillDetails []ActiveNcbStruct `json:"tBillDetail"`
	SdlDetails   []ActiveNcbStruct `json:"sdlDetail"`
	Status       string            `json:"status"`
	ErrMsg       string            `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Active NcbDetailStruct in our database table ....
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

Author: KAVYA DHARSHANI
Date: 10OCT2023
*/

func GetNcbMaster(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNcbMaster(+)", r.Method)
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
		// var lRespRec NcbStruct
		var lRespRec NcbSeriesStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
		if lErr1 != nil {
			log.Println("NGNM01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGNM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGNM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("NGNM02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// lRespArr, lErr2 := GetNcbdetail(lClientId)
		lGsecRespArr, lTbilRespArr, lSdlRespArr, lErr2 := GetNcbdetail(lClientId)
		if lErr2 != nil {
			log.Println("NGNM03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGNM03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGNM03", "Error Occur in getting Datas.."))
			return
		} else {
			if lGsecRespArr != nil || lTbilRespArr != nil || lSdlRespArr != nil {
				lRespRec.GSecDetails = lGsecRespArr
				log.Println("lRespRec.GSecDetails", lRespRec.GSecDetails)
				lRespRec.TBillDetails = lTbilRespArr
				log.Println("lRespRec.TBillDetails", lRespRec.TBillDetails)
				lRespRec.SdlDetails = lSdlRespArr
				log.Println("lRespRec.SdlDetails", lRespRec.SdlDetails)
			}
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NGNM04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("NGNM04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetNcbMaster (-)", r.Method)
	}
}

func GetNcbdetail(pClientId string) ([]ActiveNcbStruct, []ActiveNcbStruct, []ActiveNcbStruct, error) {
	log.Println("GetNcbdetail (+)")

	var lNcbMasterRec ActiveNcbStruct
	// var lNcbMasterArr []ActiveNcbStruct

	var lNcbGsecArr []ActiveNcbStruct
	var lNcbTbillArr []ActiveNcbStruct
	var lNcbSdlArr []ActiveNcbStruct

	// var lUnit string
	// Calling LocalDbConect method in ftdb to estabish the database connection

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GNGND01", lErr1)
		return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := ` select tab.id, tab.Symbol,tab.Series,tab.Name,tab.minBidqty,CONVERT(tab.maxQty, SIGNED) AS maxQty, tab.isin, tab.minprice, tab.maxprice, tab.endDate,tab.formatted_datetime, 
		(case when d.activityType = 'cancel' || tab.Flag = 'N' then "N" else "Y" end ) Flag,tab.Lotsize, CAST((tab.Lotsize / 100) AS SIGNED) AS ModifiedLotsize,
		CAST((tab.Lotsize / 100) AS SIGNED) lotvalue, tab.FaceValue, tab.CutoffPrice,
		  CONCAT(tab.minBidqty, ' - ', CONVERT(tab.maxQty, SIGNED)) AS minBidMaxQty,
		lower(tab.status),	nvl(d.OrderNo,0), nvl(d.price,0) price,  nvl(d.activityType,''), nvl(d.Unit,0)
                   from a_ncb_orderdetails d
                         right join (select nvl(h.Id,0) headerId, master.id,
				            master.Symbol,master.Series, master.Name,master.minBidqty, master.maxQty,master.isin,master.minprice,master.maxprice, master.startDate,
				            master.endDate,master.formatted_datetime, master.Lotsize, master.FaceValue, master.CutoffPrice,
		            (case when h.MasterId = master.Id and h.CancelFlag = 'N' and h.status = 'success' then 'Y' else 'N' end) Flag,
		            (case when h.MasterId = master.Id and h.status = "success" then h.status else "" end) status
	               from (select nm.id Id,nm.Symbol Symbol,nm.Name Name,nm.Series Series,
	               nm.MinBidQuantity minBidqty,nm.MaxQuantity maxQty,	nm.Isin isin,nm.MinPrice minprice,nm.MaxPrice maxprice, nm.BiddingStartDate startDate,
	               nm.Lotsize , nm.FaceValue , nm.CutoffPrice ,
	        concat(DATE_FORMAT(nm.BiddingStartDate, '%d %b %y'),' -  ',DATE_FORMAT(nm.BiddingEndDate, '%d %b %y'))  as endDate,					
	         CONCAT( case WHEN DAY(nm.BiddingEndDate) % 10 = 1 AND DAY(nm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(nm.BiddingEndDate), 'st')
				          WHEN DAY(nm.BiddingEndDate) % 10 = 2 AND DAY(nm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(nm.BiddingEndDate), 'nd')
				           WHEN DAY(nm.BiddingEndDate) % 10 = 3 AND DAY(nm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(nm.BiddingEndDate), 'rd')
			          ELSE CONCAT(DAY(nm.BiddingEndDate), 'th')
						end,' ',
	          DATE_FORMAT(nm.BiddingEndDate, '%b %Y'),' | ',
	          TIME_FORMAT(nm.DailyEndTime , '%h:%i%p')) AS formatted_datetime
            from a_ncb_master nm
	                 where nm.BiddingEndDate >= curdate()  and nm.Exchange = 'NSE' and not exists (
	        select 1
	       from a_ncb_master n	
            where n.BiddingEndDate = date(now()) and n.id = nm.id and n.DailyEndTime <= time(now())) ) master
	        LEFT JOIN a_ncb_orderheader h 
				  on master.Id = h.MasterId
					and h.CancelFlag = 'N'
					and h.ClientId = ?
					and h.status is not null
				    and h.status <> 'failed'
				   and h.cancelFlag = 'N' 
				    group by master.Symbol) tab
					on tab.headerid = d.HeaderId
					group by tab.id,tab.startDate ,tab.Symbol`

		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
		if lErr2 != nil {
			log.Println("GNGND02", lErr2)
			return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lNcbMasterRec.Id, &lNcbMasterRec.Symbol, &lNcbMasterRec.Series, &lNcbMasterRec.Name, &lNcbMasterRec.MinBidQuantity, &lNcbMasterRec.MaxQuantity, &lNcbMasterRec.Isin, &lNcbMasterRec.MinPrice, &lNcbMasterRec.MaxPrice, &lNcbMasterRec.CloseDate, &lNcbMasterRec.DateTime, &lNcbMasterRec.Flag, &lNcbMasterRec.LotSize, &lNcbMasterRec.ModifiedLotSize, &lNcbMasterRec.Lotvalue, &lNcbMasterRec.FaceValue, &lNcbMasterRec.CutoffPrice, &lNcbMasterRec.TotalQuantity, &lNcbMasterRec.Status, &lNcbMasterRec.OrderNo, &lNcbMasterRec.Amount, &lNcbMasterRec.ActivityType, &lNcbMasterRec.Unit)
				log.Println("lNcbMasterRec", lNcbMasterRec.MaxQuantity)

				if lErr3 != nil {
					log.Println("GNGND03", lErr3)
					return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr3
				} else {
					// lNcbMasterRec.Unit, _ = strconv.Atoi(lUnit)

					// Append Upi End Point in lRespRec.UpiArr array
					// lNcbMasterArr = append(lNcbMasterArr, lNcbMasterRec)
					if lNcbMasterRec.Series == "GS" {
						lNcbGsecArr = append(lNcbGsecArr, lNcbMasterRec)
						log.Println("lNcbGsecArr--->GS", lNcbGsecArr)
					} else if lNcbMasterRec.Series == "TB" {
						lNcbTbillArr = append(lNcbTbillArr, lNcbMasterRec)
						log.Println("lNcbTbillArr--->TB", lNcbTbillArr)
					} else {
						lNcbSdlArr = append(lNcbSdlArr, lNcbMasterRec)
						log.Println("lNcbSdlArr--->TB", lNcbSdlArr)
					}

				}
			}
			// log.Println(lNcbMasterArr)
		}

	}
	log.Println("GetNcbdetail (-)")
	return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, nil
}
