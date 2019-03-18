package microservice.demo.training21days.provider.service;

import lombok.Data;
import lombok.ToString;

@ToString
@Data
public class Person {
    private String name;
    private Gender gender;
}
