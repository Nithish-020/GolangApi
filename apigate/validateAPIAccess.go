package apigate

import (
	"fcs23pkg/ftdb"
	"log"
	"net/url"
)

// --------------------------------------------------------------------
// validate if the requestor is valid is access our api
// --------------------------------------------------------------------
func IsAuthorized(reqDtl RequestorDetails) bool {
	log.Println("IsAuthorized+")
	whitelistedIP := "N"
	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()
		realIP := ""
		if reqDtl.RealIP != "" {
			url, err := url.Parse(reqDtl.RealIP)
			if err != nil {
				//log.Fatal(err)
			} else {
				realIP = url.Hostname()
			}
			//fmt.Println(url.Hostname())
		}
		sqlString := `select 'Y' from xxapi_whitelistedIP
						where 1=1
						and endpoint  = '` + reqDtl.EndPoint + `' 
						and nvl(realip,'')  = (case when nvl(realip,'')='' then '' else '` + realIP + `' end )
						and nvl(forwardedip,'')  = (case when nvl(forwardedip,'')='' then '' else '` + reqDtl.ForwardedIP + `' end )
						and nvl(host,'')  = (case when nvl(host,'')='' then '' else '` + reqDtl.Host + `' end )
						and nvl(remoteaddr,'')  = (case when nvl(remoteaddr,'')='' then '' else substr('` + reqDtl.RemoteAddr + `',1,instr('` + reqDtl.RemoteAddr + `',':')-1) end ) 
						and nvl(method,'')  = (case when nvl(method,'')='' then '' else '` + reqDtl.Method + `' end ) 
						and current_date() between date_format(startDate,'%Y-%m-%d') and nvl(date_format(endDate ,'%Y-%m-%d') ,current_date()) `
		rows, err := db.Query(sqlString)
		if err != nil {
			log.Println("token select error", err.Error())
		} else {
			//-----------Before Looping records----------
			for rows.Next() {
				err := rows.Scan(&whitelistedIP)
				if err != nil {
					log.Println("token log record loop", err.Error())
				}
			}
		}
	}
	if whitelistedIP == "Y" {
		log.Println("Authorized :" + reqDtl.EndPoint + " > " + reqDtl.RealIP + " > " + reqDtl.ForwardedIP + " > " + reqDtl.Host + " > " + reqDtl.RemoteAddr + " > " + reqDtl.Method)
		log.Println("IsAuthorized-")
		return true
	} else {
		log.Println("NOT Authorized :" + reqDtl.EndPoint + " > " + reqDtl.RealIP + " > " + reqDtl.ForwardedIP + " > " + reqDtl.Host + " > " + reqDtl.RemoteAddr + " > " + reqDtl.Method)
		log.Println("IsAuthorized-")
		return false
	}

}
