package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/niluwats/azuresdk_go/clusters/aks"
	"github.com/niluwats/azuresdk_go/internal/util"
)

var (
	aksClusterName      string = "clustertest1"
	aksUsername         string = "azureuser"
	aksAgentPoolCount   int32  = 2
	aksSSHPublicKeyPath        = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
)

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Hour*1))
	defer cancel()

	// _, err := resources.CreateGroup(ctx, config.GroupName())
	// if err != nil {
	// 	util.LogAndPanic(err)
	// }

	// _, err = aks.CreateAKS(ctx, aksClusterName, config.Location(), config.GroupName(), aksUsername, aksSSHPublicKeyPath, config.ClientID(), config.ClientSecret(), aksAgentPoolCount)
	// if err != nil {
	// 	util.LogAndPanic(err)
	// }
	// util.PrintAndLog("created AKS cluster")

	c, err := aks.GetAKS(ctx, "kubernew", "clustern1")
	if err != nil {
		util.LogAndPanic(err)
	}

	// fmt.Println(c)
	fmt.Println("ID--", *c.ID)
	fmt.Println("Name--", *c.Name)
	fmt.Println("Type--", *c.Type)
	fmt.Println("Location--", *c.Location)
	fmt.Println("Tags--", &c.Tags)

	fmt.Println("_______________________")
	fmt.Println("--c.ManagedClusterProperties--")

	fmt.Println("ProvisioningState   ------", *c.ManagedClusterProperties.ProvisioningState)
	fmt.Println("DNSPrefix  ------", *c.ManagedClusterProperties.DNSPrefix)
	fmt.Println("Fqdn  ------", *c.ManagedClusterProperties.Fqdn)
	fmt.Println("KubernetesVersion  ------", *c.ManagedClusterProperties.KubernetesVersion)
	fmt.Println("LinuxProfile  ------", *c.ManagedClusterProperties.LinuxProfile)
	fmt.Println("ServicePrincipalProfile  ------", *c.ManagedClusterProperties.ServicePrincipalProfile)

	fmt.Println("_______________________")
	fmt.Println("c.ManagedClusterProperties.AgentPoolProfiles_______________________")

	for _, v := range *c.ManagedClusterProperties.AgentPoolProfiles {
		fmt.Println("Name -   ", *v.Name)
		fmt.Println("Count -   ", *v.Count)
		fmt.Println("VMSize -   ", v.VMSize)
		fmt.Println("OsDiskSizeGB -   ", *v.OsDiskSizeGB)
		fmt.Println("DNSPrefix -   ", v.DNSPrefix)
		fmt.Println("Fqdn -   ", v.Fqdn)
		fmt.Println("Ports -   ", v.Ports)
		fmt.Println("StorageProfile -   ", v.StorageProfile)
		fmt.Println("VnetSubnetID -   ", v.VnetSubnetID)
		fmt.Println("OsType -   ", v.OsType)
	}
}

