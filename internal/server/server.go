package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/env"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/handler"
)

type Server struct {
	handler handler.HandlerInterface
	mux     *http.ServeMux
	env     *env.Env
}

func NewServer(handler handler.HandlerInterface, env *env.Env) *Server {
	return &Server{
		handler: handler,
		mux:     http.NewServeMux(),
		env:     env,
	}
}

func (s *Server) Run() {
	s.routes()
	fmt.Println("🚀 Servidor iniciado com sucesso na porta " + s.env.SERVER_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.env.SERVER_PORT), s.mux))
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /", s.handler.Home)
	s.mux.HandleFunc("GET /concurso", s.handler.GetContest)
	s.mux.HandleFunc("GET /concurso/{concurso}", s.handler.GetNumberContest)
	s.mux.HandleFunc("GET /resultado/", s.handler.GenerateDailyContest)
	s.mux.HandleFunc("GET /jogo", s.handler.GenerateDailyGame)
}
