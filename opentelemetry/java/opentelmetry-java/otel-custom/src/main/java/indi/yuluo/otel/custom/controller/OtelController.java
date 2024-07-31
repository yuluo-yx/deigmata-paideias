package indi.yuluo.otel.custom.controller;

import indi.yuluo.otel.custom.service.OtelService;
import jakarta.annotation.Resource;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/")
public class OtelController {

	@Resource
	private OtelService service;

	@GetMapping("/otel")
	public String hello() {

		return service.hello();
	}

}
