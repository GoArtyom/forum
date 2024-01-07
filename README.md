![01-edu-system-blue](https://raw.githubusercontent.com/GoArtyom/study/1a66b22c5b511ccce94b582481a45dfd7f001d3a/alem.svg)

# Forum-Advanced-Features

This project follows the same [principles](https://01.alem.school/git/root/public/src/branch/master/subjects/forum/README.md) as the first project (forum).

In `forum advanced features`, we have implemented the following features :

- We create notify users when their posts are :

  - liked/disliked
  - commented

- We have created an activity page that tracks the activity of the user himself. In other words, a page that :

  - Shows the user created posts
  - Shows where the user left a like or a dislike
  - Shows where and what the user has been commenting. For this, the comment will have to be shown, as well as the post commented

- Created a section where you will be able to Edit/Remove posts and comments.


## ER Diagram

![ERD](https://github.com/ArtEmerged/study/blob/main/db_forum.png?raw=true)

## Clone repository

```bash
git clone git@git.01.alem.school:aecheist/forum-advanced-features.git
```

## Move to the direcroty

```bash
cd forum-advanced-features
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
    https://localhost:8080
```

## Audit list

<a href="https://01.alem.school/git/root/public/src/branch/master/subjects/forum/advanced-features/audit" target="_blank">forum-advanced-features audit</a>
