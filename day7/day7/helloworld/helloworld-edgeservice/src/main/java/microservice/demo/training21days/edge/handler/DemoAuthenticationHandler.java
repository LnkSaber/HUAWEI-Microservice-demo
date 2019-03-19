package microservice.demo.training21days.edge.handler;

import microservice.demo.training21days.edge.filter.DemoAuthenticationFilter;
import org.apache.servicecomb.core.Handler;
import org.apache.servicecomb.core.Invocation;

import org.apache.servicecomb.swagger.invocation.AsyncResponse;
import org.apache.servicecomb.swagger.invocation.exception.InvocationException;

import javax.ws.rs.core.Response;

public class DemoAuthenticationHandler implements Handler {
    @Override
    public void handle(Invocation invocation, AsyncResponse asyncResp) throws Exception {
        String username = invocation.getContext(DemoAuthenticationFilter.USERNAME);
        String password = invocation.getContext(DemoAuthenticationFilter.PASSWORD);
        if (null == username || !username.equals(password)) {
            asyncResp.consumerFail(new InvocationException(Response.Status.UNAUTHORIZED, "Wrong authentication information"));
            return;
        }
        invocation.next(asyncResp);
    }
}
