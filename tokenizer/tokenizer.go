package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/makeworld-the-better-one/go-isemoji"
)

type TOKEN_KIND = string

const (
	TOKEN_KIND_STRING_LIT           = "TOKEN_KIND_STRING_LIT"
	TOKEN_KIND_NUM_LIT              = "TOKEN_KIND_NUM_LIT"
	TOKEN_KIND_SLASH                = "TOKEN_KIND_SLASH"
	TOKEN_KIND_COMMA                = "TOKEN_KIND_COMMA"
	TOKEN_KIND_EQUAL                = "TOKEN_KIND_EQUAL"
	TOKEN_KIND_ID                   = "TOKEN_KIND_ID"
	TOKEN_KIND_TYPE                 = "TOKEN_KIND_TYPE"
	TOKEN_KIND_IF                   = "TOKEN_KIND_IF"
	TOKEN_KIND_ELSE                 = "TOKEN_KIND_ELSE"
	TOKEN_KIND_ELIF                 = "TOKEN_KIND_ELIF"
	TOKEN_KIND_FUNC                 = "TOKEN_KIND_FUNC"
	TOKEN_KIND_FOR                  = "TOKEN_KIND_FOR"
	TOKEN_KIND_NEW_LINE             = "TOKEN_KIND_NEW_LINE"
	TOKEN_KIND_BANG                 = "TOKEN_KIND_BANG"
	TOKEN_KIND_ARROW                = "TOKEN_KIND_ARROW"
	TOKEN_KIND_HASH                 = "TOKEN_KIND_HASH"
	TOKEN_KIND_TRUE                 = "TOKEN_KIND_TRUE"
	TOKEN_KIND_FALSE                = "TOKEN_KIND_FALSE"
	TOKEN_KIND_MINUS                = "TOKEN_KIND_MINUS"
	TOKEN_KIND_PLUS                 = "TOKEN_KIND_PLUS"
	TOKEN_KIND_STAR                 = "TOKEN_KIND_STAR"
	TOKEN_KIND_BAR                  = "TOKEN_KIND_BAR"
	TOKEN_KIND_AND                  = "TOKEN_KIND_AND"
	TOKEN_KIND_PERCENTAGE           = "TOKEN_KIND_PERCENTAGE"
	TOKEN_KIND_DB_PLUS              = "TOKEN_KIND_DB_PLUS"
	TOKEN_KIND_DB_MINUS             = "TOKEN_KIND_DB_MINUS"
	TOKEN_KIND_SEMI_COLON           = "TOKEN_KIND_SEMI_COLON"
	TOKEN_KIND_DB_COLON             = "TOKEN_KIND_DB_COLON"
	TOKEN_KIND_COLON                = "TOKEN_KIND_COLON"
	TOKEN_KIND_DB_EQUAL             = "TOKEN_KIND_DB_EQUAL"
	TOKEN_KIND_INFERIOR_SIGN        = "TOKEN_KIND_INFERIOR_SIGN"
	TOKEN_KIND_INFERIOR_OR_EQ_SIGN  = "TOKEN_KIND_INFERIOR_OR_EQ_SIGN"
	TOKEN_KIND_SUPERIOR_SIGN        = "TOKEN_KIND_SUPERIOR_SIGN"
	TOKEN_KIND_SUPERIOR_OR_EQ_SIGN  = "TOKEN_KIND_SUPERIOR_OR_EQ_SIGN"
	TOKEN_KIND_LEFT_BRACKET         = "TOKEN_KIND_LEFT_BRACKET"
	TOKEN_KIND_RIGHT_BRACKET        = "TOKEN_KIND_RIGHT_BRACKET"
	TOKEN_KIND_LEFT_SQUARE_BRACKET  = "TOKEN_KIND_LEFT_SQUARE_BRACKET"
	TOKEN_KIND_RIGHT_SQUARE_BRACKET = "TOKEN_KIND_RIGHT_SQUARE_BRACKET"
	TOKEN_KIND_LEFT_CURCLY          = "TOKEN_KIND_LEFT_CURCLY"
	TOKEN_KIND_RIGHT_CURLY          = "TOKEN_KIND_RIGHT_CURLY"
	TOKEN_KIND_PRINT                = "TOKEN_KIND_PRINT"
	TOKEN_KIND_RETURN               = "TOKEN_KIND_RETURN"
	TOKEN_KIND_LEN                  = "TOKEN_KIND_LEN"
	TOKEN_KIND_NEW                  = "TOKEN_KIND_NEW"
	TOKEN_KIND_DOT                  = "TOKEN_KIND_DOT"
)

