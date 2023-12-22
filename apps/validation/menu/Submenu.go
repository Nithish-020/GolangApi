package menu

import (
	"encoding/json"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type SubMenuStruct struct {
	SubMenuArr []RouterStruct `json:"subMenuArr"`
	Status     string         `json:"status"`
	ErrMsg     string         `json:"errMsg"`
}

/*
Purpose:This api Method is used to get SubMenu Id,Name and path
Request:

Header Value: ID

Response:
=========
*On Sucess
=========

	{
		"subMenuArr" :[
			{
			"routerId" : 1
			"router" : "Ipo"
			"path" : "/ipo"
			},
			{
			"routerId" : 2
			"router" : "Sgb"
			"path" : "/sgb"
			}
		],
	"status" : "S",
	"errMsg" : ""
	}

=========
!On Error
=========

	{
		"subMenuArr" :[],
		"status": "E",
		"errMsg": "Can't able to get data from database"
	}

Author: Nithish Kumar
Date: 16AUG2023
*/
func GetSubMenu(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSubMenu(+)", r.Method)
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Origin", origin)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "ID,NO,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		// Create an instance for SubMenuStruct
		var lRespRec SubMenuStruct

		lRespRec.Status = common.SuccessCode

		// reading Header Value from the Request
		lParentId := r.Header.Get("ID")

		// log.Println("header value", lMasterId, lAppNo)
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lSessionClientId := ""
		lClientId := ""
		lLoggedBy := ""
		var lErr1 error
		lSessionClientId, lErr1 = appsso.ValidateAndGetClientDetails2(r, common.ABHIAppName, common.ABHICookieName)
		if lErr1 != nil {
			log.Println("MGSM01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "MGSM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("MGSM01", "UserDetails Not Found!"))
			return
		} else {
			if lSessionClientId != "" {
				//get the detail for the client for whome we need to work
				lClientId = common.GetSetClient(lSessionClientId)
				//get the staff who logged
				lLoggedBy = common.GetLoggedBy(lSessionClientId)
				log.Println(lLoggedBy, lClientId)
			} else {
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lRouterParentId, _ := strconv.Atoi(lParentId)
		//call getHistoryDetails method to get the application Details
		lRespArr, lErr2 := GetSubRouters(lRouterParentId, lClientId)
		if lErr2 != nil {
			log.Println("MGSM02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "MGSM02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("MGSM02", "Unable to fetch your Records"))
			return
		} else {
			lRespRec.SubMenuArr = lRespArr
			lRespRec.Status = common.SuccessCode
			// log.Println("SubMenu ", lRespRec.SubMenuArr)
		}
		// Marshal the response structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("MGSM03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("MGSM03", "Issue in Getting your SubMenus !"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetSubMenu (-)", r.Method)
	}
}

/*
Pupose:This method used to retrieve the RouterId,Router name,path  from the database.
Parameters:
PClientId,pMasterId
Response:
==========
*On Sucess
==========

	[
		{
			"routerId" : 1
			"router" : "Ipo"
			"path" : "/ipo"
		},
		{
			"routerId" : 2
			"router" : "Sgb"
			"path" : "/sgb"
		}
	],

==========
!On Error
==========

	[],error

Author:Nithish Kumar
Date: 16AUG2023
*/
func GetSubRouters(pParentId int, pClientId string) ([]RouterStruct, error) {
	log.Println("getSubRouters (+)")

	// This Variable is used to store the Each Modify bid Records
	var lRespRec RouterStruct
	var lRouterArr []RouterStruct
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MGSR01", lErr1)
		return lRouterArr, lErr1
	} else {
		defer lDb.Close()
		lParentId := strconv.Itoa(pParentId)
		lCoreString := `select va.routerId ,va.routerName,va.routerLink
						from v_ipo_usertaskmenu va
						where va.ParentId = ` + lParentId + `
						and va.clientId = '` + pClientId + `'`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("MGSR02", lErr2)
			return lRouterArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in RouterStructArray
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRespRec.RouterId, &lRespRec.Router, &lRespRec.Path)
				if lErr3 != nil {
					log.Println("MGSR03", lErr3)
					return lRouterArr, lErr3
				} else {
					// Append the Bid Records in lRespArr
					lRouterArr = append(lRouterArr, lRespRec)
				}
			}
		}
	}
	log.Println("getSubRouters (-)")
	return lRouterArr, nil
}
