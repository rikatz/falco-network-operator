package utils

import (
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
)

func TestEgress2CondListed(t *testing.T) {
	type args struct {
		npname     string
		egressRule networkingv1.NetworkPolicyEgressRule
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Egress2CondListed(tt.args.npname, tt.args.egressRule); got != tt.want {
				t.Errorf("Egress2CondListed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetPol2FalcoCond(t *testing.T) {
	type args struct {
		netpol networkingv1.NetworkPolicy
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NetPol2FalcoCond(tt.args.netpol); got != tt.want {
				t.Errorf("NetPol2FalcoCond() = %v, want %v", got, tt.want)
			}
		})
	}
}
