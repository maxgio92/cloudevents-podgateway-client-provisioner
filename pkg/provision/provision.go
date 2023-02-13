package provision

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type Provisioner struct {
	Kubernetes    *kubernetes.Clientset
	ClientOptions ClientOptions
}

type ClientOptions struct {
	Namespace string
	Image     string
	Command   []string
}

// NewProvisioner returns a new Provisioner object to provision client pods.
func NewProvisioner(kubernetes *kubernetes.Clientset, clientNamespace, clientImage string, clientCommand []string) *Provisioner {
	return &Provisioner{
		Kubernetes: kubernetes,
		ClientOptions: ClientOptions{
			Namespace: clientNamespace,
			Image:     clientImage,
			Command:   clientCommand,
		},
	}
}

func (p *Provisioner) ReceiveCloudEvent(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	log.Printf("Event %s received. \n%s\n", event.Type(), event)

	switch event.Type() {
	case TypeClientPending:
		eventData := &ClientPending{}
		if err := event.DataAs(eventData); err != nil {
			return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
		}
		log.Printf("A new Client request is pending: %s", eventData.Message)

		gatewayName := string(eventData.GatewayName)

		// Schedule the podgateway-client.
		if gatewayName == "" {
			return sendSchedulingFailedEvent(fmt.Errorf("A Client scheduling failed: gatewayName name is empty"))
		}

		podName, err := p.ScheduleClient(gatewayName)
		if err != nil {
			return sendSchedulingFailedEvent(fmt.Errorf("A Client scheduling failed: %s", err))
		}

		return sendSchedulingDoneEvent(podName, p.ClientOptions.Namespace, fmt.Sprintf("A Client on %s gatewayName has been scheduled.", gatewayName))
	}
	log.Printf("Event type %s not supported.", event.Type())

	return nil, nil
}

func (p *Provisioner) ScheduleClient(gatewayName string) (string, error) {
	podName := uuid.New().String()

	err := ensureClientPod(p.Kubernetes, podName, p.ClientOptions.Namespace, p.ClientOptions.Image, p.ClientOptions.Command, gatewayName)
	if err != nil {
		return "", err
	}

	log.Printf("created pod %s/%s", p.ClientOptions.Namespace, podName)

	return podName, nil
}
