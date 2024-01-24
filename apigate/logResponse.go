package apigate

import (
	"fcs23pkg/ftdb"
	"log"
	"net/http"
)

// --------------------------------------------------------------------
// log request details
// --------------------------------------------------------------------
func LogResponse(req *http.Request, respStatus int, respData []byte, requestID string) {
	reqDtl := GetRequestorDetail(req)
	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()
		//insert token
		insertString := "insert into xxapi_resp_log(response,responseStatus,requesteddate,realip,forwardedip,method,path,host,remoteaddr,header,body,endpoint,requestid) values (?,?,NOW(),?,?,?,?,?,?,?,?,?,?)"
		_, err := db.Exec(insertString, string(respData), respStatus, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host, reqDtl.RemoteAddr, reqDtl.Header, reqDtl.Body, reqDtl.EndPoint, requestID)
		if err != nil {
			log.Println("api log insert error", err.Error())
		}
	}

}
