// テスト実行コマンド：cd backend/; go test model/*
package model

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRecordsByUid_positive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	data := Record{
		ID:        1,
		Uid:       "uid",
		Issue:     0,
		Type:      0,
		Comment:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()}

	columns := []string{"id", "uid", "issue", "type", "comment", "created_at", "updated_at"}
	mock.ExpectQuery("SELECT").
		WithArgs(data.Uid).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(data.ID, data.Uid, data.Issue, data.Type, data.Comment, data.CreatedAt, data.UpdatedAt))

	actual, err := GetRecordsByUid(db, data.Uid)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, data.ID, actual[0].ID)
	assert.Equal(t, data.Uid, actual[0].Uid)
	assert.Equal(t, data.Issue, actual[0].Issue)
	assert.Equal(t, data.Type, actual[0].Type)
	assert.Equal(t, data.Comment, actual[0].Comment)
	assert.Equal(t, data.CreatedAt, actual[0].CreatedAt)
	assert.Equal(t, data.UpdatedAt, actual[0].UpdatedAt)
}
