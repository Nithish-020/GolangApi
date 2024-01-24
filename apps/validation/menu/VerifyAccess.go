package menu

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RouterStruct struct {
	RouterId int    `json:"routerId"`
	Router   string `json:"router"`
	Path     string `json:"path"`
}

type RouterRespStruct struct {
	RouterArr []RouterStruct `json:"routerArr"`
	Status    string         `json:"status"`
	ErrMsg    string         `json:"errmsg"`
}

/*
Pupose:This Function is used to know the logged in user is admin or not“
Parameters:

	not applicable

Response:

	*On Sucess
	=========
	{
		"flag": "admin",
		"status": "success"
		"errmsg":""
	}

	!On Error
	========
	{	"flag":"NA",
		"status": "failed",
		"errmsg": "Issue in Fetching Your Datas"
	}

Author: Nithish Kumar
Date: 19JULY2023
*/
func VerifyAccess(w http.ResponseWriter, r *http.Request) {
	log.Println("VerifyAccess(+)", r.Method)

	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr1 error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			// TO get brokerId
			lBrokerId, lErr1 = brokers.GetBrokerId(origin)
			if lErr1 != nil {
				log.Println("Cross-Orign-Error", lErr1.Error())
			}
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {
		var lRespRec RouterRespStruct

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lSessionClientId := ""
		lClientId := ""
		lLoggedBy := ""
		lFlag := ""
		var lErr3 error

		lSessionClientId, lErr3 = appsso.ValidateAndGetClientDetails2(r, common.ABHIAppName, common.ABHICookieName)
		if lErr3 != nil {
			log.Println("AVA01", lErr3)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "AVA01" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH01", "UserDetails not Found"))
			return
		} else {
			if lSessionClientId != "" {
				log.Println("lSessionClientId", lSessionClientId)
				//get the detail for the client for whome we need to work
				lClientId = common.GetSetClient(lSessionClientId)
				//get the staff who logged
				lLoggedBy = common.GetLoggedBy(lSessionClientId)
				log.Println(lLoggedBy, lClientId)
			} else {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "AVA02" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("AVA02", "Issue in Fetching Your Datas!"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lId, lErr2 := brokers.GetMemberDetailId(lBrokerId)
		if lErr2 != nil {
			log.Println("AVA02", lErr2.Error())
		} else {
			if lId != 0 {
				lFlag = "Y"
			} else {
				lFlag = "N"
			}
			log.Println("lFlag := ", lFlag)
			// lRouterArr, lErr2 := GetRouter(lClientId, "")
			lRouterArr, lCode, lErr3 := SendConfig(lClientId, lFlag, lBrokerId) // Condition check for broker Router
			if lErr3 != nil {
				log.Println("AVA03", lErr3)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "AVA03" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("AVA03", "Issue in Fetching Your Datas!"))
				return
			} else {
				if lCode != "E" {
					lRespRec.RouterArr = lRouterArr
					lRespRec.Status = common.SuccessCode
					// log.Println("Finally RouterArr", lRespRec.RouterArr)
				} else {
					fmt.Fprintf(w, helpers.GetErrorString("ASC03", "Access Restricted !"))
					return
				}
			}
		}
		//----------------------
		lData, lErr5 := json.Marshal(lRespRec)
		if lErr5 != nil {
			log.Println("AVA04", lErr5)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "AVA04" + lErr5.Error()
			fmt.Fprintf(w, helpers.GetErrorString("AVA04", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("VerifyAccess (-)", r.Method)
	}
}

func GetClientRouter(pClientId string, pBrokerId int) ([]RouterStruct, error) {
	log.Println("GetClientRouter (+)")
	// To store the
	var lRouterRec RouterStruct
	var lRouterArr []RouterStruct

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AGR01", lErr1)
		return lRouterArr, lErr1
	} else {
		defer lDb.Close()
		// lCoreString2 := `
		// select aro.id, aro.routername,aro.router
		// from a_ipo_role air , a_ipo_router aro, a_ipo_taskrole ait, a_ipo_userauth aiu
		// where  air.id = ait.roleid
		// and ait.routerid = aro.id
		// and aiu.roleid = air.id
		// and aiu.type = 'Y'
		// and aiu.type = ait.type
		// and ait.type = aro.type
		// and aro.type = air.type
		// and aiu.clientid = ?
		// and aiu.brokerMasterId  = ?
		// and aro.ParentId is null
		// `

		lBrokerId := strconv.Itoa(pBrokerId)
		lCoreString2 := `select routerId, routerName,routerLink
		from v_ipo_usermainmenu vm
		where vm.clientId = '` + pClientId + `'
		 and vm.brokerId = ` + lBrokerId

		lCoreString := lCoreString2

		// lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("AGR02", lErr2)
			return lRouterArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRouterRec.RouterId, &lRouterRec.Router, &lRouterRec.Path)
				if lErr3 != nil {
					log.Println("AGR03", lErr3)
					return lRouterArr, lErr3
				} else {
					lRouterArr = append(lRouterArr, lRouterRec)
				}
			}
		}
	}
	log.Println("GetClientRouter (-)")
	return lRouterArr, nil
}

// func GetRouter(pClientId string, path string) ([]RouterStruct, error) {
// 	log.Println("GetRouter (+)")
// 	// To store the
// 	var lRouterRec RouterStruct
// 	var lRouterArr []RouterStruct

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("AGR01", lErr1)
// 		return lRouterArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		// lCoreString := `select r.RouterName ,r.Router
// 		// from a_ipo_router r,(select tr.RouterId Id
// 		// from a_ipo_taskrole tr
// 		// where tr.Type = 'Y' and tr.RoleId = (select case when 1 = count(1) then air.Id else id.role end
// 		// from a_ipo_role air, (select nvl(
// 		// (select a.RoleId
// 		// from a_ipo_userauth a
// 		// where a.Type = 'Active' and a.ClientId = ?),'Client')role)id
// 		// where air.Role = id.role)) rout
// 		// where r.Id = rout.Id
// 		// order by r.Id `

// 		lCoreString1 := `
// 		select aro.id, aro.routername,aro.router
// 		from a_ipo_role air , a_ipo_router aro, a_ipo_taskrole ait
// 		where air.role = 'Client'
// 		and air.id = ait.roleid
// 		and ait.routerid = aro.id
// 		and ait.type = 'Y'
// 		and ait.type = aro.type
// 		and aro.type = air.type
// 		`
// 		lCoreString2 := `
// 		union
// 		select aro.id, aro.routername,aro.router
// 		from a_ipo_role air , a_ipo_router aro, a_ipo_taskrole ait, a_ipo_userauth aiu
// 		where  air.id = ait.roleid
// 		and ait.routerid = aro.id
// 		and aiu.roleid = air.id
// 		and aiu.type = 'Y'
// 		and aiu.type = ait.type
// 		and ait.type = aro.type
// 		and aro.type = air.type
// 		and aiu.clientid = ?
// 		and aro.ParentId is null
// 		`
// 		// and aro.ParentId is null

// 		ladditionalCond1 := ""
// 		ladditionalCond2 := ""

// 		if path != "" {
// 			ladditionalCond1 = ` and 1=2`
// 			ladditionalCond2 = ` and aro.router = '` + path + `'`
// 		} else {
// 			ladditionalCond2 = `and aro.ParentId is null`
// 		}

// 		lCoreString := lCoreString1 + ladditionalCond1 + lCoreString2 + ladditionalCond2

// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
// 		if lErr2 != nil {
// 			log.Println("AGR02", lErr2)
// 			return lRouterArr, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lRouterRec.RouterId, &lRouterRec.Router, &lRouterRec.Path)
// 				if lErr3 != nil {
// 					log.Println("AGR03", lErr3)
// 					return lRouterArr, lErr3
// 				} else {
// 					lRouterArr = append(lRouterArr, lRouterRec)
// 				}

// 			}

// 		}
// 	}
// 	log.Println("GetRouter (-)")
// 	return lRouterArr, nil
// }

// ================================================================//
// This method is used to send the config Router if flag is N     //
// ==============================================================//

func SendConfig(pClientid string, pFlag string, pBrokerId int) ([]RouterStruct, string, error) {
	log.Println("sendConfig (+)")
	var lconfigArr []RouterStruct
	var lRouterArr []RouterStruct
	// To indicated the execution
	lCode := "S"

	lRole, lErr1 := GetClientRole(pClientid, pBrokerId)
	if lErr1 != nil {
		log.Println("ASC01", lErr1)
		return lRouterArr, lCode, lErr1
	} else {
		var lErr2 error
		lRouterArr, lErr2 = GetClientRouter(pClientid, pBrokerId)
		if lErr2 != nil {
			log.Println("ASC02", lErr2)
			return lRouterArr, lCode, lErr2
		} else {
			// log.Println("pFlag ", pFlag, lRole, lRouterArr)
			if pFlag == "N" {
				if lRole == "SuperAdmin" || lRole == "BrokerSuperAdmin" {
					// log.Println("Check SuperAdmin || BrokerSuperAdmin", pFlag, lRole, lRouterArr)
					for _, Router := range lRouterArr {
						if Router.Router == "Config" {

							lconfigArr = append(lconfigArr, Router)
							lRouterArr = lconfigArr

							// log.Println("Condition satisfied ", lRouterArr)
							return lRouterArr, lCode, nil
						}
					}
				} else if lRole != "SuperAdmin" || lRole != "BrokerSuperAdmin" {
					log.Println("ASC03", "No Access to the Application for : "+lRole+" Access Restricted !")
					lRouterArr = nil
					lCode = "E"
					return lRouterArr, lCode, nil
				}
			} else {
				lconfigArr, lErr3 := DefaultRouter(pClientid, pBrokerId)
				if lErr3 != nil {
					log.Println("ASC04", lErr3)
					return lRouterArr, lCode, lErr3
				} else {
					lRouterArr = lconfigArr
					// log.Println("lRouterArr************", lRouterArr, lconfigArr)
				}
			}
		}
	}

	log.Println("sendConfig (-)")
	return lRouterArr, lCode, nil
}

// =============================================================================================//
//  DefaultRouter ---> This method is used to get the defaullt values for clients as per broker//
// ============================================================================================//

func DefaultRouter(pClientId string, pBrokerId int) ([]RouterStruct, error) {
	log.Println("DefaultRouter (+)")
	var lDefaultRouterArr []RouterStruct //To get string value of Allow modules from Db
	// var Temp []RouterStruct
	var lRouterRec RouterStruct

	// strings.Split() this method split the value and return it in Array format
	lSplitRouter, lErr1 := GetAllowModules(pBrokerId)
	if lErr1 != nil {
		log.Println("ADR04", lErr1)
		return lDefaultRouterArr, lErr1
	} else {
		// If the login user is Client send defalut modules
		lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr1 != nil {
			log.Println("ADR01", lErr1)
			return lDefaultRouterArr, lErr1
		} else {
			defer lDb.Close()

			//COMMENTED BY NITHISH
			// lCoreString3 := `select aro.id, aro.routername,aro.router
			// 	from a_ipo_role air , a_ipo_router aro, a_ipo_taskrole ait
			// 	where  air.id = ait.roleid
			// 	and ait.routerid = aro.id
			// 	and ait.type = aro.type
			// 	and aro.type = air.type
			// 	and aro.RouterName  in `

			//ADDED BY NITHISH
			// THE CONDITION WHERE HANDLED USING VIEW

			lBrokerId := strconv.Itoa(pBrokerId)
			lCoreString3 := `	select va.routerId ,va.routerName,va.routerLink
								from v_ipo_usermainmenu va
								where va.clientId = '` + pClientId + `'
								and va.brokerId = ` + lBrokerId + `
								or va.routerName in`

			var lFinalQuery string
			for lIdx := 0; lIdx < len(lSplitRouter); lIdx++ {
				lHelper := "'" + lSplitRouter[lIdx] + "',"
				lFinalQuery = lFinalQuery + lHelper
			}
			// lCoreString3 = lCoreString3 + "(" + lFinalQuery[0:len(lFinalQuery)-1] + ") group by aro.Id"
			lCoreString3 = lCoreString3 + "(" + lFinalQuery[0:len(lFinalQuery)-1] + ") group by va.routerId"
			lRows1, lErr2 := lDb.Query(lCoreString3)
			if lErr2 != nil {
				log.Println("ADR02", lErr2)
				return lDefaultRouterArr, lErr2
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows1.Next() {
					lErr3 := lRows1.Scan(&lRouterRec.RouterId, &lRouterRec.Router, &lRouterRec.Path)
					if lErr3 != nil {
						log.Println("ADR03", lErr3)
						return lDefaultRouterArr, lErr3
					} else {
						lDefaultRouterArr = append(lDefaultRouterArr, lRouterRec)
					}
				}
			}
		}
		//COMMENTED BY NITHISH
		// THIS BELOW STEP ARE HANDLED IN THE ABOVE QUERY
		// //menu.GetRouter method Returns the all Router value from DB
		// lGetAllRouter, lErr2 := GetClientRouter(pClientId, pBrokerId)
		// if lErr2 != nil {
		// 	log.Println("ADR04", lErr2)
		// 	return lDefaultRouterArr, lErr2

		// } else {
		// 	log.Println("lGetAllRouter , lSplitRouter", lGetAllRouter, lSplitRouter)
		// 	// log.Println(lSplitRouter, "lsplit + clientRouterArr", lGetAllRouter)
		// 	for _, Router := range lGetAllRouter {
		// 		for _, DefaultRouter := range lSplitRouter {

		// 			if Router.Router == DefaultRouter {
		// 				lDefaultRouterArr = append(lDefaultRouterArr, Router)
		// 			}
		// 		}
		// 	}
		// 	// log.Println("llDefaultRouterArr", lDefaultRouterArr, lGetAllRouter)

		// 	if lDefaultRouterArr != nil {
		// 		// log.Println(lDefaultRouterArr, lGetAllRouter, "Array Check")
		// 		for _, Router := range lGetAllRouter {

		// 			if "IPO" != Router.Router && "SGB" != Router.Router {
		// 				lDefaultRouterArr = append(lDefaultRouterArr, Router)
		// 			}
		// 		}

		// 		// resultArray := removeDuplicates(lDefaultRouterArr)
		// 		// log.Println("Result Array Check := ", resultArray)
		// 	} else {
		// log.Println("lDefaultRouterArr := ", lDefaultRouterArr)
		// 	}
		// }

	}
	log.Println("DefaultRouter (-)")
	return lDefaultRouterArr, nil
}

func GetClientRole(pClientId string, pBrokerId int) (string, error) {
	log.Println("GetClientRole (+)")
	// To store the clients role
	var lRole string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MGCR01", lErr1)
		return lRole, lErr1
	} else {
		defer lDb.Close()
		lCoreString1 := `select case when count(1) >0 then air.Role else 'NA' end
		from a_ipo_userauth aiu ,a_ipo_role air 
		where air.Id = aiu.RoleId 
		and aiu.ClientId = ?
		and aiu.brokerMasterId = ?
		and air.type = aiu.Type`
		lCoreString := lCoreString1
		// log.Println("lCoreString := ", lCoreString, pClientId, pBrokerId)
		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pBrokerId)
		if lErr2 != nil {
			log.Println("MGCR02", lErr2)
			return lRole, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRole)
				if lErr3 != nil {
					log.Println("MGCR03", lErr3)
					return lRole, lErr3
				} else {
					log.Println("Role", lRole)
				}
			}
		}
	}
	log.Println("GetClientRole (-)")
	return lRole, nil

}

