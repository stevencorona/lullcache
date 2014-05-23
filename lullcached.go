package main

func main() {
	cache := NewCacheServer("127.0.0.1:11211")
	cache.Start()
}
