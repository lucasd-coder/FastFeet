package github.com.lucasdcoder.accessauthservice.domain;

import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;

import github.com.lucasdcoder.accessauthservice.resources.exceptions.ValidationException;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Getter
@NoArgsConstructor
public enum Roles {
	ADMIN(Collections.singleton("ADMIN"), "admin"),
	USER(Collections.singleton("USER"), "user");

	private Collection<String> field;

	private String authority;

	Roles(Collection<String> field, String authority) {
		this.field = field;
		this.authority = authority;
	}

	public static Roles toRoles(String role) {
		return Arrays.stream(Roles.values())
				.filter(s -> s.field.contains(role))
				.findFirst()
				.orElseThrow(
						() -> new ValidationException(String.format("Unsupported operation to role %s", role),
								"authority"));
	}
}
