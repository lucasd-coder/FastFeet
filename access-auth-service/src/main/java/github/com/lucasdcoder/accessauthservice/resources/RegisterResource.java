package github.com.lucasdcoder.accessauthservice.resources;

import github.com.lucasdcoder.accessauthservice.domain.User;
import github.com.lucasdcoder.accessauthservice.services.UserService;
import io.quarkus.logging.Log;
import jakarta.annotation.security.PermitAll;
import jakarta.inject.Inject;
import jakarta.validation.Valid;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

@Path("/api/register")
public class RegisterResource {

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
