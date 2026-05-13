package openapi

import (
	"encoding/json"
	"testing"
)

func TestEmbeddedSpecIsValidJSON(t *testing.T) {
	var decoded map[string]any
	if err := json.Unmarshal(spec, &decoded); err != nil {
		t.Fatalf("OpenAPI spec is not valid JSON: %v", err)
	}
	if decoded["openapi"] == "" {
		t.Fatal("OpenAPI spec must include openapi version")
	}
	if decoded["paths"] == nil {
		t.Fatal("OpenAPI spec must include paths")
	}
}

func TestEmbeddedSpecDocumentsCurrentRoutes(t *testing.T) {
	var decoded struct {
		Paths map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(spec, &decoded); err != nil {
		t.Fatalf("OpenAPI spec is not valid JSON: %v", err)
	}

	expected := []string{
		"/health",
		"/api/v1/version",
		"/api/v1/openapi.json",
		"/api/v1/tarot/cards",
		"/api/v1/tarot/cards/{sourceCode}",
		"/api/v1/tarot/spreads",
		"/api/v1/tarot/spreads/{code}",
		"/api/v1/tarot/readings",
		"/api/v1/tarot/readings/{id}",
		"/api/v1/wallet",
		"/api/v1/coin-transactions",
		"/api/v1/check-ins",
		"/api/v1/lucky-colors/today",
		"/api/v1/lucky-foods/today",
		"/api/v1/lucky-items/today",
		"/api/v1/avoidance/today",
		"/api/v1/daily-insights/today",
		"/api/v1/moon/today",
		"/api/v1/moon/birthday",
		"/api/v1/moon/reports/{id}",
	}

	for _, path := range expected {
		if _, ok := decoded.Paths[path]; !ok {
			t.Fatalf("OpenAPI spec is missing path %s", path)
		}
	}
}
