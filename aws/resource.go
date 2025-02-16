package aws

type Result struct {
	Resource Resource
	Error    error
}

type Resource struct {
	Region   string
	Service  string
	TypeName string
	Id       string
}
