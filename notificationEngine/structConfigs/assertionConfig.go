package structConfigs

type AssertionConfig struct {
	Attributes []struct {
		Name           string `json:"Name"`
		Value          string `json:"Value"`
		ThresholdValue string `json:"ThresholdValue"`
		NameSpace      string `json:"NameSpace"`
		DimensionName  string `json:"dimensionName"`
	} `json:"Attributes"`
}
