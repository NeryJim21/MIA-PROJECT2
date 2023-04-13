package main

import (
	"p2/analizador"
	"fmt"
)

func main() {
	menu()
}

// Función para menú
func menu() {
	fmt.Println("\n\n ")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # #         S I S T E M A   D E   A R C H I V O S   E X T 2         # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # #                                                                 # # #")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	fmt.Println("\n\n\n ")

	for {
		fmt.Print("-> ")
		analizador.GetCommand()
	}

}