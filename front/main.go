package front

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	rq "github.com/Samollo/maain/request"
)

type PageVariables struct {
	PageTitle string
	V         []Variable
}

var RequestServ = rq.InitializeRequest("wordpages")

func LaunchFront() {
	http.HandleFunc("/", DisplayHomePage)
	http.HandleFunc("/result", UserSelected)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func DisplayHomePage(w http.ResponseWriter, r *http.Request) {
	title := "Home"

	pv := PageVariables{
		PageTitle: title,
	}

	t, err := template.ParseFiles("front/select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, pv)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

type Variable struct {
	Value string
}

func UserSelected(w http.ResponseWriter, r *http.Request) {
	title := "Search"
	r.ParseForm()

	mail := r.FormValue("valueEntered")

	k := RequestServ.ReturnFoundPages(mail)
	element := make([]Variable, 0)

	for _, value := range k {
		element = append(element, Variable{value})
	}

	fmt.Printf("%v", k)

	pv := PageVariables{
		PageTitle: title,
		V:         element,
	}

	t, err := template.ParseFiles("front/select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, pv)
	if err != nil {
		log.Print("template executing error: ", err)
	}
	/*
		fmt.Println(mail)

		for _, p := range k {
			fmt.Fprintf(w, "<a href=\"https://fr.wikipedia.org/wiki/"+p+"\">"+p+"</a><br/>")
		} */
}
