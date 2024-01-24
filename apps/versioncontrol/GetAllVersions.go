package versioncontrol

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

type novoAppStruct struct {
	Id          int    `json:"id"`
	OS          string `json:"os"`
	ForceUpdate string `json:"forceUpdate"`
	Version     string `json:"version"`
	AppStatus   string `json:"appStatus"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedDate string `json:"updatedDate"`
}

type getAllVersionDetails struct {
	AndroidVersionList []novoAppStruct `json:"androidVersionList"`
	IosVersionList     []novoAppStruct `json:"iosVersionList"`
	Status             string          `json:"status"`
	ErrMsg             string          `json:"errMsg"`
}

func GetAllVersions(w http.ResponseWriter, r *http.Request) {
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

		log.Println("GetAllVersions (+)")
		var lGetNovoVersion getAllVersionDetails
		lGetNovoVersion.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/versioncontrol")
		if lErr1 != nil {
			log.Println("VGAV01", lErr1)
			lGetNovoVersion.Status = common.ErrorCode
			lGetNovoVersion.ErrMsg = "VGAV01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("VGAV01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("VGAV02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lGetAndroidList, lErr2 := getAndroidVersion()
		if lErr2 != nil {
			log.Println("VGAV03", lErr2)
			lGetNovoVersion.Status = common.ErrorCode
			lGetNovoVersion.ErrMsg = "VGAV03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("VGAV03", "Unable To Get an Android Version Details"))
			return
		} else {
			lGetIosList, lErr3 := getIosVersion()
			if lErr3 != nil {
				log.Println("VGAV04", lErr2)
				lGetNovoVersion.Status = common.ErrorCode
				lGetNovoVersion.ErrMsg = "VGAV04" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("VGAV04", "Unable To Get an ios Version Details"))
				return
			} else {
				lGetNovoVersion.AndroidVersionList = lGetAndroidList
				lGetNovoVersion.IosVersionList = lGetIosList
			}

		}

		lData, lErr4 := json.Marshal(lGetNovoVersion)
		if lErr4 != nil {
			log.Println("VGAV05", lErr4)
			fmt.Fprintf(w, helpers.GetErrorString("VGAV05", "Error Occur in Marshalling Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("GetRoleTask (-)")
}

func getAndroidVersion() ([]novoAppStruct, error) {
	log.Println("getAndroidVersion(+) ")
	var lAndroid novoAppStruct
	var lAndroidList []novoAppStruct

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("VGAV01", lErr1)
		return lAndroidList, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `SELECT nvl(id,0),nvl(upper(os),''),nvl(forceUpdate,''),nvl(version,''),nvl(Status,''),nvl(createdBy,''),nvl( DATE_FORMAT(createdDate, '%d %b %Y %l:%i %p') ,'')AS createdDate,nvl(updatedBy,''), nvl( DATE_FORMAT(updatedDate , '%d %b %Y %l:%i %p') ,'')AS updatedDate 
		FROM a_novo_version_controller
		where os = 'ANDROID'
		order by id Desc`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("VGAV02", lErr2)
			return lAndroidList, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lAndroid.Id, &lAndroid.OS, &lAndroid.ForceUpdate, &lAndroid.Version, &lAndroid.AppStatus, &lAndroid.CreatedBy, &lAndroid.CreatedDate, &lAndroid.UpdatedBy, &lAndroid.UpdatedDate)
				if lErr3 != nil {
					log.Println("VGAV03", lErr3)
					return lAndroidList, lErr3
				} else {

					lAndroidList = append(lAndroidList, lAndroid)
				}
			}
		}
	}
	log.Println("getAndroidVersion(-) ")
	return lAndroidList, nil

}

func getIosVersion() ([]novoAppStruct, error) {
	log.Println("getIosVersion(+) ")
	var lIosVersion novoAppStruct
	var lIosVersionList []novoAppStruct

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("VGIV01", lErr1)
		return lIosVersionList, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `SELECT nvl(id,0),nvl(upper(os),''),nvl(forceUpdate,''),nvl(version,''),nvl(Status,''),nvl(createdBy,''),nvl( DATE_FORMAT(createdDate, '%d %b %Y %l:%i %p') ,'')AS createdDate,nvl(updatedBy,''), nvl( DATE_FORMAT(updatedDate , '%d %b %Y %l:%i %p') ,'')AS updatedDate 
		FROM a_novo_version_controller
		where os = 'IOS'
		order by id Desc`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("VGIV02", lErr2)
			return lIosVersionList, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lIosVersion.Id, &lIosVersion.OS, &lIosVersion.ForceUpdate, &lIosVersion.Version, &lIosVersion.AppStatus, &lIosVersion.CreatedBy, &lIosVersion.CreatedDate, &lIosVersion.UpdatedBy, &lIosVersion.UpdatedDate)
				if lErr3 != nil {
					log.Println("VGIV03", lErr3)
					return lIosVersionList, lErr3
				} else {

					lIosVersionList = append(lIosVersionList, lIosVersion)
				}
			}
		}
	}
	log.Println("getIosVersion(-) ")
	return lIosVersionList, nil

}
