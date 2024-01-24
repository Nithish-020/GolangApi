package apigate

import (
	"crypto/rand"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
	"time"
)

// --------------------------------------------------------------------
// function to get random number
// --------------------------------------------------------------------
func tokenGenerator(clientid string, reqDtl RequestorDetails) string {
	log.Println("tokenGenerator+")
	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()
		validclientid := "N"
		sqlString := "select 'Y' from xxapi_clients where enabled='Y' and clientid = '" + clientid + "'"
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println("api client select error", err.Error())
		} else {
			//-----------Before Looping records----------
			for rows.Next() {
				err := rows.Scan(&validclientid)
				if err != nil {
					log.Println("api client record loop", err.Error())
				}
			}
			if validclientid == "Y" {
				b := make([]byte, 32)
				rand.Read(b)
				token := fmt.Sprintf("%x."+"%x", time.Now().UnixNano(), b)
				log.Println(token)
				//insert token
				insertString := "insert into xxapi_tokens(clientid,token,createdtime,expiretime,realip,forwardedip,method,path,host,remoteaddr) values (?,?,now() ,ADDTIME(now(), '00:02:00.999998'),?,?,?,?,?,?)"
				_, err = db.Exec(insertString, clientid, token, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host, reqDtl.RemoteAddr)
				if err != nil {
					log.Println("token insert error", err.Error())
				} else {
					log.Println("tokenGenerator-")
					return token
				}
			}
		}
	}
	log.Println("tokenGenerator-")
	return ""
}

// ----------------------------------------------------------------------------------------
// function used to generate token that can used to access api from third party application
// ----------------------------------------------------------------------------------------

// func CreateAppAPIAccessToken(clientid string, appkey string, expirtyTime string, reqDtl RequestorDetails) (string, string) {
// 	log.Println("CreateAppAPIAccessToken+")
// 	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if err != nil {
// 		log.Println(err)
// 		log.Println("CreateAppAPIAccessToken-")
// 		return "", "Error:10201: Please try after some time"
// 	} else {
// 		defer db.Close()
// 		b := make([]byte, 32)
// 		rand.Read(b)
// 		token := fmt.Sprintf("%x."+"%x", time.Now().UnixNano(), b)
// 		//insert token
// 		insertString := "insert into xxapi_tokens(clientid,appkey,token,createdtime,expiretime,realip,forwardedip,method,path,host,remoteaddr) values (?,?,?,now() ,ADDTIME(now(), '" + expirtyTime + "'),?,?,?,?,?,?)"
// 		_, err = db.Exec(insertString, clientid, appkey, token, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host, reqDtl.RemoteAddr)
// 		if err != nil {
// 			log.Println("token insert error", err.Error())
// 			log.Println("CreateAppAPIAccessToken-")
// 			return "", "Error:10202: Please try after some time"
// 		} else {
// 			log.Println("CreateAppAPIAccessToken-")
// 			return token, ""
// 		}
// 	}
// 	log.Println("CreateAppAPIAccessToken-")
// 	return "", "Error:10204: Please try after some time"
// }

// --------------------------------------------------------------------
// handler processes messages sent by unsubscribe
// --------------------------------------------------------------------
// func GenerateToken(w http.ResponseWriter, r *http.Request) {
// 	log.Println("GenerateToken+")

// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	w.WriteHeader(200)
// 	body, _ := ioutil.ReadAll(r.Body)
// 	reqDtl := GetRequestorDetail(r)
// 	reqDtl.Body = string(body)
// 	LogRequest("", reqDtl)

// 	authstring := strings.Fields(r.Header.Get("Authorization"))
// 	//log.Println(reqDtl)

// 	if len(authstring) > 1 {
// 		if authstring[0] == "Flattrade-oauthtoken" {
// 			if authstring[1] != "" {
// 				switch r.Method {
// 				case "GET":
// 					token := tokenGenerator(authstring[1], reqDtl)
// 					fmt.Fprintf(w, token)
// 				default:
// 					fmt.Fprintf(w, "ERROR: only GET method is supported.")
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GenerateToken-")
// }
