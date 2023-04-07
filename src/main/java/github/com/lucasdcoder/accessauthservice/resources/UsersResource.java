package github.com.lucasdcoder.accessauthservice.resources;

import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.validation.constraints.Email;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;

import org.jboss.resteasy.reactive.NoCache;

import github.com.lucasdcoder.accessauthservice.resources.response.GetUserResponse;
import github.com.lucasdcoder.accessauthservice.services.UserService;
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
        return userService.findUserByEmail(email);
    }

}