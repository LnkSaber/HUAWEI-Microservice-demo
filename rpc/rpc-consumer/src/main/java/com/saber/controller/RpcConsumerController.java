package com.saber.controller;

import com.saber.service.RpcService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class RpcConsumerController {
    @Autowired
    private RpcService rpcService;

    @RequestMapping("/rpc")
    public void rpcInvoke(){
        System.out.println(rpcService.sayRpc("servicercomb rpc"));
    }
}
