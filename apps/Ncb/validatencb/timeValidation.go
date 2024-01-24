package validatencb

import (
	"fcs23pkg/ftdb"
	"log"
)

/*
Pupose: This method is used to get the ncb endtime
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
Date: 12OCT2023
*/
func CheckNcbEndDate(pMasterId int) (string, error) {
	log.Println("CheckNcbEndDate (+)")

	var lId int
	// this variable is used to get the flag from the database
	var lFlag string
	// return the indicator
	var lIndicator string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("TCNED01", lErr1)
		return lIndicator, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select (case when count(1) > 0 then n.id else 0 end) id, (case when n.DailyEndTime > '15:30:00' then 'T' else 'F' end) flag  
		                from a_ncb_master n
		                where n.BiddingEndDate = date(now())
		                and n.id = ?`

		lRows, lErr2 := lDb.Query(lCoreString, pMasterId)
		if lErr2 != nil {
			log.Println("TCNED02", lErr2)
			return lIndicator, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lId, &lFlag)
				if lErr3 != nil {
					log.Println("TCNED03", lErr3)
					return lIndicator, lErr3
				} else {
					log.Println("pMasterId", pMasterId)
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
		log.Println("lIndicator", lIndicator)
	}
	log.Println("CheckNcbEndDate (-)")
	return lIndicator, nil
}
