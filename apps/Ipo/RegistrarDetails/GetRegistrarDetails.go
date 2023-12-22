package registrardetails

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

type RegistrarDetails struct {
	Id            int    `json:"id"`
	RegistrarName string `json:"registrarName"`
	RegistrarLink string `json:"registrarLink"`
	CreatedBy     string `json:"createdBy"`
	CreatedDate   string `json:"createdDate"`
	UpdatedBy     string `json:"updatedBy"`
	UpdatedDate   string `json:"updatedDate"`
}

type RegistrarResp struct {
	RegistrarList []RegistrarDetails `json:"registrarList"`
	Status        string             `json:"status"`
	ErrMsg        string             `json:"errMsg"`
}

func GetRegistrarDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("GetRegistrarDetails (+)", r.Method)
	origin := r.Header.Get("Origin")
	// var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			_, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		lClientId := ""
		var lRegister RegistrarResp
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipoinfo")
		if lErr1 != nil {
			lRegister.Status = common.ErrorCode
			lRegister.ErrMsg = "RGRD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRD01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lRegister.Status = common.ErrorCode
				lRegister.ErrMsg = "RGRD02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RGRD02", "Access restricted"))
				return
			}
		}
		lRegister.Status = common.SuccessCode
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lRegistrarList, lErr3 := getRegistrars()
		if lErr3 != nil {
			lRegister.Status = common.ErrorCode
			lRegister.ErrMsg = "RGRD03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRD03", "Registrars List not Found"))
			return
		} else {
			lRegister.RegistrarList = lRegistrarList
		}

		lData, lErr4 := json.Marshal(lRegister)
		if lErr4 != nil {
			log.Println("RGRD04", lErr4)
			lRegister.Status = common.ErrorCode
			lRegister.ErrMsg = "RGRD04" + lErr4.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RGRD04", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetRegistrarDetails (-)", r.Method)
	}
}

func getRegistrars() ([]RegistrarDetails, error) {
	log.Println("getRegistrars(+)")
	var lRegistrarsList []RegistrarDetails
	var lRegistrar RegistrarDetails
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RGR01", lErr1.Error())
		return lRegistrarsList, lErr1
	} else {
		defer lDb.Close()
		lSqlString := `SELECT nvl(Id,0),nvl(Trim(RegistrarName),''), nvl(Trim(RegistrarLink),''),nvl(createdBy,'') ,nvl(DATE_FORMAT(createdDate, '%d %b %Y %l:%i %p'),'')as createdDate ,nvl(updatedBy,'') ,nvl(DATE_FORMAT(updatedDate , '%d %b %Y %l:%i %p'),'')as updatedDate
		FROM a_ipo_registrars
		order by Id Desc;`
		lRows, lErr2 := lDb.Query(lSqlString)
		if lErr2 != nil {
			log.Println("RGR02", lErr2)
			return lRegistrarsList, lErr2

		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lRegistrar.Id, &lRegistrar.RegistrarName, &lRegistrar.RegistrarLink, &lRegistrar.CreatedBy, &lRegistrar.CreatedDate, &lRegistrar.UpdatedBy, &lRegistrar.UpdatedDate)
				if lErr3 != nil {
					log.Println("RGR03", lErr3)
					return lRegistrarsList, lErr3
				} else {
					lRegistrarsList = append(lRegistrarsList, lRegistrar)
				}
			}
		}
	}

	log.Println("getRegistrars(-)")
	return lRegistrarsList, nil

}

// // still symbol meed to filter Active is pending
// func getSymbolsList() ([]string, error) {
// 	log.Println("getSymbolsList(+)")
// 	var lSymbolList []string
// 	var symbols string
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("LGCPD01", lErr1.Error())
// 		return lSymbolList, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lSqlString := `SELECT nvl(symbol,'')
// 		FROM a_ipo_master
// 		order by Id Desc;`
// 		lRows, lErr2 := lDb.Query(lSqlString)
// 		if lErr2 != nil {
// 			log.Println("LGCPD02", lErr2)
// 			return lSymbolList, lErr2

// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&symbols)
// 				if lErr3 != nil {
// 					log.Println("LGCPD03", lErr3)
// 					return lSymbolList, lErr3
// 				} else {
// 					lSymbolList = append(lSymbolList, symbols)
// 				}
// 			}
// 		}
// 	}

// 	log.Println("getSymbolsList(-)")
// 	return lSymbolList, nil

// }
