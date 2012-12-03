package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ItemType int

const eof = rune(-1)

const (
	ItemError          ItemType = iota
	ItemEOF                     //  1: EOF
	ItemHtmlText                //  2: html text
	ItemHtmlLiteral             //  2: html literal ends with >, =" or ={{
	ItemHtmlComment             //  3: html comment
	ItemHtmlAttr                //  4: ={{
	ItemHtmlAttrInterp          //  5: ="
	ItemLeftMeta3               //  5: {{{
	ItemLeftMetaPound           //  6: {{#
	ItemLeftMetaSlash           //  7: {{/
	ItemLeftMetaBang            //  8: {{!
	ItemLeftMeta2               //  9: {{
	ItemRightMeta3              // 11: }}}
	ItemRightMeta2              // 12: }}
	ItemDot                     // 13: .
	ItemEqual                   // 14: =
	ItemNumber                  // 15: 42
	ItemString                  // 16: "hello"
	ItemIdentifier              // 17: hello
)

const (
	inHtmlText int = iota
	inHtmlLiteral
	inAttr
	inAttrInterp
)

const (
	leftMeta3        = "{{{"
	leftMetaPound    = "{{#"
	leftMetaSlash    = "{{/"
	leftMetaBang     = "{{!"
	leftMeta2        = "{{"
	leftXmlDeclTag   = "<?xml"
	leftDocTypeTag   = "<!DOCTYPE"
	leftCommentTag   = "<!--"
	leftCloseTag     = "</"
	leftTag          = "<"
	rightMeta3       = "}}}"
	rightMeta2       = "}}"
	rightTag         = ">"
	attrAssignInterp = "=\""
	attrAssign       = "={{"
)

type Item struct {
	Type  ItemType
	Value string
}

func Lex(name, input string) <-chan Item {
	return lex(name, input)
}

type stateFn func(*lexer) stateFn

type lexer struct {
	name    string    // used only for error reports.
	input   string    // the string being scanned.
	start   int       // start position of this item.
	pos     int       // current position in the input.
	width   int       // width of last rune read from input.
	context int       // context of actions
	items   chan Item // channel of scanned items.
}

func lex(name, input string) <-chan Item {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan Item),
	}
	go l.run() // Concurrently run state machine.
	return l.items
}

