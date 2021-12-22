package app

import (
	"github.com/alexreagan/rabbit/g"
	h "github.com/alexreagan/rabbit/server/helper"
	"github.com/alexreagan/rabbit/server/model/app"
	"github.com/alexreagan/rabbit/server/model/node"
	"github.com/alexreagan/rabbit/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIGetHostGroupTreeInputs struct {
	ID   int64  `json:"id" form:"id"`
	Type string `json:"type" form:"type"`
}

type APIGetHostGroupTreeOutputs struct {
	HostGroups []*node.HostGroup `json:"hostGroups"`
}

type APIGetHostGroupNode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type APIGetHostGroupNodeOutputs struct {
	HostGroups []*APIGetHostGroupNode `json:"hostGroups"`
}

// @Summary 获取host group信息
// @Description
// @Produce json
// @Param id query int64 true "根据ID获取host group树状信息"
// @Success 200 {object} APIGetHostGroupTreeOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tree/children [get]
func TreeChildren(c *gin.Context) {
	var inputs APIGetHostGroupTreeInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	// 根节点
	if inputs.ID == 0 {
		resp, _ := node.HostGroup{}.BuildTree(inputs.ID)
		h.JSONR(c, http.StatusOK, resp)
		return
	}

	// 群组节点
	if inputs.Type == "vmGroup" || inputs.Type == "containerGroup" {
		resp := node.HostGroup{ID: inputs.ID}.GetChildren()
		if len(resp) > 0 {
			// 有子群组，返回子群组
			h.JSONR(c, http.StatusOK, resp)
			return
		} else {
			// 没有子群组，返回群组内的节点
			switch inputs.Type {
			case "vmGroup":
				// 虚拟机类型的群组
				hosts := node.HostGroup{ID: inputs.ID}.RelatedHosts()
				// 转换显示名字
				for _, host := range hosts {
					host.Name = host.IP
				}
				h.JSONR(c, http.StatusOK, hosts)
				return
			case "containerGroup":
				// 容器类型的群组
				hostGroup := node.HostGroup{ID: inputs.ID}
				g.Con().Portal.Where(hostGroup).First(&hostGroup)
				pods := hostGroup.RelatedPods()
				h.JSONR(c, http.StatusOK, pods)
				return
			}
		}
	}

	// 叶子节点(机器)
	h.JSONR(c, http.StatusOK, []string{})
	return
}

// @Summary 重建群组树
// @Description
// @Produce json
// @Success 200 {object} []node.HostGroup
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tree/rebuild [get]
func TreeRebuild(c *gin.Context) {
	resp := service.TagService.ReBuildGraph()
	h.JSONR(c, http.StatusOK, resp)
	return
}

//// @Summary 获取某个节点的所有子节点
//// @Description
//// @Produce json
//// @Param id query int64 true "获取某个节点的所有子节点"
//// @Success 200 {object} APIGetHostGroupTreeOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/host_group/children [get]
//func HostGroupChildren(c *gin.Context) {
//	var inputs APIGetHostGroupTreeInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	resp := node.HostGroup{Tag: inputs.Tag}.GetChildren()
//	h.JSONR(c, http.StatusOK, resp)
//	return
//}
//
//type APIGetHostGroupHostsInputs struct {
//	ID int64 `json:"id" form:"id"`
//}
//
//// @Summary 获取某个节点的所有机器
//// @Description
//// @Produce json
//// @Param id query int64 true "获取某个节点的所有机器"
//// @Success 200 {object} APIGetHostGroupTreeOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/host_group/hosts [get]
//func HostGroupHosts(c *gin.Context) {
//	var inputs APIGetHostGroupHostsInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	var resp []*node.Host
//	hostGroup := &node.HostGroup{
//		Tag: inputs.ID,
//	}
//	dt := g.Con().Portal.Table(hostGroup.TableName()).Where(hostGroup).Find(&hostGroup)
//	if dt.Error != nil {
//		h.JSONR(c, h.ExpecStatus, dt.Error)
//		return
//	}
//	hosts := hostGroup.RelatedHosts()
//	for _, host := range hosts {
//		host.ServiceName = host.IP
//		host.IsWarning = host.MeetWarningCondition()
//		resp = append(resp, host)
//	}
//	h.JSONR(c, http.StatusOK, resp)
//	return
//}
//
//type APIGetHostGroupMoveInputs struct {
//	ID       int64  `json:"id" form:"id"`
//	ServiceName     string `json:"name" form:"name"`
//	ParentID int64  `json:"parentID" form:"parentID"`
//}
//
//// @Summary 将节点ID移动到parent，作为parent的父节点
//// @Description
//// @Produce json
//// @Param APIGetHostGroupMoveInputs query APIGetHostGroupMoveInputs true "获取某个节点的所有机器"
//// @Success 200 {object} APIGetHostGroupTreeOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/host_group/move [post]
//func HostGroupMove(c *gin.Context) {
//	var inputs APIGetHostGroupMoveInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	db := g.Con().Portal.Table(node.HostGroup{}.TableName())
//	db.Where(node.HostGroup{Tag: inputs.ID}).Updates(node.HostGroup{ParentID: inputs.ParentID})
//	node.HostGroup{}.ReBuildGraph()
//
//	h.JSONR(c, http.StatusOK, "ok")
//	return
//}

