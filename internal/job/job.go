package job

import (
	"hyacinth-backend/internal/repository"
	"hyacinth-backend/pkg/jwt"
	"hyacinth-backend/pkg/log"
	"hyacinth-backend/pkg/sid"
)

type Job struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
