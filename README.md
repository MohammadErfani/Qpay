# Qpay

## About the project

The template is used to create golang project. All golang projects must follow the conventions in the
template. Calling for exceptions must be brought up in the engineering team.

## Table of Contents

1. [Getting Started](#getting-started)
    - [Layout](#layout)
2. [How to Run](#how-to-run)
    - [Clone Project](#clone-project)
3. [Documentation](#documentation)
    - [API Documentation](#API-documentation)
    - [Control Version](#control-version)
        - [Branches](#branches)
        - [Branch Naming](#branch-naming)
        - [Commit Message](#commit-message)
4. [Notes](#notes)

## Getting started

Below we describe the conventions or tools specific to this project.

### Layout

```tree
.
├── cmd
│   ├── migrate.go
│   ├── root.go
│   └── seed.go
├── config
│   └── config.go
├── database
│   ├── migration
│   │   ├── 000001_init_migration.down.sql
│   │   ├── 000001_init_migration.up.sql
│   │   ├── 000002_complete_transaction.down.sql
│   │   ├── 000002_complete_transaction_migration.up.sql
│   │   ├── 000003_change_phone_number_type_in_transaction.down.sql
│   │   └── 000003_change_phone_number_type_in_transaction.up.sql
│   ├── postgres.go
├── docs
│   ├── docs.go
│   ├── packages
│   │   ├── Cobra.md
│   │   ├── Echo.md
│   │   ├── Echo-Swagger.md
│   │   ├── Golang-Migrate.md
│   │   ├── Gorm.md
│   │   ├── JWT.md
│   │   ├── Monkey.md
│   │   ├── SQL-Mock.md
│   │   ├── Testify.md
│   │   └── Viper.md
│   ├── swagger.json
│   └── swagger.yaml
├── models
│   ├── bank.go
│   ├── bankaccount.go
│   ├── commission.go
│   ├── gateway.go
│   ├── transaction.go
│   └── user.go
├── server
│   └── server.go
│   ├── handlers
│   │   ├── adminHandler.go
│   │   ├── authHandler.go
│   │   ├── bankaccountHandler.go
│   │   ├── gatewayHandler.go
│   │   ├── handler.go
│   │   ├── transactionHandler.go
│   │   └── userHandler.go
│   ├── middlewares
│   │   ├── adminMiddleware.go
│   │   └── authMiddleware.go
│   └── routes
│       ├── admin.go
│       ├── auth.go
│       ├── bankaccount.go
│       ├── gateway.go
│       ├── init_routes_v1.go
│       ├── payment.go
│       ├── transaction.go
│       └── user.go
├── services
│   ├── admin.go
│   ├── bank.go
│   ├── bankaccount.go
│   ├── bankaccount_test.go
│   ├── commission.go
│   ├── gateway.go
│   ├── gateway_test.go
│   ├── login.go
│   ├── transaction.go
│   ├── transaction_test.go
│   ├── user.go
│   └── user_test.go
└── utils
    ├── auth.go
    ├── gateway.go
    ├── mock_sheba.go
    ├── mock_transcation.go
    ├── password.go
    ├── password_test.go
    ├── validation.go
    └── validation_test.go
├── README.md
├── config.yaml
├── sample_config.yaml
```

A brief description of the layout:
#### cmd

- `migrate.go`: Database migration script.
- `root.go`: Root command for the application.
- `seed.go`: Script for seeding initial data.

#### config

- `config.go`: Configuration config settings for the project.

#### database

- `postgres.go`: Database connection setup.
- `migration`: Directory containing database migration scripts.

#### docs

- `docs.go`: Documentation generation script.
- `packages`: Directory containing individual markdown files for packages used in project.
- `swagger.json`: Swagger JSON file.
- `swagger.yaml`: Swagger YAML file.

#### models

- `bank.go`: Model for the bank entity.
- `bankaccount.go`: Model for bank accounts.
- `commission.go`: Model for commissions.
- `gateway.go`: Model for payment gateways.
- `transaction.go`: Model for transactions.
- `user.go`: Model for users.

#### server

- `handlers`: Directory containing HTTP request handlers.
- `middlewares`: Directory containing middleware functions.
- `routes`: Directory containing route definitions.
- `server.go`: Main server setup.

#### services

- Business logic and service layer of the application.

#### utils

- Various utility functions and helpers.

## How to Run
#### Clone Project
```bash
git clone https://github.com/GoFellas/Qpay.git
```
#### Before Start
Inside the project folder
```bash
cp sample_config.yaml config.yaml
```
Set your database config, your server, your desire jwt secret key and your initial admin

```bash
go build
./Qpay migrate
./Qpay seed
./Qpay
```
Use --help for more instruction if needed 
## Documentation

### API Documentation
After running project checkout for api documentation on:
```
your_host:your_ip/doc/index.html
```
You can also use swagger.yaml for api documentation

### Control Version

Project use git for control version and  Git Flow model for manage branches and branches will be merged to get log of all commit in all branch and pull request will be without commit message

The Git Flow model consists of two main branches: master and develop.

in this approach all branch will merge to dev and after test will be merge to master

#### Branches

master: This branch represents the production-ready code. It should only contain code that has been thoroughly tested and is ready to be deployed to production.

develop: This branch is used to develop new features. It should contain the latest development changes and should be the base branch for all feature branches.

In addition to these two main branches, Git Flow defines three types of supporting branches

feature: These branches are used to develop new features. They are based on develop and are merged back into develop once the feature is complete.

release: These branches are used to prepare the code for a new production release. They are based on develop and are merged back into both develop and master once the release is complete.

hotfix: These branches are used to quickly fix issues in the production code. They are based on master and are merged back into both develop and master.

####  Branch naming

You can name a branch in Git using the command git branch <branch-name>, where <branch-name> is the name you want to give to the branch. For example, to create a new branch called "feature/add_login_page", you can run the following command:

git branch feature/add_login_page

This will create a new branch with the name feature/add-login-page based on the current branch you are on.

#### Commit message

ConventionalCommit message is a specific format for writing commit messages that provides a standardized way of conveying information about changes made to code in a repository.

The format consists of three parts:

    A type that describes the kind of change being made, such as feat for a new feature, fix for a bug fix, docs for documentation updates, refactor for code refactoring, and so on.

    A scope that describes the part of the codebase being modified, such as a specific module, component, or function.

    A short description that summarizes the changes made in the commit.

Optionally, the commit message can also include a longer description that provides more detailed information about the changes, as well as references to related issues, pull requests, or other relevant information.

For example, a conventional commit message for a bug fix in the authentication module of an application might look like this:

fix(auth): Validate user input before authentication

This commit fixes a bug where the authentication module could accept invalid user input, leading to security vulnerabilities. The fix adds input validation checks to the authentication process to ensure that only valid user input is processed.

Closes #123

By using a conventional commit message format, developers can more easily understand the nature and purpose of changes made to code in a repository, which can help improve collaboration, code quality, and maintenance of the codebase over time.




## Notes
