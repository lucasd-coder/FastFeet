package github.com.lucasdcoder.accessauthservice.resources.exceptions.mappers;

import java.time.Instant;

import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.core.Response.Status;
import jakarta.ws.rs.ext.ExceptionMapper;
import jakarta.ws.rs.ext.Provider;

import github.com.lucasdcoder.accessauthservice.resources.exceptions.ValidationError;
import github.com.lucasdcoder.accessauthservice.resources.exceptions.ValidationException;

@Provider
public class ValidationExceptionMapper implements
                ExceptionMapper<ValidationException> {

        @Override
        public Response toResponse(ValidationException ex) {

                Status status = Response.Status.BAD_REQUEST;
                ValidationError err = ValidationError.builder()
                                .timestamp(Instant.now())
                                .status(status.getStatusCode())
                                .error("Validation exception")
                                .message(ex.getMessage())
                                .build();

                err.addError(ex.getField(), ex.getMessage());

                return Response.status(status)
                                .entity(err)
                                .build();
        }
}
