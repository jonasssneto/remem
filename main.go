package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/martinlindhe/notify"
)

type Reminder struct {
	Title    string
	Message  string
	Minutes  int
	Repeat   int
	NotifyAt time.Time
}

func (rem Reminder) TimeRemaining() string {
	remaining := time.Until(rem.NotifyAt)
	if remaining < 0 {
		return "Lembrete já disparado"
	}
	return remaining.Round(time.Second).String()
}

var reminders []Reminder

var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Erro ao carregar o template: %v", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addReminderHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Servidor rodando em http://localhost:6606")
	log.Fatal(http.ListenAndServe(":6606", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, reminders)
	if err != nil {
		http.Error(w, "Erro ao renderizar a página", http.StatusInternalServerError)
	}
}

func addReminderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	message := r.FormValue("message")
	minutes, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		http.Error(w, "Tempo inválido", http.StatusBadRequest)
		return
	}
	repeat, err := strconv.Atoi(r.FormValue("repeat"))
	if err != nil {
		http.Error(w, "Repetição inválida", http.StatusBadRequest)
		return
	}

	notifyAt := time.Now().Add(time.Duration(minutes) * time.Minute)

	reminder := Reminder{Title: title, Message: message, Minutes: minutes, Repeat: repeat, NotifyAt: notifyAt}
	reminders = append(reminders, reminder)

	go scheduleReminder(reminder)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func scheduleReminder(rem Reminder) {
	time.Sleep(time.Until(rem.NotifyAt))
	notify.Notify("Lembrete", rem.Title, rem.Message, "")

	if rem.Repeat > 0 {
		for {
			time.Sleep(time.Duration(rem.Repeat) * time.Minute)
			notify.Notify("Lembrete", rem.Title, rem.Message, "")
		}
	}
}