type TokenLoc struct {
	Line  int
	Start int
	End   int
}

type Token struct {
	Kind    TOKEN_KIND
	Loc     TokenLoc
	Content string
}

type Tokenizer struct {
	chars  []string
	line   int
	cursor int
	col    int
}

func NewTokenizer(code string) Tokenizer {
	return Tokenizer{
		chars: strings.Split(code, ""),
		line:  1,
		col:   -1,
	}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	tokens := []Token{}
	for {
		c, ok := t.peek()
		if !ok {
			break
		}
		r, _ := utf8.DecodeRuneInString(c)
		switch true {
		case r == '\n':
			t.line++
			previous, ok := t.peekAt(-1)
			t.consume()
			if ok && previous != "\n" {
				tokens = append(tokens, Token{
					Kind: TOKEN_KIND_NEW_LINE,
					Loc: TokenLoc{
						Start: t.col,
						End:   t.col,
						Line:  t.line,
					},
					Content: c,
				})
			}
			t.col = 0
		// white space
		case unicode.IsSpace(r):
			t.consumeSpace()
		case r == '+':
			t.consume()
			start := t.col
			end := t.col
			next, ok := t.peek()
			kind := TOKEN_KIND_PLUS
			if ok && next == "+" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_DB_PLUS
			}
			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: start,
					End:   end,
					Line:  t.line,
				},
				Content: c,
			})
		case r == ':':
			t.consume()
			start := t.col
			end := t.col
			next, ok := t.peek()
			kind := TOKEN_KIND_COLON
			if ok && next == ":" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_DB_COLON
			}
			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: start,
					End:   end,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '=':
			t.consume()
			start := t.col
			end := t.col
			next, ok := t.peek()
			kind := TOKEN_KIND_EQUAL
			if ok && next == "=" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_DB_EQUAL
			}
			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: start,
					End:   end,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '*':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_STAR,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '<':
			t.consume()
			start := t.col
			end := t.col
			next, ok := t.peek()
			kind := TOKEN_KIND_INFERIOR_SIGN
			if ok && next == "=" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_INFERIOR_OR_EQ_SIGN
			}
			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: start,
					End:   end,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '>':
			t.consume()
			start := t.col
			end := t.col
			next, ok := t.peek()
			kind := TOKEN_KIND_SUPERIOR_SIGN
			if ok && next == "=" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_SUPERIOR_OR_EQ_SIGN
			}
			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: start,
					End:   end,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '/':
			t.consume()
			t2, _ := t.peek()
			if t2 == "/" {
				t.consume()
				for {
					t3, ok := t.peek()
					if !ok {
						break
					}
					t.consume()
					if t3 == "\n" {
						break
					}
				}
			} else {
				tokens = append(tokens, Token{
					Kind: TOKEN_KIND_SLASH,
					Loc: TokenLoc{
						Start: t.col,
						End:   t.col,
						Line:  t.line,
					},
					Content: c,
				})
			}
		case r == '!':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_BANG,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '#':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_HASH,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '(':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_LEFT_BRACKET,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '|':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_BAR,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '&':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_AND,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '{':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_LEFT_CURCLY,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '}':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_RIGHT_CURLY,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == ')':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_RIGHT_BRACKET,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == ']':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_RIGHT_SQUARE_BRACKET,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '%':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_PERCENTAGE,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '[':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_LEFT_SQUARE_BRACKET,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '-':
			t.consume()
			start := t.col
			end := t.col
			next, ok := t.peek()
			kind := TOKEN_KIND_MINUS
			if ok && next == "-" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_DB_MINUS
			}
			if ok && next == ">" {
				t.consume()
				end = t.col
				c += next
				kind = TOKEN_KIND_ARROW
			}
			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: start,
					End:   end,
					Line:  t.line,
				},
				Content: c,
			})
		case r == ',':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_COMMA,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == ';':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_SEMI_COLON,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case r == '.':
			t.consume()
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_DOT,
				Loc: TokenLoc{
					Start: t.col,
					End:   t.col,
					Line:  t.line,
				},
				Content: c,
			})
		case isQuote(r):
			string := c
			t.consume()
			tokStart := t.col
			opener := r
			for {
				stop := false
				c, ok = t.peek()
				if !ok {
					break
				}
				r, _ = utf8.DecodeRuneInString(c)
				if r == opener {
					previous, _ := t.peekAt(-1)
					if previous != "\\" {
						stop = true
					}
				}

				string += c
				t.consume()
				if stop {
					break
				}

			}
			tokEnd := t.col
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_STRING_LIT,
				Loc: TokenLoc{
					Start: tokStart,
					End:   tokEnd,
					Line:  t.line,
				},
				Content: string,
			})
		case unicode.IsNumber(r):
			number := c
			t.consume()
			hasDecimal := false
			tokStart := t.col
			for {
				c, ok = t.peek()
				if !ok {
					break
				}
				r, _ = utf8.DecodeRuneInString(c)
				isDecimalPoint := len(c) > 0 && !hasDecimal && r == '.'
				if !unicode.IsNumber(r) && !isDecimalPoint {
					break
				}
				if isDecimalPoint {
					hasDecimal = true
				}
				number += c
				t.consume()
			}
			tokEnd := t.col
			tokens = append(tokens, Token{
				Kind: TOKEN_KIND_NUM_LIT,
				Loc: TokenLoc{
					Start: tokStart,
					End:   tokEnd,
					Line:  t.line,
				},
				Content: number,
			})
		case unicode.IsLetter(r) || r == '_' || isemoji.IsEmoji(c):
			id := c
			t.consume()
			tokStart := t.col
			for {
				c, ok = t.peek()
				if !ok {
					break
				}
				r, _ = utf8.DecodeRuneInString(c)
				if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && !isemoji.IsEmoji(c) {
					break
				}
				id += c
				t.consume()
			}
			tokEnd := t.col
			kind := TOKEN_KIND_ID
			if strings.ToLower(id) == "true" {
				kind = TOKEN_KIND_TRUE
			}
			if strings.ToLower(id) == "false" {
				kind = TOKEN_KIND_FALSE
			}
			if strings.ToLower(id) == "if" {
				kind = TOKEN_KIND_IF
			}
			if strings.ToLower(id) == "type" {
				kind = TOKEN_KIND_TYPE
			}
			if strings.ToLower(id) == "else" {
				kind = TOKEN_KIND_ELSE
			}
			if strings.ToLower(id) == "return" {
				kind = TOKEN_KIND_RETURN
			}
			if strings.ToLower(id) == "func" {
				kind = TOKEN_KIND_FUNC
			}
			if strings.ToLower(id) == "elif" {
				kind = TOKEN_KIND_ELIF
			}
			if strings.ToLower(id) == "len" {
				kind = TOKEN_KIND_LEN
			}
			if strings.ToLower(id) == "for" {
				kind = TOKEN_KIND_FOR
			}
			if strings.ToLower(id) == "new" {
				kind = TOKEN_KIND_NEW
			}
			if strings.ToLower(id) == "print" {
				kind = TOKEN_KIND_PRINT
			}

			tokens = append(tokens, Token{
				Kind: kind,
				Loc: TokenLoc{
					Start: tokStart,
					End:   tokEnd,
					Line:  t.line,
				},
				Content: id,
			})
		default:
			return tokens, fmt.Errorf("unexpected token %s at position %d:%d", c, t.line, t.cursor)
		}
	}

	return tokens, nil
}

func isQuote(str rune) bool {
	return str == '\'' || str == '"'
}

func (t *Tokenizer) peek() (string, bool) {
	if t.cursor >= len(t.chars) {
		return "", false
	}
	return t.chars[t.cursor], true
}
func (t *Tokenizer) peekAt(pos int) (string, bool) {
	if (t.cursor+pos < 0) || (t.cursor+pos) >= len(t.chars) {
		return "", false
	}
	return t.chars[t.cursor+pos], true
}

func (t *Tokenizer) consume() {
	t.cursor++
	t.col++
}

func (t *Tokenizer) consumeSpace() {
	t.cursor++
}

func Tokenize(code string) ([]Token, error) {
	t := NewTokenizer(code)
	return t.Tokenize()
}
