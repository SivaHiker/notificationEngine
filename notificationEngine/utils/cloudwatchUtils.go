package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/cloudwatch"
)

func getCloudWatchInstance() *cloudwatch.CloudWatch {
	region := aws.Regions["ap-southeast-1"]
	// namespace := "AWS/EC2"
	// dimension := &cloudwatch.Dimension{
	// Name:  "InstanceId",
	// Value: "i-0a4864ae",
	// }
	// metricName := "CPUUtilization"
	// now := time.Now()
	// prev := now.Add(time.Duration(600) * time.Second * -1) // 600 secs = 10 minutes

	auth, err := aws.GetAuth("AKIAJ6IRJMJ6ONSHZQXQ", "burqVrRbkiYozJDQ3BYnAFFf9vsvzSs7y7y9DHAX", "", time.Now())
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}

	cw, err := cloudwatch.NewCloudWatch(auth, region.CloudWatchServicepoint)
	if err == nil {
		return cw
	} else {
		return nil
	}
}

func getCloudWatchMetrics(cw *cloudwatch.CloudWatch, dimension *cloudwatch.Dimension, Namespace string, metricName string, statistics []string, startTime time.Time, endTime time.Time) *cloudwatch.GetMetricStatisticsResponse {
	request := &cloudwatch.GetMetricStatisticsRequest{
		Dimensions: []cloudwatch.Dimension{*dimension},
		EndTime:    endTime,
		StartTime:  startTime,
		MetricName: metricName,
		Period:     3600,
		Statistics: statistics,
		Namespace:  Namespace,
	}

	response, err := cw.GetMetricStatistics(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
	} else {
		fmt.Printf("Error: %+v\n", err)
	}
	return response
}
