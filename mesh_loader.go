package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadVec3(line string, indices int) (float32, float32, float32) {
	var result [3]float32

	line = strings.Replace(line, "\n", "", 1)
	line = strings.Replace(line, "\r", "", 1)
	data := strings.Split(line, " ")

	for i := range indices {
		d, _ := strconv.ParseFloat(data[i+1], 32)
		result[i] = float32(d)
	}

	return result[0], result[1], result[2]
}

func loadTriangleMetaData(indexes string) (int, int, int) {
	var vertex, normal, uv int

	idxs := strings.Split(indexes, "/")
	switch len(idxs) {
	case 1:
		vertex, _ = strconv.Atoi(idxs[0])
	case 2:
		vertex, _ = strconv.Atoi(idxs[0])
		uv, _ = strconv.Atoi(idxs[1])
	case 3:
		vertex, _ = strconv.Atoi(idxs[0])
		uv, _ = strconv.Atoi(idxs[1])
		normal, _ = strconv.Atoi(idxs[2])
	}

	return vertex - 1, normal - 1, uv - 1
}

func LoadTriangle(line string) []Triangle {
	tris := make([]Triangle, 0)

	line = strings.Replace(line, "\n", "", 1)
	line = strings.Replace(line, "\r", "", 1)

	data := strings.Split(line, " ")
	data = data[1:]

	cursor := 3
	offset := 0
	origin := data[0]
	for {
		tmpData := data[offset:cursor]
		t := Triangle{}
		t.v1, t.n1, t.u1 = loadTriangleMetaData(origin)
		t.v2, t.n2, t.u2 = loadTriangleMetaData(tmpData[1])
		t.v3, t.n3, t.u3 = loadTriangleMetaData(tmpData[2])
		tris = append(tris, t)
		offset++
		cursor++
		if cursor > len(data) {
			break
		}
	}

	return tris
}

func LoadMeshFromFile(modelPath string, texturePath string) (MeshData, error) {
	file, err := os.Open(modelPath)
	if err != nil {
		log.Println(err)
		return MeshData{}, err
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	verts := make([]Vec3, 0)
	normals := make([]Vec3, 0)
	uvs := make([]Vec3, 0)
	tris := make([]Triangle, 0)

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
		if len(line) > 2 {
			switch line[:2] {
			case "v ":
				x, y, z := LoadVec3(line, 3)
				verts = append(verts, NewVec3(x, y, z))
			case "vn":
				x, y, z := LoadVec3(line, 3)
				normals = append(normals, NewVec3(x, y, z))
			case "vt":
				x, y, _ := LoadVec3(line, 2)
				uvs = append(uvs, NewVec3(x, y, 0))
			case "f ":
				t := LoadTriangle(line)
				tris = append(tris, t...)
			}
		}
	}


func LoadTexture(path string) (*Texture, error) {
	file, err := os.Open(path)

	if err != nil {
		log.Println("Error: File could not be opened")
		return nil, err
	}

	defer file.Close()

	pixels, w, h, err := getPixels(file)

	if err != nil {
		log.Println("Error: Image could not be decoded")
		return nil, err
	}

	return &Texture{
		width:  w,
		height: h,
		pixels: pixels,
	}, nil
}

// Will panic if default not present
func LoadDefaultTexture() *Texture {
	img, err := LoadTexture("./assets/default.jpg")
	if err != nil {
		panic(err)
	}
	return img
}
