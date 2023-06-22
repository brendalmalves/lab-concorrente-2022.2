#!/bin/bash

cd src
go build sequencial.go
#go build concurrent.go

if [ $? -eq 0 ]; then
    echo "Build concluído com sucesso!"
else
    echo "Erro ao realizar o build."
fi