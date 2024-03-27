package indi.yuluo.gptapiserver.service;

import indi.yuluo.gptapiserver.entity.ServiceEntity;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:12
 * @description: some desc
 */
public interface RequestService {

	String request(ServiceEntity serviceEntity);

}
