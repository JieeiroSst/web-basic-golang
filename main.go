package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	//"github.com/GeertJohan/go.rice"
	"html/template"
	"net/http"
	"encoding/gob"
)

type Account struct {
	nameUser string `json:"name_user"`
	password string `json:"password"`
}

type M map[string]interface{}

var (
	Accounts []Account
	store = sessions.NewCookieStore([]byte("new-authentication-key"),
		[]byte("new-encryption-key"),
		[]byte("old-authentication-key"),
		[]byte("old-encryption-key"),
		)
)

func init(){
	gob.Register(&Account{})
	gob.Register(&M{})
}

func checkPassword(pass,repass string) bool{
	if pass==repass{
		return true
	}
	return false
}
func checkEmprty(strings string) bool{
	if len(strings)<0 {
		return true
	}
	return false
}
func MyHandler(w http.ResponseWriter,r* http.Request){
	seesion,err:=store.Get(r,"session-name")
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	seesion.Values["foo"]="bar"
	seesion.Values[42]=43
	err =seesion.Save(r,w)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderHTML(w http.ResponseWriter,r*http.Request){
	link:="./views/index.html"
	account:=Account{
		nameUser: r.FormValue("nameUser"),
		password: r.FormValue("password"),
	}
	Accounts=append(Accounts,account)
	_ = sessions.Options{
		Path:     "/login",
		Domain:   "localhost:4000/login",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	}
	tpl,err:=template.ParseFiles(link)
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
	}
	_ = tpl.Execute(w, account)
}
func renderHTMLSignUp(w http.ResponseWriter,r*http.Request){
	link:="./views/sinup.html"
	account:=Account{
		nameUser: r.FormValue("nameUser"),
		password: r.FormValue("password"),
	}
	Accounts=append(Accounts,account)
	_ = sessions.Options{
		Path:     "/login",
		Domain:   "localhost:4000/sign-up",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	}
	tpl,err:=template.ParseFiles(link)
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
	}
	_ = tpl.Execute(w, account)
}
func renderPage(w http.ResponseWriter, r* http.Request){
	link:="./views/page.html"
	
	tpl,err:=template.ParseFiles(link)
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
	}
	tpl.Execute(w,nil)
}

func main(){
	r:=mux.NewRouter()
	
	//box := rice.MustFindBox("cssfiles")
	//cssFileServer := http.StripPrefix("/css/", http.FileServer(box.HTTPBox()))
	//r.Handle("/static/",cssFileServer)
	
	r.HandleFunc("/login.html",renderHTML)
	r.HandleFunc("/page.html",renderPage)
	r.HandleFunc("/sinup.html",renderHTMLSignUp)
	
	_ = http.ListenAndServe(":4000", r)
}