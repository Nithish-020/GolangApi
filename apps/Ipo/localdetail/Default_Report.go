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
)

/*
Purpose: this method is used to get reports
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
		=========``
		{
			"status": "E",
			"errMsg": "Can't able to get data from database"
		}

Author: Nithish Kumar
Date: 07JUNE2023
*/
func DefaultReport(w http.ResponseWriter, r *http.Request) {
	log.Println("DefaultReport (+)", r.Method)
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

		// This variable helps to store the response and send it to front
		var lDRespRec RespStruct

		lDRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/report")
		if lErr1 != nil {
			log.Println("LDR01", lErr1)
			lDRespRec.Status = common.ErrorCode
			lDRespRec.ErrMsg = "LDR01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LDR01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LDR02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// Read the Request body in lBody
		lOutPutResp, lErr2 := GetReportDefault(lBrokerId)
		if lErr2 != nil {
			log.Println("LDR02", lErr2)
			lDRespRec.Status = common.ErrorCode
			lDRespRec.ErrMsg = "LDR01" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LDR01", "UserDetails not Found"))
			return
		} else {
			lDRespRec = lOutPutResp
		}

		ldata, lErr3 := json.Marshal(lDRespRec)
		if lErr3 != nil {
			log.Println("LDR03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LDR06", "Issue in Getting your Reports!"))
			return
		} else {
			fmt.Fprintf(w, string(ldata))
		}
		log.Println("DefaultReport (-)", r.Method)
	}
}

func GetReportDefault(pBrokerId int) (RespStruct, error) {
	log.Println("GetReportDefault (+)")

	var lIpoRec, lSgbRec ReportRespStruct
	var lDRespRec RespStruct

	lDRespRec.Status = common.SuccessCode

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGRD01", lErr1)
		return lDRespRec, lErr1
	} else {
		defer lDb.Close()

		lCoreString1 := `SELECT h.MasterId,h.Symbol,h.applicationNo,nvl(DATE_FORMAT(h.CreatedDate , '%d-%b-%Y'),'') applyDate,nvl(TIME_FORMAT(h.CreatedDate,'%h:%i:%s %p'),'') applytime,(case when h.cancelFlag = 'Y' then 'user cancelled' else nvl(h.status ,'') end) flag,h.clientId,nvl(h.exchange,''),nvl(category,'')
		from a_ipo_order_header h 
		where h.CreatedDate between concat (date(now()), ' 00:00:00.000') and concat (date(now()), ' 23:59:59.000')
		and h.brokerId = ?`
		lRows1, lErr4 := lDb.Query(lCoreString1, pBrokerId)
		if lErr4 != nil {
			log.Println("LGRD02", lErr4)
			return lDRespRec, lErr4
		} else {
			//Reading the Records
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lIpoRec.MasterId, &lIpoRec.Symbol, &lIpoRec.ApplicationNo, &lIpoRec.ApplyDate, &lIpoRec.AppliedTime, &lIpoRec.Status, &lIpoRec.ClientId, &lIpoRec.Exchange, &lIpoRec.Category)
				if lErr3 != nil {
					log.Println("LGRD03", lErr3)
					return lDRespRec, lErr3
				} else {
					lDRespRec.IpoArr = append(lDRespRec.IpoArr, lIpoRec)
				}

			}
			// COMMENTED BY NITHISH BECAUSE THE ORDER NO COLUMN WAS INCORRECT

			// lCoreString2 := `select oh.MasterId,m.Symbol ,od.OrderNo,nvl(DATE_FORMAT(oh.CreatedDate , '%d-%b-%Y'),'') applyDate,
			// nvl(TIME_FORMAT(oh.CreatedDate, '%h:%i:%s %p'),'') applytime,
			// (case when oh.cancelFlag = 'Y' then 'user cancelled' else nvl(oh.status ,'') end) flag,oh.clientId,nvl(oh.exchange,'')
			// from a_sgb_orderheader oh ,a_sgb_orderdetails od ,a_sgb_master m
			// where m.id = oh.MasterId and od.HeaderId =oh.Id
			// and od.CreatedDate between concat (date(now()), ' 00:00:00.000')
			// and concat (date(now()), ' 23:59:59.000')
			// and oh.brokerId = ?`

			lCoreString2 := `select oh.MasterId,m.Symbol ,od.RespOrderNo,nvl(DATE_FORMAT(oh.CreatedDate , '%d-%b-%Y'),'') applyDate,
			nvl(TIME_FORMAT(oh.CreatedDate, '%h:%i:%s %p'),'') applytime,
			(case when oh.cancelFlag = 'Y' then 'user cancelled' else nvl(oh.status ,'') end) flag,oh.clientId,nvl(oh.exchange,'')
			from a_sgb_orderheader oh ,a_sgb_orderdetails od ,a_sgb_master m 
			where m.id = oh.MasterId and od.HeaderId =oh.Id 
			and od.CreatedDate between concat (date(now()), ' 00:00:00.000') 
			and concat (date(now()), ' 23:59:59.000')
			and oh.brokerId = ?`
			lRows2, lErr4 := lDb.Query(lCoreString2, pBrokerId)
			if lErr4 != nil {
				log.Println("LGRD04", lErr4)
				return lDRespRec, lErr4
			} else {
				//Reading the Records
				for lRows2.Next() {
					lErr5 := lRows2.Scan(&lSgbRec.MasterId, &lSgbRec.Symbol, &lSgbRec.ApplicationNo, &lSgbRec.ApplyDate, &lSgbRec.AppliedTime, &lSgbRec.Status, &lSgbRec.ClientId, &lSgbRec.Exchange)
					if lErr5 != nil {
						log.Println("LGRD05", lErr5)
						return lDRespRec, lErr5
					} else {
						lDRespRec.SgbArr = append(lDRespRec.SgbArr, lSgbRec)
					}
				}
			}
		}
	}
	log.Println("GetReportDefault (-)")
	return lDRespRec, nil
}
