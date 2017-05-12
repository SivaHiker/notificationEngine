package structConfigs

type MachinesConfig struct {
	CloudMachines []struct {
		Type          string `json:"type"`
		Consider      bool   `json:"consider"`
		Mqttinstances []struct {
			InstanceIP       string `json:"instance_ip"`
			InstanceID       string `json:"instance_id"`
			InstanceSpec     string `json:"instance_Spec"`
			ServicesDeployed string `json:"services_deployed"`
		} `json:"mqttinstances"`
		Httpinstances []struct {
			InstanceIP       string `json:"instance_ip"`
			InstanceID       string `json:"instance_id"`
			InstanceSpec     string `json:"instance_Spec"`
			ServicesDeployed string `json:"services_deployed"`
		} `json:"httpinstances"`
	} `json:"cloudMachines"`
}
