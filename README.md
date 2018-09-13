#health-checker

## build

```
sudo docker run -ti -v $(pwd):/go/src/github.com/leocomelli/health-checker leocomelli/health-checker:build bash -c 'go get && go build -o dist/health-checker'
```

