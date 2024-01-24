package scheduler

import (
	"bytes"
	"encoding/json"
	"fcs23pkg/apps/SGB/sgbschedule"
	"fcs23pkg/common"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

// type sgbIssueList struct {
// 	Symbol                string
// 	JvSuccesssCount       int
// 	JvFailedCount         int
// 	ExchangeSuccessCount  int
// 	ExchangeFailedCount   int
// 	ReverseJvSuccessCount int
// 	ReverseJvFailedCount  int
// }
type sgbOrderval struct {
	SgbOrderList []sgbschedule.SgbIssueList `json:"sgbOrderList"`
	Status       string                     `json:"status"`
	ErrMsg       string                     `json:"errMsg"`
}

func SgbEndDate(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("SgbEndDate (+)", r.Method)
	if r.Method == "POST" {
		var lRespRec sgbOrderval
		// Read the request body values in lBody variable
		lBody, lErr1 := ioutil.ReadAll(r.Body)
		log.Println(string(lBody))
		if lErr1 != nil {
			log.Println("SSED01", lErr1.Error())
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = lErr1.Error()
		} else {
			config := common.ReadTomlConfig("toml/SgbConfig.toml")
			lApiKey := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_API_KEY"])
			// SGB_API_KEY Validation
			lBodyString := string(lBody)
			if lBodyString == lApiKey {
				lSgbIssuesList, lErr2 := sgbschedule.PlacingSgbOrder(r)
				if lErr2 != nil {
					log.Println("SED02", lErr2)
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = lErr2.Error()
					var lSgbRec sgbschedule.SgbIssueList
					lSgbRec.Symbol = "SED02" + lErr2.Error()
					lSgbIssuesList = append(lSgbIssuesList, lSgbRec)
					// lSgbOrders, lErr1 := ConstructMailSgborders(lSgbIssuesList)
					// if lErr1 != nil {
					// 	log.Println("Error in ConstructMailSgborders")
					// } else {
					// 	// log.Println("lSgbOrders Mail Content", lSgbOrders)
					// 	lErr2 := emailUtil.SendEmail(lSgbOrders, "SgbEndDayReport")
					// 	if lErr2 != nil {
					// 		log.Println("lErr2", lErr2)
					// 	}
					// }
				} else {
					if lSgbIssuesList == nil {
						//commented by pavithra no need to send mail in this API Scheduler will send this Email to appsupport
						// lSgbOrders, lErr2 := ConstructMailSgborders(lSgbIssuesList)
						// if lErr2 != nil {
						// 	log.Println("SED02", lErr2)
						// 	lRespRec.Status = common.ErrorCode
						// 	lRespRec.ErrMsg = lErr2.Error()
						// 	log.Println("Error in ConstructMailSgborders")
						// } else {
						// 	// log.Println("lSgbOrders Mail Content", lSgbOrders)
						// 	lErr2 := emailUtil.SendEmail(lSgbOrders, "SgbEndDayReport")
						// 	if lErr2 != nil {
						// 		log.Println("lErr2", lErr2)
						// 	}
						// }
						// } else {
						var lSgbRec sgbschedule.SgbIssueList
						lSgbRec.Symbol = "No SGB Issues ending today"
						lRespRec.Status = common.ErrorCode
						lSgbIssuesList = append(lSgbIssuesList, lSgbRec)

						//commented by pavithra no need to send mail in this API Scheduler will send this Email to appsupport
						// lSgbOrders, lErr1 := ConstructMailSgborders(lSgbIssuesList)
						// if lErr1 != nil {
						// 	log.Println("SED03", lErr1)
						// 	lRespRec.Status = common.ErrorCode
						// 	lRespRec.ErrMsg = lErr1.Error()
						// 	log.Println("Error in ConstructMailSgborders")
						// } else {
						// 	// log.Println("lSgbOrders Mail Content", lSgbOrders)
						// 	lErr2 := emailUtil.SendEmail(lSgbOrders, "SgbEndDayReport")
						// 	if lErr2 != nil {
						// 		log.Println("lErr2", lErr2)
						// 	}
						// }
					}
				}
				lRespRec.SgbOrderList = lSgbIssuesList
			} else {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "Authentication Restricted to call API"
			}
		}
		log.Println("response", lRespRec)
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("SED03", lErr3)
			fmt.Fprintf(w, "SED03"+lErr3.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
	}

	log.Println("SgbEndDate (-)", r.Method)
}

func ConstructMailSgborders(pSgbList []sgbschedule.SgbIssueList) (emailUtil.EmailInput, error) {
	log.Println("ConstructMailSgborders (+)")

	var lEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("toml/emailconfig.toml")

	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	lEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["ToEmailId"])
	var lDynamicEmailArr sgbOrderval

	lEmailContent.Subject = "SGB Orders Report"
	html := "html/SgbOrderList.html"

	lTemp, lErr := template.ParseFiles(html)
	if lErr != nil {
		log.Println("CMSO01", lErr)
		return lEmailContent, lErr
	} else {
		var lTpl bytes.Buffer
		//   ================================================||

		// lEmailContent.ToEmailId = pSgbClientDetails.Mail
		// lEmailContent.ToEmailId = "prashanth.s@fcsonline.co.in"

		// lDynamicEmailArr.SgbOrderList = pspSgbList
		// log.Println("pSgbList", pSgbList)

		lDynamicEmailArr.SgbOrderList = pSgbList
		// log.Println("lDynamicEmailArr.SgbOrderList ", lDynamicEmailArr.SgbOrderList)
		// log.Println("lDynamicEmailArr", lDynamicEmailArr)
		lTemp.Execute(&lTpl, lDynamicEmailArr)
		lEmailbody := lTpl.String()

		lEmailContent.Body = lEmailbody
	}
	log.Println("ConstructMailSgborders (-)")
	return lEmailContent, nil

}
