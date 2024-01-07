package report

import (
	"database/sql"

	"forum/internal/models"
)

type ReportSqlite struct {
	db *sql.DB
}

func NewReportSqlite(db *sql.DB) *ReportSqlite {
	return &ReportSqlite{db: db}
}

func (r *ReportSqlite) Create(report *models.CreateReport) error {
	query := "INSERT INTO posts_reports (post_id, content, moderator_id, moderator_name, create_at) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.Exec(query, report.PostId, report.Content, report.ModeratorId, report.ModeratorName, report.CreateAt)
	if err != nil && err.Error() == models.IncorRequest {
		return models.ErrPost
	}

	return nil
}

func (r *ReportSqlite) GetById(reportId int) (*models.Report, error) {
	report := &models.Report{}
	query := "SELECT * FROM posts_reports WHERE id = ?"
	err := r.db.QueryRow(query, reportId).Scan(&report.Id, &report.PostId, &report.Content,
		&report.ModeratorId, &report.ModeratorName, &report.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrReport
		}
		return nil, err
	}
	return report, nil
}

func (r *ReportSqlite) GetAll() ([]*models.Report, error) {
	query := "SELECT * FROM posts_reports"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	reports := make([]*models.Report, 0)
	for rows.Next() {
		report := new(models.Report)
		err := rows.Scan(&report.Id, &report.PostId, &report.Content,
			&report.ModeratorId, &report.ModeratorName, &report.CreateAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportSqlite) DeleteById(reportId int) error {
	query := "DELETE FROM posts_reports WHERE id = ?"

	res, err := r.db.Exec(query, reportId)
	if err != nil {
		return err
	}
	
	countRes, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if countRes == 0 {
		return models.ErrReport
	}
	return nil
}
