package bsesgb

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bseipo"
	"fcs23pkg/util/apiUtil"
	"fmt"
	"log"
)

func BseSgbToken(pUser string, pBrokerId int) (bseipo.BseLoginRespStruct, error) {
	log.Println("BseSgbToken (+)")
	// Create instance for Parameter struct
	// var lLogInputRec Function.ParameterStruct
	// Create instance for loinRespStruct
	var lApiRespRec bseipo.BseLoginRespStruct
	// create instance to hold the last inserted id
	// var lId int

	//To link the toml file  added new comments
	lConfigFile := common.ReadTomlConfig("./toml/config.toml")
	lUrl := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseLogin"])
	// To get the details for v1/login from database
	lDetail, lErr := bseipo.AccessBseCredential(pBrokerId)
	if lErr != nil {
		log.Println("BSBST01", lErr)
		return lApiRespRec, lErr
	} else {
		// Marshalling the structure for LogEntry method
		// lRequest, lErr := json.Marshal(lDetail)
		// if lErr != nil {
		// 	log.Println("BSBST01", lErr)
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
		// 		log.Println("BSBST02", lErr)
		// 		return lApiRespRec, lErr
		// 	} else {
		// TokenApi method used to call exchange API
		lResp, lErr := BseSgbTokenApi(lDetail, lUrl)
		if lErr != nil {
			log.Println("BSBST03", lErr)
			return lApiRespRec, lErr
		} else {
			lApiRespRec = lResp
		}
		// Store thre Response in Log table
		// lResponse, lErr := json.Marshal(lResp)
		// if lErr != nil {
		// 	log.Println("BSBST04", lErr)
		// 	return lApiRespRec, lErr
		// } else {
		// 	lLogInputRec.Response = string(lResponse)
		// 	lLogInputRec.LastId = lId
		// 	lLogInputRec.Flag = common.UPDATE

		// 	lId, lErr = Function.LogEntry(lLogInputRec)
		// 	if lErr != nil {
		// 		log.Println("BSBST05", lErr)
		// 		return lApiRespRec, lErr
		// 	}
		// }
		// }
		// }
	}
	log.Println("BseSgbToken (-)")
	return lApiRespRec, nil
}

func BseSgbTokenApi(pUser bseipo.BseloginStruct, pUrl string) (bseipo.BseLoginRespStruct, error) {
	log.Println("BseSgbTokenApi (+)")
	//create instance for loginResp struct
	var lUserRespRec bseipo.BseLoginRespStruct
	//create array of instance to store the key value pairs
	var lHeaderArr []apiUtil.HeaderDetails
	// create a new instance of HeaderDetails struct
	var lConsHeadRec apiUtil.HeaderDetails

	lConsHeadRec.Key = "Content-Type"
	lConsHeadRec.Value = "application/json"
	lHeaderArr = append(lHeaderArr, lConsHeadRec)

	// Marshall the structure parameter into json
	ljsonData, lErr1 := json.Marshal(pUser)
	if lErr1 != nil {
		log.Println("BSBSTA01", lErr1)
		return lUserRespRec, lErr1
	} else {
		// convert json data into string
		lReqstring := string(ljsonData)
		lResp, lErr2 := apiUtil.Api_call(pUrl, "POST", lReqstring, lHeaderArr, "BseToken")
		if lErr2 != nil {
			log.Println("BSBSTA02", lErr2)
			return lUserRespRec, lErr2
		} else {
			// Unmarshalling json to struct
			lErr3 := json.Unmarshal([]byte(lResp), &lUserRespRec)
			if lErr3 != nil {
				log.Println("BSBSTA03", lErr3)
				return lUserRespRec, lErr3
			}
		}
	}
	log.Println("BseSgbTokenApi (-)")
	return lUserRespRec, nil
}

func AccessDetail() (bseipo.BseloginStruct, error) {
	log.Println("AccessDetail (+)")
	// Create instance for loginStruct
	var lUserRec bseipo.BseloginStruct
	//Establish a Datatbase Connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BSAD01", lErr1)
		return lUserRec, lErr1
	} else {
		defer lDb.Close()
		lCoreString1 := `select dir.Member ,dir.LoginId ,dir.Password ,dir.ibbsid 
		from a_ipo_directory dir
		where dir.Status  = 'Y'
		and dir.Stream = 'BSE'`

		lRows, lErr2 := lDb.Query(lCoreString1)
		if lErr2 != nil {
			log.Println("BSAD02", lErr2)
			return lUserRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lUserRec.MemberCode, &lUserRec.LoginId, &lUserRec.Password, &lUserRec.IbbsId)
				if lErr3 != nil {
					log.Println("BSAD03", lErr3)
					return lUserRec, lErr3
				} else {
					lDecoded, lErr4 := common.DecodeToString(lUserRec.Password)
					if lErr4 != nil {
						log.Println("BSAD04", lErr4)
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
