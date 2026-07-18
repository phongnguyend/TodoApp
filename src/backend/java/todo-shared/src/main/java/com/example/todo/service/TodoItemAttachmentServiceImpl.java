package com.example.todo.service;

import com.example.todo.dto.SaveTodoItemAttachmentRequest;
import com.example.todo.dto.TodoItemAttachmentResponse;
import com.example.todo.entity.TodoItemAttachment;
import com.example.todo.repository.FileRepository;
import com.example.todo.repository.TodoItemAttachmentRepository;
import com.example.todo.repository.TodoItemRepository;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TodoItemAttachmentServiceImpl implements TodoItemAttachmentService {

    private final TodoItemAttachmentRepository repository;
    private final TodoItemRepository todoItemRepository;
    private final FileRepository fileRepository;

    private void requireTodoItem(Long todoItemId) {
        if (!todoItemRepository.existsById(todoItemId)) {
            throw new EntityNotFoundException("Todo item " + todoItemId + " not found.");
        }
    }

    private void requireFile(Long fileId) {
        if (!fileRepository.existsById(fileId)) {
            throw new EntityNotFoundException("File " + fileId + " not found.");
        }
    }

    private TodoItemAttachment requireAttachment(Long todoItemId, Long attachmentId) {
        return repository.findByIdAndTodoItemId(attachmentId, todoItemId)
                .orElseThrow(() -> new EntityNotFoundException(
                        "Attachment " + attachmentId + " not found for todo item " + todoItemId + "."));
    }

    @Override
    public List<TodoItemAttachmentResponse> getAll(Long todoItemId) {
        requireTodoItem(todoItemId);
        return repository.findByTodoItemIdOrderByCreatedAtAsc(todoItemId).stream()
                .map(TodoItemAttachmentResponse::from)
                .toList();
    }

    @Override
    public TodoItemAttachmentResponse getById(Long todoItemId, Long attachmentId) {
        requireTodoItem(todoItemId);
        return TodoItemAttachmentResponse.from(requireAttachment(todoItemId, attachmentId));
    }

    @Override
    @Transactional
    public TodoItemAttachmentResponse create(Long todoItemId, SaveTodoItemAttachmentRequest request) {
        requireTodoItem(todoItemId);
        requireFile(request.fileId());
        return repository.findByTodoItemIdAndFileId(todoItemId, request.fileId())
                .map(TodoItemAttachmentResponse::from)
                .orElseGet(() -> TodoItemAttachmentResponse.from(
                        repository.save(new TodoItemAttachment(todoItemId, request.fileId()))));
    }

    @Override
    @Transactional
    public TodoItemAttachmentResponse update(
            Long todoItemId, Long attachmentId, SaveTodoItemAttachmentRequest request) {
        requireTodoItem(todoItemId);
        requireFile(request.fileId());
        TodoItemAttachment attachment = requireAttachment(todoItemId, attachmentId);

        return repository.findByTodoItemIdAndFileId(todoItemId, request.fileId())
                .filter(existing -> !existing.getId().equals(attachmentId))
                .map(TodoItemAttachmentResponse::from)
                .orElseGet(() -> {
                    attachment.setFileId(request.fileId());
                    return TodoItemAttachmentResponse.from(repository.save(attachment));
                });
    }

    @Override
    @Transactional
    public void delete(Long todoItemId, Long attachmentId) {
        requireTodoItem(todoItemId);
        repository.delete(requireAttachment(todoItemId, attachmentId));
    }
}
