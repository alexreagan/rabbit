package caas

import (
	"github.com/alexreagan/rabbit/server/model/gtime"
)

type Port struct {
	ID            int64       `json:"id" gorm:"primary_key;column:id"`
	ContainerPort string      `json:"containerPort" gorm:"column:container_port;type:string;size:128;comment:"`
	Host          string      `json:"host" gorm:"column:host;type:string;size:128;comment:"`
	NodePort      string      `json:"nodePort" gorm:"column:node_port;type:string;size:128;comment:"`
	Protocol      string      `json:"protocol" gorm:"column:protocol;type:string;size:128;comment:"`
	UpdateTime    gtime.GTime `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
}

func (this Port) TableName() string {
	return "caas_port"
}

type Service struct {
	ID          int64  `json:"id" gorm:"primary_key;column:id"`
	Type        string `json:"type" gorm:"column:type;type:string;size:128;comment:"`
	AppID       int64  `json:"appId" gorm:"column:app_id;comment:"`
	AppName     string `json:"appName" gorm:"column:app_name;type:string;size:128;comment:"`
	ServiceName string `json:"serviceName" gorm:"column:service_name;type:string;size:128;comment:"`
	//Ports        []*Port `json:"ports,omitempty" gorm:"many2many:caas_service_port_rel;foreignkey:service_name;association_foreignkey:host;association_jointable_foreignkey:host;jointable_foreignkey:service_name;"`
	Ports []*Port `json:"ports,omitempty" gorm:"-"`
	//Envs         *Envs  `json:"envs,omitempty" gorm:"foreignkey:app_id;association_foreignkey:id;association_autoupdate:false"`
	ImageName    string `json:"imageName" gorm:"column:image_name;type:string;size:128;comment:"`
	ImageTag     string `json:"imageTag" gorm:"column:image_tag;type:string;size:128;comment:"`
	Replicas     int64  `json:"replicas" gorm:"column:replicas;comment:"`
	NowReplicas  int64  `json:"nowReplicas" gorm:"column:now_replicas;comment:"`
	Cpu          int64  `json:"cpu" gorm:"column:cpu;comment:"`
	Gpu          int64  `json:"gpu" gorm:"column:gpu;comment:"`
	Memory       int64  `json:"memory" gorm:"column:memory;comment:"`
	AffinityType string `json:"affinityType" gorm:"column:affinity_type;type:string;size:128;comment:"`
	//NodeSelectorLabel     string    `json:"nodeSelectorLabel" gorm:"column:node_selector_label;type:string;size:128;comment:"`
	//Zones                 string    `json:"zones" gorm:"column:zones;type:string;size:128;comment:"`
	//ReadinessProbeInfo    string    `json:"readinessProbeInfo" gorm:"column:readiness_probe_info;type:string;size:128;comment:"`
	//LivenessProbeInfo     string    `json:"livenessProbeInfo" gorm:"column:liveness_probe_info;type:string;size:128;comment:"`
	HeadlessName string `json:"headlessName" gorm:"column:headless_name;type:string;size:128;comment:"`
	PullPolicy   string `json:"pullPolicy" gorm:"column:pull_policy;type:string;size:128;comment:"`
	//VolumeList            string    `json:"volumeList" gorm:"column:volume_list;type:string;size:128;comment:"`
	//LocalVolumeList       string    `json:"localVolumeList" gorm:"column:local_volume_list;type:string;size:128;comment:"`
	//ConfigMapVolumeList   string    `json:"configMapVolumeList" gorm:"column:config_map_volume_list;type:string;size:128;comment:"`
	SrvLbType string `json:"srvLbType" gorm:"column:srv_lb_type;type:string;size:128;comment:"`
	//Command               string    `json:"command" gorm:"column:command;type:string;size:128;comment:"`
	//Arg                   string    `json:"arg" gorm:"column:arg;type:string;size:128;comment:"`
	Completions           int64       `json:"completions" gorm:"column:completions;comment:"`
	Parallelism           int64       `json:"parallelism" gorm:"column:parallelism;comment:"`
	ActiveDeadlineSeconds int64       `json:"activeDeadlineSeconds" gorm:"column:active_deadline_seconds;comment:"`
	ClusterIP             string      `json:"clusterIP" gorm:"column:cluster_ip;type:string;size:128;comment:"`
	CreateTime            gtime.GTime `json:"createTime" gorm:"column:create_time"`
	UpdateTime            gtime.GTime `json:"updateTime" gorm:"column:update_time;default:null;comment:"`
	FinishTime            gtime.GTime `json:"finishTime" gorm:"column:finish_time"`
	Duration              string      `json:"duration" gorm:"column:duration;type:string;size:128;comment:"`
	Status                string      `json:"status" gorm:"column:status;type:string;size:128;comment:"`
	//HostAliases           string    `json:"hostAliases" gorm:"column:host_aliases;type:string;size:128;comment:"`
	ExternalTrafficPolicy string `json:"externalTrafficPolicy" gorm:"column:external_traffic_policy;type:string;size:128;comment:"`
	Owner                 string `json:"owner" gorm:"column:owner;type:string;size:128;comment:负责人"`
}

func (this Service) TableName() string {
	return "caas_service"
}
