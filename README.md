![01-edu-system-blue](https://raw.githubusercontent.com/GoArtyom/study/1a66b22c5b511ccce94b582481a45dfd7f001d3a/alem.svg)

# Forum-Security

This project follows the same [principles](https://01.alem.school/git/root/public/src/branch/master/subjects/forum/README.md) as the first project (forum).


## ER Diagram

![ERD](https://images-ext-1.discordapp.net/external/_YbIrYIAmICRZ9kqceoEctlY-ASq2zldW0SUUs_LKEs/https/i.imgur.com/TfY6UA4.png?format=webp&quality=lossless&width=1041&height=863)

## Clone repository

```bash
git clone git@git.01.alem.school:aecheist/forum-security.git
```

## Move to the direcroty

```bash
cd forum-image-security
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
    docker run -p 8080:8080 forum
```

server will run on the next route

```
    http://localhost:8080
```

## Audit list

<a href="https://01.alem.school/git/root/public/src/branch/master/subjects/forum/security/audit" target="_blank">forum-security audit</a>
