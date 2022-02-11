package service

var (
	UserService        *userService
	InstService        *instService
	NodeService        *nodeService
	TagService         *tagService
	TagCategoryService *tagCategoryService
	ParamService       *paramService
	BucketService      *bucketService
	CaasService        *caasService
	TemplateService    *templateService
	ProcService        *procService
	WfeService         *wfeService
)

func init() {
	UserService = newUserService()
	InstService = newInstService()
	NodeService = newNodeService()
	TagService = newTagService()
	TagCategoryService = newTagCategoryService()
	ParamService = newParamService()
	BucketService = newBucketService()
	CaasService = newCaasService()
	TemplateService = newTemplateService()
	ProcService = newProcService()
	WfeService = newWfeService()
}
