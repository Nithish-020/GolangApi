package main

import (
	"io"
	"log"
	"net/textproto"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

// Bond represents the information about a bond.
type Bond struct {
	SecurityCategory string
	ISIN             string
	MaturityDate     string
	AllotmentDate    string
	SettlementDate   string
	IndicativeYield  string
}

func processHTML(html string, ch chan<- Bond) {
	// Implement your logic to extract information from the HTML part
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// Extract information using goquery selectors
	doc.Find("table.x_140845973m_6914350450650205689table tbody tr").Each(func(i int, s *goquery.Selection) {
		// Implement your logic to extract information from each row of the table
		securityCategory := s.Find("td").Eq(0).Text()
		ISIN := s.Find("td").Eq(3).Text()

		// Add more logic to extract other fields as needed
		// Example: Extracting Maturity Date
		maturityDate := s.Find("td").Eq(4).Text()
		allotmentDate := s.Find("td").Eq(7).Text()
		settlementDate := s.Find("td").Eq(8).Text()
		indicativeYield := s.Find("td").Eq(12).Text()

		// Send the extracted data to the channel
		ch <- Bond{
			SecurityCategory: securityCategory,
			ISIN:             ISIN,
			MaturityDate:     maturityDate,
			AllotmentDate:    allotmentDate,
			SettlementDate:   settlementDate,
			IndicativeYield:  indicativeYield,
		}
	})
}

func main() {
	// Connect to the server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to log out
	defer c.Logout()

	// Authenticate
	if err := c.Login("karthik2768@gmail.com", "jfsgigvcjhcdrbgh"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	keywords := []string{"RBI Retail Direct - Auction Announcement Notification"}
	// Search for unread messages
	criteria := imap.NewSearchCriteria()
	criteria.Header = textproto.MIMEHeader{"From": {"karthikraja.k@fcsonline.co.in"}}

	for _, keyword := range keywords {
		criteria.Text = append(criteria.Text, keyword)
	}

	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal("Error searching for emails:", err)
	}

	if len(uids) == 0 {
		log.Println("No unread mails")
		return
	}

	// Find the most recent email
	var recentEmailUID uint32
	var recentEmailDate time.Time

	for _, uid := range uids {
		seqset := new(imap.SeqSet)
		seqset.AddNum(uid)
		items := []imap.FetchItem{imap.FetchEnvelope}

		messages := make(chan *imap.Message, 1)
		done := make(chan struct{})

		go func() {
			defer close(done)
			if err := c.Fetch(seqset, items, messages); err != nil {
				log.Println("Error fetching email:", err)
			}
		}()

		msg := <-messages
		emailDate := msg.Envelope.Date

		if emailDate.After(recentEmailDate) {
			recentEmailDate = emailDate
			recentEmailUID = uid
		}
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(recentEmailUID)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	messages := make(chan *imap.Message, 10)
	done := make(chan struct{})
	bondDataCh := make(chan Bond)

	go func() {
		defer close(done)
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Println("Error fetching email body:", err)
		}
	}()

	// Use a separate goroutine to handle the bondData slice
	go func() {
		for msg := range messages {
			r := msg.GetBody(section)
			if r == nil {
				log.Println("Server didn't return message body")
				continue
			}

			// Create a mail reader
			mr, err := mail.CreateReader(r)
			if err != nil {
				log.Println("Error creating mail reader:", err)
				continue
			}

			// Print mail attributes and body
			header := mr.Header
			if date, err := header.Date(); err == nil {
				log.Println("Date:", date)
			}
			if from, err := header.AddressList("From"); err == nil {
				log.Println("From:", from)
			}
			if subject, err := header.Subject(); err == nil {
				log.Println("Subject:", subject)
			}

			for {
				_, err := mr.NextPart()
				if err == io.EOF {
					log.Println("Error reading next part:", err)
					break
				} else if err != nil {
					log.Fatal(err)
				}

				htmlPart, err := mr.NextPart()
				if err != nil {
					log.Println("Error reading HTML part:", err)
					log.Fatal(err)
				}
				htmlBody, err := io.ReadAll(htmlPart.Body)
				if err != nil {
					log.Println("Error reading HTML body:", err)
					log.Fatal(err)
				}

				// processHTML function is used to extract required content from the body in html
				processHTML(string(htmlBody), bondDataCh)

				// switch h := p.Header.(type) {
				// case *mail.InlineHeader:
				// 	// This is the message's text (either plain text or HTML)
				// 	b, _ := io.ReadAll(p.Body)
				// 	log.Println("Got text:", string(b))
				// case *mail.AttachmentHeader:
				// 	// This is an attachment
				// 	filename, _ := h.Filename()
				// 	log.Println("Got attachment:", filename)
				// }
			}
		}
	}()

	// Wait for all messages to be processed
	go func() {
		for bond := range bondDataCh {
			// Now 'bondData' contains the structured information about each bond
			// You can send this data to a database or perform other operations as needed
			log.Printf("Security Category: %s, ISIN: %s\n", bond.SecurityCategory, bond.ISIN)
			log.Printf("Maturity Date: %s\n", bond.MaturityDate)
			log.Printf("Allotment Date: %s\n", bond.AllotmentDate)
			log.Printf("Settlement Date: %s\n", bond.SettlementDate)
			log.Printf("Indicative Yield: %s\n", bond.IndicativeYield)
			log.Println("--------------------------")
			log.Println("bond", bond)
		}
	}()

	// Wait for all messages to be processed
	<-done
}
