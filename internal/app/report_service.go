package app

import (
	"mime/multipart"

	"github.com/elginbrian/FILKOMPLAIN-BE/config"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/repository"
)

type ReportService struct {
	Repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{Repo: repo}
}

func (s *ReportService) ListReports() ([]model.Report, error) {
	return s.Repo.GetAllReports()
}

func (s *ReportService) GetReport(id uint) (*model.Report, error) {
	return s.Repo.GetReportByID(id)
}

func (s *ReportService) CreateReport(report *model.Report) error {
	return s.Repo.CreateReport(report)
}

func (s *ReportService) UpdateReport(report *model.Report) error {
	return s.Repo.UpdateReport(report)
}

func (s *ReportService) DeleteReport(id uint) error {
	return s.Repo.DeleteReport(id)
}

func (s *ReportService) ResolveReportStatus(id uint, status string) error {
	return s.Repo.UpdateReportStatus(id, status)
}

func (s *ReportService) ReplyReport(id uint, reply string) error {
	return s.Repo.UpdateReportReply(id, reply)
}

func (s *ReportService) UploadAttachment(file *multipart.FileHeader) (string, error) {
	supabaseStorage := config.NewSupabaseStorage()
	
	fileURL, err := supabaseStorage.UploadFile(file)
	if err != nil {
		return "", err
	}
	
	return fileURL, nil
}
