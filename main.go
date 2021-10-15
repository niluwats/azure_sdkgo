package main

import (
	"context"

	// "strings"
	"time"

	"github.com/niluwats/azuresdk_go/insights"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	resId := "/subscriptions/dcb8d550-31da-4e1b-afd6-b41939ab339c/resourcegroups/vmtest/providers/microsoft.compute/virtualmachines/vmfirst1"
	// // metrics, err := insights.ListMetricDefinitions(resId)
	// if err != nil {
	// 	util.LogAndPanic(err)
	// }
	// util.PrintAndLog("available metrics:")
	// util.PrintAndLog(strings.Join(metrics, "\n"))

	insights.GetMetricsData(ctx, resId, []string{"Network In Total", "Network Out Total", "Percentage CPU", "Available Memory Bytes"})
}
