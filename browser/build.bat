set GOOS=js
set GOARCH=wasm

mkdir dist

go build -o dist/data-channels.wasm
rem tinygo build -o dist/data-channels.wasm

pause

copy dist/data-channels.wasm  C:\nginx-1.16.1\html\
copy index.html C:\nginx-1.16.1\html\

pause