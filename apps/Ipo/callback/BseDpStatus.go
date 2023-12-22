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

type bseDpStruct struct {
	Symbol          string `json:"symbol"`
	ApplicationNo   string `json:"applicationNumber"`
	DpVerStatusFlag string `json:"dpVerStatusFlag"`
	DpVerReason     string `json:"dpVerReason"`
	DpVerFailCode   string `json:"dpVerFailCode"`
}

type bseDpRespStruct struct {
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

Author: KAVYA DHARSHANI
Date: 15SEP2023
*/
func BseDpStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("BseDpStatus (+)", r.Method)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {
		// This Variable is used to get the Request object from the Exchange
		var lBseDpReqRec bseDpStruct

		// This Variable is used to send Response to Exchange
		var lBseDpRespRec bseDpRespStruct

		//This Variable is used Pass the details to Log method
		var lLogRec Function.ParameterStruct

		// This Variable is used to hold the last inserted id
		var lId int

		lLogRec.EndPoint = "/v1/dpstatus"
		lLogRec.Method = "POST"
		lLogRec.ClientId = "BSE-Exchange"

		lBody, lErr1 := ioutil.ReadAll(r.Body)
		if lErr1 != nil {
			log.Println("CBDS01", lErr1)
			lBseDpRespRec.Status = "failed"
			lBseDpRespRec.Reason = "CBDS01" + lErr1.Error()
		} else {
			// Unmarshalling json to struct
			lErr2 := json.Unmarshal(lBody, &lBseDpReqRec)
			if lErr2 != nil {
				log.Println("CBDS02", lErr2)
				lBseDpRespRec.Status = "failed"
				lBseDpRespRec.Reason = "CBDS02" + lErr2.Error()
			} else {

				//Marshaling the Dp Request to give the parameter to LogEntry method
				req, lErr3 := json.Marshal(lBseDpReqRec)
				if lErr3 != nil {
					log.Println("CBDS03", lErr3)
					lBseDpRespRec.Status = "failed"
					lBseDpRespRec.Reason = "CBDS03" + lErr3.Error()
				} else {
					lLogRec.Request = string(req)
					lLogRec.Flag = common.INSERT

					var lErr4 error
					// insert the reques struct in Log
					lId, lErr4 = Function.LogEntry(lLogRec)
					if lErr4 != nil {
						log.Println("CBDS04", lErr4)
						lBseDpRespRec.Status = "failed"
						lBseDpRespRec.Reason = "CBDS04" + lErr4.Error()
					} else {
						// Get AppNoId from database using below method
						lAppNoId, lErr5 := GetAppNo(lBseDpReqRec.ApplicationNo)
						if lErr5 != nil {
							log.Println("CBDS05", lErr5)
							lBseDpRespRec.Status = "failed"
							lBseDpRespRec.Reason = "CBDS05" + lErr5.Error()
						} else {
							// check if the 2 struct values are equal
							if lAppNoId != 0 {
								lDb, lErr6 := ftdb.LocalDbConnect(ftdb.IPODB)
								if lErr6 != nil {
									log.Println("CBDS06", lErr6)
									lBseDpRespRec.Status = "failed"
									lBseDpRespRec.Reason = "CBDS06" + lErr6.Error()
								} else {
									defer lDb.Close()
									lSqlString := `Update a_ipo_order_header oh 
									set oh.dpVerStatusFlag = ?, oh.dpVerReason = ?, oh.dpVerFailCode = ? ,oh.UpdatedBy = ?,oh.UpdatedDate = now()
									where oh.Id = ?	`

									_, lErr7 := lDb.Exec(lSqlString, &lBseDpReqRec.DpVerStatusFlag, &lBseDpReqRec.DpVerReason, &lBseDpReqRec.DpVerFailCode, common.AUTOBOT, &lAppNoId)
									if lErr7 != nil {
										log.Println("CBDS07", lErr7)
										lBseDpRespRec.Status = "failed"
										lBseDpRespRec.Reason = "CBDS07" + lErr7.Error()
									}
								}
								lBseDpRespRec.Status = "success"
							} else {
								lBseDpRespRec.Status = "failed"
								lBseDpRespRec.Reason = "Application no does not exist"
							}
						}
					}
				}
			}
			//Marshaling the Dp Response to give the parameter to LogEntry method
			resp, lErr8 := json.Marshal(lBseDpRespRec)
			if lErr8 != nil {
				log.Println("CBDS08", lErr8)
				lBseDpRespRec.Status = "failed"
				lBseDpRespRec.Reason = "CBDS08" + lErr8.Error()
			} else {
				lLogRec.Response = string(resp)
				lLogRec.LastId = lId
				lLogRec.Flag = common.UPDATE
				var lErr9 error
				// update the response structure in log
				lId, lErr9 = Function.LogEntry(lLogRec)
				if lErr9 != nil {
					log.Println("CBDS09", lErr9)
					lBseDpRespRec.Status = "failed"
					lBseDpRespRec.Reason = "CBDS09" + lErr9.Error()
				}
			}
		}
		// Marshall the structure into json
		lData, lErr10 := json.Marshal(lBseDpRespRec)
		if lErr10 != nil {
			log.Println("CBDS10", lErr10)
			lBseDpRespRec.Status = "failed"
			lBseDpRespRec.Reason = "CBDS10" + lErr10.Error()
		} else {
			fmt.Fprintf(w, "%s\n", string(lData))
		}
		log.Println("BseDpStatus (-)", r.Method)
	}
}
