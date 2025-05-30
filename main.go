package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type PageData struct {
	Error string
	Step  string
}

var (
	bot, _ = createBot("")
	ver    = true
	code   string
)

func main() {
	handler()
}

func handler() {
	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("discord_code.html", "header.html", "discord_login.html")

	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET")
		discordLGET(w, r, t, "discord", "")
	case http.MethodPost:
		fmt.Println("POST")
		discordLPOST(w, r, t)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	/*
		if r.Method == http.MethodGet {
			data := PageData{}
			tmpl.Execute(w, data)
			return
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
				return
			}

			username := r.FormValue("username")
			password := r.FormValue("password")

			storedPassword, ok := users[username]

			if !ok || storedPassword != password {
				data := PageData{
					Error: "Неверное имя пользователя или пароль",
				}
				tmpl.Execute(w, data)
				return
			}

			w.Write([]byte("Добро пожаловать, " + username + "!"))
			return
		}
	*/
}

func discordLGET(w http.ResponseWriter, r *http.Request, t *template.Template, step string, err string) {
	data := PageData{
		Error: err,
		Step:  step}
	t.ExecuteTemplate(w, "header", data)
	return
}

func discordLPOST(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	if ver == true {
		ver = false
		code = randomString(6)
	}

	source := r.FormValue("form_source")
	switch source {
	case "discord_form":

		userID := r.FormValue("discord_id")

		err = bot.messageCode(userID, code)
		if err != nil {
			discordLGET(w, r, t, "discord", err.Error())
			return
		}

		discordLGET(w, r, t, "code", "")

	case "code_form":
		codeN := r.FormValue("verification_code")

		if codeN != code {
			discordLGET(w, r, t, "code", "Код не совпадает")
			return
		}
		ver = true
		fmt.Fprintf(w, "Верификация пройдена")

	default:
		fmt.Println("default")
		return
	}

}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
