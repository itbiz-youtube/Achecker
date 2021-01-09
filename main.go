package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"net"

	exists "github.com/ashkan90/golang-in_array"
	"golang.org/x/net/idna"
)

type Email struct {
	Email string
}
type EmailValidationResult struct {
	Email        string
	ErrorCodeNum string
	ErrorText    string
}

// type DNSRecord struct {
// 	DNSMX string
// 	DNSA  string
// }

const EmailMaxLength int = 254

func emailSet(w http.ResponseWriter, r *http.Request) {
	var p Email

	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		// fmt.Print(mr.msg)

		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)

		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
		return
	}
	// fmt.Printf("%v", p)
	clearemail := strings.ReplaceAll(p.Email, "{", "")

	// Вызов проверки длины email
	resultMaxLengthCheck := emailLengthValidation(clearemail)
	if resultMaxLengthCheck == "error" {

		log.Println("Error send response" + resultMaxLengthCheck)
	} else if resultMaxLengthCheck == "0" {
		log.Println("OK;" + clearemail + ";go to next check")

	} else {
		js, err := json.Marshal(resultMaxLengthCheck)
		var val []byte = []byte(js)
		s, _ := strconv.Unquote(string(val))
		if err == nil {
			log.Println("Success send response" + resultMaxLengthCheck)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(s))

		}

	}

	resultCheckLocalOrReservedDNS := emailDNSValidation(clearemail)
	// fmt.Println(resultCheckLocalOrReservedDNS)
	if resultCheckLocalOrReservedDNS == "error" {

		log.Println("Error send response" + resultCheckLocalOrReservedDNS)
	} else if resultCheckLocalOrReservedDNS == "0" {
		log.Println("OK;" + clearemail + ";go to next check")

	} else {
		js, err := json.Marshal(resultCheckLocalOrReservedDNS)
		var val []byte = []byte(js)
		s, _ := strconv.Unquote(string(val))
		if err == nil {
			log.Println("Success send response" + resultCheckLocalOrReservedDNS)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(s))

		}

	}
	resultCheckExistingDNSRecords := emailDNSCheck(clearemail)
	// fmt.Println(resultCheckLocalOrReservedDNS)
	if resultCheckExistingDNSRecords == "error" {

		log.Println("Error send response" + resultCheckExistingDNSRecords)
	} else if resultCheckExistingDNSRecords == "0" {
		log.Println("OK;" + clearemail + ";go to next check")

	} else {
		js, err := json.Marshal(resultCheckExistingDNSRecords)
		var val []byte = []byte(js)
		s, _ := strconv.Unquote(string(val))
		if err == nil {
			log.Println("Success send response" + resultCheckExistingDNSRecords)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(s))

		}

	}

}

//TODO:Сделать пакет из этой функции
//TODO:Привести логи к стандарту
//TODO:Добавить обработчики условий
//TODO:Переписать логику формирования ответа и присвоение переменных внутри if/else
func emailLengthValidation(email string) string {

	lengthofstring := len(email)
	errorcode := "1000"
	errortext := "bigger then max(254) email length"
	// jsstring := *new(string)
	if lengthofstring > EmailMaxLength { //check length of email
		log.Println(email + ";" + errortext + ";errorcodenum=" + errorcode)
		emailvarresult := EmailValidationResult{
			email, errorcode, errortext,
		}
		js, err := json.Marshal(emailvarresult)
		if err == nil {
			jsstring := string(js)
			// fmt.Printf(jsstring)
			return (jsstring)
		} else {
			log.Println("Error make json when check email" + email)
			// fmt.Printf(jsstring)
			return "error"
		}
	} else {
		return "0"
	}

}

