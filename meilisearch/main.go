package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Overview    string   `json:"overview"`
	Genres      []string `json:"genres"`
	Poster      string   `json:"poster"`
	ReleaseDate int      `json:"release_date"`
}

func main() {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://127.0.0.1:7700",
		APIKey: "soulteary",
	})

	// 创建一个名为 'movies' 的索引，用来存储后续的数据
	index := client.Index("movies")

	// 如果索引 'movies' 不存在，Meilisearch 会在第一次添加文档时创建它
	// documents := []map[string]interface{}{
	// 	{"id": 1, "title": "Carol", "genres": []string{"Romance", "Drama"}},
	// 	{"id": 2, "title": "Wonder Woman", "genres": []string{"Action", "Adventure"}},
	// 	{"id": 3, "title": "Life of Pi", "genres": []string{"Adventure", "Drama"}},
	// }

	buf, err := os.ReadFile("data/movies.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var documents []Movie
	err = json.Unmarshal(buf, &documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(task.TaskUID)
}
