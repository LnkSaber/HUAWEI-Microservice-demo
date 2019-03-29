package com.service.saber.controller;



import static org.junit.Assert.*;
import org.junit.Before;
import org.junit.Test;

public class TestSaber {

        SaberDelegate saberDelegate = new SaberDelegate();


    @Test
    public void testhelloworld(){

        String expactReturnValue = "hello"; // You should put the expect String type value here.

        String returnValue = saberDelegate.helloworld("hello");

        assertEquals(expactReturnValue, returnValue);
    }

}