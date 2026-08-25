package transaction

import (
	"fmt"
	"sort"
	"strings"
)

type Vocabulary struct {
	itemToID map[Item]int
	idToItem []Item
}

func NewVocabulary(txs []Transaction) *Vocabulary {
	unique := make(map[Item]struct{})
	for _, tx := range txs {
		for _, it := range tx {
			unique[it] = struct{}{}
		}
	}
	items := make([]Item, 0, len(unique))
	for it := range unique {
		items = append(items, it)
	}
	sort.Strings(items)

	v := &Vocabulary{
		itemToID: make(map[Item]int, len(items)),
		idToItem: items,
	}
	for i, it := range items {
		v.itemToID[it] = i
	}
	return v
}

func (v *Vocabulary) Size() int { return len(v.idToItem) }

func (v *Vocabulary) Encode(item Item) int {
	id, ok := v.itemToID[item]
	if !ok {
		return -1
	}
	return id
}

func (v *Vocabulary) Decode(id int) (Item, error) {
	if id < 0 || id >= len(v.idToItem) {
		return "", fmt.Errorf("id %d out of range [0, %d)", id, len(v.idToItem))
	}
	return v.idToItem[id], nil
}

func (v *Vocabulary) EncodeTransaction(tx Transaction) []int {
	ids := make([]int, 0, len(tx))
	for _, it := range tx {
		if id, ok := v.itemToID[it]; ok {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)
	return ids
}

func (v *Vocabulary) DecodeTransaction(ids []int) Transaction {
	tx := make(Transaction, 0, len(ids))
	for _, id := range ids {
		if id >= 0 && id < len(v.idToItem) {
			tx = append(tx, v.idToItem[id])
		}
	}
	return tx
}

func (v *Vocabulary) EncodeAll(txs []Transaction) [][]int {
	encoded := make([][]int, len(txs))
	for i, tx := range txs {
		encoded[i] = v.EncodeTransaction(tx)
	}
	return encoded
}

func (v *Vocabulary) String() string {
	var b strings.Builder
	for id, item := range v.idToItem {
		fmt.Fprintf(&b, "%d: %s\n", id, item)
	}
	return b.String()
}
