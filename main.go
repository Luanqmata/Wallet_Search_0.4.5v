package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"wallet_search/encoding"
	"wallet_search/style"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

const (
	memBufferSize = 2 * 1024 * 1024 * 1024 // (metodo 1)
)

var (
	contador          int
	encontrado        bool
	mu                sync.Mutex
	wg                sync.WaitGroup
	ultimaChaveGerada string
	memBuffer         = make([]byte, memBufferSize) //metodo 1 usando a const obs : bate 468k em 1 minuto
	//memBuffer          = make([]byte, 4*1024*1024*1024) // metodo 2 apagando a const obs : nao bate 468k só fica em 465k
	startTime          time.Time
	tamanhoChave       int
	carteiras          map[int]string
	ranges             map[int]string
	minRange           *big.Int
	maxRange           *big.Int
	carteira_escolhida string
)

func init() {
	carteiras = make(map[int]string)
	ranges = make(map[int]string)
}

func converte_Ranges() {
	rangeStr, existe := ranges[tamanhoChave]
	if !existe {
		log.Fatalf("Só existem 160 chaves. %d este número não é aceito.", tamanhoChave)
	}

	valores := strings.Split(rangeStr, "-")
	minRange = new(big.Int)
	minRange.SetString(valores[0], 16)
	maxRange = new(big.Int)
	maxRange.SetString(valores[1], 16)

}

func gerarChavePrivada_aleatoria() string {
	chaveGerada, _ := rand.Int(rand.Reader, new(big.Int).Sub(maxRange, minRange))
	chaveGerada.Add(chaveGerada, minRange)

	chaveHex := hex.EncodeToString(chaveGerada.Bytes())
	if len(chaveHex) < 64 {
		chaveHex = strings.Repeat("0", 64-len(chaveHex)) + chaveHex
	}

	return chaveHex
}

func worker(id int) { //id inicia o multiprocesso
	defer wg.Done()

	for {
		mu.Lock()
		if encontrado {
			mu.Unlock()
			return
		}
		mu.Unlock()

		chave := gerarChavePrivada_aleatoria()
		pubKeyHash := encoding.CreatePublicHash160(chave)
		address := encoding.EncodeAddress(pubKeyHash)
		mu.Lock()
		contador++
		ultimaChaveGerada = chave
		// fmt.Print("\r PV KEY : ", chave, " Carteira: ", address) // isso quebra o codigo porem usado para checkagem
		if address == carteira_escolhida {
			output := fmt.Sprintf("\n\t\t|--------------%s----------------|\n", address) +
				"\t\t|----------------------ATENÇÃO-PRIVATE-KEY-----------------------|\n" +
				fmt.Sprintf("\t\t|%s|\n", chave)

			fmt.Print(output)

			file, err := os.OpenFile("carteira_encontradas.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Erro ao abrir arquivo: %v", err)
			}
			defer file.Close()

			_, err = file.WriteString(output)
			if err != nil {
				log.Printf("Erro ao escrever no arquivo: %v", err)
			}

			encontrado = true
			mu.Unlock()
			return
		}
		mu.Unlock()
	}
}

func main() {
	encoding.CarregarRangesDoArquivo("enderecos/ranges.txt", ranges)
	encoding.CarregarRangesDoArquivo("enderecos/carteiras.txt", carteiras)
	//----------------------------------EXIBE O LOGO---------------------------------------
	fmt.Print("\n\n\n\n\n\n")
	style.Logo_menu()
	//-----------------------------------VERIFICA N° THREADS---------------------------------------
	numCPUs := runtime.NumCPU()
	cpuInfo, _ := cpu.Info()
	cpuModelName := "Desconhecido"
	if len(cpuInfo) > 0 {
		cpuModelName = cpuInfo[0].ModelName
	}
	fmt.Printf("\n\t	Obs: O Seu Computador tem %d threads. (Processador: %s)\n", numCPUs, cpuModelName)

	//---------------------------------PROCURA A CARTEIRA------------------------------------------
	fmt.Print("\n\n\tDigite qual Carteira você vai querer procurar: ")
	var escolha_carteira_chave int
	fmt.Scanln(&escolha_carteira_chave)

	if carteira, ok := carteiras[escolha_carteira_chave]; ok {
		carteira_escolhida = carteira
		fmt.Printf("\n\tCarteira escolhida: %s (bits: %d)\n", carteira_escolhida, escolha_carteira_chave)
	} else { //tratamento de erros.
		fmt.Printf("Número de carteira %d não suportado. Escolha um valor entre 1 e %d.\n", escolha_carteira_chave, len(carteiras))
		return
	}
	//------------------------------------CONVERTE O RANGE-----------------------------------------
	tamanhoChave = escolha_carteira_chave
	converte_Ranges()

	if _, ok := ranges[tamanhoChave]; !ok { //tratamento de erros.
		fmt.Printf("	Tamanho de chave %d não suportado. Escolha um valor entre 1 e %d.\n", tamanhoChave, len(ranges))
		return
	}

	//------------------------------------MENSAGEM DE MODOS----------------------------------------

	style.Modos()
	var escolha_modo int
	fmt.Scanln(&escolha_modo)

	var numThreads int
	switch escolha_modo {
	case 1:
		numThreads = 3
	case 2:
		numThreads = 5
	case 3:
		numThreads = 11
	case 4:
		numThreads = 17
	case 5:
		numThreads = 24
	case 6:
		numThreads = runtime.NumCPU()
	default:
		fmt.Println("	Escolha inválida. Usando o SECURE MODE...  (20%) - CPU 58°C - 117K Chaves P/seg.")
		numThreads = 4
	}
	style.Mensagem_iniciando(numThreads)

	runtime.GOMAXPROCS(numThreads)
	startTime = time.Now()

	// Inicia goroutines
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go worker(i)
	}
	//--------------------------------------EXIBIÇÃO EM TEMPO REAL--------------------------------------

	style.MonitorarChaves(&mu, &contador, &encontrado, startTime, &ultimaChaveGerada)

	//-------------------------------------------MENSAGEM FINAL----------------------------------

	wg.Wait()
	style.Mensagem_final(numThreads, contador)
}
