package registrardetails

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

type RegResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetRegisterDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("SetRegisterDetails (+)")
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
		var lRegResp RegResp
		var lRegReq RegistrarDetails
		lRegResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipoinfo")
		if lErr1 != nil {
			lRegResp.Status = common.ErrorCode
			lRegResp.ErrMsg = "RSRD01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSRD01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				var lErr2 error
				lRegResp.Status = common.ErrorCode
				lRegResp.ErrMsg = "RSRD02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RSRD02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lRegResp.Status = common.SuccessCode

		lBody, lErr3 := ioutil.ReadAll(r.Body)
		if lErr3 != nil {
			lRegResp.Status = common.ErrorCode
			lRegResp.ErrMsg = "RSRD03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSRD03", "Request cannot be empty.Please try after sometime"))
			return
		} else {
			lErr4 := json.Unmarshal(lBody, &lRegReq)
			if lErr4 != nil {
				lRegResp.Status = common.ErrorCode
				lRegResp.ErrMsg = "RSRD04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("RSRD04", "Unable to proccess your request right now.Please try after sometime"))
				return
			} else {
				lFlag, lErr5 := CheckRegistrarExist(lRegReq)
				if lErr5 != nil {
					lRegResp.Status = common.ErrorCode
					lRegResp.ErrMsg = "RSRD05" + lErr4.Error()
					fmt.Fprintf(w, helpers.GetErrorString("RSRD05", "Unable to Check the Registrar"))
					return

				} else {
					if lRegReq.Id == 0 {
						if lFlag == "Y" {
							lErr6 := InsertRegistrar(lRegReq, lClientId)
							if lErr6 != nil {
								lRegResp.Status = common.ErrorCode
								lRegResp.ErrMsg = "RSRD06" + lErr4.Error()
								fmt.Fprintf(w, helpers.GetErrorString("RSRD06", "Unable to Insert your request right now.Please try after sometime"))
								return
							} else {
								lRegResp.Status = common.SuccessCode
								lRegResp.ErrMsg = "Registrar Added SuccessFully"
							}
						} else {
							lRegResp.Status = common.ErrorCode
							lRegResp.ErrMsg = "Registrar Already Exist"
						}
					} else {
						lErr7 := updateRegistrar(lRegReq, lClientId)
						if lErr7 != nil {
							lRegResp.Status = common.ErrorCode
							lRegResp.ErrMsg = "RSRD07" + lErr4.Error()
							fmt.Fprintf(w, helpers.GetErrorString("RSRD07", "Unable to Update your request right now.Please try after sometime"))
							return
						} else {
							lRegResp.Status = common.SuccessCode
							lRegResp.ErrMsg = "Registrar Updated SuccessFully"
						}
					}
				}

			}
		}
		lData, lErr8 := json.Marshal(lRegResp)
		if lErr8 != nil {
			log.Println("RSRD08", lErr8)
			lRegResp.Status = common.ErrorCode
			lRegResp.ErrMsg = "RSRD08" + lErr8.Error()
			fmt.Fprintf(w, helpers.GetErrorString("RSRD08", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SetRegisterDetails(-)")
	}

}

func InsertRegistrar(pRegReq RegistrarDetails, pClientId string) error {
	log.Println("InsertRegistrar(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RIR01", lErr1.Error())
		return lErr1
	} else {
		defer lDb.Close()
		lSqlString := `INSERT INTO a_ipo_registrars
		(RegistrarName, RegistrarLink, createdBy, createdDate, updatedBy, updatedDate)
		VALUES( upper(TRIM(?)),Trim(?), ?, now(), ? , now());`
		_, lErr2 := lDb.Exec(lSqlString, pRegReq.RegistrarName, pRegReq.RegistrarLink, pClientId, pClientId)
		if lErr2 != nil {
			log.Println("RIR02", lErr2)
			return lErr2
		}
	}

	log.Println("InsertRegistrar(-)")
	return nil
}

func updateRegistrar(pRegReq RegistrarDetails, pClientId string) error {
	log.Println("updateRegistrar(+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RUR01", lErr1.Error())
		return lErr1
	} else {
		defer lDb.Close()
		lSqlString := `UPDATE a_ipo_registrars
		SET RegistrarName = upper(Trim(?)), RegistrarLink=Trim(?), updatedBy = ?, updatedDate= now()
		WHERE id= ?;`
		_, lErr2 := lDb.Exec(lSqlString, pRegReq.RegistrarName, pRegReq.RegistrarLink, pClientId, pRegReq.Id)
		if lErr2 != nil {
			log.Println("RUR02", lErr2)
			return lErr2
		}
	}

	log.Println("updateRegistrar(-)")
	return nil
}

func CheckRegistrarExist(pReg RegistrarDetails) (string, error) {
	log.Println("CheckRegistrarExist(+)")
	var lFlag string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("RCRE01", lErr1.Error())
		return lFlag, lErr1
	} else {
		defer lDb.Close()
		lSqlString := `SELECT count(RegistrarName)
		FROM a_ipo_registrars
		where RegistrarName = ?`
		lRows, lErr2 := lDb.Query(lSqlString, pReg.RegistrarName)
		if lErr2 != nil {
			log.Println("RCRE02", lErr2)
			return lFlag, lErr2

		} else {
			lcount := 0
			for lRows.Next() {
				lErr3 := lRows.Scan(&lcount)
				if lErr3 != nil {
					log.Println("RCRE03", lErr3)
					return lFlag, lErr3
				}
			}
			if lcount == 0 {
				lFlag = "Y"
			} else {
				lFlag = "N"
			}
		}
	}

	log.Println("CheckRegistrarExist(-)")
	return lFlag, nil
}
