package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
)

type Peer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Pin struct {
	CID       string `json:"cid"`
	Name      string `json:"name"`
	ShardSize uint64 `json:"shard_size"`
}

func main() {
	router := gin.Default()
	ipfs, _ := client.NewDefaultClient(&client.Config{
		Host: os.Getenv("CLUSTER_HOST"),
	})
	router.POST("/pin/:pin", func(c *gin.Context) {
		pin := c.Param("pin")
		ci, err := cid.Parse(pin)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "bad pin",
			})
			return
		}
		name := c.Query("name")
		if name == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "pin name is required",
			})
			return
		}

		_, err = ipfs.Pin(c, ci, api.PinOptions{
			ReplicationFactorMin: 1,
			ReplicationFactorMax: 2,
			Name:                 name,
		})
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
	})

	router.GET("/peers", func(c *gin.Context) {
		peers, err := ipfs.Peers(c)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
		}

		ids := make([]Peer, 0, len(peers))
		for _, p := range peers {
			ids = append(ids, Peer{
				ID:   p.ID.Pretty(),
				Name: p.Peername,
			})
		}
		c.JSON(200, gin.H{
			"peers": ids,
		})
	})

	router.GET("/metrics", func(c *gin.Context) {
		metrics, err := ipfs.MetricNames(c)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
		}
		metricsMap := make(map[string]interface{}, len(metrics))

		for _, metric := range metrics {
			m, _ := ipfs.Metrics(c, metric)
			metricsMap[metric] = m
		}
		c.JSON(200, metricsMap)
	})

	router.GET("/pins", func(c *gin.Context) {
		allocations, err := ipfs.Allocations(c, api.AllType)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
		}

		ids := make([]Pin, 0, len(allocations))
		for _, a := range allocations {
			ids = append(ids, Pin{
				CID:       a.Cid.String(),
				Name:      a.Name,
				ShardSize: a.ShardSize,
			})
		}
		c.JSON(200, gin.H{
			"pins": ids,
		})
	})

	router.Run()
}
