``` 
aydin@pop-os:~/Belgeler/go/chat-moderation-service$ go run cmd/main.go 
{"level":"INFO","ts":"2025-06-18T04:31:49.519+0300","caller":"app/badwords.go:51","msg":"BadWords listesi yüklendi","count":638}
{"level":"INFO","ts":"2025-06-18T04:31:49.519+0300","caller":"app/moderation.go:88","msg":"Moderasyon izleme başladı"}
[KÜFÜR ALGILANDI] ayd1ndem1irci: "amcık"
{"level":"WARN","ts":"2025-06-18T04:32:28.689+0300","caller":"app/moderation.go:136","msg":"Küfür algılandı","player_name":"ayd1ndem1irci","message":"amcık"}
[KÜFÜR ALGILANDI] ayd1ndem1irci: "amcık mısın"
{"level":"WARN","ts":"2025-06-18T04:32:44.751+0300","caller":"app/moderation.go:136","msg":"Küfür algılandı","player_name":"ayd1ndem1irci","message":"amcık mısın"}
[KÜFÜR ALGILANDI] ayd1ndem1irci: "amcık-mısın"
{"level":"WARN","ts":"2025-06-18T04:32:49.772+0300","caller":"app/moderation.go:136","msg":"Küfür algılandı","player_name":"ayd1ndem1irci","message":"amcık-mısın"}
[KÜFÜR ALGILANDI] ayd1ndem1irci: "amcık mısın"
{"level":"WARN","ts":"2025-06-18T04:32:55.798+0300","caller":"app/moderation.go:136","msg":"Küfür algılandı","player_name":"ayd1ndem1irci","message":"amcık mısın"}
[KÜFÜR ALGILANDI] ayd1ndem1irci: "senin ananı sikyim"
{"level":"WARN","ts":"2025-06-18T04:33:05.834+0300","caller":"app/moderation.go:136","msg":"Küfür algılandı","player_name":"ayd1ndem1irci","message":"senin ananı sikyim"}
```

```yaml
mongo_uri: "mongodb://localhost:27017"
database_name: "foxlogger"
collection_name: "logs"
badwords_file: "badwords.txt"

allowed_players:
  - ayd1ndemirci

```
