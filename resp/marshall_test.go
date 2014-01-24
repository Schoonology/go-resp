package resp

import (
	"testing"
)

func TestMarshallInt(t *testing.T) {
	str, err := Marshall(42)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != ":42\r\n" {
		t.Fatal("Expected Marshall(42) to be :42")
	}
}

func TestMarshallInt8(t *testing.T) {
	str, err := Marshall(int8(42))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != ":42\r\n" {
		t.Fatal("Expected Marshall(42) to be :42")
	}
}

func TestMarshallString(t *testing.T) {
	str, err := Marshall("foobar")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != "$6\r\nfoobar\r\n" {
		t.Fatal("Expected Marshall(42) to be $6foobar")
	}
}

func TestMarshallStatus(t *testing.T) {
	str, err := Marshall(Status{"some status"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != "+some status\r\n" {
		t.Fatal("Expected Marshall(Status) to be +Status")
	}
}

func TestMarshallError(t *testing.T) {
	str, err := Marshall(Error{"BADERR", "some error"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != "-BADERR some error\r\n" {
		t.Fatal("Expected Marshall(Error) to be -BADERR")
	}
}

func TestMarshallMulti(t *testing.T) {
	str, err := Marshall([]interface{}{42, 23, "foo", "bar", true})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != "*5\r\n:42\r\n:23\r\n$3\r\nfoo\r\n$3\r\nbar\r\n:1\r\n" {
		t.Fatalf("Expected Marshall([]) to be *5..., but it was %v", str)
	}
}

func TestMarshallMultiHelper(t *testing.T) {
	str, err := Marshall(NewMulti(42, 23, "foo", "bar", true))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if str != "*5\r\n:42\r\n:23\r\n$3\r\nfoo\r\n$3\r\nbar\r\n:1\r\n" {
		t.Fatalf("Expected Marshall([]) to be *5..., but it was %v", str)
	}
}

func TestUnmarshall(t *testing.T) {
	val, err := Unmarshall(":42")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if val != int64(42) {
		t.Fatal("Expected Unmarshall to return %v", int64(42))
	}
}
