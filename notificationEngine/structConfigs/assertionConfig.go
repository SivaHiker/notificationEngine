package main

type assertionConfig struct {
	Attributes []struct {
		Name           string `json:"Name"`
		Value          string `json:"Value"`
		ThresholdValue string `json:"ThresholdValue"`
		NameSpace      string `json:"NameSpace"`
	} `json:"Attributes"`
}
