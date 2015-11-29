// © 2015 Steve McCoy. See LICENSE for details.

// The formate program formats text into comfortable line lengths.
// Blank lines and lines beginning with non-letter characters are treated literally.
// All other lines are combined or split in order to fit the lines within the minimum
// and maximum lengths (45 and 75 by default).
//
// The input text is expected to be in UTF-8 or a subset.
// Lines beginning with a non-UTF-8 byte sequence will be treated literally.
// Lines containing a non-UTF-8 byte sequence may be combined in ugly ways.
package main

import (
	"bufio"
	"bytes"
	"os"
	"unicode"
	"unicode/utf8"
)

var minLen = 45
var maxLen = 75

func main() {
	r := bufio.NewScanner(os.Stdin)
	for {
		para, more := scanPara(r)
		for i := 0; i < len(para); i++ {
			line := para[i]

			if isLiteral(line) {
				os.Stdout.Write(line)
				os.Stdout.Write([]byte{'\n'})
				continue
			}

			n := utf8.RuneCount(line)
			for n < minLen {
				if i+1 == len(para) || isLiteral(para[i+1]) {
					// nothing to join with
					break
				}

				if line[len(line)-1] != ' ' {
					line = append(line, ' ')
				}
				line = append(line, para[i+1]...)
				i++
				n = utf8.RuneCount(line)
			}

			if n > maxLen {
				var rs []rune
				for _, r := range string(line) {
					rs = append(rs, r)
				}
				i := maxLen
				for ; i >= 0; i-- {
					if rs[i] == ' ' {
						break
					}
				}
				first := encodeRunes(rs[:i])
				rest := encodeRunes(rs[i+1:])
				line = first
				if i+1 < len(para) {
					if isLiteral(para[i+1]) {
						// next line is literal, so insert rest before it
						para = append(para[i+1:], append(para[:i+1], rest)...)
					} else {
						if rest[len(rest)-1] != ' ' {
							rest = append(rest, ' ')
						}
						para[i+1] = append(rest, para[i+1]...)
					}
				} else {
					para = append(para, rest)
				}
			}

			os.Stdout.Write(line)
			os.Stdout.Write([]byte{'\n'})
		}

		if !more {
			break
		}

		os.Stdout.Write([]byte{'\n'})
	}
	if err := r.Err(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
}

func scanPara(r *bufio.Scanner) ([][]byte, bool) {
	var para [][]byte
	for r.Scan() {
		line := r.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			return para, true
		}
		para = append(para, append([]byte(nil), line...))
	}
	return para, false
}

func isLiteral(line []byte) bool {
	first, _ := utf8.DecodeRune(line)
	return first == utf8.RuneError || !unicode.IsLetter(first)
}

func encodeRunes(rs []rune) []byte {
	n := 0
	for _, r := range rs {
		n += utf8.RuneLen(r)
	}
	bs := make([]byte, n)
	i := 0
	for _, r := range rs {
		i += utf8.EncodeRune(bs[i:], r)
	}
	return bs
}
