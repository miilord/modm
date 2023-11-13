# MODM

[![Go Report Card](https://goreportcard.com/badge/github.com/miilord/modm)](https://goreportcard.com/report/github.com/miilord/modm)
[![Go](https://github.com/miilord/modm/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/miilord/modm/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/miilord/modm.svg)](https://pkg.go.dev/github.com/miilord/modm)
[![GitHub](https://img.shields.io/github/license/miilord/modm)](https://github.com/miilord/modm/blob/main/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/miilord/modm/badge.svg?branch=main)](https://coveralls.io/github/miilord/modm?branch=main)

MODM is a MongoDB wrapper built on top of the mongo-go-driver, leveraging the powerful features of Go generics. It provides a convenient interface for CRUD operations, allowing direct return of structured data without the need for code generation.

## Features

- **Structured Data for CRUD Operations:** Directly return structured data from CRUD operations, providing a seamless development experience.

- **No Code Generation Required:** Leverage Go 1.18's generics features to minimize code volume and enhance development efficiency.

- **Flexible Hooks:** Support automated field updates, providing a more adaptable approach to data handling.

- **Easy-to-Use Index Management:** Effortlessly create MongoDB indexes through code for streamlined management.

- **Simple Transactions:** Simplify transaction usage, empowering developers to effortlessly employ MongoDB transactions.

## Requirements

- **Go 1.18 and Above:** MODM is designed to take full advantage of the features introduced in Go 1.18 and later versions.

## Installation

```bash
go get github.com/miilord/modm
```

## Getting Started

### Connecting to the database

To use MODM, you only need to pass \*mongo.Collection in mongo-go-driver, so MODM is compatible with all libraries based on the official driver, e.g. qmgo.

Below is an example of the official driver:

```go
ctx := context.Background()
client, err := mongo.Connect(ctx, options.Client().ApplyURI("your mongodb uri"))
if err != nil {
	panic(err)
}
defer client.Disconnect(ctx)
database := client.Database("test")
```

### Importing MODM

```go
type User struct {
	modm.DefaultField `bson:",inline"`
	Name              string `bson:"name,omitempty" json:"name"`
	Age               int    `bson:"age,omitempty" json:"age"`
}

type DB struct {
	Users *modm.Repo[*User]
}

func main() {
	...
	db := DB{
		Users: modm.NewRepo[*User](database.Collection("users")),
	}
	db.Users.InsertOne(ctx, &User{Name: "gooooo", Age: 6})

	// To query for documents containing zero values, use bson.M, bson.D, or a map.
	// Alternatively, consider using pointers for the appropriate fields.
	user, _ := db.Users.FindOne(ctx, &User{Name: "gooooo"})
	fmt.Println(user.Age) // 6

	// Find() returns ([]*User, error)
	users, _ := db.Users.Find(ctx, &User{Age: 6})
}
```

## Contributions

Contributions are welcome! Feel free to open issues, submit pull requests, or provide suggestions to improve MODM.

## License

MODM is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
