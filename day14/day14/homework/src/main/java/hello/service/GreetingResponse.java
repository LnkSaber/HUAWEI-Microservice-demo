package hello.service;

import java.util.Date;

public class GreetingResponse {
    private String msg;
    private Date timestamp;

    public String getMsg() { return msg; }

    public void setMsg(String msg) { this.msg = msg; }

    public Date getTimestamp() { return timestamp; }

    public void setTimestamp(Date timestamp) { this.timestamp = timestamp; }
}
