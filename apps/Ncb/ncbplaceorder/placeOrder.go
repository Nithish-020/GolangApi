package ncbplaceorder

import (
	"crypto/rand"
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/Ncb/validatencb"
	"fcs23pkg/apps/abhilogin"
	"fcs23pkg/apps/clientDetail"
	"fcs23pkg/apps/validation/adminaccess"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fcs23pkg/integration/nse/nsencb"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// this struct is used to get the bid details.
type NcbReqStruct struct {
	MasterId      int     `json:"masterId"`
	Unit          int     `json:"unit"`
	OldUnit       int     `json:"oldUnit"`
	ActionType    string  `json:"actionType"`
	Price         float64 `json:"price"`
	UnitPrice     float64 `json:"unitPrice"`
	OrderNo       int     `json:"orderNo"`
	ApplicationNo string  `json:"applicationNo"`
	Symbol        string  `json:"symbol"`
	Series        string  `json:"series"`
	Amount        float64 `json:"amount"`
	SIText        string  `json:"SItext"`
	SIValue       bool    `json:"SIvalue"`
}

// this struct is used to get the bid details.
type NcbRespStruct struct {
	OrderStatus string `json:"orderStatus"`
	Status      string `json:"status"`
	ErrMsg      string `json:"errMsg"`
}

// // this struct is used to get application number from the  a_ipo_order_header table and  a_ipo_master table
// type NcbOrderStruct struct {
// 	ApplicationNumber string `json:"applicationNumber"`
// 	OrderNo           int    `json:"orderno"`
// 	ClientEmail       string `json:"clientEmail"`
// 	ClientId          string `json:"clientId"`
// 	CancelFlag        string `json:"cancelFlag"`
// 	ClientName        string `json:"clientName"`
// 	Symbol            string `json:"symbol"`
// 	Name              string `json:"name"`
// 	Status            string `json:"status"`
// 	OrderDate         string `json:"orderdate"`
// 	Unit              int    `json:"unit"`
// 	Amount            int    `json:"amount"`
// 	ActivityType      string `json:"activityType"`
// }

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

Author:Kavya Dharshani
Date: 20SEP2023
*/

func NcbPlaceOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("NcbPlaceOrder(+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, lErr = brokers.GetBrokerId(origin)
			log.Println(lErr, origin)
			break
		}
	}
	log.Println("lBrokerId", lBrokerId)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "POST" {

		// create a instance of array type OrderReqStruct.This variable is used to store the request values.
		var lReqRec NcbReqStruct
		// create a instance of array type OrderResStruct.This variable is used to store the response values.
		var lRespRec NcbRespStruct

		lRespRec.Status = common.SuccessCode

		lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
		lProcessingMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_Immediate_Flag"])

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------

		lClientId, lErr1 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "PNPO01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("PNPO01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lErr = common.CustomError("UserDetails not Found")
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "PNPO01" + lErr.Error()
				fmt.Fprintf(w, helpers.GetErrorString("PNPO01", "UserDetails not Found"))
				return
			}
		}

		//added by naveen:to fetch the source (from where mobile or web)by cookie name

		source, lErr2 := abhilogin.GetSourceOfUser(r, common.ABHICookieName)
		if lErr2 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "PNPO02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("PNPO02", "Unable to get source"))
			return
		} else {
			log.Println("source in api", source)

			//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
			// Read the request body values in lBody variable
			lBody, lErr3 := ioutil.ReadAll(r.Body)
			if lErr3 != nil {
				log.Println("PNPO03", lErr3.Error())
				lRespRec.Status = common.ErrorCode
				fmt.Fprintf(w, helpers.GetErrorString("PNPO03", "Unable to get your request now!!Try Again.."))
				return
			} else {
				// Unmarshal the request body values in lReqRec variable
				lErr4 := json.Unmarshal(lBody, &lReqRec)
				if lErr4 != nil {

					log.Println("PNPO04", lErr4.Error())
					lRespRec.Status = common.ErrorCode
					fmt.Fprintf(w, helpers.GetErrorString("PNPO04", "Unable to get your request now. Please try after sometime"))
					return
				} else {

					// lErr5 := json.Unmarshal(lBody, &lDynamicStruct)
					// if lErr5 != nil {
					// 	log.Println("PNPO05", lErr5.Error())
					// 	lRespRec.Status = common.ErrorCode
					// 	fmt.Fprintf(w, helpers.GetErrorString("PNPO05", "Unable to process your request now. Please try after sometime"))
					// 	return
					// } else {
					lAcceptOrder, lErrMsg, lMasterId, lErr6 := NcbAcceptClientToOrder(lReqRec, lClientId, lBrokerId)
					if lErr6 != nil {
						log.Println("PNPO06", lErr6.Error())
						lRespRec.Status = common.ErrorCode
						fmt.Fprintf(w, helpers.GetErrorString("PNPO06", "Unable to process your request now. Please try after sometime"))
						return
					} else {

						if lAcceptOrder == true {
							lTimeToApply, lErr7 := NcbTimeToAcceptOrder(lMasterId)
							if lErr7 != nil {
								log.Println("PNPO07", lErr7.Error())
								lRespRec.Status = common.ErrorCode
								fmt.Fprintf(w, helpers.GetErrorString("PNPO07", "Unable to process your request now. Please try after sometime"))
								return
							} else {
								log.Println("lTimeToApply", lTimeToApply)

								if lTimeToApply == "Y" {
									lRespRec.Status = common.ErrorCode
									fmt.Fprintf(w, helpers.GetErrorString("PNPO07", "Bond Closed,No More Orders are Acceptable"))
									return
								} else if lTimeToApply == "N" {

									lTodayAvailable, lErr8 := validatencb.CheckNcbEndDate(lReqRec.MasterId)
									if lErr8 != nil {
										log.Println("PNPO08", lErr8.Error())
										lRespRec.Status = common.ErrorCode
										fmt.Fprintf(w, helpers.GetErrorString("PNPO08", "Unable to process your request now. Please try after sometime"))
										return
									} else {
										if lTodayAvailable == "True" {
											if (lReqRec.SIValue == true && lProcessingMode == "L") ||
												(lReqRec.SIValue == false && lProcessingMode == "I") {

												lExchangeReq, lEmailId, lErr9 := ConstructNCBReqStruct(lReqRec, lClientId)
												if lErr9 != nil {
													log.Println("PNPO09", lErr9.Error())
													lRespRec.Status = common.ErrorCode
													fmt.Fprintf(w, helpers.GetErrorString("PNPO09", "Unable to process your request now. Please try after sometime"))
													return
												} else {

													lExchange, lErr10 := adminaccess.NcbFetchDirectory(lBrokerId)
													if lErr10 != nil {
														lRespRec.Status = common.ErrorCode
														lRespRec.ErrMsg = "PNPO10" + lErr10.Error()
														fmt.Fprintf(w, helpers.GetErrorString("PNPO10", "Directory Not Found. Please try after sometime"))
														return
													} else {
														if lProcessingMode != "I" {
															lErr11 := LocalUpdate(lReqRec, lExchangeReq, lClientId, lExchange, lBrokerId, source, lEmailId)
															if lErr11 != nil {
																log.Println("PNPO11", lErr11.Error())
																lRespRec.Status = common.ErrorCode
																fmt.Fprintf(w, helpers.GetErrorString("PNPO11", "Unable to process your request now. Please try after sometime"))
																return
															} else {
																if lReqRec.ActionType == "N" {
																	lRespRec.OrderStatus = "Order Placed Successfully"
																	lRespRec.Status = common.SuccessCode
																} else if lReqRec.ActionType == "M" {
																	lRespRec.OrderStatus = "Order Modified Successfully"
																	lRespRec.Status = common.SuccessCode
																} else if lReqRec.ActionType == "D" {
																	lRespRec.OrderStatus = "Order Deleted Successfully"
																	lRespRec.Status = common.SuccessCode
																}
															}
														} else {
															lRespRec.Status = common.ErrorCode
															lRespRec.ErrMsg = "Immediate order cannot be processed"
															fmt.Fprintf(w, helpers.GetErrorString("PNPO11", "Immediate order cannot be processed"))
															return
														}
													}
												}
											} else {
												lRespRec.Status = common.ErrorCode
												lRespRec.ErrMsg = "Your request cannot be processed. Please accept the policy before submitting one."
												log.Println("Your request cannot be processed. Please accept the policy before submitting one.")
											}

										} else if lTodayAvailable == "False" {
											lRespRec.Status = common.ErrorCode
											lRespRec.ErrMsg = "Bond Closed,No More Orders are Acceptable"
											log.Println("Bond Closed,No More Orders are Acceptable")
										} else {
											lRespRec.Status = common.ErrorCode
											lRespRec.ErrMsg = lErrMsg
											log.Println(lErrMsg)
										}

									}
								}
							}

						} else {
							lRespRec.Status = common.ErrorCode
							lRespRec.ErrMsg = "Not eligible for placing NCB order"
							log.Println("Not eligible for placing NCB order")
						}

					}
					// }

				}
			}
		}

		lData, lErr8 := json.Marshal(lRespRec)
		if lErr8 != nil {
			log.Println("PNPO8", lErr8)
			fmt.Fprintf(w, helpers.GetErrorString("PNPO8", "Unable to getting response.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}

	}

	log.Println("NcbPlaceOrder(-)", r.Method)
}

//-----------------------------------------------------------------------------
//Req Update
//-----------------------------------------------------------------------------

func LocalUpdate(pReqRec NcbReqStruct, pExchangeReq nsencb.NcbAddReqStruct, pClientId string, pExchange string, pBrokerId int, pSource string, pEmailId string) error {
	log.Println("LocalUpdate (+)")

	lErr1 := NcbInsertBidTrack(pReqRec, pExchangeReq, pClientId, pExchange, pBrokerId, pSource)
	if lErr1 != nil {
		log.Println("NCBPO01", lErr1)
		return lErr1
	} else {
		// lClientDetails, lErr2 := clientDetail.GetClientEmailId(pClientId)
		// if lErr2 != nil {
		// 	log.Println("NCBPO02", lErr2.Error())
		// 	return lErr2
		// } else {
		lErr3 := NcbInsertHeader(pReqRec, pExchangeReq, pClientId, pExchange, pEmailId, pBrokerId, pSource)
		if lErr3 != nil {

			log.Println("NCBPO03", lErr3.Error())
			return lErr3
		} else {
			log.Println("Updated Successfully in DB")
		}

		// }
	}

	log.Println("LocalUpdate (-)")

	return nil
}

/*
Purpose:This method is used to construct the exchange header values.
Parameters:

	pReqRec,pClientId

Response:

		==========
		*On Sucess
		==========

			Success:
			{
              "symbol": "TEST",
              "investmentValue": 100,
              "applicationNumber": "1200299929020",
              "price": 55,
              "physicalDematFlag": "D",
              "pan": "AFAKA2323L",
              "depository": "NSDL",
              "dpId": "33445566",
              "clientBenId": "12345678",
              "activityType": "N",
              "clientRefNumber": "MYREF0001",
            }

		==========
		*On Error
		==========
		[],error

Author: KAVYADHARSHANI
Date: 12OCT2023
*/
func ConstructNCBReqStruct(pReqRec NcbReqStruct, pClientId string) (nsencb.NcbAddReqStruct, string, error) {
	log.Println("ConstructNCBReqStruct(+)")

	// create an instance of lReqRec type bseNcb.NcbReqStruct.
	var lReqRec nsencb.NcbAddReqStruct

	//emailId of The client
	var lEmailId string

	// call the getDpId method to get the lDpId and lClientName
	lClientDetails, lErr := clientDetail.GetClientEmailId(pClientId)
	if lErr != nil {
		log.Println("NPCNRS01", lErr)
		return lReqRec, lEmailId, lErr
	} else {
		lEmailId = lClientDetails.EmailId
		lReqRec.Pan = lClientDetails.Pan_no
		lReqRec.ClientBenId = lClientDetails.Client_dp_code
	}
	log.Println("lClientName", lClientDetails)

	// call the getApplication method to get the lAppNo
	lAppNo, lRefno, lOrderNo, lErr := getNcbApplicationNo(pReqRec, pClientId)
	// log.Println("lSymbol", lSymbol)
	log.Println("pReqRec.Symbol", pReqRec.Symbol)
	if lErr != nil {
		log.Println("NPCNRS03", lErr)
		return lReqRec, lEmailId, lErr
	} else {
		log.Println("ApplicationNo", lAppNo, lOrderNo)
		// If lAppNo is nil,generate new application no or else pass the lAppNo
		if lAppNo == "" || lAppNo == "nil" {
			lTime := time.Now()
			lPresentTime := lTime.Format("150405")
			lReqRec.ApplicationNumber = pClientId + lPresentTime

			//referenceNumber
			lByte := make([]byte, 8)
			rand.Read(lByte)
			lReqRec.ClientRefNumber = fmt.Sprintf("%x", lByte)

		} else {
			lReqRec.ApplicationNumber = lAppNo
			lReqRec.ClientRefNumber = lRefno
		}

		if lOrderNo == 0 {
			var lTrimmedString string
			lTime := time.Now()
			lUnixTime := lTime.Unix()
			lUnixTimeString := fmt.Sprintf("%d", lUnixTime)
			if len(pClientId) >= 5 {
				lTrimmedString = pClientId[len(pClientId)-5:]
			}
			// lReqRec.OrderNumber = lUnixTimeString + lTrimmedString
			orderNumberStr := lUnixTimeString + lTrimmedString
			orderNumber, err := strconv.Atoi(orderNumberStr)
			if err != nil {
				log.Println("error", err)
			}
			lReqRec.OrderNumber = orderNumber

		} else {
			lReqRec.OrderNumber = lOrderNo
		}

	}

	lReqRec.Symbol = pReqRec.Symbol
	lReqRec.InvestmentValue = pReqRec.Unit
	lReqRec.Price = float64(pReqRec.Price)
	lReqRec.Depository = "CDSL"
	lReqRec.DpId = ""
	lReqRec.PhysicalDematFlag = "D"
	lReqRec.ActivityType = pReqRec.ActionType

	log.Println("ConstructNCBReqStruct(-)")
	return lReqRec, lEmailId, nil
}

/*
Pupose:This method is used to get the Application Number From the DataBase.
Parameters:

	pReqRec,pClientId

Response:

	==========
	*On Sucess
	==========
	FT000006912345,nil

	==========
	*On Error
	==========
	"",error

Author:Kavya Dharshani
Date: 12OCT2023
*/
func getNcbApplicationNo(pReqRec NcbReqStruct, pClientId string) (string, string, int, error) {
	log.Println("getNcbApplicationNo(+)")

	var lAppNo, lRefNo string
	var lOrderNo int

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NGAN01", lErr1)
		return lAppNo, lRefNo, lOrderNo, lErr1
	} else {
		defer lDb.Close()
		log.Println("trying to get application no", lAppNo, lRefNo, pReqRec.OrderNo, lErr1)
		lCoreString := `select (CASE WHEN count(1) > 0 THEN d.ReqOrderNo ELSE 0 END) orderno,
		                       (CASE WHEN count(1) > 0 THEN h.ReqapplicationNo ELSE 'nil' END) AppNo,
		                       (CASE WHEN count(1) > 0 THEN h.ClientRefNumber ELSE 'nil' END) Refno,nvl(n.Symbol,'') symbol
                        FROM a_ncb_orderheader h, a_ncb_master n, a_ncb_orderdetails d
                        WHERE n.id = h.MasterId
                        AND h.Id = d.HeaderId
                        AND h.ClientId = ?
                        AND h.Symbol  = ?
                        and d.ReqOrderNo  = ?
                        and h.ReqapplicationNo  = ?
						and n.id = ? `

		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pReqRec.Symbol, pReqRec.OrderNo, pReqRec.ApplicationNo, pReqRec.MasterId)
		log.Println("lReqvalue", pReqRec.Symbol)
		if lErr2 != nil {
			log.Println("NGAN02", lErr2)
			return lAppNo, lRefNo, lOrderNo, lErr2
		} else {
			for lRows.Next() {
				// lErr3 := lRows.Scan(&pReqRec.Symbol, &lOrderNo, &lAppNo, &lRefNo)
				lErr3 := lRows.Scan(&lOrderNo, &lAppNo, &lRefNo, &pReqRec.Symbol)
				if lErr3 != nil {
					log.Println("NGAN03", lErr3)
					return lAppNo, lRefNo, lOrderNo, lErr3
				}
			}
		}
	}
	log.Println("getNcbApplicationNo(-)")
	return lAppNo, lRefNo, lOrderNo, nil
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

