package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
	"github.com/mikail-tommard/2025-12-10-links/internal/usecase"
)

type Server struct {
	mux           *http.ServeMux
	linksService  *usecase.LinksService
	reportService *usecase.ReportService
}

type linkRequest struct {
	LinksList []string `json:"links_list"`
}

type reportRequest struct {
	LinksNum []int `json:"links_num"`
}

type linkResultResponse struct {
	URL    string `json:"url"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type linkResponse struct {
	LinksNum int                  `json:"links_num"`
	Results []linkResultResponse `json:"results"`
}

func NewServer(link *usecase.LinksService, report *usecase.ReportService) *Server {
	s := &Server{
		mux:           http.NewServeMux(),
		linksService:  link,
		reportService: report,
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/links", s.handleLinks)
	s.mux.HandleFunc("/report", nil)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req linkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(req.LinksList) == 0 {
		http.Error(w, "link_list must not be empty", http.StatusBadRequest)
		return
	}
	
	ctx := r.Context()
	
	batch, err := s.linksService.CreateAndCheckBatch(ctx, req.LinksList)
	if err != nil {
		http.Error(w, "Failed to create batch", http.StatusInternalServerError)
		return
	}

	respResult := make([]linkResultResponse, 0, len(batch.Results))
	for _, r := range batch.Results {
		respResult = append(respResult, linkResultResponse{
			URL: r.Link.URL,
			Status: string(r.Status),
			Error: r.Error,
		})
	}

	resp := linkResponse{
		LinksNum: int(batch.ID),
		Results: respResult,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req reportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.LinksNum) == 0 {
		http.Error(w, "links_num must not be empty", http.StatusBadRequest)
		return
	}

	ids := make([]domain.BatchID, 0, len(req.LinksNum))
	for _, v := range req.LinksNum {
		ids = append(ids, domain.BatchID(v))
	}
	
	ctx := r.Context()
	bytes, err := s.reportService.GenerateReportForBatches(ctx, ids)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(bytes); err != nil {
		return
	}
}