package ncblocaldetails

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// type ActiveNcbStruct struct {
// 	Id              int     `json:"id"`
// 	Symbol          string  `json:"symbol"`
// 	Series          string  `json:"series"`
// 	Name            string  `json:"name"`
// 	MinBidQuantity  int     `json:"minBidQuantity"`
// 	MaxQuantity     string  `json:"maxQuantity"`
// 	TotalQuantity   string  `json:"totalQuantity"`
// 	Isin            string  `json:"isin"`
// 	MinPrice        float64 `json:"minPrice"`
// 	MaxPrice        float64 `json:"maxPrice"`
// 	CloseDate       string  `json:"closeDate"`
// 	DateTime        string  `json:"dateTime"`
// 	CutOffFlag      string  `json:"cutOffFlag"`
// 	ActivityType    string  `json:"activityType"`
// 	Flag            string  `json:"flag"`
// 	Unit            int     `json:"unit"`
// 	OrderNo         int     `json:"orderNo"`
// 	Status          string  `json:"status"`
// 	LotSize         float32 `json:"lotSize"`
// 	ModifiedLotSize int     `json:"modifiedLotSize"`
// 	Lotvalue        int     `json:"lotValue"`
// 	FaceValue       float64 `json:"faceValue"`
// 	CutoffPrice     float64 `json:"cutoffPrice"`
// 	Amount          float64 `json:"amount"`
// }

type ActiveNcbStruct struct {
	Id                int    `json:"id"`
	Symbol            string `json:"symbol"`
	Series            string `json:"series"`
	Name              string `json:"name"`
	MinBidQuantity    int    `json:"minBidQuantity"`
	MaxQuantity       int    `json:"maxQuantity"`
	Multiples         int    `json:"multiples"`
	TotalQuantity     string `json:"totalQuantity"`
	Isin              string `json:"isin"`
	UnitPrice         int    `json:"unitPrice"`
	Amount            int    `json:"amount"`
	DateRange         string `json:"dateRange"`
	StartDateWithTime string `json:"startDateWithTime"`
	EndDateWithTime   string `json:"endDateWithTime"`
	DiscountAmt       int    `json:"discountAmt"`
	DiscountText      string `json:"discountText"`
	ActionFlag        string `json:"actionFlag"`
	ButtonText        string `json:"buttonText"`
	AppliedUnit       int    `json:"appliedUnit"`
	DisableActionBtn  bool   `json:"diableActionBtn"`
	CancelAllowed     bool   `json:"cancelAllowed"`
	ModifyAllowed     bool   `json:"modifyAllowed"`
	OrderNo           int    `json:"orderNo"`
	ApplicationNo     string `json:"applicationNo"`
	ShowSI            bool   `json:"showSI"`
	SIText            string `json:"SItext"`
	SIRefundText      string `json:"SIrefundText"`
	SIValue           bool   `json:"SIvalue"`
	InfoText          string `json:"infoText"`
}

type NcbSeriesStruct struct {
	GSecDetails      []ActiveNcbStruct `json:"gSecDetail"`
	TBillDetails     []ActiveNcbStruct `json:"tBillDetail"`
	SdlDetails       []ActiveNcbStruct `json:"sdlDetail"`
	Disclaimer       string            `json:"disclaimer"`
	GoiMasterFound   string            `json:"goimasterFound"`
	GoiNoDataText    string            `json:"goinoDataText"`
	TbillMasterFound string            `json:"tbillmasterFound"`
	TbillNoDataText  string            `json:"tbillnoDataText"`
	SdlMasterFound   string            `json:"sdlmasterFound"`
	SdlNoDataText    string            `json:"sdlnoDataText"`
	InvestCount      int               `json:"investCount"`
	Status           string            `json:"status"`
	ErrMsg           string            `json:"errMsg"`
}

type DisclaimerStruct struct {
	Disclaimer   string
	SItext       string
	SIrefundText string
}

/*
Pupose:This Function is used to Get the Active NcbDetailStruct in our database table ....
Parameters:

not Applicable

Response:

*On Sucess
=========

	{
		"gSecDetail": [
			{
				"id": 18,
				"symbol": "GJ20392502",
				"series": "GS",
				"minBidQuantity": 1,
			    "maxQuantity": 100,
				"isin": "IN0020150085",
				"unitPrice": 100,
				"dateRange": "14 Sep 2013 -  19 Dec 2024",
			    "dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
			    "discountAmt": 1,
			    "discountText": "Discount ₹",
			    "actionFlag": "N",
			    "buttonText": "New",
			    "appliedUnit": 2,
			    "diableActionBtn": false,
			    "cancelAllowed": true,
			    "modifyAllowed": true,
			    "orderNo": "170375639931286",
			    "SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
			    "SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
			    "SIvalue": true
			},

		],
		"disclaimer": "Your order will be placed on the exchange (NSE) at the end of the subscription period. Ensure to keep sufficient balance in your Trading account on the last day of the issue. Credit from \r\nstocks sold on the closing day of the issue will not be considered towards the purchase of the NCB. ",
		"masterFound": "",
		"noDataText": "",
		"investCount": 1,
		"status": "S",
		"errMsg": ""
	}

!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYA DHARSHANI
Date: 10OCT2023
*/

