package ncbschedule

import (
	"bytes"
	"fcs23pkg/apps/Ncb/ncbplaceorder"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nsencb"
	"fcs23pkg/integration/techexcel"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type NcbStruct struct {
	MasterId int
	Symbol   string
	Exchange string
}

type JvStatusStruct struct {
	JvAmount    string `json:"jvAmount"`
	JvStatus    string `json:"jvStatus"`
	JvStatement string `json:"jvStatement"`
	JvType      string `json:"jvType"`
}

type JvDataStruct struct {
	Unit              string `json:"unit"`
	Price             string `json:"price"`
	JVamount          string `json:"jvAmount"`
	OrderNo           int    `json:"orderNo"`
	ClientId          string `json:"clientId"`
	ActionCode        string `json:"actionCode"`
	Transaction       string `json:"transaction"`
	ApplicationNumber string `json:"applicationNumber"`
}

//----------------------------------------------------------------
// this method is used to get the Brokerid from database
//----------------------------------------------------------------
func GetNcbBrokers(pExchange string) ([]int, error) {
	log.Println("GetNcbBrokers(+)")

	var lBrokerArr []int
	var lBrokerRec int

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GNB01", lErr1)
		return lBrokerArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select  md.BrokerId 
		                from a_ipo_memberdetails md
		                where md.AllowedModules like '%Gsec%'                     
		                and md.Flag = 'Y'
		                and md.OrderPreference = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pExchange)
		if lErr2 != nil {
			log.Println("GNB02", lErr2)
			return lBrokerArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lBrokerRec)
				if lErr3 != nil {
					log.Println("GNB03", lErr3)
					return lBrokerArr, lErr3
				} else {
					lBrokerArr = append(lBrokerArr, lBrokerRec)
				}
			}
		}
	}
	log.Println("GetNcbBrokers(-)")
	return lBrokerArr, nil
}

