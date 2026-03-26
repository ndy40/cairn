# https://just.systems

default:
    echo 'Hello, world!'

build:
    go build -o cairn  ./cmd/cairn/

_docs:
    hugo build

update-docs: _docs
    git add docs/
    git commit -m 'docs: Updated documentation'

publish-docs: update-docs
    git push

serve-docs: update-docs
    hugo serve



