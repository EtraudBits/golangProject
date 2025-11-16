package main

import "fmt"

func main() {
	
	//variaveis gerais
	var salMinimo float64
	var faturadoMes float64
	
	// leitura dos dados do usuario
	fmt.Println("Qual o Salario Mínimo: ")
	fmt.Scan(&salMinimo)

	fmt.Println("Qual foi o Faturamento Mensal: ")
	fmt.Scan(&faturadoMes)

	//variaveis dos impostos
	valorInss := faturadoMes * 0.2 //20% do faturamento
	dasMei := salMinimo * 0.05 //5% do salario minimo
	dasSimples := faturadoMes * 0.06 //6% do faturamento
	pisCofins := faturadoMes * 0.0365 //3,65% do faturamento
	irCsll := faturadoMes * 0.0768 //7,68% do faturamento
	iss := faturadoMes * 0.05 //5% do faturamento
	
	//calcula o inss conforme o faturamento mensal informado pelo o usuario
	fmt.Println("======INSS PF======")
	if faturadoMes < 7507.49 { //se o valor do faturamento for menor que 7507,49
		fmt.Printf("O valor do INSS será R$ %.2f\n", valorInss) //retorna valorInss
	} else { //senão
		fmt.Println("O valor do INSS será R$ 1.501,49") //retorna
	}

	fmt.Println("======DAS DO MEI======")
		fmt.Printf("Seu DAS (MEI) será: %.2f\n", dasMei)
	
	fmt.Println("=====DAS DO SIMPLES=====")
		fmt.Printf("Seu DAS (Simples) será: R$ %.2f\n", dasSimples)

	fmt.Println ("=====PIS-COFINS=====")
		fmt.Printf("Seu Pis-Cofins (Lucro Presumido) será: R$ %.2f\n", pisCofins)

	fmt.Println("=====IR-CSLL=====")
		fmt.Printf("Seu IR-CSLL (Lucro Presumido) será: R$ %.2f\n", irCsll)
	fmt.Println("=====ISS=====")
		fmt.Printf("Seu ISS (Lucro Presumido) será: R$ %.2f\n", iss)

	//variaveis dos totais
	 inssTotal := valorInss 
	 meiTotal := dasMei 
	 simplesTotal := dasSimples 
	 lucroPresTotal := pisCofins + irCsll + iss 

	fmt.Println("====TOTAIS====")

	fmt.Printf("Total Pessoa Fisica: R$ %.2f + IR\n", inssTotal )

	fmt.Printf("Total MEI: R$ %.2f\n", meiTotal)

	fmt.Printf("Total SIMPLES: R$ %.2f\n", simplesTotal)

	fmt.Printf("Total LUCRO PRESUMIDO: R$ %.2f\n", lucroPresTotal)
}