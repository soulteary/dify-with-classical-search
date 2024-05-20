// base on: https://github.com/soulteary/dify-simple-rag-with-wp/blob/main/main.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryPayload struct {
	Queries []QueryBody `json:"queries"`
}

type QueryBody struct {
	IndexUID string `json:"indexUid"`
	Q        string `json:"q"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

type SearchResults struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	IndexUID string `json:"indexUid"`
	Hits     []struct {
		ID          int      `json:"id"`
		Title       string   `json:"title"`
		Overview    string   `json:"overview"`
		Genres      []string `json:"genres"`
		Poster      string   `json:"poster"`
		ReleaseDate int      `json:"release_date"`
	} `json:"hits"`
	Query              string `json:"query"`
	ProcessingTimeMs   int    `json:"processingTimeMs"`
	Limit              int    `json:"limit"`
	Offset             int    `json:"offset"`
	EstimatedTotalHits int    `json:"estimatedTotalHits"`
}

func GetSearchResult(search string, count int, indexes string, page int, token string) (result SearchResults, err error) {
	client := &http.Client{}

	var queryBody QueryBody
	queryBody.IndexUID = indexes
	queryBody.Q = search
	queryBody.Limit = count
	queryBody.Offset = page

	var queryPayload QueryPayload
	queryPayload.Queries = append(queryPayload.Queries, queryBody)

	payload, err := json.Marshal(queryPayload)
	if err != nil {
		return result, err
	}

	var data = strings.NewReader(string(payload))
	req, err := http.NewRequest("POST", "http://localhost:7700/multi-search", data)
	if err != nil {
		return result, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

type ExtensionPointRequest struct {
	Point  string `json:"point"`
	Params struct {
		AppID        string                 `json:"app_id"`
		ToolVariable string                 `json:"tool_variable"`
		Inputs       map[string]interface{} `json:"inputs"`
		Query        string                 `json:"query"`
	} `json:"params"`
}

type ExtensionPointResponse struct {
	Result string `json:"result"`
}

func main() {
	router := gin.Default()

	router.POST("/new-api-for-dify", func(c *gin.Context) {
		var req ExtensionPointRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Point == "ping" {
			c.JSON(http.StatusOK, ExtensionPointResponse{Result: "pong"})
			return
		}

		if req.Point != "app.external_data_tool.query" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid point"})
			return
		}

		keywords, exist := req.Params.Inputs["keywords"]
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing keyword"})
			return
		}

		s := strings.TrimSpace(fmt.Sprintf("%s", keywords))
		if s == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty keyword"})
			return
		}

		movies, err := GetSearchResult(s, 3, "movies", 0, "soulteary")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var result string
		for _, movie := range movies.Results {
			for _, hit := range movie.Hits {
				result += fmt.Sprintf("- 标题：%s\n", hit.Title)
				result += fmt.Sprintf("- 简介：%s\n\n", hit.Overview)
			}
		}

		c.JSON(http.StatusOK, ExtensionPointResponse{Result: result})
	})

	router.Run(":8084")
}
