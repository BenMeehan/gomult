package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benmeehan/gomult/internal/ai"
	"github.com/benmeehan/gomult/internal/config"
	"github.com/benmeehan/gomult/internal/sandbox"
)

type compileRequest struct {
	Code     string `json:"code"`
	Input    string `json:"input"`
	Language string `json:"language"`
	Analyze  bool   `json:"analyze"`
}

type Server struct {
	httpServer  *http.Server
	engine      *sandbox.Engine
	aiClient    *ai.Client
	maxBodySize int64
}

func New(cfg *config.Config, engine *sandbox.Engine, aiClient *ai.Client) *Server {
	mux := http.NewServeMux()
	s := &Server{
		engine:      engine,
		aiClient:    aiClient,
		maxBodySize: cfg.Server.MaxCodeSize * 2,
		httpServer: &http.Server{
			Addr:         formatAddr(cfg.Server.Port),
			Handler:      mux,
			ReadTimeout:  cfg.Server.ReadTimeoutDuration(),
			WriteTimeout: cfg.Server.WriteTimeoutDuration(),
			IdleTimeout:  120 * time.Second,
		},
	}

	mux.HandleFunc("/compile", s.handleCompile)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/languages", s.handleLanguages)

	return s
}

func (s *Server) ListenAndServe() error {
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleCompile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)

	var req compileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		http.Error(w, "language is required", http.StatusBadRequest)
		return
	}
	if req.Code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	log.Printf("compile: lang=%s code_len=%d", req.Language, len(req.Code))

	result := s.engine.Execute(req.Language, req.Code, req.Input)

	if s.aiClient != nil && req.Analyze {
		log.Printf("ai analysis: lang=%s status=%s", req.Language, result.Status.String())
		analysis, err := s.aiClient.Analyze(req.Code, req.Language, req.Input, result.Output, result.Status.String())
		if err != nil {
			log.Printf("ai analysis failed: %v", err)
			result.Output += "\n\n[AI ANALYSIS FAILED: " + err.Error() + "]"
		} else {
			result.Output += "\n\n[AI ANALYSIS]\n" + analysis
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch result.Status {
	case sandbox.StatusOK:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result.Output))
	case sandbox.StatusCompileError:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[COMPILE ERROR]\n" + result.Output))
	case sandbox.StatusTimeLimitExceeded:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[TIME LIMIT EXCEEDED]"))
	case sandbox.StatusRuntimeError:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(runtimeErrorOutput(result)))
	case sandbox.StatusUnsupportedLanguage, sandbox.StatusCodeTooLarge:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(result.Output))
	default:
		http.Error(w, result.Output, http.StatusInternalServerError)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Get("detail") == "true" {
		json.NewEncoder(w).Encode(s.engine.LanguagesDetail())
	} else {
		json.NewEncoder(w).Encode(s.engine.Languages())
	}
}

func formatAddr(port int) string {
	return fmt.Sprintf(":%d", port)
}

func runtimeErrorOutput(result *sandbox.ExecuteResult) string {
	if result.ExitCode != 0 && result.ExitCode != -1 {
		return fmt.Sprintf("[RUNTIME ERROR]\nexit code: %d\n%s", result.ExitCode, result.Output)
	}
	return "[RUNTIME ERROR]\n" + result.Output
}