func (l *lexer) run() {
	for state := lexHtmlText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{
		ItemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lexHtmlText(l *lexer) stateFn {
	l.context = inHtmlText

	for {
		if strings.HasPrefix(l.input[l.pos:], leftMeta3) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexLeftMeta3 // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftMetaPound) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexLeftMetaPound // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftMetaSlash) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexLeftMetaSlash // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftMetaBang) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexLeftMetaBang // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftMeta2) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexLeftMeta2 // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftXmlDeclTag) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexHtmlLiteral // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftDocTypeTag) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexHtmlLiteral // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftCommentTag) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexHtmlLiteral // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftCloseTag) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexHtmlLiteral // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftTag) {
			if l.pos > l.start {
				l.emit(ItemHtmlText)
			}
			return lexHtmlLiteral // Next state.
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(ItemHtmlText)
	}
	l.emit(ItemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

func lexLeftMeta3(l *lexer) stateFn {
	l.pos += len(leftMeta3)
	l.emit(ItemLeftMeta3)
	return lexInsideAction // Now inside {{{ }}}.
}

func lexLeftMetaPound(l *lexer) stateFn {
	l.pos += len(leftMetaPound)
	l.emit(ItemLeftMetaPound)
	return lexInsideAction // Now inside {{# }}.
}

func lexLeftMetaSlash(l *lexer) stateFn {
	l.pos += len(leftMetaSlash)
	l.emit(ItemLeftMetaSlash)
	return lexInsideTerminator // Now inside {{/ }}.
}

func lexLeftMetaBang(l *lexer) stateFn {
	l.pos += len(leftMetaBang)
	l.emit(ItemLeftMetaBang)
	return lexInsideComment // Now inside {{! }}.
}

func lexLeftMeta2(l *lexer) stateFn {
	l.pos += len(leftMeta2)
	l.emit(ItemLeftMeta2)
	return lexInsideAction // Now inside {{ }}.
}

func lexRightMeta3(l *lexer) stateFn {
	l.pos += len(rightMeta3)
	l.emit(ItemRightMeta3)
	switch l.context {
	case inHtmlText:
		return lexHtmlText // Now outside {{{ }}}.
	case inHtmlLiteral:
		return lexHtmlLiteral
	case inAttrInterp:
		return lexHtmlLiteral
	case inAttr:
		return lexHtmlLiteral
	}
	panic("not reached")
}

func lexRightMeta2(l *lexer) stateFn {
	l.pos += len(rightMeta2)
	l.emit(ItemRightMeta2)
	switch l.context {
	case inHtmlText:
		return lexHtmlText // Now outside {{ }}.
	case inHtmlLiteral:
		return lexHtmlLiteral
	case inAttrInterp:
		return lexHtmlLiteral
	case inAttr:
		return lexHtmlLiteral
	}
	panic("not reached")
}

func lexAttrAssign(l *lexer) stateFn {
	l.pos += len(attrAssign)
	l.emit(ItemHtmlAttr)
	return lexInsideAction // Now outside ={{ }}.
}

func lexAttrAssignInterp(l *lexer) stateFn {
	l.pos += len(attrAssignInterp)
	l.emit(ItemHtmlAttrInterp)
	return lexInsideAttrAssignInterp // Now outside =" ".
}

func lexInsideAttrAssignInterp(l *lexer) stateFn {
	l.context = inAttrInterp

	for {
		if strings.HasPrefix(l.input[l.pos:], leftMeta3) {
			if l.pos > l.start {
				l.emit(ItemHtmlLiteral)
			}
			return lexLeftMeta3 // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], leftMeta2) {
			if l.pos > l.start {
				l.emit(ItemHtmlLiteral)
			}
			return lexLeftMeta2 // Next state.
		}

		r := l.next()

		if r == eof {
			break
		}

		if r == '\\' {
			if l.next() == eof {
				break
			}
			continue
		}

		if r == '"' {
			l.emit(ItemHtmlLiteral)
			return lexHtmlLiteral
		}
	}

	return l.errorf("unclosed string")
}

func lexInsideAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate and are ignored.
	// Equal symbols separate and are emitted.
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta3) {
			return lexRightMeta3
		}
		if strings.HasPrefix(l.input[l.pos:], rightMeta2) {
			return lexRightMeta2
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case isSpace(r):
			l.ignore()
		case r == '=':
			l.emit(ItemEqual)
		case r == '.':
			l.emit(ItemDot)
		case r == '"':
			return lexQuote
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		default:
			l.errorf("unexpected character '%s'", string([]rune{r}))
		}
	}

	panic("not reached")
}

func lexHtmlLiteral(l *lexer) stateFn {
	l.context = inHtmlLiteral

	for {
		if strings.HasPrefix(l.input[l.pos:], attrAssignInterp) {
			if l.pos > l.start {
				l.emit(ItemHtmlLiteral)
			}
			return lexAttrAssignInterp
		}

		if strings.HasPrefix(l.input[l.pos:], attrAssign) {
			if l.pos > l.start {
				l.emit(ItemHtmlLiteral)
			}
			return lexAttrAssign
		}

		switch l.next() {
		case eof:
			return l.errorf("unclosed html literal")
		case '>':
			l.emit(ItemHtmlLiteral)
			return lexHtmlText
		default:
			// continue
		}
	}

	panic("not reached")
}

func lexInsideTerminator(l *lexer) stateFn {
	// Ether empty or identifier.
	// Spaces are ignored.
	for {
		r := l.next()
		if r == eof {
			return l.errorf("unclosed terminator")
		} else if isSpace(r) {
			l.ignore()
		} else {
			l.backup()
			break
		}
	}

	cont := lexIdentifier(l)
	if cont == nil {
		return cont
	}

	for {
		r := l.next()
		if r == eof {
			return l.errorf("unclosed terminator")
		} else if isSpace(r) {
			l.ignore()
		} else {
			l.backup()
			break
		}
	}

	return lexRightMeta2
}

func lexInsideComment(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta2) {
			if l.pos > l.start {
				l.emit(ItemHtmlComment)
			}
			return lexRightMeta2
		}

		if l.next() == eof {
			break
		}
	}

	return l.errorf("unclosed comment")
}

func lexIdentifier(l *lexer) stateFn {
	l.accept("abcdefghijklmnopqrstuvwxyz.")
	l.acceptRun("abcdefghijklmnopqrstuvwxyz0123456789_.")
	l.emit(ItemIdentifier)
	return lexInsideAction
}

func lexQuote(l *lexer) stateFn {
	l.ignore()

	for {
		r := l.next()

		if r == eof {
			break
		}

		if r == '\\' {
			if l.next() == eof {
				break
			}
			continue
		}

		if r == '"' {
			l.backup()
			l.emit(ItemString)
			l.next()
			return lexInsideAction
		}
	}

	return l.errorf("unclosed string")
}

func lexNumber(l *lexer) stateFn {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(ItemNumber)
	return lexInsideAction
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

func isAlphaNumeric(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
