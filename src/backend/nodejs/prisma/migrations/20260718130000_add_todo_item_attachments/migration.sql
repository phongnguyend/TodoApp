CREATE TABLE "todo_item_attachments" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "todoItemId" INTEGER NOT NULL,
    "fileId" INTEGER NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" DATETIME,
    CONSTRAINT "todo_item_attachments_todoItemId_fkey" FOREIGN KEY ("todoItemId") REFERENCES "todo_items" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "todo_item_attachments_fileId_fkey" FOREIGN KEY ("fileId") REFERENCES "files" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX "todo_item_attachments_fileId_idx" ON "todo_item_attachments"("fileId");
CREATE UNIQUE INDEX "todo_item_attachments_todoItemId_fileId_key" ON "todo_item_attachments"("todoItemId", "fileId");
