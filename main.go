// Sometimes ELB marks an instance as failed and after resolving the issue
// you want to put the node back "InService".
//
// To install:
//   go get github.com/davedash/elb-reregister
//
// Usage:
//   elb-reregister elb-name instance-id
//
// This will unregister a node if it is not "InService" and then register it.
// If the node does not exist in ELB it will be added to the load balancer.
package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/elb"
	"log"
	"os"
)

func usage() {
	fmt.Printf("Usage: %s elb-name instance-id\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	auth, err := aws.GetAuth("", "")
	if err != nil {
		log.Fatal(err)
	}

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	elbName := flag.Arg(0)
	instanceQuery := flag.Arg(1)

	client := elb.New(auth, aws.USEast)

	resp, err := client.DescribeInstanceHealth(&elb.DescribeInstanceHealth{
		LoadBalancerName: elbName})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Instances in %s:\n", elbName)

	instanceExists := false
	instanceHealthy := true

	for _, state := range resp.InstanceStates {
		fmt.Printf(" - %s (%s)\n", state.InstanceId, state.State)
		if state.InstanceId == instanceQuery {
			instanceExists = true
			if state.State != "InService" {
				instanceHealthy = false
			}
		}
	}

	fmt.Println()

	if !instanceHealthy {
		fmt.Printf("%s is not healthy, I should de-register\n", instanceQuery)
		_, err := client.DeregisterInstancesFromLoadBalancer(
			&elb.DeregisterInstancesFromLoadBalancer{
				LoadBalancerName: elbName,
				Instances:        []string{instanceQuery}})

		if err != nil {
			log.Fatal(err)
		}
	}

	if !(instanceExists && instanceHealthy) {
		// Register instance
		fmt.Printf("Registering %s\n", instanceQuery)
		_, err := client.RegisterInstancesWithLoadBalancer(
			&elb.RegisterInstancesWithLoadBalancer{
				LoadBalancerName: elbName,
				Instances:        []string{instanceQuery}})

		if err != nil {
			log.Fatal(err)
		}
	}

	return

}
