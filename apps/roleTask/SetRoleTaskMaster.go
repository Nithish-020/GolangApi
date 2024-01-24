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

type setRoleTask struct {
	TaskRoleId int    `json:"taskRoleId"`
	RoleId     int    `json:"roleId"`
	TaskId     int    `json:"taskId"`
	Status     string `json:"status"`
}

type setRoleTaskResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetRoleTaskMaster(w http.ResponseWriter, r *http.Request) {
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
		log.Println("SetRoleTaskConnector (+)")
		var lSetRoleTask setRoleTask
		var lSetRoleTaskResp setRoleTaskResp

		lSetRoleTaskResp.Status = common.SuccessCode

		// /-----------START OF GETTING CLIENT AND STAFF DETAILS----------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/assigntask")
		if lErr1 != nil {
			log.Println("RSRTM01", lErr1)
			lSetRoleTaskResp.Status = common.ErrorCode
			lSetRoleTaskResp.ErrMsg = "RSRTM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSRTM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				log.Println("RSRTM02", lErr2)
				lSetRoleTaskResp.Status = common.ErrorCode
				lSetRoleTaskResp.ErrMsg = "RSRTM02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RSRTM02", "Access restricted"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lSetRoleTaskResp.Status = common.ErrorCode
			lSetRoleTaskResp.ErrMsg = "RSRTM03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSRTM03", "Request cannot be empty.Please try after sometime"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lSetRoleTask)
			if lErr4 != nil {
				lSetRoleTaskResp.Status = common.ErrorCode
				lSetRoleTaskResp.ErrMsg = "RSRTM04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RSRTM04", "Unable to proccess your request right now.Please try after sometime"))
				return
			} else {

				if lSetRoleTask.TaskRoleId == 0 {
					lRowcount, lErr5 := insertRoleTask(lSetRoleTask, lClientId)
					if lErr5 != nil {
						log.Println("RSRTM05", lErr5)
						fmt.Fprintf(w, helpers.GetErrorString("RSRTM05", "Error in Inserting  Role"))
						return
					} else {
						if lRowcount == 0 {
							lSetRoleTaskResp.Status = common.ErrorCode
							lSetRoleTaskResp.ErrMsg = "Role and Task Already Assigned"
						} else {
							lSetRoleTaskResp.Status = common.SuccessCode
							lSetRoleTaskResp.ErrMsg = "Role and Task added successfully"
						}

					}
				} else {

					lErr6 := updateRoleTask(lSetRoleTask, lClientId)
					if lErr6 != nil {
						log.Println("RSRTM06", lErr6)
						fmt.Fprintf(w, helpers.GetErrorString("RSRTM06", "Error in Updating Role"))
						return
					} else {
						lSetRoleTaskResp.Status = common.SuccessCode
						lSetRoleTaskResp.ErrMsg = "Role and Task updated successfully"

					}
				}
			}
		}
		lData, lErr7 := json.Marshal(lSetRoleTaskResp)
		if lErr7 != nil {
			log.Println("RSRTM07", lErr7)
			fmt.Fprintf(w, helpers.GetErrorString("RSRTM07", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SetRoleTaskConnector (-)")
	}
}

func insertRoleTask(pRoleTask setRoleTask, pClientId string) (int64, error) {
	log.Println("insertRoleTask(+)")
	var lCount int64
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RIRT01", lErr1)
		return lCount, lErr1
	} else {
		defer lDb.Close()
		//  Query is not Having BrokerMaster Id If you Want Add it
		lCoreString := `insert into a_ipo_taskrole (RoleId,RouterId,type,CreatedBy,CreatedDate,UpdatedBy,
			UpdatedDate)
			select ?,?,?,?,now(),?,now()
			where not exists (select  * from a_ipo_taskrole where RoleId = ? and RouterId = ?)
			limit 1`
		lROWS, lErr2 := lDb.Exec(lCoreString, pRoleTask.RoleId, pRoleTask.TaskId, pRoleTask.Status, pClientId, pClientId, pRoleTask.RoleId, pRoleTask.TaskId)
		if lErr2 != nil {
			log.Println("RIRT02", lErr2)
			return lCount, lErr2
		} else {
			var lErr3 error
			lCount, lErr3 = lROWS.RowsAffected()

			if lErr3 != nil {
				log.Println("RIRT03", lErr3)
				return lCount, lErr3
			}

		}
	}
	log.Println("insertRoleTask (-)")
	return lCount, nil

}

func updateRoleTask(pRoleTask setRoleTask, pClientId string) error {
	log.Println("updateRoleTask(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RURT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lCoreString := `update a_ipo_taskrole  
		set RoleId = ? ,RouterId = ?,type = ?,UpdatedBy = ?,UpdatedDate = now()
		where Id = ? `
		_, lErr2 := lDb.Exec(lCoreString, pRoleTask.RoleId, pRoleTask.TaskId, pRoleTask.Status, pClientId, pRoleTask.TaskRoleId)
		if lErr2 != nil {
			log.Println("RURT01", lErr2)
			return lErr2
		}
	}
	log.Println("updateRoleTask(-)")

	return nil

}
