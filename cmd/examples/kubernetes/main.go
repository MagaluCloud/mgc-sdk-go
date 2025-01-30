package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/kubernetes"
)

func main() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	k8sClient := kubernetes.New(c)

	// To slow, use this to create a new cluster
	// comNodePool := ExampleCreateCluster(k8sClient)
	idComNodePool := "948970cb-d8e5-4193-a9c7-c34257b02284"

	// semNodePool := ExampleCreateClusterWithoutNodepool(k8sClient)
	// idSemNodePool := "a5eac088-296e-4293-a842-a33e3f1db074"

	ExampleListClusters(k8sClient)
	WaitClusterRunning(k8sClient, idComNodePool)
	ExampleManageCluster(k8sClient, idComNodePool)
	// ExampleNodePoolOperations(k8sClient, clusterID)
	// ExampleListFlavorsAndVersions(k8sClient)
	// ExampleDeleteCluster(k8sClient, clusterID)
}

func WaitClusterRunning(k8sClient *kubernetes.KubernetesClient, clusterID string) {
	for {
		cluster, err := k8sClient.Clusters().Get(context.Background(), clusterID, nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(cluster.Status.State)

		if cluster.Status.State == "Running" {
			break
		}

		time.Sleep(10 * time.Second)
	}

	fmt.Println("Cluster running")
}

func ExampleCreateClusterWithoutNodepool(k8sClient *kubernetes.KubernetesClient) string {

	// Criar um novo cluster
	createReq := kubernetes.ClusterRequest{
		Name:         "my-kubernetes-cluster-" + strconv.FormatInt(time.Now().Unix(), 10),
		Version:      "v1.30.2",
		Description:  "Cluster de exemplo",
		NodePools:    []kubernetes.CreateNodePoolRequest{},
		AllowedCIDRs: []string{"192.168.0.0/24"},
	}

	cluster, err := k8sClient.Clusters().Create(context.Background(), createReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cluster criado com ID: %s\n", cluster.ID)
	return cluster.ID
}

func ExampleCreateCluster(k8sClient *kubernetes.KubernetesClient) string {

	// Criar um novo cluster
	createReq := kubernetes.ClusterRequest{
		Name:        "my-kubernetes-cluster-" + strconv.FormatInt(time.Now().Unix(), 10),
		Version:     "v1.30.2",
		Description: "Cluster de exemplo",
		NodePools: []kubernetes.CreateNodePoolRequest{
			{
				Name:     "default-pool",
				Flavor:   "cloud-k8s.gp1.small",
				Replicas: 3,
			},
		},
		AllowedCIDRs: []string{"192.168.0.0/24"},
	}

	cluster, err := k8sClient.Clusters().Create(context.Background(), createReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cluster criado com ID: %s\n", cluster.ID)
	return cluster.ID
}

func ExampleListClusters(k8sClient *kubernetes.KubernetesClient) {

	// Listar clusters com paginação
	clusters, err := k8sClient.Clusters().List(context.Background(), kubernetes.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{"node_pools"},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nClusters listados:")

	for _, cluster := range clusters {
		fmt.Printf("%s - %s (%s) - %s\n", cluster.ID, cluster.Name, cluster.Flavor, cluster.Status.State)
	}
}

func ExampleManageCluster(k8sClient *kubernetes.KubernetesClient, clusterID string) {
	ctx := context.Background()

	// Obter detalhes do cluster
	cluster, err := k8sClient.Clusters().Get(ctx, clusterID, []string{"node_pools"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nDetalhes do Cluster %s:\n", clusterID)
	fmt.Printf(" - Versão: %s\n", cluster.Version)
	fmt.Printf(" - Status: %s\n", cluster.Status.State)
	fmt.Printf(" - Node Pools: %d\n", len(cluster.NodePools))

	// Atualizar CIDRs permitidos
	updateReq := kubernetes.AllowedCIDRsUpdateRequest{
		AllowedCIDRs: []string{"192.168.0.0/24", "10.0.0.0/16"},
	}

	updatedCluster, err := k8sClient.Clusters().Update(ctx, clusterID, updateReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nCIDRs atualizados:", updatedCluster.AllowedCIDRs)

	// Obter kubeconfig
	kubeconfig, raw, err := k8sClient.Clusters().GetKubeConfig(ctx, clusterID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nKubeconfig (primeiras 100 caracteres):", kubeconfig.CurrentContext)
	fmt.Println("\nRawKubeconfig (primeiras 100 caracteres):", (*raw)[:100])
}

func ExampleNodePoolOperations(k8sClient *kubernetes.KubernetesClient, clusterID string) {

	ctx := context.Background()

	// Criar novo node pool
	poolReq := kubernetes.CreateNodePoolRequest{
		Name:     "gpu-pool",
		Flavor:   "cloud-k8s.gp1.large.gpu",
		Replicas: 2,
		Tags:     []string{"gpu", "ai"},
	}

	newPool, err := k8sClient.Nodepools().Create(ctx, clusterID, poolReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNode Pool criado: %s (%s)\n", newPool.Name, newPool.ID)

	// Listar node pools
	pools, err := k8sClient.Nodepools().List(ctx, clusterID, kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pools:")
	for _, pool := range pools {
		fmt.Printf(" - %s (%d replicas)\n", pool.Name, pool.Replicas)
	}

	// Atualizar node pool
	updateReq := kubernetes.PatchNodePoolRequest{
		Replicas: helpers.IntPtr(3),
	}

	updatedPool, err := k8sClient.Nodepools().Update(ctx, clusterID, newPool.ID, updateReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNode Pool atualizado: %d replicas\n", updatedPool.Replicas)

	// Deletar node pool
	err = k8sClient.Nodepools().Delete(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNode Pool deletado com sucesso")
}

func ExampleListFlavorsAndVersions(k8sClient *kubernetes.KubernetesClient) {

	// Listar versões disponíveis
	versions, err := k8sClient.Info().ListVersions(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nVersões disponíveis:")
	for _, v := range versions {
		fmt.Printf(" - %s (Deprecated: %v)\n", v.Version, v.Deprecated)
	}

	// Listar flavors disponíveis
	flavors, err := k8sClient.Flavors().List(context.Background(), kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nFlavors para Node Pools:")
	for _, f := range flavors.NodePool {
		fmt.Printf(" - %s (%d vCPUs, %dMB RAM)\n", f.Name, f.VCPU, f.RAM)
	}
}

func ExampleDeleteCluster(k8sClient *kubernetes.KubernetesClient, clusterID string) {

	// Esperar cluster ficar estável antes de deletar
	err := waitForClusterStatus(context.Background(), k8sClient, clusterID, "active")
	if err != nil {
		log.Fatal(err)
	}

	err = k8sClient.Clusters().Delete(context.Background(), clusterID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCluster %s deletado com sucesso\n", clusterID)
}

// Helper function para esperar status do cluster
func waitForClusterStatus(ctx context.Context, client *kubernetes.KubernetesClient, clusterID, targetStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout esperando status do cluster")
		case <-ticker.C:
			cluster, err := client.Clusters().Get(ctx, clusterID, nil)
			if err != nil {
				return err
			}

			fmt.Printf("Status atual do cluster: %s\n", cluster.Status.State)

			if cluster.Status.State == targetStatus {
				return nil
			}

			if cluster.Status.State == "error" {
				return fmt.Errorf("cluster em estado de erro: %s", cluster.Status.Message)
			}
		}
	}
}
