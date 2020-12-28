# resource-collector
使用dcgm-exporter，将GPU资源过滤导出(基于kubernetes平台之上)


参考[gpu-monitor-tools](https://github.com/NVIDIA/gpu-monitoring-tools#quickstart-on-kubernetes)
使用helm将`dcgm-exporter`部署到kubernetes的default ns中。
在部署过程中，注意[Profiling is not supported for this group of GPUs or GPU](https://github.com/NVIDIA/gpu-monitoring-tools/issues/119#issuecomment-722885536)
因为目前profiling metrics只支持TeslaV100和T4

`dcgm-exporter.yaml`也部署后，使用`port-forward`把端口导出

比如我的案例：
```
kubectl port-forward dcgm-exporter-1608712831-t75wd --address=172.18.29.80 9407:9400 &
...

并且之前也给节点打过label
kubectl label --overrides nodes node4 accelerator=A100
```

参考之前的项目[k8sMLer-client-go中commit](https://github.com/ReyRen/k8sMLer-client-go/commit/751c5605fad2cbfc30fc0da95787e441ee5f95de)
这个功能点。

在该项目中也会有对应的代码部分[根据标签名指定端口](https://github.com/ReyRen/resource-collector/blob/288be2c289c88a377711a35bc598e76767041612/common.go#L45)

### 注意
因为使用port-forward作为后台程序发现总会被OOM了。目前没有找到合适的办法，后期空闲时候会进行更改。workaround是写了个脚本
`crontab.sh`然后加入到crontab中
```
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
*/30 * * * * /root/go/src/resource-collector/crontab.sh
```
所以目前如果DCGM-exporter名字发生变化需要进行脚本的修改