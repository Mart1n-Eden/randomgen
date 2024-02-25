package server

import (
	"net/http"
	"encoding/json"
	"math/rand"
	"github.com/google/uuid"

	"github.com/Mart1n-Eden/randomgen/internal/database"
)

type GenRequest struct {
	Type   string `json:"type"`
	Length int    `json:"length"`
}

type GenResponse struct {
	ID  string `json:"id"`
	Val string    `json:"value"`
}

func Run() {
	database.Init(GenResponse{})

	http.HandleFunc("/api/generate", Generate)
	http.HandleFunc("/api/retrieve", Retrieve)

	http.ListenAndServe(":8080", nil)
}

func Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req GenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	value := req.GenValue()

	res := GenResponse {
		ID: uuid.NewString(),
		Val: value,
	}

	res.PutIntoDB()

	json.NewEncoder(w).Encode(res)
}

func (p *GenRequest) GenValue() string {
	var charset string
	switch p.Type {
	case "alpha":
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "alphanumeric":
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case "numeric":
		charset = "0123456789"
	case "guid":
		charset = "abcdef0123456789"
	default:
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	value := make([]byte, p.Length)
	for i := range value {
		value[i] = charset[rand.Intn(len(charset))]
	}

	return string(value)
}

func (p *GenResponse) PutIntoDB() {
	if err := database.AddItem(*p); err != nil {
		panic(err)
	}
}

func (p *GenResponse) GetFromDB() {
	if err := database.TakeItem(p.ID, p); err != nil {
		panic(err)
	}
}

func Retrieve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
	
	var res GenResponse
	res.ID = r.URL.Query().Get("id")
	res.GetFromDB()

	json.NewEncoder(w).Encode(res)
}