func GetNcbMaster(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNcbMaster(+)", r.Method)
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
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// create the instance for IpoStruct
		// var lRespRec NcbStruct
		var lRespRec NcbSeriesStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ncb")
		if lErr1 != nil {
			log.Println("NGNM01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGNM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGNM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("NGNM02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// lRespArr, lErr2 := GetNcbdetail(lClientId)
		lRespArr, lErr2 := GetNcbdetail(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("NGNM03", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "NGNM03" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("NGNM03", "Error Occur in getting Datas.."))
			return
		} else {
			log.Println("lRespArr", lRespArr)

			lRespRec = lRespArr
		}

		// Marshal the Response Structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("NGNM04", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("NGNM04", "Error Occur in getting Datas.."))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetNcbMaster (-)", r.Method)
	}
}

/*
Pupose:This Function is used to construct the final NCB master response
Parameters:

ClientId ,BrokerId

Response:

*On Sucess
=========
   {
		"gSecDetail": [
			{
				"id": 18,
				"symbol": "GJ20392502",
				"series": "GS",
				"minBidQuantity": 1,
			    "maxQuantity": 100,
				"isin": "IN0020150085",
				"unitPrice": 100,
				"dateRange": "14 Sep 2013 -  19 Dec 2024",
			    "dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
			    "discountAmt": 1,
			    "discountText": "Discount ₹",
			    "actionFlag": "N",
			    "buttonText": "New",
			    "appliedUnit": 2,
			    "diableActionBtn": false,
			    "cancelAllowed": true,
			    "modifyAllowed": true,
			    "orderNo": "170375639931286",
			    "SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
			    "SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
			    "SIvalue": true
			},

		],
		"disclaimer": "Your order will be placed on the exchange (NSE) at the end of the subscription period. Ensure to keep sufficient balance in your Trading account on the last day of the issue. Credit from \r\nstocks sold on the closing day of the issue will not be considered towards the purchase of the NCB. ",
		"masterFound": "",
		"noDataText": "",
		"investCount": 1,
		"status": "S",
		"errMsg": ""
	}


!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYA DHARSHANI M
Date: 02JAN2024
*/

func GetNcbdetail(pClientId string, pBrokerId int) (NcbSeriesStruct, error) {
	log.Println("GetNcbdetail(+)")

	// var lNcbmaster ActiveNcbStruct
	var lNcbRespRec NcbSeriesStruct
	lNcbRespRec.Status = common.SuccessCode

	lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
	lGSMasterNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_GSMasterNoDataText"])
	lTbMasterNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_TBMasterNoDataText"])
	lSGMasterNoDataTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_SGMasterNoDataText"])

	lGsNcbMasterArr, lTBNcbMasterArr, lSgNcbMasterArr, lErr1 := NcbMasterRec("getNcbMaster", pBrokerId, pClientId)
	if lErr1 != nil {
		log.Println("LGND01", lErr1)
		lNcbRespRec.Status = common.ErrorCode
	} else {
		lDisclaimerResp, lErr2 := GetNcbDisclaimer(pBrokerId)
		if lErr2 != nil {
			log.Println("LGSD02", lErr2)
			lNcbRespRec.Status = common.ErrorCode
		} else {

			lNcbRespRec.GSecDetails = lGsNcbMasterArr
			lNcbRespRec.TBillDetails = lTBNcbMasterArr
			lNcbRespRec.SdlDetails = lSgNcbMasterArr
			lNcbRespRec.Disclaimer = lDisclaimerResp.Disclaimer

			log.Println("GS", lNcbRespRec.GSecDetails)
			log.Println("TB", lNcbRespRec.TBillDetails)
			log.Println("SG", lNcbRespRec.SdlDetails)

			//GSEC
			if len(lGsNcbMasterArr) > 0 {
				lNcbRespRec.GoiMasterFound = "Y"

				for _, lCount := range lGsNcbMasterArr {
					// When the client placed a order successfully the the investCount will be increased
					// based on the no.of bid they placed.
					if lCount.ActionFlag == "M" || lCount.ActionFlag == "A" || (lCount.ActionFlag == "C" &&
						lCount.AppliedUnit > 0) {
						lNcbRespRec.InvestCount++
					}
				}
			} else {
				lNcbRespRec.GoiMasterFound = "N"
				lNcbRespRec.GoiNoDataText = lGSMasterNoDataTxt
			}
			log.Println("Tbill", lTBNcbMasterArr)
			//Tbill
			if len(lTBNcbMasterArr) > 0 {
				lNcbRespRec.TbillMasterFound = "Y"

				for _, lCount := range lTBNcbMasterArr {
					// When the client placed a order successfully the the investCount will be increased
					// based on the no.of bid they placed.
					if lCount.ActionFlag == "M" || lCount.ActionFlag == "A" || (lCount.ActionFlag == "C" &&
						lCount.AppliedUnit > 0) {
						lNcbRespRec.InvestCount++
					}
				}
			} else {
				lNcbRespRec.TbillMasterFound = "N"
				lNcbRespRec.TbillNoDataText = lTbMasterNoDataTxt
			}

			//Sdl
			if len(lSgNcbMasterArr) > 0 {
				lNcbRespRec.SdlMasterFound = "Y"

				for _, lCount := range lSgNcbMasterArr {
					// When the client placed a order successfully the the investCount will be increased
					// based on the no.of bid they placed.
					if lCount.ActionFlag == "M" || lCount.ActionFlag == "A" || (lCount.ActionFlag == "C" &&
						lCount.AppliedUnit > 0) {
						lNcbRespRec.InvestCount++
					}
				}
			} else {
				lNcbRespRec.SdlMasterFound = "N"
				lNcbRespRec.SdlNoDataText = lSGMasterNoDataTxt
			}

			log.Println("empty1111", lNcbRespRec.TbillNoDataText, lNcbRespRec.TbillMasterFound)
			log.Println("others", lNcbRespRec.GoiNoDataText, lNcbRespRec.SdlNoDataText)

		}
	}

	log.Println("GetNcbdetail(-)")
	return lNcbRespRec, nil
}

/*
Pupose:This Function gives the information about the active NCB master detail along with the client applied status
Parameters:

	Method - (Need to tell where the method calling from "getNcbMaster" / "getNcbOrderHistory" ),
			  ClientId ,BrokerId
Response:

*On Sucess
=========
   {
		"gSecDetail": [
			{
				"id": 18,
				"symbol": "GJ20392502",
				"series": "GS",
				"minBidQuantity": 1,
			    "maxQuantity": 100,
				"isin": "IN0020150085",
				"unitPrice": 100,
				"dateRange": "14 Sep 2013 -  19 Dec 2024",
			    "dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
			    "discountAmt": 1,
			    "discountText": "Discount ₹",
			    "actionFlag": "N",
			    "buttonText": "New",
			    "appliedUnit": 2,
			    "diableActionBtn": false,
			    "cancelAllowed": true,
			    "modifyAllowed": true,
			    "orderNo": "170375639931286",
			    "SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
			    "SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
			    "SIvalue": true
			},

		],
	}

!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYA DHARSHANI
Date: 02JAN2024
*/

func NcbMasterRec(pMethod string, pBrokerId int, pClientId string) ([]ActiveNcbStruct, []ActiveNcbStruct, []ActiveNcbStruct, error) {

	log.Println("NcbMasterRec (+)")

	var lNcbMasterRec ActiveNcbStruct
	var lNcbGsecArr []ActiveNcbStruct
	var lNcbTbillArr []ActiveNcbStruct
	var lNcbSdlArr []ActiveNcbStruct
	var lNcbRespRec NcbSeriesStruct

	lEndDate := ""

	lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CloseTime"])
	lGsecMaxQty := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_GsecMAXQUANTITY"])
	lTbillMaxQty := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_TbillMAXQUANTITY"])
	lSdlMaxQty := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_SdlMAXQUANTITY"])
	lDiscountTxt := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_DiscountText"])
	lShowDiscount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_ShowDiscount"])
	lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CancelAllowed"])

	lProcessingMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_Immediate_Flag"])
	lDefaultRefundText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_DefaultRefundText"])

	log.Println("lShowDiscount", lShowDiscount)

	lDisclaimerResp, lErr1 := GetNcbDisclaimer(pBrokerId)
	if lErr1 != nil {
		log.Println("LNMR01", lErr1)
		lNcbRespRec.Status = common.ErrorCode
	} else {

		// Calling LocalDbConect method in ftdb to estabish the database connection
		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr2 != nil {
			log.Println("LNMR02", lErr2)
			return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr2
		} else {
			defer lDb.Close()

			// Get the current time
			lCurrentTime := time.Now()

			// Format the time as a string
			lTimeString := lCurrentTime.Format("15:04:05")
			// lTimeString := lCurrentTime.Format("25:04:05")
			// Trim leading and trailing spaces
			lCurrentTimeString := strings.TrimSpace(lTimeString)

			log.Println("lCurrentTimeString", lCurrentTimeString)
			ladditionalString := ""

			lCoreString := `select n.id, n.Symbol, n.Series, n.Name, CAST((n.MinBidQuantity  / 100) AS SIGNED) MinBidQuantity, 
			                  CAST((n.MinBidQuantity  / 100) AS SIGNED) multiples,
			                       CONCAT(n.MinBidQuantity, ' - ', CONVERT(n.MaxQuantity, SIGNED)) AS TotalQuantity, n.Isin,n.CutoffPrice unitprice,
			                      concat(DATE_FORMAT(n.BiddingStartDate, '%d %b %Y'),' -  ',DATE_FORMAT(n.BiddingEndDate, '%d %b %Y')) as dateRange,
			                      concat(DATE_FORMAT(n.BiddingStartDate, '%d %b %Y'),' ',TIME_FORMAT(n.DailyStartTime , '%h:%i%p')) as startDateWithTime,
			                      concat(DATE_FORMAT(n.BiddingEndDate, '%d %b %Y'),' ',TIME_FORMAT('` + lCloseTime + `' , '%h:%i%p')) as endDateWithTime,
			                      (case when date_sub(n.BiddingStartDate,interval 1 day ) = curdate() then 'P'
				                        when n.BiddingStartDate > curdate() then 'U'
				                        when n.BiddingEndDate = Date(now()) and '` + lCurrentTimeString + `' > '` + lCloseTime + `' then 'C'
				                        else ''
			                      end) ActionFlag,
			                     (case when  n.BiddingEndDate = Date(now()) and '` + lCurrentTimeString + `' > '` + lCloseTime + `' then 1 else 0 end ) as DisableActionBtn,
			                     (case when '` + lCancelAllowed + `' = 'Y' then 0 else 1 end) CancelledAllowed,
			                     (case when '` + lShowDiscount + `' = 'Y' then n.MaxPrice-n.MinPrice else 0 end )discount,n.BiddingEndDate
	                        from a_ncb_master n
	                        where n.Exchange = 'NSE'`

			//below condition will added when this method is called for master data that are currently active and upcoming
			if pMethod == "getNcbMaster" {
				ladditionalString = `and n.BiddingEndDate >= curdate()
				                     and not exists (
					                 select 1
					                 from a_ncb_master n1
					                 where n1.BiddingEndDate = date(now())
				                     and n1.id = n.id
					                 and n1.DailyEndTime <= time(now()) )`
			}

			lCoreFinalString := lCoreString + ladditionalString

			lRows, lErr3 := lDb.Query(lCoreFinalString)
			if lErr3 != nil {
				log.Println("LNMR03", lErr3)
				lNcbRespRec.Status = common.ErrorCode
				return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr3
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows.Next() {
					lErr4 := lRows.Scan(&lNcbMasterRec.Id, &lNcbMasterRec.Symbol, &lNcbMasterRec.Series, &lNcbMasterRec.Name, &lNcbMasterRec.MinBidQuantity, &lNcbMasterRec.Multiples, &lNcbMasterRec.TotalQuantity, &lNcbMasterRec.Isin, &lNcbMasterRec.UnitPrice, &lNcbMasterRec.DateRange, &lNcbMasterRec.StartDateWithTime, &lNcbMasterRec.EndDateWithTime, &lNcbMasterRec.ActionFlag, &lNcbMasterRec.DisableActionBtn, &lNcbMasterRec.CancelAllowed, &lNcbMasterRec.DiscountAmt, &lEndDate)
					if lErr4 != nil {
						log.Println("LNMR04", lErr4)
						lNcbRespRec.Status = common.ErrorCode
						return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr4
					} else {

						if lNcbMasterRec.DiscountAmt < 0 {
							lNcbMasterRec.DiscountAmt = 0
						}
						lNcbMasterRec.DiscountText = lDiscountTxt

						//if the processing of the order is set as immediate
						if lProcessingMode == "I" {
							//standing instruction to be displayed while cancelling the order
							lNcbMasterRec.SIRefundText = lDisclaimerResp.SIrefundText

							//if the processing of the order is set as lastday
						} else {
							//standing instruction to be displayed while pacing the order
							lNcbMasterRec.SIText = lDisclaimerResp.SItext
							// Show the default text if client try to cancel the bid
							lNcbMasterRec.SIRefundText = lDefaultRefundText
						}

						if pMethod == "getNcbMaster" {

							lNcbRec, lErr5 := NcbOrderRec(lNcbMasterRec, pClientId, pBrokerId, lEndDate)
							if lErr5 != nil {
								log.Println("LNMR05", lErr5)
								lNcbRespRec.Status = common.ErrorCode
								return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr5
							} else {
								log.Println("Seies", lNcbRec.Series)
								if lNcbRec.Series == "GS" {

									if lGsecMaxQty != "" {
										lNcbMasterRec.MaxQuantity, _ = strconv.Atoi(lGsecMaxQty)
									}

									lNcbGsecArr = append(lNcbGsecArr, lNcbRec)

									log.Println("GS master -->", lNcbGsecArr)

								} else if lNcbRec.Series == "TB" {
									if lTbillMaxQty != "" {
										lNcbMasterRec.MaxQuantity, _ = strconv.Atoi(lTbillMaxQty)
									}
									lNcbTbillArr = append(lNcbTbillArr, lNcbRec)
								} else {
									if lSdlMaxQty != "" {
										lNcbMasterRec.MaxQuantity, _ = strconv.Atoi(lSdlMaxQty)
									}
									lNcbSdlArr = append(lNcbSdlArr, lNcbRec)
								}
							}
						} else if pMethod == "getNcbOrderHistory" {
							log.Println("pMethod", pMethod)
							log.Println("lNcbMasterRec.Series", lNcbMasterRec.Series)
							if lNcbMasterRec.Series == "GS" {
								lNcbGsecArr = append(lNcbGsecArr, lNcbMasterRec)
								log.Println("history in master dfdsf", lNcbGsecArr)
							} else if lNcbMasterRec.Series == "TB" {
								lNcbTbillArr = append(lNcbTbillArr, lNcbMasterRec)
							} else {
								lNcbSdlArr = append(lNcbSdlArr, lNcbMasterRec)
							}
						}
					}
				}
			}

		}
		log.Println("Master to history", lNcbGsecArr)
	}

	log.Println("NcbMasterRec (-)")
	return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, nil
}

/*
Pupose:This Function is used to Get NCB disclamier,SIText and SIRefundText accoring to the Broker
Parameters:

	BrokerId

Response:

*On Sucess
=========

	{
		"SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
		"SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
		"disclaimer": "Your order will be placed on the exchange (NSE) at the end of the subscription period. Ensure to keep sufficient balance
					   in your Trading account on the last day of the issue. Credit from \r\nstocks sold on the closing day of the issue will not be considered towards the purchase of the NCB. ",
	}

!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYADHARSHANI M
Date: 02JAN2024
*/
func GetNcbDisclaimer(pBrokerId int) (DisclaimerStruct, error) {
	log.Println("GetNcbDisclaimer (+)")
	var lNcbDisclaimer DisclaimerStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGSD01", lErr1)
		return lNcbDisclaimer, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select nvl(bm.Ncb_SItext, '') Ncb_SItext, nvl(bm.Ncb_SIrefundtext,'') Ncb_SIrefundtext,nvl(bm.NcbDisclaimer,'') NcbDisclaimer
		                from a_ipo_brokermaster bm 
		                where bm.Id = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("LGND02", lErr2)
			return lNcbDisclaimer, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in respective structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lNcbDisclaimer.SItext, &lNcbDisclaimer.SIrefundText, &lNcbDisclaimer.Disclaimer)
				if lErr3 != nil {
					log.Println("LGND03", lErr3)
					return lNcbDisclaimer, lErr3
				}
			}
		}
	}
	log.Println("GetNcbDisclaimer (-)")
	return lNcbDisclaimer, nil
}

