package sgbplaceorder

import (
	// "fcs23pkg/apigate"

	"fcs23pkg/common"
	"fcs23pkg/integration/techexcel"
	"fmt"
	"log"
	"net/http"
	"time"
)

type JvDataStruct struct {
	MasterId    int    `json:"masterId"`
	Unit        string `json:"unit"`
	Price       string `json:"price"`
	JVamount    string `json:"jvAmount"`
	OrderNo     string `json:"orderNo"`
	ClientId    string `json:"clientId"`
	ActionCode  string `json:"actionCode"`
	Transaction string `json:"transaction"`
	Flag        string `json:"flag"`
	BidId       string `json:"bidId"`
}

type JvStatusStruct struct {
	JvAmount    string `json:"jvAmount"`
	JvStatus    string `json:"jvStatus"`
	JvStatement string `json:"jvStatement"`
	JvType      string `json:"jvType"`
}

type FailedJvStruct struct {
	SINo        int
	ClientId    string
	Amount      string
	Transaction string
}
type DynamicEmailStruct struct {
	Date        string
	FailedJvArr []FailedJvStruct
}

/*
Pupose:This method inserting the order head values in order header table.
Parameters:

	pReqArr,pMasterId,PClientId

Response:

		==========
		*On Sucess
		==========

		==========
		*On Error
		==========

Author:Pavithra
Date: 12JUNE2023
*/

// func ReverseProcess(pExchangeReq bsesgb.SgbReqStruct, pExchangeResp bsesgb.SgbRespStruct, pReqRec SgbReqStruct, pClientId string, pJvStatusRec JvStatusStruct, pJvDataStruct JvDataStruct, pClientEmail string, pExchange string, pBrokerId int) (string, error) {
// 	log.Println("ReverseProcess (+)")

// 	var lStatus string
// 	var lBidId string
// 	var lOrderNo string
// 	lCheckRev := "R"
// 	// log.Println("pExchangeReq", pExchangeReq)
// 	lExchangeReq, lErr1 := ConstructSGBReqStruct(pReqRec, pClientId)
// 	if lErr1 != nil {
// 		log.Println("SGBPRP01", lErr1.Error())
// 		lStatus = common.ErrorCode
// 		return lStatus, lErr1
// 	} else {
// 		log.Println("lExchangeReq", lExchangeReq)
// 		// lExchangeResp, lErr2 := sgbexchangecall.ApplySgb(pExchangeReq, pClientId)
// 		lJVprocessData, lErr2 := getJvData(pClientId, pReqRec, pExchangeReq, lCheckRev)
// 		if lErr2 != nil {
// 			log.Println("SGBPRP02", lErr2.Error())
// 			lStatus = common.ErrorCode
// 			return lStatus, lErr2
// 		} else {
// 			for lIdx := 0; lIdx < len(pExchangeResp.Bids); lIdx++ {
// 				lBidId = pExchangeResp.Bids[lIdx].BidId
// 				lOrderNo = pExchangeResp.Bids[lIdx].OrderNo
// 			}
// 			for lIdx := 0; lIdx < len(lExchangeReq.Bids); lIdx++ {
// 				lExchangeReq.Bids[lIdx].ActionCode = lJVprocessData.ActionCode
// 				if lJVprocessData.ActionCode == "D" {
// 					lExchangeReq.Bids[lIdx].BidId = lBidId
// 					lExchangeReq.Bids[lIdx].OrderNo = lOrderNo
// 				} else {
// 					lExchangeReq.Bids[lIdx].BidId = lJVprocessData.BidId
// 					lExchangeReq.Bids[lIdx].OrderNo = lJVprocessData.OrderNo
// 					lExchangeReq.Bids[lIdx].SubscriptionUnit = lJVprocessData.Unit
// 				}
// 				lExchangeReq.Bids[lIdx].Rate = lJVprocessData.Price
// 			}

