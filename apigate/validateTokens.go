package apigate

// --------------------------------------------------------------------
// check token validity.
// In this method we will not check if requestor is a valid to call our api
// --------------------------------------------------------------------

// func CheckTokenValidity(token string, r *http.Request) string {
// 	log.Println("CheckTokenValidity+")
// 	reqDtl := GetRequestorDetail(r)
// 	LogRequest(token, reqDtl)
// 	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
// 	validToken := "N"
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		defer db.Close()

// 		sqlString := "select NVL(min('Y'),'N') from xxapi_tokens where NOW() between createdtime and expiretime and token  = '" + token + "'"
// 		rows, err := db.Query(sqlString)
// 		if err != nil {
// 			log.Println("token select error", err.Error())
// 		} else {
// 			//-----------Before Looping records----------
// 			for rows.Next() {
// 				err := rows.Scan(&validToken)
// 				if err != nil {
// 					log.Println("token log record loop", err.Error())
// 				}
// 			}
// 		}
// 	}
// 	if validToken == "N" {
// 		log.Println("CheckTokenValidity-")
// 		return "Invalid Token"
// 	}
// 	log.Println("CheckTokenValidity-")
// 	return validToken
// }

// --------------------------------------------------------------------
// check token validity
// In this method we will check if requestor is a valid to call our api
// --------------------------------------------------------------------

// func CheckTokenValidity2(token string, r *http.Request, body string) string {
// 	log.Println("CheckTokenValidity2+")
// 	reqDtl := GetRequestorDetail(r)
// 	reqDtl.Body = body
// 	LogRequest(token, reqDtl)
// 	validToken := "N"

// 	if IsAuthorized(reqDtl) {

// 		db, err := ftdb.LocalDbConnect(ftdb.IPODB)
// 		if err != nil {
// 			log.Println(err)
// 		} else {
// 			defer db.Close()

// 			sqlString := "select NVL(min('Y'),'N') from xxapi_tokens where NOW() between createdtime and expiretime and token  = '" + token + "'"
// 			rows, err := db.Query(sqlString)
// 			if err != nil {
// 				log.Println("token select error", err.Error())
// 			} else {
// 				//-----------Before Looping records----------
// 				for rows.Next() {
// 					err := rows.Scan(&validToken)
// 					if err != nil {
// 						log.Println("token log record loop", err.Error())
// 					}
// 				}
// 			}
// 		}
// 		if validToken == "N" {
// 			log.Println("CheckTokenValidity2-")
// 			return "Invalid Token"
// 		}
// 	}
// 	log.Println("CheckTokenValidity2-")
// 	return validToken
// }

// --------------------------------------------------------------------
// check token validity
// --------------------------------------------------------------------
// func CheckTokenValidity3(token string) string {
// 	log.Println("CheckTokenValidity3+")
// 	validToken := "N"
// 	db, err := ftdb.LocalDbConnect(ftdb.IPODB)

// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		defer db.Close()

// 		sqlString := "select NVL(min('Y'),'N') from xxapi_tokens where NOW() between createdtime and expiretime and token  = '" + token + "'"
// 		rows, err := db.Query(sqlString)
// 		if err != nil {
// 			log.Println("token select error", err.Error())
// 		} else {
// 			//-----------Before Looping records----------
// 			for rows.Next() {
// 				err := rows.Scan(&validToken)
// 				if err != nil {
// 					log.Println("token log record loop", err.Error())
// 				}
// 			}
// 		}
// 	}

// 	if validToken == "N" {
// 		log.Println("Token NOT Authorized :" + token)
// 		log.Println("CheckTokenValidity3-")
// 		return "Invalid Token"
// 	}
// 	log.Println("CheckTokenValidity3-")

// 	return validToken
// }

// --------------------------------------------------------------------
// Read token from the header details shared
// --------------------------------------------------------------------

// func Readtoken(authorization string) string {
// 	log.Println("Readtoken+")
// 	log.Println(authorization)
// 	authstring := strings.Fields(authorization)
// 	if len(authstring) > 1 {
// 		if authstring[0] == "Flattrade-oauthtoken" {
// 			if authstring[1] != "" {
// 				log.Println("Readtoken-")
// 				return authstring[1]
// 			}
// 		}
// 	}
// 	log.Println("Readtoken-")
// 	return ""
// }
