package model

type APIGetVariableItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type APIGetVariableInputs struct {
	Name string `json:"name"`
}

type APIGetVariableOutputs struct {
	List       []*APIGetVariableItem `json:"list"`
	TotalCount int64                 `json:"totalCount"`
}
