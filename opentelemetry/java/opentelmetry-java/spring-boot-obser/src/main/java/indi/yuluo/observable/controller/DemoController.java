package indi.yuluo.observable.controller;

import java.util.Random;

import indi.yuluo.observable.component.HealthComponent;
import jakarta.annotation.Resource;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
public class DemoController {

	@Resource
	private HealthComponent component;

	@GetMapping("/sb")
	public String sb() {

		// mock service call.
		component.count();

		return String.valueOf(new Random().nextInt(1000));
	}

}
