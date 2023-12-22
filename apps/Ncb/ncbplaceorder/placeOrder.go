package ncbplaceorder

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/Ncb/validatencb"
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
	ActionCode    string  `json:"actionCode"`
	OrderNo       int     `json:"orderNo"`
	ApplicationNo string  `json:"applicationNo"`
	Symbol        string  `json:"symbol"`
	Amount        float64 `json:"amount"`
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

		lExchange, lErr1 := adminaccess.NcbFetchDirectory(lBrokerId)
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "PNPO01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("PNPO01", "Directory Not Found. Please try after sometime"))
			return
		} else {
			lExchange = "NSE"
			log.Println("lExchange", lExchange)
			lRespRec.Status = common.SuccessCode

			//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------

			lClientId, lErr2 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
			if lErr2 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "PNPO02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("PNPO02", "UserDetails not Found"))
				return
			} else {
				if lClientId == "" {
					lRespRec.Status = common.ErrorCode
					fmt.Fprintf(w, helpers.GetErrorString("PNPO02", "UserDetails not Found"))
					return
				}
			}
			//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
			// Read the request body values in lBody variable
			lBody, lErr3 := ioutil.ReadAll(r.Body)
			log.Println(string(lBody))
			if lErr3 != nil {
				log.Println("PNPO03", lErr3.Error())
				lRespRec.Status = common.ErrorCode
				fmt.Fprintf(w, helpers.GetErrorString("PNPO03", "Unable to get your request now!!Try Again.."))
				return
			} else {
				// Unmarshal the request body values in lReqRec variable
				lErr4 := json.Unmarshal(lBody, &lReqRec)
				log.Println("lReq", lReqRec)
				if lErr4 != nil {

					log.Println("PNPO04", lErr4.Error())
					lRespRec.Status = common.ErrorCode
					fmt.Fprintf(w, helpers.GetErrorString("PNPO04", "Unable to get your request now. Please try after sometime"))
					return
				} else {
					lExchangeReq, lErr5 := ConstructNCBReqStruct(lReqRec, lClientId)
					if lErr5 != nil {
						log.Println("PNPO05", lErr5.Error())
						lRespRec.Status = common.ErrorCode
						fmt.Fprintf(w, helpers.GetErrorString("PNPO05", "Unable to process your request now. Please try after sometime"))
						return
					} else {
						log.Println(lExchangeReq)
						log.Println(lReqRec.MasterId, "lReqRec.MasterId")
						lTodayAvailable, lErr6 := validatencb.CheckNcbEndDate(lReqRec.MasterId)
						if lErr6 != nil {
							log.Println("PNPO06", lErr6.Error())
							lRespRec.Status = common.ErrorCode
							fmt.Fprintf(w, helpers.GetErrorString("PNPO06", "Unable to process your request now. Please try after sometime"))
							return
						} else {
							if lTodayAvailable == "True" {

								lErr7 := LocalUpdate(lReqRec, lExchangeReq, lClientId, lExchange, r, lBrokerId)
								if lErr7 != nil {
									log.Println("PNPO07", lErr7.Error())
									lRespRec.Status = common.ErrorCode
									fmt.Fprintf(w, helpers.GetErrorString("PNPO07", "Unable to process your request now. Please try after sometime"))
									return
								} else {
									if lReqRec.ActionCode == "N" {
										lRespRec.OrderStatus = "Order Placed Successfully"
										lRespRec.Status = common.SuccessCode
									} else if lReqRec.ActionCode == "M" {
										lRespRec.OrderStatus = "Order Modified Successfully"
										lRespRec.Status = common.SuccessCode
									} else if lReqRec.ActionCode == "D" {
										lRespRec.OrderStatus = "Order Deleted Successfully"
										lRespRec.Status = common.SuccessCode
									}
								}

							} else {
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "Timing Closed for NCB"
								log.Println("Timing Closed for NCB")
							}
						}
					}
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
		log.Println("PlaceOrder (-)", r.Method)
	}

	log.Println("NcbPlaceOrder(-)", r.Method)
}

//-----------------------------------------------------------------------------
//Req Update
//-----------------------------------------------------------------------------

