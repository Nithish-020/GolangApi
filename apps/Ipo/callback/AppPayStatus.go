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

type payStruct struct {
	Symbol               string  `json:"symbol"`
	ApplicationNo        string  `json:"applicationNumber"`
	UpiPaymentStatusFlag int     `json:"upiPaymentStatusFlag"`
	UpiAmtBlocked        float32 `json:"upiAmtBlocked"`
	UpiPayReason         string  `json:"upiPayReason"`
}

type payRespStruct struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

/*
Pupose:This Function is used to Get the Application Upi status from the NSE-Exchange and update the
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

Author: Nithish Kumar
Date: 05JUNE2023
*/
func AppPayStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("AppPayStatus (+)", r.Method)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {
		// This Variable is used to get the Request object from the Exchange
		var lPayReqRec payStruct
		// This Variable is used to store the data from  GetDbData method
		// var lGetPayData payStruct
		// This Variable is used to send Response to Exchange
		var lPayResRec payRespStruct
		//This Variable is used Pass the details to Log method
		var lLogRec Function.ParameterStruct
		// This Variable is used to hold the last inserted id
		var lId int

		lLogRec.EndPoint = "/v1/appPayStatus"
		lLogRec.Method = "POST"
		lLogRec.ClientId = "NSE-Exchange"

		lBody, lErr1 := ioutil.ReadAll(r.Body)
		if lErr1 != nil {
			log.Println("CAPS01", lErr1)
			lPayResRec.Status = "failed"
			lPayResRec.Reason = "CAPS01" + lErr1.Error()
		} else {
			// Unmarshalling json to struct
			lErr2 := json.Unmarshal(lBody, &lPayReqRec)
			if lErr2 != nil {
				log.Println("CAPS02", lErr2)
				lPayResRec.Status = "failed"
				lPayResRec.Reason = "CAPS02" + lErr2.Error()
			} else {

				//Marshaling the Pay Request to give the parameter to LogEntry method
				lReq, lErr3 := json.Marshal(lPayReqRec)
				if lErr3 != nil {
					log.Println("CAPS03", lErr3)
					lPayResRec.Status = "failed"
					lPayResRec.Reason = "CAPS03" + lErr3.Error()
				} else {
					lLogRec.Request = string(lReq)
					lLogRec.Flag = common.INSERT

					var lErr4 error
					// insert the reques struct in Log
					lId, lErr4 = Function.LogEntry(lLogRec)
					if lErr4 != nil {
						log.Println("CAPS04", lErr4)
						lPayResRec.Status = "failed"
						lPayResRec.Reason = "CAPS04" + lErr4.Error()
					} else {

						lAppNoId, lErr5 := GetAppNo(lPayReqRec.ApplicationNo)
						if lErr5 != nil {
							log.Println("CAPS05", lErr5)
							lPayResRec.Status = "failed"
							lPayResRec.Reason = "CAPS05" + lErr5.Error()
						} else {
							// Check if the 2 struct values are equal
							if lAppNoId != 0 {
								lDb, lErr6 := ftdb.LocalDbConnect(ftdb.IPODB)
								if lErr6 != nil {
									log.Println("CAPS06", lErr6)
									lPayResRec.Status = "failed"
									lPayResRec.Reason = "CAPS06" + lErr6.Error()
								} else {
									defer lDb.Close()
									lSqlString := `Update a_ipo_order_header oh 
									set oh.upiPaymentStatusFlag = ? , oh.upiAmtBlocked = ? , oh.upiPayReason = ? ,oh.UpdatedBy = ?,oh.UpdatedDate = now()
									where oh.Id = ?`

									_, lErr7 := lDb.Exec(lSqlString, &lPayReqRec.UpiPaymentStatusFlag, &lPayReqRec.UpiAmtBlocked, &lPayReqRec.UpiPayReason, common.AUTOBOT, &lAppNoId)
									if lErr7 != nil {
										log.Println("CAPS07", lErr7)
										lPayResRec.Status = "failed"
										lPayResRec.Reason = "CAPS07" + lErr7.Error()
									} else {
										lPayResRec.Status = "success"
									}
								}
							} else {
								lPayResRec.Status = "failed"
								lPayResRec.Reason = "Application no does not exist"
							}
						}
					}

				}
			}
		}
		//Marshaling the Pay Response to give the parameter to LogEntry method
		lResp, lErr8 := json.Marshal(lPayResRec)
		if lErr8 != nil {
			log.Println("CAPS08", lErr8)
			lPayResRec.Status = "failed"
			lPayResRec.Reason = "CAPS08" + lErr8.Error()
		} else {
			lLogRec.Response = string(lResp)
			lLogRec.LastId = lId
			lLogRec.Flag = common.UPDATE

			var lErr9 error
			// update the response structure in log
			lId, lErr9 = Function.LogEntry(lLogRec)
			if lErr9 != nil {
				log.Println("CAPS09", lErr9)
				lPayResRec.Status = "failed"
				lPayResRec.Reason = "CAPS09" + lErr9.Error()
			}
		}
		// Marshal the response structure into a json
		lData, lErr10 := json.Marshal(lPayResRec)
		if lErr10 != nil {
			log.Println("CAPS10", lErr10)
			lPayResRec.Status = "failed"
			lPayResRec.Reason = "CAPS10" + lErr10.Error()
		} else {
			fmt.Fprintf(w, "%s\n", string(lData))
		}
		log.Println("AppPayStatus (-)", r.Method)
	}
}

/*
Pupose: This method returns the collection of data from the  a_ipo_oder_header database
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the payStructArr data
		from the a_ipo_oder_header Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetAppNo(pAppNo string) (int, error) {
	log.Println("GetAppNo  (+)")
	// Create the instance for payStruct
	var lAppNoId int

	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr != nil {
		log.Println("IGPD01", lErr)
		return lAppNoId, lErr
	} else {
		defer lDb.Close()
		lCoreString := `select oh.Id 
		from a_ipo_order_header oh
		where oh.applicationNo = ? `

		lRows, lErr := lDb.Query(lCoreString, pAppNo)
		if lErr != nil {
			log.Println("IGPD02", lErr)
			return lAppNoId, lErr
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr := lRows.Scan(&lAppNoId)
				// log.Println(token)
				if lErr != nil {
					log.Println("IGPD03", lErr)
					return lAppNoId, lErr
				}
			}
		}
	}
	log.Println("GetAppNo (-)")
	return lAppNoId, nil
}
