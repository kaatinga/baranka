# baranka

A simple and flexible helper for preparing SQL queries with positional parameters for SQL databases like PostgreSQL, MySQL, or SQLite.

## Features

- Generates SQL value blocks with correct placeholders (`$1`, `$2`, ... or `?`)
- Supports custom value templates (e.g., `(%s)`)
- Collects arguments for use with database drivers
- Supports SQL expressions with placeholders and arguments
- Easily extensible with options

## Installation

```sh
go get github.com/kaatinga/baranka
```

## Usage

### Basic Example

```go
import "github.com/kaatinga/baranka"

b := baranka.New()
b.Add(1, "foo")
b.Add(2, "bar")

query := "INSERT INTO my_table (id, name) VALUES " + b.Values()
args := b.Args()
// Use query and args with your database/sql driver
```

### Using Options

You can customize the placeholder format and value template:

```go
b := baranka.New(
    baranka.WithIncludeTemplate("(%s)"),
    baranka.WithPlaceholderFormat(baranka.PlaceholderFormatQuestionMark),
)

b.Add(1, "foo")
b.Add(2, "bar")

query := "INSERT INTO my_table (id, name) VALUES " + b.Values()
// VALUES (?,?), (?,?)
args := b.Args()
// [1, "foo", 2, "bar"]
```

### Using Expressions

You can embed SQL expressions with their own placeholders and arguments using the `Expression` type.

#### How to use `Expression`

- The `Expression` type allows you to inject SQL fragments with their own arguments.
- The `template` field must be a valid `fmt.Sprintf`-style template string, using `%s` for each argument you want to substitute with a placeholder.
- The number of `%s` in the template **must match** the number of elements in the `args` slice.
- Placeholders (`$1`, `$2`, `?`, etc.) will be substituted for each `%s` in the template, and the corresponding values will be appended to the argument list.

#### Example

```go
b := baranka.New(baranka.WithPlaceholderFormat(baranka.PlaceholderFormatDollar))

b.Add(
    baranka.Expression{
        template: "POINT(%s %s)", // two %s for two arguments
        args:     []any{10.1, 20.2},
    },
)
b.Add(
    baranka.Expression{
        template: "POINT(%s %s)",
        args:     []any{11.1, 21.2},
    },
)

query := "INSERT INTO points (geom) VALUES " + b.Values()
// Resulting VALUES: (POINT($1 $2)), (POINT($3 $4))
args := b.Args()
// [10.1, 20.2, 11.1, 21.2]
```

**Note:**  
If you provide a template with more or fewer `%s` than arguments, the resulting SQL will be invalid.

## API

### `New(opts ...option) *Baranka`

Creates a new Baranka helper with optional configuration.

### `(*Baranka) Add(args ...any)`

Adds a new block of values and collects arguments. Supports `Expression` for SQL fragments.

### `(*Baranka) Args() []any`

Returns the collected arguments in order.

### `(*Baranka) Values() string`

Returns the SQL value blocks as a string, e.g. `($1,$2),\n($3,$4)` or `(?,?)`.

### Options

- `WithIncludeTemplate(template string)`  
  Sets the template for value blocks (default: `(%s)`).

- `WithPlaceholderFormat(format PlaceholderFormat)`  
  Sets the placeholder format:  
  - `PlaceholderFormatDollar` (default, for PostgreSQL: `$1`, `$2`, ...)  
  - `PlaceholderFormatQuestionMark` (for MySQL/SQLite: `?`)

- `WithBlocksLength(length int)`  
  Pre-allocates the argument slice for performance.

## Example: Bulk Insert

```go
b := baranka.New(
    baranka.WithIncludeTemplate("(%s)"),
    baranka.WithPlaceholderFormat(baranka.PlaceholderFormatDollar),
)

for _, row := range rows {
    b.Add(row.ID, row.Name)
}

query := "INSERT INTO users (id, name) VALUES " + b.Values()
args := b.Args()
// db.Exec(query, args...)
```

## Testing

Run tests with:

```sh
go test ./...
```

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.
