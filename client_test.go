package gotyrant

import (
	"bytes"
	"testing"
	"time"
)

var client *Client

func TestPut(t *testing.T) {
	err := client.Put([]byte("Hello"), []byte("World"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPutKeep(t *testing.T) {
	key := []byte("Hello")
	value1 := []byte("World")
	value2 := []byte("There")
	err := client.Put(key, value1)
	if err != nil {
		t.Fatal(err)
	}
	err = client.PutKeep(key, value2)
	if err == nil {
		t.Fatal(err)
	}
	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !e.IsExist() {
		t.Fatalf("Unexpected error: %s", e)
	}
}

func TestGet(t *testing.T) {
	key := []byte("Hello")
	value := []byte("World")
	err := client.Put(key, value)
	if err != nil {
		t.Fatal(err)
	}
	val, err := client.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(value, val) != 0 {
		t.Fatalf("Value not match")
	}
}

func TestOut(t *testing.T) {
	key := []byte("Hello")
	err := client.Out(key)
	if err != nil {
		e, ok := err.(*Error)
		if !ok {
			t.Fatalf("Unexpected error: %s", err)
		}
		if !e.IsNotExist() {
			t.Fatalf("Unexpected error: %s", err)
		}
	}

	err = client.Out(key)
	if err == nil {
		t.Fatalf("Expect not exist error")
	}
	e, ok := err.(*Error)
	if !ok {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !e.IsNotExist() {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func InitClient() {
	var err error
	client, err = NewClient(Config{
		Addr:    "10.125.36.71:60000",
		Timeout: time.Second * 1,
	})
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	InitClient()
	defer client.Close()
	m.Run()
}
