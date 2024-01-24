package sgbplaceorder

// type ErrorStruct struct {
// 	ErrCode string `json:"errCode"`
// 	ErrMsg  string `json:"errMsg"`
// }

// func BsePlaceOrder(pExchangeReq bsesgb.SgbReqStruct, pClientId string, pReqRec SgbReqStruct, r *http.Request, pExchange string) (bsesgb.SgbRespStruct, string, error) {
// 	log.Println("NsePlaceOrder (+)")

// 	var lRespRec string
// 	// var lErrorRec error
// 	var lError ErrorStruct
// 	// var lExchangeResp nsesgb.SgbAddResStruct

// 	lRespRec = "S"

// 	lExchangeResp, lRespRec, lErr1 := ProcessSgbOrder(pExchangeReq, pClientId, pReqRec, r)
// 	if lErr1 != nil {
// 		lRespRec = "E"
// 		lError.ErrCode = "PBPO01"
// 		lError.ErrMsg = "Exchange Server is Busy right now,Try After Sometime."
// 		return lExchangeResp, lRespRec, lErr1
// 	}

// 	log.Println("lExchangeResp", lExchangeResp)
// 	log.Println("lRespRec", lRespRec)
// 	log.Println("lErr1", lErr1)
// 	log.Println("NsePlaceOrder (-)")
// 	return lExchangeResp, lRespRec, lErr1
// }
