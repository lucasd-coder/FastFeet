run-application:
	QUARKUS_KEYCLOAK_DEVSERVICES_ENABLED=false ./gradlew quarkusDev

test:
	./gradlew test