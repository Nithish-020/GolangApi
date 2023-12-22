package adminaccess

import (
	"fcs23pkg/apps/Ipo/placeorder"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
)

// Response Structure for FetchDirectory API
type DirectoryRespStruct struct {
	Stream string `json:"stream"`
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

/*
Pupose: This method is used to get the currently choosen orderPreference based on the broker
Parameters:
	not applicable
Response:
	    ==========
	    *On Sucess
	    ==========

			lDirectory: "NSE" || "BSE",

	    =========
	    !On Error
	    =========

			lDirectory: "",
			error: "Error in getting the orderPreference",

Author: KAVYA DHARSHANI M
Date: 06NOV2023
*/
func NcbFetchDirectory(pBrokerId int) (string, error) {
	log.Println("NcbFetchDirectory(+)")
	var lDirectory string

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NFD01", lErr1)
		return lDirectory, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select md.OrderPreference
		                from a_ipo_memberdetails md
		                where md.Flag = 'Y'
		                and md.BrokerId = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("NFD02", lErr2)
			return lDirectory, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lDirectory)
				if lErr3 != nil {
					log.Println("NFD03", lErr3)
					return lDirectory, lErr3
				}
			}
			log.Println("Current directory := ", string(lDirectory))
		}
	}

	log.Println("NcbFetchDirectory(-)")
	return lDirectory, nil
}

/*
Pupose: This method is used to get the currently choosen orderPreference based on the broker
Parameters:
	not applicable
Response:
	    ==========
	    *On Sucess
	    ==========

			lDirectory: "NSE" || "BSE",

	    =========
	    !On Error
	    =========

			lDirectory: "",
			error: "Error in getting the orderPreference",

Author: Nithish
Date: 05JUNE2023
*/
// func FetchDirectory(w http.ResponseWriter, r *http.Request) {
func SGBFetchDirectory(pBrokerId int) (string, error) {

	// log.Println("FetchDirectory (+)", r.Method)
	// origin := r.Header.Get("Origin")
	// for _, allowedOrigin := range common.ABHIAllowOrigin {
	// 	if allowedOrigin == origin {
	// 		w.Header().Set("Access-Control-Allow-Origin", origin)
	// 		log.Println(origin)
	// 		break
	// 	}
	// }

	// (w).Header().Set("Access-Control-Allow-Credentials", "true")
	// (w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	// (w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	// lBrokerId := 1
	var lDirectory string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AFD01", lErr1)
		return lDirectory, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select md.OrderPreference
		from a_ipo_memberdetails md
		where md.Flag = 'Y'
		and md.BrokerId = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("AFD02", lErr2)
			return lDirectory, lErr2
		} else {
			//This for loop is used to collect the record from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lDirectory)
				if lErr3 != nil {
					log.Println("AFD03", lErr3)
					return lDirectory, lErr3
				}
			}
			log.Println("Current directory := ", string(lDirectory))
		}
		// // Marshall the structure into json
		// lData, lErr5 := json.Marshal(lRespRec)
		// if lErr5 != nil {
		// 	log.Println("AFD05", lErr5)
		// 	fmt.Fprintf(w, helpers.GetErrorString("AFD05", "Unable to process your request now. Please try after sometime"))
		// 	return
		// } else {
		// 	fmt.Fprintf(w, string(lData))
		// }
		// log.Println("FetchDirectory (-)", r.Method)
	}
	log.Println("FetchDirectory (-)")
	return lDirectory, nil
}

/*
Pupose: This method is used to get the currently choosen orderPreference based on the broker
Parameters:
    not applicable
Response:
        ==========
        *On Sucess
        ==========
            lDirectory: "NSE" || "BSE",
        =========
        !On Error
        =========
            lDirectory: "",
            error: "Error in getting the orderPreference",
Author: Nithish
Date: 05JUNE2023
*/
// func FetchDirectory(w http.ResponseWriter, r *http.Request) {
func FetchDirectory(pReqRec placeorder.OrderReqStruct) (string, error) {
	var lDirectory string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AFD01", lErr1)
		return lDirectory, lErr1
	} else {
		defer lDb.Close()
		if pReqRec.ApplicationNo == "" {
			config := common.ReadTomlConfig("./toml/debug.toml")
			lCheckDirectory := fmt.Sprintf("%v", config.(map[string]interface{})["CurrentDirectory"])
			if lCheckDirectory == "" {
				lCoreString := `select (case when (select m.Isin  from a_ipo_master m where m.Symbol = ? and m.Exchange = "NSE") =
									(select m.Isin  from a_ipo_master m where m.Symbol = ? and m.Exchange = "BSE") then 'BSE' else am.Exchange  end )Exchange
									from a_ipo_master am
									where am.Symbol = ?
									and am.Id = ?`
				lRows, lErr2 := lDb.Query(lCoreString, pReqRec.Symbol, pReqRec.Symbol, pReqRec.Symbol, pReqRec.MasterId)
				if lErr2 != nil {
					log.Println("AFD02", lErr2)
					return lDirectory, lErr2
				} else {
					//This for loop is used to collect the record from the database and store them in structure
					for lRows.Next() {
						lErr3 := lRows.Scan(&lDirectory)
						if lErr3 != nil {
							log.Println("AFD03", lErr3)
							return lDirectory, lErr3
						}
					}
					log.Println("Current IPO directory := ", lDirectory)
				}
			} else {
				lDirectory = lCheckDirectory
			}
		} else {
			//Commnented by lakshmanan on 20OCT2023 to fix the bug to take order level exchange
			//below commented query always takes master level exchange
			// lCoreString := `select (case when count(1) > 0 then m.Exchange else (select m.Exchange from a_ipo_master m where m.Id = ?) end) ex
			// from a_ipo_master m,a_ipo_order_header h
			// where m.Symbol = ?
			// and m.Id = h.MasterId
			// and h.applicationNo = ?`
			// lRows, lErr2 := lDb.Query(lCoreString, pReqRec.MasterId, pReqRec.Symbol, pReqRec.ApplicationNo)
			//below query added by lakshmanan on 20OCT2023 to take order level exchange
			lCoreString := `select (case when nvl(h.exchange,'') <> '' then h.exchange else m.Exchange end) ex
								from a_ipo_master m,a_ipo_order_header h
								where  m.Id = h.MasterId
								and h.applicationNo = ? `
			lRows, lErr2 := lDb.Query(lCoreString, pReqRec.ApplicationNo)
			if lErr2 != nil {
				log.Println("AFD04", lErr2)
				return lDirectory, lErr2
			} else {
				//This for loop is used to collect the record from the database and store them in structure
				for lRows.Next() {
					lErr3 := lRows.Scan(&lDirectory)
					if lErr3 != nil {
						log.Println("AFD05", lErr3)
						return lDirectory, lErr3
					}
				}
				log.Println("Current IPO directory := ", lDirectory)
			}
		}
		// // Marshall the structure into json
		// lData, lErr5 := json.Marshal(lRespRec)
		// if lErr5 != nil {
		//  log.Println("AFD05", lErr5)
		//  fmt.Fprintf(w, helpers.GetErrorString("AFD05", "Unable to process your request now. Please try after sometime"))
		//  return
		// } else {
		//  fmt.Fprintf(w, string(lData))
		// }
		// log.Println("FetchDirectory (-)", r.Method)
	}
	log.Println("FetchDirectory (-)")
	return lDirectory, nil
}
