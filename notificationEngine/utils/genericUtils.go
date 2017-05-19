package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/sivahiker/notificationEngine/structConfigs"
)

var buildStatusFlag = true

func GetTimeinSeconds(seconds string) time.Time {
	// i, err := strconv.ParseInt(seconds, 10, 64)
	// if err != nil {
	// panic(err)
	// }
	// tm := time.Unix(i, 0)
	// fmt.Println(tm)
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, seconds)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	return t
}

func SendMail(startTime time.Time, endTime time.Time, machines structConfigs.MachinesConfig, attributes structConfigs.AssertionConfig, runType string, runId int, cloudType string, previousRunid int) {
	currentWD, err := os.Getwd()
	var mailData []byte
	if "MQTT" == runType {
		data, err := ioutil.ReadFile(currentWD + "/resources/mailTemplate-MQTT.html")
		if err != nil {
			fmt.Println(err)
		}
		mailData = data
	} else {
		data, err := ioutil.ReadFile(currentWD + "/resources/mailTemplate-HTTP.html")
		if err != nil {
			fmt.Println(err)
		}
		mailData = data
	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(mailData))
	var dataInString = string(mailData)
	var startDurationInString = dataTimeForGraphs(startTime)
	var endDurationInString = dataTimeForGraphs(endTime)
	dataInString = strings.Replace(dataInString, "${TestRunType}", runType, -1)
	dataInString = strings.Replace(dataInString, "${environment}", "PreProd", -1)
	dataInString = strings.Replace(dataInString, "${starttime}", startTime.String(), -1)
	dataInString = strings.Replace(dataInString, "${endtime}", endTime.String(), -1)
	dataInString = strings.Replace(dataInString, "${startDuration}", startDurationInString, -1)
	dataInString = strings.Replace(dataInString, "${endDuration}", endDurationInString, -1)
	var htmlTableBody = MakeHtmlBody(dataInString, machines, attributes, runType, runId, cloudType, previousRunid)
	dataInString = strings.Replace(dataInString, "${rows}", htmlTableBody, -1)
	var htmlSubject = makeSubject(runType, runId)
	msgBody := []byte(dataInString)
	e := email.NewEmail()
	e.From = "Tools-Engg <automationreports@hike.in>"
	e.To = []string{"siva@hike.in"}
	e.Cc = []string{"tools@hike.in"}
	e.Subject = htmlSubject
	e.HTML = msgBody
	// tlsconfig := &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	ServerName:         "smtp-relay.gmail.com",
	// }
	fmt.Println(string(msgBody))
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "automationreports@hike.in", "Bharti@123", "smtp.gmail.com"))
	// println(err.Error())

}

