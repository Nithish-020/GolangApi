package exchangecall

import (
	"fcs23pkg/common"
	"fcs23pkg/ftdb"
	"fcs23pkg/integration/bse/bseipo"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// func ChangeBsePassword(pUser string, pToken string, pBrokerId int) error {
// 	log.Println("ChangeBsePassword (+)")
// 	var lPasswordReq bseipo.PswdReqStruct
// 	// var lAutopswdchanger bseipo.PswdRespStruct
// 	paswordDate, Password, lErr1 := checkPasswordExpiry(pBrokerId)
// 	if lErr1 != nil {
// 		log.Println("BCBP01", lErr1)
// 		return lErr1
// 	} else {
// 		if paswordDate > 13 {
// 			lDecodepswd, lErr2 := common.DecodeToString(Password)
// 			if lErr2 != nil {
// 				log.Println("BCBP02", lErr2)
// 				return lErr2
// 			} else {
// 				lPasswordReq.OldPswd = lDecodepswd
// 				lCompanycode := "FTC"
// 				code := generateRandomCode(lCompanycode)
// 				lPasswordReq.NewPswd = code
// 				lPasswordReq.ConfirmPswd = code

// 				lAutopswdchanger, lErr3 := bseipo.BsePassword(lPasswordReq, pUser, pToken, pBrokerId)
// 				if lErr3 != nil {
// 					return lErr3
// 				} else {
// 					if lAutopswdchanger.ErrorCode != "E" {
// 						lEncryptPswd := common.EncodeToString(code)
// 						lErr5 := updateBsePswd(lEncryptPswd)
// 						if lErr5 != nil {
// 							log.Println("BCBP05", lErr5)
// 							return lErr5
// 						}
// 					}
// 				}

// 			}
// 		}

// 		log.Println("ChangeBsePassword (-)")
// 	}
// 	return nil
// }

// func checkPasswordExpiry(pBrokerId int) (int, string, error) {
// 	log.Println("checkPasswordExpiry (+)")
// 	var lcheckPassword int
// 	var loldPassword string

// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("BCPE01", lErr1)
// 		return lcheckPassword, loldPassword, lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `select datediff (now(),d.UpdatedDate)as DateDifference,Password
// 						from a_ipo_directory d
// 						where d.Stream  = 'BSE'
// 						and d.brokerMasterId = ?`

// 		lRows1, lErr2 := lDb.Query(lCoreString, pBrokerId)
// 		if lErr2 != nil {
// 			log.Println("BCPE02", lErr2)
// 			return lcheckPassword, loldPassword, lErr2
// 		} else {
// 			//This for loop is used to collect the records from the database and store them in structure
// 			for lRows1.Next() {
// 				lErr3 := lRows1.Scan(&lcheckPassword, &loldPassword)
// 				if lErr3 != nil {
// 					log.Println("BCPE03", lErr3)
// 					return lcheckPassword, loldPassword, lErr3
// 				}
// 			}
// 		}
// 	}
// 	log.Println("checkPasswordExpiry (-)")
// 	return lcheckPassword, loldPassword, nil
// }

// func generateRandomCode(prefix string) string {
// 	log.Println("generateRandomCode (+)")
// 	var Randomcode string
// 	// Seed the random number generator with the current time
// 	rand.Seed(time.Now().UnixNano())
// 	// Define the set of characters to choose from
// 	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@$_&."
// 	// Generate a random 5-character code
// 	code := make([]byte, 5)
// 	for i := range code {
// 		code[i] = characters[rand.Intn(len(characters))]
// 	}
// 	Randomcode = fmt.Sprintf("%s%s", prefix, string(code))
// 	//  This if condition is used to check the value is unique or not
// 	log.Println("generateRandomCode (-)")
// 	return Randomcode
// }
// func updateBsePswd(pNewPswd string) error {
// 	log.Println("UpdatedDateBsePswd (+)")
// 	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
// 	if lErr1 != nil {
// 		log.Println("BUDBP01", lErr1)
// 		return lErr1
// 	} else {
// 		defer lDb.Close()

// 		lCoreString := `update a_ipo_directory
// 		set Password = ?,UpdatedDate = now()
// 		where Stream  = 'BSE' `

// 		_, lErr2 := lDb.Exec(lCoreString, pNewPswd)
// 		if lErr2 != nil {
// 			log.Println("BUDBP01", lErr2)
// 			return lErr2
// 		}

// 	}
// 	log.Println("UpdatedDateBsePswd (-)")
// 	return nil
// }

// this method changing the password in exchange
func ChangeBsePassword(pUser string, pToken string, pBrokerId int) (string, error) {
	log.Println("ChangeBsePassword (+)")
	var lPasswordReq bseipo.PswdReqStruct
	var lErrCode string
	// var lAutopswdchanger bseipo.PswdRespStruct
	// paswordDate, Password, lErr1 := checkPasswordExpiry(pBrokerId)
	Password, lErr1 := checkPasswordExpiry(pBrokerId)
	if lErr1 != nil {
		log.Println("BCBP01", lErr1)
		return lErrCode, lErr1
	} else {
		// if paswordDate >= 13 {
		if Password != "" {
			lDecodepswd, lErr2 := common.DecodeToString(Password)
			if lErr2 != nil {
				log.Println("BCBP02", lErr2)
				return lErrCode, lErr2
			} else {
				lPasswordReq.OldPswd = lDecodepswd
				lCompanycode := "FtC"
				code := generateRandomCode(lCompanycode)
				lPasswordReq.NewPswd = code
				lPasswordReq.ConfirmPswd = code

				lAutopswdchanger, lErr3 := bseipo.BsePassword(lPasswordReq, pUser, pToken, pBrokerId)
				if lErr3 != nil {
					log.Println("BCBP03", lErr3)
					return lErrCode, lErr3
				} else {
					log.Println("Change pswd resp", lAutopswdchanger)
					if lAutopswdchanger.ErrorCode == "0" {
						// lErrCode = "0"
						lErrCode = lAutopswdchanger.ErrorCode
						lEncryptPswd := common.EncodeToString(code)
						lErr4 := UpdateBsePswd(lEncryptPswd, pBrokerId)
						if lErr4 != nil {
							log.Println("BCBP04", lErr4)
							return lErrCode, lErr4
						}
					}
				}
			}
		}
	}
	// }
	log.Println("ChangeBsePassword (-)")
	return lErrCode, nil
}

// func checkPasswordExpiry(pBrokerId int) (int, string, error) {
func checkPasswordExpiry(pBrokerId int) (string, error) {
	log.Println("checkPasswordExpiry (+)")
	// var lcheckPassword int
	var loldPassword string

	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BCPE01", lErr1)
		// return lcheckPassword, loldPassword, lErr1
		return loldPassword, lErr1

	} else {
		defer lDb.Close()

		// this query checking the update date difference for password ,failed case so it will be commented
		// lCoreString := `select datediff (now(),d.UpdatedDate)as DateDifference,Password
		// 				from a_ipo_directory d
		// 				where d.Stream  = 'BSE'
		// 				and d.Status = 'Y'
		// 				and d.brokerMasterId = ?`

		lCoreString := `select d.Password
		from a_ipo_directory d
		where d.Stream  = 'BSE'
		and d.Status = 'Y'
		and d.brokerMasterId = ?`

		lRows1, lErr2 := lDb.Query(lCoreString, pBrokerId)
		if lErr2 != nil {
			log.Println("BCPE02", lErr2)
			// return lcheckPassword, loldPassword, lErr2
			return loldPassword, lErr2
		} else {
			//This for loop is used to collect the records from the database and store them in structure
			for lRows1.Next() {
				// lErr3 := lRows1.Scan(&lcheckPassword, &loldPassword)
				lErr3 := lRows1.Scan(&loldPassword)
				if lErr3 != nil {
					log.Println("BCPE03", lErr3)
					// return lcheckPassword, loldPassword, lErr3
					return loldPassword, lErr3
				}
			}
		}
	}
	log.Println("checkPasswordExpiry (-)")
	// return lcheckPassword, loldPassword, nil
	return loldPassword, nil
}

