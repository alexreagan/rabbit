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
	tx := g.Con().Portal.Model(caas.NameSpace{})
	tx = tx.Select("max(update_time)")
	tx = tx.Find(&latest)
	return latest
}

func (s *caasService) DeleteNameSpaceBeforeTime(timestamp time.Time) {
	var namespaces []*caas.NameSpace
	tx := g.Con().Portal.Begin()
	tx = tx.Model(caas.NameSpace{})
	tx = tx.Where("update_time < ?", timestamp)
	tx = tx.Find(&namespaces)
	if len(namespaces) == 0 {
		tx.Rollback()
		return
	}

	var ids []int64
	for _, ns := range namespaces {
		ids = append(ids, ns.ID)
	}

	var rels []*caas.NamespaceServiceRel
	tx = tx.Model(caas.NamespaceServiceRel{})
	tx = tx.Where("namespace_id in (?)", ids)
	if err := tx.Delete(&rels).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}

	tx = tx.Model(caas.NameSpace{})
	tx = tx.Where("id in (?)", ids)
	if err := tx.Delete(&namespaces).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}

	tx.Commit()
}

func (s *caasService) GetServiceLatestTime() time.Time {
	var latest time.Time
	tx := g.Con().Portal.Model(caas.Service{})
	tx = tx.Select("max(update_time)")
	tx = tx.Find(&latest)
	return latest
}

func (s *caasService) DeleteServiceBeforeTime(timestamp time.Time) {
	tx := g.Con().Portal.Begin()

	var services []*caas.Service
	tx = tx.Model(caas.Service{})
	tx = tx.Where("update_time < ?", timestamp)
	tx = tx.Find(&services)
	if len(services) == 0 {
		tx.Rollback()
		return
	}

	var ids []int64
	for _, service := range services {
		ids = append(ids, service.ID)
	}
	var rels []*caas.NamespaceServiceRel
	tx = tx.Model(caas.NamespaceServiceRel{})
	tx = tx.Where("service in (?)", ids)
	if err := tx.Delete(&rels).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}

	var servicePodRels []*caas.ServicePodRel
	tx = tx.Model(caas.ServicePodRel{})
	tx = tx.Where("service in (?)", ids)
	if err := tx.Delete(&servicePodRels).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}

	var servicePortRels []*caas.ServicePortRel
	tx = tx.Model(caas.ServicePortRel{})
	tx = tx.Where("service in (?)", ids)
	if err := tx.Delete(&servicePortRels).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}

	tx = tx.Model(caas.Service{})
	tx = tx.Where("id in (?)", ids)
	if err := tx.Delete(&services).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}

	tx.Commit()
}

func (s *caasService) GetPodLatestTime() time.Time {
	var latest time.Time
	tx := g.Con().Portal.Model(caas.Pod{})
	tx = tx.Select("max(update_time)")
	tx = tx.Find(&latest)
	return latest
}

func (s *caasService) DeletePodBeforeTime(timestamp time.Time) {
	var pods []*caas.Pod
	tx := g.Con().Portal.Begin()
	tx = tx.Model(caas.Pod{})
	tx = tx.Where("update_time < ?", timestamp)
	if err := tx.Delete(&pods).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}
	tx.Commit()
}

func (s *caasService) GetPortLatestTime() time.Time {
	var latest time.Time
	tx := g.Con().Portal.Model(caas.Port{})
	tx = tx.Select("max(update_time)")
	tx = tx.Find(&latest)
	return latest
}

func (s *caasService) DeletePortBeforeTime(timestamp time.Time) {
	var ports []*caas.Port
	tx := g.Con().Portal.Begin()
	tx = tx.Model(caas.Port{})
	tx = tx.Where("update_time < ?", timestamp)
	if err := tx.Delete(&ports).Error; err != nil {
		tx.Rollback()
		log.Error(err)
		return
	}
	tx.Commit()
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
	tx := g.Con().Portal.Model(caas.Pod{}).Debug()
	tx = tx.Joins("left join `caas_service_pod_rel` on `caas_service_pod_rel`.`pod`=`caas_pod`.`id`")
	tx = tx.Joins("left join `caas_service_tag_rel` on `caas_service_tag_rel`.`service`=`caas_service_pod_rel`.`service`")
	tx = tx.Where("`caas_service_tag_rel`.`tag` in (?)", tagIDs)
	tx = tx.Group("`caas_service_pod_rel`.`pod`")
	tx = tx.Having("group_concat(`caas_service_tag_rel`.`tag` order by `caas_service_tag_rel`.`tag`) = ?", strings.Join(tmp, ","))
	tx.Find(&pods)

	pods.Sort()
	return pods
}

func (s *caasService) GetAllService() []*caas.Service {
	var services []*caas.Service
	tx := g.Con().Portal.Model(caas.Service{})
	tx.Find(&services)
	return services
}

func (this caasService) GetServiceRelatedTags(service *caas.Service) []*app.Tag {
	var tags []*app.Tag
	tx := g.Con().Portal.Model(app.Tag{}).Debug()
	tx = tx.Select("`tag`.*, `tag_category`.name as category_name")
	tx = tx.Joins("left join `caas_service_tag_rel` on `caas_service_tag_rel`.`tag` = `tag`.id")
	tx = tx.Joins("left join `tag_category` on `tag_category`.id = `tag`.category_id")
	tx = tx.Where("`caas_service_tag_rel`.`service` = ?", service.ID)
	tx.Find(&tags)
	return tags
}

func newCaasService() *caasService {
	return &caasService{}
}
