package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"crypto/md5"
	"encoding/hex"
)

type user struct {
	Id        int     	`sql:"AUTO_INCREAMENT" gorm:"primary_key"`
	User      string     	`sql:"size:255;unique;index"`
	Password  string     	`sql:"size:255;unique;index"`
	Type      string     	`sql:"size:255;unique;index"`
}

type data struct {
	Id  		int	`sql:"AUTO_INCREAMENT" gorm:"primary_key"`
	Name 		string	`sql:"size:255;unique;index"`
	Surname 	string	`sql:"size:255;unique;index"`
	Phonenumber 	string	`sql:"size:255;unique;index"`
	Email 		string	`sql:"size:255;unique;index"`
	Domicile 	string	`sql:"size:255;unique;index"`
	Region 		string	`sql:"size:255;unique;index"`
	School 		string	`sql:"size:255;unique;index"`
	Fullort 	string	`sql:"size:255;unique;index"`
	Mathort 	string	`sql:"size:255;unique;index"`
}

type session struct {
	In string
}

var database []*data = make([]*data,0)
var users map[string]*user = make(map[string]*user,0)
var db *gorm.DB
var err error
var sessionData data = data{Id:-1}
var sessionUser user

var in session = session{"out"}

// index page where starts application////////////////////////////////////
func indexHandler(rdr render.Render, r *http.Request)  {

	rdr.HTML(200,"index", in)
}

//login page where users log in////////////////////////////////
func loginHandler(rdr render.Render)  {
	rdr.HTML(200, "login", in)
}
func registrationHandler(rdr render.Render){
	rdr.HTML(200, "registration", nil)
}
// input page students enter information about them////////////
func inputHandler(rdr render.Render)  {
	if sessionData.Id == -1 {
		rdr.HTML(200, "input", nil)
	}else {
		rdr.HTML(200, "input", sessionData)
	}
}

// output page where we show data//////////////////////////////
func showHandler(rdr render.Render)  {

	db.Find(&database)
	rdr.HTML(200, "show", database)
}

// checkInput page checks our input/////////////////////////////
func checkInputHandler(rdr render.Render, r *http.Request)  {
	data := data{0,r.FormValue("name"),r.FormValue("surname"),
		r.FormValue("phonenumber"),r.FormValue("email"),
		r.FormValue("domicile"),r.FormValue("region"),
		r.FormValue("school"),r.FormValue("fullort"),
		r.FormValue("mathort")}

	fmt.Println(data)

	sessionData = data
	rdr.HTML(200,"form-send-confirm",data)
}

// checkLogin page checks username and password/////////////////
func checkLoginHandler(rdr render.Render, r *http.Request)  {
	username := r.FormValue("user")
	password := GetMD5Hash(r.FormValue("password"))

	u := []user{}
	db.Find(&u)
	for _,value := range u{
		if username == value.User && password == value.Password {
			in = session{""}
			rdr.HTML(200, "index", in)
		}else {
			rdr.HTML(200, "login", in)
		}
	}
}
//////////////////////////\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
/////////////////// MAIN FUNCTION \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
//////////////////////////\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
func main()  {
	fmt.Println("Application starts!")

	connect()



	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory: "html",
		Extensions:[]string{".tmpl",".html"},
		Charset: "UTF-8",
		IndentJSON:true,
	}))

	staticOption := martini.StaticOptions{Prefix:"js"}
	m.Use(martini.Static("js", staticOption))
	staticOption2 := martini.StaticOptions{Prefix:"css"}
	m.Use(martini.Static("css", staticOption2))
	m.Get("/",indexHandler)
	m.Get("/login", loginHandler)
	m.Get("/input", inputHandler)
	m.Get("/show", showHandler)
	m.Get("/registration", registrationHandler)
	m.Post("/",indexHandler)

	/////////////////CHECK INPUT////////////////////
	m.Get("/form-send-confirm",checkInputHandler)
	m.Post("/form-send-confirm",checkInputHandler)
	m.Get("/enterData", func(rdr render.Render) {
		if sessionData.Id != -1 {
			db.Create(&sessionData)
			sessionData.Id = -1
		}
		rdr.HTML(200, "index", in)
	})
	m.Get("/destroy", func(rdr render.Render) {
		sessionData.Id = -1
		rdr.HTML(200, "input", nil)
	})

	//////////////////CHECK LOGIN//////////////////////
	m.Post("/checkLogin", checkLoginHandler)
	m.Get("/logout", func(rdr render.Render) {
		in = session{"out"}
		rdr.HTML(200, "index", in)
	})

	/////////////////CHECK REGIST///////////////////////
	m.Post("/checkRegist", func(rdr render.Render, r *http.Request) {
		username := r.FormValue("user")
		password := r.FormValue("password")
		if password == r.FormValue("copypassword") {
			u := user{User: username,Password: GetMD5Hash(password),Type: "admin"}
			db.Create(&u)
			rdr.HTML(200,"index",in)
		}else {
			rdr.HTML(200, "registration", in)
		}
	})
	m.Run()
}


/////////////////////////////////////////////////////////////////////////////////////
///////////////////////// DATABASE CONNECTION ///////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////
func connect()  {

	db, err = gorm.Open("postgres","user=postgres password=root dbname=godb host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.CreateTable(user{})
	db.CreateTable(data{})
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}