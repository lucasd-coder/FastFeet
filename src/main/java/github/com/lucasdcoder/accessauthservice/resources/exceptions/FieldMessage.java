package github.com.lucasdcoder.accessauthservice.resources.exceptions;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class FieldMessage implements Serializable {

	private static final long serialVersionUID = 1L;

	private String fieldName;
	private String message;
}
