webserver stuff with go

## get

```
curl http://localhost:8080
```

## upload

```
curl -i -F filedata=@foo.zip -F filedata=@bar.zip http://localhost:8080/post
```
