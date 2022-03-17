package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	s := NewStripper()
	for scanner.Scan() {
		line := scanner.Bytes()
		s(line)
		fmt.Println(string(line))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func NewStripper() func([]byte) {
	var inStr, inMLComment bool
	return func(jsonc []byte) {
		var inSLComment, escape, star, slash bool
		for i, c := range jsonc {
			if inStr {
				if !escape {
					if c == '"' {
						inStr = false
					} else if c == '\\' {
						escape = true
					}
				} else {
					escape = false
				}
			} else if inMLComment {
				if star {
					if c == '/' {
						inMLComment = false
					}
					star = false
				} else if c == '*' {
					star = true
				}
				jsonc[i] = ' '
			} else if inSLComment {
				jsonc[i] = ' '
			} else if slash {
				if c == '/' {
					inSLComment = true
					jsonc[i] = ' '
					jsonc[i-1] = ' '
				} else if c == '*' {
					inMLComment = true
					jsonc[i] = ' '
					jsonc[i-1] = ' '
				}
				slash = false
			} else if c == '/' {
				slash = true
			} else if c == '"' {
				inStr = true
			}
		}
	}
}

func StripJSONC(jsonc string) string {
	s := NewStripper()
	stripped := ""
	for _, lineStr := range strings.Split(jsonc, "\n") {
		line := []byte(lineStr)
		s(line)
		stripped += string(line) + "\n"
	}
	return stripped[:len(stripped)-1]
}
