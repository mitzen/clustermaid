package util

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeObject struct {
	c *kubernetes.Clientset
}

func (n *KubeObject) NewKubeObject(clientSet *kubernetes.Clientset) {
	n.c = clientSet
}

func (n *KubeObject) ListAllNamespace() (*v1.NamespaceList, error) {
	return n.c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (n *KubeObject) ListAllNodes() (*v1.NodeList, error) {
	return n.c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func (n *KubeObject) ListDaemonset(namespace string) (*appsv1.DaemonSetList, error) {
	return n.c.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (n *KubeObject) ListAllPods(namespace string) (*v1.PodList, error) {
	return n.c.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}
