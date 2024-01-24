package exchangecall

import (
	"database/sql"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nsencb"
	"log"
	"time"
)

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

Author: Kavya Dharshani M
Date: 04 OCT 2023
*/
func FetchNCBMaster(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchNCBMaster (+)")

	lString := common.ErrorCode
	lTokenNse, lErr1 := FetchNcbMasterNSE(pUser, pBrokerId)
	if lErr1 != nil {
		// log.Println("Error fetching NCB master records:", lErr1)
		log.Println("NFNM01", lErr1)
		return lString, lErr1
	} else {
		if lTokenNse != common.ErrorCode {
			lString = common.SuccessCode
		}
		// log.Println("lTokenNse", lTokenNse, lString)
		// log.Println("NCB master records fetched and processed successfully!")
	}

	log.Println("FetchNCBMaster (-)")
	return lString, nil
}

/*
Pupose: This method used to fetch NCB master records from exchange
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

Author: Kavya Dharshani
Date: 04 OCT 2023
*/
func FetchNcbMasterNSE(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchNcbMasterNSE (+)")

	lNoToken := common.ErrorCode
	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("FNM01", lErr1)
		return lNoToken, lErr1
	} else {
		// log.Println("lToken", lToken)
		if lToken != "" {
			lErr2 := getNcbMaster(lToken, pUser)
			if lErr2 != nil {
				log.Println("FNM02", lErr2)
				return lNoToken, lErr2
			} else {
				lNoToken = common.SuccessCode
			}
			// log.Println("lNoToken", lNoToken)
		}
	}
	log.Println("FetchNcbMasterNSE(-)")
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

Author: Kavya Dharshani M
Date: 03 Oct 2023
*/
func getNcbMaster(pToken string, pUser string) error {
	log.Println("getNcbMaster (+)")
	// log.Println("pToken: getNcbMaster", pToken)

	//create instance for NcbResponseStruct
	var lNcbRespRec []nsencb.NcbDetailStruct
	lNcbResp, lErr1 := nsencb.NcbMaster(pToken, pUser)
	if lErr1 != nil {
		log.Println("SGNM01", lErr1)
		return lErr1
	} else {
		// log.Println("pToken", pToken)
		lNcbRespRec = lNcbResp.Data

		lErr2 := InsertNcbData(lNcbRespRec, pUser, common.NSE)
		if lErr2 != nil {
			log.Println("SGNM02", lErr2)
			return lErr2
		} else {
			log.Println("Data inserted")
		}
	}
	log.Println("getNcbMaster (-)")
	return nil
}

/*
Pupose: This method helpsto insert new records in a_ncb_master  datatable
Parameters:
	pNcbResp,pUser
Response:
	    ==========
	    *On Sucess
	    ==========
		nil
	    =========
	    !On Error
	    =========
		error

Author: Kavya Dharshani
Date: 03 Oct 2023
*/
func InsertNcbData(pNcbResp []nsencb.NcbDetailStruct, pUser string, pExchange string) error {
	log.Println("InsertNcbData (+)")

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NMIND01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pNcbResp != nil {
			for Idx := 0; Idx < len(pNcbResp); Idx++ {

				lexistId, lErr2 := CheckNcbDataExists(pNcbResp[Idx], pExchange)
				if lErr2 != nil {
					log.Println("NMIND02", lErr2)
					return lErr2
				} else {
					if lexistId == 0 {
						lErr3 := InsertNcbRecord(pNcbResp[Idx], pUser, pExchange)
						if lErr3 != nil {
							log.Println("NMIND03", lErr3)
							return lErr3
						} else {
							log.Println("NCB Fetch Master Inserted successfully", Idx)
						}
					} else {
						lErr4 := UpdateNcbRecord(pNcbResp[Idx], lexistId, pUser, pExchange)
						if lErr4 != nil {
							log.Println("NMIND04", lErr4)
							return lErr4
						} else {
							log.Println("NCB Fetch Master Updated successfully", Idx)
						}
					}
				}
			}
		}
	}
	log.Println("InsertNcbData (-)")
	return nil
}

/*
Pupose: This method helps to check the NCB record is already exist or not
Parameters:
	pNcbRespRec , pNcbrespidx
Response:
	    ==========
	    *On Sucess
	    ==========
	    1,nil
	    =========
	    !On Error
	    =========
	    0,error

Author: Kavya Dharshani
Date: 03 Oct 2023
*/
func CheckNcbDataExists(pNcbRespRec nsencb.NcbDetailStruct, pExchange string) (int, error) {
	log.Println("CheckNcbDataExists(+)")

	lId := 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NCNDE01", lErr1)
		return lId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select n.id  from a_ncb_master n 
		                where n.Symbol  = ?
						and n.Series =?
		                and n.Isin  = ?
		                and n.Exchange  = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pNcbRespRec.Symbol, pNcbRespRec.Series, pNcbRespRec.Isin, pExchange)
		if lErr2 != nil {
			log.Println("NCNDE02", lErr2)
			return lId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lId)
				if lErr3 != nil {
					log.Println("NCNDE03", lErr3)
					return lId, lErr3
				}
			}
		}
	}
	// log.Println(lId, "lId")
	log.Println("CheckNcbDataExists(-)")
	return lId, nil
}

