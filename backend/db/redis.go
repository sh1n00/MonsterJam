package db

import (
	"backend/Services"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

type RedisHandler struct {
	Conn *redis.Client
}

var RedisClient = NewRedisHandler()

func init() {
	if err := RedisClient.Conn.FlushDB().Err(); err != nil {
		log.Fatalln(err)
	}

	instructions := []string{
		"サイバネティクス・リアリティ工学研究室は、情報科学A棟3階にあります。",
		"インタラクティブメディア設計学研究室は、情報科学B棟7階にあります。",
		"自然言語処理学研究室は、情報科学A棟7階にあります。",
		"A棟のエレベーターはこの場所から左手奥に進んだところにあります。",
		"B棟のエレベーターはこの場所から右手に進んだところにあります。",
	}

	for _, instruction := range instructions {
		embedding, err := Services.Embeddings(instruction)
		if err != nil {
			log.Println(err)
		}
		embeddingBytes, err := json.Marshal(embedding)
		if err != nil {
			log.Println(err)
		}
		if err = RedisClient.Conn.Set(instruction, embeddingBytes, 0).Err(); err != nil {
			log.Println(err)
		}
	}
}

func NewRedisHandler() *RedisHandler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatalln(err)
	}
	return &RedisHandler{Conn: redisClient}
}
