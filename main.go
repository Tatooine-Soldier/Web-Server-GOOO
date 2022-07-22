package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var tpl *template.Template
var ErrServeHTMLFile = errors.New("failed to retrive webpage")

type HelloWorld struct {
	Message string `json:"message"`
}

type Person struct {
	UserName string
	Password string
}

func main() {
	// e := echo.New()
	// // e.GET("/", Home)
	// e.GET("/contact", Contact)

	// g := e.Group("/user")
	// g.GET("/:name", Tom)
	// // g.POST("/:name", createUser)
	// g.Use(middleware.Logger())

	// http.Handle("/", http.FileServer(http.Dir("./assets")))
	// http.Handle("/login", http.FileServer(http.Dir("./assets ")))

	// http.ListenAndServe(":3000", nil)

	// e.GET("/params/:data", getParams)
	// e.Logger.Fatal(e.Start(":1323"))

	http.Handle("/", http.FileServer(http.Dir("./assets")))
	http.Handle("/login", http.FileServer(http.Dir("./assets ")))
	http.HandleFunc("/process", process)      //login process
	http.HandleFunc("/signup", processSignup) //signup process
	http.HandleFunc("/loggin", processLogin)  //signup process

	http.ListenAndServe(":3000", nil)
}

func init() {
	tpl = template.Must(template.ParseGlob("assets/*.gohtml"))
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func getParams(c echo.Context) error {
	person := Person{}

	err := c.Bind(&person)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed processing request")
	}

	datatype := c.Param("data")

	return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid parameter type: %v", datatype))
}

func process(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing...")
	parseForm(w, r)
}

func processLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing Login...")
	usr, err := parseForm(w, r)
	if err != nil {
		fmt.Println(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}

	usersCollection := client.Database("usernames").Collection("users")
	doc := bson.D{{"uid", usr.UserName}, {"username", usr.Password}}

	var results []bson.D
	cursor, err := usersCollection.Find(context.TODO(), doc)
	if err != nil {
		fmt.Println(err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

	// err = tpl.ExecuteTemplate(w, "in.gohtml", usr)
	// if err != nil {
	// 	return err
	// }

}

func processSignup(w http.ResponseWriter, r *http.Request) {
	usr, err := parseForm(w, r)
	if err != nil {
		fmt.Println(err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}

	usersCollection := client.Database("usernames").Collection("users")
	start := time.Now().Second()

	rand.Seed(int64(start))
	rn := rand.Intn(10000)
	fmt.Printf("%v", rn)

	user := bson.D{{"uid", usr.UserName}, {"password", usr.Password}}
	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.InsertedID)

}

func parseForm(w http.ResponseWriter, r *http.Request) (Person, error) {
	usr := r.FormValue("username")
	pw := r.FormValue("password")

	person := Person{
		UserName: usr,
		Password: pw,
	}

	if person.UserName == "" || person.Password == "" {
		return person, errors.New("username or password cannot be empty")
	}

	return person, nil
}

func Tom(c echo.Context) error {
	person := Person{}

	// datatype := c.Param("data")
	// if datatype != "json" {
	// 	return c.String(http.StatusBadRequest, "Error: invalid datatype")
	// }

	// status := c.QueryParam("status")
	// if status == "" {
	// 	return c.String(http.StatusBadRequest, "Error: no username recieved")
	// }

	err := c.Bind(&person)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed processing request")
	}

	return c.String(http.StatusOK, fmt.Sprintf("This is the param name you sent us '%s'", person))

}

func Home(c echo.Context) error {
	crud := c.QueryParam("action")
	if crud == "" {
		return c.String(http.StatusBadRequest, "Error: no action recieved")
	}

	err := connectDB(c, crud)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("err->  %v", err))
	}
	return nil
}

func Contact(c echo.Context) error {
	return c.String(http.StatusOK, "Contact!")
}

func connectDB(c echo.Context, crud string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return c.String(http.StatusFailedDependency, fmt.Sprintf("Error: cannot connect to MongoDB: %v", err))
	}

	switch crud {
	case "inserted":
		insertUsers(c, client)
		// case "deleted":
		// 	deletedUsers(c, client)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Successfully performed '%v' in MongoDB", crud))

}

func ping(client *mongo.Client) {
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
}

func insertUsers(c echo.Context, client *mongo.Client) error {
	usersCollection := client.Database("usernames").Collection("users")

	param := c.QueryParam("username")
	if param == "" {
		return c.String(http.StatusBadRequest, "Error: no username specified")
	}

	start := time.Now().Second()

	rand.Seed(int64(start))
	rn := rand.Intn(10000)
	user := bson.D{{"uid", rn}, {"username", param}}
	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.InsertedID)
	return nil
}

// func deletedUsers(c echo.Context, client *mongo.Client) error {
// 	usersCollection := client.Database("usernames").Collection("users")

// 	param := c.QueryParam("username")
// 	if param == "" {
// 		return c.String(http.StatusBadRequest, "Error: no username specified")
// 	}

// 	result := usersCollection.FindOneAndDelete(context.TODO(), bson.M{"username": param})

// 	fmt.Printf("Found %v", result.Decode())
// 	return nil
// }
