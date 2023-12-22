package adminaccess

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

type chgPwdReqStruct struct {
	Type     string `json:"type"`
	Member   string `json:"member"`
	LoginId  string `json:"loginId"`
	Password string `json:"password"`
}
type chgPwdRspStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	log.Println("ChangePassword (+)")
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {

		var lReqRec chgPwdReqStruct
		var lRespRec chgPwdRspStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/config")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NCP01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NCP01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("NCP01", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lBody, lErr2 := ioutil.ReadAll(r.Body)
		if lErr2 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NCP02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NCP02", "Cannot able to proccess your request right now.Please try after sometime"))
			return
		} else {
			lErr3 := json.Unmarshal(lBody, &lReqRec)
			log.Println("lReqRec", lReqRec)
			if lErr3 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "NCP03" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("NCP03", "Cannot able to proccess your request right now.Please try after sometime"))
				return
			} else {

				lDb, lErr5 := ftdb.LocalDbConnect(ftdb.IPODB)
				if lErr5 != nil {
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "NCP05" + lErr5.Error()
					fmt.Fprintf(w, helpers.GetErrorString("NCP05", "Cannot able to proccess your request right now.Please try after sometime"))
					return
				} else {
					defer lDb.Close()

					lEncoder := common.EncodeToString(lReqRec.Password)
					sqlString := `update a_ipo_directory set Member = ?,LoginId = ?,Password = ?,UpdatedBy = ?,UpdatedDate = now()  where Stream = ?`

					_, lErr6 := lDb.Exec(sqlString, lReqRec.Member, lReqRec.LoginId, lEncoder, lClientId, lReqRec.Type)
					if lErr6 != nil {
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "NCP06" + lErr6.Error()
						fmt.Fprintf(w, helpers.GetErrorString("NCP06", "Unable to update right now.Please try after sometime"))
						return
					} else {
						lRespRec.Status = common.SuccessCode
						lRespRec.ErrMsg = "Changes were successfully saved"
					}
				}
			}
		}
		// Marshal the response structure into lData
		lData, lErr := json.Marshal(lRespRec)
		if lErr != nil {
			log.Println("NCP07", lErr)
			fmt.Fprintf(w, helpers.GetErrorString("NCP07", "Error found!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}
	log.Println("ChangePassword (-)")
}
