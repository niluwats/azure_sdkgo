package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/davecgh/go-spew/spew"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var deploymentName = "VMDeployQuickstart3"
var mutex sync.Mutex

const (
	resourceGroupName     = "GoVMQuickstart-1"
	resourceGroupLocation = "eastus"

	templateFile   = "vm-quickstart-template.json"
	parametersFile = "vm-quickstart-params.json"
)

type VmLogin struct {
	VmName     string `bson:"vm_name"`
	VmUserName string `bson:"vm_username"`
	VmPassword string `bson:"vm_password"`
	IpAdd      string `bson:"vm_ip"`
}
type ResourceGroup struct {
	Name     string    `bson:"resourcegroup_name"`
	Region   string    `bson:"region"`
	LoginDet []VmLogin `bson:"virtual_machine"`
}

type Organization struct {
	OrgName       string          `bson:"org_name"`
	ResourceGroup []ResourceGroup `bson:"resourcegroup"`
}

// Information loaded from the authorization file to identify the client
type clientInfo struct {
	SubscriptionID string
	VMPassword     string
}

var (
	ctx        = context.Background()
	clientData clientInfo
	authorizer autorest.Authorizer
)

// Authenticate with the Azure services using file-based authentication
func init() {
	var err error
	fmt.Println(os.Getenv("AZURE_AUTH_LOCATION"))
	authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Failed to get OAuth config: %v", err)
	}

	authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		log.Fatalf("Failed to read JSON: %+v", err)
	}
	clientData.SubscriptionID = (*authInfo)["subscriptionId"].(string)
	clientData.VMPassword = (*authInfo)["clientSecret"].(string)
}

func main() {
	var (
		vmname     = "QuickstartVM3"
		vmnetname  = "QuickstartNIC3"
		vmipname   = "QuickstartIP3"
		vmname2    = "QuickstartVM1"
		vmnetname2 = "QuickstartNIC1"
		vmipname2  = "QuickstartIP1"
	)
	// var wg sync.WaitGroup

	group, err := createGroup()
	if err != nil {
		log.Fatalf("failed to create group: %v", err)
	}
	log.Printf("Created group: %v", *group.Name)

	log.Printf("Starting deployment: %s", deploymentName)
	result, err := createDeployment(vmname, vmipname, vmnetname)
	result2, err2 := createDeployment(vmname2, vmipname2, vmnetname2)

	if err2 != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	if result.Name != nil {
		log.Printf("Completed deployment %v: %v", deploymentName, *result.Properties.ProvisioningState)
	} else {
		log.Printf("Completed deployment %v (no data returned to SDK)", deploymentName)
	}
	if result2.Name != nil {
		log.Printf("Completed deployment %v: %v", deploymentName, *result.Properties.ProvisioningState)
	} else {
		log.Printf("Completed deployment %v (no data returned to SDK)", deploymentName)
	}

	vmLoginCred := getLogin(vmname, vmipname, vmnetname)

	resGrp := ResourceGroup{
		Name:   *group.Name,
		Region: *group.Location,
		LoginDet: []VmLogin{
			vmLoginCred,
		},
	}
	org := Organization{
		OrgName:       "fcx",
		ResourceGroup: []ResourceGroup{resGrp},
	}

	db(org, resGrp, vmLoginCred)
	// time.Sleep(25 * time.Second)

	// for i := 0; i < 1; i++ {
	// 	wg.Add(1)

	// 		deploymentName = "VMDeployQuickstart1"
	// 		var (
	// 			vmname    = "QuickstartVM1"
	// 			vmnetname = "QuickstartNIC1"
	// 			vmipname  = "QuickstartIP1"
	// 		)
	// 		// modifyJson()

	// 		log.Printf("Starting deployment: %s", deploymentName)
	// 		result, err := createDeployment(vmname, vmipname, vmnetname)
	// 		if err != nil {
	// 			log.Fatalf("Failed to deploy: %v", err)
	// 		}
	// 		if result.Name != nil {
	// 			log.Printf("Completed deployment %v: %v", deploymentName, *result.Properties.ProvisioningState)
	// 		} else {
	// 			log.Printf("Completed deployment %v (no data returned to SDK)", deploymentName)
	// 		}

	// 		vmLoginCred := getLogin(vmname, vmipname, vmnetname)

	// 		resGrp := ResourceGroup{
	// 			Name:   *group.Name,
	// 			Region: *group.Location,
	// 			LoginDet: []VmLogin{
	// 				vmLoginCred,
	// 			},
	// 		}
	// 		org := Organization{
	// 			OrgName:       "fcx",
	// 			ResourceGroup: []ResourceGroup{resGrp},
	// 		}

	// 		db(org, resGrp, vmLoginCred)
	// 		defer wg.Done()

	// }
	// wg.Wait()
}

