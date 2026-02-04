package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClickerBasicStructure struct {
	Record   int    `json:"record"`
	Nickname string `json:"nickname"`
}

// just for testing before real database
type TestingGetRequest struct {
	Names   string `json:"names"`
	Records int    `json:"records"`
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Api is working!",
		})
	})
	// router for put info in "ClickerBasicStructure":) kurwa!
	router.POST("/postInformation", func(c *gin.Context) {
		var clickerBaseST ClickerBasicStructure

		if err := c.ShouldBindJSON(&clickerBaseST); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad JSON request",
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "All good!",
		})
	})

	router.GET("/getDatabase", func(c *gin.Context) {
		var testingGT TestingGetRequest

		testingGT.Names = "John deer"
		testingGT.Records = 1785983

		c.JSON(200, gin.H{
			"Names: ":   testingGT.Names,
			"Records: ": testingGT.Records,
		})
	})

	router.Run()
}
