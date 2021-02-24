package cache

const (
	ErrorGetRedisClient    = "Get redis client error"
	ErrorCheckRedisHealth  = "Check redis health error"
	ErrorGenerateRedisCr   = "Generate redis cr error"
	ErrorSetOwnerReference = "Set owner reference error"
	// nolint:gosec
	ErrorCreateRedisSecret            = "Create redis secret error"
	ErrorCreateRedisCr                = "Create redis cr error"
	ErrorDefaultUnstructuredConverter = "Default unstructured converter error"
)
