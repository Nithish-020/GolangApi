package brokers

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

type brokerStruct struct {
	Id          int    `json:"id"`
	BrokerName  string `json:"brokerName"`
	Domain      string `json:"domainName"`
	AdminCount  int    `json:"adminCount"`
	Status      string `json:"status"`
	Editable    bool   `json:"editable"`
	Type        string `json:"type"`
	RawDomain   string `json:"rawDomain"`
	AppName     string `json:"appName"`
	AuthURL     string `json:"authURL"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
}

type ResultRespStruct struct {
	BrokerListArr []brokerStruct `json:"brokerListArr"`
	Status        string         `json:"status"`
	ErrMsg        string         `json:"errMsg"`
}

func GetBrokerList(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBrokerList (+)", r.Method)
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// This variable is used to store the records
		var lBrokerRec brokerStruct
		// This variable is used to send resp to frontend
		var lRespRec ResultRespStruct
		lRespRec.Status = common.SuccessCode
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/domainsetup")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BGBL01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BGBL01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "BGBL02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("BGBL02", "Access restricted"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lDb, lErr3 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr3 != nil {
			log.Println("BGBL03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("BGBL03", "Unable to establish db connection"))
			return
		} else {
			defer lDb.Close()

			// ! To get the Symbol, ApplicationNo, Price and Quantity
			lCoreString2 := `SELECT
			bm.Id, bm.BrokerName,bm.DomainName, bm.Status,bm.Type,
			bm.RawDomain,bm.AppName,bm.AuthURL,bm.CreatedBy,bm.CreatedDate,
			COUNT(CASE WHEN u.RoleId = 4 THEN 1 ELSE NULL END) AS AdminCount
			from a_ipo_brokerMaster bm
			LEFT join a_ipo_userauth u ON bm.Id = u.brokerMasterId 
			GROUP BY  bm.Id,bm.BrokerName,bm.DomainName, bm.Status,
			bm.Type,bm.RawDomain, bm.AppName, bm.AuthURL,bm.CreatedBy,
			bm.CreatedDate;`

			// and m.BiddingEndDate < curdate() Its goes befor client id in query
			lRows1, lErr4 := lDb.Query(lCoreString2)
			if lErr4 != nil {
				log.Println("BGBL04", lErr4)
				fmt.Fprintf(w, helpers.GetErrorString("BGBL04", "UserDetails not Found"))
				return
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows1.Next() {
					lErr5 := lRows1.Scan(&lBrokerRec.Id, &lBrokerRec.BrokerName, &lBrokerRec.Domain, &lBrokerRec.Status, &lBrokerRec.Type, &lBrokerRec.RawDomain, &lBrokerRec.AppName, &lBrokerRec.AuthURL, &lBrokerRec.CreatedBy, &lBrokerRec.CreatedDate, &lBrokerRec.AdminCount)
					lBrokerRec.Editable = false
					if lErr5 != nil {
						log.Println("BGBL05", lErr5)
						fmt.Fprintf(w, helpers.GetErrorString("BGBL05", "UserDetails not Found"))
						return
					} else {
						// Append the history Records into lhistoryArr Array
						lRespRec.BrokerListArr = append(lRespRec.BrokerListArr, lBrokerRec)

					}
				}
			}
		}
		lData, lErr6 := json.Marshal(lRespRec)
		if lErr6 != nil {
			log.Println("BGBL06", lErr6)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "BGBL06" + lErr6.Error()
			fmt.Fprintf(w, helpers.GetErrorString("BGBL06", "Issue in Getting Your Datas.Try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetBrokerList (-)", r.Method)
	}
}

// =================================================//
// This method will check the broker present or not //
// =================================================//

func GetBrokerId(pDomainName string) (int, error) {
	log.Println("getBrokerId (+)")
	var lBrokerId int
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AGBI01", lErr1)
		return lBrokerId, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `select nvl(bm.Id,0) id
			from a_ipo_brokerMaster bm
			where bm.DomainName = ? 
			and bm.Status = 'Y'`

		lRows1, lErr2 := lDb.Query(lCoreString2, pDomainName)
		if lErr2 != nil {
			log.Println("AGBI02", lErr2)
			return lBrokerId, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lBrokerId)
				if lErr3 != nil {
					log.Println("BGBL05", lErr3)
					return lBrokerId, lErr3
				} else {
					// common.ABHIBrokerId = lBrokerId
					log.Println("Broker Id := ", lBrokerId)
				}
			}

		}
	}
	log.Println("getBrokerId (-)")
	return lBrokerId, nil
}

func GetAppName(pBrokerId int) (string, error) {
	log.Println("GetAppName (+)")
	var lAppName string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BGAN01", lErr1)
		return lAppName, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `select nvl(bm.AppName,'') name
			from a_ipo_brokerMaster bm
			where bm.id = ? 
			and bm.Status = 'Y'`

		lRows1, lErr2 := lDb.Query(lCoreString2, pBrokerId)
		if lErr2 != nil {
			log.Println("BGAN02", lErr2)
			return lAppName, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lAppName)
				if lErr3 != nil {
					log.Println("BGAN03", lErr3)
					return lAppName, lErr3
				} else {
					log.Println("lAppName", lAppName)
				}
			}

		}
	}
	log.Println("GetAppName (-)")
	return lAppName, nil
}

func GetDomain(pBrokerId int) (string, error) {
	log.Println("GetDomainName (+)")
	var lDomain string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BGAN01", lErr1)
		return lDomain, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `select nvl(bm.RawDomain,'') domain
			from a_ipo_brokerMaster bm
			where bm.id = ? 
			and bm.Status = 'Y'`

		lRows1, lErr2 := lDb.Query(lCoreString2, pBrokerId)
		if lErr2 != nil {
			log.Println("BGAN02", lErr2)
			return lDomain, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lDomain)
				if lErr3 != nil {
					log.Println("BGAN03", lErr3)
					return lDomain, lErr3
				} else {
					log.Println("lDomain", lDomain)
				}
			}

		}
	}
	log.Println("GetDomainName (-)")
	return lDomain, nil
}
