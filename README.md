
# Forum

This project consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

## ER Diagram

![ERD](https://images-ext-1.discordapp.net/external/_YbIrYIAmICRZ9kqceoEctlY-ASq2zldW0SUUs_LKEs/https/i.imgur.com/TfY6UA4.png?format=webp&quality=lossless&width=1041&height=863)

## Clone repository

```bash
    git clone git@git.01.alem.school:aecheist/forum.git
```

## Move to the direcroty

```bash
    cd forum
```

## Run Locally

with makefile

```bash
    make build
    make run
```

without docker

```bash
    go run cmd/main.go
```

with docker

```bash
    docker build -t forum .
    docker run -p 8081:8081 forum
```

server will run on the next route

```
    http://localhost:8081
```

## Audit list

<a href="https://01.alem.school/git/root/public/src/branch/master/subjects/forum/audit" target="_blank">forum audit</a>