package apigate

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"log"
	"net/http"
	"strings"
)

type RequestorDetails struct {
	RealIP      string
	ForwardedIP string
	Method      string
	Path        string
	Host        string
	RemoteAddr  string
	Header      string
	Body        string
	EndPoint    string
}

// --------------------------------------------------------------------
// get request header details
// --------------------------------------------------------------------
func GetHeaderDetails(r *http.Request) string {
	log.Println("GetHeaderDetails+")
	value1 := ""
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			value1 = value1 + " " + name + "-" + value
		}
	}

	// To Get the Active Domain from the Database
	origin := GetAllowOrigin()
	log.Println("Origin: ", origin)
	common.ABHIAllowOrigin = origin

	log.Println("GetHeaderDetails-")
	return value1
}

// --------------------------------------------------------------------
// function reads the API requestor details and send return them
// as structure to the caller
// --------------------------------------------------------------------
func GetRequestorDetail(r *http.Request) RequestorDetails {
	log.Println("GetRequestorDetail+")

	var reqDtl RequestorDetails
	reqDtl.RealIP = r.Header.Get("Referer")
	reqDtl.ForwardedIP = r.Header.Get("X-Forwarded-For")
	reqDtl.Method = r.Method
	reqDtl.Path = r.URL.Path + "?" + r.URL.RawQuery
	reqDtl.Host = r.Host
	reqDtl.RemoteAddr = r.RemoteAddr
	if strings.Contains(r.URL.Path, "/order/placeorder/") {
		reqDtl.EndPoint = r.URL.Path[:len("/order/placeorder/")]
	} else if strings.Contains(r.URL.Path, "/deals/count/") {
		reqDtl.EndPoint = r.URL.Path[:len("/deals/count/")]
	} else {
		reqDtl.EndPoint = r.URL.Path
	}

	reqDtl.Header = GetHeaderDetails(r)
	//body, _ := ioutil.ReadAll(r.Body)
	//reqDtl.Body = string(body)
	log.Println("GetRequestorDetail-")

	return reqDtl
}

// --------------------------------------------------------------------
// Copy for Sgbapicopy brach
// --------------------------------------------------------------------

func GetAllowOrigin() []string {
	log.Println("GetAllowOrigin (+)")

	var lAllowOrigin []string
	var lDomain string
	lDb, lErr := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr != nil {
		log.Println("BGAO01", lErr)
		return lAllowOrigin
	} else {
		defer lDb.Close()

		// ! To get the active DomainNames from the DB
		lCoreString2 := `select bm.DomainName
			from a_ipo_brokerMaster bm
			where bm.Status = 'Y'`

		lRows1, lErr := lDb.Query(lCoreString2)
		if lErr != nil {
			log.Println("BGAO02", lErr)
			return lAllowOrigin
		} else {
			//This for loop is used to collect the records from the database and assign them to lDomain
			for lRows1.Next() {
				lErr := lRows1.Scan(&lDomain)
				if lErr != nil {
					log.Println("BGAO03", lErr)
					return lAllowOrigin
				} else {
					// Append the domains into the lArrowedOrigin array[]
					lAllowOrigin = append(lAllowOrigin, lDomain)
				}
			}
		}
	}
	log.Println("GetAllowOrigin (-)")
	return lAllowOrigin
}
