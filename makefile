run-application:
	./gradlew quarkusDev

test:
	./gradlew test

build:
	./gradlew build -Dquarkus.package.type=native