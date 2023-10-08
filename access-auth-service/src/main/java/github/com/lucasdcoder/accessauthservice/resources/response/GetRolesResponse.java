package github.com.lucasdcoder.accessauthservice.resources.response;

import java.util.List;

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
public class GetRolesResponse {
    private List<String> roles;
}
