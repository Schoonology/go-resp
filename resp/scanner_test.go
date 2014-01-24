package resp

import (
	"strings"
	"testing"
)

func testSingle(t *testing.T, input string, test func(interface{}) (bool, interface{})) {
	s := NewScanner(strings.NewReader(input))

	if !s.Scan() {
		t.Fatal("Expected Scan() to return true the first time.")
	}

	match, expected := test(s.Obj())

	if !match {
		t.Fatalf("Expected Obj() to return %v, but it returned %v.", expected, s.Obj())
	}

	if s.Scan() {
		t.Fatal("Expected Scan() to return false the second time.")
	}
}

func TestScanStatusReply(t *testing.T) {
	testSingle(t, "+OK", func(val interface{}) (match bool, expected interface{}) {
		expected = Status{"OK"}
		status, match := val.(Status)

		switch {
		case !match, status.Message != "OK":
			match = false
		}

		return
	})
}

func TestScanErrorReply(t *testing.T) {
	testSingle(t, "-BADERR Some message", func(val interface{}) (match bool, expected interface{}) {
		expected = Error{"BADERR", "Some message"}
		err, match := val.(Error)

		switch {
		case !match:
		case err.Type != expected.(Error).Type:
			match = false
		case err.Message != expected.(Error).Message:
			match = false
		}

		return
	})
}

func TestScanIntegerReply(t *testing.T) {
	testSingle(t, ":123", func(val interface{}) (bool, interface{}) {
		return val == int64(123), int64(123)
	})
}

func TestScanBulkReply(t *testing.T) {
	testSingle(t, "$3\r\nfoo", func(val interface{}) (bool, interface{}) {
		return val == "foo", "foo"
	})
}

func TestScanNullBulkReply(t *testing.T) {
	testSingle(t, "$-1", func(val interface{}) (bool, interface{}) {
		return val == nil, nil
	})
}

func TestScanNewlineBulkReply(t *testing.T) {
	testSingle(t, "$6\r\n12\n456\r\n", func(val interface{}) (bool, interface{}) {
		return val == "12\n456", "12\\n456"
	})
}

func TestScanMultilineBulkReply(t *testing.T) {
	testSingle(t, "*3\r\n$-1\r\n:42\r\n$4\r\nawef", func(val interface{}) (match bool, expected interface{}) {
		expected = []interface{}{int64(42), "awef"}
		arr, match := val.([]interface{})

		switch {
		case !match:
		case arr[0] != nil:
			match = false
		case arr[1] != int64(42):
			match = false
		case arr[2] != "awef":
			match = false
		}

		return
	})
}

func TestScanPastBadData(t *testing.T) {
	testSingle(t, "!!!:42", func(val interface{}) (bool, interface{}) {
		return val == int64(42), int64(42)
	})
}

func TestScanPastEmptyLines(t *testing.T) {
	testSingle(t, "\r\n\n\n\r\n:42", func(val interface{}) (bool, interface{}) {
		return val == int64(42), int64(42)
	})
}
