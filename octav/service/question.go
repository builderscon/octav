package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Question) populateRowForCreate(vdb *db.Question, payload model.CreateQuestionRequest) error {
	vdb.EID = tools.UUID()
	vdb.SessionID = payload.SessionID
	vdb.UserID = payload.UserID
	vdb.Body = payload.Body

	return nil
}

func (v *Question) populateRowForUpdate(vdb *db.Question, payload model.UpdateQuestionRequest) error {
	if payload.SessionID.Valid() {
		vdb.SessionID = payload.SessionID.String
	}

	if payload.UserID.Valid() {
		vdb.UserID = payload.UserID.String
	}

	if payload.Body.Valid() {
		vdb.Body = payload.Body.String
	}

	return nil
}