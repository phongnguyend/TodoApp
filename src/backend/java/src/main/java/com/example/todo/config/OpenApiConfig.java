package com.example.todo.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * OpenAPI / Swagger configuration — mirrors services.AddSwaggerGen() in ASP.NET Core.
 * SpringDoc auto-registers Swagger UI at /swagger-ui.html and JSON at /v3/api-docs.
 */
@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI openAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("Todo API")
                        .description("RESTful API for managing todo items")
                        .version("1.0.0"));
    }
}
