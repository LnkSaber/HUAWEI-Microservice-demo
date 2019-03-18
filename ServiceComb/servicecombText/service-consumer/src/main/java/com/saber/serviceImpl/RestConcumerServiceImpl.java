package com.saber.serviceImpl;

import com.saber.service.RestService;
import org.apache.servicecomb.provider.springmvc.reference.RestTemplateBuilder;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class RestConcumerServiceImpl implements RestService {

    private final RestTemplate restTemplate = RestTemplateBuilder.create();
    @Override
    public String sayRest(String name) {
        String provideName = "provider";
        //URL : cse:// 是协议； + 微服务名称  + 具体访问路径
        return  restTemplate.getForObject("cse://"+provideName+"/hello/hello?name="+name, String.class);

    }
}