/*
Pupose: This method helps to check the NCB record is already exist or not
Parameters:

	(ncbResponseStruct , ncbResponseidx )

Response:
	    ==========
	    *On Sucess
	    ==========
	    In case of a successful execution of this method, you will get the masterId data
		from the a_ncb_master Data Table
	    =========
	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Kavya Dharshani
Date: 03 Oct 2023
*/

func InsertNcbRecord(pNcbRespRec nsencb.NcbDetailStruct, pUser string, pExchange string) error {
	log.Println("InsertNcbRecord (+)")

	var lTime sql.NullString

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NINR01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		if pExchange == "NSE" {
			pNcbRespRec = ChangeDateFormats(pNcbRespRec)
		}

		// log.Println("Time VALUE11111", pNcbRespRec.LastDayBiddingEndTime)

		if pNcbRespRec.LastDayBiddingEndTime == "" {
			pNcbRespRec.LastDayBiddingEndTime = "1970-01-01 00:00:00"

			// log.Println("Using default formattedTime", pNcbRespRec.LastDayBiddingEndTime)

			lCoreString := `insert into a_ncb_master(Symbol,Series,Name,Lotsize,FaceValue,MinBidQuantity,MinPrice,MaxPrice,TickSize,CutoffPrice,
			BiddingStartDate,BiddingEndDate,DailyStartTime,DailyEndTime,T1ModStartDate,T1ModEndDate,T1ModStartTime,T1ModEndTime,
		   Isin,IssueSize,IssueValueSize,MaxQuantity,AllotmentDate,LastDayBiddingEndTime,Exchange,CreatedBy,CretaedDate,UpdatedBy,UpdatedDate)
		   values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`
			_, lErr2 := lDb.Exec(lCoreString, pNcbRespRec.Symbol, pNcbRespRec.Series, pNcbRespRec.Name, pNcbRespRec.LotSize, pNcbRespRec.FaceValue, pNcbRespRec.MinBidQuantity, pNcbRespRec.MinPrice, pNcbRespRec.MaxPrice, pNcbRespRec.TickSize, pNcbRespRec.CutoffPrice, pNcbRespRec.BiddingStartDate, pNcbRespRec.BiddingEndDate, pNcbRespRec.DailyStartTime, pNcbRespRec.DailyEndTime, pNcbRespRec.T1ModStartDate, pNcbRespRec.T1ModEndDate, pNcbRespRec.T1ModStartTime, pNcbRespRec.T1ModEndTime, pNcbRespRec.Isin, pNcbRespRec.IssueSize, pNcbRespRec.IssueValueSize, pNcbRespRec.MaxQuantity, pNcbRespRec.AllotmentDate, lTime, pExchange, pUser, pUser)

			if lErr2 != nil {
				log.Println("NINR02", lErr2)
				return lErr2
			} else {
				log.Println("NCB Inserted..")

			}

		} else {

			// log.Println("Time VALUE2222", pNcbRespRec.LastDayBiddingEndTime)

			parsedTime, err := time.Parse("2006-01-02T15:04:05-0700", pNcbRespRec.LastDayBiddingEndTime)
			if err != nil {
				// log.Println("Error parsing LastDayBiddingEndTime:", err)
				return err
			}

			pNcbRespRec.LastDayBiddingEndTime = parsedTime.Format("2006-01-02 15:04:05")
			// log.Println("Formatted LastDayBiddingEndTime:", pNcbRespRec.LastDayBiddingEndTime)

			lCoreString := `insert into a_ncb_master(Symbol,Series,Name,Lotsize,FaceValue,MinBidQuantity,MinPrice,MaxPrice,TickSize,CutoffPrice,BiddingStartDate,BiddingEndDate,DailyStartTime,DailyEndTime,T1ModStartDate,T1ModEndDate,T1ModStartTime,T1ModEndTime,
		   Isin,IssueSize,IssueValueSize,MaxQuantity,AllotmentDate,LastDayBiddingEndTime,Exchange,CreatedBy,CretaedDate,UpdatedBy,UpdatedDate)
		   values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`
			_, lErr2 := lDb.Exec(lCoreString, pNcbRespRec.Symbol, pNcbRespRec.Series, pNcbRespRec.Name, pNcbRespRec.LotSize, pNcbRespRec.FaceValue, pNcbRespRec.MinBidQuantity, pNcbRespRec.MinPrice, pNcbRespRec.MaxPrice, pNcbRespRec.TickSize, pNcbRespRec.CutoffPrice, pNcbRespRec.BiddingStartDate, pNcbRespRec.BiddingEndDate, pNcbRespRec.DailyStartTime, pNcbRespRec.DailyEndTime, pNcbRespRec.T1ModStartDate, pNcbRespRec.T1ModEndDate, pNcbRespRec.T1ModStartTime, pNcbRespRec.T1ModEndTime, pNcbRespRec.Isin, pNcbRespRec.IssueSize, pNcbRespRec.IssueValueSize, pNcbRespRec.MaxQuantity, pNcbRespRec.AllotmentDate, pNcbRespRec.LastDayBiddingEndTime, pExchange, pUser, pUser)

			if lErr2 != nil {
				log.Println("NINR02", lErr2)
				return lErr2
			} else {
				log.Println("NCB Inserted..")

			}
		}

	}
	log.Println("InsertNcbRecord (-)")
	return nil
}

