package main

import (
	"fmt"
	"os"
	"strings"
)

func ascii(input, banner string) (string, error) {
	filename := "banners/" + banner
	content, err := os.ReadFile(filename)
	if err != nil{
		return "", nil
	}

	words := strings.Split(string(content), "\n")
	lines := strings.Split(input, "\\n")

	for _, char := range lines{
		if char == ""{
			continue
		}
		for row := 0; row < 9; row++{
			for _, c := range char{
				ascii := int(c)
				index:= (ascii - 32) * 9 + 1
				fmt.Print(words[index - row])
			}
		}
	}
	re
}