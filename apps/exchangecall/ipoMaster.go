package exchangecall

import (
	"database/sql"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bseipo"
	"fcs23pkg/integration/nse/nseipo"
	"log"
	"strconv"
	"time"
)

// func FetchIPOmaster(pUser string, pBrokerId int) error {
// 	log.Println("FetchIPOmaster (+)")

// 	lErr1 := FetchNseIPOmaster(pUser, pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("EFIM01", lErr1)
// 		return lErr1
// 	} else {
// 		lErr2 := FetchBseIPOmaster(pUser, pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("EFIM02", lErr2)
// 			return lErr2
// 		}
// 	}
// 	log.Println("FetchIPOmaster (-)")
// 	return nil
// }

func FetchIPOmaster(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchIPOmaster (+)")
	lString := common.ErrorCode
	lTokenNse, lErr1 := FetchNseIPOmaster(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EFIM01", lErr1)
		return lString, lErr1
	} else {
		if lTokenNse != common.ErrorCode {
			lString = common.SuccessCode
		}
		lTokenBse, lErr2 := FetchBseIPOmaster(pUser, pBrokerId)
		if lErr2 != nil {
			log.Println("EFIM02", lErr2)
			return lString, lErr2
		} else {
			if lTokenBse != common.ErrorCode {
				lString = common.SuccessCode
			}
		}
	}
	log.Println("FetchIPOmaster (-)")
	return lString, nil
}

