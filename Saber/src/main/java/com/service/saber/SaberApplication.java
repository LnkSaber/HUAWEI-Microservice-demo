package com.service.saber;

import org.apache.servicecomb.springboot.starter.provider.EnableServiceComb;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@EnableServiceComb
public class SaberApplication {
    public static void main(String[] args) {
         SpringApplication.run(SaberApplication.class,args);
    }
}
