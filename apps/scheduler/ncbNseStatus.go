package scheduler

import (
	"encoding/json"
	"fcs23pkg/apps/Ncb/ncbschedule"
	"fcs23pkg/common"
	"fmt"
	"log"
	"net/http"
)

func NcbStatusScheduler(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	log.Println("NcbStatusScheduler (+)", r.Method)

	if r.Method == "GET" {

		var lRespRec ncbschedule.NcbSchRespStruct

		lRespRec.Status = common.SuccessCode

		lNcbBrokers, lErr1 := ncbschedule.NcbBrokerList()
		if lErr1 != nil {
			log.Println("NCBNSS01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NCBNSS01" + lErr1.Error()
		} else {
			if len(lNcbBrokers) != 0 {
				for _, BrokerStream := range lNcbBrokers {

					if BrokerStream.Exchange == common.NSE {
						lStatusFlag, lErr2 := ncbschedule.NseNcbFetchStatus(BrokerStream.BrokerId, common.AUTOBOT)

						if lErr2 != nil {
							log.Println("NCBNSS02", lErr2)
							lRespRec.Status = common.ErrorCode
							lRespRec.ErrMsg = "NCBNSS02" + lErr2.Error()
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
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NCBNSS03", lErr3)
			fmt.Fprintf(w, "NCBNSS03"+lErr3.Error())
		} else {
			fmt.Fprintf(w, string(lData))
		}

		log.Println("NcbStatusScheduler (-)", r.Method)
	}

}
