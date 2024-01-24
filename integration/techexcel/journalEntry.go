// package incentive
package techexcel

import (
	"encoding/json"
	"fcs23pkg/apigate"
	"fcs23pkg/ftdb"
	"fcs23pkg/genpkg"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type JvInputStruct struct {
	URL               string
	COCD              string
	VoucherDate       string
	AccountCode       string
	Amount            string
	Narration         string
	JvNarration       string
	TDSNarration      string
	TDSPct            string
	EMPINCAccount     string
	TDSAccount        string
	IncSegment        string
	CounterAccount    string
	BillNo            string
	WithGST           string
	UrlUserName       string
	UrlPassword       string
	UrlDatabase       string
	UrlDataYear       string
	TeXLVoucherDate   string
	TeXLBillNo        string
	TeXLVoucherNo     string
	TeXLMessage       string
	TeXLMessageType   string
	TexXLResponseBody string
	SourceTable       string
	SourceTableKey    string
}

type jvResponseErrType struct {
	Message string `json:"MESSAGE"`
	Type    string `json:"TYPE"`
}
type jvResponseType struct {
	Columns []string          `json:"COLUMNS"`
	Data    [][]string        `json:"DATA"`
	Message jvResponseErrType `json:"MES"`
}

// --------------------------------------------------------------------
//
//	function to log JV call details
//
// --------------------------------------------------------------------
func logJVEntry(reqDtl apigate.RequestorDetails, jvData JvInputStruct, Url string, respJson jvResponseType) {
	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()

		inputJson, err := json.Marshal(jvData)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		responseJson, err := json.Marshal(respJson)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		insertString := `insert into xxtechxl_jv_log (Narration, COCD,VoucherDate,AccountCode,Amount
													,CounterAccount,SourceTable, SourceTable_Key, realip
													,forwardedip,methods,paths,host,remoteaddr
													,header,body,endpoint, TeXLMessage, TeXLMessageType,  Input_JSON, Request_JSON, Response_JSON,  createddate)
												values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())`

		_, inserterr := db.Exec(insertString, jvData.Narration, jvData.COCD, jvData.VoucherDate, jvData.AccountCode, jvData.Amount,
			jvData.CounterAccount, jvData.SourceTable, jvData.SourceTableKey, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path, reqDtl.Host,
			reqDtl.RemoteAddr, reqDtl.Header, jvData.TexXLResponseBody, reqDtl.EndPoint, jvData.TeXLMessage, jvData.TeXLMessageType, string(inputJson), Url, string(responseJson))
		if inserterr != nil {
			log.Println(inserterr)
			log.Println("logJVEntry", "Error inserting xxtechxl_jv_log msg : "+inserterr.Error())
		}
	}
}

