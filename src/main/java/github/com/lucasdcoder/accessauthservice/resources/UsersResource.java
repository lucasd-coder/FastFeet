package github.com.lucasdcoder.accessauthservice.resources;

import static github.com.lucasdcoder.accessauthservice.utils.Constants.UUID_REGEX;

import jakarta.inject.Inject;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.Pattern;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.core.Response;

import jakarta.annotation.security.RolesAllowed;
import org.jboss.resteasy.reactive.NoCache;

import github.com.lucasdcoder.accessauthservice.resources.response.GetUserResponse;
import github.com.lucasdcoder.accessauthservice.services.UserService;
import io.quarkus.logging.Log;
import io.quarkus.security.Authenticated;
import io.quarkus.security.identity.SecurityIdentity;

@Path("/api/users")
@Authenticated
public class UsersResource {

    @Inject
    SecurityIdentity identity;

    @Inject
    UserService userService;

    @GET
    @RolesAllowed("admin")
    @Path("/{email}")
    @NoCache
    public GetUserResponse findUserByEmail(@PathParam("email") @Email String email) {
        Log.info("Received request GET findUserByEmail");
        return userService.findUserByEmail(email);
    }

    @GET
    @RolesAllowed("admin")
    @Path("roles/{id}")
    @NoCache
    public Response getRoles(@PathParam("id") @Pattern(regexp = UUID_REGEX) String id) {
        Log.infof("Received request GET on with id: %s", id);
        return userService.getRoles(id);
    }

    @GET
    @RolesAllowed("admin")
    @Path("is-active/{id}")
    @NoCache
    public Response isActiveUser(@PathParam("id") @Pattern(regexp = UUID_REGEX) String id) {
        Log.infof("Received request GET on with id: %s", id);
        return userService.isActiveUser(id);
    }
}