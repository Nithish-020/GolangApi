package emailUtil

import (
	"errors"
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// VoucherNo
// VoucherDate
// ClientId
// Amount

type EmailInput struct {
	//From        string
	FromRaw      string
	FromDspName  string
	ReplyTo      string
	ToEmailId    string
	Subject      string
	Body         string
	CampaignName string
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func SendEmail(input EmailInput, ReqSource string) error {
	log.Println("SendEmail (+)")

	config := common.ReadTomlConfig("../emailconfig.toml")

	//input.ToEmailId = "prabhaharan.s@fcsonline.co.in"

	account := fmt.Sprintf("%v", config.(map[string]interface{})["Account"])
	pwd := fmt.Sprintf("%v", config.(map[string]interface{})["Pwd"])
	//	ConfigurationSet := fmt.Sprintf("%v", config.(map[string]interface{})["ConfigurationSet"])
	// from := fmt.Sprintf("%v", config.(map[string]interface{})["From"])
	// fromraw := fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
	// replyto := fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
	// //subject := fmt.Sprintf("%v", config.(map[string]interface{})["Subject"])
	url := fmt.Sprintf("%v", config.(map[string]interface{})["Url"])

	// BCC mail for All Novo Emails
	lBccMailId := fmt.Sprintf("%v", config.(map[string]interface{})["Bcc_emailId"])
	Bcc := "Bcc: " + lBccMailId

	//mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	FcsRefer := ""
	if input.CampaignName != "" {
		FcsRefer = "FCS-REFERENCE-ID: " + input.CampaignName + "\n"
		log.Println(FcsRefer, "fcs")
	}

	// msg := "From: " + input.FromDspName + "<" + input.FromRaw + ">\n" +
	// 	"To: " + input.ToEmailId + "\n" +
	// 	"reply-to: " + input.ReplyTo + "\n" +
	// 	//"X-SES-CONFIGURATION-SET: " + ConfigurationSet + "\n" +
	// 	FcsRefer +
	// 	"Subject: " + input.Subject + "\n" + mime +
	// 	input.Body

	//--------------------------------------
	to := strings.Split(input.ToEmailId, ",")
	toHeader := "To: " + strings.Join(to, ",") + "\n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// commented by pavithra, added BCC in the below code
	// msg := "From: " + input.FromDspName + "\n" +
	// 	toHeader +
	// 	//bcc +
	// 	//cc +
	// 	"reply-to: " + input.ReplyTo + "\n" +
	// 	"Subject: " + input.Subject + "\n" + mime +
	// 	input.Body
	msg := "From: " + input.FromDspName + "\n" +
		toHeader + Bcc + "\n" +
		"reply-to: " + input.ReplyTo + "\n" +
		"Subject: " + input.Subject + "\n" + mime +
		input.Body

		// ===============Added Bcc Mail Content By Prashanth==============================

	if lBccMailId != "" {
		vcc := strings.Split(lBccMailId, ",")
		if len(vcc) > 0 {
			for i := 0; i < len(vcc); i++ {
				to = append(to, vcc[i])
			}
		}
	}
	// ===========================================================================

	// log.Println(msg)
	//--------------------------------------
	auth := LoginAuth(account, pwd)

	//err := smtp.SendMail(url, auth, input.FromRaw, []string{input.ToEmailId}, []byte(msg))
	err := smtp.SendMail(url, auth, input.FromRaw, to, []byte(msg))
	if err != nil {
		log.Println("EUSE01", err)
		return err
	} else {
		err := emailLog(input, ReqSource, url)
		if err != nil {
			log.Println("EUSE02", err)
			return err
		}
	}
	log.Println("SendEmail (-)")
	return nil
}

func emailLog(input EmailInput, ReqSource string, Url string) error {
	log.Println("emailLog (+)")

	db, err := ftdb.LocalDbConnect(ftdb.IPODB)
	if err != nil {
		log.Println("EUEL01", err)
		return err
	} else {
		defer db.Close()

		// FromDspName := strings.Split(fmt.Sprintf("%v", config.(map[string]interface{})["From"]), " ")[0]
		// FromRaw := fmt.Sprintf("%v", config.(map[string]interface{})["FromRaw"])
		// ReplyTo := fmt.Sprintf("%v", config.(map[string]interface{})["ReplyTo"])
		//	Subject := fmt.Sprintf("%v", config.(map[string]interface{})["Subject"])
		//Url := fmt.Sprintf("%v", config.(map[string]interface{})["Url"])

		sqlString := `Insert into  emaillog (FromId,ToId ,
			Subject ,Body , CreationDate, SentDate ,Status, EmailServer, FromDspName, Requested_Source, ReplyTo)
		VALUES (?, ?, ?, ?, NOW(),  NOW(),'SENT', ?,?, ?, ?)`
		_, err := db.Exec(sqlString, input.FromRaw, input.ToEmailId, input.Subject, input.Body, Url, input.FromDspName, ReqSource, input.ReplyTo)
		if err != nil {
			log.Println(err)
			return err
		} else {
			log.Println("Inserted Successfully")
		}

	}

	log.Println("emailLog (-)")
	return nil

}
