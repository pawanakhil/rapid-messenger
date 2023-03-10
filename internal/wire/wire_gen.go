// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/pawanakhil/rapid-messenger/pkg/chat"
	"github.com/pawanakhil/rapid-messenger/pkg/common"
	"github.com/pawanakhil/rapid-messenger/pkg/config"
	"github.com/pawanakhil/rapid-messenger/pkg/infra"
	"github.com/pawanakhil/rapid-messenger/pkg/match"
	"github.com/pawanakhil/rapid-messenger/pkg/uploader"
	"github.com/pawanakhil/rapid-messenger/pkg/user"
	"github.com/pawanakhil/rapid-messenger/pkg/web"
)

// Injectors from wire.go:

func InitializeWebServer(name string) (*common.Server, error) {
	httpLogrus := common.NewHttpLogrus()
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	engine := web.NewGinServer(name, httpLogrus)
	httpServer := web.NewHttpServer(name, httpLogrus, configConfig, engine)
	router := web.NewRouter(httpServer)
	infraCloser := web.NewInfraCloser()
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, router, infraCloser, observabilityInjector)
	return server, nil
}

func InitializeMatchServer(name string) (*common.Server, error) {
	httpLogrus := common.NewHttpLogrus()
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	engine := match.NewGinServer(name, httpLogrus, configConfig)
	melodyMatchConn := match.NewMelodyMatchConn()
	universalClient, err := infra.NewRedisClient(configConfig)
	if err != nil {
		return nil, err
	}
	userClientConn, err := match.NewUserClientConn(configConfig)
	if err != nil {
		return nil, err
	}
	chatClientConn, err := match.NewChatClientConn(configConfig)
	if err != nil {
		return nil, err
	}
	userRepo := match.NewUserRepo(userClientConn, chatClientConn)
	userService := match.NewUserService(userRepo)
	matchSubscriber := match.NewMatchSubscriber(configConfig, universalClient, melodyMatchConn, userService)
	redisCache := infra.NewRedisCache(universalClient)
	matchingRepo := match.NewMatchingRepo(redisCache)
	channelRepo := match.NewChannelRepo(chatClientConn)
	matchingService := match.NewMatchingService(matchingRepo, channelRepo)
	httpServer := match.NewHttpServer(name, httpLogrus, configConfig, engine, melodyMatchConn, matchSubscriber, userService, matchingService)
	router := match.NewRouter(httpServer)
	infraCloser := match.NewInfraCloser()
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, router, infraCloser, observabilityInjector)
	return server, nil
}

func InitializeChatServer(name string) (*common.Server, error) {
	httpLogrus := common.NewHttpLogrus()
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	engine := chat.NewGinServer(name, httpLogrus, configConfig)
	melodyChatConn := chat.NewMelodyChatConn(configConfig)
	universalClient, err := infra.NewRedisClient(configConfig)
	if err != nil {
		return nil, err
	}
	messageSubscriber := chat.NewMessageSubscriber(configConfig, universalClient, melodyChatConn)
	redisCache := infra.NewRedisCache(universalClient)
	userRepo := chat.NewUserRepo(redisCache)
	userService := chat.NewUserService(userRepo)
	messageRepo := chat.NewMessageRepo(configConfig, redisCache)
	idGenerator, err := common.NewSonyFlake()
	if err != nil {
		return nil, err
	}
	messageService := chat.NewMessageService(messageRepo, userRepo, idGenerator)
	channelRepo := chat.NewChannelRepo(redisCache)
	channelService := chat.NewChannelService(channelRepo, userRepo, idGenerator)
	httpServer := chat.NewHttpServer(name, httpLogrus, configConfig, engine, melodyChatConn, messageSubscriber, userService, messageService, channelService)
	grpcLogrus := common.NewGrpcLogrus()
	grpcServer := chat.NewGrpcServer(grpcLogrus, configConfig, userService, channelService)
	router := chat.NewRouter(httpServer, grpcServer)
	infraCloser := chat.NewInfraCloser()
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, router, infraCloser, observabilityInjector)
	return server, nil
}

func InitializeUploaderServer(name string) (*common.Server, error) {
	httpLogrus := common.NewHttpLogrus()
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	engine := uploader.NewGinServer(name, httpLogrus, configConfig)
	httpServer := uploader.NewHttpServer(name, httpLogrus, configConfig, engine)
	router := uploader.NewRouter(httpServer)
	infraCloser := uploader.NewInfraCloser()
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, router, infraCloser, observabilityInjector)
	return server, nil
}

func InitializeUserServer(name string) (*common.Server, error) {
	httpLogrus := common.NewHttpLogrus()
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	engine := user.NewGinServer(name, httpLogrus, configConfig)
	universalClient, err := infra.NewRedisClient(configConfig)
	if err != nil {
		return nil, err
	}
	redisCache := infra.NewRedisCache(universalClient)
	userRepo := user.NewUserRepo(redisCache)
	idGenerator, err := common.NewSonyFlake()
	if err != nil {
		return nil, err
	}
	userService := user.NewUserService(userRepo, idGenerator)
	httpServer := user.NewHttpServer(name, httpLogrus, configConfig, engine, userService)
	grpcLogrus := common.NewGrpcLogrus()
	grpcServer := user.NewGrpcServer(grpcLogrus, configConfig, userService)
	router := user.NewRouter(httpServer, grpcServer)
	infraCloser := user.NewInfraCloser()
	observabilityInjector := common.NewObservabilityInjector(configConfig)
	server := common.NewServer(name, router, infraCloser, observabilityInjector)
	return server, nil
}
