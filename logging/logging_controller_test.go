package logging

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChangeLogLevel(t *testing.T) {
	testLogger := GetLogger("c_test")
	body := "{\n    \"lvl\":\"error\",\n    \"packageName\":\"c_test\"\n}"
	request, err := http.NewRequest("GET", "/test", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	ChangeLogLevel(response, request)
	actualLogger := createTestLogger(1, "c_test")
	loggersEqual(t, &actualLogger, testLogger)
}

func TestChangeLogLevel_WithBrokenJson(t *testing.T) {
	body := "{\n \": }"
	request, err := http.NewRequest("GET", "/test", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	ChangeLogLevel(response, request)
	errStr := "{\"error\":\"Invalid request payload\"}"
	assert.JSONEq(t, errStr, response.Body.String())
}

func TestChangeLogLevel_WithNonExistingLogger(t *testing.T) {
	testLogger := GetLogger("NotExist")
	nBody := "{\n    \"lvl\":\"error\",\n    \"packageName\":\"something\"\n}"
	request, err := http.NewRequest("GET", "/test", strings.NewReader(nBody))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	ChangeLogLevel(response, request)
	actualLogger := createTestLogger(3, "NotExist")
	loggersEqual(t, &actualLogger, testLogger)
}
