package com.saber.serviceImpl;

import com.saber.service.RpcService;
import org.apache.servicecomb.provider.pojo.RpcReference;
import org.springframework.stereotype.Component;

@Component
public class RpcConsumerServiceImpl implements RpcService {
    @RpcReference(microserviceName = "start.servicecomb.io:provider-rpc",schemaId = "helloRpc")
    private RpcService rpcService;
    @Override
    public String sayRpc(String name) {
        return rpcService.sayRpc(name);
    }
}
