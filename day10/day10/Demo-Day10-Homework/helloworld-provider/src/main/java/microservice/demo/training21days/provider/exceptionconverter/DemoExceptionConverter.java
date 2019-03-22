package microservice.demo.training21days.provider.exceptionconverter;

import javax.ws.rs.core.Response.Status;

import org.apache.servicecomb.swagger.invocation.Response;
import org.apache.servicecomb.swagger.invocation.SwaggerInvocation;
import org.apache.servicecomb.swagger.invocation.exception.CommonExceptionData;
import org.apache.servicecomb.swagger.invocation.exception.ExceptionToProducerResponseConverter;
import org.apache.servicecomb.swagger.invocation.exception.InvocationException;

public class DemoExceptionConverter implements ExceptionToProducerResponseConverter<IllegalArgumentException> {
  /**
   * 定义优先级
   */
  @Override
  public int getOrder() {
    return 0;
  }

  /**
   * 该方法的返回值表明DemoExceptionConverter处理IllegalArgumentException及其子类的异常。
   */
  @Override
  public Class<IllegalArgumentException> getExceptionClass() {
    return IllegalArgumentException.class;
  }

  /**
   * 当业务代码抛出的IllegalArgumentException被捕获后，会传入该方法进行处理。
   * @param swaggerInvocation 本次业务调用相关的信息
   * @param e 被捕获的异常
   * @return 转换后的响应消息，将会发送给调用方
   */
  @Override
  public Response convert(SwaggerInvocation swaggerInvocation, IllegalArgumentException e) {
    return Response.consumerFailResp(
        new InvocationException(Status.BAD_REQUEST,
            new CommonExceptionData(
                swaggerInvocation.getInvocationQualifiedName() + " gets illegal param: " + e.getMessage()))
    );
  }
}
