package model

import (
	"database/sql"
	"time"
)

type Record struct {
	ID        int       `json:"id"`
	Uid       string    `json:"uid"`
	Issue     int       `json:"issue"`
	Type      int       `json:"type"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetRecordsByUid(sql *sql.DB, uid string) ([]Record, error) {
	if sql == nil {
		panic("===== GetRecordsByUid : record RecordHander panic")
	}
	rows, err := sql.Query("SELECT id, uid, issue, type, comment, created_at, updated_at FROM records WHERE uid=?", uid)
	if err != nil {
		return nil, err
	}
	records := []Record{}
	for rows.Next() {
		var model Record
		if err := rows.Scan(&model.ID, &model.Uid, &model.Issue, &model.Type, &model.Comment, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, model)
	}
	return records, nil
}

func (model *Record) Get(sql *sql.DB) error {
	return sql.QueryRow(
		"SELECT uid, issue, type, comment, created_at, updated_at FROM records WHERE id=?",
		model.ID,
	).Scan(&model.Uid, &model.Issue, &model.Type, &model.Comment, &model.CreatedAt, &model.UpdatedAt)
}

func (model *Record) Update(sql *sql.DB) error {
	_, err := sql.Exec(
		"UPDATE records SET issue=?, type=?, comment=? WHERE id=?",
		model.Issue, model.Type, model.Comment, model.ID,
	)
	return err
}
func (model *Record) Delete(sql *sql.DB) error {
	_, err := sql.Exec("DELETE FROM records WHERE id=$1", model.ID)
	return err
}
func (model *Record) Create(sql *sql.DB) error {
	err := sql.QueryRow(
		"INSERT INTO records(uid, issue, type, comment) VALUES(?, ?, ?, ?) RETURNING id",
		model.Uid, model.Issue, model.Type, model.Comment,
	).Scan(&model.ID)
	if err != nil {
		return err
	}
	return nil
}
