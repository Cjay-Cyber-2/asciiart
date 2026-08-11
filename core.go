package handler

import (
	"os"
	"strings"
)

func ascii(input, banner string) (string, error) {
	filename := "banners/" + banner + ".txt"
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	words := strings.Split(string(content), "\n")
	lines := strings.Split(input, "\\n")

	var result strings.Builder

	for _, char := range lines{
		if char == ""{
            result.WriteString("\n")
			continue
		}
		for row := 0; row < 9; row++{
			for _, c := range char{
				ascii := int(c)
				if ascii >= 32 && ascii <= 126 {
					index := (ascii - 32) * 9 + 1
					result.WriteString(words[index + row])
				}
			}
            if row < 8 {
			    result.WriteString("\n")
            }
		}
	}
	return result.String(), nil
}