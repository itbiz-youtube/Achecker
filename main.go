package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Email struct {
	Email string
}
type EmailValidationResult struct {
	Email        string
	ErrorCodeNum string
	ErrorText    string
}

const EmailMaxLength int = 254

func emailSet(w http.ResponseWriter, r *http.Request) {
	var p Email

	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *malformedRequest
		fmt.Print(mr.msg)

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

	//Вызов проверки длины email
	resultMaxLengthCheck := emailLengthValidation(clearemail)
	if resultMaxLengthCheck != "0" {
		js, err := json.Marshal(resultMaxLengthCheck)
		var val []byte = []byte(js)
		s, _ := strconv.Unquote(string(val))
		if err == nil {
			log.Println("Success send response" + resultMaxLengthCheck)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(s))
		} else {
			log.Println("Error send response" + resultMaxLengthCheck)

		}

	}
}

//TODO:Сделать пакет из этой функции
//TODO:Привести логи к стандарту
//TODO:Добавить обработчики условий
func emailLengthValidation(email string) string {

	lengthofstring := len(email)
	errorcode := "1000"
	errortext := "bigger then max(254) email length"
	jsstring := *new(string)
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
			return (jsstring)
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
