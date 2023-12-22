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

// To store IPO History Records
type HistoryStruct struct {
	Id                int     `json:"id"`
	MasterId          int     `json:"masterId"`
	Symbol            string  `json:"symbol"`
	ApplicationNo     string  `json:"applicationNo"`
	Price             float64 `json:"price"`
	Quantity          int     `json:"quantity"`
	CutOff            string  `json:"cutOff"`
	Status            string  `json:"status"`
	Date              string  `json:"date"`
	ApplicationStatus string  `json:"applicationStatus"`
	UpiStatus         string  `json:"upiStatus"`
	CancelFlag        string  `json:"cancelFlag"`
	StartDate         string  `json:"startDate"`
	EndDate           string  `json:"endDate"`
	Allotment         string  `json:"allotment"`
	Refund            string  `json:"refund"`
	Demat             string  `json:"demat"`
	Listing           string  `json:"listing"`
}

// Response Structure for GetIpoHistory API
type HistoryRespStruct struct {
	HistoryArray []HistoryStruct `json:"history"`
	Status       string          `json:"status"`
	ErrMsg       string          `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Application Upi status from the NSE-Exchange and update the
changes in our database table ....
Parameters:

	{
		"symbol": "TEST",
		"applicationNumber": "1200299929020",
		"dpVerStatusFlag": "S",
		"dpVerReason": null,
		"dpVerFailCode": null
	}

Response:

	*On Sucess
	=========
	{
		"status": "success"
	}

	!On Error
	========
	{
		"status": "failed",
		"reason": "Application no does not exist"
	}

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetIpoHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetIpoHistory (+)", r.Method)
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
		// This variable is used to store the history structure
		var lHistoryRespRec HistoryRespStruct
		lHistoryRespRec.Status = common.SuccessCode

		/***********ADMIN PAGE
		//-----------START TO GET CLIENT  DETAILS
		// This variable is used to store the ClientId
		lClientId, err := appsso.ValidateAndGetClientDetails(r, common.ABHIAppName, common.ABHICookieName)
		if err != nil {
			log.Println("IGIH01.1", err)
			fmt.Fprintf(w, helpers.GetErrorString("IGIH01.1", "The server is busy right now. Please try after sometime"))
			return
		}
		//-----------END OF GETTING CLIENT  DETAILS
		*********************/

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		// lClientId := ""
		// var lErr1 error
		lClientId, lErr1 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("LGIH01", lErr1)
			lHistoryRespRec.Status = common.ErrorCode
			lHistoryRespRec.ErrMsg = "LGIH01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGIH02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// Call GetData method to get Token from Database
		lhistoryArr, lErr2 := GetData(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("LGIH02", lErr2)
			lHistoryRespRec.Status = common.ErrorCode
			lHistoryRespRec.ErrMsg = "LGIH02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH02", "Issue in Fetching Your Datas!"))
			return
		} else {
			// Assign lhistoryArr to HistoryRespStruct Value
			lHistoryRespRec.HistoryArray = lhistoryArr
		}
		// Marshal the reponse structure into lDatas
		lData, lErr3 := json.Marshal(lHistoryRespRec)
		if lErr3 != nil {
			log.Println("LGIH03", lErr3)
			lHistoryRespRec.Status = common.ErrorCode
			lHistoryRespRec.ErrMsg = "LGIH03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH03", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetIpoHistory (-)", r.Method)
	}
}

/*
Pupose: This method returns the collection of data from the  database
Parameters:

	not applicable

Response:

	    *On Sucess
	    ==========
	    In case of a successful execution of this method, you will get the historyStruct data
		from database

	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetData(pClientId string, pBrokerId int) ([]HistoryStruct, error) {
	log.Println("GetData (+)")
	// This variable is used to store the history structure from the database
	var lhistoryRec HistoryStruct
	// This variable is used to store the history structure from the database n the Array
	var lhistoryArr []HistoryStruct
	// This variable is used to store the stepper structure from the database
	// var lStepperRec stepperStruct

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LHGD01", lErr1)
		return lhistoryArr, lErr1
	} else {
		defer lDb.Close()

		// ! To get the Symbol, ApplicationNo, Price and Quantity
		lCoreString2 :=
			//  `select nvl(d.Id,''),m.Id,
			// nvl(h.Symbol,'') ,
			// nvl(h.ApplicationNo,'') ,
			// nvl(d.req_price,'') ,
			// nvl(d.req_quantity,''),
			// nvl((case when d.req_price = m.MaxPrice then "CutOff" else 'NA' end),'') CutOff ,
			// nvl(h.status,''),
			// nvl(DATE_FORMAT(h.CreatedDate , '%d-%b-%Y'),''),
			// nvl(
			// 	(select xld.description
			// 	from xx_lookup_details xld,xx_lookup_header xlh
			// 	where xld.headerid = xlh.id
			// 	and xlh.Code = 'IpoDp'
			// 	and xld.Code = h.dpVerStatusFlag) , 'N/A'
			// 	) ApplicationStatus,
			// nvl(
			// 	(select (case when h.status = 'success' then xld.description else 'N/A' end)
			// 	from xx_lookup_details xld,xx_lookup_header xlh
			// 	where xld.headerid = xlh.id
			// 	and xlh.Code = 'IpoPay'
			// 	and xld.Code = h.upiPaymentStatusFlag), 'N/A'
			// ) UPIStatus,
			// nvl(h.cancelFlag, 'N/A') CancelFlag	,
			// from a_ipo_orderdetails d ,a_ipo_order_header h ,a_ipo_master m
			// where d.HeaderId = h.Id
			// and m.Id = h.MasterId
			// and h.ClientId = ?
			// and h.cancelFlag <> 'Y'
			// group by h.applicationNo
			// order by d.Id desc`

			`select mas.detailId,
			mas.Id,
			mas.Symbol,
			mas.ApplicationNo,
			mas.req_price,
			mas.req_quantity,
			mas.CutOff,
			mas.status,
			nvl(DATE_FORMAT(mas.createdDate , '%d-%b-%Y'),'') createdDate,
			mas.ApplicationStatus,
			mas.UPIStatus,
			mas.CancelFlag,
			mas.startdate,
			mas.enddate,
			nvl(id.allotmentFinal,'') ,nvl(id.refundInitiate,'') ,nvl(id.dematTransfer,''),nvl(id.listingDate,'')
			from (
			select nvl(d.Id,'') detailId,
			m.Id Id,
			nvl(h.Symbol,'') Symbol,
			nvl(h.ApplicationNo,'') ApplicationNo,
			nvl(d.req_price,'') req_price,
			nvl(d.req_quantity,'') req_quantity,
			nvl((case when d.req_price = m.MaxPrice then "CutOff" else 'NA' end),'') CutOff ,
			nvl(h.status,'') status,
			nvl(h.createdDate,'') createdDate,
			nvl(
				(select xld.description 
				from xx_lookup_details xld,xx_lookup_header xlh
				where xld.headerid = xlh.id
				and xlh.Code = 'IpoDp'
				and xld.Code = h.dpVerStatusFlag) , 'N/A'
				) ApplicationStatus,
			nvl(
				(select (case when h.status = 'success' then xld.description else 'N/A' end)
				from xx_lookup_details xld,xx_lookup_header xlh
				where xld.headerid = xlh.id
				and xlh.Code = 'IpoPay'
				and xld.Code = h.upiPaymentStatusFlag), 'N/A'
			) UPIStatus,
			nvl(h.cancelFlag, 'N/A') CancelFlag	,
			nvl(m.BiddingStartDate,'') startdate,
			nvl(m.BiddingEndDate,'') enddate,
			m.Isin isin
			from a_ipo_orderdetails d ,a_ipo_order_header h ,a_ipo_master m
			where d.HeaderId = h.Id 
			and m.Id = h.MasterId 
			and h.ClientId = ?
			and h.brokerId = ?
			group by h.applicationNo 
			) mas left join a_ipo_details id 
			on mas.isin = id.Isin 
			order by mas.createdDate desc`

		// and m.BiddingEndDate < curdate() Its goes befor client id in query
		lRows2, lErr2 := lDb.Query(lCoreString2, pClientId, pBrokerId)
		if lErr2 != nil {
			log.Println("LHGD02", lErr2)
			return lhistoryArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows2.Next() {
				lErr3 := lRows2.Scan(&lhistoryRec.Id, &lhistoryRec.MasterId, &lhistoryRec.Symbol, &lhistoryRec.ApplicationNo, &lhistoryRec.Price,
					&lhistoryRec.Quantity, &lhistoryRec.CutOff, &lhistoryRec.Status, &lhistoryRec.Date, &lhistoryRec.ApplicationStatus, &lhistoryRec.UpiStatus,
					&lhistoryRec.CancelFlag, &lhistoryRec.StartDate, &lhistoryRec.EndDate, &lhistoryRec.Allotment, &lhistoryRec.Refund, &lhistoryRec.Demat, &lhistoryRec.Listing)
				if lErr3 != nil {
					log.Println("LHGD03", lErr3)
					return lhistoryArr, lErr3
				} else {
					// Append the history Records into lhistoryArr Array
					lhistoryArr = append(lhistoryArr, lhistoryRec)

					// log.Println("History Resp :=", lhistoryRec)
				}
			}
		}
	}
	log.Println("GetData (-)")
	return lhistoryArr, nil
}
