package style

import (
	"fmt"
	"sync"
	"time"
)

func Logo_menu() {
	fmt.Println(`                                              	        			-_~ BEM VINDO ~_-
						________         __  __         __       _______                            __
						|  |  |  |.---.-.|  ||  |.-----.|  |_    |     __|.-----..---.-..----..----.|  |--.
						|  |  |  ||  _  ||  ||  ||  -__||   _|   |__     ||  -__||  _  ||   _||  __||     |
						|________||___._||__||__||_____||____|   |_______||_____||___._||__|  |____||__|__|  ~ 0.3.9v  By:Ch!iNa ~
									  
									    -_~ :Carteira Puzzle: ~_-
	`)
}

func Modos() {
	time.Sleep(1 * time.Second)
	fmt.Println("\n\n\n\t\t\t\t- - - RANDOM MODE - - -")

	fmt.Println("\n\t\tModo 1:  Easy   (15%)  | CPU 60°C  | 125 RPM  |  86k Chaves P/seg")
	fmt.Println("\t\tModo 2:  Seguro (25%)  | CPU 65°C  | 170 RPM  | 150k Chaves P/seg")
	fmt.Println("\t\tModo 3:  Medium (50%)  | CPU 78°C  | 220 RPM  | 322k Chaves P/seg")
	fmt.Println("\t\tModo 4:  Hard   (75%)  | CPU 83°C  | 255 RPM  | 415k Chaves P/seg")

	fmt.Println("\n\t\t\t- - - - - - - A partir daqui já está fritando o CPU - - - - - - -")

	fmt.Println("\n\t\tModo 5:  Um pouquinho pra daqui a pouco !! PERIGO !! (94%) | CPU 80°C | 450k Chaves P/seg")
	fmt.Println("\t\tModo 6:  Extreme (100%) | CPU 85°C | 468k Chaves P/seg")

	fmt.Print("\n\n\t\tInput: Escolha o modo de acordo com o número correspondente: ")
}

func Mensagem_iniciando(numThreads int) {
	fmt.Printf("\r\n\n  		Iniciando com %d Threads .", numThreads)
	time.Sleep(1 * time.Second)
	fmt.Printf("\r  		Iniciando com %d Threads . .", numThreads)
	time.Sleep(1 * time.Second)
	fmt.Printf("\r  		Iniciando com %d Threads . . .\n\n\n\n\n", numThreads)
	time.Sleep(1 * time.Second)
}

func MonitorarChaves(mu *sync.Mutex, contador *int, encontrado *bool, startTime time.Time, ultimaChaveGerada *string) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mu.Lock()
			if *encontrado {
				mu.Unlock()
				break
			}
			duration := time.Since(startTime).Seconds()
			chavesPorSegundo := float64(*contador) / duration
			fmt.Printf("\r\t\t| Chaves Geradas: %d | Chaves Por Segundo: %.2f | Última Chave Gerada: %s |", *contador, chavesPorSegundo, *ultimaChaveGerada)
			mu.Unlock()
		}
	}()
}

func Mensagem_final(numThreads int, contador int) {
	fmt.Print("\t\t|-----------------------China-LOOP-MENU------------------------- |")
	fmt.Printf("\n\t\t|		Threads usados: %d		                 |", numThreads)
	fmt.Printf("\n\t\t|		Chaves Analisadas:	%d ", contador)
	fmt.Print("\n\t\t|________________________________________________________________|")
}
