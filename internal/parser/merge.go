package parser

import (
	"LogSentry/internal/models"
)

func MergeReports(dst *models.LogReport, src models.LogReport){
	dst.Errors = append(dst.Errors, src.Errors...)
	dst.Warns = append(dst.Warns, src.Warns...)
	dst.Infos = append(dst.Infos, src.Infos...)
	dst.Unknown = append(dst.Unknown, src.Unknown...)
	dst.Counts["ERROR"] += src.Counts["ERROR"]
	dst.Counts["WARN"] += src.Counts["WARN"]
	dst.Counts["INFO"] += src.Counts["INFO"]
	dst.Counts["DEFAULT"] += src.Counts["DEFAULT"]

}