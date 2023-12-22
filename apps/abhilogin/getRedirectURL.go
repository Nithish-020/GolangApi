package abhilogin

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RedirectURLStruct struct {
	Host    string `json:"host"`
	AppName string `json:"appName"`
	Url     string `json:"url"`
}
type RedirectResp struct {
	UrlResp RedirectURLStruct `json:"redirectUrl"`
	Status  string            `json:"status"`
	ErrMsg  string            `json:"errMsg"`
}

func GetRedirectURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GetRedirectURL(+) " + r.Method)
	var lRedirectRec RedirectURLStruct
	var lRespRec RedirectResp

	// It helps to get origin & host from the request

	origin := r.Header.Get("Origin")
	lReferrer := r.Header.Get("Referer")
	log.Println("lReferrer", lReferrer)
	var lHost string

	if lReferrer != "" {
		hostParts1 := strings.Split(lReferrer, "://")
		log.Println("hostParts", hostParts1)
		if len(hostParts1) > 1 {
			lHost = hostParts1[1]

			// TO REMOVE THE SPECIAL CHAR IN THE REFERER WHEN THE CALL FROM HTTP REQUEST
			// BASED ON THE CONDITIONS FOR BOTH PROD AND DEVP
			if strings.Contains(lHost, ":") {
				hostParts2 := strings.Split(lHost, ":")
				lHost = hostParts2[0]
			} else {
				hostParts2 := strings.Split(lHost, "/")
				lHost = hostParts2[0]
			}

			// COMMENTED BY NITHISH.
			// THIS BELOW CODE WORKS ONLY FOR LOCAL DEVELOPMENT
			// hostParts2 := strings.Split(lHost, ":")
			// lHost = hostParts2[0]
			log.Println("lHost", lHost)
		}
	}

	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr1 := brokers.GetBrokerId(origin)
			if lErr1 != nil {
				lRespRec.Status = common.LoginFailure
				log.Println("ALGRU01", lErr1)
				fmt.Fprintf(w, helpers.GetErrorString("ALGRU01", "Unable to fetch your Records"))
				return
			}
			log.Println("lBrokerId", lBrokerId)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	switch r.Method {
	case "GET":
		lRespRec.Status = common.SuccessCode
		// To Establish A database connection,call LocalDbConnect Method
		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr2 != nil {
			log.Println("ALGRU02", lErr2)
			fmt.Fprintf(w, helpers.GetErrorString("ALGRU02", "Unable to reach the server right now"))
			return

		} else {
			defer lDb.Close()

			lCoreString := `select bm.RawDomain,bm.AppName ,bm.AuthURL
			from a_ipo_brokermaster bm
			where bm.RawDomain = ?
			`
			lRows, lErr3 := lDb.Query(lCoreString, lHost)
			if lErr3 != nil {
				lRespRec.Status = common.ErrorCode
				log.Println("ALGRU03", lErr3)
				fmt.Fprintf(w, helpers.GetErrorString("ALGRU03", "request denied"))
				return

			} else {
				for lRows.Next() {
					lErr4 := lRows.Scan(&lRedirectRec.Host, &lRedirectRec.AppName, &lRedirectRec.Url)
					if lErr4 != nil {
						lRespRec.Status = common.ErrorCode
						log.Println("ALGRU04", lErr4)
						fmt.Fprintf(w, helpers.GetErrorString("ALGRU04", "You have no acces to this website"))
						return

					} else {
						if lRedirectRec.Host != "" {
							lRespRec.UrlResp = lRedirectRec
							// common.ABHIAppName = lRedirectRec.AppName
							// common.ABHIDomain = lRedirectRec.Host
							lRespRec.Status = common.SuccessCode

							log.Println("lRedirect succesful", lRedirectRec)
						} else {
							lRespRec.Status = common.LoginFailure
						}
					}
				}
			}

		}
		data, lErr5 := json.Marshal(lRespRec)
		if lErr5 != nil {
			// fmt.Fprintf(w, "Error taking data"+lErr5.Error())
			fmt.Fprintf(w, helpers.GetErrorString("ALGRU05", "Error while reaching server"))
			return
		} else {
			fmt.Fprintf(w, string(data))
		}
		log.Println("GetRedirectURL(-) ")
	}

}

// func wrapper(pHostName string) (string, error) {
// 	log.Println("Wrapper (+)")
// 	var lFlag string
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("GRUW01", lErr1)
// 		return lFlag, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select nvl((select  case when RawDomain = ? then 'Y' else 'I' end as Domain
// 		from a_ipo_brokermaster aib
// 		where RawDomain = ? ),'I') as NoDomain`
// 		lRows, lErr2 := lDb.Query(lCoreString, pHostName, pHostName)
// 		if lErr2 != nil {
// 			log.Println("GRUW02", lErr2)
// 			return lFlag, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lFlag)
// 				if lErr3 != nil {
// 					log.Println("GRUW03", lErr3)
// 					return lFlag, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("Wrapper (-)")
// 	return lFlag, nil
// }

//added by naveen : to fetch source(mobile or web) of user from novo token based on cookie name
func GetSourceOfUser(pReq *http.Request, pPublicTokenCookieName string) (string, error) {
	log.Println("GetSourceOfUser(+)")
	source := ""

	if pPublicTokenCookieName != "" {
		publicTokenCookie, lErr1 := pReq.Cookie(pPublicTokenCookieName)
		if lErr1 != nil {
			log.Print("ALGS01", lErr1)
		} else {
			lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
			if lErr1 != nil {
				log.Print("ALGS02", lErr1)
				return source, lErr1
			} else {
				defer lDb.Close()

				lCoreString := `select source 
		             from novo_token
		           where Token=?
		        `
				lRows, lErr3 := lDb.Query(lCoreString, publicTokenCookie.Value)
				if lErr3 != nil {
					log.Print("ALGS03", lErr1)
					return source, lErr3

				} else {
					for lRows.Next() {
						lErr4 := lRows.Scan(&source)
						if lErr4 != nil {
							log.Print("ALGS04", lErr1)
							return source, lErr4
						}
					}
				}
			}
		}
	}
	log.Println("GetSourceOfUser(-)")
	return source, nil

}
