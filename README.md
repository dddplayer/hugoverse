# Hugoverse
DDD SSG (Static Site Generator) Hugo.

## Development

**Setup development environment**

```shell
go install
```

**Build**

```shell
go build -o hugov
```

**Run**

```shell
➜  hugoverse git:(main) ✗ ./hugov 
Usage:
  hugov [command]

Commands:
    build:  generate static site for Hugo project
   server:  start the headless CMS server
     demo:  create demo Hugo project
  version:  show hugoverse command version

Example:
  hugov build -p pathspec/to/your/hugo/project
```

## Create Hugo Demo Project

```shell
➜  hugoverse git:(main) ✗ ./hugov demo
demo dir: /var/folders/rt/bg5xpyj51f98w79j6s80wcr40000gn/T/hugoverse-temp-dir782641825
```

**Result:**

```shell
➜  hugoverse git:(main) ✗ cd /var/folders/rt/bg5xpyj51f98w79j6s80wcr40000gn/T/hugoverse-temp-dir782641825
➜  hugoverse-temp-dir782641825 tree
.
├── config.toml
├── layouts
│   ├── _default
│   │   └── single.html
│   └── index.html
├── mycontent
│   └── blog
│       └── post.md
├── myproject.txt
└── themes
    └── mytheme
        └── mytheme.txt

7 directories, 6 files

```