func techExcelGLAPI(APIMeta JvInputStruct) (JvInputStruct, string, jvResponseType) {
	log.Println("techExcelGLAPI (+)")
	var jvInput JvInputStruct
	var jvResponse jvResponseType
	jvInput = APIMeta
	client := http.DefaultClient

	// url := APIMeta.URL + "COCD=" + APIMeta.COCD + "&VOUCHER_DATE=" + APIMeta.VoucherDate + "&ACCOUNTCODE=" + APIMeta.AccountCode + "&AMOUNT=" + APIMeta.Amount + "&NARRATION=" + APIMeta.Narration + "&CntACCOUNTCODE=" + APIMeta.CounterAccount + "&BILLNO=" + APIMeta.BillNo + "&withgst=" + APIMeta.WithGST + "&UrlUserName=" + APIMeta.UrlUserName + "&UrlPassword=" + APIMeta.UrlPassword + "&UrlDatabase=" + APIMeta.UrlDatabase + "&UrlDataYear=" + APIMeta.UrlDataYear
	// log.Println("URL: " + url)

	// url = strings.ReplaceAll(url, " ", "")

	// req, err := http.NewRequest(http.MethodGet, url, nil)

	//log.Println("URL", APIMeta.URL)

	req, err := http.NewRequest(http.MethodGet, APIMeta.URL, nil)

	//Param Construction
	q := req.URL.Query()
	q.Add("COCD", APIMeta.COCD)
	q.Add("VOUCHER_DATE", APIMeta.VoucherDate)
	q.Add("ACCOUNTCODE", APIMeta.AccountCode)
	q.Add("AMOUNT", APIMeta.Amount)
	q.Add("NARRATION", APIMeta.Narration)
	q.Add("CntACCOUNTCODE", APIMeta.CounterAccount)
	q.Add("BILLNO", APIMeta.BillNo)

	q.Add("withgst", APIMeta.WithGST)
	q.Add("UrlUserName", APIMeta.UrlUserName)
	q.Add("UrlPassword", APIMeta.UrlPassword)
	q.Add("UrlDatabase", APIMeta.UrlDatabase)
	q.Add("UrlDataYear", APIMeta.UrlDataYear)

	req.URL.RawQuery = q.Encode()

	url := req.URL.RawQuery

	log.Println("Voucher Date: ", APIMeta.VoucherDate)

	log.Println("URL", url)

	//resp, err := http.Get("http://192.168.150.21:8686/techexcelapi/index.cfm/JOURNAL/JOURNAL?COCD=NSE_CASH&VOUCHER_DATE=26/11/2021&ACCOUNTCODE=JBMU15&AMOUNT=1000&NARRATION=FUND TRANSFER FROM NSECASH To MTF&CntACCOUNTCODE=MTF_NSE_CASH&BILLNO=MTFFundTrf&withgst=N&UrlUserName=api&UrlPassword=api@123&UrlDatabase=capsfo&UrlDataYear=2021")
	if err != nil {
		log.Println(err)
		jvInput.TeXLMessage = err.Error()
		jvInput.TeXLMessageType = "ERROR"
	} else {
		response, err := client.Do(req)
		if err != nil {
			jvInput.TeXLMessage = err.Error()
			jvInput.TeXLMessageType = "ERROR"
		} else {
			//We Read the response body on the line below.
			body, err := ioutil.ReadAll(response.Body)
			//log.Println(body)
			if err != nil {
				log.Println(err)
				jvInput.TeXLMessage = err.Error()
				jvInput.TeXLMessageType = "ERROR"
			} else {
				log.Println(string(body[:]))
				err := json.Unmarshal(body, &jvResponse)
				if err != nil {
					jvInput.TexXLResponseBody = string(body[:])
					log.Println("Unable to Unmarshal request")
					log.Println(err)

					htmlErrMsg := ParseHTMLError(string(body[:]))
					if htmlErrMsg != "" {
						jvInput.TeXLMessage = htmlErrMsg
						log.Println("htmlErrMsg", htmlErrMsg)
						jvInput.TeXLMessageType = "ERROR"
					} else {
						jvInput.TeXLMessage = err.Error()
						jvInput.TeXLMessageType = "ERROR"
					}
				} else {

					for i := 0; i < len(jvResponse.Data); i++ {
						for j := 0; j < len(jvResponse.Data[i]); j++ {
							if jvResponse.Columns[j] == "VOUCHER_DATE" {

								jvInput.TeXLVoucherDate = jvResponse.Data[i][j]
							}
							if jvResponse.Columns[j] == "BILLNO" {

								jvInput.TeXLBillNo = jvResponse.Data[i][j]
							}
							if jvResponse.Columns[j] == "VNO" {

								jvInput.TeXLVoucherNo = jvResponse.Data[i][j]
							}
						}
					}

					jvInput.TeXLMessage = jvResponse.Message.Message
					jvInput.TeXLMessageType = jvResponse.Message.Type
				}
			}
		}

		//	{"COLUMNS":["VOUCHER_DATE","BILLNO","VNO"],"DATA":[["26/11/2021","MTFFundTrf","JV07317048"]]}
		//	{"MES":{"MESSAGE":"Please Defiend DataYear","TYPE":"error"}}
	}
	log.Println("techExcelGLAPI (-)")
	return jvInput, url, jvResponse
}

// --------------------------------------------------------------------
//
// public function to log journal entry
//
// --------------------------------------------------------------------
func JVPostingtoTechxl_UAT(jvMeta JvInputStruct, reqDtl apigate.RequestorDetails) error {
	techXLconfig := genpkg.ReadTomlConfig("./toml/techXLAPI_UAT.toml")
	jvMeta.URL = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["JournalAPI"])
	jvMeta.UrlUserName = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlUserName"])
	jvMeta.UrlPassword = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlPassword"])
	jvMeta.UrlDatabase = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDatabase"])
	jvMeta.UrlDataYear = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDataYear"])

	log.Println("jvMeta", jvMeta)
	FirstjvResponse, reqUrl, respJson := techExcelGLAPI(jvMeta)
	log.Println(FirstjvResponse)
	logJVEntry(reqDtl, FirstjvResponse, reqUrl, respJson)
	if strings.ToUpper(FirstjvResponse.TeXLMessageType) == "ERROR" {
		return fmt.Errorf("Something went wrong during jv posting (0x127) for " + jvMeta.AccountCode)
	}

	return nil
}

// --------------------------------------------------------------------
//
// public function to log journal entry
//
// --------------------------------------------------------------------
func JVPostingtoTechxl(jvMeta JvInputStruct, reqDtl apigate.RequestorDetails) error {
	techXLconfig := genpkg.ReadTomlConfig("./toml/techXLAPI.toml")
	jvMeta.URL = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["JournalAPI"])
	jvMeta.UrlUserName = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlUserName"])
	jvMeta.UrlPassword = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlPassword"])
	jvMeta.UrlDatabase = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDatabase"])
	jvMeta.UrlDataYear = fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDataYear"])

	FirstjvResponse, reqUrl, respJson := techExcelGLAPI(jvMeta)
	log.Println(FirstjvResponse)
	logJVEntry(reqDtl, FirstjvResponse, reqUrl, respJson)
	if strings.ToUpper(FirstjvResponse.TeXLMessageType) == "ERROR" {
		return fmt.Errorf("Something went wrong during jv posting (0x127) for " + jvMeta.AccountCode)
	}

	return nil
}
