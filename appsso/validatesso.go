package appsso

import (
	"encoding/base64"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
	"net/http"
)

/*
-----------------------------------------------------------------------------------
function used to validate session and publick token cookie value of the web application
-----------------------------------------------------------------------------------
*/
func ValidateAndGetClientDetails(req *http.Request, app string, publicTokenCookieName string) (string, error) {
	log.Println("ValidateAndGetClientDetails+")
	if publicTokenCookieName != "" {
		publicTokenCookie, err := req.Cookie(publicTokenCookieName)
		if err != nil {
			log.Println("error reading publicToken cookie")
			log.Println(err)
			log.Println("ValidateAndGetClientDetails-")
			return "", fmt.Errorf("Invalid Session (0x901)")
		} else {
			clientId := CheckPageTokenValidity2(publicTokenCookie.Value, app)
			if clientId != "" {
				//check if client is the staff
				//check if any client set in client coookie
				//if set then return the client set in cookie
				//else return the original client
				log.Println("ValidateAndGetClientDetails-")
				return clientId, nil
			} else {
				log.Println("ValidateAndGetClientDetails-")
				return "", fmt.Errorf("Invalid Session (0x902)")
			}

		}
	} else {
		log.Println("ValidateAndGetClientDetails-")
		return "", fmt.Errorf("Invalid Session (0x904)")
	}
}

/*
-----------------------------------------------------------------------------------
function used to validate session and publick token cookie value of the web application for MyAccount page
-----------------------------------------------------------------------------------
*/
func ValidateAndGetClientDetails2(req *http.Request, app string, publicTokenCookieName string) (string, error) {
	log.Println("ValidateAndGetClientDetails2+")

	var LoggedBy string // logged in ClientId
	//var SetClientId string //  Set ClientId

	var clientId string
	if publicTokenCookieName != "" {
		publicTokenCookie, err := req.Cookie(publicTokenCookieName)
		if err != nil {
			log.Println("error reading publicToken cookie")
			log.Println(err)
			log.Println("ValidateAndGetClientDetails2-")
			return clientId, fmt.Errorf("Invalid Session (0x901)")
		} else {
			LoggedBy = CheckPageTokenValidity2(publicTokenCookie.Value, app)

			if LoggedBy != "" {
				clientId = LoggedBy + "," + LoggedBy
				//check if client is the staff
				//check if any client set in client coookie
				//if set then return the client set in cookie
				//else return the original client
				log.Println(LoggedBy)
				isStaff, err := CheckIsStaff(LoggedBy)
				if err != nil {
					log.Println(err)
					log.Println("ValidateAndGetClientDetails2-")
					return clientId, fmt.Errorf("Invalid Session (0x901-1)")
				} else {
					log.Println("isStaff", isStaff)
					if isStaff == "Y" {
						var publicCookie *http.Cookie
						publicCookie, err = req.Cookie(common.ABHIClientCookieName) // when using comman it cause error
						if err != nil {
							log.Println(err)
							log.Println("ValidateAndGetClientDetails2-")
							//return "", fmt.Errorf("Invalid Session (0x901)")
						} else {
							if publicCookie.Value != "" {
								SetClientId, err := base64.StdEncoding.DecodeString(publicCookie.Value)
								if err != nil {
									log.Println(err)
									log.Println("ValidateAndGetClientDetails2-")
									return clientId, fmt.Errorf("Invalid Session (0x901)")
								} else {
									clientId = string(SetClientId) + "," + LoggedBy
									log.Println("ValidateAndGetClientDetails2-")
									return clientId, nil
								}

							} else {
								//clientId = SetClientId + "," + LoggedBy
								log.Println("ValidateAndGetClientDetails2-")
								return clientId, nil
							}
						}
					} else {
						//clientId = LoggedBy + "," + LoggedBy
						log.Println("ValidateAndGetClientDetails2-")
						return clientId, nil
					}
				}

			} else {
				log.Println("ValidateAndGetClientDetails2-")
				return clientId, fmt.Errorf("Invalid Session (0x902)")
			}

		}
	} else {
		log.Println("ValidateAndGetClientDetails2-")
		return clientId, fmt.Errorf("Invalid Session (0x904)")
	}

	log.Println("ValidateAndGetClientDetails2-")
	return clientId, nil

}

func CheckIsStaff(clientId string) (string, error) {
	log.Println("checkIsStaff+")

	var IsStaff string

	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	//if any error when opening db connection
	if err != nil {
		log.Println(err)
		return IsStaff, err
	} else {
		defer db.Close()
		CoreString := `select (case when clientId = '' then 'N' else 'Y' end) validate
					from clientStaff_Mapping
					where ClientId = ?`

		rows, err := db.Query(CoreString, clientId)
		if err != nil {
			log.Println(err)
			return IsStaff, err
		} else {
			for rows.Next() {
				err := rows.Scan(&IsStaff)
				if err != nil {
					log.Println(err)
					return IsStaff, err
				}
			}
		}
	}
	log.Println("checkIsStaff-")
	return IsStaff, nil

}

