package validatesgb

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
)

/*
Pupose: This method is used to get the Currentime is available for placing the order
Parameters:

	pMasterId

Response:

	    ==========
	    *On Sucess
	    ==========
		True,nil or False,nil
	    ==========
	    !On Error
	    ==========
		"",error

Author: Pavithra
Date: 07JUNE2023
*/
// func GetSgbTime(pMasterId int) (string, error) {
// 	log.Println("GetSgbTime (+)")

// 	var lIndication string
// 	//----------------------------------------------------------------
// 	lCurrentTime := time.Now()
// 	lCurrentTimeUnix := lCurrentTime.Unix()

// 	//---------------------------
// 	lTmrwTen := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day()+1, 10, 0, 0, 0, lCurrentTime.Location())
// 	lTmrwTenUnix := lTmrwTen.Unix()

// 	//======================================
// 	lTodayFive := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), 20, 0, 0, 0, lCurrentTime.Location())
// 	lTodayFiveUnix := lTodayFive.Unix()

// 	//---------------------------
// 	lTodayTen := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), 10, 0, 0, 0, lCurrentTime.Location())
// 	lTodayTenUnix := lTodayTen.Unix()

// 	//--------------------------------------------------------------------
// 	lCurrentDayOfWeek := lCurrentTime.Weekday()
// 	lDayOfWeekString := lCurrentDayOfWeek.String()

// 	//-----------------------------------------------------------------------

// 	if lDayOfWeekString == "Saturday" || lDayOfWeekString == "Sunday" {
// 		lIndication = "False"
// 	} else {
// 		if lCurrentTimeUnix > lTodayFiveUnix && lCurrentTimeUnix < lTmrwTenUnix {
// 			lIndication = "True"
// 			//to check the placed SGB current time is available or not
// 			lGoodTime, lErr1 := GoodTimeToApplySGB(pMasterId)
// 			if lErr1 != nil {
// 				log.Println("FGT01", lErr1)
// 				return lIndication, lErr1
// 			} else {
// 				if lGoodTime == "True" {
// 					lIndication = "False"
// 				} else if lGoodTime == "False" {
// 					lIndication = "True"
// 				}
// 			}
// 		} else if lCurrentTimeUnix > lTodayTenUnix && lCurrentTimeUnix < lTodayFiveUnix {
// 			lIndication = "False"
// 		} else {
// 			lIndication = "True"
// 		}
// 	}
// 	log.Println("GetSgbTime (+)")
// 	return lIndication, nil
// }

/*
Pupose: This method is used to get the ipo endtime is goodtime for applying
Parameters:

	pMasterId

Response:

	    ==========
	    *On Sucess
	    ==========
		True,nil or False,nil
	    ==========
	    !On Error
	    ==========
		"",error

Author: Pavithra
Date: 07JUNE2023
*/
// func GoodTimeToApplySGB(pMasterId int) (string, error) {
// 	log.Println("GoodTimeToApplySGB (+)")

// 	// this variable is used to get the flag from the database
// 	var lFlag string
// 	// get the id from database
// 	var lId int
// 	// return the indicator
// 	var lIndicator string

// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("", lErr1)
// 		return lIndicator, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := `select (case when count(1) > 0 then  'Y' else 'N' end ) a,id
// 						from a_sgb_master asm
// 						where  time(now()) between DailyStartTime and DailyEndTime
// 						and Date(now()) between  BiddingStartDate  and BiddingEndDate
// 						and id = ? `
// 		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
// 		if lErr2 != nil {
// 			log.Println("FGTA02", lErr2)
// 			return lIndicator, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lFlag, &lId)
// 				if lErr3 != nil {
// 					log.Println("", lErr3)
// 					return lIndicator, lErr3
// 				} else {
// 					if pMasterId == lId {
// 						if lFlag == "Y" {
// 							lIndicator = "True"
// 						} else if lFlag == "N" {
// 							lIndicator = "False"
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("GoodTimeToApplySGB (-)")
// 	return lIndicator, nil
// }

/*
Pupose: This method is used to get the ipo endtime
Parameters:

	pMasterId

Response:

	    ==========
	    *On Sucess
	    ==========
		True,nil or False,nil
	    ==========
	    !On Error
	    ==========
		"",error

Author: Pavithra
Date: 29JULY2023
*/
func CheckSgbEndDate(pMasterId int) (string, error) {
	log.Println("CheckSgbEndDate (+)")

	// lCurrentTime := time.Now()

	// Convert time to string using Format method
	// lDateString := lCurrentTime.Format("2006-01-02")
	// lUpdatedTime := lCurrentTime.Add(5 * time.Minute)
	// lTimeString := lUpdatedTime.Format("15:04:05")

	// log.Println(lTimeString)
	// get the id from database
	var lId int
	// this variable is used to get the flag from the database
	var lFlag string
	// return the indicator
	var lIndicator string

	lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("FCED01", lErr1)
		return lIndicator, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select (case when count(1) > 0 then m.id else 0 end) id ,(case when m.DailyEndTime > '` + lCloseTime + `' then 'T' else 'F' end) flag
						from a_sgb_master m
						where m.BiddingEndDate = date(now())
						and m.id = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
		if lErr2 != nil {
			log.Println("FCED02", lErr2)
			return lIndicator, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lId, &lFlag)
				if lErr3 != nil {
					log.Println("FCED03", lErr3)
					return lIndicator, lErr3
				} else {
					// log.Println("pMasterId", pMasterId)
					// log.Println("lId", lId)
					if lId == 0 && lFlag == "F" {
						lIndicator = "True"
					} else if lId == pMasterId && lFlag == "T" {
						lIndicator = "True"
					} else if lId == pMasterId && lFlag == "F" {
						lIndicator = "False"
					}
				}
			}
		}
	}
	log.Println("CheckSgbEndDate (-)")
	return lIndicator, nil
}
