package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nseipo"
	"log"
)

func NseIpoMktCatwise(pUser string, pBrokerId int) (string, RespCount, error) {
	log.Println("NseIpoMktData (+)")
	lNoToken := common.ErrorCode
	var lResp RespCount
	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("ENIMD01", lErr1)
		return lNoToken, lResp, lErr1
	} else {
		lIpoDataArr, lErr2 := GetActiveIpos()
		log.Println("lIpoDataArr :=", lIpoDataArr)
		if lErr2 != nil {
			log.Println("ENIMD01", lErr2)
			return lNoToken, lResp, lErr2
		} else {
			for _, lIpo := range lIpoDataArr {
				lResp.TotalCount++
				lMktDemad, lErr3 := nseipo.IpoMktCatwise(lToken, pUser, lIpo.Symbol)
				if lErr3 != nil {
					log.Println("ENIMD01", lErr3)
					return lNoToken, lResp, lErr3
				} else {
					if lMktDemad.Status == common.SUCCESS {
						lResp.SuccessCount++
						if lMktDemad.Symbol == lIpo.Symbol {

							lMktDemad.Isin = lIpo.Isin
							lErr4 := CheckAndInsertMktCatwise(lMktDemad)
							if lErr4 != nil {
								log.Println("ENIMD01", lErr4)
								return lNoToken, lResp, lErr4
							} else {
								lNoToken = common.SuccessCode
							}
						}
					} else {
						lResp.ErrCount++
					}
				}
			}
		}
	}
	log.Println("NseIpoMktData (-)")
	return lNoToken, lResp, nil
}

func CheckAndInsertMktCatwise(pMarketDemand nseipo.IpoMktCatwiseRespStruct) error {
	log.Println("CheckAndInsertMktCatwise (+)")
	lIndicator := ""

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("CAIMC01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for _, lData := range pMarketDemand.Data {

			lCoreString := `select (case when count(1) = 1 then "Y" else "N" end) indicator
						from a_ipo_mktcatwise c
						where c.symbol = ?
						and c.isin = ?
						and c.Category = ?
						and c.Subcategory = ?`

			lRows1, lErr2 := lDb.Query(lCoreString, pMarketDemand.Symbol, pMarketDemand.Isin, lData.Category, lData.Subcategory)
			if lErr2 != nil {
				log.Println("CAIMC02", lErr2)
				return lErr2
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows1.Next() {
					lErr3 := lRows1.Scan(&lIndicator)
					if lErr3 != nil {
						log.Println("CAIMC03", lErr3)
						return lErr3
					} else {
						if lIndicator == "Y" {
							lErr4 := UpdateMktCatwiseDetail(pMarketDemand, lData)
							if lErr4 != nil {
								log.Println("CAIMC04", lErr4)
								return lErr4
							}
						} else {
							lErr5 := InsertMktCatwiseDetail(pMarketDemand, lData)
							if lErr5 != nil {
								log.Println("CAIMC05", lErr5)
								return lErr5
							}
						}
					}
				}
			}
		}
	}
	log.Println("CheckAndInsertMktCatwise (-)")
	return nil
}

func InsertMktCatwiseDetail(pResp nseipo.IpoMktCatwiseRespStruct, pData nseipo.IpoMktCatwiseStruct) error {
	log.Println("InsertMktCatwiseDetail (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IMCD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lSqlString := ` insert into a_ipo_mktcatwise (Symbol,Isin,Category,SubCategory,Quantity,BidCount,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
						values( ?,?,?,?,?,?,'AUTOBOT',now(),'AUTOBOT',now())`

		_, lErr2 := lDb.Exec(lSqlString, pResp.Symbol, pResp.Isin, pData.Category, pData.Subcategory, pData.Quantity, pData.BidCount)
		if lErr2 != nil {
			log.Println("IMCD02", lErr2)
			return lErr2
		}
	}

	log.Println("InsertMktCatwiseDetail (-)")
	return nil
}

func UpdateMktCatwiseDetail(pResp nseipo.IpoMktCatwiseRespStruct, pData nseipo.IpoMktCatwiseStruct) error {
	log.Println("UpdateMktCatwiseDetail (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("UMCD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lSqlString := ` update a_ipo_mktcatwise set Category = ?,Subcategory = ?,Quantity = ?,BidCount = ?,UpdatedBy = 'AUTOBOT',UpdatedDate = now()
						where Symbol = ? and Category = ? and Subcategory = ?`

		_, lErr2 := lDb.Exec(lSqlString, pData.Category, pData.Subcategory, pData.Quantity, pData.BidCount, pResp.Symbol, pData.Category, pData.Subcategory)
		if lErr2 != nil {
			log.Println("UMCD02", lErr2)
			return lErr2
		}
	}
	log.Println("UpdateMktCatwiseDetail (-)")
	return nil
}
