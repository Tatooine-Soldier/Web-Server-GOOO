package main

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HelloWorld struct {
	Message string `json:"message"`
}

type Person struct {
	Name string `json:"name"`
}

const (
	DB_HOST = "127.0.0.1"
	DB_USER = "root"
	DB_PASS = "pass"
	DB_NAME = "name"
)

func main() {
	e := echo.New()
	e.GET("/", Home)
	e.GET("/contact", Contact)

	g := e.Group("/user")
	g.GET("/:name", Tom)
	g.Use(middleware.Logger())

	e.GET("/params/:data", getParams)
	e.Logger.Fatal(e.Start(":1323"))
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

	return c.String(http.StatusOK, fmt.Sprintf("This is the param name you sent us '%s'", person.Name))

}

func Home(c echo.Context) error {
	err := connectDB(c)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("err->  %v", err))
	}
	return c.String(http.StatusOK, "Welcome Home!")
}

func Contact(c echo.Context) error {
	return c.String(http.StatusOK, "Contact!")
}

func connectDB(c echo.Context) error {

	// db, err := sql.Open("mysql", "root:YES@tcp(127.0.0.1:3306)/users")
	// if err != nil {
	// 	log.Fatal(err)
	// 	fmt.Print("error")
	// }
	// defer db.Close()

	// insert, err := db.Query("INSERT INTO users VALUES(1, 'goman', 'gains')")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// defer insert.Close()

	// fmt.Print("Successfully connected to database")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	usersCollection := client.Database("testing").Collection("users")

	user := bson.D{{"id", 0}, {"fname", "TEST"}, {"lname", "TEST"}}
	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
	return nil
}
