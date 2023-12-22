package localdetail

import (
	"bytes"
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/Ipo/iposchedule"
	"fcs23pkg/apps/scheduler"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/helpers"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// this struct is used to give response for ManualFetch API
type manualOfflineStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

/*
Pupose:This API is used to run the OfflineSchedular from the screen
Parameters:
	not Applicable
Response:
	==========
	*On Sucess
	==========
	{
		"status": "S",
		"errMsg": ""
	}
	=========
	!On Error
	=========
	{
		"status": "E",
		"reason": "Can't able to get the ipo details"
		"errMsg": "Error"
	}
Author: Nithish Kumar
Date: 20NOV2023
*/
func ManualOfflineFetch(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}
	log.Println("lBrokerId", lBrokerId)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {
		log.Println("ManualOfflineFetch (+)", r.Method)

		//create instance for manualOfflineStruct
		var lManualResp manualOfflineStruct

		var lRespRec scheduler.ScheduleStruct
		var lRespArr []scheduler.ScheduleStruct
		var lEmailRec EmailStruct

		lManualResp.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/fetchmaster")
		if lErr1 != nil {
			log.Println("LMOF01", lErr1)
			lManualResp.Status = common.ErrorCode
			lManualResp.ErrMsg = "LMOF01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LMOF01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lManualResp.Status = common.ErrorCode
				lManualResp.ErrMsg = "LMOF02 / UserDetails not Found"
				fmt.Fprintf(w, helpers.GetErrorString("LMOF02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		//This Methods is used for OfflineAppsch
		lNseOfflinesch := iposchedule.NseOfflineSch(lClientId, lBrokerId)
		log.Println("lNseOfflinesch", lNseOfflinesch)
		if lNseOfflinesch.Status == "S" {
			lRespRec.TotalCount = lNseOfflinesch.TotalCount
			lRespRec.SuccessCount = lNseOfflinesch.SuccessCount
			lRespRec.ErrCount = lNseOfflinesch.ErrCount
			lRespRec.Status = common.SuccessCode
			lRespRec.ErrMsg = common.SUCCESS
			lRespRec.SINo = "1"
			lRespRec.Method = "NSE- IPO Offline Applications"
			lRespArr = append(lRespArr, lRespRec)
		} else {
			lRespRec.TotalCount = lNseOfflinesch.TotalCount
			lRespRec.SuccessCount = lNseOfflinesch.SuccessCount
			lRespRec.ErrCount = lNseOfflinesch.ErrCount
			lRespRec.Status = lNseOfflinesch.Status
			lRespRec.ErrMsg = lNseOfflinesch.ErrMsg
			lRespRec.SINo = "1"
			lRespRec.Method = "NSE- IPO Offline Applications"
			lRespArr = append(lRespArr, lRespRec)
		}

		lBseOfflinesch := iposchedule.BseOfflineSch(lClientId, lBrokerId)
		log.Println("lBseOfflinesch", lBseOfflinesch)
		if lBseOfflinesch.Status == "S" {
			lRespRec.TotalCount = lBseOfflinesch.TotalCount
			lRespRec.SuccessCount = lBseOfflinesch.SuccessCount
			lRespRec.ErrCount = lBseOfflinesch.ErrCount
			lRespRec.Status = common.SuccessCode
			lRespRec.ErrMsg = common.SUCCESS
			lRespRec.SINo = "2"
			lRespRec.Method = "BSE- IPO Offline Applications"
			lRespArr = append(lRespArr, lRespRec)
		} else {
			lRespRec.TotalCount = lBseOfflinesch.TotalCount
			lRespRec.SuccessCount = lBseOfflinesch.SuccessCount
			lRespRec.ErrCount = lBseOfflinesch.ErrCount
			lRespRec.Status = lBseOfflinesch.Status
			lRespRec.ErrMsg = lBseOfflinesch.ErrMsg
			lRespRec.SINo = "2"
			lRespRec.Method = "BSE- IPO Offline Applications"
			lRespArr = append(lRespArr, lRespRec)
		}
		lEmailRec.ScheduleReport = lRespArr

		lErr3 := LogEmailSummary(lEmailRec)
		if lErr3 != nil {
			log.Println("LMOF03", lErr3)
			lManualResp.Status = common.ErrorCode
			lManualResp.ErrMsg = "LMOF03 / Error While Emailing"
			fmt.Fprintf(w, helpers.GetErrorString("LMOF02", "Error While Emailing"))
			return
		}

		// Marshal the Response Structure into lData
		lData, lErr4 := json.Marshal(lManualResp)
		if lErr4 != nil {
			log.Println("LMOF04", lErr4)
			fmt.Fprintf(w, helpers.GetErrorString("LMOF04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("ManualOfflineFetch (-)", r.Method)

	}
}

type EmailStruct struct {
	ScheduleReport []scheduler.ScheduleStruct
}

//------------------------------------------
// function to log campaign wise emails
//=====================================

func LogEmailSummary(pSummary EmailStruct) error {
	log.Println("LogEmailSummary (+)")

	var lTpl bytes.Buffer
	lTemp, lErr1 := template.ParseFiles("./html/emailTemplate.html")
	if lErr1 != nil {
		log.Println("ESLES01", lErr1)
		return lErr1
	} else {
		// log.Println("pSummary", pSummary.SchedularArr)
		lTemp.Execute(&lTpl, pSummary)
		lEmailbody := lTpl.String()

		var lEmailRec emailUtil.EmailInput
		lEmailRec.Body = lEmailbody
		//emailRec.Action = constant.INSERT

		//fetch details from toml
		lConfig := common.ReadTomlConfig("./toml/emailconfig.toml")
		lEmailRec.FromDspName = fmt.Sprintf("%v", lConfig.(map[string]interface{})["From"])
		lEmailRec.FromRaw = fmt.Sprintf("%v", lConfig.(map[string]interface{})["FromRaw"])
		//emailRec.EmailServer = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["EmailServer"])
		lEmailRec.ToEmailId = fmt.Sprintf("%v", lConfig.(map[string]interface{})["ToEmailId"])
		lEmailRec.ReplyTo = fmt.Sprintf("%v", lConfig.(map[string]interface{})["ReplyTo"])
		lEmailRec.Subject = fmt.Sprintf("%v", lConfig.(map[string]interface{})["Subject"])

		dt := time.Now().Format("02/Jan/2006 3:04:05 PM")
		lEmailRec.Subject = lEmailRec.Subject + " " + dt
		//fetch details from coresettings
		// emailRec.Subject = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, emailRec.Subject)
		// emailRec.FromRaw = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, emailRec.FromRaw)
		// emailRec.FromDspName = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, emailRec.FromDspName)
		//emailRec.EmailServer = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, emailRec.EmailServer)
		// emailRec.ToEmailId = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, emailRec.ToEmailId)
		// emailRec.ReplyTo = coresettings.GetCoreSettingValue(ftdb.MariaFTPRD, emailRec.ReplyTo)

		lErr2 := emailUtil.SendEmail(lEmailRec, "IPO-OfflineScheduler")
		if lErr2 != nil {
			log.Println("ESLES02", lErr2)
			return lErr2
		}
	}
	log.Println("LogEmailSummary (-)")
	return nil
}
