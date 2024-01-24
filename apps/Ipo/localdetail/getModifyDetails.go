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

// this Structure is used to get the application details
type ModifyBidStruct struct {
	Id           int     `json:"id"`
	ActivityType string  `json:"activityType"`
	BidRefNo     string  `json:"bidReferenceNo"`
	Price        float64 `json:"price"`
	Amount       float64 `json:"amount"`
	Quantity     int     `json:"quantity"`
	CutOff       bool    `json:"cutOff"`
}

// Response Sturcture for getModifyDetails API
type ModifyRespStruct struct {
	MasterId      int               `json:"masterId"`
	ApplicationNo string            `json:"appNo"`
	Upi           string            `json:"upi"`
	Category      string            `json:"category"`
	Total         float64           `json:"total"`
	ModifyDetails []ModifyBidStruct `json:"modifyDetails"`
	Status        string            `json:"status"`
	ErrMsg        string            `json:"errMsg"`
}

/*
Purpose:This Method is used to get the bid detail records
Request:

Header Value: ID

Response:

=========
*On Sucess
=========

	{
		"modifyDetails" :[

			"appNo" : "1234567890"
			"upi" : "test@ybl"
			"category" : "IND"
			{
				"bidRefNo" : "45687843213233"
				"price" : 755
				"amount" : 14567
				"quantity" : 19
				"cutOff" : true
			},
			{
				"bidRefNo" : "45687843213233"
				"price" : 755
				"amount" : 14567
				"quantity" : 19
				"cutOff" : true
			}
		],
		"status" : "S",
		"errMsg" : ""
	}

=========
!On Error
=========

	{
		"modifyDetails" :[],
		"status": "E",
		"errMsg": "Can't able to get data from database"
	}

Author: Nithish Kumar
Date: 07JUNE2023
*/
func GetModifyDetails(w http.ResponseWriter, r *http.Request) {
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
	(w).Header().Set("Access-Control-Allow-Headers", "ID,CATEGORY,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	log.Println("GetModifyDetails(+)", r.Method)

	if r.Method == "GET" {

		// Create an instance for ModifyBid
		var lRespRec ModifyRespStruct
		lRespRec.Status = common.SuccessCode
		// reading Header Value from the Request
		lMasterId := r.Header.Get("ID")
		lCategory := r.Header.Get("CATEGORY")
		log.Println("Request := ", lMasterId, lCategory)

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("LGMD01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGMD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGMD01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGMD03", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		//call getdetails method to get the application Details
		lRespArr, lErr2 := getDetails(lClientId, lMasterId, lBrokerId, lCategory)
		if lErr2 != nil {
			log.Println("LGMD03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGMD03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGMD03", "Unable to fetch your Records"))
			return
		} else {
			lRespRec = lRespArr
			lRespRec.Status = common.SuccessCode
			log.Println(lRespRec)
		}
		// Marshal the response structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGMD04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGMD04", "Issue in Getting your Application Details!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetModifyDetails (-)", r.Method)
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
[
{
"appNo" : "1234567890"
"upi" : "test@ybl"
"category" : "IND"
"bidRefNo" : "45687843213233"
"price" : 755
"amount" : 14567
"quantity" : 19
"cutOff" : true
},
{
"appNo" : "1234575790"
"upi" : "test@ybl"
"category" : "IND"
"bidRefNo" : "45687843213233"
"price" : 755
"amount" : 14567
"quantity" : 19
"cutOff" : true
}
]
==========
!On Error
==========
[],error

Author:Pavithra
Date: 12JUNE2023
*/
func getDetails(pClientId string, pMasterId string, pBrokerId int, pCategory string) (ModifyRespStruct, error) {
	log.Println("getDetails (+)")
	// This Array Variable is used to store the Modify Application Records
	// var lRespArr []ModifyBidStruct
	// This Variable is used to store the Each Modify bid Records
	var lRespRec ModifyBidStruct

	var lAmount float64

	var lResponseRec ModifyRespStruct
	var lLotSize, lQty int
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LMGD01", lErr1)
		return lResponseRec, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select h.MasterId,h.applicationNo,d.bidReferenceNo,h.upi,
						h.category,d.id,d.req_price,d.req_amount,d.req_quantity,
						m.LotSize,d.atCutOff,d.activityType
						from a_ipo_order_header h,a_ipo_orderdetails d,
						a_ipo_master m
						where m.Id = h.MasterId 
						and h.Id = d.headerId
						and d.status <> 'failed'	
						and h.cancelFlag = 'N'
						and d.activityType <> 'cancel'	
						and h.clientId = ?
						and h.MasterId = ? 
						and h.brokerId = ?
						and h.category = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pMasterId, pBrokerId, pCategory)
		if lErr2 != nil {
			log.Println("LMGD02", lErr2)
			return lResponseRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in ModifyRespStruct
			for lRows.Next() {
				lErr3 := lRows.Scan(&lResponseRec.MasterId, &lResponseRec.ApplicationNo, &lRespRec.BidRefNo, &lResponseRec.Upi, &lResponseRec.Category, &lRespRec.Id, &lRespRec.Price, &lRespRec.Amount, &lQty, &lLotSize, &lRespRec.CutOff, &lRespRec.ActivityType)
				if lErr3 != nil {
					log.Println("LMGD03", lErr3)
					return lResponseRec, lErr3
				} else {
					// Append the Bid Records in lRespArr
					lAmount = lRespRec.Amount
					lResponseRec.Total += lAmount
					lRespRec.Quantity = int(lQty) / int(lLotSize)
					lResponseRec.ModifyDetails = append(lResponseRec.ModifyDetails, lRespRec)
				}
			}
		}
	}
	log.Println("getDetails (-)")
	return lResponseRec, nil
}