/*
Pupose:This Function is used to Get client order information for each active NCB
Parameters:

	ActiveNcbStruct, ClientId , BrokerId , EndDate

Response:

*On Sucess
=========
	 {
		"gSecDetail": [
			{
				"id": 18,
				"symbol": "GJ20392502",
				"series": "GS",
				"minBidQuantity": 1,
			    "maxQuantity": 100,
				"isin": "IN0020150085",
				"unitPrice": 100,
				"dateRange": "14 Sep 2013 -  19 Dec 2024",
			    "dateRangeWithTime": "14 Sep 2013 09:00AM -  19 Dec 2024 11:00PM",
			    "discountAmt": 1,
			    "discountText": "Discount ₹",
			    "actionFlag": "N",
			    "buttonText": "New",
			    "appliedUnit": 2,
			    "diableActionBtn": false,
			    "cancelAllowed": true,
			    "modifyAllowed": true,
			    "orderNo": "170375639931286",
			    "SItext": " Place order to the extent of the available balance in my account in case of insufficient balance.",
			    "SIrefundText": " Refund will be issued to your Trading account, after you confirm the order cancellation.",
			    "SIvalue": true
			},

		],
	},

!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYA DHARSHANI
Date: 02JAN2024
*/

func NcbOrderRec(pNcbMasterRec ActiveNcbStruct, pClientId string, pBrokerId int, pEndDate string) (ActiveNcbStruct, error) {

	log.Println("NcbOrderRec (+)")

	var lUnit, lProcessFlag, lScheduleFlag, lCancelFlag, lStatus, lSIText string
	var lReqNo int
	var lSIValue bool
	var lAmount float64

	lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
	lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CancelAllowed"])
	lInfoText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_InfoText"])
	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CloseTime"])
	lProcessingMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_Immediate_Flag"])
	lModifyAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_ModifyAllowed"])

	if lModifyAllowed == "Y" {
		pNcbMasterRec.ModifyAllowed = true
	} else {
		pNcbMasterRec.ModifyAllowed = false
	}

	//if order exists or not, set the defaults for the record
	pNcbMasterRec.CancelAllowed = false
	// pNcbMasterRec.ModifyAllowed = true
	pNcbMasterRec.DisableActionBtn = true
	pNcbMasterRec.ShowSI = true
	pNcbMasterRec.AppliedUnit = 0

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LNOR01", lErr1)
		return pNcbMasterRec, lErr1
	} else {
		defer lDb.Close()

		// Get the current time
		lCurrentTime := time.Now()

		// Format the time as a string
		lTimeString := lCurrentTime.Format("15:04:05")
		// lTimeString := lCurrentTime.Format("24:04:05")

		// Get the current date
		lCurrentDate := time.Now().Format("2006-01-02") // YYYY-MM-DD

		// Trim leading and trailing spaces
		lCurrentTimeString := strings.TrimSpace(lTimeString)

		lCorestring := `select (case when d.ReqUnit  = '' then 0 else d.ReqUnit  end)unit,
		                       nvl(h.ProcessFlag,'N') ProcessFlag,nvl(h.ScheduleStatus,'N') ScheduleStatus,
		                       nvl(h.CancelFlag,'N') CancelFlag,nvl(h.Status,'') Status,nvl(d.ReqOrderNo,'') orderNo,
		                       nvl(h.ReqapplicationNo,''),nvl(h.SIvalue,0) siValue,nvl(h.SItext,'') siText, d.ReqAmount
                        from a_ncb_orderdetails d, a_ncb_orderheader h 
                        where h.id = d.HeaderId
                        and h.brokerId = ?
                        and h.clientId = ?
                        and h.MasterId = ?
                        and nvl(h.CancelFlag,'N') = 'N'
                        and h.Status = 'success'`

		lRows, lErr2 := lDb.Query(lCorestring, pBrokerId, pClientId, pNcbMasterRec.Id)
		if lErr2 != nil {
			log.Println("LNOR02", lErr2)
			return pNcbMasterRec, lErr2
		} else {

			//This for loop is used to collect the records from the database and store them in respective structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lUnit, &lProcessFlag, &lScheduleFlag, &lCancelFlag, &lStatus, &lReqNo, &pNcbMasterRec.ApplicationNo, &lSIValue, &lSIText, &lAmount)
				if lErr3 != nil {
					log.Println("LNOR03", lErr3)
					return pNcbMasterRec, lErr3
				} else {

					pNcbMasterRec.Amount = int(math.Round(lAmount))

					log.Println("Amount", pNcbMasterRec.Amount)
					//It gives the client applied unit.
					pNcbMasterRec.AppliedUnit, _ = strconv.Atoi(lUnit)

					// &lUnit, &pNcbMasterRec.ActionFlag, &pNcbMasterRec.OrderNo, &pNcbMasterRec.SIValue, &pNcbMasterRec.SIText
					pNcbMasterRec.SIText = lSIText
					pNcbMasterRec.SIValue = lSIValue
					pNcbMasterRec.OrderNo = lReqNo

					// Manipulate the flag based on the conditions
					// CASE 1: when the order is placed to exchange then lProcessFlag and lScheduleFlag will be Y in the database
					if lProcessFlag == "Y" && lScheduleFlag == "Y" {
						// Set the flag A when the Jv process is completed
						//Mark the record as Applied
						pNcbMasterRec.ActionFlag = "A"
						// CASE 2: when the order is not placed to exchange, but eligble for modify, then show modify
					} else if pNcbMasterRec.ActionFlag != "C" {
						pNcbMasterRec.ActionFlag = "M"
					}

					//check if the global configuration is set to disable the cancel option
					if lCancelAllowed == "N" {
						pNcbMasterRec.CancelAllowed = false
					} else {
						//if the processing of the order is set as immediate
						if lProcessingMode == "I" {
							pNcbMasterRec.ModifyAllowed = false

							//incase the order is placed in immediate mode, then user can cancel and reapply
							pNcbMasterRec.InfoText = lInfoText

							pNcbMasterRec.CancelAllowed = false

							pNcbMasterRec.ShowSI = false
						} else {
							//if the record is already placed to exchange
							if lProcessFlag != "N" && lScheduleFlag != "N" {
								pNcbMasterRec.CancelAllowed = false
								pNcbMasterRec.ModifyAllowed = false

								// if today is the last day of the bid and current time is less than close time
								// then allow to cancel
							} else if pEndDate == lCurrentDate && lCurrentTimeString < lCloseTime {
								pNcbMasterRec.CancelAllowed = true
								pNcbMasterRec.ModifyAllowed = true

								//if today is not the last day of the bid
							} else if pEndDate > lCurrentDate {
								pNcbMasterRec.CancelAllowed = true
								pNcbMasterRec.ModifyAllowed = true

							} else {
								pNcbMasterRec.CancelAllowed = false
								pNcbMasterRec.ModifyAllowed = false
							}
						}
					}

				}
			}

			if pNcbMasterRec.ActionFlag == "" {
				pNcbMasterRec.ActionFlag = "B"
			}
			//get the description of the action button
			lButtonText, lErr4 := NcbDiscription(pNcbMasterRec.ActionFlag)
			if lErr4 != nil {
				log.Println("LSOR03", lErr4)
				return pNcbMasterRec, lErr4
			} else {
				pNcbMasterRec.ButtonText = lButtonText
			}

			if pNcbMasterRec.ActionFlag == "P" || pNcbMasterRec.ActionFlag == "B" || pNcbMasterRec.ActionFlag == "M" {
				pNcbMasterRec.DisableActionBtn = false
			}

			if lProcessingMode == "I" {
				//if the processing of the order is set as immediate
				pNcbMasterRec.ShowSI = false
			}

		}
	}

	log.Println("NcbOrderRec (-)")
	return pNcbMasterRec, nil
}

