package server

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"github.com/google/uuid"

	"github.com/Mart1n-Eden/randomgen/internal/database"
)

type GenRequest struct {
	Type   string `json:"type"`
	Length int    `json:"length"`
}

type GenResponse struct {
	ID  string `json:"id"`
	Val string `json:"value"`
}

func Run() {
	database.Init(GenResponse{})

	http.HandleFunc("/api/generate", Generate)
	http.HandleFunc("/api/retrieve", Retrieve)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit
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

	value, err := req.GenValue()
	if err != nil {
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}

	res := GenResponse{
		ID:  uuid.NewString(),
		Val: value,
	}

	res.PutIntoDB()

	json.NewEncoder(w).Encode(res)
}

func (p *GenRequest) GenValue() (string, error) {
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
		// charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		return string(""), errors.New("Invalid type")
	}

	value := make([]byte, p.Length)
	for i := range value {
		value[i] = charset[rand.Intn(len(charset))]
	}

	return string(value), nil
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
