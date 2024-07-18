package indi.yuluo.langchain4j.service.aiservice;

import java.util.List;

import dev.langchain4j.service.SystemMessage;
import dev.langchain4j.service.UserMessage;
import dev.langchain4j.service.V;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

public interface SystemAndUserMessage {

	@SystemMessage("You are a professional translator into {{language}}")
	@UserMessage("Translate the following text: {{text}}")
	String translate(@V("text") String text, @V("language") String language);

	@SystemMessage("Summarize every message from user in {{n}} bullet points. Provide only bullet points.")
	List<String> summarize(@UserMessage String text, @V("n") int n);

}
