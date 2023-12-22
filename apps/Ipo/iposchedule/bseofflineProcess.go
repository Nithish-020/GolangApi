package iposchedule

import (
	"fcs23pkg/apps/Ipo/placeorder"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/integration/nse/nseipo"
	"fcs23pkg/util/emailUtil"
	"log"
)

// func BseOfflineSch(w http.ResponseWriter, r *http.Request) {
// 	log.Println("BseOfflineSch (+)", r.Method)
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-lHeaders", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "GET" {

// 		lUser := r.Header.Get("USER")
// log.Println("lUser", lUser)
// // This variable is used to store the response from PendingList method

// func BseOfflineSch(lUser string) OfflineAppStruct {
// 	log.Println("BseOfflineSch (+)")
// 	var lRespRec OfflineAppStruct

// 	lRespRec.Status = common.SuccessCode

// 	lExchange := "BSE"

// 	// Call PendingList method to fetch the PendingApplications from database.
// 	lPendingBidsArr, lIdValuesArr, lErr1 := PendingListSch(lExchange)
// 	if lErr1 != nil {
// 		log.Println("ISBOS01", lErr1)
// 		lRespRec.Status = common.ErrorCode
// 		lRespRec.ErrMsg = "ISBOS01" + lErr1.Error()
// 	} else {
// 		// log.Println("lIdValuesArr", lIdValuesArr)

// 		//call GetBidTrackIdSch method to get the bidtracking details from the database
// 		lBidTrackIdArr, lErr2 := GetBidTrackIdSch(lPendingBidsArr)
// 		if lErr2 != nil {
// 			log.Println("ISBOS02", lErr2)
// 			lRespRec.Status = common.ErrorCode
// 			lRespRec.ErrMsg = "ISBOS02" + lErr2.Error()
// 		} else {
// 			lFinalReqArr, lErr3 := ConstructPendingSch(lPendingBidsArr)
// 			if lErr3 != nil {
// 				log.Println("ISBOS03", lErr3)
// 				lRespRec.Status = common.ErrorCode
// 				lRespRec.ErrMsg = "ISBOS03" + lErr3.Error()
// 			} else {
// 				// ---------------->>>>> get the application no. ------------>>>>>>>>

// 				lAppArr, lErr4 := GetApiMaster(lFinalReqArr)
// 				if lErr4 != nil {
// 					log.Println("ISBOS04", lErr4)
// 					lRespRec.Status = common.ErrorCode
// 					lRespRec.ErrMsg = "ISBOS04" + lErr4.Error()
// 				} else {
// 					log.Println(lAppArr)
// 					log.Println(lFinalReqArr, "lFinalReqArr")
// 					if lFinalReqArr != nil {
// 						//-------------------------------------------------------------->>>>>>>>>>>

// 						for lReqIdx := 0; lReqIdx < len(lFinalReqArr); lReqIdx++ {
// 							//
// 							var lNseArr []nseipo.ExchangeReqStruct
// 							lNseArr = append(lNseArr, lFinalReqArr[lReqIdx])
// 							// Call ApplyIpo to pass the applications to Exchange.
// 							lRespArr, lErr5 := BseOrderPlace(lNseArr, lUser)
// 							if lErr5 != nil {
// 								log.Println("ISBOS05", lErr5)
// 								lRespRec.Status = common.ErrorCode
// 								lRespRec.ErrMsg = "ISBOS05" + lErr5.Error()
// 							} else {
// 								//  this for loop is used to avoid panic error added by prashanth
// 								for _, lResp := range lRespArr {
// 									lBseRespArr = append(lBseRespArr, lResp)

