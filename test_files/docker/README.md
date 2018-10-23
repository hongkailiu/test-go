# Dockerfile

## http

```sh
$ buildah bud --format=docker -f Dockerfile -t quay.io/hongkailiu/test-go:http-0.0.1 .
$ buildah push --creds=hongkailiu d58cbf2a06aa docker://quay.io/hongkailiu/test-go:http-0.0.1
```
