package github.com.lucasdcoder.accessauthservice.resources;

import javax.annotation.security.PermitAll;
import javax.inject.Inject;
import javax.validation.Valid;
import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

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
