package main

import (
	"context"
	"fmt"
	"math/rand"
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

	rand.Seed(20)
	rn := rand.Intn(1000)
	user := bson.D{{"uid", rn}, {"username", param}}
	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.InsertedID)
	return nil
}
