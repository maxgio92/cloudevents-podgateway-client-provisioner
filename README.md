# Event-driven pod-gateway client provisioner

> Disclaimer: this is a proof of concept.

This service provisions pods for which traffic needs to be routed through a `pod-gateway`.

The provisioning is event-driven using CloudEvents.

The event must be of type `io.podgateway.client.pending` and its content must specify the fields:
- `gateway_name` (string)

The value must reference a valid [pod-gateway](https://github.com/angelnu/pod-gateway/)'s `setGateway` label/annotation value, as configured in its [gateway-admission-controller](https://github.com/angelnu/gateway-admision-controller).

> More on the admission controller configuration [here](https://github.com/angelnu/gateway-admision-controller/blob/main/internal/config/config.go).

```ascii
                            ┌───────────────────┐
                            │                   │
                            │  gtw mutating     │
                            │  admission        │
                            └─────┬─┬─┬────┬─┬─┬┘
                                  │ │ │    │ │ │
                                ┌─▼─┴─┴──┐ │ │ │   ┌────────────┐
                                │client  │ │ │ │   │            │
                           ┌───►│        │ ▼ │ │   │ gateway    │
┌──────┐                   │    │gtw=foo ├───┴─┴───┤ foo        │
│      │                   │    │        │ tunnel  │            ├────►
│events│  ┌────────────┐   │    └───┬─┬──┘   │ │   │            │
│      │  │            │   │        │ │      │ │   │            │
│      │  │ provisioner├───┘    ┌───▼─┴──┐   │ │   └────────────┘
│      ├─►│            │        │client  │   │ │
│      │  │            ├───────►│        │   ▼ │   ┌────────────┐
│      │  │            │        │gtw=bar ├─────┴───┤            │
│      │  │            ├───┐    │        │ tunnel  │ gateway    │
│      │  │            │   │    └─────┬──┘     │   │ bar        ├────►
│      │  └────────────┘   │          │        │   │            │
│      │                   │    ┌─────▼──┐     ▼   │            │
└──────┘                   │    │client  ├─────────┤            │
                           └───►│        │ tunnel  └────────────┘
                                │gtw=bar │
                                │        │
                                └────────┘
```

## Usage

```shell
cloudevents-podgateway-client-provisioner [--client-namespace=<client namespace>] [--client-command=<command>] [--client-image=<client container image>]
```

## Quickstart

As the only supported event spec is CloudEvents, a quickstart setup can be configured with Knative.

All of that will run in a local Kubernetes cluster.

Deploy a KinD cluster with Knative Eventing and Service components locally:

```shell
kn quickstart kind
```

Deploy pod-gateways (e.g. named `foo` and `bar`):

```shell

helm upgrade --install -n gateway-system --create-namespace pod-gateway-foo angelnu/pod-gateway -f $deploydir/pod-gateway-foo-values.yaml --version 6.1.0
helm upgrade --install -n gateway-system --create-namespace pod-gateway-bar angelnu/pod-gateway -f $deploydir/pod-gateway-bar-values.yaml --version 6.1.0
```

Deploy a Knative Broker for the CloudEvents:

```shell
kubectl apply -f deploy/namespace.yaml
kubectl apply -f deploy/broker.yaml
```

Deploy the provisioner as a Knative Service:

```shell
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/service.yaml
```

(optional) Deploy a CloudEvents dashboard:

```shell
kubectl apply -f deploy/cloudevents-player.yaml
```

and open the browser at http://cloudevents-player.client-system.127.0.0.1.sslip.io.

You can now send events of Type `io.podgateway.client.pending`, and specify the gateway for the client in a field `gateway_name`.
The value must reference an installed pod-gateway.

> Specifically, the value must match the pod-gateway [admission controller](https://github.com/angelnu/gateway-admision-controller)'s `setGatewayLabelValue`/`setGatewayAnnotationValue` flag.

For example:

- ID: *generated*
- Type: `io.podgateway.client.pending`
- Source: `mySource`
- SpecVersion: 1.0
- Message:
  ```json
  {
   "gateway_name": "foo"
  }
  ```

The event will trigger (see [here](deploy/service.yaml)) the provisioner Service that will create the client pod.

Finally, the provisioner will notify the success or failure of that operation, with a `io.podgateway.client.scheduling.done` or `io.podgateway.client.scheduling.failed` event.

On success, The `io.podgateway.client.scheduling.done` event will contain Data of the client Pod just created, such as:
- `pod_name`
- `namespace`

On failure, The `io.podgateway.client.scheduling.failed` event will contain the related error message.

## Development

### Build

```shell
make build
```

### Publish the OCI image to a local registry

```shell
make publish/local
```

