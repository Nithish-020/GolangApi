package memberdetail

import (
	"fcs23pkg/apps/validation/menu"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
)

type SegmentDetailsRec struct {
	Segments string `json:"segments"`
	Nse      bool   `json:"nse"`
	NseShare int    `json:"nseShare"`
	Bse      bool   `json:"bse"`
	BseShare int    `json:"bseShare"`
}

//----------------------------------------------------------------------------
// This method will return the Flag of Client must need to be for NSE or BSE
//----------------------------------------------------------------------------

func BseNsePercentCalc(pBrokerId int, pSegment string) (string, error) {
	log.Println("BseNseShareCalculator (+)")

	var lPercentage float64
	var lExchangeFlag string

	lNse, lBse, lErr1 := BrokersShareDetails(pBrokerId, pSegment)
	if lErr1 != nil {
		log.Println("MBNPC01", lErr1)
		return lExchangeFlag, lErr1
	} else {
		// if percantage is Fully Nse
		if lNse == 100 && lBse == 0 {
			lExchangeFlag = common.NSE
			// Percentage Fully Bse
		} else if lBse == 100 && lNse == 0 {
			lExchangeFlag = common.BSE
			// if Nse and Bse Both Shares Are Equal  it will call Exchange Vice versa
		} else if lBse == lNse {
			if lExchangeFlag == common.BSE {
				lExchangeFlag = common.NSE
			} else {
				lExchangeFlag = common.BSE
			}

		} else {
			lCount, lErr2 := clientOrderCount(pSegment, pBrokerId)
			if lErr2 != nil {
				log.Println("MBNPC02", lErr2)
				return lExchangeFlag, lErr2
			} else {
				// lcount =0 is Done only For the First Client in the Bilateral Percentage of Nse and Bse
				if lCount == 0 {
					if lNse > lBse {
						lExchangeFlag = common.NSE
					} else if lBse > lNse {
						lExchangeFlag = common.BSE
						// Call BSE method For First Client
					}
				} else if lCount > 0 {
					lExchange, lErr3 := MaxofOrderinExchange(pBrokerId)
					if lErr3 != nil {
						log.Println("MBNPC03", lErr3)
						return lExchangeFlag, lErr3
					} else {
						//  if First client lExchange is Nse it will Call Bes Like Viceversa
						if lExchange == common.NSE {
							lExchangeCount, lErr4 := CountPlacedInExchange(pBrokerId, common.BSE)
							if lErr4 != nil {
								log.Println("MBNPC04", lErr4)
								return lExchangeFlag, lErr4
							} else {
								//  Based on Exchage Count the Percentage Will be calculated
								lPercentage = RatioFinder(lBse, lExchangeCount)
								// IF Bse Share is less than Zero Check For Nse Share
								// (OR) Bse share is greater than Zero it will go for else if
								if lPercentage <= 0 {
									lPercentage = RatioFinder(lNse, lExchangeCount)
									if lPercentage > 0 {
										lExchangeFlag = common.NSE
									} else {
										//  If  NSE and BSE Both the percentages are zero
										// it will restart percentage Calculation From Begining
										lPercentage = restartRatioCounter(lBse, lExchangeCount)
										if lPercentage > 0 {
											lExchangeFlag = common.BSE
										} else {
											lExchangeFlag = common.NSE
										}
									}
								} else if lPercentage > 0 {
									lExchangeFlag = common.BSE
								}
							}
						} else if lExchangeFlag == common.BSE {
							//  Else part is Fully For NSE
							lExchangeCount, lErr5 := CountPlacedInExchange(pBrokerId, common.NSE)
							if lErr5 != nil {
								log.Println("MBNPC05", lErr5)
								return lExchangeFlag, lErr5
							} else {
								lPercentage = RatioFinder(lNse, lExchangeCount)
								//  if NSE share is greater than Zero it will go for else if
								// IF NSE Share is Zero Check For BSE Share
								if lPercentage <= 0 {
									lPercentage = RatioFinder(lBse, lExchangeCount)
									if lPercentage > 0 {
										lExchangeFlag = common.BSE
									} else {
										//  if  NSE and BSE Both the percentages are zero
										// it will restart percentage Calculation From Begining
										lPercentage = restartRatioCounter(lNse, lExchangeCount)
										if lPercentage > 0 {
											lExchangeFlag = common.NSE
										} else {
											lExchangeFlag = common.BSE
										}
									}
								} else if lPercentage > 0 {
									lExchangeFlag = common.NSE
								}
							}
						}
					}
				}
			}

		}
	}
	log.Println("BseNseShareCalculator (-)")
	return lExchangeFlag, nil

}

