# Membership Analytics

This is the search and analytics repository for the membership system.

# Getting Started
1. Create a new `config.dev.json` and `config.prod.json` from `config.example.json`
2. Make sure you correctly populate the fields on the config files. Dev uses `members-test` as index as oppose to `members` index
3. Run `go run application.go` on your terminal
4. For deployment, you must set the env variable of `ENV` to `production` or `prod` to use the correct config

# Test
Run `go test $(go list ./... | grep -v /vendor/)` on the root of your directory to run app `_test.go` files excluding the `/vendor`