package service

var (
	NodeService        *nodeService
	TagService         *tagService
	TagCategoryService *tagCategoryService
	ParamService       *paramService
	BucketService      *bucketService
	CaasService        *caasService
	TemplateService    *templateService
	ProcService *procService
)

func init() {
	NodeService = newNodeService()
	TagService = newTagService()
	TagCategoryService = newTagCategoryService()
	ParamService = newParamService()
	BucketService = newBucketService()
	CaasService = newCaasService()
	TemplateService = newTemplateService()
	ProcService = newProcService()
}
