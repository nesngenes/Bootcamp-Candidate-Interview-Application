package delivery

import (
	"fmt"
	"interview_bootcamp/config"
	"interview_bootcamp/delivery/controller/api"
	"interview_bootcamp/manager"
	"interview_bootcamp/utils/execptions"

	"github.com/gin-gonic/gin"
)

type Server struct {
	useCaseManager manager.UseCaseManager
	engine         *gin.Engine
	host           string
}

func (s *Server) Run() {
	s.setupControllers()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}

func (s *Server) setupControllers() {
	// semua controller disini
	api.NewCandidateController(s.engine, s.useCaseManager.CandidateUseCase())
	api.NewBootcampController(s.engine, s.useCaseManager.BootcampUseCase())
	api.NewResumeController(s.engine, s.useCaseManager.ResumeUseCase())
	api.NewInterviewerController(s.engine, s.useCaseManager.InterviewerUseCase())
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	execptions.CheckErr(err)
	infraManager, _ := manager.NewInfraManager(cfg)
	repoManager := manager.NewRepoManager(infraManager)
	useCaseManager := manager.NewUseCaseManager(repoManager)
	engine := gin.Default()
	host := fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort)
	return &Server{
		useCaseManager: useCaseManager,
		engine:         engine,
		host:           host,
	}
}
