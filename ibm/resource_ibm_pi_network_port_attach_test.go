// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
)

func TestAccIBMPINetworkPortAttachbasic(t *testing.T) {
	t.Skip()
	name := fmt.Sprintf("tf-pi-network-port-attach-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMPINetworkPortAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPINetworkPortAttachConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPINetworkPortAttachExists("ibm_pi_network_port_attach.power_network_port_attach"),
					resource.TestCheckResourceAttr(
						"ibm_pi_network_port_attach.power_network_port_attach", "pi_network_name", name),
				),
			},
		},
	})
}
func testAccCheckIBMPINetworkPortAttachDestroy(s *terraform.State) error {

	sess, err := testAccProvider.Meta().(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_pi_network_port_attach" {
			continue
		}
		parts, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID := parts[0]
		networkID := parts[1]
		portID := parts[2]
		networkC := st.NewIBMPINetworkClient(context.Background(), sess, cloudInstanceID)
		_, err = networkC.GetPort(networkID, portID)
		if err == nil {
			return fmt.Errorf("PI Network Port still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPINetworkPortAttachExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		sess, err := testAccProvider.Meta().(ClientSession).IBMPISession()
		if err != nil {
			return err
		}
		parts, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID := parts[0]
		networkID := parts[1]
		portID := parts[2]
		client := st.NewIBMPINetworkClient(context.Background(), sess, cloudInstanceID)

		network, err := client.GetPort(networkID, portID)
		if err != nil {
			return err
		}
		parts[1] = *network.PortID
		return nil

	}
}

func testAccCheckIBMPINetworkPortAttachConfig(name string) string {
	return testAccCheckIBMPINetworkPortConfig(name) + fmt.Sprintf(`
	resource "ibm_pi_network_port_attach" "power_network_port_attach" {
		pi_cloud_instance_id  = "%s"
		pi_network_name       = ibm_pi_network.power_networks.pi_network_name
		pi_network_port_description = "IP Reserved for Test UAT"
		port_id=ibm_pi_network_port.power_network_port.id
	}
	`, pi_cloud_instance_id)
}
