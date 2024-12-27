package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/marcom4rtinez/terraform-registry/model"
)

func GetVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	filePath := filepath.Join("data/providers", namespace, name+".json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
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

	response := gin.H{"versions": []gin.H{}}
	for _, elem := range data.Versions {
		version := gin.H{
			"version":   elem.Version,
			"protocols": elem.Protocols,
			"platforms": []gin.H{},
		}

		for _, platform := range elem.Platforms {
			platformData := gin.H{
				"os":                    platform.OS,
				"arch":                  platform.Arch,
				"filename":              platform.Filename,
				"download_url":          platform.DownloadURL,
				"shasums_url":           platform.ShasumsURL,
				"shasums_signature_url": platform.ShasumsSignatureURL,
				"shasum":                platform.Shasum,
				"signing_keys": gin.H{
					"gpg_public_keys": []gin.H{},
				},
			}

			for _, key := range platform.SigningKeys.GPGPublicKeys {
				platformData["signing_keys"].(gin.H)["gpg_public_keys"] = append(platformData["signing_keys"].(gin.H)["gpg_public_keys"].([]gin.H), gin.H{
					"key_id":          key.KeyID,
					"ascii_armor":     key.AsciiArmor,
					"trust_signature": key.TrustSignature,
					"source":          key.Source,
					"source_url":      key.SourceURL,
				})
			}

			version["platforms"] = append(version["platforms"].([]gin.H), platformData)
		}

		response["versions"] = append(response["versions"].([]gin.H), version)
	}

	c.JSON(http.StatusOK, response)
}
