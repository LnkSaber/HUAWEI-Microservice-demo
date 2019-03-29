package hello.service;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/provider/v0")
public class ProviderController {

    @RequestMapping("/hello/{name}")
    public String ConsumerHello(@PathVariable(value = "name") String name) {
        return "hello , " + name;
    }

}
