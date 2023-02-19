package fmap

type HashMap[K comparable, V any] struct {
	hasher func(K) uint64
	buckets []*bucket[K, V]
	buckets_size int
	map_size int
	collisions int
	total_cost int
	longest_probe uint64
}

func (hm *HashMap[K, V]) GetCollisions() int {
	return hm.collisions
}

func (hm *HashMap[K, V]) GetMapSize() int {
	return hm.map_size
}

func (hm *HashMap[K, V]) GetLongestProbe() uint64 {
	return hm.longest_probe
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

	//double the map size if it is 85% full or greater
	if float64(hm.map_size) / float64(hm.buckets_size) > 0.85 {
		hm.grow(hm.buckets_size * 2)
	}

	var collisions, probeposition = insert(&bucket[K, V]{key: key, value: value, hash: hash}, hm.buckets, hm.buckets_size)

	hm.collisions += collisions

	if probeposition > hm.longest_probe {
		hm.longest_probe = probeposition
	}

	hm.map_size += 1
}

func insert[K comparable, V any](b *bucket[K, V], buckets []*bucket[K, V], length int) (int, uint64) {
	index := modulo_index(b.hash, uint64(length))
	
	var probeposition uint64 = 0
	var collisions int = 0;
	for buckets[index] != nil {
		collisions += 1
		probeposition += 1
		index = modulo_index(b.hash + probeposition, uint64(length))
	}

	buckets[index] = b

	return collisions, probeposition
}

func (hm *HashMap[K, V]) grow(new_size int) {
	var new_buckets = make([]*bucket[K, V], new_size)

	var new_longest_probe uint64 = 0
	var new_collision_count int = 0
	for _, v := range hm.buckets {
		if v == nil {
			continue
		}

		var collisions, probeposition = insert(v, new_buckets, new_size)
		new_collision_count += collisions

		if probeposition > new_longest_probe {
			new_longest_probe = probeposition
		}
	}

	hm.longest_probe = new_longest_probe
	hm.collisions = new_collision_count
	hm.buckets = new_buckets
	hm.buckets_size = new_size
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