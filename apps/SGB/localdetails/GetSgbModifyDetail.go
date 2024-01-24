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

type SgbModifyStruct struct {
	Id          int    `json:"id"`
	Symbol      string `json:"name"`
	BidId       string `json:"bidId"`
	OrderNo     string `json:"orderNo"`
	Isin        string `json:"isin"`
	Unit        int    `json:"unit"`
	Price       int    `json:"price"`
	DateTime    string `json:"dateTime"`
	Total       int    `json:"total"`
	ExchOrderNo string `json:"exchOrderNo"`
}

// Response Structure for GetSgbMaster API
type SgbModifyResp struct {
	SgbModify SgbModifyStruct `json:"sgbModify"`
	Status    string          `json:"status"`
	ErrMsg    string          `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Modify Sgb Details in our database table based on masterId
Parameters:

not Applicable

Response:

*On Sucess
=========

	{
		"sgbModify":
		{
			"id": 18,
			"bidId": 1,
			"orderNo": "20230822",
			"unit": "2023-06-30",
			"price": "1000 - 2000",
		},
		"status": "S",
		"errMsg": ""
	}

!On Error
========

	{
		"status": E,
		"errMsg": "Can't able to get the requested Data"
	}

Author: Nithish Kumar
Date: 22AUG2023
*/
func GetSgbModifyDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSgbModifyDetail (+)", r.Method)
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
	(w).Header().Set("Access-Control-Allow-Headers", "ID,ORDERNO,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// create the instance for SgbModifyResp
		var lRespRec SgbModifyResp

		lRespRec.Status = common.SuccessCode

		lMasterId := r.Header.Get("ID")
		lOrderNo := r.Header.Get("ORDERNO")

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
		if lErr1 != nil {
			log.Println("LGSMD01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSMD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSMD01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGSMD02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lConvId, _ := strconv.Atoi(lMasterId)
		log.Println("lClientId, lConvId, lOrderNo, lBrokerId", lClientId, lConvId, lOrderNo, lBrokerId)

		lRespStruct, lErr2 := GetModifyDetail(lClientId, lConvId, lOrderNo, lBrokerId)
		if lErr2 != nil {
			log.Println("LGSMD03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGSMD03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGSMD03", "Error Occur in getting Datas.."))
			return
		} else {
			lRespRec.SgbModify = lRespStruct
			// log.Println("lRespRec", lRespRec.SgbModify)
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGSMD04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGSMD04", "Can't able to get the requested Data"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetSgbModifyDetail (-)", r.Method)
	}
}

func GetModifyDetail(pClientId string, pMasterId int, pOrderNo string, pBrokerId int) (SgbModifyStruct, error) {
	log.Println("GetModifyDetail (+)")

	var lSgbModifyRec SgbModifyStruct
	var lUnit, lPrice string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGMD01", lErr1)
		return lSgbModifyRec, lErr1
	} else {
		defer lDb.Close()
		//COMMENTED BY NITHISH BECAUSE THE ORDERNO COLUMN WAS CAHNAGED AND ALSO NOI NEEDED CLIENTID

		// lCoreString := `select sh.Id,sm.Symbol,sd.BidId  ,sd.OrderNo ,sd.ReqSubscriptionunit ,sd.ReqRate ,sm.Isin isin,
		// 	CONCAT( case
		// 		WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
		// 		WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
		// 		WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
		// 		ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
		// 		end,' ',
		// 		DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
		// 		TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime ,
		// 		(sd.ReqSubscriptionunit * sd.ReqRate)
		// 	from a_sgb_orderdetails sd,a_sgb_orderheader sh, a_sgb_master sm
		// 	where sh.Id = sd.HeaderId
		// 	and sm.Id = sh.MasterId
		//	and sh.clientId = ?
		// 	and sh.MasterId = ?
		// 	and sd.OrderNo = ?
		// 	and sh.brokerId = ?`
		lCoreString := `select sh.Id,sm.Symbol,sd.BidId  ,sd.ReqOrderNo ,sd.ReqSubscriptionunit ,sd.ReqRate ,sm.Isin isin,
		CONCAT( case
			WHEN DAY(sm.BiddingEndDate) % 10 = 1 AND DAY(sm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(sm.BiddingEndDate), 'st')
			WHEN DAY(sm.BiddingEndDate) % 10 = 2 AND DAY(sm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(sm.BiddingEndDate), 'nd')
			WHEN DAY(sm.BiddingEndDate) % 10 = 3 AND DAY(sm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(sm.BiddingEndDate), 'rd')
			ELSE CONCAT(DAY(sm.BiddingEndDate), 'th')
			end,' ',
			DATE_FORMAT(sm.BiddingEndDate, '%b %Y'),' | ',
			TIME_FORMAT(sm.DailyEndTime , '%h:%i%p')) AS formatted_datetime,
			sd.RespOrderNo 
		from a_sgb_orderdetails sd,a_sgb_orderheader sh, a_sgb_master sm
		where sh.Id = sd.HeaderId 
		and sm.Id = sh.MasterId 
		and sh.MasterId = ?
		and sd.ReqOrderNo = ?
		and sh.brokerId = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pOrderNo, pBrokerId)
		if lErr2 != nil {
			log.Println("LGMD02", lErr2)
			return lSgbModifyRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				var lOrderNo []uint8
				lErr3 := lRows.Scan(&lSgbModifyRec.Id, &lSgbModifyRec.Symbol, &lSgbModifyRec.BidId, &lOrderNo, &lUnit, &lPrice, &lSgbModifyRec.Isin, &lSgbModifyRec.DateTime, &lSgbModifyRec.ExchOrderNo)
				if lErr3 != nil {
					log.Println("LGMD03", lErr3)
					return lSgbModifyRec, lErr3
				} else {
					lSgbModifyRec.OrderNo = string(lOrderNo)

					lSgbModifyRec.Unit, lErr2 = strconv.Atoi(lUnit)
					if lErr2 != nil {
						log.Println("LGMD04", lErr2)
						return lSgbModifyRec, lErr2
					}
					lSgbModifyRec.Price, lErr2 = strconv.Atoi(lPrice)
					if lErr2 != nil {
						log.Println("LGMD05", lErr2)
						return lSgbModifyRec, lErr2
					}
					lSgbModifyRec.Total = lSgbModifyRec.Unit * lSgbModifyRec.Price

				}
			}
		}
	}
	log.Println("GetModifyDetail (-)")
	return lSgbModifyRec, nil
}
