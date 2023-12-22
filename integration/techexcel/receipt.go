package techexcel

import (
	"encoding/json"
	"errors"
	"fcs23pkg/apigate"
	"fcs23pkg/ftdb"
	"fcs23pkg/genpkg"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func CreateReceiptUAT(receiptRec ReceiptStruct, reqDtl apigate.RequestorDetails) error {
	// fmt.Println("config Data : ", config)
	techXLconfig := genpkg.ReadTomlConfig("./toml/techXLAPI_UAT.toml")
	ReceiptAPI := fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["ReceiptAPI"])
	request, err := http.NewRequest("POST", ReceiptAPI, nil)
	// request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Println("EC05", err)
		return errors.New("EC05")
	}
	q := request.URL.Query()
	q.Add("VoucherDate", receiptRec.VoucherDate)
	q.Add("AccountCode", receiptRec.AccountCode)
	q.Add("COMPANYCODE", receiptRec.COMPANYCODE)
	q.Add("PAYMENTREFERENCENUMBER", receiptRec.PAYMENTREFERENCENUMBER)
	q.Add("Amount", receiptRec.Amount)
	q.Add("PostingBankAccount", receiptRec.PostingBankAccount)
	q.Add("BankAccountNumber", receiptRec.BankAccountNumber)
	q.Add("NARRATION", receiptRec.NARRATION)
	q.Add("ENTRYTYPE", receiptRec.ENTRYTYPE)
	q.Add("UrlUserName", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlUserName"]))
	q.Add("UrlPassword", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlPassword"]))
	q.Add("UrlDatabase", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDatabase"]))
	q.Add("UrlDataYear", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDataYear"]))

	request.URL.RawQuery = q.Encode()
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("EC06", err)
		return errors.New("EC06")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	LogReceiptAPICall(receiptRec, request.URL.RawQuery, string(body), reqDtl)
	if err != nil {
		log.Println("EC08", err)
		return errors.New("EC08")
	}
	var receiptResponse ReceiptResponseStruct
	err = json.Unmarshal([]byte(string(body)), &receiptResponse)
	if err != nil {
		log.Println("EC09", err)
		return errors.New("EC09")
	}
	// Sucess:0,Error:Entry already punched - receiptResponse.Data[0][0]
	if receiptResponse.Message.Type == "ERROR" {
		log.Println("EC18", receiptResponse.Message)
		return errors.New("EC18")
	} else if len(receiptResponse.Data) > 0 {
		if strings.Contains(receiptResponse.Data[0][0], "Entry already punched") {
			log.Println("EC10", err)
			return errors.New("EC19")
		}
	}

	return nil
}

func CreateReceipt(receiptRec ReceiptStruct, reqDtl apigate.RequestorDetails) error {
	// fmt.Println("config Data : ", config)
	techXLconfig := genpkg.ReadTomlConfig("./toml/techXLAPI.toml")
	ReceiptAPI := fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["ReceiptAPI"])
	request, err := http.NewRequest("POST", ReceiptAPI, nil)
	// request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Println("EC05", err)
		return errors.New("EC05")
	}
	q := request.URL.Query()
	q.Add("VoucherDate", receiptRec.VoucherDate)
	q.Add("AccountCode", receiptRec.AccountCode)
	q.Add("COMPANYCODE", receiptRec.COMPANYCODE)
	q.Add("PAYMENTREFERENCENUMBER", receiptRec.PAYMENTREFERENCENUMBER)
	q.Add("Amount", receiptRec.Amount)
	q.Add("PostingBankAccount", receiptRec.PostingBankAccount)
	q.Add("BankAccountNumber", receiptRec.BankAccountNumber)
	q.Add("NARRATION", receiptRec.NARRATION)
	q.Add("ENTRYTYPE", receiptRec.ENTRYTYPE)
	q.Add("UrlUserName", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlUserName"]))
	q.Add("UrlPassword", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlPassword"]))
	q.Add("UrlDatabase", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDatabase"]))
	q.Add("UrlDataYear", fmt.Sprintf("%v", techXLconfig.(map[string]interface{})["UrlDataYear"]))

	request.URL.RawQuery = q.Encode()
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("EC06", err)
		return errors.New("EC06")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	LogReceiptAPICall(receiptRec, request.URL.RawQuery, string(body), reqDtl)
	if err != nil {
		log.Println("EC08", err)
		return errors.New("EC08")
	}
	var receiptResponse ReceiptResponseStruct
	err = json.Unmarshal([]byte(string(body)), &receiptResponse)
	if err != nil {
		log.Println("EC09", err)
		return errors.New("EC09")
	}
	// Sucess:0,Error:Entry already punched - receiptResponse.Data[0][0]
	if receiptResponse.Message.Type == "ERROR" {
		log.Println("EC18", receiptResponse.Message)
		return errors.New("EC18")
	} else if len(receiptResponse.Data) > 0 {
		if strings.Contains(receiptResponse.Data[0][0], "Entry already punched") {
			log.Println("EC10", err)
			return errors.New("EC19")
		}
	}

	return nil
}

func LogReceiptAPICall(receiptRec ReceiptStruct, Url string, respJson string, reqDtl apigate.RequestorDetails) {

	db, err := ftdb.LocalDbConnect(ftdb.MariaFTPRD)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()
		inputJson, err := json.Marshal(receiptRec)
		if err != nil {
			log.Println("EC07", err)
		}

		strQry := `insert into xxtechxl_receipt_log(SourceTable,SourceTable_Key,Narration,COMPANYCODE,VoucherDate,Amount,PostingBankAccount,
				BankAccountNumber,ENTRYTYPE,realip,forwardedip,methods,paths,host,remoteaddr,header,body,endpoint,Request_JSON,
				Response_JSON,createdBy,createdDate,createdProgram,updatedBy,updatedDate,updatedProgram)
		 		values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,'App',now(),'ATOM','App',now(),'ATOM')`

		_, err = db.Exec(strQry, receiptRec.SourceTable, receiptRec.SourceKeyId, receiptRec.NARRATION, receiptRec.COMPANYCODE, receiptRec.VoucherDate, receiptRec.Amount,
			receiptRec.PostingBankAccount, receiptRec.BankAccountNumber, receiptRec.ENTRYTYPE, reqDtl.RealIP, reqDtl.ForwardedIP, reqDtl.Method, reqDtl.Path,
			reqDtl.Host, reqDtl.RemoteAddr, reqDtl.Header, reqDtl.Body, reqDtl.EndPoint, inputJson, respJson)

		if err != nil {
			log.Println("EC08", err)
		}
	}
}
