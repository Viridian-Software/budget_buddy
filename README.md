# Viridian Budget Buddy - Backend
## Requirements
This application uses SQLC to generate queries for database interactions, and Goose to handle SQL migrations.
## Style
Variables and constants should be declared using snake_case wherever possible.
## Layout
Authentication and all other utility functions are placed in the internal package or otherwise. The main http handlers are grouped by route and placed in the main folder in separate go files. 
