# Build stage - Docker build context: java/ directory
#   docker build -f todo-worker/Dockerfile -t todo-worker .
FROM maven:3.9-eclipse-temurin-21 AS build
WORKDIR /src

# Copy all POM files first to leverage Docker layer caching for dependency downloads
COPY pom.xml .
COPY todo-shared/pom.xml todo-shared/pom.xml
COPY todo-api/pom.xml todo-api/pom.xml
COPY todo-worker/pom.xml todo-worker/pom.xml
RUN mvn dependency:go-offline -B

# Copy source for modules needed by the worker
COPY todo-shared/src todo-shared/src
COPY todo-worker/src todo-worker/src
RUN mvn package -DskipTests -B -pl todo-worker -am

# Runtime stage
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

COPY --from=build /src/todo-worker/target/*.jar app.jar

ENTRYPOINT ["java", "-jar", "app.jar"]
