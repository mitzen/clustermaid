package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mitzen/kubeconfig/config"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"istio.io/istio/istioctl/pkg/writer/envoy/configdump"
	"istio.io/istio/pkg/kube"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

const (
	istioSystemNamespace = "istio-system"
)

type IstioClient struct {
	IstioClient         *versionedclient.Clientset
	namespace           string
	IstioExtendedClient kube.CLIClient
}

func (i *IstioClient) NewIstioClient(config *rest.Config, namespace string) {

	ic, err := versionedclient.NewForConfig(config)
	i.initializeIstioClient()

	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}

	i.IstioClient = ic
	i.namespace = namespace
}

func (i *IstioClient) GetIstioControlVersion() string {

	mvi, err := i.IstioExtendedClient.GetIstioVersions(context.TODO(), istioSystemNamespace)

	if err != nil {
		fmt.Printf("Unable to get version istiod")
	}

	for _, v := range *mvi {
		if v.Info.Version != "" {
			return v.Info.Version
		}
	}
	return ""
}

func (i *IstioClient) GetIstioPod(namespace string) string {

	mvi, err := i.IstioExtendedClient.GetIstioPods(context.TODO(), namespace, map[string]string{})

	if err != nil {
		fmt.Println("error getting pods")
	}

	for _, v := range mvi {

		for _, a := range v.Spec.Containers {
			if strings.Contains(a.Name, config.IstioProxyImage) {
				ss := strings.Split(a.Image, ":")
				istioProxyVersion := ss[len(ss)-1]
				return istioProxyVersion
			}
		}
	}
	return ""
}

func (i *IstioClient) GetGateways() (*v1alpha3.GatewayList, error) {
	return i.IstioClient.NetworkingV1alpha3().Gateways(i.namespace).List(context.TODO(), v1.ListOptions{})
}

func (i *IstioClient) GetVirtualServices() {
	i.IstioClient.NetworkingV1alpha3().VirtualServices(i.namespace)
}

func (i *IstioClient) GetDesinationRules() {
	i.IstioClient.NetworkingV1alpha3().DestinationRules(i.namespace)
}

func (i *IstioClient) GetRouteDestination(podName string, ns string, name string, outputFormat string) error {

	var configWriter *configdump.ConfigWriter
	var err error

	configWriter, err = i.setupPodConfigdumpWriter(podName, ns, false, nil)

	if err != nil {
		return err
	}

	filter := configdump.RouteFilter{
		Name:    name,
		Verbose: true,
	}

	switch outputFormat {
	case "summary":
		return configWriter.PrintRouteSummary(filter)
	case "detailed":
		return configWriter.PrintRouteDump(filter, outputFormat)
	default:
		return fmt.Errorf("output format %q not supported", outputFormat)
	}
}

func (i *IstioClient) initializeIstioClient() {
	client, err := kube.NewCLIClient(kube.BuildClientCmd("", ""), "")
	if err != nil {
		fmt.Println("Unable to create istio extended client")
	}

	i.IstioExtendedClient = client
}

func (i *IstioClient) setupPodConfigdumpWriter(podName, podNamespace string, includeEds bool, out io.Writer) (*configdump.ConfigWriter, error) {
	debug, err := i.extractConfigDump(podName, podNamespace, includeEds)
	if err != nil {
		return nil, err
	}
	return i.setupConfigdumpEnvoyConfigWriter(debug, out)
}

func (i *IstioClient) extractConfigDump(podName, podNamespace string, eds bool) ([]byte, error) {

	path := "config_dump"
	if eds {
		path += "?include_eds=true"
	}
	debug, err := i.IstioExtendedClient.EnvoyDo(context.TODO(), podName, podNamespace, "GET", path)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command on %s.%s sidecar: %v", podName, podNamespace, err)
	}
	return debug, err
}

func (i *IstioClient) setupConfigdumpEnvoyConfigWriter(debug []byte, out io.Writer) (*configdump.ConfigWriter, error) {
	cw := &configdump.ConfigWriter{Stdout: out}
	err := cw.Prime(debug)
	if err != nil {
		return nil, err
	}
	return cw, nil
}
