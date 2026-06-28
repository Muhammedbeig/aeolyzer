package middleware

func ResolveContext(raw map[string]interface{}) map[string]string {
	return ExtractSanitizedContext(raw)
}

func GetContextFlag(context map[string]string, key string) bool {
	val, ok := context[key]
	return ok && val == "true"
}
