package sgbplaceorder

import (
	"database/sql"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
	"strconv"
	"time"
)

func AcceptClientToOrder(pRequest SgbReqStruct, pClientId string, pBrokerId int) (bool, string, int, error) {
	log.Println("AcceptClientToOrder (+)")
	lAllowToOrder := false
	lErrMsg := ""
	log.Println("pRequest:=" + pRequest.SIText + "=")
	var lMasterId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("SPAC01", lErr1)
		return lAllowToOrder, lErrMsg, lMasterId, lErr1
	} else {
		defer lDb.Close()

		lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
		lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])
		lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CancelAllowed"])
		lModifyAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_ModifyAllowed"])
		lMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_Immediate_Flag"])

		lInvalidInput := VerifyValidInput(lDb, pRequest)
		if lInvalidInput {
			lErrMsg = "SPAC02/ Invalid Input, Unable to process your request."
			return lAllowToOrder, lErrMsg, lMasterId, nil
		} else {

			lTime, lErr1 := time.Parse("15:04:05", lCloseTime)
			if lErr1 != nil {
				log.Println("SPAC03", lErr1)
				lErrMsg = "SPAC03/ Unable to process your request, Please try after sometime."
				return lAllowToOrder, lErrMsg, lMasterId, lErr1
			} else {
				// Get current date
				lSystemDate := time.Now().Local()

				// Get the current date for comparing with the EndDate
				lCurrentDate := time.Now().Format("2006-01-02") // YYYY-MM-DD

				// Set the date component of lTime to today's date
				lConfigTime1 := time.Date(lSystemDate.Year(), lSystemDate.Month(), lSystemDate.Day(), lTime.Hour(), lTime.Minute(), lTime.Second(), 0, time.Local)
				lConfigUnixTime := lConfigTime1.Unix()

				// Set the date component of lTime to today's date
				lConfigTime2 := time.Date(lSystemDate.Year(), lSystemDate.Month(), lSystemDate.Day(), lSystemDate.Hour(), lSystemDate.Minute(), lSystemDate.Second(), 0, time.Local)
				lCurrentUnixTime := lConfigTime2.Unix()

				lIndicator, lEndDate, lErr2 := checkOrderApplied(lDb, pRequest, pClientId, pBrokerId)
				if lErr2 != nil {
					log.Println("SPAC04", lErr2)
					lErrMsg = "SPAC04/ Unable to process your request, Please try after sometime."
					return lAllowToOrder, lErrMsg, lMasterId, lErr2
				} else {

					if lIndicator == "Y" && pRequest.ActionCode == "N" {
						lAllowToOrder = false
						lErrMsg = "SPAC05/ Already you placed an order"
					} else {

						// Parse the string dates into time.Time objects
						endDate, _ := time.Parse("2006-01-02", lEndDate)
						currentDate, _ := time.Parse("2006-01-02", lCurrentDate)

						// if today is the last day of the bid and current time is less than close time
						// then allow to order,Modify and cancel
						if endDate.Equal(currentDate) && lCurrentUnixTime < lConfigUnixTime {
							lAllowToOrder = true
							//if today is not the last day of the bid
							// then allow to order,Modify and cancel
						} else if endDate.After(currentDate) || endDate.Equal(currentDate) {
							lAllowToOrder = true
						} else {
							lErrMsg = "SPAC06/ Timing Closed for SGB"
							lAllowToOrder = false
						}

						// if pRequest.ActionCode == "M" || pRequest.ActionCode == "D" {
						// 	if pRequest.OrderNo == lOrderNo {
						// 		lAllowToOrder = true
						// 	} else {
						// 		lErrMsg = "Order Doesn't exist"
						// 		lAllowToOrder = false
						// 	}
						// }

						if lMode == "I" {
							// If the user willing to Modify the order in immediate mode
							// then check the Modification is allowed or not
							if pRequest.ActionCode == "M" {

								// // If the cancel falg hold the value N means
								// modification is prohibited
								if lModifyAllowed == "N" {
									lErrMsg = "SPAC07/ Modification Not Allowed"
									lAllowToOrder = false
								} else {
									lAllowToOrder = true
								}
								// If the user willing to cancel the order in immediate mode
								// then check the cancellation is allowed or not
							} else if pRequest.ActionCode == "D" {

								// If the cancel falg hold the value N means
								// cancellation is prohibited
								if lCancelAllowed == "N" {
									lErrMsg = "SPAC08/ Cancellation not allowed"
									lAllowToOrder = false
								} else {
									lAllowToOrder = true
								}
							} else {
								lAllowToOrder = true
							}
						} else {

							if pRequest.ActionCode == "M" {
								if lModifyAllowed == "Y" {
									lProcessFlag, lScheduleStatus, lErr := CheckEligibleToModify(lDb, pRequest)
									if lErr != nil {
										log.Println("SPAC09", lErr)
										lErrMsg = "SPAC09/ Unable to process your request, Please try after sometime."
										return lAllowToOrder, lErrMsg, lMasterId, lErr
									} else {
										if lProcessFlag == "N" && lScheduleStatus == "N" {
											// lAllowToOrder = true
											lBrokerMasterText, lErr := GetSgbSItext(lDb, pBrokerId)
											if lErr != nil {
												log.Println("SPAC10", lErr)
												lErrMsg = "SPAC10/ Unable to process your request, Please try after sometime."
												return lAllowToOrder, lErrMsg, lMasterId, lErr
											} else {
												if pRequest.SIText == lBrokerMasterText && pRequest.SIValue {
													lAllowToOrder = true
												} else {
													lErrMsg = "SPAC11/ Policy Doesn't match"
													lAllowToOrder = false
												}
											}
										} else {
											lErrMsg = "SPAC12/ Unable to Modify this order"
											lAllowToOrder = false
										}

									}
								} else {
									lErrMsg = "SPAC13/ Modification not allowed"
									lAllowToOrder = false
								}
							} else if pRequest.ActionCode == "D" {
								if lCancelAllowed == "Y" {
									lProcessFlag, lScheduleStatus, lErr := CheckEligibleToModify(lDb, pRequest)
									if lErr != nil {
										log.Println("SPAC14", lErr)
										lErrMsg = "SPAC14/ Unable to process your request, Please try after sometime."
										return lAllowToOrder, lErrMsg, lMasterId, lErr
									} else {
										if lProcessFlag == "N" && lScheduleStatus == "N" {
											lAllowToOrder = true
										} else {
											lErrMsg = "SPAC15/ Unable to Cancel this order"
											lAllowToOrder = false
										}
									}
								} else {
									lErrMsg = "SPAC16/ Cancellation not allowed"
									lAllowToOrder = false
								}
							} else {
								lBrokerMasterText, lErr := GetSgbSItext(lDb, pBrokerId)
								if lErr != nil {
									log.Println("SPAC17", lErr)
									lErrMsg = "SPAC17/ Unable to process your request, Please try after sometime."
									return lAllowToOrder, lErrMsg, lMasterId, lErr
								} else {
									if pRequest.SIText == lBrokerMasterText && pRequest.SIValue {
										lAllowToOrder = true
									} else {
										lErrMsg = "SPAC18/ Policy Doesn't match"
										lAllowToOrder = false
									}
								}
							}
						}
					}
					lMasterId, lErr2 = GetMasterId(lDb, pRequest)
					if lErr2 != nil {
						log.Println("SPAC19", lErr2)
						lErrMsg = "SPAC19/ Unable to process your request, Please try after sometime."
						return lAllowToOrder, lErrMsg, lMasterId, lErr2
					} else {
						if pRequest.ActionCode == "D" {
							if pRequest.MasterId != lMasterId {
								lErrMsg = "SPAC20/ Scrip not found"
								lAllowToOrder = false
							}
						} else if pRequest.ActionCode == "M" {
							if pRequest.MasterId != 0 {
								lMasterId = pRequest.MasterId
							}
							// lAllowToOrder = true
						} else {
							lMasterId = pRequest.MasterId
							// lAllowToOrder = true
						}
					}
				}
			}
		}
	}
	log.Println("Eligible to Apply Order := ", lAllowToOrder, lErrMsg, lMasterId)
	log.Println("AcceptClientToOrder (-)")
	return lAllowToOrder, lErrMsg, lMasterId, nil
}