// 			//to insert details in sgb bidtracking Table
// 			lErr3 := InsertSgbBidTrack(lExchangeReq, pClientId, pExchange)
// 			if lErr3 != nil {
// 				log.Println("SGBPRP03", lErr3.Error())
// 				lStatus = common.ErrorCode
// 				return lStatus, lErr3
// 			} else {

// 				log.Println("lJVprocessData", lJVprocessData)
// 				// lJVprocessData, lErr3 := getJvData(pClientId, pReqRec, pExchangeReq)
// 				lExchangeResp, lErr4 := exchangecall.ApplyBseSgb(lExchangeReq, pClientId, pB)
// 				if lErr4 != nil {
// 					log.Println("SGBPRP04", lErr4.Error())
// 					lStatus = common.ErrorCode
// 					return lStatus, lErr4
// 				} else {
// 					for lIdx := 0; lIdx < len(lExchangeResp.Bids); lIdx++ {
// 						if lExchangeResp.Bids[lIdx].ErrorCode == "0" && lExchangeResp.StatusCode == "0" {
// 							pJvStatusRec.JvStatus = "S"
// 							lErr5 := ResponseUpdate(lExchangeResp, lExchangeReq, pReqRec, pClientId, pJvStatusRec, pClientEmail, pExchange)
// 							if lErr5 != nil {
// 								log.Println("SGBPRP05", lErr5.Error())
// 								lStatus = common.ErrorCode
// 								return lStatus, lErr5
// 							} else {
// 								// lErr6 := ResponseUpdate(pExchangeResp, pExchangeReq, pReqRec, pInsertedId, "", pClientId, pJvStatusRec)
// 								lErr6 := UpdateSgbBidTrack(pExchangeResp, 0, pClientId, pJvStatusRec, pExchange)
// 								if lErr6 != nil {
// 									log.Println("SGBPRP06", lErr6.Error())
// 									lStatus = common.ErrorCode
// 									return lStatus, lErr6
// 								} else {
// 									lStatus = common.SuccessCode
// 									log.Println("Updated BID track Successfully Successfully")
// 								}
// 							}
// 						} else {
// 							pJvStatusRec.JvStatus = "E"
// 							lStatus = common.SuccessCode
// 							//-----Mail To Accounts Team-------
// 							lFailedJvMail, lErr7 := constructJVMail(pJvStatusRec, pJvDataStruct)
// 							if lErr7 != nil {
// 								log.Println("SGBPRP07", lErr7.Error())
// 								lStatus = common.ErrorCode
// 								return lStatus, lErr7
// 							} else {
// 								lString := "JV Failed in SGB Order"
// 								lErr8 := emailUtil.SendEmail(lFailedJvMail, lString)
// 								if lErr8 != nil {
// 									log.Println("SGBPRP08", lErr8.Error())
// 									lStatus = common.ErrorCode
// 									return lStatus, lErr8
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 	}
// 	log.Println("ReverseProcess (-)")
// 	return lStatus, nil
// }

