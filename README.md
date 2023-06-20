## Descrição da API
API de Normalização das Informações dos Pedidos do Sistema Legado.

Está API tem como objetivo carregar os pedidos do sistema legado a partir de arquivos no formato TXT desnormalizado e devolver as informações normalizadas no formato JSON.

Por ser uma API apenas para normalização das informações, a API mantém apenas as informações do último arquivo importado, ou seja, cada vez que um arquivo for importado será descartada as informações do arquivo importanto anteriormente.

Para garantir a consistência das informações caso ocorra algum erro durante a importação do arquivo, todo o conteúdo do arquivo será desconsiderado e será retornado a descrição do erro. 

A formatação e validação de todo o conteúdo dos registros do arquivo são validados previamente e serão retornados os registros com problema e os seus respectivos problemas.

A documentação dos endpoints estão disponíveis na própria API [localhost:9000/api/docs](localhost:9000/api/docs).

## Etapas para poder executar a API na máquina local
1. Docker instalado. [Documentação](https://docs.docker.com/engine/install/)
2. Docker-compose instalado. [Documentação](https://docs.docker.com/compose/install/linux/)
3. Go instalado. [Documentação](https://go.dev/doc/install)

#### Obs: Os comandos a seguir devem ser executados na pasta raiz do projeto.

4. Subir os serviços de banco de dados e de cache utilizados pela API
    ```
    make docker-compose-up
   
    ou
    
    docker-compose up -d
    ```

5. Executar a API

    - Direto na máquina local
        ```
        make go-run

        ou
        
        go run server.go
        ```

    - Dentro do Docker
        ```
        make docker-build
        make docker-run
        ```

6. Endpoint da API [localhost:9000](localhost:9000)

## Descrição Técnica
1. Arquitetura: Utilização dos principios de Clean Architecture que torna a API altamente testável, flexível e independente de frameworks, banco de dados, etc.
2. Repositório: Por ser uma API apenas para normalização das informações utilizei o conceito de banco em memória, portanto se a API for reiniciada as informações do arquivo importado serão perdidas sendo necessário importar novamente. Para demostrar conhecimento na utilização de banco de dados, deixei implementado a possibilidade de utilizar o banco de dados Postgres que foi escolhido por ser um serviço de banco relacional robusto, completo e open source que atende perfeitamente desde pequenas aplicações até aplicações robustas e compatível com serviços de banco de dados em nuvem.
3. Cache: Por ser uma API apenas para normalização das informações não vejo muita necessidade em utilizar cache principalmente utilizando o conceito de banco em memória mas acabei utilizando para demostrar conhecimento. Realizei a implementação de cache na consulta de pedido por ID. Sempre que um novo arquivo for importado será realizado a limpeza do cache. A utilização de cache aumenta a complexidade da API para manter a consistência dos dados do cache. Utilizei o Redis por ser um serviço de cache robusto, open source, amplamente utilizado e compatível com serviço de cache em nuvem como o memory store do GCP.
4. Migration: Para versionamento de alterações no banco de dados.
5. Health Check: [localhost:9000/api/healthz](localhost:9000/api/healthz) para monitorar se a aplicação está no ar e se os serviços de banco de dados e cache estão funcionando.
6. Documentação da API: Utilizado swagger para manter a documentação da API atualizada de forma automática utilizando tags no código. Documentação disponível na própria API [localhost:9000/api/docs](localhost:9000/api/docs).
7. CORS: Para poder permitir requisições de origem diferente da API.
8. Middleware Logger: para logar o resultado de todas as requisições contendo informações da origem da requisição e request id para ajudar no troubleshooting da aplicação. O request id ajuda a rastrear uma mesma requisição por diversos micro serviços e as informações da origem ajudam a identificar se o problema está relacionado a uma origem ou dispositivo específico.
9. Variáveis de Ambiente: Carregamento de configurações da API por meio de variáveis de ambiente ou arquivo de configuração "config.env" na pasta raiz da aplicação.
10. Testes Unitários e de Integração: Gosto de utilizar TDD para iniciar o desenvolvimento a partir da camada de regras de negócio e depois o desenvolvimento das demais camadas.
11. Desenho da API: Utilizado Mermaid Markdown para versionamento de alterações do desenho. Desenho disponível na pasta /docs.
12. O projeto já contém um arquivo config.env com todas as configurações necessárias para poder executar a API no ambiente local.
13. Docker-compose: Para poder subir os serviços de banco de dados e cache para poder rodas a API no ambiente local.
14. Dockerfile: Para poder realizar o build da API e gerar imagem docker para rodar no ambiente local.
15. Github Action: Para validar PR e Push para a branch main iniciando um processo de CI/CD contemplando a execução dos testes unitários, testes de integração e build do projeto.
16. Makefile: Para poder executar de forma simples diversos comandos.


## Geração da Documentação da API - Swagger

1. Instalar na máquina o [go-swagger](https://goswagger.io/install.html)
```
go get github.com/go-swagger/go-swagger/cmd/swagger@latest
```

2. Especificação das tags para geração automática da documentação [aqui](https://goswagger.io/use/spec.html)

3. Executar o comando abaixo na pasta raiz da API para atualizar a documentação
    ```
    make swagger
    ```

## Observação
1. Levando em consideração a Notação Big O para criação de algoritmos mais eficientes, utilizei alguns recursos como a utilização de maps para minimizar a utilização de loops.
2. Para poder utilizar o banco de dados Postgres precisa alterar apenas três alterações no arquivo server.go.

    Na linha 29 substituir a importação do package
    ```
    "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository/in_memory" 
    ```
    por
    ```
    "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository/postgres"
    ```

    Tirar os comentário das linhas 60 a 67

    Na linha 70 substituir
    ```
    repository, err := repository.NewInMemory(config)
    ```
    por
    ```
    repository, err := repository.NewPostgres(config)
    ```
3. Apesar de não ter aplicado nesse projeto também tenho conhecimento do padrão conventional commits.

    #### **Obs:** Com certeza tem mais melhorias a ser feita tanto no código quanto na documentação. Melhoria contínua deve fazer parte da vida útil de toda aplicação.

# Espero que gostem bastante do projeto que entreguei ;)
