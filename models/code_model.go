package models


var AllowedLanguages map[string]bool = map[string]bool{
	"cpp":       true,
	"java":      true,
	"python":    true,
	"javascript": true,
	"go":        true,
}


type CodeData struct {
	Body string `json:"body" binding:"required"`
	Language   string `json:"language" binding:"required"`
}

