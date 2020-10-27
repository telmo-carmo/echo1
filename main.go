/*
go get -u github.com/labstack/echo/...


https://echo.labstack.com/cookbook/jwt
import "github.com/dgrijalva/jwt-go"
e.Use(middleware.JWT([]byte("secret key")))
--

https://github.com/labstack/echo/blob/master/middleware/logger.go#L137

e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
  Format: "method=${method}, uri=${uri}, status=${status}\n",
}))

DefaultLoggerConfig = LoggerConfig{
  Skipper: DefaultSkipper,
  Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
    `"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
    `"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
    `"bytes_out":${bytes_out}}` + "\n",
  CustomTimeFormat: "2006-01-02 15:04:05.00000",
  Output: os.Stdout
}


see also:  https://github.com/ribice/gorsk

*/

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"tac/echo1/dal"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-playground/validator/v10"
)

const (
	CT_JSON = "application/json"
)

// H is like a echo.Map interface
type H map[string]interface{}

/****
// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
***/

func getDbPath() string {
	dbPath := os.Getenv("LOCAL_DB_PATH")
	if dbPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dbPath = filepath.Join(cwd, "sqlite_scott.db")
	}
	return dbPath
}

func test1() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.Printf("args len is %d\n", len(os.Args))
	repo, err := dal.NewRepo(getDbPath())
	if err != nil {
		log.Panic(err)
	}
	v := repo.GetVersion()
	log.Println("Ending!! ver: " + v)
	repo.Close()
	os.Exit(0)
	fmt.Printf("The end")
}

var repo *dal.Repo

func health(c echo.Context) error {
	n := c.QueryParam("n")
	log.Printf("n is '%v'", n)
	v := repo.GetVersion()
	c.Logger().Warnf("v is %s", v)
	return c.String(http.StatusOK, "OK-SQLITE:"+v)
}

func bonusOneHd(c echo.Context) error {
	id := c.Param("id")
	log.Printf("id is '%v'", id)
	v := repo.GetBonus(id)
	c.Logger().Infof("bonus[] is %s", v)
	c.Response().Header().Set("Content-Type", CT_JSON)
	c.Response().WriteHeader(http.StatusOK)
	_, err := c.Response().Write([]byte(v))
	return err
}

func bonusHd(c echo.Context) error {
	v := repo.GetBonus("")
	c.Logger().Infof("bonus[] is %s", v)
	c.Response().Header().Set("Content-Type", CT_JSON)
	c.Response().WriteHeader(http.StatusOK)
	_, err := c.Response().Write([]byte(v))
	return err
}

func HelloHd(c echo.Context) error {
	return c.JSON(http.StatusOK, H{
		"code":    1002,
		"message": "user successfully updated",
		"timeUtc": time.Now().UTC(),
	})
}

func chart1Hd(c echo.Context) error {
	return c.Render(http.StatusOK, "chart1", echo.Map{
		"title":  "Chart1 Title!",
		"Label1": "Receitas",
		"Vals1":  []int{10, 9, 8, 7, 6, 4, 7, 8},
		"Vals2":  []int{1, 3, 6, 5, 4, 2, 3, 5},
	})
}

func indexHd(c echo.Context) error {
	c.Logger().Warn("in index Handler")
	return c.Render(http.StatusOK, "index", echo.Map{
		"title": "Index Title!",
		"name":  "Olá Mundo Novo"})
}

func upload(c echo.Context) error {
	// Read form fields
	name := c.FormValue("name")

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fnPath := filepath.Join(cwd, "docs", file.Filename)
	log.Printf("[INFO] Uploading file '%s' to path '%s'", file.Filename, fnPath)

	// Destination
	dst, err := os.Create(fnPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	html1 := fmt.Sprintf("<html><p>File %s uploaded successfully with fields name=%s.</p><p>See file <a href='/docs/%s'>here</a></p></html>",
		file.Filename, name, file.Filename)
	return c.HTML(http.StatusOK, html1)
}

type (
	User1 struct {
		Name  string `json:"name" form:"name" validate:"required"`
		Email string `json:"email" form:"email" validate:"required,email"`
		Age   uint8  `json:"age" form:"age" validate:"gte=0,lte=100"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func form1Hd(c echo.Context) error {
	u := new(User1)

	if err := c.Bind(u); err != nil {
		return err
	}
	log.Printf("[INFO] user1 is %v", *u)

	if err := c.Validate(u); err != nil {
		log.Printf("[ERROR] Validate: %v", err)
		return c.JSON(http.StatusOK, H{
			"type":  "validate",
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, u)
}

func routes(e *echo.Echo) {
	e.Logger.Info("Creating routes")
	e.GET("/idx", indexHd)
	e.GET("/page", func(c echo.Context) error {
		//render only file, must full name with extension
		return c.Render(http.StatusOK, "index1.html", echo.Map{"title": "Page1 Title!!", "name": "Page 1"})
	})
	e.GET("/chart1", chart1Hd)
	e.GET("/hello", HelloHd)
	e.POST("/upload", upload)
	e.GET("/health", health)
	e.GET("/api/bonus", bonusHd)
	e.GET("/api/bonus/:id", bonusOneHd)
	e.GET("/rt", func(c echo.Context) error {
		return c.JSON(http.StatusOK, H{"rt": e.Routes()})
	})
	e.POST("/form1", form1Hd)
}

func main() {
	var err error
	var aPort string

	flag.StringVar(&aPort, "port", "5000", "server port number")
	verbP := flag.Bool("verb", false, "be verbose log")

	flag.Parse()

	repo, err = dal.NewRepo(getDbPath())
	if err != nil {
		log.Panic(err)
	}

	e := echo.New()
	if *verbP {
		e.Use(middleware.Logger())
	} else {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format:           "${time_custom}: ${method} ${uri} ${status}, latency_h=${latency_human}\n",
			CustomTimeFormat: "15:04:05.000",
		}))
	}

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Debug = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Static("/", "static")
	e.Static("/docs", "docs")
	//
	e.Renderer = echoview.New(goview.Config{
		Root:      "templates", // default: "views"
		Extension: ".html",
		Master:    "layouts/master",
		Funcs: template.FuncMap{
			"copyDt": func() string {
				return time.Now().Format("02-01-2006") //t.Format("2006-01-02T15:04:05.999999-07:00"))
			},
		},
		DisableCache: false,
	})
	//e.Renderer = echoview.Default()

	routes(e)

	e.Logger.Fatal(e.Start(":" + aPort)) // e.StartTLS(":8443", "cert.pem", "key.pem")
}
