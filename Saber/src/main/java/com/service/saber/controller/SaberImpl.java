package com.service.saber.controller;


import javax.ws.rs.core.MediaType;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import org.apache.servicecomb.provider.rest.common.RestSchema;

@javax.annotation.Generated(value = "io.swagger.codegen.languages.CseSpringDemoCodegen", date = "2019-03-26T01:23:27.134Z")

@RestSchema(schemaId = "saber")
@RequestMapping(path = "/Saber", produces = MediaType.APPLICATION_JSON)
public class SaberImpl {

    @Autowired
    private SaberDelegate userSaberDelegate;


    @RequestMapping(value = "/helloworld",
        produces = { "application/json" }, 
        method = RequestMethod.GET)
    public String helloworld( @RequestParam(value = "name", required = true) String name){

        return userSaberDelegate.helloworld(name);
    }

}