//---------------------------------------------------------------------------------
// this method is used to filter the Jv posted records
//---------------------------------------------------------------------------------
func PostJvForOrder(pNcbReqRec nsencb.NcbAddReqStruct, pJvDetailRec JvNcbReqStruct, r *http.Request, pValidNcb NcbStruct, pBrokerId int) (int, int, int, int, int, error) {
	log.Println("PostJvForOrder (+)")

	var lJvsuccess int
	var lJvfailed int
	var lExchangesuccess int
	var lExchangeFailed int
	var lReversedJv int
	lVerify := "VERIFY"
	lBlockClient := "BLOCKCLIENT"
	lInsufficient := "INSUFFICIENT"

	lClientFund, lErr1 := ncbplaceorder.VerifyFundDetails(pJvDetailRec.ClientId)
	log.Println("lClientFund", lClientFund)

	if lErr1 != nil {
		log.Println("PJFO01", lErr1)
		lJvfailed++

		lVerifyClient, lErr2 := ConstructMailforJvProcess(pJvDetailRec, lVerify)
		if lErr2 != nil {
			log.Println("Error: 2", lErr2)
		} else {
			lErr3 := emailUtil.SendEmail(lVerifyClient, "JV Verifying Client Account Details")

			if lErr3 != nil {
				log.Println("Error: 3", lErr3)
			}
		}

	} else {
		if lClientFund.Status == common.SuccessCode {

			if lClientFund.AccountBalance < float64(int(pNcbReqRec.Price)*pNcbReqRec.InvestmentValue) {
				log.Println("Low Account Balance")

				lPricePerQuantity := int(pNcbReqRec.Price) / pNcbReqRec.InvestmentValue

				if lClientFund.AccountBalance < float64(lPricePerQuantity) {

					lInsufficientClient, lErr3 := ConstructMailforJvProcess(pJvDetailRec, lInsufficient)

					if lErr3 != nil {
						log.Println("Error:3", lErr3)
					} else {
						lJvfailed++
						lErr4 := emailUtil.SendEmail(lInsufficientClient, "Insufficient Fund in Client Account for NCB")
						if lErr4 != nil {
							log.Println("Error:4", lErr4.Error())
						}
					}
				} else {
					log.Println("unit changed")

					MaximumUnits := lClientFund.AccountBalance / float64(lPricePerQuantity)
					pJvDetailRec.Unit = fmt.Sprintf("%f", MaximumUnits)

					Jvstatus := BlockClientFund(pJvDetailRec, r, "C")
					log.Println("Jvstatus", Jvstatus)
					pJvDetailRec.JvAmount = Jvstatus.JvAmount
					pJvDetailRec.JvStatement = Jvstatus.JvStatement
					pJvDetailRec.JvStatus = Jvstatus.JvStatus
					pJvDetailRec.JvType = Jvstatus.JvType

					if Jvstatus.JvStatus == common.SuccessCode {
						lJvsuccess++

						if pValidNcb.Exchange == common.BSE {
							log.Println("pValidNcb.Exchange", pValidNcb.Exchange)
						} else if pValidNcb.Exchange == common.NSE {
							log.Println("pValidNcb.Exchange", pValidNcb.Exchange)

							lRespRec, lErr5 := exchangecall.ApplyNseNcb(pNcbReqRec, common.AUTOBOT, pBrokerId)
							if lErr5 != nil {
								log.Println("Error:5", lErr5)
								lExchangeFailed++
							} else {

								if lRespRec.Status != "success" {
									lExchangeFailed++
									lReversedJv++
									pJvDetailRec = ReverseProcess(pJvDetailRec, r)
								} else {

									UpdateNseRecord(lRespRec, pJvDetailRec, pBrokerId)
									lExchangesuccess++

									SuccessMail, lErr6 := constructSuccessmail(pJvDetailRec, "S")
									if lErr6 != nil {
										log.Println("Error:6", lErr6)
									} else {
										lErr7 := emailUtil.SendEmail(SuccessMail, "Order Placed Succesfully on NCB")
										if lErr7 != nil {
											log.Println("Error:7", lErr7.Error())
										}
									}
								}
							}

						}
					} else {
						lJvfailed++
						lBlockFund, lErr8 := ConstructMailforJvProcess(pJvDetailRec, lBlockClient)
						if lErr8 != nil {
							log.Println("Error:8", lErr8)
						} else {
							lErr9 := emailUtil.SendEmail(lBlockFund, "JV issue to Deduct Amount From Client Account ")
							if lErr9 != nil {
								log.Println("Error:9", lErr9.Error())
							}
						}
					}
				}

			} else {
				Jvstatus := BlockClientFund(pJvDetailRec, r, "C")
				log.Println("Jvstatus", Jvstatus)
				pJvDetailRec.JvAmount = Jvstatus.JvAmount
				pJvDetailRec.JvStatement = Jvstatus.JvStatement
				pJvDetailRec.JvStatus = Jvstatus.JvStatus
				pJvDetailRec.JvType = Jvstatus.JvType

				if Jvstatus.JvStatus == common.SuccessCode {
					lJvsuccess++
					if pValidNcb.Exchange == common.BSE {
						log.Println("pValidNcb.Exchange", pValidNcb.Exchange)

					} else if pValidNcb.Exchange == common.NSE {

						lRespRec, lErr10 := exchangecall.ApplyNseNcb(pNcbReqRec, common.AUTOBOT, pBrokerId)
						if lErr10 != nil {
							log.Println("Error :10", lErr10)
							lExchangeFailed++
						} else {

							if lRespRec.Status != "success" {
								lExchangeFailed++
								pJvDetailRec = ReverseProcess(pJvDetailRec, r)
								lReversedJv++
							} else {

								UpdateNseRecord(lRespRec, pJvDetailRec, pBrokerId)
								lExchangesuccess++

								SuccessMail, lErr11 := constructSuccessmail(pJvDetailRec, "S")
								if lErr11 != nil {
									log.Println("Error:11", lErr11)
								} else {
									lErr12 := emailUtil.SendEmail(SuccessMail, "Order Placed Succesfully on NCB")
									if lErr12 != nil {
										log.Println("Error:12", lErr12.Error())
									}
								}
							}

						}
					}
				} else {
					lJvfailed++

					lBlockFund, lErr13 := ConstructMailforJvProcess(pJvDetailRec, lBlockClient)
					if lErr13 != nil {
						log.Println("Error:13", lErr13)
					} else {
						lErr14 := emailUtil.SendEmail(lBlockFund, "JV issue to Deduct Amount From Client Account ")
						if lErr14 != nil {
							log.Println("Error:14", lErr14.Error())
						}
					}
				}
			}
		} else {
			lJvfailed++
			lVerifyClient, lErr15 := ConstructMailforJvProcess(pJvDetailRec, lVerify)
			if lErr15 != nil {
				log.Println("Error:15", lErr15)
			} else {
				lErr16 := emailUtil.SendEmail(lVerifyClient, "JV Verifying Client Account Details")
				if lErr16 != nil {
					log.Println("Error:16", lErr16.Error())
				}
			}
		}

	}

	log.Println("PostJvForOrder (-)")
	return lJvsuccess, lJvfailed, lExchangesuccess, lExchangeFailed, lReversedJv, nil
}

