package brokers

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type addBrokerStruct struct {
	Id         int    `json:"id"`
	BrokerName string `json:"brokerName"`
	Domain     string `json:"domainName"`
	Type       string `json:"type"`
	RawDomain  string `json:"rawDomain"`
	AppName    string `json:"appName"`
	AuthURL    string `json:"authURL"`
	Status     string `json:"status"`
}
type addBrokerResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetBroker(w http.ResponseWriter, r *http.Request) {
	log.Println("SetBroker (+)")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {

		var lReqRec addBrokerStruct
		var lRespRec addBrokerResp

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/domainsetup")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BSB01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BSB01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BSB02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("BSB02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BSB03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BSB03", "Request cannot be empty.Please try after sometime"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lReqRec)
			log.Println("lReqRec", lReqRec)
			if lErr4 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BSB04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("BSB04", "Unable to proccess your request right now.Please try after sometime"))
				return
			} else {

				lDb, lErr5 := ftdb.LocalDbConnect(ftdb.IPODB)
				if lErr5 != nil {
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "BSB05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("BSB05", "Error in connecting server.Please try after sometime"))
					return
				} else {
					defer lDb.Close()

					if lReqRec.Id == 0 {
						sqlString := `Insert into a_ipo_brokerMaster (BrokerName,DomainName,Status,Type,AuthURL,RawDomain,AppName,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate) 
								  values (?,?,?,?,?,?,?,?,now(),?,now())`
						_, lErr6 := lDb.Exec(sqlString, lReqRec.BrokerName, lReqRec.Domain, lReqRec.Status, lReqRec.Type, lReqRec.AuthURL, lReqRec.RawDomain, lReqRec.AppName, lClientId, lClientId)
						if lErr6 != nil {
							log.Println("BAB06", lErr6)
							fmt.Fprintf(w, helpers.GetErrorString("BSB06", "Error in connecting server.Please try after sometime"))
							return
						} else {
							lRespRec.Status = common.SuccessCode
							lRespRec.ErrMsg = "Record added successfully"

						}
					} else {

						sqlString := `update a_ipo_brokerMaster bm
						set bm.BrokerName = ?,bm.DomainName = ?,bm.Status = ?,bm.Type = ?,bm.AuthURL = ?,bm.RawDomain =?,
						bm.AppName = ?,bm.UpdatedBy = ?,bm.UpdatedDate = now()
						where bm.Id = ?`

						_, lErr7 := lDb.Exec(sqlString, lReqRec.BrokerName, lReqRec.Domain, lReqRec.Status, lReqRec.Type, lReqRec.AuthURL, lReqRec.RawDomain, lReqRec.AppName, lClientId, lReqRec.Id)
						if lErr7 != nil {
							log.Println("BAB03", lErr7)
							fmt.Fprintf(w, helpers.GetErrorString("BSB05", "Error in connecting server.Please try after sometime"))
							return
						} else {
							lRespRec.Status = common.SuccessCode
							lRespRec.ErrMsg = "Record updated successfully"

						}
					}
				}
			}
		}
		// Marshal the response structure into lData
		lData, lErr8 := json.Marshal(lRespRec)
		if lErr8 != nil {
			log.Println("BSB07", lErr8)
			fmt.Fprintf(w, helpers.GetErrorString("BSB07", "Error found!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("SetBroker (-)")
}
