# MODM

MODM is a MongoDB wrapper built on top of the mongo-go-driver, leveraging the powerful features of Go generics. It provides a convenient interface for CRUD operations, allowing direct return of structured data without the need for code generation.

## Features

- **Structured Data for CRUD Operations:** Directly return structured data from CRUD operations, providing a seamless development experience.

- **No Code Generation Required:** MODM eliminates the need for code generation, making development more straightforward.

- **Flexible Hooks:** MODM supports flexible hooks for insert, find, and update operations, enabling developers to add custom logic easily.

- **Automatic Field Handling:** MODM handles fields automatically, reducing the need for manual intervention in the data structure.

- **Synchronized Index Management:** Ensure your MongoDB indexes stay in sync with your Go code, making it easier to manage your data.

- **User-Friendly Transactions:** MODM simplifies the use of transactions, making it easy for developers to work with MongoDB transactions seamlessly.

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
