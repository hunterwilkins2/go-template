# Go HTML Project Template

Go template project for a HTML go template, HTMX, and tailwind css server.

## Requirements

- Go 1.20+
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)
- [tailwindcss standalone cli](https://tailwindcss.com/blog/standalone-cli)

## Getting Started

1. Fork this repo
2. Change the go package name
3. Update the git remote origin
4. Run `make run/live` to start server

## Features

- Autocompile and browser hot reloading on save
- zap structured logging
- Middlewares
  - logging
  - csrf protection
  - secure headers on all requests
  - panic recovery
- Out of the box nested Go Templates
- HTMX
- Tailwind css
