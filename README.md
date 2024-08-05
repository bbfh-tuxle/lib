# Tuxle Lib

Tuxle library for reading & writing custom data formats. Used by both Tuxle server & Tuxle client.

<!-- vim-markdown-toc GFM -->

* [📃 Protocol](#-protocol)
    * [Letter Format](#letter-format)
    * [Read a Letter](#read-a-letter)
    * [Write Letter](#write-letter)
* [🗳️ Fields](#-fields)
    * [Parameters](#parameters)
    * [Reading parameters](#reading-parameters)
    * [Writing parameters](#writing-parameters)
* [📁 Channels](#-channels)
    * [File interface](#file-interface)
    * [List entry](#list-entry)
    * [Reading an entry](#reading-an-entry)
    * [Writing entry](#writing-entry)
    * [List file](#list-file)
    * [Read list file entry](#read-list-file-entry)
    * [Write to the list file](#write-to-the-list-file)
    * [Database](#database)
    * [Read a chunk](#read-a-chunk)
    * [Append chunk to database](#append-chunk-to-database)
    * [Overwrite a chunk in database](#overwrite-a-chunk-in-database)

<!-- vim-markdown-toc -->

# 📃 Protocol

## Letter Format

A letter is used to communicate between programs. It has the following format:

```
[TYPE] [endpoint]
<uint32 number of parameters>[parameters][body]\r
```

- `[TYPE]` — The type of the message, can be anything. `ERROR` and `OKAY` are reserved. Conventionally should be all-uppercase.
- `[endpoint]` — The endpoint of the message, can be anything. Used as a sub-title for `[type]`, to specify what message is about.
- `<uint32 number of parameters>` — Describes the number of parameters to parse. Set to `0` if the parameters list is empty.
- `[parameters]` — A list of parameters, every parameter has the following format: `[Key]=[value]\n`. Conventionally all keys should begin with an uppercase letter, values are always strings, but can be parsed to represent anything.
- `[body]` — The body of the letter, can contain any character except for `\r` (use POSIX newline representation `\n`).
- `\r` — Represents the end of the letter, every letter must end with it.

In code it's represented as the following struct:

```go
type Letter struct {
    Type       string
    Endpoint   string
    Parameters fields.Parameters  // Alias for map[string]string
    Body       string
}
```

## Read a Letter

For routing letters to appropriate handlers you can read the letter type using `protocol.ReadLetterType`. This allows you to cancel request without wasting resources parsing the entire letter.

```go
func ReadLetterType(reader *bufio.Reader) (string, error)
```

To parse the entire letter you can use `protocol.ReadLetter`.

```go
// Reads (parses) the entire Letter.
//
// Returns EOF if Parameters format is invalid.
func ReadLetter(reader *bufio.Reader) (*Letter, error)
```

It returns `nil` if an error had occured, otherwise a pointer to `Letter`.

## Write Letter

To send letter you need to encode it first. For this, use the `Letter.Write` method:

```go
// Write the letter to an io.Writter in the correct format.
func (letter *Letter) Write(buffer io.Writer) error
```

Example usage:

```go
letter := protocol.Letter{
    Type:       "TEST",
    Endpoint:   "message",
    Parameters: fields.Parameters{},
    Body:       "",
}

var buffer bytes.Buffer
letter.Write(&buffer)
```

# 🗳️ Fields

Fields contains helpful types that are usually used as fields of other structs.

## Parameters

Parameters is an alias to `map[string]string` with extra convinient functions as a part of Tuxle format.

## Reading parameters

When reading a parameters file, you can use `fields.ReadAllParameters`:

```go
// Read (parse) Parameters. until io.EOF.
//
// err == io.EOF if parameters have invalid format
func ReadAllParameters(reader *bufio.Reader) (Parameters, error)
```

For reading parameters embedded in a letter (with a predefined number of parameters) use `fields.ReadParameters`:

```go
// Read (parse) Parameters.
//
// `count` — The amount of lines to expect. Use ReadAllParameters() to read until io.EOF.
//
// err == io.EOF if parameters have invalid format or count > number of parameters.
func ReadParameters(reader *bufio.Reader, count int) (Parameters, error)
```

## Writing parameters

To encode the Parameters use `Parameters.Write` method:

```go
// Write the Parameters to an io.Writter in the correct format.
//
// No-op when Parameters are empty.
func (params Parameters) Write(buffer io.Writer)
```

# 📁 Channels

This module handles reading/writing to custom file formats.

## File interface

A file that channels will read/write from/to. An `*os.File` will always satisfy this interface.

```go
type File interface {
    Read([]byte) (int, error)
    ReadAt([]byte, int64) (int, error)
    WriteAt([]byte, int64) (int, error)
}
```

## List entry

An entry is a single "line" in the list:

```go
type Entry struct {
    Timestamp uint64  // When was the message sent?
    ChunkId   byte  // What chunk is it stored in?
    ChunkLine uint64  // What's the index with-in the chunk?
    UserId    string  // Who sent it? (max length: 47 bytes)
}
```

## Reading an entry

To read the entry use `channels.ReadEntry`, a single entry is **64 bytes**.

```go
func ReadEntry(reader io.Reader) (*Entry, error)
```

## Writing entry

To write entry into **64 bytes** use `Entry.Write` method.

```go
func (entry Entry) Write(buffer io.Writer) error
```

## List file

Open a file (e.g. `/tmp/messages.list`) and then create a `channels.ListFile`. It will write the header to the file if the file is empty to ensure it uses valid format.

```go
func NewListFile(file File) (*ListFile, error)
```

## Read list file entry

You can read a `channels.Entry` from the list file using either of the following methods:

```go
// Reads entry at a certain index starting from the OLDEST entry.
//
// Entry is nil if an error occured.
func (list ListFile) ReadOldestEntry(index int64) (*Entry, error)
```

or

```go
// Reads entry at a certain index starting from the NEWEST entry.
//
// Returns an error if index is out of bounds.
//
// Entry is nil if an error occured.
func (list ListFile) ReadNewestEntry(index int64) (*Entry, error)
```

## Write to the list file

> TODO

## Database

In [ListEntry](#list-entry) you can see the `ChunkId` variables, it's used to describe which database the message is stored in (255 possible values).

The reason multiple databases are used is because every database uses a different fixed message size.

> **Example (used in Tuxle Server):**
> - `tiny.db` uses 16 bytes per message (4-16 characters in UTF-8).
> - `small.db` uses 256 bytes per message (32-256 characters in UTF-8).
> - `medium.db` uses 1024 bytes per message (128-1024 characters in UTF-8).
> - `large.db` uses 8192 bytes per message (1024-8192 characters in UTF-8).

This allows for a more efficient storage usage (short messages stored separately from large messages).

To create a database, open a file (e.g. `/tmp/tiny.db`) and then use `channels.NewDatabase`. It will write the header to the file if the file is empty to ensure it uses valid format.

```go
func NewDatabase(file File, chunkSize int64) (*Database, error)
```

## Read a chunk

Once the database is open, you can read chunk at a specific index:

```go
// Reads chunk at a specific index.
//
// Returns io.EOF when reading out of bounds data
func (db *Database) ReadChunk(index int64) (string, error)
```

## Append chunk to database

To add a new chunk to the end of the database file use `*Database.AppendChunk` method:

```go
// Appends a chunk to the end of the database.
func (db *Database) AppendChunk(chunk string) error
```

## Overwrite a chunk in database

To edit any chunk you can use the `*Database.OverwriteChunk` method:

```go
// Writes a chunk to certain position in the database.
func (db *Database) OverwriteChunk(chunk string, index int64) error
```
