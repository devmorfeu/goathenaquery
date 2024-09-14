package goathenaquery

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

// QueryConfig Estrutura para encapsular as configurações de execução da query
type QueryConfig struct {
	WaitInterval int    // Tempo de espera entre as tentativas (segundos)
	MaxAttempts  int    // Número máximo de tentativas
	Region       string // Região AWS para a sessão
}

// AthenaQueryExecutor Struct para o executor de queries no Athena
type AthenaQueryExecutor struct {
	svc *athena.Athena
}

// NewAthenaQueryExecutor Função para inicializar o executor com uma nova sessão AWS
func NewAthenaQueryExecutor(config QueryConfig) (*AthenaQueryExecutor, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})
	if err != nil {
		return nil, err
	}

	return &AthenaQueryExecutor{
		svc: athena.New(sess),
	}, nil
}

// ExecuteQuery Função para executar a query no Athena
func (e *AthenaQueryExecutor) ExecuteQuery(input *athena.StartQueryExecutionInput, config QueryConfig) ([][]string, error) {
	// Inicia a execução da query no Athena
	result, err := e.svc.StartQueryExecution(input)
	if err != nil {
		return nil, err
	}

	// Espera pelos resultados da query
	return e.waitForResults(result.QueryExecutionId, config)
}

// Função que espera pelos resultados da query
func (e *AthenaQueryExecutor) waitForResults(queryExecutionID *string, config QueryConfig) ([][]string, error) {
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		statusInput := &athena.GetQueryExecutionInput{
			QueryExecutionId: queryExecutionID,
		}

		result, err := e.svc.GetQueryExecution(statusInput)
		if err != nil {
			return nil, err
		}

		status := *result.QueryExecution.Status.State
		if status == athena.QueryExecutionStateSucceeded {
			return e.fetchResults(queryExecutionID)
		} else if status == athena.QueryExecutionStateFailed || status == athena.QueryExecutionStateCancelled {
			return nil, errors.New("query failed or was cancelled")
		}

		// Aguarda pelo intervalo definido
		time.Sleep(time.Duration(config.WaitInterval) * time.Second)
	}

	return nil, errors.New("query did not complete within the maximum number of attempts")
}

// Função para buscar os resultados da query Athena
func (e *AthenaQueryExecutor) fetchResults(queryExecutionID *string) ([][]string, error) {
	input := &athena.GetQueryResultsInput{
		QueryExecutionId: queryExecutionID,
	}

	results, err := e.svc.GetQueryResults(input)
	if err != nil {
		return nil, err
	}

	// Converte os resultados para [][]string
	return parseResults(results), nil
}

// Função auxiliar para converter os resultados em uma matriz de strings
func parseResults(results *athena.GetQueryResultsOutput) [][]string {
	var queryResults [][]string

	for _, row := range results.ResultSet.Rows {
		var rowValues []string
		for _, data := range row.Data {
			if data.VarCharValue != nil {
				rowValues = append(rowValues, *data.VarCharValue)
			} else {
				rowValues = append(rowValues, "")
			}
		}
		queryResults = append(queryResults, rowValues)
	}

	return queryResults
}
