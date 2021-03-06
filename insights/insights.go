// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package insights

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/niluwats/azuresdk_go/internal/config"
	"github.com/niluwats/azuresdk_go/internal/iam"
)
// ListMetricDefinitions returns the list of metrics available for the specified resource in the form "Localized Name (metric name)".
func ListMetricDefinitions(resourceURI string) ([]string, error) {
	a, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		return nil, err
	}

	metricsDefClient := insights.NewMetricDefinitionsClient(config.SubscriptionID())
	metricsDefClient.Authorizer = a
	metricsDefClient.AddToUserAgent(config.UserAgent())
	result, err := metricsDefClient.List(context.Background(), resourceURI, "")
	if err != nil {
		return nil, err
	}
	metrics := make([]string, len(*result.Value))
	for i := range *result.Value {
		metrics[i] = fmt.Sprintf("%s (%s)", *(*result.Value)[i].Name.LocalizedValue, *(*result.Value)[i].Name.Value)
	}
	return metrics, nil
}

// GetMetricsData returns the specified metric data points for the specified resource ID spanning the last five minutes.
func GetMetricsData(ctx context.Context, resourceID string, metrics []string) ([]string, error) {
	a, err := iam.GetResourceManagementAuthorizer()
	if err != nil {
		return nil, err
	}
	metricsClient := insights.NewMetricsClient(config.SubscriptionID())
	metricsClient.Authorizer = a
	metricsClient.AddToUserAgent(config.UserAgent())

	endTime := time.Now().UTC()
	startTime := endTime.Add(time.Duration(-2) * time.Minute)
	timespan := fmt.Sprintf("%s/%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	resp, err := metricsClient.List(context.Background(), resourceID, timespan, nil, strings.Join(metrics, ","), "Average,count", nil, "", "", insights.Data, "")
	fmt.Println("resp-", resp.Value)

	if err != nil {
		return nil, err
	}
	var metricData []string

	for _, v := range *resp.Value {
		for _, t := range *v.Timeseries {
			for _, mv := range *t.Data {
				fmt.Println("mv----", mv)
				avg := 0.0
				cou := 0.0
				if mv.Average != nil {
					avg = byteToMb(*mv.Average)
				}
				if mv.Count != nil {
					cou = *mv.Count
				}
				metricData = append(metricData, fmt.Sprintf("%s @ %s - avg: %f, cou: %f", *v.Name.LocalizedValue, *mv.TimeStamp, avg, cou))
			}
		}
	}

	fmt.Println("met", metricData)
	return metricData, nil
}
func byteToMb(i float64) float64 {
	return i / (1024 * 1024)
}

//ctx context.Context
//resourceURI string
//timespan string
//interval *string,
// metricnames string,
//aggregation string,
// top *int32,
// orderby string,
// filter string,
// resultType ResultType,
// metricnamespace string
