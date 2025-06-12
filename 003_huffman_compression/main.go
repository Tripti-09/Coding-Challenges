package main

import (
	"container/heap"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"strings"
)

// -------------------------
// Huffman Tree Node Struct
// -------------------------
type Node struct {
	Char   byte
	Freq   int
	Left   *Node
	Right  *Node
	IsLeaf bool
}

// -------------------------
// Priority Queue (Min Heap)
// -------------------------
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Freq < pq[j].Freq }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*Node)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// -------------------------
// Build Huffman Tree
// -------------------------
func BuildHuffmanTree(freqMap map[byte]int) *Node {
	pq := &PriorityQueue{}
	heap.Init(pq)

	for ch, freq := range freqMap {
		heap.Push(pq, &Node{Char: ch, Freq: freq, IsLeaf: true})
	}

	for pq.Len() > 1 {
		left := heap.Pop(pq).(*Node)
		right := heap.Pop(pq).(*Node)

		merged := &Node{
			Freq:   left.Freq + right.Freq,
			Left:   left,
			Right:  right,
			IsLeaf: false,
		}

		heap.Push(pq, merged)
	}

	return heap.Pop(pq).(*Node)
}

// -------------------------
// Generate Prefix Codes
// -------------------------
func GenerateCodes(node *Node, prefix string, table map[byte]string) {
	if node == nil {
		return
	}
	if node.IsLeaf {
		table[node.Char] = prefix
		return
	}
	GenerateCodes(node.Left, prefix+"0", table)
	GenerateCodes(node.Right, prefix+"1", table)
}

// -------------------------
// Encode Text Using Code Table
// -------------------------
func EncodeText(data []byte, codeTable map[byte]string) []byte {
	var builder strings.Builder

	for _, b := range data {
		builder.WriteString(codeTable[b])
	}
	bitString := builder.String()

	var encoded []byte
	for i := 0; i < len(bitString); i += 8 {
		byteStr := bitString[i:]
		if len(byteStr) > 8 {
			byteStr = byteStr[:8]
		}
		var b byte
		for j := 0; j < len(byteStr); j++ {
			b <<= 1
			if byteStr[j] == '1' {
				b |= 1
			}
		}
		encoded = append(encoded, b)
	}
	return encoded
}

// -------------------------
// Decode Text from Huffman Tree
// -------------------------
func DecodeText(encoded []byte, root *Node) []byte {
	var result []byte
	current := root

	for _, b := range encoded {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			if bit == 0 {
				current = current.Left
			} else {
				current = current.Right
			}
			if current.IsLeaf {
				result = append(result, current.Char)
				current = root
			}
		}
	}
	return result
}

// -------------------------
// Save to File (.huff)
// -------------------------
func SaveToFile(filename string, freq map[byte]int, encoded []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(freq)
	if err != nil {
		return err
	}

	_, err = file.Write(encoded)
	return err
}

// -------------------------
// Load from File (.huff)
// -------------------------
func LoadFromFile(filename string) (map[byte]int, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var freq map[byte]int
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&freq)
	if err != nil {
		return nil, nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	return freq, data, nil
}

// -------------------------
// MAIN — COMPRESS + DECOMPRESS
// -------------------------
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  Compress:   go run main.go -c <input.txt>")
		fmt.Println("  Decompress: go run main.go -d <input.huff>")
		return
	}

	mode := os.Args[1]
	filename := os.Args[2]

	fmt.Println("Reading file...")

	if mode == "-c" {
		// COMPRESSION
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		fmt.Println("Counting frequencies...")
		// Step 1: Frequency count
		freq := make(map[byte]int)
		for _, b := range data {
			freq[b]++
		}

		fmt.Println("Building tree...")
		// Step 2: Build Tree
		root := BuildHuffmanTree(freq)

		fmt.Println("Generating codes...")
		// Step 3: Generate codes
		codeTable := make(map[byte]string)
		GenerateCodes(root, "", codeTable)

		
		fmt.Println("Encoding...")
		// Step 4: Encode
		encoded := EncodeText(data, codeTable)
		fmt.Println("Encoded...")

		// Step 5: Save
		err = SaveToFile("output.huff", freq, encoded)
		if err != nil {
			fmt.Println("Error saving compressed file:", err)
			return
		}
		fmt.Println("✅ Compression complete: output.huff")

	} else if mode == "-d" {
		// DECOMPRESSION
		freq, encoded, err := LoadFromFile(filename)
		if err != nil {
			fmt.Println("Error loading file:", err)
			return
		}

		// Step 6: Rebuild tree
		root := BuildHuffmanTree(freq)

		// Step 7: Decode
		decoded := DecodeText(encoded, root)

		// Step 8: Write decoded text
		err = os.WriteFile("output_decoded.txt", decoded, 0644)
		if err != nil {
			fmt.Println("Error writing decoded file:", err)
			return
		}
		fmt.Println("✅ Decompression complete: output_decoded.txt")
	} else {
		fmt.Println("Invalid mode. Use -c to compress or -d to decompress.")
	}
}
