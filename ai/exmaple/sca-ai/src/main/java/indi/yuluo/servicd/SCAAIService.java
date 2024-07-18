package indi.yuluo.servicd;

import java.util.Map;

import com.alibaba.cloud.ai.tongyi.chat.TongYiChatModel;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class SCAAIService {

	@Resource
	private TongYiChatModel chatModel;

	public Map<String, String> chat(String msg) {

		return Map.of(msg, chatModel.call(msg));
	}

}
