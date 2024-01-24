package placeorder

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/abhilogin"
	"fcs23pkg/apps/clientDetail"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fcs23pkg/integration/nse/nseipo"
	"fcs23pkg/util/emailUtil"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// this struct is used to get the request parameters.
type OrderReqStruct struct {
	UpiId         string           `json:"upiId"`
	UpiEndPoint   string           `json:"upiEndPoint"`
	Category      string           `json:"category"`
	Symbol        string           `json:"symbol"`
	MasterId      int              `json:"masterId"`
	ApplicationNo string           `json:"applicationNo"`
	PreApply      string           `json:"preApply"`
	BidDetails    []OrderBidStruct `json:"bids"`
}

// this struct is used to get the bid details.
type OrderBidStruct struct {
	ActivityType   string `json:"activityType"`
	BidReferenceNo string `json:"bidReferenceNo"`
	CutOff         bool   `json:"cutOff"`
	Price          int    `json:"price"`
	Quantity       int    `json:"quantity"`
	LotSize        int    `json:"lotSize"`
}

// this struct is used to send Response.
type OrderResStruct struct {
	AppStatus string `json:"appStatus"`
	AppReason string `json:"appReason"`
	Status    string `json:"status"`
	ErrMsg    string `json:"errMsg"`
}

// this struct is used to get application number from the  a_ipo_order_header table and  a_ipo_master table
type IpoOrderStruct struct {
	ApplicationNumber string `json:"applicationNumber"`
	ClientEmail       string `json:"clientEmail"`
	ClientId          string `json:"clientId"`
	ClientName        string `json:"clientName"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	OrderDate         string `json:"orderdate"`
	Amount            int    `json:"amount"`
	ActivityType      string `json:"activityType"`
}

type BidTrack struct {
	Id       int
	BidRefNo int
}

/*
Pupose:This API method is used to place a order
Request:

	lReqRec

Response:

	==========
	*On Sucess
	==========
	{
		"appResponse":[
			{
				"applicationNo":"FT000069130109",
				"bidRefNo":"2023061400000028",
				"bidStatus":"success",
				"reason":"",
				"appStatus":"success",
				"appReason":""
			}
		],
		"status":"S",
		"errMsg":""
	}
	==========
	*On Error
	==========

		{
		"appResponse":[],
		"status":"E",
		"errMsg":""
	}

Author:Pavithra
Date: 12JUNE2023
*/
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("PlaceOrder (+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
			log.Println(lErr, origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	// Checking the Api Method
	if r.Method == "POST" {
		// create a instance of array type OrderReqStruct.This variable is used to store the request values.
		var lReqRec OrderReqStruct
		// create a instance of array type OrderResStruct.This variable is used to store the response values.
		var lRespRec OrderResStruct
		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "PPO01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("PPO01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lErr2 := common.CustomError("UserDetails not Found")
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "PPO02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("PPO02", "UserDetails not Found"))
				return
			}
		}
		//added by naveen:to fetch the source (from where mobile or web)by cookie name
		source, lErr3 := abhilogin.GetSourceOfUser(r, common.ABHICookieName)
		if lErr3 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "PPO03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("PPO03", "Unable to get source"))
			return
		}

		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// Read the request body values in lBody variable
		lBody, lErr4 := ioutil.ReadAll(r.Body)
		log.Println(string(lBody))
		if lErr4 != nil {
			log.Println("PPO04", lErr4.Error())
			lRespRec.Status = common.ErrorCode
			fmt.Fprintf(w, helpers.GetErrorString("PPO04", "Unable to get your request now!!Try Again.."))
			return
		} else {
			// Unmarshal the request body values in lReqRec variable
			lErr5 := json.Unmarshal(lBody, &lReqRec)
			if lErr5 != nil {
				log.Println("PPO05", lErr5.Error())
				lRespRec.Status = common.ErrorCode
				fmt.Fprintf(w, helpers.GetErrorString("PPO05", "Unable to get your request now. Please try after sometime"))
				return
			} else {
				lEligibleToOrder, lErr6 := Ipo_EligibleToOrder(lReqRec, lBrokerId, lClientId)
				if lErr6 != nil {
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "PPO06" + lErr6.Error()
					fmt.Fprintf(w, helpers.GetErrorString("PPO06", "Unable to get your request now. Please try after sometime"))
					return
				} else {
					if !lEligibleToOrder {
						lErr7 := common.CustomError("You have already applied")
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "PPO07" + lErr7.Error()
						fmt.Fprintf(w, helpers.GetErrorString("PPO07", "You have already applied"))
						return
					} else if lEligibleToOrder {

						lExchange, lErr8 := FetchDirectory(lReqRec)
						// lExchange, lErr4 := memberdetail.BseNsePercentCalc(lBrokerId, "/ipo")
						// log.Println("lExchange", lExchange)
						if lErr8 != nil {
							lRespRec.Status = common.ErrorCode
							lRespRec.ErrMsg = "PPO08" + lErr8.Error()
							fmt.Fprintf(w, helpers.GetErrorString("PPO08", "Directory Not Found. Please try after sometime"))
							return
						} else {
							log.Println("exchange getting in fetch directory", lExchange)

							// lRespRec.Status = common.SuccessCode
							var lErrorRec ErrorStruct
							var lExchangeResp []nseipo.ExchangeRespStruct
							// commented by pavithra
							// log.Println("lReqRec", lReqRec)
							lExchangeReq, lClientMailId, lErr9 := ConstructReqStruct(lReqRec, lClientId)
							if lErr9 != nil {
								log.Println("PPO09", lErr9.Error())
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "PPO09" + lErr9.Error()
								fmt.Fprintf(w, helpers.GetErrorString("PPO09", "Unable to process your request now. Please try after sometime"))
								return
							} else {
								// commented by pavithra
								// log.Println(lReqRec.PreApply)
								// This method is Not Providing Mail ID Because Client Data is not inserting on DB
								// lClientMailId, lErr9 := clientDetail.GetClientEmailId(r, lClientId)
								//  Temperory Alternate method to create here to Get an Client Mail ID
								// lClientMailId, lErr7 := clientDetail.GetClientEmailId(lClientId)
								// if lErr7 != nil {
								// 	log.Println("PPO07", lErr7.Error())
								// 	lRespRec.Status = common.ErrorCode
								// 	fmt.Fprintf(w, helpers.GetErrorString("PPO07", "Unable to process your request!"))
								// 	return
								// } else {
								// check the placing order is a preapply applictaion
								if lReqRec.PreApply == "pre" {
									// for pending status updating

									//added by naveen:add one additional parameter source to insert in ipo header&and bid tracking
									//lErr8 := UpdatePending(lExchangeReq, lClientId, lReqRec, lExchange, lClientMailId, lBrokerId)
									lErr10 := UpdatePending(lExchangeReq, lClientId, lReqRec, lExchange, lClientMailId, lBrokerId, source)
									if lErr10 != nil {
										log.Println("PPO10", lErr10.Error())
										lRespRec.Status = common.ErrorCode
										lRespRec.ErrMsg = "PP010" + lErr10.Error()
										fmt.Fprintf(w, helpers.GetErrorString("PPO10", "Unable to Updating in local.. "))
										return
									} else {
										lRespRec.AppStatus = "Pending"
										lRespRec.AppReason = "Application saved in offline,it will be process when bidding will start."
										lRespRec.Status = common.SuccessCode
									}

								} else {

									lTodayAvailable, lErr11 := Function.CheckEndDate(lReqRec.MasterId)
									if lErr11 != nil {
										log.Println("PPO11", lErr11.Error())
										lRespRec.Status = common.ErrorCode
										lRespRec.ErrMsg = "PPO11" + lErr11.Error()
										fmt.Fprintf(w, helpers.GetErrorString("PPO11", "Please try after sometime"))
										return
									} else {

										if lTodayAvailable == "True" {

											if lExchangeReq != nil {
												//For Check the Current timestamp value
												lIndicator, lErr12 := Function.GetTime(lReqRec.MasterId)
												// commented by pavithra
												// log.Println("lIndicator", lIndicator)
												if lErr12 != nil {
													log.Println("PPO12", lErr12.Error())
													lRespRec.Status = common.ErrorCode
													lRespRec.ErrMsg = "PPO12" + lErr12.Error()
													fmt.Fprintf(w, helpers.GetErrorString("PPO12", "Unable to process request,Please try after sometime"))
													return
												} else {
													//-----------------------------------------
													if lIndicator == "True" {
														//added by naveen:add one additional parameter source to insert in ipo header
														//lErr11 := UpdatePending(lExchangeReq, lClientId, lReqRec, lExchange, lClientMailId, lBrokerId)
														lErr13 := UpdatePending(lExchangeReq, lClientId, lReqRec, lExchange, lClientMailId, lBrokerId, source)
														if lErr13 != nil {
															log.Println("PPO13", lErr13.Error())
															lRespRec.Status = common.ErrorCode
															lRespRec.ErrMsg = "PPO13" + lErr13.Error()

															fmt.Fprintf(w, helpers.GetErrorString("PPO13", "Unable to Updating in local.. "))
															return
														} else {
															lRespRec.AppStatus = "Pending"
															lRespRec.AppReason = "Application Saved in Offline,it will be process on next working day!"
															lRespRec.Status = common.SuccessCode
														}

													} else if lIndicator == "False" {
														if lExchange == common.BSE {
															//added by naveen:add one argument source to insert in bidtracking table
															//lRespRec, lExchangeResp, lErrorRec = BsePlaceOrder(lExchangeReq, lReqRec, lClientId, lClientMailId, lBrokerId)
															lRespRec, lExchangeResp, lErrorRec = BsePlaceOrder(lExchangeReq, lReqRec, lClientId, lClientMailId, lBrokerId, source)
															if lErrorRec.ErrCode != "" {
																// commented by pavithra
																// log.Println("lRespRec", lRespRec)
																// log.Println("lExchangeResp", lExchangeResp)
																// log.Println("lErrorRec", lErrorRec)
																fmt.Fprintf(w, helpers.GetErrorString(lErrorRec.ErrCode, lErrorRec.ErrMsg))
																return
															}
														} else if lExchange == common.NSE {
															//added by naveen:add one argument source to insert in bidtracking table
															//lRespRec, lExchangeResp, lErrorRec = NsePlaceOrder(lExchangeReq, lReqRec, lClientId, lClientMailId, lBrokerId)
															lRespRec, lExchangeResp, lErrorRec = NsePlaceOrder(lExchangeReq, lReqRec, lClientId, lClientMailId, lBrokerId, source)
															if lErrorRec.ErrCode != "" {
																// commented by pavithra
																// log.Println("lRespRec", lRespRec)
																// log.Println("lExchangeResp", lExchangeResp)
																// log.Println("lErrorRec", lErrorRec)
																fmt.Fprintf(w, helpers.GetErrorString(lErrorRec.ErrCode, lErrorRec.ErrMsg))
																return
															}
														}

														if lRespRec.Status != common.ErrorCode {
															//----------------->>>> Get Application Number <<<<<--------------------
															lGetApp, lErr14 := GetApplication(lExchangeResp)
															if lErr14 != nil {
																log.Println("PPO14", lErr14.Error())
																lRespRec.Status = common.ErrorCode
																lRespRec.ErrMsg = "PPO14" + lErr14.Error()
																fmt.Fprintf(w, helpers.GetErrorString("PPO14", "Invalid Application Number."))
																return
															} else {

																if lGetApp.ClientEmail != "" {
																	for lRespIdx := 0; lRespIdx < len(lExchangeResp); lRespIdx++ {
																		lEmail, lErr15 := ConstructMail(lGetApp, lExchangeResp[lRespIdx].Status)
																		if lErr15 != nil {
																			log.Println("PPO15", lErr15.Error())
																			lRespRec.Status = common.ErrorCode
																			lRespRec.ErrMsg = "PPO15" + lErr15.Error()
																			fmt.Fprintf(w, helpers.GetErrorString("PPO15", "Unable to process mail."))
																			return
																		} else {
																			// fmt.Println(lEmail)
																			lErr16 := emailUtil.SendEmail(lEmail, lExchangeResp[lRespIdx].Status)
																			if lErr16 != nil {
																				log.Println("PPO16", lErr16.Error())
																				lRespRec.Status = common.ErrorCode
																				lRespRec.ErrMsg = "PPO16" + lErr16.Error()
																				fmt.Fprintf(w, helpers.GetErrorString("PPO16", "Unable to sent mail."))
																				return
																			}

																		}
																	}
																}
															}
														}
													} else {
														lErr17 := common.CustomError("Unable to proceed Application!!!")
														lRespRec.Status = common.ErrorCode
														lRespRec.ErrMsg = "PPO17" + lErr17.Error()
														log.Println("PPO17", "Unable to proceed Application!!!")
														fmt.Fprintf(w, helpers.GetErrorString("PPO17", "Unable to proceed Application!!!"))
														return
													}
												}
											} else {
												lErr18 := common.CustomError("No Records Found for Exchange")
												lRespRec.Status = common.ErrorCode
												lRespRec.ErrMsg = "PPO18" + lErr18.Error()
												log.Println("PPO18", "No Records found for process application")
											}
										} else {
											lErr19 := common.CustomError("Timing Closed for IPO")
											lRespRec.Status = common.ErrorCode
											lRespRec.ErrMsg = "PPO19" + lErr19.Error()
											log.Println("PPO19", "Timing Closed for IPO")
										}
									}
								}
							}
							// }
						}
					}
				}
			}
		}
		// Marshal the response values of lRespRec in lData
		lData, lErr20 := json.Marshal(lRespRec)
		if lErr20 != nil {
			log.Println("PPO20", lErr20)
			fmt.Fprintf(w, helpers.GetErrorString("PPO20", "Unable to getting response.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("PlaceOrder (-)", r.Method)
	}
}

/*
Pupose:This method is used to process the application request to exchange
Parameters:

	pExchangeReq,pMasterId,PClientId

Response:

	==========
	*On Sucess
	==========
	2,[1,2,3],[1,2,3],nil

	==========
	*On Error
	==========
	0,[],[],error

Author:Pavithra
Date: 12JUNE2023
*/
//added by naveen:add one parameter source to insert in bidtracking table
//func ProcessNseReq(pExchangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int) ([]nseipo.ExchangeRespStruct, error) {
func ProcessNseReq(pExchangeReq []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pMailId string, pBrokerId int, pSource string) ([]nseipo.ExchangeRespStruct, error) {
	log.Println("ProcessNseReq (+)")

	// for retruning value for this method
	var lRespArr []nseipo.ExchangeRespStruct

	//added by naveen:add one parameter source to insert in bidtracking table
	//lBidTrackArr, lErr1 := InsertBidTrack(pExchangeReq, pClientId, common.NSE, pBrokerId,pSource)
	lBidTrackArr, lErr1 := InsertBidTrack(pExchangeReq, pClientId, common.NSE, pBrokerId, pSource)
	if lErr1 != nil {
		log.Println("PPNR01", lErr1)
		return lRespArr, lErr1
	} else {
		// ------------------------>>>>>>>>>>>>>>>>>>
		// Call the ApplyIpo method to Process the Application Request.
		lResponse, lErr2 := exchangecall.ApplyNseIpo(pExchangeReq, pClientId, pBrokerId)
		if lErr2 != nil {
			log.Println("PPNR02", lErr2.Error())
			return lRespArr, lErr2
		} else {
			if lResponse == nil {
				return lRespArr, lErr2
			} else {
				lRespArr = lResponse
				// update the bid tracking table when application status success or failed
				lErr3 := UpdateBidTrack(lResponse, pClientId, lBidTrackArr, pExchangeReq, common.NSE, pBrokerId)
				if lErr3 != nil {
					log.Println("PPNR03", lErr3)
					return lRespArr, lErr3
				} else {
					for lRespIdx := 0; lRespIdx < len(lResponse); lRespIdx++ {
						// check whether the application status is success or not
						if lResponse[lRespIdx].Status == "success" {
							// if the status is success update the response in header and details table.
							lErr4 := InsertHeader(lResponse, pExchangeReq, pReqRec, pClientId, common.NSE, pMailId, pBrokerId)
							if lErr4 != nil {
								log.Println("PPNR04", lErr4)
								return lRespArr, lErr4
							}
						}
					}
				}
			}
		}
	}
	log.Println("ProcessNseReq (-)")
	return lRespArr, nil
}

/*
Pupose:This method inserting the order head values in order header table.
Parameters:

	pReqArr,pMasterId,PClientId

Response:

		==========
		*On Sucess
		==========
			{
			"symbol": "UPIEMP",
			"applicationNumber": "FT000069172957",
			"category": "IND",
			"clientName": "LAKSHMANAN ASHOK KUMAR",
			"depository": "CDSL",
			"dpId": "",
			"clientBenId": "1208030000262661",
			"nonASBA": false,
			"pan": "AGMPA8575C",
			"referenceNumber": "bab61fa65db7c901",
			"allotmentMode": "demat",
			"upiFlag": "Y",
			"upi": "test@kmbl",
			"bankCode": "null",
			"locationCode": "null",
			"timestamp": "",
			"bids": [
				{
					"activityType": "modify",
					"bidReferenceNumber": "123456",
					"quantity": 10,
					"atCutOff": true,
					"price": 10,
					"amount": 100,
					"remark": "FT000069172951",
					"lotSize": 10
	        	},
				{
					"activityType": "cancel",
					"bidReferenceNumber": "2023051900001001",
					"quantity": 10,
					"atCutOff": false,
					"price": 9,
					"amount": 90,
					"remark": "FT000069172951",
					"lotSize": 10
				}
	   		]
		},nil

		==========
		*On Error
		==========
		[],error

Author:Pavithra
Date: 12JUNE2023
*/
func InsertHeader(pRespArr []nseipo.ExchangeRespStruct, pReqArr []nseipo.ExchangeReqStruct, pReqRec OrderReqStruct, pClientId string, pExchange string, pMailId string, pBrokerId int) error {
	log.Println("InsertHeader (+)")

	//get the application no id in table
	var lHeaderId int
	//set cancel Flag as N
	lCancelFlag := "N"
	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PIH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		// Looping the pReqArr Array value to inserting the request datas in databse.
		for lpRespIdx := 0; lpRespIdx < len(pRespArr); lpRespIdx++ {

			// Changing the date format DD-MM-YYYY into YYYY-MM-DD of timestamp value.
			if pRespArr[lpRespIdx].TimeStamp != "" {
				lDate := strings.Split(pRespArr[lpRespIdx].TimeStamp, " ")
				lTimeStamp, _ := time.Parse("02-01-2006", lDate[0])
				pRespArr[lpRespIdx].TimeStamp = lTimeStamp.Format("2006-01-02") + " " + lDate[1]
			} else if pRespArr[lpRespIdx].TimeStamp == "" {
				lTime := time.Now()
				pRespArr[lpRespIdx].TimeStamp = lTime.Format("2006-01-02 15:04:05")
			}

			//  call the getAppNoId method, To get Application No Id from the database
			lAppNoId, lErr2 := getAppNoId(pReqArr[lpRespIdx].ApplicationNo, pClientId)
			if lErr2 != nil {
				log.Println("PIH02", lErr2)
				return lErr2
			} else {
				lHeaderId = lAppNoId
				//if appliaction no is not present in the database,insert the new application details in orderheader table
				if lAppNoId == 0 {
					lSqlString := `insert into a_ipo_order_header (MasterId,Symbol,applicationNo,category,clientName,
						depository,dpId,clientBenId,nonASBA,chequeNo ,pan,referenceNumber,allotmentMode,upiFlag,upi,bankCode,
						locationCode,bankAccount ,ifsc ,subBrokerCode ,time_Stamp,status ,dpVerStatusFlag ,dpVerFailCode ,
						dpVerReason ,upiPaymentStatusFlag ,upiAmtBlocked ,reasonCode ,reason ,clientId,CreatedBy,CreatedDate,cancelFlag,Exchange,ClientEmail,
						UpdatedBy,UpdatedDate,brokerId) 
						values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?,?,?,now(),?)`

					lInsertedHeaderId, lErr3 := lDb.Exec(lSqlString, pReqRec.MasterId, pRespArr[lpRespIdx].Symbol, pRespArr[lpRespIdx].ApplicationNo,
						pRespArr[lpRespIdx].Category, pRespArr[lpRespIdx].ClientName, pRespArr[lpRespIdx].Depository, pRespArr[lpRespIdx].DpId,
						pRespArr[lpRespIdx].ClientBenId, pRespArr[lpRespIdx].NonASBA, pRespArr[lpRespIdx].ChequeNo, pRespArr[lpRespIdx].Pan,
						pRespArr[lpRespIdx].ReferenceNo, pRespArr[lpRespIdx].AllotmentMode, pRespArr[lpRespIdx].UpiFlag, pRespArr[lpRespIdx].Upi,
						pRespArr[lpRespIdx].BankCode, pRespArr[lpRespIdx].LocationCode, pRespArr[lpRespIdx].BankAccount, pRespArr[lpRespIdx].IFSC,
						pRespArr[lpRespIdx].SubBrokerCode, pRespArr[lpRespIdx].TimeStamp, pRespArr[lpRespIdx].Status, pRespArr[lpRespIdx].DpVerStatusFlag,
						pRespArr[lpRespIdx].DpVerFailCode, pRespArr[lpRespIdx].DpVerReason, pRespArr[lpRespIdx].UpiPaymentStatusFlag,
						pRespArr[lpRespIdx].UpiAmtBlocked, pRespArr[lpRespIdx].ReasonCode, pRespArr[lpRespIdx].Reason, pClientId, pClientId, lCancelFlag, pExchange, pMailId, pClientId, pBrokerId)
					if lErr3 != nil {
						log.Println("PIH03", lErr3)
						return lErr3
					} else {
						// get lastinserted id in lReturnId and converet them into int ,store it in lHeaderId
						lReturnId, _ := lInsertedHeaderId.LastInsertId()
						lHeaderId = int(lReturnId)

						// call InsertDetails method to inserting the order details in order details table
						lErr4 := InsertDetails(pRespArr[lpRespIdx].Bids, pReqArr[lpRespIdx].Bids, lHeaderId, pClientId)
						if lErr4 != nil {
							log.Println("PIH04", lErr4)
							return lErr4
						}
					}
				} else if len(pRespArr[lpRespIdx].Bids) == 0 && pRespArr[lpRespIdx].Status == "success" {
					lCancelFlag = "Y"
					lSqlString := `update a_ipo_order_header h
									set h.Symbol = ?,h.category = ?,h.clientName = ?,h.depository = ?,
									h.dpId = ?,h.clientBenId = ?,h.nonASBA = ?,h.chequeNo = ?,h.pan = ?,h.referenceNumber = ?,
									h.allotmentMode = ?,h.upiFlag = ?,h.upi = ?,h.bankCode = ?,h.locationCode = ?,
									h.bankAccount = ?,h.ifsc = ?,h.subBrokerCode = ?,h.time_Stamp =?,
									h.status = ?,h.dpVerStatusFlag = ?,h.dpVerFailCode = ?,h.dpVerReason = ?,
									h.upiPaymentStatusFlag = ?,h.upiAmtBlocked = ?,h.reasonCode = ?,h.reason = ?,
									h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag = ?
									where h.applicationNo = ? 
									and h.clientId = ?
									and h.Id = ?
									and h.Exchange = ?
									and h.brokerId = ? `

					_, lErr5 := lDb.Exec(lSqlString, pRespArr[lpRespIdx].Symbol, pRespArr[lpRespIdx].Category, pRespArr[lpRespIdx].ClientName,
						pRespArr[lpRespIdx].Depository, pRespArr[lpRespIdx].DpId, pRespArr[lpRespIdx].ClientBenId, pRespArr[lpRespIdx].NonASBA,
						pRespArr[lpRespIdx].ChequeNo, pRespArr[lpRespIdx].Pan, pRespArr[lpRespIdx].ReferenceNo, pRespArr[lpRespIdx].AllotmentMode,
						pRespArr[lpRespIdx].UpiFlag, pRespArr[lpRespIdx].Upi, pRespArr[lpRespIdx].BankCode, pRespArr[lpRespIdx].LocationCode,
						pRespArr[lpRespIdx].BankAccount, pRespArr[lpRespIdx].IFSC, pRespArr[lpRespIdx].SubBrokerCode, pRespArr[lpRespIdx].TimeStamp,
						pRespArr[lpRespIdx].Status, pRespArr[lpRespIdx].DpVerStatusFlag, pRespArr[lpRespIdx].DpVerFailCode, pRespArr[lpRespIdx].DpVerReason,
						pRespArr[lpRespIdx].UpiPaymentStatusFlag, pRespArr[lpRespIdx].UpiAmtBlocked, pRespArr[lpRespIdx].ReasonCode, pRespArr[lpRespIdx].Reason,
						pClientId, lCancelFlag, pRespArr[lpRespIdx].ApplicationNo, pClientId, lAppNoId, pExchange, pBrokerId)
					if lErr5 != nil {
						log.Println("PIH05", lErr5)
						return lErr5
					} else {
						// call InsertDetails method to inserting the order details in order details table
						lErr6 := InsertDetails(pRespArr[lpRespIdx].Bids, pReqArr[lpRespIdx].Bids, lHeaderId, pClientId)
						if lErr6 != nil {
							log.Println("PIH06", lErr6)
							return lErr6
						}
					}
				} else {
					log.Println("else")
					lSqlString := `update a_ipo_order_header h
									set h.Symbol = ?,h.category = ?,h.clientName = ?,h.depository = ?,
									h.dpId = ?,h.clientBenId = ?,h.nonASBA = ?,h.chequeNo = ?,h.pan = ?,h.referenceNumber = ?,
									h.allotmentMode = ?,h.upiFlag = ?,h.upi = ?,h.bankCode = ?,h.locationCode = ?,
									h.bankAccount = ?,h.ifsc = ?,h.subBrokerCode = ?,h.time_Stamp =?,
									h.status = ?,h.dpVerStatusFlag = ?,h.dpVerFailCode = ?,h.dpVerReason = ?,
									h.upiPaymentStatusFlag = ?,h.upiAmtBlocked = ?,h.reasonCode = ?,h.reason = ?,
									h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag = ?
									where h.applicationNo = ? 
									and h.clientId = ?
									and h.Id = ? 
									and h.Exchange = ?
									and h.brokerId = ?`

					_, lErr5 := lDb.Exec(lSqlString, pRespArr[lpRespIdx].Symbol, pRespArr[lpRespIdx].Category, pRespArr[lpRespIdx].ClientName,
						pRespArr[lpRespIdx].Depository, pRespArr[lpRespIdx].DpId, pRespArr[lpRespIdx].ClientBenId, pRespArr[lpRespIdx].NonASBA,
						pRespArr[lpRespIdx].ChequeNo, pRespArr[lpRespIdx].Pan, pRespArr[lpRespIdx].ReferenceNo, pRespArr[lpRespIdx].AllotmentMode,
						pRespArr[lpRespIdx].UpiFlag, pRespArr[lpRespIdx].Upi, pRespArr[lpRespIdx].BankCode, pRespArr[lpRespIdx].LocationCode,
						pRespArr[lpRespIdx].BankAccount, pRespArr[lpRespIdx].IFSC, pRespArr[lpRespIdx].SubBrokerCode, pRespArr[lpRespIdx].TimeStamp,
						pRespArr[lpRespIdx].Status, pRespArr[lpRespIdx].DpVerStatusFlag, pRespArr[lpRespIdx].DpVerFailCode, pRespArr[lpRespIdx].DpVerReason,
						pRespArr[lpRespIdx].UpiPaymentStatusFlag, pRespArr[lpRespIdx].UpiAmtBlocked, pRespArr[lpRespIdx].ReasonCode, pRespArr[lpRespIdx].Reason,
						pClientId, lCancelFlag, pRespArr[lpRespIdx].ApplicationNo, pClientId, lAppNoId, pExchange, pBrokerId)
					if lErr5 != nil {
						log.Println("PIH07", lErr5)
						return lErr5
					} else {
						// call InsertDetails method to inserting the order details in order details table
						lErr6 := InsertDetails(pRespArr[lpRespIdx].Bids, pReqArr[lpRespIdx].Bids, lHeaderId, pClientId)
						if lErr6 != nil {
							log.Println("PIH08", lErr6)
							return lErr6
						}
					}
				}
			}
		}
	}
	log.Println("InsertHeader (-)")
	return nil
}

/*
Pupose:This method inserting the bid details in order detail table.
Parameters:

	pReqBidArr,pHeaderId,PClientId

Response:

	==========
	*On Sucess
	==========
	[1,2,3],[1,2,3],nil

	==========
	*On Error
	==========
	[],[],error

Author:Pavithra
Date: 12JUNE2023
*/
func InsertDetails(pRespBidArr []nseipo.ResponseBidStruct, pReqBidArr []nseipo.RequestBidStruct, pHeaderId int, pClientId string) error {
	log.Println("InsertDetails (+)")

	// To store the Inserted datas Id of Order Detail Table in database
	// var lDetailIdArr []int
	// To store the Inserted datas Id of Bid Tracking Table in database
	// var lBidTrackIdArr []int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PID01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		// Looping the pReqBidArr array to insert the Request bid details in database
		for lReqBidIdx := 0; lReqBidIdx < len(pReqBidArr); lReqBidIdx++ {
			if pReqBidArr[lReqBidIdx].ActivityType == "cancel" {

				// ! ----------------
				lSqlString2 := `update a_ipo_orderdetails d
									set d.activityType = ?,d.status = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
									where d.headerId = ? 
									and d.bidReferenceNo = ?`
				// log.Println("lSqlString2", lSqlString2)

				_, lErr2 := lDb.Exec(lSqlString2, pReqBidArr[lReqBidIdx].ActivityType, "success",
					pClientId, pHeaderId, pReqBidArr[lReqBidIdx].BidReferenceNo)
				if lErr2 != nil {
					log.Println("Error updating into database (details)")
					log.Println("PID02", lErr2)
					return lErr2
				}
			} else {
				for lRespBidIdx := 0; lRespBidIdx < len(pRespBidArr); lRespBidIdx++ {
					log.Println("Inside Second For")

					// Check whether the requested bid is activity type is new,if it is new insert the bid details in order detail table
					if pReqBidArr[lReqBidIdx].ActivityType == "new" {
						if pRespBidArr[lRespBidIdx].ActivityType == pReqBidArr[lReqBidIdx].ActivityType &&
							pRespBidArr[lRespBidIdx].Remark == pReqBidArr[lReqBidIdx].Remark {
							lSqlString1 := `insert into a_ipo_orderdetails(headerId,activityType,bidReferenceNo,req_quantity,
							atCutOff,req_price,req_amount,remark,series,lotSize,resp_quantity ,resp_price ,resp_amount ,status ,
							reasonCode ,reason ,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`
							// log.Println("lSqlString1", lSqlString1)

							_, lErr3 := lDb.Exec(lSqlString1, pHeaderId, pRespBidArr[lRespBidIdx].ActivityType, pRespBidArr[lRespBidIdx].BidReferenceNo,
								pReqBidArr[lReqBidIdx].Quantity, pRespBidArr[lRespBidIdx].AtCutOff, pReqBidArr[lReqBidIdx].Price, pReqBidArr[lReqBidIdx].Amount,
								pRespBidArr[lRespBidIdx].Remark, pRespBidArr[lRespBidIdx].Series, pReqBidArr[lReqBidIdx].LotSize, pRespBidArr[lRespBidIdx].Quantity,
								pRespBidArr[lRespBidIdx].Price, pRespBidArr[lRespBidIdx].Amount, pRespBidArr[lRespBidIdx].Status, pRespBidArr[lRespBidIdx].ReasonCode,
								pRespBidArr[lRespBidIdx].Reason, pClientId, pClientId)
							if lErr3 != nil {
								log.Println("PID03", lErr3)
								return lErr3
							} else {
								log.Println("Details Inserted Successfully", lRespBidIdx)

							}
						}
						// Check whether the requested bid is activity type is modify or cancel,then update the bid details in order detail table
					} else if pReqBidArr[lReqBidIdx].ActivityType == "modify" {

						//Check whether the requested bid and Response bid activity type same (modify or cancel)
						if pRespBidArr[lRespBidIdx].ActivityType == pReqBidArr[lReqBidIdx].ActivityType && pRespBidArr[lRespBidIdx].BidReferenceNo == pReqBidArr[lReqBidIdx].BidReferenceNo {
							// ! ----------------
							lSqlString2 := `update a_ipo_orderdetails d
										set d.activityType = ?,d.req_quantity = ?,d.atCutOff = ?,d.req_price = ?,d.req_amount = ?,
										d.remark = ?,d.series = ?,d.lotSize = ?,d.resp_quantity = ?,d.resp_price = ?,d.resp_amount = ?,
										d.status = ?,d.reasonCode = ?,d.reason = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
										where d.headerId = ? 
										and d.bidReferenceNo = ?`
							// log.Println("lSqlString2", lSqlString2)
							_, lErr4 := lDb.Exec(lSqlString2, pRespBidArr[lRespBidIdx].ActivityType, pReqBidArr[lReqBidIdx].Quantity,
								pRespBidArr[lRespBidIdx].AtCutOff, pReqBidArr[lReqBidIdx].Price, pReqBidArr[lReqBidIdx].Amount,
								pRespBidArr[lRespBidIdx].Remark, pRespBidArr[lRespBidIdx].Series, pReqBidArr[lReqBidIdx].LotSize,
								pRespBidArr[lRespBidIdx].Quantity, pRespBidArr[lRespBidIdx].Price, pRespBidArr[lRespBidIdx].Amount,
								pRespBidArr[lRespBidIdx].Status, pRespBidArr[lRespBidIdx].ReasonCode, pRespBidArr[lRespBidIdx].Reason,
								pClientId, pHeaderId, pRespBidArr[lRespBidIdx].BidReferenceNo)
							if lErr4 != nil {
								log.Println("PID04", lErr4)
								return lErr4
							}
						}
					}
				}
				// break
			}
			// break
		}
	}
	log.Println("InsertDetails (-)")
	return nil
}

/*
Pupose:This method inserting the bid details in bid tracking table.
Parameters:

	pReqBidRec,pHeaderId,PClientId

Response:

	==========
	*On Sucess
	==========
	2,nil

	==========
	*On Error
	==========
	0,error

Author:Pavithra
Date: 12JUNE2023
*/
//added by naveen:add one parameter source to insert in bidtracking table
//func InsertBidTrack(pReqRec []nseipo.ExchangeReqStruct, pClientId string, pExchange string, pBrokerId int)
func InsertBidTrack(pReqRec []nseipo.ExchangeReqStruct, pClientId string, pExchange string, pBrokerId int, pSource string) ([]BidTrack, error) {
	log.Println("InsertBidTrack (+)")
	// To store the Inserted Bid Track Details ID in database.
	var lBidTrackRec BidTrack
	var lBidTrackArr []BidTrack
	log.Println("sourc ibd", pSource)
	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PIB01", lErr1)
		return lBidTrackArr, lErr1
	} else {
		defer lDb.Close()

		for lpReqIdx := 0; lpReqIdx < len(pReqRec); lpReqIdx++ {
			lBid := pReqRec[lpReqIdx].Bids
			for lBidReqIdx := 0; lBidReqIdx < len(lBid); lBidReqIdx++ {

				lSqlString := `insert into a_bidtracking_table(applicationNo,bidRefNo,activityType,quantity,price,clientId,brokerId,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,Exchange,source)
							values(?,?,?,?,?,?,?,?,now(),?,now(),?,?)`

				lInsertedBidId, lErr2 := lDb.Exec(lSqlString, pReqRec[lpReqIdx].ApplicationNo, lBid[lBidReqIdx].BidReferenceNo, lBid[lBidReqIdx].ActivityType,
					lBid[lBidReqIdx].Quantity, lBid[lBidReqIdx].Price, pClientId, pBrokerId, pClientId, pClientId, pExchange, pSource)

				log.Println("lInsertedBidId", lInsertedBidId)
				if lErr2 != nil {
					log.Println("Error inserting into database (bidtrack)")
					log.Println("PIB02", lErr2)
					return lBidTrackArr, lErr2
				} else {
					// get lastinserted id in lReturnId and converet them into int ,store it in lBidTrackId
					lReturnId, _ := lInsertedBidId.LastInsertId()

					lBidTrackRec.Id = int(lReturnId)
					lBidTrackRec.BidRefNo = lBid[lBidReqIdx].BidReferenceNo
					lBidTrackArr = append(lBidTrackArr, lBidTrackRec)
					// commented by pavithra
					// log.Println("lBidTrackRec.BidRefNo ", lBidTrackRec.BidRefNo)

				}
			}
		}
	}
	log.Println("InsertBidTrack (-)")
	return lBidTrackArr, nil
}

/*
Pupose:This method updating the bid details in bid tracking table.
Parameters:

	pRespBidRec,lClientId,pBidTrackIdRec

Response:

	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:Pavithra
Date: 12JUNE2023
*/
func UpdateBidTrack(pRespRec []nseipo.ExchangeRespStruct, pClientId string, pBidTrackArr []BidTrack, pExchangeReq []nseipo.ExchangeReqStruct, pExchange string, pBrokerId int) error {
	log.Println("UpdateBidTrack (+)")

	var lBidArr []nseipo.RequestBidStruct
	var lBidResp []nseipo.ResponseBidStruct
	var lStatus string

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PUB01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		// commented by pavithra
		// log.Println("req", pExchangeReq)
		// log.Println("resp", pRespRec)
		// log.Println("bidtrack", pBidTrackArr)

		for lpRespIdx := 0; lpRespIdx < len(pRespRec); lpRespIdx++ {
			lBidResp = pRespRec[lpRespIdx].Bids
			lStatus = pRespRec[lpRespIdx].Status
		}
		for lReqIdx := 0; lReqIdx < len(pExchangeReq); lReqIdx++ {
			lBidArr = pExchangeReq[lReqIdx].Bids
		}

		for pBidTrackIdx := 0; pBidTrackIdx < len(pBidTrackArr); pBidTrackIdx++ {
			// commented by pavithra
			// log.Println("bidtrackarr", pBidTrackIdx)
			for lReqBidIdx := 0; lReqBidIdx < len(lBidArr); lReqBidIdx++ {
				// commented by pavithra
				// log.Println("reqbididx", lReqBidIdx)
				if lBidArr[lReqBidIdx].ActivityType == "cancel" {
					// lRemarkValue, _ := strconv.Atoi(lBid[pBidTrackIdx].Remark)
					if pBidTrackArr[pBidTrackIdx].BidRefNo == lBidArr[lReqBidIdx].BidReferenceNo {
						lSqlString := `update a_bidtracking_table b
										set b.bidRefNo = ?,b.status = ?,b.applicationStatus = ?,b.UpdatedBy = ?,b.UpdatedDate = now() 
										where b.Id = ? and b.clientId = ? and b.brokerId = ?`

						// log.Println("lSqlString", lSqlString)

						_, lErr2 := lDb.Exec(lSqlString, lBidArr[lReqBidIdx].BidReferenceNo, lStatus, lStatus,
							pClientId, pBidTrackArr[pBidTrackIdx].Id, pClientId, pBrokerId)
						if lErr2 != nil {
							log.Println("PUB02", lErr2)
							return lErr2
						} else {
							log.Println("cancel updated")
						}
					}
				}

				if len(lBidResp) > 0 {
					// commented by pavithra
					// log.Println("greater than 0")
					for lBidRespIdx := 0; lBidRespIdx <= pBidTrackIdx; lBidRespIdx++ {

						// commented by pavithra
						// log.Println("lBidResp[pBidTrackIdx].Remark", lBidResp[pBidTrackIdx].Remark)
						lRemarkValue, _ := strconv.Atoi(lBidResp[pBidTrackIdx].Remark)
						// commented by pavithra
						// log.Println("outerloop index", lBidRespIdx)
						if lBidArr[lReqBidIdx].ActivityType == "new" {
							if pExchange == common.BSE {
								// commented by pavithra
								// log.Println("inside update bidtrack bse")

								// for pBidTrackIdx := 0; pBidTrackIdx < len(pBidTrackArr); pBidTrackIdx++ {
								// lRemarkValue, _ := strconv.Atoi(lBidResp[pBidTrackIdx].Remark)
								// commented by pavithra
								// log.Println("lBidResp[pBidTrackIdx].BidReferenceNo", lBidResp[lBidRespIdx].BidReferenceNo)
								// log.Println("lRemarkValue", lRemarkValue, lBidResp[lBidRespIdx].Remark)

								// if lBidArr[i].BidReferenceNo == pBidTrackArr[pBidTrackIdx].BidRefNo {
								// }
								if pBidTrackArr[pBidTrackIdx].BidRefNo == lRemarkValue {
									lSqlString := `update a_bidtracking_table b
												set b.bidRefNo = ?,b.status = ?,b.applicationStatus = ?,b.UpdatedBy = ?,b.UpdatedDate = now() 
												where b.Id = ? and b.clientId = ? and b.brokerId = ?`

									// log.Println("lSqlString", lSqlString)

									_, lErr2 := lDb.Exec(lSqlString, lBidResp[lBidRespIdx].BidReferenceNo, lBidResp[lBidRespIdx].Status, lStatus,
										pClientId, pBidTrackArr[pBidTrackIdx].Id, pClientId, pBrokerId)
									if lErr2 != nil {
										log.Println("PUB02", lErr2)
										return lErr2
									} else {
										log.Println("if updated in bse")
									}
								}
							} else {
								// commented by pavithra
								// log.Println("inside update bidtrack nse----------------")

								// for pBidTrackIdx := 0; pBidTrackIdx < len(pBidTrackArr); pBidTrackIdx++ {
								// lRemarkValue, _ := strconv.Atoi(lBidResp[pBidTrackIdx].Remark)
								// commented by pavithra
								// log.Println("lBidResp[pBidTrackIdx].BidReferenceNo", lBidResp[lBidRespIdx].BidReferenceNo)
								// log.Println("lRemarkValue", lRemarkValue, pBidTrackArr[pBidTrackIdx].BidRefNo)
								if pBidTrackArr[pBidTrackIdx].BidRefNo == lRemarkValue {
									log.Println("inside nse new ", pBidTrackArr[pBidTrackIdx].BidRefNo, "=", lRemarkValue)
									lSqlString := `update a_bidtracking_table b
												set b.bidRefNo = ?,b.status = ?,b.applicationStatus = ?,b.UpdatedBy = ?,b.UpdatedDate = now()
												where b.Id = ? and b.clientId = ? and b.brokerId = ?`

									_, lErr2 := lDb.Exec(lSqlString, lBidResp[lBidRespIdx].BidReferenceNo, lBidResp[lBidRespIdx].Status, lStatus,
										pClientId, pBidTrackArr[pBidTrackIdx].Id, pClientId, pBrokerId)
									if lErr2 != nil {
										log.Println("PUB02", lErr2)
										return lErr2
									} else {
										log.Println("else updated in nse")
										// break
									}
								}
							}

						} else {
							// commented by pavithra
							// log.Println("inside update bidtrack else")
							// for lBidRespIdx := 0; lBidRespIdx < len(lBid); lBidRespIdx++ {
							// for pBidTrackIdx := 0; pBidTrackIdx < len(pBidTrackArr); pBidTrackIdx++ {
							// commented by pavithra
							// log.Println(pBidTrackArr[pBidTrackIdx].BidRefNo, "=", lBidResp[lBidRespIdx].BidReferenceNo)
							if pBidTrackArr[pBidTrackIdx].BidRefNo == lBidResp[lBidRespIdx].BidReferenceNo {
								lSqlString := `update a_bidtracking_table b
								set b.bidRefNo = ?,b.status = ?,b.applicationStatus = ?,b.UpdatedBy = ?,b.UpdatedDate = now() 
								where b.Id = ? and b.clientId = ? and b.brokerId = ?`

								_, lErr2 := lDb.Exec(lSqlString, lBidResp[lBidRespIdx].BidReferenceNo, lBidResp[lBidRespIdx].Status, lStatus,
									pClientId, pBidTrackArr[pBidTrackIdx].Id, pClientId, pBrokerId)
								if lErr2 != nil {
									log.Println("PUB02", lErr2)
									return lErr2
								} else {
									log.Println("else updated")
								}
								// }
							}
						}
					}
				}
			}
		}
	}
	log.Println("UpdateBidTrack (-)")
	return nil
}

/*
Purpose:This method updating the Status As Pending in details table.
Parameters:

	pRespArr,lClientId,pHeaderId,pDetailIdArr

Response:

	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:Pavithra
Date: 12JUNE2023
*/

//added by naveen: add  one parameter psource for  insert in ipo order header
// func UpdatePending(pReqArr []nseipo.ExchangeReqStruct, pClientId string, pReqRec OrderReqStruct, pExchange string, pMailId string, pBrokerId int) error {
func UpdatePending(pReqArr []nseipo.ExchangeReqStruct, pClientId string, pReqRec OrderReqStruct, pExchange string, pMailId string, pBrokerId int, pSource string) error {
	log.Println("UpdatePending (+)")

	// This Variable is used to insert null value in database
	var lTimeStampNull sql.NullString

	// for Set the cancel flag as N
	lCancelFlag := "N"
	// for set the status as pending
	lPendingValue := "pending"
	// establish a database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PUP01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// loop the req array to insert in database
		for lpReqArrIdx := 0; lpReqArrIdx < len(pReqArr); lpReqArrIdx++ {

			lSqlString := `insert into a_ipo_order_header (MasterId,Symbol,applicationNo,category,clientName,
			depository,dpId,clientBenId,nonASBA,pan,referenceNumber,allotmentMode,upiFlag,upi,bankCode,
			locationCode,time_Stamp,clientId,CreatedBy,CreatedDate,cancelFlag,status,Exchange,ClientEmail,UpdatedBy,UpdatedDate,brokerId,source) 
			values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?,?,?,?,now(),?,?)`

			lInsertedHeaderId, lErr2 := lDb.Exec(lSqlString, pReqRec.MasterId, pReqArr[lpReqArrIdx].Symbol, pReqArr[lpReqArrIdx].ApplicationNo, pReqArr[lpReqArrIdx].Category,
				pReqArr[lpReqArrIdx].ClientName, pReqArr[lpReqArrIdx].Depository, pReqArr[lpReqArrIdx].DpId, pReqArr[lpReqArrIdx].ClientBenId, pReqArr[lpReqArrIdx].NonASBA,
				pReqArr[lpReqArrIdx].Pan, pReqArr[lpReqArrIdx].ReferenceNo, pReqArr[lpReqArrIdx].AllotmentMode, pReqArr[lpReqArrIdx].UpiFlag, pReqArr[lpReqArrIdx].Upi,
				pReqArr[lpReqArrIdx].BankCode, pReqArr[lpReqArrIdx].LocationCode, lTimeStampNull, pClientId, pClientId, lCancelFlag, lPendingValue, pExchange, pMailId, pClientId, pBrokerId, pSource)
			if lErr2 != nil {
				log.Println("Error updating pending status into database (header)")
				log.Println("PUP02", lErr2)
				return lErr2
			} else {
				// get lastinserted id in lReturnId and converet them into int ,store it in lHeaderId
				lReturnId, _ := lInsertedHeaderId.LastInsertId()
				lHeaderId := int(lReturnId)

				// Looping the Bid Details Array to Update the Application Response Values in Order Details Table
				lBidArr := pReqArr[lpReqArrIdx].Bids
				for lBidArrIdx := 0; lBidArrIdx < len(lBidArr); lBidArrIdx++ {

					// Check whether the requested bid is activity type is new,if it is new insert the bid details in order detail table
					if lBidArr[lBidArrIdx].ActivityType == "new" {
						lSqlString := `insert into a_ipo_orderdetails(headerId,activityType,bidReferenceNo,req_quantity,
							atCutOff,req_price,req_amount,remark,series,lotSize,CreatedBy,CreatedDate,status,UpdatedBy,UpdatedDate)
							values (?,?,?,?,?,?,?,?,?,?,?,now(),?,?,now())`

						_, lErr3 := lDb.Exec(lSqlString, lHeaderId, lBidArr[lBidArrIdx].ActivityType,
							lBidArr[lBidArrIdx].BidReferenceNo, lBidArr[lBidArrIdx].Quantity, lBidArr[lBidArrIdx].AtCutOff,
							lBidArr[lBidArrIdx].Price, lBidArr[lBidArrIdx].Amount, lBidArr[lBidArrIdx].Remark,
							lBidArr[lBidArrIdx].Series, lBidArr[lBidArrIdx].LotSize, pClientId, lPendingValue, pClientId)
						if lErr3 != nil {
							log.Println("PUP03", lErr3)
							return lErr3
						} else {
							log.Println("Details Inserted Successfully", lBidArrIdx)
							lSqlString := `insert into a_bidtracking_table (applicationNo,bidRefNo,activityType,quantity,
								price,clientId,applicationStatus,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,brokerId)
								values(?,?,?,?,?,?,?,?,?,now(),?,now(),?)`

							_, lErr4 := lDb.Exec(lSqlString, pReqArr[lpReqArrIdx].ApplicationNo, lBidArr[lBidArrIdx].BidReferenceNo,
								lBidArr[lBidArrIdx].ActivityType, lBidArr[lBidArrIdx].Quantity, lBidArr[lBidArrIdx].Price, pClientId,
								lPendingValue, pSource, pClientId, pClientId, pBrokerId)
							if lErr4 != nil {
								log.Println("PUP04", lErr4)
								return lErr4
							}
						}
						// Check whether the requested bid is activity type is modify or cancel,then update the bid details in order detail table
					} else if lBidArr[lBidArrIdx].ActivityType == "modify" || lBidArr[lBidArrIdx].ActivityType == "cancel" {

						SqlString := `update a_ipo_orderdetails d
						set d.activityType = ?,d.series = ?,d.req_quantity = ?,d.req_price = ?,
						d.req_amount = ?,d.atCutOff = ?,d.remark = ?,d.UpdatedBy = ?,d.UpdatedDate = now(),d.status = ?
						where d.headerId = ? and d.bidReferenceNo = ? `

						_, lErr5 := lDb.Exec(SqlString, lBidArr[lBidArrIdx].ActivityType, lBidArr[lBidArrIdx].Series,
							lBidArr[lBidArrIdx].Quantity, lBidArr[lBidArrIdx].Price, lBidArr[lBidArrIdx].Amount,
							lBidArr[lBidArrIdx].AtCutOff, lBidArr[lBidArrIdx].Remark, pClientId, lHeaderId, lBidArr[lBidArrIdx].BidReferenceNo, lPendingValue)
						if lErr5 != nil {
							log.Println("PUP05", lErr5)
							return lErr5
						} else {
							lSqlString := `insert into a_bidtracking_table (applicationNo,bidRefNo,activityType,quantity,
								price,clientId,applicationStatus,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,brokerId)
								values(?,?,?,?,?,?,?,?,?,now(),?,now(),?)`

							_, lErr6 := lDb.Exec(lSqlString, pReqArr[lpReqArrIdx].ApplicationNo, lBidArr[lBidArrIdx].BidReferenceNo,
								lBidArr[lBidArrIdx].ActivityType, lBidArr[lBidArrIdx].Quantity, lBidArr[lBidArrIdx].Price, pClientId,
								lPendingValue, pSource, pClientId, pClientId, pBrokerId)
							if lErr6 != nil {
								log.Println("PUP06", lErr6)
								return lErr6
							}
						}
					}
				}
			}
		}
	}
	log.Println("UpdatePending (-)")
	return nil
}

/*
Purpose:This method is used to construct the exchange header values.
Parameters:

	lReqRec,lClientId

Response:

		==========
		*On Sucess
		==========
	 	[
			{
				"symbol":"JDIAL",
				"applicationNumber":"FT000069130109",
				"category":"IND",
				"clientName":"LAKSHMANAN ASHOK KUMAR",
				"depository":"CDSL",
				"dpId":"",
				"clientBenId":"1208030000262661",
				"nonASBA":false,
				"pan":"AGMPA8575C",
				"referenceNumber":"2c2e506472cbb2f9",
				"allotmentMode":"demat",
				"upiFlag":"Y",
				"upi":"test@kmbl",
				"bankCode":"",
				"locationCode":"",
				"bankAccount":"",
				"ifsc":"",
				"subBrokerCode":"",
				"timestamp":"0000-00-00",
				"bids":[
					{
						"activityType":"new",
						"bidReferenceNumber":"123",
						"series":"",
						"quantity":19,
						"atCutOff":true,
						"price":755,
						"amount":14345,
						"remark":"FT000069130109",
						"lotSize":19
					}
				]
			}
		]
		==========
		*On Error
		==========
		[],error

Author:Pavithra
Date: 12JUNE2023
*/
func ConstructReqStruct(pReqRec OrderReqStruct, lClientId string) ([]nseipo.ExchangeReqStruct, string, error) {
	log.Println("ConstructReqStruct (+)")

	// create an instance of lReqRec type nse.ExchangeReqStruct.
	var lReqRec nseipo.ExchangeReqStruct
	// create an instance of lReqArr Array type nse.ExchangeReqStruct.
	var lReqArr []nseipo.ExchangeReqStruct
	var lEmailId string

	// call the getDpId method to get the lDpId and lClientName
	lClientdetails, lErr := clientDetail.GetClientEmailId(lClientId)
	if lErr != nil {
		log.Println("PCRS01", lErr)
		return nil, lEmailId, lErr
	} else {
		// lReqRec.Pan = lClientdetails.Pan_no
		// lReqRec.ClientName = lClientdetails.Client_dp_name
		// lReqRec.ClientBenId = lClientdetails.Client_dp_code
		// lEmailId = lClientdetails.EmailId

		lEmailId = lClientdetails.EmailId
		if lClientdetails.Pan_no != "" {
			lReqRec.Pan = lClientdetails.Pan_no
		} else {
			return nil, lEmailId, common.CustomError("Client PAN No not found")
		}
		if lClientdetails.Client_dp_name != "" {
			lReqRec.ClientName = lClientdetails.Client_dp_name
		} else {
			return nil, lEmailId, common.CustomError("Client Name not found")
		}
		if lClientdetails.Client_dp_code != "" {
			lReqRec.ClientBenId = lClientdetails.Client_dp_code
		} else {
			return nil, lEmailId, common.CustomError("Client DP ID not found")
		}

	}

	// call the getPanNO method to get the lPanNo
	// lPanNo, lErr := clientDetail.GetClientEmailId(lClientId)
	// if lErr != nil {
	// 	log.Println("PCRS02", lErr)
	// 	return nil, lErr
	// }

	// call the getApplication method to get the lAppNo
	lAppNo, lRefNo, lErr := getApplicationNo(pReqRec.Symbol, pReqRec.ApplicationNo, lClientId)
	if lErr != nil {
		log.Println("PCRS03", lErr)
		return nil, lEmailId, lErr
	} else {
		// If lAppNo is nil,generate new application no or else pass the lAppNo
		if lAppNo == "nil" {
			lTime := time.Now()
			lPresentTime := lTime.Format("150405")
			lReqRec.ApplicationNo = lClientId + lPresentTime

			lByte := make([]byte, 8)
			rand.Read(lByte)
			lReqRec.ReferenceNo = fmt.Sprintf("%x", lByte)
		} else {
			lReqRec.ApplicationNo = lAppNo
			lReqRec.ReferenceNo = lRefNo
		}
	}

	lReqRec.Symbol = pReqRec.Symbol
	lReqRec.Category = pReqRec.Category
	// lReqRec.ClientName = lClientdetails.Client_dp_name
	lReqRec.Depository = "CDSL"
	lReqRec.DpId = ""
	// lReqRec.ClientBenId = lClientdetails.Client_dp_code
	lReqRec.NonASBA = false
	// lReqRec.ChequeNo = ""
	// lReqRec.Pan = lClientdetails.Pan_no
	lReqRec.AllotmentMode = "demat"
	lReqRec.Upi = pReqRec.UpiId + pReqRec.UpiEndPoint
	lReqRec.UpiFlag = "Y"
	// lReqRec.BankCode = ""
	// lReqRec.LocationCode = ""
	// lReqRec.BankAccount = ""
	lReqRec.IFSC = ""
	lReqRec.SubBrokerCode = ""
	lReqRec.TimeStamp = ""

	//
	lReqRec.Bids = ConstructBidDetail(pReqRec.BidDetails)
	lReqArr = append(lReqArr, lReqRec)

	log.Println("lReqArr", lReqArr)

	log.Println("ConstructReqStruct (-)")
	return lReqArr, lEmailId, nil
}

/*
Pupose:This method is used to get the Application Number From the DataBase.
Parameters:

	pSymbol,pAppNo,PClientId

Response:

	==========
	*On Sucess
	==========
	FT000006912345,nil

	==========
	*On Error
	==========
	"",error

Author:Pavithra
Date: 12JUNE2023
*/
func getApplicationNo(pSymbol string, pAppNo string, lClientId string) (string, string, error) {
	log.Println("getApplicationNo (+)")
	// this variable is used to get the application no and reference no from the database
	var lAppNo, lRefNo string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PGA01", lErr1)
		return lAppNo, lRefNo, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select (case when count(1) > 0 then h.applicationNo else 'nil' end) AppNo,nvl(h.referenceNumber,'nil') 
						from a_ipo_order_header h
						where h.Symbol = ?
						and h.clientId = ?
						and h.applicationNo = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pSymbol, lClientId, pAppNo)
		if lErr2 != nil {
			log.Println("PGA02", lErr2)
			return lAppNo, lRefNo, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lAppNo, &lRefNo)
				if lErr3 != nil {
					log.Println("PGA03", lErr3)
					return lAppNo, lRefNo, lErr3
				}
			}
		}
	}
	log.Println("getApplicationNo (-)")
	return lAppNo, lRefNo, nil
}

/*
Pupose:This method is used to construct the exchange bid value Array
Parameters:

	pBidDetail,pAppNo

Response:

	==========
	*On Sucess
	==========
	[
		{
			"activityType":"new",
			"bidReferenceNumber":"123",
			"series":"",
			"quantity":19,
			"atCutOff":true,
			"price":755,
			"amount":14345,
			"remark":"FT000069130109",
			"lotSize":19
		}
	]

	==========
	*On Error
	==========

Author:Pavithra
Date: 12JUNE2023
*/
func ConstructBidDetail(pBidDetail []OrderBidStruct) []nseipo.RequestBidStruct {
	log.Println("ConstructBidDetail (+)")
	// create a instance of array type requestBidStruct
	var lBidArr []nseipo.RequestBidStruct
	// create a instance for requestBidStruct
	var lBidRec nseipo.RequestBidStruct

	lReqArr := pBidDetail
	// We use this loop to append the requestBidStruct inside requestBidStruct Array
	for lReqidx := 0; lReqidx < len(lReqArr); lReqidx++ {
		lBidRec.ActivityType = lReqArr[lReqidx].ActivityType

		// COMMENTED BY NITHISH BECAUSE THE RANDOM NUMBER METHOD
		// GIVES THE SAME FOR MORE THAN ONE BID

		if lReqArr[lReqidx].ActivityType == "new" {
			lRandomString, lErr1 := getSequenceNo()
			if lErr1 != nil {
				log.Println("IPCBD01", lErr1)
				return lBidArr
			} else {
				lBidRec.BidReferenceNo = lRandomString
				lBidRec.Remark = strconv.Itoa(lRandomString)
			}
		} else {
			lBidRec.BidReferenceNo, _ = strconv.Atoi(lReqArr[lReqidx].BidReferenceNo)
			lBidRec.Remark = lReqArr[lReqidx].BidReferenceNo
		}

		// commented by pavithra
		// ADDED BY NITHISH AND IT RECEIVES THE BIDREFNO GIVEN FROM THE UI SIDE HAS IT IS
		// lBidRec.BidReferenceNo, _ = strconv.Atoi(lReqArr[lReqidx].BidReferenceNo)
		// lBidRec.Remark = lReqArr[lReqidx].BidReferenceNo

		lBidRec.Series = ""
		lBidRec.Quantity = lReqArr[lReqidx].Quantity * lReqArr[lReqidx].LotSize
		lBidRec.AtCutOff = lReqArr[lReqidx].CutOff
		lBidRec.Price = lReqArr[lReqidx].Price
		lBidRec.Amount = lReqArr[lReqidx].Quantity * lReqArr[lReqidx].LotSize * lReqArr[lReqidx].Price

		lBidRec.LotSize = lReqArr[lReqidx].LotSize
		// append the bid value in lBidArr
		lBidArr = append(lBidArr, lBidRec)
	}
	log.Println("ConstructBidDetail (-)")
	return lBidArr
}

/*
Pupose:This method is used to get sequence no for unique bid ref no
Parameters:
nil
Response:

	==========
	*On Sucess
	==========
	123456789012,nil

	==========
	*On Error
	==========
	0,error

Author:Pavithra
Date: 29DEC2023
*/
func getSequenceNo() (int, error) {
	log.Println("getSequenceNo (+)")

	// this variables is used to get DpId and ClientName from the database.
	var lSequenceNo int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IPGSN01", lErr1)
		return lSequenceNo, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `SELECT NEXT VALUE FOR a_ipo_reference_s;`

		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("IPGSN02", lErr2)
			return lSequenceNo, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSequenceNo)
				if lErr3 != nil {
					log.Println("IPGSN03", lErr3)
					return lSequenceNo, lErr3
				}
			}
		}
	}
	log.Println("getSequenceNo (-)")
	return lSequenceNo, nil
}

/*
Pupose:This method is used to get the client BenId and ClientName for order input.
Parameters:

	PClientId

Response:

	==========
	*On Sucess
	==========
	123456789012,Lakshmanan Ashok Kumar,nil

	==========
	*On Error
	==========
	"","",error

Author:Pavithra
Date: 12JUNE2023
*/
func getDpId(lClientId string) (string, string, error) {
	log.Println("getDpId (+)")

	// this variables is used to get DpId and ClientName from the database.
	var lDpId string
	var lClientName string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
	if lErr1 != nil {
		log.Println("PGDI01", lErr1)
		return lDpId, lClientName, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select  idm.CLIENT_DP_CODE, idm.CLIENT_DP_NAME
						from   TECHEXCELPROD.CAPSFO.DBO.IO_DP_MASTER idm
						where idm.CLIENT_ID = ?
						and DEFAULT_ACC = 'Y'
						and DEPOSITORY = 'CDSL' `
		lRows, lErr2 := lDb.Query(lCoreString, lClientId)
		if lErr2 != nil {
			log.Println("PGDI02", lErr2)
			return lDpId, lClientName, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lDpId, &lClientName)
				if lErr3 != nil {
					log.Println("PGDI03", lErr3)
					return lDpId, lClientName, lErr3
				}
			}
		}
	}
	log.Println("getDpId (-)")
	return lDpId, lClientName, nil
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

