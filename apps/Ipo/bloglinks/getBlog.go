package bloglinks

import (
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
)

// This Structure is used to get UPI records
type blogLinkStruct struct {
	Id             int    `json:"id"`
	Symbol         string `json:"symbol"`
	Isin           string `json:"isin"`
	BlogLink       string `json:"blogLink"`
	DRHPLink       string `json:"drhpLink"`
	Sme            bool   `json:"sme"`
	CreatedBy      string `json:"createdBy"`
	CreatedDate    string `json:"createdDate"`
	AllotmentFinal string `json:"allotmentFinal"`
	RefundInitiate string `json:"refundInitiate"`
	DematTransfer  string `json:"dematTransfer"`
	ListingDate    string `json:"listingDate"`
}

// Response Structure for GetUpi API
type blogLinkRespStruct struct {
	BlogLinkArr []blogLinkStruct `json:"blogLinkArr"`
	Status      string           `json:"status"`
	ErrMsg      string           `json:"errMsg"`
}

/*
Pupose: This method returns the Blog Links from the  a_ipo_details data table
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    {
			upiArr:[
				{
					"id": "1",
					"isin": "INE192R01011",
					"blogLink": "www.example.com",
					"dhrpLink": "www.Example.com",
				},
			]
		}

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 25JULY2023
*/
func GetBlogLink(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBlogLink (+)", r.Method)
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		// create a instance of  upi struct
		var lBlogRec blogLinkStruct
		// create a instance of upiResponse struct
		var lRespRec blogLinkRespStruct

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		// lSessionClientId := ""
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipoInfo")
		if lErr1 != nil {
			log.Println("LGBL01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGBL01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGBL01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGBL01", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		// Calling LocalDbConect method in ftdb to estabish the database connection
		lDb, lErr2 := ftdb.LocalDbConnect(ftdb.IPODB)
		if lErr2 != nil {
			log.Println("LGBL02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGBL02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGBL02", "Issue in Fetching your Details.Please try after sometime"))
			return
		} else {
			defer lDb.Close()

			var lCategory int

			lCoreString := `select aid.Id ,nvl(aid.Symbol,'') ,aid.Isin ,aid.detailsLink ,aid.drhpLink ,aid.category,nvl(aid.createdBy,''),nvl(date(aid.createdDate),''),nvl(aid.allotmentFinal ,'') ,nvl(aid.refundInitiate ,'') ,nvl(aid.dematTransfer ,'') ,nvl(aid.listingDate ,'')
			from a_ipo_details aid
			order by aid.createdDate desc 	 `
			lRows, lErr3 := lDb.Query(lCoreString)
			if lErr3 != nil {
				log.Println("LGBL03", lErr3)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "LGBL03" + lErr3.Error()
				fmt.Fprintf(w, helpers.GetErrorString("LGBL03", "Unable to process your request now. Please try after somettime"))
				return
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows.Next() {
					lErr4 := lRows.Scan(&lBlogRec.Id, &lBlogRec.Symbol, &lBlogRec.Isin, &lBlogRec.BlogLink, &lBlogRec.DRHPLink, &lCategory, &lBlogRec.CreatedBy, &lBlogRec.CreatedDate, &lBlogRec.AllotmentFinal, &lBlogRec.RefundInitiate, &lBlogRec.DematTransfer, &lBlogRec.ListingDate)
					if lErr4 != nil {
						log.Println("LGBL04", lErr4)
						lRespRec.Status = common.ErrorCode
						lRespRec.ErrMsg = "LGBL04" + lErr4.Error()
						fmt.Fprintf(w, helpers.GetErrorString("LGBL04", "Unable to get Blog Links. Please try after somettime"))
						return
					} else {
						if lCategory == 1 {
							lBlogRec.Sme = true
						} else {
							lBlogRec.Sme = false
						}
						lRespRec.BlogLinkArr = append(lRespRec.BlogLinkArr, lBlogRec)
					}
				}
			}
		}
		// Marshall the structure into json
		lData, lErr5 := json.Marshal(lRespRec)
		if lErr5 != nil {
			log.Println("LGBL05", lErr5)
			fmt.Fprintf(w, helpers.GetErrorString("LGBL05", "Server busy right now. Please try after somettime"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetBlogLink (-)", r.Method)
	}
}