// 								}
// 								// this line creates the panic server error so it was commented by prashanth
// 								// lBseRespArr = append(lBseRespArr, lRespArr[lReqIdx])
// 							}
// 						}
// 							} else {
// 								log.Println(lRespArr)
// 								//call StatusValidationSch method to split the records by status and its give the finalized Array Value
// 								lFinalRespArr, lErr6 := StatusValidationSch(lRespArr)
// 								if lErr6 != nil {
// 									log.Println("ISBOS06", lErr6)
// 									lRespRec.Status = common.ErrorCode
// 									lRespRec.ErrMsg = "ISBOS06" + lErr6.Error()
// 								} else {
// 									log.Println("lFinalRespArr", lFinalRespArr)
// 									// this looping is used to request array and response array to update the records
// 									for lpReqArrIdx := 0; lpReqArrIdx < len(lPendingBidsArr); lpReqArrIdx++ {
// 										for lpRespArrIdx := 0; lpRespArrIdx < len(lFinalRespArr); lpRespArrIdx++ {
// 											// check whether the Req application no and response application no is equal or not
// 											if lPendingBidsArr[lpReqArrIdx].ApplicationNo == lFinalRespArr[lpRespArrIdx].ApplicationNo {
// 												//Call the UpdateRecordSch to Update the Response array
// 												lErr7 := UpdateRecordSch(lPendingBidsArr[lpReqArrIdx], lFinalRespArr[lpRespArrIdx], lIdValuesArr, lBidTrackIdArr, lUser)
// 												if lErr7 != nil {
// 													log.Println("ISBOS07", lErr7)
// 													lRespRec.Status = common.ErrorCode
// 													lRespRec.ErrMsg = "ISBOS07" + lErr7.Error()
// 												} else {
// 													lRespRec.Status = common.SuccessCode
// 													lRespRec.ErrMsg = "Processed Successfully...."
// 												}
// 											}
// 										}
// 									}
// 								}
// 								//-------------------------------------Mail sending process ----------------

// 								for i := 0; i < len(lFinalRespArr); i++ {
// 									for j := 0; j < len(lAppArr); j++ {
// 										if lFinalRespArr[i].ApplicationNo == lAppArr[j].ApplicationNumber {
// 											if lFinalRespArr[i].Status == "success" {
// 												lString := "IpoOrderMail"
// 												lStatus := "success"
// 												lIpoEmailContent, lErr8 := placeorder.ConstructMail(lAppArr[j], lStatus)
// 												if lErr8 != nil {
// 													log.Println("ISOAS08", lErr8)
// 													lRespRec.Status = common.ErrorCode
// 													lRespRec.ErrMsg = "ISOAS08" + lErr8.Error()

// 												} else {
// 													lErr9 := emailUtil.SendEmail(lIpoEmailContent, lString)
// 													if lErr9 != nil {
// 														log.Println("ISOAS09", lErr9)
// 														lRespRec.Status = common.ErrorCode
// 														lRespRec.ErrMsg = "ISOAS09" + lErr9.Error()
// 													}
// 												}
// 											} else {

// 												lString := "IpoOrderMail"
// 												lStatus := "failed"
// 												lIpoEmailContent, lErr10 := placeorder.ConstructMail(lAppArr[j], lStatus)
// 												if lErr10 != nil {
// 													log.Println("ISOAS10", lErr10)
// 													lRespRec.Status = common.ErrorCode
// 													lRespRec.ErrMsg = "ISOAS10" + lErr10.Error()
// 												} else {
// 													lErr11 := emailUtil.SendEmail(lIpoEmailContent, lString)
// 													// EmailInput.Subject = IpoOrderStruct.OrderDate
// 													if lErr11 != nil {
// 														log.Println("ISOAS11", lErr11)
// 														lRespRec.Status = common.ErrorCode
// 														lRespRec.ErrMsg = "ISOAS11" + lErr11.Error()
// 													}
// 												}
// 											}
// 										}
// 									}
// 								}
// 							}

// 					} else {
// 						lRespRec.ErrMsg = "No Records found for process application"
// 						lRespRec.Status = common.ErrorCode
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("BseOfflineSch (-)")
// 	return lRespRec
// }

// lData, lErr12 := json.Marshal(lRespRec)
// if lErr12 != nil {
// 	log.Println("ISOAS12", lErr12)
// 	fmt.Fprintf(w, "ISOAS12"+lErr12.Error())
// } else {
// 	fmt.Fprintf(w, string(lData))
// }
// log.Println("BseOfflineSch (-)", r.Method)

// }
// }

// func BseOrderPlace(pFinalReqArr []nseipo.ExchangeReqStruct, pUser string, pBrokerId int) ([]nseipo.ExchangeRespStruct, error) {
// 	log.Println("BseReqConstruct (+)")
// 	var lNseResArr []nseipo.ExchangeRespStruct

