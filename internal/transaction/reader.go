package transaction

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Item = string

type Transaction = []Item

func ReadText(r io.Reader) ([]Transaction, error) {
	var txs []Transaction
	sc := bufio.NewScanner(r)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		items := strings.Fields(line)
		if len(items) == 0 {
			continue
		}
		txs = append(txs, items)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("transaction: read line %d: %w", lineNo, err)
	}
	return txs, nil
}

func ReadCSV(r io.Reader) ([]Transaction, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	var txs []Transaction
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("transaction: csv: %w", err)
		}
		var items []Item
		for _, f := range row {
			f = strings.TrimSpace(f)
			if f != "" {
				items = append(items, f)
			}
		}
		if len(items) > 0 {
			txs = append(txs, items)
		}
	}
	return txs, nil
}

func ReadJSON(r io.Reader) ([]Transaction, error) {
	var raw [][]string
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, fmt.Errorf("transaction: json: %w", err)
	}
	txs := make([]Transaction, 0, len(raw))
	for _, items := range raw {
		if len(items) > 0 {
			txs = append(txs, items)
		}
	}
	return txs, nil
}

func WriteText(w io.Writer, txs []Transaction) error {
	bw := bufio.NewWriter(w)
	for _, tx := range txs {
		if _, err := bw.WriteString(strings.Join(tx, " ")); err != nil {
			return err
		}
		if err := bw.WriteByte('\n'); err != nil {
			return err
		}
	}
	return bw.Flush()
}

func WriteJSON(w io.Writer, txs []Transaction) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(txs)
}
