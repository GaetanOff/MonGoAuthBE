package main

import (
	"flag"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// defaultMongoURI ... Mongo uri to use if no uri is passed as a command line arg.
const defaultMongoURI = "mongodb://localhost:27017"

// defaultHTTPPort ... Default port for the api server to run on.
const defaultHTTPPort = 8080

// Args ... Command line arguments for the server.
type Args struct {
	HTTPPort int    `json:"HTTPPort"`
	MongoURI string `json:"MongoURI"`
	Debug    bool   `json:"debug"`
}

func parseCLIArgs() Args {
	debug := flag.Bool("debug", false, "Set to true when developing to enable readable logging.")
	mongoURI := flag.String("mongo-uri", defaultMongoURI, "Set the URI of the mongodatabase this server should use.")
	httpPort := flag.Int("http-port", defaultHTTPPort, "Set the Port number the server should use.")
	flag.Parse()

	return Args{
		HTTPPort: *httpPort,
		MongoURI: *mongoURI,
		Debug:    *debug,
	}
}

func configureLogging(debug bool) *zap.Logger {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "@timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, _ := config.Build()
	return logger
}

func main() {
	// Parse environment variable arguments.
	args := parseCLIArgs()

	// Setup global structured logging...
	logManager := configureLogging(args.Debug)
	zap.ReplaceGlobals(logManager)
	// Make sure everthing gets logged before exit.
	defer logManager.Sync()

	// Application start print.
	zap.S().Infow("Login Servier Application Start", zap.Any("args", args))

	// Create and start the server to run forever.
	server := NewLoginServer(args.HTTPPort, args.MongoURI)
	zap.L().Fatal("Login Server Application Exit", zap.Error(server.Start()))
}
