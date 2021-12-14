package service

var (
	HostService        *hostService
	TagService         *tagService
	TagCategoryService *tagCategoryService
	ParamService       *paramService
	BucketService      *bucketService
	CaasService        *caasService
)

func init() {
	HostService = newHostService()
	TagService = newTagService()
	TagCategoryService = newTagCategoryService()
	ParamService = newParamService()
	BucketService = newBucketService()
	CaasService = newCaasService()
}
