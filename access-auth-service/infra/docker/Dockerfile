FROM quay.io/quarkus/ubi-quarkus-graalvmce-builder-image:22.3.3-java17 AS build
USER root
RUN microdnf install findutils make gcc
COPY --chown=quarkus:quarkus gradlew /code/gradlew
COPY --chown=quarkus:quarkus gradle /code/gradle
COPY --chown=quarkus:quarkus build.gradle /code/
COPY --chown=quarkus:quarkus settings.gradle /code/
COPY --chown=quarkus:quarkus gradle.properties /code/
RUN mkdir /musl && \
    curl -L -o musl.tar.gz https://more.musl.cc/11.2.1/x86_64-linux-musl/x86_64-linux-musl-native.tgz && \
    tar -xvzf musl.tar.gz -C /musl --strip-components 1 && \
    curl -L -o zlib.tar.gz https://www.zlib.net/current/zlib.tar.gz && \
    mkdir zlib && tar -xvzf zlib.tar.gz -C zlib --strip-components 1 && \
    cd zlib && ./configure --static --prefix=/musl && \
    make && make install && \
    cd .. && rm -rf zlib && rm -f zlib.tar.gz && rm -f musl.tar.gz
ENV PATH="/musl/bin:${PATH}"
USER quarkus
WORKDIR /code
COPY src /code/src
RUN ./gradlew build -Dquarkus.package.type=native -Dquarkus.profile=prod -Dquarkus.native.debug.enabled=true -Dquarkus.otel.traces.exporter=cdi,otlp

FROM quay.io/quarkus/quarkus-micro-image:2.0
WORKDIR /work/
COPY --from=build /code/build/*-runner /work/application

RUN chmod 775 /work
EXPOSE 8080

ENTRYPOINT [ "sh", "-c" ,"./application $JAVA_OPTS"]