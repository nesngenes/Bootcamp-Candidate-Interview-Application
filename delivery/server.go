package delivery

import (
	"fmt"
	"interview_bootcamp/config"
	"interview_bootcamp/delivery/controller/api"
	"interview_bootcamp/delivery/middleware"
	"interview_bootcamp/manager"
	"interview_bootcamp/utils/execptions"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	useCaseManager manager.UseCaseManager
	engine         *gin.Engine
	host           string
	log            *logrus.Logger
}

func (s *Server) Run() {
	s.setupControllers()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}

func initializeCloudinaryInstance(cfg *config.Config) (*cloudinary.Cloudinary, error) {
	cloudinaryInstance, err := cloudinary.NewFromParams(
		cfg.CloudinaryConfig.CloudinaryCloudName,
		cfg.CloudinaryConfig.CloudinaryAPIKey,
		cfg.CloudinaryConfig.CloudinaryAPISecret,
	)
	if err != nil {
		return nil, err
	}
	return cloudinaryInstance, nil
}

func (s *Server) setupControllers() {
	s.engine.Use(middleware.LogRequestMiddleware(s.log))
	cfg, err := config.NewConfig()
	cloudinaryInstance, err := initializeCloudinaryInstance(cfg)
	if err != nil {
		panic("Failed to initialize Cloudinary")
	}
	// semua controller disini
	api.NewCandidateController(s.engine, s.useCaseManager.CandidateUseCase(), s.useCaseManager.BootcampUseCase(), cloudinaryInstance)
	api.NewBootcampController(s.engine, s.useCaseManager.BootcampUseCase())
	api.NewInterviewerController(s.engine, s.useCaseManager.InterviewerUseCase())
	api.NewInterviewProcessController(s.engine, s.useCaseManager.InterviewProcessUseCase())
	api.NewInterviewResultController(s.engine, s.useCaseManager.InterviewResultUseCase())
	api.NewResultController(s.engine, s.useCaseManager.ResultUseCase())
	api.NewUserRoleController(s.engine, s.useCaseManager.UserRolesUseCase())                //user role controller
	api.NewUserController(s.engine, s.useCaseManager.UsersUseCase())
	api.NewHRRecruitmentController(s.engine, s.useCaseManager.HRRecruitmentUsecase())
	api.NewStatusController(s.engine, s.useCaseManager.StatusUseCase())
	api.NewAuthController(s.engine, s.useCaseManager.AuthUseCase())
	api.NewFormController(s.engine, s.useCaseManager.FormUseCase(), cloudinaryInstance)
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	cloudinaryInstance, err := initializeCloudinaryInstance(cfg)
	if err != nil {
		panic("Failed to initialize Cloudinary")
	}
	execptions.CheckErr(err)
	infraManager, _ := manager.NewInfraManager(cfg)
	repoManager := manager.NewRepoManager(infraManager)
	useCaseManager := manager.NewUseCaseManager(repoManager)
	useCaseManager.SetCloudinaryInstance(cloudinaryInstance)
	engine := gin.Default()
	host := fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort)
	return &Server{
		useCaseManager: useCaseManager,
		engine:         engine,
		host:           host,
		log:            logrus.New(),
	}
}
