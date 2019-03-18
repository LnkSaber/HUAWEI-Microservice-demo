package com.itheima.helloworld.consumer;

import org.apache.servicecomb.springboot.starter.provider.EnableServiceComb;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * @author Administrator
 * @version 1.0
 * @create 2018-08-13 10:03
 **/
@EnableServiceComb
@SpringBootApplication
public class HelloWorldConsumerApplication {
    public static void main(String[] args) {
        SpringApplication.run(HelloWorldConsumerApplication.class);
    }
}
