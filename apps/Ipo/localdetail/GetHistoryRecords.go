package localdetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/appsso"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

type HistoryRecStruct struct {
	MasterId      int               `json:"masterId"`
	Symbol        string            `json:"symbol"`
	Name          string            `json:"name"`
	ApplicationNo string            `json:"appNo"`
	IssueDate     string            `json:"issueDate"`
	IssueSize     int64             `json:"issueSize"`
	IssuePrice    string            `json:"issuePrice"`
	LotSize       int               `json:"lotSize"`
	Upi           string            `json:"upi"`
	Category      string            `json:"category"`
	Total         float64           `json:"total"`
	DpStatus      string            `json:"dpStatus"`
	UpiStatus     string            `json:"upiStatus"`
	ErrReason     string            `json:"errReason"`
	ModifyDetails []ModifyBidStruct `json:"modifyDetails"`
	Status        string            `json:"status"`
	ErrMsg        string            `json:"errMsg"`
	Sme           bool              `json:"sme"`
	RegistrarLink string            `json:"registrarLink"`
	Discount      string            `json:"discount"`
}

/*
Purpose:This api Method is used to get Application Details
Request:

Header Value: ID

Response:
=========
*On Sucess
=========

	{
		"modifyDetails" :[
		"appNo" : "1234567890"
		"upi" : "test@ybl"
		"category" : "IND"
			{
			"bidRefNo" : "45687843213233"
			"price" : 755
			"amount" : 14567
			"quantity" : 19
			"cutOff" : true
			},
			{
			"bidRefNo" : "45687843213233"
			"price" : 755
			"amount" : 14567
			"quantity" : 19
			"cutOff" : true
			}
		],
	"status" : "S",
	"errMsg" : ""
	}

=========
!On Error
=========

	{
		"modifyDetails" :[],
		"status": "E",
		"errMsg": "Can't able to get data from database"
	}

Author: Nithish Kumar
Date: 07JUNE2023
*/
func GetHistoryRecords(w http.ResponseWriter, r *http.Request) {
	log.Println("GetHistoryRecords(+)", r.Method)
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
	(w).Header().Set("Access-Control-Allow-Headers", "ID,NO,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		// Create an instance for ModifyBid
		var lRespRec HistoryRecStruct
		lRespRec.Status = common.SuccessCode

		// reading Header Value from the Request
		lMasterId := r.Header.Get("ID")
		lAppNo := r.Header.Get("NO")

		// log.Println("header value", lMasterId, lAppNo)
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lSessionClientId := ""
		lClientId := ""
		lLoggedBy := ""
		var lErr1 error
		lSessionClientId, lErr1 = appsso.ValidateAndGetClientDetails2(r, common.ABHIAppName, common.ABHICookieName)
		if lErr1 != nil {
			log.Println("LGHR01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGHR01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGHR01", "UserDetails Not Found!"))
			return
		} else {
			if lSessionClientId != "" {
				//get the detail for the client for whome we need to work
				lClientId = common.GetSetClient(lSessionClientId)
				//get the staff who logged
				lLoggedBy = common.GetLoggedBy(lSessionClientId)
				log.Println(lLoggedBy, lClientId)
			} else {
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lMasterId1, _ := strconv.Atoi(lMasterId)
		log.Println("fdgg", lMasterId1, lAppNo, lBrokerId)
		//call getHistoryDetails method to get the application Details
		lRespResult, lErr2 := getHistoryDetails(lMasterId1, lAppNo, lBrokerId)
		if lErr2 != nil {
			log.Println("LGHR02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGHR02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGHR02", "Unable to fetch your Records"))
			return
		} else {
			lRespRec = lRespResult
			lRespRec.Status = common.SuccessCode
		}
		// Marshal the response structure into lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGHR03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGHR03", "Issue in Getting your Application Details!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetHistoryRecords (-)", r.Method)
	}
}

