package html2templateLiteral
import{
	"errors" //required for erros
	"fmt" // for utility functions
	"io" // for input files
	"golang.org/x/net/html" // converting HTML
	"strings" 
	"unicode"
}
var CouldNotParseErr = errors.New("could not parse html")
var CouldNotConvertErr = errors.New("could not convert html")

func Convert(r io.Reader, w io.Writer) error {
	tokenizer := html.NewTokenizer(r)

	for {
		tokenizer.Next()
		token := tokenizer.Token()
		err := tokenizer.Err()
		if err == io.EOF {
			return nil
		}
		if err != nil && err != io.EOF {
			return fmt.Errorf("%w: %s", CouldNotParseErr, err.Error())
		}
		jsxToken, err := jsxToken(token)
		if err != nil {
			return fmt.Errorf("%w: %s", CouldNotConvertErr, err.Error())
		}
		w.Write([]byte(jsxToken))
	}
}

func jsxToken(token html.Token) (string, error) {
	switch token.Type {
	case html.StartTagToken:
		jsxToken, err := jsxStartTag(token, false)
		if err != nil {
			return "", err
		}
		return jsxToken, nil
	case html.EndTagToken:
		return jsxEndTag(token), nil
	case html.SelfClosingTagToken:
		jsxToken, err := jsxStartTag(token, true)
		if err != nil {
			return "", err
		}
		return jsxToken, nil
	case html.TextToken:
		return token.Data, nil
	case html.CommentToken:
		fallthrough
	case html.DoctypeToken:
		return "", nil
	}
	return "", fmt.Errorf("unexpected token type %s encountered", token.Type)
}

type Attribute struct {
	Key, Value string
}

const StartTag, EndTag, Slash, Equals, Space, Quote = '<', '>', '/', '=', ' ', '"'

func jsxStartTag(t html.Token, selfClosing bool) (string, error) {
	attrs, err := jsxAttributes(t)
	if err != nil {
		return "", fmt.Errorf("could not build jsx compatible start tag: %w", err)
	}

	var b strings.Builder

	b.WriteRune(StartTag)
	b.WriteString(t.Data)
	for i := range attrs {
		b.WriteRune(Space)
		b.WriteString(attrs[i].Key)
		b.WriteRune(Equals)
		b.WriteRune(Quote)
		b.WriteString(attrs[i].Value)
		b.WriteRune(Quote)
	}
	if selfClosing {
		b.WriteRune(Slash)
	}
	b.WriteRune(EndTag)
	return b.String(), nil
}
func jsxEndTag(t html.Token) string {
	var b strings.Builder
	b.WriteRune(StartTag)
	b.WriteRune(Slash)
	b.WriteString(t.Data)
	b.WriteRune(EndTag)
	return b.String()
}

func jsxAttributes(t html.Token) ([]Attribute, error) {
	var attrs = make([]Attribute, len(t.Attr))
	for i := range t.Attr {
		attr, err := convertAttribute(t.Attr[i])
		if err != nil {
			return nil, fmt.Errorf("could not convert attributes: %w", err)
		}
		attrs[i] = attr
	}
	return attrs, nil
}

func convertAttribute(attr html.Attribute) (Attribute, error) {
	var jsxAttr Attribute
	jsxKey, err := convertAttributeKey(attr.Key)
	if err != nil {
		return Attribute{}, fmt.Errorf("could not convert attribute : %w", err)
	}
	jsxAttr.Key = jsxKey
	jsxAttr.Value = attr.Val
	return jsxAttr, nil
}

func convertAttributeKey(key string) (string, error) {
	// aria Attributes are the exception when it comes to camelCasing Attributes in JSX
	if strings.HasPrefix(key, "aria-") {
		return key, nil
	}
	if key == "class" {
		return "className", nil
	} else {
		jsxKey, err := kebapToCamel(key)
		if err != nil {
			return "", fmt.Errorf("could not convert key : %w", err)
		}
		return jsxKey, nil
	}
}

var ErrorMultipleHyphens = errors.New("multiple hypens in attribute name detected")

func kebapToCamel(kebap string) (string, error) {
	var hyphenIndex int
	var b strings.Builder
	for i, r := range kebap {
		if r == '-' { // kebap stick detected
			if hyphenIndex != 0 && hyphenIndex < i {
				return "", ErrorMultipleHyphens
			}
			hyphenIndex = i
			continue
		}

		if hyphenIndex != 0 && i == hyphenIndex+1 { // rune after the hypen
			b.WriteRune(unicode.ToUpper(r))
			continue
		}

		b.WriteRune(r)
	}
	return b.String(), nil
}