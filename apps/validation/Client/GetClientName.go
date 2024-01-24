package client

import (
	"encoding/json"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/helpers"
	"fcs23pkg/util/clientUtil"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ClientStruct struct {
	ClientName string `json:"clientName"`
	Status     string `json:"status"`
	ErrMsg     string `json:"errMsg"`
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------
/*
Pupose:This Function is used to know the logged in user is admin or not“
Parameters:

	not applicable

Response:

	*On Sucess
	=========
	{
		"ClientName": " Nithish Kumar",
		"status": "S",
		"errmsg":""
	}

	!On Error
	========
	{	"ClientName":"",
		"status": "E",
		"errmsg": "Client Id not found"
	}

Author: Nithish Kumar
Date: 20SEP2023
*/
func GetClientName(w http.ResponseWriter, r *http.Request) {
	log.Println("GetClientName(+)", r.Method)
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
		var lRespRec ClientStruct
		lRespRec.Status = common.SuccessCode
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lSessionClientId := ""
		lClientId := ""
		lLoggedBy := ""
		var lErr2 error
		lSessionClientId, lErr2 = appsso.ValidateAndGetClientDetails2(r, common.ABHIAppName, common.ABHICookieName)
		if lErr2 != nil {
			log.Println("CGCN01", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "CGCN01" + lErr2.Error()
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
				lRespRec.ErrMsg = "CGCN02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("CGCN02", "Issue in Fetching Your Datas!"))
				return
			}
		}

		lClientName, lErr3 := clientUtil.GetClientName(lClientId)
		if lErr3 != nil {
			log.Println("CGCN03", lErr3)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "CGCN03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH01", "Client details not Found"))
			return
		} else {
			// lRespRec.ClientName = lClientName
			// log.Println("ClientName", lClientName, lClientId)

			// Split the name into words
			words := strings.Fields(lClientName)

			// Extract the first letter from each word
			var lInitials string
			for _, word := range words {
				lInitials += string(word[0])
			}
			lRespRec.ClientName = lInitials

			// log.Println("Final ClientName", lRespRec.ClientName)

		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lData, lErr4 := json.Marshal(lRespRec)
		if lErr4 != nil {
			log.Println("CGCN04", lErr4)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "CGCN04" + lErr4.Error()
			fmt.Fprintf(w, helpers.GetErrorString("CGCN04", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetClientName (-)", r.Method)
	}
}
