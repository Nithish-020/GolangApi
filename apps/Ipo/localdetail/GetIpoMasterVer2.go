package localdetail

import (
	"encoding/json"
	"fcs23pkg/apps/Ipo/Function"
	"fcs23pkg/apps/Ipo/brokers"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Response Structure for GetIpoMaster API
type IpoRespStruct struct {
	IpoDetails  []IpoDetailStruct `json:"ipoDetail"`
	PolicyText  string            `json:"policyText"`
	SuggestUPI  string            `json:"suggestUPI"`
	NoDataText  string            `json:"masterNoDataText"`
	MasterFound string            `json:"masterFound"`
	InvestCount int               `json:"investCount"`
	Disclaimer  []string          `json:"disclaimer"`
	Status      string            `json:"status"`
	ErrMsg      string            `json:"errMsg"`
}

// This Structure is used to Collect the Active IPO informations
type IpoDetailStruct struct {
	Id                int               `json:"id"`
	Name              string            `json:"name"`
	Symbol            string            `json:"symbol"`
	DateRange         string            `json:"dateRange"`
	PriceRange        string            `json:"priceRange"`
	MinBidQuantity    int               `json:"minBidQuantity"`
	MinPrice          float64           `json:"minPrice"`
	CutOffPrice       float64           `json:"cutOffPrice"`
	LotSize           float64           `json:"lotSize"`
	IssueSize         float64           `json:"issueSize"`
	IssueSizeWithText string            `json:"issueSizeWithText"`
	CutOffFlag        bool              `json:"cutOffFlag"`
	ActionFlag        string            `json:"actionFlag"`
	PurchaseFlag      string            `json:"ipoPurchased"`
	ButtonText        string            `json:"buttonText"`
	Sme               bool              `json:"sme"`
	SmeText           string            `json:"smeText"`
	BlogLink          string            `json:"blogLink"`
	DrhpLink          string            `json:"drhpLink"`
	DisableActionBtn  bool              `json:"disableActionBtn"`
	Exchange          string            `json:"exchange"`
	CategoryList      []IpoCategoryList `json:"categoryList"`
}
type IpoCategoryList struct {
	Category      string        `json:"category"`
	Code          string        `json:"code"`
	PurchaseFlag  string        `json:"purchaseFlag"`
	ApplicationNo string        `json:"applicationNo"`
	MaxValue      int           `json:"maxValue"`
	DiscountText  string        `json:"discountText"`
	DiscountPrice int           `json:"discountPrice"`
	DiscountType  string        `json:"discountType"`
	ModifyAllowed bool          `json:"modifyAllowed"`
	CancelAllowed bool          `json:"cancelAllowed"`
	ShowDiscount  bool          `json:"showDiscount"`
	Infotext      string        `json:"infoText"`
	AppliedDetail AppliedDetail `json:"appliedDetail"`
}

type AppliedDetail struct {
	AppliedUPI    string      `json:"appliedUPI"`
	AppliedAmount float64     `json:"appliedAmount"`
	AppliedBids   []BidDetail `json:"appliedBids"`
}

type BidDetail struct {
	Id             int     `json:"id"`
	ActivityType   string  `json:"activityType"`
	BidReferenceNo string  `json:"bidReferenceNo"`
	Quantity       int     `json:"quantity"`
	Price          int     `json:"price"`
	Amount         float64 `json:"amount"`
	CutOff         bool    `json:"cutOff"`
}

/*
Pupose:This Function is used to Get the Active Ipo Details in our database table ....
Parameters:

not Applicable

Response:

*ON Sucess
=========

	{
		"IpoDetails": [
			{
				"id": 18,
				"symbol": "MMIPO26",
				"startDate": "2023-06-02",
				"endDate": "2023-06-30",
				"priceRange": "1000 - 2000",
				"cutOffPrice": 2000,
				"minBidQuantity": 10,
				"applicationStatus": "Pending",
				"upiStatus": "Accepted BY Investor"
			},
			{
				"id": 10,
				"symbol": "fixed",
				"startDate": "2023-05-10",
				"endDate": "2023-08-29",
				"priceRange": "755 - 755",
				"cutOffPrice": 755,
				"minBidQuantity": 100,
				"applicationStatus": "-",
				"upiStatus": "-"
			}
		],
		"status": "S",
		"errMsg": ""
	}

!ON Error
========

	{
		"status": E,
		"reason": "Can't able to get the data FROM database"
	}

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetIpoMaster2(w http.ResponseWriter, r *http.Request) {
	log.Println("GetIpoMaster(+)", r.Method)
	origin := r.Header.Get("Origin")
	var lBrokerId int
	// var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, _ = brokers.GetBrokerId(origin) // TO get brokerId
			// log.Println(lErr, origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		// create the instance for IpoStruct
		var lRespRec IpoRespStruct
		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("LGIM01", lErr1)
			lRespRec.Status = common.LoginFailure
			lRespRec.ErrMsg = "LGIM01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIM01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				log.Println("LGIM02", lErr1)
				lRespRec.Status = common.LoginFailure
				lRespRec.ErrMsg = "LGIM02" + lErr1.Error()
				fmt.Fprintf(w, helpers.GetErrorString("LGIM02", "Session Expired,Please logout and Login again !!"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		lIpoResp, lErr2 := GetIpoDetail(lClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("LGIM02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIM02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIM02", "Unable to fetch the master details"))
			return
		} else {
			lRespRec = lIpoResp
		}

		// Marshaling the response structure to lData
		lData, lErr3 := json.Marshal(lRespRec)
		if lErr3 != nil {
			log.Println("LGIM03", lErr3)
			fmt.Fprintf(w, helpers.GetErrorString("LGIM03", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetIpoMaster (-)", r.Method)
	}
}

func GetIpoDetail(pClientId string, pBrokerId int) (IpoRespStruct, error) {
	log.Println("GetIpoDetail (+)")
	var lFinalIpoArr []IpoDetailStruct
	var lIpoResp IpoRespStruct

	lIpoResp.Status = common.SuccessCode

	lConfigFile := common.ReadTomlConfig("toml/IpoConfig.toml")
	lSmeText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_SmeText"])
	lNoDataText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_MasterNoDataText"])
	lPolicyText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_PolicyText"])

	lIpoArr, lErr1 := GetActiveIpo(pClientId, pBrokerId)
	if lErr1 != nil {
		lIpoResp.Status = common.ErrorCode
		log.Println("LGID01", lErr1)
		return lIpoResp, lErr1
	} else {

		for _, lDetail := range lIpoArr {

			lOrderRec, lErr2 := getClientOrderRec(lDetail, pBrokerId, pClientId)
			if lErr2 != nil {
				lIpoResp.Status = common.ErrorCode
				log.Println("LGID02", lErr2)
				return lIpoResp, lErr2
			} else {

				lBtnText, lErr3 := getIpoDiscription(lOrderRec.ActionFlag)
				if lErr3 != nil {
					lIpoResp.Status = common.ErrorCode
					log.Println("LGID03", lErr3)
					return lIpoResp, lErr3
				} else {
					lOrderRec.ButtonText = lBtnText
				}

				// Display the sme text only if it is SME and
				// change the shares count in the text dynamically.
				lStringText := strings.ReplaceAll(lSmeText, "0", strconv.FormatFloat(lOrderRec.LotSize, 'f', -1, 64))
				if lOrderRec.Sme == true {
					lOrderRec.SmeText = lStringText
				}
				lFinalIpoArr = append(lFinalIpoArr, lOrderRec)
			}
			// Count how many Ipo does a client invested in !
			if lOrderRec.PurchaseFlag == "Y" {
				lIpoResp.InvestCount++
			}
		}

		// Custom sorting function
		sort.Slice(lFinalIpoArr, func(i, j int) bool {
			orderMap := map[string]int{"M": 1, "PA": 2, "U": 3}
			orderI, existsI := orderMap[lFinalIpoArr[i].ActionFlag]
			orderJ, existsJ := orderMap[lFinalIpoArr[j].ActionFlag]

			// If both elements have valid ActionFlag, compare them
			if existsI && existsJ {
				return orderI < orderJ
			}

			// If one element has valid ActionFlag, it comes first
			if existsI {
				return true
			}
			if existsJ {
				return false
			}

			// If neither element has valid ActionFlag, maintain the original order
			return i < j
		})

		// Assign the final Ipo details to show the master in UI
		lIpoResp.IpoDetails = lFinalIpoArr

		// Display the Policy before applying the new Ipo
		lIpoResp.PolicyText = lPolicyText

		// Get the lastly used UPI id by the client
		// to prefill the UPI field in the Bid dialog
		lUPI, lErr4 := getLastlyUsedUPI(pClientId)
		if lErr4 != nil {
			lIpoResp.Status = common.ErrorCode
			log.Println("LGID04", lErr4)
			return lIpoResp, lErr4
		} else {
			lIpoResp.SuggestUPI = lUPI
		}
		// Get the disclaimer given by the broker
		// to know the detail discription about the IPO
		lDisclaimerArr, lErr5 := getIpoDisclaimer(pBrokerId)
		if lErr5 != nil {
			lIpoResp.Status = common.ErrorCode
			log.Println("LGID05", lErr5)
			return lIpoResp, lErr5
		} else {
			lIpoResp.Disclaimer = lDisclaimerArr
		}

		if lIpoResp.IpoDetails != nil {
			lIpoResp.MasterFound = "Y"
		} else {
			lIpoResp.MasterFound = "N"
			lIpoResp.NoDataText = lNoDataText
		}

	}
	log.Println("GetIpoDetail (-)")
	return lIpoResp, nil
}

func GetActiveIpo(pClientId string, pBrokerId int) ([]IpoDetailStruct, error) {
	log.Println("GetActiveIpo (+)")
	// create the instance for IpoDetailStruct
	var lIpoDataRec IpoDetailStruct
	var lIpoDataArr []IpoDetailStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMD01", lErr1)
		return lIpoDataArr, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `
						SELECT m.Id, m.Name, m.Symbol,
						CONCAT(DATE_FORMAT(m.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(m.BiddingEndDate , '%d %b')) AS dateRange,
						concat(m.MinPrice," - ", m.MaxPrice) PriceRange,m.MinBidQuantity,
						m.MinPrice, m.MaxPrice, m.LotSize,m.IssueSize,
						(CASE
							WHEN m.IssueSize >= 10000000 THEN CONCAT(FORMAT(m.IssueSize / 10000000, 2), ' Crores')
							WHEN m.IssueSize >= 100000 THEN CONCAT(FORMAT(m.IssueSize / 100000, 2), ' Lakhs')
							ELSE m.IssueSize
						end) AS IssueSize,
						(CASE WHEN m.Id = s.MasterId AND s.AllowCutOff = 1 THEN 1 else 0 END)AllowCutOff,
						(case
							when date_sub(m.BiddingStartDate,interval 1 day ) = curdate() then 'PA'
							when m.BiddingStartDate > curdate() then 'U'
							when m.BiddingEndDate = Date(now()) AND Time(now()) > '17:00:00' then 'CL'
							else ''
						end) ActionFlag,
						NVL(
							(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.category else 1 END)
							FROM a_ipo_details ad
							WHERE ad.Isin  = m.Isin	),1
						)SubType,
						NVL(
							(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.detailsLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
						)dLink,
						NVL(
							(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.drhpLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
							)drhpLink, m.Exchange
						FROM a_ipo_master m ,a_ipo_categories c,a_ipo_subcategory s
							WHERE m.Id= s.MasterId AND m.Id = c.MasterId AND m.IssueType = "EQUITY"
						AND c.code= "RETAIL" AND s.CaCode = "RETAIL" AND s.SubCatCode = "IND"
						AND s.AllowUpi = 1 AND m.BiddingEndDate >= Curdate() AND m.Exchange = 'NSE'
						AND not exists (
							SELECT 1
							FROM (
								SELECT NVL(c.EndTime, m.DailyEndTime) endTime, c.MasterId
								FROM a_ipo_master m,a_ipo_categories c
								WHERE m.Id = c.MasterId
								AND c.Code = 'RETAIL'
							AND m.BiddingEndDate = date(now())
							) a, a_ipo_master m1
						WHERE m1.Id = a.masterId
						AND a.masterId = m.Id
						AND a.endTime <= time(now())
						)
						union
						SELECT m.Id, m.Name, m.Symbol,
						CONCAT(DATE_FORMAT(m.BiddingStartDate , '%d %b'), ' - ', DATE_FORMAT(m.BiddingEndDate , '%d %b')) AS dateRange,
						concat(m.MinPrice," - ", m.MaxPrice) PriceRange,m.MinBidQuantity,
						m.MinPrice, m.MaxPrice, m.LotSize,m.IssueSize,
						(CASE
							WHEN m.IssueSize >= 10000000 THEN CONCAT(FORMAT(m.IssueSize / 10000000, 2), ' Crores')
							WHEN m.IssueSize >= 100000 THEN CONCAT(FORMAT(m.IssueSize / 100000, 2), ' Lakhs')
							ELSE m.IssueSize
						end) AS IssueSize,
						(CASE WHEN m.Id = s.MasterId AND s.AllowCutOff = 1 THEN 1 else 0 END)AllowCutOff,
						(case
							when date_sub(m.BiddingStartDate,interval 1 day ) = curdate() then 'PA'
							when m.BiddingStartDate > curdate() then 'U'
							when m.BiddingEndDate = Date(now()) AND Time(now()) > '17:00:00' then 'CL'
							else ''
						end) ActionFlag,
						NVL(
							(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.category else 1 END)
							FROM a_ipo_details ad
							WHERE ad.Isin  = m.Isin	),1
						)SubType,
						NVL(
							(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.detailsLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
						)dLink,
						NVL(
							(SELECT(CASE WHEN  m.Isin = ad.Isin THEN ad.drhpLink else '' END)
						FROM a_ipo_details ad
						WHERE ad.Isin  = m.Isin	),""
							)drhpLink, m.Exchange
						FROM a_ipo_master m ,a_ipo_categories c,a_ipo_subcategory s
							WHERE m.Id= s.MasterId AND m.Id = c.MasterId AND m.IssueType = "EQUITY"
						AND c.code= "RETAIL" AND s.CaCode = "RETAIL" AND s.SubCatCode = "IND"
						AND s.AllowUpi = 1 AND m.BiddingEndDate >= Curdate() AND m.Exchange = 'BSE'
						AND not exists (
							SELECT 1
							FROM (
								SELECT NVL(c.EndTime, m.DailyEndTime) endTime, c.MasterId
								FROM a_ipo_master m,a_ipo_categories c
								WHERE m.Id = c.MasterId
								AND c.Code = 'RETAIL'
							AND m.BiddingEndDate = date(now())
							) a, a_ipo_master m1
						WHERE m1.Id = a.masterId
						AND a.masterId = m.Id
						AND a.endTime <= time(now())
						)`

		lRows, lErr3 := lDb.Query(lCoreString)
		if lErr3 != nil {
			log.Println("LGIMD03", lErr3)
			return lIpoDataArr, lErr3

		} else {
			//This for loop is used to collect the records FROM the database AND store them in structure
			for lRows.Next() {
				lErr4 := lRows.Scan(&lIpoDataRec.Id, &lIpoDataRec.Name, &lIpoDataRec.Symbol, &lIpoDataRec.DateRange, &lIpoDataRec.PriceRange, &lIpoDataRec.MinBidQuantity, &lIpoDataRec.MinPrice, &lIpoDataRec.CutOffPrice, &lIpoDataRec.LotSize, &lIpoDataRec.IssueSize, &lIpoDataRec.IssueSizeWithText, &lIpoDataRec.CutOffFlag, &lIpoDataRec.ActionFlag, &lIpoDataRec.Sme, &lIpoDataRec.BlogLink, &lIpoDataRec.DrhpLink, &lIpoDataRec.Exchange)
				if lErr4 != nil {
					log.Println("LGIMD04", lErr4)
					return lIpoDataArr, lErr4

				} else {
					// Append the IPO Records in lRespRec.IpoDetails Array
					lIpoDataArr = append(lIpoDataArr, lIpoDataRec)
				}
			}
		}
	}
	log.Println("GetActiveIpo (-)")
	return lIpoDataArr, nil
}

func getClientOrderRec(pDetail IpoDetailStruct, pBrokerId int, pClientId string) (IpoDetailStruct, error) {
	log.Println("getClientOrderRec (+)")
	var lCategoryRec IpoCategoryList
	var lCategoryArr []IpoCategoryList
	var lCancelFlag, lStatus, lUPI string
	lHeaderId := 0

	lConfigFile := common.ReadTomlConfig("toml/IpoConfig.toml")
	lShowDiscount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_ShowDiscount"])
	lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_CancelAllowed"])
	lModifyAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_ModifyAllowed"])
	lDiscountText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_DiscountText"])
	lInfoText := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["IPO_InfoText"])

	//if order exists or not, set the defaults for the record
	pDetail.DisableActionBtn = false
	pDetail.PurchaseFlag = "N"

	// Make the default value of modifyAllowed using the toml
	if lModifyAllowed == "Y" {
		lCategoryRec.ModifyAllowed = true
	} else {
		lCategoryRec.ModifyAllowed = false
	}
	// Make the default value of cancelAllowed using the toml
	if lCancelAllowed == "Y" {
		lCategoryRec.CancelAllowed = true
	} else {
		lCategoryRec.CancelAllowed = false
	}
	// Make the default value of showDiscount using the toml
	if lShowDiscount == "Y" {
		lCategoryRec.ShowDiscount = true
	} else {
		lCategoryRec.ShowDiscount = false
	}

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMD01", lErr1)
	} else {
		defer lDb.Close()

		lMasterId := strconv.Itoa(pDetail.Id)
		lBrokerId := strconv.Itoa(pBrokerId)

		lCoreString := `select s.SubCatCode ,nvl(tab.headerId,0),nvl(tab.cancelFlag,''),nvl(tab.status,''),
						nvl(tab.applicationNo,''),s.maxValue,s.discountPrice,nvl(s.discountType,'')discuntType,
						nvl(tab.upi,'') upi
						from v_ipo_subcategory s
						left join
						(
							select h.id headerId, h.masterId masterId, h.category category,
							h.cancelFlag cancelFlag,h.status status,
							h.applicationNo applicationNo,
							h.upi upi
							from a_ipo_order_header h
							where h.cancelFlag <> 'Y'
							and h.status <> 'failed'
							and h.clientId = '` + pClientId + `'
							and h.brokerId = ` + lBrokerId + `
						)tab
						on s.id = tab.masterId
						and s.SubCatCode = tab.category
						where s.id = ` + lMasterId + `
						order by 
						(CASE
						WHEN s.SubCatCode  = 'IND' THEN 1
						WHEN s.SubCatCode = 'EMP' THEN 2
						WHEN s.SubCatCode = 'SHA'  THEN 3
						else 4 END)`
		lRows, lErr3 := lDb.Query(lCoreString)
		if lErr3 != nil {
			log.Println("LGIMD03", lErr3)
			return pDetail, lErr3

		} else {
			//This for loop is used to collect the records FROM the database and store them in structure
			for lRows.Next() {
				lErr4 := lRows.Scan(&lCategoryRec.Code, &lHeaderId, &lCancelFlag, &lStatus, &lCategoryRec.ApplicationNo, &lCategoryRec.MaxValue, &lCategoryRec.DiscountPrice, &lCategoryRec.DiscountType, &lUPI)
				if lErr4 != nil {
					log.Println("LGIMD04", lErr4)
					return pDetail, lErr4

				} else {
					lInfoArr := strings.Split(lInfoText, ".")

					if lShowDiscount == "Y" {
						if lCategoryRec.DiscountType == "P" {
							lCategoryRec.DiscountText = strings.ReplaceAll(lDiscountText, "d", "%")
						} else {
							lCategoryRec.DiscountText = strings.ReplaceAll(lDiscountText, "d", "₹")
						}
					}
					// If the symbol is upcoming
					if pDetail.ActionFlag == "U" || pDetail.ActionFlag == "CL" {
						pDetail.DisableActionBtn = true
						lCategoryRec.PurchaseFlag = "N"
						pDetail.PurchaseFlag = "N"

						// If the symbol is preapply or current
					} else {

						if lCancelFlag == "N" && (lStatus == common.SUCCESS || lStatus == common.PENDING) {
							lCategoryRec.PurchaseFlag = "Y"
							pDetail.ActionFlag = "M"
						} else {
							lCategoryRec.PurchaseFlag = "N"
						}

						lIndicator, lErr5 := Function.GetTime(pDetail.Id)
						if lErr5 != nil {
							log.Println("LGIMD05", lErr5)
							return pDetail, lErr5

						} else {
							// If the Incicator holds a value "False" then
							// we can apply the bid direclty in the exchange
							if lIndicator == "False" {

								// if the client places the bid in market period then
								// action button needs to show modify and allow the user to modify or cancel
								if lCategoryRec.PurchaseFlag == "Y" {
									lCategoryRec.ModifyAllowed = true
									lCategoryRec.CancelAllowed = true
									lCategoryRec.Infotext = lInfoArr[0] + "," + lInfoArr[1]

									// if the client places the bid in offline then
									// action button needs to show modify and the application dosen't processed on the exchange
									// then allow the user to modify or cancel
									if lStatus == common.PENDING {
										lCategoryRec.ModifyAllowed = false
										lCategoryRec.CancelAllowed = false
										lCategoryRec.Infotext = lInfoArr[0] + "," + lInfoArr[1]
									}

									// if the client places dosen't places the bid then
									// allow the user to modify but not to cancel
								} else {
									lCategoryRec.ModifyAllowed = true
									lCategoryRec.CancelAllowed = false
									lCategoryRec.Infotext = lInfoArr[0]

									// If the exchange time has open to appply then
									// show the action button has Bid
									if pDetail.ActionFlag == "" {
										pDetail.ActionFlag = "B"
									}
								}

								// If the Incicator holds a value "True" then
								// we can apply the bid only in the local database
							} else {

								// if the client places already places the bid then the
								// action button needs to show modify and dosen't allow the user to modify or cancel
								if lCategoryRec.PurchaseFlag == "Y" {
									lCategoryRec.ModifyAllowed = false
									lCategoryRec.CancelAllowed = false
									lCategoryRec.Infotext = lInfoArr[2]

									// if the client places dosen't places the bid then the
									// allow the user to modify but not to cancel
								} else {
									lCategoryRec.ModifyAllowed = true
									lCategoryRec.CancelAllowed = false
									lCategoryRec.Infotext = lInfoArr[0]

									// If the exchange time has remain closed then
									// show the action button has Offline
									if pDetail.ActionFlag == "" {
										pDetail.ActionFlag = "O"
									}
								}
							}
						}

						// Collecting the bids detail based on the client applied
						lBidArr, lAmount, lErr7 := getClientOrderedBid(lHeaderId, int(pDetail.LotSize))
						if lErr7 != nil {
							log.Println("LGIMD07", lErr7)
							return pDetail, lErr7

						} else {
							lCategoryRec.AppliedDetail.AppliedUPI = lUPI
							lCategoryRec.AppliedDetail.AppliedAmount = lAmount

							// Collecting the bids deatil according to the category
							lCategoryRec.AppliedDetail.AppliedBids = lBidArr

							if lCategoryRec.AppliedDetail.AppliedBids == nil {
								lCategoryRec.AppliedDetail.AppliedBids = []BidDetail{}
							}

						}
						lCategory, lErr6 := getIpoDiscription(lCategoryRec.Code)
						if lErr6 != nil {
							log.Println("LGIMD06", lErr6)
						} else {
							lCategoryRec.Category = lCategory
						}
						// Append the category wise Records in lCategoryArr Array
						lCategoryArr = append(lCategoryArr, lCategoryRec)
					}
				}
			}

			// This loop helps to get the overall purchase flag for the particular Ipo
			for _, lFlag := range lCategoryArr {
				if lFlag.PurchaseFlag == "Y" {
					pDetail.PurchaseFlag = lFlag.PurchaseFlag
				}
			}

			pDetail.CategoryList = lCategoryArr

			if pDetail.CategoryList == nil {
				pDetail.CategoryList = []IpoCategoryList{}
			}
		}
	}
	log.Println("getClientOrderRec (-)")
	return pDetail, nil
}

func getClientOrderedBid(pHeaderId int, pLotSize int) ([]BidDetail, float64, error) {
	var lBidDetailArr []BidDetail
	var lBidDetailRec BidDetail
	lAmount := 0.0
	if pHeaderId != 0 {
		// Calling LocalDbConect method in ftdb to estabish the database connection
		lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr1 != nil {
			log.Println("LGCOB01", lErr1)
			return lBidDetailArr, lAmount, lErr1
		} else {
			defer lDb.Close()
			lCoreString := `select d.id,d.activityType ,d.bidReferenceNo ,d.req_quantity ,d.req_price,d.req_amount ,d.atCutOff 
						from a_ipo_order_header h,a_ipo_orderdetails d
						where h.id = d.headerId
						and h.id = ?`
			lRows, lErr2 := lDb.Query(lCoreString, pHeaderId)
			if lErr2 != nil {
				log.Println("LGCOB02", lErr2)
				return lBidDetailArr, lAmount, lErr2
			} else {
				//This for loop is used to collect the records from the database and store them in respective structure
				for lRows.Next() {
					lErr3 := lRows.Scan(&lBidDetailRec.Id, &lBidDetailRec.ActivityType, &lBidDetailRec.BidReferenceNo, &lBidDetailRec.Quantity, &lBidDetailRec.Price, &lBidDetailRec.Amount, &lBidDetailRec.CutOff)
					if lErr3 != nil {
						log.Println("LGCOB03", lErr3)
						return lBidDetailArr, lAmount, lErr3
					} else {
						lBidDetailRec.Quantity = lBidDetailRec.Quantity / pLotSize
						lBidDetailArr = append(lBidDetailArr, lBidDetailRec)
					}
				}

				// To get the Amount payable based on the highest value of amount
				if len(lBidDetailArr) == 1 {
					lAmount = lBidDetailArr[0].Amount
				} else if len(lBidDetailArr) == 2 {
					lAmount = math.Max(lBidDetailArr[0].Amount, lBidDetailArr[1].Amount)
				} else if len(lBidDetailArr) == 3 {
					lAmount = math.Max(math.Max(lBidDetailArr[0].Amount, lBidDetailArr[1].Amount), lBidDetailArr[2].Amount)
				}
			}
		}
	}
	return lBidDetailArr, lAmount, nil
}

func getIpoDiscription(pActionFlag string) (string, error) {
	lDisciptionText := ""
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMD01", lErr1)
		return lDisciptionText, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `
					Select nvl(
						(select ld.description
						from xx_lookup_details ld,xx_lookup_header lh 
						where lh.id = ld.headerId and lh.Code = 'IpoAction' and ld.code = '` + pActionFlag + `'),''
					) ActionBtn`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("LSD02", lErr2)
			return lDisciptionText, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in respective structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lDisciptionText)
				if lErr3 != nil {
					log.Println("LSD03", lErr3)
					return lDisciptionText, lErr3
				}
			}
		}
	}
	return lDisciptionText, nil
}

func getIpoDisclaimer(pBrokerId int) ([]string, error) {
	lDisclaimerText := ""
	var lDisclaimerArray []string
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMD01", lErr1)
		return lDisclaimerArray, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `SELECT nvl(trim(bm.IpoDisclaimer),'') disclaimer
						FROM a_ipo_brokermaster bm 
						WHERE bm.id = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("LSD02", lErr2)
			return lDisclaimerArray, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in respective structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lDisclaimerText)
				if lErr3 != nil {
					log.Println("LSD03", lErr3)
					return lDisclaimerArray, lErr3
				} else {
					lDisclaimerArray = strings.Split(lDisclaimerText, ".")
				}
			}
		}
	}
	return lDisclaimerArray, nil
}

func getLastlyUsedUPI(pClientId string) (string, error) {
	lUPI := ""
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGLUU01", lErr1)
		return lUPI, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select h.upi from a_ipo_order_header h
							where h.clientId = ?
							order by id desc
							limit 1`
		lRows, lErr2 := lDb.Query(lCoreString, pClientId)
		if lErr2 != nil {
			log.Println("LGLUU02", lErr2)
			return lUPI, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in respective structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lUPI)
				if lErr3 != nil {
					log.Println("LGLUU03", lErr3)
					return lUPI, lErr3
				}
			}
		}
	}
	return lUPI, nil
}
