package roleTask

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

type getRoleTaskConnector struct {
	Id       int    `json:"id"`
	RoleId   int    `json:"roleId"`
	RoleName string `json:"roleName"`
	// Router Id is named as Task Id
	TaskId      int    `json:"taskId"`
	RouterName  string `json:"routerName"`
	Type        string `json:"type"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedDate string `json:"updatedDate"`
}

type getConectorResp struct {
	RoleTaskArr   []getRoleTaskConnector `json:"roleTaskArr"`
	TaskConnector []Task                 `json:"taskConnector"`
	RoleConnector []role                 `json:"roleConnector"`
	Status        string                 `json:"status"`
	ErrMsg        string                 `json:"errMsg"`
}

func GetRoleTaskMaster(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		log.Println("GetRoleTaskMaster (+)")
		var lGetConnector getConectorResp
		lGetConnector.Status = common.SuccessCode

		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/assigntask")
		if lErr1 != nil {
			log.Println("RGRTM01", lErr1)
			lGetConnector.Status = common.ErrorCode
			lGetConnector.ErrMsg = "RGRTM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRTM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				log.Println("RGRTM02", lErr2)
				lGetConnector.Status = common.ErrorCode
				lGetConnector.ErrMsg = "RGRTM02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RGRTM02", "Access restricted"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lGetRoleTaskConnector, lErr3 := getRoleTAskArr()
		if lErr3 != nil {
			log.Println("RGRTM03", lErr3)
			lGetConnector.Status = common.ErrorCode
			lGetConnector.ErrMsg = "RGRTM03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRTM03", "UserDetails not Found"))
			return

		} else {
			lGetConnector.RoleTaskArr = lGetRoleTaskConnector
			lRoleList, lErr4 := GetRoleList()
			if lErr4 != nil {
				log.Println("RGRTM04", lErr4)
				lGetConnector.Status = common.ErrorCode
				lGetConnector.ErrMsg = "RGRTM04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RGRTM04", "Role List not Found"))
				return
			} else {
				lGetConnector.RoleConnector = lRoleList
				lTaskList, lErr5 := GetTaskList()
				if lErr5 != nil {
					log.Println("RGRTM05", lErr5)
					lGetConnector.Status = common.ErrorCode
					lGetConnector.ErrMsg = "RGRTM05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("RGRTM05", "Task List not Found"))
					return
				} else {
					lGetConnector.TaskConnector = lTaskList
					lGetConnector.Status = common.SuccessCode

				}
			}
		}

		lData, lErr6 := json.Marshal(lGetConnector)
		if lErr6 != nil {
			log.Println("RGRTM06", lErr6)
			lGetConnector.Status = common.ErrorCode
			lGetConnector.ErrMsg = "RGRTM06" + lErr6.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRTM06", "Issue in Getting Your Datas.Try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetRoleTaskMaster (-)")
	}
}

func getRoleTAskArr() ([]getRoleTaskConnector, error) {
	log.Println("getRoleTAskArr (+)")
	var lGetRoleTask getRoleTaskConnector
	var lRoleTaskArr []getRoleTaskConnector
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RGRTA01", lErr1)
		return lRoleTaskArr, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `select T.id,R.Id,R.role,Ro.Id,Ro.RouterName,T.type,nvl(T.CreatedBy,''),nvl(date_format(T.CreatedDate, '%d-%b-%y, %l:%i %p') ,''),nvl(T.UpdatedBy ,''),
		nvl(date_format(T.UpdatedDate , '%d-%b-%y, %l:%i %p'),'')
		from a_ipo_taskrole T ,a_ipo_role R ,a_ipo_router Ro
		where T.RoleId =R.Id and T.RouterId = Ro.Id `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("RGRTA02", lErr2)
			return lRoleTaskArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lGetRoleTask.Id, &lGetRoleTask.RoleId, &lGetRoleTask.RoleName, &lGetRoleTask.TaskId, &lGetRoleTask.RouterName, &lGetRoleTask.Type, &lGetRoleTask.CreatedBy, &lGetRoleTask.CreatedDate, &lGetRoleTask.UpdatedBy, &lGetRoleTask.UpdatedDate)
				if lErr3 != nil {
					log.Println("RGRTA03", lErr3)
					return lRoleTaskArr, lErr3
				} else {

					lRoleTaskArr = append(lRoleTaskArr, lGetRoleTask)
				}
			}
		}
	}
	log.Println("getRoleTAskArr (+)")
	return lRoleTaskArr, nil

}