/*
Pupose: This method returns the collection of data from the  a_ipo_oder_header database
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the dpStructArr data
		from the a_ipo_oder_header Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
// func FetchNseIPOmaster(pUser string, pBrokerId int) error {
// 	log.Println("FetchIPOmasterNSE (+)")

// 	lToken, lErr1 := GetToken(pUser, pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("EFNIM01", lErr1)
// 		return lErr1
// 	} else {
// 		if lToken != "" {
// 			log.Println(lToken, "Token Fetched successfully...")
// 			lErr2 := getIpoMaster(lToken, pUser)
// 			if lErr2 != nil {
// 				log.Println("EFNIM02", lErr2)
// 				return lErr2
// 			}
// 		}
// 	}
// 	log.Println("FetchIPOmasterNSE (-)")
// 	return nil
// }
func FetchNseIPOmaster(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchIPOmasterNSE (+)")
	lNoToken := common.ErrorCode
	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EFNIM01", lErr1)
		lNoToken = common.ErrorCode
		return lNoToken, lErr1
	} else {
		if lToken != "" {
			log.Println(lToken, "Token Fetched successfully...")
			lErr2 := getIpoMaster(lToken, pUser)
			if lErr2 != nil {
				log.Println("EFNIM02", lErr2)
				return lNoToken, lErr2
			}
			lNoToken = common.SuccessCode
		}
	}
	log.Println("FetchIPOmasterNSE (-)")
	return lNoToken, nil
}

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will insert the ipo details
		to  a_ipo_master Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func getIpoMaster(pToken string, pUser string) error {
	log.Println("getIpoMaster (+)")
	//create instance for IpoResponseStruct
	// var lIpoResRec nseipo.IpoResponseStruct

	lIporesp, lErr1 := nseipo.IpoMaster(pToken, pUser)
	if lErr1 != nil {
		log.Println("EGIM01", lErr1)
		return lErr1
	} else {
		// log.Println("IPO response", lIpoResponse)
		lErr2 := InsertDatas(lIporesp, pUser, common.NSE)
		if lErr2 != nil {
			log.Println("EGIM02", lErr2)
			return lErr2
		}
	}
	log.Println("getIpoMaster (-)")
	return nil
}

/*
Pupose: This method helpsto insert new records in  a_ipo_master datatable
Parameters:

	(ipoResponseStruct)

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will insert the ipoResponseStruct data
		to the a_ipo_master Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func InsertDatas(pIpoResponse nseipo.IpoResponseStruct, pUser string, pExchange string) error {
	log.Println("InsertDatas (-)")

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EID01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for lIporespidx := 0; lIporespidx < len(pIpoResponse.Data); lIporespidx++ {
			// to check the ipo is already exist or not exist
			lexistMasterId, lErr2 := CheckIfDataExists(pIpoResponse, lIporespidx, pExchange)
			if lErr2 != nil {
				log.Println("EID02", lErr2)
				return lErr2
			} else {

				if lexistMasterId == 0 {
					lErr3 := InsertNewRecord(pIpoResponse, lIporespidx, pUser, pExchange)
					if lErr3 != nil {
						log.Println("EID03", lErr3)
						return lErr3
					}
				} else {
					lErr4 := UpdateRecord(pIpoResponse, lIporespidx, lexistMasterId, pUser, pExchange)
					if lErr4 != nil {
						log.Println("EID04", lErr4)
						return lErr4
					}
				}
			}
		}
	}
	log.Println("InsertDatas (-)")
	return nil
}

/*
Pupose: This method helps to update Ipo master if the exchange Updated any Ipo records
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	*On Sucess
	=========
	In case of a successful execution of this method, you will update the a_ipo_master Data Table

	!On Error
	========
	In case of any exception during the execution of this method you will get the
	error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func UpdateRecord(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string, pExchange string) error {
	log.Println("UpdateRecord (+)")

	lErr1 := UpdateIpoMaster(pIpoResponse, pIporespidx, pMasterId, pUser, pExchange)
	if lErr1 != nil {
		log.Println("EUR01", lErr1)
		return lErr1
	}
	log.Println("UpdateRecord (-)")
	return lErr1
}

/*
Pupose: This method helps to Update the master table in database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	*On Sucess
	=========
	In case of a successful execution of this method, you will update master the a_ipo_master Data Table

	!On Error
	========
	In case of any exception during the execution of this method you will get the
	error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func UpdateIpoMaster(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string, pExchange string) error {
	log.Println("UpdateIpoMaster (+)")
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EUIM01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		//dateformats

		if pExchange == common.NSE {

			DateFormats(pIpoResponse, pIporespidx)
		}

		lSqlString := `update a_ipo_master aim
						set aim.Symbol = ?, aim.Name = ?,aim.BiddingStartDate = ?,aim.BiddingEndDate = ?, aim.DailyStartTime = ?,
						aim.DailyEndTime = ?,aim.MaxPrice = ?,aim.MinPrice = ?,aim.MinBidQuantity = ?,aim.LotSize = ?,
						aim.Registrar = ?,aim.T1ModStartDate = ?,aim.T1ModEndDate = ?,aim.T1ModStartTime = ?,
						aim.T1ModEndTime = ?,aim.TickSize = ?, aim.FaceValue = ?,aim.IssueSize = ?,aim.CutOffPrice = ?,
						aim.Isin = ?,aim.IssueType = ?, aim.UpdatedBy = ?, aim.UpdatedDate = now(),aim.SubType = ?
						where aim.Id = ?`

		_, lErr2 := lDb.Exec(lSqlString, pIpoResponse.Data[pIporespidx].Symbol, pIpoResponse.Data[pIporespidx].Name,
			pIpoResponse.Data[pIporespidx].BiddingStartDate, pIpoResponse.Data[pIporespidx].BiddingEndDate,
			pIpoResponse.Data[pIporespidx].DailyStartTime, pIpoResponse.Data[pIporespidx].DailyEndTime,
			pIpoResponse.Data[pIporespidx].MaxPrice, pIpoResponse.Data[pIporespidx].MinPrice,
			pIpoResponse.Data[pIporespidx].MinBidQuantity, pIpoResponse.Data[pIporespidx].LotSize,
			pIpoResponse.Data[pIporespidx].Registrar, pIpoResponse.Data[pIporespidx].T1ModStartDate,
			pIpoResponse.Data[pIporespidx].T1ModEndDate, pIpoResponse.Data[pIporespidx].T1ModStartTime,
			pIpoResponse.Data[pIporespidx].T1ModEndTime, pIpoResponse.Data[pIporespidx].TickSize,
			pIpoResponse.Data[pIporespidx].FaceValue, pIpoResponse.Data[pIporespidx].IssueSize,
			pIpoResponse.Data[pIporespidx].CutOffPrice, pIpoResponse.Data[pIporespidx].ISIN,
			pIpoResponse.Data[pIporespidx].IssueType, pUser, pIpoResponse.Data[pIporespidx].SubType, pMasterId)

		if lErr2 != nil {
			log.Println("EUIM02", lErr2)
			return lErr2
		} else {
			lErr3 := UpdateCategory(pIpoResponse, pIporespidx, pMasterId, pUser)
			if lErr3 != nil {
				log.Println("EUIM03", lErr3)
				return lErr3
			} else {
				lErr4 := UpdateSubCategory(pIpoResponse, pIporespidx, pMasterId, pUser)
				if lErr4 != nil {
					log.Println("EUIM04", lErr4)
					return lErr4
				} else {
					lErr5 := UpdateSeries(pIpoResponse, pIporespidx, pMasterId, pUser)
					if lErr5 != nil {
						log.Println("EUIM05", lErr5)
						return lErr5
					}
				}
			}
		}
	}
	log.Println("UpdateIpoMaster (-)")
	return nil
}

/*
Pupose: This method helps to Update the category table in database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	*On Sucess
	=========
	In case of a successful execution of this method, you will update the  a_ipo_category Data Table

	!On Error
	========
	In case of any exception during the execution of this method you will get the
	error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func UpdateCategory(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string) error {
	log.Println("UpdateCategory (+)")

	var lTime sql.NullString
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EUC01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if len(pIpoResponse.Data[pIporespidx].CategoryDetailsArr) != 0 {
			// local variable assinging category values
			lCategoryArr := pIpoResponse.Data[pIporespidx].CategoryDetailsArr

			for lCategoryidx := 0; lCategoryidx < len(lCategoryArr); lCategoryidx++ {

				if lCategoryArr[lCategoryidx].StartTime == "" && lCategoryArr[lCategoryidx].EndTime == "" {

					lCategoryId, lErr2 := FindCategory(lCategoryArr[lCategoryidx], pMasterId)
					if lErr2 != nil {
						log.Println("EUC02", lErr2)
						return lErr2
					} else {
						if lCategoryId != 0 {

							lSqlString := `update a_ipo_categories aic
									set aic.Code = ?, aic.StartTime = ?, aic.EndTime = ?,aic.UpdatedBy = ?, aic.UpdatedDate = now() 
									where aic.MasterId = ? and aic.Id = ?`

							_, lErr3 := lDb.Exec(lSqlString, lCategoryArr[lCategoryidx].Code, lTime, lTime, pUser, pMasterId, lCategoryId)
							if lErr3 != nil {
								log.Println("EUC03", lErr3)
								return lErr3
							}
						} else {
							lSqlString := `insert into a_ipo_categories (MasterId,Code,StartTime,EndTime,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values(?,?,?,?,?,now(),?,now())`

							_, lErr3 := lDb.Exec(lSqlString, pMasterId, lCategoryArr[lCategoryidx].Code, lTime, lTime, pUser, pUser)
							if lErr3 != nil {
								log.Println("EIC03", lErr3)
								return lErr3
							}
						}
					}
				} else {

					lCategoryId, lErr4 := FindCategory(lCategoryArr[lCategoryidx], pMasterId)
					if lErr4 != nil {
						log.Println("EUC04", lErr4)
						return lErr4
					} else {
						if lCategoryId != 0 {

							lSqlString := `update a_ipo_categories aic
									set aic.Code = ?, aic.StartTime = ?, aic.EndTime = ?,aic.UpdatedBy = ?, aic.UpdatedDate = now() 
									where aic.MasterId = ? and aic.Id = ?`

							_, lErr5 := lDb.Exec(lSqlString, lCategoryArr[lCategoryidx].Code, lCategoryArr[lCategoryidx].StartTime,
								lCategoryArr[lCategoryidx].EndTime, pUser, pMasterId, lCategoryId)
							if lErr5 != nil {
								log.Println("EUC05", lErr5)
								return lErr5
							}
						} else {
							lSqlString := `insert into a_ipo_categories (MasterId,Code,StartTime,EndTime,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values(?,?,?,?,?,now(),?,now())`

							_, lErr3 := lDb.Exec(lSqlString, pMasterId, lCategoryArr[lCategoryidx].Code, lCategoryArr[lCategoryidx].StartTime,
								lCategoryArr[lCategoryidx].EndTime, pUser, pUser)
							if lErr3 != nil {
								log.Println("EIC03", lErr3)
								return lErr3
							}
						}
					}
				}

			}
		}
	}
	log.Println("UpdateCategory (-)")
	return nil
}

/*
Pupose: This method returns the  Category id from the database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the CategoryId
		from the a_ipo_Category Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func FindCategory(pCategoryRec nseipo.CategoryDetailsStruct, pMasterId int) (int, error) {
	log.Println("FindCategory (+)")

	lCategoryId := 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EFC01", lErr1)
		return lCategoryId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select aic.Id 
					from a_ipo_categories aic 
					where aic.MasterId = ? and aic.Code = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pCategoryRec.Code)
		if lErr2 != nil {
			log.Println("EFC02", lErr2)
			return lCategoryId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lCategoryId)
				if lErr3 != nil {
					log.Println("EFC03", lErr3)
					return lCategoryId, lErr3
				}
			}
		}
	}
	log.Println("FindCategory (-)")
	return lCategoryId, nil
}

/*
Pupose: This method helps to Update the  Subcategory details in the database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will successfully updated the
		record of  a_ipo_Subcategory Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func UpdateSubCategory(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string) error {
	log.Println("UpdateSubCategory (+)")
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EUSC01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if len(pIpoResponse.Data[pIporespidx].SubCategoryDetailsArr) != 0 {
			// local variable assinging category values
			lSubCategoryArr := pIpoResponse.Data[pIporespidx].SubCategoryDetailsArr

			for lsubcategoryidx := 0; lsubcategoryidx < len(lSubCategoryArr); lsubcategoryidx++ {

				lSubCatId, lErr2 := FindSubCategory(lSubCategoryArr[lsubcategoryidx], pMasterId)
				if lErr2 != nil {
					log.Println("EUSC02", lErr2)
					return lErr2
				} else {
					if lSubCatId != 0 {

						lSqlString := `update a_ipo_subcategory ais
							set ais.CaCode = ?, ais.SubCatCode = ?, ais.Min_Value = ?, ais.Max_Value = ?, ais.AllowCutOff = ?,
							ais.DiscountType = ?, ais.DiscountPrice = ?, ais.AllowUpi = ?, ais.MaxQuantity = ?, 
							ais.MaxPrice = ?, ais.MaxUpiLimit = ?, ais.UpdatedBy = ?, ais.UpdatedDate = now()
							where ais.MasterId = ? and ais.Id = ?`

						_, lErr3 := lDb.Exec(lSqlString, lSubCategoryArr[lsubcategoryidx].CaCode,
							lSubCategoryArr[lsubcategoryidx].SubCatCode, lSubCategoryArr[lsubcategoryidx].MinValue,
							lSubCategoryArr[lsubcategoryidx].MaxValue, lSubCategoryArr[lsubcategoryidx].AllowCutOff,
							lSubCategoryArr[lsubcategoryidx].DiscountType, lSubCategoryArr[lsubcategoryidx].DiscountPrice,
							lSubCategoryArr[lsubcategoryidx].AllowUpi, lSubCategoryArr[lsubcategoryidx].MaxQuantity,
							lSubCategoryArr[lsubcategoryidx].MaxPrice, lSubCategoryArr[lsubcategoryidx].MaxUpiLimit, pUser, pMasterId, lSubCatId)

						if lErr3 != nil {
							log.Println("EUSC03", lErr3)
							return lErr3
						}
					} else {
						maxQuantityValue := lSubCategoryArr[lsubcategoryidx].MaxQuantity
						switch maxQuantityValue.(type) {
						case int:
							lMaxQtyInt := int(maxQuantityValue.(int))
							log.Println("lMaxQty (int):", lMaxQtyInt)
							lSubCategoryArr[lsubcategoryidx].MaxQuantity = lMaxQtyInt
						case float64:
							// Handle float64 value
							lMaxQtyFloat := float64(maxQuantityValue.(float64))
							log.Println("maxQuantity (float64):", lMaxQtyFloat)
							lSubCategoryArr[lsubcategoryidx].MaxQuantity = int(lMaxQtyFloat)
						default:
							// lMaxQtyInt := int(maxQuantityValue.(int))
							// log.Println("lMaxQty (int):", lMaxQtyInt)
							// lSubCategoryArr[lSubcategoryidx].MaxQuantity = lMaxQtyInt
							// Assuming maxQuantityValue is the variable you're trying to convert to an int.
							if maxQuantityValue != nil {
								lIntValue, ok := maxQuantityValue.(int)
								if !ok {
									// Handle the case where the conversion is not successful.
									log.Println("Conversion to int failed.", lIntValue)
								} else {
									log.Println("lMaxQty (int):", lIntValue)
									lSubCategoryArr[lsubcategoryidx].MaxQuantity = lIntValue
								}
							} else {
								// Handle the case where maxQuantityValue is nil.
								log.Println("maxQuantityValue is nil.", maxQuantityValue)
							}

						}

						lSqlString := `insert into a_ipo_subcategory (MasterId,CaCode,SubCatCode,Min_Value,Max_Value,AllowCutOff,
							DiscountType,DiscountPrice,AllowUpi,MaxQuantity,MaxPrice,MaxUpiLimit,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
							values(?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

						_, lErr2 := lDb.Exec(lSqlString, pMasterId, lSubCategoryArr[lsubcategoryidx].CaCode,
							lSubCategoryArr[lsubcategoryidx].SubCatCode, lSubCategoryArr[lsubcategoryidx].MinValue,
							lSubCategoryArr[lsubcategoryidx].MaxValue, lSubCategoryArr[lsubcategoryidx].AllowCutOff,
							lSubCategoryArr[lsubcategoryidx].DiscountType, lSubCategoryArr[lsubcategoryidx].DiscountPrice,
							lSubCategoryArr[lsubcategoryidx].AllowUpi, lSubCategoryArr[lsubcategoryidx].MaxQuantity,
							lSubCategoryArr[lsubcategoryidx].MaxPrice, lSubCategoryArr[lsubcategoryidx].MaxUpiLimit, pUser, pUser)

						if lErr2 != nil {
							log.Println("EISC02", lErr2)
							return lErr2
						}

					}

				}
			}
		}
	}
	log.Println("UpdateSubCategory (-)")
	return nil
}

/*
Pupose: This method returns the  SubCategory id from the database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the SubCategoryId
		from the a_ipo_SubCategory Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func FindSubCategory(pSubCatRec nseipo.SubCategorySettingsStruct, pMasterId int) (int, error) {
	log.Println("FindSubCategory (+)")

	lSubCatId := 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EFSC01", lErr1)
		return lSubCatId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select ais.Id 
						from a_ipo_subcategory ais 
						where ais.MasterId = ? and ais.CaCode = ? and ais.SubCatCode = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pSubCatRec.CaCode, pSubCatRec.SubCatCode)
		if lErr2 != nil {
			log.Println("EFSC02", lErr2)
			return lSubCatId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSubCatId)
				if lErr3 != nil {
					log.Println("EFSC03", lErr3)
					return lSubCatId, lErr3
				}
			}
		}
	}
	log.Println("FindSubCategory (-)")
	return lSubCatId, nil
}

/*
Pupose: This method helps to Update the  series data in the database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will successfully updated the
		record of  a_ipo_series Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func UpdateSeries(pIporesp nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string) error {
	log.Println("UpdateSeries (+)")
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EUS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if len(pIporesp.Data[pIporespidx].SeriesDetailsArr) != 0 {
			// local variable assinging category values
			lSeriesArr := pIporesp.Data[pIporespidx].SeriesDetailsArr

			for lSeriesidx := 0; lSeriesidx < len(lSeriesArr); lSeriesidx++ {

				lSeriesId, lErr2 := FindSeries(lSeriesArr[lSeriesidx], pMasterId)
				if lErr2 != nil {
					log.Println("EUS02", lErr2)
					return lErr2
				} else {

					lSqlString := `update a_ipo_series a
								set a.Code = ?, a.Description = ?,a.UpdatedBy = ?, a.UpdatedDate = now()
								where a.MasterId = ? and a.Id = ?`

					_, lErr3 := lDb.Exec(lSqlString, lSeriesArr[lSeriesidx].Code, lSeriesArr[lSeriesidx].Desc, pUser, pMasterId, lSeriesId)
					if lErr3 != nil {
						log.Println("EUS03", lErr3)
						return lErr3
					}
				}
			}
		}
	}
	log.Println("UpdateSeries (-)")
	return nil
}

/*
Pupose: This method returns the  series id from the database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the seriesId
		from the a_ipo_series Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func FindSeries(pSeriesRec nseipo.SeriesDetailsStruct, pMasterId int) (int, error) {
	log.Println("FindSeries (+)")

	lSeriesId := 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EFS01", lErr1)
		return lSeriesId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select ais.Id 
						from a_ipo_series ais
						where ais.MasterId = ? and ais.Code = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId, pSeriesRec.Code)
		if lErr2 != nil {
			log.Println("EFS02", lErr2)
			return lSeriesId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lSeriesId)
				if lErr3 != nil {
					log.Println("EFS03", lErr3)
					return lSeriesId, lErr3
				}
			}
		}
	}
	log.Println("FindSeries (-)")
	return lSeriesId, nil
}

// =====================Update code Finished --- >>> Inserting records code==============================

/*
Pupose: This method helps to changing the date formats received from exchange to the sql date format
Parameters:

	(ipoResponseStruct , ipoResponseidx )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will converted the date format from
		"DD/MM/YYYY" to "YYYY-MM-DD"

Author: Pavithra
Date: 04JUNE2023
*/
func DateFormats(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int) {
	log.Println("DateFormats (+)")

	// for bidding start date
	if pIpoResponse.Data[pIporespidx].BiddingStartDate != "" {
		lStartDate, _ := time.Parse("02-01-2006", pIpoResponse.Data[pIporespidx].BiddingStartDate)
		pIpoResponse.Data[pIporespidx].BiddingStartDate = lStartDate.Format("2006-01-02")
		// log.Println(pIpoResponse.Data[pIporespidx].BiddingStartDate, "BiddingStartDate")
	} else {
		pIpoResponse.Data[pIporespidx].BiddingStartDate = "0000-00-00"
	}

	// for bidding end date
	if pIpoResponse.Data[pIporespidx].BiddingEndDate != "" {
		lEndDate, _ := time.Parse("02-01-2006", pIpoResponse.Data[pIporespidx].BiddingEndDate)
		pIpoResponse.Data[pIporespidx].BiddingEndDate = lEndDate.Format("2006-01-02")
		// log.Println(pIpoResponse.Data[pIporespidx].BiddingEndDate, "BiddingEndDate")
	} else {
		pIpoResponse.Data[pIporespidx].BiddingEndDate = "0000-00-00"
	}

	// for t1 start date
	if pIpoResponse.Data[pIporespidx].T1ModStartDate != "" {
		lT1StartDate, _ := time.Parse("02-01-2006", pIpoResponse.Data[pIporespidx].T1ModStartDate)
		pIpoResponse.Data[pIporespidx].T1ModStartDate = lT1StartDate.Format("2006-01-02")
		// log.Println(pIpoResponse.Data[pIporespidx].T1ModStartDate, "T1ModStartDate")
	} else {
		pIpoResponse.Data[pIporespidx].T1ModStartDate = "0000-00-00"
	}

	// for t1 end date
	if pIpoResponse.Data[pIporespidx].T1ModEndDate != "" {
		lT1EndDate, _ := time.Parse("02-01-2006", pIpoResponse.Data[pIporespidx].T1ModEndDate)
		pIpoResponse.Data[pIporespidx].T1ModEndDate = lT1EndDate.Format("2006-01-02")
		// log.Println(pIpoResponse.Data[pIporespidx].T1ModEndDate, "T1ModEndDate")
	} else {
		pIpoResponse.Data[pIporespidx].T1ModEndDate = "0000-00-00"
	}

	log.Println("DateFormats (-)")
}

