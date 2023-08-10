package middleware

import (
	"interview_bootcamp/config"
	"interview_bootcamp/model"
	"interview_bootcamp/utils/execptions"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogRequestMiddleware(log *logrus.Logger) gin.HandlerFunc {
	cfg, err := config.NewConfig()
	execptions.CheckErr(err)
	file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	execptions.CheckErr(err)
	log.SetOutput(file)

	startTime := time.Now()

	return func(c *gin.Context) {
		c.Next()

		endTime := time.Since(startTime)
		requestLog := model.RequestLog{
			StartTime:  startTime,
			EndTime:    endTime,
			StatusCode: c.Writer.Status(),
			ClientIP:   c.ClientIP(),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			UserAgent:  c.Request.UserAgent(),
		}

		switch {
		case c.Writer.Status() >= 500:
			log.Error(requestLog)
		case c.Writer.Status() >= 400:
			log.Warn(requestLog)
		default:
			log.Info(requestLog)
		}
	}
}
