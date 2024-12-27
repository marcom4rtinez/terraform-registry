package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/marcom4rtinez/terraform-registry/model"
)

func DownloadProvider(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	version := c.Param("version")
	osParam := c.Param("os")
	arch := c.Param("arch")
	filePath := filepath.Join(model.DataPath, namespace, name+".json")

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

	var provider *model.Platform
	for _, elem := range data.Versions {
		if elem.Version == version {
			for _, platform := range elem.Platforms {
				if platform.OS == osParam && platform.Arch == arch {
					provider = &model.Platform{
						OS:                  platform.OS,
						Arch:                platform.Arch,
						Filename:            platform.Filename,
						DownloadURL:         platform.DownloadURL,
						ShasumsURL:          platform.ShasumsURL,
						ShasumsSignatureURL: platform.ShasumsSignatureURL,
						Shasum:              platform.Shasum,
						SigningKeys:         platform.SigningKeys,
					}
					break
				}
			}

			if provider != nil {
				downloadData := gin.H{
					"os":                    provider.OS,
					"arch":                  provider.Arch,
					"filename":              provider.Filename,
					"download_url":          provider.DownloadURL,
					"shasums_url":           provider.ShasumsURL,
					"shasums_signature_url": provider.ShasumsSignatureURL,
					"shasum":                provider.Shasum,
					"signing_keys": gin.H{
						"gpg_public_keys": []gin.H{},
					},
				}

				for _, key := range provider.SigningKeys.GPGPublicKeys {
					downloadData["signing_keys"].(gin.H)["gpg_public_keys"] = append(downloadData["signing_keys"].(gin.H)["gpg_public_keys"].([]gin.H), gin.H{
						"key_id":          key.KeyID,
						"ascii_armor":     key.AsciiArmor,
						"trust_signature": key.TrustSignature,
						"source":          key.Source,
						"source_url":      key.SourceURL,
					})
				}

				c.JSON(http.StatusOK, downloadData)
				return
			}

		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
}
