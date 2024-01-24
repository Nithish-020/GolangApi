package ncbplaceorder

import (
	"database/sql"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
	"strconv"
	"time"
)

func NcbAcceptClientToOrder(pRequest NcbReqStruct, pClientId string, pBrokerId int) (bool, string, int, error) {

	log.Println("NcbAcceptClientToOrder (+)")

	lAllowToOrder := false
	lErrMsg := ""
	log.Println("pRequest: " + pRequest.SIText + "=")
	var lMasterId int

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AKEFO01", lErr1)
		return lAllowToOrder, lErrMsg, lMasterId, lErr1
	} else {
		defer lDb.Close()

		lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
		lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CloseTime"])
		lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_CancelAllowed"])
		lModifyAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_ModifyAllowed"])
		lMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_Immediate_Flag"])

		lInvalidInput := NcbVerifyValidInput(lDb, pRequest)
		log.Println("lInvalidInput", lInvalidInput)
		if lInvalidInput {
			lErrMsg = "AKEFO02 / Invalid Input, Unable to process your request."
			return lAllowToOrder, lErrMsg, lMasterId, nil
		} else {

			lTime, lErr3 := time.Parse("15:04:05", lCloseTime)
			if lErr3 != nil {
				log.Println("NNACTO03", lErr3)
				lErrMsg = "Unable to process your request, Please try after sometime."
				return lAllowToOrder, lErrMsg, lMasterId, lErr3
			} else {
				lSystemDate := time.Now().Local()
				lCurrentDate := time.Now().Format("2006-01-02")
				lConfigTime1 := time.Date(lSystemDate.Year(), lSystemDate.Month(), lSystemDate.Day(), lTime.Hour(), lTime.Minute(), lTime.Second(), 0, time.Local)
				lConfigUnixTime := lConfigTime1.Unix()

				lConfigTime2 := time.Date(lSystemDate.Year(), lSystemDate.Month(), lSystemDate.Day(), lSystemDate.Hour(), lSystemDate.Minute(), lSystemDate.Second(), 0, time.Local)
				lCurrentUnixTime := lConfigTime2.Unix()

				lIndicator, lEndDate, lErr4 := NcbcheckOrderApplied(lDb, pRequest, pClientId, pBrokerId)
				if lErr4 != nil {
					log.Println("NNACTO04", lErr4)
					lErrMsg = "Unable to process your request, Please try after sometime."
					return lAllowToOrder, lErrMsg, lMasterId, lErr4
				} else {
					log.Println("lIndicator,lEndDate,lOrderNo", lEndDate, lIndicator)

					if lIndicator == "Y" && pRequest.ActionType == "N" {
						lAllowToOrder = false
						lErrMsg = "NNACTO05/ Already you placed an order"
					} else {

						endDate, _ := time.Parse("2006-01-02", lEndDate)
						currentDate, _ := time.Parse("2006-01-02", lCurrentDate)

						if endDate.Equal(currentDate) && lCurrentUnixTime < lConfigUnixTime {

							lAllowToOrder = true
						} else if endDate.After(currentDate) || endDate.Equal(currentDate) {
							log.Println("elseif", endDate, currentDate)
							lAllowToOrder = true
						} else {
							lErrMsg = "NNACTO06 / Timing Closed for NCB"
							lAllowToOrder = false
						}

						if lMode == "I" {
							if pRequest.ActionType == "M" {

								if lModifyAllowed == "N" {
									lErrMsg = "NNACTO07/ Modification Not Allowed"
									lAllowToOrder = false
								} else {
									lAllowToOrder = true
								}
							} else if pRequest.ActionType == "D" {

								// If the cancel falg hold the value N means
								// cancellation is prohibited
								if lCancelAllowed == "N" {
									lErrMsg = "NNACTO08/ Cancellation not allowed"
									lAllowToOrder = false
								} else {
									lAllowToOrder = true
								}

							} else {
								lAllowToOrder = true
							}

						} else {

							if pRequest.ActionType == "M" {
								if lModifyAllowed == "Y" {
									lProcessFlag, lScheduleStatus, lErr9 := NcbCheckEligibleToModify(lDb, pRequest)
									if lErr9 != nil {
										log.Println("NNACTO09", lErr9)
										lErrMsg = "NNACTO09/ Unable to process your request, Please try after sometime."
										return lAllowToOrder, lErrMsg, lMasterId, lErr9
									} else {
										if lProcessFlag == "N" && lScheduleStatus == "N" {
											lBrokerMasterText, lErr10 := GetNcbSItext(lDb, pBrokerId)
											if lErr10 != nil {
												log.Println("NNACTO10", lErr10)
												lErrMsg = "NNACTO10/ Unable to process your request, Please try after sometime."
												return lAllowToOrder, lErrMsg, lMasterId, lErr10
											} else {
												if pRequest.SIText == lBrokerMasterText && pRequest.SIValue {
													lAllowToOrder = true
												} else {
													lErrMsg = "NNACTO11/ Policy Doesn't match"
													lAllowToOrder = false
												}
											}
										} else {
											lErrMsg = "NNACTO12/ Unable to Modify this order"
											lAllowToOrder = false
										}

									}
								} else {
									lErrMsg = "NNACTO13/ Modification not allowed"
									lAllowToOrder = false
								}
							} else if pRequest.ActionType == "D" {
								if lCancelAllowed == "Y" {
									lProcessFlag, lScheduleStatus, lErr14 := NcbCheckEligibleToModify(lDb, pRequest)
									if lErr14 != nil {
										log.Println("NNACTO14", lErr14)
										lErrMsg = "NNACTO14/ Unable to process your request, Please try after sometime."
										return lAllowToOrder, lErrMsg, lMasterId, lErr14
									} else {
										if lProcessFlag == "N" && lScheduleStatus == "N" {
											lAllowToOrder = true
										} else {
											lErrMsg = "NNACTO15/ Unable to Cancel this order"
											lAllowToOrder = false
										}
									}
								} else {
									lErrMsg = "SPAC16/ Cancellation not allowed"
									lAllowToOrder = false
								}
							} else {
								lBrokerMasterText, lErr17 := GetNcbSItext(lDb, pBrokerId)
								if lErr17 != nil {
									log.Println("NNACTO17", lErr17)
									lErrMsg = "NNACTO17/ Unable to process your request, Please try after sometime."
									return lAllowToOrder, lErrMsg, lMasterId, lErr17
								} else {
									if pRequest.SIText == lBrokerMasterText && pRequest.SIValue {
										lAllowToOrder = true
									} else {
										lErrMsg = "NNACTO18/ Policy Doesn't match"
										lAllowToOrder = false
									}
								}
							}

						}

					}
					lMasterId, lErr4 = GetNcbMasterId(lDb, pRequest)
					if lErr4 != nil {
						log.Println("NNACTO19", lErr4)
						lErrMsg = "NNACTO19/ Unable to process your request, Please try after sometime."
						return lAllowToOrder, lErrMsg, lMasterId, lErr4
					} else {
						if pRequest.ActionType == "D" {
							if pRequest.MasterId != lMasterId {
								lErrMsg = "NNACTO20/ Scrip not found"
								lAllowToOrder = false
							}
						} else if pRequest.ActionType == "M" {
							if pRequest.MasterId != 0 {
								lMasterId = pRequest.MasterId
							}
						} else {
							lMasterId = pRequest.MasterId
						}
					}
				}
			}
		}
	}
	log.Println("Eligible to Apply Order := ", lAllowToOrder)

	log.Println("NcBAcceptClientToOrder (-)")
	return lAllowToOrder, lErrMsg, lMasterId, nil
}

