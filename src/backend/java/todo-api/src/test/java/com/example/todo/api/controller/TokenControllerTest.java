package com.example.todo.api.controller;

import com.example.todo.dto.TokenResponse;
import com.example.todo.exception.UnauthorizedException;
import com.example.todo.service.UserService;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.context.annotation.Import;
import com.example.todo.api.config.SecurityConfig;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@WebMvcTest(TokenController.class)
@Import(SecurityConfig.class)
class TokenControllerTest {
    @Autowired MockMvc mockMvc;
    @MockitoBean UserService service;

    @Test
    void createsTokenWithNoStoreHeaders() throws Exception {
        when(service.createToken(any())).thenReturn(new TokenResponse("header.payload.signature", "Bearer", 3600));

        mockMvc.perform(post("/api/tokens").contentType(MediaType.APPLICATION_JSON)
                        .content("{\"email\":\"alice@example.com\",\"password\":\"password123\"}"))
                .andExpect(status().isOk())
                .andExpect(header().string("Cache-Control", "no-store"))
                .andExpect(header().string("Pragma", "no-cache"))
                .andExpect(jsonPath("$.access_token").value("header.payload.signature"))
                .andExpect(jsonPath("$.token_type").value("Bearer"))
                .andExpect(jsonPath("$.expires_in").value(3600));
    }

    @Test
    void invalidCredentialsReturnNonDisclosingUnauthorizedResponse() throws Exception {
        when(service.createToken(any())).thenThrow(new UnauthorizedException("Invalid email or password."));

        mockMvc.perform(post("/api/tokens").contentType(MediaType.APPLICATION_JSON)
                        .content("{\"email\":\"missing@example.com\",\"password\":\"wrong\"}"))
                .andExpect(status().isUnauthorized())
                .andExpect(header().string("WWW-Authenticate", "Bearer"))
                .andExpect(jsonPath("$.error").value("Invalid email or password."));
    }

    @Test
    void invalidRequestReturnsBadRequest() throws Exception {
        mockMvc.perform(post("/api/tokens").contentType(MediaType.APPLICATION_JSON)
                        .content("{\"email\":\"invalid\"}"))
                .andExpect(status().isBadRequest());
    }
}
