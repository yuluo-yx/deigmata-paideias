package indi.yuluo.observable.controller;

import indi.yuluo.observable.component.CustomObservation;
import jakarta.annotation.Resource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@RestController
public class TraceController {

	private static final Logger logger = LoggerFactory.getLogger(TraceController.class);

	@Resource
	private CustomObservation customObservation;

	@RequestMapping("/")
	public String home() {

		String s = customObservation.productTrace();

		logger.info("home() has been called");
		return "Hello Tracing! " + s;
	}

}