Author:Pavithra
Date: 12JUNE2023
*/
func getPanNO(lClientId string) (string, error) {
	log.Println("getPanNO (+)")

	// this variables is used to get Pan number of the client from the database.
	var lPanNo string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
	if lErr1 != nil {
		log.Println("PGPN01", lErr1)
		return lPanNo, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select pan_no 
						from TECHEXCELPROD.CAPSFO.DBO.client_details
						where client_Id = ? `
		lRows, lErr2 := lDb.Query(lCoreString, lClientId)
		if lErr2 != nil {
			log.Println("PGPN02", lErr2)
			return lPanNo, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lPanNo)
				if lErr3 != nil {
					log.Println("PGPN03", lErr3)
					return lPanNo, lErr3
				}
			}
		}
	}
	log.Println("getPanNO (-)")
	return lPanNo, nil
}

/*
Pupose:This method used to retrieve the application No Id from the database.
Parameters:

	pAppNo,PClientId

Response:

	==========
	*On Sucess
	==========
	10,nil

	==========
	*On Error
	==========
	0,,error

Author:Pavithra
Date: 12JUNE2023
*/
func getAppNoId(pAppNo string, pClientId string) (int, error) {
	log.Println("getAppNoId (+)")

	// this variable is used to get the application number ID from the database
	var lAppNoId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PGAI01", lErr1)
		return lAppNoId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select (case when count(1) > 0 then h.Id else 0 end) AppNo
						from a_ipo_order_header h
						where h.applicationNo = ?
						and h.clientId = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pAppNo, pClientId)
		if lErr2 != nil {
			log.Println("PGAI02", lErr2)
			return lAppNoId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lAppNoId)
				if lErr3 != nil {
					log.Println("PGAI03", lErr3)
					return lAppNoId, lErr3
				}
			}
		}
	}
	log.Println("getAppNoId (-)")
	return lAppNoId, nil
}

// func CancelRecordUpdates(pReqArr []nse.ExchangeReqStruct, pRespArr []nse.ExchangeRespStruct) ([]nse.ExchangeRespStruct, string) {
// 	log.Println("CancelRecordUpdates (+)")

// 	var lBidRec nse.ResponseBidStruct
// 	lCancelFlag := "N"

// 	for lIdx := 0; lIdx < len(pReqArr); lIdx++ {

// 		if len(pReqArr[lIdx].Bids) > len(pRespArr[lIdx].Bids) {
// 			lBidArr := pReqArr[lIdx].Bids
// 			for lBidIdx := 0; lBidIdx < len(lBidArr); lBidIdx++ {
// 				if lBidArr[lBidIdx].ActivityType == "cancel" && pRespArr[lIdx].Status == "success" {
// 					lBidRec.BidReferenceNo, _ = strconv.Atoi(lBidArr[lBidIdx].BidReferenceNo)
// 					lBidRec.Status = "Success"
// 					lCancelFlag = "Y"
// 					pRespArr[lIdx].Bids = append(pRespArr[lIdx].Bids, lBidRec)
// 				}
// 			}
// 		}

// 	}
// 	log.Println("CancelRecordUpdates (-)")
// 	return pRespArr, lCancelFlag
// }

// func UpdateRecords(pReqArr []nse.ExchangeReqStruct, pRespArr []nse.ExchangeRespStruct, pHeaderId int, lClientId string, pDetailIdArr []int, pBidTrackIdArr []int, pCancelFlag string) error {
// 	log.Println("UpdateRecords (+)")
// 	for lpRespArrIdx := 0; lpRespArrIdx < len(pRespArr); lpRespArrIdx++ {
// 		log.Println(pRespArr[lpRespArrIdx].Status)
// 		if pRespArr[lpRespArrIdx].Status == "success" {
// 			// pReqArr
// 			lErr1 := UpdateHeader(pReqArr, pRespArr, pHeaderId, lClientId, pDetailIdArr, pBidTrackIdArr, pCancelFlag)
// 			if lErr1 != nil {
// 				log.Println("PUR01", lErr1)
// 				// fmt.Fprintf(w, helpers.GetErrorString("PPO08", "Error in Updating Details. Please try after sometime"))
// 				return lErr1
// 			}
// 		} else if pRespArr[lpRespArrIdx].Status == "failed" {
// 			for lReqIdx := 0; lReqIdx < len(pReqArr[lpRespArrIdx].Bids); lReqIdx++ {
// 				log.Println("inside first for")
// 				for lBidIdx := 0; lBidIdx < len(pRespArr[lpRespArrIdx].Bids); lBidIdx++ {
// 					log.Println("inside second for")
// 					if pReqArr[lpRespArrIdx].Bids[lReqIdx].Remark == pRespArr[lpRespArrIdx].Bids[lBidIdx].Remark {
// 						// updating bid tracking table
// 						lErr2 := UpdateBidTrack(pRespArr[lpRespArrIdx].Bids[lpRespArrIdx], lClientId, pBidTrackIdArr[lpRespArrIdx])
// 						if lErr2 != nil {
// 							log.Println("PUR02", lErr2)
// 							return lErr2
// 						} else {
// 							log.Println("BidTrack updated Successfully")
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("UpdateRecords (-)")
// 	return nil
// }

/*
Pupose:This method updating the order head values in order header table.
Parameters:
	pRespArr,pHeaderId,PClientId,pDetailIdArr,pBidTrackIdArr
Response:
	==========
	*On Sucess
	==========
	nil

	==========
	*On Error
	==========
	error

Author:Pavithra
Date: 12JUNE2023
*/
// func UpdateHeader(pReqArr []nse.ExchangeReqStruct, pRespArr []nse.ExchangeRespStruct, pHeaderId int, lClientId string, pDetailIdArr []int, pBidTrackIdArr []int, pCancelFlag string) error {
// 	log.Println("UpdateHeader (+)")

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr != nil {
// 		log.Println("IUH01", lErr)
// 		return lErr
// 	} else {
// 		defer lDb.Close()
// 		log.Println("pRespArr = ", pRespArr)
// 		// looping the pRespArr to update the application response details in database.
// 		for lpRespArrIdx := 0; lpRespArrIdx < len(pRespArr); lpRespArrIdx++ {

// 			// Changing the date format DD-MM-YYYY into YYYY-MM-DD of timestamp value.
// 			if pRespArr[lpRespArrIdx].TimeStamp != "" {
// 				lDate := strings.Split(pRespArr[lpRespArrIdx].TimeStamp, " ")
// 				lTimeStamp, _ := time.Parse("02-01-2006", lDate[0])
// 				pRespArr[lpRespArrIdx].TimeStamp = lTimeStamp.Format("2006-01-02") + " " + lDate[1]
// 				log.Println(pRespArr[lpRespArrIdx].TimeStamp, "If TimeStamp")
// 			} else if pRespArr[lpRespArrIdx].TimeStamp == "" {
// 				lTime := time.Now()
// 				pRespArr[lpRespArrIdx].TimeStamp = lTime.Format("2006-01-02 15:04:05")
// 				log.Println(pRespArr[lpRespArrIdx].TimeStamp, "Else if TimeStamp")
// 			}

// 			lSqlString := `update a_ipo_order_header h
// 							set h.Symbol = ?,h.category = ?,h.clientName = ?,h.depository = ?,
// 							h.dpId = ?,h.clientBenId = ?,h.nonASBA = ?,h.pan = ?,h.referenceNumber = ?,
// 							h.allotmentMode = ?,h.upiFlag = ?,h.upi = ?,h.bankCode = ?,h.locationCode = ?,
// 							h.bankAccount = ?,h.ifsc = ?,h.subBrokerCode = ?,h.time_Stamp =?,
// 							h.status = ?,h.dpVerStatusFlag = ?,h.dpVerFailCode = ?,h.dpVerReason = ?,
// 							h.upiPaymentStatusFlag = ?,h.upiAmtBlocked = ?,h.reasonCode = ?,h.reason = ?,
// 							h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag = ?
// 							where h.Id = ?`

// 			_, lErr := lDb.Exec(lSqlString, pRespArr[lpRespArrIdx].Symbol,
// 				pRespArr[lpRespArrIdx].Category, pRespArr[lpRespArrIdx].ClientName, pRespArr[lpRespArrIdx].Depository,
// 				pRespArr[lpRespArrIdx].DpId, pRespArr[lpRespArrIdx].ClientBenId, pRespArr[lpRespArrIdx].NonASBA,
// 				pRespArr[lpRespArrIdx].Pan, pRespArr[lpRespArrIdx].ReferenceNo, pRespArr[lpRespArrIdx].AllotmentMode,
// 				pRespArr[lpRespArrIdx].UpiFlag, pRespArr[lpRespArrIdx].Upi, pRespArr[lpRespArrIdx].BankCode,
// 				pRespArr[lpRespArrIdx].LocationCode, pRespArr[lpRespArrIdx].BankAccount, pRespArr[lpRespArrIdx].IFSC,
// 				pRespArr[lpRespArrIdx].SubBrokerCode, pRespArr[lpRespArrIdx].TimeStamp, pRespArr[lpRespArrIdx].Status,
// 				pRespArr[lpRespArrIdx].DpVerStatusFlag, pRespArr[lpRespArrIdx].DpVerFailCode, pRespArr[lpRespArrIdx].DpVerReason,
// 				pRespArr[lpRespArrIdx].UpiPaymentStatusFlag, pRespArr[lpRespArrIdx].UpiAmtBlocked,
// 				pRespArr[lpRespArrIdx].ReasonCode, pRespArr[lpRespArrIdx].Reason, lClientId, pCancelFlag, pHeaderId)
// 			if lErr != nil {
// 				log.Println("IUH02", lErr)
// 				return lErr
// 			} else {
// 				log.Println("Header Updated Sucessfully")
// 				// Call the UpdateDetails method ,to Update the Bid Details Reponse in Order Details Table
// 				lErr := UpdateDetails(pReqArr[lpRespArrIdx].Bids, pRespArr[lpRespArrIdx].Bids, lClientId, pDetailIdArr, pBidTrackIdArr)
// 				if lErr != nil {
// 					log.Println("IUH03", lErr)
// 					return lErr
// 				} else {
// 					log.Println("Details updated Successfully")
// 				}
// 			}

// 		}
// 	}
// 	log.Println("UpdateHeader (-)")
// 	return nil
// }

// /*
// Pupose:This method updating the bid details in order detail table.
// Parameters:
// 	pRespBidArr,lClientId,pDetailIdArr,pBidTrackIdArr
// Response:
// 	==========
// 	*On Sucess
// 	==========
// 	nil

// 	==========
// 	*On Error
// 	==========
// 	error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func UpdateDetails(pReqBidArr []nse.RequestBidStruct, pRespBidArr []nse.ResponseBidStruct, lClientId string, pDetailIdArr []int, pBidTrackIdArr []int) error {
// 	log.Println("UpdateDetails (+)")

// 	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr != nil {
// 		log.Println("IUD01", lErr)
// 		return lErr
// 	} else {
// 		defer lDb.Close()
// 		// Looping the Bid Details Array to Update the Application Response Values in Order Details Table.
// 		log.Println(pRespBidArr)
// 		log.Println(pReqBidArr)
// 		for lpRespBidIdx := 0; lpRespBidIdx < len(pRespBidArr); lpRespBidIdx++ {
// 			for lpReqBidIdx := 0; lpReqBidIdx < len(pReqBidArr); lpReqBidIdx++ {

// 				if pRespBidArr[lpRespBidIdx].ActivityType == "modify" {

// 					if pRespBidArr[lpRespBidIdx].ActivityType == pReqBidArr[lpReqBidIdx].ActivityType && strconv.Itoa(pRespBidArr[lpRespBidIdx].BidReferenceNo) == pReqBidArr[lpReqBidIdx].BidReferenceNo {
// 						// ! ----------------
// 						lSqlString := `update a_ipo_orderdetails d
// 									set d.activityType = ?,d.bidReferenceNo = ?,d.series = ?,d.atCutOff = ?,
// 									d.resp_quantity = ?,d.resp_price = ?,d.resp_amount = ?,d.remark = ?,
// 									d.status = ?,d.reasonCode = ?,d.reason = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
// 									where d.id = ?`

// 						_, lErr := lDb.Exec(lSqlString, pRespBidArr[lpRespBidIdx].ActivityType, pRespBidArr[lpRespBidIdx].BidReferenceNo,
// 							pRespBidArr[lpRespBidIdx].Series, pRespBidArr[lpRespBidIdx].AtCutOff, pRespBidArr[lpRespBidIdx].Quantity,
// 							pRespBidArr[lpRespBidIdx].Price, pRespBidArr[lpRespBidIdx].Amount, pRespBidArr[lpRespBidIdx].Remark,
// 							pRespBidArr[lpRespBidIdx].Status, pRespBidArr[lpRespBidIdx].ReasonCode, pRespBidArr[lpRespBidIdx].Reason,
// 							lClientId, pDetailIdArr[lpRespBidIdx])
// 						if lErr != nil {
// 							log.Println("Error updating into database (details)")
// 							log.Println("IUD02", lErr)
// 							return lErr
// 						} else {
// 							// Call the UpdateBidTrack method ,to Update the Bid Details Reponse in Bid Tracking Table
// 							lErr := UpdateBidTrack(pRespBidArr[lpRespBidIdx], lClientId, pBidTrackIdArr[lpRespBidIdx])
// 							if lErr != nil {
// 								log.Println("IUD03", lErr)
// 								return lErr
// 							} else {
// 								log.Println("BidTrack updated Successfully")
// 							}
// 						}
// 					}
// 				} else {
// 					lSqlString := `update a_ipo_orderdetails d
// 									set d.activityType = ?,d.bidReferenceNo = ?,d.series = ?,d.atCutOff = ?,
// 									d.resp_quantity = ?,d.resp_price = ?,d.resp_amount = ?,d.remark = ?,
// 									d.status = ?,d.reasonCode = ?,d.reason = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
// 									where d.id = ?`

// 					_, lErr := lDb.Exec(lSqlString, pRespBidArr[lpRespBidIdx].ActivityType, pRespBidArr[lpRespBidIdx].BidReferenceNo,
// 						pRespBidArr[lpRespBidIdx].Series, pRespBidArr[lpRespBidIdx].AtCutOff, pRespBidArr[lpRespBidIdx].Quantity,
// 						pRespBidArr[lpRespBidIdx].Price, pRespBidArr[lpRespBidIdx].Amount, pRespBidArr[lpRespBidIdx].Remark,
// 						pRespBidArr[lpRespBidIdx].Status, pRespBidArr[lpRespBidIdx].ReasonCode, pRespBidArr[lpRespBidIdx].Reason,
// 						lClientId, pDetailIdArr[lpRespBidIdx])
// 					if lErr != nil {
// 						log.Println("Error updating into database (details)")
// 						log.Println("IUD02", lErr)
// 						return lErr
// 					} else {
// 						// Call the UpdateBidTrack method ,to Update the Bid Details Reponse in Bid Tracking Table
// 						lErr := UpdateBidTrack(pRespBidArr[lpRespBidIdx], lClientId, pBidTrackIdArr[lpRespBidIdx])
// 						if lErr != nil {
// 							log.Println("IUD03", lErr)
// 							return lErr
// 						} else {
// 							log.Println("BidTrack updated Successfully")
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("UpdateDetails (-)")
// 	return nil
// }

/*
Pupose:This method is used to get the placed application status.
Parameters:
	pRespArr
Response:
	==========
	*On Sucess
	==========
	[
		{
			"applicationNo":"FT000069130109",
			"bidRefNo":"2023061400000028",
			"bidStatus":"success",
			"reason":"",
			"appStatus":"success",
			"appReason":""
		}
	],
	nil

	==========
	*On Error
	==========
	[],error

Author:Pavithra
Date: 12JUNE2023
*/
// func OrderResponse(pAppNo string, pSymbol string, lClientId string) (OrderResStruct, error) {
// 	log.Println("OrderResponse (+)")
// 	var lRespAppRec OrderResStruct

// 	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr != nil {
// 		log.Println("IOR01", lErr)
// 		return lRespAppRec, lErr
// 	} else {
// 		defer lDb.Close()

// 		lAppNo, lErr := getAppResp(pSymbol, pAppNo, lClientId)
// 		if lErr != nil {
// 			log.Println("IOR02", lErr)
// 			return lRespAppRec, lErr
// 		} else {
// 			log.Println("*******", lAppNo)

// 			lCoreString := `select h.applicationNo,nvl(h.status,'') ,nvl(h.reason ,'')
// 							from a_ipo_order_header h,a_ipo_orderdetails d
// 							where h.Id = d.headerId
// 							and h.applicationNo = ?`

// 			lRows, lErr := lDb.Query(lCoreString, lAppNo)
// 			log.Println("****", lRows)
// 			if lErr != nil {
// 				log.Println("IOR03", lErr)
// 				return lRespAppRec, lErr
// 			} else {
// 				for lRows.Next() {
// 					lErr := lRows.Scan(&lRespAppRec.ApplicationNo, &lRespAppRec.AppStatus, &lRespAppRec.AppReason)
// 					if lErr != nil {
// 						log.Println("IOR04", lErr)
// 						return lRespAppRec, lErr
// 					} else {
// 						// lResultArr = append(lResultArr, lResultRec)
// 						lRespAppRec.Status = common.SuccessCode
// 					}
// 				}
// 			}

// 		}
// 	}
// 	log.Println("OrderResponse (-)")
// 	return lRespAppRec, nil
// }

/*
Pupose:This method is used to get the Application Number From the DataBase.
Parameters:
	pSymbol,pAppNo,PClientId
Response:
	==========
	*On Sucess
	==========
	FT000006912345,nil

	==========
	*On Error
	==========
	"",error

Author:Pavithra
Date: 12JUNE2023
*/
// func getAppResp(pSymbol string, pAppNo string, lClientId string) (string, error) {
// 	log.Println("getAppResp (+)")
// 	log.Println(pSymbol, pAppNo, lClientId)

// 	// this variable is used to get the application no from the database
// 	var lAppNo string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr != nil {
// 		log.Println("IGA01", lErr)
// 		return lAppNo, lErr
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select (case when count(1) > 0 then h.applicationNo else 'nil' end) AppNo
// 						from a_ipo_order_header h
// 						where h.Symbol = ?
// 						and h.clientId = ?
// 						and h.applicationNo = ?`
// 		lRows, lErr := lDb.Query(lCoreString, pSymbol, lClientId, pAppNo)
// 		if lErr != nil {
// 			log.Println("IGA02", lErr)
// 			return lAppNo, lErr
// 		} else {
// 			for lRows.Next() {
// 				lErr := lRows.Scan(&lAppNo)
// 				if lErr != nil {
// 					log.Println("IGA03", lErr)
// 					return lAppNo, lErr
// 				} else {
// 					log.Println("Application No", lAppNo)
// 				}
// 			}
// 		}
// 	}
// 	log.Println("getAppResp (-)")
// 	return lAppNo, nil
// }

/*
Pupose:This method fetch the application number from  a_ipo_order_header table
Parameters:
    pRespArr


Response:

	==========
	*On Sucess
	==========


	==========
	*On Error
	==========


Author:KAVYA DHARSHANI M
Date: 08SEP2023
*/

func GetApplication(pRespArr []nseipo.ExchangeRespStruct) (IpoOrderStruct, error) {
	log.Println("GetApplication(+)")

	//This variable is used to get  the ApplicationNumber  in Array
	var lAppRec IpoOrderStruct

	//This variable is used to append the array in the IpoOrderStruct Structure
	// var lorderNoArr []iposchedule.IpoOrderStruct

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)

	if lErr1 != nil {
		log.Println("IPOGA01", lErr1)
		return lAppRec, lErr1
	} else {
		defer lDb.Close()

		for lAppIdx := 0; lAppIdx < len(pRespArr); lAppIdx++ {
			log.Println("pPendingArr[lpendingIdx].ApplicationNo", pRespArr[lAppIdx].ApplicationNo)

			lCoreString := `select h.applicationNo, h.ClientEmail, h.clientId  , h.clientName, h.status , m.Symbol ,d.activityType ,
			max(d.req_amount),date(d.CreatedDate)  
			from a_ipo_order_header h , a_ipo_master m,a_ipo_orderdetails d
			where h.MasterId  = m.Id
			and h.Id = d.headerId 
			and h.applicationNo  = ?
			group by d.headerId `

			lRows, lErr2 := lDb.Query(lCoreString, pRespArr[lAppIdx].ApplicationNo)
			if lErr2 != nil {
				log.Println("IPOGA02", lErr2)
				return lAppRec, lErr2
			} else {
				for lRows.Next() {
					lErr3 := lRows.Scan(&lAppRec.ApplicationNumber, &lAppRec.ClientEmail, &lAppRec.ClientId, &lAppRec.ClientName, &lAppRec.Status, &lAppRec.Symbol, &lAppRec.ActivityType, &lAppRec.Amount, &lAppRec.OrderDate)
					if lErr3 != nil {
						log.Println("IPOGA03", lErr3)
						return lAppRec, lErr3
					}
				}
			}
		}
	}
	log.Println("GetApplication(-)")
	return lAppRec, nil
}

/*

Purpose:This method is used to send mail based on the Application Number
Parameters:

	lAppArr,pStatus


Response:

	==========
	*On Sucess
	==========


	==========
	*On Error
	==========

Author:Kavya Dharshani
Date: 07SEP2023


*/
func ConstructMail(lAppArr IpoOrderStruct, pStatus string) (emailUtil.EmailInput, error) {
	type dynamicEmailStruct struct {
		Name              string
		Symbol            string
		ApplicationNumber string
		OrderDate         string
		Unit              string
		Price             string
		Amount            int
		Activity          string
		Status            string
		EmailHeader       string
		EmailBody         string
	}

	var lIpoEmailContent emailUtil.EmailInput
	config := common.ReadTomlConfig("./toml/emailconfig.toml")
	lIpoEmailContent.FromDspName = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	lIpoEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	lIpoEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	// commented by pavithra
	// lIpoEmailContent.Subject = "IPO Order"
	// html := "html/IpoOrderTemplate.html"

	// newly added to get the subject and htnl path in toml file
	lNovoConfig := common.ReadTomlConfig("./toml/novoConfig.toml")
	lIpoEmailContent.Subject = fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["IPO_ClientEmail_Subject"])
	html := fmt.Sprintf("%v", lNovoConfig.(map[string]interface{})["IPO_ClientEmail_html"])

	lTemp, lErr := template.ParseFiles(html)
	if lErr != nil {
		log.Println("IMP03", lErr)
		return lIpoEmailContent, lErr
	} else {
		var lTpl bytes.Buffer
		var lDynamicEmailVal dynamicEmailStruct

		lDynamicEmailVal.Name = lAppArr.ClientName
		lDynamicEmailVal.Symbol = lAppArr.Symbol
		lDynamicEmailVal.Amount = lAppArr.Amount
		lDynamicEmailVal.OrderDate = lAppArr.OrderDate
		lDynamicEmailVal.ApplicationNumber = lAppArr.ApplicationNumber
		lDynamicEmailVal.Activity = lAppArr.ActivityType
		lDynamicEmailVal.Status = pStatus

		lHeader, lErr2 := GetEmailHeader(lAppArr.ActivityType)
		if lErr2 != nil {
			log.Println("IPCM02", lErr2)
			return lIpoEmailContent, lErr2
		} else {
			// get the email content header for Ipo Order Report from the lookup table

			lDynamicEmailVal.EmailHeader = lHeader

			lBody, lErr3 := GetEmailBody(lAppArr.ActivityType)
			if lErr3 != nil {
				log.Println("IPCM03", lErr3)
				return lIpoEmailContent, lErr3
			} else {
				// get the email body for Ipo Order Report from the lookup table
				lDynamicEmailVal.EmailBody = lBody
			}
		}

		// commented by pavithra
		//  This Temperory Method Add to Send Mail For Client
		// lClientEmail, lErr1 := clientDetail.GetClientEmailId(lAppArr.ClientId)
		// if lErr1 != nil {
		// return lIpoEmailContent, lErr1
		// }
		// log.Println("lAppArr.ClientEmail", lAppArr.ClientEmail, "lClientEmail", lClientEmail)
		// lIpoEmailContent.ToEmailId = "prashanth.s@fcsonline.co.in"
		// lIpoEmailContent.ToEmailId = lClientEmail

		lIpoEmailContent.ToEmailId = lAppArr.ClientEmail

		// log.Println(lDynamicEmailVal, "dynamic struct")
		lTemp.Execute(&lTpl, lDynamicEmailVal)
		lEmailbody := lTpl.String()

		lIpoEmailContent.Body = lEmailbody
	}
	return lIpoEmailContent, nil
}

