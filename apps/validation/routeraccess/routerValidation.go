package routeraccess

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/menu"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

type RouterVerify struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
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
Date: 21JULY2023
/*/
func RouterAccess(w http.ResponseWriter, r *http.Request) {
	log.Println("RouterAccess(+)", r.Method)
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
	(w).Header().Set("Access-Control-Allow-Headers", "ROUTER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		var lRespRec RouterVerify

		lRouter := r.Header.Get("ROUTER")
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lSessionClientId := ""
		lClientId := ""
		lLoggedBy := ""
		var lErr1 error
		lSessionClientId, lErr1 = appsso.ValidateAndGetClientDetails2(r, common.ABHIAppName, common.ABHICookieName)
		if lErr1 != nil {
			log.Println("RRA01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "RRA01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH01", "UserDetails not Found"))
			return
		} else {
			if lSessionClientId != "" {
				//get the detail for the client for whome we need to work
				lClientId = common.GetSetClient(lSessionClientId)
				//get the staff who logged
				lLoggedBy = common.GetLoggedBy(lSessionClientId)
				log.Println(lLoggedBy, lClientId)
			} else {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "RRA02" + lErr1.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RRA02", "Issue in Fetching Your Datas!"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// lRouterArr, lErr2 := adminaccess.GetRouter(lClientId, lRouter)
		// lRouterArr, lErr2 := menu.GetRouter(lClientId, lRouter)
		lRouterArr, lErr2 := menu.DefaultRouter(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("RRA03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "RRA03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RRA03", "Issue in Fetching Your Datas!"))
			return
		} else {
			for _, v := range lRouterArr {
				if v.Path == "" {
					lTemp, lErr3 := menu.GetSubRouters(v.RouterId, lClientId)
					if lErr3 != nil {
						log.Println("RRA03", lErr3)
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "RRA03" + lErr3.Error()
						fmt.Fprintf(w, helpers.GetErrorString("RRA03", "Issue in Fetching Your Datas!"))
						return
					} else {
						for _, T := range lTemp {

							lRouterArr = append(lRouterArr, T)
						}
						lTemp = nil
					}
				}

			}
			// if lRouterArr != nil {
			if len(lRouterArr) > 0 {
				lRespRec.Status = common.SuccessCode
				for lIdx := 0; lIdx < len(lRouterArr); lIdx++ {

					if lRouterArr[lIdx].Path == lRouter {
						lRespRec.Status = common.SuccessCode
						break
					} else {
						lRespRec.Status = common.ErrorCode
					}
				}
			} else {
				lRespRec.Status = common.ErrorCode
			}
		}
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("RRA04", lErr3)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "RRA04" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RRA04", "Issue in Routing"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("RouterAccess (-)", r.Method)
	}
}