func GetAllowModules(pBrokerId int) ([]string, error) {
	log.Println("GetAllowModules(+)")
	var lgetModules string
	var lgetAllModules []string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ADR01", lErr1)
		return lgetAllModules, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `select nvl(upper(m.AllowedModules),'')
			from a_ipo_memberdetails m
			where m.Brokerid = ?
			and m.Flag = 'Y'`

		lRows1, lErr2 := lDb.Query(lCoreString2, pBrokerId)
		if lErr2 != nil {
			log.Println("ADR02", lErr2)
			return lgetAllModules, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lgetModules)
				if lErr3 != nil {
					log.Println("ADR03", lErr3)
					return lgetAllModules, lErr3
				}
			}
			// log.Println("Inside if lDefaultRouter")

			// strings.Split() this method split the value and return it in Array format
			lgetAllModules = strings.Split(lgetModules, "/")

			log.Println("GetAllowModules(-)")
		}
	}
	return lgetAllModules, nil
}

func removeDuplicates(inputArray []RouterStruct) []RouterStruct {
	uniqueMap := make(map[string]bool)
	result := []RouterStruct{}

	for _, value := range inputArray {
		if _, exists := uniqueMap[value.Router]; !exists {
			uniqueMap[value.Router] = true
			result = append(result, value)
		}
	}

	return result
}
