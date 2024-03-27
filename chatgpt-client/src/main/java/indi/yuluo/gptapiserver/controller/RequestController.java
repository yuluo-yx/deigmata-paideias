package indi.yuluo.gptapiserver.controller;

import indi.yuluo.gptapiserver.entity.ServiceEntity;
import indi.yuluo.gptapiserver.service.RequestService;
import jakarta.annotation.Resource;
import lombok.extern.slf4j.Slf4j;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.Mapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:03
 * @description: some desc
 */

@Slf4j
@RestController
public class RequestController {

	@Resource
	private RequestService requestService;

	@PostMapping("/question")
	public String request(@RequestBody ServiceEntity serviceEntity) {

		return requestService.request(serviceEntity);
	}

}
