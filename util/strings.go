package util

// Copied from github.com/pressly/goose

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

type camelSnakeStateMachine int

const ( //                                           _$$_This is some text, OK?!
	idle          camelSnakeStateMachine = iota // 0 ↑                     ↑   ↑
	firstAlphaNum                               // 1     ↑    ↑  ↑    ↑     ↑
	alphaNum                                    // 2      ↑↑↑  ↑  ↑↑↑  ↑↑↑   ↑
	delimiter                                   // 3         ↑  ↑    ↑    ↑   ↑
)

func (s camelSnakeStateMachine) next(r rune) camelSnakeStateMachine {
	switch s {
	case idle:
		if IsAlphaNumeric(r) {
			return firstAlphaNum
		}
	case firstAlphaNum:
		if IsAlphaNumeric(r) {
			return alphaNum
		}
		return delimiter
	case alphaNum:
		if !IsAlphaNumeric(r) {
			return delimiter
		}
	case delimiter:
		if IsAlphaNumeric(r) {
			return firstAlphaNum
		}
		return idle
	}
	return s
}

func CamelCase(str string) string {
	var b strings.Builder

	stateMachine := idle
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		stateMachine = stateMachine.next(r)
		switch stateMachine {
		case firstAlphaNum:
			b.WriteRune(unicode.ToUpper(r))
		case alphaNum:
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func SnakeCase(str string) string {
	var b bytes.Buffer

	stateMachine := idle
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		stateMachine = stateMachine.next(r)
		switch stateMachine {
		case firstAlphaNum, alphaNum:
			b.WriteRune(unicode.ToLower(r))
		case delimiter:
			b.WriteByte('_')
		}
	}
	if stateMachine == idle {
		return string(bytes.TrimSuffix(b.Bytes(), []byte{'_'}))
	}
	return b.String()
}

func IsAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}