func checkOrderApplied(pDb *sql.DB, pRequest SgbReqStruct, pClientId string, pBrokerId int) (string, string, error) {
	log.Println("checkOrderApplied (+)")
	var lIndicator, lEndDate string

	lBrokerId := strconv.Itoa(pBrokerId)

	lCoreString1 := `select nvl(m.BiddingEndDate,'') EndDate
						from  a_sgb_master m 
						where m.id = ?`

	lRows1, lErr2 := pDb.Query(lCoreString1, pRequest.MasterId)
	if lErr2 != nil {
		log.Println("AKEFO02", lErr2)
		return lIndicator, lEndDate, lErr2
	} else {

		for lRows1.Next() {
			lErr3 := lRows1.Scan(&lEndDate)
			if lErr3 != nil {
				log.Println("AKEFO03", lErr3)
				return lIndicator, lEndDate, lErr3
			}
		}

		lCoreString2 := `select (case when count(1) > 0 then 'Y' else 'N' end)OrderFound
						from a_sgb_orderheader oh ,a_sgb_orderdetails od ,a_sgb_master m 
						where m.id = oh.MasterId and od.HeaderId = oh.Id 
						and oh.cancelFlag = 'N' and oh.Status = 'success'
						and m.id = ? and oh.ClientId = ?
						and oh.brokerId = ? `
		if pRequest.ActionCode != "N" {
			lCoreString2 = lCoreString2 + ` and od.ReqOrderNo = '` + pRequest.OrderNo + `'`
		}

		lRows2, lErr2 := pDb.Query(lCoreString2, pRequest.MasterId, pClientId, lBrokerId)
		if lErr2 != nil {
			log.Println("AKEFO02", lErr2)
			return lIndicator, lEndDate, lErr2
		} else {
			for lRows2.Next() {
				lErr3 := lRows2.Scan(&lIndicator)
				if lErr3 != nil {
					log.Println("AKEFO03", lErr3)
					return lIndicator, lEndDate, lErr3
				}
			}
		}
	}
	log.Println("checkOrderApplied (-)")
	return lIndicator, lEndDate, nil
}

