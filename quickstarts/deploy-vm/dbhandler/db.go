package dbhandler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

type Metrics struct {
	MetricId          primitive.ObjectID     `bson:"_id,omitempty"`
	VMID              interface{}            `bson:"virtualmachine,omitempty"`
	NetworkIn         map[string]interface{} `bson:"network_in_total"`
	NetworkOut        map[string]interface{} `bson:"network_out_total"`
	PercentCpu        map[string]interface{} `bson:"percantage_cpu"`
	AvailableMemBytes map[string]interface{} `bson:"available_memory_bytes"`
}
type VirtualMachine struct {
	VMID       primitive.ObjectID `bson:"_id,omitempty"`
	ResGrpId   interface{}        `bson:"resourcegroup,omitempty"`
	VmName     string             `bson:"vm_name" json:"virtual_machine_name"`
	VmUserName string             `bson:"vm_username" json:"username,omitempty"`
	VmPassword string             `bson:"vm_password" json:"password,omitempty"`
	IpAdd      string             `bson:"vm_ip" json:"public_ip_address"`
}
type ResourceGroup struct {
	ResGrpId primitive.ObjectID `bson:"_id,omitempty"`
	OrgId    interface{}        `bson:"organization,omitempty"`
	Name     string             `bson:"resourcegroup_name" json:"resourcegroup_name"`
	Region   string             `bson:"region" json:"region"`
}

type Organization struct {
	OrgId   primitive.ObjectID `bson:"_id,omitempty"`
	OrgName string             `bson:"org_name" json:"organization_name"`
}

func SaveVm(lg *VirtualMachine) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err) 
	}
	defer client.Disconnect(ctx)
	database := client.Database("bethel_dashboard")
	orgCollection := database.Collection("organizations")
	resCollection := database.Collection("resourcegroups")
	vmCollection := database.Collection("virtualmachines")
	//metCollection := database.Collection("metrics")

	or := Organization{
		OrgId:   primitive.NewObjectID(),
		OrgName: "fcx",
	}

	rg := ResourceGroup{
		Name:   "vmres",
		Region: "westus",
	}
	vm := VirtualMachine{
		VmName:     lg.VmName,
		VmUserName: lg.VmUserName,
		VmPassword: lg.VmPassword,
		IpAdd:      lg.IpAdd,
	}

	var filtered []bson.M
	filterCursor, err1 := orgCollection.Find(ctx, bson.M{"org_name": "fcx"})
	if err1 != nil {
		log.Fatal(err1)
	}
	if err = filterCursor.All(ctx, &filtered); err != nil {
		log.Fatal(err)
	}
	if filtered != nil {
		c := filtered[0]
		rg.OrgId = c["_id"]
	}
	if filtered == nil {
		insertResult, err := orgCollection.InsertOne(ctx, or)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertResult.InsertedID)
		rg.OrgId = insertResult.InsertedID
	}

	var filtered1 []bson.M
	filterCursor1, err1 := resCollection.Find(ctx, bson.M{"resourcegroup_name": "vmres"})
	if err1 != nil {
		log.Fatal(err1) 
	}
	if err = filterCursor1.All(ctx, &filtered1); err != nil {
		log.Fatal(err)
	}
	if filtered1 != nil {
		c := filtered1[0]
		vm.ResGrpId = c["_id"]
	}
	if filtered1 == nil {
		insertResult, err := resCollection.InsertOne(ctx, rg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertResult.InsertedID)
		vm.ResGrpId = insertResult.InsertedID
	}

	var Filtered2 []bson.M
	filterCursor2, err1 := vmCollection.Find(ctx, bson.M{"vm_name": "vm1"})
	if err1 != nil {
		log.Fatal(err1)
	}
	if err = filterCursor2.All(ctx, &Filtered2); err != nil {
		log.Fatal(err)
	}
	if Filtered2 == nil {
		insertResult, err := vmCollection.InsertOne(ctx, vm)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertResult.InsertedID)

	}
	// countOrg, err := orgCollection.CountDocuments(ctx, bson.M{"org_name": "fcx"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if countOrg == 0 {
	// 	fmt.Println("not found")

	// 	orgRes, err0 = orgColl.InsertOne(ctx, &or)
	// 	if err0 != nil {
	// 		log.Fatal(err0)
	// 	}
	// }
	// countRG, err1 := resGrpColl.CountDocuments(ctx, bson.M{"resourcegroup_name": resourceGroupName})
	// if err1 != nil {
	// 	log.Fatal(err1)
	// }
	// if countRG == 0 {
	// 	fmt.Println("no such resource group before")
	// 	rg := ResourceGroup{
	// 		OrgId:  orgRes.InsertedID,
	// 		Name:   resourceGroupName,
	// 		Region: resourceGroupLocation,
	// 	}

	// 	rgRes, err0 = resGrpColl.InsertOne(ctx, &rg)
	// 	if err0 != nil {
	// 		log.Fatal(err0)
	// 	}
	// }

	// vmRes, err0 = vmColl.InsertOne(ctx, &vm)
	// if err0 != nil {
	// 	log.Fatal(err0)
	// }

}
