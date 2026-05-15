package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"math"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	_ "image/jpeg"
	_ "image/png"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/moodora/moodora/apps/api/internal/platform/config"
	"github.com/moodora/moodora/apps/api/internal/platform/database"
	"github.com/moodora/moodora/apps/api/internal/platform/storage"
)

const deckCode = "rider_waite"

type variant struct {
	format string
	size   string
	width  int
}

var variants = []variant{
	{format: "webp", size: "thumb", width: 180},
	{format: "webp", size: "medium", width: 480},
	{format: "webp", size: "large", width: 720},
	{format: "jpg", size: "medium", width: 480},
}

type appConfig struct {
	databaseURL string
	storage     config.StorageConfig
}

type processedAsset struct {
	sourceCode string
	cardID     string
	format     string
	size       string
	width      int
	height     int
	fileSize   int64
	localPath  string
	objectKey  string
	url        string
	isDefault  bool
}

type sourceImages struct {
	files   map[string]string
	missing []string
	unknown []string
}

type pipelineSummary struct {
	sourceImages   int
	processedCount int
	uploadedCount  int
	upsertedCount  int
	missingImages  []string
	unknownFiles   []string
}

func main() {
	inputDir := flag.String("input", "../../local-assets/tarot", "source image directory")
	outputDir := flag.String("output", "processed-assets/tarot/rider_waite", "processed image output directory")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := loadConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	if _, err := exec.LookPath("cwebp"); err != nil {
		logger.Error("cwebp is required for WebP generation", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.databaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	storageClient, err := storage.Connect(cfg.storage)
	if err != nil {
		logger.Error("failed to connect to object storage", "error", err)
		os.Exit(1)
	}

	cardIDs, err := loadTarotCards(ctx, db)
	if err != nil {
		logger.Error("failed to load tarot cards", "error", err)
		os.Exit(1)
	}

	sourceImages, err := findSourceImages(*inputDir, cardIDs)
	if err != nil {
		logger.Error("failed to validate source images", "error", err)
		os.Exit(1)
	}
	if len(sourceImages.files) == 0 {
		logger.Error("no source images found", "input", *inputDir)
		os.Exit(1)
	}

	bucketReady, err := storageClient.BucketExists(ctx, cfg.storage.Bucket)
	if err != nil {
		logger.Error("failed to check storage bucket", "error", err)
		os.Exit(1)
	}
	if !bucketReady {
		logger.Error("storage bucket does not exist", "bucket", cfg.storage.Bucket)
		os.Exit(1)
	}

	summary := pipelineSummary{
		sourceImages:  len(sourceImages.files),
		missingImages: sourceImages.missing,
		unknownFiles:  sourceImages.unknown,
	}
	for _, sourceCode := range sortedKeys(sourceImages.files) {
		sourcePath := sourceImages.files[sourceCode]
		assets, err := processSource(ctx, sourceCode, cardIDs[sourceCode], sourcePath, *outputDir, cfg.storage, storageClient)
		if err != nil {
			logger.Error("failed to process source image", "sourceCode", sourceCode, "error", err)
			os.Exit(1)
		}
		upserted, err := upsertAssets(ctx, db, assets)
		if err != nil {
			logger.Error("failed to upsert tarot card assets", "sourceCode", sourceCode, "error", err)
			os.Exit(1)
		}
		summary.processedCount += len(assets)
		summary.uploadedCount += len(assets)
		summary.upsertedCount += upserted
	}

	printSummary(summary)
	logger.Info(
		"tarot asset pipeline completed",
		"sourceImages", summary.sourceImages,
		"processed", summary.processedCount,
		"uploaded", summary.uploadedCount,
		"upserted", summary.upsertedCount,
		"missingImages", len(summary.missingImages),
		"unknownFiles", len(summary.unknownFiles),
	)
}

func loadConfig() (appConfig, error) {
	cfg := appConfig{
		databaseURL: os.Getenv("DATABASE_URL"),
		storage: config.StorageConfig{
			Endpoint:      getEnv("S3_ENDPOINT", "http://localhost:9000"),
			Region:        getEnv("S3_REGION", "auto"),
			Bucket:        os.Getenv("S3_BUCKET"),
			AccessKey:     os.Getenv("S3_ACCESS_KEY"),
			SecretKey:     os.Getenv("S3_SECRET_KEY"),
			PublicBaseURL: strings.TrimRight(os.Getenv("S3_PUBLIC_BASE_URL"), "/"),
		},
	}
	if cfg.databaseURL == "" {
		return appConfig{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.storage.Bucket == "" {
		return appConfig{}, fmt.Errorf("S3_BUCKET is required")
	}
	if cfg.storage.AccessKey == "" {
		return appConfig{}, fmt.Errorf("S3_ACCESS_KEY is required")
	}
	if cfg.storage.SecretKey == "" {
		return appConfig{}, fmt.Errorf("S3_SECRET_KEY is required")
	}
	if cfg.storage.PublicBaseURL == "" {
		return appConfig{}, fmt.Errorf("S3_PUBLIC_BASE_URL is required")
	}
	return cfg, nil
}

func loadTarotCards(ctx context.Context, db *pgxpool.Pool) (map[string]string, error) {
	rows, err := db.Query(ctx, `
		SELECT source_code, id::text
		FROM tarot_cards
		ORDER BY source_code
	`)
	if err != nil {
		return nil, fmt.Errorf("query tarot cards: %w", err)
	}
	defer rows.Close()

	cards := map[string]string{}
	for rows.Next() {
		var sourceCode string
		var id string
		if err := rows.Scan(&sourceCode, &id); err != nil {
			return nil, err
		}
		cards[sourceCode] = id
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("tarot_cards is empty; run the tarot card importer first")
	}
	return cards, nil
}

func findSourceImages(inputDir string, cardIDs map[string]string) (sourceImages, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return sourceImages{}, fmt.Errorf("read input directory %q: %w", inputDir, err)
	}

	sourceFiles := map[string]string{}
	var unknown []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		sourceCode := sourceCodeFromImageName(entry.Name())
		if _, ok := cardIDs[sourceCode]; !ok {
			unknown = append(unknown, entry.Name())
			continue
		}
		sourceFiles[sourceCode] = filepath.Join(inputDir, entry.Name())
	}

	var missing []string
	for sourceCode := range cardIDs {
		if _, ok := sourceFiles[sourceCode]; !ok {
			missing = append(missing, sourceCode)
		}
	}
	sort.Strings(missing)
	sort.Strings(unknown)
	return sourceImages{files: sourceFiles, missing: missing, unknown: unknown}, nil
}

func sourceCodeFromImageName(name string) string {
	baseName := strings.TrimSuffix(name, filepath.Ext(name))
	if sourceCode, ok := riderWaiteSmithSourceCode(baseName); ok {
		return sourceCode
	}
	return baseName
}

func riderWaiteSmithSourceCode(baseName string) (string, bool) {
	parts := strings.Split(baseName, "-")
	if len(parts) != 3 || parts[0] != "RWSa" {
		return "", false
	}

	switch parts[1] {
	case "T":
		return "ar" + parts[2], true
	case "W":
		return "wa" + riderWaiteSmithRank(parts[2]), true
	case "C":
		return "cu" + riderWaiteSmithRank(parts[2]), true
	case "S":
		return "sw" + riderWaiteSmithRank(parts[2]), true
	case "P":
		return "pe" + riderWaiteSmithRank(parts[2]), true
	default:
		return "", false
	}
}

func riderWaiteSmithRank(rank string) string {
	switch rank {
	case "0A":
		return "01"
	case "J1":
		return "11"
	case "J2":
		return "12"
	case "QU":
		return "13"
	case "KI":
		return "14"
	default:
		return rank
	}
}

func processSource(ctx context.Context, sourceCode string, cardID string, sourcePath string, outputDir string, cfg config.StorageConfig, storageClient *minio.Client) ([]processedAsset, error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer source.Close()

	img, _, err := image.Decode(source)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	assets := make([]processedAsset, 0, len(variants))
	for _, variant := range variants {
		resized := resizeToWidth(img, variant.width)
		localPath := filepath.Join(outputDir, variant.format, variant.size, sourceCode+"."+variant.format)
		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			return nil, err
		}

		switch variant.format {
		case "jpg":
			if err := writeJPEG(localPath, resized); err != nil {
				return nil, err
			}
		case "webp":
			if err := writeWebP(localPath, resized); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported output format %q", variant.format)
		}

		info, err := os.Stat(localPath)
		if err != nil {
			return nil, err
		}
		objectKey := fmt.Sprintf("tarot/%s/%s/%s/%s.%s", deckCode, variant.format, variant.size, sourceCode, variant.format)
		contentType := mime.TypeByExtension("." + variant.format)
		if variant.format == "jpg" {
			contentType = "image/jpeg"
		}
		if variant.format == "webp" {
			contentType = "image/webp"
		}

		file, err := os.Open(localPath)
		if err != nil {
			return nil, err
		}
		_, putErr := storageClient.PutObject(ctx, cfg.Bucket, objectKey, file, info.Size(), minio.PutObjectOptions{ContentType: contentType})
		closeErr := file.Close()
		if putErr != nil {
			return nil, fmt.Errorf("upload %s: %w", objectKey, putErr)
		}
		if closeErr != nil {
			return nil, closeErr
		}

		bounds := resized.Bounds()
		assets = append(assets, processedAsset{
			sourceCode: sourceCode,
			cardID:     cardID,
			format:     variant.format,
			size:       variant.size,
			width:      bounds.Dx(),
			height:     bounds.Dy(),
			fileSize:   info.Size(),
			localPath:  localPath,
			objectKey:  objectKey,
			url:        cfg.PublicBaseURL + "/" + objectKey,
			isDefault:  variant.format == "webp" && variant.size == "medium",
		})
	}
	return assets, nil
}

func resizeToWidth(src image.Image, width int) image.Image {
	bounds := src.Bounds()
	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()
	if sourceWidth <= 0 || sourceHeight <= 0 {
		return src
	}
	height := int(math.Round(float64(sourceHeight) * float64(width) / float64(sourceWidth)))
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		sourceY := bounds.Min.Y + int(float64(y)*float64(sourceHeight)/float64(height))
		for x := 0; x < width; x++ {
			sourceX := bounds.Min.X + int(float64(x)*float64(sourceWidth)/float64(width))
			dst.Set(x, y, src.At(sourceX, sourceY))
		}
	}
	return dst
}

