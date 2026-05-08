package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		key       string
		wantError bool
	}{
		{"Valid encryption", "Hello World", "secret123", false},
		{"Empty text", "", "secret123", false},
		{"Short key", "test", "12345", true},
		{"Unicode text", "سلام دنیا", "mykey123", false},
		{"Long text", string(make([]byte, 10000)), "securekey", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := encrypt(tt.text, tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("encrypt() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			if tt.wantError {
				return
			}
			
			decrypted, err := decrypt(encrypted, tt.key)
			if err != nil {
				t.Errorf("decrypt() error = %v", err)
				return
			}
			
			if decrypted != tt.text {
				t.Errorf("decrypt() = %v, want %v", decrypted, tt.text)
			}
		})
	}
}

func TestEncryptHandler(t *testing.T) {
	req := Request{Text: "test message", Key: "secret123"}
	body, _ := json.Marshal(req)
	
	httpReq := httptest.NewRequest("POST", "/encrypt", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	encryptHandler(w, httpReq)
	
	resp := w.Result()
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDecryptHandler(t *testing.T) {
	original := "test message"
	key := "secret123"
	
	encrypted, _ := encrypt(original, key)
	
	req := Request{Text: encrypted, Key: key}
	body, _ := json.Marshal(req)
	
	httpReq := httptest.NewRequest("POST", "/decrypt", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	decryptHandler(w, httpReq)
	
	resp := w.Result()
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestInvalidKey(t *testing.T) {
	encrypted := "test"
	errMsg := "decryption failed - wrong key or corrupted data"
	
	req := Request{Text: encrypted, Key: "wrongkey"}
	body, _ := json.Marshal(req)
	
	httpReq := httptest.NewRequest("POST", "/decrypt", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	decryptHandler(w, httpReq)
	
	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	
	if resp.Error == "" {
		t.Error("Expected error for invalid key")
	}
	
	if resp.Error != errMsg && resp.Error != "invalid encrypted text format" {
		t.Logf("Got error: %s", resp.Error)
	}
}
