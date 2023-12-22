package ncbschedule

import (
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/nse/nsencb"
	"log"
	"net/http"
)

type JvNcbReqStruct struct {
	OrderNo           int    `json:"orderno"`
	ApplicationNumber string `json:"applicationnumber"`
	ClientId          string `json:"clientid"`
	JvStatus          string `json:"jvstatus"`
	JvStatement       string `json:"Jvstatement"`
	JvAmount          string `json:"jvamount"`
	JvType            string `json:"jvtype"`
	Unit              string `json:"unit"`
	Price             string `json:"price"`
	ActionCode        string `json:"actioncode"`
	Symbol            string `json:"symbol"`
	OrderDate         string `json:"orderdate"`
	Amount            string `json:"amount"`
	Mail              string `json:"mail"`
	ClientName        string `json:"clientname"`
}

func PlacingNcbOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("PlacingNcbOrder(+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	lValidNcbArr, lErr1 := GoodTimeForApplyNcb()
	log.Println("lValidNcbArr", lValidNcbArr)
	if lErr1 != nil {
		log.Println("PNO01", lErr1)
	} else {
		log.Println("NCB Valid Bond array", lValidNcbArr)
		if len(lValidNcbArr) != 0 {

			for lNcbIdx := 0; lNcbIdx < len(lValidNcbArr); lNcbIdx++ {

				lBrokerIdArr, lErr2 := GetNcbBrokers(lValidNcbArr[lNcbIdx].Exchange)
				if lErr2 != nil {
					log.Println("PNO02", lErr2)
				} else {
					log.Println("NCB Valid Bond array", lBrokerIdArr)

					for lBrokerIdx := 0; lBrokerIdx < len(lBrokerIdArr); lBrokerIdx++ {

						jvsuccessCount, jvfailedCount, ExchangeSuccess, ExchangeFailed, ReverJv, lErr2 := ProcessNcbOrder(lValidNcbArr[lNcbIdx], lBrokerIdArr[lBrokerIdx], r)
						log.Println("jvsuccessCount, jvfailedCount, ExchangeSuccess, ExchangeFailed", jvsuccessCount, jvfailedCount, ExchangeSuccess, ExchangeFailed, ReverJv, lErr2)
					}
				}

			}

		} else {
			log.Println("No NCB ending today")
		}

	}

	log.Println("PlacingNcbOrder(-)")

}

//----------------------------------------------------------------
// this method is used to get the valid NCB record from database
//----------------------------------------------------------------
func GoodTimeForApplyNcb() ([]NcbStruct, error) {
	log.Println("GoodTimeForApplyNcb(+)")

	var lValidNcbArr []NcbStruct
	var lValidNcbRec NcbStruct

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("GTAN01", lErr1)
		return lValidNcbArr, lErr1
	} else {
		defer lDb.Close()

		lCoreString := `select n.id, n.Symbol, n.Exchange
		                from a_ncb_master n                     
						where n.BiddingEndDate  = curdate() 
						and n.DailyEndTime > now()`
		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("GTAN02", lErr2)
			return lValidNcbArr, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&lValidNcbRec.MasterId, &lValidNcbRec.Symbol, &lValidNcbRec.Exchange)
				if lErr3 != nil {
					log.Println("GTAN03", lErr3)
					return lValidNcbArr, lErr3
				} else {
					lValidNcbArr = append(lValidNcbArr, lValidNcbRec)
					log.Println("lValidNcbArr123123", lValidNcbArr)
				}
			}
		}
	}

	log.Println("GoodTimeForApplyNcb(-)")
	return lValidNcbArr, nil
}

