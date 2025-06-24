# Comprehensive Naming Convention Documentation for Projects

## General Naming Rules
1. **Use CamelCase for Structs and Interfaces**:
   - Structs represent data models or entities.
     - Example: `FriendshipService`, `Expense`
   - Interfaces define contracts or behaviors.
     - Example: `IFriendshipService`, `IExpenseRepository`

2. **Use CamelCase for Variables and Parameters**:
   - Variables store data or state.
     - Example: `UserId`, `FriendUserId`
   - Parameters are inputs to functions.
     - Example: `Ctx`, `RetryCount`

3. **Use PascalCase for Functions**:
   - Functions perform actions or calculations.
     - Example: `Create`, `Filter`, `GetAllMembers`, `Unfriend`

4. **Use UpperCamelCase for Constants**:
   - Constants represent fixed values.
     - Example: `FriendshipStatusAccepted`, `ErrNotFound`

5. **Use snake_case only for database column names and collection fields**:
   - Example: `user_id`, `friend_user_id`

6. **Use descriptive names**:
   - Avoid single-letter names except for loop variables.
     - Example: `FriendshipRepository`, `UserRepository`

## Specific Rules for the Project

### Structs and Interfaces
- Prefix interfaces with `I`.
  - Example: `IFriendshipService`
- Use meaningful names that describe the purpose.
  - Example: `FriendshipService`
- Use CamelCase for struct fields.
  - Example: `UserID`, `FriendUserID`
- Group related fields logically.
  - Example:
    ```go
    type Friendship struct {
        UserID        string
        FriendUserID  string
        CreatedAt     time.Time
    }
    ```

### Functions
- Use verbs for function names.
  - Example: `Create`, `Filter`, `Unfriend`
- Use `ctx` as the first parameter for context.
  - Example: `func (s *FriendshipService) Create(ctx context.Context, userId, friendUserId string)`
- Use descriptive names for methods.
  - Example: `GetFriendshipStatus`, `UpdateFriendship`
- Follow Go idioms for receiver names.
  - Example: `func (s *Service) Create()`

### Variables
- Use CamelCase for local variables and parameters.
  - Example: `UserId`, `FriendUserId`
- Use plural names for collections.
  - Example: `FriendIds`, `Users`
- Use meaningful names for flags and counters.
  - Example: `IsActive`, `RetryCount`
- Avoid abbreviations unless widely understood.
  - Example: `UserId` instead of `Uid`

### Constants
- Use UpperCamelCase for constants.
  - Example: `FriendshipStatusAccepted`
- Prefix error codes with `Err`.
  - Example: `ErrNotFound`, `ErrSystem`
- Use descriptive names for configuration constants.
  - Example: `MaxRetryCount`, `DefaultTimeout`
- Group constants logically.
  - Example:
    ```go
    const (
        FriendshipStatusPending   = "pending"
        FriendshipStatusAccepted = "accepted"
        FriendshipStatusRejected = "rejected"
    )
    ```

### Packages
- Use lowercase for package names.
  - Example: `service`, `adapters`, `constants`
- Avoid underscores in package names.
  - Example: `eventbus` instead of `event_bus`
- Use meaningful names that reflect functionality.
  - Example: `repositories`, `helpers`

### File Names
- Use lowercase and separate words with underscores.
  - Example: `friend.go`, `user.go`
- Use meaningful names for test files.
  - Example: `friend_test.go`, `user_test.go`
- Group files logically within folders.
  - Example:
    - `handlers/expense.go`
    - `models/expense.go`

### Logging
- Use descriptive messages for logging.
  - Example: `logging.Info().Msgf("friendship created %s", friendship)`
- Include context in log messages.
  - Example: `logging.Error().Msgf("failed to create friendship for user %s", userID)`
- Use structured logging where possible.
  - Example:
    ```go
    logging.WithFields(logging.Fields{
        "user_id": userID,
        "friend_id": friendID,
    }).Info("friendship created")
    ```

### Error Handling
- Use `Err` prefix for error variables.
  - Example: `ErrInvalidInput`, `ErrDatabaseConnection`
- Use meaningful names for error messages.
  - Example: `"invalid user ID"`, `"database connection failed"`
- Group error definitions logically.
  - Example:
    ```go
    var (
        ErrInvalidInput = errors.New("invalid input")
        ErrNotFound     = errors.New("not found")
    )
    ```

### Testing
- Use CamelCase for test function names.
  - Example: `TestCreateFriendship`, `TestGetUser`
- Use meaningful names for mock objects.
  - Example: `MockUserRepository`, `MockFriendshipService`
- Group tests logically.
  - Example:
    ```go
    func TestFriendshipService_Create(t *testing.T) {
        // Test logic here
    }
    ```

### Comments
- Use complete sentences for comments.
  - Example: `// Create initializes a new friendship between two users.`
- Use `TODO` for incomplete implementations.
  - Example: `// TODO: Add validation for user input.`
- Document exported functions and types.
  - Example:
    ```go
    // FriendshipService provides methods for managing friendships.
    type FriendshipService struct {
        // Fields here
    }
    ```

### Import Rules
1. **Group Imports Logically**:
   - Separate standard library imports, third-party imports, and project-specific imports.
     - Example:
       ```go
       import (
           "fmt"
           "time"

           "github.com/some/package"
           "project/internal/repositories"
       )
       ```

2. **Use Aliases for Conflicting or Long Import Names**:
   - Use short and meaningful aliases for imports.
     - Example:
       ```go
       import (
           repo "project/internal/repositories"
           utils "project/pkg/helpers/utils"
       )
       ```

3. **Use `_` for unnamed imports**:
   - If an import is not named, use `_` to indicate it is intentionally unused.
     - Example:
       ```go
       import (
           _ "github.com/some/package"
       )
       ```

4. **Avoid Unused Imports**:
   - Remove imports that are not used in the file.

5. **Use Full Paths for Project-Specific Imports**:
   - Always use the full path for project-specific imports.
     - Example:
       ```go
       import "project/pkg/helpers/adapters"
       ```

## Additional Guidelines
- **Consistency**:
  - Ensure consistent naming and organization across all layers of the application.
  - Use linters and formatters to enforce style guidelines.

This document serves as a comprehensive guide for naming conventions and code organization in projects. Adhering to these rules ensures readability, maintainability, and consistency across the codebase.
