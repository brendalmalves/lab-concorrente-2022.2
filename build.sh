#!/bin/bash

go build sequencial.go
#go build concorrente.go

if [ $? -eq 0 ]; then
    echo "Build concluído com sucesso!"
else
    echo "Erro ao realizar o build."
fi