/*
Pupose: This method helps to check the ipo record is already exist or not
Parameters:

	(ipoResponseStruct , ipoResponseidx )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the masterId data
		from the a_ipo_master Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func CheckIfDataExists(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pExchange string) (int, error) {
	log.Println("CheckIfDataExists (+)")

	lmasterId := 0
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("ECIDE01", lErr1)
		return lmasterId, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select aim.Id 
						from a_ipo_master aim 
						where aim.Symbol = ? 
						and aim.Isin = ?
						and aim.Exchange = ?`
		lRows, lErr2 := lDb.Query(lCoreString, pIpoResponse.Data[pIporespidx].Symbol, pIpoResponse.Data[pIporespidx].ISIN, pExchange)
		if lErr2 != nil {
			log.Println("ECIDE02", lErr2)
			return lmasterId, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lmasterId)
				if lErr3 != nil {
					log.Println("ECIDE03", lErr3)
					return lmasterId, lErr3
				}
			}
		}
	}
	log.Println("CheckIfDataExists (-)")
	return lmasterId, nil
}

/*
Pupose:  This method is used to inserting new Ipo Details to the  a_ipo_master database
Parameters:

	(ipoResponseStruct , ipoResponseidx )

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will inserted the master data
		to a_ipo_master Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func InsertNewRecord(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pUser string, pExchange string) error {
	log.Println("InsertNewRecord (+)")

	lmasterId, lErr1 := InsertIpoMaster(pIpoResponse, pIporespidx, pUser, pExchange)
	if lErr1 != nil {
		log.Println("EINR01", lErr1)
		return lErr1
	} else {
		lErr2 := InsertCategory(pIpoResponse, pIporespidx, lmasterId, pUser)
		if lErr2 != nil {
			log.Println("EINR02", lErr2)
			return lErr2
		} else {
			lErr3 := InsertSubCategory(pIpoResponse, pIporespidx, lmasterId, pUser)
			if lErr3 != nil {
				log.Println("EINR03", lErr3)
				return lErr3
			} else {
				lErr4 := InsertSeries(pIpoResponse, pIporespidx, lmasterId, pUser)
				if lErr4 != nil {
					log.Println("EINR04", lErr4)
					return lErr4
				}
			}
		}
	}
	log.Println("InsertNewRecord (-)")
	return nil
}

/*
Pupose:  This method is used to inserting the collection of data to the  a_ipo_master database
Parameters:

	(ipoResponseStruct , ipoResponseidx)

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will inserted the master data
		to a_ipo_master Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func InsertIpoMaster(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pUser string, pExchange string) (int, error) {
	log.Println("InsertIpoMaster (+)")

	lmasterId := 0
	//dateformats
	if pExchange == common.NSE {
		DateFormats(pIpoResponse, pIporespidx)
	}
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EIIM01", lErr1)
		return lmasterId, lErr1
	} else {
		defer lDb.Close()
		lSqlString := `insert into a_ipo_master (Symbol,Name,BiddingStartDate,BiddingEndDate,DailyStartTime,DailyEndTime,
		MaxPrice,MinPrice,MinBidQuantity,LotSize,Registrar,T1ModStartDate,T1ModEndDate,T1ModStartTime,
		T1ModEndTime,TickSize,FaceValue,IssueSize,CutOffPrice,Isin,IssueType,SubType,CreatedBy,CreatedDate,Exchange,UpdatedBy,UpdatedDate)
		values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?,now())`

		lInsertedMasterId, lErr2 := lDb.Exec(lSqlString, pIpoResponse.Data[pIporespidx].Symbol, pIpoResponse.Data[pIporespidx].Name,
			pIpoResponse.Data[pIporespidx].BiddingStartDate, pIpoResponse.Data[pIporespidx].BiddingEndDate,
			pIpoResponse.Data[pIporespidx].DailyStartTime, pIpoResponse.Data[pIporespidx].DailyEndTime,
			pIpoResponse.Data[pIporespidx].MaxPrice, pIpoResponse.Data[pIporespidx].MinPrice,
			pIpoResponse.Data[pIporespidx].MinBidQuantity, pIpoResponse.Data[pIporespidx].LotSize,
			pIpoResponse.Data[pIporespidx].Registrar, pIpoResponse.Data[pIporespidx].T1ModStartDate,
			pIpoResponse.Data[pIporespidx].T1ModEndDate, pIpoResponse.Data[pIporespidx].T1ModStartTime,
			pIpoResponse.Data[pIporespidx].T1ModEndTime, pIpoResponse.Data[pIporespidx].TickSize,
			pIpoResponse.Data[pIporespidx].FaceValue, pIpoResponse.Data[pIporespidx].IssueSize,
			pIpoResponse.Data[pIporespidx].CutOffPrice, pIpoResponse.Data[pIporespidx].ISIN,
			pIpoResponse.Data[pIporespidx].IssueType, pIpoResponse.Data[pIporespidx].SubType, pUser, pExchange, pUser)

		if lErr2 != nil {
			log.Println("EIIM02", lErr2)
			return lmasterId, lErr2
		} else {
			lreturnId, _ := lInsertedMasterId.LastInsertId()
			lmasterId = int(lreturnId)
		}
	}
	log.Println("InsertIpoMaster (-)")
	return lmasterId, nil
}

/*
Pupose:  This method is used to inserting the collection of data to the  a_ipo_category database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId)

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will inserted the category data
		to a_ipo_category Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func InsertCategory(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string) error {
	log.Println("InsertCategory (+)")

	var lTime sql.NullString

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EIC01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if len(pIpoResponse.Data[pIporespidx].CategoryDetailsArr) != 0 {
			// local variable assinging category values
			lCategoryArr := pIpoResponse.Data[pIporespidx].CategoryDetailsArr

			for lCategoryidx := 0; lCategoryidx < len(lCategoryArr); lCategoryidx++ {

				if lCategoryArr[lCategoryidx].StartTime == "" && lCategoryArr[lCategoryidx].EndTime == "" {
					lSqlString := `insert into a_ipo_categories (MasterId,Code,StartTime,EndTime,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
								values(?,?,?,?,?,now(),?,now())`

					_, lErr2 := lDb.Exec(lSqlString, pMasterId, lCategoryArr[lCategoryidx].Code, lTime, lTime, pUser, pUser)
					if lErr2 != nil {
						log.Println("EIC02", lErr2)
						return lErr2
					}
				} else {

					lSqlString := `insert into a_ipo_categories (MasterId,Code,StartTime,EndTime,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
								values(?,?,?,?,?,now(),?,now())`

					_, lErr3 := lDb.Exec(lSqlString, pMasterId, lCategoryArr[lCategoryidx].Code, lCategoryArr[lCategoryidx].StartTime,
						lCategoryArr[lCategoryidx].EndTime, pUser, pUser)
					if lErr3 != nil {
						log.Println("EIC03", lErr3)
						return lErr3
					}
				}
			}
		}
	}
	log.Println("InsertCategory (-)")
	return nil
}

/*
Pupose: This method is used to inserting the collection of data to the  a_ipo_Subcategory database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId)

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will inserted the Subcategory data
		to a_ipo_Subcategory Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func InsertSubCategory(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string) error {
	log.Println("InsertSubCategory (+)")
	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EISC01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if len(pIpoResponse.Data[pIporespidx].SubCategoryDetailsArr) != 0 {
			// local variable assinging category values
			lSubCategoryArr := pIpoResponse.Data[pIporespidx].SubCategoryDetailsArr

			for lSubcategoryidx := 0; lSubcategoryidx < len(lSubCategoryArr); lSubcategoryidx++ {

				maxQuantityValue := lSubCategoryArr[lSubcategoryidx].MaxQuantity

				// var lMaxQty int // Declare lMaxQty outside of the switch

				switch maxQuantityValue.(type) {
				case int:
					lMaxQtyInt := int(maxQuantityValue.(int))
					log.Println("lMaxQty (int):", lMaxQtyInt)
					lSubCategoryArr[lSubcategoryidx].MaxQuantity = lMaxQtyInt
				case float64:
					// Handle float64 value
					lMaxQtyFloat := float64(maxQuantityValue.(float64))
					log.Println("maxQuantity (float64):", lMaxQtyFloat)
					lSubCategoryArr[lSubcategoryidx].MaxQuantity = int(lMaxQtyFloat)
				default:
					// lMaxQtyInt := int(maxQuantityValue.(int))
					// log.Println("lMaxQty (int):", lMaxQtyInt)
					// lSubCategoryArr[lSubcategoryidx].MaxQuantity = lMaxQtyInt
					// Assuming maxQuantityValue is the variable you're trying to convert to an int.
					if maxQuantityValue != nil {
						lIntValue, ok := maxQuantityValue.(int)
						if !ok {
							// Handle the case where the conversion is not successful.
							log.Println("Conversion to int failed.", lIntValue)
						} else {
							log.Println("lMaxQty (int):", lIntValue)
							lSubCategoryArr[lSubcategoryidx].MaxQuantity = lIntValue
						}
					} else {
						// Handle the case where maxQuantityValue is nil.
						log.Println("maxQuantityValue is nil.", maxQuantityValue)
					}

				}

				lSqlString := `insert into a_ipo_subcategory (MasterId,CaCode,SubCatCode,Min_Value,Max_Value,AllowCutOff,
				DiscountType,DiscountPrice,AllowUpi,MaxQuantity,MaxPrice,MaxUpiLimit,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
				values(?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,now())`

				_, lErr2 := lDb.Exec(lSqlString, pMasterId, lSubCategoryArr[lSubcategoryidx].CaCode,
					lSubCategoryArr[lSubcategoryidx].SubCatCode, lSubCategoryArr[lSubcategoryidx].MinValue,
					lSubCategoryArr[lSubcategoryidx].MaxValue, lSubCategoryArr[lSubcategoryidx].AllowCutOff,
					lSubCategoryArr[lSubcategoryidx].DiscountType, lSubCategoryArr[lSubcategoryidx].DiscountPrice,
					lSubCategoryArr[lSubcategoryidx].AllowUpi, lSubCategoryArr[lSubcategoryidx].MaxQuantity,
					lSubCategoryArr[lSubcategoryidx].MaxPrice, lSubCategoryArr[lSubcategoryidx].MaxUpiLimit, pUser, pUser)

				if lErr2 != nil {
					log.Println("EISC02", lErr2)
					return lErr2
				}
			}
		}
	}
	log.Println("InsertSubCategory (-)")
	return nil
}

/*
Pupose: This method is used to insert series dewtails into database
Parameters:

	(ipoResponseStruct , ipoResponseidx , masterId)

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will inserted the series data
		to a_ipo_series Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 04JUNE2023
*/
func InsertSeries(pIpoResponse nseipo.IpoResponseStruct, pIporespidx int, pMasterId int, pUser string) error {
	log.Println("InsertSeries (+)")

	// Calling LocalDbConect method in ftdb to estabish the database connection
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("EIS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if len(pIpoResponse.Data[pIporespidx].SeriesDetailsArr) != 0 {
			// local variable assinging category values
			lSeriesArr := pIpoResponse.Data[pIporespidx].SeriesDetailsArr

			for lSeriesidx := 0; lSeriesidx < len(lSeriesArr); lSeriesidx++ {

				lSqlString := `insert into a_ipo_series (MasterId,Code,Description,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
								values(?,?,?,?,now(),?,now())`

				_, lErr2 := lDb.Exec(lSqlString, pMasterId, lSeriesArr[lSeriesidx].Code, lSeriesArr[lSeriesidx].Desc, pUser, pUser)
				if lErr2 != nil {
					log.Println("EIS02", lErr2)
					return lErr2
				}
			}
		}
	}
	log.Println("InsertSeries (-)")
	return nil
}

