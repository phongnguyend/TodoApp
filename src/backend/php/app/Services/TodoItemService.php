<?php

namespace App\Services;

use App\Http\Requests\CreateTodoItemRequest;
use App\Http\Requests\UpdateTodoItemRequest;
use App\Models\TodoItem;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Services\Contracts\TodoItemServiceInterface;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Http\UploadedFile;
use Illuminate\Pagination\LengthAwarePaginator;
use PhpOffice\PhpSpreadsheet\IOFactory;
use PhpOffice\PhpSpreadsheet\Spreadsheet;
use PhpOffice\PhpSpreadsheet\Writer\Xlsx;

/**
 * Business-logic layer.
 * Analogous to a scoped service registered via builder.Services.AddScoped<ITodoItemService, TodoItemService>().
 */
class TodoItemService implements TodoItemServiceInterface
{
    public function __construct(
        private readonly TodoItemRepositoryInterface $repository
    ) {}

    // ── Queries ───────────────────────────────────────────────────────────────

    public function getAll(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->repository->paginate($page, $perPage);
    }

    public function getIncomplete(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->repository->paginateIncomplete($page, $perPage);
    }

    public function getById(int $id): TodoItem
    {
        $todo = $this->repository->findById($id);

        if ($todo === null) {
            throw new ModelNotFoundException("Todo item {$id} not found.");
        }

        return $todo;
    }

    // ── Commands ──────────────────────────────────────────────────────────────

    public function create(CreateTodoItemRequest $request): TodoItem
    {
        return $this->repository->create($request->validated());
    }

    public function update(int $id, UpdateTodoItemRequest $request): TodoItem
    {
        $todo = $this->getById($id);
        return $this->repository->update($todo, $request->validated());
    }

    public function delete(int $id): void
    {
        $todo = $this->getById($id);
        $this->repository->delete($todo);
    }

    public function markComplete(int $id): TodoItem
    {
        $todo = $this->getById($id);
        return $this->repository->update($todo, ['is_completed' => true]);
    }

    // ── CSV import/export ─────────────────────────────────────────────────────

    public function importCsv(UploadedFile $file): array
    {
        // Strip a UTF-8 BOM (common in Excel-exported CSV files) before parsing.
        $contents = (string) preg_replace('/^\xEF\xBB\xBF/', '', (string) file_get_contents($file->getRealPath()));

        $handle = fopen('php://temp', 'r+');
        fwrite($handle, $contents);
        rewind($handle);

        $header = fgetcsv($handle);
        if ($header === false) {
            fclose($handle);
            return ['imported' => 0, 'failed' => 0, 'errors' => []];
        }
        $header = array_map(static fn ($column) => strtolower(trim((string) $column)), $header);
        $columnCount = count($header);

        $imported = 0;
        $errors = [];
        $rowNumber = 1; // the header occupies row 1

        while (($row = fgetcsv($handle)) !== false) {
            $rowNumber++;

            if ($row === [null]) {
                continue; // skip blank lines
            }

            $row = array_pad(array_slice($row, 0, $columnCount), $columnCount, null);
            $data = array_combine($header, $row);

            $title = trim((string) ($data['title'] ?? ''));
            if ($title === '') {
                $errors[] = ['row' => $rowNumber, 'error' => 'Title is required.'];
                continue;
            }

            $description = trim((string) ($data['description'] ?? ''));

            $this->repository->create([
                'title'        => $title,
                'description'  => $description !== '' ? $description : null,
                'is_completed' => $this->parseBool($data['is_completed'] ?? null),
            ]);

            $imported++;
        }

        fclose($handle);

        return [
            'imported' => $imported,
            'failed'   => count($errors),
            'errors'   => $errors,
        ];
    }

    public function exportCsv(): string
    {
        $items = $this->repository->getAllOrdered();

        $handle = fopen('php://temp', 'r+');
        fputcsv($handle, ['id', 'title', 'description', 'is_completed', 'created_at', 'updated_at']);

        foreach ($items as $item) {
            fputcsv($handle, [
                $item->id,
                $item->title,
                $item->description ?? '',
                $item->is_completed ? 'true' : 'false',
                $item->created_at?->toIso8601String() ?? '',
                $item->updated_at?->toIso8601String() ?? '',
            ]);
        }

        rewind($handle);
        $content = (string) stream_get_contents($handle);
        fclose($handle);

        return $content;
    }

    // ── Excel import/export ───────────────────────────────────────────────────

    public function importExcel(UploadedFile $file): array
    {
        $rows = IOFactory::load($file->getRealPath())->getActiveSheet()->toArray(null, true, true, false);

        if (empty($rows)) {
            return ['imported' => 0, 'failed' => 0, 'errors' => []];
        }

        $header = array_map(static fn ($column) => strtolower(trim((string) $column)), array_shift($rows));
        $columnCount = count($header);

        $imported = 0;
        $errors = [];
        $rowNumber = 1; // the header occupies row 1

        foreach ($rows as $row) {
            $rowNumber++;

            if ($this->isBlankRow($row)) {
                continue;
            }

            $row = array_pad(array_slice($row, 0, $columnCount), $columnCount, null);
            $data = array_combine($header, $row);

            $title = trim((string) ($data['title'] ?? ''));
            if ($title === '') {
                $errors[] = ['row' => $rowNumber, 'error' => 'Title is required.'];
                continue;
            }

            $description = trim((string) ($data['description'] ?? ''));

            $this->repository->create([
                'title'        => $title,
                'description'  => $description !== '' ? $description : null,
                'is_completed' => $this->parseBool($data['is_completed'] ?? null),
            ]);

            $imported++;
        }

        return [
            'imported' => $imported,
            'failed'   => count($errors),
            'errors'   => $errors,
        ];
    }

    public function exportExcel(): string
    {
        $items = $this->repository->getAllOrdered();

        $spreadsheet = new Spreadsheet();
        $sheet = $spreadsheet->getActiveSheet();
        $sheet->setTitle('Todo Items');
        $sheet->fromArray(['id', 'title', 'description', 'is_completed', 'created_at', 'updated_at'], null, 'A1');

        $rowIndex = 2;
        foreach ($items as $item) {
            $sheet->fromArray([
                $item->id,
                $item->title,
                $item->description ?? '',
                $item->is_completed,
                $item->created_at?->toIso8601String() ?? '',
                $item->updated_at?->toIso8601String() ?? '',
            ], null, "A{$rowIndex}");
            $rowIndex++;
        }

        $stream = fopen('php://temp', 'r+');
        (new Xlsx($spreadsheet))->save($stream);
        rewind($stream);
        $content = (string) stream_get_contents($stream);
        fclose($stream);

        return $content;
    }

    // ── Helpers ───────────────────────────────────────────────────────────────

    private function isBlankRow(?array $row): bool
    {
        if ($row === null) {
            return true;
        }

        foreach ($row as $value) {
            if ($value !== null && trim((string) $value) !== '') {
                return false;
            }
        }

        return true;
    }

    private function parseBool(mixed $value): bool
    {
        if (is_bool($value)) {
            return $value;
        }

        return in_array(strtolower(trim((string) $value)), ['1', 'true', 'yes', 'y'], true);
    }
}
