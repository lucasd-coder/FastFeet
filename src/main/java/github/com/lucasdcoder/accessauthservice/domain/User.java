package github.com.lucasdcoder.accessauthservice.domain;

import static github.com.lucasdcoder.accessauthservice.utils.Constants.REGEX_DEFAULT;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class User {

    @Pattern(regexp = REGEX_DEFAULT)
    private String firstName;

    @Pattern(regexp = REGEX_DEFAULT)
    private String lastName;

    @NotBlank
    @Email
    @Pattern(regexp = REGEX_DEFAULT)
    private String username;

    @Size(min = 8)
    @NotNull
    @Pattern(regexp = REGEX_DEFAULT)
    private String password;

    @Pattern(regexp = REGEX_DEFAULT)
    @NotNull
    private String authority;

    public Roles getAuthority() {
        return Roles.toRoles(this.authority);
    }

}
