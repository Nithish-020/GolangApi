package clientfund

import (
	"encoding/json"
	"fcs23pkg/apigate"
	"fcs23pkg/integration/techexcel"
	"fcs23pkg/util/apiUtil"
	"log"
)

/*
Pupose:This method is used to process JV .
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	AGMPA45767,nil

	==========
	*On Error
	==========
	"",error

Author:Pavithra
Date: 29AUG2023
*/
//commented by pavithra, changed this method name as BOProcessJV
// func ProcessJV(pJvReq techexcel.JvInputStruct, pReqDtl apigate.RequestorDetails) error {
// 	log.Println("ProcessJV (+)")

// 	lErrFromJV := techexcel.JVPostingtoTechxl_UAT(pJvReq, pReqDtl)
// 	if lErrFromJV != nil {
// 		log.Println("SPOPJV01", lErrFromJV)
// 		return lErrFromJV
// 	} else {
// 		log.Println("JV Successfully processed")
// 	}
// 	log.Println("ProcessJV (-)")
// 	return nil
// }

/*
Pupose:This method is used to process BackOffice JV .
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	AGMPA45767,nil

	==========
	*On Error
	==========
	"",error

Author:Pavithra
Date: 29AUG2023
*/
func BOProcessJV(pJvReq techexcel.JvInputStruct, pReqDtl apigate.RequestorDetails) error {
	log.Println("BOProcessJV (+)")

	lErrFromJV := techexcel.JVPostingtoTechxl_UAT(pJvReq, pReqDtl)
	if lErrFromJV != nil {
		log.Println("SPOPJV01", lErrFromJV)
		return lErrFromJV
	} else {
		log.Println("JV Successfully processed")
	}
	log.Println("BOProcessJV (-)")
	return nil
}

type FOJvRespStruct struct {
	RequestTime string `json:"request_time"`
	Status      string `json:"stat"`
	ErrMsg      string `json:"emsg"`
}

/*
Pupose:This method is used to process JV .
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	AGMPA45767,nil

	==========
	*On Error
	==========
	"",error

Author:Nithish Kumar
Date: 21NOV2023
*/

func FOProcessJV(pURL string, pBody string) (FOJvRespStruct, error) {
	log.Println("FOProcessJV (+)")
	//Instance for Front office response
	var lFORespStruct FOJvRespStruct
	// Instance for API request header
	var lHeaderArr []apiUtil.HeaderDetails

	lResp, lErr1 := apiUtil.Api_call(pURL, "POST", pBody, lHeaderArr, "FOProcessJV")
	if lErr1 != nil {
		log.Println("CFPJ01", lErr1)
		return lFORespStruct, lErr1
	} else {
		log.Println("LResponse", string(lResp))
		// Unmarshalling json to struct
		lErr2 := json.Unmarshal([]byte(lResp), &lFORespStruct)
		if lErr2 != nil {
			log.Println("CCFA02", lErr2)
			return lFORespStruct, lErr2
		}
	}
	log.Println("FOProcessJV (-)")
	return lFORespStruct, nil
}