//----------------------------------------------------------------------------
//  This Get Segments method Return The Segment Percentage as Per pBroker Id to display in the UI
//----------------------------------------------------------------------------

func GetSegments(pBrokerId int) ([]SegmentDetailsRec, error) {
	log.Println("getSegments (+)")
	var lSegmentData SegmentDetailsRec
	var lSegmentsArr []SegmentDetailsRec
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MGS01", lErr1)
		return lSegmentsArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString2 := `SELECT UPPER(R.RouterName), NSE, BSE
		FROM a_broker_segments BS,a_ipo_brokermaster BM ,a_ipo_router R
		WHERE BS.BrokerId =BM.Id and R.Id = BS.Segments and bs.BrokerId = ? and BS.Type = 'Y'`

		lRows1, lErr2 := lDb.Query(lCoreString2, pBrokerId)
		if lErr2 != nil {
			log.Println("MGS02", lErr2)
			return lSegmentsArr, lErr2
		} else {

			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {

				lErr3 := lRows1.Scan(&lSegmentData.Segments, &lSegmentData.NseShare, &lSegmentData.BseShare)
				if lErr3 != nil {
					log.Println("MGS03", lErr3)
					return lSegmentsArr, lErr3
				} else {
					lSegmentData.Nse = false
					lSegmentData.Bse = false
					if lSegmentData.NseShare > 0 {
						lSegmentData.Nse = true
					}
					if lSegmentData.BseShare > 0 {
						lSegmentData.Bse = true
					}
					lSegmentData.Segments = common.CapitalizeText(lSegmentData.Segments)
					lSegmentsArr = append(lSegmentsArr, lSegmentData)
				}
			}
		}
	}
	log.Println("getSegments (-)")
	return lSegmentsArr, nil
}

// This Method is used to Order The Structure in Array if Segment is Added TempArr Segment Also needs To be Addded
func ConstructArrayOrder(pSegmentArr []SegmentDetailsRec, pBrokerId int) ([]SegmentDetailsRec, error) {
	log.Println("ConstructArrayOrder (+)")
	var lOrgainzedDetailsArr []SegmentDetailsRec
	lAllowModules, lErr3 := menu.GetAllowModules(pBrokerId)
	if lErr3 != nil {
		log.Println("MCAO03", lErr3)
		return lOrgainzedDetailsArr, lErr3
	} else {
		for _, Temp := range lAllowModules {
			Temp = common.CapitalizeText(Temp)
			for _, Segments := range pSegmentArr {
				if Segments.Segments == Temp {
					lOrgainzedDetailsArr = append(lOrgainzedDetailsArr, Segments)
				}
			}
		}
	}
	log.Println("ConstructArrayOrder (-)")
	return lOrgainzedDetailsArr, nil
}

