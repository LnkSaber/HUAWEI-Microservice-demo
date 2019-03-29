package hello.service;


import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestTemplate;


@RestController
@RequestMapping("/consumer/v0")
public class ConsumerController {
    @RequestMapping("/hello/{name}")
    public String ProviderHello(@PathVariable(value = "name") String name) {
        RestTemplate restTemplate = new RestTemplate();
        String s = restTemplate.getForObject("http://provider/provider/v0/hello/" + name, String.class);
        return s;
    }


    @RequestMapping("/greeting")
    public GreetingResponse ConsumerGreeting(@RequestBody Person p) {
        RestTemplate restTemplate = new RestTemplate();
        GreetingResponse response = restTemplate.postForObject("http://provider/provider/v0/greeting", p, GreetingResponse.class);
        return response;
    }


}
