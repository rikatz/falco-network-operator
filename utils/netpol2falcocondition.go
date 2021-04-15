package utils

import (
	"fmt"
	"strings"

	networkingv1 "k8s.io/api/networking/v1"
)

func NetPol2FalcoCond(netpol networkingv1.NetworkPolicy) string {

	var podSelector string
	var labelIndex int

	var condition strings.Builder

	for label, value := range netpol.Spec.PodSelector.MatchLabels {
		if labelIndex > 0 {
			podSelector = podSelector + " and "
		}
		podSelector = podSelector + fmt.Sprintf("k8s.pod.label.%s=\"%s\"", label, value)
		labelIndex++
	}

	condition.WriteString(podSelector)

	if len(netpol.Spec.Egress) > 0 {
		condition.WriteString(" and (")
		for k, np := range netpol.Spec.Egress {
			if k > 0 {
				condition.WriteString(" or ")
			}
			// TODO: Completely prone to bugs when no ports and to are defined, need to fix
			condition.WriteString("(")
			npname := fmt.Sprintf("%s-%s-%d", netpol.Namespace, netpol.Name, k)
			condition.WriteString(Egress2CondListed(npname, np))
			condition.WriteString(")")

		}
		condition.WriteString(")")
	}

	return condition.String()
}

// Egress2CondListed generates a condition based on an EgressRule
// It will use List in Ports and Destination instead of parsing those fields
// Returned format: fd.sport in (npname-port) and fd.sip in (npname-cidr)
func Egress2CondListed(npname string, egressRule networkingv1.NetworkPolicyEgressRule) string {
	if len(egressRule.To) > 0 && len(egressRule.Ports) > 0 {
		return fmt.Sprintf("fd.sport in (%s-ports) and fd.sip in (%s-cidr)", npname, npname)
	}

	if len(egressRule.Ports) > 0 {
		return fmt.Sprintf("fd.sport in (%s-ports)", npname)
	}

	if len(egressRule.To) > 0 {
		return fmt.Sprintf("fd.sip in (%s-cidr)", npname)
	}

	return ""
}
