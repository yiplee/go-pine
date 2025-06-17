package pine

var cache map[string]ValueSeries = make(map[string]ValueSeries)

func getCache(key string) ValueSeries {
	return cache[key]
}

func setCache(key string, v ValueSeries) {
	cache[key] = v
}
