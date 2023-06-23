# lab-concorrente-2022.2
Laboratório da disciplina de Concorrente do período 2022.2 - UFCG.


## Executando o programa

Inicialmente, para buildar a aplicação, mude a permissão dos seguintes arquivos:

``` 
chmod +x build.sh
chmod +x sequential.sh
chmod +x concurrent.sh
```

Em seguida, execute:

```
./build.sh
```

Para executar a versão sequencial do programa, use este comando:

```
./sequencial.sh
```

Por último, para executar a versão concorrente do programa:

```
./concurrent.sh
```

## Estratégia de concorrência 

Para aumentar o desempenho do sistema e torná-lo concorrente, usamos uma estratégia de adaptação do padrão produtor-consumidor. Nossa adaptação consiste na criação de 3 canais: um para alocação dos ids de atores, um para alocação dos atores após a requisição ao servidor e outro para atores cuja média por filme já tenha sido calculada. São criadas 3 goroutines com atividades principais: uma para adicionar IDS ao canal de ids, uma para fazer as requisições de ator cujos ids já estejam no canal de ids e outra para o cálculo da média dos atores já alocados no canal de atores. Buscando aumentar ainda mais a eficiência do sistema, nas funções responsáveis por fazer requisição de ator e cálculo de média, foram adicionadas mais algumas goroutines, estas são sincronizadas com waitgroups, uma vez que há a necessidade de sincronização do fechamento dos canais sendo utilizados. Uma outra estratégia para aumentar a eficiência do programa foi criar uma goroutine para fazer a requisição de cada id de filme presente no ator, estas também sincronizadas com o uso de waitgroup.