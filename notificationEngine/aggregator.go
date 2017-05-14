package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goamz/goamz/cloudwatch"
	conf "github.com/sivahiker/notificationEngine/structConfigs"
	"github.com/sivahiker/notificationEngine/utils"
)

func main() {

	var startTime time.Time
	var endTime time.Time
	// var st = "2017-05-12T09:10:33"
	// var et = "2017-05-12T09:15:36"
	var runid = utils.GetTotalRecordsCount() + 1
	if len(os.Args) == 3 {
		startTime = utils.GetTimeinSeconds(os.Args[1])
		endTime = utils.GetTimeinSeconds(os.Args[2])
		fmt.Println(startTime)
		fmt.Println(endTime)
	}
	// startTime = utils.GetTimeinSeconds(st)
	// endTime = utils.GetTimeinSeconds(et)
	// fmt.Println(startTime)
	// fmt.Println(endTime)
	workingDirPath, err := os.Getwd()
	machinesData, err := ioutil.ReadFile(workingDirPath + "/resources/machinesConfig.json")
	attributesData, err := ioutil.ReadFile(workingDirPath + "/resources/assertionConfig.json")
	machines := conf.MachinesConfig{}
	attributes := conf.AssertionConfig{}
	machinesErr := json.Unmarshal(machinesData, &machines)
	attributesErr := json.Unmarshal(attributesData, &attributes)
	if err != nil || machinesErr != nil || attributesErr != nil {
		log.Fatalln(err)
		log.Fatalln(machinesErr)
		log.Fatalln(attributesErr)
	}
	log.Println(machines)
	var cw = utils.GetCloudWatchInstance()
	var statistics = []string{"Average", "Maximum", "Minimum"}
	var subrunid = 0
	for i := 0; i < len(machines.CloudMachines); i++ {
		if machines.CloudMachines[i].Consider {
			if machines.CloudMachines[i].Type == "AWS" {
				for j := 0; j < len(machines.CloudMachines[i].Mqttinstances); j++ {
					subrunid++
					utils.WriteDBRecordsIntoMachineTable(runid, subrunid, machines.CloudMachines[i].Mqttinstances[j].InstanceIP,
						machines.CloudMachines[i].Mqttinstances[j].InstanceSpec, "MQTT")
					for k := 0; k < len(attributes.Attributes); k++ {
						dimension := &cloudwatch.Dimension{
							Name:  attributes.Attributes[k].DimensionName,
							Value: machines.CloudMachines[i].Mqttinstances[j].InstanceID,
						}
						metricName := attributes.Attributes[k].Name
						Namespace := attributes.Attributes[k].NameSpace
						resp := utils.GetCloudWatchMetrics(cw, dimension, Namespace, metricName, statistics, startTime, endTime)
						if len(resp.GetMetricStatisticsResult.Datapoints) > 0 {
							utils.UpdateRecordsIntoDB(attributes.Attributes[k].Name, runid, subrunid, strconv.FormatFloat(resp.GetMetricStatisticsResult.Datapoints[0].Average, 'f', -1, 64), strconv.FormatFloat((resp.GetMetricStatisticsResult.Datapoints[0].Maximum), 'f', -1, 64), strconv.FormatFloat(resp.GetMetricStatisticsResult.Datapoints[0].Minimum, 'f', -1, 64))
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Average)
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Maximum)
							fmt.Println(resp.GetMetricStatisticsResult.Datapoints[0].Minimum)
						}
					}
				}
				for x := 0; x < len(machines.CloudMachines[i].Httpinstances); x++ {
					subrunid++
					utils.WriteDBRecordsIntoMachineTable(runid, subrunid, machines.CloudMachines[i].Httpinstances[x].InstanceIP,
						machines.CloudMachines[i].Httpinstances[x].InstanceSpec, "HTTP")
					for y := 0; y < len(attributes.Attributes); y++ {
						dimension := &cloudwatch.Dimension{
							Name:  attributes.Attributes[y].DimensionName,
							Value: machines.CloudMachines[i].Httpinstances[x].InstanceID,
						}
						metricName := attributes.Attributes[y].Name
						Namespace := attributes.Attributes[y].NameSpace
						resp := utils.GetCloudWatchMetrics(cw, dimension, Namespace, metricName, statistics, startTime, endTime)
						if len(resp.GetMetricStatisticsResult.Datapoints) > 0 {
							utils.UpdateRecordsIntoDB(attributes.Attributes[y].Name, runid, subrunid, strconv.FormatFloat(resp.GetMetricStatisticsResult.Datapoints[0].Average, 'f', -1, 64), strconv.FormatFloat((resp.GetMetricStatisticsResult.Datapoints[0].Maximum), 'f', -1, 64), strconv.FormatFloat(resp.GetMetricStatisticsResult.Datapoints[0].Minimum, 'f', -1, 64))
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
