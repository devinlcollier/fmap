package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"log"
	"encoding/csv"
	"strconv"
)

type HashMap[K comparable, V any] struct {
	hasher func(K) uint64
	buckets []*bucket[K, V]
	buckets_size int
	map_size int
	collisions int
	total_cost int
	longest_probe uint64
}

func New[K comparable, V any](hasher_func func(K) uint64) *HashMap[K, V] {
	ret_map := new(HashMap[K, V])
	ret_map.hasher = hasher_func
	ret_map.buckets_size = 1000
	ret_map.map_size = 0
	ret_map.longest_probe = 0
	ret_map.buckets = make([]*bucket[K, V], ret_map.buckets_size)

	return ret_map
}

func (hm *HashMap[K, V]) Put(key K, value V) {
	hash := hm.hasher(key)
	//index := fibonacci_index(hash, 55)
	hm.insert(&bucket[K, V]{key: key, value: value, hash: hash})
}

func (hm *HashMap[K, V]) insert(b *bucket[K, V]) {
	index := modulo_index(b.hash, uint64(len(hm.buckets)))
	
	var probeposition uint64 = 0
	for hm.buckets[index] != nil {
		hm.collisions += 1
		probeposition += 1
		index = modulo_index(b.hash + probeposition, uint64(len(hm.buckets)))
	}

	if probeposition > hm.longest_probe {
		hm.longest_probe = probeposition
	}

	fmt.Println("key:", b.key, index)

	hm.buckets[index] = b
	hm.map_size += 1
}

func (hm *HashMap[K, V]) Get(key K) V {
	hash := hm.hasher(key)
	index := hm.search(hash)
	return hm.buckets[index].value
}

func (hm *HashMap[K, V]) search(hash uint64) uint64 {
	index := modulo_index(hash, uint64(len(hm.buckets)))

	var probeposition uint64 = 0

	for probeposition <= hm.longest_probe && hm.buckets[index] != nil {
		if hm.buckets[index].hash == hash {
			return index
		}
		probeposition += 1
		index = modulo_index(hash + probeposition, uint64(len(hm.buckets)))
	}

	return 0
}

type bucket[K comparable, V any] struct {
	hash uint64
	key K
	value V
}

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
	for i := 0; i < 400; i++ {
		//fmt.Println(unique_words[i], ",", unique_word_map[unique_words[i]])
		hm.Put(unique_words[i], unique_word_map[unique_words[i]])
	}

	for i := 0; i < 10; i++ {
		fmt.Println("hm.Get(unique_words[" + strconv.Itoa(i) + "])",unique_words[i] ,hm.Get(unique_words[i]))
	}

	fmt.Println("collisions", hm.collisions)
	fmt.Println("map size", hm.map_size)
	fmt.Println("longest_probe", hm.longest_probe)
}