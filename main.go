package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"meugo/encoding"
	"runtime"
	"sync"
	"time"
)

const (
	prefix        = "00000000000000000000000000000000000000000000000" // Prefixo da chave
	memBufferSize = 2 * 1024 * 1024 * 1024                            // usando 2gb pro buffer
)

var chaves_desejadas = map[string]bool{
	"1BY8GQbnueYofwSuFAT3USAhGjPrkxDdW9": true, //67
	"1MVDYgVaSN6iKKEsbzRUAYFrYJadLYZvvZ": true, //68
}

// 4 byte xx xx xx xx

var (
	contador          int
	encontrado        bool
	mu                sync.Mutex
	wg                sync.WaitGroup
	ultimaChaveGerada string
	memBuffer         = make([]byte, memBufferSize)
)

func gerarChavePrivada() string {
	suffix := make([]byte, 9) // se for carteira impar coloca 1 byte a mais e no enconde ajusta com o tamanho de carcters da key
	_, err := rand.Read(suffix)
	if err != nil {
		log.Fatalf("Falha ao gerar chave: %v", err)
	}
	chaveGerada := prefix + hex.EncodeToString(suffix)[:17] //linha usada para gerar somente 17carcters

	// Manipulação eficiente do buffer (exemplo genérico)
	copy(memBuffer[:len(suffix)], suffix)

	return chaveGerada
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
		//fmt.Print("\n", ultimaChaveGerada) //mostrar Private keys
		//fmt.Print("\n", address) // mostrar Carteira
		if chaves_desejadas[address] {
			fmt.Printf("\n\n|--------------%s----------------|\n", address)
			fmt.Printf("|----------------------ATENÇÃO-PRIVATE-KEY-----------------------|")
			fmt.Printf("\n|%s|\n", chave)
			encontrado = true
			mu.Unlock()
			return
		}
		mu.Unlock()
	}
}

func main() {
	fmt.Print("\n\n\n\n\n\n")
	fmt.Println(`                                              	        			-_~ BEM VINDO ~_-
						________         __  __         __       _______                            __
						|  |  |  |.---.-.|  ||  |.-----.|  |_    |     __|.-----..---.-..----..----.|  |--.
						|  |  |  ||  _  ||  ||  ||  -__||   _|   |__     ||  -__||  _  ||   _||  __||     |
						|________||___._||__||__||_____||____|   |_______||_____||___._||__|  |____||__|__|  ~ 0.3.9v  By:Ch!iNa ~

							
									  -_~ Carteira Puzzle: 66 67 ~_-
	`)
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n	Obs: Computador do criador tem 28 threads... (KIT XEON | 2.40 ~ 3.00 GHz | E5 2680 v4)")
	fmt.Println("\n\n\n						- Random Mode -")
	fmt.Println("\n\n		Modo 1: Easy   (15%) - CPU 56°C - 125 RPM  - 86k  Chaves P/seg")
	fmt.Println("		Modo 2: Seguro (25%) - CPU 60°C - 170 RPM  - 150k Chaves P/seg")
	fmt.Println("		Modo 3: Medium (50%) - CPU 73°C - 220 RPM+ - 322k Chaves P/seg")
	fmt.Println("		Modo 4: Hard   (75%) - CPU 78°C - 255 RPM+ - 415k Chaves P/seg")
	fmt.Println("\n\n					- - - A partir daqui ja está fritando o CPU - - -")
	fmt.Println("\n\n		Modo 5: Um pouquinho pra daqui apoco !!PERIGO!! (94%) - CPU 80°C. - 450k Chaves P/seg")
	fmt.Println("		Modo 6: Não Estou nem Aí para meu computador quero que ele queime (100%) - CPU 85°C. - 468k Chaves P/seg")
	fmt.Print("\n\n	Input: Escolha o modo de acordo com o número correspondente: ")

	var escolha int
	fmt.Scanln(&escolha)

	var numThreads int
	switch escolha {
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
	fmt.Printf("\n  		Iniciando modo %d .", escolha)
	time.Sleep(1 * time.Second)
	fmt.Printf("\n  		Iniciando modo %d ..", escolha)
	time.Sleep(1 * time.Second)
	fmt.Printf("\n  		Iniciando modo %d ...\n\n", escolha)
	time.Sleep(1 * time.Second)

	runtime.GOMAXPROCS(numThreads)

	// Inicia goroutines
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go worker(i)
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			mu.Lock()
			if encontrado {
				mu.Unlock()
				break
			}
			fmt.Printf("\r 	N° Threads Usados: %d | Chaves Geradas: %d | ", numThreads, contador)
			mu.Unlock()
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			mu.Lock()
			if encontrado {
				mu.Unlock()
				break
			}
			fmt.Print("Ultima Chave Gerada: ", ultimaChaveGerada, " |")
			mu.Unlock()
		}
	}()

	wg.Wait()

	fmt.Print("\n\n	|--------------------------------------------------by-Luan-BSC---|")
	fmt.Print("\n	|-----------------------China-LOOP-MENU------------------------- |")
	fmt.Printf("\n	|		Threads usados: %d		                 |", numThreads)
	fmt.Print("\n	|		Chaves Analisadas:	", contador)
	fmt.Print("\n	|________________________________________________________________|")
}
