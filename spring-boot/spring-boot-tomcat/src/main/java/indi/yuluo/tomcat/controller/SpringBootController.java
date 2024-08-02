package indi.yuluo.tomcat.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

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