func NcbcheckOrderApplied(pDb *sql.DB, pRequest NcbReqStruct, pClientId string, pBrokerId int) (string, string, error) {
	log.Println("NcbcheckOrderApplied (+)")
	var lIndicator, lEndDate string

	lBrokerId := strconv.Itoa(pBrokerId)

	lCoreString1 := `select nvl(n.BiddingEndDate,'') EndDate
		                 from  a_ncb_master n
		                 where n.id = ?`

	lRows1, lErr2 := pDb.Query(lCoreString1, pRequest.MasterId)
	if lErr2 != nil {
		log.Println("NCOA02", lErr2)
		return lIndicator, lEndDate, lErr2
	} else {

		for lRows1.Next() {
			lErr3 := lRows1.Scan(&lEndDate)
			if lErr3 != nil {
				log.Println("NCOA03", lErr3)
				return lIndicator, lEndDate, lErr3
			}
		}

		lCoreString2 := `select (case when count(1) > 0 then 'Y' else 'N' end) OrderFound
						 from  a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
						 where n.id = h.MasterId
						 and d.headerId = h.Id
						 and h.cancelFlag = 'N'
						 and h.status  = 'success'
						 and n.id = ?
						 and h.clientId = ?
						 and h.brokerId  = ?`

		if pRequest.ActionType != "N" {

			lReqOrderNo := strconv.Itoa(pRequest.OrderNo)

			lCoreString2 = lCoreString2 + ` and d.ReqOrderNo = '` + lReqOrderNo + `'`
		}
		lRows2, lErr4 := pDb.Query(lCoreString2, pRequest.MasterId, pClientId, lBrokerId)
		if lErr4 != nil {
			log.Println("NCOA04", lErr4)
			return lIndicator, lEndDate, lErr4
		} else {
			for lRows2.Next() {
				lErr5 := lRows2.Scan(&lIndicator)
				if lErr5 != nil {
					log.Println("NCOA05", lErr5)
					return lIndicator, lEndDate, lErr5
				}
			}
		}
	}

	log.Println("NcbcheckOrderApplied (-)")
	return lIndicator, lEndDate, nil
}

