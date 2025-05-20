package repository

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"gorm.io/gorm"
)

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) GetAllReports() ([]model.Report, error) {
	var reports []model.Report
	err := r.DB.Find(&reports).Error
	return reports, err
}

func (r *ReportRepository) GetReportByID(id uint) (*model.Report, error) {
	var report model.Report
	err := r.DB.First(&report, id).Error
	return &report, err
}

func (r *ReportRepository) CreateReport(report *model.Report) error {
	return r.DB.Create(report).Error
}

func (r *ReportRepository) UpdateReport(report *model.Report) error {
	return r.DB.Save(report).Error
}

func (r *ReportRepository) DeleteReport(id uint) error {
	return r.DB.Delete(&model.Report{}, id).Error
}

func (r *ReportRepository) UpdateReportStatus(id uint, status string) error {
	return r.DB.Model(&model.Report{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ReportRepository) UpdateReportReply(id uint, reply string) error {
	return r.DB.Model(&model.Report{}).Where("id = ?", id).Update("reply", reply).Error
}