// func constructmail(lorderArr []AppOrderResStruct, status string) (emailutil.EmailInput, error) {
// 	type dynamicEmailStruct struct {
// 		Name              string
// 		Symbol            string
// 		ApplicationNumber string
// 		OrderDate         string
// 		Unit              string
// 		Price             string
// 		Amount            int
// 		Activity          string
// 		Status            string
// 	}

// 	var lIpoEmailContent emailutil.EmailInput
// 	config := common.ReadTomlConfig("./toml/emailconfig.toml")
// 	lIpoEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["From"])
// 	lIpoEmailContent.FromRaw = fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
// 	lIpoEmailContent.ReplyTo = fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
// 	lIpoEmailContent.Subject = "IPO Order"
// 	html := "html/IpoOrderTemplate.html"

// 	lTemp, lErr := template.ParseFiles(html)
// 	if lErr != nil {
// 		log.Println("POCM01", lErr)
// 		return lIpoEmailContent, lErr
// 	} else {
// 		var lTpl bytes.Buffer
// 		var lDynamicEmailVal dynamicEmailStruct

// 		for i := 0; i < len(lorderArr); i++ {

// 			if status == lorderArr[i].Status {
// 				log.Println("Success")
// 				lDynamicEmailVal.Name = lorderArr[i].ClientName
// 				lDynamicEmailVal.Symbol = lorderArr[i].Symbol
// 				lUnit, _ := strconv.Atoi(lorderArr[i].Unit)
// 				lPrice, _ := strconv.Atoi(lorderArr[i].Price)
// 				lDynamicEmailVal.Amount = lUnit * lPrice
// 				lDynamicEmailVal.Unit = lorderArr[i].Unit
// 				lDynamicEmailVal.Price = lorderArr[i].Price
// 				lDynamicEmailVal.OrderDate = lorderArr[i].OrderDate
// 				lDynamicEmailVal.ApplicationNumber = lorderArr[i].ApplicationNumber
// 				lDynamicEmailVal.Activity = lorderArr[i].ActivityType
// 				lDynamicEmailVal.Status = status
// 				lIpoEmailContent.ToEmailId = lorderArr[i].ClientEmail
// 			}

