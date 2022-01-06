package chart

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

type TagNodeChart struct {
	TagID   int64  `json:"tagID"`
	TagName string `json:"tagName"`
	//GroupPath      string `json:"groupPath"`
	//GroupPathArray string `json:"groupPathArray"`
	NodeCount string `json:"nodeCount"`
}

type TagCpuChart struct {
	TagID   int64  `json:"tagID"`
	TagName string `json:"tagName"`
	//GroupPath      string `json:"groupPath"`
	//GroupPathArray string `json:"groupPathArray"`
	CpuCount string `json:"cpuCount"`
}

// @Summary 按照tag统计tag下的CPU个数
// @Description
// @Produce json
// @Success 200 {object} APIGetChartBarOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/chart/bar [get]
func ChartBar(c *gin.Context) {
	var tnCharts []TagNodeChart
	tx := g.Con().Portal.Debug().Model(node.NodeTagRel{})
	tx = tx.Select("`node_tag_rel`.`tag` as tag_id, `tag`.`name` as tag_name, count(*) as node_count")
	tx = tx.Joins(" left join `tag` on `node_tag_rel`.tag = `tag`.id")
	tx = tx.Group("`node_tag_rel`.tag")
	tx = tx.Having("count(*)>0")
	tx = tx.Order("`node_tag_rel`.`tag`")
	tx = tx.Find(&tnCharts)

	var tcCharts []TagCpuChart
	tx = g.Con().Portal.Debug().Model(node.NodeTagRel{})
	tx = tx.Select("`node_tag_rel`.`tag` as tag_id, `tag`.`name` as tag_name, sum(`node`.`cpu_number`) as cpu_count")
	tx = tx.Joins(" left join `tag` on `node_tag_rel`.tag = `tag`.id")
	tx = tx.Joins(" left join `node` on `node_tag_rel`.node = `node`.id")
	tx = tx.Group("`node_tag_rel`.`tag`")
	tx = tx.Having("count(*)>0")
	tx = tx.Order("`node_tag_rel`.`tag`")
	tx = tx.Find(&tcCharts)

	var tagChartDataMap map[string]int64
	tagChartDataMap = make(map[string]int64)
	for _, tnChart := range tnCharts {
		cnt, _ := strconv.ParseInt(tnChart.NodeCount, 10, 64)
		tagChartDataMap[tnChart.TagName] = cnt
	}

	var tcChartDataMap map[string]int64
	tcChartDataMap = make(map[string]int64)
	for _, tcChart := range tcCharts {
		cnt, _ := strconv.ParseInt(tcChart.CpuCount, 10, 64)
		tcChartDataMap[tcChart.TagName] = cnt
	}

	var xAxisData []string
	var tnChartData []int64
	var tcChartData []int64
	for k, v := range tagChartDataMap {
		xAxisData = append(xAxisData, k)
		tnChartData = append(tnChartData, v)
		tcChartData = append(tcChartData, tcChartDataMap[k])
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
		Data: tnChartData,
	})
	series = append(series, APIGetChartBarSeries{
		Name: "CPU数",
		Type: "bar",
		Data: tcChartData,
	})

	h.JSONR(c, http.StatusOK, APIGetChartBarOutputs{
		Legend: APIGetChartBarLegend{
			Data: []string{"机器数", "CPU数"},
		},
		XAxis:  xAxis,
		Series: series,
	})
	return
}

