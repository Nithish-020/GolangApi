package brokers

import (
	"fcs23pkg/ftdb"
	"log"
)

func InsertCookieDetails(pCookieValue string, pBrokerId int, pCliendId string) (string, error) {
	log.Println("InsertCookieDetails (+)")
	var lFlag string
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BICD01", lErr1)
		return lFlag, lErr1
	} else {
		defer lDb.Close()

		// To get MemberDetailsId
		lMemberDetailId, lErr2 := GetMemberDetailId(pBrokerId)
		if lErr2 != nil {
			log.Println("BICD01", lErr2)
			return lFlag, lErr2
		} else {
			lSqlString := `insert into a_ipo_cookiedetails (BrokerId ,MemberDetailId,CookieValue,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
				values (?,?,?,?,now(),?,now())
				`
			_, lErr4 := lDb.Exec(lSqlString, pBrokerId, lMemberDetailId, pCookieValue, pCliendId, pCliendId)
			if lErr4 != nil {
				log.Println("BAB03", lErr4)
				return lFlag, lErr4
			} else {
				if lMemberDetailId != 0 {
					lFlag = "Y"
				} else {
					lFlag = "N"
				}
			}
		}
	}
	log.Println("InsertCookieDetails (-)")
	return lFlag, nil
}

func GetMemberDetailId(pBrokerId int) (int, error) {
	log.Println("GetMemberDetailId (+)")
	lIndicator := 0
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BICD01", lErr1)
		return lIndicator, lErr1
	} else {
		defer lDb.Close()
		// ! To get 0 if MemberDetailsId is not present,If present get MemberDetailsId
		lCoreString2 := `select (case when 1 = count(1) then cd.Id  else 0 end )
		from a_ipo_memberdetails cd
		where cd.BrokerId  = ? and cd.Flag = 'Y'`

		lRows1, lErr2 := lDb.Query(lCoreString2, pBrokerId)
		if lErr2 != nil {
			log.Println("AGBI02", lErr2)
			return lIndicator, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in a variable
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lIndicator)
				if lErr3 != nil {
					log.Println("BGBL05", lErr3)
					return lIndicator, lErr3
				}
			}
		}
	}
	log.Println("GetMemberDetailId (-)")
	return lIndicator, nil
}
