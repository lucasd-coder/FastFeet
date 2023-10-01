package github.com.lucasdcoder.accessauthservice.resources.response;

import java.io.Serializable;

import io.quarkus.runtime.annotations.RegisterForReflection;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@RegisterForReflection
@NoArgsConstructor
public class GetUserResponse implements Serializable {

    private static final long serialVersionUID = 1L;

    private String id;

    private String username;

    private Boolean enabled;

    private String email;
}
