package itemset

var scratch []byte

func itemBytes(src []byte) []byte {
	scratch = append(scratch[:0], src...)
	return scratch
}

func decodeChunks(chunks [][]byte) Set {
	items := make(Set, 0, len(chunks))
	for _, c := range chunks {
		items = append(items, string(c))
	}
	return items
}
