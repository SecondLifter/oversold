## oversold

> k8s资源超售，通过调整节点资源的超卖比，在可用资源上乘上一个系数来呈现一个超卖之后的节点。目前 Kubernetes 是没有这样的接口，通过 Kubernetes 的扩展机制，Hook 了 Kubelet 向上汇报的过程，实现了该功能

### 一、简介

基于[goadmission](https://github.com/mritd/goadmission) 一个 Kubernetes 动态准入控制的脚手架，进行开发完成



### 二 、背景

公司当前给业务容器进行资源分配都是采 用预分配加上 Cgroup 限制的手段，但是在实际进行服务利用率统计的时候，发现大多数业务并不了解自己的服务是什么样子，这就会导致一些问题：

1. 实际资源利用低。根据统计两周内对服务的最高资源使用数据表示，百分之八十的业务都存在百分之五十以上的资源浪费 (CPU 和内存) 。
2. 特定资源集群的利用率低，公司内有专门的开发,测试,预发等集群，这部分集群的资源利用在闲时能有 80% 的资源是处于空闲状态。

### 三、如何使用

克隆本项目到本地，在本地使用docker 对代码进行容器构建。再到k8s集群进行部署
> 前置条件 
>需要cfssl 命令行工具

```bash
cd oversold
docker build -t goadmission:v0.9 . #构建镜像
kubectl create ns oversold  #创建命名空间

sh oversold/deploy/cfssl/create.sh #生成密钥
cd oversold/cfssl/mutatingwebhook/
kubectl apply -f .
```

### 四、如何开启

 ```shell
 kubectl label --overwrite node --all kubernetes.io/oversold=oversold  #若只需要部分节点进行超售。只需要对需要超售节点进行打标
 kubectl label --overwrite node --all kubernetes.io/overcpu=3  ### 超卖倍数
 kubectl label --overwrite node --all kubernetes.io/overmem=3  #### 超卖倍数 可只超卖cpu 或meme

 ```
 ### 
