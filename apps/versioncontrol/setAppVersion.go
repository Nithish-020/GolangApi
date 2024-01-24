package versioncontrol

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type novoResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

// Addin Router Name is Pending For Verifiy Access to the Client
func SetAppVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("SetAppVersion (+)")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "PUT" {
		var lAppReq novoAppStruct
		var lAppResp novoResp
		lAppResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/versioncontrol")
		if lErr1 != nil {
			lAppResp.Status = common.ErrorCode
			lAppResp.ErrMsg = "VSNA01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("VSNA01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lAppResp.Status = common.ErrorCode
				lAppResp.ErrMsg = "VSNA02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("VSNA02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lAppResp.Status = common.ErrorCode
			lAppResp.ErrMsg = "VSNA03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("VSNA03", "Unable To Read The Body"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lAppReq)

			if lErr4 != nil {
				lAppResp.Status = common.ErrorCode
				lAppResp.ErrMsg = "VSNA04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("VSNA04", "Unable to unmarshal the Data"))
				return
			} else {

				lFlag, lErr5 := checkVersionExist(lAppReq)
				if lErr5 != nil {
					lAppResp.Status = common.ErrorCode
					lAppResp.ErrMsg = "VSNA05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("VSNA05", "Unable to Check an Version in DB"))
					return
				} else {
					if lAppReq.Id == 0 {
						if lFlag == "Y" {
							lErr6 := insertNewVersion(lAppReq, lClientId)
							if lErr6 != nil {
								lAppResp.Status = common.ErrorCode
								lAppResp.ErrMsg = "VSNA06" + lErr6.Error()
								fmt.Fprintf(w, helpers.GetErrorString("VSNA06", "Unable to Insert an New Version"))
								return
							} else {
								lAppResp.Status = common.SuccessCode
								lAppResp.ErrMsg = "Version Inserted SuccesFully"
							}
						} else {
							lAppResp.Status = common.ErrorCode
							lAppResp.ErrMsg = "Version Already Exist"

						}
					} else {
						lErr7 := updateVersion(lAppReq, lClientId)
						if lErr7 != nil {
							lAppResp.Status = common.ErrorCode
							lAppResp.ErrMsg = "VSNA06" + lErr7.Error()
							fmt.Fprintf(w, helpers.GetErrorString("VSNA07", "Unable to update Version Right Now "))
							return
						} else {
							lAppResp.Status = common.SuccessCode
							lAppResp.ErrMsg = "Version Updated SuccesFully"
						}
					}
				}

			}
		}
		lData, lErr2 := json.Marshal(lAppResp)
		if lErr2 != nil {
			log.Println("VSNA01", lErr2)
			lAppResp.Status = common.ErrorCode
			lAppResp.ErrMsg = "VSNA01" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("VSNA01", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SetAppVersion (-)", r.Method)
	}
}

func insertNewVersion(pNovo novoAppStruct, pClientId string) error {
	log.Println("insertVersion(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("VINV01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		lCoreString := `INSERT INTO a_novo_version_controller
			(os, forceUpdate, version, Status,createdBy, createdDate, updatedBy, updatedDate)
			values( ?, ?, ?,?,?, now(),?,now())`

		_, lErr2 := lDb.Exec(lCoreString, pNovo.OS, pNovo.ForceUpdate, pNovo.Version, pNovo.AppStatus, pClientId, pClientId)
		if lErr2 != nil {
			log.Println("VINV02", lErr2)
			return lErr2
		}
	}
	log.Println("insertVersion(-)")
	return nil
}

func updateVersion(pNovo novoAppStruct, pClientId string) error {
	log.Println("updateRole(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("VUV01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lCoreString := `UPDATE a_novo_version_controller
		SET  forceUpdate = ?, Status= ?, updatedBy = ?, updatedDate= now()
		WHERE id= ?; `
		_, lErr2 := lDb.Exec(lCoreString, pNovo.ForceUpdate, pNovo.AppStatus, pClientId, pNovo.Id)
		if lErr2 != nil {
			log.Println("VUV02", lErr2)
			return lErr2
		}
	}
	log.Println("updateRole(-)")

	return nil

}

func checkVersionExist(pNovo novoAppStruct) (string, error) {
	log.Println("checkVersionExist (+)")
	var lFlag string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("VCVE01", lErr1)
		return lFlag, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select count(version)  from a_novo_version_controller 
		where version = ? and os = ?  `
		lRows, lErr2 := lDb.Query(lCoreString, pNovo.Version, pNovo.OS)
		if lErr2 != nil {
			log.Println("VCVE02", lErr2)
			return lFlag, lErr2
		} else {
			Count := 0
			for lRows.Next() {
				lErr3 := lRows.Scan(&Count)
				if lErr3 != nil {
					log.Println("VCVE03", lErr3)
					return lFlag, lErr3
				} else {
					if Count == 0 {
						lFlag = "Y"

					} else {
						lFlag = "N"
					}

				}
			}

		}
	}
	log.Println("checkVersionExist(-)")

	return lFlag, nil

}
