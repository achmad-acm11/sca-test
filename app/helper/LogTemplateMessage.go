package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sca-integrator/app/shareVar"
	"time"
)

var timeElapsed *logTime

type StandartLog struct {
	TypeLayer shareVar.LayerName
	Name      shareVar.NameEntity
	NameFunc  string
}

func NewStandardLog(nameEntity shareVar.NameEntity, typeLayer shareVar.LayerName) *StandartLog {
	return &StandartLog{
		TypeLayer: typeLayer,
		Name:      nameEntity,
	}
}

func (log *StandartLog) StartFunction(request interface{}) {
	if gin.Mode() != gin.TestMode {
		timeElapsed = NewLogTime(time.Now())
		shareVar.Logger.WithField(string(shareVar.Request), request).Info(fmt.Sprintf("Start %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}

func (log *StandartLog) WarningFunction(message interface{}) {
	if gin.Mode() != gin.TestMode {
		shareVar.Logger.WithField("message", message).Warning(fmt.Sprintf("Warning %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}

func (log *StandartLog) InfoFunction(message interface{}) {
	if gin.Mode() != gin.TestMode {
		shareVar.Logger.WithField("message", message).Info(fmt.Sprintf("Warning %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}

func (log *StandartLog) ErrorFunction(message interface{}) {
	if gin.Mode() != gin.TestMode {
		shareVar.Logger.WithField("message", message).Error(fmt.Sprintf("Warning %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}

func (log *StandartLog) EndFunction(response interface{}) {
	if gin.Mode() != gin.TestMode {
		shareVar.Logger.WithFields(logrus.Fields{
			string(shareVar.Response):    response,
			string(shareVar.TimeElapsed): timeElapsed.GetTimeSince(),
		}).Info(fmt.Sprintf("End %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}
