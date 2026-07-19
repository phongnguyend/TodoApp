package com.example.todo.api.controller;

import com.example.todo.dto.*;
import com.example.todo.exception.UserConflictException;
import com.example.todo.exception.UnauthorizedException;
import com.example.todo.security.UserTokenCodec;
import com.example.todo.service.UserService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.time.Instant;
import java.util.List;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(UserController.class)
class UserControllerTest {
    @Autowired MockMvc mockMvc;
    @Autowired ObjectMapper objectMapper;
    @MockitoBean UserService service;
    @MockitoBean UserTokenCodec tokenCodec;

    private UserResponse response() {
        return new UserResponse(1L, "alice", "alice@example.com", true, Instant.now(), null);
    }

    @Test
    void managementEndpointsReturnExpectedStatuses() throws Exception {
        when(service.getAll(1, 20)).thenReturn(new PaginatedResponse<>(List.of(response()), 1, 1, 20, 1));
        when(service.getById(1L)).thenReturn(response());
        when(service.create(any())).thenReturn(response());
        when(service.update(eq(1L), any())).thenReturn(response());
        when(service.setActive(eq(1L), anyBoolean())).thenReturn(response());

        mockMvc.perform(get("/api/users")).andExpect(status().isOk()).andExpect(jsonPath("$.items[0].email").value("alice@example.com"));
        mockMvc.perform(get("/api/users/1")).andExpect(status().isOk());
        mockMvc.perform(post("/api/users").contentType(MediaType.APPLICATION_JSON)
                .content("{\"username\":\"alice\",\"email\":\"alice@example.com\",\"password\":\"password123\"}"))
                .andExpect(status().isCreated()).andExpect(jsonPath("$.password").doesNotExist());
        mockMvc.perform(put("/api/users/1").contentType(MediaType.APPLICATION_JSON).content("{\"username\":\"new-name\"}"))
                .andExpect(status().isOk());
        mockMvc.perform(patch("/api/users/1/activate")).andExpect(status().isOk());
        mockMvc.perform(patch("/api/users/1/deactivate")).andExpect(status().isOk());
    }

    @Test
    void createValidationAndConflictAreMapped() throws Exception {
        mockMvc.perform(post("/api/users").contentType(MediaType.APPLICATION_JSON)
                .content("{\"username\":\"\",\"email\":\"bad\",\"password\":\"short\"}"))
                .andExpect(status().isBadRequest());

        when(service.create(any())).thenThrow(new UserConflictException("Email is already in use."));
        mockMvc.perform(post("/api/users").contentType(MediaType.APPLICATION_JSON)
                .content("{\"username\":\"alice\",\"email\":\"alice@example.com\",\"password\":\"password123\"}"))
                .andExpect(status().isConflict());
    }

    @Test
    void signupAndPasswordResetFlowsReturnContractStatuses() throws Exception {
        when(service.signup(any())).thenReturn(response());

        mockMvc.perform(post("/api/users/signup").contentType(MediaType.APPLICATION_JSON)
                .content("{\"username\":\"alice\",\"email\":\"alice@example.com\",\"password\":\"password123\"}"))
                .andExpect(status().isCreated());
        mockMvc.perform(post("/api/users/password/reset").contentType(MediaType.APPLICATION_JSON)
                .content("{\"email\":\"missing@example.com\"}"))
                .andExpect(status().isAccepted()).andExpect(jsonPath("$.message").value(org.hamcrest.Matchers.containsString("account exists")));
        mockMvc.perform(post("/api/users/password/confirm").contentType(MediaType.APPLICATION_JSON)
                .content("{\"token\":\"token\",\"newPassword\":\"new-password123\"}"))
                .andExpect(status().isNoContent());
    }

    @Test
    void profileAndChangePasswordUseAuthenticatedSubject() throws Exception {
        when(tokenCodec.authenticatedUserId("Bearer token")).thenReturn(7L);
        when(service.getProfile(7L)).thenReturn(response());
        when(service.updateProfile(eq(7L), any())).thenReturn(response());

        mockMvc.perform(get("/api/users/profile").header("Authorization", "Bearer token")).andExpect(status().isOk());
        mockMvc.perform(put("/api/users/profile").header("Authorization", "Bearer token")
                .contentType(MediaType.APPLICATION_JSON).content("{\"username\":\"new-name\"}"))
                .andExpect(status().isOk());
        mockMvc.perform(post("/api/users/password/change").header("Authorization", "Bearer token")
                .contentType(MediaType.APPLICATION_JSON)
                .content("{\"currentPassword\":\"password123\",\"newPassword\":\"new-password123\"}"))
                .andExpect(status().isNoContent());

        verify(service).changePassword(eq(7L), any());
    }

    @Test
    void profileWithoutValidBearerTokenReturnsUnauthorized() throws Exception {
        when(tokenCodec.authenticatedUserId(null))
                .thenThrow(new UnauthorizedException("A valid bearer token is required."));

        mockMvc.perform(get("/api/users/profile"))
                .andExpect(status().isUnauthorized());
        verifyNoInteractions(service);
    }
}
