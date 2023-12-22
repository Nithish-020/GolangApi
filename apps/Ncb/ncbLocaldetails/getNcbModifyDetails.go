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

type NcbModifyStruct struct {
	Id            int     `json:"id"`
	Symbol        string  `json:"name"`
	ActivityType  string  `json:"activityType"`
	Flag          string  `json:"flag"`
	OrderNo       int     `json:"orderNo"`
	ApplicationNo string  `json:"applicationNo"`
	LotSize       int     `json:"lotSize"`
	Unit          int     `json:"unit"`
	Price         float64 `json:"price"`
}

// Response Structure for GetSgbMaster API
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
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
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

		log.Println("lClientId", lClientId)
		lRespStruct, lErr2 := GetModifyDetails(lClientId, lConvId, lIdOrder)
		if lErr2 != nil {
			log.Println("NLGNMD03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NLGNMD03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NLGNMD03", "Error Occur in getting Datas.."))
			return
		} else {
			lRespRec.NcbModify = lRespStruct
			log.Println("lRespRec", lRespRec.NcbModify)
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

==========
!On Error
==========
[],error

Author:KAVYADHARSHANI
Date: 21OCT2023
*/
func GetModifyDetails(pClientId string, pMasterId int, lOrderNo int) (NcbModifyStruct, error) {
	log.Println("GetModifyDetails (+)")

	// var lAmount float64

	// var lLotSize, lUnit int
	// var lPrice int

	// var lOrderNo sql.NullInt64
	var lNcbModifyRec NcbModifyStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LMGD01", lErr1)
		return lNcbModifyRec, lErr1
	} else {
		defer lDb.Close()

		// 		lCoreString := `select h.MasterId,d.OrderNo	, d.price , d.activityType flag, d.Unit
		// 		from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
		// 		where n.id = h.MasterId
		// 		and h.Id  = d.headerId
		// 		and d.status <> 'failed' and h.cancelFlag = 'N' and d.activityType <> 'cancel'
		// and h.clientId = ? and h.MasterId = ?`
		// and d.OrderNo = ?
		lCoreString := `select h.MasterId	, d.price , d.activityType flag, CAST((h.Investmentunit /n.Lotsize) AS SIGNED) Lotsize, d.unit unit, h.Symbol,h.cancelFlag flag, h.applicationNo 
		from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h	
		where h.id = d.headerId 
		and n.id  = h.MasterId 
		and d.status <> 'failed' and h.cancelFlag <> 'Y'
             and h.clientId = ? and h.MasterId = ? and d.OrderNo = ?`

		log.Println("pClientId", pClientId, pMasterId, "pMasterId")
		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pMasterId, lOrderNo)
		if lErr2 != nil {
			log.Println("LMGD02", lErr2)
			return lNcbModifyRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in ModifyRespStruct
			for lRows.Next() {
				lErr3 := lRows.Scan(&lNcbModifyRec.Id, &lNcbModifyRec.Price, &lNcbModifyRec.ActivityType, &lNcbModifyRec.LotSize, &lNcbModifyRec.Unit, &lNcbModifyRec.Symbol, &lNcbModifyRec.Flag, &lNcbModifyRec.ApplicationNo)

				if lErr3 != nil {
					log.Println("LMGD03", lErr3)
					return lNcbModifyRec, lErr3
				} else {
					// if lOrderNo.Valid {
					// 	lNcbModifyRec.OrderNo = int(lOrderNo.Int64)
					// } else {

					// 	lNcbModifyRec.OrderNo = 0
					// }
				}
			}
			log.Println("lOrderNo", lOrderNo, lNcbModifyRec.OrderNo)
			log.Println("lNcbModifyRec", lNcbModifyRec)
		}
	}
	log.Println("GetModifyDetails (-)")
	return lNcbModifyRec, nil
}
