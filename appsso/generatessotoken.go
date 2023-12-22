package appsso

//--------------------------------------------------------------------
//  function to insert sso token
//--------------------------------------------------------------------

// func GenerateSSOToken(db *sql.DB, clientid string, appKey string, reqDtl apigate.RequestorDetails) string {

// 	b := make([]byte, 32)
// 	rand.Read(b)
// 	token := fmt.Sprintf("%x."+"%x", time.Now().UnixNano(), b)

// 	//insert token
// 	insertString := "insert into xxapi_ssotokens(app,clientid,token,createdtime,expiretime,realip,forwardedip,method,path,host,remoteaddr) values (?,?,?,now() ,ADDTIME(now(), '00:02:00.999998'),?,?,?,?,?,?)"
// 	_, err := db.Exec(insertString, appKey, clientid, token, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host, reqDtl.RemoteAddr)
// 	if err != nil {
// 		log.Println("GenerateSSOToken insert error", err.Error())
// 	} else {
// 		return token
// 	}

// 	return ""
// }
