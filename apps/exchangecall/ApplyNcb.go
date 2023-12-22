package exchangecall

import (
	"fcs23pkg/integration/nse/nsencb"
	"log"
)

func ApplyNseNcb(pReqJson nsencb.NcbAddReqStruct, pUser string, pBrokerId int) (nsencb.NcbAddResStruct, error) {
	log.Println("ApplyNseNcb(+)")

	var lRespJsonRec nsencb.NcbAddResStruct

	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EANN01", lErr1)
		return lRespJsonRec, lErr1
	} else {
		if lToken != "" {
			lResp, lErr2 := nsencb.NcbAddOrder(lToken, pReqJson, pUser)
			if lErr2 != nil {
				log.Println("EANN02", lErr2)
				return lRespJsonRec, lErr2
			} else {
				lRespJsonRec = lResp
				log.Println("lRespJsonRec", lRespJsonRec, lResp, "lResp")
			}
		}
	}
	log.Println("ApplyNseNcb(-)")
	return lRespJsonRec, nil
}
