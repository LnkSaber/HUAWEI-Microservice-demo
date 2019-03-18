package microservice.demo.training21days.provider.service;



import lombok.Data;
import lombok.ToString;

import java.util.Date;

@Data
@ToString
public class GreetingResponse {
    private String msg;

    private Date timestamp;


}
