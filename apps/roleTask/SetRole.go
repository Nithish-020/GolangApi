package roleTask

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

type setRole struct {
	Id       int    `json:"id"`
	RoleName string `json:"roleName"`
	Status   string `json:"status"`
}

type setRoleResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetRole(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "PUT" {
		log.Println("setRole (+)")
		var lRoleReq setRole
		var lRoleResp setRoleResp

		lRoleResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/createtask")
		if lErr1 != nil {
			lRoleResp.Status = common.ErrorCode
			lRoleResp.ErrMsg = "RSR01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSR01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lRoleResp.Status = common.ErrorCode
				lRoleResp.ErrMsg = "RSR02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RSR02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lRoleResp.Status = common.ErrorCode
			lRoleResp.ErrMsg = "RSR03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSR03", "Request cannot be empty.Please try after sometime"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lRoleReq)
			if lErr4 != nil {
				lRoleResp.Status = common.ErrorCode
				lRoleResp.ErrMsg = "RSR04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RSR04", "Unable to proccess your request right now.Please try after sometime"))
				return
			} else {

				if lRoleReq.Id == 0 {
					lRowcount, lErr5 := insertRole(lRoleReq, lClientId)
					if lErr5 != nil {
						log.Println("RSR05", lErr5)
						fmt.Fprintf(w, helpers.GetErrorString("RSR05", "Error in Inserting  Role"))
						return
					} else {
						if lRowcount == 0 {
							lRoleResp.Status = common.ErrorCode
							lRoleResp.ErrMsg = "Role Already Exist"
						} else {
							lRoleResp.Status = common.SuccessCode
							lRoleResp.ErrMsg = "Role added successfully"
						}

					}
				} else {

					lErr6 := updateRole(lRoleReq, lClientId)
					if lErr6 != nil {
						log.Println("RSR06", lErr6)
						fmt.Fprintf(w, helpers.GetErrorString("RSR06", "Error in Updating Role"))
						return
					} else {
						lRoleResp.Status = common.SuccessCode
						lRoleResp.ErrMsg = "Record updated successfully"

					}
				}
			}
		}
		lData, lErr7 := json.Marshal(lRoleResp)
		if lErr7 != nil {
			log.Println("RSR07", lErr7)
			fmt.Fprintf(w, helpers.GetErrorString("RSR07", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetRoleAndTAsk (-)")
	}
}

func insertRole(pRoleReq setRole, pClientId string) (int64, error) {
	log.Println("insertRole(+)")
	var lCount int64
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RIR01", lErr1)
		return lCount, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `insert into a_ipo_role (role,type,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
		select ?,?,?,now(),?,now()
		where not exists (select * from a_ipo_role  where  Role = ? )
		 limit 1;`
		lROWS, lErr2 := lDb.Exec(lCoreString, pRoleReq.RoleName, pRoleReq.Status, pClientId, pClientId, pRoleReq.RoleName)
		if lErr2 != nil {
			log.Println("RIR02", lErr2)
			return lCount, lErr2
		} else {
			var lErr3 error
			lCount, lErr3 = lROWS.RowsAffected()

			if lErr3 != nil {
				log.Println("RIR03", lErr3)
				return lCount, lErr3
			}

		}
	}
	log.Println("insertRole (-)")
	return lCount, nil
}

func updateRole(pRoleReq setRole, pClientId string) error {
	log.Println("updateRole(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RUR01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lCoreString := `update a_ipo_role 
			set role= ? ,type = ?,UpdatedBy =?,UpdatedDate =now()
			where Id =? `
		_, lErr2 := lDb.Exec(lCoreString, pRoleReq.RoleName, pRoleReq.Status, pClientId, pRoleReq.Id)
		if lErr2 != nil {
			log.Println("RUR02", lErr2)
			return lErr2
		}
	}
	log.Println("updateRole(-)")

	return nil

}
