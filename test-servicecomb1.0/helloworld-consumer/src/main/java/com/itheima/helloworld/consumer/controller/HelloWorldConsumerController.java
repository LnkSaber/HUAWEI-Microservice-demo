package com.itheima.helloworld.consumer.controller;

import com.itheima.helloworld.api.HelloWorldInterface;
import com.itheima.model.helloworld.Student;
import org.apache.servicecomb.provider.pojo.RpcReference;
import org.apache.servicecomb.provider.rest.common.RestSchema;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

/**
 * @author Administrator
 * @version 1.0
 * @create 2018-08-19 18:51
 **/
@RestSchema(schemaId = "helloWorldConsumer")
@RequestMapping("/")
public class HelloWorldConsumerController   {

    //基于rpc的方式调用远程rest接口
    @RpcReference(microserviceName="helloworld-provider",schemaId = "helloworld")
    HelloWorldInterface helloWorldInterface;

    @GetMapping(path = "request")
    public Student request(String name){
        //远程调用helloworld-provider的/hello方法
        Student student = helloWorldInterface.hello(name);
        return student;
    }

}
