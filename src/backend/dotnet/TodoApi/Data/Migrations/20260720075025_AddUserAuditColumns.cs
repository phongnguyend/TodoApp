using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace TodoApi.Data.Migrations
{
    /// <inheritdoc />
    public partial class AddUserAuditColumns : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<int>(
                name: "CreatedByUserId",
                table: "Users",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "UpdatedByUserId",
                table: "Users",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "CreatedByUserId",
                table: "TodoItems",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "UpdatedByUserId",
                table: "TodoItems",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "CreatedByUserId",
                table: "TodoItemAttachments",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "UpdatedByUserId",
                table: "TodoItemAttachments",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "CreatedByUserId",
                table: "Files",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "UpdatedByUserId",
                table: "Files",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "CreatedByUserId",
                table: "EmailLogs",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "UpdatedByUserId",
                table: "EmailLogs",
                type: "INTEGER",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_Users_CreatedByUserId",
                table: "Users",
                column: "CreatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_Users_UpdatedByUserId",
                table: "Users",
                column: "UpdatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_TodoItems_CreatedByUserId",
                table: "TodoItems",
                column: "CreatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_TodoItems_UpdatedByUserId",
                table: "TodoItems",
                column: "UpdatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_TodoItemAttachments_CreatedByUserId",
                table: "TodoItemAttachments",
                column: "CreatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_TodoItemAttachments_UpdatedByUserId",
                table: "TodoItemAttachments",
                column: "UpdatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_Files_CreatedByUserId",
                table: "Files",
                column: "CreatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_Files_UpdatedByUserId",
                table: "Files",
                column: "UpdatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_EmailLogs_CreatedByUserId",
                table: "EmailLogs",
                column: "CreatedByUserId");

            migrationBuilder.CreateIndex(
                name: "IX_EmailLogs_UpdatedByUserId",
                table: "EmailLogs",
                column: "UpdatedByUserId");

            migrationBuilder.AddForeignKey(
                name: "FK_EmailLogs_Users_CreatedByUserId",
                table: "EmailLogs",
                column: "CreatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_EmailLogs_Users_UpdatedByUserId",
                table: "EmailLogs",
                column: "UpdatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_Files_Users_CreatedByUserId",
                table: "Files",
                column: "CreatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_Files_Users_UpdatedByUserId",
                table: "Files",
                column: "UpdatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_TodoItemAttachments_Users_CreatedByUserId",
                table: "TodoItemAttachments",
                column: "CreatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_TodoItemAttachments_Users_UpdatedByUserId",
                table: "TodoItemAttachments",
                column: "UpdatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_TodoItems_Users_CreatedByUserId",
                table: "TodoItems",
                column: "CreatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_TodoItems_Users_UpdatedByUserId",
                table: "TodoItems",
                column: "UpdatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_Users_Users_CreatedByUserId",
                table: "Users",
                column: "CreatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);

            migrationBuilder.AddForeignKey(
                name: "FK_Users_Users_UpdatedByUserId",
                table: "Users",
                column: "UpdatedByUserId",
                principalTable: "Users",
                principalColumn: "Id",
                onDelete: ReferentialAction.SetNull);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_EmailLogs_Users_CreatedByUserId",
                table: "EmailLogs");

            migrationBuilder.DropForeignKey(
                name: "FK_EmailLogs_Users_UpdatedByUserId",
                table: "EmailLogs");

            migrationBuilder.DropForeignKey(
                name: "FK_Files_Users_CreatedByUserId",
                table: "Files");

            migrationBuilder.DropForeignKey(
                name: "FK_Files_Users_UpdatedByUserId",
                table: "Files");

            migrationBuilder.DropForeignKey(
                name: "FK_TodoItemAttachments_Users_CreatedByUserId",
                table: "TodoItemAttachments");

            migrationBuilder.DropForeignKey(
                name: "FK_TodoItemAttachments_Users_UpdatedByUserId",
                table: "TodoItemAttachments");

            migrationBuilder.DropForeignKey(
                name: "FK_TodoItems_Users_CreatedByUserId",
                table: "TodoItems");

            migrationBuilder.DropForeignKey(
                name: "FK_TodoItems_Users_UpdatedByUserId",
                table: "TodoItems");

            migrationBuilder.DropForeignKey(
                name: "FK_Users_Users_CreatedByUserId",
                table: "Users");

            migrationBuilder.DropForeignKey(
                name: "FK_Users_Users_UpdatedByUserId",
                table: "Users");

            migrationBuilder.DropIndex(
                name: "IX_Users_CreatedByUserId",
                table: "Users");

            migrationBuilder.DropIndex(
                name: "IX_Users_UpdatedByUserId",
                table: "Users");

            migrationBuilder.DropIndex(
                name: "IX_TodoItems_CreatedByUserId",
                table: "TodoItems");

            migrationBuilder.DropIndex(
                name: "IX_TodoItems_UpdatedByUserId",
                table: "TodoItems");

            migrationBuilder.DropIndex(
                name: "IX_TodoItemAttachments_CreatedByUserId",
                table: "TodoItemAttachments");

            migrationBuilder.DropIndex(
                name: "IX_TodoItemAttachments_UpdatedByUserId",
                table: "TodoItemAttachments");

            migrationBuilder.DropIndex(
                name: "IX_Files_CreatedByUserId",
                table: "Files");

            migrationBuilder.DropIndex(
                name: "IX_Files_UpdatedByUserId",
                table: "Files");

            migrationBuilder.DropIndex(
                name: "IX_EmailLogs_CreatedByUserId",
                table: "EmailLogs");

            migrationBuilder.DropIndex(
                name: "IX_EmailLogs_UpdatedByUserId",
                table: "EmailLogs");

            migrationBuilder.DropColumn(
                name: "CreatedByUserId",
                table: "Users");

            migrationBuilder.DropColumn(
                name: "UpdatedByUserId",
                table: "Users");

            migrationBuilder.DropColumn(
                name: "CreatedByUserId",
                table: "TodoItems");

            migrationBuilder.DropColumn(
                name: "UpdatedByUserId",
                table: "TodoItems");

            migrationBuilder.DropColumn(
                name: "CreatedByUserId",
                table: "TodoItemAttachments");

            migrationBuilder.DropColumn(
                name: "UpdatedByUserId",
                table: "TodoItemAttachments");

            migrationBuilder.DropColumn(
                name: "CreatedByUserId",
                table: "Files");

            migrationBuilder.DropColumn(
                name: "UpdatedByUserId",
                table: "Files");

            migrationBuilder.DropColumn(
                name: "CreatedByUserId",
                table: "EmailLogs");

            migrationBuilder.DropColumn(
                name: "UpdatedByUserId",
                table: "EmailLogs");
        }
    }
}