Author:KAVYA DHARSHANI
Date: 14OCT2023
*/
func NcbInsertBidTrack(pReqResp NcbReqStruct, pReqRec nsencb.NcbAddReqStruct, pClientId string, pExchange string, pBrokerId int, pSource string) error {

	log.Println("NcbInsertBidTrack(+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NIBT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		log.Println("pReqRec.ActivityType", pReqRec.ActivityType)

		if pReqRec.ActivityType == "M" || pReqRec.ActivityType == "N" || pReqRec.ActivityType == "D" {

			lSqlString := `insert into a_Ncbbidtracking_table(brokerId,applicationNo,Series,orderNo,activityType,Unit,price,Amount,clientId,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,Exchange)
								values(?,?,?,?,?,?,?,?,?,?,?,now(),?,now(),?)`
			_, lErr2 := lDb.Exec(lSqlString, pBrokerId, pReqRec.ApplicationNumber, pReqResp.Series, pReqRec.OrderNumber, pReqRec.ActivityType, pReqRec.InvestmentValue, pReqResp.Price, pReqRec.Price, pClientId, pSource, pClientId, pClientId, pExchange)

			if lErr2 != nil {
				log.Println("Error inserting into database (bidtrack)")
				log.Println("NIBT02", lErr2)
				return lErr2
			}
		}

	}

	log.Println("NcbInsertBidTrack(-)")
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
			Success: {
                "symbol": "TEST",
                "orderNumber": 2019042500000003,
                "series": "GS",
                "applicationNumber": "1200299929020",
                "investmentValue": 100,
                "price": 10500,
                "totalAmountPayable": 10500,
                "physicalDematFlag": "D",
                "pan": "AFAKA2323L",
                "depository": "NSDL",
                "dpId": "33445566",
                "clientBenId": "12345678",
                "clientRefNumber": "MYREF0001",
                "orderStatus ": "ES",
                "rejectionReason ": null,
                "enteredBy ": "samir",
                "entryTime ": "25-04-2019 12:39:01",
                "verificationStatus ": "P",
                "verificationReason ": null,
                "clearingStatus ": "FP",
                "clearingReason ": "",
                "LastActionTime ": "25-04-2019 12:39:01",
                "status" : "success"
            }

		==========
		*On Error
		==========
		[],error

Author:KAVYA DHARSHANI M
Date: 14OCT2023
*/
func NcbInsertHeader(pReqResp NcbReqStruct, pReqRec nsencb.NcbAddReqStruct, pClientId string, pExchange string, pMailId string, pBrokerId int, pSource string) error {
	log.Println("InsertHeader (+)")

	//get the application no id in table
	var lHeaderId int
	//set cancel Flag as N
	lCancelFlag := "N"
	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NIH01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// Changing the date format DD-MM-YYYY into YYYY-MM-DD of timestamp value.

		//  call the getAppNoId method, To get Application No Id from the database
		lHeadId, lErr2 := getAppNoId(pReqResp, pReqRec, pReqResp.MasterId, pClientId)
		if lErr2 != nil {
			log.Println("NIH02", lErr2)
			return lErr2
		} else {
			if lHeadId != 0 {
				lHeaderId = lHeadId
			} else {
				lHeaderId = pReqResp.MasterId
			}

			if lHeadId == 0 {

				processFlag := "N"
				Schstatus := "N"

				lSqlString1 := `insert into a_ncb_orderheader (brokerId,MasterId, Symbol,Series,ReqInvestmentunit,ReqapplicationNo,RespapplicationNo,pan,PhysicalDematFlag,depository ,dpId, ClientBenId, clientRefNumber, clientId, CreatedBy,CreatedDate,cancelFlag,Exchange,status,ClientEmail,source, ProcessFlag,ScheduleStatus, UpdatedBy,UpdatedDate, SItext,SIvalue )
							values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?,?,?,?,?,?,?,now(),?,?)`
				lInsertedHeaderId, lErr3 := lDb.Exec(lSqlString1, pBrokerId, lHeaderId, pReqRec.Symbol, pReqResp.Series, pReqRec.InvestmentValue,
					pReqRec.ApplicationNumber, pReqRec.ApplicationNumber, pReqRec.Pan, pReqRec.PhysicalDematFlag, pReqRec.Depository, pReqRec.DpId,
					pReqRec.ClientBenId, pReqRec.ClientRefNumber, pClientId, pClientId, lCancelFlag, pExchange, common.SUCCESS, pMailId, pSource, processFlag, Schstatus, pClientId, pReqResp.SIText, pReqResp.SIValue)
				if lErr3 != nil {
					log.Println("NIH03", lErr3)
					return lErr3
				} else {
					lReturnId, _ := lInsertedHeaderId.LastInsertId()
					lHeaderId = int(lReturnId)

					lErr4 := InsertDetails(pReqRec, pReqResp, lHeaderId, pClientId, pExchange)
					if lErr4 != nil {
						log.Println("NIH04", lErr4)
						return lErr4
					} else {
						log.Println("header inserted successfully")
					}
				}

			} else if pReqRec.ActivityType == "M" {
				lSqlString3 := `update a_ncb_orderheader h
				                set h.Symbol = ?,h.Series=?,h.ReqInvestmentunit =?,h.pan  = ?,h.PhysicalDematFlag=?, h.depository  =?, h.dpId =?, h.clientBenId  =?, h.clientRefNumber =? ,h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag= ?,h.status = ?,h.clientId=?, h.SItext = ?,h.SIvalue = ?
				                where h.clientId = ?
				                and h.Id = ?
				                and h.Exchange = ?
				                and h.brokerId = ?
								and h.ReqapplicationNo = ?`
				_, lErr7 := lDb.Exec(lSqlString3, pReqRec.Symbol, pReqResp.Series, pReqResp.Unit, pReqRec.Pan, pReqRec.PhysicalDematFlag, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenId, pReqRec.ClientRefNumber,
					pClientId, lCancelFlag, common.SUCCESS, pClientId, pReqResp.SIText, pReqResp.SIValue, pClientId, lHeaderId, pExchange, pBrokerId, pReqRec.ApplicationNumber)
				if lErr7 != nil {
					log.Println("NIH07", lErr7)
					return lErr7
				} else {
					// call InsertDetails method to inserting the order details in order details table
					lErr8 := InsertDetails(pReqRec, pReqResp, lHeaderId, pClientId, pExchange)
					if lErr8 != nil {
						log.Println("NIH08", lErr8)
						return lErr8
					} else {
						log.Println("header updated successfully")
					}
				}
			} else if pReqRec.ActivityType == "D" {
				lCancelFlag = "Y"
				lSqlString2 := `update a_ncb_orderheader h
				                set h.Symbol = ?,h.Series=?,h.ReqInvestmentunit =?,h.pan  = ?,h.PhysicalDematFlag=?, h.depository  =?, h.dpId =?, h.clientBenId  =?, h.clientRefNumber =? ,h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag= ?,h.status = ?,h.clientId=?,h.SItext = ?,h.SIvalue = ?
				                where h.clientId = ?
				                and h.Id = ?
				                and h.Exchange = ?
				                and h.brokerId = ?
								and h.ReqapplicationNo = ?`
				_, lErr5 := lDb.Exec(lSqlString2, pReqRec.Symbol, pReqResp.Series, pReqResp.Unit, pReqRec.Pan, pReqRec.PhysicalDematFlag, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenId, pReqRec.ClientRefNumber,
					pClientId, "Y", common.SUCCESS, pClientId, pReqResp.SIText, pReqResp.SIValue, pClientId, lHeaderId, pExchange, pBrokerId,
					pReqRec.ApplicationNumber)
				if lErr5 != nil {
					log.Println("NIH05", lErr5)
					return lErr5
				} else {
					lErr6 := InsertDetails(pReqRec, pReqResp, lHeaderId, pClientId, pExchange)
					if lErr6 != nil {
						log.Println("NIH06", lErr6)
						return lErr6
					} else {
						log.Println("lCancelFlag", lCancelFlag, pReqRec.ActivityType)
						log.Println("header cancel updated successfully")
					}

				}
			}

		}

	}
	log.Println("InsertHeader (-)")
	return nil
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

Author:KAVYADHARSHANI
Date: 14OCT2023
*/

func getAppNoId(pReqResp NcbReqStruct, pReqRec nsencb.NcbAddReqStruct, pMasterId int, pClientId string) (int, error) {
	log.Println("getAppNoId(+)")

	var lHeaderId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NGAI01", lErr1)
		return lHeaderId, lErr1
	} else {
		defer lDb.Close()

		if pReqRec.ApplicationNumber == pReqResp.ApplicationNo {
			lCoreString := `select (case when count(1) > 0 then h.Id else 0 end) AppNo
			                from a_ncb_orderheader h
		                    where h.clientId = ? 
			                and h.MasterId = ?
			                and h.ReqapplicationNo =? `
			lRows, lErr2 := lDb.Query(lCoreString, pClientId, pMasterId, pReqResp.ApplicationNo)
			if lErr2 != nil {
				log.Println("NGAI02", lErr2)
				return lHeaderId, lErr2
			} else {
				for lRows.Next() {
					lErr3 := lRows.Scan(&lHeaderId)
					if lErr3 != nil {
						log.Println("NGAI03", lErr3)
						return lHeaderId, lErr3
					}
				}
			}
			log.Println(lHeaderId)
		}
	}
	log.Println("getAppNoId(-)")
	return lHeaderId, nil
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

Author:KAVYADHARSHANI
Date: 14OCT2023
*/

func InsertDetails(pReqRec nsencb.NcbAddReqStruct, pReqResp NcbReqStruct, pHeaderId int, pClientId string, pExchange string) error {
	log.Println("InsertDetails(+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PID01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		if pReqRec.ActivityType == "N" {
			lSqlString1 := `insert into a_ncb_orderdetails(headerId,activityType,clientReferenceNo,ReqOrderNo,RespOrderNo,symbol,ReqAmount,Reqprice,ReqUnit,status,exchange ,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values (?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

			_, lErr3 := lDb.Exec(lSqlString1, pHeaderId, pReqRec.ActivityType, pReqRec.ClientRefNumber, pReqRec.OrderNumber, pReqRec.OrderNumber, pReqRec.Symbol, pReqRec.Price, pReqResp.Price, pReqRec.InvestmentValue, common.SUCCESS, pExchange, pClientId, pClientId)
			if lErr3 != nil {
				log.Println("PID03", lErr3)
				return lErr3
			} else {
				log.Println("Price", pReqRec.Price, pReqResp.Price)
				log.Println("Details Inserted Successfully")

			}
		} else if pReqRec.ActivityType == "M" {
			lSqlString2 := `update a_ncb_orderdetails d
			set d.activityType = ?,d.clientReferenceNo =?,d.ReqAmount=?,d.Reqprice = ?,d.ReqUnit =?,d.UpdatedBy = ?,d.UpdatedDate = now()
			where d.headerId = ?
			and d.exchange = ?
			and d.symbol =?`
			_, lErr4 := lDb.Exec(lSqlString2, pReqRec.ActivityType, pReqRec.ClientRefNumber, pReqRec.Price, pReqResp.Price, pReqResp.Unit, pClientId, pHeaderId, pExchange, pReqRec.Symbol)
			if lErr4 != nil {
				log.Println("NID04", lErr4)
				return lErr4
			} else {
				log.Println("Details Modified Successfully")

			}
		} else if pReqRec.ActivityType == "D" {
			lSqlString3 := `update a_ncb_orderdetails d
			set d.activityType = ?,d.clientReferenceNo =?,d.ReqAmount=?,d.Reqprice = ?,d.ReqUnit =?,d.UpdatedBy = ?,d.UpdatedDate = now()
			where d.headerId = ?
			and d.exchange = ?
			and d.symbol =?`
			_, lErr5 := lDb.Exec(lSqlString3, pReqRec.ActivityType, pReqRec.ClientRefNumber, pReqRec.Price, pReqResp.Price, pReqResp.Unit, pClientId, pHeaderId, pExchange, pReqRec.Symbol)
			if lErr5 != nil {
				log.Println("NID05", lErr5)
				return lErr5
			} else {
				log.Println("Details Deleted Successfully")

			}
		}

	}

	log.Println("InsertDetails(-)")
	return nil

}

// this method is used to Accept the order Based on Time
func NcbTimeToAcceptOrder(pMasterId int) (string, error) {
	log.Println("NcbTimeToAcceptOrder (+)")
	var lorderAcceptFlag string

	lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CloseTime"])

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NTTAO01", lErr1)
		return lorderAcceptFlag, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `select (case when( BiddingEndDate  = Date(now()) and Time(now()) > '` + lCloseTime + `' )  then 'Y' else 'N'end ) as LastDay
		                 from a_ncb_master
		                 where BiddingStartDate <= curdate()
		                 and BiddingEndDate >= curdate() 
		                 and id = ?`

		lRows1, lErr2 := lDb.Query(lCoreString2, pMasterId)
		if lErr2 != nil {
			log.Println("NTTAO02", lErr2)
			return lorderAcceptFlag, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lorderAcceptFlag)
				if lErr3 != nil {
					log.Println("NTTAO03", lErr3)
					return lorderAcceptFlag, lErr3
				} else {
					// common.ABHIBrokerId = lBrokerId
					log.Println("lorderAcceptFlag  := ", lorderAcceptFlag)
				}
			}

		}
	}
	log.Println("NcbTimeToAcceptOrder (-)")
	return lorderAcceptFlag, nil
}