// this Method Will Insert Segment in Array Format if Data Alredy Exist it will NOt Insert
func InsertSegments(pBrokerId int, pSegmentToUpdate []SegmentDetailsRec) error {
	log.Println("InsertSegments (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MIS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		log.Println("Insert", pBrokerId, pSegmentToUpdate)
		for _, Segments := range pSegmentToUpdate {

			lSqlString1 := `INSERT INTO a_broker_segments (Segments, NSE, BSE, BrokerId,Type,CreatedBy,CreatedDate,UpdatedBy,UpdatedDate)
			select(select air.id	from a_ipo_router air where RouterName = ?),?,?,?,"Y",?,now(),?,now()
			where not exists (select * from a_broker_segments where Segments = (select air.id	from a_ipo_router air where RouterName = ?)
			and BrokerId = ? and Type = 'Y')
			limit 1
			`
			_, lErr2 := lDb.Exec(lSqlString1, Segments.Segments, Segments.NseShare, Segments.BseShare, pBrokerId, pBrokerId, pBrokerId, Segments.Segments, pBrokerId)
			if lErr2 != nil {
				log.Println("MIS02", lErr2)
				return lErr2
			}
		}

	}

	log.Println("InsertSegments (-)")
	return nil
}

// This Method Will insert an data if Segment Data is empty for Broker id it will Directly Insert the data
// if already Data is Present in Segments and It will call For insert and update method
func InsertAndUpdateSegments(pBrokerId int, pMemberDetail memberDetailStruct) ([]string, error) {
	log.Println("InsertAndUpdateSegments (+)")
	var lNoAcceptSegments []string
	lSegmentsArr, lErr1 := GetSegments(pBrokerId)
	if lErr1 != nil {
		log.Println("MIAUS01")
		return lNoAcceptSegments, lErr1
	} else {
		log.Println(" pMemberDetail memberDetailStruct", pMemberDetail.SegmentShares)
		var lTempSegments []SegmentDetailsRec
		// This for loop will Remove the value of BSE and NSE Both the Share are 0 as per segments
		for _, InsertSegments := range pMemberDetail.SegmentShares {
			if InsertSegments.BseShare != 0 || InsertSegments.NseShare != 0 {
				if (InsertSegments.NseShare + InsertSegments.BseShare) == 100 {

					lTempSegments = append(lTempSegments, InsertSegments)
				} else {
					lNoAcceptSegments = append(lNoAcceptSegments, InsertSegments.Segments)
				}
			}
		}

		log.Println("lTempSegments", lTempSegments)
		if lSegmentsArr == nil {
			lErr1 := InsertSegments(pBrokerId, lTempSegments)
			if lErr1 != nil {
				log.Println("MIAUS02")
				return lNoAcceptSegments, lErr1
			}

		} else {
			lErr2 := InsertSegments(pBrokerId, lTempSegments)
			if lErr2 != nil {
				log.Println("MIAUS03", lErr2)
				return lNoAcceptSegments, lErr2
			} else {
				lErr3 := updateSegments(pBrokerId, lTempSegments)
				if lErr3 != nil {
					log.Println("MIAUS04")
					return lNoAcceptSegments, lErr3
				}

			}
		}
	}

	log.Println("InsertAndUpdateSegments (+)")
	return lNoAcceptSegments, nil
}

// This Method Will Update Data in Array based on SegmentID ,BrokerId and Type
func updateSegments(pBrokerId int, pSegmentToUpdate []SegmentDetailsRec) error {
	log.Println("updateSegments (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MUS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		log.Println("pSegmentToUpdate", pSegmentToUpdate)
		for _, Update := range pSegmentToUpdate {
			lSqlString1 := `UPDATE a_broker_segments
			SET  NSE= ?, BSE= ?,UpdatedBy= ?, UpdatedDate= now()
			WHERE Segments = (select air.id	from a_ipo_router air where RouterName = ?) and BrokerId = ? and Type = 'Y'`
			_, lErr2 := lDb.Exec(lSqlString1, Update.NseShare, Update.BseShare, pBrokerId, Update.Segments, pBrokerId)
			if lErr2 != nil {
				log.Println("MUS02", lErr2)
				return lErr2

			}

		}
	}

	log.Println("updateSegments (-)")
	return nil
}

