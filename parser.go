// Copyright 2015 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package csgo_cfg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type tokenType int

const (
	_TOKEN_INVALID tokenType = iota
	_TOKEN_COMMENT
	_TOKEN_SECTION
	_TOKEN_KEY
)

type parser struct {
	buf     *bufio.Reader
	isEOF   bool
	count   int
	comment *bytes.Buffer
}

func newParser(r io.Reader) *parser {
	return &parser{
		buf:     bufio.NewReader(r),
		count:   1,
		comment: &bytes.Buffer{},
	}
}

// BOM handles header of BOM-UTF8 format.
// http://en.wikipedia.org/wiki/Byte_order_mark#Representations_of_byte_order_marks_by_encoding
func (p *parser) BOM() error {
	mask, err := p.buf.Peek(3)
	if err != nil && err != io.EOF {
		return err
	} else if len(mask) < 3 {
		return nil
	} else if mask[0] == 239 && mask[1] == 187 && mask[2] == 191 {
		p.buf.Read(mask)
	}
	return nil
}

func (p *parser) readUntil(delim byte) ([]byte, error) {
	data, err := p.buf.ReadBytes(delim)
	if err != nil {
		if err == io.EOF {
			p.isEOF = true
		} else {
			return nil, err
		}
	}
	return data, nil
}

func cleanComment(in []byte) ([]byte, bool) {
	i := bytes.IndexAny(in, "#;")
	if i == -1 {
		return nil, false
	}
	return in[i:], true
}

func readKeyName(in []byte) (string, int, error) {
	line := string(in)

	// Get out key name
	endIdx := strings.Index(line, " ")
	if endIdx < 0 {
		return "", -1, ErrDelimiterNotFound{line}
	}
	return strings.TrimSpace(line[0:endIdx]), endIdx + 1, nil
}

// hasSurroundedQuote check if and only if the first and last characters
// are quotes \" or \'.
// It returns false if any other parts also contain same kind of quotes.
func hasSurroundedQuote(in string, quote byte) bool {
	return len(in) > 2 && in[0] == quote && in[len(in)-1] == quote &&
		strings.IndexByte(in[1:], quote) == len(in)-2
}

// value, isString, comment, error
func (p *parser) readValue(in []byte) (string, bool, string, error) {
	line := strings.TrimLeftFunc(string(in), unicode.IsSpace)

	// Line empty, or nothing before the comment
	if len(line) == 0 || (len(line) > 1 && line[0:1] == "//") {
		return "", true, "", nil
	}

	// Drop everything after space, unless space is within a string
	var valueEndIndex int
	insideQuote := false
	for _, c := range line {
		if c == ' ' && !insideQuote {
			break
		}
		if c == '"' {
			insideQuote = !insideQuote
		}
		valueEndIndex++
	}

	// Ended the string without matching quote. Invalid
	if insideQuote {
		return "", false, "", ErrDelimiterNotFound{line}
	}

	// Remove
	valueLine := strings.TrimSpace(line[:valueEndIndex])
    comment := ""
    if len(in) > valueEndIndex {
        comment = line[valueEndIndex:]
    }    

	if hasSurroundedQuote(valueLine, '\'') ||
		hasSurroundedQuote(valueLine, '"') {
		return valueLine[1 : len(valueLine)-1], true, comment, nil
	}

	// Should be a number at this point
	return valueLine, false, comment, nil
}

// parse parses data through an io.Reader.
func (f *File) parse(reader io.Reader) (err error) {
	p := newParser(reader)
	if err = p.BOM(); err != nil {
		return fmt.Errorf("BOM: %v", err)
	}

	section, _ := f.NewSection(DEFAULT_SECTION)

	var line []byte
	for !p.isEOF {
		line, err = p.readUntil('\n')
		if err != nil {
			return err
		}

		line = bytes.TrimLeftFunc(line, unicode.IsSpace)
		if len(line) == 0 {
			continue
		}

		kname, offset, err := readKeyName(line)
		if err != nil {
			return err
		}

		key, err := section.NewKey(kname, "")
		if err != nil {
			return err
		}

		value, isString, comment, err := p.readValue(line[offset:])
		if err != nil {
			return err
		}

		key.SetValue(value)
		key.isString = isString

		// Comments
        comment = strings.TrimLeftFunc(comment, unicode.IsSpace)
		if len(comment) > 1 && comment[0:2] == "//" {
			key.Comment = strings.TrimSpace(comment)
		}
	}
	return nil
}
