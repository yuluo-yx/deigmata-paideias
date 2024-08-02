# Spring Boot Web 应用请求响应太慢排查

线上接口响应太慢，Spring Boot + 内嵌的 Tomcat 服务器，Tomcat 处理请求的最大线程数普通情况是 150 左右，最大是 200，所以当同时处理的请求过多，并且每个请求一直没有处理完成。所有的线程都在繁忙，没有办法处理新的请求，就会导致新的请求排队等待处理，从而造成了迟迟无法响应的线上事故，用户体验太差。

## 场景复现

一个简单的 controller 接口

```java
@RestController
public class SpringBootController {

	@GetMapping("/")
	public String index() throws InterruptedException {

		// 每个线程都睡眠，模拟业务接口响应慢的问题
		Thread.sleep(60 * 30 * 1000);
		System.out.println("current thread: " + Thread.currentThread().getName());

		return "Hello Spring Boot!";
	}

}
```

随后设置 tomcat 服务器的最大线程为 5 ：

```yml
server:
  port: 8080
  tomcat:
    threads:
      max: 5

spring:
  application:
    name: simple-application
    
```

## 测试

之后我们在单元测试中模拟 10 个线程对接口发起请求

```java
public class RequestTest {

	public static void main(String[] args) {
		
		for (int i = 0; i < 10; i++) {
			new Thread(new RequestTask()).start();
		}
	}

	private static class RequestTask implements Runnable {
		@Override
		public void run() {
			try {
				URL url = new URL("http://localhost:8080/");
				HttpURLConnection connection = (HttpURLConnection) url.openConnection();
				connection.setRequestMethod("GET");
				int responseCode = connection.getResponseCode();
				System.out.println("Response Code: " + responseCode);
			}
			catch (IOException e) {
				e.printStackTrace();
			}
		}
	}
	
}
```

首先，我们注释掉 controller 中的睡眠代码，运行测试得到以下的输出：

```java
spring boot web 应用响应输出：
2024-08-01T15:56:25.369+08:00  INFO 7304 --- [simple-application] [nio-8080-exec-2] o.s.web.servlet.DispatcherServlet        : Completed initialization in 0 ms
current thread: http-nio-8080-exec-4
current thread: http-nio-8080-exec-5
current thread: http-nio-8080-exec-1
current thread: http-nio-8080-exec-2
current thread: http-nio-8080-exec-3
current thread: http-nio-8080-exec-5
current thread: http-nio-8080-exec-2
current thread: http-nio-8080-exec-4
current thread: http-nio-8080-exec-1
current thread: http-nio-8080-exec-3
    
测试输出：
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
Response Code: 200
```

我们看到线程在响应输出。接下来，打开注释代码，模拟事故现场。将会在控制台看不到任何输出，测试方法也会卡住不动。

## 排查

我们使用 JVM 的一些命令来进行排查：

```shell
# 首先查看 java 的进程 id
C:\Users\Administrator>jps -l
18400 indi.yuluo.tomcat.SpringBootApplication
10324 jdk.jcmd/sun.tools.jps.Jps
16872 org.jetbrains.jps.cmdline.Launcher
18072 org.jetbrains.idea.maven.server.RemoteMavenServer36
9624
7372 inid.yuluo.tomcat.RequestTest
```

我们看到 SpringBootApplication 和 RequestTest 在运行，之后使用 jstack 命令生成线程快照，并保存为文件。

```shell
C:\Users\Administrator>jstack 18400 > C:\Users\Administrator\Desktop\spring-boot-thread_dump.txt
```

打开线程快照文件，搜索 `http-nio` 就能看到 Tomcat 的请求处理线程，所有的请求处理线程状态都是 `TIMED_WAITING` ，表示线程正在等待另一个线程执行特定的动作，但是有一个指定的等待时间。而且能直接看到请求是阻塞在了哪个代码位置。

```shell
"http-nio-8080-exec-2" #28 daemon prio=5 os_prio=0 cpu=0.00ms elapsed=187.14s tid=0x0000015ffe6089f0 nid=0x1ee4 waiting on condition  [0x00000026eb5fc000]
						# `TIMED_WAITING`
   java.lang.Thread.State: TIMED_WAITING (sleeping)
	at java.lang.Thread.sleep(java.base@17.0.1/Native Method)
	# 阻塞位置：
	at indi.yuluo.tomcat.controller.SpringBootController.index(SpringBootController.java:18)
```

## 解决方案

适当修改 Tomcat 的最大线程数，可以增加并发请求的处理能力。

适当调大 Tomcat 的最小空闲线程数，可以确保在并发高峰时刻，Tomcat 能迅速响应新的请求，而不需要重新创建线程。

修改值需要对用户体量做出预估之后，进行测试之后确定。