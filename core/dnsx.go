package core

import (
	"fmt"
	"os"
	"os/exec"
)

// RunDNSX executa o dnsx com input e output definidos
func RunDNSX(inputPath, outputPath string, extraArgs ...string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir input: %v", err)
	}
	defer inputFile.Close()

	// Cria output
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro ao criar output: %v", err)
	}
	defer outputFile.Close()

	cmdArgs := []string{"-silent", "-a", "-resp", "-nc"}
	cmdArgs = append(cmdArgs, extraArgs...)

	cmd := exec.Command("dnsx", cmdArgs...)
	cmd.Stdin = inputFile
	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	fmt.Println("[+] Rodando dnsx...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar dnsx: %v", err)
	}

	fmt.Println("[+] dnsx finalizado! Resultado salvo em:", outputPath)
	return nil
}
