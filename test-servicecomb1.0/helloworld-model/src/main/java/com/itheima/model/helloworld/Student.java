package com.itheima.model.helloworld;

import lombok.Data;
import lombok.ToString;

import java.util.Date;

/**
 * @author Administrator
 * @version 1.0
 * @create 2018-08-13 10:37
 **/
@Data
@ToString
public class Student {
    String name;
    int age;
    String address;
    Date birthday;
}
