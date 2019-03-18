package com.itheima.hellworld.controller;

import com.itheima.helloworld.api.HelloWorldInterface;
import com.itheima.model.helloworld.Student;
import org.apache.servicecomb.provider.rest.common.RestSchema;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

/**
 * @author Administrator
 * @version 1.0
 * @create 2018-08-19 18:40
 **/
@RestSchema(schemaId = "helloworld")
@RequestMapping("/")
public class HelloWorldController implements HelloWorldInterface {

    @GetMapping(path = "hello")
    @Override
    public Student hello(String name) {
        Student student = new Student();
        student.setName(name);
        return student;
    }
}