func VerifyValidInput(pDb *sql.DB, pRequest SgbReqStruct) bool {
	log.Println("VerifyValidInput (+)")
	var lMasterId, lUnitPrice, lOrderedUnit int
	var lUnit string
	lVerifiedInput := true
	lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
	lMaxQty1 := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_MAXQUANTITY"])
	lShowDiscount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_ShowDiscount"])
	lMaxQty2, _ := strconv.Atoi(lMaxQty1)

	if pRequest.ActionCode == "M" || pRequest.ActionCode == "D" {
		// check the given OrderNo is Valid or Not
		if pRequest.OrderNo != "" {

			lCoreString := `select nvl(m.id,0),
							(case when count(1) > 0 then
							(case when '` + lShowDiscount + `' = 'Y' then m.MaxPrice - (m.MaxPrice-m.MinPrice)  else m.MinPrice  end ) 
							else 0 end ) unitPrice ,
							nvl(d.ReqSubscriptionUnit,'') unit
							from a_sgb_orderdetails d,a_sgb_orderheader h,a_sgb_master m
							where m.id = h.MasterId 
							and h.Id = d.HeaderId 
							and d.ReqOrderNo = ?`

			lRows, lErr2 := pDb.Query(lCoreString, pRequest.OrderNo)
			if lErr2 != nil {
				log.Println("SPGST02", lErr2)
				return lVerifiedInput
			} else {
				for lRows.Next() {
					lErr3 := lRows.Scan(&lMasterId, &lUnitPrice, &lUnit)
					if lErr3 != nil {
						log.Println("SPGST03", lErr3)
						return lVerifiedInput
					} else {
						lOrderedUnit, _ = strconv.Atoi(lUnit)
						// log.Println("lMasterId,lUnitPrice,lOrderdUnit", lMasterId, lUnitPrice)
						if pRequest.ActionCode == "D" {

							lVerifiedInput = false
							return lVerifiedInput
						} else if pRequest.ActionCode == "M" {
							if pRequest.OldUnit > 0 && pRequest.Unit > 0 {

								if pRequest.Unit == pRequest.OldUnit && pRequest.Unit <= lMaxQty2 && pRequest.Unit == lOrderedUnit {
									return lVerifiedInput
								} else {
									if pRequest.Price == lUnitPrice {
										lVerifiedInput = false
										return lVerifiedInput
									} else {
										return lVerifiedInput
									}
								}
							} else {
								return lVerifiedInput
							}
						} else {
							return lVerifiedInput
						}
					}
				}
			}
		} else {
			return lVerifiedInput
		}

	} else if pRequest.ActionCode == "N" {
		//verify the Unit and price greater than 0
		if pRequest.Unit > 0 && pRequest.Price > 0 && pRequest.Unit <= lMaxQty2 && pRequest.MasterId > 0 {
			log.Println("Inside if ", pRequest.Unit, pRequest.Price, pRequest.MasterId)
			lCoreString := `select nvl(m.id,0),
							(case when count(1) > 0 then
							(case when '` + lShowDiscount + `' = 'Y' then m.MaxPrice - (m.MaxPrice-m.MinPrice)  else m.MinPrice  end ) 
							else 0 end ) unitPrice 
							from a_sgb_master m
							where  m.id = ?`

			lRows, lErr2 := pDb.Query(lCoreString, pRequest.MasterId)
			if lErr2 != nil {
				log.Println("SPGST02", lErr2)
				return lVerifiedInput
			} else {
				for lRows.Next() {
					lErr3 := lRows.Scan(&lMasterId, &lUnitPrice)
					if lErr3 != nil {
						log.Println("SPGST03", lErr3)
						return lVerifiedInput
					} else {
						log.Println("pRequest.Price", pRequest.Price, lUnitPrice)
						if pRequest.Price == lUnitPrice {
							lVerifiedInput = false
							return lVerifiedInput
						} else {
							return lVerifiedInput
						}
					}
				}
			}
		} else {
			return lVerifiedInput
		}
	} else {
		return lVerifiedInput
	}
	log.Println("VerifyValidInput (-)")
	return lVerifiedInput
}

