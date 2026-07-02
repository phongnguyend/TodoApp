namespace TodoApi.Services;

/// <summary>
/// Thrown when an uploaded file exceeds the configured maximum upload size.
/// </summary>
public class FileTooLargeException(string message) : Exception(message)
{
}
