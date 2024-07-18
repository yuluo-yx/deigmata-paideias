package indi.yuluo.langchain4j.service.prompt;

import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.model.input.Prompt;
import dev.langchain4j.model.input.structured.StructuredPrompt;
import dev.langchain4j.model.input.structured.StructuredPromptProcessor;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class PromptService {

	// 结构体注解，用于格式化回复
	@StructuredPrompt({
			"""
			请查询演员{{name}} 的信息，
			回答格式如下：
				{
					姓名:...,
					简介:...,
				}     
			"""
	})

	static class CreatActorPrompt{
		String name;

		public CreatActorPrompt(String name) {
			this.name = name;
		}

	}

	@Resource
	private ChatLanguageModel chatLanguageModel;

	public String chat() {

		CreatActorPrompt prompt_ = new CreatActorPrompt("孙悟空");
		Prompt prompt = StructuredPromptProcessor.toPrompt(prompt_);

		return chatLanguageModel.generate(prompt.text());
	}

}


