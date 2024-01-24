package sgbschedule

import (
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bsesgb"
	"log"
	"sync"
)

type SchRespStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}
type SgbBrokers struct {
	BrokerId int
	Exchange string
}

// func SgbDownStatusSch(w http.ResponseWriter, r *http.Request) {
// 	log.Println("SgbDownStatusSch (+)", r.Method)
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "GET" {
// 		//   var lSgbData bsesgb.SgbStatusReqStruct
// 		//   var lSgbDataArr [] bsesgb.SgbStatusReqStruct
// 		var lRespRec SchRespStruct
// 		var lSgbReqData bsesgb.SgbStatusReqStruct

// 		lRespRec.Status = common.SuccessCode
// 		// var lSgbBseRespDataArr []bsesgb.SgbDownDataRespStruct

// 		lUser := r.Header.Get("USER")
// 		// lUser := common.AUTOBOT
// 		lSgbDataArr, lErr1 := FetchOrderData()
// 		if lErr1 != nil {
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "SS01" + lErr1.Error()
// 			fmt.Fprintf(w, helpers.GetErrorString("SS01", "error in executing query"))
// 			return
// 		} else {
// 			lToken, lErr2 := exchangecall.BseGetToken(lUser)
// 			if lErr2 != nil {
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = "SS02" + lErr2.Error()
// 				fmt.Fprintf(w, helpers.GetErrorString("SS02", "error in getting token"))
// 				return
// 			} else {
// 				for lIdxOut := 0; lIdxOut < len(lSgbDataArr); lIdxOut++ {
// 					for lIdxIn := 0; lIdxIn < len(lSgbDataArr[lIdxOut].PanNo); lIdxIn++ {
// 						lSgbReqData.ScripId = lSgbDataArr[lIdxOut].ScripId
// 						lSgbReqData.PanNo = lSgbDataArr[lIdxOut].PanNo[lIdxIn]
// 						lSgbBseRespStruct, lErr3 := bsesgb.BseSgbDownOrder(lToken, lUser, lSgbReqData)
// 						for lIdx := 0; lIdx < len(lSgbBseRespStruct); lIdx++ {
// 							if "0" == lSgbBseRespStruct[lIdx].ErrorCode {
// 								lErr3 := updateSgbOrderHeader(lSgbBseRespStruct[lIdx], lUser)
// 								if lErr3 != nil {
// 									lRespRec.Status = common.ErrorCode
// 									lRespRec.ErrMsg = "SS03" + lErr2.Error()
// 									fmt.Fprintf(w, helpers.GetErrorString("SS03", "error in getting response from Bse"))
// 									return
// 								} else {
// 									lErr3 = updateSgbOrderDetails(lSgbBseRespStruct[lIdx], lUser)
// 								}
// 							}
// 							if lErr3 != nil {
// 								lRespRec.Status = common.ErrorCode
// 								lRespRec.ErrMsg = "SS03" + lErr2.Error()
// 								fmt.Fprintf(w, helpers.GetErrorString("SS03", "error in getting response from Bse"))
// 								return
// 							} else {
// 								// lSgbBseRespDataArr = append(lSgbBseRespDataArr, lSgbBseRespStruct)
// 								log.Println("response struct from bse", lSgbBseRespStruct)
// 							}
// 						}
// 					}

// 				}
// 			}
// 		}
// 		lData, lErr2 := json.Marshal(lRespRec)
// 		if lErr2 != nil {
// 			log.Println("SS04", lErr2)
// 			fmt.Fprintf(w, "SS04"+lErr2.Error())
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("SgbDownStatusSch (-)", r.Method)
// 	}
// }

func FetchOrderData(pBrokerId int) ([]bsesgb.SgbDataStruct, error) {
	log.Println("FetchOrderData (+)")
	var lSgbData bsesgb.SgbDataStruct
	var lSgbDataArr []bsesgb.SgbDataStruct

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)

	if lErr1 != nil {
		log.Println("SSFO01", lErr1)
		return lSgbDataArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select  aso.PanNo,aso.ScripId
		from a_sgb_orderheader aso,a_sgb_master asm,a_ipo_brokermaster BM 
		WHERE CreatedDate >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
		and cancelFlag!="Y" and aso.MasterId =asm.id and asm.Exchange ="BSE" and BM.Id = aso.brokerId  and aso.brokerId = ? `

		lRows, lErr1 := lDb.Query(lCoreString, pBrokerId)
		if lErr1 != nil {
			log.Println("SSFO01", lErr1)
			return lSgbDataArr, lErr1
		} else {
			for lRows.Next() {
				lPanNo := ""
				lErr1 = lRows.Scan(&lPanNo, &lSgbData.ScripId)

				if lErr1 != nil {
					log.Println("SSFO02", lErr1)
					return lSgbDataArr, lErr1
				} else {

					if len(lSgbDataArr) == 0 { //if arrary is empty push the panno and scripthid
						lSgbData.PanNo = append(lSgbData.PanNo, lPanNo)
						lSgbDataArr = append(lSgbDataArr, lSgbData)
					} else { //check the previous script id if scripotid match push panno in scriptid struct
						count := 0
						for i := 0; i < count; i++ {
							if lSgbData.ScripId == lSgbDataArr[i].ScripId {
								lSgbDataArr[i].PanNo = append(lSgbDataArr[i].PanNo, lPanNo)
								count++
							}
						}
						if count == 0 { //if count is therefore scriptid not already exist
							lSgbData.PanNo = append(lSgbData.PanNo, lPanNo)
							lSgbDataArr = append(lSgbDataArr, lSgbData)
						}
					}

				}

			}

		}

	}
	log.Println("FetchOrderData (-)")
	return lSgbDataArr, nil
}

func updateSgbOrderHeader(pSgbBseRespData bsesgb.SgbDownDataRespStruct, pUser string, pBrokerId int) error {
	log.Println("updateSgbOrderHeader (+)")

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSUSO1", lErr1)
		return lErr1
	} else {
		defer lDb.Close()
		for i := 0; i < len(pSgbBseRespData.Bids); i++ {
			lCoreString := `update a_sgb_orderheader  h
			set h.DpRemarks=?,h.DpStatus=?,h.RbiStatus=?,h.RbiRemarks=?,h.RbiInvstId=?,h.UpdatedDate=now(),h.UpdatedBy=?
			where h.Id in (select d.HeaderId
			from a_sgb_orderdetails d,a_sgb_orderheader h ,a_ipo_brokermaster BM 
			where d.HeaderId =h.Id and d.OrderNo =? and BM.id = h.brokerId and h.brokerId = ? ) and ScripId=? and PanNo=?`
			// log.Println(pSgbBseRespData.DpStatus, "dpstatus")
			_, lErr1 = lDb.Exec(lCoreString, pSgbBseRespData.DpRemarks, pSgbBseRespData.DpStatus, pSgbBseRespData.RbiStatus, pSgbBseRespData.RbiInvstId, pSgbBseRespData.RbiInvstId, pUser, pSgbBseRespData.Bids[i].OrderNo, pBrokerId, pSgbBseRespData.ScripId, pSgbBseRespData.PanNo)
			if lErr1 != nil {
				log.Println("SSUSO2", lErr1)
				return lErr1
			}
		}
	}
	log.Println("updateSgbOrderHeader (-)")
	return nil

}
func updateSgbOrderDetails(pSgbBseRespData bsesgb.SgbDownDataRespStruct, pUser string, pBrokerId int) error {
	log.Println("updateSgbOrderDetials (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSUSDO1", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		for i := 0; i < len(pSgbBseRespData.Bids); i++ {
			lCoreString := `update a_sgb_orderdetails d
			set d.AddedDate=?,d.ModifiedDate=?,UpdatedDate=now(),UpdatedBy=?
			where d.HeaderId in (select h.Id from a_sgb_orderdetails d,a_sgb_orderheader h,a_ipo_brokermaster BM
			where d.HeaderId =h.Id  and  BM.id = h.brokerId and h.brokerId = ?  ) and d.OrderNo =? `

			_, lErr1 = lDb.Exec(lCoreString, pSgbBseRespData.Bids[i].Addeddate, pSgbBseRespData.Bids[i].ModifiedDate, pUser, pBrokerId, pSgbBseRespData.Bids[i].OrderNo)
			if lErr1 != nil {
				log.Println("SSUSDO2", lErr1)
				return lErr1
			}
		}
	}
	log.Println("updateSgbOrderDetails (-)")
	return nil

}

func SgbDownStatusSch(lwg *sync.WaitGroup, pBrokerId int, pUser string) {
	log.Println("SgbDownStatusSch (+)")

	defer lwg.Done()
	var lSgbReqData bsesgb.SgbStatusReqStruct
	// dataChannel := make(chan error)

	lSgbDataArr, lErr1 := FetchOrderData(pBrokerId)
	if lErr1 != nil {
		log.Println("SSS01", lErr1)
		// return lErr1
	} else {
		lToken, lErr2 := exchangecall.BseGetToken(pUser, pBrokerId)
		if lErr2 != nil {
			log.Println("SSS02", lErr2)
			// return lErr2
		} else {
			for lIdxOut := 0; lIdxOut < len(lSgbDataArr); lIdxOut++ {
				for lIdxIn := 0; lIdxIn < len(lSgbDataArr[lIdxOut].PanNo); lIdxIn++ {
					lSgbReqData.ScripId = lSgbDataArr[lIdxOut].ScripId
					lSgbReqData.PanNo = lSgbDataArr[lIdxOut].PanNo[lIdxIn]
					lSgbBseRespStruct, lErr3 := bsesgb.BseSgbDownOrder(lToken, pUser, lSgbReqData, pBrokerId)
					for lIdx := 0; lIdx < len(lSgbBseRespStruct); lIdx++ {
						if "0" == lSgbBseRespStruct[lIdx].ErrorCode {
							lErr3 := updateSgbOrderHeader(lSgbBseRespStruct[lIdx], pUser, pBrokerId)
							if lErr3 != nil {
								log.Println("lErr3", lErr3)
							} else {
								lErr3 = updateSgbOrderDetails(lSgbBseRespStruct[lIdx], pUser, pBrokerId)
							}
						}
						if lErr3 != nil {
							log.Println("SSS03", lErr3)
							// return lErr3
						} else {
							// lSgbBseRespDataArr = append(lSgbBseRespDataArr, lSgbBseRespStruct)
							log.Println("response struct from bse", lSgbBseRespStruct)
						}
					}
				}

			}
		}
	}
	log.Println("SgbDownStatusSch (-)")
	// return nil
}

//  commented by pavithra
// func SgbStatusScheduler() {
// 	log.Println("SgbStatusScheduler (+)")
// 	var lRespRec SchRespStruct
// 	var lBrokerList []SgbBrokers
// 	var lWg sync.WaitGroup

// 	lSgbBrokers, lErr1 := SgbBrokerList()
// 	if lErr1 != nil {
// 		log.Println("SSSS01", lErr1)
// 		lRespRec.Status = common.ErrorCode
// 		lRespRec.ErrMsg = "SSSS01" + lErr1.Error()
// 	} else {
// 		lBrokerList = lSgbBrokers
// 		for _, BrokerStream := range lBrokerList {
// 			lWg.Add(1)
// if BrokerStream.Exchange == common.BSE {
// 	go SgbDownStatusSch(&lWg, BrokerStream.BrokerId, common.AUTOBOT)
// } else
// 			if BrokerStream.Exchange == common.NSE {
// 				go NseSgbFetchStatus(&lWg, BrokerStream.BrokerId, common.AUTOBOT)
// 			}
// 		}
// 	}
// 	log.Println("SgbStatusScheduler(-)")
// }

func SgbBrokerList() ([]SgbBrokers, error) {
	log.Println("SgbBrokerList (+)")
	var Brokers SgbBrokers
	var BrokersList []SgbBrokers
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SSUSO1", lErr1)
		return BrokersList, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select Bm.Id ,d.Stream 
						from a_ipo_brokermaster bm ,a_ipo_directory d,a_ipo_memberdetails m 
						where bm.Id = d.brokerMasterId 
						and bm.Id = m.BrokerId 
						and m.AllowedModules like '%Sgb%'
						and bm.Status = 'Y' 
						and d.Status ='Y'
						and m.Flag = 'Y'`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("SGBFSC02", lErr2)
			return BrokersList, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&Brokers.BrokerId, &Brokers.Exchange)

				if lErr3 != nil {
					log.Println("SGBFSC03", lErr3)
					return BrokersList, lErr3
				} else {
					BrokersList = append(BrokersList, Brokers)
				}
			}
		}
	}
	log.Println("SgbBrokerList (-)")
	return BrokersList, nil
}
