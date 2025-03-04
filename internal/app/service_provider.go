package app

import (
	"context"
	"log"

	"auth/internal/config"
	"auth/internal/jwt"
	"auth/pkg/client/db"
	"auth/pkg/client/db/pg"
	"auth/pkg/client/db/transaction"
	"auth/pkg/closer"

	authApi "auth/internal/api/auth"
	authService "auth/internal/service/auth"

	userRepository "auth/internal/repository/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	jwtConfig  config.JWTConfig

	dbClient  db.Client
	txManager db.TxManager

	authRepository     authService.IAuthRepository
	jwt                authService.IJWT
	authService        authApi.IAuthService
	authImplementation *authApi.AuthImplementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) authService.IAuthRepository {
	if s.authRepository == nil {
		s.authRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) JWT(ctx context.Context) authService.IJWT {
	if s.jwt == nil {
		s.jwt = jwt.NewJSONWebToken(s.JWTConfig().GetSecretKey())
	}

	return s.jwt
}

func (s *serviceProvider) AuthService(ctx context.Context) authApi.IAuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx), s.JWT(ctx), s.TxManager(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthImplementation(ctx context.Context) *authApi.AuthImplementation {
	if s.authImplementation == nil {
		s.authImplementation = authApi.NewAuthImplementation(s.AuthService(ctx))
	}

	return s.authImplementation
}
