package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadMeshFromFile(modelPath string, texturePath string) (MeshData, error) {
	file, err := os.Open(modelPath)
	if err != nil {
		log.Println(err)
		return MeshData{}, err
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	verts := make([]Vec3, 0)

	for {
		line, err := scanner.ReadString('\n')
		if err == io.EOF {
			if len(line) != 0 {
				log.Println(line, "last line")
			}
			break
		}
		if err != nil {
			log.Println(err)
			return MeshData{}, err
		}
		if line[0:2] == "v " {
			line = strings.Replace(line, "\n", "", 1)
			line = strings.Replace(line, "\r", "", 1)
			data := strings.Split(line, " ")
			x, _ := strconv.ParseFloat(data[1], 32)
			y, _ := strconv.ParseFloat(data[2], 32)
			z, _ := strconv.ParseFloat(data[3], 32)
			verts = append(verts, NewVec3(float32(x), float32(y), float32(z)))
		}
	}

	log.Println(verts)

	return MeshData{
		verts: verts,
	}, nil
}
