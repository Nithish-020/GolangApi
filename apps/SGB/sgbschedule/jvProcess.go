package sgbschedule

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fcs23pkg/apigate"
// 	"fcs23pkg/apps/SGB/sgbplaceorder"
// 	"fcs23pkg/apps/exchangecall"
// 	"fcs23pkg/common"
// 	"fcs23pkg/ftdb"
// 	"fcs23pkg/integration/bse/bsesgb"
// 	"fcs23pkg/integration/clientfund"
// 	"fcs23pkg/integration/nse/nsesgb"
// 	"fcs23pkg/integration/techexcel"
// 	"fcs23pkg/util/emailUtil"
// 	"fmt"
// 	"log"
// 	"math"
// 	"net/http"
// 	"strconv"
// 	"text/template"
// 	"time"
// )

// type SgbClientDetail struct {
// 	OrderNo    string `json:"orderno"`
// 	Unit       string `json:"unit"`
// 	Price      string `json:"price"`
// 	Symbol     string `json:"symbol"`
// 	OrderDate  string `json:"orderdate"`
// 	Amount     int    `json:"amount"`
// 	Mail       string `json:"mail"`
// 	ClientName string `json:"clientname"`
// 	Activity   string `json:"activity"`
// }

// type JvDataStruct struct {
// 	Unit        string `json:"unit"`
// 	Price       string `json:"price"`
// 	JVamount    string `json:"jvAmount"`
// 	OrderNo     string `json:"orderNo"`
// 	ClientId    string `json:"clientId"`
// 	ActionCode  string `json:"actionCode"`
// 	Transaction string `json:"transaction"`
// 	BidId       string `json:"bidId"`
// }

// type JvStatusStruct struct {
// 	BoJvStatus    string `json:"boJvstatus"`
// 	BoJvStatement string `json:"boJvstatement"`
// 	BoJvAmount    string `json:"boJvamount"`
// 	BoJvType      string `json:"boJvtype"`
// 	FoJvStatus    string `json:"foJvstatus"`
// 	FoJvStatement string `json:"foJvstatement"`
// 	FoJvAmount    string `json:"foJvamount"`
// 	FoJvType      string `json:"foJvtype"`
// }

// type DynamicEmailStruct struct {
// 	Date        string
// 	FailedJvArr []FailedJvStruct
// }

// type SgbStruct struct {
// 	MasterId int
// 	Symbol   string
// 	Exchange string
// }

// // this struct is used to get the Frontoffice JV details.
// type FOJvStruct struct {
// 	UserId       string `json:"uid"`
// 	AccountId    string `json:"actid"`
// 	Amount       string `json:"amt"`
// 	SourceUserId string `json:"src_uid"`
// 	Remarks      string `json:"remarks"`
// }

// /*
// Pupose:This method is used to get the pan number for order input.
// Parameters:

// 	PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	AGMPA45767,nil

// 	==========
// 	*On Error
// 	==========
// 	"",error

// Author:Pavithra
// Date: 29AUG2023
// */
// func BlockBOClientFund(pJvDetailRec JvReqStruct, pRequest *http.Request, pFlag string) JvStatusStruct {
// 	log.Println("BlockClientFund (+)")

// 	// this variables is used to get Status
// 	var lJVReq techexcel.JvInputStruct
// 	var lReqDtl apigate.RequestorDetails
// 	var lJvStatusRec JvStatusStruct
// 	pJvProcessRec, lErr1 := JvConstructor(pJvDetailRec, pFlag)
// 	if lErr1 != nil {
// 		log.Println("SPOPJV01", lErr1.Error())
// 		lJvStatusRec.BoJvStatus = "E"
// 		return lJvStatusRec
// 	} else {

// 		log.Println("flag", pFlag)

// 		lReqDtl = apigate.GetRequestorDetail(pRequest)
// 		lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
// 		//commented by lakshmanan on 21 DEC 2023
// 		// lCocd := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["PaymentCode"])
// 		//commented by lakshmanan on 21 DEC 2023
// 		// lJVReq.COCD = lCocd
// 		//added by lakshmanan on 21 DEC 2023
// 		lJVReq.COCD, _ = sgbplaceorder.PaymentCode(pJvProcessRec.ClientId)
// 		lJVReq.VoucherDate = time.Now().Format("02/01/2006")
// 		lJVReq.BillNo = "SGB" + pJvProcessRec.ClientId
// 		lJVReq.SourceTable = "a_sgb_orderdetails"
// 		lJVReq.SourceTableKey = pJvProcessRec.OrderNo
// 		lJVReq.Amount = pJvProcessRec.JVamount
// 		lJVReq.WithGST = "N"

// 		// lConfigFile := common.ReadTomlConfig("toml/config.toml")
// 		lFTCaccount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["BseSGBorderAccount"])

// 		if pJvProcessRec.Transaction == "F" {
// 			lJVReq.AccountCode = lFTCaccount
// 			lJVReq.CounterAccount = pJvProcessRec.ClientId

// 			lJVReq.Narration = "Failed of SGB Order Refund  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
// 			lJvStatusRec.BoJvStatement = lJVReq.Narration
// 			lJvStatusRec.BoJvType = "C"
// 		} else if pJvProcessRec.Transaction == "C" {
// 			lJVReq.AccountCode = pJvProcessRec.ClientId
// 			lJVReq.CounterAccount = lFTCaccount
// 			lJVReq.Narration = "SGB Purchase  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
// 			lJvStatusRec.BoJvStatement = lJVReq.Narration
// 			lJvStatusRec.BoJvType = "D"
// 		}

// 		lJvStatusRec.BoJvStatus = "S"
// 		lJvStatusRec.BoJvAmount = pJvProcessRec.JVamount
// 		log.Println("JV Record:", lJVReq)
// 		//JV processing method
// 		//commented by pavithra --method name changed
// 		// lErr1 := clientfund.ProcessJV(lJVReq, lReqDtl)
// 		lErr1 := clientfund.BOProcessJV(lJVReq, lReqDtl)
// 		if lErr1 != nil {
// 			log.Println("SPOPJV02", lErr1.Error())
// 			lJvStatusRec.BoJvStatus = "E"
// 			return lJvStatusRec
// 		} else {
// 			lJvStatusRec.BoJvStatus = "S"
// 		}
// 	}
// 	log.Println("BlockClientFund (-)")
// 	return lJvStatusRec
// }

// /*
// Purpose: This function is used to construct the details for the BO Jv API

// 	Parameters:
// 		{
// 			JvAmount    string `json:"jvAmount"`
// 			JvStatus    string `json:"jvStatus"`
// 			JvStatement string `json:"jvStatement"`
// 			JvType      string `json:"jvType"`
// 		}, "FT000069"

// 	Response:

// 	On Success :
// 	===============

// 	===============

// 	On Error:
// 	===============
// 		{},error
// 	===============

// Author:Nithish kumar
// Date: 21Nov2023
// */
// func BlockFOClientFund(pRequest JvReqStruct, pStatusFlag string) (JvStatusStruct, error) {
// 	log.Println("BlockFOClientFund (+)")
// 	var lFOJvReqStruct FOJvStruct
// 	var lJvStatusRec JvStatusStruct
// 	var lUrl string

// 	config := common.ReadTomlConfig("./toml/config.toml")
// 	lFOJvReqStruct.UserId = fmt.Sprintf("%v", config.(map[string]interface{})["VerifyUser"])
// 	lFOJvReqStruct.SourceUserId = lFOJvReqStruct.UserId

// 	lToken, lErr1 := sgbplaceorder.GetFOToken()
// 	if lErr1 != nil {
// 		log.Println("SBFCF01", lErr1.Error())
// 		lJvStatusRec.FoJvStatus = common.ErrorCode
// 		return lJvStatusRec, lErr1
// 	} else {
// 		lPrice, lErr2 := strconv.Atoi(pRequest.Price)
// 		if lErr2 != nil {
// 			log.Println("SBFCF02", lErr2)
// 			lJvStatusRec.FoJvStatus = common.ErrorCode
// 			return lJvStatusRec, lErr2
// 		}

// 		lUnit, lErr3 := strconv.Atoi(pRequest.Unit)
// 		if lErr3 != nil {
// 			log.Println("SBFCF03", lErr3)
// 			lJvStatusRec.FoJvStatus = common.ErrorCode
// 			return lJvStatusRec, lErr3
// 		}
// 		// Choose the url for the api based on the mode Debit or Credit
// 		if pRequest.BoJvType == "D" {
// 			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayoutUrl"])
// 			lJvStatusRec.FoJvType = pRequest.BoJvType
// 			if pStatusFlag != "R" {
// 				lFOJvReqStruct.Remarks = "AMOUNT HOLD FOR SGB ORDER"
// 				lJvStatusRec.FoJvStatement = lFOJvReqStruct.Remarks
// 				lFOJvReqStruct.Amount = "-" + strconv.Itoa(lPrice*lUnit)
// 			} else {
// 				lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM SGB ORDER"
// 				lJvStatusRec.FoJvStatement = lFOJvReqStruct.Remarks
// 				lFOJvReqStruct.Amount = strconv.Itoa(lPrice * lUnit)
// 			}
// 		} else if pRequest.BoJvType == "C" {
// 			lUrl = fmt.Sprintf("%v", config.(map[string]interface{})["PayinUrl"])
// 			lJvStatusRec.FoJvType = pRequest.BoJvType
// 			lFOJvReqStruct.Remarks = "AMOUNT RELEASE FROM SGB ORDER"
// 			lFOJvReqStruct.Amount = strconv.Itoa(lPrice * lUnit)
// 		}
// 		lFOJvReqStruct.AccountId = pRequest.ClientId
// 		lJvStatusRec.FoJvAmount = lFOJvReqStruct.Amount

