package clientDetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/common"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

/*
Pupose:This method is used to get the email id  .
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	AGMPA45767,nil

	==========
	*On Error
	==========
	"",error

Author:Pavithra
Date: 29AUG2023
*/
// this method is commented by prashanth because of Login data is not inserting on DB So Cant Get Client Mail ID For Alternate purpose Another GetClent MAil fuction is created
// func GetClientEmailId(r *http.Request, pClientId string) (string, error) {
// 	log.Println("GetClientEmailId (+)")

// 	// this variables is used to get Pan number of the client from the database.
// 	var lEmailId string

// 	publicTokenCookie, lErr1 := r.Cookie(common.ABHICookieName)
// 	if lErr1 != nil {
// 		log.Println("CDGCE01", lErr1)
// 		return lEmailId, lErr1
// 	}

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr2 != nil {
// 		log.Println("CDGCE02", lErr2)
// 		return lEmailId, lErr2
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select UserMailId
// 						from novo_token
// 						where Token = ?
// 						and UserId = ? `
// 		lRows, lErr3 := lDb.Query(lCoreString, publicTokenCookie.Value, pClientId)
// 		if lErr3 != nil {
// 			log.Println("CDGCE03", lErr3)
// 			return lEmailId, lErr3
// 		} else {
// 			for lRows.Next() {
// 				lErr3 = lRows.Scan(&lEmailId)
// 				if lErr3 != nil {
// 					log.Println("CDGCE04", lErr3)
// 					return lEmailId, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetClientEmailId (-)")
// 	return lEmailId, nil
// }

type EmailStruct struct {
	EmailId string `json:"emailId"`
}

func GetClientEmailId(pClientId string) (string, error) {
	log.Println("GetClientEmailId (+)")
	config := common.ReadTomlConfig("./toml/emailconfig.toml")
	loginurl := fmt.Sprintf("%v", config.(map[string]interface{})["LoginUrl"])

	log.Println("clientid", pClientId)
	login := loginurl + pClientId
	var ClientDetails []EmailStruct
	var lEmailId string
	var lHeaderArr []apiUtil.HeaderDetails
	// ==========================
	//create parameters struct for LogEntry method
	var lLogInputRec Function.ParameterStruct
	lLogInputRec.Request = pClientId
	lLogInputRec.EndPoint = "/GetClientEmailId"
	lLogInputRec.Flag = common.INSERT
	lLogInputRec.ClientId = pClientId
	lLogInputRec.Method = "POST"

	// ! LogEntry method is used to store the Request in Database
	lId, lErr1 := Function.LogEntry(lLogInputRec)
	if lErr1 != nil {
		log.Println("CDGCE01", lErr1)
		return lEmailId, lErr1
	} else {
		EmailId, lErr2 := apiUtil.Api_call(login, "GET", "", lHeaderArr, "Emailfetch")
		if lErr2 != nil {
			log.Println("CDGCE02", lErr2)
			return lEmailId, lErr2
		} else {
			lErr3 := json.Unmarshal([]byte(EmailId), &ClientDetails)
			if lErr2 != nil {
				log.Println("NESM02", lErr3)
				return lEmailId, lErr3
			} else {
				if ClientDetails != nil {
					for _, Mail := range ClientDetails {
						// log.Println("Mail", Mail)
						lEmailId = Mail.EmailId
					}
				} else {

					log.Println("No Details Available for ClienMail")
					//  Return Error Message Needed
				}
				lLogInputRec.Response = lEmailId
				lLogInputRec.LastId = lId
				lLogInputRec.Flag = common.UPDATE
				// create instance to hold errors
				var lErr4 error
				lId, lErr4 = Function.LogEntry(lLogInputRec)
				if lErr4 != nil {
					log.Println("NSM04", lErr4)
					return lEmailId, lErr4
				}
			}
		}
	}
	log.Println("GetClientEmailId (-)")
	return lEmailId, nil
}
