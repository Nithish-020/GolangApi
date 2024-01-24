package apiaccess

import (
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"log"
	"net/http"
)

/*
Pupose: This method is used to get the currently active directory
Parameters:
	not applicable
Response:
	    ==========
	    *On Sucess
	    ==========
	    {

			clientId: "IPOA001",
			status: "S",
			errMSg: ""
		}
	    =========
	    !On Error
	    =========
	       {
			clientId: "",
			errMSg: "Error"
		}
Author: Nithish Kumar
Date: 24JULY2023
*/

func VerifyApiAccess(pR *http.Request, pApp string, pPublicTokenCookieName string, pPath string) (string, error) {
	log.Println("VerifyApiAccess (+)")
	// TO capture the request Endpoint
	//! lRequestURL := pR.URL.Path
	lSessionClientId := ""
	lClientId := ""
	lLoggedBy := ""
	lAccess := ""
	var lErr error
	lSessionClientId, lErr = appsso.ValidateAndGetClientDetails2(pR, pApp, pPublicTokenCookieName)
	if lErr != nil {
		return lClientId, lErr
	} else {
		if lSessionClientId != "" {
			//get the detail for the client for whome we need to work
			lClientId = common.GetSetClient(lSessionClientId)
			//get the staff who logged
			lLoggedBy = common.GetLoggedBy(lSessionClientId)
			log.Println(lLoggedBy, lClientId)

			//! TO check the Api eligible for the clients role
			// lEligible, lErr := CheckApiEligible(lRequestURL, lClientId)
			// if lEligible != "N" {
			// TO check the router permission for the clients role
			lAccess, lErr = CheckRoute(pPath, lClientId)
			if lErr != nil {
				return lClientId, lErr

			} else {
				if lAccess == "N" {
					lClientId = ""
					return lClientId, lErr
				}
			}
			// } else {
			// 	lClientId = ""
			// 	return lClientId, lErr
			// }
		}
	}
	log.Println("VerifyApiAccess (-)")
	return lClientId, nil
}

func CheckRoute(pPath string, pClientId string) (string, error) {
	log.Println("CheckRoute (+)")
	var lSignal string
	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	//if any error when opening db connection
	if err != nil {
		log.Println(err)
		return lSignal, err
	} else {
		defer db.Close()
		CoreString := `select case when aro.Router <> '' then 'Y' else 'N' end
		from a_ipo_role air , a_ipo_router aro, a_ipo_taskrole ait, a_ipo_userauth aiu
			where  air.id = ait.roleid
			and ait.routerid = aro.id
			and aiu.roleid = air.id
			and aiu.type = 'Y'
			and aiu.type = ait.type
			and ait.type = aro.type
			and aro.type = air.type
			and aro.Router = ?
			and aiu.clientid = ?
		union
			select case when aro.Router <> '' then 'Y' else 'N' end
			from a_ipo_role air , a_ipo_router aro, a_ipo_taskrole ait
			where air.role = 'Client'
			and air.id = ait.roleid
			and ait.routerid = aro.id
			and ait.type = 'Y'
			and ait.type = aro.type
			and aro.type = air.type	
			and aro.Router = ?
		`

		rows, err := db.Query(CoreString, pPath, pClientId, pPath)
		if err != nil {
			log.Println(err)
			return lSignal, err
		} else {
			// lSignal = "N"
			for rows.Next() {
				err := rows.Scan(&lSignal)
				if err != nil {
					log.Println(err)
					return lSignal, err
				} else {
					log.Println("Router Access: " + lSignal)
				}
			}
		}
	}
	log.Println("CheckRoute (-)")
	return lSignal, nil
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------

// func CheckApiEligible(pPath string, pClientId string) (string, error) {
// 	log.Println("CheckApiEligible (+)")
// 	var lEligible string
// 	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
// 	//if any error when opening db connection
// 	if err != nil {
// 		log.Println(err)
// 		return lEligible, err
// 	} else {
// 		defer db.Close()
// 		CoreString := `select case when count(1) = 1 then 'Y' else 'N' end
// 		from a_ipo_role air ,a_ipo_userauth aiu ,a_ipo_apiendpoints aia
// 		where air.Id = aiu.RoleId
// 		and air.Id = aia.RoleId
// 		and air.type = aia.Type
// 		and aia.EndPoint = ?
// 		and aiu.ClientId = ? `

// 		rows, err := db.Query(CoreString, pPath, pClientId)
// 		if err != nil {
// 			log.Println(err)
// 			return lEligible, err
// 		} else {
// 			for rows.Next() {
// 				err := rows.Scan(&lEligible)
// 				if err != nil {
// 					log.Println(err)
// 					return lEligible, err
// 				}
// 			}
// 		}
// 	}
// 	log.Println("CheckApiEligible (-)")
// 	return lEligible, nil

// }
