package main

import (
	"context"
	"log"
	"path/filepath"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"

	"github.com/maxgio92/cloudevents-podgateway-client-provisioner/pkg/k8s"
	"github.com/maxgio92/cloudevents-podgateway-client-provisioner/pkg/provision"
)

var (
	clientNamespace      string
	clientContainerImage string
	clientCommand        []string
	kubeconfig           string
)

func main() {
	flag.StringVarP(&clientNamespace, "client-namespace", "n", provision.DefaultNamespace, "The namespace where pod-gateway client will be created")
	flag.StringVarP(&clientContainerImage, "client-image", "i", provision.DefaultContainerImage, "The container image of the pod-gateway client")
	flag.StringSliceVarP(&clientCommand, "client-command", "c", provision.DefaultContainerCommand, "The command of the pod-gateway client container")
	flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	log.Printf("Provisioner started for pod-gateway clients with namespace=%s, image=%s, command=%s.", clientNamespace, clientContainerImage, clientCommand)

	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	k, err := k8s.BuildClient(kubeconfig)
	if err != nil {
		log.Fatalf("error building Kubernetes client: %s", err)
	}

	provisioner := provision.NewProvisioner(k, clientNamespace, clientContainerImage, clientCommand)

	log.Fatal(c.StartReceiver(context.Background(), provisioner.ReceiveCloudEvent))
}
