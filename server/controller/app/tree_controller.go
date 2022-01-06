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

type APIGetNodeGroupTreeInputs struct {
	ID   int64  `json:"id" form:"id"`
	Type string `json:"type" form:"type"`
}

type APIGetNodeGroupTreeOutputs struct {
	NodeGroups []*node.NodeGroup `json:"nodeGroups"`
}

type APIGetNodeGroupNode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type APIGetNodeGroupNodeOutputs struct {
	NodeGroups []*APIGetNodeGroupNode `json:"nodeGroups"`
}

// @Summary 获取node children信息
// @Description
// @Produce json
// @Param id query int64 true "根据ID获取node children信息"
// @Success 200 {object} APIGetNodeGroupTreeOutputs
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tree/children [get]
func TreeChildren(c *gin.Context) {
	var inputs APIGetNodeGroupTreeInputs

	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, h.BadStatus, err)
		return
	}

	// 根节点
	if inputs.ID == 0 {
		resp, _ := node.NodeGroup{}.BuildTree(inputs.ID)
		h.JSONR(c, http.StatusOK, resp)
		return
	}

	// 群组节点
	if inputs.Type == "vmGroup" || inputs.Type == "containerGroup" {
		resp := node.NodeGroup{ID: inputs.ID}.GetChildren()
		if len(resp) > 0 {
			// 有子群组，返回子群组
			h.JSONR(c, http.StatusOK, resp)
			return
		} else {
			// 没有子群组，返回群组内的节点
			switch inputs.Type {
			case "vmGroup":
				// 虚拟机类型的群组
				nodes := node.NodeGroup{ID: inputs.ID}.RelatedNodes()
				// 转换显示名字
				for _, n := range nodes {
					n.Name = n.IP
				}
				h.JSONR(c, http.StatusOK, nodes)
				return
			case "containerGroup":
				// 容器类型的群组
				nodeGroup := node.NodeGroup{ID: inputs.ID}
				g.Con().Portal.Where(nodeGroup).First(&nodeGroup)
				pods := nodeGroup.RelatedPods()
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
// @Success 200 {object} []node.NodeGroup
// @Failure 400 "bad arguments"
// @Failure 417 "internal error"
// @Router /api/v1/tree/rebuild [get]
func TreeRebuild(c *gin.Context) {
	resp := service.TagService.ReBuildGraphV2()
	h.JSONR(c, http.StatusOK, resp)
	return
}

type APIPostNodeGroupDraggingInputs struct {
	DraggingNodeID   int64  `json:"draggingNodeID" form:"draggingNodeID"`
	DraggingNodeType string `json:"draggingNodeType" form:"draggingNodeType"`
	DropNodeID       int64  `json:"dropNodeID" form:"dropNodeID"`
	DropNodeType     string `json:"dropNodeType" form:"dropNodeType"`
}

type APIGetV2TreeInputs struct {
	CategoryIDs []int64 `json:"categoryIDs[]" form:"categoryIDs[]"`
	TagIDs      []int64 `json:"tagIDs[]" form:"tagIDs[]"`
}

// @Summary 根据tags获取tags下所有的机器
// @Description
// @Produce json
// @Param id query int64 true "根据tags获取tags下所有的机器"
// @Success 200 {object} APIGetNodeGroupTreeOutputs
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
			h.JSONR(c, http.StatusOK, service.BuildChildrenInformation(globalTagGraphNode))
			return
		}

		// 其他节点
		// 找到inputs的末级节点
		graphNode := globalTagGraphNode
		for _, id := range inputs.TagIDs {
			graphNode = graphNode.Next[id]
		}

		h.JSONR(c, h.OKStatus, service.BuildChildrenInformation(graphNode))
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
		//	for _, x := range graphNode.UnTaggedNodes {
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
		//	for _, x := range graphNode.RelatedNodes {
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
		var nodes []*node.Node
		if len(inputs.TagIDs) > 0 {
			nodes = service.NodeService.NodesHavingTagIDs(inputs.TagIDs)
		} else {
			nodes = service.NodeService.NodesRelatedTags()
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
			tagMap, unTaggedNodes := service.BucketService.Sort(nodes, tags)
			var nodeTags app.Tags
			for _, x := range tagMap {
				nodeTags = append(nodeTags, x)
			}
			nodeTags.Sort()

			var resp []interface{}
			for _, x := range nodeTags {
				resp = append(resp, x)
			}

			for _, x := range unTaggedNodes {
				resp = append(resp, x)
			}
			h.JSONR(c, h.OKStatus, resp)
			return
		} else {
			// 已经拿全category，返回所有机器
			h.JSONR(c, http.StatusOK, nodes)
			return
		}
	}
}
