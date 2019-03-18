package com.itheima.hellworld;

import org.apache.servicecomb.springboot.starter.provider.EnableServiceComb;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * @author Administrator
 * @version 1.0
 * @create 2018-08-19 18:45
 **/
@EnableServiceComb
@SpringBootApplication
public class HelloWorldProviderApplication {
    public static void main(String[] args) {
        SpringApplication.run(HelloWorldProviderApplication.class);
    }
}
