package controllers

import (
    "testing"
    "net/http"
    "net/http/httptest"
)

func TestBaseIndexCallsRender(t *testing.T){
    r := httptest.NewRequest(http.MethodGet, "/", nil)
    w := httptest.NewRecorder()
    Helper = MockHelper{}
    BaseIndex(w, r) 
    if w.Header().Get("RenderTemplateCalled") != "true" {
        t.Fatalf("Render Template was not called")
    }
}
