

.PHONY: run-application build-application test run-native-agent

run-application:
	./gradlew quarkusDev

build-application:
	./gradlew build -Dquarkus.package.type=native -Dquarkus.profile=dev

test:
	./gradlew test

run-native-agent:
	./gradlew build
	@java -Dquarkus.profile=dev -agentlib:native-image-agent=access-filter-file=src/main/resources/access-filter.json,config-output-dir=src/main/resources/native-image/ \
     -jar build/quarkus-app/quarkus-run.jar