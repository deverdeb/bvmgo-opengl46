package assetsmngr

import (
	"fmt"
	"testing"
)

func buildManager() *Manager[string] {
	manager := CreateManager[string]()
	manager.Register("first",
		func() (*string, error) { msg := "hello first"; return &msg, nil },
		func(s *string) {})
	manager.Register("second",
		func() (*string, error) { msg := "hello second"; return &msg, nil },
		func(s *string) {})
	manager.Register("error",
		func() (*string, error) { return nil, fmt.Errorf("big test error") },
		func(s *string) {})
	return &manager
}

func TestManager_Get(t *testing.T) {
	manager := buildManager()
	asset, err := manager.Get("first")
	if err != nil {
		t.Fatalf("Get() error=\"%v\", want error=nil", err)
	}
	if asset == nil {
		t.Fatalf("Get() asset=nil, want asset=\"hello first\"")
	}
	if *asset != "hello first" {
		t.Fatalf("Get() asset=\"%v\", want asset=\"hello first\"", *asset)
	}

	asset, err = manager.Get("error")
	if err == nil {
		t.Fatalf("Get() error=nil, want error=\"%v\"", fmt.Errorf("big test error"))
	}
}

func TestManager_LoadAll(t *testing.T) {
	manager := buildManager()
	err := manager.LoadAll()
	if err == nil {
		t.Fatalf("LoadAll() error=nil, want error=\"%v\"", fmt.Errorf("big test error"))
	}

	asset, err := manager.Get("first")
	if err != nil {
		t.Fatalf("Get() error=\"%v\", want error=nil", err)
	}
	if asset == nil {
		t.Fatalf("Get() asset=nil, want asset=\"hello first\"")
	}
	if *asset != "hello first" {
		t.Fatalf("Get() asset=\"%v\", want asset=\"hello first\"", *asset)
	}
}
