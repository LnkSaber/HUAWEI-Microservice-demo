package microservice.demo.training21days.edge.filter;

import javax.ws.rs.core.Response.Status;

import org.apache.servicecomb.common.rest.filter.HttpServerFilter;
import org.apache.servicecomb.core.Invocation;
import org.apache.servicecomb.foundation.vertx.http.HttpServletRequestEx;
import org.apache.servicecomb.swagger.invocation.Response;
import org.apache.servicecomb.swagger.invocation.exception.CommonExceptionData;
import org.apache.servicecomb.swagger.invocation.exception.InvocationException;
import org.springframework.util.StringUtils;


public class DemoAuthenticationFilter implements HttpServerFilter {
    public static final String USERNAME ="Username";
    public static final String PASSWORD = "Password";

    @Override
    public int getOrder() {
        return 0;
    }

    @Override
    public Response afterReceiveRequest(Invocation invocation, HttpServletRequestEx requestEx) {
        String username = requestEx.getHeader(USERNAME);
        String password = requestEx.getHeader(PASSWORD);
        if (StringUtils.isEmpty(username) || StringUtils.isEmpty(password)) {
            return Response.consumerFailResp(
                    new InvocationException(Status.UNAUTHORIZED, new CommonExceptionData("Lack of authentication information"))
            );
        }
        invocation.addContext(USERNAME, username);
        invocation.addContext(PASSWORD, password);
        return null;
    }
}
