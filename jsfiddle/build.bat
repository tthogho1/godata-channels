set GOOS=js
set GOARCH=wasm

go build -o dist/demo.wasm

pause

copy dist/demo.wasm C:\nginx-1.16.1\html
copy demo.css C:\nginx-1.16.1\html
copy offer.html C:\nginx-1.16.1\html
copy wasm_exec.js C:\nginx-1.16.1\html

pause