package github.com.lucasdcoder.accessauthservice.resources.exceptions.mappers;

import java.time.Instant;
import java.util.Arrays;

import javax.ws.rs.core.Response;
import javax.ws.rs.ext.ExceptionMapper;
import javax.ws.rs.ext.Provider;

import org.jboss.resteasy.reactive.ClientWebApplicationException;

import github.com.lucasdcoder.accessauthservice.resources.exceptions.StandardError;

@Provider
public class HandlerClientWebApplicationExceptionMapper implements
                ExceptionMapper<ClientWebApplicationException> {

        @Override
        public Response toResponse(ClientWebApplicationException exception) {
                StandardError err = StandardError.builder()
                                .timestamp(Instant.now())
                                .status(exception.getResponse().getStatus())
                                .error(Arrays.toString(exception.getStackTrace()))
                                .message(exception.getMessage())
                                .build();

                return Response.status(err.getStatus())
                                .entity(err)
                                .build();
        }

}
