package govatar

import (
	"errors"
	"hash/fnv"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var errUnknownGender = errors.New("Unknown gender")

type person struct {
	Clothes []string
	Eye     []string
	Face    []string
	Hair    []string
	Mouth   []string
}

type store struct {
	Background []string
	Male       person
	Female     person
	Monster    person
}

var assetsStore *store

// Gender represents gender type
type Gender int

// Male and female constants
const (
	MALE Gender = iota
	FEMALE
	MONSTER
)

func init() {
	male := getPerson(MALE)
	female := getPerson(FEMALE)
	monster := getPerson(MONSTER)
	assetsStore = &store{Background: readAssetsFrom("data/background"), Male: male, Female: female, Monster: monster}
	rand.Seed(time.Now().UTC().UnixNano())
}

// Generate generates random avatar
func Generate(gender Gender) (image.Image, error) {
	switch gender {
	case MALE:
		return randomAvatar(assetsStore.Male, time.Now().UnixNano())
	case FEMALE:
		return randomAvatar(assetsStore.Female, time.Now().UnixNano())
	case MONSTER:
		return randomAvatar(assetsStore.Monster, time.Now().UnixNano())
	default:
		return nil, errUnknownGender
	}
}

// GenerateFile generates random avatar and save it to specified file.
// Image format depends on file extension (jpeg, jpg, png, gif). Default is png
func GenerateFile(gender Gender, filePath string) error {
	img, err := Generate(gender)
	if err != nil {
		return err
	}
	return saveToFile(img, filePath)
}

// GenerateFromUsername generates avatar from string
func GenerateFromUsername(gender Gender, username string) (image.Image, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(username))
	if err != nil {
		return nil, err
	}
	switch gender {
	case MALE:
		return randomAvatar(assetsStore.Male, int64(h.Sum32()))
	case FEMALE:
		return randomAvatar(assetsStore.Female, int64(h.Sum32()))
	case MONSTER:
		return randomAvatar(assetsStore.Monster, int64(h.Sum32()))
	default:
		return nil, errUnknownGender
	}
}

// GenerateFileFromUsername generates avatar from string and save it to specified file.
// Image format depends on file extension (jpeg, jpg, png, gif). Default is png
func GenerateFileFromUsername(gender Gender, username string, filePath string) error {
	img, err := GenerateFromUsername(gender, username)
	if err != nil {
		return err
	}
	return saveToFile(img, filePath)
}

func saveToFile(img image.Image, filePath string) error {
	outFile, err := os.Create(filePath)
	defer outFile.Close()
	if err != nil {
		return err
	}
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".jpeg", ".jpg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
	case ".gif":
		err = gif.Encode(outFile, img, nil)
	default:
		err = png.Encode(outFile, img)
	}
	return err
}

func randomAvatar(p person, seed int64) (image.Image, error) {
	rnd := rand.New(rand.NewSource(seed))
	avatar := image.NewRGBA(image.Rect(0, 0, 400, 400))
	var err error
	err = drawImg(avatar, randSliceString(rnd, assetsStore.Background), err)
	err = drawImg(avatar, randSliceString(rnd, p.Face), err)
	err = drawImg(avatar, randSliceString(rnd, p.Clothes), err)
	err = drawImg(avatar, randSliceString(rnd, p.Mouth), err)
	err = drawImg(avatar, randSliceString(rnd, p.Hair), err)
	err = drawImg(avatar, randSliceString(rnd, p.Eye), err)
	return avatar, err
}

func drawImg(dst draw.Image, asset string, err error) error {
	if err != nil {
		return err
	}
	infile, err := os.Open(asset)
	if err != nil {
		// replace this with real error handling
		panic(err)
	}
	defer infile.Close()
	src, _, err := image.Decode(infile) //bindata.MustAsset(asset)))
	if err != nil {
		return err
	}
	draw.Draw(dst, dst.Bounds(), src, image.Point{0, 0}, draw.Over)
	return nil
}

func getPerson(gender Gender) person {
	var genderPath string

	switch gender {
	case FEMALE:
		genderPath = "female"
	case MALE:
		genderPath = "male"
	case MONSTER:
		genderPath = "monster"
	}

	return person{
		Clothes: readAssetsFrom("data/" + genderPath + "/clothes"),
		Eye:     readAssetsFrom("data/" + genderPath + "/eye"),
		Face:    readAssetsFrom("data/" + genderPath + "/face"),
		Hair:    readAssetsFrom("data/" + genderPath + "/hair"),
		Mouth:   readAssetsFrom("data/" + genderPath + "/mouth"),
	}
}

func readAssetsFrom(dir string) (assets []string) {

	files, err := ioutil.ReadDir("./" + dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, asset := range files {
		if asset.Name() == ".DS_Store" {
			continue
		}

		assets = append(assets, filepath.Join(dir, asset.Name()))
	}
	sort.Sort(naturalSort(assets))
	return assets
}
