package microservice.demo.training21days.provider.service;
import java.util.Date;


public class GreetingResponse {
    private String msg;

    private Date timestamp;

    @Override
    public String toString() {
        return "GreetingResponse{" +
                "msg='" + msg + '\'' +
                ", timestamp=" + timestamp +
                '}';
    }

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    public Date getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(Date timestamp) {
        this.timestamp = timestamp;
    }
}
