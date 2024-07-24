# prettyql

Feed it an ugly and poorly-written query, it spits out a pretty and poorly-written query

```shell
echo 'label_replace(instance:node_memory_utilisation:ratio{cluster="mycluster"}, "node", "$1", "instance", "(.+):metrics") * on(node) group_left(customer) sum by(node) (label_replace(kube_pod_info{cluster="mycluster",pod=~".+mypod.+"},"customer", "$1", "namespace", "(.+)"))' | ./prettyql
  label_replace(
    instance:node_memory_utilisation:ratio{cluster="mycluster"},
    "node",
    "$1",
    "instance",
    "(.+):metrics"
  )
* on (node) group_left (customer)
  sum by (node) (
    label_replace(
      kube_pod_info{cluster="mycluster",pod=~".+mypod.+"},
      "customer",
      "$1",
      "namespace",
      "(.+)"
    )
  )
```