func ConstructMailforJvProcess(pJV JvNcbReqStruct, pStatus string) (emailUtil.EmailInput, error) {
	type dynamicEmailStruct struct {
		Date     string `json:"date"`
		ClientId string `json:"clientId"`
		OrderNo  int    `json:"orderNo"`
		JvUnit   string `json:"jvUnit"`
		JvAmount string `json:"jvAmount"`
	}
	var lJVEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("./toml/emailconfig.toml")
	lJVEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lJVEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lJVEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	// Mail Id for It support is Not Added in Toml
	lJVEmailContent.ToEmailId = fmt.Sprintf("%v", config.(map[string]interface{})["ToEmailId"])
	lJVEmailContent.Subject = "JV Orders"
	var html string

	if pStatus == "VERIFY" {
		html = "html/VerifyClientFund.html"
	} else if pStatus == "BLOCKCLIENT" {
		html = "html/BlockClientMail.html"
	} else if pStatus == "INSUFFICIENT" {
		html = "html/InsufficientAmount.html"
	} else if pStatus == "REVERSEJV" {
		html = "html/Exchagestatus.html"
	}
	currentTime := time.Now()
	currentDate := currentTime.Format("02-01-2006")

	lTemp, lErr := template.ParseFiles(html)
	if lErr != nil {
		log.Println("CMFJP01", lErr)
		return lJVEmailContent, lErr
	} else {

		var lTpl bytes.Buffer
		var lDynamicEmailVal dynamicEmailStruct

		if pStatus == "VERIFY" {
			//  IT Dept To verify client Account
			lDynamicEmailVal.Date = currentDate
			lDynamicEmailVal.ClientId = pJV.ClientId
			lDynamicEmailVal.OrderNo = pJV.OrderNo
		} else if pStatus == "BLOCKCLIENT" {
			//  IT Dept To Deducting Amount from client Account
			lDynamicEmailVal.Date = currentDate
			lDynamicEmailVal.ClientId = pJV.ClientId
			lDynamicEmailVal.OrderNo = pJV.OrderNo
			lDynamicEmailVal.JvUnit = pJV.Unit
			lDynamicEmailVal.JvAmount = pJV.Amount
		} else if pStatus == "INSUFFICIENT" {
			lDynamicEmailVal.Date = currentDate
			lDynamicEmailVal.ClientId = pJV.ClientId
			lDynamicEmailVal.OrderNo = pJV.OrderNo
			lDynamicEmailVal.JvUnit = pJV.Unit
			lDynamicEmailVal.JvAmount = pJV.Amount
			//  for No Balance Amount Directly to client
			lJVEmailContent.ToEmailId = pJV.Mail
		} else if pStatus == "REVERSEJV" {
			// IT Dept failed During Exchange processs
			lDynamicEmailVal.Date = currentDate
			lDynamicEmailVal.ClientId = pJV.ClientId
			lDynamicEmailVal.OrderNo = pJV.OrderNo
			lDynamicEmailVal.JvUnit = pJV.Unit
			lDynamicEmailVal.JvAmount = pJV.Amount
		}
		lTemp.Execute(&lTpl, lDynamicEmailVal)
		lEmailbody := lTpl.String()

		lJVEmailContent.Body = lEmailbody
	}

	return lJVEmailContent, nil
}

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

