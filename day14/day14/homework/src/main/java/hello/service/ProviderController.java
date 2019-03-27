package hello.service;

import org.springframework.web.bind.annotation.*;

import java.util.Date;

@RestController
@RequestMapping("/provider/v0")
public class ProviderController {

    @RequestMapping("/hello/{name}")
    public String ConsumerHello(@PathVariable(value = "name") String name) {
        return "hello , " + name;
    }

    @PostMapping("/greeting")
    public GreetingResponse providerGreeting(@RequestBody Person person) {
        GreetingResponse greetingResponse = new GreetingResponse();
        if (Gender.MALE.equals(person.getGender())) {
            greetingResponse.setMsg("ConsumerHello, Mr." + person.getName());
        } else {
            greetingResponse.setMsg("ConsumerHello, Ms." + person.getName());
        }
        greetingResponse.setTimestamp(new Date());

        return greetingResponse;
    }
}
