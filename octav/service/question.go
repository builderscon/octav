package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/internal/context"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *QuestionSvc) Init() {}

func (v *QuestionSvc) populateRowForCreate(ctx context.Context, vdb *db.Question, payload *model.CreateQuestionRequest) error {
	vdb.EID = tools.UUID()
	vdb.SessionID = payload.SessionID
	vdb.UserID = context.GetUserID(ctx)
	vdb.Body = payload.Body

	return nil
}

func (v *QuestionSvc) populateRowForUpdate(ctx context.Context, vdb *db.Question, payload *model.UpdateQuestionRequest) error {
	if payload.SessionID.Valid() {
		vdb.SessionID = payload.SessionID.String
	}

	if payload.Body.Valid() {
		vdb.Body = payload.Body.String
	}

	return nil
}