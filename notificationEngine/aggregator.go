package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/goamz/goamz/cloudwatch"
)

func main() {

	var startTime time.Time
	var endTime time.Time
	var st = "1494533966"
	var et = "1494537566"
	getTotalRecordsCount("sd", "asd")
	if len(os.Args) == 2 {
		startTime = getTimeinSeconds(st)
		endTime = getTimeinSeconds(et)
		fmt.Println(startTime)
		fmt.Println(endTime)
	}
	startTime = getTimeinSeconds(st)
	endTime = getTimeinSeconds(et)
	fmt.Println(startTime)
	fmt.Println(endTime)
	workingDirPath, err := os.Getwd()
	machinesData, err := ioutil.ReadFile(workingDirPath + "/machinesConfig.json")
	attributesData, err := ioutil.ReadFile(workingDirPath + "/assertionConfig.json")
	machines := machinesConfig{}
	attributes := assertionConfig{}
	machinesErr := json.Unmarshal(machinesData, &machines)
	attributesErr := json.Unmarshal(attributesData, &attributes)
	if err != nil || machinesErr != nil || attributesErr != nil {
		log.Fatalln(err)
		log.Fatalln(machinesErr)
		log.Fatalln(attributesErr)
	}
	log.Println(machines)
	var cw = getCloudWatchInstance()
	var statistics = []string{"Average"}
	for i := 0; i < len(machines.CloudMachines); i++ {
		if machines.CloudMachines[i].Consider {
			if machines.CloudMachines[i].Type == "AWS" {
				for j := 0; j < len(machines.CloudMachines[i].Mqttinstances); j++ {
					for k := 0; k < len(attributes.Attributes); k++ {
						dimension := &cloudwatch.Dimension{
							Name:  "InstanceId",
							Value: machines.CloudMachines[i].Mqttinstances[j].InstanceID,
						}
						metricName := attributes.Attributes[k].Name
						Namespace := attributes.Attributes[k].NameSpace
						// for l := 0; l <= len(attributes.Attributes); l++ {
						// 	statistics[l] = attributes.Attributes[l].Value
						// }
						resp := getCloudWatchMetrics(cw, dimension, Namespace, metricName, statistics, startTime, endTime)
						if len(resp.GetMetricStatisticsResult.Datapoints) != 0 {
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Average)
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Maximum)
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Minimum)
						}
					}
				}
				for x := 0; x < len(machines.CloudMachines[i].Httpinstances); x++ {
					for y := 0; y < len(attributes.Attributes); y++ {
						dimension := &cloudwatch.Dimension{
							Name:  "InstanceId",
							Value: machines.CloudMachines[i].Httpinstances[x].InstanceID,
						}
						metricName := attributes.Attributes[y].Name
						Namespace := attributes.Attributes[y].NameSpace
						// for l := 0; l <= len(attributes.Attributes); l++ {
						// 	statistics[l] = attributes.Attributes[l].Value
						// }
						resp := getCloudWatchMetrics(cw, dimension, Namespace, metricName, statistics, startTime, endTime)
						if len(resp.GetMetricStatisticsResult.Datapoints) != 0 {
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Average)
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Maximum)
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Minimum)
						}
					}
				}
			} else if machines.CloudMachines[i].Type == "GCP" {
				// Code to be implemented for GCP
			}
		}
	}
	// fmt.Println(response.GetMetricStatisticsResult.Datapoints[0].Maximum)
}
