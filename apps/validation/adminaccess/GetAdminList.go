package adminaccess

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

type adminStruct struct {
	Id         int    `json:"id"`
	BrokerId   int    `json:"brokerId"`
	BrokerName string `json:"brokerName"`
	ClientId   string `json:"clientId"`
	RoleId     int    `json:"roleId"`
	RoleName   string `json:"roleName"`
	Status     string `json:"status"`
	Editable   bool   `json:"editable"`
}

type role struct {
	RoleId   int    `json:"roleId"`
	RoleName string `json:"roleName"`
}

type Broker struct {
	BrokerId   int    `json:"brokerId"`
	BrokerName string `json:"brokerName"`
}

type ResultRespStruct struct {
	AdminListArr  []adminStruct `json:"adminListArr"`
	RoleListArr   []role        `json:"roleListArr"`
	BrokerNameArr []Broker      `json:"brokerNameArr"`
	Status        string        `json:"status"`
	ErrMsg        string        `json:"errMsg"`
}

func GetAdminList(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAdminList (+)", r.Method)
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
		// This variable is used to store the records
		var lAdminRec adminStruct
		// This variable is used to send resp to frontend
		var lRespRec ResultRespStruct
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/config")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "AGAL01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("AGAL01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("AGAL01", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lRespRec.Status = common.SuccessCode
		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr2 != nil {
			log.Println("AGAL02", lErr2)
			fmt.Fprintf(w, helpers.GetErrorString("AGAL02", "Unable to establish db connection"))
			return
		} else {
			defer lDb.Close()

			// ! To get the Symbol, ApplicationNo, Price and Quantity
			lCoreString2 := `select u.id ,u.brokerMasterId ,b.BrokerName ,u.ClientId,u.RoleId,r.Role,u.Type 
			from a_ipo_userauth u, a_ipo_role r,a_ipo_brokermaster b 
			where u.RoleId = r.Id  and b.Id =u.brokerMasterId and r.Role = 'BrokerSuperAdmin'
			order by u.id `

			// and m.BiddingEndDate < curdate() Its goes befor client id in query
			lRows1, lErr3 := lDb.Query(lCoreString2)
			if lErr3 != nil {
				log.Println("AGAL03", lErr3)
				fmt.Fprintf(w, helpers.GetErrorString("AGAL03", "UserDetails not Found"))
				return
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows1.Next() {
					lAdminRec.Editable = false
					lErr4 := lRows1.Scan(&lAdminRec.Id, &lAdminRec.BrokerId, &lAdminRec.BrokerName, &lAdminRec.ClientId, &lAdminRec.RoleId, &lAdminRec.RoleName, &lAdminRec.Status)
					if lErr4 != nil {
						log.Println("AGAL04", lErr4)
						fmt.Fprintf(w, helpers.GetErrorString("AGAL04", "UserDetails not Found"))
						return
					} else {
						// Append the history Records into lhistoryArr Array
						lRespRec.AdminListArr = append(lRespRec.AdminListArr, lAdminRec)
						lRespRec.Status = common.SuccessCode
					}
				}
				lRoleArr, lErr5 := GetRoleList()
				if lErr5 != nil {
					log.Println("AGAL04", lErr5)
					fmt.Fprintf(w, helpers.GetErrorString("AGAL04", "UserDetails not Found"))
					return
				} else {
					for _, Role := range lRoleArr {
						// fmt.Println("Role", Role, Role.RoleName)
						if Role.RoleName == "BrokerSuperAdmin" {

							lRespRec.RoleListArr = append(lRespRec.RoleListArr, Role)
						}

					}
					// log.Println("lRespRec.RoleListArr", lRespRec.RoleListArr)
					// lRespRec.RoleListArr = lRoleArr
					lBrokerArr, lErr6 := GetBrokerList()
					if lErr6 != nil {
						log.Println("AGAL04", lErr6)
						fmt.Fprintf(w, helpers.GetErrorString("AGAL04", "UserDetails not Found"))
						return
					} else {
						lRespRec.BrokerNameArr = lBrokerArr
					}
				}
			}

		}
		lData, lErr5 := json.Marshal(lRespRec)
		if lErr5 != nil {
			log.Println("AGAL05", lErr5)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "AGAL05" + lErr5.Error()
			fmt.Fprintf(w, helpers.GetErrorString("AGAL05", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetAdminList (-)", r.Method)
	}
}

func GetRoleList() ([]role, error) {
	var lRole role
	var lRoleArr []role

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ARL01", lErr1)

		return lRoleArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select id,Role
			from a_ipo_role air `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("ARL02", lErr2)

			return lRoleArr, lErr2
		} else {
			//This for loop is used to collect the record from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRole.RoleId, &lRole.RoleName)
				if lErr3 != nil {
					log.Println("ARL03", lErr3)
					return lRoleArr, lErr3
				} else {
					lRoleArr = append(lRoleArr, lRole)
				}
			}
		}
	}
	return lRoleArr, nil
}

func GetBrokerList() ([]Broker, error) {
	var lBroker Broker
	var lBrokerArr []Broker

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AGBL01", lErr1)

		return lBrokerArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select aib.Id ,aib.BrokerName 
from a_ipo_brokermaster aib `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("AGBL02", lErr2)

			return lBrokerArr, lErr2
		} else {
			//This for loop is used to collect the record from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lBroker.BrokerId, &lBroker.BrokerName)
				if lErr3 != nil {
					log.Println("AGBL03", lErr3)

					return lBrokerArr, lErr3
				} else {
					lBrokerArr = append(lBrokerArr, lBroker)
				}
			}
		}
	}
	return lBrokerArr, nil
}
