package service

import (
	"github.com/alexreagan/rabbit/g"
	"testing"
)

func init() {
	g.ParseConfig("../../config/cfg.json")
	g.InitDBPool()

	//router := gin.Default()
	//utils.SetRouter(router)
}

func TestClean(t *testing.T) {
	latestTime := CaasService.GetNameSpaceLatestTime()
	oneDayBeforeLatestTime := latestTime.AddDate(0, 0, -1)
	CaasService.DeleteNameSpaceBeforeTime(oneDayBeforeLatestTime)
	CaasService.DeleteServiceBeforeTime(oneDayBeforeLatestTime)
	CaasService.DeletePodBeforeTime(oneDayBeforeLatestTime)
	CaasService.DeletePortBeforeTime(oneDayBeforeLatestTime)
}