func LocalUpdate(pReqRec NcbReqStruct, pExchangeReq nsencb.NcbAddReqStruct, pClientId string, pExchange string, r *http.Request, pBrokerId int) error {
	log.Println("LocalUpdate (+)")

	lErr1 := NcbInsertBidTrack(pExchangeReq, pClientId, pExchange)
	if lErr1 != nil {
		log.Println("NCBPO01", lErr1)
		return lErr1
	} else {
		lClientMailId, lErr2 := clientDetail.GetClientEmailId(pClientId)
		if lErr2 != nil {
			log.Println("NCBPO02", lErr2.Error())
			return lErr2
		} else {
			lErr3 := NcbInsertHeader(pReqRec, pExchangeReq, pClientId, pExchange, lClientMailId, pBrokerId)
			if lErr3 != nil {

				log.Println("NCBPO03", lErr3.Error())
				return lErr3
			} else {
				log.Println("Updated Successfully in DB")
			}

		}
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
func ConstructNCBReqStruct(pReqRec NcbReqStruct, pClientId string) (nsencb.NcbAddReqStruct, error) {
	log.Println("ConstructNCBReqStruct(+)")

	// create an instance of lReqRec type bsesgb.SgbReqStruct.
	var lReqRec nsencb.NcbAddReqStruct
	// var orderNo int

	// call the getDpId method to get the lDpId and lClientName
	lDpId, lClientName, lErr := getDpId(pClientId)
	if lErr != nil {
		log.Println("NPCNRS01", lErr)
		return lReqRec, lErr
	}
	log.Println("lClientName", lClientName)
	// call the getPanNO method to get the lPanNo
	lPan, lErr := getPan(pClientId)
	if lErr != nil {
		log.Println("NPCNRS02", lErr)
		return lReqRec, lErr
	}

	// call the getApplication method to get the lAppNo
	lAppNo, lRefno, lOrderNo, lErr := getNcbApplicationNo(pReqRec, pClientId)
	// log.Println("lSymbol", lSymbol)
	log.Println("pReqRec.Symbol", pReqRec.Symbol)
	if lErr != nil {
		log.Println("NPCNRS03", lErr)
		return lReqRec, lErr
	} else {
		log.Println("ApplicationNo", lAppNo, lOrderNo)
		// If lAppNo is nil,generate new application no or else pass the lAppNo
		if lAppNo == "" || lAppNo == "nil" {
			lTime := time.Now()
			lPresentTime := lTime.Format("150405")
			lReqRec.ApplicationNumber = pClientId + lPresentTime

			log.Println("lReqRec.ApplicationNumber", lReqRec.ApplicationNumber)

			if lOrderNo == 0 {
				var lTrimmedString string
				// lTime := time.Now()
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
				log.Println("lReqRec.OrderNumber1", lReqRec.OrderNumber)

			} else {

				lReqRec.OrderNumber = lOrderNo
				log.Println("lReqRec.OrderNumber1", lReqRec.OrderNumber)
			}
		} else {
			lReqRec.ApplicationNumber = lAppNo
		}

		log.Println("lReqRec.ApplicationNumber1111", lReqRec.ApplicationNumber)

	}

	lReqRec.Symbol = pReqRec.Symbol
	log.Println("pReqRec.Symbol", pReqRec.Symbol)
	lReqRec.InvestmentValue = pReqRec.Unit
	// lReqRec.ApplicationNumber = lAppno
	lReqRec.Price = pReqRec.Amount
	lReqRec.Pan = lPan
	lReqRec.Depository = "CDSL"
	lReqRec.DpId = ""
	lReqRec.PhysicalDematFlag = "D"
	lReqRec.ClientBenId = lDpId
	lReqRec.ClientRefNumber = lRefno
	lReqRec.ActivityType = pReqRec.ActionCode
	log.Println("pReqRec.ActionCode", pReqRec.ActionCode)
	log.Println("lReqRec.ActivityType", lReqRec.ActivityType)

	log.Println("ConstructNCBReqStruct(-)")
	return lReqRec, nil
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

Author:KAVYADHARSHANI M
Date: 12OCT2023
*/

func getDpId(lClientId string) (string, string, error) {
	log.Println("getDpId (+)")

	// this variables is used to get DpId and ClientName from the database.
	var lDpId string
	var lClientName string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
	if lErr1 != nil {
		log.Println("NGDI01", lErr1)
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
			log.Println("NGDI02", lErr2)
			return lDpId, lClientName, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lDpId, &lClientName)
				if lErr3 != nil {
					log.Println("NGDI03", lErr3)
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

Author:KAVYA DHARSHANI
Date: 12OCT2023
*/
func getPan(pClientId string) (string, error) {
	log.Println("getPan(+)")

	// this variables is used to get Pan number of the client from the database.
	var lPanNo string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.ClientDB)
	if lErr1 != nil {
		log.Println("NGP01", lErr1)
		return lPanNo, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select pan_no
						from TECHEXCELPROD.CAPSFO.DBO.client_details
						where client_Id = ? `
		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
		if lErr2 != nil {
			log.Println("NGP02", lErr2)
			return lPanNo, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lPanNo)
				if lErr3 != nil {
					log.Println("NGP03", lErr3)
					return lPanNo, lErr3
				}
			}
		}
	}
	log.Println("getPan(-)")
	return lPanNo, nil
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
	// var lReqvalue NcbOrderStruct
	// var lReqsym nsencb.NcbAddReqStruct
	var lOrderNo int

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NGAN01", lErr1)
		return lAppNo, lRefNo, lOrderNo, lErr1
	} else {
		defer lDb.Close()
		log.Println("trying to get application no", lAppNo, lRefNo, pReqRec.OrderNo, lErr1)
		lCoreString := `SELECT
                               (CASE WHEN count(1) > 0 THEN d.OrderNo ELSE 0 END) orderno,
                               (CASE WHEN count(1) > 0 THEN h.applicationNo ELSE 'nil' END) AppNo,
                               nvl(h.ClientRefNumber, 'nil') refno,nvl(n.Symbol,'') symbol
                         FROM a_ncb_orderheader h, a_ncb_master n, a_ncb_orderdetails d
                         WHERE n.id = h.MasterId
                         AND h.Id = d.HeaderId
                         AND h.ClientId = ?
						 AND h.Symbol  = ?
						 and d.OrderNo = ?
						 and h.applicationNo  =?`

		//  AND h.ClientRefNumber = ?lRefNo, lOrderNo
		//  AND d.OrderNo = ?
		log.Println("where adat-------------------", pClientId, pReqRec.Symbol, pReqRec.OrderNo)
		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pReqRec.Symbol, pReqRec.OrderNo, pReqRec.ApplicationNo)
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
	log.Println("appno------------------------", lAppNo)
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
func NcbInsertBidTrack(pReqRec nsencb.NcbAddReqStruct, pClientId string, pExchange string) error {

	log.Println("NcbInsertBidTrack(+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NIBT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lSqlString := `insert into a_Ncbbidtracking_table(brokerId,applicationNo,orderNo,activityType,price,clientId,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,Exchange)
							values(?,?,?,?,?,?,?,now(),?,now(),?)`
		_, lErr2 := lDb.Exec(lSqlString, 4, pReqRec.ApplicationNumber, pReqRec.OrderNumber, pReqRec.ActivityType,
			pReqRec.Price, pClientId, pClientId, pClientId, pExchange)

		if lErr2 != nil {
			log.Println("Error inserting into database (bidtrack)")
			log.Println("NIBT02", lErr2)
			return lErr2
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
func NcbInsertHeader(pReqResp NcbReqStruct, pReqRec nsencb.NcbAddReqStruct, pClientId string, pExchange string, pMailId string, pBrokerId int) error {
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

				lSqlString1 := `insert into a_ncb_orderheader (brokerId,MasterId, Symbol,Investmentunit,applicationNo,pan, PhysicalDematFlag,depository ,dpId, ClientBenId, clientRefNumber, clientId, CreatedBy,CreatedDate,cancelFlag,Exchange,status,ClientEmail,
					UpdatedBy,UpdatedDate )
							values (?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?,?,?,?,now())`
				lInsertedHeaderId, lErr3 := lDb.Exec(lSqlString1, pBrokerId, lHeaderId, pReqRec.Symbol, pReqRec.InvestmentValue,
					pReqRec.ApplicationNumber, pReqRec.Pan, pReqRec.PhysicalDematFlag, pReqRec.Depository, pReqRec.DpId,
					pReqRec.ClientBenId, pReqRec.ClientRefNumber, pClientId, pClientId, lCancelFlag, pExchange, common.SUCCESS, pMailId, pClientId)
				if lErr3 != nil {
					log.Println("NIH03", lErr3)
					return lErr3
				} else {
					lReturnId, _ := lInsertedHeaderId.LastInsertId()
					lHeaderId = int(lReturnId)
					log.Println("lHeaderId", lHeaderId)

					lErr4 := InsertDetails(pReqRec, pReqResp, lHeaderId, pClientId, pExchange)
					if lErr4 != nil {
						log.Println("NIH04", lErr4)
						return lErr4
					} else {
						log.Println("header inserted successfully")
					}
				}

			} else if pReqRec.ActivityType == "M" {
				log.Println("else in header1", pReqRec.ActivityType, pReqResp.OldUnit, pReqResp.Unit)
				lSqlString3 := `update a_ncb_orderheader h
				set h.Symbol = ?,h.Investmentunit=?,h.pan  = ?,h.PhysicalDematFlag=?, h.depository  =?, h.dpId =? , h.clientBenId  =?, h.clientRefNumber =? ,h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag= ?,h.status = ?,h.clientId=?
			  where h.applicationNo = ?
				  and h.clientId = ?
				  and h.Id = ?
				  and h.Exchange = ?
				  and h.brokerId = ?`
				_, lErr7 := lDb.Exec(lSqlString3, pReqRec.Symbol, pReqResp.OldUnit, pReqRec.Pan, pReqRec.PhysicalDematFlag, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenId, pReqRec.ClientRefNumber,
					pClientId, lCancelFlag, common.SUCCESS, pClientId, pReqRec.ApplicationNumber, pClientId, lHeaderId, pExchange, pBrokerId)
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
			} else {
				log.Println("else in header2", pReqRec.ActivityType, pReqResp.OldUnit, pReqResp.Unit)
				lCancelFlag = "Y"
				lSqlString2 := `update a_ncb_orderheader h
		                          set h.Symbol = ?,h.Investmentunit=?,h.pan  = ?,h.PhysicalDematFlag=?, h.depository  =?, h.dpId =? , h.clientBenId  =?, h.clientRefNumber =? ,h.UpdatedBy = ?,h.UpdatedDate = now(),h.cancelFlag= ?,h.status = ?
                                where h.applicationNo = ?
									and h.clientId = ?
									and h.Id = ?
									and h.Exchange = ?
									and h.brokerId = ?`
				_, lErr5 := lDb.Exec(lSqlString2, pReqRec.Symbol, pReqResp.OldUnit, pReqRec.Pan, pReqRec.PhysicalDematFlag, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenId, pReqRec.ClientRefNumber,
					pClientId, "Y", common.SUCCESS, pReqRec.ApplicationNumber, pClientId, lHeaderId, pExchange, pBrokerId)
				if lErr5 != nil {
					log.Println("NIH05", lErr5)
					return lErr5
				} else {
					log.Println("lCancelFlag", lCancelFlag)
					lErr6 := InsertDetails(pReqRec, pReqResp, lHeaderId, pClientId, pExchange)
					if lErr6 != nil {
						log.Println("NIH06", lErr6)
						return lErr6
					} else {
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

		if pReqRec.ApplicationNumber == pReqRec.ApplicationNumber {
			lCoreString := `select (case when count(1) > 0 then h.Id else 0 end) AppNo
		                from a_ncb_orderheader h
		               where h.applicationNo = ?
		               and h.clientId = ?	 
		                and h.MasterId = ?`
			lRows, lErr2 := lDb.Query(lCoreString, pReqResp.ApplicationNo, pClientId, pMasterId)
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

		log.Println("pReqRec", pReqRec)

		if pReqRec.ActivityType == "N" {
			lSqlString1 := `insert into a_ncb_orderdetails(headerId,activityType,clientReferenceNo,OrderNo,symbol,price,Unit,status,exchange ,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values (?,?,?,?,?,?,?,?,?,?,now(),?,now())`

			_, lErr3 := lDb.Exec(lSqlString1, pHeaderId, pReqRec.ActivityType, pReqRec.ClientRefNumber, pReqRec.OrderNumber, pReqRec.Symbol, pReqRec.Price, pReqRec.InvestmentValue, common.SUCCESS, pExchange, pClientId, pClientId)
			if lErr3 != nil {
				log.Println("PID03", lErr3)
				return lErr3
			} else {
				log.Println("Details Inserted Successfully")

			}
		} else if pReqRec.ActivityType == "M" {
			lSqlString2 := `update a_ncb_orderdetails d
			set d.activityType = ?,d.clientReferenceNo =?,d.price = ?,d.Unit=?,d.UpdatedBy = ?,d.UpdatedDate = now()
			where d.headerId = ?
			and d.exchange = ?
			and d.symbol =?`
			_, lErr4 := lDb.Exec(lSqlString2, pReqRec.ActivityType, pReqRec.ClientRefNumber, pReqRec.Price, pReqResp.OldUnit, pClientId, pHeaderId, pExchange, pReqRec.Symbol)
			if lErr4 != nil {
				log.Println("NID04", lErr4)
				return lErr4
			} else {
				log.Println("Details Inserted Successfully")

			}
		} else if pReqRec.ActivityType == "D" {
			lSqlString3 := `update a_ncb_orderdetails d
			set d.activityType = ?,d.clientReferenceNo =?,d.price = ?,d.Unit=?,d.UpdatedBy = ?,d.UpdatedDate = now()
			where d.headerId = ?
			and d.exchange = ?
			and d.symbol =?`
			_, lErr5 := lDb.Exec(lSqlString3, pReqRec.ActivityType, pReqRec.ClientRefNumber, pReqRec.Price, pReqResp.OldUnit, pClientId, pHeaderId, pExchange, pReqRec.Symbol)
			if lErr5 != nil {
				log.Println("NID05", lErr5)
				return lErr5
			} else {
				log.Println("Details Inserted Successfully")

			}
		}

	}

	log.Println("InsertDetails(-)")
	return nil

}
