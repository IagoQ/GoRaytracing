package main

import (
	"fmt"
	"io"
	"os"
)

type ObjModel struct {
	Vecs    []Vector
	Normals []Vector
	Faces   [][]Vector
}

func LoadShape(filepath string, mat Material) (IShape, error) {
	model, err := parseFile(filepath)
	if err != nil {
		return nil, err
	}

	triangles := make([]IShape, 0)
	for _, f := range model.Faces {
		t := CreateTriangle(f[0], f[1], f[2], mat)

		triangles = append(triangles, &t)
	}

	fmt.Println("faces", len(triangles))

	return CreateBvh(triangles, 0, len(triangles)), nil
}

func parseFile(filepath string) (ObjModel, error) {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(filepath)
	if err != nil {
		panic("failed to open file")
	}

	// Don't forget to close the file reader.
	defer objFile.Close()

	// Create a model to store stuff.
	model := ObjModel{}

	// Read the file and get it's contents.
	for {
		var lineType string

		// Scan the type field.
		_, err := fmt.Fscanf(objFile, "%s", &lineType)

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Check the type.
		switch lineType {
		// VERTICES.
		case "v":
			// Create a vec to assign digits to.
			vec := Vector{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f%f%f\n", &vec.x, &vec.y, &vec.z)

			// Add the vector to the model.
			model.Vecs = append(model.Vecs, vec.scalarMult(0.1))

		// NORMALS.
		case "vn":
			// Create a vec to assign digits to.
			vec := Vector{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec.x, &vec.y, &vec.z)

			// Add the vector to the model.
			model.Normals = append(model.Normals, vec)

		// TEXTURE VERTICES.
		// case "vt":
		// 	// Create a Uv pair.
		// 	vec := mgl32.Vec2{}
		//
		// 	// Get the digits from the file.
		// 	fmt.Fscanf(objFile, "%f %f\n", &vec[0], &vec[1])
		//
		// 	// Add the uv to the model.
		// 	model.Uvs = append(model.Uvs, vec)

		// INDICES.
		case "f":
			// Create a vec to assign digits to.
			vertices := make([]int, 3)

			// Get the digits from the file.
			matches, _ := fmt.Fscanf(objFile, "%d%d%d\n", &vertices[0], &vertices[1], &vertices[2])

			if matches != 3 {
				fmt.Println("matchs", matches)
				panic("Cannot read your file")
			}

			face := []Vector{model.Vecs[vertices[0]-1], model.Vecs[vertices[1]-1], model.Vecs[vertices[2]-1]}
			model.Faces = append(model.Faces, face)

			// Add the numbers to the model.
		}
	}

	// Return the newly created Model.
	return model, nil
}