// 		}
// 		log.Println(lDynamicEmailVal, "dynamic struct")

// 		lTemp.Execute(&lTpl, lDynamicEmailVal)
// 		lEmailbody := lTpl.String()

// 		lIpoEmailContent.Body = lEmailbody

// 	}
// 	return lIpoEmailContent, nil
// }
//

//--------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------

// func ProcessReq(pExchangeReq []nseipo.ExchangeReqStruct, pMasterId int, pClientId string) ([]nseipo.ExchangeRespStruct, error) {
// 	log.Println("ProcessReq (+)")

// 	// for retruning value for this method
// 	var lRespArr []nseipo.ExchangeRespStruct

// 	// lBidTrackId, lErr1 := InsertBidTrack(pExchangeReq, pClientId)
// 	// if lErr1 != nil {
// 	// 	log.Println("PPR01", lErr1)
// 	// 	return lRespArr, lErr1
// 	// } else {
// 	// ------------------------>>>>>>>>>>>>>>>>>>
// 	// Call the ApplyIpo method to Process the Application Request.
// 	lResponse, lErr2 := exchangecall.ApplyIpo(pExchangeReq, pClientId)
// 	if lErr2 != nil {
// 		log.Println("PPR02", lErr2.Error())
// 		return lRespArr, lErr2
// 	} else {
// 		if lResponse == nil {
// 			return lRespArr, lErr2
// 		} else {
// 			lRespArr = lResponse
// 			// update the bid tracking table when application status success or failed
// 			// lErr3 := UpdateBidTrack(lResponse, pClientId, lBidTrackId)
// 			// if lErr3 != nil {
// 			// 	log.Println("PPR03", lErr3)
// 			// 	return lRespArr, lErr3
// 			// } else {
// 			for lRespIdx := 0; lRespIdx < len(lResponse); lRespIdx++ {
// 				// check whether the application status is success or not
// 				if lResponse[lRespIdx].Status == "success" {
// 					// if the status is success update the response in header and details table.
// 					// lErr4 := InsertHeader(lResponse, pExchangeReq, pMasterId, pClientId)
// 					// if lErr4 != nil {
// 					// 	log.Println("PPR04", lErr4)
// 					// 	return lRespArr, lErr4
// 					// }
// 				}
// 			}
// 			// }
// 		}
// 	}
// 	// }
// 	// }
// 	log.Println("ProcessReq (-)")
// 	return lRespArr, nil
// }
func FetchDirectory(pReqRec OrderReqStruct) (string, error) {
	log.Println("FetchDirectory (+)")

	var lDirectory string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AFD01", lErr1)
		return lDirectory, lErr1
	} else {
		defer lDb.Close()

		if pReqRec.ApplicationNo == "" {

			config := common.ReadTomlConfig("./toml/debug.toml")
			lCheckDirectory := fmt.Sprintf("%v", config.(map[string]interface{})["CurrentDirectory"])

			if lCheckDirectory == "" {
				// Commnented by nithish kumar on 10JAN2024 to fix the bug to take atleast one exchange of the ipo
				// below commented query fails when we check the isin with the exchange of the IPO
				// lCoreString := `select (case when (select m.Isin  from a_ipo_master m where m.Symbol = ? and m.Exchange = "NSE") =
				// 					(select m.Isin  from a_ipo_master m where m.Symbol = ? and m.Exchange = "BSE") then 'BSE' else am.Exchange  end )Exchange
				// 					from a_ipo_master am
				// 					where am.symbol = ?
				// 					 and am.Id = ?`

				//below query added by nithish on 10Jan2024 to take atleast one exchange of the ipo
				lCoreString := `select (case when count(distinct m.Exchange) > 1 then 'BSE' else m.Exchange end )Exchange 
								from a_ipo_master m 
								where m.Symbol = ?
								and m.Exchange in('NSE','BSE')`
				lRows, lErr2 := lDb.Query(lCoreString, pReqRec.Symbol)
				if lErr2 != nil {
					log.Println("AFD02", lErr2)
					return lDirectory, lErr2
				} else {
					//This for loop is used to collect the record from the database and store them in structure
					for lRows.Next() {
						lErr3 := lRows.Scan(&lDirectory)
						if lErr3 != nil {
							log.Println("AFD03", lErr3)
							return lDirectory, lErr3
						}
					}
					log.Println("Current IPO directory := ", lDirectory)
				}
			} else {
				lDirectory = lCheckDirectory
			}
		} else {
			//Commnented by lakshmanan on 20OCT2023 to fix the bug to take order level exchange
			//below commented query always takes master level exchange
			// lCoreString := `select (case when count(1) > 0 then m.Exchange else (select m.Exchange from a_ipo_master m where m.Id = ?) end) ex
			// from a_ipo_master m,a_ipo_order_header h
			// where m.Symbol = ?
			// and m.Id = h.MasterId
			// and h.applicationNo = ?`
			// lRows, lErr2 := lDb.Query(lCoreString, pReqRec.MasterId, pReqRec.Symbol, pReqRec.ApplicationNo)
			//below query added by lakshmanan on 20OCT2023 to take order level exchange
			lCoreString := `select (case when nvl(h.exchange,'') <> '' then h.exchange else m.Exchange end) ex    
							from a_ipo_master m,a_ipo_order_header h
							where  m.Id = h.MasterId 
							and h.applicationNo = ? `
			lRows, lErr2 := lDb.Query(lCoreString, pReqRec.ApplicationNo)
			if lErr2 != nil {
				log.Println("AFD04", lErr2)
				return lDirectory, lErr2
			} else {
				//This for loop is used to collect the record from the database and store them in structure
				for lRows.Next() {
					lErr3 := lRows.Scan(&lDirectory)
					if lErr3 != nil {
						log.Println("AFD05", lErr3)
						return lDirectory, lErr3
					}
				}
				log.Println("Current IPO directory := ", lDirectory)
			}
		}
		// // Marshall the structure into json
		// lData, lErr5 := json.Marshal(lRespRec)
		// if lErr5 != nil {
		//  log.Println("AFD05", lErr5)
		//  fmt.Fprintf(w, helpers.GetErrorString("AFD05", "Unable to process your request now. Please try after sometime"))
		//  return
		// } else {
		//  fmt.Fprintf(w, string(lData))
		// }
		// log.Println("FetchDirectory (-)", r.Method)
	}
	log.Println("FetchDirectory (-)")
	return lDirectory, nil
}

