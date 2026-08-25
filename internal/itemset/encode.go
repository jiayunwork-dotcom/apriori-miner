package itemset

import (
	"fmt"
	"strings"
)

func Key(s Set) string {
	return strings.Join(s, "\x00")
}

func FromKey(k string) Set {
	if k == "" {
		return nil
	}
	return strings.Split(k, "\x00")
}

func Display(s Set) string {
	if len(s) == 0 {
		return "{}"
	}
	return "{" + strings.Join(s, ", ") + "}"
}

func DisplayRule(antecedent, consequent Set) string {
	return Display(antecedent) + " => " + Display(consequent)
}

func ParseDisplay(s string) Set {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var items []Item
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			items = append(items, p)
		}
	}
	return New(items...)
}

func Encode(s Set) []byte {
	var buf []byte
	for _, item := range s {
		b := []byte(item)
		if len(b) > 255 {
			b = b[:255]
		}
		buf = append(buf, byte(len(b)))
		buf = append(buf, b...)
	}
	return buf
}

func Decode(data []byte) (Set, error) {
	var items []Item
	i := 0
	for i < len(data) {
		if i >= len(data) {
			return nil, fmt.Errorf("truncated at offset %d", i)
		}
		length := int(data[i])
		i++
		if i+length > len(data) {
			return nil, fmt.Errorf("item at offset %d extends past end", i)
		}
		items = append(items, string(data[i:i+length]))
		i += length
	}
	return items, nil
}