// 	lReqRec := placeorder.BseReqConstruct(pFinalReqArr)

// 	lRespRec, lErr1 := exchangecall.ApplyBseIpo(lReqRec, pUser, pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("ISBOP01", lErr1)
// 		return lNseResArr, lErr1
// 	} else {
// 		lNseResArr = placeorder.BseRespConstruct(lRespRec)
// 	}
// 	log.Println("BseReqConstruct (+)")
// 	return lNseResArr, nil
// }

func BseOfflineSch(lUser string, pBrokerId int) ScheduleStruct {
	log.Println("BseOfflineSch (+)")
	var lRespRec ScheduleStruct

	lRespRec.Status = common.SuccessCode

	lExchange := common.BSE
	// Call PendingList method to fetch the PendingApplications from database.
	lPendingBidsArr, lIdValuesArr, lErr1 := PendingListSch(pBrokerId, lExchange)
	if lErr1 != nil {
		log.Println("ISBOS01", lErr1)
		lRespRec.Status = common.ErrorCode
		lRespRec.ErrMsg = "ISBOS01" + lErr1.Error()
	} else {
		// log.Println("lIdValuesArr", lIdValuesArr)

		//call GetBidTrackIdSch method to get the bidtracking details from the database
		lBidTrackIdArr, lErr2 := GetBidTrackIdSch(pBrokerId, lPendingBidsArr)
		if lErr2 != nil {
			log.Println("ISBOS02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "ISBOS02" + lErr2.Error()
		} else {
			// log.Println("lBidTrackIdArr", lBidTrackIdArr)
			// log.Println("lPendingBidsArr", lPendingBidsArr)
			//todo:-------------------------------
			lFinalReqArr, lErr3 := ConstructPendingSch(lPendingBidsArr)
			if lErr3 != nil {
				log.Println("ISBOS03", lErr3)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "ISBOS03" + lErr3.Error()
			} else {
				// ---------------->>>>> get the application no. ------------>>>>>>>>

				lAppArr, lErr4 := GetApiMaster(lFinalReqArr)
				if lErr4 != nil {
					log.Println("ISBOS04", lErr4)
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "ISBOS04" + lErr4.Error()
				} else {
					// log.Println(lAppArr)
					lRespRec.TotalCount = len(lFinalReqArr)
					// log.Println(lFinalReqArr, "lFinalReqArr")
					if lFinalReqArr != nil {
						//-------------------------------------------------------------->>>>>>>>>>>
						var lBseRespArr []nseipo.ExchangeRespStruct
						for lReqIdx := 0; lReqIdx < len(lFinalReqArr); lReqIdx++ {
							//
							var lNseArr []nseipo.ExchangeReqStruct
							lNseArr = append(lNseArr, lFinalReqArr[lReqIdx])
							// Call ApplyIpo to pass the applications to Exchange.
							lRespArr, lErr5 := BseOrderPlace(lNseArr, lUser, pBrokerId)
							if lErr5 != nil {
								log.Println("ISBOS05", lErr5)
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "ISBOS05" + lErr5.Error()
							} else {
								//  this for loop is used to avoid panic error added by prashanth
								for _, lResp := range lRespArr {
									lBseRespArr = append(lBseRespArr, lResp)

								}
								// this line creates the panic server error so it was commented by prashanth
								// lBseRespArr = append(lBseRespArr, lRespArr[lReqIdx])
							}
						}
						//call StatusValidationSch method to split the records by status and its give the finalized Array Value
						lFinalRespArr, lSuccessCount, lErrCount, lErr6 := StatusValidationSch(lBseRespArr, common.BSE)
						if lErr6 != nil {
							log.Println("ISOAS06", lErr6)
							lRespRec.ErrCount = lErrCount
							lRespRec.SuccessCount = lSuccessCount
							lRespRec.Status = common.ErrorCode
							lRespRec.ErrMsg = "ISOAS06" + lErr6.Error()
						} else {
							lRespRec.ErrCount = lErrCount
							lRespRec.SuccessCount = lSuccessCount
							// log.Println("lFinalRespArr", lFinalRespArr)
							// this looping is used to request array and response array to update the records
							for lpReqArrIdx := 0; lpReqArrIdx < len(lPendingBidsArr); lpReqArrIdx++ {
								for lpRespArrIdx := 0; lpRespArrIdx < len(lFinalRespArr); lpRespArrIdx++ {
									// check whether the Req application no and response application no is equal or not
									if lPendingBidsArr[lpReqArrIdx].ApplicationNo == lFinalRespArr[lpRespArrIdx].ApplicationNo {
										//Call the UpdateRecordSch to Update the Response array
										lErr7 := UpdateRecordSch(lPendingBidsArr[lpReqArrIdx], lFinalRespArr[lpRespArrIdx], lIdValuesArr, lBidTrackIdArr, lUser)
										if lErr7 != nil {
											log.Println("ISBOS07", lErr7)
											lRespRec.Status = common.ErrorCode
											lRespRec.ErrMsg = "ISBOS07" + lErr7.Error()
										} else {
											lRespRec.Status = common.SuccessCode
											lRespRec.ErrMsg = "Processed Successfully...."
										}
									}
								}
							}
						}
						//-------------------------------------Mail sending process ----------------

						for i := 0; i < len(lFinalRespArr); i++ {
							for j := 0; j < len(lAppArr); j++ {
								if lFinalRespArr[i].ApplicationNo == lAppArr[j].ApplicationNumber {
									//added email validation
									if lAppArr[j].ClientEmail != "" {
										if lFinalRespArr[i].Status == "success" {
											lString := "IpoOrderMail"
											lStatus := "success"
											lIpoEmailContent, lErr8 := placeorder.ConstructMail(lAppArr[j], lStatus)
											if lErr8 != nil {
												log.Println("ISOAS08", lErr8)
												lRespRec.Status = common.ErrorCode
												lRespRec.ErrMsg = "ISOAS08" + lErr8.Error()

											} else {
												lErr9 := emailUtil.SendEmail(lIpoEmailContent, lString)
												if lErr9 != nil {
													log.Println("ISOAS09", lErr9)
													lRespRec.Status = common.ErrorCode
													lRespRec.ErrMsg = "ISOAS09" + lErr9.Error()
												}
											}
										} else {

											lString := "IpoOrderMail"
											lStatus := "failed"
											lIpoEmailContent, lErr10 := placeorder.ConstructMail(lAppArr[j], lStatus)
											if lErr10 != nil {
												log.Println("ISOAS10", lErr10)
												lRespRec.Status = common.ErrorCode
												lRespRec.ErrMsg = "ISOAS10" + lErr10.Error()
											} else {
												lErr11 := emailUtil.SendEmail(lIpoEmailContent, lString)
												// EmailInput.Subject = IpoOrderStruct.OrderDate
												if lErr11 != nil {
													log.Println("ISOAS11", lErr11)
													lRespRec.Status = common.ErrorCode
													lRespRec.ErrMsg = "ISOAS11" + lErr11.Error()
												}
											}
										}
									}
								}
							}
						}
					} else {
						lRespRec.ErrMsg = "No Records found for process application"
						lRespRec.Status = common.ErrorCode
					}
				}
			}
		}
	}
	log.Println("BseOfflineSch (-)")
	return lRespRec
}

// lData, lErr12 := json.Marshal(lRespRec)
// if lErr12 != nil {
// 	log.Println("ISOAS12", lErr12)
// 	fmt.Fprintf(w, "ISOAS12"+lErr12.Error())
// } else {
// 	fmt.Fprintf(w, string(lData))
// }
// log.Println("BseOfflineSch (-)", r.Method)

// }
// }

func BseOrderPlace(pFinalReqArr []nseipo.ExchangeReqStruct, pUser string, pBrokerId int) ([]nseipo.ExchangeRespStruct, error) {
	log.Println("BseReqConstruct (+)")
	var lNseResArr []nseipo.ExchangeRespStruct

	lReqRec := placeorder.BseReqConstruct(pFinalReqArr)

	lRespRec, lErr1 := exchangecall.ApplyBseIpo(lReqRec, pUser, pBrokerId)
	if lErr1 != nil {
		log.Println("ISBOP01", lErr1)
		return lNseResArr, lErr1
	} else {
		lNseResArr = placeorder.BseRespConstruct(lRespRec)
	}
	log.Println("BseReqConstruct (+)")
	return lNseResArr, nil
}
