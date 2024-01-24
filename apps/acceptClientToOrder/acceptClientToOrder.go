package acceptclienttoorder

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
	"strconv"
	"time"
)

func AcceptClientToOrder(pRequest map[string]interface{}, pMethod string, pClientId string, pBrokerId int) (bool, string, error) {
	log.Println("AcceptClientToOrder (+)")
	lAllowToOrder := false
	lErrMsg := ""
	log.Println("pRequest: ", pRequest)

	lConfigFile := common.ReadTomlConfig("toml/SgbConfig.toml")
	lCloseTime := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CloseTime"])
	lCancelAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_CancelAllowed"])
	// lModifyAllowed := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_ModifyAllowed"])
	lMode := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["SGB_Immediate_Flag"])

	lInvalidInput := VerifyValidInput(pRequest, lMode)
	if lInvalidInput {
		lErrMsg = "Unable to process your request, Please try after sometime."
		return lAllowToOrder, lErrMsg, nil
	} else {

		lTime, lErr1 := time.Parse("15:04:05", lCloseTime)
		if lErr1 != nil {
			log.Println("SSFP01", lErr1)
			lErrMsg = "Unable to process your request, Please try after sometime."
			return lAllowToOrder, lErrMsg, lErr1
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

			lIndicator, lEndDate, lOrderNo, lErr2 := checkOrderApplied(pMethod, pRequest, pClientId, pBrokerId)
			if lErr2 != nil {
				log.Println("SSFP02", lErr2)
				lErrMsg = "Unable to process your request, Please try after sometime."
				return lAllowToOrder, lErrMsg, lErr2
			} else {
				log.Println("lIndicator,lEndDate,lOrderNo", lEndDate, lIndicator, lOrderNo)
				if lIndicator == "Y" && fmt.Sprintf("%v", pRequest["actionCode"]) == "N" {
					lAllowToOrder = false
					lErrMsg = "Already you placed an order"
				} else {

					// Parse the string dates into time.Time objects
					endDate, _ := time.Parse("2006-01-02", lEndDate)
					currentDate, _ := time.Parse("2006-01-02", lCurrentDate)

					// if today is the last day of the bid and current time is less than close time
					// then allow to order,Modify and cancel
					if endDate.Equal(currentDate) && lCurrentUnixTime < lConfigUnixTime {
						log.Println("if", endDate, currentDate, lCurrentUnixTime, lConfigUnixTime)

						lAllowToOrder = true

						//if today is not the last day of the bid
						// then allow to order,Modify and cancel
					} else if endDate.After(currentDate) || endDate.Equal(currentDate) {
						log.Println("elseif", endDate, currentDate)
						lAllowToOrder = true
					} else {
						log.Println("else")
						lAllowToOrder = false
					}

					// if fmt.Sprintf("%v", pRequest["actionCode"]) == "M" || fmt.Sprintf("%v", pRequest["actionCode"]) == "C" {
					// 	if fmt.Sprintf("%v", pRequest["orderNo"]) == lOrderNo {
					// 		lAllowToOrder = true
					// 	} else {
					// 		lAllowToOrder = false
					// 	}
					// } else {
					// 	lAllowToOrder = false
					// }

					if lMode == "I" {
						// If the user willing to Modify the order in immediate mode
						// then check the Modification is allowed or not
						if fmt.Sprintf("%v", pRequest["actionCode"]) == "M" {

							// // If the cancel falg hold the value N means
							// modification is prohibited
							// if lModifyAllowed == "N" {
							// lAllowToOrder = false
							// } else {
							// 	lAllowToOrder = true
							// }
							lAllowToOrder = false

							// If the user willing to cancel the order in immediate mode
							// then check the cancellation is allowed or not
						} else if fmt.Sprintf("%v", pRequest["actionCode"]) == "C" {

							// If the cancel falg hold the value N means
							// cancellation is prohibited
							if lCancelAllowed == "N" {
								lAllowToOrder = false
							} else {
								lAllowToOrder = true
							}
						} else {
							lAllowToOrder = true
						}
					}
				}
			}
			log.Println("Eligible to Apply Order := ", lAllowToOrder)
		}
	}

	log.Println("AcceptClientToOrder (-)")
	return lAllowToOrder, lErrMsg, nil
}

