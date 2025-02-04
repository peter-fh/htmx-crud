package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"todo/Database"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var db *mongo.Collection

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var tmpl *template.Template



func init() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	check(err)
	
	godotenv.Load()
	mongodb_user := os.Getenv("DATABASE_USERNAME")
	mongodb_password := os.Getenv("DATABASE_PASSWORD")
	mongodb_uri := os.Getenv("DATABASE_URI")
	Database.DropAll(mongodb_user, mongodb_password, mongodb_uri)
	db = Database.GetTasksTable(mongodb_user, mongodb_password, mongodb_uri)
	
}

func Homepage (w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	task := r.FormValue("task")
	err := Database.AddTask(db, task)
	if err != nil {
		log.Fatalf("Failed to add task: %v", err)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
}



func TaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("hi")
	switch r.Method {
	case http.MethodPost:
		AddTask(w, r)
	default:
		log.Fatalln("Incorrect method to /task")
	}
	allTasks := Database.GetAllTasks(db)
	for _, task := range allTasks {
		fmt.Println(task.Name)
	}
}

func main() {

	http.HandleFunc("/", Homepage)
	http.HandleFunc("/task", TaskHandler)
	http.ListenAndServe(":4000", nil)
}
