package model

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	KubeDeployment  = "deployments"
	KubeStatefulset = "statefulsets"
	KubeDaemonset   = "daemonsets"
	KubePod         = "pods"
	KubeService     = "services"
	KubeIngress     = "ingresses"
)

type KubeBase struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type KubeApp struct {
	KubeBase
	Kind       string             `json:"kind"`
	Replicas   *int32             `json:"replicas"`
	Containers []corev1.Container `json:"containers"`
}
