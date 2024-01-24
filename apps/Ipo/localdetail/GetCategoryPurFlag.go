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
	"net/http"
	"strconv"
)

type OrderedCategory struct {
	Category      string `json:"text"`
	Code          string `json:"code"`
	Flag          string `json:"flag"`
	ApplicationNo string `json:"applicationNo"`
}
type OrderCatgyResp struct {
	OrderedCategory []OrderedCategory `json:"orderedCategory"`
	Status          string            `json:"status"`
	ErrMsg          string            `json:"errMsg"`
}

/*
Pupose: This method returns the  names from the  upi_endnames data table
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    {
			upiArr:[
				{
					"Category": "Individual",
					"name": "@okaxis"
				},
			]
		}

	    !On Error
	    ========
	    In CASE of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Nithish Kumar
Date: 05JUNE2023
*/
func GetCategoryPurFlag(w http.ResponseWriter, r *http.Request) {
	lBrokerId := 0
	origin := r.Header.Get("Origin")
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			lBrokerId, _ = brokers.GetBrokerId(origin)
			log.Println(origin)
			break
		}
	}
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "ID,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {
		log.Println("GetCategoryPurFlag (+)", r.Method)
		var lRespRec OrderCatgyResp

		lId := r.Header.Get("ID")
		lMasterId, _ := strconv.Atoi(lId)
		lRespRec.Status = common.SuccessCode
		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		// lSessionClientId := ""
		lClientId := ""
		var lErr1 error
		lClientId, lErr1 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr1 != nil {
			log.Println("LGCPF01", lErr1)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGCPF01" + lErr1.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGCPF01", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGCPF02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------

		lOrderFlagArr, lErr3 := getCategoryPlacedDetails(lBrokerId, lMasterId, lClientId)
		if lErr3 != nil {
			log.Println("LGCPF03", lErr3.Error())
			fmt.Fprintf(w, helpers.GetErrorString("LGCPF03", "Error while getting category placed details"))
			return
		} else {
			lRespRec.OrderedCategory = lOrderFlagArr
		}

		// Marshall the structure into json
		lData, lErr4 := json.Marshal(lRespRec)
		if lErr4 != nil {
			log.Println("LGCPF04", lErr4)
			fmt.Fprintf(w, helpers.GetErrorString("LGCPF04", "Unable to process your request now. Please try after somettime"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetCategoryPurFlag (-)")
	}
}

func getCategoryPlacedDetails(pBrokerId int, pMasterId int, pClientId string) ([]OrderedCategory, error) {
	log.Println("GetCategoryPlacedDetails (+)")
	var lCategoryRec OrderedCategory
	var lOrderedCategory []OrderedCategory

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGCPD01", lErr1.Error())
		return lOrderedCategory, lErr1
	} else {
		defer lDb.Close()
		// IND EMP SHA mandatory To display in Catogary Dialog using V-for so I created an Left Join
		// COALESCE is the updated version of NVL
		lSqlString := `
					SELECT  
					CASE WHEN Tab1.code = "IND" then "Individual Investor"
					WHEN Tab1.code = "EMP" then "Employee"
					ELSE "Existing Shareholder" END AS category,
					Tab1.code Code,
					COALESCE(Tab2.flag, 'N') AS flag,
					COALESCE(ApplicationNo, '') AS applicationNo
					FROM
					(
						SELECT DISTINCT
							COALESCE(s.SubCatCode, '') AS code
						FROM
							a_ipo_subcategory s
							JOIN a_ipo_master m ON m.Id = s.MasterId
						WHERE
							m.Id = ?
							AND s.MaxUpiLimit > 0
							AND s.SubCatCode IN ('IND', 'SHA', 'EMP')
							AND s.CaCode IN ('RETAIL', 'EMPRET', 'SHARET')
					) Tab1
					LEFT JOIN
					(
						SELECT
							COALESCE(category, '') AS category,
							COALESCE(Flag, 'N') AS flag,
							COALESCE(ApplicationNo, 'N') AS applicationNo
						FROM
							(
								select oh.category,
									CASE
										WHEN oh.cancelFlag = 'N' AND oh.status = 'success' or oh.status = 'pending' THEN 'Y'
										ELSE 'N'
									END AS Flag,
									ROW_NUMBER() OVER (PARTITION BY oh.category ORDER BY oh.Id DESC) AS RowNum,
									CASE oh.category
										WHEN 'IND' THEN 1
										WHEN 'EMP' THEN 2
										WHEN 'SHA' THEN 3
										ELSE 4
									END AS CategoryOrder,
									oh.ApplicationNo AS ApplicationNo
								FROM
									a_ipo_order_header oh
									JOIN a_ipo_master M ON M.Id = oh.MasterId
								WHERE
									oh.brokerId = ?
									AND M.Id = ?
									AND oh.clientId = ?
									AND oh.category IN ('IND', 'EMP', 'SHA')
							) AS OrderedRows
						WHERE RowNum = 1
						ORDER BY CategoryOrder
					) Tab2
					ON Tab1.code = Tab2.category;`
		lRows, lErr2 := lDb.Query(lSqlString, pMasterId, pBrokerId, pMasterId, pClientId)
		if lErr2 != nil {
			log.Println("LGCPD02", lErr2)
			return lOrderedCategory, lErr2

		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lCategoryRec.Category, &lCategoryRec.Code, &lCategoryRec.Flag, &lCategoryRec.ApplicationNo)
				if lErr3 != nil {
					log.Println("LGCPD03", lErr3)
					return lOrderedCategory, lErr3
				} else {
					// Append the IPO Records in lRespRec.IpoDetails Array
					lOrderedCategory = append(lOrderedCategory, lCategoryRec)
				}
			}

		}
	}
	log.Println("GetCategoryPlacedDetails (-)")
	return lOrderedCategory, nil
}
