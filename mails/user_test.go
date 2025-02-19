package mails_test

import (
	"strings"
	"testing"

	"github.com/blinkinglight/pocketbase-mysql/mails"
	"github.com/blinkinglight/pocketbase-mysql/tests"
)

func TestSendUserPasswordReset(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// ensure that action url normalization will be applied
	testApp.Settings().Meta.AppUrl = "http://localhost:8090////"

	user, _ := testApp.Dao().FindUserByEmail("test@example.com")

	err := mails.SendUserPasswordReset(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/#/users/confirm-password-reset/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastHtmlBody, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastHtmlBody)
		}
	}
}

func TestSendUserVerification(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.Dao().FindUserByEmail("test@example.com")

	err := mails.SendUserVerification(testApp, user)
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/#/users/confirm-verification/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastHtmlBody, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastHtmlBody)
		}
	}
}

func TestSendUserChangeEmail(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	user, _ := testApp.Dao().FindUserByEmail("test@example.com")

	err := mails.SendUserChangeEmail(testApp, user, "new_test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if testApp.TestMailer.TotalSend != 1 {
		t.Fatalf("Expected one email to be sent, got %d", testApp.TestMailer.TotalSend)
	}

	expectedParts := []string{
		"http://localhost:8090/#/users/confirm-email-change/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
	}
	for _, part := range expectedParts {
		if !strings.Contains(testApp.TestMailer.LastHtmlBody, part) {
			t.Fatalf("Couldn't find %s \nin\n %s", part, testApp.TestMailer.LastHtmlBody)
		}
	}
}
