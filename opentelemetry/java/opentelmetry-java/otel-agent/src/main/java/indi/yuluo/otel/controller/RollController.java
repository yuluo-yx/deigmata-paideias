package indi.yuluo.otel.controller;

import java.util.Random;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
public class RollController {

	@GetMapping("/index")
	public String index() {

		return String.valueOf(new Random().nextInt(100));
	}
}

