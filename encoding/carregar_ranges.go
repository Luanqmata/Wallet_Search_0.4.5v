package encoding

import (
	"fmt"
	"log"
	"os"
)

func CarregarRangesDoArquivo(nomeDoArquivo string, ranges map[int]string) {
	file, err := os.Open(nomeDoArquivo)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo de ranges: %v", err)
	}
	defer file.Close()

	for {
		var bits int
		var rangeStr string
		_, err := fmt.Fscanf(file, "%d %s\n", &bits, &rangeStr)
		if err != nil {
			break
		}
		ranges[bits] = rangeStr
	}
}


