package indi.yuluo.controller;

import java.util.Map;

import indi.yuluo.service.SpringAIService;
import jakarta.annotation.Resource;

import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
@RequestMapping("/spring-ai")
public class SpringAIController {

	@Resource
	private SpringAIService springAIService;


	@GetMapping("/chat")
	public Map<String, ChatResponse> chat(@RequestParam(value = "msg", defaultValue = "Tell me a joke") String msg) {

		return springAIService.chat(msg);
	}

}
