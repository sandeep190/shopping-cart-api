# Shopping Cart Api

# Getting started
1. go get https://github.com/sandeep190/shopping-cart-api 
1. Change the .env.example as you need(see warning below)
1. Rename .env.example to .env
1. Seed the database passing "create seed" as arguments to the app(read main.go to understand what I mean)

## WARNING
The recommended database to use is Postgresql, the other database backends may not work as expected.
Unfortunately the MySQL does not work as expected, for example the BeforeSave Hook for User is not able to retrieve
the Role model if using MySQL, the same code does work if SQLite, it is weird, because the SQL query generated is valid and it
returns a row, but somehow the driver is not able to map it to the user.

# Features
- Authentication / Authorization
- JWT middleware for authentication
- Multi file upload
- Database seed
- Paging with Limit and Offset using GORM (Golang ORM framework)
- CRUD operations on products, comments, tags, categories, orders
- Orders, guest users may place an order


# What you will learn
- Golang
- Golang Go-Gonic web framework
- JWT
- Controllers
- Middlewares
- JWT Authentication
- Role based authorization
- GORM
    - associations: ManyToMany, OneToMany, ManyToOne
    - virtual fields
    - Select specific columns
    - Eager loading
    - Count related association
    