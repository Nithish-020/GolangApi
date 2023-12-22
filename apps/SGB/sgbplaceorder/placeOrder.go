package sgbplaceorder

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/SGB/validatesgb"
	"fcs23pkg/apps/abhilogin"
	"fcs23pkg/apps/clientDetail"
	"fcs23pkg/apps/validation/adminaccess"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fcs23pkg/integration/bse/bsesgb"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// this struct is used to get the bid details.
type SgbReqStruct struct {
	MasterId   int    `json:"masterId"`
	BidId      string `json:"bidId"`
	Unit       int    `json:"unit"`
	OldUnit    int    `json:"oldUnit"`
	Price      int    `json:"price"`
	ActionCode string `json:"actionCode"`
	OrderNo    string `json:"orderNo"`
	Amount     int    `json:"amount"`
	PreApply   string `json:"preApply"`
}

type SgbClientDetail struct {
	OrderNo    string `json:"orderno"`
	Unit       string `json:"unit"`
	Price      string `json:"price"`
	Symbol     string `json:"symbol"`
	OrderDate  string `json:"orderdate"`
	Amount     int    `json:"amount"`
	Mail       string `json:"mail"`
	ClientName string `json:"clientname"`
	Activity   string `json:"activity"`
}

// this struct is used to get the bid details.
type SgbRespStruct struct {
	OrderStatus string `json:"orderStatus"`
	Status      string `json:"status"`
	ErrMsg      string `json:"errMsg"`
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
Date: 20SEP2023
*/
func SgbPlaceOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("SgbPlaceOrder (+)", r.Method)
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
	log.Println("lBrokerId", lBrokerId)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	// Checking the Api Method
	if r.Method == "POST" {
		// create a instance of array type OrderReqStruct.This variable is used to store the request values.
		var lReqRec SgbReqStruct
		// create a instance of array type OrderResStruct.This variable is used to store the response values.
		var lRespRec SgbRespStruct
		//
		lRespRec.Status = common.SuccessCode

		lExchange, lErr1 := adminaccess.SGBFetchDirectory(lBrokerId)
		// lExchange, lErr1 := memberdetail.BseNsePercentCalc(lBrokerId, "/sgb")
		if lErr1 != nil {
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "SGBPO01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("SGBPO01", "Directory Not Found. Please try after sometime"))
			return
		} else {
			lRespRec.Status = common.SuccessCode
			// var lErrorRec error
			// var lFlag string
			// var lExchangeResp bsesgb.SgbRespStruct

			//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------

			lClientId, lErr2 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/sgb")
			if lErr2 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "SGBPO02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("SGBPO02", "UserDetails not Found"))
				return
			} else {
				if lClientId == "" {
					lRespRec.Status = common.ErrorCode
					fmt.Fprintf(w, helpers.GetErrorString("SGBPO03", "UserDetails not Found"))
					return
				}
			}
			//added by naveen:to fetch the source (from where mobile or web)by cookie name
			source, lErr3 := abhilogin.GetSourceOfUser(r, common.ABHICookieName)
			if lErr3 != nil {
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "SGBPO04" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("SGBPO04", "Unable to get source"))
				return
			}
			log.Println("source in api", source)

			//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
			// Read the request body values in lBody variable
			lBody, lErr4 := ioutil.ReadAll(r.Body)
			if lErr4 != nil {
				log.Println("SGBPO05", lErr4.Error())
				lRespRec.Status = common.ErrorCode
				fmt.Fprintf(w, helpers.GetErrorString("SGBPO05", "Unable to get your request now!!Try Again.."))
				return
			} else {
				// Unmarshal the request body values in lReqRec variable
				lErr5 := json.Unmarshal(lBody, &lReqRec)
				if lErr5 != nil {
					log.Println("SGBPO06", lErr5.Error())
					lRespRec.Status = common.ErrorCode
					fmt.Fprintf(w, helpers.GetErrorString("SGBPO06", "Unable to get your request now. Please try after sometime"))
					return
				} else {
					log.Println("lReqRec", lReqRec)
					lTimeToApply, lErr6 := TimeToAcceptOrder(lReqRec)
					if lErr6 != nil {
						log.Println("SGBPO07", lErr6.Error())
						lRespRec.Status = common.ErrorCode
						fmt.Fprintf(w, helpers.GetErrorString("SGBPO07", "Unable to Accept your Order. Please try after sometime"))
						return
					} else {
						log.Println("lTimeToApply", lTimeToApply)
						if lTimeToApply == "Y" {
							lRespRec.Status = common.ErrorCode
							fmt.Fprintf(w, helpers.GetErrorString("SGBPO08", "Bond Closed,No More Orders are Acceptable"))
							return
						} else if lTimeToApply == "N" {

							//this method is used to construct the Req struct for exchange
							lExchangeReq, lErr7 := ConstructSGBReqStruct(lReqRec, lClientId)
							if lErr7 != nil {
								log.Println("SGBPO09", lErr7.Error())
								lRespRec.Status = common.ErrorCode
								fmt.Fprintf(w, helpers.GetErrorString("SGBPO09", "Unable to process your request now. Please try after sometime"))
								return
							} else {
								log.Println(lExchangeReq)
								// log.Println(lReqRec.MasterId, "lReqRec.MasterId")
								lTodayAvailable, lErr8 := validatesgb.CheckSgbEndDate(lReqRec.MasterId)
								if lErr8 != nil {
									log.Println("SGBPO010", lErr8.Error())
									lRespRec.Status = common.ErrorCode
									fmt.Fprintf(w, helpers.GetErrorString("SGBPO010", "Unable to process your request now. Please try after sometime"))
									return
								} else {

									// check the symbol is available or not
									if lTodayAvailable == "True" {

										//added by naveen:add one additional parameter source to insert in sgb orderheader
										//lErr8 := UpdateToLocal(lReqRec, lExchangeReq, lClientId, lExchange, lBrokerId, r)
										lErr9 := UpdateToLocal(lReqRec, lExchangeReq, lClientId, lExchange, lBrokerId, r, source)
										if lErr9 != nil {
											log.Println("SGBPO11", lErr9.Error())
											lRespRec.Status = common.ErrorCode
											fmt.Fprintf(w, helpers.GetErrorString("SGBPO11", "Unable to process your request now. Please try after sometime"))
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
									} else if lTodayAvailable == "False" {
										lRespRec.Status = common.ErrorCode
										lRespRec.ErrMsg = "Bond Closed,No More Orders are Acceptable"
										log.Println("Bond Closed,No More Orders are Acceptable")
									} else {
										lRespRec.Status = common.ErrorCode
										lRespRec.ErrMsg = "Timing closed for Soverign Gold Bond"
										log.Println("Timing closed for Soverign Gold Bond")
									}
								}
							}
						}

					}
				}
			}
		}
		lData, lErr10 := json.Marshal(lRespRec)
		if lErr10 != nil {
			log.Println("SGBPO12", lErr10)
			fmt.Fprintf(w, helpers.GetErrorString("SGBPO12", "Unable to getting response.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("SgbPlaceOrder (-)", r.Method)
	}
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
func GetOrderId(pRequest SgbReqStruct, pReqRec bsesgb.SgbReqStruct, pMasterId int, pClientId string) (int, error) {
	log.Println("GetOrderId (+)")

	var lHeaderId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPGOI01", lErr1)
		return lHeaderId, lErr1
	} else {
		defer lDb.Close()

		for lReqIdx := 0; lReqIdx < len(pReqRec.Bids); lReqIdx++ {
			// for lRespIdx := 0; lRespIdx < len(pRespRec.Bids); lRespIdx++ {
			if pReqRec.Bids[lReqIdx].OrderNo == pRequest.OrderNo {
				log.Println("OrderNo", pReqRec.Bids[lReqIdx].OrderNo, pRequest.OrderNo)

				lCoreString := `select (case when count(1) > 0 then h.Id  else 0 end) Id
									from a_sgb_orderheader h,a_sgb_orderdetails d
									where h.Id = d.HeaderId 
									and h.MasterId = ?
									and h.ClientId = ?
									and d.ReqOrderNo = ?`
				lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pClientId, pRequest.OrderNo)
				if lErr2 != nil {
					log.Println("SPGOI02", lErr2)
					return lHeaderId, lErr2
				} else {
					for lRows.Next() {
						lErr3 := lRows.Scan(&lHeaderId)
						if lErr3 != nil {
							log.Println("SPGOI03", lErr3)
							return lHeaderId, lErr3
						}
					}
				}
				//
				log.Println(lHeaderId)
			}
		}
	}
	// }
	log.Println("GetOrderId (-)")
	return lHeaderId, nil
}

/*
Purpose:This method is used to construct the exchange header values.
Parameters:

	pReqRec,pClientId

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
func ConstructSGBReqStruct(pReqRec SgbReqStruct, pClientId string) (bsesgb.SgbReqStruct, error) {
	log.Println("ConstructSGBReqStruct (+)")

	// create an instance of lReqRec type bsesgb.SgbReqStruct.
	var lReqRec bsesgb.SgbReqStruct

	// create an instance of lReqRec type bsesgb.SgbReqStruct.
	var lBidReqRec bsesgb.ReqSgbBidStruct

	// call the getDpId method to get the lDpId and lClientName
	lDpId, lClientName, lErr := getDpId(pClientId)
	if lErr != nil {
		log.Println("SPCSRS01", lErr)
		return lReqRec, lErr
	}

	// call the getPanNO method to get the lPanNo
	lPanNo, lErr := getPanNO(pClientId)
	if lErr != nil {
		log.Println("SPCSRS02", lErr)
		return lReqRec, lErr
	}

	// call the getApplication method to get the lAppNo
	lSymbol, lBidId, lOrderNo, lErr := getSGBApplication(pReqRec, pClientId)
	if lErr != nil {
		log.Println("SPCSRS03", lErr)
		return lReqRec, lErr
	} else {
		// If lAppNo is nil,generate new application no or else pass the lAppNo
		if lOrderNo == "0" {
			var lTrimmedString string
			lTime := time.Now()
			lUnixTime := lTime.Unix()
			lUnixTimeString := fmt.Sprintf("%d", lUnixTime)
			if len(pClientId) >= 5 {
				lTrimmedString = pClientId[len(pClientId)-5:]
			}
			lBidReqRec.OrderNo = lUnixTimeString + lTrimmedString
			lBidReqRec.BidId = strconv.Itoa(common.GetRandomNumber())
		} else {
			lBidReqRec.OrderNo = pReqRec.OrderNo
			lBidReqRec.BidId = lBidId
		}
	}

	lReqRec.ScripId = lSymbol
	lReqRec.InvestorCategory = "CTZ"
	lReqRec.PanNo = lPanNo
	lReqRec.ApplicantName = lClientName
	lReqRec.Depository = "CDSL"
	lReqRec.DpId = "0"
	lReqRec.ClientBenfId = lDpId
	lReqRec.GuardianName = ""
	lReqRec.GuardianPanno = ""
	lReqRec.GuardianRelation = ""
	//bid details
	lBidReqRec.SubscriptionUnit = strconv.Itoa(pReqRec.Unit)
	lBidReqRec.Rate = strconv.Itoa(pReqRec.Price)
	lBidReqRec.ActionCode = pReqRec.ActionCode

	//append to the array of bse structs
	lReqRec.Bids = append(lReqRec.Bids, lBidReqRec)

	log.Println("ConstructSGBReqStruct (-)")
	return lReqRec, nil
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
func getSGBApplication(pReqRec SgbReqStruct, pClientId string) (string, string, string, error) {
	log.Println("getSGBApplication (+)")

	// this variable is used to get the application no and reference no from the database
	var lSymbol, lBidId, lOrderNo string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("PGSA01", lErr1)
		return lSymbol, lBidId, lOrderNo, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select nvl(m.Symbol,'') symbol,(case when count(1) > 0 then d.ReqOrderNo else 0 end) orderno,(case when count(1) > 0 then d.BidId else 0 end) bidid
						from a_sgb_orderheader h,a_sgb_master m,a_sgb_orderdetails d
						where m.id = h.MasterId 
						and h.Id = d.HeaderId 
						and h.ClientId = ?
						and d.ReqOrderNo = ?		
						and m.id = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pClientId, pReqRec.OrderNo, pReqRec.MasterId)
		if lErr2 != nil {
			log.Println("PGSA02", lErr2)
			return lSymbol, lBidId, lOrderNo, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSymbol, &lOrderNo, &lBidId)
				if lErr3 != nil {
					log.Println("PGSA03", lErr3)
					return lSymbol, lBidId, lOrderNo, lErr3
				}
			}
		}
	}
	log.Println("-----------------", lSymbol, lBidId, lOrderNo)
	log.Println("getSGBApplication (-)")
	return lSymbol, lBidId, lOrderNo, nil
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
func getDpId(pClientId string) (string, string, error) {
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
		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
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
func getPanNO(pClientId string) (string, error) {
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
		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
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

Author:Prashanth
Date: 24aug2023
*/
//added by naveen:add one additional parameter pSource to insert in sgbtacking
// func InsertSgbBidTrack(pReqRec bsesgb.SgbReqStruct, pClientId string, pExchange string) error {
// commented by pavithra
// func InsertSgbBidTrack(pReqRec bsesgb.SgbReqStruct, pClientId string, pExchange string, pSource string) error {
func InsertSgbBidTrack(pReqRec bsesgb.SgbReqStruct, pClientId string, pExchange string, pSource string, pBrokerId int) error {
	log.Println("InsertSgbBidTrack (+)")
	// var Id int
	// Id = 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPIBT01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for lBidIdx := 0; lBidIdx < len(pReqRec.Bids); lBidIdx++ {

			//commented by pavithra
			// lSqlString := `insert  into a_sgbtracking_table (OrderNo,Bidid,ActivityType,Unit,Price,ClientId,Exchange,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate) values (?,?,?,?,?,?,?,?,?,now(),?,now())`
			lSqlString := `insert  into a_sgbtracking_table (OrderNo,Bidid,ActivityType,Unit,Price,ClientId,Exchange,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,brokerId) values (?,?,?,?,?,?,?,?,?,now(),?,now(),?)`

			//commented by pavithra
			// _, lErr2 := lDb.Exec(lSqlString, pReqRec.Bids[lBidIdx].OrderNo, pReqRec.Bids[lBidIdx].BidId, pReqRec.Bids[lBidIdx].ActionCode, pReqRec.Bids[lBidIdx].SubscriptionUnit, pReqRec.Bids[lBidIdx].Rate, pClientId, pExchange, pSource, pClientId, pClientId)
			_, lErr2 := lDb.Exec(lSqlString, pReqRec.Bids[lBidIdx].OrderNo, pReqRec.Bids[lBidIdx].BidId, pReqRec.Bids[lBidIdx].ActionCode, pReqRec.Bids[lBidIdx].SubscriptionUnit, pReqRec.Bids[lBidIdx].Rate, pClientId, pExchange, pSource, pClientId, pClientId, pBrokerId)
			if lErr2 != nil {
				log.Println("SPIBT02", lErr2)
				return lErr2
			} else {
				// lTrack, _ := lInserted.LastInsertId()
				// Id = int(lTrack)
				log.Println("Inserted Successfully")
			}
		}
	}
	log.Println("InsertSgbBidTrack (-)")
	return nil
}

