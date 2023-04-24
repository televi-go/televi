package televi

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/televi-go/televi/models/pages"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Scene interface {
	View(ctx BuildContext)
}
type BuildContext = pages.PageBuildContext
type TransitPolicy = pages.TransitPolicy

func ForEach[T any](context BuildContext, data []T, runner func(element T) pages.View) {
	for _, datum := range data {
		runner(datum).Build(context)
	}
}

func Range(start, end int) []int {
	result := make([]int, end-start)
	for i := start; i < end; i++ {
		result[i-start] = i
	}
	return result
}

type ImageAsset interface {
	Embed(c pages.PhotoConsumer) pages.ImageOptionsSetter
	EmbedAnimation(c pages.AnimationConsumer)
}

type imageAsset struct {
	preloadedBytes []byte
	path           string
}

func (imageAsset imageAsset) EmbedAnimation(c pages.AnimationConsumer) {
	c.Animation(imageAsset.path, bytes.NewReader(imageAsset.preloadedBytes), filepath.Base(imageAsset.path))
}

func (imageAsset imageAsset) Embed(c pages.PhotoConsumer) pages.ImageOptionsSetter {
	return c.Image(imageAsset.path, bytes.NewReader(imageAsset.preloadedBytes))
}

type ImageAssetGroupLoader interface {
	Add(path string, asset *ImageAsset) ImageAssetGroupLoader
	Load() error
}
type imageAssetLoadingEntry struct {
	path  string
	asset *ImageAsset
}

type imageAssetGroupLoader struct {
	tasks []imageAssetLoadingEntry
}

func (imageAssetGroupLoader *imageAssetGroupLoader) Add(path string, asset *ImageAsset) ImageAssetGroupLoader {
	imageAssetGroupLoader.tasks = append(imageAssetGroupLoader.tasks, imageAssetLoadingEntry{
		path:  path,
		asset: asset,
	})
	return imageAssetGroupLoader
}

func (imageAssetGroupLoader *imageAssetGroupLoader) Load() error {
	wg := sync.WaitGroup{}
	errorChannel := make(chan error, len(imageAssetGroupLoader.tasks))
	wg.Add(len(imageAssetGroupLoader.tasks))
	for _, task := range imageAssetGroupLoader.tasks {
		go func(task imageAssetLoadingEntry) {
			defer wg.Done()

			if task.asset == nil {
				errorChannel <- fmt.Errorf("asset specified by path %s cannot be loaded", task.path)
				return
			}

			if imgAsset, isImgAsset := (*task.asset).(imageAsset); isImgAsset {
				if len(imgAsset.preloadedBytes) != 0 {
					errorChannel <- nil
					return
				}
			}

			contents, err := os.ReadFile(task.path)
			*task.asset = imageAsset{
				preloadedBytes: contents,
				path:           task.path,
			}
			errorChannel <- err
		}(task)
	}
	wg.Wait()
	var assetErrorsStrBuilder strings.Builder
	for i := 0; i < len(imageAssetGroupLoader.tasks); i++ {
		err := <-errorChannel
		if err != nil {
			assetErrorsStrBuilder.WriteString(fmt.Sprintf("Error loading asset %s: %v\n", imageAssetGroupLoader.tasks[i].path, err))
		}
	}
	result := assetErrorsStrBuilder.String()
	if result == "" {
		return nil
	}
	return errors.New(result)
}

func NewAssetLoader() ImageAssetGroupLoader {
	return &imageAssetGroupLoader{}
}
