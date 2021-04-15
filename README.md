WIP :D

This Controller will read Egress Network Policies and generate Falco Rules based on them, so
one can know when an application is trying to connect to some non authorized endpoint.

There's a huge TODO/Backlog here:

* Adjust the rules format to better optimization
* Verify if we want to write the rules as ConfigMap vs CRDs
* Create a sidecar for Falco, to download the rules and verify before applying/kill -HUP in Falco
* Deal with the NetworkPolicyPeer reconciliation 
 * In case of CIDRs, update the list
 * In case of a selected Pod that changes its IP, reconcile and update the list
* Implement the port list creation
* Unit tests for the helper functions :D
* Improve the output / description messages with at least the Pod and Namespace

Maybe in a v2:
* How this can help with Ingress Policies? Need to figure out.
