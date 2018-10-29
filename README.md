# JVessel

京舰（JVessel）系统是基于京东云 IaaS 基础服务，为 PaaS 产品提供统一调度管理平台。


## 1. Usage

这里介绍paas开发者如何使用JVessel  

**对于paas开发者，只需要看tutorial即可**
[tutorial](./tutorial )


## 2. Architecture

内部模块介绍

- [matrix](./doc/matrix_api.md )
- [controller](./doc/controler.md)
- [worker](./doc/worker.md )
- [probe](./doc/probe.md )
- [dns](./doc/dns.md )
- [log system](./doc/log_system.md)

## 3. Design 

一些跨不同module共同实现的功能的设计  

- [container实现标准](./doc/container_spec.md)
- [关于az](./doc/az-ops.md)
- [etcd存储](./doc/etcd_storage.md)
- [可用性设计](./doc/design_availablility.md)
- [云盘设计](./doc/design_cloud_volume.md)
- [安全组设计](./doc/design.md)
- [管理流量透传](./doc/underlayentry_design.png)
- [az设计](./doc/az_design.png)
- [pipeline设计](./doc/pipeline.png)

## 4. Test

- 单元测试, 函数代码级别测试 
   make unit-test
- 资源测试，对内部各种资源调度策略的测试，在单机上启动进程测试即可  
   [resource_test](./test/resource_test)
- 单独模块测试，独立模块的集成测试  
   [module_test](./test/module_test)
- 系统集成测试，整个系统运行健壮性的测试  
   [system_test](./test/system_test)
- smoke test，保证系统基本功能都正常运行  
   [smoke_test](./test/smoke_test)
- stability test，持续验证系统的正确性和稳定性  
   [stability_test](./test/stability_test)
- 性能测试  
   [benchmark](./test/benchmark)

## 5. Deployment && Ops

系统部署描述：[deployment](./doc/deployment.md)

az部署描述：[az_deploy.md](./doc/az_deploy.md)

系统运维： [ops](./doc/ops.md)

线上运维工具path： /export/Instance/xxx-jvessel-paas/x.xxx-jvessel-paas/runtime/bin/jvesctl

## 6. PASS产品接入Jvessel系统流程文档
http://cf.jd.com/pages/viewpage.action?pageId=106287929
