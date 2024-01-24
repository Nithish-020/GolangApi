package sgbplaceorder

import (
	"fcs23pkg/integration/bse/bsesgb"
	"log"
	"strconv"
)

// // func NsePlaceOrder(pExchangeReq bsesgb.SgbReqStruct, pClientId string, pReqRec SgbReqStruct, r *http.Request) (bsesgb.SgbRespStruct, string, error) {
// // log.Println("NsePlaceOrder (+)")

// // var lRespRec string
// // // var lErrorRec error
// // var lError ErrorStruct
// // // var lExchangeResp nsesgb.SgbAddResStruct

// // lRespRec = "S"

// // lExchangeResp, lRespRec, lErr1 := ProcessSgbOrder(pExchangeReq, pClientId, pReqRec, r)
// // if lErr1 != nil {
// // 	lRespRec = "E"
// // 	lError.ErrCode = "PBPO01"
// // 	lError.ErrMsg = "Exchange Server is Busy right now,Try After Sometime."
// // 	return lExchangeResp, lRespRec, lErr1
// // }
// // log.Println("lExchangeResp", lExchangeResp)
// // log.Println("NsePlaceOrder (-)")
// // return lExchangeResp, lRespRec, lErr1
// // }

// // func SgbReqConstruct(pSNseReq nsesgb.SgbAddReqStruct) bsesgb.SgbReqStruct {
// // 	log.Println("SgbReqConstruct(+)")

// // 	var lSBseReq bsesgb.SgbReqStruct
// // 	var lSBse bsesgb.ReqSgbBidStruct

// // 	lSBseReq.ScripId = pSNseReq.Symbol
// // 	lSBseReq.Depository = pSNseReq.Depository
// // 	lSBseReq.DpId = pSNseReq.DpId
// // 	lSBseReq.ClientBenfId = pSNseReq.ClientBenId

// // 	if pSNseReq.ActivityType == "ER" {
// // 		lSBse.ActionCode = "Entry request"
// // 	} else if pSNseReq.ActivityType == "MR" {
// // 		lSBse.ActionCode = "Modify Request"
// // 	} else {
// // 		lSBse.ActionCode = "Cancel Request"
// // 	}

// // 	lSBse.BidId = ""

// // 	lSubscriptionUnit := strconv.Itoa(pSNseReq.Quantity)
// // 	lSBse.SubscriptionUnit = lSubscriptionUnit

// // 	lOrderNumber := strconv.Itoa(pSNseReq.OrderNumber)
// // 	lSBse.OrderNo = lOrderNumber

// // 	lPrice := strconv.FormatFloat(float64(pSNseReq.Price), 'f', -1, 32)
// // 	lSBse.Rate = lPrice

// // 	lSBseReq.Bids = append(lSBseReq.Bids, lSBse)

// // 	log.Println(lSBseReq.Bids, "Request")

// // 	log.Println(lSBseReq, "Request")

// // 	log.Println("SgbReqConstruct(-)")
// // 	return lSBseReq
// // }

// // func SgbRespConstruct(pSBseResp bsesgb.SgbRespStruct) nsesgb.SgbAddResStruct {
// // 	log.Println("SgbReqConstruct(-)")

// // 	var lSNseResp nsesgb.SgbAddResStruct

// // 	lSNseResp.Symbol = pSBseResp.ScripId
// // 	lSNseResp.Depository = pSBseResp.Depository
// // 	lSNseResp.DpId = pSBseResp.DpId
// // 	lSNseResp.ClientBenId = pSBseResp.ClientBenfId
// // 	lSNseResp.Status = pSBseResp.StatusCode
// // 	lSNseResp.Reason = pSBseResp.StatusMessage

// // 	for ldx := 0; ldx < len(pSBseResp.Bids); ldx++ {

// // 		if pSBseResp.Bids[ldx].ActionCode == "Entry request" {
// // 			pSBseResp.Bids[ldx].ActionCode = "ER"
// // 		} else if pSBseResp.Bids[ldx].ActionCode == "Modify Request" {
// // 			pSBseResp.Bids[ldx].ActionCode = "MR"
// // 		} else {
// // 			pSBseResp.Bids[ldx].ActionCode = "CR"
// // 		}

// // 		lSubscriptionUnit, lErr1 := strconv.Atoi(pSBseResp.Bids[ldx].SubscriptionUnit)
// // 		if lErr1 != nil {
// // 			log.Println("Error:", lErr1)
// // 		} else {

// // 			lSNseResp.Quantity = lSubscriptionUnit
// // 		}

// // 		lRate, lErr2 := common.ConvertStringToFloat(pSBseResp.Bids[ldx].Rate)
// // 		if lErr2 != nil {
// // 			log.Println("Error:", lErr2)
// // 		} else {
// // 			lSNseResp.Price = lRate
// // 		}

// // 		OrderNo, lErr3 := strconv.Atoi(pSBseResp.Bids[ldx].OrderNo)
// // 		if lErr3 != nil {
// // 			log.Println("Error:", lErr3)
// // 		} else {
// // 			lSNseResp.OrderNumber = OrderNo
// // 		}

// // 		lSNseResp.Status = pSBseResp.Bids[ldx].ErrorCode
// // 		lSNseResp.Reason = pSBseResp.Bids[ldx].Message

// // 		log.Println(lSNseResp, "Response")

// // 	}

// // 	log.Println("SgbReqConstruct(-)")
// // 	return lSNseResp

// // }

// //----------------------------
// package sgbplaceorder

// import (
// 	"encoding/json"
// 	"fcs23pkg/apps/SGB/validatesgb"
// 	"fcs23pkg/apps/clientDetail"
// 	"fcs23pkg/apps/exchangecall"
// 	"fcs23pkg/apps/validation/adminaccess"
// 	"fcs23pkg/apps/validation/apiaccess"
// 	"fcs23pkg/common"
// 	"fcs23pkg/ftdb"
// 	"fcs23pkg/helpers"
// 	"fcs23pkg/integration/bse/bsesgb"

// 	"fcs23pkg/util/emailUtil"

// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"time"
// )

// // this struct is used to get the bid details.
// type SgbReqStruct struct {
// 	MasterId   int    `json:"masterId"`
// 	BidId      string `json:"bidId"`
// 	Unit       int    `json:"unit"`
// 	OldUnit    int    `json:"oldUnit"`
// 	Price      int    `json:"price"`
// 	ActionCode string `json:"actionCode"`
// 	OrderNo    string `json:"orderNo"`
// 	Amount     int    `json:"amount"`
// 	PreApply   string `json:"preApply"`
// }

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
// 	MasterId    int    `json:"masterId"`
// 	Unit        string `json:"unit"`
// 	Price       string `json:"price"`
// 	JVamount    string `json:"jvAmount"`
// 	OrderNo     string `json:"orderNo"`
// 	ClientId    string `json:"clientId"`
// 	ActionCode  string `json:"actionCode"`
// 	Transaction string `json:"transaction"`
// 	Flag        string `json:"flag"`
// 	BidId       string `json:"bidId"`
// }

// // this struct is used to get the bid details.
// type SgbRespStruct struct {
// 	OrderStatus string `json:"orderStatus"`
// 	Status      string `json:"status"`
// 	ErrMsg      string `json:"errMsg"`
// }

// type JvStatusStruct struct {
// 	JvAmount    string `json:"jvAmount"`
// 	JvStatus    string `json:"jvStatus"`
// 	JvStatement string `json:"jvStatement"`
// 	JvType      string `json:"jvType"`
// }

// /*
// Pupose:This API method is used to place a order
// Request:

