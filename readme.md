
# File Renaming With Path Handling


```
go build main.go
./main 
Flags:
-input index.html -output index.html -replace "js/app.js" -workdir . 
```

```
- name: Run File Renamer Action
  uses: pavansh/generate-hash-for-caching@v1.0
  with:
    input: 'path/to/input.html'
    output: 'path/to/output.html'
    replace: 'app.js,styles.css'
    workdir: ./build

```
