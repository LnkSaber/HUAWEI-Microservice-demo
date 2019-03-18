package microservice.demo.training21days.provider.service;

import org.apache.servicecomb.provider.rest.common.RestSchema;
import org.springframework.web.bind.annotation.*;

import java.util.Date;

@RestSchema(schemaId = "hello")        // 该注解声明这是一个REST接口类，CSEJavaSDK会扫描到这个类，根据它的代码生成接口契约
@RequestMapping(path = "/provider/v0") // @RequestMapping是Spring的注解，这里在使用Spring MVC风格开发REST接口
public class HelloService {
  @RequestMapping(path = "/hello", method = RequestMethod.GET)
  public String sayHello(@RequestParam(value = "name") String name) {
    return "Hello," + name;
  }


  @PostMapping(path = "/greeting")
  public GreetingResponse greeting(@RequestBody Person person){
    GreetingResponse greetingResponse = new GreetingResponse();
    if (Gender.MALE.equals(person.getGender())) {
      greetingResponse.setMsg("Hello, MR." + person.getName());
    }
    else {
      greetingResponse.setMsg("Hello, MS." + person.getName());
    }
    greetingResponse.setTimestamp(new Date());
    return greetingResponse;
  }


}