// 	lReqRec

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	{
// 		"appResponse":[
// 			{
// 				"applicationNo":"FT000069130109",
// 				"bidRefNo":"2023061400000028",
// 				"bidStatus":"success",
// 				"reason":"",
// 				"appStatus":"success",
// 				"appReason":""
// 			}
// 		],
// 		"status":"S",
// 		"errMsg":""
// 	}
// 	==========
// 	*On Error
// 	==========

// 		{
// 		"appResponse":[],
// 		"status":"E",
// 		"errMsg":""
// 	}

// Author:Kavya Dharshani
// Date: 20SEP2023
// */
// func SgbPlaceOrder(w http.ResponseWriter, r *http.Request) {
// 	log.Println("SgbPlaceOrder (+)", r.Method)
// 	// origin := r.Header.Get("Origin")
// 	// for _, allowedOrigin := range common.ABHIAllowOrigin {
// 	// 	if allowedOrigin == origin {
// 	// 		w.Header().Set("Access-Control-Allow-Origin", origin)
// 	// 		log.Println(origin)
// 	// 		break
// 	// 	}
// 	// }
// 	(w).Header().Set("Access-Control-Allow-Origin", common.ABHIAllowOrigin)
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	// Checking the Api Method
// 	if r.Method == "POST" {
// 		// create a instance of array type OrderReqStruct.This variable is used to store the request values.
// 		var lReqRec SgbReqStruct
// 		// create a instance of array type OrderResStruct.This variable is used to store the response values.
// 		var lRespRec SgbRespStruct
// 		//
// 		lRespRec.Status = common.SuccessCode

// 		lExchange, lErr1 := adminaccess.FetchDirectory()
// 		if lErr1 != nil {
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "SGBPO01" + lErr1.Error()
// 			fmt.Fprintf(w, helpers.GetErrorString("SGBPO01", "Directory Not Found. Please try after sometime"))
// 			return
// 		} else {
// 			lExchange = "BSE"
// 			lRespRec.Status = common.SuccessCode
// 			var lErrorRec error
// 			var lFlag string
// 			var lExchangeResp bsesgb.SgbRespStruct

// 			//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------

// 			lClientId, lErr2 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
// 			if lErr2 != nil {
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = "SGBPO02" + lErr2.Error()
// 				fmt.Fprintf(w, helpers.GetErrorString("SGBPO02", "UserDetails not Found"))
// 				return
// 			} else {
// 				if lClientId == "" {
// 					lRespRec.Status = common.ErrorCode
// 					fmt.Fprintf(w, helpers.GetErrorString("SGBPO02", "UserDetails not Found"))
// 					return
// 				}
// 			}
// 			//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
// 			// Read the request body values in lBody variable
// 			lBody, lErr3 := ioutil.ReadAll(r.Body)
// 			log.Println(string(lBody))
// 			log.Println("lBody", lBody)
// 			if lErr3 != nil {
// 				log.Println("SGBPO03", lErr3.Error())
// 				lRespRec.Status = common.ErrorCode
// 				fmt.Fprintf(w, helpers.GetErrorString("SGBPO03", "Unable to get your request now!!Try Again.."))
// 				return
// 			} else {
// 				// Unmarshal the request body values in lReqRec variable
// 				lErr4 := json.Unmarshal(lBody, &lReqRec)
// 				if lErr4 != nil {
// 					log.Println("SGBPO04", lErr4.Error())
// 					lRespRec.Status = common.ErrorCode
// 					fmt.Fprintf(w, helpers.GetErrorString("SGBPO04", "Unable to get your request now. Please try after sometime"))
// 					return
// 				} else {
// 					log.Println("lReqRec", lReqRec)
// 					//this method is used to construct the Req struct for exchange
// 					lExchangeReq, lErr5 := ConstructSGBReqStruct(lReqRec, lClientId, lExchange)
// 					if lErr5 != nil {
// 						log.Println("SGBPO05", lErr5.Error())
// 						lRespRec.Status = common.ErrorCode
// 						fmt.Fprintf(w, helpers.GetErrorString("SGBPO05", "Unable to process your request now. Please try after sometime"))
// 						return
// 					} else {
// 						log.Println(lExchangeReq)
// 						if lReqRec.PreApply == "pre" {
// 							lResp, lErr6 := UpdateSgbPending(lExchangeReq, lClientId, lReqRec, r, lExchange)
// 							if lErr6 != nil {
// 								log.Println("SGBPO06", lErr6.Error())
// 								lRespRec.Status = common.ErrorCode
// 								fmt.Fprintf(w, helpers.GetErrorString("SGBPO06", "Unable to process your request now. Please try after sometime"))
// 								return
// 							} else {
// 								lRespRec.OrderStatus = "Order saved in offline,it will be process when bidding will start."
// 								lRespRec.Status = lResp.Status
// 								lRespRec.ErrMsg = lResp.ErrMsg
// 							}
// 						} else {
// 							log.Println(lReqRec.MasterId, "lReqRec.MasterId")
// 							lTodayAvailable, lErr7 := validatesgb.CheckSgbEndDate(lReqRec.MasterId)
// 							if lErr7 != nil {
// 								log.Println("SGBPO07", lErr7.Error())
// 								lRespRec.Status = common.ErrorCode
// 								fmt.Fprintf(w, helpers.GetErrorString("SGBPO07", "Unable to process your request now. Please try after sometime"))
// 								return
// 							} else {

// 								// check the symbol is available or not
// 								if lTodayAvailable == "True" {

// 									lIndicator, lErr8 := validatesgb.GetSgbTime(lReqRec.MasterId)
// 									log.Println("lIndicator", lIndicator)
// 									if lErr8 != nil {
// 										log.Println("SGBPO08", lErr8.Error())
// 										lRespRec.Status = common.ErrorCode
// 										fmt.Fprintf(w, helpers.GetErrorString("SGBPO08", "Unable to process your request now. Please try after sometime"))
// 										return
// 									} else {
// 										//----------------------------------------------------------------
// 										if lIndicator == "True" {
// 											lResp, lErr9 := UpdateSgbPending(lExchangeReq, lClientId, lReqRec, r, lExchange)
// 											if lErr9 != nil {
// 												log.Println("SGBPO09", lErr9.Error())
// 												lRespRec.Status = common.ErrorCode
// 												fmt.Fprintf(w, helpers.GetErrorString("SGBPO09", "Unable to process your request now. Please try after sometime"))
// 												return
// 											} else {
// 												lRespRec.OrderStatus = "Order Saved in Offline,it will be process on next working day!"
// 												lRespRec.Status = lResp.Status
// 												lRespRec.ErrMsg = lResp.ErrMsg
// 											}
// 											//-----------------------------------------------------------
// 										} else if lIndicator == "False" {

// 											// if lExchange == common.BSE {
// 											// log.Println("if ", lExchange)
// 											// log.Println("lExchangeReq ", lExchangeReq)
// 											// lExchangeResp, lFlag, lErrorRec = BsePlaceOrder(lExchangeReq, lClientId, lReqRec, r, lExchange)
// 											// if lErrorRec != nil {
// 											// 	log.Println("SGBPO10", lErrorRec)
// 											// 	lRespRec.Status = common.ErrorCode
// 											// 	fmt.Fprintf(w, helpers.GetErrorString("SGBPO10", "Unable to process your request now. Please try after sometime"))
// 											// 	return
// 											// }
// 											// } else if lExchange == common.NSE {
// 											// log.Println("else if ", lExchange)
// 											// log.Println("lExchangeReq ", lExchangeReq)
// 											// lExchangeResp, lFlag, lErrorRec = NsePlaceOrder(lExchangeReq, lClientId, lReqRec, r)
// 											// if lErrorRec != nil {
// 											// 	log.Println("SGBPO11", lErrorRec)
// 											// 	lRespRec.Status = common.ErrorCode
// 											// 	fmt.Fprintf(w, helpers.GetErrorString("SGBPO11", "Unable to process your request now. Please try after sometime"))
// 											// 	return
// 											// } else {
// 											// 	log.Println("lExchangeResp", lExchangeResp)
// 											// }
// 											// }

// 											lExchangeResp, lFlag, lErr8 = ProcessSgbOrder(lExchangeReq, lClientId, lReqRec, r, lExchange)
// 											if lErrorRec != nil {
// 												log.Println("SGBPO10", lErrorRec)
// 												lRespRec.Status = common.ErrorCode
// 												fmt.Fprintf(w, helpers.GetErrorString("SGBPO10", "Unable to process your request now. Please try after sometime"))
// 												return
// 											} else {

// 												if lFlag == "S" {
// 													lRespRec = getOrderStatus(lExchangeResp)
// 													log.Println("lRespRec", lRespRec)
// 													if lRespRec.Status == "S" {
// 														sgbClientDetails, lErr12 := fetchSgbClientorder(lClientId, lExchangeResp)
// 														if lErr12 != nil {
// 															log.Println("SGBPO12", lErr12.Error())
// 															lRespRec.Status = common.ErrorCode
// 															fmt.Fprintf(w, helpers.GetErrorString("SGBPO10", "Client Details Not found"))
// 															return
// 														} else {
// 															lSucessMail, lErr13 := constructSuccessmail(sgbClientDetails, lRespRec.Status)
// 															if lErr13 != nil {
// 																log.Println("SGBPO13", lErr13.Error())
// 																lRespRec.Status = common.ErrorCode
// 																fmt.Fprintf(w, helpers.GetErrorString("SGBPO13", "Error in Getting Datas"))
// 																return
// 															} else {
// 																lString := "SGBOrder"
// 																lErr14 := emailUtil.SendEmail(lSucessMail, lString)
// 																if lErr14 != nil {
// 																	log.Println("SGBPO14", lErr14.Error())
// 																	lRespRec.Status = common.ErrorCode
// 																	fmt.Fprintf(w, helpers.GetErrorString("SGBPO14", "Error in Sending Mail"))
// 																	return
// 																}
// 															}
// 														}
// 													}
// 												} else if lFlag == "R" {
// 													lRespRec.ErrMsg = "Order Failed,Please try after sometime!"
// 													lRespRec.Status = common.ErrorCode
// 												} else if lFlag == "E" {
// 													lRespRec.ErrMsg = "Unable to proceed ,Please try after sometime!"
// 													lRespRec.Status = common.ErrorCode
// 												} else if lFlag == "F" {
// 													log.Println("order failed ********************************")
// 													lRespRec = getOrderStatus(lExchangeResp)
// 												}
// 											}
// 										}
// 									}
// 								} else {
// 									lRespRec.Status = common.ErrorCode
// 									lRespRec.ErrMsg = "Timing Closed for SGB"
// 									log.Println("Timing Closed for SGB")
// 								}
// 							}
// 						}

// 					}
// 				}
// 			}
// 		}
// 		lData, lErr12 := json.Marshal(lRespRec)
// 		if lErr12 != nil {
// 			log.Println("SGBPO12", lErr12)
// 			fmt.Fprintf(w, helpers.GetErrorString("SGBPO12", "Unable to getting response.."))
// 			return
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("SgbPlaceOrder (-)", r.Method)
// 	}
// }

// func ProcessSgbOrder(pExchangeReq bsesgb.SgbReqStruct, pClientId string, pReqRec SgbReqStruct, r *http.Request, pExchange string) (bsesgb.SgbRespStruct, string, error) {
// 	log.Println("ProcessSgbOrder (+)")

// 	var lStatus string
// 	var lExchangeResp bsesgb.SgbRespStruct
// 	var lJvStatusRec JvStatusStruct
// 	//Call InsertSgbBidTrack to insert the request
// 	lInsertedId, lErr1 := InsertSgbBidTrack(pExchangeReq, pClientId, pExchange)
// 	if lErr1 != nil {
// 		log.Println("SGBPPSO01", lErr1.Error())
// 		return lExchangeResp, lStatus, lErr1
// 	} else {
// 		log.Println("lInsertedId", lInsertedId)
// 		lJVprocessData, lErr2 := getJvData(pClientId, pReqRec, pExchangeReq, "")
// 		if lErr2 != nil {
// 			log.Println("SGBPPSO02", lErr2.Error())
// 			return lExchangeResp, lStatus, lErr2
// 		} else {
// 			log.Println("lJVprocessData", lJVprocessData)
// 			lClientEmail, lErr3 := clientDetail.GetClientEmailId(r, pClientId)
// 			if lErr3 != nil {
// 				log.Println("SGBPPSO03", lErr3.Error())
// 				return lExchangeResp, lStatus, lErr3
// 			} else {

// 				log.Println("lClientEmail", lClientEmail)
// 				var lErr4 error
// 				// this method is used to call the exchange and update the details in database
// 				lExchangeResp, lErr4 = exchangecall.ApplySgb(pExchangeReq, pClientId, pExchange)
// 				if lErr4 != nil {
// 					log.Println("SGBPPSO04", lErr4.Error())
// 					return lExchangeResp, lStatus, lErr4
// 				} else {
// 					// log.Println(lExchangeResp, "lExchangeResp", lExchangeResp.Bids, "lExchangeResp.Bids", lExchangeResp.StatusCode, "lExchangeResp.StatusCode", lExchangeResp.Bids[0].ErrorCode, "lExchangeResp.Bids[lIdx].ErrorCode")

// 					//--------------success or fail
// 					for lIdx := 0; lIdx < len(lExchangeResp.Bids); lIdx++ {
// 						if lExchangeResp.Bids[lIdx].ErrorCode == "0" && lExchangeResp.StatusCode == "0" {

// 							lJvStatusRec = BlockClientFund(lJVprocessData, r)
// 							if lJvStatusRec.JvStatus == "E" {
// 								log.Println("-------------------------------------JVEROR----------------------------------------")
// 								lErr5 := ResponseUpdate(lExchangeResp, pExchangeReq, pReqRec, pClientId, lJvStatusRec, lClientEmail, pExchange)
// 								if lErr5 != nil {
// 									log.Println("SGBPPSO05", lErr5.Error())
// 									return lExchangeResp, lStatus, lErr5
// 								}
// 								// } else {
// 								lFlag, lErr4 := ReverseProcess(pExchangeReq, lExchangeResp, pReqRec, pClientId, lJvStatusRec, lJVprocessData, lClientEmail, pExchange)
// 								if lErr4 != nil {
// 									log.Println("SGBPPSO06", lErr4.Error())
// 									return lExchangeResp, lStatus, lErr4
// 								} else {
// 									log.Println("lFlag", lFlag)
// 									if lFlag == "E" {
// 										lStatus = common.ErrorCode
// 									} else if lFlag == "S" {
// 										lStatus = "R"
// 									}
// 								}
// 							} else {
// 								lErr5 := ResponseUpdate(lExchangeResp, pExchangeReq, pReqRec, pClientId, lJvStatusRec, lClientEmail, pExchange)
// 								if lErr5 != nil {
// 									log.Println("SGBPPSO07", lErr5.Error())
// 									return lExchangeResp, lStatus, lErr5
// 								} else {
// 									lStatus = common.SuccessCode
// 									log.Println("lStatus", lStatus)
// 								}
// 							}
// 						} else {
// 							lStatus = "F"
// 						}
// 						lErr2 := UpdateSgbBidTrack(lExchangeResp, lInsertedId, pClientId, lJvStatusRec, pExchange)
// 						if lErr2 != nil {
// 							log.Println("SGBPRU02", lErr2.Error())
// 							return lExchangeResp, lStatus, lErr2
// 						} else {
// 							log.Println("Updated Successfully in bid track")
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~", lStatus)
// 	log.Println("ProcessSgbOrder (-)")
// 	return lExchangeResp, lStatus, nil
// }

// /*
// Purpose:This method updating the Status As Pending in details table.
// Parameters:

// 	pRespArr,pClientId,pHeaderId,pDetailIdArr

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
// func UpdateSgbPending(pExchangeReq bsesgb.SgbReqStruct, pClientId string, pReqRec SgbReqStruct, r *http.Request, pExchange string) (SgbRespStruct, error) {
// 	log.Println("UpdateSgbPending (+)")

// 	var lRespRec SgbRespStruct

// 	lJVprocessData, lErr1 := getJvData(pClientId, pReqRec, pExchangeReq, "")
// 	if lErr1 != nil {
// 		log.Println("SGBPUSO01", lErr1.Error())
// 		lRespRec.Status = common.ErrorCode
// 		lRespRec.ErrMsg = lErr1.Error()
// 		return lRespRec, lErr1
// 	} else {
// 		// log.Println("lJVprocessData in pending", lJVprocessData)
// 		//JV
// 		lJvStatusRec := BlockClientFund(lJVprocessData, r)
// 		if lJvStatusRec.JvStatus == "E" {
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "Unable to deduct amount in your trading account,Please try after sometime!"
// 			return lRespRec, nil
// 		} else {
// 			lClientEmail, lErr2 := clientDetail.GetClientEmailId(r, pClientId)
// 			if lErr2 != nil {
// 				log.Println("SGBPUSO02", lErr2.Error())
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = lErr2.Error()
// 				return lRespRec, lErr2
// 			} else {
// 				log.Println("lClientEmail", lClientEmail)
// 				lErr2 = UpdatePendingToLocal(pExchangeReq, pClientId, pReqRec.MasterId, lClientEmail, lJvStatusRec, pExchange)
// 				if lErr2 != nil {
// 					log.Println("SGBPUSO03", lErr2.Error())
// 					lRespRec.Status = common.ErrorCode
// 					lRespRec.ErrMsg = lErr2.Error()
// 					return lRespRec, lErr2
// 				} else {
// 					lRespRec.Status = common.SuccessCode
// 					lRespRec.ErrMsg = "Order Saved in Offline,it will be process on next working day!"
// 				}
// 			}
// 		}
// 	}
// 	log.Println("UpdateSgbPending (-)")
// 	return lRespRec, nil
// }

// /*
// Pupose:This method inserting the order head values in order header table.
// Parameters:

// 	pReqArr,pMasterId,PClientId

// Response:

// 		==========
// 		*On Sucess
// 		==========

// 		==========
// 		*On Error
// 		==========

// Author:Pavithra
// Date: 12JUNE2023
// */
// func UpdatePendingToLocal(pExchangeReq bsesgb.SgbReqStruct, pClientId string, pMasterId int, pCLientEmail string, pJvRec JvStatusStruct, pExchange string) error {
// 	log.Println("UpdatePendingToLocal (+)")

// 	// for Set the cancel flag as N
// 	lCancelFlag := "N"
// 	// establish a database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SGBPUP01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		// loop the req array to insert in database
// 		// for lpReqArrIdx := 0; lpReqArrIdx < len(pReqArr); lpReqArrIdx++ {

// 		lSqlString := `insert into a_sgb_orderheader (MasterId,ScripId ,PanNo,InvestorCategory ,ApplicantName ,Depository ,
// 					DpId ,ClientBenfId ,GuardianName ,GuardianPanNo ,GuardianRelation ,Status ,ClientId ,cancelFlag ,ClientEmail,Exchange ,CreatedBy ,
// 					CreatedDate ,UpdatedBy ,UpdatedDate )
// 					values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

// 		lInsertedHeaderId, lErr2 := lDb.Exec(lSqlString, pMasterId, pExchangeReq.ScripId, pExchangeReq.PanNo, pExchangeReq.InvestorCategory, pExchangeReq.ApplicantName, pExchangeReq.Depository, pExchangeReq.DpId, pExchangeReq.ClientBenfId, pExchangeReq.GuardianName, pExchangeReq.GuardianPanno, pExchangeReq.GuardianRelation, common.PENDING, pClientId, lCancelFlag, pCLientEmail, pExchange, pClientId, pClientId)
// 		if lErr2 != nil {
// 			log.Println("Error updating pending status into database (header)")
// 			log.Println("SGBPUP02", lErr2)
// 			return lErr2
// 		} else {
// 			// get lastinserted id in lReturnId and converet them into int ,store it in lHeaderId
// 			lReturnId, _ := lInsertedHeaderId.LastInsertId()
// 			lHeaderId := int(lReturnId)

// 			// Looping the Bid Details Array to Update the Application Response Values in Order Details Table
// 			lBidArr := pExchangeReq.Bids
// 			for lBidArrIdx := 0; lBidArrIdx < len(lBidArr); lBidArrIdx++ {

// 				// Check whether the requested bid is activity type is new,if it is new insert the bid details in order detail table
// 				if lBidArr[lBidArrIdx].ActionCode == "N" {
// 					lSqlString := `insert into a_sgb_orderdetails (HeaderId,BidId,OrderNo,ActionCode,ReqSubscriptionUnit,ReqRate,
// 							JvStatus,JvAmount,JvStatement,JvType,Exchange,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
// 							values(?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

// 					_, lErr3 := lDb.Exec(lSqlString, lHeaderId, lBidArr[lBidArrIdx].BidId, lBidArr[lBidArrIdx].OrderNo, lBidArr[lBidArrIdx].ActionCode, lBidArr[lBidArrIdx].SubscriptionUnit, lBidArr[lBidArrIdx].Rate, pJvRec.JvStatus, pJvRec.JvAmount, pJvRec.JvStatement, pJvRec.JvType, common.BSE, pClientId, pClientId)
// 					if lErr3 != nil {
// 						log.Println("SGBPUP03", lErr3)
// 						return lErr3
// 					} else {
// 						log.Println("Details Inserted Successfully", lBidArrIdx)
// 						lSqlString := `insert into a_sgbtracking_table (OrderNo,BidId,ActivityType,Unit,Price,ApplicationStatus,ClientId,Exchange,
// 								JvAmount,JvStatus,JvStatement,JvType,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
// 								values(?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

// 						_, lErr4 := lDb.Exec(lSqlString, lBidArr[lBidArrIdx].OrderNo, lBidArr[lBidArrIdx].BidId, lBidArr[lBidArrIdx].ActionCode, lBidArr[lBidArrIdx].SubscriptionUnit, lBidArr[lBidArrIdx].Rate, common.PENDING, pClientId, common.BSE, pJvRec.JvAmount, pJvRec.JvStatus, pJvRec.JvStatement, pJvRec.JvType, pClientId, pClientId)
// 						if lErr4 != nil {
// 							log.Println("SGBPUP04", lErr4)
// 							return lErr4
// 						}
// 					}
// 					// Check whether the requested bid is activity type is modify or cancel,then update the bid details in order detail table
// 				} else if lBidArr[lBidArrIdx].ActionCode == "M" || lBidArr[lBidArrIdx].ActionCode == "D" {

// 					SqlString := `update a_sgb_orderdetails d
// 						set d.BidId = ?,d.OrderNo = ?,d.ActionCode = ?,d.ReqSubscriptionUnit = ?,d.ReqRate = ?,
// 						d.JvAmount = ?,d.JvStatus = ?,d.JvStatement = ?,d.JvType = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
// 						where d.HeaderId = ? and d.Exchange = ?`

// 					_, lErr5 := lDb.Exec(SqlString, lBidArr[lBidArrIdx].BidId, lBidArr[lBidArrIdx].OrderNo, lBidArr[lBidArrIdx].ActionCode, lBidArr[lBidArrIdx].SubscriptionUnit, lBidArr[lBidArrIdx].Rate, pJvRec.JvAmount, pJvRec.JvStatus, pJvRec.JvStatement, pJvRec.JvType, pClientId, pMasterId, common.BSE)
// 					if lErr5 != nil {
// 						log.Println("SGBPUP05", lErr5)
// 						return lErr5
// 					} else {
// 						lSqlString := `insert into a_sgbtracking_table (OrderNo,BidId,ActivityType,Unit,Price,ApplicationStatus,ClientId,Exchange,
// 								JvAmount,JvStatus,JvStatement,JvType,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
// 								values(?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

// 						_, lErr6 := lDb.Exec(lSqlString, lBidArr[lBidArrIdx].OrderNo, lBidArr[lBidArrIdx].BidId, lBidArr[lBidArrIdx].ActionCode, lBidArr[lBidArrIdx].SubscriptionUnit, lBidArr[lBidArrIdx].Rate, common.PENDING, pClientId, common.BSE, pJvRec.JvAmount, pJvRec.JvStatus, pJvRec.JvStatement, pJvRec.JvType, pClientId, pClientId)
// 						if lErr6 != nil {
// 							log.Println("SGBPUP06", lErr6)
// 							return lErr6
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("UpdatePendingToLocal (-)")
// 	return nil
// }

// /*
// Pupose:This method inserting the order head values in order header table.
// Parameters:

// 	pReqArr,pMasterId,PClientId

// Response:

// 		==========
// 		*On Sucess
// 		==========

// 		==========
// 		*On Error
// 		==========

// Author:Pavithra
// Date: 12JUNE2023
// */
// func ResponseUpdate(pExchangeResp bsesgb.SgbRespStruct, pExchangeReq bsesgb.SgbReqStruct, pReqRec SgbReqStruct, pClientId string, pJvStatusRec JvStatusStruct, pEmailId string, pExchange string) error {
// 	log.Println("ResponseUpdate (+)")

// 	for lIdx := 0; lIdx < len(pExchangeResp.Bids); lIdx++ {
// 		if pExchangeResp.Bids[lIdx].ErrorCode == "0" && pExchangeResp.StatusCode == "0" {
// 			lErr1 := InsertSGBHeader(pExchangeResp, pExchangeReq, pReqRec.MasterId, pClientId, pJvStatusRec, pEmailId, pExchange)
// 			if lErr1 != nil {
// 				log.Println("SGBPRU01", lErr1.Error())
// 				return lErr1
// 			} else {
// 				log.Println("Values Updated Successfully")
// 			}
// 		}
// 	}
// 	log.Println("ResponseUpdate (-)")
// 	return nil
// }

// /*

// Pupose:This method is used to get the client BenId and ClientName for order input.
// Parameters:

// 	PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	123456789012,Lakshmanan Ashok Kumar,nil

// 	==========
// 	*On Error
// 	==========
// 	"","",error

// Author:Pavithra
// Date: 12JUNE2023
// */
func getJvData(pClientId string, pReqRec SgbReqStruct, pExchangeReq bsesgb.SgbReqStruct, pCheckFlag string) (JvDataStruct, error) {
	log.Println("getJvData (+)")

	// this variable is used to get the JV process struct
	var lJvData JvDataStruct

	//check the action code if it is "N" directly assign the values or else call getOrderDetails method
	if pReqRec.ActionCode == "N" {

		//check the flag is empty or not to do reverse process
		if pCheckFlag != "" {
			lJvData.ActionCode = "D"
			lJvData.ClientId = pClientId
			lJvData.MasterId = pReqRec.MasterId
			lJvData.Unit = strconv.Itoa(pReqRec.OldUnit)
			lJvData.Price = strconv.Itoa(pReqRec.Price)
			lJvData.BidId = pReqRec.BidId
			lJvData.JVamount = strconv.Itoa(pReqRec.OldUnit * pReqRec.Price)
			lJvData.OrderNo = pReqRec.OrderNo
		} else {
			lJvData.ActionCode = pReqRec.ActionCode
			lJvData.BidId = ""
			lJvData.ClientId = pClientId
			lJvData.MasterId = pReqRec.MasterId
			lJvData.Unit = strconv.Itoa(pReqRec.Unit)
			lJvData.Price = strconv.Itoa(pReqRec.Price)
			lJvData.Transaction = "C"
			lJvData.JVamount = strconv.Itoa(pReqRec.Unit * pReqRec.Price)
			lJvData.Flag = "Purchase"
			for lIdx := 0; lIdx < len(pExchangeReq.Bids); lIdx++ {
				lJvData.OrderNo = pExchangeReq.Bids[lIdx].OrderNo
			}
		}

	} else {
		// lJvReq, lErr1 := getOrderedDetails(pClientId, pReqRec, pCheckFlag)
		// if lErr1 != nil {
		// 	log.Println("SPGJD01", lErr1)
		// 	return lJvData, lErr1
		// } else {
		// 	lJvData = lJvReq
		// }
	}

	log.Println("getJvData (-)")
	return lJvData, nil
}

// /*
// Pupose:This method is used to get the client BenId and ClientName for order input.
// Parameters:

// 	PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	123456789012,Lakshmanan Ashok Kumar,nil

// 	==========
// 	*On Error
// 	==========
// 	"","",error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func getOrderStatus(pExchangeResp bsesgb.SgbRespStruct) SgbRespStruct {
// 	log.Println("getOrderStatus (+)")

// 	var lOrderResp SgbRespStruct

// 	for lIdx := 0; lIdx < len(pExchangeResp.Bids); lIdx++ {
// 		if pExchangeResp.StatusCode == "0" && pExchangeResp.Bids[lIdx].ErrorCode == "0" {
// 			if pExchangeResp.Bids[lIdx].ErrorCode == "0" && pExchangeResp.Bids[lIdx].ActionCode == "N" {
// 				lOrderResp.OrderStatus = "Order Placed Successfully"
// 				lOrderResp.Status = "S"
// 			} else if pExchangeResp.Bids[lIdx].ErrorCode == "0" && pExchangeResp.Bids[lIdx].ActionCode == "M" {
// 				lOrderResp.OrderStatus = "Order Modified Successfully"
// 				lOrderResp.Status = "S"
// 			} else if pExchangeResp.Bids[lIdx].ErrorCode == "0" && pExchangeResp.Bids[lIdx].ActionCode == "D" {
// 				lOrderResp.OrderStatus = "Order Deleted Successfully"
// 				lOrderResp.Status = "S"
// 			}
// 		} else {
// 			log.Println("inside else ")
// 			// for lIdx := 0; lIdx < len(pExchangeResp.Bids); lIdx++ {
// 			if pExchangeResp.Bids[lIdx].ErrorCode != "0" && pExchangeResp.StatusCode != "0" {
// 				log.Println("inside else if ")
// 				lOrderResp.OrderStatus = pExchangeResp.Bids[lIdx].Message + "/" + pExchangeResp.StatusMessage
// 				lOrderResp.Status = "E"
// 			} else if pExchangeResp.Bids[lIdx].ErrorCode == "0" && pExchangeResp.StatusCode != "0" {
// 				log.Println("inside else if 1")
// 				lOrderResp.OrderStatus = pExchangeResp.StatusMessage
// 				lOrderResp.Status = "E"
// 			} else if pExchangeResp.Bids[lIdx].ErrorCode != "0" && pExchangeResp.StatusCode == "0" {
// 				log.Println("inside else if 2")
// 				lOrderResp.OrderStatus = pExchangeResp.Bids[lIdx].Message
// 				lOrderResp.Status = "E"
// 			} else {

// 			}
// 		}
// 	}
// 	log.Println(lOrderResp)
// 	log.Println("getOrderStatus (-)")
// 	return lOrderResp
// }

// /*
// Pupose:This method is used to get the client BenId and ClientName for order input.
// Parameters:

// 	PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	123456789012,Lakshmanan Ashok Kumar,nil

// 	==========
// 	*On Error
// 	==========
// 	"","",error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func getOrderedDetails(pClientId string, pReqRec SgbReqStruct, pCheckFlag string) (JvDataStruct, error) {
// 	log.Println("getOrderedDetails (+)")

// 	// this variables is used to get DpId and ClientName from the database.
// 	var lJvData JvDataStruct
// 	var lJvUnit int

// 	//check the flag is empty or not to do reverse process
// 	if pCheckFlag != "" {
// 		log.Println("--------------INSIDE REVERSE PROCESS-------------------------")
// 		lJvData.MasterId = pReqRec.MasterId
// 		lJvData.OrderNo = pReqRec.OrderNo
// 		lJvData.ClientId = pClientId
// 		lJvData.Price = strconv.Itoa(pReqRec.Price)
// 		//check the requested unit is less than or equal to already placed unit to addon the qty in reverse process
// 		if pReqRec.Unit <= pReqRec.OldUnit {
// 			lJvData.Transaction = "F"
// 			lJvData.Unit = strconv.Itoa(pReqRec.OldUnit)
// 			lJvData.JVamount = strconv.Itoa(pReqRec.OldUnit * pReqRec.Price)
// 			lJvData.ActionCode = "M"
// 			lJvData.BidId = pReqRec.BidId
// 			lJvData.Flag = "Addon"

// 			//if the unit equals to zero then it is a cancel type then its is changes to reverse process as new type
// 			if pReqRec.OldUnit == pReqRec.Unit {

// 				lJvData.Unit = strconv.Itoa(pReqRec.OldUnit)
// 				lJvData.JVamount = strconv.Itoa(pReqRec.OldUnit * pReqRec.Price)
// 				lJvData.ActionCode = "N"
// 				lJvData.BidId = ""
// 				lTime := time.Now()
// 				lUnixTime := lTime.Unix()
// 				lUnixTimeString := fmt.Sprintf("%d", lUnixTime)
// 				var lTrimmedString string
// 				if len(pClientId) >= 5 {
// 					lTrimmedString = pClientId[len(pClientId)-5:]
// 				}
// 				lJvData.OrderNo = lUnixTimeString + lTrimmedString
// 				lJvData.Flag = "Purchase"
// 			}
// 		} else {
// 			// it checks the unit is added in placed
// 			lJvData.Transaction = "C"
// 			lJvData.Unit = strconv.Itoa(pReqRec.OldUnit)
// 			lJvData.JVamount = strconv.Itoa(pReqRec.OldUnit * pReqRec.Price)
// 			lJvData.ActionCode = "M"
// 			lJvData.BidId = pReqRec.BidId
// 			lJvData.Flag = "Reversal"
// 		}
// 	} else {

// 		lJvData.MasterId = pReqRec.MasterId
// 		lJvData.OrderNo = pReqRec.OrderNo
// 		lJvData.ClientId = pClientId
// 		lJvData.Price = strconv.Itoa(pReqRec.Price)
// 		lJvData.ActionCode = pReqRec.ActionCode

// 		//if the requested qty is less than the placed qty then it is a reversal process
// 		if pReqRec.Unit <= pReqRec.OldUnit {
// 			lJvData.Transaction = "F"
// 			lJvUnit = pReqRec.OldUnit - pReqRec.Unit
// 			lJvData.Unit = strconv.Itoa(lJvUnit)
// 			lJvData.JVamount = strconv.Itoa(lJvUnit * pReqRec.Price)
// 			lJvData.Flag = "Reversal"

// 			//if the unit equals to zero then it is a cancel type
// 			if lJvUnit == 0 {
// 				lJvUnit = pReqRec.Unit
// 				lJvData.Unit = strconv.Itoa(lJvUnit)
// 				lJvData.JVamount = strconv.Itoa(lJvUnit * pReqRec.Price)
// 				lJvData.Flag = "Cancel"
// 			}
// 		} else {
// 			// the requested qty is greater than the placed qty it is a Addon Process
// 			lJvData.Transaction = "C"
// 			lJvUnit = pReqRec.Unit - pReqRec.OldUnit
// 			lJvData.Unit = strconv.Itoa(lJvUnit)
// 			lJvData.JVamount = strconv.Itoa(lJvUnit * pReqRec.Price)
// 			lJvData.Flag = "Addon"
// 		}
// 	}
// 	log.Println("getOrderedDetails (-)")
// 	return lJvData, nil
// }

// /*
// Pupose:This method inserting the order head values in order header table.
// Parameters:

// 	pReqArr,pMasterId,PClientId

// Response:

// 		==========
// 		*On Sucess
// 		==========

// 		==========
// 		*On Error
// 		==========

// Author:Pavithra
// Date: 12JUNE2023
// */
// func InsertSGBHeader(pRespRec bsesgb.SgbRespStruct, pReqRec bsesgb.SgbReqStruct, pMasterId int, pClientId string, pJvStatusRec JvStatusStruct, pEmailId string, pExchange string) error {
// 	log.Println("InsertSGBHeader (+)")

// 	// get the application no id in table
// 	var lHeaderId int
// 	//set cancel Flag as N
// 	lCancelFlag := "N"
// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("PIH01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		// Looping the pReqArr Array value to inserting the request datas in databse.

// 		//  call the getAppNoId method, To get Application No Id from the database
// 		lHeadId, lErr2 := GetOrderId(pRespRec, pReqRec, pMasterId, pClientId)
// 		if lErr2 != nil {
// 			log.Println("PIH02", lErr2)
// 			return lErr2
// 		} else {
// 			if lHeadId != 0 {
// 				lHeaderId = lHeadId
// 			} else {
// 				lHeaderId = pMasterId
// 			}

// 			for lBidIdx := 0; lBidIdx < len(pRespRec.Bids); lBidIdx++ {

// 				//if appliaction no is not present in the database,insert the new application details in orderheader table
// 				if lHeadId == 0 {
// 					lSqlString1 := `insert into a_sgb_orderheader (MasterId ,ScripId,PanNo,InvestorCategory,ApplicantName,Depository,
// 									DpId,ClientBenfId,GuardianName,GuardianPanNo,GuardianRelation,StatusCode,StatusMessage,ErrorCode,ErrorMessage,
// 									Status,ClientId,cancelFlag,Exchange,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,ClientEmail)
// 									values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now(),?)`

// 					lInsertedHeaderId, lErr3 := lDb.Exec(lSqlString1, lHeaderId, pRespRec.ScripId, pRespRec.PanNo, pRespRec.InvestorCategory, pRespRec.ApplicantName, pRespRec.Depository, pRespRec.DpId, pRespRec.ClientBenfId, pRespRec.GuardianName, pRespRec.GuardianPanno, pRespRec.GuardianRelation, pRespRec.StatusCode, pRespRec.StatusMessage, pRespRec.ErrorCode, pRespRec.ErrorMessage, common.SUCCESS, pClientId, lCancelFlag, pExchange, pClientId, pClientId, pEmailId)
// 					if lErr3 != nil {
// 						log.Println("PIH03", lErr3)
// 						return lErr3
// 					} else {
// 						// get lastinserted id in lReturnId and converet them into int ,store it in lHeaderId
// 						lReturnId, _ := lInsertedHeaderId.LastInsertId()
// 						lHeaderId = int(lReturnId)
// 						log.Println("lHeaderId", lHeaderId)

// 						// call InsertDetails method to inserting the order details in order details table
// 						lErr4 := InsertSGBDetail(pRespRec.Bids, pReqRec.Bids, lHeaderId, pClientId, pJvStatusRec, pExchange)
// 						if lErr4 != nil {
// 							log.Println("PIH04", lErr4)
// 							return lErr4
// 						} else {
// 							log.Println("header inserted successfully")
// 						}
// 					}
// 				} else if pRespRec.Bids[lBidIdx].ActionCode == "D" && pRespRec.StatusCode == "0" {
// 					lCancelFlag = "Y"
// 					lSqlString2 := `update a_sgb_orderheader h
// 									set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
// 									h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.StatusCode = ?,
// 									h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,h.Status = ?,h.ClientId = ?,h.cancelFlag = ?,
// 									h.UpdatedBy = ?,h.UpdatedDate = now()
// 									where h.Id = ?
// 									and h.ClientId = ?
// 									and h.Exchange = ?`

// 					_, lErr5 := lDb.Exec(lSqlString2, pRespRec.ScripId, pRespRec.PanNo, pRespRec.InvestorCategory, pRespRec.ApplicantName, pRespRec.Depository, pRespRec.DpId, pRespRec.ClientBenfId, pRespRec.GuardianName, pRespRec.GuardianPanno, pRespRec.GuardianRelation, pRespRec.StatusCode, pRespRec.StatusMessage, pRespRec.ErrorCode, pRespRec.ErrorMessage, common.SUCCESS, pClientId, lCancelFlag, pClientId, lHeaderId, pClientId, pExchange)
// 					if lErr5 != nil {
// 						log.Println("PIH05", lErr5)
// 						return lErr5
// 					} else {
// 						// call InsertDetails method to inserting the order details in order details table
// 						lErr6 := InsertSGBDetail(pRespRec.Bids, pReqRec.Bids, lHeaderId, pClientId, pJvStatusRec, pExchange)
// 						if lErr6 != nil {
// 							log.Println("PIH06", lErr6)
// 							return lErr6
// 						} else {
// 							log.Println("header cancel updated successfully")
// 						}
// 					}
// 				} else {
// 					log.Println("else")
// 					lSqlString3 := `update a_sgb_orderheader h
// 								set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
// 								h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,h.StatusCode = ?,
// 								h.StatusMessage = ?,h.ErrorCode = ?,h.ErrorMessage = ?,h.Status = ?,h.ClientId = ?,h.cancelFlag = ?,
// 								h.UpdatedBy = ?,h.UpdatedDate = now()
// 								where h.Id = ?
// 								and h.ClientId = ?
// 								and h.Exchange = ?`

// 					_, lErr5 := lDb.Exec(lSqlString3, pRespRec.ScripId, pRespRec.PanNo, pRespRec.InvestorCategory, pRespRec.ApplicantName, pRespRec.Depository, pRespRec.DpId, pRespRec.ClientBenfId, pRespRec.GuardianName, pRespRec.GuardianPanno, pRespRec.GuardianRelation, pRespRec.StatusCode, pRespRec.StatusMessage, pRespRec.ErrorCode, pRespRec.ErrorMessage, common.SUCCESS, pClientId, lCancelFlag, pClientId, lHeaderId, pClientId, pExchange)
// 					if lErr5 != nil {
// 						log.Println("PIH07", lErr5)
// 						return lErr5
// 					} else {
// 						// call InsertDetails method to inserting the order details in order details table
// 						lErr6 := InsertSGBDetail(pRespRec.Bids, pReqRec.Bids, lHeaderId, pClientId, pJvStatusRec, pExchange)
// 						if lErr6 != nil {
// 							log.Println("PIH08", lErr6)
// 							return lErr6
// 						} else {
// 							log.Println("header updated successfully")
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("InsertSGBHeader (-)")
// 	return nil
// }

// /*
// Pupose:This method inserting the bid details in order detail table.
// Parameters:

// 	pReqBidArr,pHeaderId,PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	[1,2,3],[1,2,3],nil

// 	==========
// 	*On Error
// 	==========
// 	[],[],error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func InsertSGBDetail(pRespBidArr []bsesgb.RespSgbBidStruct, pReqBidArr []bsesgb.ReqSgbBidStruct, pHeaderId int, pClientId string, pJvStatusRec JvStatusStruct, pExchange string) error {
// 	log.Println("InsertSGBDetail (+)")

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SPISD01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		log.Println("pReqBidArr", pReqBidArr)
// 		log.Println("pRespBidArr", pRespBidArr)
// 		// Looping the pReqBidArr array to insert the Request bid details in database
// 		for lReqBidIdx := 0; lReqBidIdx < len(pReqBidArr); lReqBidIdx++ {

// 			for lRespBidIdx := 0; lRespBidIdx < len(pRespBidArr); lRespBidIdx++ {

// 				// Check whether the requested bid is activity type is new,if it is new insert the bid details in order detail table
// 				if pRespBidArr[lRespBidIdx].ActionCode == "N" {
// 					log.Println("pRespBidArr[lRespBidIdx].ActionCode", pRespBidArr[lRespBidIdx].ActionCode)

// 					log.Println("pRespBidArr[lRespBidIdx].OrderNo", pRespBidArr[lRespBidIdx].OrderNo)
// 					log.Println("pReqBidArr[lReqBidIdx].OrderNo ", pReqBidArr[lReqBidIdx].OrderNo)
// 					if pExchange == common.BSE {
// 						if pRespBidArr[lRespBidIdx].OrderNo == pReqBidArr[lReqBidIdx].OrderNo {

// 							lSqlString1 := `insert into a_sgb_orderdetails (HeaderId,BidId,OrderNo,ActionCode,ReqSubscriptionUnit,ReqRate,
// 							RespSubscriptionunit,RespRate,ErrorCode,Message,Exchange,JvAmount,JvStatus,JvStatement,JvType,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
// 							values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

// 							_, lErr3 := lDb.Exec(lSqlString1, pHeaderId, pRespBidArr[lRespBidIdx].BidId, pRespBidArr[lRespBidIdx].OrderNo, pRespBidArr[lRespBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pRespBidArr[lRespBidIdx].SubscriptionUnit, pRespBidArr[lRespBidIdx].Rate, pRespBidArr[lRespBidIdx].ErrorCode, pRespBidArr[lRespBidIdx].Message, pExchange, pJvStatusRec.JvAmount, pJvStatusRec.JvStatus, pJvStatusRec.JvStatement, pJvStatusRec.JvType, pClientId, pClientId)
// 							if lErr3 != nil {
// 								log.Println("SPISD02", lErr3)
// 								return lErr3
// 							} else {
// 								log.Println("Details Inserted Successfully in bse")
// 							}
// 						}
// 					} else {
// 						lSqlString1 := `insert into a_sgb_orderdetails (HeaderId,BidId,OrderNo,ActionCode,ReqSubscriptionUnit,ReqRate,
// 							RespSubscriptionunit,RespRate,ErrorCode,Message,Exchange,JvAmount,JvStatus,JvStatement,JvType,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
// 							values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

// 						_, lErr3 := lDb.Exec(lSqlString1, pHeaderId, pRespBidArr[lRespBidIdx].BidId, pRespBidArr[lRespBidIdx].OrderNo, pRespBidArr[lRespBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pRespBidArr[lRespBidIdx].SubscriptionUnit, pRespBidArr[lRespBidIdx].Rate, pRespBidArr[lRespBidIdx].ErrorCode, pRespBidArr[lRespBidIdx].Message, pExchange, pJvStatusRec.JvAmount, pJvStatusRec.JvStatus, pJvStatusRec.JvStatement, pJvStatusRec.JvType, pClientId, pClientId)
// 						if lErr3 != nil {
// 							log.Println("SPISD03", lErr3)
// 							return lErr3
// 						} else {
// 							log.Println("Details Inserted Successfully in nse")
// 						}
// 					}
// 					// Check whether the requested bid is activity type is modify or cancel,then update the bid details in order detail table
// 				} else if pRespBidArr[lRespBidIdx].ActionCode == "M" {

// 					//Check whether the requested bid and Response bid activity type same (modify or cancel)
// 					if pRespBidArr[lRespBidIdx].OrderNo == pReqBidArr[lReqBidIdx].OrderNo {
// 						// ! ----------------
// 						lSqlString2 := `update a_sgb_orderdetails d
// 										set d.OrderNo = ?,d.ActionCode = ?,d.ReqSubscriptionUnit = ?,d.ReqRate = ?,
// 										d.RespSubscriptionunit = ?,d.RespRate = ?,d.ErrorCode = ?,d.Message = ?,
// 										d.JvAmount = ?,d.JvStatus = ?,d.JvStatement = ?,d.JvType = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
// 										where d.HeaderId = ?
// 										and d.Exchange = ?
// 										and d.OrderNo = ?`

// 						_, lErr4 := lDb.Exec(lSqlString2, pRespBidArr[lRespBidIdx].OrderNo, pRespBidArr[lRespBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pRespBidArr[lRespBidIdx].SubscriptionUnit, pRespBidArr[lRespBidIdx].Rate, pRespBidArr[lRespBidIdx].ErrorCode, pRespBidArr[lRespBidIdx].Message, pJvStatusRec.JvAmount, pJvStatusRec.JvStatus, pJvStatusRec.JvStatement, pJvStatusRec.JvType, pClientId, pHeaderId, pExchange, pRespBidArr[lRespBidIdx].OrderNo)
// 						if lErr4 != nil {
// 							log.Println("SPISD03", lErr4)
// 							return lErr4
// 						} else {
// 							log.Println("Details updated Successfully")
// 						}
// 					}
// 				} else if pRespBidArr[lRespBidIdx].ActionCode == "D" {

// 					//Check whether the requested bid and Response bid activity type same (modify or cancel)
// 					if pRespBidArr[lRespBidIdx].OrderNo == pReqBidArr[lReqBidIdx].OrderNo {
// 						// ! ----------------
// 						lSqlString3 := `update a_sgb_orderdetails d
// 										set d.OrderNo = ?,d.ActionCode = ?,d.ReqSubscriptionUnit = ?,d.ReqRate = ?,
// 										d.RespSubscriptionunit = ?,d.RespRate = ?,d.ErrorCode = ?,d.Message = ?,
// 										d.JvAmount = ?,d.JvStatus = ?,d.JvStatement = ?,d.JvType = ?,d.UpdatedBy = ?,d.UpdatedDate = now()
// 										where d.HeaderId = ?
// 										and d.Exchange = ?
// 										and d.OrderNo = ?`

// 						_, lErr4 := lDb.Exec(lSqlString3, pRespBidArr[lRespBidIdx].OrderNo, pRespBidArr[lRespBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pRespBidArr[lRespBidIdx].SubscriptionUnit, pRespBidArr[lRespBidIdx].Rate, pRespBidArr[lRespBidIdx].ErrorCode, pRespBidArr[lRespBidIdx].Message, pJvStatusRec.JvAmount, pJvStatusRec.JvStatus, pJvStatusRec.JvStatement, pJvStatusRec.JvType, pClientId, pHeaderId, pExchange, pRespBidArr[lRespBidIdx].OrderNo)
// 						if lErr4 != nil {
// 							log.Println("SPISD04", lErr4)
// 							return lErr4
// 						} else {
// 							log.Println("Details updated Successfully")
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("InsertSGBDetail (-)")
// 	return nil
// }

// /*
// Pupose:This method is used to get the Application Number From the DataBase.
// Parameters:

// 	pSymbol,pAppNo,PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	FT000006912345,nil

// 	==========
// 	*On Error
// 	==========
// 	"",error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func GetOrderId(pRespRec bsesgb.SgbRespStruct, pReqRec bsesgb.SgbReqStruct, pMasterId int, pClientId string) (int, error) {
// 	log.Println("GetOrderId (+)")

// 	var lMasterId int

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SPGOI01", lErr1)
// 		return lMasterId, lErr1
// 	} else {
// 		defer lDb.Close()

// 		for lReqIdx := 0; lReqIdx < len(pReqRec.Bids); lReqIdx++ {
// 			for lRespIdx := 0; lRespIdx < len(pRespRec.Bids); lRespIdx++ {
// 				if pReqRec.Bids[lReqIdx].OrderNo == pRespRec.Bids[lRespIdx].OrderNo {
// 					log.Println("OrderNo", pReqRec.Bids[lReqIdx].OrderNo, pRespRec.Bids[lRespIdx].OrderNo)

// 					lCoreString := `select (case when count(1) > 0 then h.Id  else 0 end) Id
// 									from a_sgb_orderheader h,a_sgb_orderdetails d
// 									where h.Id = d.HeaderId
// 									and h.MasterId = ?
// 									and h.ClientId = ?
// 									and d.OrderNo = ?`
// 					lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pClientId, pRespRec.Bids[lRespIdx].OrderNo)
// 					if lErr2 != nil {
// 						log.Println("SPGOI02", lErr2)
// 						return lMasterId, lErr2
// 					} else {
// 						for lRows.Next() {
// 							lErr3 := lRows.Scan(&lMasterId)
// 							if lErr3 != nil {
// 								log.Println("SPGOI03", lErr3)
// 								return lMasterId, lErr3
// 							}
// 						}
// 					}
// 					//
// 					// log.Println(lMasterId)
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GetOrderId (-)")
// 	return lMasterId, nil
// }

// /*
// Purpose:This method is used to construct the exchange header values.
// Parameters:

// 	pReqRec,pClientId

// Response:

// 		==========
// 		*On Sucess
// 		==========
// 	 	[
// 			{
// 				"symbol":"JDIAL",
// 				"applicationNumber":"FT000069130109",
// 				"category":"IND",
// 				"clientName":"LAKSHMANAN ASHOK KUMAR",
// 				"depository":"CDSL",
// 				"dpId":"",
// 				"clientBenId":"1208030000262661",
// 				"nonASBA":false,
// 				"pan":"AGMPA8575C",
// 				"referenceNumber":"2c2e506472cbb2f9",
// 				"allotmentMode":"demat",
// 				"upiFlag":"Y",
// 				"upi":"test@kmbl",
// 				"bankCode":"",
// 				"locationCode":"",
// 				"bankAccount":"",
// 				"ifsc":"",
// 				"subBrokerCode":"",
// 				"timestamp":"0000-00-00",
// 				"bids":[
// 					{
// 						"activityType":"new",
// 						"bidReferenceNumber":"123",
// 						"series":"",
// 						"quantity":19,
// 						"atCutOff":true,
// 						"price":755,
// 						"amount":14345,
// 						"remark":"FT000069130109",
// 						"lotSize":19
// 					}
// 				]
// 			}
// 		]
// 		==========
// 		*On Error
// 		==========
// 		[],error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func ConstructSGBReqStruct(pReqRec SgbReqStruct, pClientId string, pExchange string) (bsesgb.SgbReqStruct, error) {
// 	log.Println("ConstructSGBReqStruct (+)")

// 	// create an instance of lReqRec type bsesgb.SgbReqStruct.
// 	var lReqRec bsesgb.SgbReqStruct

// 	// create an instance of lReqRec type bsesgb.SgbReqStruct.
// 	var lBidReqRec bsesgb.ReqSgbBidStruct

// 	// call the getDpId method to get the lDpId and lClientName
// 	lDpId, lClientName, lErr := getDpId(pClientId)
// 	if lErr != nil {
// 		log.Println("SPCSRS01", lErr)
// 		return lReqRec, lErr
// 	}

// 	// call the getPanNO method to get the lPanNo
// 	lPanNo, lErr := getPanNO(pClientId)
// 	if lErr != nil {
// 		log.Println("SPCSRS02", lErr)
// 		return lReqRec, lErr
// 	}

// 	// call the getApplication method to get the lAppNo
// 	lSymbol, lBidId, lOrderNo, lErr := getSGBApplication(pReqRec.MasterId, pReqRec.OrderNo, pClientId)
// 	if lErr != nil {
// 		log.Println("SPCSRS03", lErr)
// 		return lReqRec, lErr
// 	} else {
// 		// If lAppNo is nil,generate new application no or else pass the lAppNo
// 		if lOrderNo == "0" {
// 			var lTrimmedString string
// 			lTime := time.Now()
// 			lUnixTime := lTime.Unix()
// 			lUnixTimeString := fmt.Sprintf("%d", lUnixTime)
// 			if len(pClientId) >= 5 {
// 				lTrimmedString = pClientId[len(pClientId)-5:]
// 			}
// 			lBidReqRec.OrderNo = lUnixTimeString + lTrimmedString
// 			lBidReqRec.BidId = ""
// 		} else {
// 			lBidReqRec.OrderNo = pReqRec.OrderNo
// 			lBidReqRec.BidId = lBidId
// 		}
// 	}

// 	lReqRec.ScripId = lSymbol
// 	lReqRec.InvestorCategory = "CTZ"
// 	lReqRec.PanNo = lPanNo
// 	lReqRec.ApplicantName = lClientName
// 	lReqRec.Depository = "CDSL"
// 	lReqRec.DpId = "0"
// 	lReqRec.ClientBenfId = lDpId
// 	lReqRec.GuardianName = ""
// 	lReqRec.GuardianPanno = ""
// 	lReqRec.GuardianRelation = ""

// 	if pExchange == common.NSE {
// 		lBidReqRec.BidId = pReqRec.BidId
// 		lBidReqRec.OrderNo = pReqRec.BidId
// 	}
// 	//bid details
// 	lBidReqRec.SubscriptionUnit = strconv.Itoa(pReqRec.Unit)
// 	lBidReqRec.Rate = strconv.Itoa(pReqRec.Price)
// 	lBidReqRec.ActionCode = pReqRec.ActionCode

// 	//append to the array of bse structs
// 	lReqRec.Bids = append(lReqRec.Bids, lBidReqRec)

// 	log.Println("ConstructSGBReqStruct (-)")
// 	return lReqRec, nil
// }

// /*
// Pupose:This method is used to get the Application Number From the DataBase.
// Parameters:

// 	pSymbol,pAppNo,PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	FT000006912345,nil

// 	==========
// 	*On Error
// 	==========
// 	"",error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func getSGBApplication(pMasterId int, pOrderNo string, pClientId string) (string, string, string, error) {
// 	log.Println("getSGBApplication (+)")

// 	// this variable is used to get the application no and reference no from the database
// 	var lSymbol, lBidId, lOrderNo string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("PGSA01", lErr1)
// 		return lSymbol, lBidId, lOrderNo, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select nvl(m.Symbol,'') symbol,(case when count(1) > 0 then d.OrderNo else 0 end) orderno,(case when count(1) > 0 then d.BidId else 0 end) bidid
// 						from a_sgb_orderheader h,a_sgb_master m,a_sgb_orderdetails d
// 						where m.id = h.MasterId
// 						and h.Id = d.HeaderId
// 						and h.ClientId = ?
// 						and d.OrderNo = ?
// 						and m.id = ?`
// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pOrderNo, pMasterId)
// 		if lErr2 != nil {
// 			log.Println("PGSA02", lErr2)
// 			return lSymbol, lBidId, lOrderNo, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lSymbol, &lOrderNo, &lBidId)
// 				if lErr3 != nil {
// 					log.Println("PGSA03", lErr3)
// 					return lSymbol, lBidId, lOrderNo, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("getSGBApplication (-)")
// 	return lSymbol, lBidId, lOrderNo, nil
// }

// /*
// Pupose:This method is used to get the client BenId and ClientName for order input.
// Parameters:

// 	PClientId

// Response:

// 	==========
// 	*On Sucess
// 	==========
// 	123456789012,Lakshmanan Ashok Kumar,nil

// 	==========
// 	*On Error
// 	==========
// 	"","",error

// Author:Pavithra
// Date: 12JUNE2023
// */
// func getDpId(pClientId string) (string, string, error) {
// 	log.Println("getDpId (+)")

// 	// this variables is used to get DpId and ClientName from the database.
// 	var lDpId string
// 	var lClientName string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
// 	if lErr1 != nil {
// 		log.Println("PGDI01", lErr1)
// 		return lDpId, lClientName, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select  idm.CLIENT_DP_CODE, idm.CLIENT_DP_NAME
// 						from   TECHEXCELPROD.CAPSFO.DBO.IO_DP_MASTER idm
// 						where idm.CLIENT_ID = ?
// 						and DEFAULT_ACC = 'Y'
// 						and DEPOSITORY = 'CDSL' `
// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
// 		if lErr2 != nil {
// 			log.Println("PGDI02", lErr2)
// 			return lDpId, lClientName, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lDpId, &lClientName)
// 				if lErr3 != nil {
// 					log.Println("PGDI03", lErr3)
// 					return lDpId, lClientName, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("getDpId (-)")
// 	return lDpId, lClientName, nil
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
// Date: 12JUNE2023
// */
// func getPanNO(pClientId string) (string, error) {
// 	log.Println("getPanNO (+)")

// 	// this variables is used to get Pan number of the client from the database.
// 	var lPanNo string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
// 	if lErr1 != nil {
// 		log.Println("PGPN01", lErr1)
// 		return lPanNo, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select pan_no
// 						from TECHEXCELPROD.CAPSFO.DBO.client_details
// 						where client_Id = ? `
// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
// 		if lErr2 != nil {
// 			log.Println("PGPN02", lErr2)
// 			return lPanNo, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lPanNo)
// 				if lErr3 != nil {
// 					log.Println("PGPN03", lErr3)
// 					return lPanNo, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("getPanNO (-)")
// 	return lPanNo, nil
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

// Author:Prashanth
// Date: 24aug2023
// */
// func UpdateSgbBidTrack(pRespRec bsesgb.SgbRespStruct, pId int, pClientId string, pJvStatusRec JvStatusStruct, pExchange string) error {
// 	log.Println("UpdateSgbBidTrack (+)")
// 	var lStatus string
// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SPUST01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()
// 		if pRespRec.StatusCode == "0" {
// 			lStatus = "success"
// 		} else {
// 			lStatus = "failed"
// 		}
// 		// log.Println(pRespRec.Bids)
// 		for lBidIdx := 0; lBidIdx < len(pRespRec.Bids); lBidIdx++ {
// 			lSqlString := `update  a_sgbtracking_table
// 							set OrderNo= ?,Bidid = ?, ActivityType =?,Unit = ?,Price =?,ApplicationStatus = ?,ErrorCode = ?,Message = ?,
// 							JvAmount = ?,JvStatus = ?,JvStatement = ?,JvType = ?,UpdatedBy = ?,UpdatedDate = now() ,Exchange = ?
// 							where id = ? `
// 			_, lErr2 := lDb.Exec(lSqlString, pRespRec.Bids[lBidIdx].OrderNo, pRespRec.Bids[lBidIdx].BidId, pRespRec.Bids[lBidIdx].ActionCode, pRespRec.Bids[lBidIdx].SubscriptionUnit, pRespRec.Bids[lBidIdx].Rate, lStatus, pRespRec.Bids[lBidIdx].ErrorCode, pRespRec.Bids[lBidIdx].Message, pJvStatusRec.JvAmount, pJvStatusRec.JvStatus, pJvStatusRec.JvStatement, pJvStatusRec.JvType, pClientId, pExchange, pId)
// 			if lErr2 != nil {
// 				log.Println("SPUST02", lErr2)
// 				return lErr2
// 			}
// 		}
// 	}
// 	log.Println("UpdateSgbBidTrack (-)")
// 	return nil
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

// Author:Prashanth
// Date: 24aug2023
// */
// func InsertSgbBidTrack(pReqRec bsesgb.SgbReqStruct, pClientId string, pExchange string) (int, error) {
// 	log.Println("InsertSgbBidTrack (+)")
// 	var Id int
// 	Id = 0
// 	// Calling LocalDbConect method in ftdb to estabish the database connection
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SPIBT01", lErr1)
// 		return 1, lErr1
// 	} else {
// 		defer lDb.Close()
// 		for lBidIdx := 0; lBidIdx < len(pReqRec.Bids); lBidIdx++ {

// 			lSqlString := `insert  into a_sgbtracking_table (OrderNo,Bidid,ActivityType,Unit,Price,ClientId,Exchange ,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate) values (?,?,?,?,?,?,?,?,now(),?,now())`
// 			lInserted, lErr2 := lDb.Exec(lSqlString, pReqRec.Bids[lBidIdx].OrderNo, pReqRec.Bids[lBidIdx].BidId, pReqRec.Bids[lBidIdx].ActionCode, pReqRec.Bids[lBidIdx].SubscriptionUnit, pReqRec.Bids[lBidIdx].Rate, pClientId, pExchange, pClientId, pClientId)
// 			if lErr2 != nil {
// 				log.Println("SPIBT02", lErr2)
// 				return Id, lErr2
// 			} else {
// 				lTrack, _ := lInserted.LastInsertId()
// 				Id = int(lTrack)
// 				// log.Println(Id)
// 			}
// 		}
// 	}
// 	log.Println("InsertSgbBidTrack (-)")
// 	return Id, nil
// }
