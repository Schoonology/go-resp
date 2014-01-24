package resp

import (
	"bufio"
	"io"
	"strconv"
)

type Scanner struct {
	scanner *bufio.Scanner
	bytes   []byte
	next    interface{}
	err     error
}

func NewScanner(reader io.Reader) *Scanner {
	self := new(Scanner)

	self.scanner = bufio.NewScanner(reader)
	self.scanner.Split(bufio.ScanLines)

	return self
}

func (self *Scanner) ident() rune {
	return rune(self.bytes[0])
}

func (self *Scanner) content() []byte {
	return self.bytes[1:]
}

func parseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func (self *Scanner) parseBytes() (bail bool, again bool) {
	if self.bytes == nil || len(self.bytes) == 0 {
		return false, true
	}

	switch self.ident() {
	// Status Reply
	case '+':
		self.next = Status{string(self.content())}
	case '-':
		self.next = parseError(string(self.content()))
	case ':':
		self.next, self.err = strconv.ParseInt(string(self.content()), 10, 64)
	case '$':
		var length int64
		length, self.err = strconv.ParseInt(string(self.content()), 10, 64)
		if self.err != nil {
			bail = true
			return true, false
		}
		if length < 0 {
			self.next = nil
			return false, false
		}

		next := make([]byte, length)
		self.scanner.Split(bufio.ScanBytes)
		for k, _ := range next {
			if !self.scanner.Scan() {
				return true, false
			}
			next[k] = self.scanner.Bytes()[0]
		}
		self.scanner.Split(bufio.ScanLines)
		self.next = string(next)

		// This is to skip the impending \r\n.
		self.scanner.Scan()
	case '*':
		var length int64
		length, self.err = strconv.ParseInt(string(self.content()), 10, 64)
		if self.err != nil {
			return true, false
		}
		next := make([]interface{}, length)
		for i := int64(0); i < length && self.Scan(); i++ {
			next[i] = self.Obj()
		}
		self.next = next
	default:
		self.bytes = self.content()
		return self.parseBytes()
	}

	return false, false
}

func (self *Scanner) Scan() bool {
	if !self.scanner.Scan() {
		return false
	}

	self.bytes = self.scanner.Bytes()
	bail, again := self.parseBytes()

	if bail {
		return false
	}

	if again {
		return self.Scan()
	}

	return self.err == nil
}

func (self *Scanner) Obj() interface{} {
	return self.next
}

func (self *Scanner) Err() error {
	return self.err
}
