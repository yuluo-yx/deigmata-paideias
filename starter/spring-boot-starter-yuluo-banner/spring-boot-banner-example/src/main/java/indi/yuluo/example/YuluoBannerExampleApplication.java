package indi.yuluo.example;

import indi.yuluo.banner.EnableYuluoBanner;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@EnableYuluoBanner
@SpringBootApplication
public class YuluoBannerExampleApplication {

	public static void main(String[] args) {

		SpringApplication.run(YuluoBannerExampleApplication.class, args);
	}

	@Bean
	@GetMapping("/")
	public String hello() {

		return "Hello, Yuluo!";
	}

}