//----------------------------------------------------------------------------------------------------
// BSE IPO MASTER
//----------------------------------------------------------------------------------------------------

/*
Pupose: This method returns the collection of data from the  a_ipo_oder_header database
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will get the dpStructArr data
		from the a_ipo_oder_header Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 12SEP2023
*/
// func FetchBseIPOmaster(pUser string, pBrokerId int) error {
// 	log.Println("FetchBseIPOmaster (+)")

// 	lToken, lErr1 := BseGetToken(pUser, pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("EFBIM01", lErr1)
// 		return lErr1
// 	} else {
// 		if lToken != "" {
// 			log.Println(lToken, "BSE Token Fetched successfully...")
// 			lErr2 := GetBseIpoMaster(lToken, pUser, pBrokerId)
// 			if lErr2 != nil {
// 				log.Println("EFBIM02", lErr2)
// 				return lErr2
// 			}
// 		}
// 	}
// 	log.Println("FetchBseIPOmaster (-)")
// 	return nil
// }
func FetchBseIPOmaster(pUser string, pBrokerId int) (string, error) {
	log.Println("FetchBseIPOmaster (+)")
	lNoToken := common.ErrorCode
	lToken, lErr1 := BseGetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EFBIM01", lErr1)
		lNoToken = common.ErrorCode
		return lNoToken, lErr1
	} else {
		if lToken != "" {
			log.Println(lToken, "BSE Token Fetched successfully...")
			lErr2 := GetBseIpoMaster(lToken, pUser, pBrokerId)
			if lErr2 != nil {
				log.Println("EFBIM02", lErr2)
				return lNoToken, lErr2
			}
			lNoToken = common.SuccessCode
		}
	}
	log.Println("FetchBseIPOmaster (-)")
	return lNoToken, nil
}

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will insert the ipo details
		to  a_ipo_master Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 12SEP2023