//-----------------------------------------------------------------------------
//Req Update
//-----------------------------------------------------------------------------

// added by naveen:add one additional parameter pSource to insert in sgb orderheader
// func UpdateToLocal(pReqRec SgbReqStruct, pExchangeReq bsesgb.SgbReqStruct, pClientId string, pExchange string, pBrokerId int, r *http.Request) error {
func UpdateToLocal(pReqRec SgbReqStruct, pExchangeReq bsesgb.SgbReqStruct, pClientId string, pExchange string, pBrokerId int, r *http.Request, pSource string) error {
	log.Println("UpdateToLocal (+)")

	log.Println("source", pSource)
	//added by naveen:add one additional argument pSource to insert in sgbtacking
	//lErr1 := InsertSgbBidTrack(pExchangeReq, pClientId, pExchange)
	//commented by pavithra : added a new parameter to Insert sgb bid track table
	// lErr1 := InsertSgbBidTrack(pExchangeReq, pClientId, pExchange, pSource)
	lErr1 := InsertSgbBidTrack(pExchangeReq, pClientId, pExchange, pSource, pBrokerId)
	if lErr1 != nil {
		log.Println("SPUTL01", lErr1.Error())
		return lErr1
	} else {
		// This method is Not Providing Mail ID Because Client Data is not inserting on DB
		// lClientEmail, lErr2 := clientDetail.GetClientEmailId(r, lClientId)
		//  Temperory Alternate method to create here to Get an Client Mail ID
		lClientEmail, lErr2 := clientDetail.GetClientEmailId(pClientId)
		if lErr2 != nil {
			log.Println("SPUTL02", lErr2.Error())
			return lErr2
		} else {

			log.Println("lClientEmail", lClientEmail)

			//added by naveen:add one additonal argument pSource to insert in sgb orderheader
			//lErr3 := InsertHeader(pReqRec, pExchangeReq, pClientId, lClientEmail, pExchange, pBrokerId)
			lErr3 := InsertHeader(pReqRec, pExchangeReq, pClientId, lClientEmail, pExchange, pBrokerId, pSource)
			if lErr3 != nil {
				log.Println("SPUTL03", lErr3.Error())
				return lErr3
			} else {
				log.Println("Updated Successfully in DB")
			}

		}
	}
	log.Println("UpdateToLocal (-)")

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


Author:Pavithra
Date: 12JUNE2023
*/
//added by naveen:add one additional parameter pSource to insert in sgb orderheader
//func InsertHeader(pRequest SgbReqStruct, pReqRec bsesgb.SgbReqStruct, pClientId string, pEmailId string, pExchange string, pBrokerId int) error {
func InsertHeader(pRequest SgbReqStruct, pReqRec bsesgb.SgbReqStruct, pClientId string, pEmailId string, pExchange string, pBrokerId int, pSource string) error {
	log.Println("InsertHeader (+)")
	log.Println("source in api", pSource)
	// get the application no id in table
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

		//  call the getAppNoId method, To get Application No Id from the database
		lHeadId, lErr2 := GetOrderId(pRequest, pReqRec, pRequest.MasterId, pClientId)
		if lErr2 != nil {
			log.Println("PIH02", lErr2)
			return lErr2
		} else {
			if lHeadId != 0 {
				lHeaderId = lHeadId
			} else {
				lHeaderId = pRequest.MasterId
			}
			// 		config := common.ReadTomlConfig("./toml/SgbConfig.toml")
			// NonProcess := fmt.Sprintf("%v", config.(map[string]interface{})["NonProcess"])

			for lBidIdx := 0; lBidIdx < len(pReqRec.Bids); lBidIdx++ {

				//if appliaction no is not present in the database,insert the new application details in orderheader table
				if lHeadId == 0 {

					processFlag := "N"
					Schstatus := "N"
					lSqlString1 := `insert into a_sgb_orderheader (brokerId,MasterId ,ScripId,PanNo,InvestorCategory,ApplicantName,Depository,
									DpId,ClientBenfId,GuardianName,GuardianPanNo,GuardianRelation,
									Status,ClientId,cancelFlag,Exchange,source,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate,ClientEmail,ProcessFlag,ScheduleStatus)
									values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now(),?,?,?)`

					lInsertedHeaderId, lErr3 := lDb.Exec(lSqlString1, pBrokerId, lHeaderId, pReqRec.ScripId, pReqRec.PanNo, pReqRec.InvestorCategory, pReqRec.ApplicantName, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenfId, pReqRec.GuardianName, pReqRec.GuardianPanno, pReqRec.GuardianRelation, common.SUCCESS, pClientId, lCancelFlag, pExchange, pSource, pClientId, pClientId, pEmailId, processFlag, Schstatus)
					if lErr3 != nil {
						log.Println("PIH03", lErr3)
						return lErr3
					} else {
						// get lastinserted id in lReturnId and converet them into int ,store it in lHeaderId
						lReturnId, _ := lInsertedHeaderId.LastInsertId()
						lHeaderId = int(lReturnId)
						log.Println("lHeaderId", lHeaderId)

						// call InsertDetails method to inserting the order details in order details table
						lErr4 := InsertDetail(pRequest, pReqRec.Bids, lHeaderId, pClientId, pExchange)
						if lErr4 != nil {
							log.Println("PIH04", lErr4)
							return lErr4
						} else {
							log.Println("header inserted successfully")
						}
					}
				} else if pReqRec.Bids[lBidIdx].ActionCode == "D" {
					lCancelFlag = "Y"
					lSqlString2 := `update a_sgb_orderheader h
									set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
									h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,
									h.GuardianRelation = ?,h.Status = ?,h.ClientId = ?,h.cancelFlag = ?,
									h.UpdatedBy = ?,h.UpdatedDate = now()
									where h.Id = ?
									and h.ClientId = ?
									and h.Exchange = ?
									and h.brokerId = ?`

					_, lErr5 := lDb.Exec(lSqlString2, pReqRec.ScripId, pReqRec.PanNo, pReqRec.InvestorCategory, pReqRec.ApplicantName, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenfId, pReqRec.GuardianName, pReqRec.GuardianPanno, pReqRec.GuardianRelation, common.SUCCESS, pClientId, lCancelFlag, pClientId, lHeaderId, pClientId, pExchange, pBrokerId)
					if lErr5 != nil {
						log.Println("PIH05", lErr5)
						return lErr5
					} else {
						// call InsertDetails method to inserting the order details in order details table
						lErr6 := InsertDetail(pRequest, pReqRec.Bids, lHeaderId, pClientId, pExchange)
						if lErr6 != nil {
							log.Println("PIH06", lErr6)
							return lErr6
						} else {
							log.Println("header cancel updated successfully")
						}
					}
				} else {
					log.Println("else in header")
					lSqlString3 := `update a_sgb_orderheader h
								set h.ScripId = ?,h.PanNo = ?,h.InvestorCategory = ?,h.ApplicantName = ?,h.Depository = ?,h.DpId = ?,
								h.ClientBenfId = ?,h.GuardianName = ?,h.GuardianPanNo = ?,h.GuardianRelation = ?,
								h.Status = ?,h.ClientId = ?,h.cancelFlag = ?,
								h.UpdatedBy = ?,h.UpdatedDate = now()
								where h.Id = ?
								and h.ClientId = ? 
								and h.Exchange = ?
								and h.brokerId = ?`

					_, lErr5 := lDb.Exec(lSqlString3, pReqRec.ScripId, pReqRec.PanNo, pReqRec.InvestorCategory, pReqRec.ApplicantName, pReqRec.Depository, pReqRec.DpId, pReqRec.ClientBenfId, pReqRec.GuardianName, pReqRec.GuardianPanno, pReqRec.GuardianRelation, common.SUCCESS, pClientId, lCancelFlag, pClientId, lHeaderId, pClientId, pExchange, pBrokerId)
					if lErr5 != nil {
						log.Println("PIH07", lErr5)
						return lErr5
					} else {
						// call InsertDetails method to inserting the order details in order details table
						lErr6 := InsertDetail(pRequest, pReqRec.Bids, lHeaderId, pClientId, pExchange)
						if lErr6 != nil {
							log.Println("PIH08", lErr6)
							return lErr6
						} else {
							log.Println("header updated successfully")
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
func InsertDetail(pRequest SgbReqStruct, pReqBidArr []bsesgb.ReqSgbBidStruct, pHeaderId int, pClientId string, pExchange string) error {
	log.Println("InsertDetail (+)")

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPISD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		// Looping the pReqBidArr array to insert the Request bid details in database
		for lReqBidIdx := 0; lReqBidIdx < len(pReqBidArr); lReqBidIdx++ {

			// Check whether the requested bid is activity type is new,if it is new insert the bid details in order detail table
			if pReqBidArr[lReqBidIdx].ActionCode == "N" {

				// if pExchange == common.BSE {
				// if pRequest.OrderNo == pReqBidArr[lReqBidIdx].OrderNo {

				lSqlString1 := `insert into a_sgb_orderdetails (HeaderId,BidId,ReqOrderNo,RespOrderNo,ActionCode,ReqSubscriptionUnit,ReqRate
							,Exchange,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values (?,?,?,?,?,?,?,?,?,now(),?,now())`

				_, lErr3 := lDb.Exec(lSqlString1, pHeaderId, pReqBidArr[lReqBidIdx].BidId, pReqBidArr[lReqBidIdx].OrderNo, pReqBidArr[lReqBidIdx].OrderNo, pReqBidArr[lReqBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pExchange, pClientId, pClientId)
				if lErr3 != nil {
					log.Println("SPISD02", lErr3)
					return lErr3
				} else {
					log.Println("Details Inserted Successfully")
				}
				// }

				// Check whether the requested bid is activity type is modify or cancel,then update the bid details in order detail table
			} else if pReqBidArr[lReqBidIdx].ActionCode == "M" {

				//Check whether the requested bid and Response bid activity type same (modify or cancel)
				// if pRequest.OrderNo == pReqBidArr[lReqBidIdx].OrderNo {
				// ! ----------------
				lSqlString2 := `update a_sgb_orderdetails d
										set d.BidId = ?,d.ReqOrderNo = ?,d.ActionCode = ?,d.ReqSubscriptionUnit = ?,d.ReqRate = ?,
										d.UpdatedBy = ?,d.UpdatedDate = now()
										where d.HeaderId = ?
										and d.Exchange = ?`

				_, lErr4 := lDb.Exec(lSqlString2, pReqBidArr[lReqBidIdx].BidId, pReqBidArr[lReqBidIdx].OrderNo, pReqBidArr[lReqBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pClientId, pHeaderId, pExchange)
				if lErr4 != nil {
					log.Println("SPISD03", lErr4)
					return lErr4
				} else {
					log.Println("Details updated Successfully")
				}
				// }
			} else if pReqBidArr[lReqBidIdx].ActionCode == "D" {

				//Check whether the requested bid and Response bid activity type same (modify or cancel)
				// if pRequest.OrderNo == pReqBidArr[lReqBidIdx].OrderNo {
				// ! ----------------
				lSqlString3 := `update a_sgb_orderdetails d
										set d.BidId = ?,d.ReqOrderNo = ?,d.ActionCode = ?,d.ReqSubscriptionUnit = ?,d.ReqRate = ?,
										d.UpdatedBy = ?,d.UpdatedDate = now()
										where d.HeaderId = ?
										and d.Exchange = ?`

				_, lErr4 := lDb.Exec(lSqlString3, pReqBidArr[lReqBidIdx].BidId, pReqBidArr[lReqBidIdx].OrderNo, pReqBidArr[lReqBidIdx].ActionCode, pReqBidArr[lReqBidIdx].SubscriptionUnit, pReqBidArr[lReqBidIdx].Rate, pClientId, pHeaderId, pExchange)
				if lErr4 != nil {
					log.Println("SPISD04", lErr4)
					return lErr4
				} else {
					log.Println("Details updated Successfully")
				}
				// }
			}
		}
	}
	log.Println("InsertDetail (-)")
	return nil
}

// this method is used to Accept the order Based on Time
func TimeToAcceptOrder(ApplyReq SgbReqStruct) (string, error) {
	log.Println("TimeToAcceptOrder (+)")
	var lorderAcceptFlag string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPTAO01", lErr1)
		return lorderAcceptFlag, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `select (case when( master.BiddingEndDate  = Date(now()) and Time(now()) > '15:30:00' )  then 'Y' else 'N'end ) as LastDay
		from a_sgb_master master
		where id = ?`

		lRows1, lErr2 := lDb.Query(lCoreString2, ApplyReq.MasterId)
		if lErr2 != nil {
			log.Println("SPTAO02", lErr2)
			return lorderAcceptFlag, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lorderAcceptFlag)
				if lErr3 != nil {
					log.Println("SPTAO03", lErr3)
					return lorderAcceptFlag, lErr3
				} else {
					// common.ABHIBrokerId = lBrokerId
					log.Println("lorderAcceptFlag  := ", lorderAcceptFlag)
				}
			}

		}
	}
	log.Println("TimeToAcceptOrder (-)")
	return lorderAcceptFlag, nil
}
