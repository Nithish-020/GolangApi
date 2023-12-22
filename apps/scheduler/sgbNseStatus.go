package scheduler

import (
	"encoding/json"
	"fcs23pkg/apps/SGB/sgbschedule"
	"fcs23pkg/common"
	"fmt"
	"log"
	"net/http"
)

func SgbStatusScheduler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("SgbStatusScheduler (+)", r.Method)

	var lRespRec sgbschedule.SchRespStruct
	var lBrokerList []sgbschedule.SgbBrokers

	lRespRec.Status = common.SuccessCode

	lSgbBrokers, lErr1 := sgbschedule.SgbBrokerList()
	if lErr1 != nil {
		log.Println("SGBSS01", lErr1)
		lRespRec.Status = common.ErrorCode
		lRespRec.ErrMsg = "SGBSS01" + lErr1.Error()
	} else {
		lBrokerList = lSgbBrokers
		for _, BrokerStream := range lBrokerList {

			// commented by pavithra
			// if BrokerStream.Exchange == common.BSE {
			// 	go SgbDownStatusSch(&lWg, BrokerStream.BrokerId, common.AUTOBOT)
			// } else
			if BrokerStream.Exchange == common.NSE {
				lErr2 := sgbschedule.NseSgbFetchStatus(BrokerStream.BrokerId, common.AUTOBOT)
				if lErr2 != nil {
					log.Println("SGBSS02", lErr2)
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "SGBSS02" + lErr2.Error()
				} else {
					lRespRec.Status = common.SuccessCode
				}
			}
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
