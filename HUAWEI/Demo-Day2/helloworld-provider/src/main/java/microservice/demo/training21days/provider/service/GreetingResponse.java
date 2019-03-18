package microservice.demo.training21days.provider.service;

import lombok.Data;
import lombok.ToString;

import java.util.Date;
@ToString
@Data
public class GreetingResponse {
    private String msg;

    private Date timestamp;


}