type APIPostHostGroupDraggingInputs struct {
	DraggingNodeID   int64  `json:"draggingNodeID" form:"draggingNodeID"`
	DraggingNodeType string `json:"draggingNodeType" form:"draggingNodeType"`
	DropNodeID       int64  `json:"dropNodeID" form:"dropNodeID"`
	DropNodeType     string `json:"dropNodeType" form:"dropNodeType"`
}

//// @Summary 拖动节点
//// @Description
//// @Produce json
//// @Param APIGetHostGroupMoveInputs query APIGetHostGroupMoveInputs true "获取某个节点的所有机器"
//// @Success 200 {object} APIGetHostGroupTreeOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
//// @Router /api/v1/tree/dragging [post]
//func TreeDragging(c *gin.Context) {
//	var inputs APIPostHostGroupDraggingInputs
//
//	if err := c.Bind(&inputs); err != nil {
//		h.JSONR(c, h.BadStatus, err)
//		return
//	}
//
//	switch inputs.DraggingNodeType {
//	case "":
//		// 拖动虚拟机
//		if inputs.DropNodeType != "" {
//			tx := g.Con().Portal.Model(node.HostGroupRel{}).Begin()
//			if dt := tx.Debug().Where(node.HostGroupRel{HostID: inputs.DraggingNodeID}).Delete(node.HostGroupRel{}); dt.Error != nil {
//				h.JSONR(c, h.ExpecStatus, dt.Error)
//				dt.Rollback()
//				return
//			}
//			if dt := tx.Debug().Create(&node.HostGroupRel{
//				HostID:  inputs.DraggingNodeID,
//				GroupID: inputs.DropNodeID,
//			}); dt.Error != nil {
//				h.JSONR(c, h.ExpecStatus, dt.Error)
//				dt.Rollback()
//				return
//			}
//			tx.Commit()
//		}
//	case "vmGroup":
//		// 拖动群组
//		if inputs.DropNodeType != "" {
//			db := g.Con().Portal.Model(node.HostGroup{}).Debug()
//			db.Where(node.HostGroup{Tag: inputs.DraggingNodeID}).Updates(node.HostGroup{ParentID: inputs.DropNodeID})
//		}
//	case "containerGroup":
//		// 拖动群组
//		if inputs.DropNodeType != "" {
//			db := g.Con().Portal.Model(node.HostGroup{}).Debug()
//			db.Where(node.HostGroup{Tag: inputs.DraggingNodeID}).Updates(node.HostGroup{ParentID: inputs.DropNodeID})
//		}
//	}
//
//	resp, _ := node.HostGroup{}.ReBuildTree()
//
//	h.JSONR(c, http.StatusOK, resp)
//	return
//}