*/
func GetBseIpoMaster(pToken string, pUser string, pBrokerId int) error {
	log.Println("GetBseIpoMaster (+)")

	lCategoryArr, lErr1 := bseipo.BseIpoCategory(pToken, pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("EGBIM01", lErr1)
		return lErr1
	} else {
		//
		lMasterArr, lErr2 := bseipo.BseIpoMaster(pToken, pUser, pBrokerId)
		if lErr2 != nil {
			log.Println("EGBIM02", lErr2)
			return lErr2
		} else {
			lFinalMasterArr, lFinalCategoryArr := FilterBseDatas(lMasterArr, lCategoryArr)

			lNseStructDatas := ConstructIpoData(lFinalMasterArr, lFinalCategoryArr)

			lErr3 := InsertDatas(lNseStructDatas, pUser, common.BSE)
			if lErr3 != nil {
				log.Println("EGBIM03", lErr3)
				return lErr3
			} else {
				log.Println("BSEIPOINSERTED SUCCESSFULLY")
			}
		}
	}
	log.Println("GetBseIpoMaster (-)")
	return nil
}

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will insert the ipo details
		to  a_ipo_master Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 12SEP2023
*/
func FilterBseDatas(pMasterArr []bseipo.BseIpoMasterStruct, pCategoryArr []bseipo.BseIpoCategoryStruct) ([]bseipo.BseIpoMasterStruct, []bseipo.BseIpoCategoryStruct) {
	log.Println("FilterBseDatas (+)")

	//this variable is used to get the final master array
	var lFinalMasterArr []bseipo.BseIpoMasterStruct

	//this variable is used to get the final category array
	var lFinalCategoryArr []bseipo.BseIpoCategoryStruct

	for lMasterIdx := 0; lMasterIdx < len(pMasterArr); lMasterIdx++ {
		// commented by pavithra added FP issuetype validation and ISIN validation
		// if pMasterArr[lMasterIdx].AsbaNonAsba == "2" && pMasterArr[lMasterIdx].Category == "IND" && pMasterArr[lMasterIdx].IssueType == "BB" {
		// 	lFinalMasterArr = append(lFinalMasterArr, pMasterArr[lMasterIdx])
		// }
		if pMasterArr[lMasterIdx].ISIN != "" {
			if pMasterArr[lMasterIdx].AsbaNonAsba == "2" && (pMasterArr[lMasterIdx].Category == "IND" || pMasterArr[lMasterIdx].Category == "EMP" || pMasterArr[lMasterIdx].Category == "SHA") && (pMasterArr[lMasterIdx].IssueType == "BB" || pMasterArr[lMasterIdx].IssueType == "FP") {
				lFinalMasterArr = append(lFinalMasterArr, pMasterArr[lMasterIdx])
			}
		}
	}

	for lCategoryIdx := 0; lCategoryIdx < len(pCategoryArr); lCategoryIdx++ {
		if pCategoryArr[lCategoryIdx].Type == "2" && pCategoryArr[lCategoryIdx].Series == "" && (pCategoryArr[lCategoryIdx].Category == "IND" || pCategoryArr[lCategoryIdx].Category == "EMP" || pCategoryArr[lCategoryIdx].Category == "SHA") {
			lFinalCategoryArr = append(lFinalCategoryArr, pCategoryArr[lCategoryIdx])
		}
	}
	// log.Println("lFinalMasterArr", lFinalMasterArr)
	log.Println("FilterBseDatas (-)")
	return lFinalMasterArr, lFinalCategoryArr
}

