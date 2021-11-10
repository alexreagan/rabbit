package node

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/caas"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type APIGetChartBarXAxis struct {
	Type string   `json:"type"`
	Data []string `json:"data"`
}

type APIGetChartBarSeries struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	Data []int64 `json:"data"`
}

type APIGetChartBarLegend struct {
	Data []string `json:"data"`
}

type APIGetChartBarOutputs struct {
	Legend APIGetChartBarLegend   `json:"legend"`
	XAxis  []APIGetChartBarXAxis  `json:"xAxis"`
	Series []APIGetChartBarSeries `json:"series"`
}

type APIGetChartPieSeriesItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type APIGetChartPieSeries struct {
	Name string                     `json:"name"`
	Data []APIGetChartPieSeriesItem `json:"data"`
}

type APIGetChartPieOutputs struct {
	Title  string               `json:"title"`
	Series APIGetChartPieSeries `json:"series"`
}

type GroupNodeChart struct {
	GroupId        int64  `json:"groupId"`
	GroupName      string `json:"groupName"`
	GroupPath      string `json:"groupPath"`
	GroupPathArray string `json:"groupPathArray"`
	NodeCount      string `json:"nodeCount"`
}

type GroupCpuChart struct {
	GroupId        int64  `json:"groupId"`
	GroupName      string `json:"groupName"`
	GroupPath      string `json:"groupPath"`
	GroupPathArray string `json:"groupPathArray"`
	CpuCount       string `json:"cpuCount"`
}

//func ArrayStringJoin(arr string) string {
//	var path []string
//	json.Unmarshal([]byte(arr), &path)
//	return strings.Join(path, "/")
//}

// @Summary 按照host_group统计host_group下的CPU个数
// @Description
// @Produce json
// @Success 200 {object} APIGetChartBarOutputs
// @Failure 400 {object} APIGetChartBarOutputs
// @Router /api/v1/chart/bar [get]
func ChartBar(c *gin.Context) {
	var gnCharts []GroupNodeChart
	db := g.Con().Portal.Debug().Table(node.HostGroupRel{}.TableName())
	db = db.Select("`host_group_rel`.`group_id`, `host_group`.`name` as group_name, " +
		"`host_group`.`path` as group_path, `host_group`.`path_array` as group_path_array, count(*) as node_count")
	db = db.Joins(" left join `host_group` on `host_group_rel`.group_id = `host_group`.id")
	db = db.Group("`host_group_rel`.group_id")
	db = db.Having("count(*)>0")
	db = db.Order("`host_group_rel`.`group_id`")
	db = db.Find(&gnCharts)

	var gcCharts []GroupCpuChart
	db = g.Con().Portal.Debug().Table(node.HostGroupRel{}.TableName())
	db = db.Select("`host_group_rel`.`group_id`, `host_group`.`name` as group_name, " +
		"`host_group`.`path_array` as group_path_array, `host_group`.`path` as group_path, " +
		"sum(`host`.`cpu_number`) as cpu_count")
	db = db.Joins(" left join `host_group` on `host_group_rel`.group_id = `host_group`.id")
	db = db.Joins(" left join `host` on `host_group_rel`.host_id = `host`.id")
	db = db.Group("`host_group_rel`.group_id")
	db = db.Having("count(*)>0")
	db = db.Order("`host_group_rel`.`group_id`")
	db = db.Find(&gcCharts)

	var gnChartDataMap map[string]int64
	gnChartDataMap = make(map[string]int64)
	for _, gnChart := range gnCharts {
		cnt, _ := strconv.ParseInt(gnChart.NodeCount, 10, 64)
		gnChartDataMap[gnChart.GroupPath] = cnt
	}

	var gcChartDataMap map[string]int64
	gcChartDataMap = make(map[string]int64)
	for _, gcChart := range gcCharts {
		cnt, _ := strconv.ParseInt(gcChart.CpuCount, 10, 64)
		gcChartDataMap[gcChart.GroupPath] = cnt
	}

	var xAxisData []string
	var gnChartData []int64
	var gcChartData []int64
	for k, v := range gnChartDataMap {
		xAxisData = append(xAxisData, k)
		gnChartData = append(gnChartData, v)
		gcChartData = append(gcChartData, gcChartDataMap[k])
	}

	var xAxis []APIGetChartBarXAxis
	xAxis = append(xAxis, APIGetChartBarXAxis{
		Type: "category",
		//BoundaryGap: false,
		Data: xAxisData,
	})

	var series []APIGetChartBarSeries
	series = append(series, APIGetChartBarSeries{
		Name: "机器数",
		Type: "bar",
		Data: gnChartData,
	})
	series = append(series, APIGetChartBarSeries{
		Name: "CPU数",
		Type: "bar",
		Data: gcChartData,
	})

	h.JSONR(c, http.StatusOK, APIGetChartBarOutputs{
		Legend: APIGetChartBarLegend{
			Data: []string{"机器数", "CPU数"},
		},
		XAxis:  xAxis,
		Series: series,
	})
	return

	//var legend []string
	//var xAxis []APIGetChartBarXAxis
	//var series []APIGetChartBarSeries
	//for _, gnChart := range gnCharts {
	//	legend = append(legend, gnChart.GroupName)
	//	data, _ := strconv.ParseInt(gnChart.NodeCount, 10, 64)
	//	series = append(series, APIGetChartBarSeries{
	//		Name: gnChart.GroupName,
	//		Type: "bar",
	//		Data: []int64{data},
	//	})
	//}
	//
	//xAxis = append(xAxis, APIGetChartBarXAxis{
	//	Type: "category",
	//	BoundaryGap: false,
	//	Data: []string{""},
	//})
	//
	//h.JSONR(c, http.StatusOK, APIGetChartBarOutputs{
	//	Legend: APIGetChartBarLegend{
	//		Data: legend,
	//	},
	//	XAxis: xAxis,
	//	Series: series,
	//})
	return
}