type APIGetV2TreeInputs struct {
	CategoryIDs []int64 `json:"categoryIDs[]" form:"categoryIDs[]"`
	TagIDs      []int64 `json:"tagIDs[]" form:"tagIDs[]"`
}

// @Summary 根据tags获取tags下所有的机器
// @Description
// @Produce json
// @Param id query int64 true "根据tags获取tags下所有的机器"
// @Success 200 {object} APIGetHostGroupTreeOutputs
//// @Failure 400 "bad arguments"
//// @Failure 417 "internal error"
// @Router /api/v2/tree/children [get]
func V2TreeChildren(c *gin.Context) {
	var inputs APIGetV2TreeInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	var method = "1"
	switch method {
	case "1":
		////////////// method 1: 路由图
		globalTagGraphNode := service.TagService.GlobalTagGraphNodeV2()
		if globalTagGraphNode == nil {
			globalTagGraphNode = service.TagService.BuildGraphV2()
		}
		// 根节点
		if len(inputs.TagIDs) == 0 {
			h.JSONR(c, http.StatusOK, globalTagGraphNode.Children)
			return
		}

		// 其他节点
		// 找到inputs的末级节点
		graphNode := globalTagGraphNode
		for _, id := range inputs.TagIDs {
			graphNode = graphNode.Next[id]
		}

		h.JSONR(c, h.OKStatus, graphNode.Children)
		return
		//下一个标签类型
		//categoryNames, err := service.ParamService.GetTreeOrder()
		//if err != nil {
		//	h.JSONR(c, h.ExpecStatus, err)
		//	return
		//}
		//if len(inputs.CategoryIDs) < len(categoryNames) {
		//
		//	// 末级节点的子节点
		//	var resp []interface{}
		//	for _, x := range graphNode.Nexts() {
		//		resp = append(resp, x)
		//	}
		//
		//	for _, x := range graphNode.UnTaggedHosts {
		//		resp = append(resp, x)
		//	}
		//
		//	for _, x := range graphNode.UnTaggedPods {
		//		resp = append(resp, x)
		//	}
		//
		//	h.JSONR(c, h.OKStatus, graphNode.Children)
		//	return
		//} else {
		//	var resp []interface{}
		//	for _, x := range graphNode.RelatedHosts {
		//		resp = append(resp, x)
		//	}
		//	for _, x := range graphNode.RelatedPods {
		//		resp = append(resp, x)
		//	}
		//	h.JSONR(c, h.OKStatus, resp)
		//	return
		//}
	case "2":
		//////////// method 2: 顺序遍历
		// 当前标签顺序下的所有机器
		var hosts []*node.Host
		if len(inputs.TagIDs) > 0 {
			hosts = service.HostService.HostsHavingTagIDs(inputs.TagIDs)
		} else {
			hosts = service.HostService.HostsRelatedTags()
		}

		// 下一个标签类型
		categoryNames, err := service.ParamService.GetTreeOrder()
		if err != nil {
			h.JSONR(c, h.ExpecStatus, err)
			return
		}
		if len(inputs.CategoryIDs) < len(categoryNames) {
			// 没有取全category，取下一个category
			nextCategoryName := categoryNames[len(inputs.CategoryIDs)]
			nextCategory := service.TagCategoryService.GetByName(nextCategoryName)
			tags := service.TagCategoryService.GetTagsByCategory(nextCategory)

			// 分桶
			tagMap, untaggedHosts := service.BucketService.Sort(hosts, tags)
			var nodeTags app.Tags
			for _, x := range tagMap {
				nodeTags = append(nodeTags, x)
			}
			nodeTags.Sort()

			var resp []interface{}
			for _, x := range nodeTags {
				resp = append(resp, x)
			}

			for _, x := range untaggedHosts {
				resp = append(resp, x)
			}
			h.JSONR(c, h.OKStatus, resp)
			return
		} else {
			// 已经拿全category，返回所有机器
			h.JSONR(c, http.StatusOK, hosts)
			return
		}
	}
}
