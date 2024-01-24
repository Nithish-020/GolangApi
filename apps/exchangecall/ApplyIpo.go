package exchangecall

import (
	"fcs23pkg/integration/bse/bseipo"
	"fcs23pkg/integration/nse/nseipo"
	"log"
)

/*
Pupose:  This method is used to inserting the collection of data to the  a_ipo_order_header ,
a_ipo_orderdetails , a_bidTracking tables in  database and also used to place the Bid in NSE
Parameters:
   (ExchangeReqStruct )
Response:
    *On Sucess
    =========
    In case of a successful execution of this method, you will apply for the Bid
	in Exchange using /v1/transaction/addbulk endpoint ad get the response struct from NSE Exchange
    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 09JUNE2023
*/
func ApplyNseIpo(pReqJson []nseipo.ExchangeReqStruct, pUser string, pBrokerId int) ([]nseipo.ExchangeRespStruct, error) {
	log.Println("ApplyNseIpo (+)")

	// var pReqJson []nse.ExchangeReqStruct
	var lRespJsonRec []nseipo.ExchangeRespStruct

	// get token from database usig below method
	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EANI01", lErr1)
		return lRespJsonRec, lErr1
	} else {
		// ReqJson := SetRequest()
		if lToken != "" {
			lResp, lErr2 := nseipo.Transaction(lToken, pReqJson, pUser)
			if lErr2 != nil {
				log.Println("EANI02", lErr2)
				return lRespJsonRec, lErr2
			} else {
				lRespJsonRec = lResp
			}
		} else {
			return lRespJsonRec, lErr1
		}
	}
	log.Println("ApplyNseIpo (-)")
	return lRespJsonRec, nil
}

/*
Pupose:  This method is used to inserting the collection of data to the  a_ipo_order_header ,
a_ipo_orderdetails , a_bidTracking tables in  database and also used to place the Bid in NSE
Parameters:
   (ExchangeReqStruct )
Response:
    *On Sucess
    =========
    In case of a successful execution of this method, you will apply for the Bid
	in Exchange using /v1/transaction/addbulk endpoint ad get the response struct from NSE Exchange
    !On Error
    ========
    In case of any exception during the execution of this method you will get the
    error details. the calling program should handle the error

Author: Pavithra
Date: 09JUNE2023
*/
func ApplyBseIpo(pReqJson bseipo.BseExchangeReqStruct, pUser string, pBrokerId int) (bseipo.BseExchangeRespStruct, error) {
	log.Println("ApplyBseIpo (+)")

	var lRespJsonRec bseipo.BseExchangeRespStruct

	// get token from database usig below method
	lToken, lErr1 := BseGetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EABI01", lErr1)
		return lRespJsonRec, lErr1
	} else {
		// ReqJson := SetRequest()
		if lToken != "" {
			lResp, lErr2 := bseipo.BseIpoOrder(lToken, pReqJson, pUser, pBrokerId)
			if lErr2 != nil {
				log.Println("EABI02", lErr2)
				return lRespJsonRec, lErr2
			} else {
				lRespJsonRec = lResp
			}
		}
	}
	log.Println("ApplyBseIpo (-)")
	return lRespJsonRec, nil
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------
// func ApplyIpo(pReqJson []nseipo.ExchangeReqStruct, pUser string) ([]nseipo.ExchangeRespStruct, error) {
// 	log.Println("ApplyIpo (+)")

// 	// var pReqJson []nse.ExchangeReqStruct
// 	var lRespJsonRec []nseipo.ExchangeRespStruct

// 	// get token from database usig below method
// 	lToken, lErr1 := GetToken(pUser)
// 	if lErr1 != nil {
// 		log.Println("EAI01", lErr1)
// 		return lRespJsonRec, lErr1
// 	} else {
// 		// ReqJson := SetRequest()
// 		if lToken != "" {
// 			lResp, lErr2 := nseipo.Transaction(lToken, pReqJson, pUser)
// 			if lErr2 != nil {
// 				log.Println("EAI02", lErr2)
// 				return lRespJsonRec, lErr2
// 			} else {
// 				lRespJsonRec = lResp
// 			}
// 		} else {
// 			return lRespJsonRec, lErr1
// 		}
// 	}
// 	log.Println("ApplyIpo (-)")
// 	return lRespJsonRec, nil
// }
