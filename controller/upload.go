package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/marcom4rtinez/terraform-registry/model"
)

func UploadProvider(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	filePath := filepath.Join(model.DataPath, namespace, name+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Couldn't create provider"})
			return
		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			err := os.WriteFile(filePath, []byte(defaultProviderContent), 0644)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Couldn't create record"})
				return
			}

		}
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	var data model.ProviderData
	if err := json.Unmarshal(file, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	var newVersion model.Version
	if err := c.BindJSON(&newVersion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if data.Versions[0].Version != "" {
		for _, elem := range data.Versions {
			if elem.Version == newVersion.Version {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate release"})
				return
			}
		}
		data.Versions = append(data.Versions, newVersion)
	} else {
		data.Versions[0] = newVersion
	}

	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize JSON"})
		return
	}

	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider version added successfully"})

}

const defaultProviderContent = `
{
	"versions": [
		{}
	]
}
`
