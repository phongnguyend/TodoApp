package com.example.todo.repository;

import com.example.todo.entity.User;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
class UserRepositoryTest {
    @Autowired UserRepository repository;

    @Test
    void findsUsersAndChecksUniquenessIgnoringCase() {
        User saved = repository.saveAndFlush(new User("Alice", "Alice@Example.com", "hash", true));

        assertThat(repository.findByEmailIgnoreCase("alice@example.com")).contains(saved);
        assertThat(repository.existsByUsernameIgnoreCase("alice")).isTrue();
        assertThat(repository.existsByEmailIgnoreCaseAndIdNot("ALICE@example.com", saved.getId())).isFalse();
    }
}