// @Summary 按照tag统计tag下的CPU个数
// @Description
// @Produce json
// @Success 200 {object} APIGetChartBarOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/chart/pie [get]
func ChartPie(c *gin.Context) {
	var tnCharts []TagNodeChart
	tx := g.Con().Portal.Debug().Model(node.NodeTagRel{})
	tx = tx.Select("`node_tag_rel`.`tag` as tag_id, `tag`.`name` as tag_name, count(*) as node_count")
	tx = tx.Joins(" left join `tag` on `node_tag_rel`.tag = `tag`.id")
	tx = tx.Group("`node_tag_rel`.tag")
	tx = tx.Having("count(*)>0")
	tx = tx.Order("`node_tag_rel`.`tag`")
	tx = tx.Find(&tnCharts)

	var seriesData []APIGetChartPieSeriesItem
	for _, tnChart := range tnCharts {
		data, _ := strconv.ParseInt(tnChart.NodeCount, 10, 64)
		seriesData = append(seriesData, APIGetChartPieSeriesItem{
			//ServiceName:  tnChart.TagName,
			Name:  tnChart.TagName,
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
	TotalNodeCount  int64  `json:"totalNodeCount"`
	UsedNodeCount   int64  `json:"usedNodeCount"`
	UnUsedNodeCount int64  `json:"unUsedNodeCount"`
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
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/chart/vm/stat [get]
func ChartVMStat(c *gin.Context) {
	// 按照物理子系统统计所有机器使用情况
	var totalStat []APIGetChartStat
	tx := g.Con().Portal.Debug().Model(node.Node{})
	tx = tx.Select("`node`.`physical_system` as name, count(*) as total_node_count, sum(`node`.`cpu_number`) as total_cpu_count")
	tx = tx.Group("`node`.`physical_system`")
	tx = tx.Order("`node`.`physical_system`")
	tx = tx.Find(&totalStat)

	var totalCpuCount int64
	var totalNodeCount int64
	var totalChartStat map[string]APIGetChartStat
	totalChartStat = make(map[string]APIGetChartStat)
	for _, s := range totalStat {
		totalCpuCount += s.TotalCpuCount
		totalNodeCount += s.TotalNodeCount
		totalChartStat[s.Name] = s
	}

	var usedStat []APIGetChartStat
	var ids []int64
	tx = g.Con().Portal.Debug().Model(node.Node{})
	tx = tx.Select("distinct(`node`.`id`)")
	tx = tx.Joins("left join `node_tag_rel` on `node`.id=`node_tag_rel`.node")
	tx = tx.Where("`node_tag_rel`.tag is not null;")
	tx = tx.Find(&ids)

	// 查询所有机器
	tx = g.Con().Portal.Debug().Model(node.Node{})
	tx = tx.Select("`node`.`physical_system` as `name`, count(*) as total_node_count, sum(`node`.`cpu_number`) as total_cpu_count")
	tx = tx.Where("id in (?)", ids)
	tx = tx.Group("`node`.`physical_system`")
	tx = tx.Order("`node`.`physical_system`")
	tx = tx.Find(&usedStat)

	var usedCpuCount int64
	var usedNodeCount int64
	var usedChartStat map[string]APIGetChartStat
	usedChartStat = make(map[string]APIGetChartStat)
	for _, s := range usedStat {
		usedCpuCount += s.TotalCpuCount
		usedNodeCount += s.TotalNodeCount
		usedChartStat[s.Name] = s
	}

	var subStat []APIGetChartStat
	for k, v := range totalChartStat {
		t := APIGetChartStat{}
		t.Name = v.Name
		t.TotalNodeCount = v.TotalNodeCount
		t.TotalCpuCount = v.TotalCpuCount
		t.UsedNodeCount = usedChartStat[k].TotalNodeCount
		t.UsedCpuCount = usedChartStat[k].TotalCpuCount
		t.UnUsedNodeCount = v.TotalNodeCount - usedChartStat[k].TotalNodeCount
		t.UnUsedCpuCount = v.TotalCpuCount - usedChartStat[k].TotalCpuCount
		subStat = append(subStat, t)
	}

	var s APIGetChartStatOutputs
	s.Name = ""
	s.TotalNodeCount = totalNodeCount
	s.TotalCpuCount = totalCpuCount
	s.UsedNodeCount = usedNodeCount
	s.UsedCpuCount = usedCpuCount
	s.UnUsedNodeCount = totalNodeCount - usedNodeCount
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
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/chart/container/stat [get]
func ChartContainerStat(c *gin.Context) {
	// 按照物理子系统统计所有机器使用情况
	var totalStat APIGetChartContainerStatOutputs
	tx := g.Con().Portal.Debug()
	tx = tx.Model(&caas.NameSpace{})
	tx = tx.Select("sum(`cpu`) as total_cpu_count, sum(`memory`) as total_memory_count, sum(`shared_volume`) as total_shared_volume, sum(`local_volume`) as total_local_volume")
	tx.Find(&totalStat)

	tx = g.Con().Portal.Debug()
	tx = tx.Model(&caas.Service{})
	tx = tx.Select("sum(`cpu` * now_replicas) / 1000 as used_cpu_count, sum(`memory` * now_replicas) / 1024 as used_memory_count")
	tx.Find(&totalStat)
	h.JSONR(c, http.StatusOK, totalStat)
	return
}