func NcbVerifyValidInput(pDb *sql.DB, pRequest NcbReqStruct) bool {
	log.Println("NcbVerifyValidInput (+)")

	var lMasterId, lUnit int
	var lUnitPrice float64

	lVerifiedInput := true
	lConfigFile := common.ReadTomlConfig("toml/NcbConfig.toml")
	// lShowDiscount := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_ShowDiscount"])
	lMaxQty1 := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_GsecMAXQUANTITY"])
	lMaxQty2 := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_TbillMAXQUANTITY"])
	lMaxQty3 := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["NCB_SdlMAXQUANTITY"])

	lGsecMaxQty4, _ := strconv.Atoi(lMaxQty1)
	lTbillMaxQty5, _ := strconv.Atoi(lMaxQty2)
	lSdlMaxQty6, _ := strconv.Atoi(lMaxQty3)

	log.Println("ActionType", pRequest.ActionType)

	if pRequest.ActionType == "M" || pRequest.ActionType == "D" {

		if pRequest.OrderNo != 0 {

			lCoreString := `select nvl(n.id,0),n.CutoffPrice unitPrice 
	                        from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
	                        where n.id = h.MasterId 
	                        and h.Id = d.HeaderId 
	                         and d.ReqOrderNo = ?`

			lRows, lErr1 := pDb.Query(lCoreString, pRequest.OrderNo)
			if lErr1 != nil {
				log.Println("NPGNT01", lErr1)
				return lVerifiedInput
			} else {
				for lRows.Next() {
					lErr2 := lRows.Scan(&lMasterId, &lUnitPrice, &lUnit)
					if lErr2 != nil {
						log.Println("NPGNT02", lErr2)
						return lVerifiedInput
					} else {

						if pRequest.ActionType == "D" {

							lVerifiedInput = false
							return lVerifiedInput

						} else if pRequest.ActionType == "M" {
							if pRequest.OldUnit > 0 && pRequest.Unit > 0 {

								if (pRequest.Unit == pRequest.OldUnit && pRequest.Unit == lUnit) && (pRequest.Unit <= lTbillMaxQty5 || pRequest.Unit <= lGsecMaxQty4 || pRequest.Unit <= lSdlMaxQty6) {
									return lVerifiedInput
								} else {
									if pRequest.UnitPrice == lUnitPrice {
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

	} else if pRequest.ActionType == "N" {
		//verify the Unit and price greater than 0
		if (pRequest.Unit > 0 && pRequest.Price > 0 && pRequest.MasterId > 0) && (pRequest.Unit <= lTbillMaxQty5 || pRequest.Unit <= lGsecMaxQty4 || pRequest.Unit <= lSdlMaxQty6) {

			log.Println("Unit", pRequest.Unit, pRequest.Price, pRequest.MasterId, lTbillMaxQty5, lGsecMaxQty4, lSdlMaxQty6)

			lCoreString := `select nvl(n.id, 0),n.CutoffPrice unitPrice
			                 from a_ncb_master n
			                 where n.id = ?`

			lRows, lErr3 := pDb.Query(lCoreString, pRequest.MasterId)
			if lErr3 != nil {
				log.Println("NPGNT03", lErr3)
				return lVerifiedInput
			} else {
				for lRows.Next() {
					lErr4 := lRows.Scan(&lMasterId, &lUnitPrice)
					if lErr4 != nil {
						log.Println("NPGNT04", lErr4)
						return lVerifiedInput
					} else {
						log.Println("pRequest.Price", pRequest.Price, lUnitPrice)
						if pRequest.UnitPrice == lUnitPrice {
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

	log.Println("NcbVerifyValidInput (-)")
	return lVerifiedInput
}

func GetNcbSItext(pDb *sql.DB, pBrokerId int) (string, error) {
	log.Println("GetNcbSItext (-)")

	var lSItext string

	lCoreString := `SELECT nvl(bm.Ncb_SItext,'')SItext
	                FROM a_ipo_brokermaster bm 
	                WHERE bm.id = ?`

	lRows, lErr2 := pDb.Query(lCoreString, pBrokerId)
	if lErr2 != nil {
		log.Println("NPGNT02", lErr2)
		return lSItext, lErr2
	} else {
		for lRows.Next() {
			lErr3 := lRows.Scan(&lSItext)
			if lErr3 != nil {
				log.Println("NPGNT03", lErr3)
				return lSItext, lErr3
			}
		}
	}
	log.Println("GetNcbSItext (-)")
	return lSItext, nil
}

func NcbCheckEligibleToModify(pDb *sql.DB, pRequest NcbReqStruct) (string, string, error) {
	log.Println("NcbCheckEligibleToModify (+)")

	var lProcessFlag, lScheduleStatus string

	lCoreString := `select  nvl(h.ProcessFlag,'') processFlag ,nvl(h.ScheduleStatus,'') scheduleStatus
	                from a_ncb_orderheader h, a_ncb_orderdetails d	
                	where h.Id = d.HeaderId 
					and h.cancelFlag <> 'Y'
	                and d.ReqOrderNo  = ?`

	lRows, lErr2 := pDb.Query(lCoreString, pRequest.OrderNo)
	if lErr2 != nil {
		log.Println("NPNCETM01", lErr2)
		return lProcessFlag, lScheduleStatus, lErr2
	} else {
		for lRows.Next() {
			lErr3 := lRows.Scan(&lProcessFlag, &lScheduleStatus)
			if lErr3 != nil {
				log.Println("NPNCETM02", lErr3)
				return lProcessFlag, lScheduleStatus, lErr2
			}
		}
	}
	log.Println("NcbCheckEligibleToModify (-)")
	return lProcessFlag, lScheduleStatus, nil
}

func GetNcbMasterId(pDb *sql.DB, pRequest NcbReqStruct) (int, error) {
	log.Println("GetNcbMasterId (-)")

	var lMasterId int

	lCoreString := `select (case when count(1) > 0 then n.id else 0 end) masterId
	                from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
	                where n.id = h.MasterId 
	                and h.Id = d.HeaderId 
	                and d.ReqOrderNo = ? `

	lRows, lErr1 := pDb.Query(lCoreString, pRequest.OrderNo)
	if lErr1 != nil {
		log.Println("NPGMI01", lErr1)
		return lMasterId, lErr1
	} else {
		for lRows.Next() {
			lErr2 := lRows.Scan(&lMasterId)
			if lErr2 != nil {
				log.Println("NPGMI02", lErr2)
				return lMasterId, lErr2
			}
		}
	}
	log.Println("GetNcbMasterId (-)")
	return lMasterId, nil
}
