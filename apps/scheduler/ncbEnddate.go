package scheduler

import (
	"bytes"
	"encoding/json"
	"fcs23pkg/apps/Ncb/ncbschedule"
	"fcs23pkg/common"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

type ncbOrderval struct {
	NcbOrderList []ncbschedule.NcbIssueList `json:"ncbOrderList"`
	Status       string                     `json:"status"`
	ErrMsg       string                     `json:"errMsg"`
}

func NcbEndDate(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("NcbEndDate (+)", r.Method)

	if r.Method == "GET" {

		var lRespRec ncbOrderval
		time.Sleep(50 * time.Second)

		lBody, lErr1 := ioutil.ReadAll(r.Body)
		log.Println(string(lBody))
		if lErr1 != nil {
			log.Println("NED01", lErr1.Error())
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = lErr1.Error()
		} else {
			lBodyString := string(lBody)
			if lBodyString == "FTC" {

				lNcbIssuesList, lErr2 := ncbschedule.PlacingNcbbOrder(r)
				if lErr2 != nil {
					log.Println("NED02", lErr2)
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = lErr2.Error()

					var lNcbRec ncbschedule.NcbIssueList
					lNcbRec.Symbol = "NED02" + lErr2.Error()
					lNcbIssuesList = append(lNcbIssuesList, lNcbRec)
				} else {
					if lNcbIssuesList == nil {
						var lNcbRec ncbschedule.NcbIssueList

						lNcbRec.Symbol = "No NCB Issues ending today"
						lRespRec.Status = common.ErrorCode
						lNcbIssuesList = append(lNcbIssuesList, lNcbRec)
					}
				}
				lRespRec.NcbOrderList = lNcbIssuesList
			} else {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "Restricted to call API"
			}

		}
		log.Println("response", lRespRec)
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NED03", lErr3)
			fmt.Fprintf(w, "NED03"+lErr3.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}

	}
	log.Println("NcbEndDate (-)", r.Method)
}

func ConstructMailNcborders(pNcbList []ncbschedule.NcbIssueList) (emailUtil.EmailInput, error) {

	log.Println("ConstructMailNcborders (+)")

	var lEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig(".toml/emailconfig.toml")

	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	//============================================================//
	lEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["ToEmailId"])
	// lEmailContent.ToEmailId = "kavyadharshani.m@fcsonline.co.in"

	var lDynamicEmailArr ncbOrderval

	lEmailContent.Subject = "NCB Orders Report"
	html := "html/NcbOrderList.html"

	lTemp, lErr1 := template.ParseFiles(html)
	if lErr1 != nil {
		log.Println("CMNO01", lErr1)
		return lEmailContent, lErr1
	} else {
		var lTpl bytes.Buffer

		lDynamicEmailArr.NcbOrderList = pNcbList
		lTemp.Execute(&lTpl, lDynamicEmailArr)
		lEmailbody := lTpl.String()

		lEmailContent.Body = lEmailbody
	}

	log.Println("ConstructMailNcborders (-)")
	return lEmailContent, nil

}