/*
Pupose: This method returns the collection of data from the  exchange
Parameters:

	not applicable

Response:

	    *On Sucess
	    =========
	    In case of a successful execution of this method, you will insert the ipo details
		to  a_ipo_master Data Table

	    !On Error
	    ========
	    In case of any exception during the execution of this method you will get the
	    error details. the calling program should handle the error

Author: Pavithra
Date: 12SEP2023
*/
func ConstructIpoData(pMasterArr []bseipo.BseIpoMasterStruct, pCategoryArr []bseipo.BseIpoCategoryStruct) nseipo.IpoResponseStruct {
	log.Println("ConstructIpoData (+)")

	//this variable is used to get the final master array
	var lMasterArr nseipo.IpoResponseStruct

	for lIdx := 0; lIdx < len(pMasterArr); lIdx++ {

		var lMasterRec nseipo.IpoDetailsStruct
		var lCategoryRec nseipo.CategoryDetailsStruct
		var lSubCategoryRec nseipo.SubCategorySettingsStruct
		var lString string
		var lFloatValue float32
		var lErr error

		lMasterRec.Symbol = pMasterArr[lIdx].Symbol
		lMasterRec.Name = pMasterArr[lIdx].Name
		lMasterRec.ISIN = pMasterArr[lIdx].ISIN

		lMasterRec.IssueType = "EQUITY"
		// lCategoryRec.Code = "RETAIL"
		// lSubCategoryRec.CaCode = "RETAIL"
		// lSubCategoryRec.SubCatCode = pMasterArr[lIdx].Category
		lSubCategoryRec.AllowCutOff = true
		lSubCategoryRec.AllowUpi = true
		lSubCategoryRec.DiscountType = pMasterArr[lIdx].DiscountType
		lSubCategoryRec.MaxUpiLimit = 200000.0
		//string to int

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].IssueSize)
		if lErr != nil {
			log.Println("Strnconv Error: 01", lErr)
		} else {
			lMasterRec.IssueSize = int(lFloatValue)
		}
		// int(floatNumber)
		if pMasterArr[lIdx].TradingLot != "" {
			lMasterRec.LotSize, lErr = strconv.Atoi(pMasterArr[lIdx].TradingLot)
			if lErr != nil {
				log.Println("Strnconv Error: 02", lErr)
			}
		} else {
			lMasterRec.LotSize = 0
		}

		if pMasterArr[lIdx].MinBidQty != "" {
			lMasterRec.MinBidQuantity, lErr = strconv.Atoi(pMasterArr[lIdx].MinBidQty)
			if lErr != nil {
				log.Println("Strnconv Error: 03", lErr)
			}
		} else {
			lMasterRec.MinBidQuantity = 0
		}

		if pMasterArr[lIdx].MaxBidQty != "" {
			lSubCategoryRec.MaxQuantity, lErr = strconv.Atoi(pMasterArr[lIdx].MaxBidQty)
			if lErr != nil {
				log.Println("Strnconv Error: 04", lErr)
			}
		} else {
			lSubCategoryRec.MaxQuantity = 0
		}

		//date format
		lString, lErr = common.ChangeTimeFormat("2006-01-02", pMasterArr[lIdx].OpenDateTime)
		// log.Println(pMasterArr[lIdx].OpenDateTime, lIdx, "format 1")
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 01", lErr)
			lMasterRec.BiddingStartDate = "0001-01-01"
		} else {
			lMasterRec.BiddingStartDate = lString
		}

		lString, lErr = common.ChangeTimeFormat("2006-01-02", pMasterArr[lIdx].CloseDateTime)
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 02", lErr)
			lMasterRec.BiddingEndDate = "0001-01-01"
		} else {
			lMasterRec.BiddingEndDate = lString
		}

		lString, lErr = common.ChangeTimeFormat("15:04:05", pMasterArr[lIdx].OpenDateTime)
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 03", lErr)
			lMasterRec.DailyStartTime = "00:00:00"
		} else {
			lMasterRec.DailyStartTime = lString
		}

		lString, lErr = common.ChangeTimeFormat("15:04:05", pMasterArr[lIdx].CloseDateTime)
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 04", lErr)
			lMasterRec.DailyEndTime = "00:00:00"
		} else {
			lMasterRec.DailyEndTime = lString
		}

		lString, lErr = common.ChangeTimeFormat("2006-01-02", pMasterArr[lIdx].TplusModificationFrom)
		log.Println(pMasterArr[lIdx].TplusModificationFrom, lIdx, "format 2")
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 05", lErr)
			lMasterRec.T1ModStartDate = "0001-01-01"
		} else {
			lMasterRec.T1ModStartDate = lString
		}

		lString, lErr = common.ChangeTimeFormat("2006-01-02", pMasterArr[lIdx].TplusModificationTo)
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 06", lErr)
			lMasterRec.T1ModEndDate = "0001-01-01"
		} else {
			lMasterRec.T1ModEndDate = lString
		}

		lString, lErr = common.ChangeTimeFormat("15:04:05", pMasterArr[lIdx].TplusModificationFrom)
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 07", lErr)
			lMasterRec.T1ModStartTime = "00:00:00"
		} else {
			lMasterRec.T1ModStartTime = lString
		}

		lString, lErr = common.ChangeTimeFormat("15:04:05", pMasterArr[lIdx].TplusModificationTo)
		if lErr != nil {
			log.Println("ChangeTimeFormat Error: 08", lErr)
			lMasterRec.T1ModEndTime = "00:00:00"
		} else {
			lMasterRec.T1ModEndTime = lString
		}

		//------------------------------------------------------------------
		//string to float value
		//------------------------------------------------------------------

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].TickPrice)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 01", lErr)
			lMasterRec.TickSize = lFloatValue
		} else {
			lMasterRec.TickSize = lFloatValue
		}

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].FloorPrice)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 02", lErr)
			lMasterRec.MinPrice = lFloatValue
		} else {
			lMasterRec.MinPrice = lFloatValue
		}

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].CeilingPrice)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 03", lErr)
			lMasterRec.MaxPrice = lFloatValue
		} else {
			lMasterRec.MaxPrice = lFloatValue
		}

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].CutOff)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 04", lErr)
			lMasterRec.CutOffPrice = lFloatValue
		} else {
			lMasterRec.CutOffPrice = lFloatValue
		}

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].MinValue)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 05", lErr)
			lSubCategoryRec.MinValue = lFloatValue
		} else {
			lSubCategoryRec.MinValue = lFloatValue
		}

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].MaxValue)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 06", lErr)
			lSubCategoryRec.MaxValue = lFloatValue
		} else {
			lSubCategoryRec.MaxValue = lFloatValue
		}

		lFloatValue, lErr = common.ConvertStringToFloat(pMasterArr[lIdx].DiscountValue)
		if lErr != nil {
			log.Println("ConvertStringToFloat Error: 07", lErr)
			lSubCategoryRec.DiscountPrice = lFloatValue
		} else {
			lSubCategoryRec.DiscountPrice = lFloatValue
		}

		for lCategoryIdx := 0; lCategoryIdx < len(pCategoryArr); lCategoryIdx++ {
			if pCategoryArr[lCategoryIdx].ScripId == pMasterArr[lIdx].Symbol {

				if pMasterArr[lIdx].Category == pCategoryArr[lCategoryIdx].Category {
					if pMasterArr[lIdx].Category == "IND" {

						lCategoryRec.Code = "RETAIL"
						lSubCategoryRec.CaCode = "RETAIL"
						lSubCategoryRec.SubCatCode = pMasterArr[lIdx].Category
					} else if pMasterArr[lIdx].Category == "EMP" {
						lCategoryRec.Code = "EMPRET"
						lSubCategoryRec.CaCode = "EMPRET"
						lSubCategoryRec.SubCatCode = pMasterArr[lIdx].Category
					} else if pMasterArr[lIdx].Category == "SHA" {
						lCategoryRec.Code = "SHARET"
						lSubCategoryRec.CaCode = "SHARET"
						lSubCategoryRec.SubCatCode = pMasterArr[lIdx].Category
					}

					lString, lErr = common.ChangeTimeFormat("15:04:05", pCategoryArr[lCategoryIdx].OpenTime)
					log.Println(pCategoryArr[lCategoryIdx].OpenTime, lIdx, "format 3")
					if lErr != nil {
						log.Println("ChangeTimeFormat Error: 09", lErr)
						lCategoryRec.StartTime = "00:00:00"
					} else {
						lCategoryRec.StartTime = lString
					}

					lString, lErr = common.ChangeTimeFormat("15:04:05", pCategoryArr[lCategoryIdx].CloseTime)
					if lErr != nil {
						log.Println("ChangeTimeFormat Error: 10", lErr)
						lCategoryRec.EndTime = "00:00:00"
					} else {
						lCategoryRec.EndTime = lString
					}
					// lMasterRec.CategoryDetailsArr = append(lMasterRec.CategoryDetailsArr, lCategoryRec)
					// lMasterRec.SubCategoryDetailsArr = append(lMasterRec.SubCategoryDetailsArr, lSubCategoryRec)
					if len(lMasterArr.Data) != 0 {

						for lArrIdx := 0; lArrIdx < len(lMasterArr.Data); lArrIdx++ {
							// log.Println("lMasterArr.Data[lArrIdx].Symbol == pCategoryArr[lCategoryIdx].ScripId", lMasterArr.Data[lArrIdx].Symbol, pCategoryArr[lCategoryIdx].ScripId)
							if lMasterArr.Data[lArrIdx].Symbol == pCategoryArr[lCategoryIdx].ScripId {
								log.Println("if -------------", lMasterRec)

								lMasterArr.Data[lArrIdx].CategoryDetailsArr = append(lMasterArr.Data[lArrIdx].CategoryDetailsArr, lCategoryRec)
								lMasterArr.Data[lArrIdx].SubCategoryDetailsArr = append(lMasterArr.Data[lArrIdx].SubCategoryDetailsArr, lSubCategoryRec)
								break
								// lMasterArr.Data = append(lMasterArr.Data, lMasterRec)
							} else if lArrIdx == len(lMasterArr.Data)-1 {
								log.Println("else -------------", lMasterRec)
								lMasterRec.CategoryDetailsArr = append(lMasterRec.CategoryDetailsArr, lCategoryRec)
								lMasterRec.SubCategoryDetailsArr = append(lMasterRec.SubCategoryDetailsArr, lSubCategoryRec)
								lMasterArr.Data = append(lMasterArr.Data, lMasterRec)
								break
							}
						}
					} else {
						lMasterRec.CategoryDetailsArr = append(lMasterRec.CategoryDetailsArr, lCategoryRec)
						lMasterRec.SubCategoryDetailsArr = append(lMasterRec.SubCategoryDetailsArr, lSubCategoryRec)
						lMasterArr.Data = append(lMasterArr.Data, lMasterRec)
					}
				}
			}
		}
	}
	log.Println("lMasterArr", lMasterArr)
	log.Println("ConstructIpoData (-)")
	return lMasterArr
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------

