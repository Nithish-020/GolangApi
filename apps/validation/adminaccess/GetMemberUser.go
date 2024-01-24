package adminaccess

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/apps/validation/menu"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

type memberStruct struct {
	Id       int    `json:"id"`
	ClientId string `json:"clientId"`
	RoleId   int    `json:"roleId"`
	RoleName string `json:"roleName"`
	Status   string `json:"status"`
}

type memberRespStruct struct {
	MemberList  []memberStruct `json:"memberListArr"`
	RoleListArr []role         `json:"roleListArr"`
	Status      string         `json:"status"`
	ErrMsg      string         `json:"errMsg"`
}

func GetMemberUser(w http.ResponseWriter, r *http.Request) {
	log.Println("GetMemberUser (+)", r.Method)
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
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// This variable is used to store the records
		var lMemberRec memberStruct
		// This variable is used to send resp to frontend
		var lRespRec memberRespStruct

		var lMemberArr []memberStruct
		var lRoleTempArr []role
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/managerole")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "AGMU01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("AGMU01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("AGMU01", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr2 != nil {
			log.Println("AGMU02", lErr2)
			fmt.Fprintf(w, helpers.GetErrorString("AGMU02", "Unable to establish db connection"))
			return
		} else {
			defer lDb.Close()

			lLoginRole, lErr3 := menu.GetClientRole(lClientId, lBrokerId)
			log.Println("lLoggedRole", lLoginRole)
			if lErr3 != nil {
				log.Println("AGMU03", lErr3)
				fmt.Fprintf(w, helpers.GetErrorString("AGMU03", "Unable to verify User,please Logout and Login again !"))
				return
			} else {

				lRoleArr, lErr4 := GetRoleList()
				if lErr4 != nil {
					log.Println("AGMU04", lErr4)
					fmt.Fprintf(w, helpers.GetErrorString("AGMU04", "Issue in getting RoleList"))
					return
				} else {

					if lLoginRole != "SuperAdmin" {
						for _, lRole := range lRoleArr {
							if lRole.RoleName != "SuperAdmin" && lRole.RoleName != "BrokerSuperAdmin" {
								lRoleTempArr = append(lRoleTempArr, lRole)
							}
						}
						lRespRec.RoleListArr = lRoleTempArr
					} else {
						lRespRec.RoleListArr = lRoleArr
					}
					log.Println("Final Role List ", lRespRec.RoleListArr)

					lCoreString2 := `select u.id ,u.ClientId,u.RoleId,r.Role,u.Type 
					from a_ipo_userauth u, a_ipo_role r,a_ipo_brokermaster b 
					where u.RoleId = r.Id 
					and b.Id =u.brokerMasterId 
					and u.brokerMasterId = ?
					order by u.id  `

					// and m.BiddingEndDate < curdate() Its goes befor client id in query
					lRows1, lErr5 := lDb.Query(lCoreString2, lBrokerId)
					if lErr5 != nil {
						log.Println("AGMU05", lErr5)
						fmt.Fprintf(w, helpers.GetErrorString("AGMU05", "UserDetails not Found"))
						return
					} else {
						//This for loop is used to collect the records from the database and store them in structure
						for lRows1.Next() {
							lErr6 := lRows1.Scan(&lMemberRec.Id, &lMemberRec.ClientId, &lMemberRec.RoleId, &lMemberRec.RoleName, &lMemberRec.Status)
							if lErr6 != nil {
								log.Println("AGMU06", lErr6)
								fmt.Fprintf(w, helpers.GetErrorString("AGMU06", "UserDetails not Found"))
								return
							} else {
								// Append the history Records into lhistoryArr Array
								lMemberArr = append(lMemberArr, lMemberRec)
							}
						}

						if lLoginRole != "SuperAdmin" {
							for _, lUsers := range lMemberArr {
								if lUsers.RoleName != "SuperAdmin" && lUsers.RoleName != "BrokerSuperAdmin" {
									lRespRec.MemberList = append(lRespRec.MemberList, lUsers)
								}
							}
						} else {
							lRespRec.MemberList = lMemberArr
						}
						log.Println("Final MemberList", lRespRec.MemberList)
						lRespRec.Status = common.SuccessCode
					}
				}
			}
		}
		lData, lErr7 := json.Marshal(lRespRec)
		if lErr7 != nil {
			log.Println("AGMU07", lErr7)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "AGMU07" + lErr7.Error()
			fmt.Fprintf(w, helpers.GetErrorString("AGMU07", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetMemberUser (-)", r.Method)
	}
}
