package models_test

import (
	"testing"

	"github.com/blinkinglight/pocketbase-mysql/models"
)

func TestParamTableName(t *testing.T) {
	m := models.Param{}
	if m.TableName() != "_params" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}
