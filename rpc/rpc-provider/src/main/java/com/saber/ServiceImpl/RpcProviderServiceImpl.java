package com.saber.ServiceImpl;

import com.saber.service.RpcService;
import org.apache.servicecomb.provider.pojo.RpcSchema;

@RpcSchema(schemaId = "helloRpc")
public class RpcProviderServiceImpl implements RpcService {
    @Override
    public String sayRpc(String name) {
        return "hello" +name;
    }
}
