package iposchedule

// type ProgramStruct struct {
// 	Sno        string `json:"sno"`
// 	MethodName string `json:"methodName"`
// 	Status     string `json:"status"`
// 	ErrMsg     string `json:"errMsg"`
// }

// // this response struct for the fetchIpoMAster API
// type FetchRespStruct struct {
// 	ResponseArr []ProgramStruct `json:"responseArr"`
// 	Status      string          `json:"status"`
// 	ErrMsg      string          `json:"errMsg"`
// }

// /*
// Pupose:This Function is used to Get the Active Ipo Details from Exchange and insert it in our database table ....
// Parameters:

// 	not Applicable

// Response:

// 	*On Sucess
// 	=========
// 	{
// 		"IpoDetails": [
// 			{
// 				"id": 18,
// 				"symbol": "MMIPO26",
// 				"startDate": "2023-06-02",
// 				"endDate": "2023-06-30",
// 				"priceRange": "1000 - 2000",
// 				"cutOffPrice": 2000,
// 				"minBidQuantity": 10,
// 				"applicationStatus": "Pending",
// 				"upiStatus": "Accepted by Investor"
// 			},
// 			{
// 				"id": 10,
// 				"symbol": "fixed",
// 				"startDate": "2023-05-10",
// 				"endDate": "2023-08-29",
// 				"priceRange": "755 - 755",
// 				"cutOffPrice": 755,
// 				"minBidQuantity": 100,
// 				"applicationStatus": "-",
// 				"upiStatus": "-"
// 			}
// 		],
// 		"status": "S",
// 		"errMsg": ""
// 	}

// 	!On Error
// 	========
// 	{
// 		"status": "E",
// 		"reason": "Can't able to get the ipo details"
// 	}

