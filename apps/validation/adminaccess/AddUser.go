package adminaccess

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// * from DB side
type AddUserReqStruct struct {
	Id       int    `json:"id"`
	BrokerId int    `json:"brokerId"`
	ClientId string `json:"clientId"`
	Status   string `json:"status"`
	RoleId   int    `json:"roleId"`
}

// response struct for AddUser API
type RespStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func Adduser(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
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
		log.Println("AddUser (+)")
		var lReqRec AddUserReqStruct
		var lRespRec RespStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/managerole")
		if lErr1 != nil {
			log.Println("IAA01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "IAA01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("IAA01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("IAA02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr2 := ioutil.ReadAll(r.Body)
		if lErr2 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "IAA03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("IAA03", "Issue in reading request.Please try after sometime"))
			return
		} else {
			lErr3 := json.Unmarshal(lBody, &lReqRec)
			if lErr3 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "IAA04" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("IAA04", "The Request was empty.Please enter valied deatil and try again"))
				return
			} else {
				// log.Println("lReqRec.ClientId", lReqRec.Id, lReqRec.ClientId, lReqRec.Status, lReqRec.RoleId, lReqRec.BrokerId)
				if lReqRec.ClientId != "" && lReqRec.Status != "" && lReqRec.RoleId != 0 {
					lReqRec.BrokerId = lBrokerId
					log.Println("lReqRec := ", lReqRec)
					lErr4 := addQuery(lClientId, lReqRec)
					if lErr4 != nil {

						// lRespRec.Status = common.ErrorCode
						// lRespRec.ErrMsg = "IAA03" + helpers.ErrPrint(lErr4)
						lDebug.Log(helpers.Elog, lErr4.Error())
						fmt.Fprintf(w, helpers.GetErrorString("IAA05", "Issue in Connecting DB.Please try after sometime"))
						return
					}
				}
			}
		}
		lData, lErr6 := json.Marshal(lRespRec)
		if lErr6 != nil {
			fmt.Fprintf(w, "Error taking data"+lErr6.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("AddUser (-)", r.Method)
	}
}

func addQuery(pClientId string, pReqRec AddUserReqStruct) error {
	log.Println("addQuery (+)")
	// TO capture the signal from select query
	// lsignal := ""
	// log.Println("pReqRec add Req", pReqRec)
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IAAQ01", lErr1)
		return helpers.ErrReturn(lErr1)
	} else {
		defer lDb.Close()
		if pReqRec.Id != 0 {
			sqlString := `update a_ipo_userauth 
						set Type = ?,brokerMasterId=?,RoleId=?,ClientId=?,UpdatedBy = ?,UpdatedDate = now()  
						where Id = ?`

			_, lErr4 := lDb.Exec(sqlString, pReqRec.Status, pReqRec.BrokerId, pReqRec.RoleId, pReqRec.ClientId, pClientId, pReqRec.Id)
			if lErr4 != nil {
				log.Println("IAAQ04", lErr4)
				return lErr4
			}
		} else {

			sqlString := `insert into a_ipo_userauth (ClientId,RoleId,Type,brokerMasterId,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate) values (?,?,?,?,?,now(),?,now())`
			log.Println("Insert Values := ", pReqRec)
			_, lErr5 := lDb.Exec(sqlString, pReqRec.ClientId, pReqRec.RoleId, pReqRec.Status, pReqRec.BrokerId, pClientId, pClientId)
			if lErr5 != nil {
				log.Println("IAAQ05", lErr5)
				return lErr5
			}
		}
	}
	log.Println("addQuery (-)")
	return nil
}
