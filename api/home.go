package handler

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func HandleHome(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tmpl, err := template.New("index").Parse("./html/index.html")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error parsing template",
		}, nil
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error executing template",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: buf.String(),
	}, nil
}
