package ncbplaceorder

type FundRespStruct struct {
	AccountBalance float64 `json:"accountBalance"`
	Status         string  `json:"status"`
	ErrMsg         string  `json:"errMsg"`
}

// commented by pavithra
// NCB fund checking function is changed
// func VerifyFundDetails(pClientId string) (FundRespStruct, error) {
// 	log.Println("VerifyFundDetails(+)")
// 	var lClientFundDetailRec clientfund.ClientFundStruct
// 	//
// 	var lRespRec FundRespStruct
// 	lClientFundDetailArr, lErr1 := clientfund.VerifyFund(pClientId)
// 	if lErr1 != nil {
// 		log.Println("NPVFD01", lErr1)
// 		return lRespRec, lErr1
// 	} else {
// 		log.Println(lClientFundDetailArr)
// 		if lClientFundDetailArr != nil {
// 			for lidx := 0; lidx < len(lClientFundDetailArr); lidx++ {
// 				lClientFundDetailRec = lClientFundDetailArr[lidx]
// 			}
// 			if lClientFundDetailRec.AccountCode == "" && lClientFundDetailRec.Amount == "" {
// 				log.Println("lClientFundDetailRec", lClientFundDetailRec)
// 				lRespRec.AccountBalance = 0.0
// 				lRespRec.Status = common.ErrorCode
// 			} else {
// 				lRespRec.AccountBalance, _ = strconv.ParseFloat(lClientFundDetailRec.Amount, 64)
// 				lRespRec.Status = common.SuccessCode
// 			}
// 		} else {
// 			lRespRec.AccountBalance = 0.0
// 			lRespRec.Status = common.ErrorCode
// 		}
// 	}
// 	log.Println("VerifyFundDetails(-)")
// 	return lRespRec, nil
// }
