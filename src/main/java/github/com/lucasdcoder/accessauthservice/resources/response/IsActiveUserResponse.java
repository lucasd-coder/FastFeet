package github.com.lucasdcoder.accessauthservice.resources.response;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class IsActiveUserResponse implements Serializable {
    private Boolean active;
}