func GetSgbSItext(pDb *sql.DB, pBrokerId int) (string, error) {
	log.Println("GetSgbSItext (-)")

	var lSItext string

	lCoreString := `SELECT nvl(trim(bm.Sgb_SItext),'')SItext
					FROM a_ipo_brokermaster bm 
					WHERE bm.id = ?`

	lRows, lErr2 := pDb.Query(lCoreString, pBrokerId)
	if lErr2 != nil {
		log.Println("SPGST02", lErr2)
		return lSItext, lErr2
	} else {
		for lRows.Next() {
			lErr3 := lRows.Scan(&lSItext)
			if lErr3 != nil {
				log.Println("SPGST03", lErr3)
				return lSItext, lErr3
			}
		}
	}
	log.Println("GetSgbSItext (-)")
	return lSItext, nil
}

func CheckEligibleToModify(pDb *sql.DB, pRequest SgbReqStruct) (string, string, error) {
	log.Println("CheckEligibleToModify (+)")

	var lProcessFlag, lScheduleStatus string

	lCoreString := `select nvl(h.ProcessFlag,'') processFlag ,nvl(h.ScheduleStatus,'') scheduleStatus
					from a_sgb_orderheader h,a_sgb_orderdetails d
					where h.Id = d.HeaderId 
					and h.cancelFlag <> 'Y'
					and d.ReqOrderNo = ?`

	lRows, lErr2 := pDb.Query(lCoreString, pRequest.OrderNo)
	if lErr2 != nil {
		log.Println("SPCETM01", lErr2)
		return lProcessFlag, lScheduleStatus, lErr2
	} else {
		for lRows.Next() {
			lErr3 := lRows.Scan(&lProcessFlag, &lScheduleStatus)
			if lErr3 != nil {
				log.Println("SPCETM02", lErr3)
				return lProcessFlag, lScheduleStatus, lErr2
			}
		}
	}
	log.Println("CheckEligibleToModify (-)")
	return lProcessFlag, lScheduleStatus, nil
}

func GetMasterId(pDb *sql.DB, pRequest SgbReqStruct) (int, error) {
	log.Println("GetMasterId (-)")

	var lMasterId int

	lCoreString := `select (case when count(1) > 0 then m.id else 0 end) masterId
					from a_sgb_master m,a_sgb_orderheader h,a_sgb_orderdetails d
					where m.id = h.MasterId 
					and h.Id = d.HeaderId 
					and d.ReqOrderNo = ?`

	lRows, lErr2 := pDb.Query(lCoreString, pRequest.OrderNo)
	if lErr2 != nil {
		log.Println("SPGST02", lErr2)
		return lMasterId, lErr2
	} else {
		for lRows.Next() {
			lErr3 := lRows.Scan(&lMasterId)
			if lErr3 != nil {
				log.Println("SPGST03", lErr3)
				return lMasterId, lErr3
			}
		}
	}
	log.Println("GetMasterId (-)")
	return lMasterId, nil
}