Author:KAVYA DHARSAHANI
Date: 08NOV2023
*/
func BlockClientFund(pJvDetailRec JvNcbReqStruct, pRequest *http.Request, pFlag string) JvStatusStruct {
	log.Println("BlockClientFund(+)")

	pJvProcessRec := JvConstructor(pJvDetailRec, pFlag)

	// this variables is used to get Status
	var lJVReq techexcel.JvInputStruct
	// var lReqDtl apigate.RequestorDetails
	var lJvStatusRec JvStatusStruct

	lConfigFile := common.ReadTomlConfig("toml/techXLAPI_UAT.toml")
	lCocd := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["PaymentCode"])

	intVal := pJvProcessRec.OrderNo
	// Convert int to string
	strVal := strconv.Itoa(intVal)

	lJVReq.COCD = lCocd
	lJVReq.VoucherDate = time.Now().Format("02/01/2006")
	lJVReq.BillNo = "NCB" + pJvProcessRec.ClientId
	lJVReq.SourceTable = "a_ncb_orderdetails"
	lJVReq.SourceTableKey = strVal
	lJVReq.Amount = pJvProcessRec.JVamount
	lJVReq.WithGST = "N"

	lFTCaccount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NseNcborderAccount"])

	if pJvProcessRec.Transaction == "F" {
		lJVReq.AccountCode = lFTCaccount
		lJVReq.CounterAccount = pJvProcessRec.ClientId

		lJVReq.Narration = "Failed of Ncb Order Refund  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		lJvStatusRec.JvStatement = lJVReq.Narration
		lJvStatusRec.JvType = "C"
	} else if pJvProcessRec.Transaction == "C" {
		lJVReq.AccountCode = pJvProcessRec.ClientId
		lJVReq.CounterAccount = lFTCaccount
		lJVReq.Narration = "NCB Purchase  " + pJvProcessRec.Unit + " * " + pJvProcessRec.Price + " on " + lJVReq.VoucherDate + " - " + pJvProcessRec.ClientId
		lJvStatusRec.JvStatement = lJVReq.Narration
		lJvStatusRec.JvType = "D"
	}

	lJvStatusRec.JvStatus = "S"
	lJvStatusRec.JvAmount = pJvProcessRec.JVamount
	log.Println("JV Record:", lJVReq)

	log.Println("BlockClientFund(-)")
	return lJvStatusRec
}

