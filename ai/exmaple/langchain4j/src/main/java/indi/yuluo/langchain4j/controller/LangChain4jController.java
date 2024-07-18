package indi.yuluo.langchain4j.controller;

import java.util.Map;

import indi.yuluo.langchain4j.service.aiservice.AIServiceImpl;
import indi.yuluo.langchain4j.service.aiservice.AiServiceImpl2;
import indi.yuluo.langchain4j.service.common.LangChain4jService;
import indi.yuluo.langchain4j.service.prompt.PromptService;
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
@RequestMapping("/langchain4j")
public class LangChain4jController {

	@Resource
	private LangChain4jService langChain4jService;

	@GetMapping("/chat")
	public Map<String, String> chat(@RequestParam(value = "message", defaultValue = "Hello") String message) {

		return langChain4jService.chat(message);
	}

	@Resource
	private AIServiceImpl aiService1;

	@GetMapping("/aiservice1")
	public void aiService1(@RequestParam(value = "message", defaultValue = "Hello") String message) {

		aiService1.chat(message);
	}

	@Resource
	private AiServiceImpl2 aiService2;

	@GetMapping("/aiservice2")
	public void aiService2() {

		aiService2.chat();
	}

	@Resource
	private PromptService promptService;

	@GetMapping("/prompt")
	public String prompt() {

		return promptService.chat();
	}

}
