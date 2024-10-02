package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"meugo/encoding"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

const (
	memBufferSize = 4 * 1024 * 1024 * 1024 // usando 4GB pro buffer
)

var (
	contador          int
	encontrado        bool
	mu                sync.Mutex
	wg                sync.WaitGroup
	ultimaChaveGerada string
	memBuffer         = make([]byte, memBufferSize) //metodo 1 usando a const
	//memBuffer          = make([]byte, 4*1024*1024*1024) // metodo 2 apagando a const
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

func setupRanges() {
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

func gerarChavePrivada() string {
	var chaveGerada *big.Int
	for {
		chaveGerada, _ = rand.Int(rand.Reader, new(big.Int).Sub(maxRange, minRange))
		chaveGerada.Add(chaveGerada, minRange)

		if chaveGerada.Cmp(minRange) >= 0 && chaveGerada.Cmp(maxRange) <= 0 {
			break
		}
	}

	chaveHex := hex.EncodeToString(chaveGerada.Bytes())
	if len(chaveHex) < 64 {
		chaveHex = strings.Repeat("0", 64-len(chaveHex)) + chaveHex
	}

	copy(memBuffer[:len(chaveHex)], chaveHex)

	return chaveHex
}

func worker(id int) {
	defer wg.Done()

	for {
		mu.Lock()
		if encontrado {
			mu.Unlock()
			return
		}
		mu.Unlock()

		chave := gerarChavePrivada()
		pubKeyHash := encoding.CreatePublicHash160(chave)
		address := encoding.EncodeAddress(pubKeyHash)
		mu.Lock()
		contador++
		ultimaChaveGerada = chave
		//fmt.Print("\n", ultimaChaveGerada) // ver as chaves geradas
		//fmt.Print("\n", address) // ver os endereços das chaves geradas
		if address == carteira_escolhida {
			fmt.Printf("\n\n	|--------------%s----------------|\n", address)
			fmt.Printf("	|----------------------ATENÇÃO-PRIVATE-KEY-----------------------|")
			fmt.Printf("\n	|%s|\n", chave)
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

	fmt.Print("\n\n\n\n\n\n")
	fmt.Println(`                                              	        			-_~ BEM VINDO ~_-
						________         __  __         __       _______                            __
						|  |  |  |.---.-.|  ||  |.-----.|  |_    |     __|.-----..---.-..----..----.|  |--.
						|  |  |  ||  _  ||  ||  ||  -__||   _|   |__     ||  -__||  _  ||   _||  __||     |
						|________||___._||__||__||_____||____|   |_______||_____||___._||__|  |____||__|__|  ~ 0.3.9v  By:Ch!iNa ~
									  
									    -_~ :Carteira Puzzle: ~_-
	`)

	// Obter o número de CPUs disponíveis e o nome do processador
	numCPUs := runtime.NumCPU()
	cpuInfo, _ := cpu.Info()
	cpuModelName := "Desconhecido"
	if len(cpuInfo) > 0 {
		cpuModelName = cpuInfo[0].ModelName
	}
	fmt.Printf("\n	Obs: O Seu Computador tem %d threads. (Processador: %s)\n", numCPUs, cpuModelName)

	//---------------------------------------------------------------------------
	fmt.Print("\n\nDigite qual Carteira você vai querer procurar: ")
	var escolha_carteira_chave int
	fmt.Scanln(&escolha_carteira_chave)

	if carteira, ok := carteiras[escolha_carteira_chave]; ok {
		carteira_escolhida = carteira
		fmt.Printf("\nCarteira escolhida: %s (bits: %d)\n", carteira_escolhida, escolha_carteira_chave)
	} else {
		fmt.Printf("Número de carteira %d não suportado. Escolha um valor entre 1 e %d.\n", escolha_carteira_chave, len(carteiras))
		return
	}

	tamanhoChave = escolha_carteira_chave
	setupRanges()

	if _, ok := ranges[tamanhoChave]; !ok {
		fmt.Printf("	Tamanho de chave %d não suportado. Escolha um valor entre 1 e %d.\n", tamanhoChave, len(ranges))
		return
	}

	//----------------------------------------------------------------------------

	time.Sleep(1 * time.Second)
	fmt.Println("\n\n\n						- Random Mode -")

	fmt.Println("\n\n		Modo 1: Easy   (15%) - CPU 60°C - 125 RPM  - 86k  Chaves P/seg")
	fmt.Println("		Modo 2: Seguro (25%) - CPU 65°C - 170 RPM  - 150k Chaves P/seg")
	fmt.Println("		Modo 3: Medium (50%) - CPU 78°C - 220 RPM (+) - 322k Chaves P/seg")
	fmt.Println("		Modo 4: Hard   (75%) - CPU 83°C - 255 RPM (+) - 415k Chaves P/seg")
	fmt.Println("\n\n					- - - A partir daqui já está fritando o CPU - - -")

	fmt.Println("\n\n		Modo 5: Um pouquinho pra daqui a pouco !!PERIGO!! (94%) - CPU 80°C. - 450k Chaves P/seg")
	fmt.Println("		Modo 6: Não Estou nem Aí para meu computador quero que ele queime (100%) - CPU 85°C. - 468k Chaves P/seg")

	fmt.Print("\n\n Input: Escolha o modo de acordo com o número correspondente: ")

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
	fmt.Printf("\r\n\n  		Threads: %d iniciando .", numThreads)
	time.Sleep(1 * time.Second)
	fmt.Printf("\r  		Threads: %d iniciando ..", numThreads)
	time.Sleep(1 * time.Second)
	fmt.Printf("\r  		Threads: %d iniciando ...\n\n", numThreads)
	time.Sleep(1 * time.Second)

	runtime.GOMAXPROCS(numThreads)
	startTime = time.Now()

	// Inicia goroutines
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go worker(i)
	}
	//----------------------------------------------------------------------------
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mu.Lock()
			if encontrado {
				mu.Unlock()
				break
			}
			duration := time.Since(startTime).Seconds()
			chavesPorSegundo := float64(contador) / duration
			fmt.Printf("\r  N° Threads Usados: %d | Chaves Geradas: %d | Chaves Por Segundo: %.2f |", numThreads, contador, chavesPorSegundo)
			fmt.Print("Ultima Chave Gerada: ", ultimaChaveGerada, " |")
			mu.Unlock()
		}
	}()
	//-----------------------------------------------------------------------------
	wg.Wait()

	fmt.Print("	|--------------------------------------------------by-Luan-BSC---|")
	fmt.Print("\n	|-----------------------China-LOOP-MENU------------------------- |")
	fmt.Printf("\n	|		Threads usados: %d		                 |", numThreads)
	fmt.Print("\n	|		Chaves Analisadas:	", contador)
	fmt.Print("\n	|________________________________________________________________|")
}