/*
Pupose:This method is used to get the pan number for order input.
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
func BlockClientFund(pJvProcessRec JvDataStruct, pRequest *http.Request) JvStatusStruct {
	log.Println("BlockClientFund (+)")

	// this variables is used to get Status
	var lJVReq techexcel.JvInputStruct
	// var lReqDtl apigate.RequestorDetails
	var lJvStatusRec JvStatusStruct

	// lReqDtl = apigate.GetRequestorDetail(pRequest)
	lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
	lCocd := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["PaymentCode"])

	lJVReq.COCD = lCocd
	lJVReq.VoucherDate = time.Now().Format("02/01/2006")
	lJVReq.BillNo = "SGB" + pJvProcessRec.ClientId
	lJVReq.SourceTable = "a_sgb_orderdetails"
	lJVReq.SourceTableKey = pJvProcessRec.OrderNo
	lJVReq.Amount = pJvProcessRec.JVamount
	lJVReq.WithGST = "N"

	// lConfigFile := common.ReadTomlConfig("toml/config.toml")
	lFTCaccount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseSGBorderAccount"])

	if pJvProcessRec.Transaction == "F" {
		lJVReq.AccountCode = lFTCaccount
		lJVReq.CounterAccount = pJvProcessRec.ClientId
		// if pJvProcessRec.Flag == "Reversal" {
		// 	lJVReq.Narration = "Reversal on SGB purchase  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		// 	// lJvStatusRec.JvStatement = "Credit to client, " + lJVReq.Narration
		// 	lJvStatusRec.JvStatement = lJVReq.Narration
		// 	lJvStatusRec.JvType = "C"
		// } else if pJvProcessRec.Flag == "Cancel" {
		// 	lJVReq.Narration = "Cancellation on SGB purchase  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		// 	// lJvStatusRec.JvStatement = "Credit to client, " + lJVReq.Narration
		// 	lJvStatusRec.JvStatement = lJVReq.Narration
		// 	lJvStatusRec.JvType = "C"
		// } else if pJvProcessRec.Flag == "Failed" {
		lJVReq.Narration = "Failed of SGB Order Refund  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		// lJvStatusRec.JvStatement = "Credit to client, " + lJVReq.Narration
		lJvStatusRec.JvStatement = lJVReq.Narration
		lJvStatusRec.JvType = "C"
		// }
	} else if pJvProcessRec.Transaction == "C" {
		lJVReq.AccountCode = pJvProcessRec.ClientId
		lJVReq.CounterAccount = lFTCaccount
		// if pJvProcessRec.Flag == "Addon" {
		lJVReq.Narration = "SGB Addditonal purchase  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		// lJvStatusRec.JvStatement = "Debit to client, " + lJVReq.Narration
		lJvStatusRec.JvStatement = lJVReq.Narration
		lJvStatusRec.JvType = "D"
		// } else if pJvProcessRec.Flag == "Purchase" {
		// 	lJVReq.Narration = "SGB Purchase  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		// 	// lJvStatusRec.JvStatement = "Debit to client, " + lJVReq.Narration
		// 	lJvStatusRec.JvStatement = lJVReq.Narration
		// 	lJvStatusRec.JvType = "D"
		// }
	}

	lJvStatusRec.JvStatus = "S"
	lJvStatusRec.JvAmount = pJvProcessRec.JVamount
	log.Println("JV Record:", lJVReq)
	//JV processing method
	// lErr1 := clientfund.ProcessJV(lJVReq, lReqDtl)
	// if lErr1 != nil {
	// 	log.Println("SPOPJV01", lErr1.Error())
	// 	lJvStatusRec.JvStatus = "E"
	// 	return lJvStatusRec
	// } else {
	// 	lJvStatusRec.JvStatus = "S"
	// }
	// log.Println("BlockClientFund (-)")
	return lJvStatusRec
}

/*
Pupose:This method is used to get the  Client order details .
Parameters:

	PClientId,PorderNo

Response:

	==========
	*On Sucess
	==========
	{
   Orderno:123556845
	unit:5
	Price :6000
	Activity :M
	Amount:24000
	OrderDate:24-09-2023
    ClientMail:'prashanth.s@fcsOnline.co.in'
	*On Error
	==========
	{},error

Author:Pavithra
Date: 29AUG2023
*/

// func GetClientEmailId(r *http.Request, pClientId string) (string, error) {
// 	log.Println("GetClientEmailId (+)")

// 	// this variables is used to get Pan number of the client from the database.
// 	var lEmailId string

