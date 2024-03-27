package indi.yuluo.gptapiserver.entity;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

/**
 * @author: yuluo
 * @date: 2023/7/21 18:09
 * @description: some desc
 */

@Configuration
public class ServiceConfig {

	@Value("${spring.gpt.service.proxy-address}")
	private String proxyAddress;

	@Value("${spring.gpt.service.api-key}")
	private String apiKey;

	public String getApiKey() {
		return apiKey;
	}

	public void setApiKey(String apiKey) {
		this.apiKey = apiKey;
	}

	public String getProxyAddress() {
		return proxyAddress;
	}

	public void setProxyAddress(String proxyAddress) {
		this.proxyAddress = proxyAddress;
	}
}