// this Method Check as Segment Details Present or Not As per Allow Modules
// if Data is Not Equal to allow Modules This Method Will call to Disable Segments
func UpdateSegmentType(pBrokerId int) error {
	log.Println("UpdateSegmentType(+)")
	lSegmentsArr, lErr1 := GetSegments(pBrokerId)
	if lErr1 != nil {
		log.Println("MUST01")
		return lErr1
	} else {
		if lSegmentsArr != nil {
			lAllowModules, lErr2 := menu.GetAllowModules(pBrokerId)
			if lErr2 != nil {
				log.Println("MUST02", lErr2)
				return lErr2
			} else {
				var filteredArray []string
				// Create a map to store the elements of array2 for faster lookup
				//  In This For We Need To add a Array Which need To Mandatory Data
				elementsInAllowModules := make(map[string]bool)
				for _, Modules := range lAllowModules {
					Modules = common.CapitalizeText(Modules)
					elementsInAllowModules[Modules] = true
				}
				//  In this For We can Remove False Vale Which Means Data Not present in UpperLoop
				for _, Segments := range lSegmentsArr {
					if !elementsInAllowModules[Segments.Segments] {
						filteredArray = append(filteredArray, Segments.Segments)
					}
				}
				lErr3 := diableSegments(pBrokerId, filteredArray)
				if lErr3 != nil {
					log.Println("MUST03", lErr3)
					return lErr3
				}
			}
		}
	}
	log.Println("UpdateSegmentType(-)")
	return nil
}

// THIS  Method will Disable The Segments As per allow Modules
func diableSegments(pBrokerId int, pSegmentToUpdate []string) error {
	log.Println("diableSegments (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MDS01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for _, Update := range pSegmentToUpdate {
			lSqlString1 := `UPDATE a_broker_segments
			SET  Type= 'N',UpdatedBy= ?, UpdatedDate= now()
			WHERE Segments = (select air.id	from a_ipo_router air where RouterName = ?) and BrokerId = ? and Type = 'Y'`
			// log.Println("UpdateData", pBrokerId, Update.Segments, pBrokerId)
			_, lErr2 := lDb.Exec(lSqlString1, pBrokerId, Update, pBrokerId)
			if lErr2 != nil {
				log.Println("MDS02", lErr2)
				return lErr2
			}
		}
	}

	log.Println("diableSegments (-)")
	return nil
}

// This Method is used to Return The Share as per BrokerId and RouterPath
func BrokersShareDetails(pBrokerId int, pRouter string) (float64, float64, error) {
	log.Println("BrokersShareDetails (+)")
	var lNse float64
	var lBse float64
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MBSD01", lErr1)
		return lNse, lBse, lErr1
	} else {
		defer lDb.Close()
		//  Based on Router Path Query will Change Dynamically

		lCoreString2 := `select  NSE, BSE
		from a_broker_segments BS ,a_ipo_router R
		where  BS.BrokerId  = ? and  R.Id = BS.Segments and R.RouterName = ? and BS.type = 'Y' `

		lRows1, lErr2 := lDb.Query(lCoreString2, pBrokerId, pRouter)
		if lErr2 != nil {
			log.Println("MBSD02", lErr2)
			return lNse, lBse, lErr2
		} else {
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lNse, lBse)
				if lErr3 != nil {
					log.Println("MBSD03", lErr3)
					return lNse, lBse, lErr3
				}
			}
		}
	}
	log.Println("BrokersShareDetails (-)")
	return lNse, lBse, nil
}

// This Method Return The count Done on Day as per Router And BrokerId
// we will not count client Based on Exchange only Based on path and how many Records in a Table
//
//	Based on Record We will seperate an Client
func clientOrderCount(pSegment string, pBrokerId int) (int, error) {
	log.Println("clientOrderCount (+)")
	var lCount int
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MCOC01", lErr1)
		return lCount, lErr1
	} else {
		defer lDb.Close()

		lCoreString1 := `select count(sh.Id)
		from a_sgb_orderheader SH
		where  SH.brokerId = ? and sh.CreatedDate  = now()`

		lCoreString2 := `select  count(OH.applicationNo)
		from a_ipo_order_header OH
		where oh.CreatedDate = now() and OH.brokerId = ?  `

		lCoreString := ""
		if pSegment == "IPO" {
			lCoreString = lCoreString2
		} else if pSegment == "SGB" {
			lCoreString = lCoreString1
		}

		lRows1, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("MCOC02", lErr2)
			return lCount, lErr2
		} else {
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lCount)
				if lErr3 != nil {
					log.Println("MCOC03", lErr3)
					return lCount, lErr3
				}
			}
		}
		fmt.Println("lCount", lCount)
	}
	log.Println("clientOrderCount (-)")
	return lCount, nil

}

