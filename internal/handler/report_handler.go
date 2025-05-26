package handler

import (
	"strconv"
	"time"

	"github.com/elginbrian/FILKOMPLAIN-BE/internal/app"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/request"
	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	Service *app.ReportService
}

func NewReportHandler(service *app.ReportService) *ReportHandler {
	return &ReportHandler{Service: service}
}

func toReportData(r *model.Report) response.ReportData {
	return response.ReportData{
		ID:          r.ID,
		UserName:    r.UserName,
		Title:       r.Title,
		Content:     r.Content,
		Place:       r.Place,
		PhoneNumber: r.PhoneNumber,
		Status:      r.Status,
		Attachment:  r.Attachment,
		Reply:       r.Reply,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *ReportHandler) GetReports(c *fiber.Ctx) error {
	reports, err := h.Service.ListReports()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	var data []response.ReportData
	for _, r := range reports {
		data = append(data, toReportData(&r))
	}
	return c.JSON(response.Response{
		Success: true,
		Data:    response.ReportListData{Reports: data},
	})
}

func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid id",
		})
	}
	report, err := h.Service.GetReport(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Success: false,
			Error:   "not found",
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Data:    toReportData(report),
	})
}

func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
    if form, err := c.MultipartForm(); err == nil {
        var report model.Report

        if len(form.Value["user_name"]) > 0 {
            report.UserName = form.Value["user_name"][0]
        }
        if len(form.Value["title"]) > 0 {
            report.Title = form.Value["title"][0]
        }
        if len(form.Value["content"]) > 0 {
            report.Content = form.Value["content"][0]
        }
        if len(form.Value["place"]) > 0 {
            report.Place = form.Value["place"][0]
        }
        if len(form.Value["phone_number"]) > 0 {
            report.PhoneNumber = form.Value["phone_number"][0]
        }
        if len(form.Value["status"]) > 0 {
            report.Status = form.Value["status"][0]
        }

        if report.UserName == "" {
            return c.Status(fiber.StatusBadRequest).JSON(response.Response{
                Success: false,
                Error:   "user_name is required",
            })
        }
        if report.PhoneNumber == "" {
             return c.Status(fiber.StatusBadRequest).JSON(response.Response{
                Success: false,
                Error:   "phone_number is required",
            })
        }
        if report.Title == "" {
            return c.Status(fiber.StatusBadRequest).JSON(response.Response{
                Success: false,
                Error:   "title is required",
            })
        }
     
        if files := form.File["attachment"]; len(files) > 0 {
            file := files[0]
            attachmentURL, err := h.Service.UploadAttachment(file) 
            if err != nil {
                return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
                    Success: false,
                    Error:   "Failed to upload attachment: " + err.Error(),
                })
            }
            report.Attachment = attachmentURL
        }
        
        if err := h.Service.CreateReport(&report); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
                Success: false,
                Error:   err.Error(),
            })
        }
        
        return c.Status(fiber.StatusCreated).JSON(response.Response{
            Success: true,
            Message: "report created",
            Data:    toReportData(&report),
        })
    }
    
    var report model.Report
    if err := c.BodyParser(&report); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.Response{
            Success: false,
            Error:   err.Error(),
        })
    }
    
    if report.UserName == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.Response{
            Success: false,
            Error:   "user_name is required",
        })
    }
    if report.PhoneNumber == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.Response{
            Success: false,
            Error:   "phone_number is required",
        })
    }
    if report.Title == "" {
        return c.Status(fiber.StatusBadRequest).JSON(response.Response{
            Success: false,
            Error:   "title is required",
        })
    }
    
    if err := h.Service.CreateReport(&report); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
            Success: false,
            Error:   err.Error(),
        })
    }
    
    return c.Status(fiber.StatusCreated).JSON(response.Response{
        Success: true,
        Message: "report created",
        Data:    toReportData(&report),
    })
}

func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid id",
		})
	}
	var report model.Report
	if err := c.BodyParser(&report); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	report.ID = uint(id)
	if err := h.Service.UpdateReport(&report); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Message: "report updated",
		Data:    toReportData(&report),
	})
}

func (h *ReportHandler) DeleteReport(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid id",
		})
	}
	if err := h.Service.DeleteReport(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Message: "report deleted",
	})
}

func (h *ReportHandler) ResolveReportStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid id",
		})
	}
	var body request.ResolveReportStatusRequest
	if err := c.BodyParser(&body); err != nil || body.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid status",
		})
	}
	if err := h.Service.ResolveReportStatus(uint(id), body.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Message: "report status updated",
	})
}

func (h *ReportHandler) ReplyReport(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid id",
		})
	}
	var body request.ReplyReportRequest
	if err := c.BodyParser(&body); err != nil || body.Reply == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Success: false,
			Error:   "invalid reply",
		})
	}
	if err := h.Service.ReplyReport(uint(id), body.Reply); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}
	return c.JSON(response.Response{
		Success: true,
		Message: "report replied",
	})
}
