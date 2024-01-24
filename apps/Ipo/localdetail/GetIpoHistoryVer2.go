package localdetail

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
)

// To store IPO History Records
type HistoryDetailStruct struct {
	Id                int               `json:"id"`
	Name              string            `json:"name"`
	Symbol            string            `json:"symbol"`
	OrderDate         string            `json:"orderDate"`
	DateRange         string            `json:"dateRange"`
	PriceRange        string            `json:"priceRange"`
	LotSize           float64           `json:"lotSize"`
	IssueSizeWithText string            `json:"issueSizeWithText"`
	Sme               bool              `json:"sme"`
	ApplicationStatus string            `json:"status"`
	DpStatus          string            `json:"dpStatus"`
	UpiStatus         string            `json:"upiStatus"`
	ErrReason         string            `json:"errReason"`
	CategoryList      IpoHistoryCatList `json:"categoryList"`
	RegistarLink      string            `json:"registarLink"`
	CancelFlag        string            `json:"cancelFlag"`
	StartDate         string            `json:"startDate"`
	EndDate           string            `json:"endDate"`
	Allotment         string            `json:"allotment"`
	Refund            string            `json:"refund"`
	Demat             string            `json:"demat"`
	Listing           string            `json:"listing"`
}
type IpoHistoryCatList struct {
	Category      string        `json:"category"`
	Code          string        `json:"code"`
	ApplicationNo string        `json:"applicationNo"`
	DiscountText  string        `json:"discountText"`
	DiscountPrice int           `json:"discountPrice"`
	DiscountType  string        `json:"discountType"`
	AppliedDetail AppliedDetail `json:"appliedDetail"`
}

