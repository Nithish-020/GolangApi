package localdetail

import (
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

// This Structures is used to Collect the Active IPO market informations
type ipoMktDemandStruct struct {
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
	Cutoff   bool   `json:"cutoff"`
}
type ipoMktCatwiseStruct struct {
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
}

// Response Structure for GetIpoMaster API
type IpoMktDataStruct struct {
	IpoMktDemandArr  []ipoMktDemandStruct  `json:"ipoMktDemandArr"`
	IpoMktCatwiseArr []ipoMktCatwiseStruct `json:"ipoMktCatwiseArr"`
	Status           string                `json:"status"`
	ErrMsg           string                `json:"errMsg"`
}

/*
Pupose:This Function is used to Get the Active Ipo Details in our database table ....
Parameters:

not Applicable

Response:

*ON Sucess
=========

	{
		"ipoMktDemandArr": [
			{
			"price": "1250",
			"quantity": 50,
			"cutoff": true
			},
			{
			"price": "300",
			"quantity": 125,
			"cutoff": false
			},
			{
			"price": "250",
			"quantity": 550,
			"cutoff": false
			}
		],
		"ipoMktCatwiseArr": [
			{
			"category": "NIBAT",
			"quantity": 0
			},
			{
			"category": "NIBBT",
			"quantity": 525
			},
			{
			"category": "RETAIL",
			"quantity": 25
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
Date: 21DEC2023
*/
func GetIpoMKtData(w http.ResponseWriter, r *http.Request) {
	log.Println("GetIpoMKtData(+)", r.Method)
	origin := r.Header.Get("Origin")
	// var lErr error
	for _, allowedOrigin := range common.ABHIAllowOrigin {
		if allowedOrigin == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}

	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "ID,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	if r.Method == "GET" {

		// create the instance for IpoStruct
		var lRespRec IpoMktDataStruct

		lId := r.Header.Get("ID")
		lMasterId, _ := strconv.Atoi(lId)

		lRespRec.Status = common.SuccessCode

		//-----------START TO GETTING CLIENT AND STAFF DETAILS--------------
		lClientId := ""
		var lErr2 error
		lClientId, lErr2 = apiaccess.VerifyApiAccess(r, common.ABHIAppName, common.ABHICookieName, "/ipo")
		if lErr2 != nil {
			log.Println("LGIMD02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIMD02" + lErr2.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIMD02", "UserDetails not Found"))
			return
		} else {
			if lClientId == "" {
				fmt.Fprintf(w, helpers.GetErrorString("LGIMD02", "UserDetails not Found"))
				return
			}
		}
		//-----------END OF GETTING CLIENT AND STAFF DETAILS----------------
		log.Println("lMasterId := ", lMasterId)
		lIpoMktCatwise, lErr3 := GetIpoMKtCatwiseDetail(lMasterId)
		if lErr3 != nil {
			log.Println("LGIMD02", lErr3)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIMD03" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIMD03", "Unable to fetch the ipo category wise details"))
			return
		} else {
			lRespRec.IpoMktCatwiseArr = lIpoMktCatwise
		}

		lIpoMktDemand, lErr3 := GetIpoMKtDemandDetail(lMasterId)
		if lErr3 != nil {
			log.Println("LGIMD04", lErr3)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "LGIMD04" + lErr3.Error()
			fmt.Fprintf(w, helpers.GetErrorString("LGIMD04", "Unable to fetch the ipo demand details"))
			return
		} else {
			lRespRec.IpoMktDemandArr = lIpoMktDemand
		}
		// Marshaling the response structure to lData
		lData, lErr5 := json.Marshal(lRespRec)
		if lErr5 != nil {
			log.Println("LGIMD05", lErr5)
			fmt.Fprintf(w, helpers.GetErrorString("LGIMD05", "Issue in Getting Datas!"))
			return
		} else {
			fmt.Fprintf(w, string(lData))
		}
		log.Println("GetIpoMKtData (-)", r.Method)
	}
}

func GetIpoMKtCatwiseDetail(pMasterId int) ([]ipoMktCatwiseStruct, error) {
	log.Println("GetIpoMKtDetail (+)")
	// create the instance for activeIpoStruct
	var lIpoDataRec ipoMktCatwiseStruct
	var lIpoDataArr []ipoMktCatwiseStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMCD01", lErr1)
		return lIpoDataArr, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `SELECT 	c.Category ,sum(c.Quantity) 	
						FROM a_ipo_master m ,a_ipo_mktcatwise c
						WHERE c.Symbol = m.Symbol
							AND m.Isin = c.Isin
							AND m.id = ?
							AND m.BiddingEndDate >= Curdate() AND m.Exchange = 'NSE'
						GROUP BY c.Category`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
		if lErr2 != nil {
			log.Println("LGIMCD03", lErr2)
			return lIpoDataArr, lErr2

		} else {
			//This for loop is used to collect the records FROM the database AND store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lIpoDataRec.Category, &lIpoDataRec.Quantity)
				if lErr3 != nil {
					log.Println("LGIMCD03", lErr3)
					return lIpoDataArr, lErr3

				} else {
					// Append the IPO Records in lRespRec.IpoDetails Array
					lIpoDataArr = append(lIpoDataArr, lIpoDataRec)
				}
			}
		}
	}
	log.Println("GetIpoMKtDetail (-)")
	return lIpoDataArr, nil
}

func GetIpoMKtDemandDetail(pMasterId int) ([]ipoMktDemandStruct, error) {
	log.Println("GetIpoMKtDemandDetail (+)")
	// create the instance for activeIpoStruct
	var lIpoDataRec ipoMktDemandStruct
	var lIpoDataArr []ipoMktDemandStruct

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("LGIMDD01", lErr1)
		return lIpoDataArr, lErr1

	} else {
		defer lDb.Close()

		lCoreString := `SELECT	d.Price, SUM(d.cumulativeQuantity),max(d.CutoffIndicator)
						FROM a_ipo_master m ,a_ipo_mktdemand d
						WHERE m.Symbol = d.Symbol
							AND m.Isin = d.Isin
							AND m.id = ?
							AND m.BiddingEndDate >= Curdate() AND m.Exchange = 'NSE'
						GROUP BY d.Price 
						ORDER BY d.Price desc`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
		if lErr2 != nil {
			log.Println("LGIMDD02", lErr2)
			return lIpoDataArr, lErr2

		} else {
			//This for loop is used to collect the records FROM the database AND store them in structure
			for lRows.Next() {
				lErr3 := lRows.Scan(&lIpoDataRec.Price, &lIpoDataRec.Quantity, &lIpoDataRec.Cutoff)
				if lErr3 != nil {
					log.Println("LGIMDD03", lErr3)
					return lIpoDataArr, lErr3

				} else {
					// Append the IPO Records in lRespRec.IpoDetails Array
					lIpoDataArr = append(lIpoDataArr, lIpoDataRec)
				}
			}
		}
	}
	log.Println("GetIpoMKtDemandDetail (-)")
	return lIpoDataArr, nil
}
