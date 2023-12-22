package sgbschedule

// this response struct for the fetchIpoMAster API
// type SgbMasterRespStruct struct {
// 	Status string `json:"status"`
// 	ErrMsg string `json:"errMsg"`
// }

/*
Pupose:This Function is used to Get the Active Ipo Details from Exchange and insert it in our database table ....
Parameters:

	not Applicable

Response:

	*On Sucess
	=========
	{
		"IpoDetails": [
			{
				"id": 18,
				"symbol": "MMIPO26",
				"startDate": "2023-06-02",
				"endDate": "2023-06-30",
				"priceRange": "1000 - 2000",
				"cutOffPrice": 2000,
				"minBidQuantity": 10,
				"applicationStatus": "Pending",
				"upiStatus": "Accepted by Investor"
			},
			{
				"id": 10,
				"symbol": "fixed",
				"startDate": "2023-05-10",
				"endDate": "2023-08-29",
				"priceRange": "755 - 755",
				"cutOffPrice": 755,
				"minBidQuantity": 100,
				"applicationStatus": "-",
				"upiStatus": "-"
			}
		],
		"status": "S",
		"errMsg": ""
	}

	!On Error
	========
	{
		"status": "E",
		"reason": "Can't able to get the ipo details"
	}

Author: Nithish Kumar
Date: 08JUNE2023
*/
// func FetchSgbMasterSch(w http.ResponseWriter, r *http.Request) {
// 	log.Println("FetchSgbMasterSch (+)", r.Method)
// 	origin := r.Header.Get("Origin")
// 	var lBrokerId int
// 	var lErr error
// 	for _, allowedOrigin := range common.ABHIAllowOrigin {
// 		if allowedOrigin == origin {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
// 			log.Println(lErr, origin)
// 			break
// }
// }
// (w).Header().Set("Access-Control-Allow-Credentials", "true")
// (w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// (w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// if r.Method == "GET" {

//create instance for manualStruct
// var lRespRec SgbMasterRespStruct

//get header value
// lUser := r.Header.Get("USER")
// validate the user and check if the user role is admin and then allow for next process
// lRespRec.Status = common.SuccessCode

// Calling the FetchIpomaster method to get the Active Ipo details From exchange and
// then store the details  in the database
// lErr1 := exchangecall.FetchSGBMaster(lUser, lBrokerId)
// if lErr1 != nil {
// 	log.Println("SSFSM01", lErr1)
// 	lRespRec.Status = common.ErrorCode
// 	lRespRec.ErrMsg = "SSFSM01" + lErr1.Error()
// } else {
// 	log.Println("SGB Details Fetched Successfully")
// }

// Marshal the Response Structure into lData
// 		lData, lErr2 := json.Marshal(lRespRec)
// 		if lErr2 != nil {
// 			log.Println("SSFSM02", lErr2)
// 			fmt.Fprintf(w, "SSFSM02"+lErr2.Error())
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("FetchSgbMasterSch (-)", r.Method)
// 	}
// }
