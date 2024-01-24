package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nsesgb"
	"log"
	"strings"
	"time"
)

/*
Pupose: This method used to fetch SGB master records from exchange
Parameters:

	pUser

Response:

	    ==========
	    *On Sucess
	    ==========
		nil
	    =========
	    !On Error
	    =========
	    error

Author: Pavithra
Date: 31 July 2023
*/
func FetchSgbMasterNSE(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchSgbMasterNSE (+)")

	lNoToken := common.ErrorCode
	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("SEFSMN01", lErr1)
		return lNoToken, lErr1
	} else {
		log.Println("lToken", lToken)
		if lToken != "" {
			lErr2 := getSgbMasterNSE(lToken, pUser)
			if lErr2 != nil {
				log.Println("SEFSMN02", lErr2)
				return lNoToken, lErr2
			}
			lNoToken = common.SuccessCode
		}
	}
	log.Println("FetchSgbMasterNSE (-)")
	return lNoToken, nil
}

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:

	pToken,pUser

Response:

	    ==========
	    *On Sucess
	    ==========
	    nil
	    =========
	    !On Error
	    =========
		error

Author: Pavithra
Date: 31 July 2023
*/
func getSgbMasterNSE(pToken string, pUser string) error {
	log.Println("getSgbMasterNSE (+)")
	//create instance for IpoResponseStruct
	var lSgbRespRec []nsesgb.SgbDetailStruct

	lSgbResp, lErr1 := nsesgb.SgbMaster(pToken, pUser)
	if lErr1 != nil {
		log.Println("SEGSMN01", lErr1)
		return lErr1
	} else {
		lSgbRespRec = lSgbResp.Data

		// log.Println("IPO response", lIpoResponse)
		lErr2 := InsertSgbDatas(lSgbRespRec, pUser, common.NSE)
		if lErr2 != nil {
			log.Println("SEGSMN02", lErr2)
			return lErr2
		} else {
			log.Println("Data inserted successfully")
		}
	}
	log.Println("getSgbMasterNSE (-)")
	return nil
}

/*
Pupose: This method helpsto insert new records in  a_ipo_master datatable
Parameters:

	pSgbResp,pUser

Response:

	    ==========
	    *On Sucess
	    ==========
		nil
	    =========
	    !On Error
	    =========
		error

Author: Pavithra
Date: 31 July 2023
*/
func InsertSgbDatas(pSgbResp []nsesgb.SgbDetailStruct, pUser string, pExchange string) error {
	log.Println("InsertSgbDatas (-)")

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SEISD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pSgbResp != nil {
			for lSgbRespidx := 0; lSgbRespidx < len(pSgbResp); lSgbRespidx++ {
				// to check the ipo is already exist or not exist
				lexistId, lErr2 := CheckDataExists(pSgbResp[lSgbRespidx], pExchange)
				if lErr2 != nil {
					log.Println("SEISD02", lErr2)
					return lErr2
				} else {
					if lexistId == 0 {
						lErr3 := InsertSgbRecord(pSgbResp[lSgbRespidx], pUser, pExchange)
						if lErr3 != nil {
							log.Println("SEISD03", lErr3)
							return lErr3
						} else {
							log.Println("SGB Fetch Master Inserted successfully", lSgbRespidx)
						}
					} else {
						lErr4 := UpdateSgbRecord(pSgbResp[lSgbRespidx], lexistId, pUser, pExchange)
						if lErr4 != nil {
							log.Println("SEISD04", lErr4)
							return lErr4
						} else {
							log.Println("SGB Fetch Master Updated successfully", lSgbRespidx)
						}
					}
				}
			}
		}
	}
	log.Println("InsertSgbDatas (-)")
	return nil
}

