# A Clean Code & DDD Adventure in Go 🚀

Embark on a journey through my first Go application, where I have embraced Clean Code principles and Domain-Driven Design (DDD) to build a robust and maintainable system.

## Project Overview

This project is a hands-on exploration of building a login system in Go. The primary goal was to create a functional application while adhering to Clean Code and DDD best practices. I aimed for simplicity, avoiding unnecessary external dependencies to focus on core design principles.

## Project Structure & Architecture

The project structure follows DDD principles to ensure separation of concerns and maintainability:

```
server-skl/
├── cmd/            # Main application entry point
├── domain/         # Core business logic and entities
├── entrypoint/     # Application execution layer
├── infra/          # Infrastructure-related code
├── permissions/    # Business rules abstraction
├── usecase/        # Application-specific business logic
├── .gitignore      # Git ignore file
├── go.mod          # Go module file
├── go.sum          # Go module checksum file
└── README.md       # Project documentation
```

### Layer Responsibilities

- **Domain:** Contains core business logic and entities, independent of external dependencies.
- **Usecase:** Implements application-specific business rules by orchestrating domain objects.
- **Entrypoint:** Initializes dependencies and runs the application.
- **Infra:** Manages infrastructure components like databases and external APIs.
- **Permissions:** Handles business rule abstraction.
- **Cmd:** Contains the executable application.


### Key Takeaways

- **DDD in Practice:** Designing with DDD improved my understanding of business modeling.
- **Clean Architecture:** The separation of concerns enhanced testability and maintainability.
- **Go's Simplicity:** The language's minimalism and concurrency support made development smooth and enjoyable.


## Future Roadmap

- **Database Integration:** Implement persistent storage with PostgreSQL or MySQL.
- **Unit Testing:** Develop a comprehensive suite of tests.
- **Authentication:** Introduce JWT or other authentication mechanisms.
- **Code Refactoring:** Improve code structure and readability.
- **Enhanced Error Handling:** Implement more robust error logging and handling.

## Self-Reflection & Areas for Improvement

- **Testing:** Increasing test coverage is a priority.
- **Error Handling:** Need to implement more informative and structured error responses.
- **Documentation:** Improving internal code documentation for better readability.