// This Method Will Split the Percentage as per client order
// if percentage Share is Zero
func RatioFinder(pShare float64, pcount int) float64 {
	log.Println("RatioFinder (+)")
	//  Constructing in Binary Tree Format
	var lPercentage float64
	if pcount%2 == 1 {
		lPercentage = pShare / float64(pcount+1)
	} else {
		lPercentage = pShare
	}
	log.Println("RatioFinder (-)")

	return lPercentage
}

// This method will restart the count by getting Remainder value of count
func restartRatioCounter(pShare float64, pcount int) float64 {
	log.Println("restartRatioCounter (+)")
	var lPercentage float64
	//  Constructing in Binary Tree Format
	if pcount%2 == 1 {
		lPercentage = pShare / float64(pcount%10+1)
	} else {
		lPercentage = pShare
	}

	log.Println("restartRatioCounter (-)")

	return lPercentage
}

func MaxofOrderinExchange(pBrokerId int) (string, error) {
	log.Println("MaxofOrderinExchange (+)")
	var lLastDoneExchange int
	var lMaxofExchange string

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MMOOIE01", lErr1)
		return lMaxofExchange, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select max(id),Exchange  
		from a_ipo_order_header 
		where brokerId  = ? and CreatedDate = now()`
		lRows1, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("MMOOIE02", lErr2)
			return lMaxofExchange, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lLastDoneExchange, &lMaxofExchange)
				if lErr3 != nil {
					log.Println("MMOOIE03", lErr3)
					return lMaxofExchange, lErr3
				}
			}
		}
	}
	log.Println("MaxofOrderinExchange (-)")
	return lMaxofExchange, nil

}

func CountPlacedInExchange(pBrokerId int, pExchange string) (int, error) {
	log.Println("CountPlacedInExchange (+)")
	var lCountonExchange int

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("MCPIE01", lErr1)
		return lCountonExchange, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select count(Exchange)
		from a_ipo_order_header 
		where brokerId  = ? and Exchange = ? and CreatedDate = now()`
		lRows1, lErr2 := lDb.Query(lCoreString, pBrokerId, pExchange)
		if lErr2 != nil {
			log.Println("MCPIE02", lErr2)
			return lCountonExchange, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lCountonExchange)
				if lErr3 != nil {
					log.Println("MCPIE03", lErr3)
					return lCountonExchange, lErr3
				}
			}
		}
	}
	log.Println("CountPlacedInExchange (-)")
	return lCountonExchange, nil

}

// func restartCounter(plastOrderTime string, pBrokerId int, pExchange string) (int, error) {
// 	log.Println("restartCounter (+)")
// 	var lCountonExchange int

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("MRC01", lErr1)
// 		return lCountonExchange, lErr1
// 	} else {
// 		lCoreString := `select count(oh.applicationNo)
// 		from a_ipo_order_header oh
// 		where  oh.CreatedDate > ? and oh.CreatedDate = now() and oh.brokerId = ? and oh.Exchange = ?`
// 		lRows1, lErr2 := lDb.Query(lCoreString, plastOrderTime, pBrokerId, pExchange)
// 		if lErr2 != nil {
// 			log.Println("MRC02", lErr2)
// 			return lCountonExchange, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows1.Next() {
// 				lErr3 := lRows1.Scan(&lCountonExchange)
// 				if lErr3 != nil {
// 					log.Println("MRC03", lErr3)
// 					return lCountonExchange, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("CountPlacedInExchange (-)")
// 	return lCountonExchange, nil

// }