func FilterBseCategory(pCategoryArr []bseipo.BseIpoCategoryStruct) []bseipo.BseIpoCategoryStruct {
	log.Println("FilterBseCategory (+)")

	//this variable is used to get the final category array
	var lFinalCategoryArr []bseipo.BseIpoCategoryStruct

	for lCategoryIdx := 0; lCategoryIdx < len(pCategoryArr); lCategoryIdx++ {
		if pCategoryArr[lCategoryIdx].Type == "2" && pCategoryArr[lCategoryIdx].Series == "" && pCategoryArr[lCategoryIdx].Category == "IND" {
			lFinalCategoryArr = append(lFinalCategoryArr, pCategoryArr[lCategoryIdx])
		}
	}
	// log.Println("lFinalCategoryArr", lFinalCategoryArr)
	log.Println("FilterBseCategory (-)")
	return lFinalCategoryArr
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy branch
// --------------------------------------------------------------------

func FilterBseMaster(pMasterArr []bseipo.BseIpoMasterStruct) []bseipo.BseIpoMasterStruct {
	log.Println("FilterBseMaster (+)")

	//this variable is used to get the final master array
	var lFinalMasterArr []bseipo.BseIpoMasterStruct

	for lMasterIdx := 0; lMasterIdx < len(pMasterArr); lMasterIdx++ {
		if pMasterArr[lMasterIdx].AsbaNonAsba == "2" && pMasterArr[lMasterIdx].Category == "IND" && pMasterArr[lMasterIdx].IssueType == "BB" {
			lFinalMasterArr = append(lFinalMasterArr, pMasterArr[lMasterIdx])
		}
	}
	// log.Println("lFinalMasterArr", lFinalMasterArr)
	log.Println("FilterBseMaster (-)")
	return lFinalMasterArr
}