// Response Structure for GetIpoMaster API
type IpoHistoryRespStruct struct {
	HistoryDetails []HistoryDetailStruct `json:"historyDetail"`
	NoDataText     string                `json:"historyNoDataText"`
	HistoryFound   string                `json:"historyFound"`
	OrderCount     int                   `json:"orderCount"`
	Status         string                `json:"status"`
	ErrMsg         string                `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Application Upi status from the NSE-Exchange and update the
changes in our database table ....
Parameters:

	{
		"symbol": "TEST",
		"applicationNumber": "1200299929020",
		"dpVerStatusFlag": "S",
		"dpVerReason": null,
		"dpVerFailCode": null
	}

Response:

	*On Sucess
	=========
	{
		"status": "success"
	}

	!On Error
	========
	{
		"status": "failed",
		"reason": "Application no does not exist"
	}

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetIpoHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetIpoHistory (+)", r.Method)
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
		// This variable is used to store the history structure
		var lHistoryRespRec IpoHistoryRespStruct

		lConfigFile := common.ReadTomlConfig("toml/IpoConfig.toml")
		lNoDataText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_HistoryNoDataText"])

		lHistoryRespRec.Status = common.SuccessCode

		/***********ADMIN PAGE
		//-----------START TO GET CLIENT  DETAILS
		// This variable is used to store the ClientId
		lClientId, err := appsso.ValidateAndGetClientDetails(r, common.ABHIAppName, common.ABHICookieName)
		if err != nil {
			log.Println("IGIH01.1", err)
			fmt.Fprintf(w, helpers.GetErrorString("IGIH01.1", "The server is busy right now. Please try after sometime"))
			return
		}
		//-----------END OF GETTING CLIENT  DETAILS
		*********************/

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		// lClientId := ""
		// var lErr1 error
		lClientId, lErr1 := apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("LGIH01", lErr1)
			lHistoryRespRec.Status = common.LoginFailure
			lHistoryRespRec.ErrMsg = "LGIH01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				lHistoryRespRec.Status = common.LoginFailure
				lHistoryRespRec.ErrMsg = "LGIH01" + ""
				fmt.Fprintf(w, helpers.GetErrorString("LGIH02", "Session Expired,Please logout and Login again !!"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// Call GetData method to get Token from Database
		lhistoryArr, lErr2 := GetData(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("LGIH02", lErr2)
			lHistoryRespRec.Status = common.ErrorCode
			lHistoryRespRec.ErrMsg = "LGIH02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH02", "Issue in Fetching Your Datas!"))
			return
		} else {

			if lhistoryArr != nil {
				lHistoryRespRec.HistoryFound = "Y"

				for _, lRecord := range lhistoryArr {
					if lRecord.CancelFlag == "N" {
						lHistoryRespRec.OrderCount++
					}

				}
				// Assign lhistoryArr to HistoryRespStruct Value
				lHistoryRespRec.HistoryDetails = lhistoryArr
			} else {
				lHistoryRespRec.HistoryFound = "N"
				lHistoryRespRec.NoDataText = lNoDataText
				lHistoryRespRec.HistoryDetails = []HistoryDetailStruct{}
			}
		}
		// Marshal the reponse structure into lDatas
		lData, lErr3 := json.Marshal(lHistoryRespRec)
		if lErr3 != nil {
			log.Println("LGIH03", lErr3)
			lHistoryRespRec.Status = common.ErrorCode
			lHistoryRespRec.ErrMsg = "LGIH03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIH03", "Issue in Getting Your Datas.Please try after sometime!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetIpoHistory (-)", r.Method)
	}
}

/*
Pupose: This method returns the collection of data from the  database
Parameters:

	not applicable

Response:

	    *On Sucess
	    ==========
	    In case of a successful execution of this method, you will get the historyStruct data
		from database

	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetData(pClientId string, pBrokerId int) ([]HistoryDetailStruct, error) {
	log.Println("GetData (+)")
	// This variable is used to store the history structure from the database
	var lhistoryRec HistoryDetailStruct
	// This variable is used to store the history structure from the database n the Array
	var lhistoryArr []HistoryDetailStruct
	// This variable is used to store the stepper structure from the database
	// var lStepperRec stepperStruct

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LHGD01", lErr1)
		return lhistoryArr, lErr1
	} else {
		defer lDb.Close()
		lBrokerId := strconv.Itoa(pBrokerId)

		// ! To get the Symbol, ApplicationNo, Price and Quantity
		lCoreString2 := `SELECT mas.Id, mas.Name, mas.Symbol,
						NVL(DATE_FORMAT(mas.createdDate , '%d-%b-%Y'),'') orderDate,
						CONCAT(DATE_FORMAT(mas.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(mas.BiddingEndDate , '%d %b')) AS dateRange,
						CONCAT(mas.MinPrice," - ", mas.MaxPrice) PriceRange,mas.lotSize,
						(CASE
							WHEN mas.IssueSize >= 10000000 THEN CONCAT(FORMAT(mas.IssueSize / 10000000, 2), ' Crores')
							WHEN mas.IssueSize >= 100000 THEN CONCAT(FORMAT(mas.IssueSize / 100000, 2), ' Lakhs')
							ELSE mas.IssueSize
						end) AS IssueSize,mas.Upi,mas.ErrReason,
						mas.Sme, mas.status, mas.DpStatus,
						mas.UPIStatus,mas.Category,mas.ApplicationNo, 
						mas.DiscountType,mas.DiscountPrice,mas.Registrar,mas.CancelFlag,
						mas.startdate, mas.enddate,
						NVL(id.allotmentFinal,''), NVL(id.refundInitiate,''), NVL(id.dematTransfer,''), NVL(id.listingDate,'')
						FROM (
							SELECT m.id, NVL(m.Name,'') Name,
							NVL(h.Symbol,'') Symbol, NVL(h.createdDate,'') createdDate,
							m.BiddingStartDate BiddingStartDate, m.BiddingEndDate BiddingEndDate,
							m.minPrice MinPrice, m.maxPrice MaxPrice,
							m.lotSize lotSize, m.IssueSize IssueSize,h.upi Upi, nvl(h.reason,'')ErrReason,
							(CASE WHEN h.cancelFlag ='N' THEN NVL(h.status,'') ELSE 'cancelled' END) status,
							NVL(
								(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.category else 1 END)
								FROM a_ipo_details ad
								WHERE ad.Isin  = m.Isin	),1
							)Sme,
							NVL(
								(SELECT xld.description 
								FROM xx_lookup_details xld,xx_lookup_header xlh
								WHERE xld.headerid = xlh.id
								AND xlh.Code = 'IpoDp'
								AND xld.Code = h.dpVerStatusFlag) , 'N/A'
								) DpStatus,
							NVL(
								(SELECT (case when h.status = 'success' then xld.description else 'N/A' end)
								FROM xx_lookup_details xld,xx_lookup_header xlh
								WHERE xld.headerid = xlh.id
								AND xlh.Code = 'IpoPay'
								AND xld.Code = h.upiPaymentStatusFlag), 'N/A'
							) UPIStatus,
							h.category Category,h.applicationNo ApplicationNo,
							s.discountType DiscountType,s.discountPrice DiscountPrice,
							NVL(
								(SELECT NVL(R.RegistrarLink, '')
								FROM a_ipo_registrars R
								WHERE R.RegistrarName = m.Registrar),
								''
							) AS Registrar,
							NVL(h.cancelFlag, 'N/A') CancelFlag	,
							NVL(m.BiddingStartDate,'') startdate,
							NVL(m.BiddingEndDate,'') enddate,
							m.Isin isin
							FROM a_ipo_orderdetails d ,a_ipo_order_header h ,a_ipo_master m,v_ipo_subcategory s
							WHERE d.HeaderId = h.Id 
							AND m.Id = h.MasterId 
							and m.id = s.id
							and h.category = s.SubCatCode
							AND h.ClientId = '` + pClientId + `'
							AND h.brokerId = ` + lBrokerId + `
							GROUP BY h.applicationNo 
						) mas LEFT JOIN a_ipo_details id 
						ON mas.isin = id.Isin 
						ORDER BY mas.createdDate desc`

		// and m.BiddingEndDate < curdate() Its goes befor client id in query
		lRows2, lErr2 := lDb.Query(lCoreString2)
		if lErr2 != nil {
			log.Println("LHGD02", lErr2)
			return lhistoryArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows2.Next() {
				lErr3 := lRows2.Scan(&lhistoryRec.Id, &lhistoryRec.Name, &lhistoryRec.Symbol, &lhistoryRec.OrderDate, &lhistoryRec.DateRange, &lhistoryRec.PriceRange,
					&lhistoryRec.LotSize, &lhistoryRec.IssueSizeWithText, &lhistoryRec.CategoryList.AppliedDetail.AppliedUPI, &lhistoryRec.ErrReason, &lhistoryRec.Sme, &lhistoryRec.ApplicationStatus, &lhistoryRec.DpStatus, &lhistoryRec.UpiStatus, &lhistoryRec.CategoryList.Code, &lhistoryRec.CategoryList.ApplicationNo, &lhistoryRec.CategoryList.DiscountType, &lhistoryRec.CategoryList.DiscountPrice, &lhistoryRec.RegistarLink, &lhistoryRec.CancelFlag, &lhistoryRec.StartDate, &lhistoryRec.EndDate, &lhistoryRec.Allotment, &lhistoryRec.Refund, &lhistoryRec.Demat, &lhistoryRec.Listing)
				if lErr3 != nil {
					log.Println("LHGD03", lErr3)
					return lhistoryArr, lErr3
				} else {

					lAppliedBids, lErr4 := getIpoAppliedBidDetails(lhistoryRec, pBrokerId)
					if lErr4 != nil {
						log.Println("LHGD04", lErr4)
						return lhistoryArr, lErr4
					} else {

						lhistoryRec.CategoryList.AppliedDetail.AppliedBids = lAppliedBids

						lResult, lErr5 := CalcAmountPayable(lhistoryRec)
						if lErr5 != nil {
							log.Println("LHGD05", lErr5)
							return lhistoryArr, lErr5
						} else {
							lhistoryRec.CategoryList.DiscountText = lResult.DiscountText
							lhistoryRec.CategoryList.AppliedDetail.AppliedAmount = lResult.AppliedDetail.AppliedAmount

							lCategory, lErr6 := getIpoDiscription(lhistoryRec.CategoryList.Code)
							if lErr6 != nil {
								log.Println("LHGD06", lErr6)
								return lhistoryArr, lErr5
							} else {
								lhistoryRec.CategoryList.Category = lCategory
							}
						}

						// Append the history Records into lhistoryArr Array
						lhistoryArr = append(lhistoryArr, lhistoryRec)
					}

				}
			}
		}
	}
	log.Println("GetData (-)")
	return lhistoryArr, nil
}

/*
Pupose:This method used to retrieve the applied Bid Details from the database based on the client.
Parameters:
	pMasterId,pApplicationNo,pBrokerId
Response:
==========
*On Sucess
==========

	[

		{
			"id": 1
			"appNo" : "1234575790"
			"upi" : "test@ybl"
			"category" : "IND"
			"bidRefNo" : "45687843213233"
			"price" : 755
			"amount" : 14567
			"quantity" : 19
			"cutOff" : true
		}
	]

==========
!On Error
==========

	[],error

Author:Nithish Kumar
Date: 19JAN2024
*/
func getIpoAppliedBidDetails(pHistoryRec HistoryDetailStruct, pBrokerId int) ([]BidDetail, error) {
	log.Println("getHistoryDetails (+)")

	// This Variable is used to store the Each applied bid Records
	var lBidRec BidDetail

	var lBidArr []BidDetail
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGHIAD01", lErr1)
		return lBidArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select d.id, d.bidReferenceNo,d.req_price,
                        d.req_amount,d.req_quantity,d.atCutOff,d.activityType
						from
						a_ipo_order_header h,a_ipo_orderdetails d, a_ipo_master m
                        where m.Id = h.MasterId 
                        and h.Id = d.headerId
                        and h.MasterId = ?
                        and h.applicationNo = ?
                        and h.brokerId = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pHistoryRec.Id, pHistoryRec.CategoryList.ApplicationNo, pBrokerId)
		if lErr2 != nil {
			log.Println("LGHIAD02", lErr2)
			return lBidArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in ModifyRespStruct
			for lRows.Next() {
				lErr3 := lRows.Scan(&lBidRec.Id, &lBidRec.BidReferenceNo, &lBidRec.Price, &lBidRec.Amount, &lBidRec.Quantity, &lBidRec.CutOff, &lBidRec.ActivityType)

				if lErr3 != nil {
					log.Println("LGHIAD03", lErr3)
					return lBidArr, lErr3
				} else {
					// Append the Bid Records in lRespArr
					log.Println("lBidRec.Quantity", lBidRec.Quantity, "pHistoryRec.LotSize", pHistoryRec.LotSize)
					lBidRec.Quantity = lBidRec.Quantity / int(pHistoryRec.LotSize)
					log.Println("lBidRec.Quantity:= ", lBidRec.Quantity)
					lBidArr = append(lBidArr, lBidRec)
					log.Println("lBidArr", lBidArr)
				}
			}
		}
	}
	log.Println("getHistoryDetails (-)")
	return lBidArr, nil
}

func CalcAmountPayable(pHistoryRec HistoryDetailStruct) (IpoHistoryCatList, error) {
	log.Println("CalcAmountPayable (+)")
	var Total float64

	lBidArr := pHistoryRec.CategoryList.AppliedDetail.AppliedBids
	lCategory := pHistoryRec.CategoryList

	lConfigFile := common.ReadTomlConfig("toml/IpoConfig.toml")
	lDiscountText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_DiscountText"])

	if lCategory.DiscountType == "A" {
		// Formula for absolute (quantity*(price-discountPrice))
		if len(lBidArr) == 1 {
			Total = float64(lBidArr[0].Quantity) * float64(lBidArr[0].Price-lCategory.DiscountPrice)

		} else if len(lBidArr) == 2 {
			Total = math.Max(float64(lBidArr[0].Quantity)*
				float64(lBidArr[0].Price-lCategory.DiscountPrice),
				float64(lBidArr[1].Quantity)*
					float64(lBidArr[1].Price-lCategory.DiscountPrice))

		} else if len(lBidArr) == 3 {
			Total = math.Max(math.Max(float64(lBidArr[0].Quantity)*
				float64(lBidArr[0].Price-lCategory.DiscountPrice),
				float64(lBidArr[1].Quantity)*
					float64(lBidArr[1].Price-lCategory.DiscountPrice)),
				float64(lBidArr[2].Quantity)*
					float64(lBidArr[2].Price-lCategory.DiscountPrice))

		}

	} else {
		// Formula for percentage quantity*(price - ((percentage/price)*100)
		if len(lBidArr) == 1 {
			Total = float64(lBidArr[0].Quantity) *
				float64(lBidArr[0].Price-(lBidArr[0].Price*(lCategory.DiscountPrice/100)))

		} else if len(lBidArr) == 2 {
			Total = math.Max(float64(lBidArr[0].Quantity)*
				float64(lBidArr[0].Price-(lBidArr[0].Price*(lCategory.DiscountPrice/100))),
				float64(lBidArr[1].Quantity)*
					float64(lBidArr[1].Price-(lBidArr[1].Price*(lCategory.DiscountPrice/100))))

		} else if len(lBidArr) == 3 {
			Total = math.Max(math.Max(float64(lBidArr[0].Quantity)*
				float64(lBidArr[0].Price-(lBidArr[0].Price*(lCategory.DiscountPrice/100))),
				float64(lBidArr[1].Quantity)*
					float64(lBidArr[1].Price-(lBidArr[1].Price*(lCategory.DiscountPrice/100)))),
				float64(lBidArr[2].Quantity)*
					float64(lBidArr[2].Price-(lBidArr[2].Price*(lCategory.DiscountPrice/100))))
		}
	}

	if lCategory.DiscountType == "P" {
		lCategory.DiscountText = strings.ReplaceAll(lDiscountText, "d", "%")
	} else {
		lCategory.DiscountText = strings.ReplaceAll(lDiscountText, "d", "₹")
	}

	lCategory.AppliedDetail.AppliedAmount = float64(Total)
	log.Println("CalcAmountPayable (-)")
	return lCategory, nil
}
