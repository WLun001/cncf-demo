package main

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var h *Handler
var ee *echo.Echo

func TestMain(m *testing.M) {
	// Write code here to run before tests

	_ = os.Setenv("MONGO_URI", "mongodb://admin:apple123@localhost:27017")
	handler, e, err := run()
	if err != nil {
		panic(err)
	}

	h = handler
	ee = e
	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	stop()
	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestAddItem(t *testing.T) {

	body := `{"name":"item 2","quantity":2}`

	req := httptest.NewRequest(
		http.MethodPost, "/items", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := ee.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.AddItem(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestGetItems(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	c := ee.NewContext(req, rec)

	if assert.NoError(t, h.GetItems(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
