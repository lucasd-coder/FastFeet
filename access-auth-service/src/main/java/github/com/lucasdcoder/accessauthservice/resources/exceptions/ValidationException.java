package github.com.lucasdcoder.accessauthservice.resources.exceptions;

import lombok.Getter;

public class ValidationException extends
        RuntimeException {

    private static final long serialVersionUID = 1L;

    @Getter
    private final String field;

    public ValidationException(String message, String field) {
        super(message);
        this.field = field;
    }
}
