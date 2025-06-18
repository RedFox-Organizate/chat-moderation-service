package app

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
)

type BadWordManager struct {
	words  []string
	mu     sync.RWMutex
	file   string
	logger *zap.Logger
}

func NewBadWordManager(file string, logger *zap.Logger) (*BadWordManager, error) {
	bwm := &BadWordManager{
		file:   file,
		logger: logger,
	}
	err := bwm.load()
	if err != nil {
		return nil, err
	}
	return bwm, nil
}

func (b *BadWordManager) load() error {
	file, err := os.Open(b.file)
	if err != nil {
		return err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}
	b.mu.Lock()
	b.words = words
	b.mu.Unlock()
	b.logger.Info("BadWords listesi yüklendi", zap.Int("count", len(words)))
	return scanner.Err()
}

func (b *BadWordManager) WatchFileChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		b.logger.Error("fsnotify watcher başlatılamadı", zap.Error(err))
		return
	}
	defer watcher.Close()

	err = watcher.Add(b.file)
	if err != nil {
		b.logger.Error("Dosya izlenemedi", zap.String("file", b.file), zap.Error(err))
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				b.logger.Info("BadWords dosyası değişti, yeniden yüklüyor...")
				err := b.load()
				if err != nil {
					b.logger.Error("BadWords dosyası yüklenirken hata", zap.Error(err))
				}
			}
		case err := <-watcher.Errors:
			b.logger.Error("Watcher hatası", zap.Error(err))
		}
	}
}

func (b *BadWordManager) GetWords() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]string(nil), b.words...)
}
