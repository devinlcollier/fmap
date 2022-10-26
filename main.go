package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"log"
	"encoding/csv"
)

// type Keyable interface {
// 	hash.Hash64
// }

// type HashTable[K Keyable, V any] struct {
// 	arraysize int
// 	buckets [12]*bucket[K, V]
// }

// type bucket[K Keyable, V any] struct {
// 	head *bucketnode[K, V]
// }

// type bucketnode[K Keyable, V any] struct {
// 	key   K
// 	value V
// 	next  *bucketnode[K, V]
// }

// func Init[K Keyable, V any]() *HashTable[K, V] {
// 	ret_table := &HashTable[K, V]{arraysize: 12}

// 	for i := range ret_table.buckets {
// 		ret_table.buckets[i] = &bucket[K, V]{}
// 	}
// 	return ret_table
// }

type HashMap[K comparable, V any] struct {
	hasher func(K) uint64
	buckets []*bucket[K, V]
	collisions int
}

func New[K comparable, V any](hasher_func func(K) uint64) *HashMap[K, V] {
	ret_map := new(HashMap[K, V])
	ret_map.hasher = hasher_func
	ret_map.buckets = make([]*bucket[K, V], 1000)

	return ret_map
}

func (hm *HashMap[K, V]) Put(key K, value V) {
	hash := hm.hasher(key)
	//index := fibonacci_index(hash, 55)
	index := modulo_index(hash, uint64(len(hm.buckets)))
	fmt.Println("key:", key, index)
	if hm.buckets[index] != nil {
		fmt.Println("collision, attempt:", key, " existing:", hm.buckets[index].key)
		hm.collisions++
	}
	hm.buckets[index] = &bucket[K, V]{key: key, value: value, hash: hash}
}

func (hm *HashMap[K, V]) Get(key K) V {
	hash := hm.hasher(key)
	//index := fibonacci_index(hash, 57)
	index := modulo_index(hash, uint64(len(hm.buckets)))
	fmt.Println("key:", key, index)
	return hm.buckets[index].value
}

type bucket[K comparable, V any] struct {
	hash uint64
	key K
	value V
}

// func (hashtable *HashTable[K, V]) Insert(key K, value V) {
// 	index := int(fibonacci_hash(key.Sum64(), 60))
// 	hashtable.buckets[index].insert(key, value)
// }

// func (hashtable *HashTable[K, V]) Search(key K) bool {
// 	index := int(fibonacci_hash(key.Sum64(), 60))
// 	return hashtable.buckets[index].search(key)
// }

// func (hashtable *HashTable[K, V]) Delete(key K) {
// 	index := int(fibonacci_hash(key.Sum64(), 60))
// 	hashtable.buckets[index].delete(key)
// }

// func (b *bucket[K, V]) insert(key K, value V) {
// 	newNode := &bucketnode[K, V]{key: key, value: value}
// 	newNode.next = b.head
// 	b.head = newNode
// }

// func (b *bucket[K, V]) search(key K) bool {
// 	currentNode := b.head
// 	for currentNode != nil {
// 		if currentNode.key == key {
// 			return true
// 		}
// 		currentNode = currentNode.next
// 	}
// 	return false
// }

// func (b *bucket[K, V]) delete(key K) {
// 	if b.head.key == key {
// 		b.head = b.head.next
// 		return
// 	}

// 	previousNode := b.head
// 	for previousNode.next != nil {
// 		if previousNode.next.key == key {
// 			previousNode.next = previousNode.next.next
// 		}
// 		previousNode = previousNode.next
// 	}
// }

func fibonacci_index(hash uint64, bits int) uint64 {
	return (hash * 11400714819323198485) >> bits
}

func modulo_index(hash uint64, size uint64) uint64 {
	return hash % size;
}

func main() {
	dictionary, err := os.Open("dictionary.csv")
	if err != nil {
		log.Fatal(err)
	}

	csv_reader := csv.NewReader(dictionary)

	lines, ln_err := csv_reader.ReadAll()
	if ln_err != nil {
		log.Fatal(ln_err)
	}

	unique_words := make([]string, len(lines))
	unique_word_map := make(map[string]string)

	for i := range lines {
		unique_word_map[lines[i][0]] = lines[i][2]
	}

	i := 0
	for k, _ := range unique_word_map {
		unique_words[i] = k
		i++
	}

	/*hasher := func(k string) uint64 {
		h := fnv.New64a()
		h.Write([]byte(k))
		return h.Sum64()
	}*/

	carter_wegman_hasher := func(k string) uint64 {
		var ret uint64 = 0
		for _, e := range k {
			//fmt.Println("e", e)
			h := fnv.New64a()
			h.Write([]byte(string(e)))
			ret = ret ^ h.Sum64()
		}
		return ret
	}

	hm := New[string, string](carter_wegman_hasher)
	for i := 0; i < 90; i++ {
		//fmt.Println(unique_words[i], ",", unique_word_map[unique_words[i]])
		hm.Put(unique_words[i], unique_word_map[unique_words[i]])
	}
	fmt.Println(hm.buckets)
	fmt.Println("collisions", hm.collisions)
}