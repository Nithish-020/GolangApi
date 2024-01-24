package nseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

type loginReqStruct struct {
	Member   string `json:"member"`
	LogId    string `json:"loginId"`
	Password string `json:"password"`
}

type loginRespStruct struct {
	CurrentTime string `json:"currentTime"`
	LoginId     string `json:"loginId"`
	Member      string `json:"member"`
	Status      string `json:"status"`
	Token       string `json:"token"`
}

func Token(pUser string, pBrokerId int) (loginRespStruct, error) {
	log.Println("Token....(+)")
	// Create instance for Parameter struct
	// var lLogInputRec Function.ParameterStruct
	// Create instance for loinRespStruct
	var lApiRespRec loginRespStruct
	// create instance to hold the last inserted id
	// var lId int
	//To link the toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Login"])
	// To get the details for v1/login from database
	lDetail, lErr := AccessNseCredential(pBrokerId)
	if lErr != nil {
		log.Println("NT01", lErr)
		return lApiRespRec, lErr
	} else {
		// Marshalling the structure for LogEntry method
		// lRequest, lErr := json.Marshal(lDetail)
		// if lErr != nil {
		// 	log.Println("NT01", lErr)
		// 	return lApiRespRec, lErr
		// } else {
		// 	lLogInputRec.Request = string(lRequest)
		// 	lLogInputRec.EndPoint = "/v1/login"
		// 	lLogInputRec.Flag = common.INSERT
		// 	lLogInputRec.ClientId = pUser
		// 	lLogInputRec.Method = "POST"

		// 	// LogEntry method is used to store the Request in Database
		// 	lId, lErr = Function.LogEntry(lLogInputRec)
		// 	if lErr != nil {
		// 		log.Println("NT02", lErr)
		// 		return lApiRespRec, lErr
		// 	} else {
		// TokenApi method used to call exchange API
		lResp, lErr := TokenApi(lDetail, lUrl)
		if lErr != nil {
			log.Println("NT03", lErr)
			return lApiRespRec, lErr
		} else {
			lApiRespRec = lResp
		}
		// Store thre Response in Log table
		// lResponse, lErr := json.Marshal(lResp)
		// if lErr != nil {
		// 	log.Println("NT04", lErr)
		// 	return lApiRespRec, lErr
		// } else {
		// 	lLogInputRec.Response = string(lResponse)
		// 	lLogInputRec.LastId = lId
		// 	lLogInputRec.Flag = common.UPDATE

		// 	lId, lErr = Function.LogEntry(lLogInputRec)
		// 	if lErr != nil {
		// 		log.Println("NT05", lErr)
		// 		return lApiRespRec, lErr
		// 	}
		// }
		// }
		// }
	}
	log.Println("Token....(-)")
	return lApiRespRec, nil
}

func TokenApi(pUser loginReqStruct, pUrl string) (loginRespStruct, error) {
	log.Println("TokenApi....(+)")
	//create instance for loginResp struct
	var lUserRespRec loginRespStruct
	//create array of instance to store the key value pairs
	var lHeaderArr []apiUtil.HeaderDetails

	// Marshall the structure parameter into json
	ljsonData, lErr := json.Marshal(pUser)
	if lErr != nil {
		log.Println("NTA01", lErr)
		return lUserRespRec, lErr
	} else {
		// convert json data into string
		lReqstring := string(ljsonData)
		lResp, lErr := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "Token")
		if lErr != nil {
			log.Println("NTA02", lErr)
			return lUserRespRec, lErr
		} else {
			// Unmarshalling json to struct
			lErr = json.Unmarshal([]byte(lResp), &lUserRespRec)
			if lErr != nil {
				log.Println("NTA03", lErr)
				return lUserRespRec, lErr
			}
		}
	}
	log.Println("TokenApi....(-)")
	return lUserRespRec, nil
}

func AccessNseCredential(pBrokerId int) (loginReqStruct, error) {
	log.Println("AccessNseCredential (+)")
	// Create instance for loginStruct
	var lUserRec loginReqStruct
	//Establish a Datatbase Connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IAAQ01", lErr1)
		return lUserRec, lErr1
	} else {
		defer lDb.Close()
		lCoreString1 := `select Member,LoginId,Password 
						from a_ipo_directory dir
						where dir.Status  = 'Y'
						and dir.Stream = 'NSE'
						and dir.brokerMasterId = ?`

		lRows, lErr2 := lDb.Query(lCoreString1, pBrokerId)
		if lErr2 != nil {
			log.Println("IAAQ02", lErr2)
			return lUserRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lUserRec.Member, &lUserRec.LogId, &lUserRec.Password)
				if lErr3 != nil {
					log.Println("IAAQ03", lErr3)
					return lUserRec, lErr3
				} else {
					lDecoded, lErr4 := common.DecodeToString(lUserRec.Password)
					if lErr4 != nil {
						log.Println("IAAQ03", lErr4)
						return lUserRec, lErr4
					} else {
						lUserRec.Password = lDecoded
					}
				}
			}
		}
	}
	log.Println("AccessNseCredential (-)")
	return lUserRec, nil
}
