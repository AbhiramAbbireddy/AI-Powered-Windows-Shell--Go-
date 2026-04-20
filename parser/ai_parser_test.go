package parser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAIParserParse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected auth header %q", got)
		}

		fmt.Fprint(w, `{
			"choices": [
				{
					"message": {
						"role": "assistant",
						"content": "{\"action\":\"list_folders\",\"target\":\"current_directory\",\"source\":\"\",\"destination\":\"\",\"requires_info\":false,\"clarification\":\"\",\"explanation\":\"Lists folders in the current directory.\"}"
					}
				}
			]
		}`)
	}))
	defer server.Close()

	parser := NewAIParser("test-key", server.URL, "llama-test", 2*time.Second)
	intent, err := parser.Parse("what are the folders in this directory")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if intent.Action != ActionListFolders {
		t.Fatalf("expected action %q, got %q", ActionListFolders, intent.Action)
	}
	if intent.Target != "current_directory" {
		t.Fatalf("expected target current_directory, got %q", intent.Target)
	}
}

func TestAIParserParseRequiresInfo(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"choices": [
				{
					"message": {
						"role": "assistant",
						"content": "{\"action\":\"delete_file\",\"target\":\"\",\"source\":\"\",\"destination\":\"\",\"requires_info\":true,\"clarification\":\"Which file should I delete?\",\"explanation\":\"Deletes the requested file.\"}"
					}
				}
			]
		}`)
	}))
	defer server.Close()

	parser := NewAIParser("test-key", server.URL, "llama-test", 2*time.Second)
	intent, err := parser.Parse("delete file")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !intent.RequiresInfo {
		t.Fatal("expected intent to require clarification")
	}
	if intent.Clarification != "Which file should I delete?" {
		t.Fatalf("unexpected clarification %q", intent.Clarification)
	}
}

func TestAIParserParseInvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"not-json"}}]}`)
	}))
	defer server.Close()

	parser := NewAIParser("test-key", server.URL, "llama-test", 2*time.Second)
	_, err := parser.Parse("something")
	if err != ErrAIResponse {
		t.Fatalf("expected %v, got %v", ErrAIResponse, err)
	}
}
