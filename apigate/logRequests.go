package apigate

import (
	"fcs23pkg/ftdb"
	"log"
)

// --------------------------------------------------------------------
// log request details
// --------------------------------------------------------------------
func LogRequest(token string, reqDtl RequestorDetails, requestID string) {
	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()
		//insert token
		insertString := "insert into xxapi_log(token,requesteddate,realip,forwardedip,method,path,host,remoteaddr,header,body,endpoint,requestid) values (?,NOW(),?,?,?,?,?,?,?,?,?,?)"
		_, err := db.Exec(insertString, token, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host, reqDtl.RemoteAddr, reqDtl.Header, reqDtl.Body, reqDtl.EndPoint, requestID)
		if err != nil {
			log.Println("api log insert error", err.Error())
		}
	}

}