// generating password for exchange
func generateRandomCode(prefix string) string {
	log.Println("generateRandomCode (+)")
	var Randomcode string

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Define the character sets
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharacters := "@"
	numbers := "0123456789"

	// Generate a random 4-character code (without a special character)
	code := make([]byte, 4)
	for i := range code {
		code[i] = characters[rand.Intn(len(characters))]
	}

	// Add one special character to the code
	specialChar := specialCharacters[rand.Intn(len(specialCharacters))]
	code = append(code, specialChar)

	numbersChar := numbers[rand.Intn(len(numbers))]
	code = append(code, numbersChar)

	// Shuffle the code to randomize the position of the special character
	rand.Shuffle(len(code), func(i, j int) {
		code[i], code[j] = code[j], code[i]
	})

	Randomcode = fmt.Sprintf("%s%s", prefix, string(code))

	//  This if condition is used to check if the value is unique or not

	log.Println("generateRandomCode (-)")
	return Randomcode
}

//updating the changed password in Db
func UpdateBsePswd(pNewPswd string, pBrokerId int) error {
	log.Println("UpdatedDateBsePswd (+)")
	lDb, lErr1 := ftdb.LocalDbConnect(ftdb.IPODB)
	if lErr1 != nil {
		log.Println("BUDBP01", lErr1)
		return lErr1
	} else {
		defer lDb.Close()

		lCoreString := `update a_ipo_directory
		set Password = ?,UpdatedDate = now()
		where Stream  = 'BSE' and Status = 'Y'
		and brokerMasterId = ?`

		_, lErr2 := lDb.Exec(lCoreString, pNewPswd, pBrokerId)
		if lErr2 != nil {
			log.Println("BUDBP02", lErr2)
			return lErr2
		}
	}
	log.Println("UpdatedDateBsePswd (-)")
	return nil
}
