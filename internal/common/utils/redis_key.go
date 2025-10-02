package utils

const (
	Prefix = "saas:"
)

func GetRedisKey(key string) string {
	return Prefix + key
}
