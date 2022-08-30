package feature

import (
	"cdx.foc/clusterMaid/pkg/kube/util"
)

type ClusterManager struct {
	Namespace string
}

const (
	unitMegabytes int64 = 1000000
)

// Get cluster cpu / memory configuration
// Get pods memory and cpu request and limits
// Total it all up

func (i *ClusterManager) Execute(namespace string) {

	//cfg := config.ClientConfig{}
	//restConfig := cfg.NewRestConfig()
	//clientset := cfg.NewClientSet(restConfig)

	//nsutil := util.KubeObject{}
	//nsutil.NewKubeObject(clientset)
	if namespace == "" {
		namespace = "default"
	}

	ic := util.IstioClient{}
	ic.NewIstioClient(namespace)

	l := ic.GetPods()

	for _, pod := range l.Items {
		ic.GetRoutesConfig(pod.Name, namespace, "", "summary")
	}

	// nodes, err := nsutil.ListAllNodes()

	// if err != nil {
	// 	fmt.Printf("Unable get nodes info.")
	// }

	// for _, node := range nodes.Items {

	// 	fmt.Printf("----------------------------------------------\n")
	// 	fmt.Printf("Node: %s \n", node.Name)
	// 	fmt.Printf("Pods: %d \n", node.Status.Capacity.Storage().ToDec().Value())
	// 	fmt.Printf("----------------------------------------------\n")

	// 	pods, err := nsutil.ListAllPods(apiv1.NamespaceAll)
	// 	if err != nil {
	// 		fmt.Printf("Unable get nodes info.")
	// 	}

	// 	var (
	// 		totalMemoryRequested, totalCPURequested, totalCPULimit, totalMemoryLimit int64
	// 	)

	// 	totalMemoryRequested = 0
	// 	totalCPURequested = 0
	// 	totalCPULimit = 0
	// 	totalMemoryLimit = 0

	// 	for _, pod := range pods.Items {
	// 		if pod.Spec.NodeName == node.Name {

	// 			fmt.Printf("Namespace: %s \n", pod.Namespace)
	// 			fmt.Printf("Pod name: %s \n", pod.Name)

	// 			for _, container := range pod.Spec.Containers {

	// 				CPURequested := container.Resources.Requests.Cpu().Value()
	// 				MemoryRequested := container.Resources.Requests.Memory().Value()
	// 				CPULimit := container.Resources.Limits.Cpu().Value()
	// 				MemoryLimit := container.Resources.Limits.Memory().Value()

	// 				fmt.Printf("Container Name: %s \n", container.Name)
	// 				fmt.Printf("Image Name: %s \n", container.Image)
	// 				fmt.Printf("Cpu request for container: %d \n", CPURequested)
	// 				fmt.Printf("Cpu limits for container: %d \n", CPULimit)
	// 				fmt.Printf("Memory request for container (M): %d \n", MemoryRequested/unitMegabytes)
	// 				fmt.Printf("Memory limits for container (M): %d \n", MemoryLimit/unitMegabytes)

	// 				totalCPURequested += CPURequested
	// 				totalMemoryRequested += MemoryRequested
	// 				totalCPULimit += CPULimit
	// 				totalMemoryLimit += MemoryLimit
	// 			}
	// 		}
	// 	}

	// 	fmt.Printf("Total cpu request for container: %d \n", totalCPURequested)
	// 	fmt.Printf("Total cpu limits for container: %d \n", totalCPULimit)
	// 	fmt.Printf("Total memory request for container:%d \n", totalMemoryRequested/unitMegabytes)
	// 	fmt.Printf("Total memory limits for container:%d \n", totalMemoryLimit/unitMegabytes)
	// 	fmt.Printf("Node CPU: %d \n", node.Status.Capacity.Cpu().ToDec().Value())
	// 	fmt.Printf("Node Memory: %d \n", node.Status.Capacity.Memory().ToDec().Value()/unitMegabytes)
	// }
}

func (i *ClusterManager) GetMetric() {

}
