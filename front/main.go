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
	Variable  string
}

var RequestServ = rq.InitializeRequest("wordpages")

func LaunchFront() {
	http.HandleFunc("/", DisplayHomePage)
	http.HandleFunc("/result", UserSelected)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func DisplayHomePage(w http.ResponseWriter, r *http.Request) {
	title := "Gogole"

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

func UserSelected(w http.ResponseWriter, r *http.Request) {
	title := "Gogole"
	r.ParseForm()

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
	k := RequestServ.ReturnFoundPages(r.Form.Get("valueEntered"))

	for _, p := range k {
		fmt.Fprintf(w, "<a href=\"https://fr.wikipedia.org/wiki/"+p+"\">"+p+"</a><br/>")
	}
}
