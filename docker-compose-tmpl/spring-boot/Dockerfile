FROM openjdk:17-oracle

LABEL author=yuluo \
    email=yuluo829@aliyun.com \
    version="2023.1.28"

ARG APP_HOME

VOLUME ${APP_HOME}
WORKDIR ${APP_HOME}

ENTRYPOINT ["java", "-jar", "app.jar", "--spring.profiles.active=docker"]