/*
Pupose: This method helps to check the Ncb record is already exist or not
Parameters:

	(ncbResponseStruct , ncbResponseidx )

Response:
	    ==========
	    *On Sucess
	    ==========
	    In case of a successful execution of this method, you will get the masterId data
		from the a_ncb_master Data Table
	    =========
	    !On Error
	    =========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 01 AUG 2023
*/
func UpdateNcbRecord(pNcbRespRec nsencb.NcbDetailStruct, pId int, pUser string, pExchange string) error {
	log.Println("UpdateNcbRecord (+)")

	var lTime sql.NullString
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NUNR01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pExchange == "NSE" {
			pNcbRespRec = ChangeDateFormats(pNcbRespRec)
		}
		if pNcbRespRec.LastDayBiddingEndTime == "" {
			pNcbRespRec.LastDayBiddingEndTime = "1970-01-01 00:00:00"

			// log.Println("Using default formattedTime", pNcbRespRec.LastDayBiddingEndTime)

			lCoreString := `update a_ncb_master n 
		                set n.Symbol = ?, n.Series = ?,n.Name = ?,n.Lotsize = ?, n.FaceValue = ?,n.MinBidQuantity = ?,n.MinPrice = ?,
		               n.MaxPrice = ?,n.TickSize = ?,n.BiddingStartDate = ?,n.BiddingEndDate = ?,n.DailyStartTime = ?,n.DailyEndTime = ?,
		               n.T1ModStartDate = ?,n.T1ModEndDate = ?,n.T1ModStartTime = ?,n.T1ModEndTime  =?,n.Isin = ?,n.IssueSize = ?,n.IssueValueSize = ?,n.MaxQuantity = ?,n.AllotmentDate = ?,n.LastDayBiddingEndTime = ?,n.UpdatedBy = ?,
					   n.UpdatedDate = now()	
					    where n.id = ? and n.Exchange = ?`

			_, lErr2 := lDb.Query(lCoreString, pNcbRespRec.Symbol, pNcbRespRec.Series, pNcbRespRec.Name, pNcbRespRec.LotSize, pNcbRespRec.FaceValue, pNcbRespRec.MinBidQuantity, pNcbRespRec.MinPrice, pNcbRespRec.MaxPrice, pNcbRespRec.TickSize, pNcbRespRec.BiddingStartDate, pNcbRespRec.BiddingEndDate, pNcbRespRec.DailyStartTime, pNcbRespRec.DailyEndTime, pNcbRespRec.T1ModStartDate, pNcbRespRec.T1ModEndDate, pNcbRespRec.T1ModStartTime, pNcbRespRec.T1ModEndTime, pNcbRespRec.Isin, pNcbRespRec.IssueSize, pNcbRespRec.IssueValueSize, pNcbRespRec.MaxQuantity, pNcbRespRec.AllotmentDate, lTime, pUser, pId, pExchange)
			if lErr2 != nil {
				log.Println("NUNR02", lErr2)
				return lErr2
			} else {
				log.Println("NCB Updated..")
			}
		} else {
			parsedTime, err := time.Parse("2006-01-02T15:04:05-0700", pNcbRespRec.LastDayBiddingEndTime)
			if err != nil {
				// log.Println("Error parsing LastDayBiddingEndTime:", err)
				return err
			}

			pNcbRespRec.LastDayBiddingEndTime = parsedTime.Format("2006-01-02 15:04:05")
			// log.Println("Formatted LastDayBiddingEndTime:", pNcbRespRec.LastDayBiddingEndTime)
			lCoreString := `update a_ncb_master n 
		                set n.Symbol = ?, n.Series = ?,n.Name = ?,n.Lotsize = ?, n.FaceValue = ?,n.MinBidQuantity = ?,n.MinPrice = ?,
		               n.MaxPrice = ?,n.TickSize = ?,n.BiddingStartDate = ?,n.BiddingEndDate = ?,n.DailyStartTime = ?,n.DailyEndTime = ?,
		               n.T1ModStartDate = ?,n.T1ModEndDate = ?,n.T1ModStartTime = ?,n.T1ModEndTime  =?,n.Isin = ?,n.IssueSize = ?,n.IssueValueSize = ?,n.MaxQuantity = ?,n.AllotmentDate = ?,n.LastDayBiddingEndTime = ?,n.UpdatedBy = ?,
					   n.UpdatedDate = now()	
					    where n.id = ? and n.Exchange = ?`

			_, lErr2 := lDb.Query(lCoreString, pNcbRespRec.Symbol, pNcbRespRec.Series, pNcbRespRec.Name, pNcbRespRec.LotSize, pNcbRespRec.FaceValue, pNcbRespRec.MinBidQuantity, pNcbRespRec.MinPrice, pNcbRespRec.MaxPrice, pNcbRespRec.TickSize, pNcbRespRec.BiddingStartDate, pNcbRespRec.BiddingEndDate, pNcbRespRec.DailyStartTime, pNcbRespRec.DailyEndTime, pNcbRespRec.T1ModStartDate, pNcbRespRec.T1ModEndDate, pNcbRespRec.T1ModStartTime, pNcbRespRec.T1ModEndTime, pNcbRespRec.Isin, pNcbRespRec.IssueSize, pNcbRespRec.IssueValueSize, pNcbRespRec.MaxQuantity, pNcbRespRec.AllotmentDate, pNcbRespRec.LastDayBiddingEndTime, pUser, pId, pExchange)
			if lErr2 != nil {
				log.Println("NUNR02", lErr2)
				return lErr2
			} else {
				log.Println("NCB Updated..")
			}
		}
	}
	log.Println("UpdateNcbRecord (-)")
	return nil
}

