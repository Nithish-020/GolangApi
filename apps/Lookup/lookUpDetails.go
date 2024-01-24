package lookup

import (
	"database/sql"
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// To store LookUp Details Records
type LookUpDetails struct {
	ID           string `json:"id"`
	User         string `json:"user"`
	Code         string `json:"code"`
	HeaderId     string `json:"headerId"`
	Description  string `json:"description"`
	Attribute    string `json:"attribute"`
	Created_By   string `json:"createdBy"`
	Created_Date string `json:"createdDate"`
	Updated_By   string `json:"updatedBy"`
	Updated_Date string `json:"updatedDate"`
}

// lResponse Structure for GetLookUpDetails API
type DetailsResponse struct {
	DetailsArr []LookUpDetails `json:"details"`
	ErrMsg     string          `json:"errMsg"`
	Status     string          `json:"status"`
}

/*
Pupose:This Function is used to fetch the LookUp details data from our database table ....
Request:

id:"1"

lResponse:

	*On Sucess
	=========
	{
	[{
		"id": "1",
		"Code" :"10",
		"HeaderId" :"50",
		"Description":"wetgrh",
		"attribute" :"2",
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
func GetLookUpDetails(w http.ResponseWriter, r *http.Request) {
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
	log.Println("GetLookUpDetails(+)", r.Method)
	if r.Method == "PUT" {
		// var lId string
		// This variable is used to store the LookUpDetails structure
		var lInput LookUpDetails
		// This variable is used to store the lResp structure
		var lResp DetailsResponse
		lResp.Status = "S"
		body, lErr1 := ioutil.ReadAll(r.Body)
		// log.Println("body", string(body))
		if lErr1 != nil {
			lResp.Status = "E"
			log.Println("Err 1", lErr1)
			lResp.ErrMsg = "Error: 1" + lErr1.Error()
		} else {
			// lErr2 := json.Unmarshal(body, &id)
			lInput.HeaderId = string(body)
			// log.Println("HeaderId", lInput.HeaderId)
			// if lErr2 != nil {
			// 	lResp.Status = "E"
			// 	log.Println("Err 2", lErr2)
			// 	lResp.ErrMsg = "Error: 2" + lErr2.Error()
			// } else {

			// To Establish A database connection,call LocalDbConnect Method
			lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
			if lErr2 != nil {
				lResp.Status = "E"
				log.Println("Err 2", lErr2)
				lResp.ErrMsg = "Error: 2" + lErr2.Error()
			} else {

				defer lDb.Close()
				// ! To get the headerId,code,description,attribute
				CoreString := `select id,Code ,description ,nvl(Attribute1, '') Attribute1,nvl(createdBy, '') createdBy ,nvl(createdDate, '') createdDate ,nvl(updatedBy, '') updatedBy ,nvl(updatedDate, '')  updatedDate
				from xx_lookup_details
				where headerid = ? `
				lRows, lErr3 := lDb.Query(CoreString, lInput.HeaderId)
				if lErr3 != nil {
					lResp.Status = "E"
					log.Println("Err 3", lResp.ErrMsg)
					lResp.ErrMsg = "Error: 3" + lErr3.Error()
				} else {
					//This for loop is used to collect the records from the database and store them in structure
					for lRows.Next() {
						lErr4 := lRows.Scan(&lInput.ID, &lInput.Code, &lInput.Description, &lInput.Attribute, &lInput.Created_By, &lInput.Created_Date, &lInput.Updated_By, &lInput.Updated_Date)
						if lErr4 != nil {
							lResp.Status = "E"
							log.Println("Err 4", lResp.ErrMsg)
							lResp.ErrMsg = "Error: 4" + lErr4.Error()
						} else {
							// Append the LookUpDetails Records into lResp DetailsArr Array
							// log.Println("lInput", lInput)
							lResp.DetailsArr = append(lResp.DetailsArr, lInput)
							// log.Println("DetailsArr", lResp.DetailsArr)
							lResp.Status = "S"
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
	log.Println("GetLookUpDetails(-)")
}

/*
Pupose:This Function is used to Add new LookUp details data into our database table ....
Request:

{
	"id": "1",
		"Code" :"10",
		"HeaderId" :"50",
		"Description":"wetgrh",
		"attribute" :"2",
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
func AddLookUpDetails(w http.ResponseWriter, r *http.Request) {
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
	log.Println("AddLookUpDetails(+)")
	if r.Method == "PUT" {
		// This variable is used to store the StampConfigMaster structure
		var lInput LookUpDetails
		// This variable is used to store the lResponse structure
		var lResp DetailsResponse
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
					log.Println("Err 2", lErr3)
					lResp.ErrMsg = "Error: 2" + lErr3.Error()
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
						// ! To insert the headerId,code,description,attribute

						// CoreString := `insert into xx_lookup_details(headerid,Code , description ,Attribute1,  createdBy, createdDate, updatedBy, updatedDate)
						// values(?,?,?,?,?,Now(),?,Now())`
						// _, lErr3 := lDb.Exec(CoreString, &lInput.HeaderId, &lInput.Code, &lInput.Description, &lInput.Attribute, &lInput.User, &lInput.User)
						// if lErr3 != nil {
						// 	lResp.Status = "E"
						// 	log.Println("Err 3", lResp.ErrMsg)
						// 	lResp.ErrMsg = "Error: 3" + lErr3.Error()
						// } else {
						// 	lResp.Status = "S"
						// }
						lExistsVal, lErr4 := ExistsHeaderId(lDb, lInput)
						if lErr4 != nil {
							lResp.Status = "E"
							log.Println("Err 3", lResp.ErrMsg)
							lResp.ErrMsg = "Error: 3" + lErr4.Error()
						} else {
							if lExistsVal == "Y" {
								lResp.Status = "E"
								log.Println("Err 4 : This Code is already exists", lResp.ErrMsg)
								lResp.ErrMsg = "This Code is already exists"
							} else if lExistsVal == "N" {
								lErr5 := InsertDetailsValue(lDb, lInput)
								if lErr5 != nil {
									lResp.Status = "E"
									log.Println("Err 4", lErr5)
									lResp.ErrMsg = "Error: 4" + lErr5.Error()
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
	log.Println("AddLookUpDetails(-)")
}
func InsertDetailsValue(lDb *sql.DB, lInput LookUpDetails) error {
	log.Println("InsertDetailsValue(+)")

	// log.Println("insert input", lInput)

	// ! To insert the LookUpHeader code,description
	CoreString := `insert into xx_lookup_details(headerid,Code , description ,Attribute1,  createdBy, createdDate, updatedBy, updatedDate)
	values(?,?,?,?,?,Now(),?,Now())`
	_, lErr := lDb.Exec(CoreString, &lInput.HeaderId, &lInput.Code, &lInput.Description, &lInput.Attribute, &lInput.User, &lInput.User)
	if lErr != nil {
		log.Println("error1")
		return lErr
	}
	log.Println("InsertDetailsValue(-)")
	return nil
}
func ExistsHeaderId(db *sql.DB, lInput LookUpDetails) (string, error) {
	log.Println("ExistsHeaderId(+)")

	var ExistVal string
	CoreString := `select (case when count(1) > 0 then 'Y' else 'N' end) isExists
	from xx_lookup_details 
	WHERE headerid = ? and Code=?`

	rows, err := db.Query(CoreString, lInput.HeaderId, lInput.Code)
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
		log.Println("ExistsHeaderId(-)")
	}
	return ExistVal, err
}

/*
Pupose:This Function is used to Update new LookUp details data into our database table ....
Request:

{
	"id": "1",
		"Code" :"10",
		"HeaderId" :"50",
		"Description":"wetgrh",
		"attribute" :"2",
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
func UpdateLookUpDetails(w http.ResponseWriter, r *http.Request) {
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
	log.Println("UpdateLookUpDetails(+)")
	if r.Method == "PUT" {
		// This variable is used to store the StampConfigMaster structure
		var lInput LookUpDetails
		// This variable is used to store the lResponse structure
		var lResp DetailsResponse
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
			log.Println("headerId value in update", lInput.HeaderId)
			if lErr2 != nil {
				lResp.Status = "E"
				log.Println("Err 2", lErr2)
				lResp.ErrMsg = "Error: 2" + lErr2.Error()
			} else {
				// lClientId, lErr3 := verifyClient(r)
				// if lErr3 != nil {
				// 	lResp.Status = "E"
				// 	log.Println("Err 2", lErr3)
				// 	lResp.ErrMsg = "Error: 2" + lErr3.Error()
				// } else {
				lInput.User = lClientId

				// To Establish A database connection,call LocalDbConnect Method
				lDb, lErr4 := ftdb.LocalDbConnect(ftdb.IPODB)
				if lErr4 != nil {
					lResp.Status = "E"
					log.Println("Err 4", lErr4)
					lResp.ErrMsg = "Error: 4" + lErr4.Error()
				} else {
					defer lDb.Close()
					// // ! To update the Code,Description,attribute
					// CoreString := `update xx_lookup_details
					// set Code = ?,description = ?,Attribute1 = ?,updatedBy = 'AutoBot',updatedDate= Now()
					// where id = ?`
					// _, lErr3 := lDb.Exec(CoreString, &lInput.Code, &lInput.Description, &lInput.Attribute, &lInput.ID)
					// if lErr3 != nil {
					// 	lResp.Status = "E"
					// 	log.Println("Err 3", lResp.ErrMsg)
					// 	lResp.ErrMsg = "Error: 3" + lErr3.Error()
					// } else {
					// 	lResp.Status = "S"
					// }
					lExistsCode, lErr5 := ExistsDetailsCode(lDb, lInput)
					if lErr5 != nil {
						lResp.Status = "E"
						log.Println("Err 5", lResp.ErrMsg)
						lResp.ErrMsg = "Error: 5" + lErr5.Error()
					} else {
						log.Println("lExistsCode", lExistsCode)
						if lExistsCode == "Y" {
							lResp.Status = "E"
							log.Println("Err 5: This Code is already exists", lResp.ErrMsg)
							lResp.ErrMsg = "This Code is already exists"
						} else if lExistsCode == "N" {
							lErr6 := UpdateDetailsValue(lDb, lInput)
							if lErr6 != nil {
								lResp.Status = "E"
								log.Println("Err 6", lResp.ErrMsg)
								lResp.ErrMsg = "Error: 4" + lErr6.Error()
							} else {
								lResp.Status = "S"
							}
						}
					}
				}
				// }
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
	log.Println("UpdateLookUpDetails(-)")
}
func ExistsDetailsCode(db *sql.DB, lInput LookUpDetails) (string, error) {
	log.Println("ExistsDetailsCode(+)")

	var ExistCode string
	CoreString := `select (case when count(1) > 0 then 'Y' else 'N' end) isExists
	from xx_lookup_details 
	WHERE id != ` + lInput.ID + ` and Code = '` + lInput.Code + `' and headerid = '` + lInput.HeaderId + `'`

	log.Println(CoreString)
	rows, err := db.Query(CoreString)
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
		log.Println("ExistsDetailsCode(-)")
	}
	return ExistCode, err
}
func UpdateDetailsValue(lDb *sql.DB, lInput LookUpDetails) error {
	log.Println("UpdateDetailsValue(+)")
	log.Println("lInput", lInput)

	// ! To update the LookUpHeader Id ,Code,Description
	CoreString := `update xx_lookup_details
	set Code = ?,description = ?,updatedBy = ?,updatedDate= Now()
	where id = ?`
	_, lErr := lDb.Exec(CoreString, &lInput.Code, &lInput.Description, &lInput.User, &lInput.ID)
	if lErr != nil {
		log.Println("error1")
		return lErr
	}
	log.Println("UpdateDetailsValue(-)")
	return nil
}
