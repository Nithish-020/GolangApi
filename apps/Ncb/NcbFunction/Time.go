package NcbFunction

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

Author: KAVYADHARSHANI
Date: 23OCT2023
*/
// func NcbGetTime(pMasterId int) (string, error) {

// 	var lIndication string
// 	//----------------------------------------------------------------
// 	lCurrentTime := time.Now()
// 	lCurrentTimeUnix := lCurrentTime.Unix()

// 	//---------------------------
// 	lTmrwTen := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day()+1, 10, 0, 0, 0, lCurrentTime.Location())
// 	lTmrwTenUnix := lTmrwTen.Unix()

// 	//======================================
// 	lTodayFive := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), 17, 0, 0, 0, lCurrentTime.Location())
// 	lTodayFiveUnix := lTodayFive.Unix()

// 	//---------------------------
// 	lTodayTen := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), 10, 0, 0, 0, lCurrentTime.Location())
// 	lTodayTenUnix := lTodayTen.Unix()

// 	//--------------------------------------------------------------------
// 	lCurrentDayOfWeek := lCurrentTime.Weekday()
// 	lDayOfWeekString := lCurrentDayOfWeek.String()

// 	//-----------------------------------------------------------------------

// 	if lDayOfWeekString == "Saturday" || lDayOfWeekString == "Sunday" {
// 		lIndication = "True"
// 	} else {
// 		if lCurrentTimeUnix > lTodayFiveUnix && lCurrentTimeUnix < lTmrwTenUnix {
// 			lIndication = "True"

// 			lGoodTime, lErr1 := NsbGoodTimeToApply(pMasterId)
// 			if lErr1 != nil {
// 				log.Println("FNGT01", lErr1)
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

Author: KAVYADHARSHANI
Date: 23OCT2023
*/
// func NsbGoodTimeToApply(pMasterId int) (string, error) {
// 	log.Println("NsbGoodTimeToApply (+)")

// 	// this variable is used to get the flag from the database
// 	var lFlag string
// 	// get the id from database
// 	var lId int
// 	// return the indicator
// 	var lIndicator string

// 	var lSymbol string
// 	// To Establish A database connection,call LocalDbConnect Method
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("FNGTA01", lErr1)
// 		return lIndicator, lErr1
// 	} else {
// 		defer lDb.Close()
// 		lCoreString := ` select
// 		                     (case when count(1) > 0 then  'Y' else 'N' end ) flag,  COALESCE(n.id, -1) AS id, nvl(n.Symbol,"")
// 		                 from
// 					        a_ncb_master n
// 			            where
// 				          time(now()) between n.DailyStartTime and n.DailyEndTime
// 						  and Date(now()) between  n.BiddingStartDate  and n.BiddingEndDate
// 						  AND n.id IS NOT NULL
// 						  and n.id = ?`
// 		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
// 		if lErr2 != nil {
// 			log.Println("FNGTA02", lErr2)
// 			return lIndicator, lErr2
// 		} else {
// 			for lRows.Next() {
// 				lErr3 := lRows.Scan(&lFlag, &lId, &lSymbol)
// 				if lErr3 != nil {
// 					log.Println("FNGTA03", lErr3)
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
// 		log.Println("lIndicator", lIndicator)
// 	}
// 	log.Println("NsbGoodTimeToApply (-)")
// 	return lIndicator, nil
// }
