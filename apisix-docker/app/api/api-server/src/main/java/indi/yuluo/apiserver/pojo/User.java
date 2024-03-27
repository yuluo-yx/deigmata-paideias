package indi.yuluo.apiserver.pojo;

/**
 * @author yuluo
 * @author 1481556636@qq.com
 */

public class User {

	public int id;

	public int age;

	public String name;

	public int getId() {
		return id;
	}

	public void setId(int id) {
		this.id = id;
	}

	public int getAge() {
		return age;
	}

	public void setAge(int age) {
		this.age = age;
	}

	public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}

	@Override
	public String toString() {
		final StringBuilder sb = new StringBuilder("User{");
		sb.append("id=").append(id);
		sb.append(", age=").append(age);
		sb.append(", name='").append(name).append('\'');
		sb.append('}');
		return sb.toString();
	}

}
