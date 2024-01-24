package scheduler

import (
	"encoding/json"
	"fcs23pkg/apps/SGB/sgbschedule"
	"fcs23pkg/common"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SgbStatusScheduler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("SgbStatusScheduler (+)", r.Method)

	if r.Method == "POST" {

		var lRespRec sgbschedule.SchRespStruct

		lRespRec.Status = common.SuccessCode
		lBody, lErr1 := ioutil.ReadAll(r.Body)
		log.Println(string(lBody))
		if lErr1 != nil {
			log.Println("SGBSS01", lErr1.Error())
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "SGBSS01" + lErr1.Error()
		} else {
			config := common.ReadTomlConfig("toml/SgbConfig.toml")
			lApiKey := fmt.Sprintf("%v", config.(map[string]interface{})["SGB_API_KEY"])
			// SGB_API_KEY Validation
			lBodyString := string(lBody)
			if lBodyString == lApiKey {
				lSgbBrokers, lErr1 := sgbschedule.SgbBrokerList()
				if lErr1 != nil {
					log.Println("SGBSS02", lErr1)
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "SGBSS02" + lErr1.Error()
				} else {
					if len(lSgbBrokers) != 0 {
						for _, BrokerStream := range lSgbBrokers {

							// commented by pavithra
							// if BrokerStream.Exchange == common.BSE {
							// 	go SgbDownStatusSch(&lWg, BrokerStream.BrokerId, common.AUTOBOT)
							// } else
							if BrokerStream.Exchange == common.NSE {
								lStatusFlag, lErr2 := sgbschedule.NseSgbFetchStatus(BrokerStream.BrokerId, common.AUTOBOT)
								if lErr2 != nil {
									log.Println("SGBSS03", lErr2)
									lRespRec.Status = common.ErrorCode
									lRespRec.ErrMsg = "SGBSS03" + lErr2.Error()
								} else {
									if lStatusFlag != common.ErrorCode {
										lRespRec.Status = common.SuccessCode
										lRespRec.ErrMsg = common.SUCCESS
									} else {
										lRespRec.Status = common.ErrorCode
										lRespRec.ErrMsg = common.FAILED
									}
								}
							}
						}
					} else {
						log.Println("No Brokers Found")
					}
				}
			} else {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "Authentication Restricted to call API"
			}
		}
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("SGBSS03", lErr3)
			fmt.Fprintf(w, "SGBSS03"+lErr3.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SgbStatusScheduler (-)", r.Method)
	}
}
