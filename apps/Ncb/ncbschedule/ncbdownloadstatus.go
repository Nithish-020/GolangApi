package ncbschedule

import (
	"fcs23pkg/ftdb"
	"log"
)

type NcbSchRespStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}
type NcbBrokers struct {
	BrokerId int
	Exchange string
}

func NcbBrokerList() ([]NcbBrokers, error) {
	log.Println("NcbBrokerList (+)")

	var Brokers NcbBrokers
	var BrokersList []NcbBrokers

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("NBLO1", lErr1)
		return BrokersList, lErr1
	} else {
		defer lDb.Close()
		lCoreString := `select Bm.Id ,d.Stream 
		                from a_ipo_brokermaster bm ,a_ipo_directory d,a_ipo_memberdetails m 
		                where bm.Id = d.brokerMasterId 
		                and bm.Id = m.BrokerId 
		                and m.AllowedModules like '%Gsec%'
		                and bm.Status = 'Y' 
		                and d.Status ='Y'
		                and m.Flag = 'Y'`

		lRows, lErr2 := lDb.Query(lCoreString)
		if lErr2 != nil {
			log.Println("NBLO2", lErr2)
			return BrokersList, lErr2
		} else {
			for lRows.Next() {
				lErr3 := lRows.Scan(&Brokers.BrokerId, &Brokers.Exchange)

				if lErr3 != nil {
					log.Println("NBLO3", lErr3)
					return BrokersList, lErr3
				} else {
					BrokersList = append(BrokersList, Brokers)
				}
			}
		}

	}

	log.Println("NcbBrokerList(-)")
	return BrokersList, nil
}
