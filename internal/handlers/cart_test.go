package handlers

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestAddItem(t *testing.T) {
    payload := `{"id":"1","name":"Test Item","price":9.99}`
    req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(payload))
    w := httptest.NewRecorder()

    AddItem(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
    }
    var item CartItem
    if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }
    if item.ID != "1" || item.Name != "Test Item" || item.Price != 9.99 {
        t.Fatalf("unexpected item returned: %+v", item)
    }
}
