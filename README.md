# Buggeon Backend

## Description

This is a backend repository of Buggeon Project. Wrote on Golang, working with MongoDB and SeewedFS.

## Repo structure

```
[marcus@neo backend]$ tree
.
├── cmd
│   └── main.go
├── config
│   └── config.go
├── docs
│   ├── database.scheme.excalidraw
│   └── entities
│       ├── board.md
│       ├── card.md
│       ├── member.md
│       ├── project.md
│       ├── team.md
│       ├── user.md
│       └── zone.md
├── go.mod
├── go.sum
├── internal
│   ├── db
│   │   └── database.go
│   ├── dto
│   │   ├── board_dto.go
│   │   ├── card_dto.go
│   │   ├── member_dto.go
│   │   ├── message_dto.go
│   │   ├── project_dto.go
│   │   └── user_dto.go
│   ├── handlers
│   │   ├── project_handler.go
│   │   ├── system_handler.go
│   │   ├── test_handler.go
│   │   └── user_handler.go
│   ├── middleware
│   │   └── auth.go
│   ├── models
│   │   ├── board_model.go
│   │   ├── card_model.go
│   │   ├── member_model.go
│   │   ├── message_model.go
│   │   ├── project_model.go
│   │   ├── token.go
│   │   └── user_model.go
│   ├── repositories
│   │   ├── board_repo.go
│   │   ├── card_repo.go
│   │   ├── member_repo.go
│   │   ├── project_repo.go
│   │   └── user_repo.go
│   ├── security
│   │   └── hash_password.go
│   └── services
│       ├── board_service.go
│       ├── card_service.go
│       ├── member_service.go
│       ├── notification_service.go
│       ├── project_service.go
│       ├── system_service.go
│       ├── token_service.go
│       └── user_service.go
├── LICENSE
└── README.md
```

## About contributing

Buggeon project distributes on AGPL