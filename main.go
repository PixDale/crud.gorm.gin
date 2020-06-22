package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func main() {
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Person{})

	r := gin.Default()
	r.GET("/people", getPeople)
	r.GET("/people/:id", getPerson)
	r.POST("/people", createPerson)
	r.PUT("/people/:id", updatePerson)
	r.DELETE("/people/:id", deletePerson)

	r.GET("/populate", populateDB)

	r.Run(":8989")
}

func getPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}
func getPerson(c *gin.Context) {
	var p Person
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&p).Error; err != nil {
		fmt.Println("Erro ao procurar pessoa: " + err.Error())
	} else {
		c.JSON(200, p)
	}
}
func createPerson(c *gin.Context) {
	var p Person
	c.BindJSON(&p)
	db.Create(&p)
	c.JSON(200, p)
}
func updatePerson(c *gin.Context) {
	var p Person
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&p).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&p)
	db.Save(&p)
	c.JSON(200, p)
}
func deletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var p Person
	d := db.Where("id = ?", id).Delete(&p)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
func populateDB(c *gin.Context) {
	people := []Person{
		{0, "Joao", "Dias"},
		{0, "Elvis", "Presley"},
		{0, "Jackie", "Chan"},
		{0, "Herbert", "Richards"},
	}
	for _, p := range people {
		db.Create(&p)
	}
	c.String(200, "Banco populado com sucesso")
}
