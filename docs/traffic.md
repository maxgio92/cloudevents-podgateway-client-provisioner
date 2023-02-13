## Inspect the trafic

The client pod's egress traffic can be inspected in the selected Pod Gateway:

```shell
kubectl exec -n gateway-system -it pod-gateway-foo -- sh
$ apk add tcpdump
$ tcpdump -i eth0 -nv
tcpdump: listening on eth0, link-type EN10MB (Ethernet), snapshot length 262144 bytes

```

and sending ICMP traffic to Internet, from the client pod.

You can retrieve the namespace and the name of the pod from the `io.podgateway.client.scheduling.done` event message.

```shell
kubectl exec -n <namespace> -it <client_pod_name> -- sh
$ ping -c1 8.8.8.8
```

it can be verified with `tcpdump` run above, that it comes from the *eth0* Pod Gateway interface:

```shell
[...]
$ tcpdump -i eth0 -nv
tcpdump: listening on eth0, link-type EN10MB (Ethernet), snapshot length 262144 bytes
14:24:21.791574 IP (tos 0x0, ttl 63, id 45165, offset 0, flags [none], proto UDP (17), length 134)
    10.244.0.31.43342 > 10.244.0.26.8472: OTV, flags [I] (0x08), overlay 0, instance 42
IP (tos 0x0, ttl 64, id 64442, offset 0, flags [DF], proto ICMP (1), length 84)
    172.16.0.184 > 8.8.8.8: ICMP echo request, id 68, seq 0, length 64
14:24:21.791618 IP (tos 0x0, ttl 63, id 64442, offset 0, flags [DF], proto ICMP (1), length 84)
    10.244.0.26 > 8.8.8.8: ICMP echo request, id 68, seq 0, length 64
14:24:21.804365 IP (tos 0x0, ttl 113, id 0, offset 0, flags [none], proto ICMP (1), length 84)
    8.8.8.8 > 10.244.0.26: ICMP echo reply, id 68, seq 0, length 64
14:24:21.804393 IP (tos 0x0, ttl 64, id 37367, offset 0, flags [none], proto UDP (17), length 134)
    10.244.0.26.43342 > 10.244.0.31.8472: OTV, flags [I] (0x08), overlay 0, instance 42
IP (tos 0x0, ttl 112, id 0, offset 0, flags [none], proto ICMP (1), length 84)
    8.8.8.8 > 172.16.0.184: ICMP echo reply, id 68, seq 0, length 64
```

