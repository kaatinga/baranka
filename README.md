# baranka

A simple and flexible helper for preparing SQL queries with positional parameters for SQL databases like PostgreSQL, MySQL or SQLite.

## Features

- Generates SQL value blocks with correct placeholders (`$1`, `$2`, ... or `?`)
- Supports custom value templates (e.g., `(%s)`)
- Collects arguments for use with database drivers
- Easily extensible with options

## Installation

```sh
go get github.com/yourusername/baranka
```

## Usage

### Basic Example

```go
import "github.com/yourusername/baranka"

b := baranka.NewBaranka()
b.Add(1, "foo")
b.Add(2, "bar")

query := "INSERT INTO my_table (id, name) VALUES " + b.Values()
args := b.Args()
// Use query and args with your database/sql driver
```

### Using Options

You can customize the placeholder format and value template:

```go
b := baranka.getOptions([]baranka.option{
    baranka.WithIncludeTemplate("(%s)"),
    baranka.WithPlaceholderFormat(baranka.PlaceholderFormatQuestionMark),
})

b.Add(1, "foo")
b.Add(2, "bar")

query := "INSERT INTO my_table (id, name) VALUES " + b.Values()
// VALUES (?,?), (?,?)
args := b.Args()
// [1, "foo", 2, "bar"]
```

## API

### `NewBaranka() *Baranka`

Creates a new Baranka helper with default settings.

### `(*Baranka) Add(args ...any)`

Adds a new block of values and collects arguments.

### `(*Baranka) Args() []any`

Returns the collected arguments in order.

### `(*Baranka) Values() string`

Returns the SQL value blocks, e.g. `($1,$2),\n($3,$4)` or `(?,?)`.

### Options

- `WithIncludeTemplate(template string)`  
  Sets the template for value blocks (default: `(%s)`).

- `WithPlaceholderFormat(format PlaceholderFormat)`  
  Sets the placeholder format:  
  - `PlaceholderFormatDollar` (default, for PostgreSQL: `$1`, `$2`, ...)  
  - `PlaceholderFormatQuestionMark` (for MySQL/SQLite: `?`)

## Example: Bulk Insert

```go
b := baranka.getOptions([]baranka.option{
    baranka.WithIncludeTemplate("(%s)"),
    baranka.WithPlaceholderFormat(baranka.PlaceholderFormatDollar),
})

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

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.
```
