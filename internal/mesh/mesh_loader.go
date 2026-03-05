package mesh

import (
	"bufio"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/alvinobarboza/go-raster/internal/transforms"
)

func LoadVec3(line string, indices int, zNegative bool) (float32, float32, float32) {
	var result [3]float32

	line = strings.Replace(line, "\n", "", 1)
	line = strings.Replace(line, "\r", "", 1)
	data := strings.Split(line, " ")

	for i := range indices {
		d, _ := strconv.ParseFloat(data[i+1], 32)
		result[i] = float32(d)
	}

	if zNegative {
		result[2] = -result[2]
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

func LoadTriangle(line string, shaderSmooth bool) []Triangle {
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
		t := Triangle{ShaderSmooth: shaderSmooth}
		t.V1, t.N1, t.U1 = loadTriangleMetaData(origin)
		t.V2, t.N2, t.U2 = loadTriangleMetaData(tmpData[1])
		t.V3, t.N3, t.U3 = loadTriangleMetaData(tmpData[2])
		tris = append(tris, t)
		offset++
		cursor++
		if cursor > len(data) {
			break
		}
	}

	return tris
}

func LoadMeshFromFile(modelPath string, texturePath string, zNegative, windingReorder, flipNormals bool) (MeshData, error) {
	file, err := os.Open(modelPath)
	if err != nil {
		log.Println(err)
		return MeshData{}, err
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	verts := make([]transforms.Vec3, 0)
	normals := make([]transforms.Vec3, 0)
	uvs := make([]transforms.Vec2, 0)
	tris := make([]Triangle, 0)
	shaderSmooth := false

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
			switch {
			case strings.HasPrefix(line, "v "):
				x, y, z := LoadVec3(line, 3, zNegative)
				verts = append(verts, transforms.NewVec3(x, y, z))
			case strings.HasPrefix(line, "vn"):
				x, y, z := LoadVec3(line, 3, zNegative)
				if flipNormals {
					x = -x
					y = -y
					z = -z
				}
				normals = append(normals, transforms.NewVec3(x, y, z))
			case strings.HasPrefix(line, "vt"):
				x, y, _ := LoadVec3(line, 2, zNegative)
				uvs = append(uvs, transforms.NewVec2(x, y))
			case strings.HasPrefix(line, "f "):
				t := LoadTriangle(line, shaderSmooth)
				tris = append(tris, t...)
			case strings.HasPrefix(line, "s 0") || strings.HasPrefix(line, "s off"):
				shaderSmooth = false
			case strings.HasPrefix(line, "s 1"):
				shaderSmooth = true
			}

		}
	}

	texture, err := LoadTexture(texturePath)
	if err != nil {
		log.Println(err)
		return MeshData{}, err
	}

	// Calculate avarege color for whole object, very hacky to be honest
	var r, g, b int
	for _, p := range texture.pixels {
		r += int(p.R)
		g += int(p.G)
		b += int(p.B)
	}

	r /= len(texture.pixels)
	g /= len(texture.pixels)
	b /= len(texture.pixels)

	var tempV, tempUV, tempN int
	for i := range len(tris) {
		if windingReorder {
			tempV = tris[i].V1
			tempUV = tris[i].U1
			tempN = tris[i].N1
			tris[i].V1 = tris[i].V3
			tris[i].U1 = tris[i].U3
			tris[i].N1 = tris[i].N3
			tris[i].V3 = tempV
			tris[i].U3 = tempUV
			tris[i].N3 = tempN
		}
	}

	return NewMesh(verts, normals, uvs, tris, texture), nil
}

func getPixels(file io.Reader) ([]color.RGBA, int, int, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([]color.RGBA, 0)
	// Upside down, since render is upside down
	for y := range height {
		yu := height - y - 1
		for x := range width {
			r, g, b, a := img.At(x, yu).RGBA()

			// From alpha pre-multiplied values
			// 0xFF00 > 0x00FF > 0xFF
			pixels = append(pixels, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}

	return pixels, width, height, nil
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
		width:      w,
		height:     h,
		fWidth:     float32(w),
		fHeight:    float32(h),
		widthMask:  w - 1,
		heightMask: h - 1,
		pixels:     pixels,
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