func checkOrderApplied(pMethod string, pRequest map[string]interface{}, pClientId string, pBrokerId int) (string, string, string, error) {
	log.Println("checkOrderApplied (+)")
	var lIndicator, lTabel1, lTabel2, lEndDate, lOrderNo string

	// To Establish A database connection,call LocalDbConnect Method
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("AKEFO01", lErr1)
		return lIndicator, lEndDate, lOrderNo, lErr1
	} else {
		defer lDb.Close()

		lBrokerId := strconv.Itoa(pBrokerId)

		if pMethod == "sgb" {
			lTabel1 = ` a_sgb_master m `
			lTabel2 = ` a_sgb_orderheader oh ,a_sgb_orderdetails od ,a_sgb_master m `
		} else if pMethod == "ncb" {
			lTabel1 = ` a_ncb_master m `
			lTabel2 = ` a_ncb_orderheader oh ,a_ncb_orderdetails od ,a_ncb_master m `
		}

		lCoreString1 := `select nvl(m.BiddingEndDate,'') EndDate
						from ` + lTabel1 + `
						where m.id = ?`

		lRows1, lErr2 := lDb.Query(lCoreString1, pRequest["masterId"])
		if lErr2 != nil {
			log.Println("AKEFO02", lErr2)
			return lIndicator, lEndDate, lOrderNo, lErr2
		} else {

			for lRows1.Next() {
				lErr3 := lRows1.Scan(&lEndDate)
				if lErr3 != nil {
					log.Println("AKEFO03", lErr3)
					return lIndicator, lEndDate, lOrderNo, lErr3
				}
			}

			lCoreString2 := `select (case when count(1) > 0 then 'Y' else 'N' end)OrderFound,nvl(od.ReqOrderNo,'') OrderNo 
						from ` + lTabel2 + `
						where m.id = oh.MasterId and od.HeaderId = oh.Id 
						and oh.cancelFlag = 'N' and oh.Status = 'success'
						and m.id = ` + fmt.Sprintf("%v", pRequest["masterId"]) + ` and oh.ClientId = ?
						and oh.brokerId = ? `
			if fmt.Sprintf("%v", pRequest["actionCode"]) != "N" {
				lCoreString2 = lCoreString2 + ` and od.ReqOrderNo = '` + fmt.Sprintf("%v", pRequest["orderNo"]) + `'`
			}

			lRows2, lErr2 := lDb.Query(lCoreString2, pClientId, lBrokerId)
			if lErr2 != nil {
				log.Println("AKEFO02", lErr2)
				return lIndicator, lEndDate, lOrderNo, lErr2
			} else {
				for lRows2.Next() {
					lErr3 := lRows2.Scan(&lIndicator, &lOrderNo)
					if lErr3 != nil {
						log.Println("AKEFO03", lErr3)
						return lIndicator, lEndDate, lOrderNo, lErr3
					}
				}
			}
		}
	}

	log.Println("checkOrderApplied (-)")
	return lIndicator, lEndDate, lOrderNo, nil
}

// MasterId   int    `json:"masterId"`
// BidId      string `json:"bidId"`
// Unit       int    `json:"unit"`
// OldUnit    int    `json:"oldUnit"`
// Price      int    `json:"price"`
// ActionCode string `json:"actionCode"`
// OrderNo    string `json:"orderNo"`
// Amount     int    `json:"amount"`
// PreApply   string `json:"preApply"`
// SIText     string `json:"SItext"`
// SIValue    bool   `json:"SIvalue"`
func VerifyValidInput(pRequest map[string]interface{}, pMode string) bool {
	log.Println("VerifyValidInput (+)")

	lVerifiedInput := true

	// if fmt.Sprintf("%v", pRequest["unit"]) != 0 {

	// }

	if fmt.Sprintf("%v", pRequest["actionCode"]) == "N" || fmt.Sprintf("%v", pRequest["actionCode"]) == "M" || fmt.Sprintf("%v", pRequest["actionCode"]) == "D" {
		if (fmt.Sprintf("%v", pRequest["actionCode"]) == "M" || fmt.Sprintf("%v", pRequest["actionCode"]) == "D") && fmt.Sprintf("%v", pRequest["orderNo"]) != "" {

			lVerifiedInput = false
			return lVerifiedInput
		} else if fmt.Sprintf("%v", pRequest["actionCode"]) == "N" {
			lVerifiedInput = false
			return lVerifiedInput
		} else {
			return lVerifiedInput
		}
	} else {
		return lVerifiedInput
	}
	log.Println("VerifyValidInput (-)")
	return lVerifiedInput
}
