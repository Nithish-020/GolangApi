package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
)

func GenerateOTP() string {
	log.Println("GenerateOTP+")
	rand.Seed(time.Now().UnixNano())
	b := rand.Intn(1000000)
	otp := fmt.Sprintf("%06d", b)
	//log.Println(token)
	log.Println("GenerateOTP-")
	return otp
}

// ---------------------------------------------------------------
// function Encrypts the given mailid
// Returns the Encryted Mailid
// ------------------------------------------------------------
func GetEncryptedemail(emailId string) (string, error) {
	log.Println("GetEncryptedemail(+)")
	idxStart := 0
	idxMiddle := strings.Index(emailId, "@")
	// log.Println(" idxStart:", idxStart)
	// log.Println(" idxEnd:", idxMiddle)
	// for printinf first letter
	firstLetter := emailId[idxStart : idxStart+1]
	//log.Println(firstLetter)

	firstHalf := emailId[idxStart+1 : idxMiddle]

	for i := range firstHalf {

		firstHalf = strings.Replace(firstHalf, string(firstHalf[i]), "*", 1)
	}
	//idxMiddle := strings.Index(emailId, "@")
	idxEnd := len(emailId)
	// log.Println(" idxStart1:", idxMiddle)
	// log.Println(" idxEnd1:", idxEnd)

	letterAfterAt := emailId[idxMiddle+2 : idxEnd]
	//log.Println(letterAfterAt)
	// for printing first letter
	SecondHalf := emailId[idxMiddle : idxMiddle+2]

	for j := range letterAfterAt {
		//log.Println(j)
		if string(letterAfterAt[j]) == "." {
			break
		}
		letterAfterAt = strings.Replace(letterAfterAt, string(letterAfterAt[j]), "*", 1)

	}

	encrytedEmaiId := firstLetter + firstHalf + SecondHalf + letterAfterAt
	//log.Println(encrytedEmaiId)
	log.Println("GetEncryptedemail(-)")
	return encrytedEmaiId, nil
}

// ---------------------------------------------------------------
// function Encrypts the given MobileNumber
// Returns the Encryted MobileNumber
// ------------------------------------------------------------
func GetEncryptedMobile(mobileNumber string) (string, error) {
	log.Println("GetEncryptedMobile(+)")
	for K := range mobileNumber {
		if K == len(mobileNumber)-4 {
			break
		}
		mobileNumber = strings.Replace(mobileNumber, string(mobileNumber[K]), "*", 1)
	}
	//log.Println(mobileNumber)
	log.Println("GetEncryptedMobile(-)")
	return mobileNumber, nil

}

func LogError(who string, ref string, msg string) {
	log.Println(who, "ERROR:("+ref+")", msg)
}

func LogDebug(who string, ref string, msg string) {
	log.Println("DEBUG: ", who, ref, msg)
}

func GetFileName_UUID_String() string {
	var id = uuid.New()
	return id.String()
}

// -------------------------------------
// Method encodes the given input
// Returns the data in string
// --------------------------------------
func EncodeToString(fileBody string) string {
	log.Println("EncodeToString(+)")

	encodedText := base64.StdEncoding.EncodeToString([]byte(fileBody))
	//log.Println("Encoded text:", encodedText)
	log.Println("EncodeToString(-)")

	return encodedText
}

// -------------------------------------
// Method decode the given input
// Returns the data in string
// --------------------------------------
func DecodeToString(fileBody string) (string, error) {
	log.Println("DecodeToString(+)")
	decodeText, err := base64.StdEncoding.DecodeString(fileBody)
	if err != nil {
		log.Println("Decoded text:", decodeText)
		LogDebug("common.DecodeToString ", "(CDS01)", err.Error())
		return string(decodeText), err
	}

	log.Println("DecodeToString(-)")
	return string(decodeText), nil
}

// ------------------------------------------
// Method read html file and return the
// file data as String
// --------------------------------------------
func HtmlFileToString(fileName string) (string, error) {
	log.Println("HtmlFileToString+")

	var htmlString string
	var tpl bytes.Buffer
	temp, err := template.ParseFiles(fileName) // change this
	if err != nil {
		LogDebug("common.HtmlFileToString ", "(CHFS01)", err.Error())
		return htmlString, err
	} else {
		temp.Execute(&tpl, "")
		htmlString = tpl.String()
	}

	log.Println("HtmlFileToString-")
	return htmlString, nil

}

func GetLoggedBy(ClientId string) string {
	log.Println("GetLoggedBy+")

	array := strings.Split(ClientId, ",")

	log.Println("GetLoggedBy-")
	return array[1]
}

func GetSetClient(ClientId string) string {
	log.Println("GetSetClient (+)")

	array := strings.Split(ClientId, ",")

	log.Println("GetSetClient (-)")
	return array[0]
}

