package microservice.demo.training21days.provider.bootevent;

import org.apache.servicecomb.core.BootListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

@Component
public class CustomBootEventListener implements BootListener {

  private static final Logger LOGGER = LoggerFactory.getLogger(CustomBootEventListener.class);

  public void onBootEvent(BootEvent bootEvent) {
    // BootEvent中的EventType有多种
    switch (bootEvent.getEventType()) {
      case AFTER_REGISTRY: // 微服务实例注册成功
        LOGGER.info("=============================");
        LOGGER.info("Service startup completed!");
        LOGGER.info("=============================");
        break;
      case BEFORE_CLOSE:   // 微服务进程即将退出
        LOGGER.info("=============================");
        LOGGER.info("JVM process is closing!");
        LOGGER.info("=============================");
        break;
      default:
    }
  }
}