func writeJPEG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, img, &jpeg.Options{Quality: 88})
}

func writeWebP(path string, img image.Image) error {
	var pngBuffer bytes.Buffer
	if err := png.Encode(&pngBuffer, img); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "moodora-tarot-*.png")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if _, err := tempFile.Write(pngBuffer.Bytes()); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	cmd := exec.Command("cwebp", "-quiet", "-q", "82", tempPath, "-o", path)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func upsertAssets(ctx context.Context, db *pgxpool.Pool, assets []processedAsset) (int, error) {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin upsert transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var upserted int
	for _, asset := range assets {
		if asset.isDefault {
			_, err := tx.Exec(ctx, `
				UPDATE tarot_card_assets
				SET is_default = false,
					updated_at = now()
				WHERE card_id = $1
					AND deck_code = $2
					AND is_default = true
					AND NOT (size = $3 AND format = $4)
			`, asset.cardID, deckCode, asset.size, asset.format)
			if err != nil {
				return 0, fmt.Errorf("clear existing default for %s: %w", asset.sourceCode, err)
			}
		}

		tag, err := tx.Exec(ctx, `
			INSERT INTO tarot_card_assets (
				card_id, deck_code, size, format, url, width, height, file_size, is_default
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (card_id, deck_code, size, format) DO UPDATE
			SET url = EXCLUDED.url,
				width = EXCLUDED.width,
				height = EXCLUDED.height,
				file_size = EXCLUDED.file_size,
				is_default = EXCLUDED.is_default,
				updated_at = now()
		`, asset.cardID, deckCode, asset.size, asset.format, asset.url, asset.width, asset.height, asset.fileSize, asset.isDefault)
		if err != nil {
			return 0, fmt.Errorf("upsert %s %s %s: %w", asset.sourceCode, asset.format, asset.size, err)
		}
		upserted += int(tag.RowsAffected())
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit upsert transaction: %w", err)
	}
	return upserted, nil
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func printSummary(summary pipelineSummary) {
	fmt.Println("Tarot asset pipeline summary")
	fmt.Printf("processed count: %d\n", summary.processedCount)
	fmt.Printf("uploaded count: %d\n", summary.uploadedCount)
	fmt.Printf("upserted count: %d\n", summary.upsertedCount)
	fmt.Printf("missing images: %d\n", len(summary.missingImages))
	if len(summary.missingImages) > 0 {
		fmt.Printf("missing image source_codes: %s\n", strings.Join(summary.missingImages, ", "))
	}
	fmt.Printf("unknown files: %d\n", len(summary.unknownFiles))
	if len(summary.unknownFiles) > 0 {
		fmt.Printf("unknown file names: %s\n", strings.Join(summary.unknownFiles, ", "))
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
