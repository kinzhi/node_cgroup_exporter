# Node Cgroup exporter

由于当前在用 exporter 中没有 `cgroup` 的 `memory.usage_in_bytes` 的相关指标；所以以此 exporter 作为补充使用。



## 快速启动
```bash
docker run -it -p 9100:9100 node_cgroup_exporter:latest 
```

## 部署在 Kubernetes
```bash
kubectl apply -f deploy/*
```

## 背景
Kubernetes 官方[文档](https://kubernetes.io/zh-cn/examples/admin/resource/memory-available.sh)关于 Kubelet 驱逐逻辑的计算公式：
```
memory.available = memory_capacity_in_bytes - memory_working_set(memory_usage_in_bytes - 
memory_total_inactive_file)
```

官方检测脚本：
```bash
#!/bin/bash
#!/usr/bin/env bash

# This script reproduces what the kubelet does
# to calculate memory.available relative to root cgroup.

# current memory usage
memory_capacity_in_kb=$(cat /proc/meminfo | grep MemTotal | awk '{print $2}')
memory_capacity_in_bytes=$((memory_capacity_in_kb * 1024))
memory_usage_in_bytes=$(cat /sys/fs/cgroup/memory/memory.usage_in_bytes)
memory_total_inactive_file=$(cat /sys/fs/cgroup/memory/memory.stat | grep total_inactive_file | awk '{print $2}')

memory_working_set=${memory_usage_in_bytes}
if [ "$memory_working_set" -lt "$memory_total_inactive_file" ];
then
    memory_working_set=0
else
    memory_working_set=$((memory_usage_in_bytes - memory_total_inactive_file))
fi

memory_available_in_bytes=$((memory_capacity_in_bytes - memory_working_set))
memory_available_in_kb=$((memory_available_in_bytes / 1024))
memory_available_in_mb=$((memory_available_in_kb / 1024))

echo "memory.capacity_in_bytes $memory_capacity_in_bytes"
echo "memory.usage_in_bytes $memory_usage_in_bytes"
echo "memory.total_inactive_file $memory_total_inactive_file"
echo "memory.working_set $memory_working_set"
echo "memory.available_in_bytes $memory_available_in_bytes"
echo "memory.available_in_kb $memory_available_in_kb"
echo "memory.available_in_mb $memory_available_in_mb"

```

### 指标

Name     | Description | Exporter | Metric
---------|-------------|-------------|----
memory_capacity_in_bytes | 节点总内存 | Node_exporter ｜ node_memory_MemTotal_bytes
memory_working_set | 节点在用内存 | 需计算 | /
memory_usage_in_bytes | 节点在用内存(包含节点进程使用、进程缓存等) | Node_cgroup_exporter(新增) | node_cgroupMem_usage
memory_total_inactive_file | 节点inactive缓存(当前exporter新增) | Node_exporter | node_memory_Inactive_file_bytes


### 告警规则

根据 kubelet 驱逐逻辑制定告警规则：
```
sum(label_replace(node_cgroupMem_usage,"nodeip","$1","instance","(.*):.*"))by(nodeip)  / (sum(label_replace(node_memory_MemTotal_bytes,"nodeip","$1","instance","(.*):.*"))by(nodeip) + sum(label_replace(node_memory_Inactive_file_bytes,"nodeip","$1","instance","(.*):.*"))by(nodeip)) > 0.8
```

## 构建

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

```
docker build --platform linux/amd64 -t harbor.harbor.cn/etc/ww/node_cgroup_exporter:v1.1 . --no-cache
```

```
docker save harbor.harbor.cn/etc/ww/node_cgroup_exporter:v1.1 -o cgroup-exporter.tar
```