/*
Pupose: This method helps to check the ipo record is already exist or not
Parameters:

	pSgbRespRec , pIporespidx

Response:

	==========
	*On Sucess
	==========
	1,nil
	=========
	!On Error
	=========
	0,error

Author: Pavithra
Date: 01 AUG 2023
*/
func CheckDataExists(pSgbRespRec nsesgb.SgbDetailStruct, pExchange string) (int, error) {
	log.Println("CheckDataExists (+)")

	lId := 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SECDE01", lErr1)
		return lId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select m.id 
						from a_sgb_master m
						where m.Symbol = ? 
						and m.Isin = ?
						and m.Exchange = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pSgbRespRec.Symbol, pSgbRespRec.ISIN, pExchange)
		if lErr2 != nil {
			log.Println("SECDE02", lErr2)
			return lId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lId)
				if lErr3 != nil {
					log.Println("SECDE03", lErr3)
					return lId, lErr3
				}
			}
		}
	}
	log.Println(lId)
	log.Println("CheckDataExists (-)")
	return lId, nil
}

/*
Pupose: This method helps to check the ipo record is already exist or not
Parameters:

	(ipoResponseStruct , ipoResponseidx )

Response:

	    ==========
	    *On Sucess
	    ==========
	    In case of a successful execution of this method, you will get the masterId data
		from the a_ipo_master Data Table
	    =========
	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 01 AUG 2023
*/
func InsertSgbRecord(pSgbRespRec nsesgb.SgbDetailStruct, pUser string, pExchange string) error {
	log.Println("InsertSgbRecord (+)")
	var lRedemption string

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SEISR01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pExchange == "NSE" {
			pSgbRespRec = ChangeDateFormat(pSgbRespRec)
		}

		//commented below by Lakshmanan on 17 Dec 2023
		//there seem to be some bug here
		/*if !containsRedemption(strings.ToUpper(pSgbRespRec.Name)) {
			lRedemption = "N"
		} else {
			lRedemption = "Y"
		}
		log.Println("Insert Redemption", lRedemption, pSgbRespRec.Symbol)

		// Filter out records containing "REDEMPTION" in the name
		if !containsRedemption(strings.ToUpper(pSgbRespRec.Name)) {
			lRedemption = "Y"
		} else {
			lRedemption = "N"
		}*/

		//added below if else by Lakshmanan on 17 Dec 2023
		//this below condition is the replacement of above commented condition
		if strings.Contains(strings.ToUpper(pSgbRespRec.Name), "REDEMPTION") {
			lRedemption = "Y"
		} else {
			lRedemption = "N"
		}

		lCoreString := `insert into a_sgb_master (Symbol,Series,Name,IssueType,Lotsize,FaceValue,MinBidQuantity,MinPrice,MaxPrice,TickSize,
						BiddingStartDate,BiddingEndDate,DailyStartTime,DailyEndTime,T1ModStartDate,T1ModEndDate,T1ModStartTime,T1ModEndTime,
						Isin,IssueSize,IssueValueSize,MaxQuantity,AllotmentDate,IncompleteModEndDate,Exchange,Redemption,CreatedBy,CretaedDate,UpdatedBy,UpdatedDate)
						values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

		_, lErr2 := lDb.Exec(lCoreString, pSgbRespRec.Symbol, pSgbRespRec.Series, pSgbRespRec.Name, pSgbRespRec.IssueType, pSgbRespRec.LotSize, pSgbRespRec.FaceValue, pSgbRespRec.MinBidQuantity, pSgbRespRec.MinPrice, pSgbRespRec.MaxPrice, pSgbRespRec.TickSize, pSgbRespRec.BiddingStartDate, pSgbRespRec.BiddingEndDate, pSgbRespRec.DailyStartTime, pSgbRespRec.DailyEndTime, pSgbRespRec.T1ModStartDate, pSgbRespRec.T1ModEndDate, pSgbRespRec.T1ModStartTime, pSgbRespRec.T1ModEndTime, pSgbRespRec.ISIN, pSgbRespRec.IssueSize, pSgbRespRec.IssueValueSize, pSgbRespRec.MaxQuantity, pSgbRespRec.AllotmentDate, pSgbRespRec.IncompleteModEndDate, pExchange, lRedemption, pUser, pUser)
		if lErr2 != nil {
			log.Println("SEISR02", lErr2)
			return lErr2
		} else {
			log.Println("SGB Inserted..")

		}

	}
	log.Println("InsertSgbRecord (-)")
	return nil
}

// Function to check if the name contains "REDEMPTION"
func containsRedemption(name string) bool {
	return !containsSubstring(name, "REDEMPTION")
}

// Function to check substring existence in a string
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

func ChangeDateFormat(pSgbRespRec nsesgb.SgbDetailStruct) nsesgb.SgbDetailStruct {
	log.Println("ChangeDateFormat (-)")

	//for bidding start date
	if pSgbRespRec.BiddingStartDate != "" {
		lStartDate, _ := time.Parse("02-01-2006", pSgbRespRec.BiddingStartDate)
		pSgbRespRec.BiddingStartDate = lStartDate.Format("2006-01-02")
		// log.Println(pSgbRespRec.BiddingStartDate, "BiddingStartDate")
	} else {
		pSgbRespRec.BiddingStartDate = "0000-00-00"
	}
	//for bidding end date
	if pSgbRespRec.BiddingEndDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pSgbRespRec.BiddingEndDate)
		pSgbRespRec.BiddingEndDate = lEndDate.Format("2006-01-02")
		// log.Println(pSgbRespRec.BiddingEndDate, "BiddingEndDate")
	} else {
		pSgbRespRec.BiddingEndDate = "0000-00-00"
	}
	//for t1 bidding start date
	if pSgbRespRec.T1ModStartDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pSgbRespRec.T1ModStartDate)
		pSgbRespRec.T1ModStartDate = lEndDate.Format("2006-01-02")
		// log.Println(pSgbRespRec.T1ModStartDate, "T1ModStartDate")
	} else {
		pSgbRespRec.T1ModStartDate = "0000-00-00"
	}
	//for t1 bidding end date
	if pSgbRespRec.T1ModEndDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pSgbRespRec.T1ModEndDate)
		pSgbRespRec.T1ModEndDate = lEndDate.Format("2006-01-02")
		// log.Println(pSgbRespRec.T1ModEndDate, "T1ModEndDate")
	} else {
		pSgbRespRec.T1ModEndDate = "0000-00-00"
	}

	log.Println("ChangeDateFormat (-)")
	return pSgbRespRec
}

/*
Pupose: This method helps to check the ipo record is already exist or not
Parameters:

	(ipoResponseStruct , ipoResponseidx )

Response:

	    ==========
	    *On Sucess
	    ==========
	    In case of a successful execution of this method, you will get the masterId data
		from the a_ipo_master Data Table
	    =========
	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 01 AUG 2023
*/
func UpdateSgbRecord(pSgbRespRec nsesgb.SgbDetailStruct, pId int, pUser string, pExchange string) error {
	log.Println("UpdateSgbRecord (+)")
	var lRedemption string

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SEUSR01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pExchange == "NSE" {
			pSgbRespRec = ChangeDateFormat(pSgbRespRec)
		}
		//commented below by Lakshmanan on 17 Dec 2023
		//there seem to be some bug here
		// if !containsRedemption(strings.ToUpper(pSgbRespRec.Name)) {
		// 	lRedemption = "N"
		// } else {
		// 	lRedemption = "Y"
		// }

		//added below if else by Lakshmanan on 17 Dec 2023
		//this below condition is the replacement of above commented condition
		if strings.Contains(strings.ToUpper(pSgbRespRec.Name), "REDEMPTION") {
			lRedemption = "Y"
		} else {
			lRedemption = "N"
		}

		log.Println("Update Redemption", lRedemption, pSgbRespRec.Symbol)
		lCoreString := `update a_sgb_master m 
						set m.Symbol = ?,m.Series = ?,m.Name = ?,m.IssueType = ?,m.Lotsize = ?,m.FaceValue = ?,m.MinBidQuantity = ?,m.MinPrice = ?,
						m.MaxPrice = ?,m.TickSize = ?,m.BiddingStartDate = ?,m.BiddingEndDate = ?,m.DailyStartTime = ?,m.DailyEndTime = ?,
						m.T1ModStartDate = ?,m.T1ModEndDate = ?,m.T1ModStartTime = ?,m.T1ModEndTime  =?,m.Isin = ?,m.IssueSize = ?,m.IssueValueSize = ?,
						m.MaxQuantity = ?,m.AllotmentDate = ?,m.IncompleteModEndDate = ?,m.Redemption = ?,m.UpdatedBy = ?,m.UpdatedDate = now()
						where m.id = ? and m.Exchange = ?`

		_, lErr2 := lDb.Query(lCoreString, pSgbRespRec.Symbol, pSgbRespRec.Series, pSgbRespRec.Name, pSgbRespRec.IssueType, pSgbRespRec.LotSize, pSgbRespRec.FaceValue, pSgbRespRec.MinBidQuantity, pSgbRespRec.MinPrice, pSgbRespRec.MaxPrice, pSgbRespRec.TickSize, pSgbRespRec.BiddingStartDate, pSgbRespRec.BiddingEndDate, pSgbRespRec.DailyStartTime, pSgbRespRec.DailyEndTime, pSgbRespRec.T1ModStartDate, pSgbRespRec.T1ModEndDate, pSgbRespRec.T1ModStartTime, pSgbRespRec.T1ModEndTime, pSgbRespRec.ISIN, pSgbRespRec.IssueSize, pSgbRespRec.IssueValueSize, pSgbRespRec.MaxQuantity, pSgbRespRec.AllotmentDate, pSgbRespRec.IncompleteModEndDate, lRedemption, pUser, pId, pExchange)
		if lErr2 != nil {
			log.Println("SEUSR02", lErr2)
			return lErr2
		} else {
			log.Println("SGB Updated..")
		}
	}
	log.Println("UpdateSgbRecord (-)")
	return nil
}

//--------------------->>>>>>>>>>>>>>>>>>BSE SGB MASTER

/*
Pupose: This method used to fetch SGB master records from exchange
Parameters:
	pUser
Response:
	    ==========
	    *On Sucess
	    ==========
		nil
	    =========
	    !On Error
	    =========
	    error

Author: Pavithra
Date: 31 July 2023
*/

// func FetchSgbMasterBSE(pUser string) error {
// 	log.Println("FetchSgbMasterBSE (+)")

// 	lToken, lErr1 := BseGetToken(pUser)
// 	if lErr1 != nil {
// 		log.Println("SEFSMB01", lErr1)
// 		return lErr1
// 	} else {
// 		log.Println("lToken", lToken)
// 		if lToken != "" {
// 			lErr2 := getSgbMasterBSE(lToken, pUser)
// 			if lErr2 != nil {
// 				log.Println("SEFSMB02", lErr2)
// 				return lErr2
// 			}
// 		}
// 	}
// 	log.Println("FetchSgbMasterBSE (-)")
// 	return nil
// }

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:
	pToken,pUser
Response:
	    ==========
	    *On Sucess
	    ==========
	    nil
	    =========
	    !On Error
	    =========
		error

Author: Pavithra
Date: 16 Aug 2023
*/

// func getSgbMasterBSE(pToken string, pUser string) error {
// 	log.Println("getSgbMasterBSE (+)")
// 	//create instance for IpoResponseStruct
// 	// var lSgbRespRec []nsesgb.SgbDetailStruct

// 	lSgbResp, lErr1 := bsesgb.BseSgbMaster(pToken, pUser)
// 	if lErr1 != nil {
// 		log.Println("SEGSMB01", lErr1)
// 		return lErr1
// 	} else {

// 		log.Println("SGB Response lSgbRespRec", lSgbResp)
// 		if lSgbResp != nil {
// 			// construct bse struct to nse struct
// 			lSgbRespRec := ConstructSGB(lSgbResp)
// 			// inseerting the details into database
// 			lErr2 := InsertSgbDatas(lSgbRespRec, pUser, common.BSE)
// 			if lErr2 != nil {
// 				log.Println("SEGSMB02", lErr2)
// 				return lErr2
// 			} else {
// 				log.Println("Data inserted")
// 			}
// 		}
// 	}
// 	log.Println("getSgbMasterBSE (-)")
// 	return nil
// }

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:
	pToken,pUser
Response:
	    ==========
	    *On Sucess
	    ==========
	    nil
	    =========
	    !On Error
	    =========
		error

Author: Pavithra
Date: 16 aug 2023
*/

// func ConstructSGB(pSgbArr []bsesgb.BseSgbMasterStruct) []nsesgb.SgbDetailStruct {
// 	log.Println("ConstructSGB (+)")
// 	var lNseSgbRec nsesgb.SgbDetailStruct
// 	var lNseSgbArr []nsesgb.SgbDetailStruct

// 	for lSgbIdx := 0; lSgbIdx < len(pSgbArr); lSgbIdx++ {
// 		lNseSgbRec.ISIN = pSgbArr[lSgbIdx].Isin
// 		lNseSgbRec.Name = pSgbArr[lSgbIdx].IsssueName
// 		lNseSgbRec.Symbol = pSgbArr[lSgbIdx].Symbol
// 		lNseSgbRec.Series = pSgbArr[lSgbIdx].SerialNo

// 		// bidding startdate
// 		if pSgbArr[lSgbIdx].OpenDate != "" {
// 			lStartDate, lErr := time.Parse("1/2/2006 3:04:05 PM", pSgbArr[lSgbIdx].OpenDate)
// 			if lErr != nil {
// 				log.Println("Error ", lErr)
// 			}
// 			lNseSgbRec.BiddingStartDate = lStartDate.Format("2006-01-02")
// 			lNseSgbRec.DailyStartTime = lStartDate.Format("15:04:05")
// 			lNseSgbRec.T1ModStartDate = lStartDate.Format("2006-01-02")
// 			lNseSgbRec.T1ModStartTime = lStartDate.Format("15:04:05")
// 			log.Println(pSgbArr[lSgbIdx].OpenDate, "pSgbArr[lSgbIdx].OpenDate")
// 			log.Println(lNseSgbRec.BiddingStartDate, "BiddingStartDate")
// 			log.Println(lNseSgbRec.DailyStartTime, "DailyStartTime")
// 			log.Println(lNseSgbRec.T1ModStartDate, "T1ModStartDate")
// 			log.Println(lNseSgbRec.T1ModStartTime, "T1ModStartTime")
// 		}

// 		// Bidding EndDate
// 		if pSgbArr[lSgbIdx].CloseDate != "" {
// 			lEndDate, lErr := time.Parse("1/2/2006 3:04:05 PM", pSgbArr[lSgbIdx].CloseDate)
// 			if lErr != nil {
// 				log.Println("Error ", lErr)
// 			}
// 			lNseSgbRec.BiddingEndDate = lEndDate.Format("2006-01-02")
// 			lNseSgbRec.T1ModEndDate = lEndDate.Format("2006-01-02")
// 			lNseSgbRec.DailyEndTime = lEndDate.Format("15:04:05")
// 			lNseSgbRec.T1ModEndTime = lEndDate.Format("15:04:05")
// 			log.Println(lNseSgbRec.T1ModEndTime, "T1ModEndTime")
// 			log.Println(lNseSgbRec.DailyEndTime, "DailyEndTime")
// 			log.Println(lNseSgbRec.T1ModEndDate, "T1ModEndDate")
// 			log.Println(lNseSgbRec.BiddingEndDate, "BiddingEndDate")

// 		}

// 		//allotment date
// 		if pSgbArr[lSgbIdx].DateOfAllotment != "" {
// 			lEndDate, lErr := time.Parse("1/2/2006 3:04:05 PM", pSgbArr[lSgbIdx].DateOfAllotment)
// 			if lErr != nil {
// 				log.Println("Error ", lErr)
// 			}
// 			lNseSgbRec.AllotmentDate = lEndDate.Format("2006-01-02 15:04:05")
// 			log.Println(lNseSgbRec.AllotmentDate, "AllotmentDate")
// 		}

// 		//min price and max price
// 		lFloatValue, _ := strconv.ParseFloat(pSgbArr[lSgbIdx].IssuePrice, 32)
// 		lPrice := float32(lFloatValue) // Convert float64 to float32
// 		lNseSgbRec.MinPrice = lPrice
// 		lNseSgbRec.MaxPrice = lPrice

// 		lNseSgbRec.MinBidQuantity, _ = strconv.Atoi(pSgbArr[lSgbIdx].MinQty)
// 		lNseSgbRec.MaxQuantity, _ = strconv.Atoi(pSgbArr[lSgbIdx].MaxQty)

// 		//array of nse sgb struct
// 		lNseSgbArr = append(lNseSgbArr, lNseSgbRec)
// 	}

// 	log.Println("ConstructSGB (-)")
// 	return lNseSgbArr
// }

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:

	pToken,pUser

Response:

	    ==========
	    *On Sucess
	    ==========
	    nil
	    =========
	    !On Error
	    =========
		error

Author: Pavithra
Date: 16 aug 2023
*/
func FetchSGBMaster(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchSGBMaster (+)")

	lNoToken, lErr1 := FetchSgbMasterNSE(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("SEFSM01", lErr1)
		return lNoToken, lErr1
		// } else {
		// lErr2 := FetchSgbMasterBSE(pUser)
		// if lErr2 != nil {
		// 	log.Println("SEFSM02", lErr2)
		// 	return lErr2
	}
	// }
	log.Println("FetchSGBMaster (-)")
	return lNoToken, nil
}

// func fetchSgbPlaced() ([]JvReqStruct, error) {
// 	var lReqJVData JvReqStruct
// 	var lReqJVDataArr []JvReqStruct
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("SSFP01", lErr1)
// 		return lReqJVDataArr, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lSqlString := `select h.ApplicantName ,d.CreatedDate,
// 		h.ClientEmail ,m.Symbol,d.JvStatus,d.JvAmount,
// 			h.ClientId,d.JvType,d.ReqSubscriptionUnit,d.ReqRate,d.ActionCode
// 			from a_sgb_orderdetails d ,a_sgb_orderheader h ,a_sgb_master m
// 			where h.MasterId = m.id
// 			and d.HeaderId = h.Id
// 			and h.Status = "pending"
// 			and m.BiddingStartDate <= curdate()
// 			and m.BiddingEndDate >= curdate()
// 			and time(now()) between m.DailyStartTime and m.DailyEndTime
// 			and h.cancelFlag != 'Y' and m.Exchange = "BSE"`

// 		lRows, lErr2 := lDb.Query(lSqlString)
// 		if lErr2 != nil {
// 			log.Println("SSFP02", lErr2)
// 			return lReqJVDataArr, lErr2
// 		} else {

// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lReqJVData.ClientName, &lReqJVData.OrderDate, &lReqJVData.Mail, &lReqJVData.Symbol, &lReqJVData.JvStatus, &lReqJVData.JvAmount, &lReqJVData.ClientId, &lReqJVData.JvType, &lReqJVData.Unit, &lReqJVData.Price,
// 					&lReqJVData.ActionCode)
// 				if lErr3 != nil {
// 					log.Println("SSFP03", lErr3)
// 					return lReqJVDataArr, lErr3
// 				} else {
// 					lReqJVDataArr = append(lReqJVDataArr, lReqJVData)
// 				}
// 			}

// 		}

// 	}
// 	log.Println("lReqJVDataArr", lReqJVDataArr)
// 	log.Println("fetchPendingSgbOrder (-)")
// 	return lReqJVDataArr, lErr1
// }
