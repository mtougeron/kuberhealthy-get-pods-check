## Kuberhealthy Get Pods Check

![Go](https://github.com/mtougeron/kuberhealthy-get-pods-check/workflows/Go/badge.svg) ![Gosec](https://github.com/mtougeron/kuberhealthy-get-pods-check/workflows/Gosec/badge.svg) [![GitHub tag](https://img.shields.io/github/tag/mtougeron/kuberhealthy-get-pods-check.svg)](https://github.com/mtougeron/kuberhealthy-get-pods-check/tags/)

The `Kuberhealthy Get Pods Check` checks if the API servers returns a list of pods under a specified time-limit.

## Thanks Comcast!

A big shout-out and thank you goes to Comcast for writing [Kuberhealthy](https://github.com/Comcast/kuberhealthy)

## Kuberhealthy AWS IAM Role Check Kube Spec Example

```yaml
apiVersion: comcast.github.io/v1
kind: KuberhealthyCheck
metadata:
  name: get-pods
spec:
  runInterval: 5m
  timeout: 1m
  podSpec:
    containers:
    - name: main
      image: ghcr.io/mtougeron/khcheck-get-pods:latest
      imagePullPolicy: IfNotPresent
      env:
        - name: NAMESPACE
          value: "kuberhealthy"
        - name: MAX_DURATION_MILLISECONDS
          value: "250"
        - name: DEBUG
          value: "1"
```
where `NAMESPACE` is the Kubernetes namespace to get the pods from.

### Installation

>Make sure you are using the latest release of Kuberhealthy 2.2.0.

Run `kubectl apply` against [example spec file](example/khcheck-get-pods.yaml)

```bash
kubectl apply -f khcheck-get-pods.yaml -n kuberhealthy
```
#### Container Image

Image is available [Github Container Registry](https://github.com/users/mtougeron/packages/container/khcheck-get-pods/)

### Licensing

This project is licensed under the Apache V2 License. See [LICENSE](LICENSE) for more information.