// 		lRequest, lErr4 := json.Marshal(lFOJvReqStruct)
// 		if lErr4 != nil {
// 			log.Println("SBFCF04", lErr4.Error())
// 			lJvStatusRec.FoJvStatus = common.ErrorCode
// 			return lJvStatusRec, lErr4
// 		} else {
// 			// construct the request body
// 			lBody := `jData=` + string(lRequest) + `&jKey=` + lToken
// 			lResp, lErr5 := clientfund.FOProcessJV(lUrl, lBody)
// 			if lErr5 != nil {
// 				log.Println("SBFCF05", lErr5.Error())
// 				lJvStatusRec.FoJvStatus = common.ErrorCode
// 				return lJvStatusRec, lErr5
// 			} else {
// 				if lResp.Status == "Ok" {
// 					log.Println("FO JV processed successfully for : ", lRequest)
// 					lJvStatusRec.FoJvStatus = common.SuccessCode
// 				} else {
// 					log.Println("FO JV processed Failed for : ", lRequest)
// 					lJvStatusRec.FoJvStatus = common.ErrorCode
// 				}
// 			}
// 		}
// 	}
// 	log.Println("BlockFOClientFund (-)")
// 	return lJvStatusRec, nil
// }

// /*
// Pupose:This method inserting the order head values in order header table.
// Parameters:

// 	pReqArr,pMasterId,PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========

// 	==========
// 	*On Error
// 	==========

// Author:Pavithra
// Date: 12JUNE2023
// */
// func ReverseProcess(pJvDetailRec JvReqStruct, r *http.Request, pExchange string, pBrokerId int, ErrorCase string, ReverseProcessError string) (JvStatusStruct, int, int) {
// 	log.Println("ReverseProcess (+)")
// 	var lReversedJvSuccess int
// 	var lReversedJvFailed int
// 	var lJvDatatoReturn JvStatusStruct

// 	Jvstatus := BlockBOClientFund(pJvDetailRec, r, "R")
// 	log.Println("Jvstatus", Jvstatus)
// 	lJvDatatoReturn.BoJvAmount = Jvstatus.BoJvAmount
// 	lJvDatatoReturn.BoJvStatement = Jvstatus.BoJvStatement
// 	lJvDatatoReturn.BoJvStatus = Jvstatus.BoJvStatus
// 	lJvDatatoReturn.BoJvType = Jvstatus.BoJvType
// 	//-----Mail To Accounts Team-------
// 	// REVERSEJS := "REVERSEJV"
// 	if Jvstatus.BoJvStatus == common.ErrorCode {

// 		lErr := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, ReverseProcessError)
// 		if lErr != nil {
// 			log.Println("updateOrderDetails ReverseProcess Error :", lErr)
// 			// return lJvDatatoReturn, lReversedJvSuccess, lReversedJvFailed
// 		} else {
// 			lReversedJvFailed++
// 			//-----------------need to send Accounts team mail
// 			//  Mail content Reverse process Failed Durin Back office Jv
// 			lFailedJvMail, lErr1 := ConstructMailforJvProcess(pJvDetailRec, "RPBOFAILED")
// 			if lErr1 != nil {
// 				log.Println("SGBPRP07", lErr1.Error())
// 			} else {
// 				lString := "JV Failed in SGB Order"
// 				lErr2 := emailUtil.SendEmail(lFailedJvMail, lString)
// 				if lErr2 != nil {
// 					log.Println("SGBPRP08", lErr2.Error())
// 				}
// 			}
// 		}
// 	} else {
// 		lReversedJvSuccess++
// 		//-----------------need to send mail to IT team
// 		lSuccessJvMail, lErr1 := ConstructMailforJvProcess(pJvDetailRec, "RPBOSUCCESS")
// 		if lErr1 != nil {
// 			log.Println("SGBPRP07", lErr1.Error())
// 		} else {
// 			lString := "JV Success in SGB Order"
// 			lErr2 := emailUtil.SendEmail(lSuccessJvMail, lString)
// 			if lErr2 != nil {
// 				log.Println("SGBPRP08", lErr2.Error())
// 			}
// 		}

// 	}

// 	log.Println("ReverseProcess (-)")
// 	return lJvDatatoReturn, lReversedJvSuccess, lReversedJvFailed
// }

// // func constructJVMail(pJvStatusRec JvStatusStruct, pJvDataStruct JvDataStruct) (emailUtil.EmailInput, error) {
// // 	log.Println("constructJVMail (+)")
// // 	var lEmailContent emailUtil.EmailInput
// // 	config := common.ReadTomlConfig("toml/emailconfig.toml")

// // 	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// // 	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// // 	lEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["AccountsTeamMail"])

// // 	lEmailContent.Subject = "Failed JV For SGB Orders"
// // 	html := "html/JvInvoice.html"

// // 	lTemp, lErr := template.ParseFiles(html)
// // 	if lErr != nil {
// // 		log.Println("SPOCJM01", lErr)
// // 		return lEmailContent, lErr
// // 	} else {
// // 		var lTpl bytes.Buffer
// // 		count := 1
// // 		var lDynamicEmailVal FailedJvStruct
// // var lDynamicEmailArr DynamicEmailStruct

// // 		lDynamicEmailVal.SINo = count
// // 		lDynamicEmailVal.ClientId = pJvDataStruct.ClientId
// // 		lDynamicEmailVal.Amount = pJvStatusRec.JvAmount
// // 		if pJvStatusRec.JvType == "D" {
// // 			lDynamicEmailVal.Transaction = "Debit"
// // 		} else if pJvStatusRec.JvType == "C" {
// // 			lDynamicEmailVal.Transaction = "Credit"
// // 		}
// // 		currentTime := time.Now()
// // 		// Format the current time in "11-09-2023" format
// // 		formattedDate := currentTime.Format("02-01-2006")
// // 		lDynamicEmailArr.Date = formattedDate
// // 		lDynamicEmailArr.FailedJvArr = append(lDynamicEmailArr.FailedJvArr, lDynamicEmailVal)
// // 		lTemp.Execute(&lTpl, lDynamicEmailArr)
// // 		lEmailbody := lTpl.String()

// // 		lEmailContent.Body = lEmailbody

// // 		log.Println("constructJVMail (-)")
// // 	}
// // 	return lEmailContent, nil
// // }

// /*

// Pupose:This method is used to get the email Input for success Sgb Place Order  .
// Parameters:

