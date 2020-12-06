package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/ldicarlo/legifrss/server/dila"
	"github.com/ldicarlo/legifrss/server/generate"
	"github.com/ldicarlo/legifrss/server/models"
	"github.com/ldicarlo/legifrss/server/token"
	"github.com/ldicarlo/legifrss/server/utils"
)

var clientId string
var clientSecret string

func init() {
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic("missing env file")
	}
	clientId = envs["client_id"]
	clientSecret = envs["client_secret"]

	if clientId == "" || clientSecret == "" {
		panic("Missing one of the env params")
	}

}

func Start() (str string, result string) {
	err, token := token.GetToken(clientId, clientSecret)
	utils.ErrCheckStr(err)

	err, res := dila.FetchJORF(token)
	utils.ErrCheckStr(err)

	var jorfContents []models.JOContainerResult
	for _, jorf := range res.Containers {
		err, nextContent := dila.FetchCont(token, jorf.Id)
		utils.ErrCheckStr(err)
		jorfContents = append(jorfContents, nextContent)
	}
	list := utils.ExtractAndConvertDILA(jorfContents)
	generate.Generate(list)
	return "", "ok"
}

func main() {
	fmt.Println(Start())
}