/*
Pupose:This Function is used to Get NCB ActionFlag discription from the database
Parameters:

	"actionFlag": "M",

Response:

*On Sucess
=========

	"buttonText": "Modify",

!On Error
========

	{
		"status": E,
		"reason": "Can't able to get the data from database"
	}

Author: KAVYA DHARSHANI
Date: 02JAN2024
*/
func NcbDiscription(pActionFlag string) (string, error) {
	log.Println("NcbDiscription (+)")
	lButtonText := ""

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LND01", lErr1)
		return lButtonText, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select nvl((select ld.description
			            from xx_lookup_details ld,xx_lookup_header lh 
			            where lh.id = ld.headerId  and lh.Code = 'SgbAction' and ld.code = '` + pActionFlag + `'),'') ActionBtn
                        from xx_lookup_details ld,xx_lookup_header lh `
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("LND02", lErr2)
			return lButtonText, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in respective structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lButtonText)
				if lErr3 != nil {
					log.Println("LND03", lErr3)
					return lButtonText, lErr3
				}
			}
		}
	}
	log.Println("NcbDiscription (-)")
	return lButtonText, nil
}

// func Oldverdsion_GetNcbdetail(pClientId string, pBrokerId int) ([]ActiveNcbStruct, []ActiveNcbStruct, []ActiveNcbStruct, error) {
// 	log.Println("GetNcbdetail (+)")

// 	var lNcbMasterRec ActiveNcbStruct

// 	var lNcbGsecArr []ActiveNcbStruct
// 	var lNcbTbillArr []ActiveNcbStruct
// 	var lNcbSdlArr []ActiveNcbStruct

// 	// var lUnit string
// 	// Calling LocalDbConect method in ftdb to estabish the database connection

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("GNGND01", lErr1)
// 		return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := ` select tab.id, tab.Symbol,tab.Series,tab.Name,tab.minBidqty,CONVERT(tab.maxQty, SIGNED) AS maxQty, tab.isin, tab.minprice, tab.maxprice, tab.endDate,tab.formatted_datetime,
// 		(case when d.activityType = 'cancel' || tab.Flag = 'N' then "N" else "Y" end ) Flag,tab.Lotsize, CAST((tab.Lotsize / 100) AS SIGNED) AS ModifiedLotsize,
// 		CAST((tab.Lotsize / 100) AS SIGNED) lotvalue, tab.FaceValue, tab.CutoffPrice,
// 		  CONCAT(tab.minBidqty, ' - ', CONVERT(tab.maxQty, SIGNED)) AS minBidMaxQty,
// 		lower(tab.status),	nvl(d.OrderNo,0), nvl(d.price,0) price,  nvl(d.activityType,''), nvl(d.Unit,0)
//                    from a_ncb_orderdetails d
//                          right join (select nvl(h.Id,0) headerId, master.id,
// 				            master.Symbol,master.Series, master.Name,master.minBidqty, master.maxQty,master.isin,master.minprice,master.maxprice, master.startDate,
// 				            master.endDate,master.formatted_datetime, master.Lotsize, master.FaceValue, master.CutoffPrice,
// 		            (case when h.MasterId = master.Id and h.CancelFlag = 'N' and h.status = 'success' then 'Y' else 'N' end) Flag,
// 		            (case when h.MasterId = master.Id and h.status = "success" then h.status else "" end) status
// 	               from (select nm.id Id,nm.Symbol Symbol,nm.Name Name,nm.Series Series,
// 	               nm.MinBidQuantity minBidqty,nm.MaxQuantity maxQty,	nm.Isin isin,nm.MinPrice minprice,nm.MaxPrice maxprice, nm.BiddingStartDate startDate,
// 	               nm.Lotsize , nm.FaceValue , nm.CutoffPrice ,
// 	        concat(DATE_FORMAT(nm.BiddingStartDate, '%d %b %y'),' -  ',DATE_FORMAT(nm.BiddingEndDate, '%d %b %y'))  as endDate,
// 	         CONCAT( case WHEN DAY(nm.BiddingEndDate) % 10 = 1 AND DAY(nm.BiddingEndDate) % 100 <> 11 THEN CONCAT(DAY(nm.BiddingEndDate), 'st')
// 				          WHEN DAY(nm.BiddingEndDate) % 10 = 2 AND DAY(nm.BiddingEndDate) % 100 <> 12 THEN CONCAT(DAY(nm.BiddingEndDate), 'nd')
// 				           WHEN DAY(nm.BiddingEndDate) % 10 = 3 AND DAY(nm.BiddingEndDate) % 100 <> 13 THEN CONCAT(DAY(nm.BiddingEndDate), 'rd')
// 			          ELSE CONCAT(DAY(nm.BiddingEndDate), 'th')
// 						end,' ',
// 	          DATE_FORMAT(nm.BiddingEndDate, '%b %Y'),' | ',
// 	          TIME_FORMAT(nm.DailyEndTime , '%h:%i%p')) AS formatted_datetime
//             from a_ncb_master nm
// 	                 where nm.BiddingEndDate >= curdate()  and nm.Exchange = 'NSE' and not exists (
// 	        select 1
// 	       from a_ncb_master n
//             where n.BiddingEndDate = date(now()) and n.id = nm.id and n.DailyEndTime <= time(now())) ) master
// 	        LEFT JOIN a_ncb_orderheader h
// 				  on master.Id = h.MasterId
// 					and h.CancelFlag = 'N'
// 					and h.ClientId = ?
// 					and h.status is not null
// 				    and h.status <> 'failed'
// 				   and h.cancelFlag = 'N'
// 				    group by master.Symbol) tab
// 					on tab.headerid = d.HeaderId
// 					group by tab.id,tab.startDate ,tab.Symbol`

// 		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
// 		if lErr2 != nil {
// 			log.Println("GNGND02", lErr2)
// 			return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lNcbMasterRec.Id, &lNcbMasterRec.Symbol, &lNcbMasterRec.Series, &lNcbMasterRec.Name, &lNcbMasterRec.MinBidQuantity, &lNcbMasterRec.MaxQuantity, &lNcbMasterRec.Isin, &lNcbMasterRec.TotalQuantity, &lNcbMasterRec.OrderNo)
// 				log.Println("lNcbMasterRec", lNcbMasterRec.MaxQuantity)

// 				if lErr3 != nil {
// 					log.Println("GNGND03", lErr3)
// 					return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, lErr3
// 				} else {
// 					// lNcbMasterRec.Unit, _ = strconv.Atoi(lUnit)

// 					// Append Upi End Point in lRespRec.UpiArr array
// 					// lNcbMasterArr = append(lNcbMasterArr, lNcbMasterRec)
// 					if lNcbMasterRec.Series == "GS" {
// 						lNcbGsecArr = append(lNcbGsecArr, lNcbMasterRec)
// 						log.Println("lNcbGsecArr--->GS", lNcbGsecArr)
// 					} else if lNcbMasterRec.Series == "TB" {
// 						lNcbTbillArr = append(lNcbTbillArr, lNcbMasterRec)
// 						log.Println("lNcbTbillArr--->TB", lNcbTbillArr)
// 					} else {
// 						lNcbSdlArr = append(lNcbSdlArr, lNcbMasterRec)
// 						log.Println("lNcbSdlArr--->TB", lNcbSdlArr)
// 					}

// 				}
// 			}
// 			// log.Println(lNcbMasterArr)
// 		}

// 	}
// 	log.Println("GetNcbdetail (-)")
// 	return lNcbGsecArr, lNcbTbillArr, lNcbSdlArr, nil
// }
