package github.com.lucasdcoder.accessauthservice.resources.exceptions.mappers;

import java.time.Instant;

import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;
import javax.ws.rs.ext.ExceptionMapper;
import javax.ws.rs.ext.Provider;

import github.com.lucasdcoder.accessauthservice.resources.exceptions.StandardError;
import github.com.lucasdcoder.accessauthservice.services.exceptions.ResourceNotFoundException;

@Provider
public class ResourceNotFoundExceptionMapper implements
                ExceptionMapper<ResourceNotFoundException> {

        @Override
        public Response toResponse(ResourceNotFoundException ex) {

                Status status = Response.Status.NOT_FOUND;
                StandardError err = StandardError.builder()
                                .timestamp(Instant.now())
                                .status(status.getStatusCode())
                                .error("Resource not found")
                                .message(ex.getMessage())
                                .build();

                return Response.status(status)
                                .entity(err)
                                .build();
        }
}
