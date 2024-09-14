# GoAthenaQuery

## Purpose

GoAthenaQuery is a Go library designed to simplify the process of querying AWS Athena. It abstracts the polling process for results stored in S3, making it easier to execute and retrieve results from Athena queries.

## How to Use

To use GoAthenaQuery, you need to follow these steps:

1. Import the library in your Go project:

```go
package main

import "github.com/devmorfeu/goathenaquery"

func main () {
	
	config := goathenaquery.QueryConfig{
		WaitInterval: 10,    // Tempo de espera entre as tentativas (segundos)
		MaxAttempts:  5,     // Número máximo de tentativas
		Region:       "us-west-2", // Região AWS para a sessão
	}

	executor, err := goathenaquery.NewAthenaQueryExecutor(config)
	if err != nil {
		log.Fatal(err)
	}

	input := &athena.StartQueryExecutionInput{
		// Seus parâmetros de consulta aqui
	}

	results, err := executor.ExecuteQuery(input, config)
	if err != nil {
		log.Fatal(err)
	}
}
```

## How to Contribute

To contribute to GoAthenaQuery, follow these steps:

1. Fork this repository.
2. Create a branch: `git checkout -b <branch_name>`.
3. Make your changes and commit them: `git commit -m '<commit_message>'`
4. Push to the original branch: `git push origin <project_name>/<location>`
5. Create the pull request.
6. Alternatively, see the GitHub documentation on creating a pull request.
7. After the pull request is merged, you can delete your branch.
8. By contributing to this project, you agree to the terms of the Apache 2.0 license.
9. For more information, see the CONTRIBUTING file.