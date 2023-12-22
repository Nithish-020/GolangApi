package callback

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
)

type bsePayStruct struct {
	Symbol               string  `json:"symbol"`
	ApplicationNo        string  `json:"applicationnumber"`
	UpiPaymentStatusFlag int     `json:"upipaymentstatusflag"`
	UpiAmtBlocked        float32 `json:"upiamtblocked"`
	UpiPayReason         string  `json:"upipayreason"`
}

type bsePayRespStruct struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

/*
Pupose:This Function is used to Get the Application Upi status from the BSE-Exchange and update the
changes in our database table ....
Parameters:

	{
		"symbol": "TEST",
		"applicationNumber": "1200299929020",
		"upiPaymentStatusFlag": 100,
		"upiAmtBlocked": 125000.00,
		"upiPayReason": null
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

Author: KAVYA DHARSHANI
Date: 15SEP2023
*/
func BsePayStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("BsePayStatus (+)", r.Method)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {
		// This Variable is used to get the Request object from the Exchange
		var lBsePayReqRec bsePayStruct

		// This Variable is used to send Response to Exchange
		var lBsePayResRec bsePayRespStruct

		//This Variable is used Pass the details to Log method
		var lLogRec Function.ParameterStruct
		// This Variable is used to hold the last inserted id
		var lId int

		lLogRec.EndPoint = "/v1/upistatus"
		lLogRec.Method = "POST"
		lLogRec.ClientId = "BSE-Exchange"

		lBody, lErr1 := ioutil.ReadAll(r.Body)
		if lErr1 != nil {
			log.Println("CBPS01", lErr1)
			lBsePayResRec.Status = "failed"
			lBsePayResRec.Reason = "CBPS01" + lErr1.Error()
		} else {
			// Unmarshalling json to struct
			lErr2 := json.Unmarshal(lBody, &lBsePayReqRec)
			if lErr2 != nil {
				log.Println("CBPS02", lErr2)
				lBsePayResRec.Status = "failed"
				lBsePayResRec.Reason = "CBPS02" + lErr2.Error()
			} else {

				//Marshaling the Pay Request to give the parameter to LogEntry method
				lReq, lErr3 := json.Marshal(lBsePayReqRec)
				if lErr3 != nil {
					log.Println("CBPS03", lErr3)
					lBsePayResRec.Status = "failed"
					lBsePayResRec.Reason = "CBPS03" + lErr3.Error()
				} else {
					lLogRec.Request = string(lReq)
					lLogRec.Flag = common.INSERT

					var lErr4 error
					// insert the reques struct in Log
					lId, lErr4 = Function.LogEntry(lLogRec)
					if lErr4 != nil {
						log.Println("CBPS04", lErr4)
						lBsePayResRec.Status = "failed"
						lBsePayResRec.Reason = "CBPS04" + lErr4.Error()
					} else {

						lAppNoId, lErr5 := GetAppNo(lBsePayReqRec.ApplicationNo)
						if lErr5 != nil {
							log.Println("CBPS05", lErr5)
							lBsePayResRec.Status = "failed"
							lBsePayResRec.Reason = "CBPS05" + lErr5.Error()
						} else {
							// Check if the 2 struct values are equal
							if lAppNoId != 0 {
								lDb, lErr6 := ftdb.LocalDbConnect(ftdb.IPODB)
								if lErr6 != nil {
									log.Println("CBPS06", lErr6)
									lBsePayResRec.Status = "failed"
									lBsePayResRec.Reason = "CBPS06" + lErr6.Error()
								} else {
									defer lDb.Close()
									lSqlString := `Update a_ipo_order_header oh 
									set oh.upiPaymentStatusFlag = ? , oh.upiAmtBlocked = ? , oh.upiPayReason = ? ,oh.UpdatedBy = ?,oh.UpdatedDate = now()
									where oh.Id = ?`

									_, lErr7 := lDb.Exec(lSqlString, &lBsePayReqRec.UpiPaymentStatusFlag, &lBsePayReqRec.UpiAmtBlocked, &lBsePayReqRec.UpiPayReason, common.AUTOBOT, &lAppNoId)
									if lErr7 != nil {
										log.Println("CBPS07", lErr7)
										lBsePayResRec.Status = "failed"
										lBsePayResRec.Reason = "CBPS07" + lErr7.Error()
									} else {
										lBsePayResRec.Status = "success"
									}
								}
							} else {
								lBsePayResRec.Status = "failed"
								lBsePayResRec.Reason = "Application no does not exist"
							}
						}
					}

				}
			}
		}
		//Marshaling the Pay Response to give the parameter to LogEntry method
		lResp, lErr8 := json.Marshal(lBsePayResRec)
		if lErr8 != nil {
			log.Println("CBPS08", lErr8)
			lBsePayResRec.Status = "failed"
			lBsePayResRec.Reason = "CBPS08" + lErr8.Error()
		} else {
			lLogRec.Response = string(lResp)
			lLogRec.LastId = lId
			lLogRec.Flag = common.UPDATE

			var lErr9 error
			// update the response structure in log
			lId, lErr9 = Function.LogEntry(lLogRec)
			if lErr9 != nil {
				log.Println("CBPS09", lErr9)
				lBsePayResRec.Status = "failed"
				lBsePayResRec.Reason = "CBPS09" + lErr9.Error()
			}
		}
		// Marshal the response structure into a json
		lData, lErr10 := json.Marshal(lBsePayResRec)
		if lErr10 != nil {
			log.Println("CBPS10", lErr10)
			lBsePayResRec.Status = "failed"
			lBsePayResRec.Reason = "CBPS10" + lErr10.Error()
		} else {
			fmt.Fprintf(w, "%s\n", string(lData))
		}
		log.Println("BsePayStatus (-)", r.Method)
	}
}