func JvConstructor(Jv JvNcbReqStruct, Flag string) JvDataStruct {
	log.Println("NseJvConstruct (+)")
	var JvData JvDataStruct
	// JvData.MasterId = Jv.,
	JvData.ActionCode = Jv.ActionCode
	JvData.ApplicationNumber = Jv.ApplicationNumber
	JvData.JVamount = Jv.Amount
	JvData.ClientId = Jv.ClientId
	JvData.OrderNo = Jv.OrderNo
	JvData.Price = Jv.Price
	JvData.Unit = Jv.Unit

	if Flag == "N" {
		JvData.Transaction = "C"
	} else if Flag == "N" {
		JvData.Transaction = "F"

	}

	log.Println("NseJvConstruct (-)")
	return JvData
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


Author:KAVYA DHARSHANI
Date: 08NOV2023
*/

func ReverseProcess(pJvDetailRec JvNcbReqStruct, r *http.Request) JvNcbReqStruct {
	log.Println("ReverseProcess (+)")

	lJvDatatoReturn := pJvDetailRec

	Jvstatus := BlockClientFund(pJvDetailRec, r, "R")
	log.Println("Jvstatus", Jvstatus)
	lJvDatatoReturn.JvAmount = Jvstatus.JvAmount
	lJvDatatoReturn.JvStatement = Jvstatus.JvStatement
	lJvDatatoReturn.JvStatus = Jvstatus.JvStatus
	lJvDatatoReturn.JvType = Jvstatus.JvType
	// //-----Mail To Accounts Team-------
	if Jvstatus.JvStatus == common.ErrorCode {
		REVERSEJS := "REVERSEJV"

		lFailedJvMail, lErr1 := ConstructMailforJvProcess(pJvDetailRec, REVERSEJS)
		if lErr1 != nil {
			log.Println("NCBRP01", lErr1.Error())
		} else {
			lString := "JV Failed in NCB Order"
			lErr2 := emailUtil.SendEmail(lFailedJvMail, lString)
			if lErr2 != nil {
				log.Println("RPRP02", lErr2.Error())
			}
		}
	}

	log.Println("ReverseProcess (-)")
	return lJvDatatoReturn
}

/*
Pupose:This method is used to get the email Input for success NCB Place Order  .
Parameters:

	pNcbClientDetails {},pStatus string

Response:

	==========
	*On Sucess
	==========
	{
     Name: clientName
     Status: S
     OrderDate: 8Nov2023
     OrderNumber: 121345687
     Symbol: Ncb test
     Unit: 5
     Price: 500
     Amount:2500
     Activity : M
	},nil

	==========
	*On Error
	==========
	"",error

Author:KAVYA DHARSHANI
Date: 8NOV2023
*/
func constructSuccessmail(pNcbClientDetails JvNcbReqStruct, pStatus string) (emailUtil.EmailInput, error) {
	log.Println("constructSuccessmail (+)")
	type dynamicEmailStruct struct {
		Name        string
		Status      string
		OrderDate   string
		Symbol      string
		OrderNumber int
		Unit        string
		Price       string
		Amount      string
		Activity    string
	}

	var lEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("toml/emailconfig.toml")

	lEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	lEmailContent.Subject = "NCB Order"
	html := "html/NcbOrderTemplate.html"

	lTemp, lErr := template.ParseFiles(html)
	if lErr != nil {
		log.Println("CSM01", lErr)
		return lEmailContent, lErr
	} else {
		var lTpl bytes.Buffer
		var lDynamicEmailVal dynamicEmailStruct
		lDynamicEmailVal.Name = pNcbClientDetails.ClientName
		lDynamicEmailVal.Amount = pNcbClientDetails.Amount
		lDynamicEmailVal.Unit = pNcbClientDetails.Unit
		lDynamicEmailVal.Price = pNcbClientDetails.Price
		lDynamicEmailVal.OrderDate = pNcbClientDetails.OrderDate
		lDynamicEmailVal.OrderNumber = pNcbClientDetails.OrderNo
		lDynamicEmailVal.Symbol = pNcbClientDetails.Symbol
		if pStatus == "S" {
			lDynamicEmailVal.Status = common.SUCCESS
		} else {
			lDynamicEmailVal.Status = common.FAILED
		}

		lEmailContent.ToEmailId = pNcbClientDetails.Mail

		lTemp.Execute(&lTpl, lDynamicEmailVal)
		lEmailbody := lTpl.String()

		lEmailContent.Body = lEmailbody
	}
	log.Println("constructSuccessmail (-)")
	return lEmailContent, nil
}

//---------------------------------------------------------------------------------
// this method is used to update the reocrds od NSE in db
//---------------------------------------------------------------------------------
func UpdateNseRecord(pRespRec nsencb.NcbAddResStruct, pJvRec JvNcbReqStruct, pBrokerId int) error {
	log.Println("UpdateNseRecord (+)")

	lErr1 := UpdateHeader(pRespRec, pJvRec, common.NSE, pBrokerId)
	if lErr1 != nil {
		log.Println("NCBOFUNR01", lErr1)
		return lErr1
	} else {
		log.Println("Nse Record Updated Successfully")
	}

	log.Println("UpdateNseRecord (-)")
	return nil
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


Author:KAVYADHARSHANI
Date: 08NOV2023
*/

func UpdateHeader(pNseRespRec nsencb.NcbAddResStruct, pJvRec JvNcbReqStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateHeader (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PUH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lHeaderId, lErr2 := GetOrderId(pJvRec, pExchange)
		if lErr2 != nil {
			log.Println("PUH02", lErr2)
			return lErr2
		} else {

			if lHeaderId != 0 {

				if pExchange == common.BSE {
					log.Println("pExchange", pExchange)
				} else if pExchange == common.NSE {
					log.Println("pExchange", pExchange)

					lSqlString := `update a_ncb_orderheader h
				                     set h.Symbol = ?,h.Investmentunit=?,h.pan  = ?,h.PhysicalDematFlag=?, h.depository  =?, h.dpId =? , h.clientBenId  =?, h.clientRefNumber =? ,h.status = ?, h.StatusCode = ?, h.StatusMessage =?,h.ErrorCode  =?, h.ErrorMessage =?,,h.UpdatedBy = ?,h.UpdatedDate = now()
								   and h.clientId = ?
								   and h.brokerId  = ?
								   and h.Id = ? `

					_, lErr3 := lDb.Exec(lSqlString, pNseRespRec.Symbol, pNseRespRec.InvestmentValue, pNseRespRec.Pan, pNseRespRec.PhysicalDematFlag, pNseRespRec.Depository, pNseRespRec.DpId, pNseRespRec.ClientBenId, pNseRespRec.ClientRefNumber, common.SUCCESS, pNseRespRec.Status, pNseRespRec.OrderStatus, pNseRespRec.Reason, pNseRespRec.RejectionReason, common.AUTOBOT, pJvRec.ClientId, pBrokerId, lHeaderId)
					if lErr3 != nil {
						log.Println("PUH03", lErr3)
						return lErr3
					} else {
						lErr4 := UpdateDetail(pJvRec, pNseRespRec, lHeaderId)
						if lErr4 != nil {
							log.Println("PUH04", lErr4)
							return lErr4
						}
					}
				}
			}
		}
	}

	log.Println("UpdateHeader (-)")
	return nil
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


Author:KAVYADHARSHAMI
Date: 12JUNE2023
*/
func UpdateDetail(pJvRec JvNcbReqStruct, pRespArr nsencb.NcbAddResStruct, pHeaderId int) error {
	log.Println("UpdateDetail (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PUH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lSqlString := `update a_ncb_orderdetails d
		                 set d.OrderNo ,d.price = ?,d.Unit=?,d.ErrorCode =?, d.ErrorMessage = ?, d.JvStatus =?, d.JvAmount =?d.JvStatement = ?, d.JvType =?, d.UpdatedBy = ?,d.UpdatedDate = now()
		               where d.headerId = ?
		               and d.OrderNo  =?`
		_, lErr2 := lDb.Exec(lSqlString, pRespArr.OrderNumber, pRespArr.Price, pRespArr.InvestmentValue, pRespArr.Reason, pRespArr.RejectionReason, pJvRec.JvStatus, pJvRec.JvAmount, pJvRec.JvStatement, pJvRec.JvType, common.AUTOBOT, pHeaderId, pJvRec.OrderNo)
		if lErr2 != nil {
			log.Println("PUH02", lErr2)
			return lErr2
		} else {
			log.Println("UpdateDetailSuccessfully")

		}
	}
	log.Println("UpdateDetail (-)")
	return nil
}

func GetOrderId(pJvRec JvNcbReqStruct, pExchange string) (int, error) {
	log.Println("GetOrderId (+)")

	var lHeaderId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GOI01", lErr1)
		return lHeaderId, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select (case when count(1) > 0 then h.Id else 0 end) Id
		                from a_ncb_orderheader h, a_ncb_orderdetails d
		                  where h.Id  = d.headerId
						  and h.clientId  = ?
						  and d.OrderNo  =?
						  and h. exchange =?
						  and h.MasterId  = (
								select n.id
								 from a_ncb_master n
								 where n.Symbol =?)`
		lRows, lErr2 := lDb.Query(lCoreString, pJvRec.ClientId, pJvRec.OrderNo, pExchange, pJvRec.Symbol)
		if lErr2 != nil {
			log.Println("GOI02", lErr2)
			return lHeaderId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lHeaderId)
				if lErr3 != nil {
					log.Println("GOI03", lErr3)
					return lHeaderId, lErr3
				}
				log.Println(lHeaderId)
			}
		}
	}

	log.Println("GetOrderId (-)")
	return lHeaderId, nil
}