func MakeHtmlBody(dataInString string, machines structConfigs.MachinesConfig, attributes structConfigs.AssertionConfig, runType string, runId int, cloudType string, previousRunid int) string {
	var outputRows = ""
	var count = 0
	if "MQTT" == runType {
		for i := 0; i < len(machines.CloudMachines[0].Mqttinstances); i++ {
			for j := 0; j < len(attributes.Attributes); j++ {
				count++
				var valueTobeFetched = GetMetricValueFromDB(attributes.Attributes[j].Value, runId, machines.CloudMachines[0].Mqttinstances[i].InstanceIP)
				var thresholdValue = getMetricValueFromBaseLine(i, j, previousRunid, attributes, machines)
				if valueTobeFetched > thresholdValue {
					outputRows = outputRows + "<tr><td align=left>" + strconv.Itoa(count) +
						"</td><td nowrap=nowrap>" + attributes.Attributes[j].Name +
						"</td><td><center> MQTT </td><td><center>" + machines.CloudMachines[0].Mqttinstances[i].InstanceIP +
						"</td><td><center> " + machines.CloudMachines[0].Mqttinstances[i].ServicesDeployed + " </td><td><b><center>" +
						strconv.FormatFloat(valueTobeFetched, 'f', 2, 64) + "  <font color=\"red\">(+" + strconv.FormatFloat(valueTobeFetched-thresholdValue, 'f', 2, 64) + ")</font></b></td></tr>"
					buildStatusFlag = false
				} else {
					outputRows = outputRows + "<tr><td align=left>" + strconv.Itoa(count) +
						"</td><td nowrap=nowrap>" + attributes.Attributes[j].Name +
						"</td><td><center> MQTT </td><td><center>" + machines.CloudMachines[0].Mqttinstances[i].InstanceIP +
						"</td><td><center> " + machines.CloudMachines[0].Mqttinstances[i].ServicesDeployed + " <td><center>" +
						strconv.FormatFloat(valueTobeFetched, 'f', 2, 64) + "  <font color=\"green\">(" + strconv.FormatFloat(valueTobeFetched-thresholdValue, 'f', 2, 64) + ")</font></b></td></tr>"
				}
			}
		}
		fmt.Println(outputRows)
	} else if "HTTP" == runType {
		var count = 0
		for i := 0; i < len(machines.CloudMachines[0].Httpinstances); i++ {
			for j := 0; j < len(attributes.Attributes); j++ {
				count++
				var valueTobeFetched = GetMetricValueFromDB(attributes.Attributes[j].Value, runId, machines.CloudMachines[0].Httpinstances[i].InstanceIP)
				var thresholdValue = getMetricValueFromBaseLine(i, j, previousRunid, attributes, machines)
				if valueTobeFetched > thresholdValue {
					outputRows = outputRows + "<tr><td align=left>" + strconv.Itoa(count) +
						"</td><td nowrap=nowrap>" + attributes.Attributes[j].Name +
						"</td><td><center> HTTP </td><td><center>" + machines.CloudMachines[0].Httpinstances[i].InstanceIP +
						"</td><td><center> " + machines.CloudMachines[0].Httpinstances[i].ServicesDeployed + " </td><td><b><center>" +
						strconv.FormatFloat(valueTobeFetched, 'f', 2, 64) + "  <font color=\"red\">(+" + strconv.FormatFloat(valueTobeFetched-thresholdValue, 'f', 2, 64) + ")</font></b></td></tr>"
					buildStatusFlag = false
				} else {
					outputRows = outputRows + "<tr><td align=left>" + strconv.Itoa(count) +
						"</td><td nowrap=nowrap>" + attributes.Attributes[j].Name +
						"</td><td><center> HTTP </td><td><center>" + machines.CloudMachines[0].Httpinstances[i].InstanceIP +
						"</td><td><center> " + machines.CloudMachines[0].Httpinstances[i].ServicesDeployed + " </td><td><center>" +
						strconv.FormatFloat(valueTobeFetched, 'f', 2, 64) + "  <font color=\"green\">(" + strconv.FormatFloat(valueTobeFetched-thresholdValue, 'f', 2, 64) + ")</font></b></td></tr>"
				}
			}
		}
		fmt.Println(outputRows)
	}
	return outputRows
}

func makeSubject(runType string, rundId int) string {
	var mailSubject string
	if buildStatusFlag {
		mailSubject = runType + " Load Test Results For the RunID " + strconv.Itoa(rundId) + "-- PASSED"
	} else {
		mailSubject = runType + " Load Test Results For the RunID " + strconv.Itoa(rundId) + "-- FAILED "
	}
	return mailSubject
}

func getMetricValueFromBaseLine(i int, j int, previousRunid int, attributes structConfigs.AssertionConfig, machines structConfigs.MachinesConfig) float64 {
	var metricThresholdValue float64
	if previousRunid == 0 {
		metricThreshold, err := strconv.ParseFloat(attributes.Attributes[j].ThresholdValue, 64)
		if err != nil {
			log.Fatalln(err)
		}
		metricThresholdValue = metricThreshold
	} else {
		metricThresholdValue = GetMetricValueFromDB(attributes.Attributes[j].Value, previousRunid, machines.CloudMachines[0].Mqttinstances[i].InstanceIP)
	}
	return metricThresholdValue
}

func dataTimeForGraphs(timeDuration time.Time) string {

	var splitTime = strings.Split(timeDuration.String(), " ")
	var date = strings.Split(splitTime[0], "-")
	var time = splitTime[1]
	var finalDate = strings.Join(date, "")
	fmt.Println(finalDate)
	var finalTime = time[3:len(time)]
	return finalTime + "_" + finalDate
}
