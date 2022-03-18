package stripjsonc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func StripJSONCStream(in, out *os.File) {
	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	strip := NewStripper()
	for scanner.Scan() {
		line := scanner.Bytes()
		strip(line)
		newlined := append(line, '\n')
		if _, err := writer.Write(newlined); err != nil {
			fmt.Fprintf(os.Stderr, "error writing to %s: %e\n", out.Name(), err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading from %s: %e\n", in.Name(), err)
	}
	if err := writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "error flushing %s: %e\n", in.Name(), err)
	}
}

func StripJSONCString(jsonc string) string {
	s := NewStripper()
	stripped, delim := "", ""
	for _, lineStr := range strings.Split(jsonc, "\n") {
		line := []byte(lineStr)
		s(line)
		stripped += delim + string(line)
		delim = "\n"
	}
	return stripped
}
