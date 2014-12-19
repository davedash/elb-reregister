# elb-reregister
--
Sometimes ELB marks an instance as failed and after resolving the issue you want
to put the node back "InService".

To install:

    go get github.com/davedash/elb-reregister

Usage:

    elb-reregister elb-name instance-id

This will unregister a node if it is not "InService" and then register it. If
the node does not exist in ELB it will be added to the load balancer.
