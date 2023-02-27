# Fibonacci Hashmap in Go

I started this project to learn about the inner workings of hashmaps as well as play around with Go's new generics.

fmap uses the Fibonacci indexing method from Donald Knuth's The Art of Computer Programming Vol.3.

The map is fully functioning, although not very efficient, with the exception of deleting entries.

The caller must provide a hash function that takes whatever concrete type the key is and returns the hash as an **uint64**

Below is an example.

```
// import "hash/fnv"
hasher := func(k string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(k))
	return h.Sum64()
}

hm := fmap.New[string, string](hasher)
hm.Put("This is a unique key", "This is the value for the key")
```