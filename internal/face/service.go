package face

import (
	"bytes"
	"encoding/base64"
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/Kagami/go-face"
)

type Service interface {
	Train(imageData []byte, personID string) error
	Compare(imageData []byte) (string, float32, error)
}

type service struct {
	rec *face.Recognizer
}

func NewService() (*service, error) {
	rec, err := face.NewRecognizer("./models")
	if err != nil {
		return nil, err
	}
	return &service{rec: rec}, nil
}

func (s *service) Train(imageData []byte, personID string) error {
	img, err := decodeImage(imageData)
	if err != nil {
		return err
	}

	faces, err := s.rec.RecognizeFile(img)
	if err != nil {
		return err
	}

	if len(faces) == 0 {
		return errors.New("no face detected in the image")
	}

	// Store the face descriptor
	s.rec.SetSampleForPerson(personID, faces[0].Descriptor)
	return nil
}

func (s *service) Compare(imageData []byte) (string, float32, error) {
	img, err := decodeImage(imageData)
	if err != nil {
		return "", 0, err
	}

	faces, err := s.rec.RecognizeFile(img)
	if err != nil {
		return "", 0, err
	}

	if len(faces) == 0 {
		return "", 0, errors.New("no face detected in the image")
	}

	// Find the closest match
	samples := s.rec.Samples
	if len(samples) == 0 {
		return "", 0, errors.New("no trained faces available")
	}

	bestMatch := ""
	bestDistance := float32(0.6) // Threshold for face similarity
	currentDistance := float32(1.0)

	for id, sample := range samples {
		distance := face.SquaredEuclideanDistance(faces[0].Descriptor, sample.Descriptor)
		if distance < currentDistance {
			currentDistance = distance
			bestMatch = id
		}
	}

	if currentDistance > bestDistance {
		return "", currentDistance, errors.New("no match found within threshold")
	}

	return bestMatch, currentDistance, nil
}

func decodeImage(data []byte) (string, error) {
	// Try to decode base64
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err == nil {
		data = decoded
	}

	// Create a temporary file for the image
	tmpfile, err := os.CreateTemp("", "face-*.jpg")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := io.Copy(tmpfile, bytes.NewReader(data)); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
