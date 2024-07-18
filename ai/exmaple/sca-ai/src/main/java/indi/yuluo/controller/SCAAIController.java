package indi.yuluo.controller;

import indi.yuluo.servicd.SCAAIService;
import jakarta.annotation.Resource;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/sca-ai")
public class SCAAIController {

	@Resource
	private SCAAIService scaaiService;

	@GetMapping("/chat")
	public String chat(@RequestParam(value = "message", defaultValue = "Hello")  String msg) {

		return scaaiService.chat(msg).toString();
	}

}
