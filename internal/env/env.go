package env

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(fileName string) {
	_ = godotenv.Load(fileName)
}

func OpenAIAPIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func GroqAPIKey() string {
	return os.Getenv("GROQ_API_KEY")
}
