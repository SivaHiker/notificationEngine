package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goamz/goamz/cloudwatch"
	conf "github.com/sivahiker/notificationEngine/structConfigs"
	"github.com/sivahiker/notificationEngine/utils"
)

func main() {

	var startTime time.Time
	var endTime time.Time
	var baseLineID int
	var loadRunType string

	// Getting the max Run Id from the DB, inorder to generate the new run id for the current Run.
	var runid = utils.GetTotalRecordsCount() + 1

	if len(os.Args) == 5 {
		startTime = utils.GetTimeinSeconds(os.Args[1])
		endTime = utils.GetTimeinSeconds(os.Args[2])
		baseLineRunID, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalln(err)
		}
		baseLineID = baseLineRunID
		loadRunType = os.Args[4]
		fmt.Println(startTime)
		fmt.Println(endTime)
	} else if len(os.Args) == 2 {
		if "--help" == os.Args[1] {
			fmt.Println("The command has the following options ./notificationEngine <startTime> <endTime> <baseLineRunId> <loadRunType> \n " +
				"startTime  = Send in UTC Time Format withing double quotes \n endTime = Send in UTC Time Format withing double quotes \n " +
				"baseLineId = This is the runId, which is the primary key in the QA DB, which has the perf readings. This will be decided " +
				"by the Infra Engg Team.Should be given in double quotes. \n loadRunType = Either HTTP or MQTT in double quotes")
			log.Fatalln("Please try again after going through the help options")
		} else {
			fmt.Println("Please use the following command for help \n ./notificationEngine --help")
			log.Fatalln("Please try again after going through the help options")
		}
	} else if len(os.Args) == 1 {
		fmt.Println("Please use the following command for help \n ./notificationEngine --help")
		log.Fatalln("Please try again after going through the help options")
	}
	// Below Commented Code Lines were written for testing purposes in local box.
	// var st = "2017-05-12T09:10:33"
	// var et = "2017-05-12T09:15:36"
	// startTime = utils.GetTimeinSeconds(st)
	// endTime = utils.GetTimeinSeconds(et)
	// baseLineID = 17
	// fmt.Println(startTime)
	// fmt.Println(endTime)

	// Coding Logic Starts Here !!!
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
						machines.CloudMachines[i].Mqttinstances[j].InstanceSpec, "MQTT", "AWS")
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
						machines.CloudMachines[i].Httpinstances[x].InstanceSpec, "HTTP", "AWS")
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
	// Reporter Part For Time Being. Eventually we will end up in writing a golang service for the same.
	if strings.ToLower(loadRunType) == "mqtt" {
		utils.SendMail(startTime, endTime, machines, attributes, "MQTT", runid, "AWS", baseLineID)
	} else {
		utils.SendMail(startTime, endTime, machines, attributes, "HTTP", runid, "AWS", baseLineID)
	}

}
