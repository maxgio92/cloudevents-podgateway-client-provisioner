package provision

const (
	DefaultNamespace = "default"

	ContainerName         = "client"
	DefaultContainerImage = "alpine/curl"

	KubernetesLabelAppName         = "app.kubernetes.io/name"
	KubernetesLabelSetGateway      = "setGateway"
	KubernetesAnnotationSetGateway = "setGateway"

	Source = "podgateway/eventing/client/provisioner"

	TypeClientPending           = "io.podgateway.client.pending"
	TypeClientSchedulingDone    = "io.podgateway.client.scheduling.done"
	TypeClientSchedulingFailure = "io.podgateway.client.scheduling.failed"

	ClientSchedulingDoneMessage = "A pod-gateway client has been scheduled."
)

var (
	DefaultContainerCommand = []string{"sleep", "infinity"}
)
