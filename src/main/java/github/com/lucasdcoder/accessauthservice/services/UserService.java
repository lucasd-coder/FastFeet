package github.com.lucasdcoder.accessauthservice.services;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Objects;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;
import javax.ws.rs.core.Response;

import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.keycloak.admin.client.CreatedResponseUtil;
import org.keycloak.admin.client.Keycloak;
import org.keycloak.admin.client.resource.RealmResource;
import org.keycloak.admin.client.resource.UserResource;
import org.keycloak.representations.idm.CredentialRepresentation;
import org.keycloak.representations.idm.RoleRepresentation;
import org.keycloak.representations.idm.UserRepresentation;

import github.com.lucasdcoder.accessauthservice.domain.Roles;
import github.com.lucasdcoder.accessauthservice.domain.User;
import github.com.lucasdcoder.accessauthservice.resources.response.GetRolesResponse;
import github.com.lucasdcoder.accessauthservice.resources.response.GetUserResponse;
import github.com.lucasdcoder.accessauthservice.services.exceptions.ResourceNotFoundException;

@ApplicationScoped
public class UserService {

    @ConfigProperty(name = "application.realm")
    String realm;

    @Inject
    Keycloak keycloak;

    public Response createUser(User user) {
        RealmResource realmResource = keycloak.realm(realm);

        List<RoleRepresentation> rolesRepresentation = realmResource.roles().list();

        List<RoleRepresentation> roleRepresentation = filterRolesRepresentations(rolesRepresentation,
                rolesIdentity(user.getAuthority()));

        CredentialRepresentation credential = new CredentialRepresentation();

        credential.setType(CredentialRepresentation.PASSWORD);
        credential.setTemporary(false);
        credential.setValue(user.getPassword());

        UserRepresentation userRepresen = new UserRepresentation();

        userRepresen.setEnabled(true);
        userRepresen.setUsername(user.getUsername());
        userRepresen.setFirstName(user.getFirstName());
        userRepresen.setLastName(user.getLastName());
        userRepresen.setEmail(user.getUsername());
        userRepresen.setCredentials(Arrays.asList(credential));

        Response response = realmResource.users().create(userRepresen);

        String userId = CreatedResponseUtil.getCreatedId(response);

        realmResource.users().get(userId).roles()
                .realmLevel().add(roleRepresentation);

        return response;
    }

    public GetUserResponse findUserByEmail(String email) {
        RealmResource realmResource = keycloak.realm(realm);

        List<UserRepresentation> userRepresentation = realmResource.users()
                .searchByEmail(email, true);

        UserRepresentation users = userRepresentation.stream()
                .filter(user -> user.getEmail().contains(email)).findFirst()
                .orElseThrow(() -> new ResourceNotFoundException("User Not Found"));

        return toGetUserResponse(users);

    }

    public UserResource findById(String id) {
        RealmResource realmResource = keycloak.realm(realm);

        return realmResource.users().get(id);
    }

    public Response getRoles(String id) {
        UserResource userResource = findById(id);

        List<RoleRepresentation> roleRepresentations = userResource.roles().realmLevel().listEffective();

        List<String> roles = roleRepresentations.stream()
                .map(RoleRepresentation::getName)
                .toList();

        GetRolesResponse resp = GetRolesResponse.builder().roles(roles).build();

        return Response.ok(resp).build();
    }

    private List<String> rolesIdentity(Roles roles) {

        List<String> aux = new ArrayList<>(3);

        switch (roles) {
            case ADMIN:
                aux.addAll(Arrays.asList(Roles.ADMIN.getAuthority(), Roles.USER.getAuthority()));
                break;
            case USER:
                aux.add(Roles.USER.getAuthority());
                break;
            default:
                break;
        }

        return aux;
    }

    private List<RoleRepresentation> filterRolesRepresentations(List<RoleRepresentation> roleRepresen,
            List<String> rolesIdentity) {

        return roleRepresen.stream().filter(roles -> rolesIdentity.contains(roles.getName())).toList();
    }

    private GetUserResponse toGetUserResponse(UserRepresentation userRepresentation) {
        if (Objects.isNull(userRepresentation)) {
            return null;
        }

        return GetUserResponse.builder()
                .id(userRepresentation.getId())
                .email(userRepresentation.getEmail())
                .username(userRepresentation.getUsername())
                .enabled(userRepresentation.isEnabled())
                .build();
    }
}
