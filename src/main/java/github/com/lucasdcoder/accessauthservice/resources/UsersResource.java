package github.com.lucasdcoder.accessauthservice.resources;

import static github.com.lucasdcoder.accessauthservice.utils.Constants.UUID_REGEX;

import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.validation.constraints.Email;
import javax.validation.constraints.Pattern;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.core.Response;

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