// 	pSgbClientDetails {},pStatus string

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	{
//      Name: clientName
//      Status: S
//      OrderDate: 8SEP2023
//      OrderNumber: 121345687
//      Symbol: sgb test
//      Unit: 5
//      Price: 500
//      Amount:2500
//      Activity : M
// 	},nil

// 	==========
// 	*On Error
// 	==========
// 	"",error

// Author:PRASHANTH
// Date: 8SEP2023
// */

// func constructSuccessmail(pSgbClientDetails JvReqStruct, pStatus string) (emailUtil.EmailInput, error) {
// 	log.Println("constructSuccessmail (+)")
// 	type dynamicEmailStruct struct {
// 		Name        string
// 		Status      string
// 		OrderDate   string
// 		Symbol      string
// 		OrderNumber string
// 		Unit        string
// 		Price       string
// 		Amount      string
// 		// Activity    string
// 	}

// 	var lEmailContent emailUtil.EmailInput
// 	config := common.ReadTomlConfig("toml/emailconfig.toml")

// 	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// 	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
// 	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 	// commented by pavithra
// 	// lEmailContent.Subject = "SGB Order"
// 	// html := "html/SgbOrderTemplate.html"

// 	lNovoConfig := common.ReadTomlConfig("toml/novoConfig.toml")
// 	lEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_ClientEmail_Subject"])
// 	html := fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_OrderSuccess_html"])

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
// 		// lDynamicEmailVal.Activity = pSgbClientDetails.ActionCode
// 		//   ================================================
// 		if pStatus == "S" {
// 			lDynamicEmailVal.Status = common.SUCCESS
// 		} else {
// 			lDynamicEmailVal.Status = common.FAILED
// 		}
// 		//   ================================================||

// 		// lEmailContent.ToEmailId = "prashanth.s@fcsonline.co.in"
// 		lEmailContent.ToEmailId = pSgbClientDetails.Mail
// 		// log.Println("lDynamicEmailVal Succes Mail", lDynamicEmailVal)
// 		lTemp.Execute(&lTpl, lDynamicEmailVal)
// 		lEmailbody := lTpl.String()

// 		lEmailContent.Body = lEmailbody
// 	}
// 	log.Println("constructSuccessmail (-)")
// 	return lEmailContent, nil
// }

// // ----------------------------------------------------------------
// // this method is used to get the Brokerid from database
// // ----------------------------------------------------------------
// func GetSgbBrokers(pExchange string) ([]int, error) {
// 	log.Println("GetSgbBrokers(+)")

// 	var lBrokerArr []int
// 	var lBrokerRec int
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("GSB01", lErr1)
// 		return lBrokerArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select md.BrokerId
// 						from a_ipo_memberdetails md
// 						where md.AllowedModules  like '%Sgb%'
// 						and md.Flag = 'Y'
// 						and md.OrderPreference = ?`
// 		lRows, lErr2 := lDb.Query(lCoreString, pExchange)
// 		if lErr2 != nil {
// 			log.Println("GSB02", lErr2)
// 			return lBrokerArr, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lBrokerRec)
// 				if lErr3 != nil {
// 					log.Println("GSB03", lErr3)
// 					return lBrokerArr, lErr3
// 				} else {
// 					lBrokerArr = append(lBrokerArr, lBrokerRec)
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetSgbBrokers(-)")
// 	return lBrokerArr, nil
// }

// /*
// Pupose:This method inserting the order head values in order header table.
// Parameters:

// 	pReqArr,pMasterId,PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========

// 	==========
// 	*On Error
// 	==========

// Author:Pavithra
// Date: 12JUNE2023
// */
// func UpdateHeader(pNseRespRec nsesgb.SgbAddResStruct, pJvRec JvReqStruct, pBseRespRec bsesgb.SgbRespStruct, pExchange string, pBrokerId int) error {
// 	log.Println("UpdateHeader (+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("PIH01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		lHeaderId, lErr2 := GetOrderId(pJvRec, pExchange)
// 		if lErr2 != nil {
// 			log.Println("SPGOI02", lErr2)
// 			return lErr2
// 		} else {
// 			if lHeaderId != 0 {
// 				if pExchange == common.BSE {
// 					lSqlString := `update a_sgb_orderheader h
// 							set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
// 							h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.Status = ?,
// 							h.StatusCode = ?,h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,
// 							h.UpdatedBy = ?,h.UpdatedDate = now()
// 							where h.ClientId = ?
// 							and h.brokerId = ?
// 							and h.Id = ?`

// 					_, lErr5 := lDb.Exec(lSqlString, pBseRespRec.ScripId, pBseRespRec.PanNo, pBseRespRec.InvestorCategory, pBseRespRec.ApplicantName, pBseRespRec.Depository, pBseRespRec.DpId, pBseRespRec.ClientBenfId, pBseRespRec.GuardianName, pBseRespRec.GuardianPanno, pBseRespRec.GuardianRelation, common.SUCCESS, pBseRespRec.StatusCode, pBseRespRec.StatusMessage, pBseRespRec.ErrorCode, pBseRespRec.ErrorMessage, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
// 					if lErr5 != nil {
// 						log.Println("PIH05", lErr5)
// 						return lErr5
// 					} else {
// 						// call InsertDetails method to inserting the order details in order details table
// 						lErr6 := UpdateDetail(pJvRec, pBseRespRec.Bids, lHeaderId)
// 						if lErr6 != nil {
// 							log.Println("PIH06", lErr6)
// 							return lErr6
// 						}
// 					}
// 				} else if pExchange == common.NSE {
// 					lSqlString := `update a_sgb_orderheader h
// 					set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
// 					h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.Status = ?,h.ClientReferenceNo = ?,
// 					h.StatusCode = ?,h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,
// 					h.UpdatedBy = ?,h.UpdatedDate = now()
// 					where h.ClientId = ?
// 					and h.brokerId = ?
// 					and h.Id = ?`

// 					_, lErr5 := lDb.Exec(lSqlString, pBseRespRec.ScripId, pBseRespRec.PanNo, pBseRespRec.InvestorCategory, pBseRespRec.ApplicantName, pBseRespRec.Depository, pBseRespRec.DpId, pBseRespRec.ClientBenfId, pBseRespRec.GuardianName, pBseRespRec.GuardianPanno, pBseRespRec.GuardianRelation, common.SUCCESS, pNseRespRec.ClientRefNumber, pBseRespRec.StatusCode, pBseRespRec.StatusMessage, pBseRespRec.ErrorCode, pBseRespRec.ErrorMessage, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
// 					if lErr5 != nil {
// 						log.Println("PIH05", lErr5)
// 						return lErr5
// 					} else {
// 						// call InsertDetails method to inserting the order details in order details table
// 						lErr6 := UpdateDetail(pJvRec, pBseRespRec.Bids, lHeaderId)
// 						if lErr6 != nil {
// 							log.Println("PIH05", lErr6)
// 							return lErr6
// 						} else {
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("UpdateHeader (-)")
// 	return nil
// }

// /*
// Pupose:This method inserting the order head values in order header table.
// Parameters:

// 	pReqArr,pMasterId,PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========

// 	==========
// 	*On Error
// 	==========

// Author:Pavithra
// Date: 12JUNE2023
// */
// func UpdateDetail(pJvRec JvReqStruct, pRespBidArr []bsesgb.RespSgbBidStruct, pHeaderId int) error {
// 	log.Println("UpdateDetail (+)")

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("PIH01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		for lIdx := 0; lIdx < len(pRespBidArr); lIdx++ {

// 			lSqlString := `update a_sgb_orderdetails d
// 							set d.RespApplicationNo = ?,d.RespOrderNo = ?,d.ActionCode = ?,d.RespSubscriptionunit = ?,d.RespRate = ?,
// 							d.ErrorCode = ?,d.Message = ?,d.BoJvStatus = ?,d.BoJvAmount = ?,d.BoJvStatement = ?,d.BoJvType = ?,
// 							d.FoJvStatus = ?,d.FoJvAmount = ?,d.FoJvStatement = ?,d.FoJvType = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
// 							where d.HeaderId = ?
// 							and d.ReqOrderNo = ?`

// 			_, lErr5 := lDb.Exec(lSqlString, pRespBidArr[lIdx].OrderNo, pRespBidArr[lIdx].BidId, pRespBidArr[lIdx].ActionCode, pRespBidArr[lIdx].SubscriptionUnit, pRespBidArr[lIdx].Rate, pRespBidArr[lIdx].ErrorCode, pRespBidArr[lIdx].Message, pJvRec.BoJvStatus, pJvRec.BoJvAmount, pJvRec.BoJvStatement, pJvRec.BoJvType, pJvRec.FoJvStatus, pJvRec.FoJvAmount, pJvRec.FoJvStatement, pJvRec.FoJvType, common.AUTOBOT, pHeaderId, pJvRec.OrderNo)
// 			if lErr5 != nil {
// 				log.Println("PIH05", lErr5)
// 				return lErr5
// 			} else {
// 				// call InsertDetails method to inserting the order details in order details table

// 			}
// 		}
// 	}
// 	log.Println("UpdateDetail (-)")
// 	return nil
// }

// func GetOrderId(pJvRec JvReqStruct, pExchange string) (int, error) {
// 	log.Println("GetOrderId (+)")

// 	var lHeaderId int

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SPGOI01", lErr1)
// 		return lHeaderId, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select (case when count(1) > 0 then h.Id  else 0 end) Id
// 						from a_sgb_orderheader h,a_sgb_orderdetails d
// 						where h.Id = d.HeaderId
// 						and h.ClientId = ?
// 						and d.ReqOrderNo = ?
// 						and h.Exchange  = ?
// 						and h.MasterId = (
// 							select m.id
// 							from a_sgb_master m
// 							where m.Symbol = ?
// 						)`
// 		lRows, lErr2 := lDb.Query(lCoreString, pJvRec.ClientId, pJvRec.OrderNo, pExchange, pJvRec.Symbol)
// 		if lErr2 != nil {
// 			log.Println("SPGOI02", lErr2)
// 			return lHeaderId, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lHeaderId)
// 				if lErr3 != nil {
// 					log.Println("SPGOI03", lErr3)
// 					return lHeaderId, lErr3
// 				}
// 				log.Println(lHeaderId)
// 			}
// 		}
// 	}
// 	log.Println("GetOrderId (-)")
// 	return lHeaderId, nil
// }

// // ---------------------------------------------------------------------------------
// // this method is used to update the reocrds od NSE in db
// // ---------------------------------------------------------------------------------
// func UpdateNseRecord(pRespRec nsesgb.SgbAddResStruct, pJvRec JvReqStruct, pBseConvertedRec bsesgb.SgbRespStruct, pBrokerId int) error {
// 	log.Println("UpdateNseRecord (+)")
// 	// for lJvReqIdx := 0; lJvReqIdx < len(pJvRecArr); lJvReqIdx++ {
// 	// if strconv.Itoa(pRespRec.OrderNumber) == pJvRecArr[lJvReqIdx].OrderNo {
// 	lErr1 := UpdateHeader(pRespRec, pJvRec, pBseConvertedRec, common.NSE, pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("SGBOFUNR01", lErr1)
// 		return lErr1
// 	} else {
// 		log.Println("Nse Record Updated Successfully")
// 	}
// 	// }
// 	// }
// 	log.Println("UpdateNseRecord (-)")
// 	return nil
// }

// //---------------------------------------------------------------------------------
// // this method is used to update the reocrds od BSE in db
// //---------------------------------------------------------------------------------

// func UpdateBseRecord(pRespRec bsesgb.SgbRespStruct, pJvRec JvReqStruct, pBrokerId int) error {
// 	log.Println("UpdateBseRecord (+)")
// 	// for lJvReqIdx := 0; lJvReqIdx < len(pJvRecArr); lJvReqIdx++ {
// 	for lBidIdx := 0; lBidIdx < len(pRespRec.Bids); lBidIdx++ {
// 		if pRespRec.Bids[lBidIdx].OrderNo == pJvRec.OrderNo {
// 			var lEmptyRec nsesgb.SgbAddResStruct
// 			lErr1 := UpdateHeader(lEmptyRec, pJvRec, pRespRec, common.BSE, pBrokerId)
// 			if lErr1 != nil {
// 				log.Println("SGBOFUBR01", lErr1)
// 				return lErr1
// 			} else {
// 				log.Println("Bse Record Updated Successfully")
// 			}
// 		}
// 	}
// 	// }
// 	log.Println("UpdateBseRecord (-)")
// 	return nil
// }

// // ---------------------------------------------------------------------------------
// // this method is used to filter the Jv posted records
// // ---------------------------------------------------------------------------------
// func PostJvForOrder(pSgbReqRec bsesgb.SgbReqStruct, pJvDetailRec JvReqStruct, r *http.Request, pValidSgb SgbStruct, pBrokerId int) (int, int, int, int, int, int, int, int, int, int, error) {
// 	log.Println("PostJvForOrder (+)")

// 	config := common.ReadTomlConfig("./toml/SgbConfig.toml")
// 	// NonProcess := fmt.Sprintf("%v", config.(map[string]interface{})["NonProcess"])
// 	// Processed := fmt.Sprintf("%v", config.(map[string]interface{})["Processed"])
// 	ErrorCase := fmt.Sprintf("%v", config.(map[string]interface{})["ErrorCase"])
// 	// Success := fmt.Sprintf("%v", config.(map[string]interface{})["Success"])
// 	VerifyError := fmt.Sprintf("%v", config.(map[string]interface{})["VerifyError"])
// 	InsufficientError := fmt.Sprintf("%v", config.(map[string]interface{})["InsufficientError"])
// 	BlockClientFundError := fmt.Sprintf("%v", config.(map[string]interface{})["BlockClientFundError"])
// 	// ReverseProcessError := fmt.Sprintf("%v", config.(map[string]interface{})["ReverseProcessError"])
// 	// FrontOfficeError := fmt.Sprintf("%v", config.(map[string]interface{})["FrontOfficeError"])
// 	// ExchangeError := fmt.Sprintf("%v", config.(map[string]interface{})["ExchangeError"])

// 	var lJvsuccess int
// 	var lJvfailed int
// 	var lExchangesuccess int
// 	var lExchangeFailed int
// 	var lReversedJvSuccess int
// 	var lReversedJvFailed int
// 	var lFoJvSuccess int
// 	var lFoJvFailed int
// 	var lFundVerifySuccess int
// 	var lFundVerifyFailed int

// 	lVerify := "VERIFY"
// 	lBlockClient := "BLOCKCLIENT"
// 	lInsufficient := "INSUFFICIENT"

// 	for lBidIdx := 0; lBidIdx < len(pSgbReqRec.Bids); lBidIdx++ {

// 		lSubscriptionUnit, lErr1 := strconv.Atoi(pSgbReqRec.Bids[lBidIdx].SubscriptionUnit)
// 		if lErr1 != nil {
// 			log.Println("SPJO01", lErr1)
// 			return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr1
// 		}
// 		lRate, lErr2 := common.ConvertStringToFloat(pSgbReqRec.Bids[lBidIdx].Rate)
// 		if lErr2 != nil {
// 			log.Println("SPJO02", lErr2)
// 			return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr2
// 		}

// 		pJvDetailRec.Amount = strconv.Itoa(int(lRate) * lSubscriptionUnit)
// 		// ===================================================Converting string To Float and int =======================================
// 		lClientFund, lErr3 := sgbplaceorder.VerifyFOFundDetails(pJvDetailRec.ClientId)
// 		log.Println("lClientFund", lClientFund)
// 		if lErr3 != nil {
// 			// ====================================while verifying the fund getting error================================
// 			lFundVerifyFailed++
// 			log.Println("SPJO03", lErr3)
// 			lErr4 := updateOrderDetails(pJvDetailRec, pValidSgb.Exchange, pBrokerId, ErrorCase, VerifyError)
// 			if lErr4 != nil {
// 				log.Println("updateOrderDetails Error SPJO04:", lErr4)
// 				// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr3
// 			} else {
// 				//immediate mail to IT team
// 				// Error in Verifiying Accounts
// 				lVerifyClient, lErr5 := ConstructMailforJvProcess(pJvDetailRec, lVerify)
// 				if lErr5 != nil {
// 					log.Println("SPJO05", lErr5)
// 					// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr5
// 				} else {
// 					lErr6 := emailUtil.SendEmail(lVerifyClient, "Error while verifying Client Acc Balance details")
// 					if lErr6 != nil {
// 						log.Println("SPJO06", lErr6)
// 						return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr6

// 					}
// 				}
// 			}
// 			return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr3
// 			// ====================================while verifying the fund getting error================================
// 		} else {
// 			//----------------------------------------------------------------
// 			//Demat Account Balance Verification Success
// 			//----------------------------------------------------------------
// 			if lClientFund.Status == common.SuccessCode {
// 				lFundVerifySuccess++

// 				if lClientFund.AccountBalance < float64(int(lRate)*lSubscriptionUnit) {

// 					// ====================================== INSUFFICIENT ACC BALANCE START =====================================================
// 					if lClientFund.AccountBalance < float64(lRate) {
// 						//  Insufficient Client Account Balance
// 						pJvDetailRec.BoJvAmount = "insufficient"
// 						lErr7 := updateOrderDetails(pJvDetailRec, pValidSgb.Exchange, pBrokerId, ErrorCase, InsufficientError)
// 						if lErr7 != nil {
// 							log.Println("updateOrderDetails Error SPJO07:", lErr7)
// 							return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr7
// 						} else {
// 							lInsufficientClient, lErr8 := ConstructMailforJvProcess(pJvDetailRec, lInsufficient)
// 							if lErr8 != nil {
// 								log.Println("SPJO08", lErr8)
// 								// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr8
// 							} else {
// 								lErr9 := emailUtil.SendEmail(lInsufficientClient, "Insufficient Fund in Client Account for SGB")
// 								if lErr9 != nil {
// 									log.Println("SPJO09", lErr9.Error())
// 									return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr9
// 								}
// 							}
// 						}
// 						// ==========================================INSUFFICIENT ACC BALANCE END===========================================
// 					} else {
// 						// ===============================INSUFFICIENT AMOUNT CHANGING UNIT================================

// 						MaximumUnits := lClientFund.AccountBalance / float64(lRate)
// 						// Convert float to int using Ceil (round up)
// 						MaximumUnitsIntFloor := int(math.Floor(MaximumUnits))
// 						pJvDetailRec.Amount = strconv.Itoa(int(lRate) * MaximumUnitsIntFloor)

// 						pJvDetailRec.Unit = strconv.Itoa(MaximumUnitsIntFloor)
// 						// Convert the integer to string
// 						for i := 0; i < len(pSgbReqRec.Bids); i++ {
// 							pSgbReqRec.Bids[i].SubscriptionUnit = pJvDetailRec.Unit
// 						}

// 						Jvstatus := BlockBOClientFund(pJvDetailRec, r, "C")
// 						log.Println("Jvstatus", Jvstatus)
// 						pJvDetailRec.BoJvAmount = Jvstatus.BoJvAmount
// 						pJvDetailRec.BoJvStatement = Jvstatus.BoJvStatement
// 						pJvDetailRec.BoJvStatus = Jvstatus.BoJvStatus
// 						pJvDetailRec.BoJvType = Jvstatus.BoJvType
// 						//----------------------------------------------------------------
// 						//JV Process while success
// 						//----------------------------------------------------------------
// 						if Jvstatus.BoJvStatus == common.SuccessCode {
// 							lJvsuccess++
// 							//----------------------------------------------------------------
// 							//exchange calling place
// 							//----------------------------------------------------------------
// 							lExchSuccess, lExchFailed, lBoRevJvSuccess, lBoRevJvFailed, lFoSuccess, lFoFailed, lErr10 := ExchangeProcess(pSgbReqRec, pJvDetailRec, pValidSgb.Exchange, pBrokerId, r)
// 							if lErr10 != nil {
// 								log.Println("SPJO10", lErr10.Error())
// 								lExchangesuccess += lExchSuccess
// 								lExchangeFailed += lExchFailed
// 								lReversedJvSuccess += lBoRevJvSuccess
// 								lReversedJvFailed += lBoRevJvFailed
// 								lFoJvSuccess += lFoSuccess
// 								lFoJvFailed += lFoFailed
// 								return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr10
// 							} else {
// 								lExchangesuccess += lExchSuccess
// 								lExchangeFailed += lExchFailed
// 								lReversedJvSuccess += lBoRevJvSuccess
// 								lReversedJvFailed += lBoRevJvFailed
// 								lFoJvSuccess += lFoSuccess
// 								lFoJvFailed += lFoFailed
// 							}
// 						} else if Jvstatus.BoJvStatus == common.ErrorCode {
// 							// =========================================================== BACK OFFICE FAILED ===================================

// 							lErr13 := updateOrderDetails(pJvDetailRec, pValidSgb.Exchange, pBrokerId, ErrorCase, BlockClientFundError)
// 							if lErr13 != nil {
// 								log.Println("updateOrderDetails Error SPJO13:", lErr13)
// 								// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr13
// 							}
// 							//immediate mail to IT Team jv error but the client account has fund to process jv
// 							// JV error
// 							lJvfailed++
// 							lBlockFund, lErr14 := ConstructMailforJvProcess(pJvDetailRec, lBlockClient)
// 							if lErr14 != nil {
// 								log.Println("SPJO14", lErr14)
// 								// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed,lFoJvSuccess,lFoJvFailed, lErr14
// 							} else {
// 								lErr15 := emailUtil.SendEmail(lBlockFund, "Unable to Deduct Amount BO From Client Account")
// 								if lErr15 != nil {
// 									log.Println("SPJO15", lErr15.Error())
// 									return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr15
// 								}
// 							}
// 							return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr13
// 							// ================================================== BACK OFFICE FAILED =======================================
// 						}
// 					}
// 					// ===============================INSUFFICIENT AMOUNT CHANGING UNIT================================
// 				} else {
// 					// =================================BACK OFFICE JV SUFFICIENT AMOUNT================================
// 					Jvstatus := BlockBOClientFund(pJvDetailRec, r, "C")
// 					log.Println("Jvstatus", Jvstatus)
// 					pJvDetailRec.BoJvAmount = Jvstatus.BoJvAmount
// 					pJvDetailRec.BoJvStatement = Jvstatus.BoJvStatement
// 					pJvDetailRec.BoJvStatus = Jvstatus.BoJvStatus
// 					pJvDetailRec.BoJvType = Jvstatus.BoJvType
// 					//----------------------------------------------------------------
// 					//JV Process while success
// 					//----------------------------------------------------------------
// 					if Jvstatus.BoJvStatus == common.SuccessCode {
// 						lJvsuccess++
// 						//----------------------------------------------------------------
// 						//exchange call
// 						//----------------------------------------------------------------
// 						lExchSuccess, lExchFailed, lBoRevJvSuccess, lBoRevJvFailed, lFoSuccess, lFoFailed, lErr16 := ExchangeProcess(pSgbReqRec, pJvDetailRec, pValidSgb.Exchange, pBrokerId, r)
// 						if lErr16 != nil {
// 							log.Println("SPJO16", lErr16.Error())
// 							lExchangesuccess += lExchSuccess
// 							lExchangeFailed += lExchFailed
// 							lReversedJvSuccess += lBoRevJvSuccess
// 							lReversedJvFailed += lBoRevJvFailed
// 							lFoJvSuccess += lFoSuccess
// 							lFoJvFailed += lFoFailed
// 							return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr16
// 						}
// 						lExchangesuccess += lExchSuccess
// 						lExchangeFailed += lExchFailed
// 						lReversedJvSuccess += lBoRevJvSuccess
// 						lReversedJvFailed += lBoRevJvFailed
// 						lFoJvSuccess += lFoSuccess
// 						lFoJvFailed += lFoFailed

// 					} else if Jvstatus.BoJvStatus == common.ErrorCode {
// 						// =========================================== BACK OFFICE FAILED ====================================
// 						lJvfailed++
// 						lErr17 := updateOrderDetails(pJvDetailRec, pValidSgb.Exchange, pBrokerId, ErrorCase, BlockClientFundError)
// 						if lErr17 != nil {
// 							log.Println("updateOrderDetails Error SPJO17 :", lErr17)
// 							// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr17
// 						}
// 						//immediate mail to IT Team
// 						//  Error in Blocking Amount From Client Account JV Error
// 						lBlockFund, lErr18 := ConstructMailforJvProcess(pJvDetailRec, lBlockClient)
// 						if lErr18 != nil {
// 							log.Println("SPJO18", lErr18)
// 							// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr18
// 						}
// 						lErr19 := emailUtil.SendEmail(lBlockFund, "Unable to Deduct Amount BO From Client Account")
// 						if lErr19 != nil {
// 							log.Println("SPJO19", lErr19.Error())
// 							return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr19
// 						}
// 						// =========================================== BACK OFFICE FAILED ====================================
// 					}
// 				}
// 				// =================================BACK OFFICE JV SUFFICIENT AMOUNT==========================================
// 			} else {
// 				// ==================================FRONT OFFICE VERIFICATION ERROR===========================================
// 				lFundVerifyFailed++
// 				lErr20 := updateOrderDetails(pJvDetailRec, pValidSgb.Exchange, pBrokerId, ErrorCase, VerifyError)
// 				if lErr20 != nil {
// 					log.Println("updateOrderDetails Error SPJO20 :", lErr20)
// 					// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr20
// 				} else {
// 					//----------------------------------------------------------------
// 					//Demat Account Balance Verification Error
// 					//----------------------------------------------------------------
// 					//immediate mail to IT team
// 					// Error in Verifiying Account Balance
// 					lVerifyClient, lErr21 := ConstructMailforJvProcess(pJvDetailRec, lVerify)
// 					if lErr21 != nil {
// 						log.Println("SPJO21", lErr21)
// 						// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr21
// 					}
// 					lErr22 := emailUtil.SendEmail(lVerifyClient, "FO Verifying Client Account Details")
// 					if lErr22 != nil {
// 						log.Println("SPJO22", lErr22.Error())
// 						return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, lErr22
// 					}
// 				}
// 				// ================================FRONT OFFICE VERIFICATION ERROR===========================================
// 			}
// 		}
// 	}
// 	log.Println("PostJvForOrder (-)")
// 	return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvFailed, lFundVerifySuccess, lFundVerifyFailed, nil
// }

// func ExchangeProcess(pSgbReqRec bsesgb.SgbReqStruct, pJvDetailRec JvReqStruct, pExchange string, pBrokerId int, r *http.Request) (int, int, int, int, int, int, error) {
// 	log.Println("ExchangeProcess (+)")

// 	var lExchangesuccess int
// 	var lExchangeFailed int
// 	var lReversedJvSuccess int
// 	var lReversedJvFailed int
// 	var lFoJvSuccess int
// 	var lFoJvError int

// 	config := common.ReadTomlConfig("./toml/SgbConfig.toml")
// 	Success := fmt.Sprintf("%v", config.(map[string]interface{})["Success"])
// 	Processed := fmt.Sprintf("%v", config.(map[string]interface{})["Processed"])

// 	ErrorCase := fmt.Sprintf("%v", config.(map[string]interface{})["ErrorCase"])
// 	ExchangeError := fmt.Sprintf("%v", config.(map[string]interface{})["ExchangeError"])
// 	ReverseProcessError := fmt.Sprintf("%v", config.(map[string]interface{})["ReverseProcessError"])
// 	FrontOfficeError := fmt.Sprintf("%v", config.(map[string]interface{})["FrontOfficeError"])

// 	if pExchange == common.BSE {
// 		lRespRec, lErr1 := exchangecall.ApplyBseSgb(pSgbReqRec, common.AUTOBOT, pBrokerId)
// 		if lErr1 != nil {
// 			lExchangeFailed++
// 			log.Println("SSEP01", lErr1)
// 			lErr2 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, ExchangeError)
// 			if lErr2 != nil {
// 				log.Println("updateOrderDetails Error SSEP02 :", lErr2)
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr2
// 			}
// 			return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr1
// 		} else {
// 			// ===============================EXCHANGE FAILED ================================================
// 			if lRespRec.ErrorCode != "0" {
// 				lExchangeFailed++
// 				lErr3 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, ExchangeError)
// 				if lErr3 != nil {
// 					log.Println("updateOrderDetails Error SSEP03 :", lErr3)
// 					// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr3
// 				}
// 				lErr3 = UpdateExchFailed(pJvDetailRec, pExchange, pBrokerId)
// 				if lErr3 != nil {
// 					log.Println("SSEP26 :", lErr3)
// 					// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr3
// 				}

// 				BoJvStatus, lsuccessJv, lFailedJv := ReverseProcess(pJvDetailRec, r, pExchange, pBrokerId, ErrorCase, ReverseProcessError)
// 				pJvDetailRec.BoJvAmount = BoJvStatus.BoJvAmount
// 				pJvDetailRec.BoJvStatement = BoJvStatus.BoJvStatement
// 				pJvDetailRec.BoJvType = BoJvStatus.BoJvType
// 				pJvDetailRec.BoJvStatus = BoJvStatus.BoJvStatus
// 				lReversedJvFailed += lFailedJv
// 				lReversedJvSuccess += lsuccessJv
// 				return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil
// 				// ===========================================================================================
// 			} else {
// 				lExchangesuccess++
// 				lFoJvRec, lErr4 := BlockFOClientFund(pJvDetailRec, "C")
// 				pJvDetailRec.FoJvAmount = lFoJvRec.FoJvAmount
// 				pJvDetailRec.FoJvStatement = lFoJvRec.FoJvStatement
// 				pJvDetailRec.FoJvType = lFoJvRec.FoJvType
// 				pJvDetailRec.FoJvStatus = lFoJvRec.FoJvStatus
// 				if lErr4 != nil {
// 					lFoJvError++
// 					log.Println("SSEP04", lErr4)
// 					lErr5 := FoJvprocessMail(pJvDetailRec)
// 					if lErr5 != nil {
// 						log.Println("FO JV error SSEP05", lErr5)
// 					}
// 					lErr6 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, FrontOfficeError)
// 					if lErr6 != nil {
// 						log.Println("updateOrderDetails Error SSEP06 :", lErr6)
// 						// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr6
// 					}
// 					return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr4
// 				} else {
// 					// ==================================FRONT OFFICE ============================================
// 					if pJvDetailRec.FoJvStatus == common.ErrorCode {
// 						lFoJvError++
// 						lErr7 := FoJvprocessMail(pJvDetailRec)
// 						if lErr7 != nil {
// 							log.Println("FO JV error SSEP07 ", lErr7)
// 						}
// 						lErr8 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, FrontOfficeError)
// 						if lErr8 != nil {
// 							log.Println("updateOrderDetails Error SSEP08 :", lErr8)
// 							// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr8
// 						}
// 						return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil
// 						// =============================================================================================

// 					} else if pJvDetailRec.FoJvStatus == common.SuccessCode {
// 						lFoJvSuccess++

// 						lErr9 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, Processed, Success)
// 						if lErr9 != nil {
// 							log.Println("updateOrderDetails Error SSEP09 :", lErr9)
// 							// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr9
// 						}
// 						return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil

// 					}
// 				}
// 				lErr10 := UpdateBseRecord(lRespRec, pJvDetailRec, pBrokerId)
// 				if lErr10 != nil {
// 					log.Println("SSEP10", lErr10)
// 					// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr10
// 				}
// 				lErr11 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, Processed, Success)
// 				if lErr11 != nil {
// 					log.Println("updateOrderDetails Error SSEP11 :", lErr11)
// 					// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr11
// 				}
// 				SuccessMail, lErr12 := constructSuccessmail(pJvDetailRec, "S")
// 				if lErr12 != nil {
// 					log.Println("SSEP12", lErr12)
// 					// return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed,lFoJvSuccess,lFoJvError, lErr12
// 				} else {
// 					lErr13 := emailUtil.SendEmail(SuccessMail, "Order Placed Succesfully on SGB")
// 					if lErr13 != nil {
// 						log.Println("SSEP13", lErr13.Error())
// 						// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr13
// 					}
// 				}
// 			}
// 		}
// 	} else if pExchange == common.NSE {
// 		lNseReqRec := exchangecall.SgbReqConstruct(pSgbReqRec)
// 		lRespRec, lErr14 := exchangecall.ApplyNseSgb(lNseReqRec, common.AUTOBOT, pBrokerId)
// 		// if lErr14 != nil {
// 		// 	log.Println("SSEP14", lErr14)
// 		// 	lExchangeFailed++
// 		// 	return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr14
// 		// } else {
// 		if lErr14 != nil || lRespRec.Status != "success" {
// 			log.Println("SSEP14", lErr14)
// 			lExchangeFailed++

// 			// ===================================================EXCHANGE FAILED ====================================
// 			lErr15 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, ExchangeError)
// 			if lErr15 != nil {
// 				log.Println("updateOrderDetails Error SSEP15 :", lErr15)
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr15
// 			}
// 			lErr15 = UpdateExchFailed(pJvDetailRec, pExchange, pBrokerId)
// 			if lErr15 != nil {
// 				log.Println("SSEP27 :", lErr15)
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr15
// 			}
// 			BoJvRec, lsuccessJv, lFailedJv := ReverseProcess(pJvDetailRec, r, pExchange, pBrokerId, ErrorCase, ReverseProcessError)
// 			pJvDetailRec.BoJvAmount = BoJvRec.BoJvAmount
// 			pJvDetailRec.BoJvStatement = BoJvRec.BoJvStatement
// 			pJvDetailRec.BoJvType = BoJvRec.BoJvType
// 			pJvDetailRec.BoJvStatus = BoJvRec.BoJvStatus
// 			lReversedJvFailed += lFailedJv
// 			lReversedJvSuccess += lsuccessJv
// 			return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil
// 			// ============================================EXCHANGE FAILED====================================
// 		} else {
// 			lExchangesuccess++
// 			lFoJvRec, lErr16 := BlockFOClientFund(pJvDetailRec, "C")
// 			pJvDetailRec.FoJvAmount = lFoJvRec.FoJvAmount
// 			pJvDetailRec.FoJvStatement = lFoJvRec.FoJvStatement
// 			pJvDetailRec.FoJvType = lFoJvRec.FoJvType
// 			pJvDetailRec.FoJvStatus = lFoJvRec.FoJvStatus
// 			if lErr16 != nil {
// 				lFoJvError++
// 				log.Println("SSEP16", lErr16)
// 				lErr17 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, FrontOfficeError)
// 				if lErr17 != nil {
// 					log.Println("updateOrderDetails Error SSEP17 :", lErr17)
// 					// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr17
// 				}
// 				lErr18 := FoJvprocessMail(pJvDetailRec)
// 				if lErr18 != nil {
// 					log.Println("FO JV error SSEP18", lErr18)
// 					// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr18
// 				}
// 				return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr16
// 			} else {
// 				// ============================================FRONT OFFICE  ==================================
// 				if pJvDetailRec.FoJvStatus == common.SuccessCode {
// 					lFoJvSuccess++
// 					lErr19 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, Processed, Success)
// 					if lErr19 != nil {
// 						log.Println("updateOrderDetails Error SSEP19 :", lErr19)
// 						// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr19
// 					}
// 				} else if pJvDetailRec.FoJvStatus == common.ErrorCode {
// 					lFoJvError++
// 					lErr20 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, ErrorCase, FrontOfficeError)
// 					if lErr20 != nil {
// 						log.Println("updateOrderDetails Error SSEP20 :", lErr20)
// 						return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr20
// 					}
// 					lErr21 := FoJvprocessMail(pJvDetailRec)
// 					if lErr21 != nil {
// 						log.Println("FO JV error SSEP21", lErr21)
// 						return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr20
// 					}
// 					return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil
// 					//=============================================================================================
// 				}
// 			}
// 			lErr22 := updateOrderDetails(pJvDetailRec, pExchange, pBrokerId, Processed, Success)
// 			if lErr22 != nil {
// 				log.Println("updateOrderDetails Error SSEP22 :", lErr22)
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr22
// 			}
// 			lBseConvertedRec := exchangecall.SgbRespConstruct(lRespRec)
// 			// Appplicant Name Was Not Come as Return
// 			lBseConvertedRec.ApplicantName = pSgbReqRec.ApplicantName
// 			lBseConvertedRec.InvestorCategory = pSgbReqRec.InvestorCategory
// 			lErr23 := UpdateNseRecord(lRespRec, pJvDetailRec, lBseConvertedRec, pBrokerId)
// 			if lErr23 != nil {
// 				log.Println("SSEP23", lErr23)
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr23
// 			}
// 			SuccessMail, lErr24 := constructSuccessmail(pJvDetailRec, common.SuccessCode)
// 			if lErr24 != nil {
// 				log.Println("SSEP24", lErr24)
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr24
// 			}
// 			lErr25 := emailUtil.SendEmail(SuccessMail, "Order Placed Succesfully on SGB")
// 			if lErr25 != nil {
// 				log.Println("SSEP25", lErr25.Error())
// 				// return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, lErr25
// 			}
// 			return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil
// 		}
// 		//	}
// 	}
// 	log.Println("ExchangeProcess (-)")
// 	return lExchangesuccess, lExchangeFailed, lReversedJvSuccess, lReversedJvFailed, lFoJvSuccess, lFoJvError, nil
// }

// func FoJvprocessMail(pJvDetailRec JvReqStruct) error {
// 	log.Println("FoJvprocessMail (+)")
// 	lFailedJvMail, lErr1 := ConstructMailforJvProcess(pJvDetailRec, "RPFOFAILED")
// 	if lErr1 != nil {
// 		log.Println("SGBPRP07", lErr1.Error())
// 		return lErr1
// 	} else {
// 		lString := "JV Failed in SGB Order"
// 		lErr2 := emailUtil.SendEmail(lFailedJvMail, lString)
// 		if lErr2 != nil {
// 			log.Println("SGBPRP08", lErr2.Error())
// 			return lErr2
// 		}
// 	}
// 	log.Println("FoJvprocessMail (-)")
// 	return nil
// }
// func ConstructMailforJvProcess(pJV JvReqStruct, pStatus string) (emailUtil.EmailInput, error) {
// 	log.Println("ConstructMailforJvProcess (+)")

// 	type dynamicEmailStruct struct {
// 		Date     string `json:"date"`
// 		ClientId string `json:"clientId"`
// 		OrderNo  string `json:"orderNo"`
// 		JvUnit   string `json:"jvUnit"`
// 		JvAmount string `json:"jvAmount"`
// 	}
// 	var lJVEmailContent emailUtil.EmailInput
// 	config := common.ReadTomlConfig("toml/emailconfig.toml")
// 	lJVEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// 	lJVEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
// 	lJVEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 	// Mail Id for It support is Not Added in Toml
// 	lJVEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["ToEmailId"])

// 	lNovoConfig := common.ReadTomlConfig("toml/novoConfig.toml")
// 	lJVEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_ClientEmail_Subject"])

// 	var html string
// 	if pStatus == "VERIFY" {
// 		// commented by pavithra
// 		// html = "html/VerifyClientFund.html"
// 		//commented by pavithra --
// 		// html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_VerifyFund_html"])
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_VerifyFOFund_html"])
// 	} else if pStatus == "BLOCKCLIENT" {
// 		// commented by pavithra
// 		// html = "html/BlockClientMail.html"
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_BlockClientFund_html"])
// 	} else if pStatus == "INSUFFICIENT" {
// 		// commented by pavithra
// 		// html = "html/InsufficientAmount.html"
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_InsufficientAmt_html"])
// 	} else if pStatus == "REVERSEJV" {
// 		// commented by pavithra
// 		// html = "html/Exchagestatus.html"
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_SuccessOrderMail_html"])
// 	} else if pStatus == "RPBOFAILED" {
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_RPBO_Failed_Mail_html"])
// 	} else if pStatus == "RPBOSUCCESS" {
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_RPBO_SUCCESS_Mail_html"])
// 	} else if pStatus == "RPFOFAILED" {
// 		html = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["SGB_RPFO_Failed_Mail_html"])
// 	}

// 	currentTime := time.Now()
// 	currentDate := currentTime.Format("02-01-2006")

// 	lTemp, lErr := template.ParseFiles(html)
// 	if lErr != nil {
// 		log.Println("SSCFJ01", lErr)
// 		return lJVEmailContent, lErr
// 	} else {
// 		var lTpl bytes.Buffer
// 		var lDynamicEmailVal dynamicEmailStruct
// 		if pStatus == "VERIFY" {
// 			//  IT Dept To verify client Account
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 		} else if pStatus == "BLOCKCLIENT" {
// 			//  IT Dept To Deducting Amount from client Account
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 			lDynamicEmailVal.JvUnit = pJV.Unit
// 			lDynamicEmailVal.JvAmount = pJV.Amount
// 		} else if pStatus == "INSUFFICIENT" {
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 			lDynamicEmailVal.JvUnit = pJV.Unit
// 			lDynamicEmailVal.JvAmount = pJV.Amount
// 			lJVEmailContent.ToEmailId = pJV.Mail
// 		} else if pStatus == "REVERSEJV" {
// 			// IT Dept failed During Exchange processs
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 			lDynamicEmailVal.JvUnit = pJV.Unit
// 			lDynamicEmailVal.JvAmount = pJV.Amount
// 		} else if pStatus == "RPBOFAILED" {
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 			lDynamicEmailVal.JvAmount = pJV.Amount
// 			lJVEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["AccountsTeamMail"])
// 		} else if pStatus == "RPBOSUCCESS" {
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 			lDynamicEmailVal.JvAmount = pJV.Amount
// 		} else if pStatus == "RPFOFAILED" {
// 			lDynamicEmailVal.Date = currentDate
// 			lDynamicEmailVal.ClientId = pJV.ClientId
// 			lDynamicEmailVal.OrderNo = pJV.OrderNo
// 			lDynamicEmailVal.JvAmount = pJV.Amount
// 			lJVEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["AccountsTeamMail"])
// 		}

// 		// lJVEmailContent.ToEmailId = "prashanth.s@fcsonline.co.in"
// 		// lJVEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 		lTemp.Execute(&lTpl, lDynamicEmailVal)
// 		lEmailbody := lTpl.String()

// 		lJVEmailContent.Body = lEmailbody
// 	}
// 	log.Println("ConstructMailforJvProcess (-)")
// 	return lJVEmailContent, nil
// }

// //  This Method is used to update Nse Sgb Pending Details
// // author prashanth

// // func updateNseSgbPending(pRespSgbdata nsesgb.SgbAddResStruct, pReqSgbdata nsesgb.SgbAddReqStruct, pUser string) error {
// // 	log.Println("updatePendingSgbStatus (+)")
// // 	lStatus := common.SUCCESS
// // 	lErr1 := updateNseSgbHeader(pRespSgbdata, pUser, lStatus)
// // 	if lErr1 != nil {
// // 		log.Println("SSUNSP01", lErr1)
// // 		return lErr1
// // 	} else {
// // 		lErr2 := updateNseSgbDetails(pRespSgbdata, pUser)
// // 		if lErr2 != nil {
// // 			log.Println("SSUNSP02", lErr2)
// // 			return lErr2
// // 		} else {
// // 			lErr3 := updateNseBidTracking(pRespSgbdata, pReqSgbdata, pUser)
// // 			if lErr3 != nil {
// // 				log.Println("SSUNSP03", lErr3)
// // 				return lErr3
// // 			}
// // 		}
// // 	}
// // 	log.Println("updatePendingSgbStatus (+)")
// // 	return nil

// // }

// // Pupose:This Method is used to Construct JvReq Struct into an JvData Struct .
// // Parameters:

// // 	Jv JvReqStruct

// // Response:

// // 	==========
// // 	*On Sucess
// // 	==========
// // 		{
// // 			MasterId:
// // 			ActionCode
// // 			BidId : 005156
// // 			JVamount : 5000
// // 			ClientId :FT000069
// // 			OrderNo : 005156
// // 			Price :8000
// // 			Transaction :
// // 		},nil

// // 	==========
// // 	*On Error
// // 	==========
// // 	"",error

// // Author:PRASHANTH
// // Date: 27SEP2023
// // */
// //
// //	This Method is used to Construct JvReq Struct into an JvData Struct
// //
// // author prashanth
// func JvConstructor(Jv JvReqStruct, Flag string) (JvDataStruct, error) {
// 	log.Println("NseJvConstruct (+)")
// 	var JvData JvDataStruct
// 	// JvData.MasterId = Jv.,
// 	JvData.ActionCode = Jv.ActionCode
// 	JvData.BidId = Jv.OrderNo
// 	// JvData.JVamount = Jv.Amount
// 	JvData.ClientId = Jv.ClientId
// 	JvData.OrderNo = Jv.OrderNo
// 	JvData.Price = Jv.Price
// 	JvData.Unit = Jv.Unit

// 	if Flag == "C" {
// 		JvData.Transaction = "C"
// 	} else if Flag == "R" {
// 		JvData.Transaction = "F"
// 	}

// 	lPrice, lErr1 := strconv.Atoi(Jv.Price)
// 	if lErr1 != nil {
// 		log.Println("SSJVC01", lErr1)
// 		return JvData, lErr1
// 	}

// 	lUnit, lErr2 := strconv.Atoi(Jv.Unit)
// 	if lErr2 != nil {
// 		log.Println("SSJVC02", lErr2)
// 		return JvData, lErr2
// 	}
// 	JvData.JVamount = strconv.Itoa(lPrice * lUnit)

// 	log.Println("NseJvConstruct (-)")
// 	return JvData, nil
// }

// /* Pupose:This Method is used to Construct JvReq Struct into an JvData Struct .
// Parameters:

// 	pRespSgbdata nsesgb.SgbAddResStruct

// re

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	{
// 		"status": S,
// 		"reason": ""

// 	},
// 	nil

// 	==========
// 	*On Error
// 	{
// 		"status": E,
// 		"reason": "Can't able to Update the data in database"
// 	}
// 	==========
// 	"",error

// Author:PRASHANTH
// Date: 27SEP2023

// author prashanth */
// //  This Method is used to Update SgbDetails For Nse

// // func updateNseSgbDetails(pRespSgbdata nsesgb.SgbAddResStruct, pUser string) error {
// // 	log.Println("updateNseSgbDetails(+)")

// // 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// // 	if lErr1 != nil {
// // 		log.Println("SSUNSDO1", lErr1)
// // 		return lErr1
// // 	} else {
// // 		defer lDb.Close()
// // 		lCoreString := `update a_sgb_orderdetails d
// // 							set d.RespSubscriptionunit = ?,d.RespRate = ?,d.ActionCode = ?,UpdatedDate = now(),UpdatedBy = ?
// // 							where d.OrderNo = ?`

// // 		_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Quantity, pRespSgbdata.Price,
// // 			pRespSgbdata.Status, pUser, pRespSgbdata.OrderNumber)
// // 		if lErr1 != nil {
// // 			log.Println("SSUNSDO2", lErr1)
// // 			return lErr1
// // 		}
// // 	}
// // 	log.Println("updateNseSgbDetails(-)")
// // 	return nil
// // }

// /* Pupose:This Method is used to Update Sgb Header For Nse.
// Parameters:

// 	pRespSgbdata nsesgb.SgbAddResStruct,user,Status

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	{
// 		"status": S,
// 		"reason": ""

// 	},
// 	nil

// 	==========
// 	*On Error
// 	{
// 		"status": E,
// 		"reason": "Can't able to Update the data from database"
// 	}
// 	==========
// 	"",error

// Author:PRASHANTH
// Date: 27SEP2023
// */
// // This Method is used to Update Sgb Header For Nse

// // func updateNseSgbHeader(pRespSgbdata nsesgb.SgbAddResStruct, pUser string, pStatus string) error {
// // 	log.Println("updateSgbHeader(+)")

// // 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// // 	if lErr1 != nil {
// // 		log.Println("SSUNSHO1", lErr1)
// // 		return lErr1
// // 	} else {
// // 		defer lDb.Close()
// // 		lCoreString := `update a_sgb_orderheader  h
// // 							set h.ScripId = ?,h.PanNo = ?,h.Depository = ?,
// // 							h.DpId=?,h.ClientBenfId=?
// // 							,h.status=?,
// // 							h.UpdatedDate = now(),h.UpdatedBy = ?
// // 							where h.Id in (select d.HeaderId
// // 							from a_sgb_orderdetails d,a_sgb_orderheader h
// // 							where d.HeaderId = h.Id and d.OrderNo = ?)`

// // 		_, lErr1 = lDb.Exec(lCoreString, pRespSgbdata.Symbol, pRespSgbdata.Pan,
// // 			pRespSgbdata.Depository, pRespSgbdata.DpId, pRespSgbdata.ClientBenId, pStatus, pUser, pRespSgbdata.OrderNumber)
// // 		if lErr1 != nil {
// // 			log.Println("SSUNSHO2", lErr1)
// // 			return lErr1
// // 		}
// // 	}

// // 	log.Println("updateSgbHeader(-)")
// // 	return nil
// // }

// /* Pupose:This Method is used to Update Sgb Header For Nse.
// Parameters:

// 	pRespSgbdata nsesgb.SgbAddResStruct,user,Status

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	{
// 		"status": S,
// 		"reason": ""

// 	},
// 	nil

// 	==========
// 	*On Error
// 	{
// 		"status": E,
// 		"reason": "Can't able to Update the data from database"
// 	}
// 	==========
// 	"",error

// Author:PRASHANTH
// Date: 27SEP2023
// */
// // This Method is used to update Bid Tracking Details of Nse

// // func updateNseBidTracking(pRespSgbdata nsesgb.SgbAddResStruct, pReqSgbData nsesgb.SgbAddReqStruct, pUser string) error {
// // 	log.Println("updateBidTracking(+)")

// // 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// // 	if lErr1 != nil {
// // 		log.Println("SSUNBO1", lErr1)
// // 		return lErr1
// // 	} else {
// // 		defer lDb.Close()
// // 		lCoreString := `update a_sgbtracking_table  b
// // 							set b.ApplicationStatus =?,UpdatedDate = now(),UpdatedBy = ?
// // 							where b.OrderNo = ? and b.ActivityType = ? and b.Unit = ?`

// // 		_, lErr1 = lDb.Exec(lCoreString, common.SUCCESS, pUser, pRespSgbdata.OrderNumber, pReqSgbData.ActivityType, pRespSgbdata.Quantity)
// // 		if lErr1 != nil {
// // 			log.Println("SSUNBO2", lErr1)
// // 			return lErr1
// // 		}
// // 	}

// // 	log.Println("updateBidTracking(-)")
// // 	return nil
// // }

// // func NseReqConstruct(BSEStruct []bsesgb.SgbReqStruct, JVStruct []JvReqStruct) (nsesgb.SgbAddReqStruct, error) {
// // 	log.Println("NseReqConstruct (+)")
// // 	var NseStruct nsesgb.SgbAddReqStruct
// // 	var NseReqArr []nsesgb.SgbAddReqStruct
// // 	for _, JV := range JVStruct {
// // 		NseStruct.Symbol = JV.Symbol
// // 		var lErr1 error
// // 		NseStruct.Price, lErr1 = common.ConvertStringToFloat(JV.Price)
// // 		if lErr1 != nil {
// // 			log.Println("lErr1", lErr1)
// // 			return NseStruct, lErr1
// // 		}
// // 		NseStruct.Quantity = strconv.Atoi(JV.Unit)
// // 		NseStruct.ActivityType = JV.ActionCode
// // 		NseStruct.PhysicalDematFlag = "D"
// // 		NseStruct.ClientCode = ""
// // 		NseStruct.OrderNumber = strconv.Atoi(JV.OrderNo)
// // 		NseStruct.ClientRefNumber = ""

// // 	}

// // 	for _, Bse := range BSEStruct {
// // 		NseStruct.ClientBenId = Bse.ClientBenfId
// // 		NseStruct.Depository = Bse.Depository
// // 		NseStruct.DpId = Bse.DpId
// // 		NseStruct.Pan = Bse.PanNo
// // 	}

// // 	NseReqArr = append(NseReqArr, NseStruct)

// // 	log.Println("NseReqConstruct (-)")
// // 	return NseStruct, nil
// // }

// // func NseReqConstruct() ([]NseReqStruct, error) {
// // 	log.Println("NseReqConstruct (+)")
// // 	var NseStruct NseReqStruct
// // 	var NseReqArr []NseReqStruct
// // 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// // 	if lErr1 != nil {
// // 		log.Println("SSFP01", lErr1)
// // 		return NseReqArr, lErr1
// // 	} else {
// // 		defer lDb.Close()

// // 		lCoreString := `select m.Symbol ,d.OrderNo ,d.ReqRate e,d.ReqSubscriptionUnit ,'D', '',h.PanNo,h.Depository ,h.ClientBenfId ,h.DpId, (CASE
// // 		WHEN d.ActionCode = 'M' THEN 'Modify'
// // 		WHEN d.ActionCode = 'N' THEN 'New'
// // 		WHEN d.ActionCode = 'D' THEN 'Delete'
// // 		ELSE d.ActionCode
// // 	END ) AS ActionDescription,h.clientId
// // 		from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
// // 		where h.MasterId = m.id
// // 		and d.HeaderId = h.Id
// // 		and h.Status = "pending"
// // 		and m.BiddingStartDate <= curdate()
// // 		or m.BiddingEndDate >= curdate()
// // 		and time(now()) between m.DailyStartTime and m.DailyEndTime
// // 		and h.cancelFlag != 'Y' and m.Exchange = "NSE"`

// // 		lRows, lErr2 := lDb.Query(lCoreString)
// // 		if lErr2 != nil {
// // 			log.Println("SSFP02", lErr2)
// // 			return NseReqArr, lErr2
// // 		} else {

// // 			for lRows.Next() {
// // 				lErr3 := lRows.Scan(&NseStruct.Symbol, &NseStruct.OrderNumber, &NseStruct.Price, &NseStruct.Quantity, &NseStruct.PhysicalDematFlag, &NseStruct.ClientCode, &NseStruct.Pan, &NseStruct.Depository, &NseStruct.ClientBenId, &NseStruct.DpId, &NseStruct.ActivityType, &NseStruct.ClientId)
// // 				if lErr3 != nil {
// // 					log.Println("SSFP03", lErr3)
// // 					return NseReqArr, lErr3
// // 				} else {

// // 					NseReqArr = append(NseReqArr, NseStruct)
// // 				}

// // 			}
// // 		}
// // 	}
// // 	log.Println("NseReqArr", NseReqArr)

// // 	log.Println("NseReqConstruct (-)")
// // 	return NseReqArr, nil
// // }
// /* Pupose:This method is used to update order status during schedule for Sgb Place Order  .
// Parameters:

// 	processFlag string, ScheduleStatus string, pHeaderId int, pBrokerId int, pJv {}

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	nil

// 	==========
// 	*On Error
// 	==========
// 	error

// Author:PRASHANTH
// Date: 19DEC2023
// */
// func updateSchStatus(processFlag string, ScheduleStatus string, pHeaderId int, pBrokerId int, pJv JvReqStruct) error {
// 	log.Println("updateSchStatus(+)")
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SUSS01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		// log.Println(" processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId", processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId)
// 		lCoreString := ` update a_sgb_orderheader
// 		                 set  ProcessFlag =  ? , ScheduleStatus = ?
// 	                     where ClientId = ?
// 	                     and  brokerId = ?
// 	                     and Id = ?`
// 		_, lErr2 := lDb.Exec(lCoreString, processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId)
// 		// log.Println("lCoreString", lCoreString)
// 		if lErr2 != nil {
// 			log.Println("SUSS02", lErr2)
// 			return lErr2

// 		}
// 	}

// 	log.Println("updateSchStatus (-)")
// 	return nil
// }

// func updateOrderDetails(pJvRec JvReqStruct, pExchange string, pBrokerId int, processFlag string, ScheduleStatus string) error {
// 	log.Println("updateOrderDetails (+)")

// 	lHeaderId, lErr1 := GetOrderId(pJvRec, pExchange)
// 	if lErr1 != nil {
// 		log.Println("SUOD01", lErr1)
// 		return lErr1
// 	} else {
// 		lErr2 := updateSchStatus(processFlag, ScheduleStatus, lHeaderId, pBrokerId, pJvRec)
// 		if lErr2 != nil {
// 			log.Println("SUOD02", lErr2)
// 			return lErr2
// 		}

// 	}

// 	log.Println("updateOrderDetails (-)")
// 	return nil
// }

// func UpdateExchFailed(pJvRec JvReqStruct, pExchange string, pBrokerId int) error {
// 	log.Println("UpdateExchFailed (+)")

// 	lHeaderId, lErr1 := GetOrderId(pJvRec, pExchange)
// 	if lErr1 != nil {
// 		log.Println("SUEF01", lErr1)
// 		return lErr1
// 	} else {
// 		lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 		if lErr1 != nil {
// 			log.Println("SUSS01", lErr1)
// 			return lErr1
// 		} else {
// 			defer lDb.Close()

// 			// log.Println(" processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId", processFlag, ScheduleStatus, pJv.ClientId, pBrokerId, pHeaderId)
// 			lCoreString := ` update a_sgb_orderheader h
// 							set h.Status = ?,h.UpdatedBy = ?,h.UpdatedDate = now()
// 							where h.ClientId = ?
// 							and h.brokerId = ?
// 							and h.Id = ?`
// 			_, lErr2 := lDb.Exec(lCoreString, common.FAILED, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
// 			// log.Println("lCoreString", lCoreString)
// 			if lErr2 != nil {
// 				log.Println("SUSS02", lErr2)
// 				return lErr2

// 			}
// 		}
// 	}
// 	log.Println("UpdateExchFailed (-)")
// 	return nil
// }
