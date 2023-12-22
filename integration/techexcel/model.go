package techexcel

type ReceiptStruct struct {
	VoucherDate            string
	AccountCode            string
	COMPANYCODE            string
	PAYMENTREFERENCENUMBER string
	Amount                 string
	PostingBankAccount     string
	BankAccountNumber      string
	NARRATION              string
	ENTRYTYPE              string
	UrlUserName            string
	UrlPassword            string
	UrlDatabase            string
	UrlDataYear            string
	SourceTable            string
	SourceKeyId            string
}

type receiptResponseErrType struct {
	Message string `json:"MESSAGE"`
	Type    string `json:"TYPE"`
}

type ReceiptResponseStruct struct {
	Columns []string               `json:"COLUMNS"`
	Data    [][]string             `json:"DATA"`
	Message receiptResponseErrType `json:"MES"`
}
