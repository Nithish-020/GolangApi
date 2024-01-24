package clientDetail

import (
	"encoding/json"
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

// type EmailStruct struct {
// 	EmailId string `json:"emailId"`
// }

// [{"client_dp_code":"1208030000741751","client_dp_name":"  KARTHIKRAJA","emailId":"KARTHIK2768@YMAIL.COM","pan_no":"EXPPK4076L"}]

type EmailStruct struct {
	EmailId        string `json:"emailId"`
	Client_dp_code string `json:"client_dp_code"`
	Client_dp_name string `json:"client_dp_name"`
	Pan_no         string `json:"pan_no"`
}

func GetClientEmailId(pClientId string) (EmailStruct, error) {
	log.Println("GetClientEmailId (+)")
	config := common.ReadTomlConfig("toml/novoConfig.toml")
	loginurl := fmt.Sprintf("%v", config.(map[string]interface{})["LoginUrl"])

	login := loginurl + pClientId
	var lClientdetails EmailStruct
	var ClientDetails []EmailStruct
	var lHeaderArr []apiUtil.HeaderDetails
	// ==========================
	//create parameters struct for LogEntry method
	// var lLogInputRec Function.ParameterStruct
	// lLogInputRec.Request = pClientId
	// lLogInputRec.EndPoint = "/?qc=CLEML&typ=v&qp1=" + pClientId
	// lLogInputRec.Flag = common.INSERT
	// lLogInputRec.ClientId = pClientId
	// lLogInputRec.Method = "POST"

	// ! LogEntry method is used to store the Request in Database
	// lId, lErr1 := Function.LogEntry(lLogInputRec)
	// if lErr1 != nil {
	// 	log.Println("CDGCE01", lErr1)
	// 	return lClientdetails, lErr1
	// } else {
	lResp, lErr2 := apiUtil.Api_call(login, "GET", "", lHeaderArr, "ClientDetails")
	if lErr2 != nil {
		log.Println("CDGCE02", lErr2)
		return lClientdetails, lErr2
	} else {
		lErr3 := json.Unmarshal([]byte(lResp), &ClientDetails)
		if lErr2 != nil {
			log.Println("CDGCE03", lErr3)
			return lClientdetails, lErr3
		} else {
			if ClientDetails == nil {
				lErr := common.CustomError("Client Details not found")
				return lClientdetails, lErr
			} else {
				for _, Mail := range ClientDetails {
					// log.Println("Mail", Mail)
					lClientdetails = Mail
				}
				// lClientdetails.Pan_no = ""
				// lResponse, lErr4 := json.Marshal(ClientDetails)
				// if lErr4 != nil {
				// 	log.Println("CDGCE04", lErr4)
				// 	return lClientdetails, lErr4
				// } else {
				// lLogInputRec.Response = string(lResponse)
				// lLogInputRec.LastId = lId
				// lLogInputRec.Flag = common.UPDATE
				// _, lErr5 := Function.LogEntry(lLogInputRec)
				// if lErr5 != nil {
				// 	log.Println("CDGCE05", lErr5)
				// 	return lClientdetails, lErr5
				// }
				// }
			}

		}
	}
	// }
	log.Println("GetClientEmailId (-)")
	return lClientdetails, nil
}