// @Summary 按照host_group统计host_group下的CPU个数
// @Description
// @Produce json
// @Success 200 {object} APIGetChartBarOutputs
// @Failure 400 {object} APIGetChartBarOutputs
// @Router /api/v1/chart/pie [get]
func ChartPie(c *gin.Context) {
	var gnCharts []GroupNodeChart
	db := g.Con().Portal.Debug().Table(node.HostGroupRel{}.TableName())
	db = db.Select("`host_group_rel`.`group_id`, `host_group`.`name` as group_name, " +
		"`host_group`.`path` as group_path, `host_group`.`path_array` as group_path_array, " +
		"count(*) as node_count")
	db = db.Joins(" left join `host_group` on `host_group_rel`.group_id = `host_group`.id")
	db = db.Group("`host_group_rel`.group_id")
	db = db.Having("count(*)>0")
	db = db.Order("`host_group_rel`.`group_id`")
	db = db.Find(&gnCharts)

	var seriesData []APIGetChartPieSeriesItem
	for _, gnChart := range gnCharts {
		data, _ := strconv.ParseInt(gnChart.NodeCount, 10, 64)
		seriesData = append(seriesData, APIGetChartPieSeriesItem{
			//Name:  gnChart.GroupName,
			Name:  gnChart.GroupPath,
			Value: data,
		})
	}

	h.JSONR(c, http.StatusOK, APIGetChartPieOutputs{
		Title: "机器分组饼状图",
		Series: APIGetChartPieSeries{
			Name: "机器数",
			Data: seriesData,
		},
	})
	return
}

type APIGetChartStat struct {
	Name            string `json:"name"`
	TotalHostCount  int64  `json:"totalHostCount"`
	UsedHostCount   int64  `json:"usedHostCount"`
	UnUsedHostCount int64  `json:"unUsedHostCount"`
	TotalCpuCount   int64  `json:"totalCpuCount"`
	UsedCpuCount    int64  `json:"usedCpuCount"`
	UnUsedCpuCount  int64  `json:"unUsedCpuCount"`
}

type APIGetChartStatOutputs struct {
	APIGetChartStat
	SubStat []APIGetChartStat `json:"subStat"`
}

