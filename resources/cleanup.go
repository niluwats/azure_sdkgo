// package resources

// import (
// 	"context"
// 	"log"

// 	"github.com/niluwats/azuresdk_go/internal/config"
// )

// // Cleanup deletes the resource group created for the sample
// func Cleanup(ctx context.Context) {
// 	if config.KeepResources() {
// 		log.Println("Hybrid resources cleanup: keeping resources")
// 		return
// 	}
// 	log.Println("Hybrid resources cleanup: deleting resources")
// 	_, _ = DeleteGroup(ctx, config.GroupName())
// }