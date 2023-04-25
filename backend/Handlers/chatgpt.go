package Handlers

import (
	"backend/Services"
	"backend/Types"
	"backend/db"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// HealthCheck ChatGPTAPI通信確認用
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	chatGPTResponse, err := Services.HealthCheck()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
	}

	if err = json.NewEncoder(w).Encode(chatGPTResponse); err != nil {
		log.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// Chat ChatGPTにMessageを送る
func Chat(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")

	chatGPTResponse, err := Services.Chat(content)
	if err != nil {
		log.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.NewEncoder(w).Encode(chatGPTResponse); err != nil {
		log.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func Embeddings(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")

	embeddings, err := Services.Embeddings(input)
	if err != nil {
		log.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.NewEncoder(w).Encode(embeddings); err != nil {
		log.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func CalcEmbeddings(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	inputEmbedded, err := Services.Embeddings(text)
	if err != nil {
		log.Println(err)
		return
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
			return
		}
		embeddingBytes, err := json.Marshal(embedding)
		if err = db.RedisClient.Conn.Set(instruction, embeddingBytes, 0).Err(); err != nil {
			log.Println(err)
			return
		}
	}

	keys, err := db.RedisClient.Conn.Keys("*").Result()
	if err != nil {
		log.Println(err)
		return
	}

	mcos := make(map[string]float64)

	for _, key := range keys {
		val, err := db.RedisClient.Conn.Get(key).Result()
		if err != nil {
			log.Println(err)
			return
		}
		var embedding *Types.Embedding
		err = json.Unmarshal([]byte(val), &embedding)
		cos, err := Services.CalcCosSimilarity(inputEmbedded, embedding)
		if err != nil {
			log.Println(err)
			return
		}
		mcos[key] = cos
	}

	mSorted := utils.InstructSortByCosin(mcos)

	var selectedInstructions []string
	for i, instruction := range mSorted {
		if i > 3 {
			break
		}
		selectedInstructions = append(selectedInstructions, instruction.Key)
	}

	prompt := utils.GenerateTemplate(strings.Join(selectedInstructions, "\n"), text)
	chatGPTResponse, err := Services.Chat(prompt)
	if err != nil {
		log.Println(err)
		return
	}

	if err = json.NewEncoder(w).Encode(chatGPTResponse); err != nil {
		log.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
