package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"apriori-miner/internal/apriori"
)

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	addr := fs.String("addr", ":8080", "listen address")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/mine", handleMine)

	fmt.Printf("apriori-miner serving on %s\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}
}

func handleMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Transactions [][]string `json:"transactions"`
		MinSupport   float64    `json:"min_support"`
		MinConfidence float64   `json:"min_confidence"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Transactions) == 0 {
		http.Error(w, "transactions required", http.StatusBadRequest)
		return
	}
	if req.MinSupport <= 0 {
		req.MinSupport = 0.2
	}
	if req.MinConfidence <= 0 {
		req.MinConfidence = 0.6
	}

	txs := make([]apriori.Transaction, len(req.Transactions))
	for i, items := range req.Transactions {
		tx := make(apriori.Transaction, len(items))
		for j, it := range items {
			tx[j] = apriori.Item(it)
		}
		txs[i] = tx
	}

	freq := apriori.Apriori(txs, req.MinSupport)
	rules := apriori.GenerateRules(freq, len(txs), req.MinConfidence)

	type ruleOut struct {
		Antecedent []string `json:"antecedent"`
		Consequent []string `json:"consequent"`
		Support    float64  `json:"support"`
		Confidence float64  `json:"confidence"`
	}
	out := make([]ruleOut, len(rules))
	for i, rl := range rules {
		ant := make([]string, len(rl.Antecedent))
		for j, it := range rl.Antecedent {
			ant[j] = string(it)
		}
		con := make([]string, len(rl.Consequent))
		for j, it := range rl.Consequent {
			con[j] = string(it)
		}
		out[i] = ruleOut{Antecedent: ant, Consequent: con, Support: rl.Support, Confidence: rl.Confidence}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"frequent_sets": len(freq),
		"rules":         out,
	})
}