// 	publicTokenCookie, lErr1 := r.Cookie(common.ABHICookieName)
// 	if lErr1 != nil {
// 		log.Println("SGBGCE01", lErr1)
// 		return lEmailId, lErr1
// 	}

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr2 != nil {
// 		log.Println("SGBGCE02", lErr2)
// 		return lEmailId, lErr2
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select UserMailId
// 						from novo_token
// 						where Token = ?
// 						and UserId = ? `
// 		lRows, lErr3 := lDb.Query(lCoreString, publicTokenCookie.Value, pClientId)
// 		if lErr3 != nil {
// 			log.Println("SGBGCE03", lErr3)
// 			return lEmailId, lErr3
// 		} else {
// 			for lRows.Next() {
// 				lErr3 = lRows.Scan(&lEmailId)
// 				if lErr3 != nil {
// 					log.Println("SGBGCE04", lErr3)
// 					return lEmailId, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetClientEmailId (-)")
// 	return lEmailId, nil
// }

/*
Pupose:This method is used to get the  Client order details .
Parameters:

	PClientId,PorderNo

Response:

	==========
	*On Sucess
	==========
	{
   Orderno:123556845
	unit:5
	Price :6000
	Activity :M
	Amount:24000
	OrderDate:24-09-2023
    ClientMail:'prashanth.s@fcsOnline.co.in'
	*On Error
	==========
	{},error

Author:Prashanth
Date: 8SEP2023
*/

// func fetchSgbClientorder(pClientId string, pExchangeResp bsesgb.SgbRespStruct) (SgbClientDetail, error) {
// 	log.Println("fetchSgbClientorder (+)")

// 	//
// 	var lClientOrder SgbClientDetail

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SGBFSC01", lErr1)
// 		return lClientOrder, lErr1
// 	} else {
// 		defer lDb.Close()

// 		for lIdx := 0; lIdx < len(pExchangeResp.Bids); lIdx++ {

// 			lCoreString := `select h.ApplicantName ,date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date ,nvl(h.ClientEmail,'') mail ,m.Symbol
// 							,d.ReqSubscriptionUnit,d.ReqRate,(d.ReqSubscriptionUnit * d.ReqRate) as Amount,
// 							(CASE
// 								WHEN d.ActionCode = 'M' THEN 'Modify'
// 								WHEN d.ActionCode = 'N' THEN 'New'
// 								WHEN d.ActionCode = 'D' THEN 'Delete'
// 								ELSE d.ActionCode
// 							END ) AS ActionDescription,d.OrderNo
// 							from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
// 							where h.MasterId = m.id
// 							and d.HeaderId = h.Id
// 							and  h.ClientId = ? and d.OrderNo = ?`

// 			lRows, lErr2 := lDb.Query(lCoreString, pClientId, pExchangeResp.Bids[lIdx].OrderNo)
// 			if lErr2 != nil {
// 				log.Println("SGBFSC02", lErr2)
// 				return lClientOrder, lErr2
// 			} else {
// 				for lRows.Next() {
// 					lErr3 := lRows.Scan(&lClientOrder.ClientName, &lClientOrder.OrderDate, &lClientOrder.Mail, &lClientOrder.Symbol, &lClientOrder.Unit, &lClientOrder.Price, &lClientOrder.Amount, &lClientOrder.Activity, &lClientOrder.OrderNo)

// 					if lErr3 != nil {
// 						log.Println("SGBFSC03", lErr3)
// 						return lClientOrder, lErr3
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("fetchSgbClientorder (-)")
// 	return lClientOrder, nil
// }

/*

Pupose:This method is used to get the email Input for success Sgb Place Order  .
Parameters:

	pSgbClientDetails {},pStatus string

Response:

	==========
	*On Sucess
	==========
	{
     Name: clientName
     Status: S
     OrderDate: 8SEP2023
     OrderNumber: 121345687
     Symbol: sgb test
     Unit: 5
     Price: 500
     Amount:2500
     Activity : M
	},nil

	==========
	*On Error
	==========
	"",error

Author:PRASHANTH
Date: 8SEP2023
*/

// func constructSuccessmail(pSgbClientDetails SgbClientDetail, pStatus string) (emailUtil.EmailInput, error) {
// 	log.Println("constructSuccessmail (+)")
// 	type dynamicEmailStruct struct {
// 		Name        string
// 		Status      string
// 		OrderDate   string
// 		Symbol      string
// 		OrderNumber string
// 		Unit        string
// 		Price       string
// 		Amount      int
// 		Activity    string
// 	}

// 	var lEmailContent emailUtil.EmailInput
// 	config := common.ReadTomlConfig("./toml/emailconfig.toml")

// 	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// 	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
// 	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 	lEmailContent.Subject = "SGB Order"
// 	html := "html/SgbOrderTemplate.html"

// 	lTemp, lErr := template.ParseFiles(html)
// 	if lErr != nil {
// 		log.Println("IMP03", lErr)
// 		return lEmailContent, lErr
// 	} else {
// 		var lTpl bytes.Buffer
// 		var lDynamicEmailVal dynamicEmailStruct
// 		lDynamicEmailVal.Name = pSgbClientDetails.ClientName
// 		lDynamicEmailVal.Amount = pSgbClientDetails.Amount
// 		lDynamicEmailVal.Unit = pSgbClientDetails.Unit
// 		lDynamicEmailVal.Price = pSgbClientDetails.Price
// 		lDynamicEmailVal.OrderDate = pSgbClientDetails.OrderDate
// 		lDynamicEmailVal.OrderNumber = pSgbClientDetails.OrderNo
// 		lDynamicEmailVal.Symbol = pSgbClientDetails.Symbol
// 		lDynamicEmailVal.Activity = pSgbClientDetails.Activity
// 		//   ================================================
// 		if pStatus == "S" {
// 			lDynamicEmailVal.Status = common.SUCCESS
// 		} else {
// 			lDynamicEmailVal.Status = common.FAILED
// 		}
// 		//   ================================================||

// 		// lEmailContent.ToEmailId = pSgbClientDetails.Mail
// 		lEmailContent.ToEmailId = "pavithra.v@fcsonline.co.in"

// 		lTemp.Execute(&lTpl, lDynamicEmailVal)
// 		lEmailbody := lTpl.String()

// 		lEmailContent.Body = lEmailbody
// 	}
// 	log.Println("constructSuccessmail (-)")
// 	return lEmailContent, nil
// }

// func constructJVMail(pJvStatusRec JvStatusStruct, pJvDataStruct JvDataStruct) (emailUtil.EmailInput, error) {
// 	log.Println("constructJVMail (+)")
// 	var lEmailContent emailUtil.EmailInput
// 	config := common.ReadTomlConfig("./toml/emailconfig.toml")

// 	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// 	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 	lEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["AccountsTeamMail"])

// 	lEmailContent.Subject = "Failed JV For SGB Orders"
// 	html := "html/JvInvoice.html"

// 	lTemp, lErr := template.ParseFiles(html)
// 	if lErr != nil {
// 		log.Println("SPOCJM01", lErr)
// 		return lEmailContent, lErr
// 	} else {
// 		var lTpl bytes.Buffer
// 		count := 1
// 		var lDynamicEmailVal FailedJvStruct
// 		var lDynamicEmailArr DynamicEmailStruct

// 		lDynamicEmailVal.SINo = count
// 		lDynamicEmailVal.ClientId = pJvDataStruct.ClientId
// 		lDynamicEmailVal.Amount = pJvStatusRec.JvAmount
// 		if pJvStatusRec.JvType == "D" {
// 			lDynamicEmailVal.Transaction = "Debit"
// 		} else if pJvStatusRec.JvType == "C" {
// 			lDynamicEmailVal.Transaction = "Credit"
// 		}
// 		currentTime := time.Now()
// 		// Format the current time in "11-09-2023" format
// 		formattedDate := currentTime.Format("02-01-2006")
// 		lDynamicEmailArr.Date = formattedDate
// 		lDynamicEmailArr.FailedJvArr = append(lDynamicEmailArr.FailedJvArr, lDynamicEmailVal)
// 		lTemp.Execute(&lTpl, lDynamicEmailArr)
// 		lEmailbody := lTpl.String()

// 		lEmailContent.Body = lEmailbody

// 		log.Println("constructJVMail (-)")
// 	}
// 	return lEmailContent, nil
// }
