https://opentelemetry.io/docs/languages/java/getting-started/

## 1. 准备 java spring boot 应用

## 2. 下载 otel-java-agent.jar

## 3. 配置 java 代理和 运行的 jar 包

```shell
export JAVA_TOOL_OPTIONS="-javaagent:D:\project\opentelemetry\java\otel-java\opentelemetry-javaagent.jar" \
OTEL_TRACES_EXPORTER=logging \
OTEL_METRICS_EXPORTER=logging \
OTEL_LOGS_EXPORTER=logging \
OTEL_METRIC_EXPORT_INTERVAL=15000
```

## 4. 运行 jar 

java -jar target/spring-boot-otel.jar

## 5. 看到如下的指标输出

```shell
[otel.javaagent 2024-04-18 12:17:44:615 +0800] [PeriodicMetricReader-1] INFO io.opentelemetry.exporter.logging.LoggingMetricExporter - Received a collection of 12 metrics for export.
[otel.javaagent 2024-04-18 12:17:44:615 +0800] [PeriodicMetricReader-1] INFO io.opentelemetry.exporter.logging.Logging
MetricExporter - metric: ImmutableMetricData{resource=Resource{schemaUrl=https://opentelemetry.io/schemas/1.24.0, attributes={host.arch="amd64", host.name="yuluo-laptop", os.description="Windows 11 10.0", os.type="windows", process.comm
and_line="D:\environment\Java\graalvm-jdk22\graalvm-community-openjdk-22+36.1\bin\java.exe -XX:ThreadPriorityPolicy=1 
-XX:+UnlockExperimentalVMOptions -XX:+EnableJVMCIProduct -XX:-UnlockExperimentalVMOptions -javaagent:D:\project\opente
lemetry\java\otel-java\opentelemetry-javaagent.jar -jar target/otel-java-2024.04.13.jar", process.executable.path="D:\
environment\Java\graalvm-jdk22\graalvm-community-openjdk-22+36.1\bin\java.exe", process.pid=3888, process.runtime.desc
ription="GraalVM Community OpenJDK 64-Bit Server VM 22+36-jvmci-b02", process.runtime.name="OpenJDK Runtime Environmen
t", process.runtime.version="22+36-jvmci-b02", service.instance.id="871864b9-55cf-4523-802c-4f5441a0dc83", service.nam
e="otel-application", telemetry.distro.name="opentelemetry-java-instrumentation", telemetry.distro.version="2.3.0", te
lemetry.sdk.language="java", telemetry.sdk.name="opentelemetry", telemetry.sdk.version="1.37.0"}}, instrumentationScop
eInfo=InstrumentationScopeInfo{name=io.opentelemetry.runtime-telemetry-java8, version=2.3.0-alpha, schemaUrl=null, att
ributes={}}, name=jvm.class.unloaded, description=Number of classes unloaded since JVM start., unit={class}, type=LONG
_SUM, data=ImmutableSumData{points=[ImmutableLongPointData{startEpochNanos=1713413849262309100, epochNanos=1713413864271759000, attributes={}, value=1, exemplars=[]}], monotonic=true, aggregationTemporality=CUMULATIVE}}
[otel.javaagent 2024-04-18 12:17:44:616 +0800] [PeriodicMetricReader-1] INFO io.opentelemetry.exporter.logging.Logging
MetricExporter - metric: ImmutableMetricData{resource=Resource{schemaUrl=https://opentelemetry.io/schemas/1.24.0, attr
ibutes={host.arch="amd64", host.name="yuluo-laptop", os.description="Windows 11 10.0", os.type="windows", process.comm
and_line="D:\environment\Java\graalvm-jdk22\graalvm-community-openjdk-22+36.1\bin\java.exe -XX:ThreadPriorityPolicy=1 
-XX:+UnlockExperimentalVMOptions -XX:+EnableJVMCIProduct -XX:-UnlockExperimentalVMOptions -javaagent:D:\project\opente
lemetry\java\otel-java\opentelemetry-javaagent.jar -jar target/otel-java-2024.04.13.jar", process.executable.path="D:\
environment\Java\graalvm-jdk22\graalvm-community-openjdk-22+36.1\bin\java.exe", process.pid=3888, process.runtime.desc
ription="GraalVM Community OpenJDK 64-Bit Server VM 22+36-jvmci-b02", process.runtime.name="OpenJDK Runtime Environmen
t", process.runtime.version="22+36-jvmci-b02", service.instance.id="871864b9-55cf-4523-802c-4f5441a0dc83", service.nam
e="otel-application", telemetry.distro.name="opentelemetry-java-instrumentation", telemetry.distro.version="2.3.0", te
lemetry.sdk.language="java", telemetry.sdk.name="opentelemetry", telemetry.sdk.version="1.37.0"}}, instrumentationScop
eInfo=InstrumentationScopeInfo{name=io.opentelemetry.runtime-telemetry-java8, version=2.3.0-alpha, schemaUrl=null, att
ributes={}}, name=jvm.memory.limit, description=Measure of max obtainable memory., unit=By, type=LONG_SUM, data=Immuta
bleSumData{points=[ImmutableLongPointData{startEpochNanos=1713413849262309100, epochNanos=1713413864271759000, attribu
tes={jvm.memory.pool.name="Compressed Class Space", jvm.memory.type="non_heap"}, value=1073741824, exemplars=[]}, Immu
tableLongPointData{startEpochNanos=1713413849262309100, epochNanos=1713413864271759000, attributes={jvm.memory.pool.na
me="CodeHeap 'non-nmethods'", jvm.memory.type="non_heap"}, value=7012352, exemplars=[]}, ImmutableLongPointData{startE
pochNanos=1713413849262309100, epochNanos=1713413864271759000, attributes={jvm.memory.pool.name="CodeHeap 'non-profile
d nmethods'", jvm.memory.type="non_heap"}, value=122355712, exemplars=[]}, ImmutableLongPointData{startEpochNanos=1713
413849262309100, epochNanos=1713413864271759000, attributes={jvm.memory.pool.name="CodeHeap 'profiled nmethods'", jvm.
memory.type="non_heap"}, value=122290176, exemplars=[]}, ImmutableLongPointData{startEpochNanos=1713413849262309100, e
pochNanos=1713413864271759000, attributes={jvm.memory.pool.name="G1 Old Gen", jvm.memory.type="heap"}, value=4234149888, exemplars=[]}], monotonic=false, aggregationTemporality=CUMULATIVE}}
[otel.javaagent 2024-04-18 12:17:44:617 +0800] [PeriodicMetricReader-1] INFO io.opentelemetry.exporter.logging.Logging
MetricExporter - metric: ImmutableMetricData{resource=Resource{schemaUrl=https://opentelemetry.io/schemas/1.24.0, attr
ibutes={host.arch="amd64", host.name="yuluo-laptop", os.description="Windows 11 10.0", os.type="windows", process.comm
and_line="D:\environment\Java\graalvm-jdk22\graalvm-community-openjdk-22+36.1\bin\java.exe -XX:ThreadPriorityPolicy=1 
-XX:+UnlockExperimentalVMOptions -XX:+EnableJVMCIProduct -XX:-UnlockExperimentalVMOptions -javaagent:D:\project\opente
lemetry\java\otel-java\opentelemetry-javaagent.jar -jar target/otel-java-2024.04.13.jar", process.executable.path="D:\
environment\Java\graalvm-jdk22\graalvm-community-openjdk-22+36.1\bin\java.exe", process.pid=3888, process.runtime.desc
ription="GraalVM Community OpenJDK 64-Bit Server VM 22+36-jvmci-b02", process.runtime.name="OpenJDK Runtime Environmen
t", process.runtime.version="22+36-jvmci-b02", service.i........................
```
