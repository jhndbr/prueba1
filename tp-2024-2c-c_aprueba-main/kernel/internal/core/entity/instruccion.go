package entity

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

// Instruction sera la entidad que represente las operaciones a realizar.
type Instruction struct {
	Code string
	Args []string
}

func readFile(filename string) ([]Instruction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error("Error al cerrar el archivo", "COD", err)
		}
	}(file)

	var instrucciones []Instruction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Fields(linea)
		if len(partes) > 0 {
			comando := partes[0]
			argumentos := partes[1:]
			instruccion := Instruction{Code: comando, Args: argumentos}
			instrucciones = append(instrucciones, instruccion)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return instrucciones, nil
}

func ObtenerInstrucciones(filename string) ([]Instruction, error) {

	instrucciones, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	return instrucciones, nil
}
