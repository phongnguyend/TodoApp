import { ValidationPipe } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { AppModule } from './app.module';

/**
 * Bootstrap function - mirrors Program.cs in ASP.NET Core.
 * Sets up the NestJS application, global pipes, and Swagger.
 */
async function bootstrap(): Promise<void> {
  const app = await NestFactory.create(AppModule);

  // Global validation pipe (like FluentValidation / Data Annotations pipeline in ASP.NET Core)
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,          // strip unknown properties
      forbidNonWhitelisted: true,
      transform: true,          // auto-transform payloads to DTO class instances
    }),
  );

  // CORS
  app.enableCors();

  // Swagger / OpenAPI (mirrors app.UseSwagger() + app.UseSwaggerUI())
  const config = new DocumentBuilder()
    .setTitle('Todo API')
    .setDescription('RESTful API for managing todo items')
    .setVersion('1.0')
    .build();
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('swagger', app, document);

  const port = process.env.PORT ?? 3000;
  await app.listen(port);
  console.log(`Application is running on: http://localhost:${port}`);
  console.log(`Swagger UI: http://localhost:${port}/swagger`);
}

bootstrap();
