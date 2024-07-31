spring boot 应用启动配置时添加如下 程序参数

```shell
-javaagent:D:/information/deigmata-paideias/opentelemetry/java/opentelmetry-java/otel-agent/opentelemetry-javaagent.jar
-Dotel.exporter.otlp.endpoint=http://localhost:4318
-Dotel.service.name=springboot-otel
```

或者是 

```shell
mvn clean install 
打 jar 包之后
java -javaagent:D:/information/deigmata-paideias/opentelemetry/java/opentelmetry-java/otel-agent/opentelemetry-javaagent.jar -Dotel.exporter.otlp.endpoint=http://192.168.2.31:4317 -Dotel.service.name=springboot-otel -jar target/otel-agent-1.0-SNAPSHOT.jar
```

```shell
-javaagent:D:/information/deigmata-paideias/opentelemetry/java/opentelmetry-java/otel-agent/opentelemetry-javaagent.jar
-Dotel.traces.exporter=jaeger
-Dotel.exporter.jaeger.endpoint=http://localhost:14250
-Dotel.exporter.jaeger.timeout=10000
-Dotel.metrics.exporter=logging
-Dotel.logs.exporter=none
```