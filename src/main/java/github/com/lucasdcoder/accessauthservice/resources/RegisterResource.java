package github.com.lucasdcoder.accessauthservice.resources;

import jakarta.inject.Inject;
import jakarta.validation.Valid;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

import jakarta.annotation.security.PermitAll;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.keycloak.admin.client.Keycloak;

import github.com.lucasdcoder.accessauthservice.domain.User;
import github.com.lucasdcoder.accessauthservice.services.UserService;
import io.quarkus.logging.Log;

@Path("/api/register")
public class RegisterResource {

    @ConfigProperty(name = "application.realm")
    String realm;

    @Inject
    Keycloak keycloak;

    @Inject
    UserService userService;

    @PermitAll
    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public Response createUser(@Valid User user) {
        Log.info("Received request POST on [] with payload");
        return userService.createUser(user);
    }

}
