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

type roleMaster struct {
	RoleId      int    `json:"roleId"`
	RoleName    string `json:"roleName"`
	Type        string `json:"type"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedDate string `json:"updatedDate"`
}
type role struct {
	RoleId   int    `json:"roleId"`
	RoleName string `json:"roleName"`
}

type Task struct {
	TaskId     int    `json:"taskId"`
	Router     string `json:"router"`
	RouterName string `json:"routerName"`
	ParentId   int    `json:"parentId"`
}

type taskMaster struct {
	// Router Id is named as Task Id
	TaskId      int    `json:"taskId"`
	RouterName  string `json:"routerName"`
	Router      string `json:"router"`
	ParentId    int    `json:"parentId"`
	Type        string `json:"type"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedDate string `json:"updatedDate"`
}
type getRoleTaskResp struct {
	RoleMasterArr []roleMaster `json:"roleMasterArr"`
	TaskMasterArr []taskMaster `json:"taskMasterArr"`
	RoleListArr   []role       `json:"roleListArr"`
	TaskListArr   []Task       `json:"taskListArr"`
	Status        string       `json:"status"`
	ErrMsg        string       `json:"errMsg"`
}

func GetRoleTask(w http.ResponseWriter, r *http.Request) {
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

		log.Println("GetRoleTask (+)")
		var lGetRoleTask getRoleTaskResp
		lGetRoleTask.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/createtask")
		if lErr1 != nil {
			log.Println("LDR01", lErr1)
			lGetRoleTask.Status = common.ErrorCode
			lGetRoleTask.ErrMsg = "LDR01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRT01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("RGRT02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lGetRole, lErr2 := getRole()
		if lErr2 != nil {
			log.Println("RGRT01", lErr2)
			lGetRoleTask.Status = common.ErrorCode
			lGetRoleTask.ErrMsg = lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRT03", "Role Details not Found"))
			return
		} else {
			lGetRoleTask.RoleMasterArr = lGetRole

			lGetTask, lErr3 := getTask()
			if lErr3 != nil {
				log.Println("RGRT02", lErr3)
				lGetRoleTask.Status = common.ErrorCode
				lGetRoleTask.ErrMsg = lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RGRT04", "Task Details not Found"))
				return

			} else {
				lGetRoleTask.TaskMasterArr = lGetTask

				lTaskList, lErr4 := GetTaskList()
				if lErr4 != nil {
					log.Println("RGRT03", lErr4)
					lGetRoleTask.Status = common.ErrorCode
					lGetRoleTask.ErrMsg = lErr4.Error()
					fmt.Fprintf(w, helpers.GetErrorString("RGRT05", "Task DropDown Details not Found"))
					return

				} else {
					lGetRoleTask.TaskListArr = lTaskList
				}

			}
		}
		lData, lErr5 := json.Marshal(lGetRoleTask)
		if lErr5 != nil {
			log.Println("RGRT06", lErr5)
			fmt.Fprintf(w, helpers.GetErrorString("RGRT03", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("GetRoleTask (-)")
}

func getRole() ([]roleMaster, error) {
	log.Println("getRole(+)")
	var lRole roleMaster
	var lRoleArr []roleMaster

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RGR01", lErr1)
		return lRoleArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select id,Role ,type,CreatedBy ,date_format(CreatedDate, '%d-%b-%y, %l:%i %p') ,UpdatedBy ,date_format(UpdatedDate, '%d-%b-%y, %l:%i %p') 
		               from a_ipo_role  `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("RGR02", lErr2)
			return lRoleArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRole.RoleId, &lRole.RoleName, &lRole.Type, &lRole.CreatedBy, &lRole.CreatedDate, &lRole.UpdatedBy, &lRole.UpdatedDate)
				if lErr3 != nil {
					log.Println("RGR03", lErr3)
					return lRoleArr, lErr3
				} else {

					lRoleArr = append(lRoleArr, lRole)
				}
			}
		}
	}
	log.Println("getRole(-)")

	return lRoleArr, nil

}

func getTask() ([]taskMaster, error) {
	log.Println("getTask(+) ")
	var lTask taskMaster
	var lTaskArr []taskMaster

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RGT01", lErr1)
		return lTaskArr, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `select id,RouterName ,nvl(Router,'') ,nvl(ParentId,0) ,type ,nvl(CreatedBy,''),
		nvl(date_format(CreatedDate, '%d-%b-%y, %l:%i %p'),''),nvl(UpdatedBy,'')  ,nvl(date_format(UpdatedDate, '%d-%b-%y, %l:%i %p'),'') 	from a_ipo_router  `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("RGT02", lErr2)
			return lTaskArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lTask.TaskId, &lTask.RouterName, &lTask.Router, &lTask.ParentId, &lTask.Type, &lTask.CreatedBy, &lTask.CreatedDate, &lTask.UpdatedBy, &lTask.UpdatedDate)
				if lErr3 != nil {
					log.Println("RGT03", lErr3)
					return lTaskArr, lErr3
				} else {

					lTaskArr = append(lTaskArr, lTask)
				}
			}
		}
	}
	log.Println("getTask(-) ")
	return lTaskArr, nil

}
func GetRoleList() ([]role, error) {
	log.Println("GetRoleTask (+)")
	var lRole role
	var lRoleArr []role

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RGRL01", lErr1)

		return lRoleArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select id,Role
			from a_ipo_role air `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("RGRL02", lErr2)

			return lRoleArr, lErr2
		} else {
			//This for loop is used to collect the record from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRole.RoleId, &lRole.RoleName)
				if lErr3 != nil {
					log.Println("RGRL03", lErr3)
					return lRoleArr, lErr3
				} else {
					lRoleArr = append(lRoleArr, lRole)
				}
			}
		}
	}
	log.Println("GetRoleTask (-)")
	return lRoleArr, nil
}

func GetTaskList() ([]Task, error) {
	log.Println(" GetTaskList(+)")
	var lTask Task
	var lTaskArr []Task

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RGTL01", lErr1)

		return lTaskArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select id,nvl(Router,""), RouterName,nvl(ParentId,0)
		from a_ipo_router 
		where type = 'Y'`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("RGTL02", lErr2)

			return lTaskArr, lErr2
		} else {
			//This for loop is used to collect the record from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lTask.TaskId, &lTask.Router, &lTask.RouterName, &lTask.ParentId)
				if lErr3 != nil {
					log.Println("RGTL03", lErr3)
					return lTaskArr, lErr3
				} else {
					lTaskArr = append(lTaskArr, lTask)
				}
			}
		}
	}
	log.Println(" GetTaskList(-)")
	return lTaskArr, nil
}