func ProcessNcbOrder(pValidNcb NcbStruct, pBrokerId int, pApiRequest *http.Request) (int, int, int, int, int, error) {
	log.Println("ProcessNcbOrder (+)")

	var ljvsuccessCount int
	var ljvfailedCount int
	var lExchangeSuccess int
	var lExchangeFailed int
	var lReverseJv int

	lNcbReqArr, lJvDetailArr, lErr1 := fetchPendingNcbOrder(pValidNcb.Exchange, pBrokerId, pValidNcb.MasterId)
	if lErr1 != nil {
		log.Println("PSO01", lErr1)

	} else {

		for lNqlReqIdx := 0; lNqlReqIdx < len(lNcbReqArr); lNqlReqIdx++ {

			for lJvReqIdx := 0; lJvReqIdx < len(lJvDetailArr); lJvReqIdx++ {

				if lNcbReqArr[lNqlReqIdx].OrderNumber == lJvDetailArr[lJvReqIdx].OrderNo {

					jvsuccessCount, jvfailedCount, ExchangeSuccess, ExchangeFailed, ReverseJv, lErr2 := PostJvForOrder(lNcbReqArr[lNqlReqIdx], lJvDetailArr[lJvReqIdx], pApiRequest, pValidNcb, pBrokerId)
					if lErr2 != nil {
						log.Println("lErr2", lErr2)
						return ljvsuccessCount, ljvfailedCount, lExchangeSuccess, lExchangeFailed, ReverseJv, lErr2
					} else {
						ljvsuccessCount += jvsuccessCount
						ljvfailedCount += jvfailedCount
						lExchangeSuccess += ExchangeSuccess
						lExchangeFailed += ExchangeFailed
						lReverseJv += ReverseJv
					}
				}
			}

		}
	}

	log.Println("ProcessNcbOrder (-)")
	return ljvsuccessCount, ljvfailedCount, lExchangeSuccess, lExchangeFailed, lReverseJv, nil
}

func fetchPendingNcbOrder(pExchange string, BrokerId int, pMasterId int) ([]nsencb.NcbAddReqStruct, []JvNcbReqStruct, error) {
	log.Println("fetchPendingNcbOrder (+)")

	var lReqNcbData nsencb.NcbAddReqStruct
	var lReqJVData JvNcbReqStruct
	var lReqNcbDataArr []nsencb.NcbAddReqStruct
	var lReqJVDataArr []JvNcbReqStruct

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("FPNO01", lErr1)
		return lReqNcbDataArr, lReqJVDataArr, lErr1
	} else {
		defer lDb.Close()

		lSqlString := `select h.Symbol , h.Investmentunit ,h.applicationNo, d.OrderNo, d.price, h.PhysicalDematFlag, h.pan, h.depository,
		h.dpId, h.ClientRefNumber,h.clientBenId , h.ClientEmail, h.clientId, date_format(d.CreatedDate,'%d-%m-%Y') AS formatted_date,n.Symbol, h.Investmentunit,
				   (CASE
					 WHEN d.activityType  = 'M' THEN 'Modify'
					 WHEN d.activityType = 'N' THEN 'New'
					 WHEN d.activityType = 'D' THEN 'Delete'
					ELSE d.activityType
					END ) AS ActionDescription
				  from a_ncb_master n, a_ncb_orderdetails d, a_ncb_orderheader h
				  where h.MasterId  = n.id
				  and d.headerId  = h.Id
				  and h.status  = "success"
				  and n.BiddingStartDate  <= curdate()
				  and n.BiddingEndDate  >= curdate()
				  and time(now()) between n.DailyStartTime and n.DailyEndTime
				  and h.cancelFlag != 'Y'
				  and h.Exchange = ?
				  and h.brokerId = ?
				  and n.id = ?`

		lRows, lErr2 := lDb.Query(lSqlString, pExchange, BrokerId, pMasterId)
		if lErr2 != nil {
			log.Println("FPNO02", lErr2)
			return lReqNcbDataArr, lReqJVDataArr, lErr2
		} else {

			for lRows.Next() {
				lErr3 := lRows.Scan(&lReqNcbData.Symbol, &lReqNcbData.InvestmentValue, &lReqNcbData.ApplicationNumber, &lReqNcbData.OrderNumber, &lReqJVData.Amount, &lReqNcbData.PhysicalDematFlag, &lReqNcbData.Pan, &lReqNcbData.Depository, &lReqNcbData.DpId, &lReqNcbData.ClientRefNumber, &lReqNcbData.ClientBenId, &lReqJVData.Mail, &lReqJVData.ClientId, &lReqJVData.OrderDate, &lReqJVData.Symbol, &lReqJVData.Unit, &lReqJVData.ActionCode)

				if lErr3 != nil {
					log.Println("FPNO03", lErr3)
					return lReqNcbDataArr, lReqJVDataArr, lErr3
				} else {
					lReqJVData.OrderNo = lReqNcbData.OrderNumber
					lReqNcbDataArr = append(lReqNcbDataArr, lReqNcbData)
					lReqJVDataArr = append(lReqJVDataArr, lReqJVData)
				}
			}

		}

	}

	log.Println("fetchPendingNcbOrder (-)")
	return lReqNcbDataArr, lReqJVDataArr, nil
}
