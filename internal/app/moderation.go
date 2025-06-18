package app

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Moderator struct {
	client         *mongo.Client
	collection     *mongo.Collection
	badWordManager *BadWordManager
	lastChecked    time.Time
	ctx            context.Context
	allowedPlayers map[string]struct{}
	logger         *zap.Logger

	mu        sync.Mutex
	csvFile   *os.File
	csvWriter *csv.Writer
}

func NewModerator(client *mongo.Client, dbName, collName string, bwm *BadWordManager, logger *zap.Logger, allowed []string) (*Moderator, error) {
	if logger == nil {
		defaultLogger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("varsayılan Zap logger oluşturulamadı: %w", err)
		}
		logger = defaultLogger
	}

	allowedMap := make(map[string]struct{}, len(allowed))
	for _, p := range allowed {
		allowedMap[p] = struct{}{}
	}

	f, err := os.OpenFile("badwords_log.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(f)

	fi, err := f.Stat()
	if err == nil && fi.Size() == 0 {
		writer.Write([]string{"PlayerName", "Timestamp", "Message"})
		writer.Flush()
	}

	return &Moderator{
		client:         client,
		collection:     client.Database(dbName).Collection(collName),
		badWordManager: bwm,
		lastChecked:    time.Now(),
		ctx:            context.Background(),
		allowedPlayers: allowedMap,
		logger:         logger,
		csvFile:        f,
		csvWriter:      writer,
	}, nil
}

func (m *Moderator) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.csvWriter.Flush()
	m.csvFile.Close()
}

func (m *Moderator) logBadWord(playerName string, timestamp time.Time, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	err := m.csvWriter.Write([]string{playerName, timestamp.Format(time.RFC3339), message})
	if err != nil {
		m.logger.Error("CSV yazma hatası", zap.Error(err)) 
	}
	m.csvWriter.Flush()
}

func (m *Moderator) StartMonitoring() {
	m.logger.Info("Moderasyon izleme başladı") 
	processedMessages := make(map[string]struct{})

	for {
		filter := bson.M{
			"event_type": "player_chat",
			"timestamp":  bson.M{"$gt": m.lastChecked},
		}

		cur, err := m.collection.Find(m.ctx, filter)
		if err != nil {
			m.logger.Error("MongoDB find hatası", zap.Error(err)) 
			time.Sleep(time.Second)
			continue
		}

		var results []struct {
			Player struct {
				Name string `bson:"name"`
			} `bson:"player"`
			Details struct {
				Message string `bson:"message"`
			} `bson:"details"`
			Timestamp time.Time `bson:"timestamp"`
		}

		if err := cur.All(m.ctx, &results); err != nil {
			m.logger.Error("MongoDB cursor decode hatası", zap.Error(err)) 
			time.Sleep(time.Second)
			continue
		}

		maxTimestamp := m.lastChecked

		for _, doc := range results {
			msgID := doc.Player.Name + doc.Timestamp.String() + doc.Details.Message
			if _, exists := processedMessages[msgID]; exists {
				continue
			}
			processedMessages[msgID] = struct{}{}

			if _, ok := m.allowedPlayers[doc.Player.Name]; ok {
				m.logger.Info("İzin verilen oyuncu (küfür serbest)", zap.String("player_name", doc.Player.Name), zap.String("message", doc.Details.Message)) 
				continue
			}

			if ContainsBadWord(doc.Details.Message, m.badWordManager.GetWords()) {
				fmt.Printf("[KÜFÜR ALGILANDI] %s: \"%s\"\n", doc.Player.Name, doc.Details.Message)
				m.logger.Warn("Küfür algılandı", zap.String("player_name", doc.Player.Name), zap.String("message", doc.Details.Message)) 
				m.logBadWord(doc.Player.Name, doc.Timestamp, doc.Details.Message)
			}

			if doc.Timestamp.After(maxTimestamp) {
				maxTimestamp = doc.Timestamp
			}
		}

		m.lastChecked = maxTimestamp

		time.Sleep(time.Second)
	}
}