/*
Pupose:This method used to retrieve the Application Details from the database.
Parameters:
PClientId,pMasterId
Response:
==========
*On Sucess
==========

	[
		{
			"appNo" : "1234567890"
			"upi" : "test@ybl"
			"category" : "IND"
			"bidRefNo" : "45687843213233"
			"price" : 755
			"amount" : 14567
			"quantity" : 19
			"cutOff" : true
		},
		{
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

Author:Pavithra
Date: 12JUNE2023
*/
func getHistoryDetails(pMasterId int, pAppNo string, pBrokerId int) (HistoryRecStruct, error) {
	log.Println("getHistoryDetails (+)")

	// This Variable is used to store the Each Modify bid Records
	var lRespRec ModifyBidStruct
	// var lAmount float64
	var lResponseRec HistoryRecStruct
	// var lLotSize, lQty int
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGHD01", lErr1)
		return lResponseRec, lErr1
	} else {
		defer lDb.Close()

		// lCoreString := `select
		// h.MasterId,
		// m.Symbol ,
		// m.Name ,
		// h.applicationNo,CONCAT(DATE_FORMAT(m.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(m.BiddingEndDate , '%d %b')) AS BidDates,
		// m.IssueSize Issue,
		// concat(m.MinPrice," - ", m.MaxPrice) PriceRange,
		// m.LotSize LotSize,
		// d.bidReferenceNo ,
		// h.upi ,
		// nvl(h.reason,''),
		// h.category ,
		// d.req_price,
		// d.req_amount,
		// d.req_quantity,
		// m.LotSize,
		// d.atCutOff,
		// d.activityType,
		// nvl(
		// 	(select xld.description
		// 	from xx_lookup_details xld,xx_lookup_header xlh
		// 	where xld.headerid = xlh.id
		// 	and xlh.Code = 'IpoDp'
		// 	and xld.Code = h.dpVerStatusFlag) , 'N/A'
		// 	) ApplicationStatus,
		// nvl(
		// 	(select (case when h.status = 'success' then xld.description else 'N/A' end)
		// 	from xx_lookup_details xld,xx_lookup_header xlh
		// 	where xld.headerid = xlh.id
		// 	and xlh.Code = 'IpoPay'
		// 	and xld.Code = h.upiPaymentStatusFlag), 'N/A'
		// ) UPIStatus,
		// nvl(
		// 	(select(case when  m.Isin = ad.Isin then ad.category else 1 end)
		// 	from a_ipo_details ad
		// 	where ad.Isin  = m.Isin	),1
		// ) sme,
		// NVL(
		// 	(SELECT NVL(R.RegistrarLink, '')
		// 	 FROM a_ipo_registrars R
		// 	 WHERE R.RegistrarName = m.Registrar),
		// 	''
		// ) AS Registrar
		// from
		// a_ipo_order_header h,a_ipo_orderdetails d, a_ipo_master m
		// where m.Id = h.MasterId
		// and h.Id = d.headerId
		// and h.MasterId = ?
		// and h.applicationNo = ?
		// and h.brokerId = ?`

		//  In this Below Query inline Query is Added by prashanth To Get Registrar Link To show in iframe By prashanth
		lCoreString := `select
                        h.MasterId,
                        m.Symbol ,
                        m.Name ,
                        h.applicationNo,CONCAT(DATE_FORMAT(m.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(m.BiddingEndDate , '%d %b')) AS BidDates,
                        m.IssueSize Issue,
                        concat(m.MinPrice," - ", m.MaxPrice) PriceRange,
                        m.LotSize LotSize,
						d.id,
                        d.bidReferenceNo ,
                        h.upi ,
                        nvl(h.reason,''),
                        h.category ,
                        d.req_price,
                        d.req_amount,
                        d.req_quantity,
                        m.LotSize,
                        d.atCutOff,
                        d.activityType,
                        nvl(
                        	(select xld.description 
                        	from xx_lookup_details xld,xx_lookup_header xlh
                        	where xld.headerid = xlh.id
                        	and xlh.Code = 'IpoDp'
                        	and xld.Code = h.dpVerStatusFlag) , 'N/A'
                        	) ApplicationStatus,
                        nvl(
                        	(select (case when h.status = 'success' then xld.description else 'N/A' end)
                        	from xx_lookup_details xld,xx_lookup_header xlh
                        	where xld.headerid = xlh.id
                        	and xlh.Code = 'IpoPay'
                        	and xld.Code = h.upiPaymentStatusFlag), 'N/A'
                        ) UPIStatus,
                        nvl(
                        	(select(case when  m.Isin = ad.Isin then ad.category else 1 end)
                        	from a_ipo_details ad
                        	where ad.Isin  = m.Isin	),1
                        ) sme,
						NVL(
							(SELECT NVL(R.RegistrarLink, '')
							 FROM a_ipo_registrars R
							 WHERE R.RegistrarName = m.Registrar),
							''
						) AS Registrar
						from
						a_ipo_order_header h,a_ipo_orderdetails d, a_ipo_master m
                        where m.Id = h.MasterId 
                        and h.Id = d.headerId
                        and h.MasterId = ?
                        and h.applicationNo = ?
                        and h.brokerId = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pAppNo, pBrokerId)
		if lErr2 != nil {
			log.Println("LGHD02", lErr2)
			return lResponseRec, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in ModifyRespStruct
			for lRows.Next() {
				lErr3 := lRows.Scan(&lResponseRec.MasterId, &lResponseRec.Symbol, &lResponseRec.Name, &lResponseRec.ApplicationNo, &lResponseRec.IssueDate, &lResponseRec.IssueSize, &lResponseRec.IssuePrice,
					&lResponseRec.LotSize, &lRespRec.Id, &lRespRec.BidRefNo, &lResponseRec.Upi, &lResponseRec.ErrReason, &lResponseRec.Category,
					&lRespRec.Price, &lRespRec.Amount, &lRespRec.Quantity, &lResponseRec.LotSize, &lRespRec.CutOff, &lRespRec.ActivityType, &lResponseRec.DpStatus, &lResponseRec.UpiStatus, &lResponseRec.Sme, &lResponseRec.RegistrarLink)

				if lErr3 != nil {
					log.Println("LGHD03", lErr3)
					return lResponseRec, lErr3
				} else {
					// Append the Bid Records in lRespArr
					// lAmount = lRespRec.Amount
					// lResponseRec.Total += lAmount
					// lRespRec.Quantity = int(lQty) / int(lLotSize)
					lResponseRec.ModifyDetails = append(lResponseRec.ModifyDetails, lRespRec)

				}
			}

			lCategoryArr, lErr4 := GetCategoryDiscount("", pMasterId)
			if lErr4 != nil {
				log.Println("LGHD04", lErr4)
				return lResponseRec, lErr4
			} else {
				lResult, lErr5 := CalcDiscount(lResponseRec, lCategoryArr)
				if lErr5 != nil {
					log.Println("LGHD05", lErr5)
					return lResponseRec, lErr5
				} else {
					// lResponseRec.Total = lTotal
					lResponseRec = lResult
				}
			}
			// if len(lResponseRec.ModifyDetails) == 1 {
			// 	lResponseRec.Total = lResponseRec.ModifyDetails[0].Amount
			// } else if len(lResponseRec.ModifyDetails) == 2 {
			// 	lResponseRec.Total = math.Max(lResponseRec.ModifyDetails[0].Amount, lResponseRec.ModifyDetails[1].Amount)
			// } else if len(lResponseRec.ModifyDetails) == 3 {
			// 	lResponseRec.Total = math.Max(math.Max(lResponseRec.ModifyDetails[0].Amount, lResponseRec.ModifyDetails[1].Amount), lResponseRec.ModifyDetails[2].Amount)
			// }

		}
	}
	log.Println("getHistoryDetails (-)")
	return lResponseRec, nil
}

func CalcDiscount(pHistoryRec HistoryRecStruct, pCategoryArr []CategoryStruct) (HistoryRecStruct, error) {
	log.Println("CalcDiscount (+)")
	var Total float64
	for _, lCategory := range pCategoryArr {
		if pHistoryRec.Category == lCategory.Value {

			if lCategory.DiscountType == "A" {
				// Formula for absolute (quantity*(price-discountPrice))
				if len(pHistoryRec.ModifyDetails) == 1 {
					Total = float64(pHistoryRec.ModifyDetails[0].Quantity) * (pHistoryRec.ModifyDetails[0].Price - lCategory.DiscountPrice)
					//
					pHistoryRec.ModifyDetails[0].Quantity = int(pHistoryRec.ModifyDetails[0].Quantity) / int(pHistoryRec.LotSize)

				} else if len(pHistoryRec.ModifyDetails) == 2 {
					Total = math.Max(float64(pHistoryRec.ModifyDetails[0].Quantity)*
						(pHistoryRec.ModifyDetails[0].Price-lCategory.DiscountPrice),
						float64(pHistoryRec.ModifyDetails[1].Quantity)*
							(pHistoryRec.ModifyDetails[1].Price-lCategory.DiscountPrice))
					//
					pHistoryRec.ModifyDetails[0].Quantity = int(pHistoryRec.ModifyDetails[0].Quantity) / int(pHistoryRec.LotSize)
					pHistoryRec.ModifyDetails[1].Quantity = int(pHistoryRec.ModifyDetails[1].Quantity) / int(pHistoryRec.LotSize)

				} else if len(pHistoryRec.ModifyDetails) == 3 {
					Total = math.Max(math.Max(float64(pHistoryRec.ModifyDetails[0].Quantity)*
						(pHistoryRec.ModifyDetails[0].Price-lCategory.DiscountPrice),
						float64(pHistoryRec.ModifyDetails[1].Quantity)*
							(pHistoryRec.ModifyDetails[1].Price-lCategory.DiscountPrice)),
						float64(pHistoryRec.ModifyDetails[2].Quantity)*
							(pHistoryRec.ModifyDetails[2].Price-lCategory.DiscountPrice))

					//
					pHistoryRec.ModifyDetails[0].Quantity = int(pHistoryRec.ModifyDetails[0].Quantity) / int(pHistoryRec.LotSize)
					pHistoryRec.ModifyDetails[1].Quantity = int(pHistoryRec.ModifyDetails[1].Quantity) / int(pHistoryRec.LotSize)
					pHistoryRec.ModifyDetails[2].Quantity = int(pHistoryRec.ModifyDetails[2].Quantity) / int(pHistoryRec.LotSize)

				}
				if lCategory.DiscountPrice == 0 {
					pHistoryRec.Discount = "N/A"
				} else {
					pHistoryRec.Discount = "₹" + strconv.FormatFloat(lCategory.DiscountPrice, 'f', -1, 64)
				}
			} else {
				// Formula for percentage quantity*(price - ((percentage/price)*100)
				if len(pHistoryRec.ModifyDetails) == 1 {
					Total = float64(pHistoryRec.ModifyDetails[0].Quantity) *
						(pHistoryRec.ModifyDetails[0].Price - (pHistoryRec.ModifyDetails[0].Price * (lCategory.DiscountPrice / 100)))
						//
					pHistoryRec.ModifyDetails[0].Quantity = int(pHistoryRec.ModifyDetails[0].Quantity) / int(pHistoryRec.LotSize)
				} else if len(pHistoryRec.ModifyDetails) == 2 {
					Total = math.Max(float64(pHistoryRec.ModifyDetails[0].Quantity)*
						(pHistoryRec.ModifyDetails[0].Price-(pHistoryRec.ModifyDetails[0].Price*(lCategory.DiscountPrice/100))),
						float64(pHistoryRec.ModifyDetails[1].Quantity)*
							(pHistoryRec.ModifyDetails[1].Price-(pHistoryRec.ModifyDetails[1].Price*(lCategory.DiscountPrice/100))))
					//
					pHistoryRec.ModifyDetails[0].Quantity = int(pHistoryRec.ModifyDetails[0].Quantity) / int(pHistoryRec.LotSize)
					pHistoryRec.ModifyDetails[1].Quantity = int(pHistoryRec.ModifyDetails[1].Quantity) / int(pHistoryRec.LotSize)

				} else if len(pHistoryRec.ModifyDetails) == 3 {
					Total = math.Max(math.Max(float64(pHistoryRec.ModifyDetails[0].Quantity)*
						(pHistoryRec.ModifyDetails[0].Price-(pHistoryRec.ModifyDetails[0].Price*(lCategory.DiscountPrice/100))),
						float64(pHistoryRec.ModifyDetails[1].Quantity)*
							(pHistoryRec.ModifyDetails[1].Price-(pHistoryRec.ModifyDetails[1].Price*(lCategory.DiscountPrice/100)))),
						float64(pHistoryRec.ModifyDetails[2].Quantity)*
							(pHistoryRec.ModifyDetails[2].Price-(pHistoryRec.ModifyDetails[2].Price*(lCategory.DiscountPrice/100))))
					//
					pHistoryRec.ModifyDetails[0].Quantity = int(pHistoryRec.ModifyDetails[0].Quantity) / int(pHistoryRec.LotSize)
					pHistoryRec.ModifyDetails[1].Quantity = int(pHistoryRec.ModifyDetails[1].Quantity) / int(pHistoryRec.LotSize)
					pHistoryRec.ModifyDetails[2].Quantity = int(pHistoryRec.ModifyDetails[2].Quantity) / int(pHistoryRec.LotSize)
				}
				if lCategory.DiscountPrice == 0 {
					pHistoryRec.Discount = "N/A"
				} else {
					pHistoryRec.Discount = strconv.FormatFloat(lCategory.DiscountPrice, 'f', -1, 64) + "%%"
				}
			}
		}
	}
	pHistoryRec.Total = Total
	log.Println("CalcDiscount (-)")
	return pHistoryRec, nil
}
