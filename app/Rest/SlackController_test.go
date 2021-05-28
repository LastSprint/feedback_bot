package Rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// mockDb is stub for testing SlackController
type mockDb struct {
	// messages is used for holding messages written by caller
	messages []string
	// err is used as return value in WriteFeedbackMessage
	err error
}

// WriteFeedbackMessage stub implementation which saves messages and return err
// for details look at mockDb declaration
func (m *mockDb) WriteFeedbackMessage(s string) error {
	m.messages = append(m.messages, s)
	return m.err
}

// --------- Test Cases ---------
// 1. POST request will be handled
// 2. GET request won't be handled
// 3. Request with empty message will return error
// 4. Db error will be returned to caller
// 5. If no error is occurred method will return 200 OK and db will have message in messages

// POST request will be handled
func TestSlackController_handleFeedbackCommand_POSTRequestWillBeHandled(t *testing.T) {
	// Arrange

	body := url.Values{
		"text": []string{"text"},
	}.Encode()

	db := mockDb{}
	controller := SlackController{DB: &db}

	// Act

	responseWritter := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/slack/cto/feedback", strings.NewReader(body))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(body)))

	controller.handleFeedbackCommand(responseWritter, request)

	response := responseWritter.Result()

	// Assert

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected %v, but got %v", http.StatusOK, response.StatusCode)
		return
	}
}

// GET request won't be handled
func TestSlackController_handleFeedbackCommand_GetRequestWontBeHandled(t *testing.T) {
	// Arrange

	body := url.Values{
		"text": []string{"text"},
	}.Encode()

	db := mockDb{}
	controller := SlackController{DB: &db}

	// Act

	responseWritter := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/slack/cto/feedback", strings.NewReader(body))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(body)))

	controller.handleFeedbackCommand(responseWritter, request)

	response := responseWritter.Result()

	// Assert

	if response.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Expected %v, but got %v", http.StatusNotImplemented, response.StatusCode)
		return
	}
}

// Request with empty message will return error
func TestSlackController_handleFeedbackCommand_EmptyMessageProduceError(t *testing.T) {
	// Arrange

	body := url.Values{
		"text": []string{""},
	}.Encode()

	db := mockDb{}
	controller := SlackController{DB: &db}

	// Act

	responseWritter := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/slack/cto/feedback", strings.NewReader(body))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(body)))

	controller.handleFeedbackCommand(responseWritter, request)

	response := responseWritter.Result()

	// Assert

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected %v, but got %v", http.StatusBadRequest, response.StatusCode)
		return
	}
}

// Db error will be returned to caller
func TestSlackController_handleFeedbackCommand_DbErrorWillBeReturned(t *testing.T) {
	// Arrange

	body := url.Values{
		"text": []string{"text"},
	}.Encode()

	db := mockDb{
		err: errors.New("err"),
	}
	controller := SlackController{DB: &db}

	// Act

	responseWritter := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/slack/cto/feedback", strings.NewReader(body))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(body)))

	controller.handleFeedbackCommand(responseWritter, request)

	response := responseWritter.Result()

	// Assert

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected %v, but got %v", http.StatusInternalServerError, response.StatusCode)
		return
	}
}

// If no error is occurred method will return 200 OK and db will have message in messages
func TestSlackController_handleFeedbackCommand_WriteRequestWillBeWrittenInDB(t *testing.T) {
	// Arrange

	body := url.Values{
		"text": []string{"text"},
	}.Encode()

	db := mockDb{}
	controller := SlackController{DB: &db}

	// Act

	responseWritter := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/slack/cto/feedback", strings.NewReader(body))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(body)))

	controller.handleFeedbackCommand(responseWritter, request)

	response := responseWritter.Result()

	// Assert

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected %v, but got %v", http.StatusOK, response.StatusCode)
		return
	}

	if len(db.messages) == 0 {
		t.Fatal("DB is empty")
		return
	}
}