func ChangeDateFormats(pNcbRespRec nsencb.NcbDetailStruct) nsencb.NcbDetailStruct {
	log.Println("ChangeDateFormat (+)")

	var lString string
	var lErr error

	//for bidding start date
	if pNcbRespRec.BiddingStartDate != "" {
		lStartDate, _ := time.Parse("02-01-2006", pNcbRespRec.BiddingStartDate)
		pNcbRespRec.BiddingStartDate = lStartDate.Format("2006-01-02")
		// log.Println(pNcbRespRec.BiddingStartDate, "BiddingStartDate")
	} else {
		pNcbRespRec.BiddingStartDate = "0000-00-00"
	}
	//for bidding end date
	if pNcbRespRec.BiddingEndDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pNcbRespRec.BiddingEndDate)
		pNcbRespRec.BiddingEndDate = lEndDate.Format("2006-01-02")
		// log.Println(pNcbRespRec.BiddingEndDate, "BiddingEndDate")
	} else {
		pNcbRespRec.BiddingEndDate = "0000-00-00"
	}
	//for t1 bidding start date
	if pNcbRespRec.T1ModStartDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pNcbRespRec.T1ModStartDate)
		pNcbRespRec.T1ModStartDate = lEndDate.Format("2006-01-02")
		// log.Println(pNcbRespRec.T1ModStartDate, "T1ModStartDate")
	} else {
		pNcbRespRec.T1ModStartDate = "0000-00-00"
	}
	//for t1 bidding end date
	if pNcbRespRec.T1ModEndDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pNcbRespRec.T1ModEndDate)
		pNcbRespRec.T1ModEndDate = lEndDate.Format("2006-01-02")
		// log.Println(pNcbRespRec.T1ModEndDate, "T1ModEndDate")
	} else {
		pNcbRespRec.T1ModEndDate = "0000-00-00"
	}

	lString, lErr = common.ChangeTimeFormat("15:04:05", pNcbRespRec.T1ModStartTime)
	if lErr != nil {
		log.Println("ChangeTimeFormat Error: 07", lErr)
		pNcbRespRec.T1ModStartTime = "00:00:00"
	} else {
		pNcbRespRec.T1ModStartTime = lString
	}
	// log.Println("pNcbRespRec.T1ModStartTime", pNcbRespRec.T1ModStartTime)

	lString, lErr = common.ChangeTimeFormat("15:04:05", pNcbRespRec.T1ModEndTime)
	if lErr != nil {
		log.Println("ChangeTimeFormat Error: 08", lErr)
		pNcbRespRec.T1ModEndTime = "00:00:00"
	} else {
		pNcbRespRec.T1ModEndTime = lString
	}
	// log.Println("pNcbRespRec.T1ModEndTime", pNcbRespRec.T1ModEndTime)

	log.Println("ChangeDateFormat (-)")
	return pNcbRespRec
}