// @Summary 统计平台收纳的所有机器数/cpu数/已分配机器数/已分配cpu数
// @Description
// @Produce json
// @Success 200 {object} APIGetChartStatOutputs
// @Failure 400 {object} APIGetChartStatOutputs
// @Router /api/v1/chart/vm/stat [get]
func ChartVMStat(c *gin.Context) {
	// 按照物理子系统统计所有机器使用情况
	var totalStat []APIGetChartStat
	db := g.Con().Portal.Debug().Model(node.Host{})
	db = db.Select("`host`.`physical_system` as name, count(*) as total_host_count, sum(`host`.`cpu_number`) as total_cpu_count")
	db = db.Group("`host`.`physical_system`")
	db = db.Order("`host`.`physical_system`")
	db = db.Find(&totalStat)

	var totalCpuCount int64
	var totalHostCount int64
	var totalChartStat map[string]APIGetChartStat
	totalChartStat = make(map[string]APIGetChartStat)
	for _, s := range totalStat {
		totalCpuCount += s.TotalCpuCount
		totalHostCount += s.TotalHostCount
		totalChartStat[s.Name] = s
	}

	var usedStat []APIGetChartStat
	// 已关联到组的机器
	//var subStatOutputs []APIGetChartStatOutputs
	var ids []int64
	db = g.Con().Portal.Debug().Model(node.Host{})
	db = db.Select("distinct(`host`.`id`)")
	db = db.Joins("left join `host_group_rel` on `host`.id=`host_group_rel`.host_id")
	db = db.Where("`host_group_rel`.group_id is not null;")
	db = db.Find(&ids)

	// 查询所有机器
	db = g.Con().Portal.Debug().Model(node.Host{})
	db = db.Select("`host`.`physical_system` as `name`, count(*) as total_host_count, sum(`host`.`cpu_number`) as total_cpu_count")
	db = db.Where("id in (?)", ids)
	db = db.Group("`host`.`physical_system`")
	db = db.Order("`host`.`physical_system`")
	db = db.Find(&usedStat)

	var usedCpuCount int64
	var usedHostCount int64
	var usedChartStat map[string]APIGetChartStat
	usedChartStat = make(map[string]APIGetChartStat)
	for _, s := range usedStat {
		usedCpuCount += s.TotalCpuCount
		usedHostCount += s.TotalHostCount
		usedChartStat[s.Name] = s
	}

	var subStat []APIGetChartStat
	for k, v := range totalChartStat {
		t := APIGetChartStat{}
		t.Name = v.Name
		t.TotalHostCount = v.TotalHostCount
		t.TotalCpuCount = v.TotalCpuCount
		t.UsedHostCount = usedChartStat[k].TotalHostCount
		t.UsedCpuCount = usedChartStat[k].TotalCpuCount
		t.UnUsedHostCount = v.TotalHostCount - usedChartStat[k].TotalHostCount
		t.UnUsedCpuCount = v.TotalCpuCount - usedChartStat[k].TotalCpuCount
		subStat = append(subStat, t)
	}

	var s APIGetChartStatOutputs
	s.Name = ""
	s.TotalHostCount = totalHostCount
	s.TotalCpuCount = totalCpuCount
	s.UsedHostCount = usedHostCount
	s.UsedCpuCount = usedCpuCount
	s.UnUsedHostCount = totalHostCount - usedHostCount
	s.UnUsedCpuCount = totalCpuCount - usedCpuCount
	s.SubStat = subStat

	h.JSONR(c, http.StatusOK, s)
	return
}

type APIGetChartContainerStatOutputs struct {
	TotalCpuCount     float64 `json:"totalCpuCount"`
	TotalMemoryCount  float64 `json:"totalMemoryCount"`
	TotalSharedVolume int64   `json:"totalSharedVolume"`
	TotalLocalVolume  int64   `json:"totalLocalVolume"`
	UsedCpuCount      float64 `json:"usedCpuCount"`
	UsedMemoryCount   float64 `json:"usedMemoryCount"`
}

// @Summary 统计平台收纳的容器使用量/cpu数/已分配机器数/已分配cpu数
// @Description
// @Produce json
// @Success 200 {object} APIGetChartStatOutputs
// @Failure 400 {object} APIGetChartStatOutputs
// @Router /api/v1/chart/container/stat [get]
func ChartContainerStat(c *gin.Context) {
	// 按照物理子系统统计所有机器使用情况
	var totalStat APIGetChartContainerStatOutputs
	db := g.Con().Portal.Debug()
	db = db.Model(&caas.NameSpace{})
	db = db.Select("sum(`cpu`) as total_cpu_count, sum(`memory`) as total_memory_count, sum(`shared_volume`) as total_shared_volume, sum(`local_volume`) as total_local_volume")
	db.Find(&totalStat)

	tx := g.Con().Portal.Debug()
	tx = tx.Model(&caas.Service{})
	tx = tx.Select("sum(`cpu` * now_replicas) / 1000 as used_cpu_count, sum(`memory` * now_replicas) / 1024 as used_memory_count")
	tx.Find(&totalStat)
	h.JSONR(c, http.StatusOK, totalStat)
	return
}
