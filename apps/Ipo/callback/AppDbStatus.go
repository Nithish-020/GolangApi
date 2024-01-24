package callback

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type dpStruct struct {
	Symbol          string `json:"symbol"`
	ApplicationNo   string `json:"applicationNumber"`
	DpVerStatusFlag string `json:"dpVerStatusFlag"`
	DpVerReason     string `json:"dpVerReason"`
	DpVerFailCode   string `json:"dpVerFailCode"`
}

type dpRespStruct struct {
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
func AppDpStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("AppDpStatus (+)", r.Method)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {
		// This Variable is used to get the Request object from the Exchange
		var lDpReqRec dpStruct
		// This Variable is used to store the data from GetDpData method
		// var lGetDbDataRec dpStruct
		// This Variable is used to send Response to Exchange
		var lDpRespRec dpRespStruct
		//This Variable is used Pass the details to Log method
		var lLogRec Function.ParameterStruct
		// This Variable is used to hold the last inserted id
		var lId int

		lLogRec.EndPoint = "/appDpStatus"
		lLogRec.Method = "POST"
		lLogRec.ClientId = "NSE-Exchange"

		lBody, lErr1 := ioutil.ReadAll(r.Body)
		if lErr1 != nil {
			log.Println("CADS01", lErr1)
			lDpRespRec.Status = "failed"
			lDpRespRec.Reason = "CADS01" + lErr1.Error()
		} else {
			// Unmarshalling json to struct
			lErr2 := json.Unmarshal(lBody, &lDpReqRec)
			if lErr2 != nil {
				log.Println("CADS02", lErr2)
				lDpRespRec.Status = "failed"
				lDpRespRec.Reason = "CADS02" + lErr2.Error()
			} else {

				//Marshaling the Dp Request to give the parameter to LogEntry method
				req, lErr3 := json.Marshal(lDpReqRec)
				if lErr3 != nil {
					log.Println("CADS03", lErr3)
					lDpRespRec.Status = "failed"
					lDpRespRec.Reason = "CADS03" + lErr3.Error()
				} else {
					lLogRec.Request = string(req)
					lLogRec.Flag = common.INSERT

					var lErr4 error
					// insert the reques struct in Log
					lId, lErr4 = Function.LogEntry(lLogRec)
					if lErr4 != nil {
						log.Println("CADS04", lErr4)
						lDpRespRec.Status = "failed"
						lDpRespRec.Reason = "CADS04" + lErr4.Error()
					} else {
						// Get AppNoId from database using below method
						lAppNoId, lErr5 := GetAppNo(lDpReqRec.ApplicationNo)
						if lErr5 != nil {
							log.Println("CADS05", lErr5)
							lDpRespRec.Status = "failed"
							lDpRespRec.Reason = "CADS05" + lErr5.Error()
						} else {
							// check if the 2 struct values are equal
							if lAppNoId != 0 {
								lDb, lErr6 := ftdb.LocalDbConnect(ftdb.IPODB)
								if lErr6 != nil {
									log.Println("CADS06", lErr6)
									lDpRespRec.Status = "failed"
									lDpRespRec.Reason = "CADS06" + lErr6.Error()
								} else {
									defer lDb.Close()
									lSqlString := `Update a_ipo_order_header oh 
									set oh.dpVerStatusFlag = ?, oh.dpVerReason = ?, oh.dpVerFailCode = ? ,oh.UpdatedBy = ?,oh.UpdatedDate = now()
									where oh.Id = ?	`

									_, lErr7 := lDb.Exec(lSqlString, &lDpReqRec.DpVerStatusFlag, &lDpReqRec.DpVerReason, &lDpReqRec.DpVerFailCode, common.AUTOBOT, &lAppNoId)
									if lErr7 != nil {
										log.Println("CADS07", lErr7)
										lDpRespRec.Status = "failed"
										lDpRespRec.Reason = "CADS07" + lErr7.Error()
									}
								}
								lDpRespRec.Status = "success"
							} else {
								lDpRespRec.Status = "failed"
								lDpRespRec.Reason = "Application no does not exist"
							}
						}
					}
				}
			}
			//Marshaling the Dp Response to give the parameter to LogEntry method
			resp, lErr8 := json.Marshal(lDpRespRec)
			if lErr8 != nil {
				log.Println("CADS08", lErr8)
				lDpRespRec.Status = "failed"
				lDpRespRec.Reason = "CADS08" + lErr8.Error()
			} else {
				lLogRec.Response = string(resp)
				lLogRec.LastId = lId
				lLogRec.Flag = common.UPDATE
				var lErr9 error
				// update the response structure in log
				lId, lErr9 = Function.LogEntry(lLogRec)
				if lErr9 != nil {
					log.Println("CADS09", lErr9)
					lDpRespRec.Status = "failed"
					lDpRespRec.Reason = "CADS09" + lErr9.Error()
				}
			}
		}
		// Marshall the structure into json
		lData, lErr10 := json.Marshal(lDpRespRec)
		if lErr10 != nil {
			log.Println("CADS10", lErr10)
			lDpRespRec.Status = "failed"
			lDpRespRec.Reason = "CADS10" + lErr10.Error()
		} else {
			fmt.Fprintf(w, "%s\n", string(lData))
		}
		log.Println("AppDpStatus (-)", r.Method)
	}
}

/*
Pupose: This method returns the collection of data from the  a_ipo_oder_header database
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the dpStructArr data
		from the a_ipo_oder_header Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 05JUNE2023s
*/
// func GetDpData() (dpStruct, error) {
// 	log.Println("GetDbData  (+)")
// 	//create a instance for dpStruct
// 	var lGetDbDataRec dpStruct

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr != nil {
// 		log.Println("IGDD01", lErr)
// 		return lGetDbDataRec, lErr
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select oh.Symbol ,oh.applicationNo ,nvl(oh.dpVerStatusFlag,0) ,nvl(oh.dpVerReason,'') ,nvl(oh.dpVerFailCode,'')
// 		from a_ipo_order_header oh `

// 		lRows, lErr := lDb.Query(lCoreString)
// 		if lErr != nil {
// 			log.Println("IGDD02", lErr)
// 			return lGetDbDataRec, lErr
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows.Next() {
// 				err := lRows.Scan(&lGetDbDataRec.Symbol, &lGetDbDataRec.ApplicationNo, &lGetDbDataRec.DpVerFailCode, &lGetDbDataRec.DpVerReason, &lGetDbDataRec.DpVerFailCode)
// 				// log.Println(token)
// 				if err != nil {
// 					log.Println("IGDD03", err)
// 					return lGetDbDataRec, err
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetDbData (-)")
// 	return lGetDbDataRec, nil
// }
