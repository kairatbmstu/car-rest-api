package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var carRepository = CarRepository{}

func main() {
	// Set up Gin router
	router := gin.Default()
	// Run the server
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))
	router.POST("/cars", CreateCar)
	router.PUT("/cars/:id", UpdateCar)
	router.DELETE("/cars/:id", DeleteCar)
	router.GET("/cars/:id", FindCardById)
	router.GET("/cars", GetAll)

	log.Println("Starting server on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}

func CreateCar(c *gin.Context) {
	var car Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	// Insert post into the database
	carCreated := carRepository.Create(car)
	c.JSON(http.StatusOK, carCreated)
}

func UpdateCar(c *gin.Context) {
	var car Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	// Insert post into the database
	carUpdated := carRepository.Update(car)
	c.JSON(http.StatusOK, carUpdated)
}

func DeleteCar(c *gin.Context) {
	idPrm := c.Param("id")
	id, err := strconv.Atoi(idPrm)

	// Check for errors during conversion
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameter"})
		return
	}

	carRepository.Delete(id)
	c.JSON(http.StatusOK, gin.H{})
}

func FindCardById(c *gin.Context) {
	idPrm := c.Param("id")
	id, err := strconv.Atoi(idPrm)

	// Check for errors during conversion
	if err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameter"})
		return
	}

	car := carRepository.FindById(id)
	c.JSON(http.StatusOK, car)
}

func GetAll(c *gin.Context) {
	if len(carRepository.cars) == 0 {
		cars := make([]Car, 0)
		c.JSON(http.StatusOK, cars)
	} else {
		c.JSON(http.StatusOK, carRepository.cars)
	}
}

type Car struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

type CarRepository struct {
	counter int
	cars    []Car
}

func (r *CarRepository) FindById(id int) *Car {
	for _, c := range r.cars {
		if c.Id == id {
			return &c
		}
	}

	return nil
}

func (r *CarRepository) Create(car Car) *Car {
	r.counter++
	car.Id = r.counter
	r.cars = append(r.cars, car)
	return &car
}

func (r *CarRepository) Update(car Car) *Car {
	for _, c := range r.cars {
		if c.Id == car.Id {
			c.Name = car.Name
			c.Year = car.Year
			return &car
		}
	}

	return nil
}

func (r *CarRepository) Delete(id int) {
	indexToDelete := -1
	for i, c := range r.cars {
		if c.Id == id {
			indexToDelete = i
		}
	}

	if indexToDelete > 0 {
		r.cars = append(r.cars[:indexToDelete], r.cars[indexToDelete+1:]...)
	}
}
