package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Name      string   `json:"name"`
	Age       uint16   `json:"age"`
	Money     float32  `json:"money"`
	AvgGrades float64  `json:"avgGrades"`
	Happiness float64  `json:"happiness"`
	Hobbies   []string `json:"hobbies"`
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("Name: %s\nAge: %d\nMoney: %.2f\nAvg Grades: %.2f\nHappiness: %.2f", u.Name, u.Age, u.Money, u.AvgGrades, u.Happiness)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

type Article struct {
	Id      uint16
	Title   string
	Anons   string
	Content string
}

// Определяем глобальные переменные
var (
	db  *sql.DB
	err error
)

// Определям первичную инициализацию
func init() {
	// Открываем файл базы данных SQLite или создаем его, если он не существует
	db, err = sql.Open("sqlite3", "./webgo.db")
	if err != nil {
		panic(err)
	}

	// Создание таблицы users (если она еще не существует)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50),
		age INTEGER
	);`)
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы articles (если она еще не существует)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS articles (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(100),
		anons VARCHAR(255),
		content TEXT
	);`)

}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.gohtml", "templates/header.gohtml", "templates/footer.gohtml")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Получаем все записи из таблицы articles
	rows, err := db.Query("SELECT * FROM articles")
	if err != nil {
		log.Fatal(err)
	}

	var posts []Article
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.Id, &article.Title, &article.Anons, &article.Content)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, article)

	}

	tmpl.ExecuteTemplate(w, "index", posts)
}

func createPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.gohtml", "templates/header.gohtml", "templates/footer.gohtml")
	if err != nil {
		fmt.Fprintf(w, err.Error())

	}
	tmpl.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	anons := r.FormValue("anons")
	content := r.FormValue("content")

	if title == "" || anons == "" || content == "" {
		fmt.Fprintf(w, "Заполните все поля!")
	} else {
		_, err = db.Exec("INSERT INTO articles (title, anons, content) VALUES (?,?,?)", title, anons, content)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {

	user := User{}
	user.Name = "John"
	user.Hobbies = []string{"Sports", "Reading", "Sleeping"}

	/*
		fmt.Fprintf(w, "Going to home page!")
		fmt.Fprintf(w, "\n\nHello, %s!\n", user.name)
		fmt.Fprintf(w, user.getAllInfo())

		user.setNewName("Jane")
		fmt.Fprintf(w, "\n\nHello, %s!\n", user.name)
		fmt.Fprintf(w, user.getAllInfo())
	*/

	tmpl, _ := template.ParseFiles("templates/indexold.gohtml")
	tmpl.Execute(w, user)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Going to about page!")
}

func postPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tmpl, err := template.ParseFiles("templates/post.gohtml", "templates/header.gohtml", "templates/footer.gohtml")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Получаем все записи из таблицы articles
	rows, err := db.Query("SELECT * FROM articles WHERE id =?", vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	var showPost []Article
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.Id, &article.Title, &article.Anons, &article.Content)
		if err != nil {
			log.Fatal(err)
		}

		showPost = append(showPost, article)

	}

	tmpl.ExecuteTemplate(w, "post", showPost)

}

func handleRoutes() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", indexPage).Methods("GET")
	rtr.HandleFunc("/home/", homePage).Methods("GET")
	rtr.HandleFunc("/about/", aboutPage).Methods("GET")
	rtr.HandleFunc("/create/", createPage).Methods("GET")
	rtr.HandleFunc("/post/{id:[0-9]+}", postPage).Methods("GET")

	rtr.HandleFunc("/save_article", save_article).Methods("POST")

	http.Handle("/", rtr)

	// Обработка статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	/* Стандартный механизм
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/home/", homePage)
	http.HandleFunc("/about/", aboutPage)

	http.HandleFunc("/create/", createPage)
	http.HandleFunc("/save_article", save_article)

	*/
}

func main() {
	handleRoutes()
	http.ListenAndServe(":5050", nil)
}
