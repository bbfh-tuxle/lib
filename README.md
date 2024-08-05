# Tuxle Lib

Tuxle library for reading & writing custom data formats. Used by both Tuxle server & Tuxle client.

All functions/structs are documented in the code. This document is a brief overview of parts that might need additional explanation.

<!-- vim-markdown-toc GFM -->

* [📃 Protocol](#-protocol)
    * [Letter Format](#letter-format)
    * [Read a Letter](#read-a-letter)
    * [Write Letter](#write-letter)
* [🗳️ Fields](#-fields)
    * [Parameters](#parameters)
* [📁 Channels](#-channels)
    * [File interface](#file-interface)
    * [List entry](#list-entry)
    * [Database](#database)

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
