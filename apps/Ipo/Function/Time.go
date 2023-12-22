package Function

import (
	"fcs23pkg/ftdb"
	"log"
	"time"
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
func GetTime(pMasterId int) (string, error) {

	var lIndication string
	//----------------------------------------------------------------
	lCurrentTime := time.Now()
	lCurrentTimeUnix := lCurrentTime.Unix()

	//---------------------------
	lTmrwTen := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day()+1, 10, 0, 0, 0, lCurrentTime.Location())
	lTmrwTenUnix := lTmrwTen.Unix()

	//======================================
	lTodayFive := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), 17, 0, 0, 0, lCurrentTime.Location())
	lTodayFiveUnix := lTodayFive.Unix()

	//---------------------------
	lTodayTen := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), 10, 0, 0, 0, lCurrentTime.Location())
	lTodayTenUnix := lTodayTen.Unix()

	//--------------------------------------------------------------------
	lCurrentDayOfWeek := lCurrentTime.Weekday()
	lDayOfWeekString := lCurrentDayOfWeek.String()

	//-----------------------------------------------------------------------

	if lDayOfWeekString == "Saturday" || lDayOfWeekString == "Sunday" {
		lIndication = "True"
	} else {
		if lCurrentTimeUnix > lTodayFiveUnix && lCurrentTimeUnix < lTmrwTenUnix {
			lIndication = "True"

			lGoodTime, lErr1 := GoodTimeToApply(pMasterId)
			if lErr1 != nil {
				log.Println("FGT01", lErr1)
				return lIndication, lErr1
			} else {
				if lGoodTime == "True" {
					lIndication = "False"
				} else if lGoodTime == "False" {
					lIndication = "True"
				}
			}

		} else if lCurrentTimeUnix > lTodayTenUnix && lCurrentTimeUnix < lTodayFiveUnix {
			lIndication = "False"
		} else {
			lIndication = "True"
		}
	}
	return lIndication, nil
}

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
func GoodTimeToApply(pMasterId int) (string, error) {
	log.Println("GoodTimeToApply (+)")

	// this variable is used to get the flag from the database
	var lFlag string
	// get the id from database
	var lId int
	// return the indicator
	var lIndicator string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("FGTA01", lErr1)
		return lIndicator, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select m.Id,(case when (time(now()) between c.StartTime and c.EndTime) then 'Y' else 
						(
						(case when (c.StartTime is null and c.EndTime is null) then 
						(case when (time(now()) between m.DailyStartTime and m.DailyEndTime) then 'I' else 'N' end) else 'S' end)
						) end) value
						from a_ipo_master m,a_ipo_categories c
						where m.Id = c.MasterId
						and c.Code = 'RETAIL'
						and m.Id = ? `
		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
		if lErr2 != nil {
			log.Println("FGTA02", lErr2)
			return lIndicator, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lId, &lFlag)
				if lErr3 != nil {
					log.Println("FGTA03", lErr3)
					return lIndicator, lErr3
				} else {
					if pMasterId == lId {
						if lFlag == "Y" || lFlag == "I" {
							lIndicator = "True"
						} else if lFlag == "N" || lFlag == "S" {
							lIndicator = "False"
						}
					}
				}
			}
		}
	}
	log.Println("GoodTimeToApply (-)")
	return lIndicator, nil
}

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
func CheckEndDate(pMasterId int) (string, error) {
	log.Println("CheckEndDate (+)")

	lCurrentTime := time.Now()

	// Convert time to string using Format method
	// lDateString := lCurrentTime.Format("2006-01-02")
	lUpdatedTime := lCurrentTime.Add(5 * time.Minute)
	lTimeString := lUpdatedTime.Format("15:04:05")

	log.Println(lTimeString)
	// get the id from database
	var lId int
	// this variable is used to get the flag from the database
	var lFlag string
	// return the indicator
	var lIndicator string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("FCED01", lErr1)
		return lIndicator, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `				
		select (case when count(1) > 0 then m1.Id else 0 end) id,(case when a.endTime > '` + lTimeString + `' then 'T' else 'F' end) timevalue
		from (
		select nvl(c.EndTime, m.DailyEndTime) endTime, c.MasterId Id
		from a_ipo_master m,a_ipo_categories c
		where m.Id = c.MasterId
		and c.Code = 'RETAIL'
		and m.BiddingEndDate = date(now())
		) a, a_ipo_master m1
		where m1.Id = a.Id
		and m1.id = ?`

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
					log.Println("pMasterId", pMasterId)
					log.Println("lId", lId)
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
	log.Println("CheckEndDate (-)")
	return lIndicator, nil
}
