package github.com.lucasdcoder.accessauthservice.domain;

import static github.com.lucasdcoder.accessauthservice.utils.Constants.REGEX_DEFAULT;

import javax.validation.constraints.Email;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Pattern;
import javax.validation.constraints.Size;

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
