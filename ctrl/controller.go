package ctrl

import "github.com/humpback/humpback-center/cluster"
import "github.com/humpback/humpback-center/etc"
import "github.com/humpback/humpback-center/repository"
import "github.com/humpback/humpback-center/storage"
import "github.com/humpback/gounits/logger"

type Controller struct {
	Configuration   *etc.Configuration
	DataStorage     *storage.DataStorage
	Cluster         *cluster.Cluster
	RepositoryCache *repository.RepositoryCache
}

func NewController(configuration *etc.Configuration) (*Controller, error) {

	mongo := configuration.Storage.Mongodb
	datastorage, err := storage.NewDataStorage(mongo.URIs)
	if err != nil {
		return nil, err
	}

	cluster, err := createCluster(configuration)
	if err != nil {
		return nil, err
	}

	repositorycache, err := createRepositoryCache(configuration)
	if err != nil {
		return nil, err
	}

	return &Controller{
		Cluster:         cluster,
		DataStorage:     datastorage,
		RepositoryCache: repositorycache,
	}, nil
}

func (c *Controller) Initialize() error {

	logger.INFO("[#ctrl#] controller initialize.....")
	if err := c.initCluster(); err != nil {
		return err
	}
	return c.startCluster()
}

func (c *Controller) UnInitialize() {

	c.stopCluster()
	c.DataStorage.Close()
	logger.INFO("[#ctrl#] controller uninitialized.")
}