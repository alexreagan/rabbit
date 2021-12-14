package service

import (
	"github.com/alexreagan/rabbit/g"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/caas"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
	"time"
)

type caasService struct {
}

func (s *caasService) GetNameSpaceLatestTime() time.Time {
	var latest time.Time
	db := g.Con().Portal.Model(caas.NameSpace{})
	db = db.Select("max(update_time)")
	db = db.Find(&latest)
	return latest
}

func (s *caasService) DeleteNameSpaceBeforeTime(timestamp time.Time) {
	var namespaces []*caas.NameSpace
	db := g.Con().Portal.Begin()
	dt := db.Model(caas.NameSpace{}).Debug()
	dt = dt.Where("update_time < ?", timestamp)
	dt = dt.Find(&namespaces)
	if len(namespaces) == 0 {
		return
	}

	var ids []int64
	for _, ns := range namespaces {
		ids = append(ids, ns.ID)
	}

	var rels []*caas.NamespaceServiceRel
	dt = db.Model(caas.NamespaceServiceRel{}).Debug()
	dt = dt.Where("namespace_id in (?)", ids)
	if dt = dt.Delete(&rels); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}

	dt = db.Model(caas.NameSpace{})
	dt = dt.Where("id in (?)", ids)
	if dt = dt.Delete(&namespaces); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}

	db.Commit()
}

func (s *caasService) GetServiceLatestTime() time.Time {
	var latest time.Time
	db := g.Con().Portal.Model(caas.Service{})
	db = db.Select("max(update_time)")
	db = db.Find(&latest)
	return latest
}

func (s *caasService) DeleteServiceBeforeTime(timestamp time.Time) {
	db := g.Con().Portal.Begin()

	var services []*caas.Service
	dt := g.Con().Portal.Model(caas.Service{})
	dt = dt.Where("update_time < ?", timestamp)
	dt = dt.Find(&services)
	if len(services) == 0 {
		return
	}

	var ids []int64
	for _, service := range services {
		ids = append(ids, service.ID)
	}
	var rels []*caas.NamespaceServiceRel
	dt = db.Model(caas.NamespaceServiceRel{})
	dt = dt.Where("service in (?)", ids)
	if dt = dt.Delete(&rels); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}

	var servicePodRels []*caas.ServicePodRel
	dt = db.Model(caas.ServicePodRel{})
	dt = dt.Where("service in (?)", ids)
	if dt = dt.Delete(&servicePodRels); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}

	var servicePortRels []*caas.ServicePortRel
	dt = db.Model(caas.ServicePortRel{})
	dt = dt.Where("service in (?)", ids)
	if dt = dt.Delete(&servicePortRels); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}

	dt = db.Model(caas.Service{})
	dt = dt.Where("id in (?)", ids)
	if dt = dt.Delete(&services); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}

	db.Commit()
}

func (s *caasService) GetPodLatestTime() time.Time {
	var latest time.Time
	db := g.Con().Portal.Model(caas.Pod{})
	db = db.Select("max(update_time)")
	db = db.Find(&latest)
	return latest
}

func (s *caasService) DeletePodBeforeTime(timestamp time.Time) {
	var pods []*caas.Pod
	db := g.Con().Portal.Begin()
	dt := db.Model(caas.Pod{})
	dt = dt.Where("update_time < ?", timestamp)
	if dt = dt.Delete(&pods); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}
	db.Commit()
}

func (s *caasService) GetPortLatestTime() time.Time {
	var latest time.Time
	db := g.Con().Portal.Model(caas.Port{})
	db = db.Select("max(update_time)")
	db = db.Find(&latest)
	return latest
}

func (s *caasService) DeletePortBeforeTime(timestamp time.Time) {
	var ports []*caas.Port
	db := g.Con().Portal.Begin()
	dt := db.Model(caas.Port{})
	dt = dt.Where("update_time < ?", timestamp)
	if dt = dt.Delete(&ports); dt.Error != nil {
		db.Rollback()
		log.Error(dt.Error)
		return
	}
	db.Commit()
}

// 含有tagIDs的Pods
func (s *caasService) PodsHavingTagIDs(tagIDs []int64) []*caas.Pod {
	var tIDs []int
	for _, i := range tagIDs {
		tIDs = append(tIDs, int(i))
	}
	sort.Ints(tIDs)

	var tmp []string
	for _, i := range tIDs {
		tmp = append(tmp, strconv.Itoa(i))
	}

	var pods caas.Pods
	db := g.Con().Portal.Model(caas.Pod{}).Debug()
	db = db.Joins("left join `caas_service_pod_rel` on `caas_service_pod_rel`.`pod`=`caas_pod`.`id`")
	db = db.Joins("left join `caas_service_tag_rel` on `caas_service_tag_rel`.`service`=`caas_service_pod_rel`.`service`")
	db = db.Where("`caas_service_tag_rel`.`tag` in (?)", tagIDs)
	db = db.Group("`caas_service_pod_rel`.`pod`")
	db = db.Having("group_concat(`caas_service_tag_rel`.`tag` order by `caas_service_tag_rel`.`tag`) = ?", strings.Join(tmp, ","))
	db.Find(&pods)

	pods.Sort()
	return pods
}

func (s *caasService) GetAllService() []*caas.Service {
	var services []*caas.Service
	db := g.Con().Portal.Model(caas.Service{})
	db.Find(&services)
	return services
}

func (this caasService) GetServiceRelatedTags(service *caas.Service) []*app.Tag {
	var tags []*app.Tag
	db := g.Con().Portal.Model(app.Tag{}).Debug()
	db = db.Select("`tag`.*, `tag_category`.name as category_name")
	db = db.Joins("left join `caas_service_tag_rel` on `caas_service_tag_rel`.`tag` = `tag`.id")
	db = db.Joins("left join `tag_category` on `tag_category`.id = `tag`.category_id")
	db = db.Where("`caas_service_tag_rel`.`service` = ?", service.ID)
	db.Find(&tags)
	return tags
}

func newCaasService() *caasService {
	return &caasService{}
}
