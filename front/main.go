package front

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Samollo/maain/request"
)

type PageVariables struct {
	PageTitle string
	V         []Variable
}

var RequestServ request.Request
var title = "Distrib"

func LaunchFront() {
	http.HandleFunc("/", DisplayHomePage)
	http.HandleFunc("/results", UserSelected)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func DisplayHomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("front/select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

type Variable struct {
	Value string
}

func UserSelected(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	requestSentence := r.FormValue("valueEntered")

	results := RequestServ.Intersection(requestSentence)
	element := make([]Variable, 0)

	for _, value := range results {
		element = append(element, Variable{strings.Title(value)})
	}

	fmt.Printf("%v", results)

	pv := PageVariables{
		PageTitle: requestSentence,
		V:         element,
	}

	t, err := template.ParseFiles("front/search.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, pv)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
