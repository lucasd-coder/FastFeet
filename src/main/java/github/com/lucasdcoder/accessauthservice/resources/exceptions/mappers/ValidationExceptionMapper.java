package github.com.lucasdcoder.accessauthservice.resources.exceptions.mappers;

import java.time.Instant;

import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;
import javax.ws.rs.ext.ExceptionMapper;
import javax.ws.rs.ext.Provider;

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
