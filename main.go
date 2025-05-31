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
	bot, _     = initbot()
	usersCodes = make(map[string]string)
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

	t, err := template.ParseFiles("template/discord_code.html", "template/header.html", "template/discord_login.html")

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

	source := r.FormValue("form_source")
	switch source {
	case "discord_form":

		userID := r.FormValue("discord_id")
		usersCodes[userID] = randomString(6)

		err = bot.messageCode(userID, usersCodes[userID])
		if err != nil {
			discordLGET(w, r, t, "discord", err.Error())
			return
		}

		cookie := &http.Cookie{
			Name:     "discord_id",
			Value:    userID,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)

		discordLGET(w, r, t, "code", "")

	case "code_form":

		cookie, err := r.Cookie("discord_id")
		if err != nil {
			discordLGET(w, r, t, "code", "Кука не найдена!")
			return
		}

		userID := cookie.Value

		codeN := r.FormValue("verification_code")
		fmt.Println("Проверка кода для:", userID, "Введенный код:", codeN, "Ожидаемый код:", usersCodes[userID])

		if codeN != usersCodes[userID] {
			discordLGET(w, r, t, "code", "Код не совпадает")
			return
		}
		delete(usersCodes, userID)
		cookie = &http.Cookie{
			Name:   "discord_id",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
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
