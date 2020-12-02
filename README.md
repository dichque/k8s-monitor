# Kubernetes Cluster Failure Detector Solution

## Control plane services that must be monitored:
1. etcd [http://127.0.0.1:2381/health] [http://127.0.0.1:2381/metrics]
2. api server [https://127.0.0.1:6443/livez]
3. controller [http://127.0.0.1:10257/healthz]
4. scheduler [https://127.0.0.1:10259/healthz]
5. DNS
6. Third party resources or service endpoints: Load balancer / proxy ?

## Scope:

### Availability
1. Service availability can be checked by validating endpoint health
2. Pods can be checked by validating status field, restartcount

### Capacity
1. CPU, Memory
2. Network
3. Storage Volume

### Performance
1. cadvisor
2. metrics server

## Solution 1: Prometheus centric
Instrument k8s cluster and ensure all k8s components export health state, capacity & performance metrics for prometheus to ingest. Automate deployment & configuration prometheus & alertmanager as part of cluster provisioning.
Leverage prometheus operator for automation.

Pros:
1. Out of the box setup, includes monitoring, alerting & dashboard using grafana.
2. Relativily less effort to rollout & implement.
3. Custom components can be monitored by instrumenting using  prometheus client libraries.

Cons:
1. Scaling & tuning prometheus is complex.
2. Will not be able to monitor external critical systems


## Solution 2: [https://kubernetes.io/docs/tasks/debug-application-cluster/events-stackdriver/](Events)
All activity of k8s objects are logged to events, which can be exported to third party log/event aggregation, messaging & alerting systems. Develop or leverage open source event exporter to filter & act on critical events

Pros:
1. Allows integration with several systems for analysis, alerting & archival
2. Ability to filter events

Cons:
1. Volume of data is huge, got to be careful of whats worth acting on


## Solution 3: Custom / Home grown
Develop in-house monitoring & alerting system. Leverage k8s controllers & custom resource mechanism for resilency & scaling. Leverage messaging system to agregate alerts & act on them. 

Pros:
1. Fully customized solution, can be focused only few components that needs monitoring