func GetEmailHeader(pActivityType string) (string, error) {
	log.Println("GetEmailHeader (+)")

	//Instance to capture the email contentHeader based on the activity type
	lEmailHeader := ""

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IPGEH01", lErr1)
		return lEmailHeader, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select nvl(ld.description,'') description  
						from xx_lookup_details ld ,xx_lookup_header lh
						where lh.id = ld.headerid 	
						and lh.Code = 'ipo_email_contentHeader'
						and ld.Code = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pActivityType)
		if lErr2 != nil {
			log.Println("IPGEH02", lErr2)
			return lEmailHeader, lErr2
		} else {
			// This for loop is used to collect the record from the database
			// and store them in variable
			for lRows.Next() {
				lErr3 := lRows.Scan(&lEmailHeader)
				if lErr3 != nil {
					log.Println("IPGEH03", lErr3)
					return lEmailHeader, lErr3
				}
			}
		}
	}
	log.Println("GetEmailHeader (-)")
	return lEmailHeader, nil
}

func GetEmailBody(pActivityType string) (string, error) {
	log.Println("GetEmailBody (+)")

	// Instance to capture the email contentBody based on the activity type
	lEmailBody := ""

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IPGEB01", lErr1)
		return lEmailBody, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select nvl(ld.description,'') description  
						from xx_lookup_details ld ,xx_lookup_header lh
						where lh.id = ld.headerid 	
						and lh.Code = 'ipo_email_contentBody'
						and ld.Code = ? `
		lRows, lErr2 := lDb.Query(lCoreString, pActivityType)
		if lErr2 != nil {
			log.Println("IPGEB02", lErr2)
			return lEmailBody, lErr2
		} else {
			// This for loop is used to collect the record from the database
			// and store them in variable
			for lRows.Next() {
				lErr3 := lRows.Scan(&lEmailBody)
				if lErr3 != nil {
					log.Println("IPGEB03", lErr3)
					return lEmailBody, lErr3
				}
			}
		}
	}
	log.Println("GetEmailBody (-)")
	return lEmailBody, nil
}