// Author: Nithish Kumar
// Date: 08JUNE2023
// */
// func FetchIpoMasterSch(w http.ResponseWriter, r *http.Request) {
// 	log.Println("FetchIpoMasterSch (+)", r.Method)
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "GET" {

// 		//create instance for manualStruct
// 		var lRespRec FetchRespStruct

// 		var lPrgrmResp ProgramStruct
// 		//get header value
// 		lUser := r.Header.Get("USER")
// 		// validate the user and check if the user role is admin and then allow for next process
// 		lRespRec.Status = common.SuccessCode

// 		// Calling the FetchIpomaster method to get the Active Ipo details From exchange and
// 		// then store the details  in the database
// 		lErr1 := exchangecall.FetchIPOmaster(lUser)
// 		if lErr1 != nil {
// 			log.Println("ISFIMS01", lErr1)
// 			lPrgrmResp.Sno = "1"
// 			lPrgrmResp.MethodName = "Fetch-IPO-Master"
// 			lPrgrmResp.Status = common.ErrorCode
// 			lPrgrmResp.ErrMsg = "ISFIMS01" + lErr1.Error()
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 		} else {
// 			lPrgrmResp.Sno = "1"
// 			lPrgrmResp.MethodName = "Fetch-IPO-Master"
// 			lPrgrmResp.Status = common.SuccessCode
// 			lPrgrmResp.ErrMsg = common.SUCCESS
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 			log.Println("IPO Details Fetched Successfully")
// 		}

// 		//sgb master fetch in bse exchange
// 		lErr2 := exchangecall.FetchSGBMaster(lUser)
// 		if lErr2 != nil {
// 			log.Println("ISFIMS02", lErr2)
// 			lPrgrmResp.Sno = "2"
// 			lPrgrmResp.MethodName = "Fetch-SGB-Master"
// 			lPrgrmResp.Status = common.ErrorCode
// 			lPrgrmResp.ErrMsg = "ISFIMS02" + lErr2.Error()
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 		} else {
// 			lPrgrmResp.Sno = "2"
// 			lPrgrmResp.MethodName = "Fetch-SGB-Master"
// 			lPrgrmResp.Status = common.SuccessCode
// 			lPrgrmResp.ErrMsg = common.SUCCESS
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 			log.Println("SGB Details Fetched Successfully")
// 		}

// 		// Marshal the Response Structure into lData
// 		lData, lErr3 := json.Marshal(lRespRec)
// 		if lErr3 != nil {
// 			log.Println("ISFIMS03", lErr3)
// 			fmt.Fprintf(w, "ISFIMS03"+lErr3.Error())
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("FetchIpoMasterSch (-)", r.Method)
// 	}
// }

//  copied From API27 Fix
type ProgramStruct struct {
	Sno        string `json:"sno"`
	MethodName string `json:"methodName"`
	Status     string `json:"status"`
	ErrMsg     string `json:"errMsg"`
}

// this response struct for the fetchIpoMAster API
type FetchRespStruct struct {
	ResponseArr []ProgramStruct `json:"responseArr"`
	Status      string          `json:"status"`
	ErrMsg      string          `json:"errMsg"`
}

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
// func FetchIpoMasterSch(w http.ResponseWriter, r *http.Request) {
// 	log.Println("FetchIpoMasterSch (+)", r.Method)
// 	origin := r.Header.Get("Origin")
// 	var lBrokerId int
// 	var lErr error
// 	w.Header().Set("Access-Control-Allow-Origin", origin)
// 	lBrokerId, lErr = brokers.GetBrokerId(origin) // TO get brokerId
// 	if lErr != nil {
// 		log.Println(lErr, origin)
// 	}
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	(w).Header().Set("Access-Control-Allow-Headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")

// 	if r.Method == "GET" {

// 		//create instance for manualStruct
// 		var lRespRec FetchRespStruct

// 		var lPrgrmResp ProgramStruct
// 		//get header value
// 		lUser := r.Header.Get("USER")
// 		// validate the user and check if the user role is admin and then allow for next process
// 		lRespRec.Status = common.SuccessCode

// 		// Calling the FetchIpomaster method to get the Active Ipo details From exchange and
// 		// then store the details  in the database
// 		lNoToken := "Access Token not found"
// 		lNSEToken, lErr1 := exchangecall.FetchNseIPOmaster(lUser, lBrokerId)
// 		if lErr1 != nil {
// 			log.Println("ISFIMS01", lErr1)
// 			lPrgrmResp.Sno = "1"
// 			lPrgrmResp.MethodName = "NSE Fetch-IPO-Master"
// 			lPrgrmResp.Status = common.ErrorCode
// 			lPrgrmResp.ErrMsg = "ISFIMS01" + lErr1.Error()
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 		} else {
// 			if lNSEToken == common.ErrorCode {
// 				lPrgrmResp.Sno = "1"
// 				lPrgrmResp.MethodName = "NSE Fetch-IPO-Master"
// 				lPrgrmResp.Status = common.ErrorCode
// 				lPrgrmResp.ErrMsg = lNoToken
// 				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 				log.Println("NSE IPO Details FETCHING ERROR")

// 			} else {

// 				lPrgrmResp.Sno = "1"
// 				lPrgrmResp.MethodName = "NSE Fetch-IPO-Master"
// 				lPrgrmResp.Status = common.SuccessCode
// 				lPrgrmResp.ErrMsg = common.SUCCESS
// 				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 				log.Println("NSE IPO Details Fetched Successfully")
// 			}
// 		}

// 		lBseToken, lErr2 := exchangecall.FetchBseIPOmaster(lUser, lBrokerId)
// 		if lErr2 != nil {
// 			log.Println("ISFIMS02", lErr2)
// 			lPrgrmResp.Sno = "2"
// 			lPrgrmResp.MethodName = "BSE Fetch-IPO-Master"
// 			lPrgrmResp.Status = common.ErrorCode
// 			lPrgrmResp.ErrMsg = "ISFIMS02" + lErr2.Error()
// 			lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 		} else {
// 			if lBseToken == common.ErrorCode {
// 				lPrgrmResp.Sno = "2"
// 				lPrgrmResp.MethodName = "BSE Fetch-IPO-Master"
// 				lPrgrmResp.Status = common.ErrorCode
// 				lPrgrmResp.ErrMsg = lNoToken
// 				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 				log.Println("BSE IPO Details Fetching Error")
// 			} else {

// 				lPrgrmResp.Sno = "2"
// 				lPrgrmResp.MethodName = "BSE Fetch-IPO-Master"
// 				lPrgrmResp.Status = common.SuccessCode
// 				lPrgrmResp.ErrMsg = common.SUCCESS
// 				lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 				log.Println("BSE IPO Details Fetched Successfully")
// 			}
// 		}

//sgb master fetch in bse exchange
// lErr3 := exchangecall.FetchSGBMaster(lUser)
// if lErr3 != nil {
// 	log.Println("ISFIMS02", lErr3)
// 	lPrgrmResp.Sno = "3"
// 	lPrgrmResp.MethodName = "BSE Fetch-SGB-Master"
// 	lPrgrmResp.Status = common.ErrorCode
// 	lPrgrmResp.ErrMsg = "ISFIMS03" + lErr3.Error()
// 	lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// } else {
// 	lPrgrmResp.Sno = "3"
// 	lPrgrmResp.MethodName = "BSE Fetch-SGB-Master"
// 	lPrgrmResp.Status = common.SuccessCode
// 	lPrgrmResp.ErrMsg = common.SUCCESS
// 	lRespRec.ResponseArr = append(lRespRec.ResponseArr, lPrgrmResp)
// 	log.Println("BSE SGB Details Fetched Successfully")
// }

// Marshal the Response Structure into lData
// 		lData, lErr4 := json.Marshal(lRespRec)
// 		if lErr4 != nil {
// 			log.Println("ISFIMS03", lErr4)
// 			fmt.Fprintf(w, "ISFIMS03"+lErr4.Error())
// 		} else {
// 			fmt.Fprintf(w, string(lData))
// 		}
// 		log.Println("FetchIpoMasterSch (-)", r.Method)
// 	}
// }
// the above code commented by pavithra fetchmaster api repeated twice
