package indi.yuluo.apiserver.controller;

import indi.yuluo.apiserver.pojo.User;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author yuluo
 * @author 1481556636@qq.com
 */

@RestController
@RequestMapping("/api")
public class ApiServerController {

	private static final User[] users = new User[10];

	static {
		for (int i = 0; i < 10; i++) {
			User user = new User();
			user.setId(i);
			user.setAge(20 + i);
			user.setName("yuluo" + i);

			users[i] = user;
		}
	}

	/**
	 * Test a1 api.
	 * @return string
	 */
	@GetMapping("/test")
	public String testA1() {

		return "test";
	}

	@GetMapping("/user/list")
	public User[] getAllUser() {

		return users;
	}

	@GetMapping("/user/{id}")
	public User getUserById(@PathVariable("id") int id) {

		for (User user : users) {
			if (user.id == id) {
				return user;
			}
		}

		return null;
	}

}
