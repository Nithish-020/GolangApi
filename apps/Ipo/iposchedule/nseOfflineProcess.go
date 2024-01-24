package iposchedule

import (
	"fcs23pkg/apps/Ipo/placeorder"
	"fcs23pkg/apps/exchangecall"
	"fcs23pkg/common"
	"fcs23pkg/util/emailUtil"
	"log"
)

// this struct is response struct for OfflineSchedulePrgm
type OfflineAppStruct struct {
	TotalCount   int    `json:"totalCount"`
	SuccessCount int    `json:"successCount"`
	ErrCount     int    `json:"errCount"`
	Status       string `json:"status"`
	ErrMsg       string `json:"errMsg"`
}

/*
Purpose:This Method is used to process the pending application to the exchange
Parameters:

        User

Response:

	==========
	*On Sucess
	==========
	{
		"status": "success"
	}
	=========
	!On Error
	=========
	{
		"status": "E",
		"errMsg": "Error : Invalid Pending Bids"
	}

Author: KavyaDharshani
Date: 12SEP2023
*/
func NseOfflineSch(lUser string, pBrokerId int) OfflineAppStruct {
	log.Println("OfflineAppSch (+)")

	//This variable is used to store the response from PendingList method
	var lRespRec OfflineAppStruct
	lRespRec.Status = common.SuccessCode
	lExchange := common.NSE
	// Call PendingList method to fetch the PendingApplications from database.
	lPendingBidsArr, lIdValuesArr, lErr1 := PendingListSch(pBrokerId, lExchange)
	if lErr1 != nil {
		log.Println("ISOAS01", lErr1)
		lRespRec.Status = common.ErrorCode
		lRespRec.ErrMsg = "ISOAS01" + lErr1.Error()
	} else {

		//call GetBidTrackIdSch method to get the bidtracking details from the database
		lBidTrackIdArr, lErr2 := GetBidTrackIdSch(pBrokerId, lPendingBidsArr)
		if lErr2 != nil {
			log.Println("ISOAS02", lErr2)
			lRespRec.Status = common.ErrorCode
			lRespRec.ErrMsg = "ISOAS02" + lErr2.Error()
		} else {
			// log.Println("lPendingBidsArr", lPendingBidsArr)
			//todo:-------------------------------
			lFinalReqArr, lErr3 := ConstructPendingSch(lPendingBidsArr)
			if lErr3 != nil {
				log.Println("ISOAS03", lErr3)
				lRespRec.Status = common.ErrorCode
				lRespRec.ErrMsg = "ISOAS03" + lErr3.Error()
			} else {
				// ---------------->>>>> get the application no. ------------>>>>>>>>
				lAppArr, lErr4 := GetApiMaster(lFinalReqArr)
				if lErr4 != nil {
					log.Println("ISOAS04", lErr4)
					lRespRec.Status = common.ErrorCode
					lRespRec.ErrMsg = "ISOAS04" + lErr4.Error()
				} else {
					log.Println(lAppArr)
					lRespRec.TotalCount = len(lFinalReqArr)
					log.Println(lFinalReqArr, "lFinalReqArr")
					if lFinalReqArr != nil {
						//-------------------------------------------------------------->>>>>>>>>>>
						// Call ApplyIpo to pass the applications to Exchange.
						lRespRec.TotalCount = len(lFinalReqArr)
						lRespArr, lErr5 := exchangecall.ApplyNseIpo(lFinalReqArr, lUser, pBrokerId)
						if lErr5 != nil {
							log.Println("ISOAS05", lErr5)
							lRespRec.Status = common.ErrorCode
							lRespRec.ErrMsg = "ISOAS05" + lErr5.Error()
						} else {
							// log.Println(lRespArr)
							//call StatusValidationSch method to split the records by status and its give the finalized Array Value
							lFinalRespArr, lSuccessCount, lErrCount, lErr6 := StatusValidationSch(lRespArr, common.NSE)
							if lErr6 != nil {
								log.Println("ISOAS06", lErr6)
								lRespRec.ErrCount = lErrCount
								lRespRec.SuccessCount = lSuccessCount
								lRespRec.Status = common.ErrorCode
								lRespRec.ErrMsg = "ISOAS06" + lErr6.Error()
							} else {
								lRespRec.ErrCount = lErrCount
								lRespRec.SuccessCount = lSuccessCount
								log.Println("lFinalRespArr", lFinalRespArr)
								// this looping is used to request array and response array to update the records
								for lpReqArrIdx := 0; lpReqArrIdx < len(lFinalReqArr); lpReqArrIdx++ {
									for lpRespArrIdx := 0; lpRespArrIdx < len(lFinalRespArr); lpRespArrIdx++ {
										// check whether the Req application no and response application no is equal or not
										if lFinalReqArr[lpReqArrIdx].ApplicationNo == lFinalRespArr[lpRespArrIdx].ApplicationNo {
											//Call the UpdateRecordSch to Update the Response array
											lErr7 := UpdateRecordSch(lFinalReqArr[lpReqArrIdx], lFinalRespArr[lpRespArrIdx], lIdValuesArr, lBidTrackIdArr, lUser)
											if lErr7 != nil {
												log.Println("ISOAS07", lErr7)
												lRespRec.Status = common.ErrorCode
												lRespRec.ErrMsg = "ISOAS07" + lErr7.Error()
											} else {
												lRespRec.Status = common.SuccessCode
												lRespRec.ErrMsg = "Processed Successfully...."
											}
										}
									}
								}
								//------------------------>>>>>>>>>>>.emailsent <<<<<<<<<<<----------------
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
	log.Println("OfflineAppSch (-)")
	return lRespRec
}
