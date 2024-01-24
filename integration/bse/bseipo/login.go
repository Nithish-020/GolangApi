package bseipo

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

// {
// 	"membercode":"1004",
// 	"loginid":"mkt01",
// 	"password":"rij03@uo",
// 	"ibbsid":"ibtest104y"
// }

type BseloginStruct struct {
	MemberCode string `json:"membercode"`
	LoginId    string `json:"loginid"`
	Password   string `json:"password"`
	IbbsId     string `json:"ibbsid"`
}

type BseLoginRespStruct struct {
	MemberCode string `json:"membercode"`
	LoginId    string `json:"loginid"`
	BranchCode string `json:"branchcode"`
	Token      string `json:"token"`
	ErrorCode  string `json:"errorcode"`
	Message    string `json:"message"`
}

func BseToken(pUser string, pBrokerId int) (BseLoginRespStruct, error) {
	log.Println("BseToken (+)")
	// Create instance for Parameter struct
	// var lLogInputRec Function.ParameterStruct
	// Create instance for loinRespStruct
	var lApiRespRec BseLoginRespStruct
	// create instance to hold the last inserted id
	// var lId int
	//To link the toml file
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseLogin"])
	// To get the details for v1/login from database
	lDetail, lErr := AccessBseCredential(pBrokerId)
	if lErr != nil {
		log.Println("BBT01", lErr)
		return lApiRespRec, lErr
	} else {
		// Marshalling the structure for LogEntry method
		// lRequest, lErr := json.Marshal(lDetail)
		// if lErr != nil {
		// 	log.Println("BBT01", lErr)
		// 	return lApiRespRec, lErr
		// } else {
		// 	lLogInputRec.Request = string(lRequest)
		// 	lLogInputRec.EndPoint = "bse/v1/login"
		// 	lLogInputRec.Flag = common.INSERT
		// 	lLogInputRec.ClientId = pUser
		// 	lLogInputRec.Method = "POST"

		// 	// LogEntry method is used to store the Request in Database
		// 	lId, lErr = Function.LogEntry(lLogInputRec)
		// 	if lErr != nil {
		// 		log.Println("BBT02", lErr)
		// 		return lApiRespRec, lErr
		// 	} else {
		// TokenApi method used to call exchange API
		lResp, lErr := BseTokenApi(lDetail, lUrl)
		if lErr != nil {
			log.Println("BBT03", lErr)
			return lApiRespRec, lErr
		} else {
			lApiRespRec = lResp
		}
		// Store thre Response in Log table
		// lResponse, lErr := json.Marshal(lResp)
		// if lErr != nil {
		// 	log.Println("BBT04", lErr)
		// 	return lApiRespRec, lErr
		// } else {
		// 	lLogInputRec.Response = string(lResponse)
		// 	lLogInputRec.LastId = lId
		// 	lLogInputRec.Flag = common.UPDATE

		// 	lId, lErr = Function.LogEntry(lLogInputRec)
		// 	if lErr != nil {
		// 		log.Println("BBT05", lErr)
		// 		return lApiRespRec, lErr
		// 	}
		// }
		// }
		// }
	}
	log.Println("BseToken (-)")
	return lApiRespRec, nil
}

func BseTokenApi(pUser BseloginStruct, pUrl string) (BseLoginRespStruct, error) {
	log.Println("BseTokenApi (+)")
	//create instance for loginResp struct
	var lUserRespRec BseLoginRespStruct
	//create array of instance to store the key value pairs
	var lHeaderArr []apiUtil.HeaderDetails

	// Marshall the structure parameter into json
	ljsonData, lErr1 := json.Marshal(pUser)
	if lErr1 != nil {
		log.Println("BBTA01", lErr1)
		return lUserRespRec, lErr1
	} else {
		// convert json data into string
		lReqstring := string(ljsonData)
		lResp, lErr2 := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "BseToken")
		if lErr2 != nil {
			log.Println("BBTA02", lErr2)
			return lUserRespRec, lErr2
		} else {
			// Unmarshalling json to struct
			lErr3 := json.Unmarshal([]byte(lResp), &lUserRespRec)
			if lErr3 != nil {
				log.Println("BBTA03", lErr3)
				return lUserRespRec, lErr3
			}
		}
	}
	log.Println("BseTokenApi (-)")
	return lUserRespRec, nil
}

func AccessBseCredential(pBrokerId int) (BseloginStruct, error) {
	log.Println("AccessDetail (+)")
	// Create instance for loginStruct
	var lUserRec BseloginStruct
	//Establish a Datatbase Connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BAD01", lErr1)
		return lUserRec, lErr1
	} else {
		defer lDb.Close()
		lCoreString1 := `select dir.Member ,dir.LoginId ,dir.Password ,dir.ibbsid 
		from a_ipo_directory dir
		where dir.Status  = 'Y'
		and dir.Stream = 'BSE'
		and dir.brokerMasterId  = ?`

		lRows, lErr2 := lDb.Query(lCoreString1, pBrokerId)
		if lErr2 != nil {
			log.Println("BAD02", lErr2)
			return lUserRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lUserRec.MemberCode, &lUserRec.LoginId, &lUserRec.Password, &lUserRec.IbbsId)
				if lErr3 != nil {
					log.Println("BAD03", lErr3)
					return lUserRec, lErr3
				} else {
					lDecoded, lErr4 := common.DecodeToString(lUserRec.Password)
					if lErr4 != nil {
						log.Println("BAD03", lErr4)
						return lUserRec, lErr4
					} else {
						lUserRec.Password = lDecoded
					}
				}
			}
		}
	}
	log.Println("AccessDetail (-)")

	return lUserRec, nil
}
