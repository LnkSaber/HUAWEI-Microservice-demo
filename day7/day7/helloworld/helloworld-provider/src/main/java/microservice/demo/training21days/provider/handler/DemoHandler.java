package microservice.demo.training21days.provider.handler;

import javax.ws.rs.core.Response.Status;

import org.apache.servicecomb.core.Handler;
import org.apache.servicecomb.core.Invocation;
import org.apache.servicecomb.swagger.invocation.AsyncResponse;
import org.apache.servicecomb.swagger.invocation.exception.CommonExceptionData;
import org.apache.servicecomb.swagger.invocation.exception.InvocationException;

public class DemoHandler implements Handler {
  @Override
  public void handle(Invocation invocation, AsyncResponse asyncResp) throws Exception {
    // 从这里可以取出本次请求调用的方法的完整名字，格式是 serviceName.schemaId.operationId
    String operationName = invocation.getOperationMeta().getMicroserviceQualifiedName();
    // 这里我们只检查sayHello方法的参数
    if ("provider.hello.sayHello".equals(operationName)) {
      Object name = invocation.getSwaggerArgument(0);
      // 如果name=stranger，则拒绝请求，返回403
      if (!"true".equalsIgnoreCase(invocation.getContext("LetStrangerPass"))
          && "stranger".equalsIgnoreCase((String) name)) {
        asyncResp.producerFail(new InvocationException(Status.FORBIDDEN, new CommonExceptionData("Don't know you :(")));
        return;
      }
    }
    // 通过检查，继续执行后面的逻辑
    invocation.next(asyncResp);
  }
}
