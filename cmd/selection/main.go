package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/selection/api/protobuf-spec/selectionpb"
	"github.com/jukeizu/selection/internal/startup"
	"github.com/jukeizu/selection/selection"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/shawntoffel/gossage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var Version = ""

var (
	flagMigrate = false
	flagVersion = false
	flagServer  = true

	grpcPort       = "50052"
	dbAddress      = "root@localhost:26257"
	serviceAddress = "localhost:" + grpcPort
)

func parseConfig() {
	flag.StringVar(&grpcPort, "grpc.port", grpcPort, "grpc port for server")
	flag.StringVar(&dbAddress, "db", dbAddress, "Database connection address")
	flag.StringVar(&serviceAddress, "service.addr", serviceAddress, "address of service if not local")
	flag.BoolVar(&flagServer, "server", flagServer, "Run as server")
	flag.BoolVar(&flagMigrate, "migrate", flagMigrate, "Run db migrations")
	flag.BoolVar(&flagVersion, "v", flagVersion, "version")

	flag.Parse()
}

func main() {
	parseConfig()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("instance", xid.New().String()).
		Str("component", "selection").
		Str("version", Version).
		Logger()

	grpcLoggerV2 := grpczerolog.New(logger.With().Str("transport", "grpc").Logger())
	grpclog.SetLoggerV2(grpcLoggerV2)

	repository, err := selection.NewRepository(dbAddress)
	if err != nil {
		logger.Error().Err(err).Caller().Msg("could not create selection repository")
		os.Exit(1)
	}

	if flagMigrate {
		gossage.Logger = func(format string, a ...interface{}) {
			msg := fmt.Sprintf(format, a...)
			logger.Info().Str("component", "migrator").Msg(msg)
		}

		err := repository.Migrate()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("could not migrate selection repository")
			os.Exit(1)
		}
	}

	g := run.Group{}

	if flagServer {
		grpcServer := newGrpcServer(logger)
		server := startup.NewServer(logger, grpcServer)

		sorter := selection.NewSorter(logger)
		batcher := selection.NewBatcher(logger)

		selectionService := selection.NewDefaultService(logger, repository, sorter, batcher)
		selectionServer := selection.NewGrpcServer(selectionService)
		selectionpb.RegisterSelectionServer(grpcServer, selectionServer)

		grpcAddr := ":" + grpcPort

		g.Add(func() error {
			return server.Start(grpcAddr)
		}, func(error) {
			server.Stop()
		})
	}

	cancel := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancel)
	}, func(error) {
		close(cancel)
	})

	logger.Info().Err(g.Run()).Msg("stopped")
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-cancel:
		return errors.New("stopping")
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	}
}

func newGrpcServer(logger zerolog.Logger) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    5 * time.Minute,
				Timeout: 10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
		startup.LoggingInterceptor(logger),
	)

	return grpcServer
}
