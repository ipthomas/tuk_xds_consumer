package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukhttp"
)

func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.SetFlags(log.Lshortfile)
	nhs := req.QueryStringParameters[tukcnst.TUK_EVENT_QUERY_PARAM_NHS]
	mhd := tukhttp.MHDRequest{URL: os.Getenv("MHD_SERVER_URL"), PID_OID: tukcnst.NHS_OID_DEFAULT, PID: nhs}
	if err := tukhttp.NewRequest(&mhd); err != nil {
		queryResponse(http.StatusInternalServerError, err.Error(), tukcnst.TEXT_PLAIN)
	}
	return queryResponse(http.StatusOK, string(mhd.Response), tukcnst.TEXT_HTML)
}
func setAwsResponseHeaders(contentType string) map[string]string {
	awsHeaders := make(map[string]string)
	awsHeaders["Server"] = "ICB Workflow Service"
	awsHeaders["Access-Control-Allow-Origin"] = "*"
	awsHeaders["Access-Control-Allow-Headers"] = "accept, Content-Type"
	awsHeaders["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
	awsHeaders[tukcnst.CONTENT_TYPE] = contentType
	return awsHeaders
}
func queryResponse(statusCode int, body string, contentType string) (*events.APIGatewayProxyResponse, error) {
	log.Println(body)
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    setAwsResponseHeaders(contentType),
		Body:       body,
	}, nil
}
