package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nseipo"
	"log"
)

/*
Purpose : This method is used to get token.
Parameters : nil
Response :

	==========
	On success
	==========

		It return 6454fb1e-4998-45ca-894a-0d03b1ae9ba3

	========
	On error
	========

		In case of any exception during the execution of this method you will get the error details.
		The calling program should handle the error

Author : NITHISH KUMAR
Date : 05-06-2023
*/
func GetToken(pUser string, pBrokerId int) (string, error) {
	log.Println("GetToken (+)")
	// Create a instance to store the token
	var lTokenValue string
	// Call method to get lToken from Database
	lToken, lErr1 := ValidToken(pBrokerId, common.NSE)
	if lErr1 != nil {
		log.Println("EGT01", lErr1)
		return lTokenValue, lErr1
	} else {
		// check is the token is null
		if lToken == "" {
			// call Token method to get Token From Exchange
			lTokenResp, lErr2 := nseipo.Token(pUser, pBrokerId)
			if lErr2 != nil {
				log.Println("EGT02", lErr2)
				return lTokenValue, lErr2
			} else {
				// check if the Exchnage response status is success
				if lTokenResp.Status == "success" {
					// call InsertToken method to insert the currently active token in database.
					lErr3 := InsertToken(lTokenResp.Token, pUser, pBrokerId, common.NSE)
					if lErr3 != nil {
						log.Println("EGT03", lErr3)
						return lTokenValue, lErr3
					} else {
						// get token from database
						lToken, lErr4 := ValidToken(pBrokerId, common.NSE)
						if lErr4 != nil {
							log.Println("EGT04", lErr4)
							return lTokenValue, lErr4
						} else {
							lTokenValue = lToken
						}
					}
				}
			}
		} else {
			lTokenValue = lToken
		}
	}
	log.Println("GetToken (-)")
	return lTokenValue, nil
}

/*
Purpose : This method is used to valid get token from Token database.
Parameters : nil
Response :

	==========
	On success
	==========

		It return 6454fb1e-4998-45ca-894a-0d03b1ae9ba3

	========
	On error
	========

		In case of any exception during the execution of this method you will get the error details.
		The calling program should handle the error.

Authorization : 'NITHISH KUMAR'
Date : '05-06-2023'
*/
func ValidToken(pBrokerId int, pExchange string) (string, error) {
	log.Println("ValidToken (+)")
	// create a new token instance
	var lToken string

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ETVT01", lErr1)
		return lToken, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select ait.Token
						from a_ipo_token ait 
						where now() between ait.ValidTime and ait.ExpireTime
						and ait.brokerId = ?
						and ait.Exchange = ?
						order by ait.id desc 
						limit 1 `
		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId, pExchange)
		if lErr2 != nil {
			log.Println("ETVT02", lErr2)
			return lToken, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lToken)
				if lErr3 != nil {
					log.Println("ETVT03", lErr3)
					return lToken, lErr3
				} else {
					log.Println("TokenCaptured", lToken)
				}
			}
		}
	}
	log.Println("ValidToken -")
	return lToken, nil
}

/*
Purpose : This method is used to cinsert the token Value in Database.
Parameter :

	pToken

Response :

	===========
	On Success:
	===========

		nil

	===========
	On Error:
	===========

		Error

Author : Pavithra
Date : 31-may-2023
*/
func InsertToken(pToken string, pUser string, pBrokerId int, pExchange string) error {
	log.Println("InsertToken (+)")

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ETIT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lSqlString := `insert into a_ipo_token (Token ,BrokerId,ValidTime,ExpireTime,Exchange,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
		values (?,?,now(),time(date_add(now(),interval 1 hour)),?,?,now(),?,now())`

		_, lErr2 := lDb.Exec(lSqlString, pToken, pBrokerId, pExchange, pUser, pUser)
		if lErr2 != nil {
			log.Println("ETIT02", lErr2)
			return lErr2
		}
	}
	log.Println("InsertToken (-)")
	return nil
}
