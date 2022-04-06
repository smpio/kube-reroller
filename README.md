# kube-reroller


Automatically rerolls ReplicaSets, Deployments, StatefulSets by modifing Pod template annotation. Use it to schedule automatic pod restarts, pull "latest" image and so on.

## Usage

Label Deployment:

```
kubectl label deploy example k8s.smp.io/reroll-every=720h
```
