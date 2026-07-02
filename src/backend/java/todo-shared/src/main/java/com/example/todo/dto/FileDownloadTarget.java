package com.example.todo.dto;

/**
 * Everything the controller needs to stream a file's content back to the client.
 *
 * @param path        absolute path to the file content on disk
 * @param name        the original file name to present to the client
 * @param contentType the MIME type to send in the response
 */
public record FileDownloadTarget(String path, String name, String contentType) {}