func emailDNSValidation(email string) string {
	// Arguable pattern to extract the domain. Not aiming to validate the domain nor the email
	lastATpos := strings.Index(email, "@")
	// fmt.Println(lastATpos)
	lengthofstring := len(email)
	// fmt.Println(lengthofstring)
	host := string(email[lastATpos+1 : lengthofstring])
	//TODO: Сплитить и посчитать на сколько элементов было разделение
	hostparts := strings.Split(host, ".")

	//TODO: Нужно раскоментить и обратботать
	reservedTopLevelDNSNames := []string{
		"test",
		"example",
		"invalid",
		"localhost",
		// mDNS
		"local",
		// Private DNS Namespaces
		"intranet",
		"internal",
		"private",
		"corp",
		"home",
		"lan",
	}

	//Проверка того что домен локальный
	islocalDomain := len(hostparts) <= 1

	// fmt.Println(hostparts[len(hostparts)-1])
	// fmt.Println(islocalDomain)
	//Проверка того что домен находится не в списке зарезервированных
	isReservedTopLevelDNSName := exists.In_array(hostparts[len(hostparts)-1], reservedTopLevelDNSNames, true)
	// fmt.Println(isReservedTopLevelDNSName)

	if islocalDomain || isReservedTopLevelDNSName {
		errorcode := "2000"
		errortext := "Local, mDNS or reserved domain (RFC2606, RFC6762)"

		log.Println(email + ";" + errortext + ";errorcodenum=" + errorcode)
		emailvarresult := EmailValidationResult{
			email, errorcode, errortext,
		}
		js, err := json.Marshal(emailvarresult)
		if err == nil {
			jsstring := string(js)
			// fmt.Printf(jsstring)
			return (jsstring)
		} else {
			log.Println("Error make json when check email" + email)
			// fmt.Printf(jsstring)
			return "error"
		}
	} else {
		return "0"
	}
}

func emailDNSCheck(email string) string {

	// Arguable pattern to extract the domain. Not aiming to validate the domain nor the email
	lastATpos := strings.Index(email, "@")

	lengthofstring := len(email)
	// fmt.Println(lengthofstring)
	host := (string(email[lastATpos+1 : lengthofstring]))
	//TODO: Сплитить и посчитать на сколько элементов было разделение

	//translate to ASCII for non latin domain names
	var p *idna.Profile
	p = idna.New()
	hostASCII, _ := p.ToASCII(host)

	hostnametrimmed := strings.TrimRight(hostASCII, ".")

	//check DNS A record
	iprecords, _ := net.LookupIP(hostnametrimmed)
	lenIPRecords := len(iprecords) //shuld be 0 if domain not exist

	// //check Name server
	// nameServer, _ := net.LookupNS(hostnametrimmed)
	// lenNameServer := len(nameServer) //shuld be 0 if domain not exist
	// check MX records
	mxrecords, _ := net.LookupMX(hostnametrimmed)
	lenMXServer := len(mxrecords) //shuld be 0 if domain not exist

	if lenIPRecords == 0 {
		errorcode := "2001"
		errortext := "No DNS_A reccord for this email domain"

		log.Println(email + ";" + errortext + ";errorcodenum=" + errorcode)
		emailvarresult := EmailValidationResult{
			email, errorcode, errortext,
		}
		js, err := json.Marshal(emailvarresult)
		if err == nil {
			jsstring := string(js)
			// fmt.Printf(jsstring)
			return (jsstring)
		} else {
			log.Println("Error make json when check email" + email)
			// fmt.Printf(jsstring)
			return "error"
		}
	} else if lenMXServer == 0 {
		errorcode := "2002"
		errortext := "No DNS_MX reccord for this email domain"

		log.Println(email + ";" + errortext + ";errorcodenum=" + errorcode)
		emailvarresult := EmailValidationResult{
			email, errorcode, errortext,
		}
		js, err := json.Marshal(emailvarresult)
		if err == nil {
			jsstring := string(js)
			// fmt.Printf(jsstring)
			return (jsstring)
		} else {
			log.Println("Error make json when check email" + email)
			// fmt.Printf(jsstring)
			return "error"
		}
	} else {
		return "0"
	}

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/email", emailSet)

	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
