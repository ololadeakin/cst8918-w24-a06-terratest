package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "ad8edcbe-5549-4328-aad3-0a825b857928"


func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
	// The path to where our Terraform code is located
	TerraformDir: "../",
	// Override the default terraform variables
	Vars: map[string]interface{}{
	"labelPrefix": "akin0098",
	},
	}
	defer terraform.Destroy(t, terraformOptions)
	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)
	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
	// Get the NIC name from terraform output
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	// Confirm NIC exists
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))
	// Confirm the VM is running the correct Ubuntu version
	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, "Canonical", *vm.StorageProfile.ImageReference.Publisher)
	assert.Equal(t, "0001-com-ubuntu-server-jammy", *vm.StorageProfile.ImageReference.Offer)
	assert.Equal(t, "22_04-lts-gen2", *vm.StorageProfile.ImageReference.Sku)
	}
