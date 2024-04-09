package llm

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

func NumberOfTokens(text string) (int, error) {
	tkm, err := tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		err = fmt.Errorf("error getting encoding for model: %v", err)
		return 0, err
	}
	return len(tkm.Encode(text, nil, nil)), nil
}

func TrimNonJSON(s string) string {
	startIndex := strings.IndexAny(s, "{[")
	endIndex := strings.LastIndexAny(s, "]}")

	if startIndex == -1 || endIndex == -1 {
		return s // Return original string if no JSON boundaries are found
	}

	return s[startIndex : endIndex+1]
}
