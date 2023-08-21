package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := BuscarCotacao()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao.Bid)

}

func BuscarCotacao() (*Cotacao, error) {
	var URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*200)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		log.Println(error.Error(err))
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(error.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(error.Error(err))
		return nil, err
	}

	cotacao := JsonToMap(body)

	var c Cotacao
	jsonStr, _ := json.Marshal(cotacao["USDBRL"])
	err = json.Unmarshal(jsonStr, &c)
	if err != nil {
		return nil, err
	}
	SaveCotacao(ctx, &c)
	return &c, nil
}

func SaveCotacao(ctx context.Context, c *Cotacao) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	db, err := gorm.Open(sqlite.Open("usdbrl.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Cotacao{})

	db.WithContext(ctx).Create(&c)
	return nil
}

func JsonToMap(body []byte) map[string]interface{} {
	var cotacao map[string]interface{}
	json.Unmarshal(body, &cotacao)
	return cotacao
}