// --------------------------------------------------------------------
// function reads the constants from the config.toml file
// --------------------------------------------------------------------
func ReadTomlConfig(filename string) interface{} {
	var f interface{}
	if _, err := toml.DecodeFile(filename, &f); err != nil {
		log.Println(err)
	}
	return f
}

// --------------------------------------------------------------------
// function convert the string to sha256 format
// --------------------------------------------------------------------
func EncodeToSHA256(input string) string {
	log.Println("EncodeToSHA256 (+)")

	hash := sha256.New()
	hash.Write([]byte(input))
	hashSum := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashSum)

	log.Println("EncodeToSHA256 (-)")
	return hashHex
}

// --------------------------------------------------------------------
// function convert the time and date format to customized format
// --------------------------------------------------------------------
func ChangeTimeFormat(pCustomizeLayout string, pInput string) (string, error) {
	log.Println("ChangeTimeFormat (+)")
	var lFormattedValue string

	Layout := ""
	length := len(pInput)
	if length == 19 {
		Layout = "02-01-2006 15:04:05"
	} else if length == 5 {
		Layout = "15:04"
	} else if length == 8 {
		Layout = "15:04:05"
	} else {
		Layout = "02-01-2006 15:04"
	}
	lTimevalue, lErr1 := time.Parse(Layout, pInput)
	if lErr1 != nil {
		log.Println("Error in Parse Timing:", lErr1)
		return lFormattedValue, lErr1
	} else {
		lFormattedValue = lTimevalue.Format(pCustomizeLayout)
	}

	log.Println("ChangeTimeFormat (-)")
	return lFormattedValue, nil
}

// --------------------------------------------------------------------
// function convert string value to float
// --------------------------------------------------------------------
func ConvertStringToFloat(pInput string) (float32, error) {
	log.Println("ConvertStringToFloat (+)")

	var lFloatValue float32

	if pInput != "" {
		float32Value, lErr1 := strconv.ParseFloat(pInput, 64)
		if lErr1 != nil {
			log.Println("Error in Parse Timing:", lErr1)
			return lFloatValue, lErr1
		} else {
			lFloatValue = float32(float32Value)
		}
	}

	log.Println("ConvertStringToFloat (-)")
	return lFloatValue, nil
}

// --------------------------------------------------------------------
// function to get unique random number
// --------------------------------------------------------------------
func GetRandomNumber() int {
	log.Println("GetRandomNumber (+)")

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random 6-digit number without leading zeros
	min := 100000 // Smallest 6-digit number
	max := 999999 // Largest 6-digit number
	randomNum := rand.Intn(max-min+1) + min

	log.Println("GetRandomNumber (-)")
	return randomNum
}

func RemoveDuplicateStrings(arr []string) []string {
	uniqueMap := make(map[string]bool)
	result := []string{}

	for _, item := range arr {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			result = append(result, item)
		}
	}

	return result
}

//----------------------------------------------------------------
// Function to CapitalizeText capitalizes the first letter of each word in a string.
//----------------------------------------------------------------
func CapitalizeText(input string) string {
	words := strings.Fields(input) // Split the input into words
	var capitalizedWords []string

	for _, word := range words {
		// Capitalize the first letter of the word
		capitalizedWord := strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		capitalizedWords = append(capitalizedWords, capitalizedWord)
	}

	// Join the capitalized words back into a string
	return strings.Join(capitalizedWords, " ")
}

//----------------------------------------------------------------
// Creating custom error
// ----------------------------------------------------------------

func CustomError(pErrorMsg string) error {
	err := errors.New(pErrorMsg)
	return err
}

//----------------------------------------------------------------
// this method generating a 16 digit unique string value,
// concatinates the clientId and unix timestamp value
// ----------------------------------------------------------------
func Generate16DigitString(pClientId string) string {
	log.Println("Generate16DigitString (+)")

	// Get Unix timestamp in seconds
	unixTimestamp := time.Now().Unix()

	// Calculate the remaining length needed for timestamp digits
	remainingLength := 16 - len(pClientId)

	// Convert timestamp to a string
	timestampString := fmt.Sprintf("%d", unixTimestamp)

	// Calculate how many characters from the timestamp are needed based on the prefix length
	charsFromTimestamp := remainingLength
	if len(pClientId) < remainingLength {
		charsFromTimestamp = remainingLength - len(pClientId)
	}

	// Get the appropriate substring of the timestamp
	trimmedTimestamp := timestampString[len(timestampString)-charsFromTimestamp:]

	log.Println("Generate16DigitString (-)")
	// Concatenate the prefix and adjusted timestamp
	return pClientId + trimmedTimestamp
}
