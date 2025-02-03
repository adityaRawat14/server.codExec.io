package handler

import (
	"context"
	"log"
	"net/http"
	"server/models"
	"server/services/util/k8s"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleCodeExecution(c *gin.Context) {
	var codeData models.CodeData
	if err := c.ShouldBindJSON(&codeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid request!!"})
		return
	}
	if _, ok := models.AllowedLanguages[codeData.Language]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "Language not supported"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	resultChan := make(chan string)

	go func() {
		k8s.Execute(codeData.Language, codeData.Body, resultChan, ctx)
	}()

	select {
	case result := <-resultChan:
		if strings.HasPrefix(result, "failed") {
			c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": result})
		} else {
			log.Println("Code executed: " + result)
			c.JSON(http.StatusOK, gin.H{"response": result})
		}
	case <-ctx.Done():
		log.Println("Request timed out")
		c.JSON(http.StatusGatewayTimeout, gin.H{"errorMessage": "Request timed out"})
	}
}
