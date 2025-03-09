package aws

type Result struct {
	Resource Resource
	Error    error
}

type Resource struct {
	TypeConfig *ResourceTypeConfig

	Region   string
	Service  string
	TypeName string
	Id       string
}
