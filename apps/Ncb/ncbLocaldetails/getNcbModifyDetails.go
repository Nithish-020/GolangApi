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
	"net/http"
	"strconv"
)

type NcbModifyStruct struct {
	Id            int    `json:"id"`
	Symbol        string `json:"symbol"`
	Flag          string `json:"flag"`
	OrderNo       int    `json:"orderNo"`
	ApplicationNo string `json:"applicationNo"`
	LotSize       int    `json:"lotSize"`
	Isin          string `json:"isin"`
	Unit          int    `json:"unit"`
	Price         int    `json:"price"`
	Amount        int    `json:"amount"`
	DateTime      string `json:"dateTime"`
	Total         int    `json:"total"`
	RespOrderNo   string `json:"respOrderNo"`
}

// Response Structure for GetNcbMaster API
type NcbModifyResp struct {
	NcbModify NcbModifyStruct `json:"NcbModify"`
	Status    string          `json:"status"`
	ErrMsg    string          `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Modify NCB Details in our database table based on masterId
Parameters:

not Applicable

Response:

*On Sucess
=========

lRespRec : {
	         "NcbModify":
		            {
                      	"id": 18,
			            "orderNo": "20230822",
			            unit": "2023-06-30",
			            "price": "1000 - 2000",
		            },
            }

!On Error
========

	{
		"status": E,
		"errMsg": "Can't able to get the requested Data"
	}

Author: KAVYA DHARSHANI M
Date: 21OCT2023
*/

func GetNcbModifyDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNcbModifyDetail(+)", r.Method)
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
		// create the instance for NcbModifyResp
		var lRespRec NcbModifyResp

		lRespRec.Status = common.SuccessCode

		lMasterId := r.Header.Get("ID")
		lOrderNo := r.Header.Get("ORDERNO")

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
		if lErr1 != nil {
			log.Println("NLGNMD01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NLGNMD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NLGNMD01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("NLGNMD02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lConvId, _ := strconv.Atoi(lMasterId)
		lIdOrder, _ := strconv.Atoi(lOrderNo)

		log.Println("Id,Orderno", lIdOrder, lConvId)

		lRespStruct, lErr2 := GetModifyDetails(lClientId, lConvId, lIdOrder, lBrokerId)
		if lErr2 != nil {
			log.Println("NLGNMD03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NLGNMD03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NLGNMD03", "Error Occur in getting Datas.."))
			return
		} else {
			lRespRec.NcbModify = lRespStruct
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NLGNMD04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("NLGNMD04", "Can't able to get the requested Data"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetNcbModifyDetail(-)", r.Method)
	}
}

/*
Pupose:This method used to retrieve the Application Details from the database.
Parameters:
PClientId,pMasterId
Response:
==========
*On Sucess
==========
lNcbModifyRec :
           {2 AP20492402 N N 0 FT032287164628 100 1000000 1.05e+08}

==========
!On Error
==========
[],error

Author:KAVYADHARSHANI
Date: 21OCT2023
*/
func GetModifyDetails(pClientId string, pMasterId int, pOrderNo int, pBrokerId int) (NcbModifyStruct, error) {
	log.Println("GetModifyDetails (+)")

	var lNcbModifyRec NcbModifyStruct

	// lOrderNo = lNcbModifyRec.OrderNo

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LMGD01", lErr1)
		return lNcbModifyRec, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select h.Id,n.Symbol,d.ReqOrderNo ,d.ReqUnit  ,d.Reqprice,d.ReqAmount, n.Isin isin,
		                       CONCAT( case 
								 WHEN DAY(n.BiddingEndDate) % 10 = 1 AND DAY(n.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(n.BiddingEndDate), 'st')
			                     WHEN DAY(n.BiddingEndDate) % 10 = 2 AND DAY(n.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(n.BiddingEndDate), 'nd')
			                     WHEN DAY(n.BiddingEndDate) % 10 = 3 AND DAY(n.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(n.BiddingEndDate), 'rd')
			                     ELSE CONCAT(DAY(n.BiddingEndDate), 'th')
			                     end,' ',
			                     DATE_FORMAT(n.BiddingEndDate, '%b %Y'),' | ',
			                  TIME_FORMAT(n.DailyEndTime , '%h:%i%p')) AS formatted_datetime, d.RespOrderNo  
                        from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
                        where h.Id = d.HeaderId 
                        and n.Id = h.MasterId
                        and h.MasterId  = ?
                        and d.ReqOrderNo = ?
                        and h.brokerId  = ?`

		log.Println("pClientId", pClientId, pMasterId, "pMasterId")
		lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pOrderNo, pBrokerId)
		if lErr2 != nil {
			log.Println("LMGD02", lErr2)
			return lNcbModifyRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in ModifyRespStruct
			for lRows.Next() {
				lErr3 := lRows.Scan(&lNcbModifyRec.Id, &lNcbModifyRec.Symbol, &lNcbModifyRec.OrderNo, &lNcbModifyRec.Unit, &lNcbModifyRec.Price, &lNcbModifyRec.Amount, &lNcbModifyRec.Isin, &lNcbModifyRec.DateTime, &lNcbModifyRec.RespOrderNo)

				if lErr3 != nil {
					log.Println("LMGD03", lErr3)
					return lNcbModifyRec, lErr3
				}

				// if lOrderNo.Valid {
				// 	lNcbModifyRec.OrderNo = int(lOrderNo.Int64)
				// } else {

				// 	lNcbModifyRec.OrderNo = 0
				// }
				lNcbModifyRec.Total = lNcbModifyRec.Unit * lNcbModifyRec.Price

			}
			log.Println("lOrderNo", pOrderNo, lNcbModifyRec.OrderNo)
			log.Println("lNcbModifyRec", lNcbModifyRec)
		}
	}
	log.Println("GetModifyDetails (-)")
	return lNcbModifyRec, nil
}
