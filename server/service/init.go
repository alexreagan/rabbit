package service

var (
	HostService        *hostService
	TagService         *tagService
	TagCategoryService *tagCategoryService
	ParamService       *paramService
	BucketService      *bucketService
	CaasService        *caasService
	TemplateService    *templateService
)

func init() {
	HostService = newHostService()
	TagService = newTagService()
	TagCategoryService = newTagCategoryService()
	ParamService = newParamService()
	BucketService = newBucketService()
	CaasService = newCaasService()
	TemplateService = newTemplateService()
}
