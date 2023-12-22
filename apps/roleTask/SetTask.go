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
	// "strconv"
)

type setTask struct {
	TaskId     int    `json:"taskId"`
	RouterName string `json:"routerName"`
	Router     string `json:"router"`
	ParentId   int    `json:"parentId"`
	Status     string `json:"status"`
}

type setTaskResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetTask(w http.ResponseWriter, r *http.Request) {
	log.Println("SetTask (+)", r.Method)

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
		var lTaskReq setTask
		var lTaskResp setTaskResp

		lTaskResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/createtask")
		if lErr1 != nil {
			lTaskResp.Status = common.ErrorCode
			lTaskResp.ErrMsg = "RST01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RST01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lTaskResp.Status = common.ErrorCode
				lTaskResp.ErrMsg = "RST02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RST02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lTaskResp.Status = common.ErrorCode
			lTaskResp.ErrMsg = "RST03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RST03", "Request cannot be empty.Please try after sometime"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lTaskReq)

			if lErr4 != nil {
				lTaskResp.Status = common.ErrorCode
				lTaskResp.ErrMsg = "RST04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RST04", "Unable to proccess your request right now.Please try after sometime"))
				return
			} else {

				if lTaskReq.TaskId == 0 {
					lRowcount, lErr5 := insertTask(lTaskReq, lClientId)
					if lErr5 != nil {
						log.Println("RST05", lErr5)
						fmt.Fprintf(w, helpers.GetErrorString("RST05", "Error in Inserting  Task Method"))
						return
					} else {

						if lRowcount == 0 {
							lTaskResp.Status = common.ErrorCode
							lTaskResp.ErrMsg = "Task Already Exist"
						} else {
							lTaskResp.Status = common.SuccessCode
							lTaskResp.ErrMsg = "Task added successfully"
							if lTaskReq.ParentId != 0 {
								lErr6 := updateRouter(lTaskReq)
								if lErr6 != nil {
									log.Println("RST06", lErr6)
									fmt.Fprintf(w, helpers.GetErrorString("RST06", "Error in Updating Router  After inserting Data"))
									return
								}
							}
						}

					}
				} else {

					lErr7 := updateTask(lTaskReq, lClientId)
					if lErr7 != nil {
						log.Println("RST07", lErr7)
						fmt.Fprintf(w, helpers.GetErrorString("RST07", "Error in Updating Task Method"))
						return
					} else {
						lTaskResp.Status = common.SuccessCode
						lTaskResp.ErrMsg = "Record updated successfully"

					}
				}
			}
		}
		lData, lErr8 := json.Marshal(lTaskResp)
		if lErr8 != nil {
			log.Println("RST08", lErr8)
			fmt.Fprintf(w, helpers.GetErrorString("RST08", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SetTask (-)")
	}
}

func insertTask(pTaskReq setTask, pClientId string) (int64, error) {
	log.Println("insertTask(+)")
	var lCountTask int64
	var ParentId *int
	var lRouter *string
	if pTaskReq.ParentId != 0 {
		ParentId = &pTaskReq.ParentId
	} else {
		ParentId = nil
	}
	if pTaskReq.Router == "" {
		lRouter = nil
	} else {
		lRouter = &pTaskReq.Router
	}

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RIT01", lErr1)
		return lCountTask, lErr1
	} else {
		defer lDb.Close()

		lCoreString := ` insert into a_ipo_router (RouterName,Router,ParentId,type,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
		select ?,?,?,?,?,now(),?,now()
		where not exists (select * from a_ipo_router where  RouterName = ? or Router = ? )
		limit 1; `

		lROWS, lErr2 := lDb.Exec(lCoreString, pTaskReq.RouterName, lRouter, ParentId, pTaskReq.Status, pClientId, pClientId, pTaskReq.RouterName, lRouter)
		if lErr2 != nil {
			log.Println("RIT02", lErr2)
			return lCountTask, lErr2
		} else {
			var lErr3 error
			lCountTask, lErr3 = lROWS.RowsAffected()
			if lErr3 != nil {
				log.Println("RIT03", lErr3)
				return lCountTask, lErr3
			}
		}
	}

	log.Println("insertTask (-)")
	return lCountTask, nil

}

func updateTask(pTaskReq setTask, pClientId string) error {
	log.Println("updateTask(+)")
	var ParentId *int
	if pTaskReq.ParentId != 0 {
		ParentId = &pTaskReq.ParentId
	} else {
		ParentId = nil
	}
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RUT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lCoreString := `update a_ipo_router 
		set RouterName = ?,Router = ?,ParentId = ? ,type = ?,UpdatedBy = ?,UpdatedDate = now()
		where id = ?`
		_, lErr2 := lDb.Exec(lCoreString, pTaskReq.RouterName, pTaskReq.Router, ParentId, pTaskReq.Status, pClientId, pTaskReq.TaskId)
		if lErr2 != nil {
			log.Println("RUT02", lErr2)
			return lErr2
		} else {
			if ParentId != nil {
				lErr3 := updateRouter(pTaskReq)
				if lErr3 != nil {
					log.Println("RUT03", lErr3)
					return lErr3
				}

			}
		}
	}
	log.Println("updateTask(-)")

	return nil
}

func updateRouter(pTaskReq setTask) error {
	log.Println("updateRouter(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RUR01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lCoreString := `update a_ipo_router 
		set Router = '' 
		where  id = ?`
		_, lErr2 := lDb.Exec(lCoreString, pTaskReq.ParentId)
		if lErr2 != nil {
			log.Println("RUR02", lErr2)
			return lErr2
		}
	}
	log.Println("updateRouter(-)")

	return nil

}
