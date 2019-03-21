package microservice.demo.training21days.edge.filter;

import org.apache.servicecomb.common.rest.filter.HttpServerFilter;
import org.apache.servicecomb.core.Invocation;
import org.apache.servicecomb.foundation.vertx.http.HttpServletRequestEx;
import org.apache.servicecomb.swagger.invocation.Response;
import org.springframework.util.StringUtils;

public class DemoFilter implements HttpServerFilter {

  private static final String LET_STRANGER_PASS = "LetStrangerPass";

  @Override
  public int getOrder() {
    return 0;
  }

  @Override
  public Response afterReceiveRequest(Invocation invocation, HttpServletRequestEx httpServletRequestEx) {
    // 从请求中取出一个header
    String letStrangerPass = httpServletRequestEx.getHeader(LET_STRANGER_PASS);
    if (!StringUtils.isEmpty(letStrangerPass)) {
      // 如果此header存在则将它存入到InvocationContext中，InvocationContext可以从上游服务自动传递给下游服务
      invocation.addContext(LET_STRANGER_PASS, letStrangerPass);
    }
    return null;
  }
}
