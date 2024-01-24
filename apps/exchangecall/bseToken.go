package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/integration/bse/bsesgb"
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
func BseGetToken(pUser string, pBrokerId int) (string, error) {
	log.Println("BseGetToken (+)")
	// Create a instance to store the token
	var lTokenValue string
	// Call method to get lToken from Database
	lToken, lErr1 := ValidToken(pBrokerId, common.BSE)
	if lErr1 != nil {
		log.Println("EBGT01", lErr1)
		return lTokenValue, lErr1
	} else {
		// check is the token is null
		if lToken == "" {
			// call Token method to get Token From Exchange
			lTokenResp, lErr2 := bsesgb.BseSgbToken(pUser, pBrokerId)
			if lErr2 != nil {
				log.Println("EBGT02", lErr2)
				return lTokenValue, lErr2
			} else {
				log.Println("lTokenResp", lTokenResp)
				//This condition executes when the token is successfully received from the exchange but
				// Also Get the error message as "Change Password" or error Code "012"
				if lTokenResp.Token != "" && lTokenResp.Message == "Change Password" || lTokenResp.ErrorCode == "012" {
					//------change pswd
					lErrorCode, lErr3 := ChangeBsePassword(pUser, lTokenResp.Token, pBrokerId)
					if lErr3 != nil {
						log.Println("EBGT03", lErr3)
						return lTokenValue, lErr3
					} else {
						if lErrorCode == "0" {
							lTokenResp, lErr4 := bsesgb.BseSgbToken(pUser, pBrokerId)
							if lErr4 != nil {
								log.Println("EBGT04", lErr4)
								return lTokenValue, lErr4
							} else {
								log.Println("lTokenResp", lTokenResp)
								// check if the Exchnage response status is success

								if lTokenResp.Token != "" {
									// call InsertToken method to insert the currently active token in database.
									lErr5 := InsertToken(lTokenResp.Token, pUser, pBrokerId, common.BSE)
									if lErr5 != nil {
										log.Println("EBGT05", lErr5)
										return lTokenValue, lErr5
									} else {
										// get token from database
										lToken, lErr6 := ValidToken(pBrokerId, common.BSE)
										if lErr6 != nil {
											log.Println("EBGT06", lErr6)
											return lTokenValue, lErr6
										} else {
											lTokenValue = lToken
										}
									}
								}
							}
						}
					}
					// check if the Exchnage response status is success
				} else if lTokenResp.Token != "" {
					// call InsertToken method to insert the currently active token in database.
					lErr7 := InsertToken(lTokenResp.Token, pUser, pBrokerId, common.BSE)
					if lErr7 != nil {
						log.Println("EBGT07", lErr7)
						return lTokenValue, lErr7
					} else {
						// get token from database
						lToken, lErr8 := ValidToken(pBrokerId, common.BSE)
						if lErr8 != nil {
							log.Println("EBGT08", lErr8)
							return lTokenValue, lErr8
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
	log.Println("BseGetToken (-)")
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
// func BseValidToken() (string, error) {
// 	log.Println("BseValidToken (+)")
// 	// create a new token instance
// 	var lToken string

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("ETVT01", lErr1)
// 		return lToken, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select ait.Token
// 						from a_ipo_token ait
// 						where now() between ait.ValidTime and ait.ExpireTime
// 						and ait.Exchange = 'BSE'
// 						order by ait.id desc
// 						limit 1 `
// 		lRows, lErr2 := lDb.Query(lCoreString)
// 		if lErr2 != nil {
// 			log.Println("ETVT02", lErr2)
// 			return lToken, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lToken)
// 				if lErr3 != nil {
// 					log.Println("ETVT03", lErr3)
// 					return lToken, lErr3
// 				} else {
// 					log.Println("TokenCaptured", lToken)
// 				}
// 			}
// 		}
// 	}
// 	log.Println("BseValidToken -")
// 	return lToken, nil
// }

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
// func BseInsertToken(pToken string, pUser string, pBrokerId int, pExchange string) error {
// 	log.Println("BseInsertToken (+)")

// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("ETIT01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		lSqlString := `insert into a_ipo_token (Token ,BrokerId,ValidTime,ExpireTime,Exchange,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
// 		values (?,?,now(),time(date_add(now(),interval 1 hour)),?,?,now(),?,now())`

// 		_, lErr2 := lDb.Exec(lSqlString, pToken, pBrokerId, common.BSE, pUser, pUser)
// 		if lErr2 != nil {
// 			log.Println("ETIT02", lErr2)
// 			return lErr2
// 		}
// 	}
// 	log.Println("BseInsertToken (-)")
// 	return nil
// }
