package main

import (
	"assignment_2/models"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/antonlindstrom/pgstore"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var store *pgstore.PGStore
var db *gorm.DB

const (
	DBHost         = "db"
	DBUserName     = "user"
	DBUserPassword = "password"
	DBName         = "user"
	DBPort         = "5432"
)

func newPostgresConnection() *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", DBHost, DBUserName, DBUserPassword, DBName, DBPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(0)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func newPostgresStore() *pgstore.PGStore {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DBUserName, DBUserPassword, DBHost, DBPort, DBName)
	authKey := []byte("my-auth-key-very-secret")
	encryptionKey := []byte("my-encryption-key-very-secret123")

	var err error
	store, err = pgstore.NewPGStore(url, authKey, encryptionKey)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	return store
}

func main() {
	gob.Register(models.User{})

	store = newPostgresStore()
	db = newPostgresConnection()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	http.HandleFunc("/", routeIndex)
	http.HandleFunc("/register", routeRegister)
	http.HandleFunc("/logout", routeLogout)

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, _ := store.Get(r, "session")
		val := session.Values["user"]
		user, ok := val.(models.User)

		// if not signed in, redirect to signin page
		if !ok {
			var tmpl = template.Must(template.New("signin").ParseFiles("view/signin.html"))
			var err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// if already signed in, go to home page
		var tmpl = template.Must(template.New("home").ParseFiles("view/home.html"))
		var err = tmpl.Execute(w, user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if r.Method == "POST" {
		// var tmpl = template.Must(template.New("home").ParseFiles("view/home.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var username = r.FormValue("username")
		var password = r.Form.Get("password")

		user := models.User{}

		// check username
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var tmpl = template.Must(template.New("signin").ParseFiles("view/signin.html"))
				var err = tmpl.Execute(w, models.View{
					Message: "username not found",
				})

				if err != nil {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}

				return
			}
		}

		// check password
		if ok := CheckPasswordHash(password, user.Password); !ok {
			var tmpl = template.Must(template.New("signin").ParseFiles("view/signin.html"))
			var err = tmpl.Execute(w, models.View{
				Message: "wrong password",
			})

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// session struct has field Values map[interface{}]interface{}
		session.Values["user"] = user
		// save before writing to response/return from handler
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("register").ParseFiles("view/register.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var firstname = r.FormValue("firstname")
		var lastname = r.Form.Get("lastname")
		var username = r.Form.Get("username")
		var password = r.Form.Get("password")

		//save to database
		hashedPassword, _ := HashPassword(password)

		user := models.User{}

		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
		}

		if user.ID > 0 {
			var tmpl = template.Must(template.New("register").ParseFiles("view/register.html"))
			var err = tmpl.Execute(w, models.View{
				Message: "username already exists",
			})

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			return
		}

		user = models.User{
			FirstName: firstname,
			LastName:  lastname,
			Password:  hashedPassword,
			Username:  username,
		}

		if err := db.Create(&user).Error; err != nil {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete the session
		session.Options.MaxAge = -1
		session.Save(r, w)

		var tmpl = template.Must(template.New("signin").ParseFiles("view/signin.html"))
		err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
