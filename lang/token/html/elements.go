package html

import (
	"strconv"
)

type Element uint

const (
	ILLEGAL Element = iota
	FOREIGN

	void_elements_beg
	AREA
	BASE
	BR
	COL
	EMBED
	HR
	IMG
	INPUT
	KEYGEN
	LINK
	MENUITEM
	META
	PARAM
	SOURCE
	TRACK
	WBR
	void_elements_end

	raw_text_elements_beg
	SCRIPT
	STYLE
	raw_text_elements_end

	rcdata_elements_beg
	TEXTAREA
	TITLE
	rcdata_elements_end

	normal_elements_beg
	A
	ABBR
	ADDRESS
	ARTICLE
	ASIDE
	AUDIO
	B
	BDI
	BDO
	BDOY
	BLOCKQUOTE
	BUTTON
	CANVAS
	CAPTION
	CITE
	CODE
	COLGROUP
	COMMAND
	DATALIST
	DD
	DEL
	DETAILS
	DFN
	DIV
	DL
	DT
	EM
	FIELDSET
	FIGCAPTION
	FIGURE
	FOOTER
	FORM
	H1
	H2
	H3
	H4
	H5
	H6
	HEAD
	HEADER
	HGROUP
	HTML
	I
	IFRAME
	INS
	KDB
	LABEL
	LEGEND
	LI
	MAP
	MARK
	MENU
	METER
	NAV
	NOSCRIPT
	OBJECT
	OL
	OPTGROUP
	OPTION
	OUTPUT
	P
	PRE
	PROGRESS
	Q
	RP
	RT
	RUBY
	S
	SAMP
	SECTION
	SELECT
	SMALL
	SPAN
	STRONG
	SUB
	SUMMARY
	SUP
	TABLE
	TBODY
	TD
	TFOOT
	TH
	THEAD
	TIME
	TR
	U
	UL
	VAR
	VIDEO
	normal_elements_end
)

var tokens = [...]string{
	AREA:     "area",
	BASE:     "base",
	BR:       "br",
	COL:      "col",
	EMBED:    "embed",
	HR:       "hr",
	IMG:      "img",
	INPUT:    "input",
	KEYGEN:   "keygen",
	LINK:     "link",
	MENUITEM: "menuitem",
	META:     "meta",
	PARAM:    "param",
	SOURCE:   "source",
	TRACK:    "track",
	WBR:      "wbr",

	SCRIPT: "script",
	STYLE:  "style",

	TEXTAREA: "textarea",
	TITLE:    "title",

	A:          "a",
	ABBR:       "abbr",
	ADDRESS:    "address",
	ARTICLE:    "article",
	ASIDE:      "aside",
	AUDIO:      "audio",
	B:          "b",
	BDI:        "bdi",
	BDO:        "bdo",
	BDOY:       "bdoy",
	BLOCKQUOTE: "blockquote",
	BUTTON:     "button",
	CANVAS:     "canvas",
	CAPTION:    "caption",
	CITE:       "cite",
	CODE:       "code",
	COLGROUP:   "colgroup",
	COMMAND:    "command",
	DATALIST:   "datalist",
	DD:         "dd",
	DEL:        "del",
	DETAILS:    "details",
	DFN:        "dfn",
	DIV:        "div",
	DL:         "dl",
	DT:         "dt",
	EM:         "em",
	FIELDSET:   "fieldset",
	FIGCAPTION: "figcaption",
	FIGURE:     "figure",
	FOOTER:     "footer",
	FORM:       "form",
	H1:         "h1",
	H2:         "h2",
	H3:         "h3",
	H4:         "h4",
	H5:         "h5",
	H6:         "h6",
	HEAD:       "head",
	HEADER:     "header",
	HGROUP:     "hgroup",
	HTML:       "html",
	I:          "i",
	IFRAME:     "iframe",
	INS:        "ins",
	KDB:        "kdb",
	LABEL:      "label",
	LEGEND:     "legend",
	LI:         "li",
	MAP:        "map",
	MARK:       "mark",
	MENU:       "menu",
	METER:      "meter",
	NAV:        "nav",
	NOSCRIPT:   "noscript",
	OBJECT:     "object",
	OL:         "ol",
	OPTGROUP:   "optgroup",
	OPTION:     "option",
	OUTPUT:     "output",
	P:          "p",
	PRE:        "pre",
	PROGRESS:   "progress",
	Q:          "q",
	RP:         "rp",
	RT:         "rt",
	RUBY:       "ruby",
	S:          "s",
	SAMP:       "samp",
	SECTION:    "section",
	SELECT:     "select",
	SMALL:      "small",
	SPAN:       "span",
	STRONG:     "strong",
	SUB:        "sub",
	SUMMARY:    "summary",
	SUP:        "sup",
	TABLE:      "table",
	TBODY:      "tbody",
	TD:         "td",
	TFOOT:      "tfoot",
	TH:         "th",
	THEAD:      "thead",
	TIME:       "time",
	TR:         "tr",
	U:          "u",
	UL:         "ul",
	VAR:        "var",
	VIDEO:      "video",
}

var elements = func() map[string]Element {
	m := make(map[string]Element, len(tokens))
	for i, name := range tokens {
		e := Element(i)
		if e.IsValidElement() {
			m[name] = e
		}
	}
	return m
}()

func Lookup(ident string) Element {
	if e, is_element := elements[ident]; is_element {
		return e
	}
	return FOREIGN
}

func (e Element) String() string {
	s := ""
	if 0 <= e && e < Element(len(tokens)) {
		s = tokens[e]
	}
	if s == "" {
		s = "element(" + strconv.Itoa(int(e)) + ")"
	}
	return s
}

func (e Element) IsValidElement() bool {
	return e.IsVoidElement() || e.IsRawTextElement() || e.IsRcdataElement() || e.IsNormalElement()
}

func (e Element) IsVoidElement() bool {
	return void_elements_beg < e && e < void_elements_end
}

func (e Element) IsRawTextElement() bool {
	return raw_text_elements_beg < e && e < raw_text_elements_end
}

func (e Element) IsRcdataElement() bool {
	return rcdata_elements_beg < e && e < rcdata_elements_end
}

func (e Element) IsNormalElement() bool {
	return normal_elements_beg < e && e < normal_elements_end
}
