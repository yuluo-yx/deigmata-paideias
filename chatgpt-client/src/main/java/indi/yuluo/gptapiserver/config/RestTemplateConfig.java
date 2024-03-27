package indi.yuluo.gptapiserver.config;

import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:15
 * @description: some desc
 */

@Component
public class RestTemplateConfig {

	@Bean
	public RestTemplate restTemplate() {

		return new RestTemplate();
	}

}
