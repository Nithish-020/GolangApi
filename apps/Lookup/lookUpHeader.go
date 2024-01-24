package lookup

import (
	"database/sql"
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// To store Look Up Header Records
type LookUpHeader struct {
	ID           string `json:"id"`
	User         string `json:"user"`
	Code         string `json:"code"`
	Description  string `json:"description"`
	Created_By   string `json:"createdBy"`
	Created_Date string `json:"createdDate"`
	Updated_By   string `json:"updatedBy"`
	Updated_Date string `json:"updatedDate"`
}

// lResponse Structure for GetLookUpHeader API
type HeaderResponse struct {
	HeaderArr []LookUpHeader `json:"header"`
	ErrMsg    string         `json:"errMsg"`
	Status    string         `json:"status"`
}

/*
Pupose:This Function is used to fetch the lookup header data from our database table ....
Request:

nil

lResponse:

	*On Sucess
	=========
	{
	[{
		"id": "1",
		"User": "xyz",
		"Code": "20",
		"Description": "hwsfuiqgu",
	}],
	status: "S",
	ErrMsg: "",
	}

	!On Error
	========
	{
		"status": "failed",
		"reason": "Application no does not exist"
	}

Author: Nithish Kumar
Date: 08Aug2023
*/
func GetLookUpHeader(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	log.Println("GetLookUpHeader(+)")
	if r.Method == "GET" {
		// This variable is used to store the LookUpHeader structure
		var lInput LookUpHeader
		// This variable is used to store the lResponse structure
		var lResp HeaderResponse
		lResp.Status = "S"
		// To Establish A database connection,call LocalDbConnect Method
		lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr1 != nil {
			lResp.Status = "E"
			lResp.ErrMsg = "Error: " + lErr1.Error()
		} else {
			defer lDb.Close()
			// ! To get the Code  and description
			CoreString := `select id,nvl(Code,''),nvl(description,'') ,nvl(createdBy, '') createdBy, nvl(createdDate, '') createdDate ,nvl(updatedBy, '') updatedBy ,nvl(updatedDate, '')  updatedDate
			from xx_lookup_header `
			lRows, lErr2 := lDb.Query(CoreString)
			if lErr2 != nil {
				lResp.Status = "E"
				log.Println("Err 1", lErr2)
				lResp.ErrMsg = "Error: " + lErr2.Error()
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows.Next() {
					lErr3 := lRows.Scan(&lInput.ID, &lInput.Code, &lInput.Description, &lInput.Created_By, &lInput.Created_Date, &lInput.Updated_By, &lInput.Updated_Date)
					if lErr3 != nil {
						lResp.Status = "E"
						lResp.ErrMsg = "Error: " + lErr3.Error()
						log.Println("Err 2", lResp.ErrMsg)
					} else {
						// Append the LookUpHeader Records into lResp HeaderArr Array
						lResp.HeaderArr = append(lResp.HeaderArr, lInput)
						lResp.Status = "S"
					}
				}
			}
		}
		// Marshal the reponse structure into lDatas
		lData, lErr4 := json.Marshal(lResp)
		if lErr4 != nil {
			fmt.Fprintf(w, "Error taking data"+lErr4.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("GetLookUpHeader(-)")
}

func verifyClient(r *http.Request) (string, error) {
	log.Println("verifyClient (+)")
	//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
	// lSessionClientId := ""
	lClientId := ""
	lLoggedBy := ""
	lSessionClientId, lErr1 := appsso.ValidateAndGetClientDetails2(r, common.ABHIAppName, common.ABHICookieName)
	if lErr1 != nil {
		log.Println("LVC01", lErr1)
		return lClientId, lErr1
	} else {
		if lSessionClientId != "" {
			log.Println("lSessionClientId", lSessionClientId)
			//get the detail for the client for whome we need to work
			lClientId = common.GetSetClient(lSessionClientId)
			//get the staff who logged
			lLoggedBy = common.GetLoggedBy(lSessionClientId)
			log.Println(lLoggedBy, lClientId)
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		log.Println("verifyClient(-)")
	}
	return lClientId, lErr1
}

/*
Pupose:This Function is used to Add new lookup header data into our database table ....
Request:

{
	"id": "1",
		"User": "xyz",
		"Code": "20",
		"Description": "hwsfuiqgu",
}

lResponse:

	*On Sucess
	=========
	{
	status: "S",
	ErrMsg: "",
	}

	!On Error
	========
	{
		"status": "failed",
		"reason": "Application no does not exist"
	}

Author: Nithish Kumar
Date: 08Aug2023
*/
func AddLookUpHeader(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	log.Println("AddLookUpHeader(+)")
	if r.Method == "PUT" {
		// This variable is used to store the LookUpHeader structure
		var lInput LookUpHeader
		// This variable is used to store the lResponse structure
		var lResp HeaderResponse
		lResp.Status = "S"
		body, lErr1 := ioutil.ReadAll(r.Body)
		// log.Println("body", string(body))
		if lErr1 != nil {
			lResp.Status = "E"
			log.Println("Err 1", lErr1)
			lResp.ErrMsg = "Error: 1" + lErr1.Error()
		} else {

			lErr2 := json.Unmarshal(body, &lInput)
			// lId = string(body)
			// log.Println("Input", lInput)
			if lErr2 != nil {
				lResp.Status = "E"
				log.Println("Err 2", lErr2)
				lResp.ErrMsg = "Error: 2" + lErr2.Error()
			} else {
				lClientId, lErr3 := verifyClient(r)
				if lErr3 != nil {
					lResp.Status = "E"
					log.Println("Err 3", lErr3)
					lResp.ErrMsg = "Error: 3" + lErr3.Error()
				} else {
					lInput.User = lClientId

					// To Establish A database connection,call LocalDbConnect Method
					lDb, lErr4 := ftdb.LocalDbConnect(ftdb.IPODB)
					if lErr4 != nil {
						lResp.Status = "E"
						log.Println("Err 4", lErr4)
						lResp.ErrMsg = "Error: 4" + lErr4.Error()
					} else {
						defer lDb.Close()
						// ! To insert the LookUpHeader code,description
						// CoreString := `insert into xx_lookup_header(Code ,description ,createdBy, createdDate, updatedBy, updatedDate)
						// values(?,?,?,Now(),?,Now())`
						// _, lErr3 := lDb.Exec(CoreString, &lInput.Code, &lInput.Description, &lInput.User, &lInput.User)
						// if lErr3 != nil {
						// 	lResp.Status = "E"
						// 	log.Println("Err 3", lResp.ErrMsg)
						// 	lResp.ErrMsg = "Error: 3" + lErr3.Error()
						// } else {
						// 	lResp.Status = "S"
						// }
						lExistsVal, lErr5 := ExistsId(lDb, lInput.Code)
						if lErr5 != nil {
							lResp.Status = "E"
							log.Println("Err 5", lResp.ErrMsg)
							lResp.ErrMsg = "Error: 5" + lErr5.Error()
						} else {
							if lExistsVal == "Y" {
								lResp.Status = "E"
								log.Println("Err 5: This Code is already exists", lResp.ErrMsg)
								lResp.ErrMsg = "This Code is already exists"
							} else if lExistsVal == "N" {
								lErr4 := InsertHeaderValue(lDb, lInput)
								if lErr4 != nil {
									lResp.Status = "E"
									log.Println("Err 4", lResp.ErrMsg)
									lResp.ErrMsg = "Error: 4" + lErr4.Error()
								} else {
									lResp.Status = "S"
								}
							}
						}
					}
				}
				// Marshal the reponse structure into lDatas
				lData, lErr5 := json.Marshal(lResp)
				if lErr5 != nil {
					fmt.Fprintf(w, "Error taking data"+lErr5.Error())
				} else {
					fmt.Fprintf(w, string(lData))
				}
			}
		}
	}
	log.Println("AddLookUpHeader(-)")
}
func ExistsId(db *sql.DB, Code string) (string, error) {
	log.Println("ExistsId(+)")

	var ExistVal string
	CoreString := `select (case when count(1) > 0 then 'Y' else 'N' end) isExists
	from xx_lookup_header 
	WHERE Code = ?`

	rows, err := db.Query(CoreString, Code)
	if err != nil {
		log.Println(err)
		return ExistVal, err
	} else {
		for rows.Next() {
			err := rows.Scan(&ExistVal)
			if err != nil {
				log.Println("error1")
				return ExistVal, err
			}
		}
		log.Println("ExistsId(-)")
	}
	return ExistVal, err
}
func InsertHeaderValue(lDb *sql.DB, lInput LookUpHeader) error {
	log.Println("InsertHeaderValue(+)")

	// log.Println("insert input", lInput)

	// ! To insert the LookUpHeader code,description
	CoreString := `insert into xx_lookup_header(Code ,description ,createdBy, createdDate, updatedBy, updatedDate)
	values(?,?,?,Now(),?,Now())`
	_, lErr := lDb.Exec(CoreString, &lInput.Code, &lInput.Description, &lInput.User, &lInput.User)
	if lErr != nil {
		log.Println("error1")
		return lErr
	}
	log.Println("InsertHeaderValue(-)")
	return nil
}

/*
Pupose:This Function is used to Update new lookup header data into our database table ....
Request:

{
	"id": "1",
	"User": "xyz",
	"Code": "20",
	"Description": "hwsfuiqgu",
}

lResponse:

	*On Sucess
	=========
	{
	status: "S",
	ErrMsg: "",
	}

	!On Error
	========
	{
		"status": "failed",
		"reason": "Application no does not exist"
	}

Author: Nithish Kumar
Date: 08Aug2023
*/
func UpdateLookUpHeader(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	log.Println("UpdateLookUpHeader(+)")
	if r.Method == "PUT" {
		// This variable is used to store the LookUpHeader structure
		var lInput LookUpHeader
		// This variable is used to store the lResponse structure
		var lResp HeaderResponse
		lResp.Status = "S"
		//added by pavithra - to update the updatedBy value in Db
		lClientId, lErr1 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/lookup")
		if lErr1 != nil {
			log.Println("SPFC01", lErr1)
			lResp.Status = common.ErrorCode
			lResp.ErrMsg = "SPFC01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("SPFC01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lResp.Status = common.ErrorCode
				lResp.ErrMsg = "SPFC02 / UserDetails not Found"
				fmt.Fprintf(w, helpers.GetErrorString("SPFC02", "UserDetails not Found"))
				return
			}
		}

		body, lErr1 := ioutil.ReadAll(r.Body)
		log.Println("body", string(body))
		if lErr1 != nil {
			lResp.Status = "E"
			log.Println("Err 1", lErr1)
			lResp.ErrMsg = "Error: 1" + lErr1.Error()
		} else {
			lErr2 := json.Unmarshal(body, &lInput)
			// lId = string(body)
			log.Println("Input", lInput)
			if lErr2 != nil {
				lResp.Status = "E"
				log.Println("Err 2", lErr2)
				lResp.ErrMsg = "Error: 2" + lErr2.Error()
			} else {
				lClientId, lErr3 := verifyClient(r)
				if lErr3 != nil {
					lResp.Status = "E"
					log.Println("Err 2", lErr3)
					lResp.ErrMsg = "Error: 2" + lErr3.Error()
				} else {
					lInput.User = lClientId

					// To Establish A database connection,call LocalDbConnect Method
					lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
					if lErr2 != nil {
						lResp.Status = "E"
						log.Println("Err 2", lErr2)
						lResp.ErrMsg = "Error: 2" + lErr2.Error()
					} else {
						defer lDb.Close()
						// // ! To update the LookUpHeader Id ,Code,Description
						// CoreString := `update xx_lookup_header
						// set Code = ?,description = ?,updatedBy = 'AutoBot',updatedDate= Now()
						// where id = ?`
						// _, lErr3 := lDb.Exec(CoreString, &lInput.Code, &lInput.Description, &lInput.ID)
						// if lErr3 != nil {
						// 	lResp.Status = "E"
						// 	log.Println("Err 3", lResp.ErrMsg)
						// 	lResp.ErrMsg = "Error: 3" + lErr3.Error()
						// } else {
						// 	lResp.Status = "S"
						// }
						lExistsCode, lErr3 := ExistsCode(lDb, lInput)
						if lErr3 != nil {
							lResp.Status = "E"
							log.Println("Err 3", lResp.ErrMsg)
							lResp.ErrMsg = "Error: 3" + lErr3.Error()
						} else {
							if lExistsCode == "Y" {
								lResp.Status = "E"
								log.Println("Err 4: This Code is already exists", lResp.ErrMsg)
								lResp.ErrMsg = "This Code is already exists"
							} else if lExistsCode == "N" {
								lErr4 := UpdateHeaderValue(lDb, lInput, lClientId)
								if lErr4 != nil {
									lResp.Status = "E"
									log.Println("Err 4", lResp.ErrMsg)
									lResp.ErrMsg = "Error: 4" + lErr4.Error()
								} else {
									lResp.Status = "S"
								}
							}
						}
					}
				}
				// Marshal the reponse structure into lDatas
				lData, lErr5 := json.Marshal(lResp)
				if lErr5 != nil {
					fmt.Fprintf(w, "Error taking data"+lErr5.Error())
				} else {
					fmt.Fprintf(w, string(lData))
				}
			}
		}
	}
	log.Println("UpdateLookUpHeader(-)")
}
func ExistsCode(db *sql.DB, lInput LookUpHeader) (string, error) {
	log.Println("ExistsCode(+)")

	var ExistCode string
	CoreString := `select (case when count(1) > 0 then 'Y' else 'N' end) isExists
	from xx_lookup_header 
	WHERE id != ? and Code = ?`

	rows, err := db.Query(CoreString, lInput.ID, lInput.Code)
	if err != nil {
		log.Println(err)
		return ExistCode, err
	} else {
		for rows.Next() {
			err := rows.Scan(&ExistCode)
			if err != nil {
				log.Println("error1")
				return ExistCode, err
			}
		}
		log.Println("ExistsCode(-)")
	}
	return ExistCode, err
}

//added client ID in updatedby
func UpdateHeaderValue(lDb *sql.DB, lInput LookUpHeader, lClientId string) error {
	log.Println("UpdateHeaderValue(+)")
	log.Println("lInput", lInput)

	// ! To update the LookUpHeader Id ,Code,Description
	CoreString := `update xx_lookup_header
	set Code = ?,description = ?,updatedBy = ?,updatedDate= Now()
	where id = ?`
	_, lErr := lDb.Exec(CoreString, &lInput.Code, &lInput.Description, &lClientId, &lInput.ID)
	if lErr != nil {
		log.Println("error1")
		return lErr
	}
	log.Println("UpdateHeaderValue(-)")
	return nil
}