// --------------------------------------------------------------------
// check sso token validity
// --------------------------------------------------------------------
func CheckPageTokenValidity2(token string, app string) string {
	log.Println("CheckPageTokenValidity2+")
	//db, err := util.Getdb(config.Database.DbType, config.Database)
	//db, err := localdbconnect(MariaAuthPRD)
	db, err := ftdb.LocalDbConnect(ftdb.SSODB)
	clientId := ""
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()

		// sqlString := "select clientid from xxapi_ssotokens where NOW() between createdtime and expiretime and  nvl(validated,'N') = 'N' and token  = '" + token + "' and app='" + app + "' "

		sqlString := "select clientid from xxapi_ssotokens where NOW() between createdtime and expiretime and  nvl(validated,'N') = 'N' and token  = ? and app=? "

		rows, err := db.Query(sqlString, token, app)
		if err != nil {
			log.Println("CheckPageTokenValidity2 > token select error", err.Error())
		} else {
			//-----------Before Looping records----------
			for rows.Next() {
				err := rows.Scan(&clientId)
				if err != nil {
					log.Println("CheckPageTokenValidity2 > soo token log record loop", err.Error())
				}
			}
		}
	}
	log.Println("CheckPageTokenValidity2-")
	return clientId
}

// --------------------------------------------------------------------
// function to validate sso token
// --------------------------------------------------------------------
func IsAuthTokenValid(app string, token string, client string) string {
	log.Println("IsAuthTokenValid+")
	tokenid := 0
	db, err := ftdb.LocalDbConnect(ftdb.SSODB)
	if err != nil {
		log.Println(err)
		log.Println("IsAuthTokenValid-")
		return "N"
	} else {
		defer db.Close()
		log.Println("app", app)
		log.Println("token", token)
		log.Println("client", client)

		sqlString := `select id
						from xxapi_ssotokens xs 
						where NOW() between xs.createdtime and xs.expiretime 
						and  nvl(xs.validated,'N') = 'N' 
						and xs.token  = '` + token + `' 
						and xs.clientid ='` + client + `' 
						and xs.app ='` + app + `'
						`
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println("isAuthTokenValid >>> token select error", err.Error())
			log.Println("IsAuthTokenValid-")
			return "N"
		} else {
			//-----------Before Looping records----------
			for rows.Next() {
				err := rows.Scan(&tokenid)
				if err != nil {
					log.Println("isAuthTokenValid >>> soo token log record loop", err.Error())
				}
			}
			//log.Println(tokenID)
			if tokenid > 0 {
				updString := "update xxapi_ssotokens set validated='Y' , validatedtime = NOW() where id = ?"
				_, err := db.Exec(updString, tokenid)
				if err != nil {
					log.Println("isTokenValid >>> sso token validated update error", err.Error())
					log.Println("IsAuthTokenValid-")
					return "N"
				} else {
					log.Println("IsAuthTokenValid-")
					return "Y"
				}
			}
		}
	}
	log.Println("IsAuthTokenValid-")
	return "N"

}

// --------------------------------------------------------------------
// check sso token validity
// --------------------------------------------------------------------
// func CheckSSOTokenValidity(token string, clientid string, reqDtl apigate.RequestorDetails) int {

// 	//logRequest(token, reqDtl)
// 	db, err := ftdb.LocalDbConnect(ftdb.SSODB)
// 	//db, err := util.Getdb(config.Database.DbType, config.Database)
// 	tokenID := 0
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		defer db.Close()

// 		sqlString := "select id from xxapi_ssotokens where NOW() between createdtime and expiretime and  nvl(validated,'N') = 'N' and token  = '" + token + "' and clientid ='" + clientid + "'"
// 		log.Println(sqlString)
// 		rows, err := db.Query(sqlString)
// 		if err != nil {
// 			log.Println("token select error", err.Error())
// 		} else {
// 			//-----------Before Looping records----------
// 			for rows.Next() {
// 				err := rows.Scan(&tokenID)
// 				if err != nil {
// 					log.Println("soo token log record loop", err.Error())
// 				}
// 			}
// 			log.Println(tokenID)
// 			if tokenID > 0 {
// 				updString := "update xxapi_ssotokens set validated='Y' , validatedtime = NOW() where id = ?"
// 				_, err := db.Exec(updString, tokenID)
// 				if err != nil {
// 					log.Println("sso token validated update error", err.Error())
// 				}
// 			}
// 		}
// 	}

// 	return tokenID
// }
