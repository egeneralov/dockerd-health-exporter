# dockerd-health-exporter

Just export this

```
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                    PORTS               NAMES
c6958fb713ca        test                "/bin/sh -c 'sleep i…"   14 minutes ago      Up 14 minutes (healthy)                       test
```

as

```
# HELP dockerd_container_health_test 1 = healthy
# TYPE dockerd_container_health_test gauge
dockerd_container_health_test{id="c6958fb713ca2dc04ced4889591022dc67a19d57442f288d0bc694c81bcb38f2",image="test",imageid="sha256:3d39e0c0c474f832f727ca934d7bb25ae4518e0aaea51876c7f6fcd56332b3b0"} 1
```

### docker hub
