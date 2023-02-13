package provision

type SocialNetworkName string

type AvatarID string

type GatewayName string

// GatewaySpec defines the Pod Gateway specification.
type GatewaySpec struct {

	// GatewayName holds the Pod Gateway name to be used for the Client's requests.
	// It must match the value of the label/annotation selector specified to the pod-gateway
	// admission controller.
	GatewayName GatewayName `json:"gateway_name"`
}

// ClientPending defines the Data of CloudEvent with type=io.podgateway.client.pending
type ClientPending struct {

	// Message holds the message from the event
	Message string `json:"message,omitempty"`

	GatewaySpec
}

// ClientSchedulingDone defines the Data of CloudEvent with type=io.podgateway.client.scheduling.done
type ClientSchedulingDone struct {

	// Message holds the message from the event
	Message string `json:"message"`

	// PodName holds the name of the Kubernetes Client Pod
	PodName string `json:"pod_name"`

	// Namespace holds the Kubernetes Namespace where the Client Pod has been created
	Namespace string `json:"namespace"`

	// Request holds the details of the Client's request
	GatewaySpec
}

// ClientSchedulingFailure defines the Data of CloudEvent with type=io.podgateway.client.scheduling.failed
type ClientSchedulingFailure struct {

	// Message holds the message from the event
	Message string `json:"message"`

	// Request holds the details of the Client's request
	GatewaySpec
}
