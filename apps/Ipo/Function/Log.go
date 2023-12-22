package Function

import (
	"fcs23pkg/ftdb"
	"log"
)

// This Structure is Used to pass parameters to LogEntry method.
type ParameterStruct struct {
	EndPoint string `json:"endPoint"`
	Request  string `json:"request"`
	Response string `json:"response"`
	Method   string `json:"method"`
	ClientId string `json:"client"`
	Flag     string `json:"flag"`
	LastId   int    `json:"lastId"`
}

/*
Pupose: This method is used to store the data for endppoint datatable
Parameters:

	send ParameterStruct as a parameter to this method

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the dpStructArr data
		from the a_ipo_oder_header Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 07JUNE2023
*/
func LogEntry(pInput ParameterStruct) (int, error) {
	log.Println("InsertLog (+)")

	// create a instace to hold last inserted Id
	var lLogId int

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("FLE01", lErr1)
		return lLogId, lErr1
	} else {
		defer lDb.Close()

		// check is the flag is Insert
		if pInput.Flag == "INSERT" {
			lSqlString1 := `insert into a_ipo_endpoint_logtable (EndPoint,ReqJson,ResJson,Method,CreatedBy,CreatedDate)
				values (?,?,?,?,?,now())`

			lInsertedId, lErr2 := lDb.Exec(lSqlString1, pInput.EndPoint, pInput.Request, pInput.Response, pInput.Method, pInput.ClientId)
			if lErr2 != nil {
				log.Println("FLE02", lErr2)
				return lLogId, lErr2
			} else {
				lLog, _ := lInsertedId.LastInsertId()
				lLogId = int(lLog)
			}
			// Check if the flag is Update
		} else if pInput.Flag == "UPDATE" {
			lSqlString2 := `Update a_ipo_endpoint_logtable SET ResJson = ?,UpdatedBy = ?,UpdatedDate = now() 
			where id = ?`

			_, lErr3 := lDb.Exec(lSqlString2, pInput.Response, pInput.ClientId, pInput.LastId)
			if lErr3 != nil {
				log.Println("FLE03", lErr3)
				return lLogId, lErr3
			}
		}
	}
	log.Println("InsertLog (-)")
	return lLogId, nil
}