// Create a resource group for the deployment.
func createGroup() (group resources.Group, err error) {
	groupsClient := resources.NewGroupsClient(clientData.SubscriptionID)
	groupsClient.Authorizer = authorizer

	return groupsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resources.Group{
			Location: to.StringPtr(resourceGroupLocation)})
}

// Create the deployment
func createDeployment(vmname, vmipname, vmnetnic string) (deployment resources.DeploymentExtended, err error) {
	// mutex.Lock()
	modifyJson(vmname, vmipname, vmnetnic)
	template, err := readJSON(templateFile)
	if err != nil {
		return
	}
	params, err := readJSON(parametersFile)
	fmt.Println((*params)["virtualMachines_QuickstartVM_name"].(map[string]interface{}))
	fmt.Println((*params)["networkInterfaces_quickstartvm_name"].(map[string]interface{}))
	fmt.Println((*params)["publicIPAddresses_QuickstartVM_ip_name"].(map[string]interface{}))

	if err != nil {
		return
	}
	(*params)["vm_password"] = map[string]string{
		"value": clientData.VMPassword,
	}

	deploymentsClient := resources.NewDeploymentsClient(clientData.SubscriptionID)
	deploymentsClient.Authorizer = authorizer

	deploymentFuture, err := deploymentsClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       resources.Incremental,
			},
		},
	)
	if err != nil {
		return
	}
	err = deploymentFuture.WaitForCompletionRef(ctx, deploymentsClient.BaseClient.Client)
	if err != nil {
		fmt.Println("deploymentFuture - ", err)
		return
	}
	// mutex.Unlock()
	return deploymentFuture.Result(deploymentsClient)
}

// Get login information by querying the deployed public IP resource.
func getLogin(vmname, ipname, netnic string) VmLogin {

	params, err := readJSON(parametersFile)

	if err != nil {
		log.Fatalf("Unable to read parameters. Get login information with `az network public-ip list -g %s", resourceGroupName)
	}

	addressClient := network.NewPublicIPAddressesClient(clientData.SubscriptionID)
	addressClient.Authorizer = authorizer
	ipName := (*params)["publicIPAddresses_QuickstartVM_ip_name"].(map[string]interface{})
	ipAddress, err := addressClient.Get(ctx, resourceGroupName, ipName["value"].(string), "")
	if err != nil {
		log.Fatalf("Unable to get IP information. Try using `az network public-ip list -g %s", resourceGroupName)
	}

	vmUser := (*params)["vm_user"].(map[string]interface{})
	vmName := (*params)["virtualMachines_QuickstartVM_name"].(map[string]interface{})
	fmt.Println(vmName)
	log.Printf("Log in with ssh: %s@%s, password: %s",
		vmUser["value"].(string),
		*ipAddress.PublicIPAddressPropertiesFormat.IPAddress,
		clientData.VMPassword)

	lg := VmLogin{
		VmName:     vmName["value"].(string),
		VmUserName: vmUser["value"].(string),
		VmPassword: clientData.VMPassword,
		IpAdd:      *ipAddress.IPAddress,
	}
	return lg
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return &contents, nil
}
func modifyJson(vmname, vmipname, vmnic string) {
	data, err := ioutil.ReadFile(parametersFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	params := make(map[string]interface{})
	_ = json.Unmarshal(data, &params)

	params["virtualMachines_QuickstartVM_name"] = vmname
	params["networkInterfaces_quickstartvm_name"] = vmnic
	params["publicIPAddresses_QuickstartVM_ip_name"] = vmipname

	err = ioutil.WriteFile(parametersFile, data, 0644)
}
func db(org Organization, resGrp ResourceGroup, vmLoginCred VmLogin) {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}
	db := session.DB("bethel_dashboard")
	col := db.C("organizations")
	defer session.Close()

	var res Organization
	err = col.Find(bson.M{"org_name": "fcx", "resourcegroup.resourcegroup_name": resGrp.Name}).One(&res)
	spew.Dump(res)
	if err == mgo.ErrNotFound {
		err = col.Insert(&org)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		pushQuery := bson.M{"resourcegroup.$.virtual_machine": vmLoginCred}
		err1 := col.Update(bson.M{"resourcegroup.resourcegroup_name": resGrp.Name}, bson.M{"$addToSet": pushQuery})
		if err1 != nil {
			fmt.Println(err1.Error())
			log.Fatal(err1)
		}
	}
}
