package versioncontrol

import (
	"encoding/json"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type novoVersionResp struct {
	Version     string `json:"version"`
	ForceUpdate string `json:"forceUpdate"`
	Status      string `json:"status"`
	ErrMsg      string `json:"errMsg"`
}

// ==============================================================Novo  Update For Client side ===========================
func GetCurVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCurVersion (+)", r.Method)

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
		var lNovoUpdate novoVersionResp
		lNovoUpdate.Status = "S"
		// user Agent will Get Request from Header  as string
		userAgent := r.UserAgent()
		log.Println("userAgent", userAgent)

		//  Get Platform Method return url as per Device
		operatingSystem := getOSFromUserAgent(userAgent)
		if operatingSystem != "Unknown" {

			Version, lFlag, lErr1 := getAppupdates(operatingSystem)
			if lErr1 != nil {
				log.Println("VCGNV01 lErr1", lErr1.Error())
				lNovoUpdate.Status = common.ErrorCode
				lNovoUpdate.ErrMsg = "VCGNV01" + lErr1.Error()
				fmt.Fprintf(w, helpers.GetErrorString("VCGNV01", "Unable to Get Version Updates"))
				return
			} else {
				lNovoUpdate.Version = Version
				lNovoUpdate.ForceUpdate = lFlag
			}
		} else {
			lNovoUpdate.Status = common.ErrorCode
			lNovoUpdate.ErrMsg = "Unable To Find an Device Name"

		}
		//  else {
		// 	lNovoUpdate.Status = common.ErrorCode
		// 	lNovoUpdate.ErrMsg = "Unable To Find an Device Name"

		// }
		lData, lErr2 := json.Marshal(lNovoUpdate)
		if lErr2 != nil {
			log.Println("VCGNV02", lErr2)
			lNovoUpdate.Status = common.ErrorCode
			lNovoUpdate.ErrMsg = "VCGNV02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("VCGNV02", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetCurVersion (-)", r.Method)

	}
}

func getOSFromUserAgent(userAgent string) string {
	log.Println("getOSFromUserAgent (+)")
	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "android") {
		return "Android"
	} else if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") || strings.Contains(userAgent, "ios") {
		return "iOS"
	}
	log.Println("getOSFromUserAgent (-)")
	return "Unknown"
}

func getAppupdates(pDeviceName string) (string, string, error) {
	log.Println("getAppupdates (+)")
	var lVersion string

	var lNovoUpdate string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("VCGAU01", lErr1)
		return lVersion, lNovoUpdate, lErr1
	} else {
		defer lDb.Close()
		lCoreString2 := `SELECT  nvl(version,''),nvl(forceUpdate,'N')
		FROM a_novo_version_controller
		WHERE os = ?  and Status = 'Y'
		order by id desc limit 1`

		lRows1, lErr2 := lDb.Query(lCoreString2, pDeviceName)
		if lErr2 != nil {
			log.Println("VCGAU02", lErr2)
			return lVersion, lNovoUpdate, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lVersion, &lNovoUpdate)
				if lErr3 != nil {
					log.Println("VCGAU03", lErr3)
					return lVersion, lNovoUpdate, lErr3
				}
			}

		}
	}

	log.Println("getAppupdates (-)")
	return lVersion, lNovoUpdate, nil
}
