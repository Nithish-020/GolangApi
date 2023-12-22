package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nseipo"
	"log"
)

type IpoDataStruct struct {
	Symbol string
	Isin   string
	Price  int
}
type RespCount struct {
	TotalCount   int `json:"totalCount"`
	SuccessCount int `json:"successCount"`
	ErrCount     int `json:"errCount"`
}

func NseIpoMktDemand(pUser string, pBrokerId int) (string, RespCount, error) {
	log.Println("NseIpoMktDemand (+)")
	lNoToken := common.ErrorCode
	var lResp RespCount

	lToken, lErr1 := GetToken(pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("ENIMD01", lErr1)
		return lNoToken, lResp, lErr1
	} else {
		lIpoDataArr, lErr2 := GetActiveIpos()
		if lErr2 != nil {
			log.Println("ENIMD01", lErr2)
			return lNoToken, lResp, lErr2
		} else {
			for _, lIpo := range lIpoDataArr {
				lResp.TotalCount++
				lMktDemad, lErr3 := nseipo.IpoMktDemand(lToken, pUser, lIpo.Symbol)
				if lErr3 != nil {
					log.Println("ENIMD01", lErr3)
					return lNoToken, lResp, lErr3
				} else {
					if lMktDemad.Status == common.SUCCESS {

						if lMktDemad.Symbol == lIpo.Symbol {

							lMktDemad.Isin = lIpo.Isin
							lMktDemad.Price = lIpo.Price

							lErr4 := CheckAndInsertMktDemand(lMktDemad)
							if lErr4 != nil {
								log.Println("ENIMD01", lErr4)
								return lNoToken, lResp, lErr4
							} else {
								lNoToken = common.SuccessCode
							}
						}
						lResp.SuccessCount++
					} else {
						lResp.ErrCount++
					}
				}
			}
		}
	}
	log.Println("NseIpoMktDemand (-)")
	return lNoToken, lResp, nil
}

func GetActiveIpos() ([]IpoDataStruct, error) {
	log.Println("GetActiveIpos (+)")
	var lIpoDataArr []IpoDataStruct
	var lIpoRec IpoDataStruct
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GSB01", lErr1)
		return lIpoDataArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `SELECT distinct NVL(TRIM(m.Symbol),''),NVL(m.Isin,''),m.MaxPrice
						FROM a_ipo_master m ,a_ipo_categories c,a_ipo_subcategory s
						WHERE m.Id= s.MasterId AND m.Id = c.MasterId AND m.IssueType = "EQUITY"
						AND c.code= "RETAIL" AND s.CaCode = "RETAIL"
						AND s.AllowUpi = 1 
						AND m.Exchange = 'NSE'
						AND curdate() BETWEEN m.BiddingStartDate AND m.BiddingEndDate 
						and m.Symbol not like '% %'`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("GSB02", lErr2)
			return lIpoDataArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lIpoRec.Symbol, &lIpoRec.Isin, &lIpoRec.Price)
				if lErr3 != nil {
					log.Println("GSB03", lErr3)
					return lIpoDataArr, lErr3
				} else {

					lIpoDataArr = append(lIpoDataArr, lIpoRec)
				}
			}
		}
	}
	log.Println("GetActiveIpos (-)")
	return lIpoDataArr, nil
}

func CheckAndInsertMktDemand(pMarketDemand nseipo.IpoMktDemandRespStruct) error {
	log.Println("CheckAndInsertMktDemand (+)")
	lIndicator := ""

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("CAIMD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for _, pData := range pMarketDemand.Data {

			if pData.CutOffIndicator == true {
				pData.Price = float32(pMarketDemand.Price)
			}

			lCoreString := `select (case when count(1) > 0 then "Y" else "N" end) indicator
							from a_ipo_mktdemand d
							where d.symbol = ?
							and d.isin = ?
							and d.Price = ?
							and d.CutOffIndicator = ?`

			lRows1, lErr2 := lDb.Query(lCoreString, pMarketDemand.Symbol, pMarketDemand.Isin, pData.Price, pData.CutOffIndicator)
			if lErr2 != nil {
				log.Println("CAIMD02", lErr2)
				return lErr2
			} else {
				//This for loop is used to collect the records from the database and store them in structure
				for lRows1.Next() {
					lErr3 := lRows1.Scan(&lIndicator)
					if lErr3 != nil {
						log.Println("CAIMD03", lErr3)
						return lErr3
					} else {
						if lIndicator == "Y" {
							lErr4 := UpdateMktDemandDetail(pMarketDemand, pData)
							if lErr4 != nil {
								log.Println("CAIMD04", lErr4)
								return lErr4
							}
						} else if lIndicator == "N" {
							lErr5 := InsertMktDemandDetail(pMarketDemand, pData)
							if lErr5 != nil {
								log.Println("CAIMD05", lErr5)
								return lErr5
							}
						}
					}
				}
			}
		}
	}
	log.Println("CheckAndInsertMktDemand (-)")
	return nil
}

func InsertMktDemandDetail(pResp nseipo.IpoMktDemandRespStruct, pData nseipo.IpoMktDemandStruct) error {
	log.Println("InsertMktDemandDetail (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("IMDD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pData.CutOffIndicator == true {
			pData.Price = float32(pResp.Price)
		}

		lSqlString := ` insert into a_ipo_mktdemand (Symbol,Isin,Price,CutOffIndicator,Series,absoluteQuantity,CumulativeQuantity,
						absoluteBidCount,CumulativeBidCount,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
						values( ?,?,?,?,?,?,?,?,?,'AUTOBOT',now(),'AUTOBOT',now())`

		_, lErr2 := lDb.Exec(lSqlString, pResp.Symbol, pResp.Isin, pData.Price, pData.CutOffIndicator, pData.Series, pData.AbsoluteQuantity, pData.CumulativeQuantity, pData.AbsoluteBidCount, pData.CumulativeBidCount)
		if lErr2 != nil {
			log.Println("IMDD02", lErr2)
			return lErr2
		}
	}

	log.Println("InsertMktDemandDetail (-)")
	return nil
}

func UpdateMktDemandDetail(pResp nseipo.IpoMktDemandRespStruct, pData nseipo.IpoMktDemandStruct) error {
	log.Println("UpdateMktDemandDetail (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("UMDD01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		if pData.CutOffIndicator == true {
			pData.Price = float32(pResp.Price)
		}
		lSqlString := ` update a_ipo_mktdemand set Price = ?,CutOffIndicator = ?,Series = ?,absoluteQuantity = ?,CumulativeQuantity = ?,
						absoluteBidCount = ?,CumulativeBidCount = ?,UpdatedBy = 'AUTOBOT',UpdatedDate = now()
						where Symbol = ? and Isin = ? and Price = ? and CutOffIndicator = ?`

		_, lErr2 := lDb.Exec(lSqlString, pData.Price, pData.CutOffIndicator, pData.Series, pData.AbsoluteQuantity, pData.CumulativeQuantity, pData.AbsoluteBidCount, pData.CumulativeBidCount, pResp.Symbol, pResp.Isin, pData.Price, pData.CutOffIndicator)
		if lErr2 != nil {
			log.Println("UMDD02", lErr2)
			return lErr2
		}
	}
	log.Println("UpdateMktDemandDetail (-)")
	return nil
}
