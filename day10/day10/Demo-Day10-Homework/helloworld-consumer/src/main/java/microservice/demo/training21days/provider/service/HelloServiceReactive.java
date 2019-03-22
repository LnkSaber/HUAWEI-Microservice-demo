package microservice.demo.training21days.provider.service;

import java.util.concurrent.CompletableFuture;

public interface HelloServiceReactive {
  CompletableFuture<String> sayHello(String name);
}
