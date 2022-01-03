package cron

import "strconv"

type tokenType int

const (
	number tokenType = iota
	slash
	dash
	comma
	star
)

// token lexed token
type token struct {
	ttype tokenType
	num   int
}

// tokenise turn a cron config element into a set of tokens.
// expecting a single times element.
func tokenise(text string) []token {
	tokens := []token{}

	for i := 0; i < len(text); i++ {
		switch c := text[i]; {
		case c == '*':
			tokens = append(tokens, token{ttype: star})
		case '0' <= c && c <= '9':
			start := i
			for ; i < len(text); i++ {
				c := text[i]
				if !('0' <= c && c <= '9') {
					break
				}
			}
			num := text[start:i]
			// FIXME: should deal with errors
			n, _ := strconv.Atoi(num)
			tokens = append(tokens, token{ttype: number, num: n})
			i--
		case c == '/':
			tokens = append(tokens, token{ttype: slash})
		case c == '-':
			tokens = append(tokens, token{ttype: dash})
		case c == ',':
			tokens = append(tokens, token{ttype: comma})
		}
	}

	return tokens
}
