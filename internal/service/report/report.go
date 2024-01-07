package report

import (
	"forum/internal/models"
	repo "forum/internal/repository"
)

type ReportServer struct {
	report repo.Report
	post   repo.Post
}

func NewReportService(repo *repo.Repository) *ReportServer {
	return &ReportServer{report: repo.Report, post: repo.Post}
}

func (s *ReportServer) Create(report *models.CreateReport) error {
	return s.report.Create(report)
}

func (s *ReportServer) GetAll() ([]*models.Report, error) {
	return s.report.GetAll()
}

func (s *ReportServer) DeleteById(reportId, resp int) error {
	if resp == models.Approve {
		report, err := s.report.GetById(reportId)
		if err != nil {
			return err
		}
		return s.post.DeleteById(report.PostId)
	}

	if resp == models.Refuse {
		return s.report.DeleteById(reportId)
	}

	return models.ErrReport
}
