package localdetail

import (
	"database/sql"
	"encoding/json"
	"fcs23pkg/apps/validation/apiaccess"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// This Structure is used to get UPI records
type CategoryStruct struct {
	Value         string  `json:"value"`
	Text          string  `json:"text"`
	Maxvalue      float64 `json:"maxvalue"`
	DiscountPrice float64 `json:"discountPrice"`
	DiscountType  string  `json:"discountType"`
}

// Response Structure for GetUpi API
type CategoryRespStruct struct {
	CategoryArr []CategoryStruct `json:"categoryArr"`
	Status      string           `json:"status"`
	ErrMsg      string           `json:"errMsg"`
}

/*
Pupose: This method returns the Upi names from the  categoryArr data table
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    {
			categoryArr:[
				{
					"value": "IND",
					"text": "Individual Investor",
					"maxvalue":50000,
					"discountprice":75,
					"discountType":"P"
				},
			]
		}

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 05DEC2023
*/
func GetCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCategory (+)", r.Method)
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Println(origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET")
	(w).Header().Set("Access-Control-Allow-Headers", "ID,PATH,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

	if r.Method == "GET" {

		// create a instance of  upi struct
		// var lCategoryRec CategoryStruct
		// create a instance of upiResponse struct
		var lRespRec CategoryRespStruct

		lId := r.Header.Get("ID")
		lMasterId, lErr1 := strconv.Atoi(lId)
		if lErr1 != nil {
			log.Println("LGC01", "Can't convert ID to int")
		} else {

			lPath := r.Header.Get("PATH")
			lRespRec.Status = common.SuccessCode
			log.Println("masterid", lId)
			log.Println("path", lPath)

			//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
			// lSessionClientId := ""
			lClientId := ""
			var lErr2 error
			lClientId, lErr2 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
			if lErr2 != nil {
				log.Println("LGC02", lErr2)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "LGC02" + lErr2.Error()
				fmt.Fprintf(w, helpers.GetErrorString("LGC02", "UserDetails not Found"))
				return
			} else {
				if lClientId == "" {
					fmt.Fprintf(w, helpers.GetErrorString("LGC03", "UserDetails not Found"))
					return
				}
			}
			//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
			lcategoryArr, lErr4 := GetCategoryDiscount(lPath, lMasterId)
			if lErr4 != nil {
				log.Println("LGC04", lErr4)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "LGC04" + lErr4.Error()
				fmt.Fprintf(w, helpers.GetErrorString("LGC04", "Unable to get category details"))
				return
			} else {
				lRespRec.CategoryArr = lcategoryArr
			}
		}
		// Marshall the structure into json
		lData, lErr5 := json.Marshal(lRespRec)
		if lErr5 != nil {
			log.Println("LGC05", lErr5)
			fmt.Fprintf(w, helpers.GetErrorString("LGC05", "Unable to process your request now. Please try after somettime"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetCategory (-)", r.Method)
	}
}

func GetCategoryDiscount(pPath string, pMasterId int) ([]CategoryStruct, error) {
	log.Println("GetCategoryDiscount (+)")
	var lCategoryRec CategoryStruct
	var lCategoryArr []CategoryStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGDC01", lErr1)
		return lCategoryArr, lErr1
	} else {
		defer lDb.Close()
		ladditionalCond1 := ""
		ladditionalCond2 := ""
		lCoreString := `select distinct nvl(s.SubCatCode,'') code ,nvl(s.MaxUpiLimit,0) Max_Value ,
	(case when (s.DiscountType = "P" or s.DiscountType = 3) then "P" else "A" end) as discountType,
	s.DiscountPrice as discountPrice
					from a_ipo_subcategory s,a_ipo_master m
					where m.Id = s.MasterId `

		if pPath != "/report" {
			ladditionalCond1 = ` and s.MaxUpiLimit > 0 
							and m.Id = ?`
		}
		ladditionalCond2 = ` and s.SubCatCode in('IND','SHA','EMP')
							and s.CaCode in ('RETAIL','EMPRET','SHARET')`

		lCoreFinal := lCoreString + ladditionalCond1 + ladditionalCond2

		var lRows *sql.Rows
		var lErr2 error
		if pPath != "/report" {
			lRows, lErr2 = lDb.Query(lCoreFinal, pMasterId)
		} else {
			lRows, lErr2 = lDb.Query(lCoreFinal)
		}

		if lErr2 != nil {
			log.Println("LGDC02", lErr2)
			return lCategoryArr, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lCategoryRec.Value, &lCategoryRec.Maxvalue, &lCategoryRec.DiscountType, &lCategoryRec.DiscountPrice)
				if lErr3 != nil {
					log.Println("LGDC03", lErr3)
					return lCategoryArr, lErr3
				} else {
					switch lCategoryRec.Value {
					case "IND":
						lCategoryRec.Text = "Individual Investor"
					case "EMP":
						lCategoryRec.Text = "Employee"
					case "SHA":
						lCategoryRec.Text = "Existing Shareholder"
					default:
						lCategoryRec.Text = "Individual Investor"
					}
					// Append Upi End Point in lRespRec.UpiArr array
					lCategoryArr = append(lCategoryArr, lCategoryRec)
				}
			}
		}
	}

	log.Println("GetCategoryDiscount (-)")
	return lCategoryArr, nil
}
