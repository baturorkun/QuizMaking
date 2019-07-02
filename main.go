package main

import (
	"QuizMaking/models"
	"QuizMaking/setting"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type PageQuizData struct {
	Quizes []models.Quiz
}

type PageQuestionsData struct {
	Count int
	Questions []models.Question
}

type PageScoreData struct {
	Quiz string
	Stats map[string]int
	ResList map[int]string
}


func cleanup() {
	fmt.Println("Cleanup")
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	setting.Setup()
	setting.Init()

	models.Setup()


	// mux serves
	mx := mux.NewRouter()
	mx.HandleFunc("/", Enter)
	mx.HandleFunc("/login", Login)
	mx.HandleFunc("/register", Register)
	mx.HandleFunc("/select", SelectQuiz)
	mx.HandleFunc("/start", StartQuiz)

	// http serves
	http.Handle("/", mx)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe( setting.ServerSetting.Port, nil)
}


func Enter(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		return
	}

	log.Println("Request ->", r.URL.Path)

	tmpl := template.Must(template.ParseFiles("templates/layouts/index.html", "templates/index.html",))


	tmpl.Execute(w, nil)

}


func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		return
	}
	log.Println("Request ->", r.URL.Path)

	if r.Method == http.MethodPost {

		r.ParseForm()
		var user models.User

		if r.Form.Get("email") == "" || r.Form.Get("password") == "" {
			http.Redirect(w, r, "/login?error=missing-credential", http.StatusPermanentRedirect)
		}

		user.Get(r.Form.Get("email"))

		if user.Password == r.Form.Get("password") {

			log.Println(">>> VALID-credential")
			session, _ := setting.Store.Get(r, "USER")

			session.Values["user"] = setting.SessionData{Email: user.Email, ID: user.ID}
			session.Save(r, w)

			http.Redirect(w, r, "/select", http.StatusMovedPermanently)
		} else {
			log.Println(">>> invalid-credential")
			http.Redirect(w, r, "/login?error=invalid-credential", http.StatusMovedPermanently)
		}
	} else {
		tmpl := template.Must(template.ParseFiles("templates/layouts/login.html", "templates/login.html", ))

		tmpl.Execute(w, nil)
	}
}


func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		return
	}

	log.Println("Request ->", r.URL.Path)

	if r.Method == http.MethodPost {

		r.ParseForm()

		user  := models.User{}
		user.Add(r.Form)

		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)

	} else {
		tmpl := template.Must(template.ParseFiles("templates/layouts/login.html", "templates/register.html", ))

		tmpl.Execute(w, nil)
	}
}

func SelectQuiz(w http.ResponseWriter, r *http.Request)  {

	checkAuth(w,r)

	tmpl := template.Must(template.ParseFiles("templates/layouts/login.html", "templates/select.html", ))

	quizes := models.Quiz{}.GetAll()

	data := PageQuizData{Quizes: quizes}

	tmpl.Execute(w, data)


}


func StartQuiz(w http.ResponseWriter, r *http.Request)  {

	checkAuth(w,r)

	r.ParseForm()


	quizId, _ := strconv.ParseInt(r.Form.Get("quiz"),10,64)

	tmpl := template.Must(template.ParseFiles("templates/layouts/login.html", "templates/start.html", ))

	questions := models.Question{Quiz: quizId }.Get()

	data := PageQuestionsData{Questions: questions, Count: len(questions)}

	tmpl.Execute(w, data)


}


func checkAuth(w http.ResponseWriter, r *http.Request) {
	session, err := setting.Store.Get(r, "USER")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if session.Values["user"] == nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}

	log.Printf("SESS > %+v", session.Values["user"])

}
